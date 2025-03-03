package config

import (
	"inventorymanager/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbURI := os.Getenv("USER_DB_URI")
    db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Enable UUID extension
    if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
        log.Fatal("Failed to enable uuid-ossp extension:", err)
    }

    // Auto migrate the schema with table name specification
    err = db.AutoMigrate(&models.Inventory{})
    if err != nil {
        log.Fatal("Failed to auto migrate:", err)
    }

    return db
}
