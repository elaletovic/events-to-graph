package main

import (
	stdSql "database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-sql/pkg/sql"

	"github.com/ThreeDotsLabs/watermill/message"
)

type mySQLSchemaAdapter struct{}

func (m mySQLSchemaAdapter) SchemaInitializingQueries(topic string) []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS ` + topic + ` (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			message_id VARCHAR(50) NOT NULL,
			user_id BIGINT NOT NULL,
			type VARCHAR(36) NOT NULL,
			payload VARCHAR(4000) NOT NULL,
			created_at BIGINT NOT NULL
		);`}
}

func (m mySQLSchemaAdapter) InsertQuery(topic string, msgs message.Messages) (string, []interface{}, error) {
	query := fmt.Sprintf("INSERT INTO %s (message_id, user_id, type, payload, created_at) VALUES %s",
		topic,
		strings.TrimRight(strings.Repeat(`(?,?,?,?,?),`, len(msgs)), ","))
	args := []interface{}{}
	for _, msg := range msgs {
		event := Event{}
		err := json.Unmarshal(msg.Payload, &event)
		if err != nil {
			return "", nil, err
		}
		args = append(args, msg.UUID, event.UserID, event.Type, string(event.Payload), event.CreatedAt)
	}

	return query, args, nil
}

func (m mySQLSchemaAdapter) SelectQuery(topic string, consumerGroup string, offsetsAdapter sql.OffsetsAdapter) (string, []interface{}) {
	nextOffsetQuery, nextOffsetArgs := offsetsAdapter.NextOffsetQuery(topic, consumerGroup)
	selectQuery := `
		SELECT id, message_id, user_id, type, payload, created_at FROM ` + topic + `
		WHERE 
			id > (` + nextOffsetQuery + `)
		ORDER BY 
			id ASC
		LIMIT 1`

	return selectQuery, nextOffsetArgs
}

func (m mySQLSchemaAdapter) UnmarshalMessage(row *stdSql.Row) (offset int, msg *message.Message, err error) {
	event := Event{}
	var id int
	err = row.Scan(&id, &event.UserID, &event.Type, &event.Payload, &event.CreatedAt)
	if err != nil {
		return 0, nil, err
	}

	payload, err := json.Marshal(&event)
	if err != nil {
		return 0, nil, err
	}

	msg = message.NewMessage(watermill.NewULID(), payload)

	return id, msg, nil
}
