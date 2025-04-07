package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON writes a JSON response
func JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}
	return nil
}

// Error writes an error response
func Error(w http.ResponseWriter, status int, message string) {
	if err := JSON(w, status, Response{
		Success: false,
		Error:   message,
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Success writes a success response
func Success(w http.ResponseWriter, data interface{}, message string) {
	if err := JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
