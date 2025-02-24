package routes

import (
	"mypackage/internal/handlers"
	"mypackage/internal/services"
	"mypackage/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(orderService *services.OrderService, inventoryService *services.InventoryService, userService *services.UserService) *gin.Engine {
	r := gin.Default()

	// Add logging middleware
	r.Use(gin.Logger())
	r.Use(middleware.RateLimitMiddleware())

	// Initialize handlers with services
	orderHandler := handlers.NewOrderHandler(orderService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	userHandler := handlers.NewUserHandler(userService)

	// Public routes (no auth required)
	public := r.Group("/api/v1")
	{
		public.POST("/login", userHandler.Login)
		public.POST("/users", userHandler.CreateUser) // Allow user creation without auth
	}

	// Protected routes (auth required)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware()) // Adds auth middleware
	{
		// User routes
		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.GET("/:id/orders", orderHandler.GetUserOrders)
		}

		// Inventory routes
		inventory := protected.Group("/inventory")
		{
			inventory.POST("", inventoryHandler.CreateInventory)
			inventory.GET("", inventoryHandler.GetInventory)
		}

		// Order routes
		orders := protected.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PATCH("/:id", orderHandler.CancelOrder)
			orders.GET("/:id/history", orderHandler.GetOrderHistory)
		}
	}

	return r
}
