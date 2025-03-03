package models

import (
	"time"

	"github.com/google/uuid"
)

type InventoryStatus string

const (
	StatusInStock    InventoryStatus = "IN_STOCK"
	StatusOutOfStock InventoryStatus = "OUT_OF_STOCK"
	StatusDeleted    InventoryStatus = "DELETED"
)

type Inventory struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	ProductName string          `json:"product_name" gorm:"type:varchar(255);not null"`
	Stock       int             `json:"stock" gorm:"not null"`
	Price       float64         `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time       `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Status      InventoryStatus `json:"status" gorm:"type:varchar(20);not null;default:'IN_STOCK'"`
}

// TableName specifies the table name for GORM
func (Inventory) TableName() string {
	return "inventories"
}
