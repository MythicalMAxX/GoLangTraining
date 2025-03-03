package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"userservices/internal/clients"
	"userservices/internal/models"
	"userservices/internal/services"
	"userservices/pkg/middleware"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService        *services.UserService
	authClient         *clients.AuthClient
	inventoryClient    *clients.InventoryClient
	orderClient        *clients.OrderClient
	paymentClient      *clients.PaymentClient
	notificationClient *clients.NotificationClient
}

func (h *UserHandler) GetAuthClient() *clients.AuthClient {
	return h.authClient
}

func NewUserHandler(userService *services.UserService) (*UserHandler, error) {
	authClient, err := clients.NewAuthClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %v", err)
	}

	inventoryClient, err := clients.NewInventoryClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory client: %v", err)
	}

	orderClient, err := clients.NewOrderClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create order client: %v", err)
	}

	paymentClient, err := clients.NewPaymentClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create payment client: %v", err)
	}

	notificationClient, err := clients.NewNotificationClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification client: %v", err)
	}

	return &UserHandler{
		userService:        userService,
		authClient:         authClient,
		inventoryClient:    inventoryClient,
		orderClient:        orderClient,
		paymentClient:      paymentClient,
		notificationClient: notificationClient,
	}, nil
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.notificationClient.NotifyError("user", "login validation", err)
		sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(&req)
	if err != nil {
		h.notificationClient.NotifyError("user", fmt.Sprintf("login attempt for %s", req.Email), err)
		sendResponse(w, "Invalid credentials", nil, http.StatusUnauthorized)
		return
	}

	token, err := h.authClient.GenerateToken(user.ID, user.Role)
	if err != nil {
		h.notificationClient.NotifyError("auth", "token generation", err)
		sendResponse(w, "Error generating token", nil, http.StatusInternalServerError)
		return
	}

	h.notificationClient.NotifySuccess("user", fmt.Sprintf("user logged in: %s", user.Email))
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

	// Check IsValid directly from ValidationResponse since it's not nested
	if validationResp.IsValid != 1 {
		sendResponse(w, "Invalid token", nil, http.StatusUnauthorized)
		return
	}

	// Create response body
	responseBody := models.ValidationBody{
		IsValid: validationResp.IsValid,
		UUID:    validationResp.UUID,
		Role:    validationResp.Role,
	}

	sendResponse(w, "Token is valid", responseBody, http.StatusOK)
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
		h.notificationClient.NotifyError("user", "registration validation", err)
		sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		h.notificationClient.NotifyError("user", "registration", err)
		sendResponse(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	userResponse := models.UserResponse{
		Name:  user.Username,
		Email: user.Email,
		UUID:  user.ID,
		Role:  user.Role,
	}

	h.notificationClient.NotifySuccess("user", fmt.Sprintf("new user registered: %s", user.Email))
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
	var req models.InventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.notificationClient.NotifyError("inventory", "create validation", err)
		sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	req.UserID = middleware.GetUserUUIDFromContext(r.Context())
	inventory, err := h.inventoryClient.CreateInventory(&req)
	if err != nil {
		h.notificationClient.NotifyError("inventory", "create", err)
		sendResponse(w, "Failed to create inventory", err.Error(), http.StatusInternalServerError)
		return
	}

	h.notificationClient.NotifySuccess("inventory", fmt.Sprintf("created new inventory: %s", inventory.ID))
	sendResponse(w, "Inventory created successfully", inventory, http.StatusCreated)
}

func (h *UserHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	inventory, err := h.inventoryClient.GetInventoryByID(id)
	if err != nil {
		if err.Error() == "inventory not found" {
			sendResponse(w, "Inventory not found", nil, http.StatusNotFound)
			return
		}
		sendResponse(w, "Failed to get inventory", err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, "Inventory retrieved successfully", inventory, http.StatusOK)
}

func (h *UserHandler) GetAllInventory(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10
	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			page = val
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil {
			pageSize = val
		}
	}

	status := r.URL.Query().Get("status")
	productName := r.URL.Query().Get("product_name")

	response, err := h.inventoryClient.GetAllInventory(page, pageSize, status, productName)
	if err != nil {
		sendResponse(w, "Failed to get inventories", err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, "Inventories retrieved successfully", response, http.StatusOK)
}

func (h *UserHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	if err := h.inventoryClient.UpdateInventory(id, updates); err != nil {
		sendResponse(w, "Failed to update inventory", err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, "Updated successfully", nil, http.StatusOK)
}

func (h *UserHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.notificationClient.NotifyError("order", "create validation", err)
		sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserUUIDFromContext(r.Context())
	inventory, err := h.inventoryClient.GetInventoryByID(req.InventoryID.String())
	if err != nil {
		h.notificationClient.NotifyError("order", "fetch inventory", err)
		sendResponse(w, "Failed to get inventory details", err.Error(), http.StatusInternalServerError)
		return
	}

	order, err := h.orderClient.CreateOrder(userID, inventory, req.OrderCount)
	if err != nil {
		h.notificationClient.NotifyError("order", "create", err)
		sendResponse(w, "Failed to create order", err.Error(), http.StatusInternalServerError)
		return
	}

	h.notificationClient.NotifySuccess("order", fmt.Sprintf("created new order: %s", order.ID))
	sendResponse(w, "Order created successfully", order, http.StatusCreated)
}

func (h *UserHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	log.Printf("Getting order with ID: %s", id)

	// Validate UUID format first
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid UUID format: %s (Error: %v)", id, err)
		sendResponse(w, "Invalid order ID format: must be a valid UUID v4", nil, http.StatusBadRequest)
		return
	}

	order, err := h.orderClient.GetOrder(parsedUUID.String())
	if err != nil {
		log.Printf("Error getting order %s: %v", id, err)
		switch {
		case strings.Contains(err.Error(), "invalid order ID"):
			sendResponse(w, "Invalid order ID", nil, http.StatusBadRequest)
		case err.Error() == "order not found":
			sendResponse(w, "Order not found", nil, http.StatusNotFound)
		default:
			sendResponse(w, "Failed to get order", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	sendResponse(w, "Order retrieved successfully", order, http.StatusOK)
}

func (h *UserHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])

	if _, err := uuid.Parse(id); err != nil {
		h.notificationClient.NotifyError("order", "cancel validation", fmt.Errorf("invalid UUID: %s", id))
		sendResponse(w, "Invalid order ID format: must be a valid UUID v4", nil, http.StatusBadRequest)
		return
	}

	if err := h.orderClient.CancelOrder(id); err != nil {
		h.notificationClient.NotifyError("order", fmt.Sprintf("cancel order %s", id), err)
		// ... error handling ...
		return
	}

	h.notificationClient.NotifySuccess("order", fmt.Sprintf("cancelled order: %s", id))
	sendResponse(w, "Order cancelled successfully", nil, http.StatusOK)
}

func (h *UserHandler) MakePayment(w http.ResponseWriter, r *http.Request) {
	var req models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.notificationClient.NotifyError("payment", "validation", err)
		sendResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	if req.PaymentMethod != models.Card && req.PaymentMethod != models.Cash {
		err := fmt.Errorf("invalid payment method: %s", req.PaymentMethod)
		h.notificationClient.NotifyError("payment", "method validation", err)
		sendResponse(w, "Invalid payment method. Must be 'card' or 'cash'", nil, http.StatusBadRequest)
		return
	}

	payment, err := h.paymentClient.MakePayment(&req)
	if err != nil {
		h.notificationClient.NotifyError("payment", fmt.Sprintf("process payment for order %s", req.OrderID), err)
		sendResponse(w, "Failed to process payment", err.Error(), http.StatusInternalServerError)
		return
	}

	h.notificationClient.NotifySuccess("payment", fmt.Sprintf("processed payment for order %s", req.OrderID))
	sendResponse(w, "Payment processed successfully", payment, http.StatusCreated)
}

func (h *UserHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	log.Printf("Getting payment with ID: %s", id)

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		log.Printf("Invalid UUID format: %s", id)
		sendResponse(w, "Invalid payment ID format", nil, http.StatusBadRequest)
		return
	}

	payment, err := h.paymentClient.GetPayment(id)
	if err != nil {
		log.Printf("Error getting payment: %v", err)
		switch {
		case err.Error() == "payment not found":
			sendResponse(w, "Payment not found", nil, http.StatusNotFound)
		default:
			sendResponse(w, "Failed to get payment details", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	sendResponse(w, "Payment details retrieved successfully", payment, http.StatusOK)
}
