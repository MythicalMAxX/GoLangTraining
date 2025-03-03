package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"orderservice/config"
	// "orderservice/config"
	"orderservice/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get environment variables
	dbURI := os.Getenv("USER_DB_URI")
	serviceName := os.Getenv("SERVICE_NAME")
	servicePort := os.Getenv("SERVICE_PORT")
	rabbitMQUrl := os.Getenv("RABBITMQ_URL")

	if servicePort == "" {
		servicePort = "8084" // default port
	}

	// Initialize database
	database, err := config.InitDB(dbURI)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to RabbitMQ
	conn, err := config.ConnectToRabbitMQ(rabbitMQUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Setup RabbitMQ exchanges and queues
	err = config.SetupRabbitMQ(ch)
	if err != nil {
		log.Fatal(err)
	}

	// Create handler
	orderHandler := handlers.NewOrderHandler(database, ch)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/orders", orderHandler.CreateOrder).Methods("POST")
	r.HandleFunc("/api/v1/orders", orderHandler.ListOrders).Methods("GET")
	r.HandleFunc("/api/v1/orders/{id}", orderHandler.GetOrder).Methods("GET")
	r.HandleFunc("/api/v1/orders/{id}", orderHandler.UpdateOrder).Methods("PUT")
	r.HandleFunc("/api/v1/orders/{id}", orderHandler.DeleteOrder).Methods("DELETE")
	// Add new route for restaurant orders
	r.HandleFunc("/api/v1/restaurants/{restaurant_id}/orders", orderHandler.GetOrdersByRestaurant).Methods("GET")

	// Start server
	addr := fmt.Sprintf(":%s", servicePort)
	log.Printf("Starting %s on %s", serviceName, addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
