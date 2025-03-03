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
	PaymentStatus PaymentStatus `json:"payment_status"`
	OrderStatus   OrderStatus   `json:"order_status"`
	Price         float64       `json:"price"`
	OrderCount    int           `json:"order_count"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type CreateOrderRequest struct {
	InventoryID uuid.UUID `json:"inventory_id"`
	UserID      uuid.UUID `json:"user_id"`
	Price       float64   `json:"price"`
	OrderCount  int       `json:"order_count"`
}
