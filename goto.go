package main

import (
	"fmt"
)

func makeChoice() int {
	var res int
	fmt.Scan(&res)
	return res
}

func main() {
	fmt.Println("Start")
	fmt.Println("In any number:")
	choice := makeChoice()
	if choice%2 == 0 {
		goto Even
	} else {
		goto Odd
	}
Even:
	fmt.Println("You've choosen Even Number")
Odd:
	fmt.Println("You've choosen Odd Number")
}
