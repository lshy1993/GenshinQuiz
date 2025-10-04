package services

import (
	"fmt"

	"genshin-quiz-backend/internal/models"
	"genshin-quiz-backend/internal/repository"
)

type QuizService struct {
	quizRepo *repository.QuizRepository
}

func NewQuizService(quizRepo *repository.QuizRepository) *QuizService {
	return &QuizService{
		quizRepo: quizRepo,
	}
}

func (s *QuizService) GetQuizzes(limit, offset int, category, difficulty string) (*models.ListResponse[models.Quiz], error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Validate category
	if category != "" && !isValidCategory(category) {
		return nil, fmt.Errorf("invalid category: %s", category)
	}

	// Validate difficulty
	if difficulty != "" && !isValidDifficulty(difficulty) {
		return nil, fmt.Errorf("invalid difficulty: %s", difficulty)
	}

	quizzes, total, err := s.quizRepo.GetAll(limit, offset, category, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to get quizzes: %w", err)
	}

	return &models.ListResponse[models.Quiz]{
		Data:   quizzes,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *QuizService) GetQuiz(id int64) (*models.Quiz, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid quiz ID")
	}

	quiz, err := s.quizRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}

	if quiz == nil {
		return nil, fmt.Errorf("quiz not found")
	}

	return quiz, nil
}

func (s *QuizService) CreateQuiz(req models.CreateQuizRequest) (*models.Quiz, error) {
	// Validate request
	if err := s.validateCreateQuizRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	quiz, err := s.quizRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create quiz: %w", err)
	}

	return quiz, nil
}

func (s *QuizService) UpdateQuiz(id int64, req models.UpdateQuizRequest) (*models.Quiz, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid quiz ID")
	}

	// Validate request
	if err := s.validateUpdateQuizRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if quiz exists
	existing, err := s.quizRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("quiz not found")
	}

	quiz, err := s.quizRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update quiz: %w", err)
	}

	return quiz, nil
}

func (s *QuizService) DeleteQuiz(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid quiz ID")
	}

	err := s.quizRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete quiz: %w", err)
	}

	return nil
}

func (s *QuizService) validateCreateQuizRequest(req models.CreateQuizRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(req.Title) > 200 {
		return fmt.Errorf("title must be 200 characters or less")
	}
	if !isValidCategory(req.Category) {
		return fmt.Errorf("invalid category: %s", req.Category)
	}
	if !isValidDifficulty(req.Difficulty) {
		return fmt.Errorf("invalid difficulty: %s", req.Difficulty)
	}
	if len(req.Questions) == 0 {
		return fmt.Errorf("at least one question is required")
	}
	if req.CreatedBy <= 0 {
		return fmt.Errorf("invalid creator ID")
	}

	// Validate questions
	for i, q := range req.Questions {
		if err := s.validateQuestionRequest(q, i+1); err != nil {
			return fmt.Errorf("question %d: %w", i+1, err)
		}
	}

	return nil
}

func (s *QuizService) validateUpdateQuizRequest(req models.UpdateQuizRequest) error {
	if req.Title != nil && len(*req.Title) > 200 {
		return fmt.Errorf("title must be 200 characters or less")
	}
	if req.Category != nil && !isValidCategory(*req.Category) {
		return fmt.Errorf("invalid category: %s", *req.Category)
	}
	if req.Difficulty != nil && !isValidDifficulty(*req.Difficulty) {
		return fmt.Errorf("invalid difficulty: %s", *req.Difficulty)
	}

	// Validate questions if provided
	if req.Questions != nil {
		if len(req.Questions) == 0 {
			return fmt.Errorf("at least one question is required")
		}
		for i, q := range req.Questions {
			if err := s.validateQuestionRequest(q, i+1); err != nil {
				return fmt.Errorf("question %d: %w", i+1, err)
			}
		}
	}

	return nil
}

func (s *QuizService) validateQuestionRequest(req models.CreateQuestionRequest, index int) error {
	if req.QuestionText == "" {
		return fmt.Errorf("question text is required")
	}
	if len(req.QuestionText) > 500 {
		return fmt.Errorf("question text must be 500 characters or less")
	}
	if !isValidQuestionType(req.QuestionType) {
		return fmt.Errorf("invalid question type: %s", req.QuestionType)
	}
	if req.CorrectAnswer == "" {
		return fmt.Errorf("correct answer is required")
	}
	if req.Points <= 0 || req.Points > 100 {
		return fmt.Errorf("points must be between 1 and 100")
	}
	if req.OrderIndex <= 0 {
		return fmt.Errorf("order index must be positive")
	}

	// Validate options for multiple choice questions
	if req.QuestionType == "multiple_choice" {
		if len(req.Options) < 2 {
			return fmt.Errorf("multiple choice questions must have at least 2 options")
		}
		// Check if correct answer is in options
		found := false
		for _, option := range req.Options {
			if option == req.CorrectAnswer {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("correct answer must be one of the options")
		}
	}

	return nil
}

func isValidCategory(category string) bool {
	validCategories := []string{"characters", "weapons", "artifacts", "lore", "gameplay"}
	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}

func isValidDifficulty(difficulty string) bool {
	validDifficulties := []string{"easy", "medium", "hard"}
	for _, valid := range validDifficulties {
		if difficulty == valid {
			return true
		}
	}
	return false
}

func isValidQuestionType(questionType string) bool {
	validTypes := []string{"multiple_choice", "true_false", "fill_in_blank"}
	for _, valid := range validTypes {
		if questionType == valid {
			return true
		}
	}
	return false
}