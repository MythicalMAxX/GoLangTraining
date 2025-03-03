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
}
