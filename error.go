package main

import (
	"errors"
	"fmt"
)

var ErrorMessage = errors.New("My third Error Message")

type MyCustomError struct {
	message string
}

func (e *MyCustomError) Error() string {
	return e.message
}

func main() {
	// Error using fmt
	err1 := fmt.Errorf("This is the first error using fmt")
	err2 := errors.New("This is my second error message using errors.New()")

	fmt.Println(err1)
	fmt.Println(err2)

	err3 := fmt.Errorf("This is third error message %w", ErrorMessage)
	if errors.Is(err3, ErrorMessage) {
		fmt.Printf("This is my third Error Message using errors.Is()\n")
	}

	err4 := fmt.Errorf("This is the fourth error message %w", &MyCustomError{"4th Error"})
	var myerr *MyCustomError
	if errors.As(err4, &myerr) {
		fmt.Println("This is my fourth error message using errors.As()")
	}

}
