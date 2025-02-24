package middleware

import (
    "net/http"
    "os"
    "strings"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
            c.Abort()
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
            c.Abort()
            return
        }

        tokenString := bearerToken[1]
        log.Printf("Received token: %s", tokenString) // Add this line to log the token

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET_KEY")), nil
        })

        if err != nil || !token.Valid {
            log.Printf("Token validation error: %v", err) // Add this line to log validation errors
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Log the parsed claims
        log.Printf("Token claims - UserID: %s, Role: %s", claims.UserID, claims.Role)

        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}