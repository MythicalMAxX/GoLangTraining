package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type node struct {
	Name string `json "name"`
	Age  int    `json "age"`
	Next *node  `json "json *node"`
}

func main() {
	Node := node{Name: "Vinamra Yadav", Age: 24, Next: nil}
	jsonData, _ := json.Marshal(Node)
	fmt.Println(string(jsonData))

	n := reflect.TypeOf(Node)
	for i := 0; i < n.NumField(); i++ {
		row := n.Field(i)
		fmt.Printf("Tag: %v		Name: %v\n", row.Tag, row.Name)
	}
}
