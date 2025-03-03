package handlers

import (
	"net/http"
	"paymentservice/db"
	"paymentservice/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRequest struct {
	OrderID       uuid.UUID            `json:"order_id" binding:"required"`
	PaymentMethod models.PaymentMethod `json:"payment_method" binding:"required"`
}

func MakePayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start database transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
		return
	}

	// Create payment
	payment := models.Payment{
		ID:            uuid.New(),
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		PaymentDate:   time.Now(),
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update order status
	updateResult := tx.Model(&models.Order{}).
		Where("id = ?", req.OrderID).
		Updates(map[string]interface{}{
			"payment_status": models.PaymentStatusPaid,
			"order_status":   models.OrderStatusProcessing,
			"updated_at":     time.Now(),
		})

	if updateResult.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}

	if updateResult.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func GetPaymentByID(c *gin.Context) {
	paymentID := c.Param("id")

	// Validate UUID format: must be exactly 36 characters (32 hex digits + 4 hyphens)
	if len(paymentID) != 36 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid UUID format",
			"details": "UUID must be in format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			"example": "123e4567-e89b-12d3-a456-426614174000",
		})
		return
	}

	id, err := uuid.Parse(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "invalid UUID format",
			"details":  err.Error(),
			"received": paymentID,
		})
		return
	}

	var payment models.Payment
	result := db.DB.First(&payment, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "payment not found",
				"id":    id,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
