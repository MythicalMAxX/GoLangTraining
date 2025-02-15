package main

// importing required packages
import (
	"bufio"
	"container/heap"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
)

// FileJob represents a file processing task with its metadata
type FileJob struct {
	filepath string // Path to the file to be processed
	size     int64  // size of the file in bytes
	index    int    // Index used by the priority Queue
}

// FilePriorityQueue implements head.Interface and holds FileJobs
type FilePriorityQueue []*FileJob

// Standard heap interface implementations for FilePriorityQueue
func (pq FilePriorityQueue) Len() int { return len(pq) }

// Less determines the priority of the files - larger files have higher priority
func (pq FilePriorityQueue) Less(i, j int) bool {
	return pq[i].size > pq[j].size
}

func (pq FilePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// Push adds a new FileJob to the PriorityQueue
func (pq *FilePriorityQueue) Push(x interface{}) {
	n := len(*pq)
	job := x.(*FileJob)
	job.index = n
	*pq = append(*pq, job)
}

// Pop removes and returns the highest priority FileJob
func (pq *FilePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	job := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return job
}

// WorkerPool manager a pool of goroutines for concurrent file processing
type WorkerPool struct {
	jobs        chan *FileJob      // channel for distributing jobs to workers
	results     chan string        // channel for collecting processing results
	logger      *zap.Logger        // structured logger
	wg          sync.WaitGroup     //WaitGroup for worker synchronization
	mutex       sync.Mutex         // Mutex for thread-safe operations
	ctx         context.Context    //context for cancellation
	cancelFunc  context.CancelFunc // function to cancel the context
	workerCount int                //current number of active workers
	maxWorkers  int                // maximum allowed workers
}

// StartWorkers initializes and starts the specified number of worker goroutines
func (wp *WorkerPool) StartWorkers(workerCount int) {
	wp.logger.Info("Starting workers", zap.Int("WorkerCount", workerCount))
	wp.workerCount = workerCount
	for i := 0; i < workerCount; i++ {
		go wp.worker(i)
	}
}

// NewWorkerPool creates and initializes a new WorkerPool instance
func NewWorkerPool(ctx context.Context, maxWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	logger, _ := zap.NewProduction()

	return &WorkerPool{
		jobs:       make(chan *FileJob, 100),
		results:    make(chan string, 100),
		logger:     logger,
		ctx:        ctx,
		cancelFunc: cancel,
		maxWorkers: maxWorkers,
	}
}

// worker is the main worker goroutine that process file
func (wp *WorkerPool) worker(id int) {
	wp.wg.Add(1)
	defer wp.wg.Done()
	wp.logger.Info("Worker started", zap.Int("ID", id))

	for {
		select {
		case fileJob, ok := <-wp.jobs:
			if !ok {
				wp.logger.Info("Worker Shutting Down", zap.Int("ID", id))
				return
			}
			wp.processFile(fileJob)
		case <-wp.ctx.Done():
			wp.logger.Warn("Worker exiting due to cancellation", zap.Int("ID", id))
			return
		}
	}
}

// processFile handles the actual file reading with retry mechanism
func (wp *WorkerPool) processFile(FileJob *FileJob) {
	const maxRetries = 5

	for attempt := 1; attempt <= maxRetries; attempt++ {
		file, err := os.Open(FileJob.filepath)
		if err != nil {
			wp.logger.Error("Failed to open file", zap.String("File", FileJob.filepath), zap.Int("Attempt", attempt), zap.Error(err))
			time.Sleep(2 * time.Second)
			continue
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			wp.results <- fmt.Sprintf("Processed: %s", line)
		}
		file.Close()

		if err := scanner.Err(); err != nil {
			wp.logger.Error("Error reading file", zap.String("File", FileJob.filepath), zap.Error(err))
		}
		return
	}
	wp.logger.Error("File failed after retries", zap.String("File", FileJob.filepath))
}

// AddFile adds a new file to the processing queue and triggers worker scaling
func (wp *WorkerPool) AddFile(fileJob *FileJob) {
	select {
	case wp.jobs <- fileJob:
		wp.logger.Info("File added to queue", zap.String("File", fileJob.filepath))
		wp.scaleWorkers()
	case <-wp.ctx.Done():
		wp.logger.Warn("File ignored due to cancellation", zap.String("File", fileJob.filepath))
	}
}

// scaleWorkers dynamically adjusts the number of workers based on the queue size
func (wp *WorkerPool) scaleWorkers() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if len(wp.jobs) > cap(wp.jobs)/2 && wp.workerCount < wp.maxWorkers {
		wp.workerCount++
		go wp.worker(wp.workerCount)
		wp.logger.Info("Scaled up worker", zap.Int("NewWorkerCount", wp.workerCount))
	}
}

// Close performs graceful shutdown of the worker pool
func (wp *WorkerPool) Close() {
	wp.cancelFunc()
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
	wp.logger.Info("Worker Pool shutdown complete")
}

// MonitorResults starts a goroutine to collect and log processing results
func (wp *WorkerPool) MonitorResults() {
	go func() {
		for result := range wp.results {
			wp.logger.Info(result)
		}
	}()
}

// processFiles is responsible for the file processing workflow
func processFiles(ctx context.Context, files []string) {
	wp := NewWorkerPool(ctx, 10)
	defer wp.Close()

	wp.StartWorkers(3)

	pq := &FilePriorityQueue{}
	heap.Init(pq)

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			wp.logger.Error("Failed to get file info", zap.String("File", file), zap.Error(err))
			continue
		}
		heap.Push(pq, &FileJob{filepath: file, size: info.Size()})
	}

	for pq.Len() > 0 {
		fileJob := heap.Pop(pq).(*FileJob)
		wp.AddFile(fileJob)
	}
	wp.MonitorResults()

	time.Sleep(3 * time.Second)
}

// main is the entry point of the application
func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	files, err := filepath.Glob("files/*.txt")

	if err != nil {
		logger.Fatal("Failed to read directory", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	processFiles(ctx, files)
}
