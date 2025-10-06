package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"genshin-quiz-backend/internal/config"
	"genshin-quiz-backend/internal/infrastructure"
	"genshin-quiz-backend/internal/migration"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Define command line flags
	var (
		action = flag.String("action", "", "Migration action (up, down, create, status)")
		name   = flag.String("name", "", "Migration name for create action")
		help   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help || *action == "" {
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

	logger.Info("Migration service starting",
		zap.String("version", version),
		zap.String("buildTime", buildTime),
		zap.String("action", *action),
	)

	// Initialize migration manager
	migrator := migration.NewMigrator(infra.DB, logger)

	// Execute migration action
	switch *action {
	case "up":
		if err := migrator.Up(); err != nil {
			logger.Fatal("Migration up failed", zap.Error(err))
		}
		logger.Info("Migrations applied successfully")

	case "down":
		if err := migrator.Down(); err != nil {
			logger.Fatal("Migration down failed", zap.Error(err))
		}
		logger.Info("Migrations rolled back successfully")

	case "create":
		if *name == "" {
			logger.Fatal("Migration name is required for create action")
		}
		if err := migrator.Create(*name); err != nil {
			logger.Fatal("Migration creation failed", zap.Error(err))
		}
		logger.Info("Migration file created successfully", zap.String("name", *name))

	case "status":
		status, err := migrator.Status()
		if err != nil {
			logger.Fatal("Failed to get migration status", zap.Error(err))
		}
		fmt.Println(status)

	default:
		logger.Fatal("Unknown migration action", zap.String("action", *action))
	}

	logger.Info("Migration operation completed")
}

func showHelp() {
	fmt.Println("Genshin Quiz Migration Tool")
	fmt.Printf("Version: %s, Build Time: %s\n\n", version, buildTime)
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/migration -action=<action> [options]")
	fmt.Println("")
	fmt.Println("Actions:")
	fmt.Println("  up                      - Apply all pending migrations")
	fmt.Println("  down                    - Rollback the last migration")
	fmt.Println("  create -name=<name>     - Create a new migration file")
	fmt.Println("  status                  - Show migration status")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/migration -action=up")
	fmt.Println("  go run ./cmd/migration -action=down")
	fmt.Println("  go run ./cmd/migration -action=create -name=add_user_profile")
	fmt.Println("  go run ./cmd/migration -action=status")
}