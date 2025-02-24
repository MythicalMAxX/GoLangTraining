package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"mypackage/internal/models"
	"mypackage/pkg/jwt"
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

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is requied")
	}
	jwt.SetSecretKey([]byte(jwtSecret))

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

	config := &gorm.Config{
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	// Initial Transaction Implementation
	// First create extensions and types in a separate transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		// Create UUID extension
		if err := tx.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
			return fmt.Errorf("failed to create uuid extension: %v", err)
		}

		// Create enum type with error handling
		if err := tx.Exec(`DO $$ 
        BEGIN 
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
                CREATE TYPE order_status AS ENUM ('pending', 'processing', 'completed', 'canceled');
            END IF;
        END $$;`).Error; err != nil {
			return fmt.Errorf("failed to create order_status enum: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize extensions: %v", err)
	}

	// Then auto-migrate models to create tables
	if err := db.AutoMigrate(
		&models.User{},
		&models.Order{},
		&models.Inventory{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate models: %v", err)
	}

	// Create indexes after tables exist
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_orders_inventory_id ON orders(inventory_id)",
		"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Get connection pool settings from env
	maxOpenConns := 5
	maxIdleConns := 5
	connMaxLifetime := time.Minute * 10
	connMaxIdleTime := time.Minute * 10

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func initMongo() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Updated MongoDB connection options for replica set
	opts := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetTLSConfig(&tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false, // Changed back to secure mode
		}).
		SetTimeout(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second).
		SetRetryWrites(true).
		SetRetryReads(true).
		SetDirect(false) // Changed: Remove direct connection for replica sets

	// Create client
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection failed: %v", err)
	}

	// Ping with longer timeout
	pingCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %v", err)
	}

	return client.Database("orderdb"), nil
}
