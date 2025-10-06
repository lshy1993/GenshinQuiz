package error

import (
	"encoding/json"
	"net/http"
	
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// Handler creates an error handling middleware
func Handler(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						zap.String("method", r.Method),
						zap.String("url", r.URL.String()),
						zap.Any("error", err),
						zap.String("request_id", r.Header.Get("X-Request-ID")),
					)
					
					writeErrorResponse(w, http.StatusInternalServerError, "Internal server error", "", "")
				}
			}()
			
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// WriteErrorResponse writes a standardized error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message, code, details string) {
	writeErrorResponse(w, statusCode, message, code, details)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message, code, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}
	
	json.NewEncoder(w).Encode(response)
}
