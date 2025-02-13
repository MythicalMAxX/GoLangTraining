package main

import (
	"fmt"
	"strconv"
)

func main() {
	var input string
	fmt.Println("Enter the value to be checked")
	fmt.Scan(&input)

	var value interface{}
	if num, err := strconv.Atoi(input); err == nil {
		value = num
	} else if num, err := strconv.ParseFloat(input, 64); err == nil {
		value = num
	} else if input == "false" || input == "true" {
		value = "input" == "true"
	} else {
		value = input
	}

	switch value.(type) {
	case int:
		fmt.Println("Type: INT")
		break
	case float64:
		fmt.Println("Type: Float64")
		break
	case bool:
		fmt.Println("Type: Boolean")
		break
	case string:
		fmt.Println("Type: String")
		break
	default:
		fmt.Println("Invalid Datatype!")
	}
}
