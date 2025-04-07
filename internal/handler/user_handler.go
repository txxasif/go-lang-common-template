package handler

import (
	"encoding/json"
	"net/http"

	"myapp/internal/model"
	"myapp/internal/pkg/response"
	"myapp/internal/pkg/validation"
	"myapp/internal/service"
)

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
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "User registration details"
// @Success 201 {object} model.RegisterResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NewServiceError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body").Write(w)
		return
	}

	// Validate request
	if errs := validation.ValidateRegisterRequest(&req); errs.HasErrors() {
		response.NewValidationError(errs.Errors).Write(w)
		return
	}

	// Register user
	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		switch err {
		case model.ErrEmailAlreadyExists:
			response.NewServiceError(http.StatusConflict, "EMAIL_EXISTS", "User with this email already exists").Write(w)
		case model.ErrUsernameAlreadyExists:
			response.NewServiceError(http.StatusConflict, "USERNAME_EXISTS", "User with this username already exists").Write(w)
		default:
			response.NewServiceError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user").Write(w)
		}
		return
	}

	response.NewSuccess(http.StatusCreated, resp).Write(w)
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "User login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NewServiceError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body").Write(w)
		return
	}

	// Validate request
	if errs := validation.ValidateLoginRequest(&req); errs.HasErrors() {
		response.NewValidationError(errs.Errors).Write(w)
		return
	}

	// Login user
	resp, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case model.ErrInvalidCredentials:
			response.NewServiceError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password").Write(w)
		case model.ErrUserNotFound:
			response.NewServiceError(http.StatusUnauthorized, "USER_NOT_FOUND", "User not found").Write(w)
		default:
			response.NewServiceError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to login").Write(w)
		}
		return
	}

	response.NewSuccess(http.StatusOK, resp).Write(w)
}

// GetProfile handles getting the user's profile
// @Summary Get user profile
// @Description Get the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /profile [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		response.NewServiceError(http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header required").Write(w)
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Get user from token
	user, err := h.authService.GetUserByToken(r.Context(), token)
	if err != nil {
		response.NewServiceError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token").Write(w)
		return
	}

	response.NewSuccess(http.StatusOK, user.ToResponse()).Write(w)
}
