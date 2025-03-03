package worker

import (
	"context"
	"log"
	"orderservice/internal/logger"
	"orderservice/internal/models"
	"orderservice/internal/service"
	"sync"
	"time"
)

type WorkerPool struct {
    config      *Config
    jobs        chan *models.Order
    orderSvc    *service.OrderService
    workerCount int
    mu          sync.RWMutex
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup  // Add WaitGroup to track workers
}

func NewWorkerPool(config *Config, orderSvc *service.OrderService) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        config:      config,
        jobs:        make(chan *models.Order, config.QueueSize),
        orderSvc:    orderSvc,
        workerCount: config.InitialWorkers,
        ctx:         ctx,
        cancel:      cancel,
        wg:          sync.WaitGroup{},
    }
}

func (p *WorkerPool) Start() {
	// Start initial workers
	for i := 0; i < p.config.InitialWorkers; i++ {
		p.startWorker()
	}

	// Start periodic processing
	go p.periodicProcessing()

	// Start auto-scaling monitor
	go p.monitorQueueAndScale()
}

func (p *WorkerPool) Stop() {
    logger.InfoLogger.Printf("Stopping worker pool...")
    p.cancel()
    close(p.jobs)
    p.wg.Wait()
    logger.InfoLogger.Printf("Worker pool stopped")
}

func (p *WorkerPool) startWorker() {
    p.mu.Lock()
    p.workerCount++
    p.mu.Unlock()

    p.wg.Add(1)
    go func() {
        defer func() {
            p.mu.Lock()
            p.workerCount--
            p.mu.Unlock()
            p.wg.Done()
        }()

        for {
            select {
            case job, ok := <-p.jobs:
                if !ok {
                    return
                }
                p.processOrder(job)
            case <-p.ctx.Done():
                return
            }
        }
    }()
}

func (p *WorkerPool) processOrder(order *models.Order) {
    logger.InfoLogger.Printf("Worker processing order: %s", order.ID)

    // Update to processing status
    err := p.orderSvc.UpdateOrderStatus(order.ID, models.OrderStatusProcessing)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to update order status to processing: %v", err)
        return
    }

    // Simulate processing time (reduced for testing)
    time.Sleep(500 * time.Millisecond)

    // Update to completed status
    err = p.orderSvc.UpdateOrderStatus(order.ID, models.OrderStatusCompleted)
    if err != nil {
        logger.ErrorLogger.Printf("Failed to complete order processing: %v", err)
        // Retry logic could be added here
        return
    }

    logger.InfoLogger.Printf("Successfully processed order: %s", order.ID)
}


func (p *WorkerPool) periodicProcessing() {
	ticker := time.NewTicker(p.config.ProcessInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.fetchAndQueuePendingOrders()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *WorkerPool) fetchAndQueuePendingOrders() {
    orders, err := p.orderSvc.GetPendingOrders()
    if err != nil {
        logger.ErrorLogger.Printf("Error fetching pending orders: %v", err)
        return
    }

    logger.InfoLogger.Printf("Found %d pending orders", len(orders))

    for _, order := range orders {
        select {
        case p.jobs <- order:
            logger.InfoLogger.Printf("Queued order %s for processing", order.ID)
        case <-p.ctx.Done():
            return
        default:
            logger.InfoLogger.Printf("Queue is full, order %s will be processed in next batch", order.ID)
        }
    }
}

func (p *WorkerPool) monitorQueueAndScale() {
	ticker := time.NewTicker(p.config.ScaleInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.adjustWorkerCount()
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *WorkerPool) adjustWorkerCount() {
	p.mu.RLock()
	currentWorkers := p.workerCount
	p.mu.RUnlock()

	queueUtilization := float64(len(p.jobs)) / float64(p.config.QueueSize)
	logger.InfoLogger.Printf("Current worker count: %d, Queue utilization: %.2f%%",
		currentWorkers, queueUtilization*100)

	switch {
	case queueUtilization > p.config.ScaleThreshold && currentWorkers < p.config.MaxWorkers:
		// Scale up
		workersToAdd := min(p.config.MaxWorkers-currentWorkers, 2)
		for i := 0; i < workersToAdd; i++ {
			p.startWorker()
		}
		log.Printf("Scaled up workers to %d", currentWorkers+workersToAdd)

	case queueUtilization < p.config.ScaleThreshold/2 && currentWorkers > p.config.MinWorkers:
		// Scale down
		// Workers will naturally scale down when they finish their current jobs
		log.Printf("Scaling down workers. Current: %d", currentWorkers)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
