package main

import "fmt"

func passbypointer(val *int) {
	fmt.Println("This is the pass by pointer function and here is pointer value:", &val)
}

func callbyvalue(val int) {
	fmt.Println("This is the call by value function and here is the passed value:", val)
}

func main() {
	var val int = 5
	passbypointer(&val)
	callbyvalue(val)
}
