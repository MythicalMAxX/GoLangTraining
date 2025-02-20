package repositories

import (
    "mypackage/internal/models"
    "gorm.io/gorm"
)

type MemberRepository struct {
    db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
    return &MemberRepository{db: db}
}

func (r *MemberRepository) FindByID(id string) (*models.Member, error) {
    var member models.Member
    result := r.db.First(&member, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &member, nil
}

func (r *MemberRepository) Create(member *models.Member) error {
    return r.db.Create(member).Error
}

func (r *MemberRepository) Update(member *models.Member) error {
    return r.db.Save(member).Error
}

func (r *MemberRepository) Delete(id string) error {
    return r.db.Delete(&models.Member{}, id).Error
}