package processors

import (
	stdSql "database/sql"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-sql/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	schema "github.com/elaletovic/events-to-graph/sql"
	driver "github.com/go-sql-driver/mysql"
)

var (
	// SqlEventsTopic --
	SqlEventsTopic = "events"
)

// EventProcessor --
type EventProcessor struct {
}

// SaveToSQLDB --
func (ep EventProcessor) SaveToSQLDB(msg *message.Message) ([]*message.Message, error) {
	log.Printf("saving event to SQL DB. Message UUID: %v; Message Payload: %s", msg.UUID, string(msg.Payload))
	return message.Messages{msg}, nil
}

// SaveToGraphDB --
func (ep EventProcessor) SaveToGraphDB(msg *message.Message) error {
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

// CreateSQLPublisher --
func CreateSQLPublisher(logger watermill.LoggerAdapter) message.Publisher {
	db := createSQLDB()

	pub, err := sql.NewPublisher(
		db,
		sql.PublisherConfig{
			SchemaAdapter:        schema.MySQLSchemaAdapter{},
			AutoInitializeSchema: true,
		},
		logger,
	)

	if err != nil {
		panic(err)
	}

	return pub
}
