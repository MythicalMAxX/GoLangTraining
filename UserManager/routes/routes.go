package routes

import (
	"userservices/internal/handlers"
	"userservices/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	// Public routes
	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Protected routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(userHandler.GetAuthClient()))

	api.HandleFunc("/validate", userHandler.ValidateToken).Methods("GET")

	// Admin only routes
	admin := api.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/users", middleware.AdminOnly(userHandler.GetUsers)).Methods("GET")

	inventory := api.PathPrefix("/inventory").Subrouter()
    inventory.HandleFunc("", middleware.AdminOnly(userHandler.CreateInventory)).Methods("POST")
    inventory.HandleFunc("", userHandler.GetAllInventory).Methods("GET")
    inventory.HandleFunc("/{id}", userHandler.GetInventoryByID).Methods("GET")
    inventory.HandleFunc("/{id}", middleware.AdminOnly(userHandler.UpdateInventory)).Methods("PATCH")

	orders := api.PathPrefix("/orders").Subrouter()
    orders.HandleFunc("", userHandler.CreateOrder).Methods("POST")
    orders.HandleFunc("/{id}", userHandler.GetOrder).Methods("GET")
    orders.HandleFunc("/{id}/cancel", userHandler.CancelOrder).Methods("PUT")

	payments := api.PathPrefix("/payments").Subrouter()
    payments.HandleFunc("", userHandler.MakePayment).Methods("POST")
    payments.HandleFunc("/{id}", userHandler.GetPayment).Methods("GET")
}
