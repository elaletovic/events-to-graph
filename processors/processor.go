package processors

import (
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
)

// EventProcessor --
type EventProcessor struct {
}

// SaveToGraphDB --
func (ep EventProcessor) SaveToGraphDB(msg *message.Message) error {
	log.Printf("saving event to graph DB. Message UUID: %v; Message Payload: %s", msg.UUID, string(msg.Payload))
	return nil
}
