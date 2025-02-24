package services

import (
	"context"
	"errors"
	"orderservice/internal/models"
	"orderservice/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	db               *gorm.DB
	orderRepo        *repositories.OrderRepository
	orderHistoryRepo *repositories.OrderHistoryRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo *repositories.OrderRepository,
	orderHistoryRepo *repositories.OrderHistoryRepository,
) *OrderService {
	return &OrderService{
		db:               db,
		orderRepo:        orderRepo,
		orderHistoryRepo: orderHistoryRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create order
	order.Status = models.OrderStatusPending
	if err := s.orderRepo.Create(order); err != nil {
		tx.Rollback()
		return err
	}

	// Create order history
	history := &models.OrderHistory{
		OrderID: order.ID.String(),
		Status:  models.OrderStatusPending,
		Log: []models.OrderLog{{
			Status:    models.OrderStatusPending,
			Timestamp: time.Now(),
		}},
	}

	if err := s.orderHistoryRepo.CreateOrderHistory(ctx, history); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *OrderService) AutoMigrate() error {
	return s.db.AutoMigrate(&models.Order{})
}

func (s *OrderService) GetUserOrders(userID string) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID string, userID string) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if order.UserID.String() != userID {
		return errors.New("unauthorized to cancel this order")
	}

	if order.Status == models.OrderStatusCanceled {
		return errors.New("order is already canceled")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.orderRepo.UpdateStatus(orderID, order.Status, models.OrderStatusCanceled); err != nil {
		tx.Rollback()
		return err
	}

	log := models.OrderLog{
		Status:    models.OrderStatusCanceled,
		Timestamp: time.Now(),
	}

	if err := s.orderHistoryRepo.UpdateOrderHistory(ctx, orderID, log); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *OrderService) GetOrderHistory(ctx context.Context, orderID string) (*models.OrderHistory, error) {
	return s.orderHistoryRepo.GetOrderHistory(ctx, orderID)
}

// GetPendingOrders retrieves all orders with pending status
func (s *OrderService) GetPendingOrders() ([]models.Order, error) {
	return s.orderRepo.FindByStatus(models.OrderStatusPending)
}

// ProcessOrder processes a single order
func (s *OrderService) ProcessOrder(ctx context.Context, orderID string) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Only update if status is still pending
	if err := s.orderRepo.UpdateStatus(orderID, models.OrderStatusPending, models.OrderStatusProcessing); err != nil {
		tx.Rollback()
		return err
	}

	log := models.OrderLog{
		Status:    models.OrderStatusProcessing,
		Timestamp: time.Now(),
	}

	if err := s.orderHistoryRepo.UpdateOrderHistory(ctx, orderID, log); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *OrderService) ResetProcessingOrders() error {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.orderRepo.ResetProcessingOrders(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// CompleteOrder marks an order as completed
func (s *OrderService) CompleteOrder(ctx context.Context, orderID string) error {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update order status to completed
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := s.orderRepo.UpdateStatus(orderID, order.Status, models.OrderStatusCompleted); err != nil {
		tx.Rollback()
		return err
	}

	// Update order history
	log := models.OrderLog{
		Status:    models.OrderStatusCompleted,
		Timestamp: time.Now(),
	}

	if err := s.orderHistoryRepo.UpdateOrderHistory(ctx, orderID, log); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
