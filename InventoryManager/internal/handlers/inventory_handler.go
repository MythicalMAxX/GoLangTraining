package handlers

import (
	"fmt"
	"inventorymanager/internal/models"
	"inventorymanager/internal/repositories"
	"inventorymanager/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryHandler struct {
	service *services.InventoryService
}

func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) Create(c *gin.Context) {
	var inventory models.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateInventory(&inventory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inventory)
}

func (h *InventoryHandler) GetAll(c *gin.Context) {
	var status *models.InventoryStatus
	if s := c.Query("status"); s != "" {
		temp := models.InventoryStatus(s)
		status = &temp
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	productName := c.Query("product_name")

	filter := repositories.FilterOptions{
		Status:      status,
		ProductName: &productName,
		Page:        page,
		PageSize:    pageSize,
	}

	inventories, total, err := h.service.GetAllInventory(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      inventories,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *InventoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid UUID format",
			"details": err.Error(),
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate allowed fields
	allowedFields := map[string]bool{
		"stock":  true,
		"price":  true,
		"status": true,
	}

	for field := range updates {
		if !allowedFields[field] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          fmt.Sprintf("Field '%s' is not allowed for update", field),
				"allowed_fields": []string{"stock", "price", "status"},
			})
			return
		}
	}

	if err := h.service.UpdateInventory(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update inventory",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated successfully",
	})
}

func (h *InventoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid UUID format",
			"details": err.Error(),
		})
		return
	}

	inventory, err := h.service.GetInventoryByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, inventory)
}
