package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"genshin-quiz-backend/internal/config"
	"genshin-quiz-backend/internal/handlers"
	"genshin-quiz-backend/internal/infrastructure"
	"genshin-quiz-backend/internal/middleware/auth"
	"genshin-quiz-backend/internal/middleware/logging"
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

	logger.Info("Server starting",
		zap.String("version", version),
		zap.String("buildTime", buildTime),
		zap.String("environment", cfg.Environment),
	)

	// Initialize JWT auth
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	// Initialize repositories
	userRepo := repository.NewUserRepository(infra.DB, logger)
	quizRepo := repository.NewQuizRepository(infra.DB, logger)

	// Initialize services
	userService := services.NewUserService(userRepo, infra.TaskClient, logger)
	quizService := services.NewQuizService(quizRepo, infra.TaskClient, logger)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, logger)
	quizHandler := handlers.NewQuizHandler(quizService, logger)

	// Setup router
	r := setupRouter(cfg, logger, tokenAuth, userHandler, quizHandler)

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Server forced to shutdown", zap.Error(err))
		}
	}()

	logger.Info("Server listening", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}

	logger.Info("Server stopped")
}

func setupRouter(cfg *config.Config, logger *zap.Logger, tokenAuth *jwtauth.JWTAuth, userHandler *handlers.UserHandler, quizHandler *handlers.QuizHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logging.Logger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"status": "ok",
			"version": "%s",
			"buildTime": "%s",
			"environment": "%s"
		}`, version, buildTime, cfg.Environment)))
	})

	// Public routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", userHandler.Login)
			r.Post("/register", userHandler.Register)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			// JWT authentication middleware
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(auth.Authenticator)

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.GetUsers)
				r.Get("/me", userHandler.GetCurrentUser)
				r.Put("/me", userHandler.UpdateCurrentUser)
				r.Get("/{id}", userHandler.GetUser)
				r.Put("/{id}", userHandler.UpdateUser)
				r.Delete("/{id}", userHandler.DeleteUser)
			})

			// Quiz routes
			r.Route("/quizzes", func(r chi.Router) {
				r.Get("/", quizHandler.GetQuizzes)
				r.Post("/", quizHandler.CreateQuiz)
				r.Get("/{id}", quizHandler.GetQuiz)
				r.Put("/{id}", quizHandler.UpdateQuiz)
				r.Delete("/{id}", quizHandler.DeleteQuiz)
				r.Post("/{id}/submit", quizHandler.SubmitQuiz)
			})
		})
	})

	return r
}