package main

import (
	"fmt"
	"sync"
	"time"
)

func chef(dishchan chan string) {
	time.Sleep(2 * time.Second)
	dish := "Pizza"
	fmt.Println("Chef prepared", dish)
	dishchan <- dish
}

func waiter(dishchan chan string) {
	dish := <-dishchan
	fmt.Println("Waiter served", dish)
}

func waitFunc() {
	for i := 0; i < 10; i++ {
		fmt.Println("Hello World!")
		time.Sleep(100 * time.Millisecond)
	}
}

func worker(id int, wg sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting.\n", id)

	for i := 0; i < 3; i++ {
		fmt.Printf("Worker %d doing job %d\n", id, i)
	}
	fmt.Printf("Worker %d done.\n", id)
}

func main() {
	go waitFunc()
	time.Sleep(2 * time.Second)
	fmt.Println("Program Executed")

	dishchan := make(chan string)
	go chef(dishchan)
	go waiter(dishchan)

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, wg)
	}

	wg.Wait()
	fmt.Println("All workers done.")

}
