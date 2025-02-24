package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"orderservice/config"
	"orderservice/internal/handlers"
	"orderservice/internal/repositories"
	"orderservice/internal/services"
	"orderservice/pkg/worker"
	"orderservice/routes"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongoDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetTLSConfig(&tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true, // Only for testing, remove in production
		}).
		SetTimeout(15 * time.Second).
		SetConnectTimeout(15 * time.Second).
		SetServerSelectionTimeout(15 * time.Second).
		SetDirect(true) // Try direct connection

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	// Verify connection with longer timeout
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("MongoDB connection established")
	return client.Database("orderdb"), nil
}

func main() {
	// Initialize databases
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	// Initialize repositories
	orderRepo := repositories.NewOrderRepository(db.PostgresDB)
	orderHistoryRepo := repositories.NewOrderHistoryRepository(db.MongoDB)

	// Initialize services
	orderService := services.NewOrderService(
		db.PostgresDB,
		orderRepo,
		orderHistoryRepo,
	)

	// Run migrations
	if err := orderService.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	// Reset any processing orders back to pending (in case of previous crash)
	if err := orderService.ResetProcessingOrders(); err != nil {
		log.Printf("Failed to reset processing orders: %v", err)
	}

	// Initialize and start order worker
	orderWorker := worker.NewOrderWorker(orderService, 4)
	orderWorker.Start()

	// Setup router
	router := routes.SetupRouter(orderHandler)

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		orderWorker.Stop()
		os.Exit(0)
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting Order service on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
