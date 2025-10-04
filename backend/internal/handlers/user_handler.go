package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"genshin-quiz-backend/internal/models"
	"genshin-quiz-backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	// Get users
	result, err := h.userService.GetUsers(limit, offset)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"users":  result.Data,
		"total":  result.Total,
		"limit":  result.Limit,
		"offset": result.Offset,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID", "User ID must be a valid integer")
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			h.renderError(w, r, http.StatusNotFound, "User not found", "The requested user could not be found")
			return
		}
		h.renderError(w, r, http.StatusInternalServerError, "Failed to get user", err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", "Failed to parse JSON request body")
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			h.renderError(w, r, http.StatusConflict, "Conflict", err.Error())
			return
		}
		h.renderError(w, r, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID", "User ID must be a valid integer")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", "Failed to parse JSON request body")
		return
	}

	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		if err.Error() == "user not found" {
			h.renderError(w, r, http.StatusNotFound, "User not found", "The requested user could not be found")
			return
		}
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			h.renderError(w, r, http.StatusConflict, "Conflict", err.Error())
			return
		}
		h.renderError(w, r, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID", "User ID must be a valid integer")
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			h.renderError(w, r, http.StatusNotFound, "User not found", "The requested user could not be found")
			return
		}
		h.renderError(w, r, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) renderError(w http.ResponseWriter, r *http.Request, status int, error, message string) {
	render.Status(r, status)
	render.JSON(w, r, models.ErrorResponse{
		Error:   error,
		Message: message,
		Code:    status,
	})
}