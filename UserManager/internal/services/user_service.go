package services

import (
    "errors"
    "userservices/internal/models"
    "userservices/internal/repositories"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "log"
)

type UserService struct {
    repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) Register(req *models.RegisterRequest) (*models.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Set default role if not provided
    role := "user"
    if req.Role != "" {
        role = req.Role
    }

    user := &models.User{
        ID:       uuid.New().String(),
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     role,
    }

    if err := s.repo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *UserService) Login(req *models.LoginRequest) (*models.User, error) {
    log.Printf("Attempting to fetch user with email: %s", req.Email)
    user, err := s.repo.GetByEmail(req.Email)
    if err != nil {
        log.Printf("Database error: Failed to get user by email %s: %v", req.Email, err)
        return nil, errors.New("invalid credentials")
    }

    log.Printf("Comparing password for user: %s", user.ID)
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        log.Printf("Password mismatch for user %s", user.ID)
        return nil, errors.New("invalid credentials")
    }

    log.Printf("User %s authenticated successfully", user.ID)
    return user, nil
}

// GetAllUsers returns all users from the database
func (s *UserService) GetAllUsers() ([]*models.User, error) {
    return s.repo.GetAll()
}