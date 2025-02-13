// Write a program that demonstrates the difference between stack and heap allocations using new and make

package main

import (
	"fmt"
)

type head struct {
	name string
	age  int
}

func stack() {
	var s head
	s.name = "Stack"
	s.age = 23
	fmt.Println(s)
}

func heap() {
	h := new(head)
	h.name = "Heap"
	h.age = 23
	fmt.Println(h)
}

func main() {
	stack()
	heap()
}
