package services

import (
    "mypackage/internal/models"
    "mypackage/internal/repositories"
    "github.com/google/uuid"
)

type UserService struct {
    repo repositories.UserRepositoryInterface
}

func NewUserService(repo repositories.UserRepositoryInterface) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
    return s.repo.Create(user)
}

func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
    return s.repo.FindByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
    return s.repo.GetAll()
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
    return s.repo.FindByEmail(email)
}