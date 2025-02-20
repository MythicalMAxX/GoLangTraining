package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Enums for Member
type MembershipType string
type MemberStatus string

const (
	Standard MembershipType = "standard"
	Premium  MembershipType = "premium"
	VIP      MembershipType = "vip"
)

const (
	Active    MemberStatus = "active"
	Inactive  MemberStatus = "inactive"
	Suspended MemberStatus = "suspended"
)

type Member struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name"`
	Email          string         `json:"email" gorm:"unique"`
	Phone          string         `json:"phone"`
	Address        string         `json:"address"`
	Membershiptype MembershipType `json:"membershipType" gorm:"type:varchar(20);check:membershiptype IN ('standard', 'premium', 'vip')"`
	JoinDate       time.Time      `json:"joinDate"`
	Status         MemberStatus   `json:"status" gorm:"type:varchar(20);check:status IN ('active', 'inactive', 'suspended')"`
	Borrows        []Borrow       `json:"borrows" gorm:"foreignKey:UserID"`
}

// Add validation methods
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if err := m.validateEnums(); err != nil {
		return err
	}
	return nil
}

func (m *Member) BeforeUpdate(tx *gorm.DB) error {
	if err := m.validateEnums(); err != nil {
		return err
	}
	return nil
}

func (m *Member) validateEnums() error {
	// Validate MembershipType
	switch m.Membershiptype {
	case Standard, Premium, VIP:
		// valid
	default:
		return fmt.Errorf("invalid membership type: %s", m.Membershiptype)
	}

	// Validate Status
	switch m.Status {
	case Active, Inactive, Suspended:
		// valid
	default:
		return fmt.Errorf("invalid status: %s", m.Status)
	}

	return nil
}

// Borrow model with proper relations
type Borrow struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"userId"`
	BorrowerName string    `json:"borrowerName"`
	Member       Member    `json:"-" gorm:"foreignKey:UserID;references:ID"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Add hooks to automatically update BorrowerName when Member name changes
func (m *Member) AfterUpdate(tx *gorm.DB) error {
	// Update BorrowerName in all related Borrow records
	return tx.Model(&Borrow{}).
		Where("user_id = ?", m.ID).
		Update("borrower_name", m.Name).
		Error
}

// Add these methods to convert strings to enum types
func StringToMembershipType(s string) (MembershipType, error) {
	switch MembershipType(s) {
	case Standard, Premium, VIP:
		return MembershipType(s), nil
	default:
		return "", fmt.Errorf("invalid membership type: %s", s)
	}
}

func StringToMemberStatus(s string) (MemberStatus, error) {
	switch MemberStatus(s) {
	case Active, Inactive, Suspended:
		return MemberStatus(s), nil
	default:
		return "", fmt.Errorf("invalid status: %s", s)
	}
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

// Add this with your other request types
type BorrowRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

type BorrowResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"userId"`
	BorrowerName string    `json:"borrowerName"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Add these with your other request/response types
type GetMembersRequest struct {
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"pageSize,default=10"`
    Status   string `form:"status"`
    Type     string `form:"type"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64      `json:"total"`
    Page       int        `json:"page"`
    PageSize   int        `json:"pageSize"`
    TotalPages int        `json:"totalPages"`
}