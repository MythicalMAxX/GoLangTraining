package services

import (
	"mypackage/internal/models"
	"mypackage/internal/repositories"
)

type InventoryService struct {
	repo repositories.InventoryRepositoryInterface
}

func NewInventoryService(repo repositories.InventoryRepositoryInterface) *InventoryService {
	return &InventoryService{repo: repo}
}

type InventoryServiceInterface interface {
	GetInventory() ([]models.Inventory, error)
	CreateInventory(inventory *models.Inventory) error
	UpdateStock(id string, quantity int) error
}

func (s *InventoryService) GetInventory() ([]models.Inventory, error) {
	return s.repo.GetAll()
}

func (s *InventoryService) CreateInventory(inventory *models.Inventory) error {
	return s.repo.Create(inventory)
}

func (s *InventoryService) UpdateStock(id string, quantity int) error {
	return s.repo.UpdateStock(id, quantity)
}
