package service

import (
	"delivery-service/config"
	"delivery-service/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type DeliveryService struct {
	db         *gorm.DB
	config     *config.Config
	conn       *amqp.Connection
	ch         *amqp.Channel
	workerPool *WorkerPool
}

func NewDeliveryService(db *gorm.DB, cfg *config.Config) *DeliveryService {
	ds := &DeliveryService{
		db:     db,
		config: cfg,
	}
	ds.workerPool = NewWorkerPool(ds)
	return ds
}

func (s *DeliveryService) Start() error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(s.config.RabbitMQURL)
	if err != nil {
		return err
	}
	s.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	s.ch = ch

	// Declare exchange
	err = ch.ExchangeDeclare(
		"orders_exchange", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		"delivery_orders_queue", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,            // queue name
		"order.#",         // routing key
		"orders_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare notification queue
	notificationQueue, err := s.ch.QueueDeclare(
		"delivery_notifications", // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return fmt.Errorf("error declaring notification queue: %v", err)
	}

	// Bind notification queue
	err = s.ch.QueueBind(
		notificationQueue.Name, // queue name
		"notification.#",       // routing key
		"orders_exchange",      // exchange
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return fmt.Errorf("error binding notification queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			s.handleMessage(d)
		}
	}()

	log.Printf(" [*] Delivery service waiting for messages")
	<-forever

	return nil
}

func (s *DeliveryService) handleMessage(d amqp.Delivery) {
	// Log the received message
	log.Printf("Received message: %s", string(d.Body))

	// Validate message is not empty
	if len(d.Body) == 0 {
		log.Printf("Empty message received")
		d.Nack(false, false)
		return
	}

	var message models.OrderMessage
	if err := json.Unmarshal(d.Body, &message); err != nil {
		log.Printf("Error unmarshaling message wrapper: %v\nMessage body: %s", err, string(d.Body))
		d.Nack(false, false)
		return
	}

	// Check if this is a new order
	if message.Type != "order.new" {
		log.Printf("Ignoring message of type: %s", message.Type)
		d.Ack(false)
		return
	}

	// Enhanced validation
	if err := s.validateOrder(&message.Order); err != nil {
		log.Printf("Order validation failed: %v", err)
		d.Nack(false, false)
		return
	}

	log.Printf("Successfully parsed order: ID=%s, CustomerID=%s, RestaurantID=%s, Status=%s",
		message.Order.ID, message.Order.CustomerID, message.Order.RestaurantID, message.Order.Status)

	// Add the order to the worker pool
	s.workerPool.AddJob(&message.Order)

	// Acknowledge the message
	d.Ack(false)
}

func (s *DeliveryService) validateOrder(order *models.Order) error {
	if order.ID == "" {
		return fmt.Errorf("order ID is empty")
	}
	if order.RestaurantID == "" {
		return fmt.Errorf("restaurant ID is empty")
	}
	if order.CustomerID == "" {
		return fmt.Errorf("customer ID is empty")
	}
	if order.Status == "" {
		return fmt.Errorf("order status is empty")
	}

	// Don't verify in database as this is a new order
	return nil
}

func (s *DeliveryService) updateOrderDeliveryPartner(orderID string, partnerID string) error {
	orderUpdate := models.OrderUpdate{
		OrderID: orderID,
		Status:  "DELIVERY_ASSIGNED",
		Message: "Delivery partner assigned",
	}

	return s.publishOrderUpdate(orderUpdate)
}

func (s *DeliveryService) updateOrderStatus(orderID string, status string) error {
	orderUpdate := models.OrderUpdate{
		OrderID: orderID,
		Status:  status,
		Message: "Order status updated to: " + status,
	}

	return s.publishOrderUpdate(orderUpdate)
}

func (s *DeliveryService) publishOrderUpdate(update models.OrderUpdate) error {
	updateBytes, err := json.Marshal(update)
	if err != nil {
		return err
	}

	return s.ch.Publish(
		"orders_exchange",
		"order.update",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        updateBytes,
		})
}

func (s *DeliveryService) assignDeliveryPartner() (*models.DeliveryPartner, error) {
	var partner models.DeliveryPartner
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Find available delivery partner
		if err := tx.Where("status = ?", models.DeliveryPartnerStatusAvailable).
			First(&partner).Error; err != nil {
			return fmt.Errorf("no available delivery partners: %v", err)
		}

		// Update partner status to busy
		if err := tx.Model(&partner).
			Where("id = ?", partner.ID).
			Update("status", models.DeliveryPartnerStatusBusy).Error; err != nil {
			return fmt.Errorf("failed to update partner status: %v", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Assigned delivery partner: ID=%s, Name=%s", partner.ID, partner.Name)
	return &partner, nil
}

func (s *DeliveryService) sendOrderStatusNotification(orderID string, notificationType models.NotificationType, message string, deliveryPartnerID ...string) error {
	notification := &models.Notification{
		Type:    notificationType,
		OrderID: orderID,
		Message: message,
	}

	if len(deliveryPartnerID) > 0 {
		notification.DeliveryID = deliveryPartnerID[0]
	}

	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error marshaling notification: %v", err)
	}

	// Publish to notification queue
	err = s.ch.Publish(
		"orders_exchange",    // exchange
		"notification.order", // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         notificationBytes,
			DeliveryMode: amqp.Persistent, // Make messages persistent
			Type:         string(notificationType),
		})

	if err != nil {
		return fmt.Errorf("error publishing notification: %v", err)
	}

	log.Printf("Published notification: Type=%s, OrderID=%s, Message=%s",
		notificationType, orderID, message)
	return nil
}

func (s *DeliveryService) createDelivery(delivery *models.OrderDelivery) error {
	// Generate UUID for delivery
	delivery.ID = uuid.New().String()

	// Start transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Verify order exists and get full order data
		var order models.Order
		if err := tx.Where("id = ?", delivery.OrderID).First(&order).Error; err != nil {
			return fmt.Errorf("order not found (ID: %s): %v", delivery.OrderID, err)
		}

		// Verify delivery partner exists
		var partner models.DeliveryPartner
		if err := tx.Where("id = ?", delivery.DeliveryPartnerID).First(&partner).Error; err != nil {
			return fmt.Errorf("delivery partner not found (ID: %s): %v",
				delivery.DeliveryPartnerID, err)
		}

		// Set restaurant ID from order if not set
		if delivery.RestaurantID == "" {
			delivery.RestaurantID = order.RestaurantID
		}

		// Create delivery record with validated data
		if err := tx.Create(delivery).Error; err != nil {
			return fmt.Errorf("failed to create delivery record: %v", err)
		}

		// Update order status
		if err := tx.Model(&order).Update("status", "OUT_FOR_DELIVERY").Error; err != nil {
			return fmt.Errorf("failed to update order status: %v", err)
		}

		log.Printf("Created delivery record: ID=%s, OrderID=%s, RestaurantID=%s, DeliveryPartnerID=%s",
			delivery.ID, delivery.OrderID, delivery.RestaurantID, delivery.DeliveryPartnerID)
		return nil
	})
}

func (s *DeliveryService) updateDelivery(delivery *models.OrderDelivery) error {
	return s.db.Save(delivery).Error
}

func (s *DeliveryService) completeDelivery(delivery *models.OrderDelivery) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update delivery record
		if err := tx.Model(delivery).Updates(map[string]interface{}{
			"status":        "DELIVERED",
			"delivery_time": time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("failed to update delivery record: %v", err)
		}

		// Update order status
		if err := tx.Model(&models.Order{}).
			Where("id = ?", delivery.OrderID).
			Update("status", "DELIVERED").Error; err != nil {
			return fmt.Errorf("failed to update order status: %v", err)
		}

		// Mark delivery partner as available again
		if err := tx.Model(&models.DeliveryPartner{}).
			Where("id = ?", delivery.DeliveryPartnerID).
			Update("status", models.DeliveryPartnerStatusAvailable).Error; err != nil {
			return fmt.Errorf("failed to update delivery partner status: %v", err)
		}

		log.Printf("Completed delivery: ID=%s, OrderID=%s",
			delivery.ID, delivery.OrderID)
		return nil
	})
}
