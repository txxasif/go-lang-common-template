package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDKey contextKey = "requestID"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a new request ID
			requestID := uuid.New().String()

			// Add the request ID to the response header
			w.Header().Set("X-Request-ID", requestID)

			// Add the request ID to the request context
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)

			// Create a logger with request ID
			requestLogger := logger.With(
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
			)

			// Log the request start
			requestLogger.Info("Request started")

			// Create a response writer that tracks the status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Start timing the request
			start := time.Now()

			// Call the next handler
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Log the request completion
			duration := time.Since(start)
			requestLogger.Info("Request completed",
				zap.Int("status_code", rw.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}

// GetRequestID returns the request ID from the context
func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// responseWriter is a wrapper around http.ResponseWriter that tracks the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
