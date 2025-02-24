package routes

import (
	"orderservice/internal/handlers"
	"orderservice/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(orderHandler *handlers.OrderHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetUserOrders)
			orders.PATCH("/:id/cancel", orderHandler.CancelOrder)
			orders.GET("/:id/history", orderHandler.GetOrderHistory)
		}

		// User-specific order routes
		users := api.Group("/users")
		{
			users.GET("/:id/orders", orderHandler.GetUserOrders)
		}
	}

	return r
}
