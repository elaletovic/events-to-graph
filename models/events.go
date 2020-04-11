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
	// ItemCreated event is triggered when a new item is added to the system
	ItemCreated = "item_created"
	// UserRegistered event is triggered when a new user is registered
	UserRegistered = "user_registered"
)

// Event --
type Event struct {
	CreatedAt int64  `json:"created_at"`
	Type      string `json:"type"`
	Payload   []byte `json:"payload"`
	Origin    string `json:"origin"`
}

// ItemViewedPayload --
type ItemViewedPayload struct {
	UserID uint64 `json:"user_id"`
	ItemID uint64 `json:"item_id"`
}

// ItemPurchasedPayload --
type ItemPurchasedPayload struct {
	ItemID uint64 `json:"item_id"`
	UserID uint64 `json:"user_id"`
}

// ItemDroppedPayload --
type ItemDroppedPayload struct {
	ItemID uint64 `json:"item_id"`
	UserID uint64 `json:"user_id"`
}

// ItemDeliveredPayload --
type ItemDeliveredPayload struct {
	ItemID  uint64 `json:"item_id"`
	UserID  uint64 `json:"user_id"`
	Address string `json:"address"`
}

// ItemNotDeliveredPayload --
type ItemNotDeliveredPayload struct {
	ItemID  uint64 `json:"item_id"`
	UserID  uint64 `json:"user_id"`
	Address string `json:"address"`
	Reason  string `json:"reason"`
}
