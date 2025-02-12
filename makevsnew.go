package main

import "fmt"

type node struct{
	name string
	age int
	next *node
}

func main(){
	newNode := new(node)
	makeMap := make(map[string]string)
	fmt.Println("newNode made using new function contains", newNode)
	fmt.Println("makeMap made using map function contains", makeMap)

	newArr := new([5]int)
	makeArr := make([]int,5)
	fmt.Println("newArr declared using new function contains", newArr)
	fmt.Println("makeArr declared using make function contains", makeArr)
}