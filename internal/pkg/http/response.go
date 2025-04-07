package http

import (
	"encoding/json"
	"net/http"

	"myapp/internal/pkg/validation"
)

// Error sends an error response
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

// ValidationError sends a validation error response
func ValidationError(w http.ResponseWriter, errors *validation.ValidationErrors) {
	JSON(w, http.StatusBadRequest, map[string]interface{}{
		"error":  "Validation failed",
		"errors": errors.Errors,
	})
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
