package main

import (
	"log"
	"mypackage/config"
	"mypackage/internal/repositories"
	"mypackage/internal/services"
	"mypackage/pkg/worker"
	"mypackage/routes"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize database connections
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	// Initialize repositories
	orderRepo := repositories.NewOrderRepository(db.PostgresDB)
	inventoryRepo := repositories.NewInventoryRepository(db.PostgresDB)
	userRepo := repositories.NewUserRepository(db.PostgresDB)
	orderHistoryRepo := repositories.NewOrderHistoryRepository(db.MongoDB)

	// Initialize services
	orderService := services.NewOrderService(
		db.PostgresDB,
		orderRepo,
		inventoryRepo,
		orderHistoryRepo,
	)

	inventoryService := services.NewInventoryService(inventoryRepo)
	userService := services.NewUserService(userRepo)

	// Initialize and start order worker
	orderWorker := worker.NewOrderWorker(orderService, 4)
	orderWorker.Start()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
        <-quit
        log.Println("Shutting down server...")
        orderWorker.Stop()
        log.Println("Server shutdown complete")
        os.Exit(0)
    }()

	// Setup router with all services
	r := routes.SetupRouter(orderService, inventoryService, userService)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
