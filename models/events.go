package models

const (
	//Nothing is an event that is not handled
	Nothing = ""
	//ItemViewed event is triggered when an item has been viewed by a user
	ItemViewed = "item_viewed"
	// ItemPurchased event is triggered when an item is purchased by a user
	ItemPurchased = "item_purchased"
	// ItemDropped event is triggered when user deletes an item from his cart
	ItemDropped = "item_dropped"
	// ItemDelivered event is triggered when an item is delivered to a user
	ItemDelivered = "item_delivered"
	// ItemNotDelivered event is triggered when an item is not delivered to a user
	ItemNotDelivered = "item_not_delivered"
	// UserAddressValidated event is triggered when user's address is validated
	UserAddressValidated = "user_address_validated"
	// UserAddressValidationFailed event is triggered when user's address valdiation fails
	UserAddressValidationFailed = "user_address_validation_failed"
)

// Event --
type Event struct {
	UserID    int    `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	Type      string `json:"type"`
	Payload   []byte `json:"payload"`
}

// ItemViewedPayload --
type ItemViewedPayload struct {
	ItemID int     `json:"item_id"`
	Price  float64 `json:"price"`
}

// ItemPurchasedPayload --
type ItemPurchasedPayload struct {
	ItemID   int     `json:"item_id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// ItemDroppedPayload --
type ItemDroppedPayload struct {
	ItemID   int     `json:"item_id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// ItemDeliveredPayload --
type ItemDeliveredPayload struct {
	ItemID  int    `json:"item_id"`
	Address string `json:"address"`
}

// ItemNotDeliveredPayload --
type ItemNotDeliveredPayload struct {
	ItemID  int    `json:"item_id"`
	Address string `json:"address"`
	Reason  string `json:"reason"`
}

// UserAddressValidatedPayload --
type UserAddressValidatedPayload struct {
	Address string `json:"address"`
}

// UserAddressValidationFailedPayload --
type UserAddressValidationFailedPayload struct {
	Address string `json:"address"`
	Reason  string `json:"reason"`
}
