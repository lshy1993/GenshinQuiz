package middleware

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"genshin-quiz/config"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

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

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// If JSON encoding fails, write a simple error message
		// We intentionally ignore the error from Write as there's nothing more we can do
		w.Write([]byte(`{"error":"Internal server error"}`)) //nolint:errcheck // fallback error writing
	}
}

func HandleBadRequestError(app *config.App) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		app.Logger.Error("Bad request error",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Error(err),
			zap.String("request_id", r.Header.Get("X-Request-ID")),
		)

		writeErrorResponse(w, http.StatusBadRequest, "Bad request", "INVALID_REQUEST", err.Error())
	}
}

func HandleResponseErrorWithLog(app *config.App) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		app.Logger.Error("Response error",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Error(err),
			zap.String("request_id", r.Header.Get("X-Request-ID")),
		)

		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error", "INTERNAL_ERROR", err.Error())
	}
}
