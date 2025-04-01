package middleware

import (
	"context"
	"myapp/internal/service"
	"net/http"
	"strings"
)

// contextKey is a custom type for context keys
type contextKey string

const userIDKey contextKey = "userID"

// Auth is a middleware that checks for a valid JWT token in the Authorization header
func Auth(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if the Authorization header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			// Get the token
			token := parts[1]
			if token == "" {
				http.Error(w, "Token required", http.StatusUnauthorized)
				return
			}

			// Verify the token
			userID, err := authService.VerifyToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add the user ID to the request context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
