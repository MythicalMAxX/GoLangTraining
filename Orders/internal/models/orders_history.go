package models

import (
    "time"
)

type OrderLog struct {
    Status    OrderStatus `json:"status" bson:"status"`
    Timestamp time.Time   `json:"timestamp" bson:"timestamp"`
}

type OrderHistory struct {
    OrderID string     `json:"order_id" bson:"order_id"`
    Status  OrderStatus `json:"status" bson:"status"`
    Log     []OrderLog  `json:"log" bson:"log"`
}