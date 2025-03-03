package models

import (
    "time"
    "github.com/google/uuid"
)

type InventoryRequest struct {
    UserID      string  `json:"user_id"`
    ProductName string  `json:"product_name"`
    Stock       int     `json:"stock"`
    Price       float64 `json:"price"`
}

type InventoryResponse struct {
    ID          uuid.UUID `json:"id"`
    UserID      uuid.UUID `json:"user_id"`
    ProductName string    `json:"product_name"`
    Stock       int       `json:"stock"`
    Price       float64   `json:"price"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Add this new struct for paginated list response
type InventoryListResponse struct {
    Data      []InventoryResponse `json:"data"`
    Total     int                `json:"total"`
    Page      int                `json:"page"`
    PageSize  int                `json:"page_size"`
}