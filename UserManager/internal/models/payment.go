package models

import (
    "time"
    "github.com/google/uuid"
)

type PaymentMethod string

const (
    Card PaymentMethod = "card"
    Cash PaymentMethod = "cash"
)

type PaymentRequest struct {
    OrderID       uuid.UUID     `json:"order_id"`
    PaymentMethod PaymentMethod `json:"payment_method"`
}

type PaymentResponse struct {
    ID            uuid.UUID     `json:"id"`
    OrderID       uuid.UUID     `json:"order_id"`
    PaymentMethod PaymentMethod `json:"payment_method"`
    PaymentDate   time.Time     `json:"payment_date"`
}