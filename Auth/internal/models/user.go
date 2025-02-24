package models

import (
    "github.com/google/uuid"
)

type Role string

const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
)

func (r Role) IsValid() bool {
    switch r {
    case RoleAdmin, RoleUser:
        return true
    }
    return false
}

type User struct {
    ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name     string    `gorm:"type:varchar(255)" json:"name"`
    Email    string    `gorm:"type:varchar(255);unique" json:"email"`
    Password string    `gorm:"type:varchar(255)" json:"-"` 
    Role     Role      `gorm:"type:varchar(20);default:'user'" json:"role"`
}