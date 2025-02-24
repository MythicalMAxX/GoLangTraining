package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name      string    `gorm:"type:varchar(255)" json:"name"`
    Email     string    `gorm:"type:varchar(255);unique" json:"email"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

