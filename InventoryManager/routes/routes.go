// routes/routes.go
package routes

import (
	"inventorymanager/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, handler *handlers.InventoryHandler) {
	v1 := r.Group("/api/v1")
	{
		inventory := v1.Group("/inventory")
		{
			inventory.POST("/", handler.Create)
			inventory.GET("/", handler.GetAll)
			inventory.GET("/:id", handler.GetByID) 
			inventory.PATCH("/:id", handler.Update)
		}
	}
}
