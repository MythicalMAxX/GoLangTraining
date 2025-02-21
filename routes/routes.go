package routes

import (
	"mypackage/internal/handlers"
	"mypackage/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(orderService *services.OrderService, inventoryService *services.InventoryService, userService *services.UserService) *gin.Engine {
	r := gin.Default()

	// Initialize handlers with services
	orderHandler := handlers.NewOrderHandler(orderService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	userHandler := handlers.NewUserHandler(userService)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.GET("/:id/orders", orderHandler.GetUserOrders) // Changed from :userId to :id
		}

		// Inventory routes
		inventory := v1.Group("/inventory")
		{
			inventory.POST("", inventoryHandler.CreateInventory)
			inventory.GET("", inventoryHandler.GetInventory)
		}

		// Order routes
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.DELETE("/:id", orderHandler.CancelOrder)
			orders.GET("/:id/history", orderHandler.GetOrderHistory)
		}
	}

	return r
}
