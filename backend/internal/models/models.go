package models

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrQuizNotFound = errors.New("quiz not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrQuizAlreadyExists = errors.New("quiz already exists")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
)

// User represents a user in the system
type User struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"` // Never include in JSON responses
	DisplayName      *string   `json:"display_name,omitempty"`
	AvatarURL        *string   `json:"avatar_url,omitempty"`
	TotalScore       int       `json:"total_score"`
	QuizzesCompleted int       `json:"quizzes_completed"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateUserRequest represents the data needed to create a new user
type CreateUserRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=50"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,max=100"`
	AvatarURL   *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

// UpdateUserRequest represents the data needed to update a user
type UpdateUserRequest struct {
	Username    *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,max=100"`
	AvatarURL   *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

// Quiz represents a quiz in the system
type Quiz struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Category    string     `json:"category"`
	Difficulty  string     `json:"difficulty"`
	Questions   []Question `json:"questions,omitempty"`
	TimeLimit   *int       `json:"time_limit,omitempty"`
	CreatedBy   int64      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateQuizRequest represents the data needed to create a new quiz
type CreateQuizRequest struct {
	Title       string                    `json:"title" validate:"required,min=1,max=200"`
	Description *string                   `json:"description,omitempty" validate:"omitempty,max=1000"`
	Category    string                    `json:"category" validate:"required,oneof=characters weapons artifacts lore gameplay"`
	Difficulty  string                    `json:"difficulty" validate:"required,oneof=easy medium hard"`
	Questions   []CreateQuestionRequest   `json:"questions" validate:"required,min=1,dive"`
	TimeLimit   *int                      `json:"time_limit,omitempty" validate:"omitempty,min=30,max=3600"`
	CreatedBy   int64                     `json:"created_by" validate:"required"`
}

// UpdateQuizRequest represents the data needed to update a quiz
type UpdateQuizRequest struct {
	Title       *string                  `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string                  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Category    *string                  `json:"category,omitempty" validate:"omitempty,oneof=characters weapons artifacts lore gameplay"`
	Difficulty  *string                  `json:"difficulty,omitempty" validate:"omitempty,oneof=easy medium hard"`
	Questions   []CreateQuestionRequest  `json:"questions,omitempty" validate:"omitempty,min=1,dive"`
	TimeLimit   *int                     `json:"time_limit,omitempty" validate:"omitempty,min=30,max=3600"`
}

// Question represents a quiz question
type Question struct {
	ID           int64     `json:"id"`
	QuizID       int64     `json:"quiz_id"`
	QuestionText string    `json:"question_text"`
	QuestionType string    `json:"question_type"`
	Options      []string  `json:"options,omitempty"`
	CorrectAnswer string   `json:"correct_answer"`
	Explanation  *string   `json:"explanation,omitempty"`
	Points       int       `json:"points"`
	OrderIndex   int       `json:"order_index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateQuestionRequest represents the data needed to create a new question
type CreateQuestionRequest struct {
	QuestionText  string   `json:"question_text" validate:"required,min=1,max=500"`
	QuestionType  string   `json:"question_type" validate:"required,oneof=multiple_choice true_false fill_in_blank"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer string   `json:"correct_answer" validate:"required,min=1"`
	Explanation   *string  `json:"explanation,omitempty" validate:"omitempty,max=1000"`
	Points        int      `json:"points" validate:"required,min=1,max=100"`
	OrderIndex    int      `json:"order_index" validate:"required,min=1"`
}

// QuizAttempt represents a user's attempt at a quiz
type QuizAttempt struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	QuizID      int64     `json:"quiz_id"`
	Score       int       `json:"score"`
	MaxScore    int       `json:"max_score"`
	TimeTaken   *int      `json:"time_taken,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserAnswer represents a user's answer to a specific question
type UserAnswer struct {
	ID           int64     `json:"id"`
	AttemptID    int64     `json:"attempt_id"`
	QuestionID   int64     `json:"question_id"`
	UserAnswer   string    `json:"user_answer"`
	IsCorrect    bool      `json:"is_correct"`
	PointsEarned int       `json:"points_earned"`
	AnsweredAt   time.Time `json:"answered_at"`
}

// ListResponse represents a paginated list response
type ListResponse[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}