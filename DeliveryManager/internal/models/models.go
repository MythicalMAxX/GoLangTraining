package models

import (
	"time"
)

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
	RestaurantID string      `json:"restaurant_id" gorm:"index"`
	Items        []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount  float64     `json:"total_amount"`
	Status       string      `json:"status" gorm:"index"`
	DeliveryAddr string      `json:"delivery_address"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type DeliveryPartner struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderUpdate struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type NotificationType string

const (
	OrderProcessed        NotificationType = "ORDER_PROCESSED"
	OrderPickedByDelivery NotificationType = "ORDER_PICKED_BY_DELIVERY"
	OrderDelivered        NotificationType = "ORDER_DELIVERED"
)

type Notification struct {
	Type       NotificationType `json:"type"`
	OrderID    string           `json:"order_id"`
	Message    string           `json:"message"`
	DeliveryID string           `json:"delivery_id,omitempty"`
}

type OrderDelivery struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OrderID           string    `json:"order_id" gorm:"index;not null"`
	RestaurantID      string    `json:"restaurant_id" gorm:"index;not null"`
	DeliveryPartnerID string    `json:"delivery_partner_id" gorm:"index;not null"`
	Status            string    `json:"status" gorm:"not null"`
	PickupTime        time.Time `json:"pickup_time"`
	DeliveryTime      time.Time `json:"delivery_time,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type OrderMessage struct {
	Type      string `json:"type"`
	Order     Order  `json:"order"`
	Timestamp string `json:"timestamp"`
}
