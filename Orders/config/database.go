package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnection struct {
	PostgresDB *gorm.DB
	MongoDB    *mongo.Database
}

func InitDB() (*DBConnection, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// Initialize PostgreSQL
	postgresDB, err := initPostgres()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	// Initialize MongoDB
	mongodb, err := initMongoDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	return &DBConnection{
		PostgresDB: postgresDB,
		MongoDB:    mongodb,
	}, nil
}

func initPostgres() (*gorm.DB, error) {
	dsn := os.Getenv("POSTGRES_URI")
	config := &gorm.Config{
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
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

	log.Println("PostgreSQL connection established")
	return db, nil
}

func initMongoDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}).
		SetTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("MongoDB connection established")
	return client.Database("orderdb"), nil
}
