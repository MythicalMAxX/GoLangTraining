package repositories

import (
	"errors"
	"mypackage/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type OrderRepositoryInterface interface {
	Create(order *models.Order) error
	FindByID(id string) (*models.Order, error)
	FindByUserID(userID string) ([]models.Order, error)
	UpdateStatus(id string, status models.OrderStatus) error
	FindByStatus(status models.OrderStatus) ([]models.Order, error)
}

// Implement the interface methods
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("Inventory").
		Where("orders.id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByUserID(userID string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("User").Preload("Inventory").
		Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateStatus(id string, status models.OrderStatus) error {
	if !status.IsValid() {
		return errors.New("invalid order status")
	}

	return r.db.Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *OrderRepository) FindByStatus(status models.OrderStatus) ([]models.Order, error) {
	if !status.IsValid() {
		return nil, errors.New("invalid order status")
	}

	var orders []models.Order
	err := r.db.Preload("User").Preload("Inventory").
		Where("status = ?", status).
		Find(&orders).Error
	return orders, err
}
