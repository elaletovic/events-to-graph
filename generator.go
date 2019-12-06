package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/ThreeDotsLabs/watermill"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/brianvoe/gofakeit"
)

var (
	initialEvents            = []string{ItemViewed, UserAddressValidated, UserAddressValidationFailed}
	itemViewedAfterEvents    = []string{ItemPurchased, ItemDropped, Nothing}
	itemPurchasedAfterEvents = []string{ItemDelivered, ItemNotDelivered}
	initialEventsTopic       = "initial_events_topic"
	checkoutTopic            = "checkout_topic"
	deliveryTopic            = "delivery_topic"
)

func generateEvents(publisher message.Publisher) {
	for {
		eventType := gofakeit.RandString(initialEvents)
		var eventObj interface{}
		switch eventType {
		case ItemViewed:
			eventObj = ItemViewedPayload{
				ItemID: gofakeit.Number(1, 10),
				Price:  gofakeit.Price(0.01, 19.98),
			}

		case UserAddressValidated:
			eventObj = UserAddressValidatedPayload{
				Address: fmt.Sprintf("%s, %s, %s", gofakeit.Street(), gofakeit.City(), gofakeit.Country()),
			}
		case UserAddressValidationFailed:
			eventObj = UserAddressValidationFailedPayload{
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

		obj := Event{
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
		publisher.Publish(initialEventsTopic, msg)

		time.Sleep(200 * time.Millisecond)
	}
}

type generatorHandler struct {
}

func (gh generatorHandler) InitialEventsHandler(msg *message.Message) ([]*message.Message, error) {
	event := Event{}
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		log.Printf("failed to unmarshal initial events. Error %v, payload: %v\n", err, string(msg.Payload))
		return nil, err
	}
	newEvent := Event{
		UserID:    event.UserID,
		CreatedAt: time.Now().Unix(),
		Type:      gofakeit.RandString(itemViewedAfterEvents),
	}
	switch event.Type {
	case ItemViewed:
		eventPayload := ItemViewedPayload{}
		err = json.Unmarshal(event.Payload, &eventPayload)
		if err != nil {
			log.Printf("failed to unmarshal event payload. Error %v, payload: %v\n", err, string(event.Payload))
			return nil, err
		}

		var newEventObj interface{}
		switch newEvent.Type {
		case ItemPurchased:
			newEventObj = ItemPurchasedPayload{
				Price:    eventPayload.Price,
				ItemID:   eventPayload.ItemID,
				Quantity: gofakeit.Number(1, 5),
			}
		case ItemDropped:
			newEventObj = ItemDroppedPayload{
				Price:    eventPayload.Price,
				ItemID:   eventPayload.ItemID,
				Quantity: gofakeit.Number(1, 5),
			}
		case Nothing:
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

func (gh generatorHandler) PurchasedEventsHandler(msg *message.Message) ([]*message.Message, error) {
	event := Event{}
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		log.Printf("PurchasedEventsHandler: failed to unmarshal initial events. Error %v, payload: %v\n", err, string(msg.Payload))
		return nil, err
	}
	newEvent := Event{
		UserID:    event.UserID,
		CreatedAt: time.Now().Unix(),
		Type:      gofakeit.RandString(itemViewedAfterEvents),
	}
	switch event.Type {
	case ItemPurchased:
		eventPayload := ItemPurchasedPayload{}
		err = json.Unmarshal(event.Payload, &eventPayload)
		if err != nil {
			log.Printf("PurchasedEventsHandler: failed to unmarshal event payload. Error %v, payload: %v\n", err, string(event.Payload))
			return nil, err
		}

		var newEventObj interface{}
		switch newEvent.Type {
		case ItemDelivered:
			newEventObj = ItemDeliveredPayload{
				Address: fmt.Sprintf("%s, %s, %s", gofakeit.Street(), gofakeit.City(), gofakeit.Country()),
				ItemID:  eventPayload.ItemID,
			}
		case ItemNotDelivered:
			newEventObj = ItemNotDeliveredPayload{
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
