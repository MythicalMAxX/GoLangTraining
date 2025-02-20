package models

import "time"

type Member struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	Membershiptype string    `json:"membershipType"`
	JoinDate       time.Time `json:"joinDate"`
	Status         string    `json:"status"`
}

// Request/Response types
type RegisterRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type RegisterResponse struct {
	UserID  uint        `json:"userID"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
	Status  int         `json:"status"`
}

type UpdateMemberRequest struct {
	Name           *string    `json:"name"`
	Email          *string    `json:"email"`
	Phone          *string    `json:"phone"`
	Address        *string    `json:"address"`
	Membershiptype *string    `json:"membershipType"`
	JoinDate       *time.Time `json:"joinDate"`
	Status         *string    `json:"status"`
}
