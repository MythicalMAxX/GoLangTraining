package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	Cash PaymentMethod = "cash"
	Card PaymentMethod = "card"
)

type Payment struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID       uuid.UUID     `json:"order_id" gorm:"type:uuid;not null"`
	PaymentDate   time.Time     `json:"payment_date"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}
