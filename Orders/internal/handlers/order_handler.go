package handlers

import (
    "net/http"
    "orderservice/internal/models"
    "orderservice/internal/services"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type OrderHandler struct {
    orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
    return &OrderHandler{orderService: orderService}
}

type CreateOrderRequest struct {
    InventoryID string  `json:"inventory_id" binding:"required"`
    Quantity    int     `json:"quantity" binding:"required"`
    Amount      float64 `json:"amount" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var req CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, err := uuid.Parse(c.GetString("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    inventoryID, err := uuid.Parse(req.InventoryID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory ID"})
        return
    }

    order := &models.Order{
        UserID:      userID,
        InventoryID: inventoryID,
        Quantity:    req.Quantity,
        Amount:      req.Amount,
        Status:      models.OrderStatusPending,
    }

    if err := h.orderService.CreateOrder(c, order); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
    userID := c.GetString("user_id")
    orders, err := h.orderService.GetUserOrders(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
    orderID := c.Param("id")
    userID := c.GetString("user_id")

    if err := h.orderService.CancelOrder(c, orderID, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "order cancelled successfully"})
}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
    orderID := c.Param("id")
    history, err := h.orderService.GetOrderHistory(c, orderID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, history)
}