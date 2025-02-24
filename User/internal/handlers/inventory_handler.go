package handlers

import (
    "net/http"
    "mypackage/internal/models"
    "mypackage/internal/services"

    "github.com/gin-gonic/gin"
)

type InventoryHandler struct {
    service *services.InventoryService
}

func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
    return &InventoryHandler{service: service}
}

func (h *InventoryHandler) CreateInventory(c *gin.Context) {
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

func (h *InventoryHandler) GetInventory(c *gin.Context) {
    inventory, err := h.service.GetAllInventory()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, inventory)
}