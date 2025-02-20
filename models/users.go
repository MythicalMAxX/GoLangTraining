package models

import (
	"time"
)

type Member struct {
	ID             uint
	Name           string
	Email          string
	Phone          string
	Address        string
	Membershiptype string
	JoinDate       time.Time
	Status         string
}

type Book struct {
	ID              uint
	Title           string
	ISBN            string
	Publisher       string
	PublicationYear string
	Category        string
	Edition         string
	Language        string
	PageCount       uint
}

type BookCopy struct {
	ID            uint
	BookID        uint
	Status        string
	ShelfLocation string
}

type Borrow struct {
	BorrowID   uint
	MemberID   uint
	CopyID     uint
	StaffID    uint
	BorrowDate time.Time
	DueDate    time.Time
	FineAmount float64
	Status     string
}

type Staff struct {
	ID          uint
	Name        string
	Role        string
	Department  string
	ContactInfo string
	HireDate    time.Time
	Shift       string
}

type Fine struct {
	ID            uint
	BorrowID      uint
	Amount        float64
	IssueDate     time.Time
	PaymentStatus string
	PaymentMethod string
}

type BookAuthor struct {
	BookID   uint
	AuthorID uint
}

type Author struct {
	ID             uint
	Name           string
	Biography      string
	Nationality    string
	BirthDate      time.Time
	Specialization string
}
