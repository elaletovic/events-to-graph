package generators

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/elaletovic/events-to-graph/models"

	"github.com/ThreeDotsLabs/watermill"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/brianvoe/gofakeit"
)

var (
	initialEvents            = []string{models.ItemViewed, models.UserAddressValidated, models.UserAddressValidationFailed}
	itemViewedAfterEvents    = []string{models.ItemPurchased, models.ItemDropped, models.Nothing}
	itemPurchasedAfterEvents = []string{models.ItemDelivered, models.ItemNotDelivered}
	// InitialEventsTopic --
	InitialEventsTopic = "initial_events_topic"
	// CheckoutTopic --
	CheckoutTopic = "checkout_topic"
	// DeliveryTopic --
	DeliveryTopic = "delivery_topic"
)

// GenerateEvents --
func GenerateEvents(publisher message.Publisher) {
	for {
		eventType := gofakeit.RandString(initialEvents)
		var eventObj interface{}
		switch eventType {
		case models.ItemViewed:
			eventObj = models.ItemViewedPayload{
				ItemID: gofakeit.Number(1, 10),
				Price:  gofakeit.Price(0.01, 19.98),
			}

		case models.UserAddressValidated:
			eventObj = models.UserAddressValidatedPayload{
				Address: fmt.Sprintf("%s, %s, %s", gofakeit.Street(), gofakeit.City(), gofakeit.Country()),
			}
		case models.UserAddressValidationFailed:
			eventObj = models.UserAddressValidationFailedPayload{
				Address: fmt.Sprintf("%s, %s, %s", gofakeit.Street(), gofakeit.City(), gofakeit.Country()),
				Reason:  gofakeit.RandString([]string{"fake", "not occupied by user"}),
			}
		default:
			continue
		}

		eventPayload, err := json.Marshal(&eventObj)
		if err != nil {
			log.Printf("generateEvents: error while marshalling event payload, error %v\n", err)
			continue
		}

		obj := models.Event{
			UserID:    gofakeit.Number(1, 10),
			CreatedAt: time.Now().Unix(),
			Type:      eventType,
			Payload:   eventPayload,
		}
		payload, err := json.Marshal(&obj)
		if err != nil {
			log.Printf("generateEvents: error while marshalling main payload, error %v\n", err)
			continue
		}
		msg := message.NewMessage(watermill.NewUUID(), payload)
		middleware.SetCorrelationID(watermill.NewUUID(), msg)
		log.Printf("pushing message ID %s with payload %v\n", msg.UUID, string(msg.Payload))
		publisher.Publish(InitialEventsTopic, msg)

		time.Sleep(200 * time.Millisecond)
	}
}

// GeneratorHandler --
type GeneratorHandler struct {
}

// InitialEventsHandler --
func (gh GeneratorHandler) InitialEventsHandler(msg *message.Message) ([]*message.Message, error) {
	event := models.Event{}
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		log.Printf("failed to unmarshal initial events. Error %v, payload: %v\n", err, string(msg.Payload))
		return nil, err
	}
	newEvent := models.Event{
		UserID:    event.UserID,
		CreatedAt: time.Now().Unix(),
		Type:      gofakeit.RandString(itemViewedAfterEvents),
	}
	switch event.Type {
	case models.ItemViewed:
		eventPayload := models.ItemViewedPayload{}
		err = json.Unmarshal(event.Payload, &eventPayload)
		if err != nil {
			log.Printf("failed to unmarshal event payload. Error %v, payload: %v\n", err, string(event.Payload))
			return nil, err
		}

		var newEventObj interface{}
		switch newEvent.Type {
		case models.ItemPurchased:
			newEventObj = models.ItemPurchasedPayload{
				Price:    eventPayload.Price,
				ItemID:   eventPayload.ItemID,
				Quantity: gofakeit.Number(1, 5),
			}
		case models.ItemDropped:
			newEventObj = models.ItemDroppedPayload{
				Price:    eventPayload.Price,
				ItemID:   eventPayload.ItemID,
				Quantity: gofakeit.Number(1, 5),
			}
		case models.Nothing:
			log.Println("returning Nothing")
			return nil, nil
		}

		newEventPayload, err := json.Marshal(&newEventObj)
		if err != nil {
			log.Printf("failed to marshal new event payload. Error %v\n", err)
			return nil, err
		}
		newEvent.Payload = newEventPayload
		payload, err := json.Marshal(&newEvent)
		if err != nil {
			log.Printf("error while marshalling main payload, error %v\n", err)
			return nil, err
		}

		newMsg := message.NewMessage(watermill.NewUUID(), payload)
		return message.Messages{newMsg}, nil
	}
	return nil, nil
}

// PurchasedEventsHandler --
func (gh GeneratorHandler) PurchasedEventsHandler(msg *message.Message) ([]*message.Message, error) {
	event := models.Event{}
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		log.Printf("PurchasedEventsHandler: failed to unmarshal initial events. Error %v, payload: %v\n", err, string(msg.Payload))
		return nil, err
	}
	newEvent := models.Event{
		UserID:    event.UserID,
		CreatedAt: time.Now().Unix(),
		Type:      gofakeit.RandString(itemViewedAfterEvents),
	}
	switch event.Type {
	case models.ItemPurchased:
		eventPayload := models.ItemPurchasedPayload{}
		err = json.Unmarshal(event.Payload, &eventPayload)
		if err != nil {
			log.Printf("PurchasedEventsHandler: failed to unmarshal event payload. Error %v, payload: %v\n", err, string(event.Payload))
			return nil, err
		}

		var newEventObj interface{}
		switch newEvent.Type {
		case models.ItemDelivered:
			newEventObj = models.ItemDeliveredPayload{
				Address: fmt.Sprintf("%s, %s, %s", gofakeit.Street(), gofakeit.City(), gofakeit.Country()),
				ItemID:  eventPayload.ItemID,
			}
		case models.ItemNotDelivered:
			newEventObj = models.ItemNotDeliveredPayload{
				ItemID:  eventPayload.ItemID,
				Address: fmt.Sprintf("%s, %s, %s", gofakeit.Street(), gofakeit.City(), gofakeit.Country()),
				Reason:  gofakeit.RandString([]string{"fake", "not occupied by user"}),
			}
		}

		newEventPayload, err := json.Marshal(&newEventObj)
		if err != nil {
			log.Printf("PurchasedEventsHandler: failed to marshal new event payload. Error %v\n", err)
			return nil, err
		}
		newEvent.Payload = newEventPayload
		payload, err := json.Marshal(&newEvent)
		if err != nil {
			log.Printf("PurchasedEventsHandler: error while marshalling main payload, error %v\n", err)
			return nil, err
		}

		newMsg := message.NewMessage(watermill.NewUUID(), payload)
		return message.Messages{newMsg}, nil
	}
	return nil, nil
}
