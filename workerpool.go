// worker pool

package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan Job, results chan<- int, wg *sync.WaitGroup) {
	for j := range jobs {
		fmt.Printf("Worker %d starting job %d\n", id, j.id)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d finished job %d\n", id, j.id)
		results <- j.id * 2
		wg.Done()
	}
}

type Job struct {
	id int
}

func main() {
	const numWorker = 5
	const numJobs = 10
	jobs := make(chan Job, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	for i := 1; i <= numWorker; i++ {
		go worker(i, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		wg.Add(1)
		jobs <- Job{id: j}
	}
	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}

	fmt.Println("All workers finished")
}
