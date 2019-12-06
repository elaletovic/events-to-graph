package main

import (
	stdSql "database/sql"
	"log"

	"github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	driver "github.com/go-sql-driver/mysql"
)

var (
	sqlEventsTopic = "events"
)

type eventProcessor struct {
}

func (ep eventProcessor) SaveToSQLDB(msg *message.Message) ([]*message.Message, error) {
	log.Printf("saving event to SQL DB. Message UUID: %v; Message Payload: %s", msg.UUID, string(msg.Payload))
	return message.Messages{msg}, nil
}

func (ep eventProcessor) SaveToGraphDB(msg *message.Message) error {
	log.Printf("saving event to graph DB. Message UUID: %v; Message Payload: %s", msg.UUID, string(msg.Payload))
	return nil
}

func createSQLDB() *stdSql.DB {
	conf := driver.NewConfig()
	conf.Net = "tcp"
	conf.User = "root"
	conf.Addr = "localhost"
	conf.DBName = "eventstore"

	db, err := stdSql.Open("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func createSQLPublisher() message.Publisher {
	db := createSQLDB()

	pub, err := sql.NewPublisher(
		db,
		sql.PublisherConfig{
			SchemaAdapter:        mySQLSchemaAdapter{},
			AutoInitializeSchema: true,
		},
		logger,
	)

	if err != nil {
		panic(err)
	}

	return pub
}
