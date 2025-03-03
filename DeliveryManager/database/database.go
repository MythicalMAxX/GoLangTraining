package database

import (
	"delivery-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbURI string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schemas
	err = db.AutoMigrate(
		&models.Order{},
		&models.OrderItem{},
		&models.DeliveryPartner{},
		&models.OrderDelivery{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
