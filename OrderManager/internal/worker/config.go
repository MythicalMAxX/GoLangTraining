package worker

import "time"

type Config struct {
	InitialWorkers  int           // Initial number of workers
	MinWorkers      int           // Minimum number of workers
	MaxWorkers      int           // Maximum number of workers
	QueueSize       int           // Size of the job queue
	ProcessInterval time.Duration // Interval to process pending orders
	ScaleInterval   time.Duration // Interval to check for scaling
	ScaleThreshold  float64       // Queue utilization threshold for scaling
}

func DefaultConfig() *Config {
	return &Config{
		InitialWorkers:  5,
		MinWorkers:      3,
		MaxWorkers:      20,
		QueueSize:       100,
		ProcessInterval: 30 * time.Second,
		ScaleInterval:   1 * time.Minute,
		ScaleThreshold:  0.7, // 70% queue utilization
	}
}
