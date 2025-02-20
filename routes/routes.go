package routes

import (
	"mypackage/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, memberHandler *handlers.MemberHandler) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/running", handlers.GetStatus)
		v1.GET("/member/:id", memberHandler.Get)
		v1.POST("/member/register", memberHandler.Register)
		v1.DELETE("/member/:id", memberHandler.Delete)
		v1.PATCH("/member/:id", memberHandler.Modify)
		v1.PUT("/member/:id", memberHandler.Update)
	}
}
