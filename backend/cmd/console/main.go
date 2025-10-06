package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"genshin-quiz-backend/internal/config"
	"genshin-quiz-backend/internal/console"
	"genshin-quiz-backend/internal/infrastructure"
	"genshin-quiz-backend/internal/repository"
	"genshin-quiz-backend/internal/services"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Define command line flags
	var (
		command = flag.String("command", "", "Command to run (seed, migrate, user, quiz)")
		help    = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help || *command == "" {
		showHelp()
		return
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize infrastructure
	infra, err := infrastructure.NewInfrastructure(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize infrastructure: %v\n", err)
		os.Exit(1)
	}
	defer infra.Close()

	logger := infra.Logger

	logger.Info("Console application starting",
		zap.String("version", version),
		zap.String("buildTime", buildTime),
		zap.String("command", *command),
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(infra.DB, logger)
	quizRepo := repository.NewQuizRepository(infra.DB, logger)

	// Initialize services
	userService := services.NewUserService(userRepo, infra.TaskClient, logger)
	quizService := services.NewQuizService(quizRepo, infra.TaskClient, logger)

	// Initialize console commands
	cmdHandler := console.NewCommandHandler(userService, quizService, logger)

	// Execute command
	args := flag.Args()
	if err := cmdHandler.Execute(*command, args); err != nil {
		logger.Error("Command execution failed", 
			zap.String("command", *command),
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Info("Command executed successfully", zap.String("command", *command))
}

func showHelp() {
	fmt.Println("Genshin Quiz Console Application")
	fmt.Printf("Version: %s, Build Time: %s\n\n", version, buildTime)
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/console -command=<command> [args...]")
	fmt.Println("")
	fmt.Println("Available commands:")
	fmt.Println("  seed                    - Seed the database with sample data")
	fmt.Println("  user create <email>     - Create a new user")
	fmt.Println("  user list               - List all users")
	fmt.Println("  user delete <id>        - Delete a user")
	fmt.Println("  quiz create <title>     - Create a new quiz")
	fmt.Println("  quiz list               - List all quizzes")
	fmt.Println("  quiz delete <id>        - Delete a quiz")
	fmt.Println("  stats                   - Show system statistics")
	fmt.Println("  cleanup                 - Run data cleanup")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/console -command=seed")
	fmt.Println("  go run ./cmd/console -command=user create user@example.com")
	fmt.Println("  go run ./cmd/console -command=stats")
}