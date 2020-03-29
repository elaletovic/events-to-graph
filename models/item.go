package models

// Item --
type Item struct {
	ID           uint64  `json:"id"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	Manufacturer string  `json:"manufacturer"`
	Origin       string  `json:"origin"`
}
