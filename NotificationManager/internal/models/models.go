package models

type Notification struct {
	Type       string `json:"type"`
	OrderID    string `json:"order_id"`
	DeliveryID string `json:"delivery_id,omitempty"`
	Message    string `json:"message"`
}
