package routes

import (
	"orderservice/internal/handlers"
	"orderservice/internal/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *mux.Router, orderHandler *handlers.OrderHandler) {
	// API version prefix
	api := router.PathPrefix("/api/v1").Subrouter()

	// Orders routes
	orders := api.PathPrefix("/orders").Subrouter()
	orders.Use(middleware.LoggingMiddleware)
	orders.Use(middleware.JSONMiddleware)

	orders.HandleFunc("", orderHandler.CreateOrder).Methods("POST")
	orders.HandleFunc("/{id}", orderHandler.GetOrder).Methods("GET")
	orders.HandleFunc("/{id}/cancel", orderHandler.CancelOrder).Methods("PUT")

	// Health check
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
}
