package repositories

import (
	"errors"
	"mypackage/internal/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindByEmail(email string) (*models.User, bool, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false, nil // User not found, but not an error
		}
		return nil, false, result.Error // Database error
	}
	return &user, true, nil // User found
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		// Check for unique constraint violation
		if strings.Contains(result.Error.Error(), "unique constraint") {
			return errors.New("email already registered")
		}
		return result.Error
	}
	return nil
}

func (r *AuthRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
