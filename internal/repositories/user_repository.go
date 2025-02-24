package repositories

import (
    "mypackage/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

type UserRepositoryInterface interface {
    Create(user *models.User) error
    FindByID(id uuid.UUID) (*models.User, error)
    FindByEmail(email string) (*models.User, error)  // Add this method
    GetAll() ([]models.User, error)
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
    var users []models.User
    err := r.db.Find(&users).Error
    return users, err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}