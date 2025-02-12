package main

import "fmt"

// represents a stack data structure
type Stack struct {
	items []int
}

// adds an element to the stack
func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

// removes and returns the top element of the stack
func (s *Stack) Pop() (int, error) {
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

func main() {
	stack := Stack{}
	stack.Push(11)
	stack.Push(22)
	stack.Push(33)
	fmt.Println("Stack:", stack.items)

	item, err := stack.Pop()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Popped:", item)
	}
	fmt.Println("Stack after pop:", stack.items)
}
