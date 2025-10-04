package services

import (
	"fmt"

	"genshin-quiz-backend/internal/models"
	"genshin-quiz-backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUsers(limit, offset int) (*models.ListResponse[models.User], error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.userRepo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return &models.ListResponse[models.User]{
		Data:   users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *UserService) GetUser(id int64) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if username already exists
	existing, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	existing, err = s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("email already exists")
	}

	user, err := s.userRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(id int64, req models.UpdateUserRequest) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Validate request
	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if user exists
	existing, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check username uniqueness if being updated
	if req.Username != nil && *req.Username != existing.Username {
		user, err := s.userRepo.GetByUsername(*req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if user != nil {
			return nil, fmt.Errorf("username already exists")
		}
	}

	// Check email uniqueness if being updated
	if req.Email != nil && *req.Email != existing.Email {
		user, err := s.userRepo.GetByEmail(*req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if user != nil {
			return nil, fmt.Errorf("email already exists")
		}
	}

	user, err := s.userRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserService) DeleteUser(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	err := s.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) validateCreateUserRequest(req models.CreateUserRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	// Add more validation as needed
	return nil
}

func (s *UserService) validateUpdateUserRequest(req models.UpdateUserRequest) error {
	if req.Username != nil {
		if len(*req.Username) < 3 || len(*req.Username) > 50 {
			return fmt.Errorf("username must be between 3 and 50 characters")
		}
	}
	// Add more validation as needed
	return nil
}