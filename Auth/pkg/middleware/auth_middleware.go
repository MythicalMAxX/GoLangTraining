package middleware

import (
    "mypackage/internal/models"
    "mypackage/pkg/jwt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware(allowedRoles ...models.Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
            c.Abort()
            return
        }

        claims, err := jwt.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // If specific roles are required, check them
        if len(allowedRoles) > 0 {
            roleMatched := false
            for _, role := range allowedRoles {
                if claims.Role == role {
                    roleMatched = true
                    break
                }
            }
            if !roleMatched {
                c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
                c.Abort()
                return
            }
        }

        // Set claims in context
        c.Set("userID", claims.UserID)
        c.Set("userRole", claims.Role)
        c.Next()
    }
}