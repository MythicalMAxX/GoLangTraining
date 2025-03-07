// Question: Design a meal ordering system where a customer can build their meal by selecting different items (e.g., burger, drink, fries) with various options (e.g., size, type).
// Scenario: Create a MealBuilder class that allows customers to customize their meal and then build it step-by-step.

package main

import "fmt"

// Item interface
type Item interface {
	Name() string
	Packing() Packing
	Price() float64
}

// Packing interface
type Packing interface {
	Pack() string
}

// Burger struct
type Burger struct{}

func (b *Burger) Name() string {
	return "Burger"
}

func (b *Burger) Packing() Packing {
	return &Wrapper{}
}

func (b *Burger) Price() float64 {
	return 25.0
}

// Wrapper struct
type Wrapper struct{}

func (w *Wrapper) Pack() string {
	return "Wrapper"
}

// Drink struct
type Drink struct{}

func (d *Drink) Name() string {
	return "Drink"
}

func (d *Drink) Packing() Packing {
	return &Bottle{}
}

func (d *Drink) Price() float64 {
	return 10.0
}

// Bottle struct
type Bottle struct{}

func (b *Bottle) Pack() string {
	return "Bottle"
}

// Meal struct
type Meal struct {
	items []Item
}

func (m *Meal) AddItem(item Item) {
	m.items = append(m.items, item)
}

func (m *Meal) GetCost() float64 {
	cost := 0.0
	for _, item := range m.items {
		cost += item.Price()
	}
	return cost
}

func (m *Meal) ShowItems() {
	for _, item := range m.items {
		fmt.Printf("Item: %s\n", item.Name())
		fmt.Printf("Packing: %s\n", item.Packing().Pack())
		fmt.Printf("Price: %.2f\n", item.Price())
	}
}

// MealBuilder struct
type MealBuilder struct{}

func (mb *MealBuilder) PrepareVegMeal() *Meal {
	meal := &Meal{}
	meal.AddItem(&Burger{})
	meal.AddItem(&Drink{})
	return meal
}

func (mb *MealBuilder) PrepareNonVegMeal() *Meal {
	meal := &Meal{}
	meal.AddItem(&Burger{})
	meal.AddItem(&Drink{})
	return meal
}

func main() {
	mealBuilder := &MealBuilder{}

	vegMeal := mealBuilder.PrepareVegMeal()
	fmt.Println("Veg Meal")
	vegMeal.ShowItems()
	fmt.Printf("Total Cost: %.2f\n", vegMeal.GetCost())

	nonVegMeal := mealBuilder.PrepareNonVegMeal()
	fmt.Println("Non-Veg Meal")
	nonVegMeal.ShowItems()
	fmt.Printf("Total Cost: %.2f\n", nonVegMeal.GetCost())
}
