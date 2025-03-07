// Question: Develop a shape creation system where different shapes (e.g., circle, square, rectangle) can be created without specifying the exact class of the object that will be created.
// Scenario: Create a ShapeFactory class that returns instances of different shapes based on the input provided.

package main

import "fmt"

// Shape interface
type Shape interface {
	Draw()
}

// Circle struct
type Circle struct{}

func (c *Circle) Draw() {
	fmt.Println("Circle")
}

// Square struct
type Square struct{}

func (s *Square) Draw() {
	fmt.Println("Square")
}

// Rectangle struct
type Rectangle struct{}

func (r *Rectangle) Draw() {
	fmt.Println("Rectangle")
}

// ShapeFactory struct
type ShapeFactory struct{}

// GetShape function
func (sf *ShapeFactory) GetShape(shapeType string) Shape {
	switch shapeType {
	case "CIRCLE":
		return &Circle{}
	case "SQUARE":
		return &Square{}
	case "RECTANGLE":
		return &Rectangle{}
	}
	return nil
}


func main() {
	shapeFactory := &ShapeFactory{}

	// Get an object of Circle and call its Draw method
	circle := shapeFactory.GetShape("CIRCLE")
	circle.Draw()

	// Get an object of Square and call its Draw method
	square := shapeFactory.GetShape("SQUARE")
	square.Draw()

	// Get an object of Rectangle and call its Draw method
	rectangle := shapeFactory.GetShape("RECTANGLE")
	rectangle.Draw()
}