package routes

import (
    "mypackage/internal/handlers"
    "mypackage/pkg/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler, inventoryHandler *handlers.InventoryHandler) *gin.Engine {
    r := gin.Default()

    api := r.Group("/api/v1")
    api.Use(middleware.AuthMiddleware())
    {
        // User routes
        users := api.Group("/users")
        {
            users.POST("/orders", userHandler.CreateOrder)
            users.GET("/orders", userHandler.GetUserOrders)
            users.PATCH("/orders/:id/cancel", userHandler.CancelOrder)
        }

        // Inventory routes (admin only)
        inventory := api.Group("/inventory")
        inventory.Use(middleware.AdminAuthMiddleware())
        {
            inventory.POST("", inventoryHandler.CreateInventory)
            inventory.GET("", inventoryHandler.GetInventory)
        }
    }

    return r
}