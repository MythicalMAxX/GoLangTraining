// Question: Design a collection of books where you can iterate through the collection to find books based on different criteria (e.g., author, genre).
// Scenario: Create an Iterator interface and concrete iterators for different types of collections.

package main

import "fmt"

// Book struct
type Book struct{
	Name string
	Author string
	Genre string 
}

// BookCollectionInterface interface
type BookCollectionInterface interface{
	CreateIterator() Iterator

}

// Iterator interface
type Iterator interface{
	HasNext() bool
	GetNext() *Book
}

// BookCollection struct
type BookCollection struct{
	Books []*Book
}

// CreateIterator function
func (bc *BookCollection) CreateIterator() Iterator{
	return &BookIterator{books: bc.Books}
}

// BookIterator struct
type BookIterator struct{
	books []*Book
	position int 
}

// HasNext function
func (bi *BookIterator) HasNext() bool{
	if bi.position < len(bi.books) {
		return true
	}
	return false
}

// GetNext function
func (bi *BookIterator) GetNext() *Book {
	book := bi.books[bi.position]
	bi.position++
	return book
}
func main(){
	book1 := &Book{Name: "Book1", Author: "Author1", Genre: "Genre1"}
	book2 := &Book{Name: "Book2", Author: "Author2", Genre: "Genre2"}
	book3 := &Book{Name: "Book3", Author: "Author3", Genre: "Genre3"}

	bookCollection := &BookCollection{Books: []*Book{book1, book2, book3}}

	iterator := bookCollection.CreateIterator()

	for iterator.HasNext(){
		book := iterator.GetNext()
		fmt.Printf("Book: %s, Author: %s, Genre: %s\r\n", book.Name, book.Author, book.Genre)
	}
}