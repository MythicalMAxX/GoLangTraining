package routes

import (
    "mypackage/internal/handlers"
    "mypackage/internal/models"
    "mypackage/pkg/middleware"

    "github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
    r := gin.Default()

    // Public routes
    api := r.Group("/api/v1")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
        }
    }

    // Protected routes example (requiring admin role)
    protected := api.Group("/admin")
    protected.Use(middleware.AuthMiddleware(models.RoleAdmin))
    {
        // Admin-only endpoints here
    }

    return r
}