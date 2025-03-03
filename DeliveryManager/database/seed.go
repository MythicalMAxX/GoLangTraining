package database

import (
	"delivery-service/internal/models"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var locations = []string{
	"Mumbai", "Delhi", "Bangalore", "Hyderabad", "Chennai",
	"Kolkata", "Pune", "Ahmedabad", "Jaipur", "Surat",
}

func SeedDeliveryPartners(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.DeliveryPartner{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // Partners already seeded
	}

	partners := make([]models.DeliveryPartner, 10)
	for i := 0; i < 10; i++ {
		partners[i] = models.DeliveryPartner{
			ID:       uuid.New().String(),
			Name:     fmt.Sprintf("Delivery Partner %d", i+1),
			Status:   models.DeliveryPartnerStatusAvailable,
			Location: locations[rand.Intn(len(locations))],
		}
	}

	return db.Create(&partners).Error
}
