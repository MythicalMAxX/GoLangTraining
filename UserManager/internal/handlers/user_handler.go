package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"userservices/internal/clients"
	"userservices/internal/models"
	"userservices/internal/services"
	"userservices/pkg/middleware"
	"log"
)

type UserHandler struct {
	userService *services.UserService
	authClient  *clients.AuthClient
}

func (h *UserHandler) GetAuthClient() *clients.AuthClient {
	return h.authClient
}

func NewUserHandler(userService *services.UserService) (*UserHandler, error) {
	authClient, err := clients.NewAuthClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %v", err)
	}

	return &UserHandler{
		userService: userService,
		authClient:  authClient,
	}, nil
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Login error: Invalid request body: %v", err)
        sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
        return
    }

    log.Printf("Login attempt for email: %s", req.Email)
    user, err := h.userService.Login(&req)
    if err != nil {
        log.Printf("Login error: Invalid credentials for email %s: %v", req.Email, err)
        sendResponse(w, "Invalid credentials", nil, http.StatusUnauthorized)
        return
    }

    log.Printf("User authenticated successfully. UserID: %s, Role: %s", user.ID, user.Role)
    token, err := h.authClient.GenerateToken(user.ID, user.Role)
    if err != nil {
        log.Printf("Token generation error for user %s: %v", user.ID, err)
        sendResponse(w, "Error generating token", nil, http.StatusInternalServerError)
        return
    }

    log.Printf("Token generated successfully for user %s", user.ID)
    loginResponse := models.LoginResponse{
        User: models.UserResponse{
            Name:  user.Username,
            Email: user.Email,
            UUID:  user.ID,
            Role:  user.Role,
        },
        Token: token,
    }

    sendResponse(w, "Login successful", loginResponse, http.StatusOK)
}

func (h *UserHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        sendResponse(w, "No token provided", nil, http.StatusUnauthorized)
        return
    }

    tokenParts := strings.Split(authHeader, " ")
    if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
        sendResponse(w, "Invalid token format", nil, http.StatusUnauthorized)
        return
    }

    // Validate token using AuthClient
    validationResp, err := h.authClient.ValidateToken(tokenParts[1])
    if err != nil {
        sendResponse(w, "Error validating token", nil, http.StatusInternalServerError)
        return
    }

    if validationResp.Body.IsValid != 1 {
        sendResponse(w, "Invalid token", nil, http.StatusUnauthorized)
        return
    }

    sendResponse(w, "Token is valid", validationResp.Body, http.StatusOK)
}

func sendResponse(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.Response{
		Message:    message,
		Body:       data,
		StatusCode: statusCode,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		sendResponse(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	userResponse := models.UserResponse{
		Name:  user.Username,
		Email: user.Email,
		UUID:  user.ID,
		Role:  user.Role,
	}

	sendResponse(w, "Registration successful", userResponse, http.StatusCreated)
}



// GetUsers handles the admin request to get all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Get user role from context
	role := r.Context().Value(middleware.UserRoleKey).(string)
	if role != "admin" {
		sendResponse(w, "Admin access required", nil, http.StatusForbidden)
		return
	}

	// Get all users
	users, err := h.userService.GetAllUsers()
	if err != nil {
		sendResponse(w, "Error getting users", nil, http.StatusInternalServerError)
		return
	}

	// Convert users to response format
	var usersResponse []models.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, models.UserResponse{
			Name:  user.Username,
			Email: user.Email,
			UUID:  user.ID,
			Role:  user.Role,
		})
	}

	sendResponse(w, "Users retrieved successfully", usersResponse, http.StatusOK)
}

func (h *UserHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
    // Check if user is admin
    userRole := middleware.GetUserRoleFromContext(r.Context())
    if userRole != "admin" {
        sendResponse(w, "Admin access required", nil, http.StatusForbidden)
        return
    }

    // Parse request body
    var req models.InventoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
        return
    }

    // Validate request
    if req.ProductName == "" || req.Stock < 0 || req.Price < 0 {
        sendResponse(w, "Invalid request parameters", nil, http.StatusBadRequest)
        return
    }

    // Make HTTP request to inventory service instead of using Kafka
    // This should be moved to a separate inventory client similar to AuthClient
    sendResponse(w, "Not implemented", nil, http.StatusNotImplemented)
}