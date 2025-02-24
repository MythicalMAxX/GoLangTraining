package handlers

import (
	"mypackage/pkg/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	orderClient *client.OrderClient
}

func NewUserHandler(orderClient *client.OrderClient) *UserHandler {
	return &UserHandler{
		orderClient: orderClient,
	}
}

func (h *UserHandler) CreateOrder(c *gin.Context) {
	var orderRequest client.CreateOrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add user ID from context
	orderRequest.UserID = c.GetString("user_id")

	if err := h.orderClient.CreateOrder(orderRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "order created successfully"})
}

func (h *UserHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	orders, err := h.orderClient.GetUserOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *UserHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id") // Get user ID from context

	if err := h.orderClient.CancelOrder(orderID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled successfully"})
}
