package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"orderservice/internal/services"
)

type OrderWorker struct {
	orderService *services.OrderService
	isRunning    bool
	interval     time.Duration
	stopChan     chan struct{}
	wg           sync.WaitGroup
	workerCount  int
	workQueue    chan struct{}
	metrics      *WorkerMetrics
	schedulerWg  sync.WaitGroup
}

func (m *WorkerMetrics) incrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCount++
}

func (m *WorkerMetrics) incrementProcessed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.processedOrders++
}

type WorkerMetrics struct {
	mu              sync.RWMutex
	processedOrders int64
	errorCount      int64
	lastScaleTime   time.Time
}

func NewOrderWorker(orderService *services.OrderService, initialWorkers int) *OrderWorker {
	return &OrderWorker{
		orderService: orderService,
		interval:     time.Minute,
		stopChan:     make(chan struct{}),
		workerCount:  initialWorkers,
		workQueue:    make(chan struct{}, 100),
		metrics:      &WorkerMetrics{lastScaleTime: time.Now()},
	}
}

func (w *OrderWorker) Start() {
	if w.isRunning {
		return
	}
	w.isRunning = true
	w.stopChan = make(chan struct{})
	w.workQueue = make(chan struct{}, 100)

	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.startSingleWorker(i)
	}

	w.schedulerWg.Add(1)
	go func() {
		defer w.schedulerWg.Done()
		w.scheduleWork()
	}()
}

func (w *OrderWorker) Stop() {
	if !w.isRunning {
		return
	}
	w.isRunning = false
	close(w.stopChan)
	close(w.workQueue)
	w.wg.Wait()
	w.schedulerWg.Wait()
}

func (w *OrderWorker) scheduleWork() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.workQueue <- struct{}{}
		case <-w.stopChan:
			return
		}
	}
}

func (w *OrderWorker) startSingleWorker(workerID int) {
	defer w.wg.Done()
	log.Printf("Order processing worker %d started", workerID)

	for {
		select {
		case _, ok := <-w.workQueue:
			if !ok {
				return
			}
			ctx := context.Background()
			if err := w.processOrders(ctx); err != nil {
				w.metrics.incrementErrors()
				log.Printf("Worker %d error: %v", workerID, err)
			} else {
				w.metrics.incrementProcessed()
			}
		case <-w.stopChan:
			return
		}
	}
}

func (w *OrderWorker) processOrders(ctx context.Context) error {
	// Multiple workers can get the same pending orders simultaneously
	orders, err := w.orderService.GetPendingOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		// Each worker tries to process the same order
		if err := w.orderService.ProcessOrder(ctx, order.ID.String()); err != nil {
			return err
		}

		// Simulate some processing work
		time.Sleep(2 * time.Second)

		// Then mark as completed
		if err := w.orderService.CompleteOrder(ctx, order.ID.String()); err != nil {
			return err
		}

		log.Printf("Order %s completed successfully", order.ID.String())
	}
	return nil
}
