// Day 1: Go Environment & Basic Syntax
// Setting up Go environment (GOPATH, Go Modules, workspace)
// Go tools ecosystem (go mod, go fmt, go vet, golint)
// Package management and visibility (public/private)
// Multiple return values and named returns
// Basic syntax and common patterns
// Arrays vs Slices (internal structure)
// Interfaces

package main

import (
	"GoLangTraining/mypackage"
	"fmt"
)

// PI is a constant value
const PI = 3.14

// Struct
type node struct {
	name string
	age  int
	next *node
}

type calculator interface {
	sums()
	subtract()
	multiply()
	divide()
}

type nums struct {
	num1, num2 int
}

func (n nums) sums() {
	fmt.Println(n.num1 + n.num2)
}

func (n nums) subtract() {
	fmt.Println(n.num1 - n.num2)
}

func (n nums) divide() {
	if n.num2 == 0 {
		fmt.Println("Zero Division Error, Num2 can't be zero")
		return
	}
	fmt.Println(n.num1 / n.num2)
}

func (n nums) multiply() {
	fmt.Println(n.num1 * n.num2)
}

// Function(single return value)
func sum(num1 int, num2 int) int {
	return num1 + num2
}

// Function(multiple return values)
func minmax(num1, num2 int) (int, int) {
	if num1 > num2 {
		return num2, num1
	}
	return num1, num2
}

// Function(Named return)
func returnName() (X string) {
	X = "Vinamra Yadav"
	return
}

func main() {
	// basic syntax
	fmt.Println("Hello World")

	// variables

	var name string = "Vinamra Yadav"
	var age int = 24
	var isCool = true
	var height float64 = 174.5
	var mynumber int64 = 9936476927
	fmt.Printf("Hi I am %s and I am %d years old\n. I am %v tall. Contact me using %v.\n", name, age, height, mynumber)
	fmt.Print("Am I cool? ", isCool)
	fmt.Println()

	// short hand
	firstname := "Vinamra"
	lastname := "Yadav"

	fmt.Println(firstname, lastname)

	// constants
	fmt.Println("Value of PI is", PI)

	// Operators
	fmt.Println("Sum of 5 and 6 is", 5+6)
	fmt.Println("Subtraction of 5 and 6 is", 5-6)
	fmt.Println("Division of 5 and 6 is", 5/6)
	fmt.Println("Multiplication of 5 and 6 is", 5*6)
	fmt.Println("Modulus of 5 and 6 is", 5%6)

	// Conditional Operators
	fmt.Println("Is 5 greater than 6?", 5 > 6)
	fmt.Println("Is 5 equal to 5", 5 == 5)
	fmt.Println("Is 20 greater than equal to 21", 20 >= 21)

	//  Bitwise Operators

	fmt.Println("Bitwise AND of 5 and 6 is", 5&6)
	fmt.Println("Bitwise OR of 5 and 6 is", 5|6)
	fmt.Println("Bitwise XOR of 5 and 6 is", 5^6)

	// Loops and if else
	for i := 0; i < 10; i++ {
		if i == 5 {
			continue
		} else if i == 9 {
			break
		} else {
			fmt.Println(i)
		}
	}

	// Iteration over an array

	var arr = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for ind, val := range arr {
		fmt.Println("Index:", ind, "Value:", val)
	}

	// Slicing

	var slicedarr = arr[1:5]
	fmt.Println(slicedarr)

	// Switch

	var operator = "/"
	var num1 int = 10
	var num2 int = 20

	switch operator {
	case "+":
		fmt.Println("Addition Selected", num1+num2)
		break

	case "-":
		fmt.Println("Subtraction Selected", num1-num2)
		break

	case "*":
		fmt.Println("Multiplication Selected", num1*num2)
		break

	case "/":
		fmt.Print("Division Selected ")
		if num2 == 0 {
			fmt.Println("Zero Division Error, Num2 can't be zero")
			break
		} else {
			fmt.Println(num1 / num2)
		}
		break
	case "%":
		fmt.Println("Modulus Selected", num1%num2)
		break
	default:
		fmt.Println("Invalid Operator")
		break

	}

	// struct

	var newnode node

	newnode.name = "Vinamra Yadav"
	newnode.age = 24
	newnode.next = nil

	fmt.Printf("Hi my name is %v and I am %v years old", newnode.name, newnode.age)

	// function
	fmt.Println("Sum of 12 and 15 is", sum(12, 15))

	// Pointers
	var myname string = "Vinamra Yadav"
	var namepointer *string = &myname
	fmt.Printf("My name is %v and its pointer is %v.\n", myname, namepointer)

	// map
	var mymap = make(map[string]string)
	mymap["firstname"] = "Vinamra"
	mymap["lastname"] = "Yadav"
	fmt.Println("My first name is", mymap["firstname"])

	mypackage.PublicFunc()

	var (
		minimum int
		maximum int
	)
	minimum, maximum = minmax(10, 20)
	fmt.Printf("Minimum Value is %v and maximum value is %v", minimum, maximum)

	fmt.Println("Value from named return function is: ", returnName())

	var calc calculator
	calc = nums{num1: 20, num2: 21}
	fmt.Print("Sum of both the numbers is: ")
	calc.sums()
	calc.subtract()
	calc.multiply()
	calc.divide()
}
