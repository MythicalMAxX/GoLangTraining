package main

import (
	"fmt"
)

func helper() {
	fmt.Println("Deffered function")
}

func main() {
	fmt.Println("Program started")
	defer helper()
	fmt.Println("All task done")
}
