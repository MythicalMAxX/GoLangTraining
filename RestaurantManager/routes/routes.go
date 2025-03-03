package routes

import (
	"github.com/gin-gonic/gin"
	"restaurantservice/internal/handlers"
)

func SetupRoutes(r *gin.Engine, handler *handlers.RestaurantHandler) {
	api := r.Group("/api/v1")
	{
		restaurants := api.Group("/restaurants")
		{
			restaurants.POST("/", handler.CreateRestaurant)
			restaurants.GET("/:id", handler.GetRestaurant)
			restaurants.PUT("/:id", handler.UpdateRestaurant)
			restaurants.DELETE("/:id", handler.DeleteRestaurant)
			restaurants.GET("/", handler.ListRestaurants)
		}
	}
}
