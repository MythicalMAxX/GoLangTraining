package config

import (
	"fmt"
	"mypackage/migrations"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// Tell Viper to look for .env file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // Read environment variables

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}
}

// DBConfig holds database configuration
type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	URI             string
}

// GetDBConfig loads database configuration from environment
func GetDBConfig() DBConfig {
	return DBConfig{
		MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONN"),
		MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONN"),
		ConnMaxLifetime: time.Duration(viper.GetInt("DB_CONN_MAX_LIFE_TIME")) * time.Second,
		ConnMaxIdleTime: time.Duration(viper.GetInt("DB_CONN_MAX_IDLE_TIME")) * time.Second,
		URI:             viper.GetString("URI"),
	}
}

// InitDB initializes the database connection with proper connection pooling
func InitDB() (*gorm.DB, error) {
	config := GetDBConfig()

	if config.URI == "" {
		return nil, fmt.Errorf("database URI is not set")
	}

	// Initialize GORM DB connection
	db, err := gorm.Open(postgres.Open(config.URI), &gorm.Config{
		PrepareStmt: true, // Enable statement cache
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	// Get underlying SQL DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return db, nil
}
