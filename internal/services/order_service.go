package services

import (
	"context"
	"errors"
	"log"
	"mypackage/internal/models"
	"mypackage/internal/repositories"
	"time"

	"gorm.io/gorm"
)

// OrderServiceInterface defines the contract for order service
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	GetUserOrders(userID string) ([]models.Order, error)
	CancelOrder(ctx context.Context, id string) error
	ProcessPendingOrders(ctx context.Context) error
	GetOrderHistory(ctx context.Context, id string) (*models.OrderHistory, error)      // Add this method
	UpdateOrderStatus(ctx context.Context, id string, status models.OrderStatus) error // Add this method
}

type OrderService struct {
	db               *gorm.DB
	orderRepo        repositories.OrderRepositoryInterface
	inventoryRepo    repositories.InventoryRepositoryInterface
	orderHistoryRepo *repositories.OrderHistoryRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo repositories.OrderRepositoryInterface,
	inventoryRepo repositories.InventoryRepositoryInterface,
	orderHistoryRepo *repositories.OrderHistoryRepository,
) *OrderService {
	return &OrderService{
		db:               db,
		orderRepo:        orderRepo,
		inventoryRepo:    inventoryRepo,
		orderHistoryRepo: orderHistoryRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	// Validate order status
	order.Status = models.OrderStatusPending
	if !order.Status.IsValid() {
		return errors.New("invalid order status")
	}

	// Start PostgreSQL transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check inventory availability
	var inventory models.Inventory
	if err := tx.Where("id = ?", order.InventoryID).First(&inventory).Error; err != nil {
		tx.Rollback()
		return err
	}

	if inventory.Stock < order.Quantity {
		tx.Rollback()
		return errors.New("insufficient stock")
	}

	// Deduct stock
	if err := tx.Model(&inventory).Update("stock", gorm.Expr("stock - ?", order.Quantity)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create order
	order.Status = models.OrderStatusPending
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create order history in MongoDB
	orderHistory := &models.OrderHistory{
		OrderID: order.ID.String(),
		Status:  models.OrderStatusPending,
		Log: []models.OrderLog{{
			Status:    models.OrderStatusPending,
			Timestamp: time.Now(),
		}},
	}

	if err := s.orderHistoryRepo.CreateOrderHistory(ctx, orderHistory); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *OrderService) GetOrder(id string) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) GetUserOrders(userID string) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
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

	// Update order status
	if err := tx.Model(&models.Order{}).Where("id = ?", id).
		Update("status", models.OrderStatusCanceled).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update order history
	log := models.OrderLog{
		Status:    models.OrderStatusCanceled,
		Timestamp: time.Now(),
	}
	if err := s.orderHistoryRepo.UpdateOrderHistory(ctx, id, log); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Add ProcessPendingOrders method to implement the interface
func (s *OrderService) ProcessPendingOrders(ctx context.Context) error {
	// Get all pending orders
	pendingOrders, err := s.orderRepo.FindByStatus(models.OrderStatusPending)
	if err != nil {
		return err
	}

	for _, order := range pendingOrders {
		// Start transaction for each order
		tx := s.db.Begin()
		if tx.Error != nil {
			return tx.Error
		}

		// Wrap in closure to use defer
		func() {
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			// Update order status to processing
			if err := tx.Model(&order).Update("status", models.OrderStatusProcessing).Error; err != nil {
				tx.Rollback()
				return
			}

			// Log status change to MongoDB
			processingLog := models.OrderLog{
				Status:    models.OrderStatusProcessing,
				Timestamp: time.Now(),
			}
			if err := s.orderHistoryRepo.UpdateOrderHistory(ctx, order.ID.String(), processingLog); err != nil {
				tx.Rollback()
				return
			}

			// Simulate processing time (optional, remove in production)
			time.Sleep(2 * time.Second)

			// Update to completed status
			if err := tx.Model(&order).Update("status", models.OrderStatusCompleted).Error; err != nil {
				tx.Rollback()
				return
			}

			// Log completion to MongoDB
			completedLog := models.OrderLog{
				Status:    models.OrderStatusCompleted,
				Timestamp: time.Now(),
			}
			if err := s.orderHistoryRepo.UpdateOrderHistory(ctx, order.ID.String(), completedLog); err != nil {
				tx.Rollback()
				return
			}

			if err := tx.Commit().Error; err != nil {
				log.Printf("Error committing transaction for order %s: %v", order.ID, err)
			} else {
				log.Printf("Successfully processed order %s", order.ID)
			}
		}()
	}

	return nil
}

// Add GetOrderHistory implementation to OrderService
func (s *OrderService) GetOrderHistory(ctx context.Context, id string) (*models.OrderHistory, error) {
	return s.orderHistoryRepo.GetOrderHistory(ctx, id)
}

// Add UpdateOrderStatus implementation to OrderService
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status models.OrderStatus) error {
	if !status.IsValid() {
		return errors.New("invalid order status")
	}

	return s.orderRepo.UpdateStatus(id, status)
}

type OrderRepositoryInterface interface {
	UpdateStatus(id string, status models.OrderStatus) error
}

type OrderRepository struct {
	db *gorm.DB
}

func (r *OrderRepository) UpdateStatus(id string, status models.OrderStatus) error {
	if !status.IsValid() {
		return errors.New("invalid order status")
	}

	return r.db.Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}
