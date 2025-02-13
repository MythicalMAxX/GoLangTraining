package main

import "fmt"

const PI = 3.14

const (
	Zero = iota
	One
	Two
	Three
	Four
	Five
	Six
)

func main() {
	fmt.Println(Zero)
	fmt.Println(One)
	fmt.Println(Two)
	fmt.Println(Three)
	fmt.Println(Four)
	fmt.Println(Five)
	fmt.Println(Six)
	fmt.Println(PI)
}
