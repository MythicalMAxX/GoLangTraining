package main

import (
	"fmt"
	"mypackage/config"
	"mypackage/internal/handlers"
	"mypackage/internal/repositories"
	"mypackage/internal/services"
	"mypackage/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	// Initialize repositories, services, and handlers
	memberRepo := repositories.NewMemberRepository(db)
	memberService := services.NewMemberService(memberRepo)
	memberHandler := handlers.NewMemberHandler(memberService)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, memberHandler)

	// Start the server
	router.Run(":8080")
}
