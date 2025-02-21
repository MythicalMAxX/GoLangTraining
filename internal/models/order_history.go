package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHistory struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	OrderID string             `bson:"order_id"`
	Status  OrderStatus        `bson:"status"`
	Log     []OrderLog         `bson:"log"`
}

type OrderLog struct {
	Status    OrderStatus `bson:"status"`
	Timestamp time.Time   `bson:"timestamp"`
}
