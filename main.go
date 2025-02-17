// Handled race issues
// Implemented worker pool for handling multiple tasks
// added a job queue using buffered channel to store jobs
// workers takes jobs from buffered queue and execute them.

package main

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Job struct
type Job struct {
	id   int
	data string
}

// Worker struct
type Worker struct {
	id       int
	jobQueue chan Job
	logger   *zap.Logger
}

// WorkerPool struct
type WorkerPool struct {
	workerCount int
	jobQueue    chan Job
	workers     []*Worker
	wg          sync.WaitGroup
	logger      *zap.Logger
}

// NewWorkerPool function
func NewWorkerPool(workerCount int, logger *zap.Logger) *WorkerPool {
	jobQueue := make(chan Job, workerCount)
	workers := make([]*Worker, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = &Worker{
			id:       i,
			jobQueue: jobQueue,
			logger:   logger,
		}
	}
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    jobQueue,
		workers:     workers,
		logger:      logger,
	}
}

// Start function
func (wp *WorkerPool) Start() {
	wp.logger.Info("Starting worker pool")
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.workers[i].start(&wp.wg)
	}
}

// Stop function
func (wp *WorkerPool) Stop() {
	wp.logger.Info("Stopping worker pool")
	close(wp.jobQueue)
	wp.wg.Wait()
}

// AddJob function
func (wp *WorkerPool) AddJob(job Job) {
	wp.logger.Info("Adding job to queue", zap.Int("jobID", job.id))
	wp.jobQueue <- job
}

// start function
func (w *Worker) start(wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range w.jobQueue {
		w.logger.Info("Worker started job", zap.Int("workerID", w.id), zap.Int("jobID", job.id))
		// Simulate job processing
		fmt.Printf("Worker with id %d started job with id %d\n", w.id, job.id)
		time.Sleep(5 * time.Second)
		fmt.Printf("Worker with id %d finished job with id %d\n", w.id, job.id)
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	wp := NewWorkerPool(3, logger)
	wp.Start()
	// wp.AddJob(Job{id: 1, data: "data1"})
	// wp.AddJob(Job{id: 2, data: "data2"})
	// wp.AddJob(Job{id: 3, data: "data3"})
	// wp.AddJob(Job{id: 4, data: "data4"})
	// wp.AddJob(Job{id: 5, data: "data5"})

	for i := 1; i <= 10; i++ {
		wp.AddJob(Job{id: i, data: "Task" + string(i)})
	}
	wp.Stop()
}
