package repositories

import (
	"mypackage/internal/models"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

type InventoryRepositoryInterface interface {
	GetAll() ([]models.Inventory, error)
	UpdateStock(id string, quantity int) error
	Create(inventory *models.Inventory) error
}

// Implement the interface methods
func (r *InventoryRepository) GetAll() ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.db.Find(&inventory).Error
	return inventory, err
}

func (r *InventoryRepository) UpdateStock(id string, quantity int) error {
	return r.db.Model(&models.Inventory{}).Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *InventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}
