package main

import (
	"fmt"
	"mypackage/config"
	"mypackage/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func updateMemberByID(c *gin.Context) {}
func modifyMemberByID(c *gin.Context) {}

func deleteMemberByID(c *gin.Context) {
	// Get id from URL parameter
	id := c.Param("id")

	// Delete member from database
	result := db.Delete(&models.Member{}, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}

// Add this with your other type definitions
type MemberResponse struct {
	UserID         uint      `json:"userID"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	Membershiptype string    `json:"membershipType"`
	JoinDate       time.Time `json:"joinDate"`
	Status         string    `json:"status"`
}

func getMemberByID(c *gin.Context) {
	// Get id from URL parameter
	id := c.Param("id")

	var member models.Member
	result := db.First(&member, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve member"})
		return
	}

	// response := MemberResponse{
	//     UserID:        member.ID,
	//     Name:          member.Name,
	//     Email:         member.Email,
	//     Phone:         member.Phone,
	//     Address:       member.Address,
	//     Membershiptype: member.Membershiptype,
	//     JoinDate:      member.JoinDate,
	//     Status:        member.Status,
	// }

	response := RegisterResponse{
		UserID:  member.ID,
		Message: "Request Successful",
		Body:    member,
		Status:  200,
	}

	c.JSON(http.StatusOK, response)
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
}

type RegisterResponse struct {
	UserID  uint        `json:"userID"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
	Status  int         `json:"status"`
}

func registerUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new member
	member := models.Member{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		Address:        req.Address,
		Membershiptype: "standard",
		JoinDate:       time.Now(),
		Status:         "active",
	}

	// Save to database
	result := db.Create(&member)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Prepare response
	response := RegisterResponse{
		UserID:  member.ID,
		Message: "User registered successfully",
		Body:    member,
		Status:  201,
	}

	c.JSON(http.StatusCreated, response)
}

type Status struct {
	Status string `json:"status"`
	Time   string `json:"timestamp"`
}

func getStatus(c *gin.Context) {
	status := Status{
		Status: "running",
		Time:   time.Now().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, status)
}

func main() {
	var err error
	db, err = config.InitDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	err = db.AutoMigrate(&models.Member{}, &models.Book{}, &models.BookCopy{}, &models.Borrow{}, &models.Staff{}, &models.Fine{})
	if err != nil {
		fmt.Println("Failed to migrate tables:", err)
		return
	}

	// Initialize Gin router
	router := gin.Default()

	// Add routes
	router.GET("/api/v1/running", getStatus)
	router.GET("/api/v1/member/:id", getMemberByID)
	router.POST("/api/v1/member/register", registerUser)
	router.DELETE("api/v1/member/:id", deleteMemberByID)
	router.PATCH("api/v1/member/:id", updateMemberByID)
	router.PUT("api/v1/member/:id", modifyMemberByID)

	// Start the server
	router.Run(":8080")
}
