package handlers

import (
    "mypackage/internal/models"
    "mypackage/internal/services"
    "net/http"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

type RegisterRequest struct {
    Name     string      `json:"name" binding:"required"`
    Email    string      `json:"email" binding:"required,email"`
    Password string      `json:"password" binding:"required,min=6"`
    Role     models.Role `json:"role"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.authService.Register(services.RegisterInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        Role:     req.Role,
    })
    if err != nil {
        // Check for specific errors and return appropriate status codes
        switch err.Error() {
        case "email already registered":
            c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "user registered successfully",
        "user": gin.H{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
            "role":  user.Role,
        },
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := h.authService.Login(services.LoginInput{
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
    })
}