package db

import (
	"log"
	"paymentservice/config"
	"paymentservice/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DBUri), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop existing table to handle the column type change
	err = DB.Migrator().DropTable(&models.Payment{})
	if err != nil {
		log.Printf("Warning during table drop: %v", err)
	}

	// AutoMigrate the schemas
	err = DB.AutoMigrate(&models.Payment{}, &models.Order{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
