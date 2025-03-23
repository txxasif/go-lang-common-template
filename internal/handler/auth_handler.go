package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"myapp/internal/model"
	"myapp/internal/service"

	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		respondWithError(w, http.StatusBadRequest, formatValidationErrors(validationErrors))
		return
	}

	// Register the user
	resp, err := h.authService.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			respondWithError(w, http.StatusConflict, "User with this email or username already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		respondWithError(w, http.StatusBadRequest, formatValidationErrors(validationErrors))
		return
	}

	// Login the user
	resp, err := h.authService.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to log in")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// GetUser retrieves the current authenticated user's information
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user info
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to get user information")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user.ToResponse())
}

// formatValidationErrors formats validation errors
func formatValidationErrors(errs validator.ValidationErrors) string {
	// For simplicity, just return a basic message
	// In a real application, you might want to format these errors more nicely
	return "Validation failed: Please check your input and try again"
}
