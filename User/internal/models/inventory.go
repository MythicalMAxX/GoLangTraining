package models

import (
    "time"
    "github.com/google/uuid"
)

type Inventory struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name      string    `gorm:"type:varchar(255)" json:"name"`
    Stock     int       `gorm:"type:int" json:"stock"`
    Price     float64   `gorm:"type:decimal(10,2)" json:"price"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ValidationError represents an inventory validation error
type ValidationError struct {
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}

// Validate checks if the inventory item is valid
func (i *Inventory) Validate() error {
    if i.Name == "" {
        return &ValidationError{Message: "name is required"}
    }
    if i.Stock < 0 {
        return &ValidationError{Message: "stock cannot be negative"}
    }
    if i.Price <= 0 {
        return &ValidationError{Message: "price must be greater than 0"}
    }
    return nil
}