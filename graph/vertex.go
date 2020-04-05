package graph

import (
	"fmt"
)

type user struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Key  string `json:"_key"`
}

type item struct {
	ID           uint64  `json:"id"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	Manufacturer string  `json:"manufacturer"`
	Origin       string  `json:"origin"`
	Key          string  `json:"_key"`
}

// AddUser creates user document and adds it to the graph
func (s *Store) AddUser(id uint64, name string, age int) error {
	u := user{
		ID:   id,
		Name: name,
		Age:  age,
		Key:  fmt.Sprintf("%d", id),
	}

	_, err := s.Users.CreateDocument(nil, u)
	return err
}

// AddItem create item document and adds it to the graph
func (s *Store) AddItem(id uint64, price float64, title, manufacturer, origin string) error {
	i := item{
		ID:           id,
		Price:        price,
		Title:        title,
		Manufacturer: manufacturer,
		Origin:       origin,
		Key:          fmt.Sprintf("%d", id),
	}

	_, err := s.Items.CreateDocument(nil, i)
	return err
}
