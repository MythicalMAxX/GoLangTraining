# 📌 Concurrent File Processing System with Dynamic Worker Pool
## 🔹 Problem Statement:
Design and implement a Concurrent File Processing System in Go that efficiently reads and processes multiple files in parallel. The system should leverage Goroutines, Worker Pools, Channels, Mutexes, and Context for Cancellation & Timeouts.
 
## 🔹 Requirements:
### 1️⃣ Concurrent File Processing:
- Read multiple large text files concurrently.
- Each line in the file should be processed independently.
- Use worker pools to manage concurrent file processing efficiently.
### 2️⃣ Worker Pool with Dynamic Scaling:
- Implement a worker pool where the number of workers adjusts dynamically based on system load.
- Ensure that new workers are spawned if the queue is growing too fast.
### 3️⃣ Rate Limiting & Synchronization:
- Use channels to control the flow of file processing.
- Implement rate limiting to avoid overwhelming system resources.
- Use sync.Mutex where necessary to prevent race conditions.
### 4️⃣ Graceful Shutdown & Cancellation:
- Use the context package to allow users to cancel processing at any time.
- Implement timeouts for long-running operations.
### 5️⃣ Logging & Monitoring:
- Replace print statements with structured logging using logrus or zap.
### 6️⃣ Optimise:
- Implement a retry mechanism for failed file reads.
- Optimize worker assignment by using a priority queue for large files. 	 	 	 	 	 	 
