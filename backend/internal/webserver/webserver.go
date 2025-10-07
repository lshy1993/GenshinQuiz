package webserver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"

	"genshin-quiz/config"
	"genshin-quiz/generated/oapi"
	"genshin-quiz/internal/webserver/handler"
	mw "genshin-quiz/internal/webserver/middleware"
)

type Server struct {
	router     *chi.Mux
	serverAddr string
}

func NewServer(app *config.App) *Server {
	// Setup router
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logger(app.Logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Error handling middleware
	r.Use(mw.Handler(app.Logger))

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
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
		if err != nil {
			app.Logger.Error("Failed to write health check response", zap.Error(err))
		}
	})

	// Setup routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.JWTAuth))
		r.Use(mw.Authenticator)

		baseURL := ""
		serverOptions := oapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  mw.HandleBadRequestError(app),
			ResponseErrorHandlerFunc: mw.HandleResponseErrorWithLog(app),
		}
		strictHandler := oapi.NewStrictHandlerWithOptions(
			handler.NewHandler(app),
			[]oapi.StrictMiddlewareFunc{},
			serverOptions,
		)
		oapi.HandlerFromMuxWithBaseURL(strictHandler, r, baseURL)
	})

	return &Server{
		router:     r,
		serverAddr: fmt.Sprintf("%s:%s", app.Server.Host, app.Server.Port),
	}
}

func (s *Server) Start() {
	log.Print("Starting server", zap.String("addr", s.serverAddr))

	const maxHeaderBytes = 1 << 20
	const readTimeout = 10 * time.Second
	const writeTimeout = 30 * time.Second
	const idleTimeout = 10 * time.Second

	// Create HTTP server
	srv := &http.Server{
		Addr:           s.serverAddr,
		Handler:        s.router,
		ReadTimeout:    readTimeout,
		MaxHeaderBytes: maxHeaderBytes,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) Router() chi.Router {
	return s.router
}
