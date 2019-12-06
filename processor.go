package main

import (
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
)

type eventProcessor struct {
}

func (ep eventProcessor) SaveToSQLDB(msg *message.Message) error {
	log.Printf("saving event to SQL DB. Message UUID: %v; Message Payload: %s", msg.UUID, string(msg.Payload))
	return nil
}

func (ep eventProcessor) SaveToGraphDB(msg *message.Message) error {
	log.Printf("saving event to graph DB. Message UUID: %v; Message Payload: %s", msg.UUID, string(msg.Payload))
	return nil
}
