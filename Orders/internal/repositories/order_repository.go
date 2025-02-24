package repositories

import (
	"errors"
	"orderservice/internal/models"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByUserID(userID string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateStatus(id string, fromStatus, toStatus models.OrderStatus) error {
	if !toStatus.IsValid() {
		return errors.New("invalid order status")
	}

	result := r.db.Model(&models.Order{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Update("status", toStatus)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order status was already changed")
	}

	return nil
}

func (r *OrderRepository) FindByStatus(status models.OrderStatus) ([]models.Order, error) {
	if !status.IsValid() {
		return nil, errors.New("invalid order status")
	}

	var orders []models.Order
	err := r.db.Where("status = ?", status).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) ResetProcessingOrders() error {
	return r.db.Model(&models.Order{}).
		Where("status = ?", models.OrderStatusProcessing).
		Update("status", models.OrderStatusPending).
		Error
}
