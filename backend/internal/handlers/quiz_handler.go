package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/getsentry/sentry-go"

	"genshin-quiz-backend/internal/models"
	"genshin-quiz-backend/internal/services"
)

type QuizHandler struct {
	quizService *services.QuizService
	logger      *log.Logger
}

func NewQuizHandler(quizService *services.QuizService, logger *log.Logger) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		logger:      logger,
	}
}

func (h *QuizHandler) GetQuizzes(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	category := r.URL.Query().Get("category")
	difficulty := r.URL.Query().Get("difficulty")

	limit := 10 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Get quizzes
	result, err := h.quizService.GetQuizzes(limit, offset, category, difficulty)
	if err != nil {
		h.logger.Printf("failed to get quizzes: %v", err)
		sentry.CaptureException(err)
		h.renderError(w, r, http.StatusBadRequest, "Failed to get quizzes", err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"quizzes": result.Data,
		"total":   result.Total,
		"limit":   result.Limit,
		"offset":  result.Offset,
	})
}

func (h *QuizHandler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid quiz ID", "Quiz ID must be a valid integer")
		return
	}

	quiz, err := h.quizService.GetQuiz(id)
	if err != nil {
		if err.Error() == "quiz not found" {
			h.renderError(w, r, http.StatusNotFound, "Quiz not found", "The requested quiz could not be found")
			return
		}
		h.renderError(w, r, http.StatusInternalServerError, "Failed to get quiz", err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, quiz)
}

func (h *QuizHandler) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var req models.CreateQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", "Failed to parse JSON request body")
		return
	}

	quiz, err := h.quizService.CreateQuiz(req)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Failed to create quiz", err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, quiz)
}

func (h *QuizHandler) UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid quiz ID", "Quiz ID must be a valid integer")
		return
	}

	var req models.UpdateQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", "Failed to parse JSON request body")
		return
	}

	quiz, err := h.quizService.UpdateQuiz(id, req)
	if err != nil {
		if err.Error() == "quiz not found" {
			h.renderError(w, r, http.StatusNotFound, "Quiz not found", "The requested quiz could not be found")
			return
		}
		h.renderError(w, r, http.StatusBadRequest, "Failed to update quiz", err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, quiz)
}

func (h *QuizHandler) DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid quiz ID", "Quiz ID must be a valid integer")
		return
	}

	err = h.quizService.DeleteQuiz(id)
	if err != nil {
		if err.Error() == "quiz not found" {
			h.renderError(w, r, http.StatusNotFound, "Quiz not found", "The requested quiz could not be found")
			return
		}
		h.renderError(w, r, http.StatusInternalServerError, "Failed to delete quiz", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuizHandler) renderError(w http.ResponseWriter, r *http.Request, status int, error, message string) {
	render.Status(r, status)
	render.JSON(w, r, models.ErrorResponse{
		Error:   error,
		Message: message,
		Code:    status,
	})
}