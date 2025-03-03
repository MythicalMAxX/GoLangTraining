package models

import "time"

type OrderItem struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	OrderID   string  `json:"order_id" gorm:"index"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID           string      `json:"id" gorm:"primaryKey"`
	CustomerID   string      `json:"customer_id" gorm:"index"`
	RestaurantID string      `json:"restaurant_id" gorm:"index"` // Add this field
	Items        []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount  float64     `json:"total_amount"`
	Status       string      `json:"status" gorm:"index"`
	DeliveryAddr string      `json:"delivery_address"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
