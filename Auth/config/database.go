package config

import (
    "fmt"
    "log"
    "mypackage/internal/models"
    "mypackage/pkg/jwt"
    "os"
    "time"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
    if err := godotenv.Load(); err != nil {
        return nil, err
    }

    // Set JWT secret key
    jwtSecret := os.Getenv("JWT_SECRET_KEY")
    if jwtSecret == "" {
        return nil, fmt.Errorf("JWT_SECRET_KEY is required")
    }
    jwt.SetSecretKey([]byte(jwtSecret))

    // Connect to database
    dsn := os.Getenv("URI")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // AutoMigrate the schema
    if err := db.AutoMigrate(&models.User{}); err != nil {
        return nil, err
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