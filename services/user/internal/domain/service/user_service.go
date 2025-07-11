package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/robrt95x/godops/services/user/internal/domain/entity"
	"github.com/robrt95x/godops/services/user/internal/domain/port"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(name, email string) (*entity.User, error) {
	// Check if user already exists
	if existingUser, _ := s.repo.GetByEmail(email); existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	
	// Create new user
	user, err := entity.NewUser(name, email)
	if err != nil {
		return nil, err
	}
	
	// Generate ID
	user.ID = uuid.New().String()
	
	// Save user
	if err := s.repo.Save(user); err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *UserService) GetUserByID(id string) (*entity.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}
	
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

func (s *UserService) GetAllUsers() ([]*entity.User, error) {
	return s.repo.GetAll()
}
