package services

import (
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

func (s *MemberService) RegisterMember(req models.RegisterRequest) (*models.Member, error) {
	member := &models.Member{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		Address:        req.Address,
		Membershiptype: "standard",
		JoinDate:       time.Now(),
		Status:         "active",
	}

	err := s.memberRepo.Create(member)
	if err != nil {
		return nil, err
	}

	return member, nil
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
		member.Membershiptype = *req.Membershiptype
	}
	if req.JoinDate != nil {
		member.JoinDate = *req.JoinDate
	}
	if req.Status != nil {
		member.Status = *req.Status
	}

	err = s.memberRepo.Update(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *MemberService) DeleteMember(id string) error {
	return s.memberRepo.Delete(id)
}
