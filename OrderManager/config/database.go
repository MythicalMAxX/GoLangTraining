package config

import (
	"fmt"
	"log"
	"orderservice/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dbURI string) (*gorm.DB, error) {
	if dbURI == "" {
		return nil, fmt.Errorf("database URI is required")
	}

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Order{}, &models.OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected and migration completed")
	return db, nil
}
