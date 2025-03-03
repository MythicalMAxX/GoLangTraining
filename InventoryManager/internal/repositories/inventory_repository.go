package repositories

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "inventorymanager/internal/models"
)

type InventoryRepository struct {
    db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
    return &InventoryRepository{db: db}
}

type FilterOptions struct {
    Status      *models.InventoryStatus
    ProductName *string
    Page        int
    PageSize    int
}

func (r *InventoryRepository) Create(inventory *models.Inventory) error {
    return r.db.Create(inventory).Error
}

func (r *InventoryRepository) GetByID(id uuid.UUID) (*models.Inventory, error) {
    var inventory models.Inventory
    err := r.db.Where("id = ? AND status != ?", id, models.StatusDeleted).First(&inventory).Error
    return &inventory, err
}

func (r *InventoryRepository) Update(inventory *models.Inventory) error {
    return r.db.Save(inventory).Error
}

func (r *InventoryRepository) GetAll(filter FilterOptions) ([]models.Inventory, int64, error) {
    var inventories []models.Inventory
    var total int64
    
    query := r.db.Model(&models.Inventory{})
    
    if filter.Status != nil {
        query = query.Where("status = ?", *filter.Status)
    }
    
    if filter.ProductName != nil {
        query = query.Where("product_name LIKE ?", "%"+*filter.ProductName+"%")
    }
    
    // Get total count
    query.Count(&total)
    
    // Apply pagination
    offset := (filter.Page - 1) * filter.PageSize
    err := query.Offset(offset).Limit(filter.PageSize).Find(&inventories).Error
    
    return inventories, total, err
}