package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	var num1 int = 1
	var num2 float64 = 10.6
	var string1 string = "12"
	var string2 string = "12.6"

	convertednum1 := strconv.Itoa(num1)
	convertednum2 := strconv.FormatFloat(num2, 'f', -1, 64)
	convertedstring1, _ := strconv.Atoi(string1)
	convertedstring2, _ := strconv.ParseFloat(string2, 64)
	fmt.Printf("After converting num1 to string: %v Type: %v.\n", convertednum1, reflect.TypeOf(convertednum1))
	fmt.Println("After converting num2(float) to string: " + convertednum2)
	fmt.Println("After converting string1 to int: ", convertedstring1)
	fmt.Println("After converting strig2 to float: ", convertedstring2)
}
