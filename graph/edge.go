package graph

import (
	"fmt"
	"log"
)

type viewed struct {
	From       string `json:"_from"`
	To         string `json:"_to"`
	ViewedDate int64  `json:"viewed_date"`
}

// AddViewed creates a viewed edge and adds it to the collection
func (s *Store) AddViewed(userID, itemID uint64, date int64) error {
	edge := viewed{
		From:       fmt.Sprintf("%s/%d", usersCollection, userID),
		To:         fmt.Sprintf("%s/%d", itemsCollection, itemID),
		ViewedDate: date,
	}

	log.Println(edge, s.UsersToItems.Collection, s.UsersToItems.Constraints)
	_, err := s.UsersToItems.Collection.CreateDocument(nil, edge)
	if err != nil {
		log.Fatalf("AddViewed: %v", err)
	}

	return err
}
