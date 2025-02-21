package handlers

import (
	"mypackage/internal/models"
	"mypackage/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

type CreateOrderRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	InventoryID string  `json:"inventory_id" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}

	inventoryID, err := uuid.Parse(req.InventoryID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid inventory ID"})
		return
	}

	order := &models.Order{
		UserID:      userID,
		InventoryID: inventoryID,
		Quantity:    req.Quantity,
		Amount:      req.Amount,
	}

	if err := h.orderService.CreateOrder(c.Request.Context(), order); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrder(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	c.JSON(200, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.orderService.CancelOrder(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "order cancelled successfully"})
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.Param("id") // Changed from "userId" to "id"
	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	id := c.Param("id")
	history, err := h.orderService.GetOrderHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, history)
}
