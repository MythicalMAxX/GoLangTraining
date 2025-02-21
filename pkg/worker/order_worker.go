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
}

func NewOrderWorker(orderService services.OrderServiceInterface, workerCount int) *OrderWorker {
	if workerCount <= 0 {
		workerCount = 1
	}
	return &OrderWorker{
		orderService: orderService,
		interval:     time.Minute,
		stopChan:     make(chan struct{}),
		workerCount:  workerCount,
	}
}

func (w *OrderWorker) startSingleWorker(workerID int) {
	defer w.wg.Done()
	log.Printf("Order processing worker %d started", workerID)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx := context.Background()
			if err := w.processOrders(ctx); err != nil {
				log.Printf("Worker %d error processing orders: %v", workerID, err)
			}
		case <-w.stopChan:
			log.Printf("Order processing worker %d stopped", workerID)
			return
		}
	}
}

func (w *OrderWorker) Start() {
	if w.isRunning {
		return
	}
	w.isRunning = true

	// Start multiple workers
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.startSingleWorker(i)
	}
}

func (w *OrderWorker) Stop() {
	if !w.isRunning {
		return
	}
	w.isRunning = false
	close(w.stopChan)
	w.wg.Wait()
}

func (w *OrderWorker) processOrders(ctx context.Context) error {
	log.Println("Processing pending orders...")
	return w.orderService.ProcessPendingOrders(ctx)
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
