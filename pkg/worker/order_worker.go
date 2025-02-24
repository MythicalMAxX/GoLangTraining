package worker

import (
	"context"
	"log"
	"mypackage/internal/services"
	"sync"
	"time"
)

type OrderWorker struct {
	orderService services.OrderServiceInterface
	isRunning    bool
	interval     time.Duration
	stopChan     chan struct{}
	wg           sync.WaitGroup
	workerCount  int
	workQueue    chan struct{}
	metrics      *WorkerMetrics
	schedulerWg  sync.WaitGroup
}

type WorkerMetrics struct {
	mu              sync.RWMutex
	processedOrders int64
	errorCount      int64
	lastScaleTime   time.Time
}

// Add these methods to handle metrics
func (m *WorkerMetrics) incrementProcessed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.processedOrders++
}

func (m *WorkerMetrics) incrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCount++
}

func (m *WorkerMetrics) getMetrics() (processed int64, errors int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.processedOrders, m.errorCount
}

func NewOrderWorker(orderService services.OrderServiceInterface, initialWorkers int) *OrderWorker {
	return &OrderWorker{
		orderService: orderService,
		interval:     time.Minute,
		stopChan:     make(chan struct{}),
		workerCount:  initialWorkers,
		workQueue:    make(chan struct{}, 100), // Buffer size for work items
		metrics:      &WorkerMetrics{lastScaleTime: time.Now()},
	}
}

func (w *OrderWorker) startSingleWorker(workerID int) {
	defer w.wg.Done()
	log.Printf("Order processing worker %d started", workerID)

	for {
		select {
		case _, ok := <-w.workQueue:
			if !ok {
				log.Printf("Worker %d stopping due to closed queue", workerID)
				return
			}
			ctx := context.Background()
			if err := w.processOrders(ctx); err != nil {
				w.metrics.incrementErrors()
				log.Printf("Worker %d error processing orders: %v", workerID, err)
			} else {
				w.metrics.incrementProcessed()
			}
		case <-w.stopChan:
			log.Printf("Order processing worker %d stopped", workerID)
			return
		}
	}
}

func (w *OrderWorker) scaleWorkers() {
	processedOrders, errorCount := w.metrics.getMetrics()

	w.metrics.mu.RLock()
	timeSinceLastScale := time.Since(w.metrics.lastScaleTime)
	w.metrics.mu.RUnlock()

	// Only scale every 5 minutes
	if timeSinceLastScale < 5*time.Minute {
		return
	}

	// Scale up if we have high throughput and low errors
	if processedOrders > 100 && errorCount < 10 && w.workerCount < 10 {
		w.addWorker()
	}

	// Scale down if we have low throughput or high errors
	if (processedOrders < 50 || errorCount > 20) && w.workerCount > 1 {
		w.removeWorker()
	}

	w.metrics.mu.Lock()
	w.metrics.lastScaleTime = time.Now()
	w.metrics.processedOrders = 0
	w.metrics.errorCount = 0
	w.metrics.mu.Unlock()
}

func (w *OrderWorker) addWorker() {
	w.workerCount++
	w.wg.Add(1)
	go w.startSingleWorker(w.workerCount)
	log.Printf("Scaled up to %d workers", w.workerCount)
}

func (w *OrderWorker) removeWorker() {
	w.workerCount--
	log.Printf("Scaled down to %d workers", w.workerCount)
}

func (w *OrderWorker) Start() {
	if w.isRunning {
		return
	}
	w.isRunning = true
	w.stopChan = make(chan struct{})       // Reset stop channel
	w.workQueue = make(chan struct{}, 100) // Reset work queue

	// Start multiple workers
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.startSingleWorker(i)
	}

	// Start the scheduler with wait group
	w.schedulerWg.Add(1)
	go func() {
		defer w.schedulerWg.Done()
		w.scheduleWork()
	}()
}

// Add this new method
func (w *OrderWorker) scheduleWork() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if w.isRunning {
				select {
				case w.workQueue <- struct{}{}:
					log.Printf("Scheduled work item for processing")
				default:
					log.Printf("Work queue full, skipping schedule")
				}
			}
		case <-w.stopChan:
			log.Println("Scheduler stopping...")
			return
		}
	}
}

func (w *OrderWorker) Stop() {
	if !w.isRunning {
		return
	}
	log.Println("Initiating worker shutdown...")
	w.isRunning = false

	// Signal stop to all workers and scheduler
	close(w.stopChan)

	// Close work queue after all workers have stopped
	log.Println("Closing work queue...")
	close(w.workQueue)

	// Wait for all workers to finish
	log.Println("Waiting for workers to finish...")
	w.wg.Wait()

	// Wait for scheduler to finish
	log.Println("Waiting for scheduler to finish...")
	w.schedulerWg.Wait()

	log.Println("Worker shutdown complete")
}

func (w *OrderWorker) processOrders(ctx context.Context) error {
	log.Println("Processing pending orders...")
	w.scaleWorkers()	// Scales Workers Periodically
	err := w.orderService.ProcessPendingOrders(ctx)
	if err != nil {
		log.Printf("Error processing orders: %v", err)
		return err
	}
	log.Println("Finished processing pending orders")
	return nil
}

func (w *OrderWorker) IsRunning() bool {
	return w.isRunning
}

func (w *OrderWorker) SetInterval(interval time.Duration) {
	w.interval = interval
}

func (w *OrderWorker) GetWorkerCount() int {
	return w.workerCount
}
