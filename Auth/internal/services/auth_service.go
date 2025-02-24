package services

import (
	"errors"
	"fmt"
	"mypackage/internal/models"
	"mypackage/internal/repositories"
	"mypackage/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.AuthRepository
}

func NewAuthService(repo *repositories.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
	Role     models.Role
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	// Check if email already exists
	_, exists, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	// Validate role
	if !input.Role.IsValid() {
		input.Role = models.RoleUser // Default to user role if invalid
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(input LoginInput) (string, error) {
	user, _, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate token with user ID and role
	token, err := jwt.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
