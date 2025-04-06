package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"myapp/internal/model"
	httputil "myapp/internal/pkg/http"
	"myapp/internal/service"
)

var validate = validator.New()

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	authService service.AuthService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(authService service.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Register user
	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			httputil.Error(w, http.StatusConflict, "User already exists")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	httputil.JSON(w, http.StatusCreated, resp)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Login user
	resp, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			httputil.Error(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}

// GetProfile handles getting the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		httputil.Error(w, http.StatusUnauthorized, "Authorization header required")
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Get user from token
	user, err := h.authService.GetUserByToken(r.Context(), token)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	httputil.JSON(w, http.StatusOK, user.ToResponse())
}
