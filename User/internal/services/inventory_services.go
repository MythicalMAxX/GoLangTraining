package services

import (
	"mypackage/internal/models"
	"mypackage/internal/repositories"
)

type InventoryService struct {
	repo *repositories.InventoryRepository
}

func NewInventoryService(repo *repositories.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) CreateInventory(inventory *models.Inventory) error {
	return s.repo.Create(inventory)
}

func (s *InventoryService) GetAllInventory() ([]models.Inventory, error) {
	return s.repo.GetAll()
}

func (s *InventoryService) UpdateStock(id string, quantity int) error {
	return s.repo.UpdateStock(id, quantity)
}
