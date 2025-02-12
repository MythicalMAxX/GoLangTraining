package main

import "fmt"

// Queue data structure
type Queue struct {
	items []int
}

// Enqueue adds an element to the queue
func (q *Queue) Enqueue(item int) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front element of the queue
func (q *Queue) Dequeue() (int, error) {
	if len(q.items) == 0 {
		return 0, fmt.Errorf("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func main() {
	queue := Queue{}
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)
	fmt.Println("Queue:", queue.items)

	item, err := queue.Dequeue()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Dequeued:", item)
	}
	fmt.Println("Queue after dequeue:", queue.items)
}
