package service

import (
	"database/sql"
	"errors"
	"orderservice/internal/logger"
	"orderservice/internal/models"
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	db *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) CreateOrder(req *models.CreateOrderRequest) (*models.Order, error) {
	logger.InfoLogger.Printf("Creating order with request: %+v", req)

	now := time.Now()
	// Calculate total price based on order count
	totalPrice := req.Price * float64(req.OrderCount)

	order := &models.Order{
		ID:            uuid.New(),
		InventoryID:   req.InventoryID,
		UserID:        req.UserID,
		PaymentStatus: models.PaymentStatusPending,
		OrderStatus:   models.OrderStatusPending,
		Price:         totalPrice,
		OrderCount:    req.OrderCount,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	query := `
        INSERT INTO orders (
            id, inventory_id, user_id, payment_status, order_status, 
            price, order_count, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING *`

	err := s.db.QueryRow(
		query,
		order.ID,
		order.InventoryID,
		order.UserID,
		order.PaymentStatus,
		order.OrderStatus,
		order.Price,
		order.OrderCount,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(
		&order.ID,
		&order.InventoryID,
		&order.UserID,
		&order.PaymentStatus,
		&order.OrderStatus,
		&order.Price,
		&order.OrderCount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		logger.ErrorLogger.Printf("Failed to create order: %v", err)
		return nil, err
	}

	// Queue the order for processing through the worker pool
	s.QueueOrderForProcessing(order)

	logger.InfoLogger.Printf("Successfully created order: %+v", order)
	return order, nil
}

func (s *OrderService) QueueOrderForProcessing(order *models.Order) {
	// Implementation in worker pool
	logger.InfoLogger.Printf("Order %s queued for processing", order.ID)
}

func (s *OrderService) CancelOrder(orderID uuid.UUID) error {
	logger.InfoLogger.Printf("Cancelling order: %s", orderID)

	query := `
        UPDATE orders 
        SET order_status = 'cancelled', updated_at = $1
        WHERE id = $2 AND order_status != 'completed'`

	result, err := s.db.Exec(query, time.Now(), orderID)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to cancel order %s: %v", orderID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get rows affected for order %s: %v", orderID, err)
		return err
	}

	if rowsAffected == 0 {
		logger.ErrorLogger.Printf("Order not found or already completed: %s", orderID)
		return errors.New("order not found or already completed")
	}

	logger.InfoLogger.Printf("Successfully cancelled order: %s", orderID)
	return nil
}

func (s *OrderService) GetOrder(orderID uuid.UUID) (*models.Order, error) {
	logger.InfoLogger.Printf("Getting order: %s", orderID)

	order := &models.Order{}
	query := `SELECT * FROM orders WHERE id = $1`

	err := s.db.QueryRow(query, orderID).Scan(
		&order.ID,
		&order.InventoryID,
		&order.UserID,
		&order.PaymentStatus,
		&order.OrderStatus,
		&order.Price,
		&order.OrderCount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		logger.ErrorLogger.Printf("Order not found: %s", orderID)
		return nil, errors.New("order not found")
	}

	if err != nil {
		logger.ErrorLogger.Printf("Failed to get order %s: %v", orderID, err)
		return nil, err
	}

	logger.InfoLogger.Printf("Successfully retrieved order: %+v", order)
	return order, nil
}

func (s *OrderService) GetPendingOrders() ([]*models.Order, error) {
	logger.InfoLogger.Printf("Getting pending orders")

	query := `
        SELECT id, inventory_id, user_id, payment_status, order_status, price, order_count, created_at, updated_at 
        FROM orders 
        WHERE order_status = 'pending'
        LIMIT 100`

	rows, err := s.db.Query(query)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get pending orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(
			&order.ID,
			&order.InventoryID,
			&order.UserID,
			&order.PaymentStatus,
			&order.OrderStatus,
			&order.Price,
			&order.OrderCount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to scan pending order: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	logger.InfoLogger.Printf("Successfully retrieved pending orders")
	return orders, nil
}

func (s *OrderService) ProcessPendingOrder(orderID uuid.UUID) error {
	logger.InfoLogger.Printf("Processing order: %s", orderID)

	// Simulate order processing
	time.Sleep(2 * time.Second)

	query := `
        UPDATE orders 
        SET order_status = 'completed', updated_at = $1
        WHERE id = $2 AND order_status = 'pending'
        RETURNING id`

	result, err := s.db.Exec(query, time.Now(), orderID)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to process order %s: %v", orderID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get rows affected for order %s: %v", orderID, err)
		return err
	}

	if rowsAffected == 0 {
		logger.ErrorLogger.Printf("Order not found or not in pending status: %s", orderID)
		return errors.New("order not found or not in pending status")
	}

	logger.InfoLogger.Printf("Successfully processed order: %s", orderID)
	return nil
}

// Add this method to OrderService
func (s *OrderService) UpdateOrderStatus(orderID uuid.UUID, status models.OrderStatus) error {
    logger.InfoLogger.Printf("Updating order %s status to %s", orderID, status)
    
    query := `
        UPDATE orders 
        SET order_status = $1, updated_at = $2
        WHERE id = $3
        RETURNING id`
        
    result, err := s.db.Exec(query, status, time.Now(), orderID)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to update order status: %v", err)
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("order not found")
    }

    return nil
}