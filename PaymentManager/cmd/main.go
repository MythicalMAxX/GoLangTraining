package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"paymentservice/config"
	"paymentservice/consul"
	"paymentservice/db"
	"paymentservice/handlers"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize database
	db.InitDB(cfg)

	// Set up router
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Add health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// Routes
	r.POST("/payment", handlers.MakePayment)
	r.GET("/payment/:id", handlers.GetPaymentByID)

	// Register service with Consul
	consulClient, err := consul.RegisterService(cfg.ServiceName, cfg.ServiceName, cfg.ServicePort)
	if err != nil {
		log.Fatal(err)
	}

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := consul.DeregisterService(consulClient, cfg.ServiceName); err != nil {
			log.Printf("Error deregistering service: %v", err)
		}
		os.Exit(0)
	}()

	// Start server
	port := fmt.Sprintf(":%s", cfg.ServicePort)
	if err := r.Run(port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
