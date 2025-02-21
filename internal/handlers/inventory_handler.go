package handlers

import (
	"mypackage/internal/models"
	"mypackage/internal/services"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

type CreateInventoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Stock int    `json:"stock" binding:"required"`
}

func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	inventory := &models.Inventory{
		Name:  req.Name,
		Stock: req.Stock,
	}

	if err := h.inventoryService.CreateInventory(inventory); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, inventory)
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	inventory, err := h.inventoryService.GetInventory()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, inventory)
}
