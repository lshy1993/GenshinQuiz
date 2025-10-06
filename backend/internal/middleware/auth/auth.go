package auth

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

// UserClaims represents the claims stored in JWT
type UserClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserContext key for storing user claims in context
type userContextKey struct{}

// Authenticator is a middleware that checks for a valid JWT token
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract user claims
		userClaims := UserClaims{
			UserID:   int64(claims["user_id"].(float64)),
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
		}

		// Add user claims to context
		ctx := context.WithValue(r.Context(), userContextKey{}, userClaims)

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly is a middleware that checks if the user has admin privileges
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims := GetUserFromContext(r.Context())
		if userClaims == nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Unauthorized"})
			return
		}

		// Check if user is admin (you might want to add admin role to UserClaims)
		// For now, we'll assume admin check based on user ID or other criteria
		// This is a placeholder - implement your actual admin logic

		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts user claims from the request context
func GetUserFromContext(ctx context.Context) *UserClaims {
	if claims, ok := ctx.Value(userContextKey{}).(UserClaims); ok {
		return &claims
	}
	return nil
}

// RequireUser is a middleware that ensures a user is authenticated
func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims := GetUserFromContext(r.Context())
		if userClaims == nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Authentication required"})
			return
		}

		next.ServeHTTP(w, r)
	})
}