package handlers

import (
	"mypackage/internal/models"
	"mypackage/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService *services.MemberService
}

func NewMemberHandler(memberService *services.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

func (h *MemberHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.memberService.RegisterMember(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{
		UserID:  member.ID,
		Message: "User registered successfully",
		Body:    member,
		Status:  201,
	})
}

func (h *MemberHandler) Get(c *gin.Context) {
	id := c.Param("id")
	member, err := h.memberService.GetMember(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	c.JSON(http.StatusOK, models.RegisterResponse{
		UserID:  member.ID,
		Message: "Request successful",
		Body:    member,
		Status:  200,
	})
}

func (h *MemberHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req models.Member
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.memberService.UpdateMember(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update member"})
		return
	}

	c.JSON(http.StatusOK, models.RegisterResponse{
		UserID:  member.ID,
		Message: "Member updated successfully",
		Body:    member,
		Status:  200,
	})
}

func (h *MemberHandler) Modify(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.memberService.ModifyMember(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to modify member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Member modified successfully",
		"data":    member,
	})
}

func (h *MemberHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.memberService.DeleteMember(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}

func (h *MemberHandler) AddBorrow(c *gin.Context) {
	var req models.BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrow, err := h.memberService.CreateBorrow(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{
		UserID:  borrow.UserID,
		Message: "Borrow record created successfully",
		Body:    borrow,
		Status:  201,
	})
}

func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
}

func (h *MemberHandler) GetMembers(c *gin.Context) {
    var req models.GetMembersRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Set default values if not provided
    if req.Page <= 0 {
        req.Page = 1
    }
    if req.PageSize <= 0 {
        req.PageSize = 10
    }
    
    response, err := h.memberService.GetMembers(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
        return
    }
    
    c.JSON(http.StatusOK, response)
}