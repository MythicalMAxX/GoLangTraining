package main

import (
	"inventorymanager/config"
	"inventorymanager/internal/handlers"
	"inventorymanager/internal/repositories"
	"inventorymanager/internal/services"
	"inventorymanager/pkg/discovery"
	"inventorymanager/routes"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB
	db := config.InitDB()

	// Initialize Repository
	repo := repositories.NewInventoryRepository(db)

	// Initialize Service Discovery
	discovery, err := discovery.NewServiceDiscovery(os.Getenv("CONSUL_HTTP_ADDR"))
	if err != nil {
		log.Fatal("Failed to create service discovery:", err)
	}

	// Initialize Service
	service := services.NewInventoryService(repo)

	// Initialize Handler
	handler := handlers.NewInventoryHandler(service)

	// Initialize Gin Router
	r := gin.Default()

	// Add health check endpoint for Consul
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Setup routes
	routes.SetupRoutes(r, handler)

	// Get service port
	port, err := strconv.Atoi(os.Getenv("SERVICE_PORT"))
	if err != nil {
		log.Fatal("Invalid port:", err)
	}

	// Register service with Consul
	err = discovery.Register(
		os.Getenv("SERVICE_NAME"),
		os.Getenv("SERVICE_NAME"),
		"localhost",
		port,
	)
	if err != nil {
		log.Fatal("Failed to register service:", err)
	}

	log.Printf("Service '%s' successfully registered with Consul at %s\n",
		os.Getenv("SERVICE_NAME"),
		os.Getenv("CONSUL_HTTP_ADDR"))

	// Start server
	r.Run(":" + os.Getenv("SERVICE_PORT"))
}
