package handler

import (
	"encoding/json"
	"myapp/internal/service"
	"net/http"
)

// Handler holds all HTTP handlers
type Handler struct {
	Auth *AuthHandler
}

// HandlerOption defines function signature for handler options
type HandlerOption func(*Handler)

// WithAuthHandler sets the auth handler
func WithAuthHandler(authService service.AuthService) HandlerOption {
	return func(h *Handler) {
		h.Auth = NewAuthHandler(authService)
	}
}

// New creates a new Handler with options
func New(opts ...HandlerOption) *Handler {
	h := &Handler{}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Internal Server Error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
