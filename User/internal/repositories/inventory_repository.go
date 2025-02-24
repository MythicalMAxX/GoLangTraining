package repositories

import (
    "mypackage/internal/models"

    // "github.com/google/uuid"
    "gorm.io/gorm"
)

type InventoryRepository struct {
    db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
    return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(inventory *models.Inventory) error {
    return r.db.Create(inventory).Error
}

func (r *InventoryRepository) GetAll() ([]models.Inventory, error) {
    var inventories []models.Inventory
    err := r.db.Find(&inventories).Error
    return inventories, err
}

func (r *InventoryRepository) UpdateStock(id string, quantity int) error {
    return r.db.Model(&models.Inventory{}).
        Where("id = ?", id).
        Update("stock", gorm.Expr("stock + ?", quantity)).Error
}