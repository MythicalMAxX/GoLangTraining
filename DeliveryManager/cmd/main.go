package main

import (
	"log"

	"delivery-service/config"
	"delivery-service/database"
	"delivery-service/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.Connect(cfg.DBUri)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Seed delivery partners
	if err := database.SeedDeliveryPartners(db); err != nil {
		log.Printf("Error seeding delivery partners: %v", err)
	}

	deliveryService := service.NewDeliveryService(db, cfg)
	if err := deliveryService.Start(); err != nil {
		log.Fatalf("Error starting delivery service: %v", err)
	}
}
