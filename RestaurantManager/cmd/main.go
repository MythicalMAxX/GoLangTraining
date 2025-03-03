package main

import (
	"log"
	"os"

	"restaurantservice/config"
	"restaurantservice/internal/handlers"
	"restaurantservice/internal/repositories"
	"restaurantservice/internal/services"
	"restaurantservice/pkg/messaging"
	"restaurantservice/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := config.ConnectDB()

	// Initialize RabbitMQ client with restaurant ID
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	// Generate or get restaurant ID
	restaurantID := uuid.New() // In practice, you'd get this from your config/database

	rabbitmq, err := messaging.NewRabbitMQClient(rabbitmqURL, restaurantID)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitmq.Close()
	log.Printf("Connected to RabbitMQ at %s", rabbitmqURL)

	// Initialize dependencies
	restaurantRepo := repositories.NewRestaurantRepository(db)
	restaurantService := services.NewRestaurantService(restaurantRepo, rabbitmq)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService)

	// Start consuming orders
	err = rabbitmq.ConsumeOrders(restaurantService.HandleNewOrder)
	if err != nil {
		log.Fatal("Failed to start order consumer:", err)
	}
	log.Println("Restaurant service is ready to process orders")

	// Setup Gin router
	r := gin.Default()
	routes.SetupRoutes(r, restaurantHandler)

	// Start server
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8083"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
