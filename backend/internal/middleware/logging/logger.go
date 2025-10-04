package logging

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// Logger creates a request logging middleware
func Logger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Wrap the response writer to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			// Process request
			next.ServeHTTP(ww, r)
			
			// Log request details
			duration := time.Since(start)
			
			logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"url":        r.URL.String(),
				"status":     ww.Status(),
				"bytes":      ww.BytesWritten(),
				"duration":   duration.String(),
				"remote_ip":  r.RemoteAddr,
				"user_agent": r.UserAgent(),
				"request_id": middleware.GetReqID(r.Context()),
			}).Info("Request processed")
		}
		return http.HandlerFunc(fn)
	}
}