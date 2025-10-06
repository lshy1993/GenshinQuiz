package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"genshin-quiz/internal/config"
	"genshin-quiz/internal/infrastructure"
	"genshin-quiz/internal/repository"
	"genshin-quiz/internal/services"
	"genshin-quiz/internal/tasks"
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

	logger.Info("Worker starting",
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

	// Initialize task processors
	taskProcessor := tasks.NewProcessor(userService, quizService, logger)

	// Create task multiplexer
	mux := asynq.NewServeMux()

	// Register task handlers
	mux.HandleFunc(tasks.TypeEmailVerification, taskProcessor.ProcessEmailVerification)
	mux.HandleFunc(tasks.TypeQuizSubmission, taskProcessor.ProcessQuizSubmission)
	mux.HandleFunc(tasks.TypeUserStatisticsUpdate, taskProcessor.ProcessUserStatisticsUpdate)
	mux.HandleFunc(tasks.TypeQuizAnalytics, taskProcessor.ProcessQuizAnalytics)
	mux.HandleFunc(tasks.TypeImageUpload, taskProcessor.ProcessImageUpload)

	// Handle shutdown gracefully
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the worker server
	go func() {
		if err := infra.TaskServer.Run(mux); err != nil {
			logger.Fatal("Worker server failed to start", zap.Error(err))
		}
	}()

	logger.Info("Worker server started",
		zap.Int("concurrency", cfg.Worker.Concurrency),
	)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Worker shutting down...")

	// Shutdown the worker server
	infra.TaskServer.Shutdown()

	logger.Info("Worker stopped")
}
