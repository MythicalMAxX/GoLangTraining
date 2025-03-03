package service

import (
	"delivery-service/internal/models"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	initialWorkers     = 10
	maxWorkers         = 50
	minWorkers         = 5
	scaleUpThreshold   = 100 // Queue size that triggers scaling up
	scaleDownThreshold = 20  // Queue size that triggers scaling down
)

type WorkerPool struct {
	jobs            chan *models.Order
	deliveryService *DeliveryService
	workersCount    int
	mutex           sync.Mutex
}

func NewWorkerPool(deliveryService *DeliveryService) *WorkerPool {
	wp := &WorkerPool{
		jobs:            make(chan *models.Order, 1000),
		deliveryService: deliveryService,
		workersCount:    initialWorkers,
	}
	wp.Start()
	go wp.monitorQueueSize()
	return wp
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workersCount; i++ {
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d started", id)
	for order := range wp.jobs {
		wp.processOrder(id, order)
	}
}

func (wp *WorkerPool) processOrder(workerID int, order *models.Order) {
	if order == nil {
		log.Printf("Worker %d received nil order", workerID)
		return
	}

	log.Printf("Worker %d starting to process order: ID=%s, RestaurantID=%s, Status=%s",
		workerID, order.ID, order.RestaurantID, order.Status)

	log.Printf("Worker %d processing order %s", workerID, order.ID)

	// Step 1: Process Order (1 minute)
	log.Printf("Worker %d: Processing order %s...", workerID, order.ID)
	time.Sleep(1 * time.Minute)

	// Send order processed notification
	if err := wp.deliveryService.sendOrderStatusNotification(order.ID, models.OrderProcessed, "Order is being processed by restaurant"); err != nil {
		log.Printf("Error sending order processed notification: %v", err)
	}

	// Step 2: Assign and notify delivery partner (1 minute)
	log.Printf("Worker %d: Assigning delivery partner for order %s...", workerID, order.ID)
	time.Sleep(1 * time.Minute)

	deliveryPartner, err := wp.deliveryService.assignDeliveryPartner()
	if err != nil {
		log.Printf("Error assigning delivery partner: %v", err)
		return
	}

	// Create delivery record with all required fields
	delivery := &models.OrderDelivery{
		OrderID:           order.ID,           // Ensure order ID is set
		RestaurantID:      order.RestaurantID, // Set restaurant ID from order
		DeliveryPartnerID: deliveryPartner.ID, // Ensure delivery partner ID is set
		Status:            "ASSIGNED",
		PickupTime:        time.Now(),
	}

	// Log delivery details before creation
	log.Printf("Creating delivery record: OrderID=%s, DeliveryPartnerID=%s",
		delivery.OrderID, delivery.DeliveryPartnerID)

	if err := wp.deliveryService.createDelivery(delivery); err != nil {
		log.Printf("Error creating delivery record: %v", err)
		return
	}

	// Send delivery partner assigned notification
	if err := wp.deliveryService.sendOrderStatusNotification(order.ID, models.OrderPickedByDelivery,
		fmt.Sprintf("Order picked up by delivery partner %s", deliveryPartner.Name),
		deliveryPartner.ID); err != nil {
		log.Printf("Error sending delivery partner assigned notification: %v", err)
	}

	// Step 3: Simulate delivery (1 minute)
	log.Printf("Worker %d: Delivering order %s...", workerID, order.ID)
	time.Sleep(1 * time.Minute)

	if err := wp.deliveryService.completeDelivery(delivery); err != nil {
		log.Printf("Error completing delivery: %v", err)
		return
	}

	// Send order delivered notification
	if err := wp.deliveryService.sendOrderStatusNotification(order.ID, models.OrderDelivered,
		"Order has been delivered successfully", deliveryPartner.ID); err != nil {
		log.Printf("Error sending order delivered notification: %v", err)
	}

	log.Printf("Worker %d completed processing order %s", workerID, order.ID)
}

func (wp *WorkerPool) AddJob(order *models.Order) {
	wp.jobs <- order
}

func (wp *WorkerPool) monitorQueueSize() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		queueSize := len(wp.jobs)
		wp.mutex.Lock()

		// Scale up if queue is getting full
		if queueSize > scaleUpThreshold && wp.workersCount < maxWorkers {
			newWorkers := 5
			for i := 0; i < newWorkers; i++ {
				go wp.worker(wp.workersCount + i)
			}
			wp.workersCount += newWorkers
			log.Printf("Scaled up to %d workers", wp.workersCount)
		}

		// Scale down if queue is small
		if queueSize < scaleDownThreshold && wp.workersCount > minWorkers {
			wp.workersCount -= 2
			log.Printf("Scaled down to %d workers", wp.workersCount)
		}

		wp.mutex.Unlock()
	}
}
