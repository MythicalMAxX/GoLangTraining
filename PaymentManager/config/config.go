package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUri           string
	ServiceName     string
	ServicePort     string
	OrderServiceURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBUri:           os.Getenv("USER_DB_URI"),
		ServiceName:     os.Getenv("SERVICE_NAME"),
		ServicePort:     os.Getenv("SERVICE_PORT"),
		OrderServiceURL: os.Getenv("ORDER_SERVICE_URL"),
	}
}
