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

func (r *MemberRepository) CreateBorrow(borrow *models.Borrow) error {
	return r.db.Create(borrow).Error
}

func (r *MemberRepository) GetBorrowsByUserID(userID uint) ([]models.Borrow, error) {
	var borrows []models.Borrow
	err := r.db.Where("user_id = ?", userID).Find(&borrows).Error
	return borrows, err
}


func (r *MemberRepository) GetMembers(page, pageSize int, filters map[string]interface{}) ([]models.Member, int64, error) {
    var members []models.Member
    var total int64
    
    query := r.db.Model(&models.Member{})
    
    // Apply filters if they exist
    for key, value := range filters {
        if value != "" {
            query = query.Where(key+" = ?", value)
        }
    }
    
    // Get total count
    err := query.Count(&total).Error
    if err != nil {
        return nil, 0, err
    }
    
    // Get paginated data
    err = query.Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&members).Error
    if err != nil {
        return nil, 0, err
    }
    
    return members, total, nil
}