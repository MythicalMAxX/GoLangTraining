package services

import (
	"fmt"
	"inventorymanager/internal/models"
	"inventorymanager/internal/repositories"

	"github.com/google/uuid"
)

type InventoryService struct {
	repo *repositories.InventoryRepository
}

func NewInventoryService(repo *repositories.InventoryRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (s *InventoryService) CreateInventory(inventory *models.Inventory) error {
	if inventory.Status == "" {
		inventory.Status = models.StatusInStock
	}
	return s.repo.Create(inventory)
}

func (s *InventoryService) UpdateInventory(id uuid.UUID, updates map[string]interface{}) error {
	inventory, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	for key, value := range updates {
		switch key {
		case "stock":
			if floatVal, ok := value.(float64); ok {
				inventory.Stock = int(floatVal)
			} else {
				return fmt.Errorf("invalid type for stock: expected number")
			}
		case "price":
			if floatVal, ok := value.(float64); ok {
				inventory.Price = floatVal
			} else {
				return fmt.Errorf("invalid type for price: expected number")
			}
		case "status":
			if strVal, ok := value.(string); ok {
				inventory.Status = models.InventoryStatus(strVal)
			} else {
				return fmt.Errorf("invalid type for status: expected string")
			}
		default:
			return fmt.Errorf("unknown field: %s", key)
		}
	}

	return s.repo.Update(inventory)
}

func (s *InventoryService) GetAllInventory(filter repositories.FilterOptions) ([]models.Inventory, int64, error) {
	return s.repo.GetAll(filter)
}

func (s *InventoryService) GetInventoryByID(id uuid.UUID) (*models.Inventory, error) {
	return s.repo.GetByID(id)
}
