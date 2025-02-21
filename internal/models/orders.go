package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

// Order status enum values
const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCanceled   OrderStatus = "canceled"
)

// IsValid checks if the status is valid
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusProcessing, OrderStatusCompleted, OrderStatusCanceled:
		return true
	}
	return false
}

// Order represents the order model
type Order struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID   `gorm:"type:uuid" json:"user_id"`
	InventoryID uuid.UUID   `gorm:"type:uuid" json:"inventory_id"`
	Quantity    int         `gorm:"type:int" json:"quantity"`
	Amount      float64     `gorm:"type:decimal(10,2)" json:"amount"`
	Status      OrderStatus `gorm:"type:order_status" json:"status"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	User        User        `gorm:"foreignKey:UserID" json:"user"`
	Inventory   Inventory   `gorm:"foreignKey:InventoryID" json:"inventory"`
}
