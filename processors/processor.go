package processors

import (
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/elaletovic/events-to-graph/graph"
	"github.com/elaletovic/events-to-graph/models"
)

// GraphProcessor --
type GraphProcessor interface {
	SaveToGraph(msg *message.Message) error
}

// eventProcessor --
type eventProcessor struct {
	Store *graph.Store
}

// NewEventProcessor --
func NewEventProcessor(store *graph.Store) GraphProcessor {
	return &eventProcessor{
		Store: store,
	}
}

// SaveToGraph --
func (ep *eventProcessor) SaveToGraph(msg *message.Message) error {
	event := models.Event{}
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		log.Printf("failed to unmarshal event in SaveToGraph. Error %v, payload: %v\n", err, string(msg.Payload))
		return err
	}
	log.Printf("received: %s %s %s\n", msg.UUID, event.Type, string(event.Payload))
	switch event.Type {
	case models.UserRegistered:
		user := new(models.User)
		err := json.Unmarshal(event.Payload, user)
		if err != nil {
			return err
		}
		return ep.Store.AddUser(user.ID, user.Name, user.Age)
	case models.ItemCreated:
		item := new(models.Item)
		err := json.Unmarshal(event.Payload, item)
		if err != nil {
			return err
		}
		return ep.Store.AddItem(item.ID, item.Price, item.Title, item.Manufacturer, item.Origin)
	case models.ItemViewed:
		itemViewed := new(models.ItemViewedPayload)
		err := json.Unmarshal(event.Payload, itemViewed)
		if err != nil {
			return err
		}
		return ep.Store.AddViewed(itemViewed.UserID, itemViewed.ItemID, event.CreatedAt)
	}
	return nil
}
