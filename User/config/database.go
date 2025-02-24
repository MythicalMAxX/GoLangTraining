package config

import (
    "fmt"
    "log"
    "os"
    "time"
    "mypackage/internal/models"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
    if err := godotenv.Load(); err != nil {
        return nil, err
    }

    dsn := os.Getenv("USER_DB_URI")
    config := &gorm.Config{
        PrepareStmt:                              true,
        SkipDefaultTransaction:                   true,
        DisableForeignKeyConstraintWhenMigrating: true,
    }

    db, err := gorm.Open(postgres.Open(dsn), config)
    if err != nil {
        return nil, err
    }

    // Create UUID extension
    if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
        return nil, fmt.Errorf("failed to create uuid extension: %v", err)
    }

    // Auto-migrate models
    if err := db.AutoMigrate(&models.User{}, &models.Inventory{}); err != nil {
        return nil, fmt.Errorf("failed to migrate models: %v", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    log.Println("Database connection established")
    return db, nil
}