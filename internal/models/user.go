package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"type:varchar(255)" json:"name"`
	Email    string    `gorm:"type:varchar(255);unique" json:"email"`
	Password string    `gorm:"type:varchar(255)" json:"-"` // "-" means don't show in JSON response
}
