// Implement a worker pool to process jobs concurrently

package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Job struct {
	id  int
	url string
}

type Result struct {
	jobID int
	body  string
	err   error
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	for j := range jobs {
		fmt.Printf("Worker %d starting job %d\n", id, j.id)
		resp, err := http.Get(j.url)
		var body string
		if err == nil {
			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body) 
			if err == nil {
				body = string(bodyBytes)
			} else {
				err = fmt.Errorf("failed to read response body: %v", err)
			}
		}
		results <- Result{jobID: j.id, body: body, err: err}
		fmt.Printf("Worker %d finished job %d\n", id, j.id)
		wg.Done()
	}
}

func main() {
	const numWorker = 5
	const numJobs = 10
	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	for i := 1; i <= numWorker; i++ {
		go worker(i, jobs, results, &wg)
	}

	urls := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
		"https://jsonplaceholder.typicode.com/posts/4",
		"https://jsonplaceholder.typicode.com/posts/5",
		"https://jsonplaceholder.typicode.com/posts/6",
		"https://jsonplaceholder.typicode.com/posts/7",
		"https://jsonplaceholder.typicode.com/posts/8",
		"https://jsonplaceholder.typicode.com/posts/9",
		"https://jsonplaceholder.typicode.com/posts/10",
	}

	for j := 1; j <= numJobs; j++ {
		wg.Add(1)
		jobs <- Job{id: j, url: urls[j-1]}
	}
	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		if result.err != nil {
			fmt.Printf("Job %d failed: %v\n", result.jobID, result.err)
		} else {
			fmt.Printf("Job %d succeeded: %s\n", result.jobID, result.body)
		}
	}

	fmt.Println("All workers finished")
}
