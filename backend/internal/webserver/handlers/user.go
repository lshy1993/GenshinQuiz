package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"genshin-quiz/internal/models"
	"genshin-quiz/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login authenticates a user and returns a JWT token
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Authenticate user
	user, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			h.renderError(w, r, http.StatusUnauthorized, "Invalid credentials", "")
			return
		}
		h.logger.Error("Authentication failed", 
			zap.String("email", req.Email),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Authentication failed", "")
		return
	}

	// Generate JWT token
	tokenAuth := jwtauth.FindTokenAuth(r.Context())
	if tokenAuth == nil {
		h.logger.Error("JWT auth not found in context")
		h.renderError(w, r, http.StatusInternalServerError, "Authentication configuration error", "")
		return
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		h.logger.Error("Failed to generate token", 
			zap.Int64("user_id", user.ID),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to generate token", "")
		return
	}

	response := AuthResponse{
		Token: tokenString,
		User:  user,
	}

	h.logger.Info("User logged in successfully", 
		zap.Int64("user_id", user.ID),
		zap.String("email", user.Email),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Register creates a new user account
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Create user
	createReq := models.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.userService.CreateUser(createReq)
	if err != nil {
		if err == models.ErrUserAlreadyExists {
			h.renderError(w, r, http.StatusConflict, "User already exists", "")
			return
		}
		h.logger.Error("Failed to create user", 
			zap.String("username", req.Username),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to create user", "")
		return
	}

	h.logger.Info("User registered successfully", 
		zap.Int64("user_id", user.ID),
		zap.String("username", user.Username),
	)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// GetUsers returns a paginated list of users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")

	limit := 10 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
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
	result, err := h.userService.GetUsers(limit, offset, search)
	if err != nil {
		h.logger.Error("Failed to get users", 
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.String("search", search),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to get users", "")
		return
	}

	render.JSON(w, r, result)
}

// GetCurrentUser returns the current authenticated user
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.GetUserFromContext(r.Context())
	if userClaims == nil {
		h.renderError(w, r, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	user, err := h.userService.GetUser(userClaims.UserID)
	if err != nil {
		if err == models.ErrUserNotFound {
			h.renderError(w, r, http.StatusNotFound, "User not found", "")
			return
		}
		h.logger.Error("Failed to get current user", 
			zap.Int64("user_id", userClaims.UserID),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to get user", "")
		return
	}

	render.JSON(w, r, user)
}

// UpdateCurrentUser updates the current authenticated user
func (h *UserHandler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userClaims := auth.GetUserFromContext(r.Context())
	if userClaims == nil {
		h.renderError(w, r, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(userClaims.UserID, req)
	if err != nil {
		if err == models.ErrUserNotFound {
			h.renderError(w, r, http.StatusNotFound, "User not found", "")
			return
		}
		h.logger.Error("Failed to update current user", 
			zap.Int64("user_id", userClaims.UserID),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to update user", "")
		return
	}

	render.JSON(w, r, user)
}

// GetUser returns a specific user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err == models.ErrUserNotFound {
			h.renderError(w, r, http.StatusNotFound, "User not found", "")
			return
		}
		h.logger.Error("Failed to get user", 
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to get user", "")
		return
	}

	render.JSON(w, r, user)
}

// UpdateUser updates a specific user (admin only)
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		if err == models.ErrUserNotFound {
			h.renderError(w, r, http.StatusNotFound, "User not found", "")
			return
		}
		h.logger.Error("Failed to update user", 
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to update user", "")
		return
	}

	render.JSON(w, r, user)
}

// DeleteUser deletes a specific user (admin only)
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, "Invalid user ID", "")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		if err == models.ErrUserNotFound {
			h.renderError(w, r, http.StatusNotFound, "User not found", "")
			return
		}
		h.logger.Error("Failed to delete user", 
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		h.renderError(w, r, http.StatusInternalServerError, "Failed to delete user", "")
		return
	}

	render.Status(r, http.StatusNoContent)
}

// renderError renders an error response
func (h *UserHandler) renderError(w http.ResponseWriter, r *http.Request, status int, message string, details string) {
	render.Status(r, status)
	errorResponse := map[string]interface{}{
		"error":   message,
		"status":  status,
	}
	
	if details != "" && h.isDevelopment() {
		errorResponse["details"] = details
	}
	
	render.JSON(w, r, errorResponse)
}

// isDevelopment checks if we're in development mode
func (h *UserHandler) isDevelopment() bool {
	// This could check an environment variable or config
	return true // Simplified for now
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