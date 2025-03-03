package config

import (
	"log"
	"os"
	"restaurantservice/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := os.Getenv("USER_DB_URI")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// First, enable the UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatal("Failed to create UUID extension:", err)
	}

	// Drop the existing table if it exists (during development)
	if err := db.Migrator().DropTable(&models.Restaurant{}); err != nil {
		log.Fatal("Failed to drop table:", err)
	}

	// Then auto-migrate with the correct UUID type
	if err := db.AutoMigrate(&models.Restaurant{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully")
	return db
}
