package services

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"genshin-quiz/internal/models"
	"genshin-quiz/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUsers(limit, offset int, search string) (*models.ListResponse[models.User], error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var users []models.User
	var total int
	var err error

	if search != "" {
		users, total, err = s.userRepo.Search(search, limit, offset)
	} else {
		users, total, err = s.userRepo.GetAll(limit, offset)
	}

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
		return nil, models.ErrInvalidInput
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get user", 
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, models.ErrInvalidInput
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.logger.Error("Failed to get user by email", 
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && err != models.ErrUserNotFound {
		s.logger.Error("Failed to check existing user", 
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, models.ErrUserAlreadyExists
	}

	// Check username uniqueness
	existing, err = s.userRepo.GetByUsername(req.Username)
	if err != nil && err != models.ErrUserNotFound {
		s.logger.Error("Failed to check existing username", 
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check existing username: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Create user
	user, err := s.userRepo.Create(req)
	if err != nil {
		s.logger.Error("Failed to create user", 
			zap.String("username", req.Username),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Queue email verification task
	taskClient := tasks.NewClient(s.taskClient)
	if err := taskClient.EnqueueEmailVerification(user.ID, user.Email, "verification-token"); err != nil {
		s.logger.Warn("Failed to queue email verification", 
			zap.Int64("user_id", user.ID),
			zap.Error(err),
		)
	}

	s.logger.Info("User created successfully", 
		zap.Int64("user_id", user.ID),
		zap.String("username", user.Username),
	)

	return user, nil
}

func (s *UserService) UpdateUser(id int64, req models.UpdateUserRequest) (*models.User, error) {
	// Validate request
	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check username uniqueness if being updated
	if req.Username != nil {
		existing, err := s.userRepo.GetByUsername(*req.Username)
		if err != nil && err != models.ErrUserNotFound {
			return nil, fmt.Errorf("failed to check username uniqueness: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("username already exists")
		}
	}

	// Check email uniqueness if being updated
	if req.Email != nil {
		existing, err := s.userRepo.GetByEmail(*req.Email)
		if err != nil && err != models.ErrUserNotFound {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// Update user
	user, err := s.userRepo.Update(id, req)
	if err != nil {
		s.logger.Error("Failed to update user", 
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Info("User updated successfully", 
		zap.Int64("user_id", id),
	)

	return user, nil
}

func (s *UserService) DeleteUser(id int64) error {
	if id <= 0 {
		return models.ErrInvalidInput
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete user
	if err := s.userRepo.Delete(id); err != nil {
		s.logger.Error("Failed to delete user", 
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("User deleted successfully", 
		zap.Int64("user_id", id),
	)

	return nil
}

func (s *UserService) Authenticate(email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, models.ErrInvalidCredentials
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if err == models.ErrUserNotFound {
			return nil, models.ErrInvalidCredentials
		}
		s.logger.Error("Failed to get user for authentication", 
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Warn("Invalid password attempt", 
			zap.String("email", email),
			zap.Int64("user_id", user.ID),
		)
		return nil, models.ErrInvalidCredentials
	}

	s.logger.Info("User authenticated successfully", 
		zap.Int64("user_id", user.ID),
		zap.String("email", email),
	)

	return user, nil
}

func (s *UserService) UpdateUserStats(userID int64, scoreIncrement int) error {
	if userID <= 0 {
		return models.ErrInvalidInput
	}

	// Get current user to calculate new stats
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Update stats
	newTotalScore := user.TotalScore + scoreIncrement
	newQuizzesCompleted := user.QuizzesCompleted + 1

	if err := s.userRepo.UpdateStats(userID, newTotalScore, newQuizzesCompleted); err != nil {
		s.logger.Error("Failed to update user stats", 
			zap.Int64("user_id", userID),
			zap.Int("score_increment", scoreIncrement),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update user stats: %w", err)
	}

	// Queue statistics update task for further processing
	taskClient := tasks.NewClient(s.taskClient)
	data := map[string]interface{}{
		"score_increment": scoreIncrement,
		"total_score":     newTotalScore,
		"quizzes_completed": newQuizzesCompleted,
	}
	
	if err := taskClient.EnqueueUserStatisticsUpdate(userID, "quiz_completion", data); err != nil {
		s.logger.Warn("Failed to queue user statistics update", 
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
	}

	s.logger.Info("User stats updated", 
		zap.Int64("user_id", userID),
		zap.Int("score_increment", scoreIncrement),
		zap.Int("new_total_score", newTotalScore),
	)

	return nil
}

func (s *UserService) validateCreateUserRequest(req models.CreateUserRequest) error {
	if req.Username == "" || len(req.Username) < 3 || len(req.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}

	if req.Email == "" {
		return fmt.Errorf("email is required")
	}

	if req.Password == "" || len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	return nil
}

func (s *UserService) validateUpdateUserRequest(req models.UpdateUserRequest) error {
	if req.Username != nil && (len(*req.Username) < 3 || len(*req.Username) > 50) {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}

	if req.Email != nil && *req.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	return nil
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