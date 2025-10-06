package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"genshin-quiz/config"
	"genshin-quiz/internal/database"
	"genshin-quiz/internal/repository"
	"genshin-quiz/internal/services"
	"genshin-quiz/internal/webserver/handlers"
	mw "genshin-quiz/internal/webserver/middleware"
)

type Server struct {
	config *config.Config
	logger *zap.Logger
	db     *database.DB
	router chi.Router
	server *http.Server
}

type Dependencies struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *database.DB
}

func New(deps Dependencies) *Server {
	return &Server{
		config: deps.Config,
		logger: deps.Logger,
		db:     deps.DB,
	}
}

func (s *Server) Initialize() error {
	// Setup router
	s.router = chi.NewRouter()

	// Setup middleware
	s.setupMiddleware()

	// Setup routes
	if err := s.setupRoutes(); err != nil {
		return fmt.Errorf("failed to setup routes: %w", err)
	}

	// Create HTTP server
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	return nil
}

func (s *Server) setupMiddleware() {
	// Basic middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(mw.Logger(s.logger))
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// Error handling middleware
	s.router.Use(mw.Handler(s.logger))

	// CORS configuration
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (s *Server) setupRoutes() error {
	// Health check endpoint
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// Initialize dependencies for handlers
	deps, err := s.initializeDependencies()
	if err != nil {
		return fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", deps.UserHandler.Login)
			r.Post("/auth/register", deps.UserHandler.Register)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.JWTAuth(s.config.JWTSecret))

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", deps.UserHandler.GetUsers)
				r.Post("/", deps.UserHandler.CreateUser)
				r.Get("/{id}", deps.UserHandler.GetUser)
				r.Put("/{id}", deps.UserHandler.UpdateUser)
				r.Delete("/{id}", deps.UserHandler.DeleteUser)
			})

			// Quiz routes
			r.Route("/quizzes", func(r chi.Router) {
				r.Get("/", deps.QuizHandler.GetQuizzes)
				r.Post("/", deps.QuizHandler.CreateQuiz)
				r.Get("/{id}", deps.QuizHandler.GetQuiz)
				r.Put("/{id}", deps.QuizHandler.UpdateQuiz)
				r.Delete("/{id}", deps.QuizHandler.DeleteQuiz)
			})
		})
	})

	return nil
}

type HandlerDependencies struct {
	UserHandler *handlers.UserHandler
	QuizHandler *handlers.QuizHandler
}

func (s *Server) initializeDependencies() (*HandlerDependencies, error) {
	// Initialize repositories (no logger needed)
	userRepo := repository.NewUserRepository(s.db.GetJetDB())
	quizRepo := repository.NewQuizRepository(s.db.GetJetDB())

	// Initialize services (no logger needed)
	userService := services.NewUserService(userRepo)
	quizService := services.NewQuizService(quizRepo)

	// Initialize handlers (no logger needed)
	userHandler := handlers.NewUserHandler(userService)
	quizHandler := handlers.NewQuizHandler(quizService)

	return &HandlerDependencies{
		UserHandler: userHandler,
		QuizHandler: quizHandler,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", zap.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}

func (s *Server) Router() chi.Router {
	return s.router
}
