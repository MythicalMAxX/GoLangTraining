package services

import (
	"fmt"
	"mypackage/internal/models"
	"mypackage/internal/repositories"
	"time"
)

type MemberService struct {
	memberRepo *repositories.MemberRepository
}

func NewMemberService(memberRepo *repositories.MemberRepository) *MemberService {
	return &MemberService{memberRepo: memberRepo}
}

func (s *MemberService) GetMember(id string) (*models.Member, error) {
	return s.memberRepo.FindByID(id)
}

func (s *MemberService) UpdateMember(id string, req *models.Member) (*models.Member, error) {
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Preserve the ID
	req.ID = member.ID

	err = s.memberRepo.Update(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *MemberService) DeleteMember(id string) error {
	return s.memberRepo.Delete(id)
}

func (s *MemberService) ModifyMember(id string, req *models.UpdateMemberRequest) (*models.Member, error) {
	member, err := s.memberRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		member.Name = *req.Name
	}
	if req.Email != nil {
		member.Email = *req.Email
	}
	if req.Phone != nil {
		member.Phone = *req.Phone
	}
	if req.Address != nil {
		member.Address = *req.Address
	}
	if req.Membershiptype != nil {
		membershipType, err := models.StringToMembershipType(*req.Membershiptype)
		if err != nil {
			return nil, fmt.Errorf("invalid membership type: %s", *req.Membershiptype)
		}
		member.Membershiptype = membershipType
	}
	if req.JoinDate != nil {
		member.JoinDate = *req.JoinDate
	}
	if req.Status != nil {
		status, err := models.StringToMemberStatus(*req.Status)
		if err != nil {
			return nil, fmt.Errorf("invalid status: %s", *req.Status)
		}
		member.Status = status
	}

	err = s.memberRepo.Update(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *MemberService) RegisterMember(req models.RegisterRequest) (*models.Member, error) {
	member := &models.Member{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		Address:        req.Address,
		Membershiptype: models.Standard, // Using enum directly
		JoinDate:       time.Now(),
		Status:         models.Active, // Using enum directly
	}

	err := s.memberRepo.Create(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *MemberService) CreateBorrow(req models.BorrowRequest) (*models.Borrow, error) {
	// First get the member to ensure they exist and get their name
	member, err := s.memberRepo.FindByID(fmt.Sprintf("%d", req.UserID))
	if err != nil {
		return nil, fmt.Errorf("member not found: %v", err)
	}

	borrow := &models.Borrow{
		UserID:       member.ID,
		BorrowerName: member.Name,
		CreatedAt:    time.Now(),
	}

	err = s.memberRepo.CreateBorrow(borrow)
	if err != nil {
		return nil, fmt.Errorf("failed to create borrow record: %v", err)
	}

	return borrow, nil
}


func (s *MemberService) GetMembers(req models.GetMembersRequest) (*models.PaginatedResponse, error) {
    // Validate and prepare filters
    filters := make(map[string]interface{})
    
    if req.Status != "" {
        status, err := models.StringToMemberStatus(req.Status)
        if err == nil { // Only add if valid
            filters["status"] = status
        }
    }
    
    if req.Type != "" {
        membershipType, err := models.StringToMembershipType(req.Type)
        if err == nil { // Only add if valid
            filters["membershiptype"] = membershipType
        }
    }
    
    // Get members with pagination
    members, total, err := s.memberRepo.GetMembers(req.Page, req.PageSize, filters)
    if err != nil {
        return nil, err
    }
    
    // Calculate total pages
    totalPages := int(total) / req.PageSize
    if int(total)%req.PageSize != 0 {
        totalPages++
    }
    
    return &models.PaginatedResponse{
        Data:       members,
        Total:      total,
        Page:       req.Page,
        PageSize:   req.PageSize,
        TotalPages: totalPages,
    }, nil
}