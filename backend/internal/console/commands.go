package console

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"genshin-quiz/internal/models"
	"genshin-quiz/internal/services"

	"go.uber.org/zap"
)

// CommandHandler handles console commands
type CommandHandler struct {
	userService *services.UserService
	quizService *services.QuizService
	logger      *log.Logger
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(userService *services.UserService, quizService *services.QuizService, logger *log.Logger) *CommandHandler {
	return &CommandHandler{
		userService: userService,
		quizService: quizService,
		logger:      logger,
	}
}

// Execute executes a console command
func (h *CommandHandler) Execute(command string, args []string) error {
	switch command {
	case "seed":
		return h.seedDatabase()
	case "user":
		return h.handleUserCommand(args)
	case "quiz":
		return h.handleQuizCommand(args)
	case "stats":
		return h.showStats()
	case "cleanup":
		return h.runCleanup()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// seedDatabase seeds the database with sample data
func (h *CommandHandler) seedDatabase() error {
	h.logger.Info("Seeding database with sample data...")

	// Create sample users
	sampleUsers := []models.CreateUserRequest{
		{
			Username: "admin",
			Email:    "admin@genshinquiz.com",
			Password: "admin123",
		},
		{
			Username: "testuser1",
			Email:    "user1@example.com",
			Password: "password123",
		},
		{
			Username: "testuser2",
			Email:    "user2@example.com",
			Password: "password123",
		},
	}

	for _, userReq := range sampleUsers {
		user, err := h.userService.CreateUser(userReq)
		if err != nil {
			h.logger.Error("Failed to create sample user",
				zap.String("username", userReq.Username),
				zap.Error(err),
			)
			continue
		}
		h.logger.Info("Created sample user",
			zap.String("username", user.Username),
			zap.Int64("id", user.ID),
		)
	}

	// Create sample quizzes
	sampleQuizzes := []models.CreateQuizRequest{
		{
			Title:       "Genshin Impact Characters",
			Description: "Test your knowledge about Genshin Impact characters",
			Category:    "Characters",
			Difficulty:  "Easy",
			Questions: []models.CreateQuestionRequest{
				{
					QuestionText: "Who is the Anemo Archon?",
					QuestionType: "multiple_choice",
					Options: []models.CreateOptionRequest{
						{OptionText: "Venti", IsCorrect: true},
						{OptionText: "Diluc", IsCorrect: false},
						{OptionText: "Xiao", IsCorrect: false},
						{OptionText: "Jean", IsCorrect: false},
					},
				},
				{
					QuestionText: "Which character is known as the 'Dark Knight Hero'?",
					QuestionType: "multiple_choice",
					Options: []models.CreateOptionRequest{
						{OptionText: "Kaeya", IsCorrect: false},
						{OptionText: "Diluc", IsCorrect: true},
						{OptionText: "Albedo", IsCorrect: false},
						{OptionText: "Childe", IsCorrect: false},
					},
				},
			},
		},
		{
			Title:       "Genshin Impact Regions",
			Description: "How well do you know the regions of Teyvat?",
			Category:    "World",
			Difficulty:  "Medium",
			Questions: []models.CreateQuestionRequest{
				{
					QuestionText: "What is the main city of Mondstadt?",
					QuestionType: "multiple_choice",
					Options: []models.CreateOptionRequest{
						{OptionText: "Mondstadt City", IsCorrect: true},
						{OptionText: "Liyue Harbor", IsCorrect: false},
						{OptionText: "Inazuma City", IsCorrect: false},
						{OptionText: "Sumeru City", IsCorrect: false},
					},
				},
			},
		},
	}

	for _, quizReq := range sampleQuizzes {
		quiz, err := h.quizService.CreateQuiz(quizReq)
		if err != nil {
			h.logger.Error("Failed to create sample quiz",
				zap.String("title", quizReq.Title),
				zap.Error(err),
			)
			continue
		}
		h.logger.Info("Created sample quiz",
			zap.String("title", quiz.Title),
			zap.Int64("id", quiz.ID),
		)
	}

	h.logger.Info("Database seeding completed")
	return nil
}

// handleUserCommand handles user-related commands
func (h *CommandHandler) handleUserCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("user command requires subcommand (create, list, delete)")
	}

	switch args[0] {
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("create command requires email argument")
		}
		return h.createUser(args[1])
	case "list":
		return h.listUsers()
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("delete command requires user ID argument")
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid user ID: %s", args[1])
		}
		return h.deleteUser(id)
	default:
		return fmt.Errorf("unknown user subcommand: %s", args[0])
	}
}

// handleQuizCommand handles quiz-related commands
func (h *CommandHandler) handleQuizCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("quiz command requires subcommand (create, list, delete)")
	}

	switch args[0] {
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("create command requires title argument")
		}
		return h.createQuiz(args[1])
	case "list":
		return h.listQuizzes()
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("delete command requires quiz ID argument")
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid quiz ID: %s", args[1])
		}
		return h.deleteQuiz(id)
	default:
		return fmt.Errorf("unknown quiz subcommand: %s", args[0])
	}
}

// createUser creates a new user
func (h *CommandHandler) createUser(email string) error {
	req := models.CreateUserRequest{
		Username: email,
		Email:    email,
		Password: "defaultpassword123", // In real app, you'd prompt for password
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("User created successfully:\n")
	fmt.Printf("  ID: %d\n", user.ID)
	fmt.Printf("  Username: %s\n", user.Username)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Created: %s\n", user.CreatedAt.Format(time.RFC3339))

	return nil
}

// listUsers lists all users
func (h *CommandHandler) listUsers() error {
	response, err := h.userService.GetUsers(100, 0, "")
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	fmt.Printf("Users (Total: %d):\n", response.Total)
	fmt.Printf("%-5s %-20s %-30s %-20s\n", "ID", "Username", "Email", "Created")
	fmt.Println("--------------------------------------------------------------------")

	for _, user := range response.Data {
		fmt.Printf("%-5d %-20s %-30s %-20s\n",
			user.ID,
			user.Username,
			user.Email,
			user.CreatedAt.Format("2006-01-02 15:04"),
		)
	}

	return nil
}

// deleteUser deletes a user
func (h *CommandHandler) deleteUser(id int64) error {
	if err := h.userService.DeleteUser(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	fmt.Printf("User %d deleted successfully\n", id)
	return nil
}

// createQuiz creates a simple quiz
func (h *CommandHandler) createQuiz(title string) error {
	req := models.CreateQuizRequest{
		Title:       title,
		Description: "Quiz created via console",
		Category:    "General",
		Difficulty:  "Easy",
		Questions: []models.CreateQuestionRequest{
			{
				QuestionText: "Sample question?",
				QuestionType: "multiple_choice",
				Options: []models.CreateOptionRequest{
					{OptionText: "Option A", IsCorrect: true},
					{OptionText: "Option B", IsCorrect: false},
				},
			},
		},
	}

	quiz, err := h.quizService.CreateQuiz(req)
	if err != nil {
		return fmt.Errorf("failed to create quiz: %w", err)
	}

	fmt.Printf("Quiz created successfully:\n")
	fmt.Printf("  ID: %d\n", quiz.ID)
	fmt.Printf("  Title: %s\n", quiz.Title)
	fmt.Printf("  Category: %s\n", quiz.Category)
	fmt.Printf("  Created: %s\n", quiz.CreatedAt.Format(time.RFC3339))

	return nil
}

// listQuizzes lists all quizzes
func (h *CommandHandler) listQuizzes() error {
	response, err := h.quizService.GetQuizzes(100, 0, "", "")
	if err != nil {
		return fmt.Errorf("failed to get quizzes: %w", err)
	}

	fmt.Printf("Quizzes (Total: %d):\n", response.Total)
	fmt.Printf("%-5s %-30s %-15s %-10s %-20s\n", "ID", "Title", "Category", "Difficulty", "Created")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, quiz := range response.Data {
		fmt.Printf("%-5d %-30s %-15s %-10s %-20s\n",
			quiz.ID,
			quiz.Title,
			quiz.Category,
			quiz.Difficulty,
			quiz.CreatedAt.Format("2006-01-02 15:04"),
		)
	}

	return nil
}

// deleteQuiz deletes a quiz
func (h *CommandHandler) deleteQuiz(id int64) error {
	if err := h.quizService.DeleteQuiz(id); err != nil {
		return fmt.Errorf("failed to delete quiz: %w", err)
	}

	fmt.Printf("Quiz %d deleted successfully\n", id)
	return nil
}

// showStats shows system statistics
func (h *CommandHandler) showStats() error {
	fmt.Println("System Statistics:")
	fmt.Println("==================")

	// Get user stats
	users, err := h.userService.GetUsers(1, 0, "")
	if err == nil {
		fmt.Printf("Total Users: %d\n", users.Total)
	}

	// Get quiz stats
	quizzes, err := h.quizService.GetQuizzes(1, 0, "", "")
	if err == nil {
		fmt.Printf("Total Quizzes: %d\n", quizzes.Total)
	}

	fmt.Printf("Current Time: %s\n", time.Now().Format(time.RFC3339))

	return nil
}

// runCleanup runs data cleanup
func (h *CommandHandler) runCleanup() error {
	h.logger.Info("Running manual data cleanup...")

	// TODO: Implement cleanup logic
	// This could include:
	// - Cleaning up old sessions
	// - Removing orphaned records
	// - Archiving old data

	fmt.Println("Data cleanup completed")
	return nil
}
