package models

import (
    "time"
    "github.com/google/uuid"
)

type OrderStatus string

const (
    OrderStatusPending    OrderStatus = "pending"
    OrderStatusProcessing OrderStatus = "processing"
    OrderStatusCompleted  OrderStatus = "completed"
    OrderStatusCanceled   OrderStatus = "canceled"
)

func (s OrderStatus) IsValid() bool {
    switch s {
    case OrderStatusPending, OrderStatusProcessing, OrderStatusCompleted, OrderStatusCanceled:
        return true
    }
    return false
}

type Order struct {
    ID          uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID      uuid.UUID   `gorm:"type:uuid" json:"user_id"`
    InventoryID uuid.UUID   `gorm:"type:uuid" json:"inventory_id"`
    Quantity    int         `gorm:"type:int" json:"quantity"`
    Amount      float64     `gorm:"type:decimal(10,2)" json:"amount"`
    Status      OrderStatus `gorm:"type:varchar(20)" json:"status"`
    CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}