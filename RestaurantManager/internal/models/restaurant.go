package models

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name" binding:"required"`
	Address   string    `json:"address" binding:"required"`
	Cuisine   string    `json:"cuisine" binding:"required"`
	Rating    float32   `json:"rating"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (restaurant *Restaurant) BeforeCreate() error {
	restaurant.ID = uuid.New()
	return nil
}
