package main

import (
	"log"
	"mypackage/config"
	handlers "mypackage/internal/handlers"
	"mypackage/internal/repositories"
	"mypackage/internal/services"
	"mypackage/pkg/client"
	"mypackage/routes"
)

func main() {
	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	inventoryRepo := repositories.NewInventoryRepository(db)

	// Initialize services
	inventoryService := services.NewInventoryService(inventoryRepo)
	orderClient := client.NewOrderClient()

	// Initialize handlers
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	userHandler := handlers.NewUserHandler(orderClient)

	// Setup router
	router := routes.SetupRouter(userHandler, inventoryHandler)

	// Start server
	log.Println("Starting User service on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
