package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func convertToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Invalid Value! Please try again with valid string.")
		os.Exit(0)
	}
	return num
}

func convertToFloat(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Invalid Value! Please try again with valid string.")
		os.Exit(0)
	}
	return num
}

func main() {
	var str string
	var choice int
	fmt.Print("Enter your string input to be converted: ")
	fmt.Scan(&str)
	fmt.Print("\nEnter your choice, \n1 for convertion to Int\n2 for convertion to Float.\n")
	fmt.Scan(&choice)

	if choice == 1 {
		res := convertToInt(str)
		fmt.Printf("String Converted to Int! Value: %d, Type: %v.\n", res, reflect.TypeOf(res))
	} else if choice == 2 {
		res := convertToFloat(str)
		fmt.Printf("String Converted to Float! Value: %f, Type: %v.\n", res, reflect.TypeOf(res))
	} else {
		fmt.Println("Invalid Choice! Please try again with valid input")
		os.Exit(0)
	}

}
