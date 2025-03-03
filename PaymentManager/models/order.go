package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string
type OrderStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"

	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
)

type Order struct {
	ID            uuid.UUID     `json:"id"`
	InventoryID   uuid.UUID     `json:"inventory_id"`
	UserID        uuid.UUID     `json:"user_id"`
	PaymentStatus PaymentStatus `json:"payment_status,omitempty"`
	OrderStatus   OrderStatus   `json:"order_status,omitempty"`
	Price         float64       `json:"price"`
	OrderCount    int           `json:"order_count"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
}
