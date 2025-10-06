package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"genshin-quiz-backend/internal/config"
	"genshin-quiz-backend/internal/cron"
	"genshin-quiz-backend/internal/infrastructure"
	"genshin-quiz-backend/internal/repository"
	"genshin-quiz-backend/internal/services"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
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

	logger.Info("Cronjob service starting",
		zap.String("version", version),
		zap.String("buildTime", buildTime),
		zap.String("environment", cfg.Environment),
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(infra.DB, logger)
	quizRepo := repository.NewQuizRepository(infra.DB, logger)

	// Initialize services
	userService := services.NewUserService(userRepo, infra.TaskClient, logger)
	quizService := services.NewQuizService(quizRepo, infra.TaskClient, logger)

	// Initialize cron scheduler
	scheduler := cron.NewScheduler(userService, quizService, infra.TaskClient, logger)

	// Start the scheduler
	scheduler.Start()

	logger.Info("Cronjob service started")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Cronjob service shutting down...")

	// Stop the scheduler
	scheduler.Stop()

	logger.Info("Cronjob service stopped")
}