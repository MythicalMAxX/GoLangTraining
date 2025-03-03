package models

import (
    "time"
    "github.com/google/uuid"
)

type OrderRequest struct {
    InventoryID uuid.UUID `json:"inventory_id"`
    OrderCount  int       `json:"order_count"`
}

type OrderResponse struct {
    ID            uuid.UUID  `json:"id"`
    UserID        uuid.UUID  `json:"user_id"`
    InventoryID   uuid.UUID  `json:"inventory_id"`
    ProductName   string     `json:"product_name"`
    Price         float64    `json:"price"`
    TotalPrice    float64    `json:"total_price"`
    OrderCount    int        `json:"order_count"`
    PaymentStatus string     `json:"payment_status"`
    OrderStatus   string     `json:"order_status"`
    CreatedAt     *time.Time `json:"created_at,omitempty"`
    UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}