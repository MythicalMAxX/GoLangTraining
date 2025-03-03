package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"orderservice/internal/database"
	"orderservice/internal/handlers"
	"orderservice/internal/service"
	"orderservice/internal/worker"
	"orderservice/routes"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
)

func registerConsul(port string) {
    consulAddress := os.Getenv("CONSUL_HTTP_ADDR")
    if consulAddress == "" {
        consulAddress = "localhost:8500"
    }

    // Create Consul client
    config := api.DefaultConfig()
    config.Address = consulAddress
    consulClient, err := api.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create Consul client: %v", err)
    }

    // Convert port to integer
    portInt, err := strconv.Atoi(port)
    if err != nil {
        log.Fatalf("Invalid SERVICE_PORT: %v", err)
    }

    // Use service name directly from environment variable
    serviceName := os.Getenv("SERVICE_NAME")
    registration := &api.AgentServiceRegistration{
        ID:      serviceName,               // Using just the service name as ID
        Name:    serviceName,
        Address: "localhost",
        Port:    portInt,
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://localhost:%s/health", port),
            Interval: "10s",
            Timeout:  "1s",
        },
    }

    if err = consulClient.Agent().ServiceRegister(registration); err != nil {
        log.Fatalf("Failed to register with Consul: %v", err)
    }
    log.Printf("Service registered with Consul: %s", serviceName)
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize database and run migrations
	db := database.InitDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}
	defer db.Close()

	// Run migrations immediately after DB connection
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services and handlers
	orderService := service.NewOrderService(db)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Initialize and start worker pool
	workerConfig := worker.DefaultConfig()
	workerPool := worker.NewWorkerPool(workerConfig, orderService)
	workerPool.Start()
	defer workerPool.Stop()

	// Initialize router
	router := mux.NewRouter()

	// Setup routes
	routes.SetupRoutes(router, orderHandler)

	// Start server
	port := os.Getenv("SERVICE_PORT")
	log.Printf("Starting server on port %s", port)

	// Register service with Consul
	registerConsul(port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
