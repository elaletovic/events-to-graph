package graph

import (
	"fmt"
	"log"
)

type link struct {
	From      string `json:"_from"`
	To        string `json:"_to"`
	CreatedAt int64  `json:"created_at"`
}

type viewed link
type purchased link
type dropped link

type delivered struct {
	Address string
	link
}

type notDelivered struct {
	Address string
	Reason  string
	link
}

// AddViewed creates a viewed edge and adds it to the collection
func (s *Store) AddViewed(userID, itemID uint64, date int64) error {
	edge := viewed{
		From:      fmt.Sprintf("%s/%d", usersCollection, userID),
		To:        fmt.Sprintf("%s/%d", itemsCollection, itemID),
		CreatedAt: date,
	}
	return s.createEdge(edge, s.ViewedEdge)
}

// AddPurchased creates a purchased edge and adds it to the collection
func (s *Store) AddPurchased(userID, itemID uint64, date int64) error {
	edge := purchased{
		From:      fmt.Sprintf("%s/%d", usersCollection, userID),
		To:        fmt.Sprintf("%s/%d", itemsCollection, itemID),
		CreatedAt: date,
	}
	return s.createEdge(edge, s.PurchasedEdge)
}

// AddDropped creates a dropped edge and adds it to the collection
func (s *Store) AddDropped(userID, itemID uint64, date int64) error {
	edge := dropped{
		From:      fmt.Sprintf("%s/%d", usersCollection, userID),
		To:        fmt.Sprintf("%s/%d", itemsCollection, itemID),
		CreatedAt: date,
	}
	return s.createEdge(edge, s.DroppedEdge)
}

// AddDelivered creates an item_delivered edge and adds it to the collection
func (s *Store) AddDelivered(userID, itemID uint64, address string, date int64) error {
	base := link{
		From:      fmt.Sprintf("%s/%d", itemsCollection, itemID),
		To:        fmt.Sprintf("%s/%d", usersCollection, userID),
		CreatedAt: date,
	}
	edge := delivered{
		Address: address,
		link:    base,
	}
	return s.createEdge(edge, s.DeliveredEdge)
}

// AddNotDelivered creates a not_delivered edge and adds it to the collection
func (s *Store) AddNotDelivered(userID, itemID uint64, address, reason string, date int64) error {
	base := link{
		From:      fmt.Sprintf("%s/%d", itemsCollection, itemID),
		To:        fmt.Sprintf("%s/%d", usersCollection, userID),
		CreatedAt: date,
	}
	edge := notDelivered{
		Address: address,
		Reason:  reason,
		link:    base,
	}
	return s.createEdge(edge, s.NotDeliveredEdge)
}

func (s *Store) createEdge(edge interface{}, edgeCol edgeCollection) error {
	_, err := edgeCol.Collection.CreateDocument(nil, edge)
	if err != nil {
		log.Fatalf("createEdge %s: error: %v, payload: %v", edgeCol.Label, err, edge)
	}
	return err
}
