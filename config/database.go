package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"mypackage/internal/models"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	PostgresDB *gorm.DB
	MongoDB    *mongo.Database
}

func InitDB() (*Database, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// PostgreSQL connection
	postgresDB, err := initPostgres()
	if err != nil {
		return nil, err
	}

	// MongoDB connection
	mongodb, err := initMongo()
	if err != nil {
		return nil, err
	}

	return &Database{
		PostgresDB: postgresDB,
		MongoDB:    mongodb,
	}, nil
}

func initPostgres() (*gorm.DB, error) {
	dsn := os.Getenv("URI")

	// Add performance configs
	config := &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	// 1. Create UUID extension first
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, err
	}

	// 2. Create enum type
	if err := db.Exec(`DO $$ BEGIN
        CREATE TYPE order_status AS ENUM ('pending', 'processing', 'completed', 'canceled');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END $$;`).Error; err != nil {
		return nil, err
	}

	// 3. Auto-migrate models to create tables
	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.Inventory{}); err != nil {
		return nil, err
	}

	// 4. Create indexes after tables exist
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_inventory_id ON orders(inventory_id)`).Error; err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	return db, nil
}

func initMongo() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	opts := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetTimeout(20 * time.Second).
		SetServerSelectionTimeout(20 * time.Second).
		SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection failed: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %v", err)
	}

	return client.Database("orderdb"), nil
}
