package main

import (
	"fmt"
)

func main() {
	var choice int
	fmt.Print("Choose the input type\n1 for Byte arr\n2 for string:")
	fmt.Scan(&choice)

	if choice == 1 {
		fmt.Println("Enter the length of slice")
		var SIZE int
		fmt.Scan(&SIZE)
		byteslice := []byte{}
		fmt.Println("Type your Byte Array")
		for i := 0; i < SIZE; i++ {
			var BYTE byte
			fmt.Scan(&BYTE)
			byteslice = append(byteslice, BYTE)
		}
		fmt.Println("After Convert Byte Slice to String: ", string(byteslice))
	} else if choice == 2 {
		fmt.Println("Enter your String")
		var str string
		fmt.Scan(&str)
		res := []byte(str)
		fmt.Println("After Converting String to Byte Slice: ", res)
	} else {
		fmt.Print("Invalid Input! Please try again later.")
	}
}
