package service

import (
	"errors"
	"myapp/internal/config"
	"myapp/internal/model"
	"myapp/internal/repository"
	"myapp/pkg/hash"
	"myapp/pkg/jwt"
)

var (
	// ErrInvalidCredentials is returned when provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserAlreadyExists is returned when a user with the given email or username already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotFound is returned when a user doesn't exist
	ErrUserNotFound = errors.New("user not found")
)

// AuthService defines the interface for authentication-related operations
type AuthService interface {
	Register(req model.RegisterRequest) (*model.AuthResponse, error)
	Login(req model.LoginRequest) (*model.AuthResponse, error)
	GetUserByID(id uint) (*model.User, error)
	VerifyToken(token string) (uint, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

// Register registers a new user
func (s *authService) Register(req model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if user with the given email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Check if user with the given username already exists
	existingUser, err = s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &model.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, s.config.JWTSecret, s.config.JWTExpiresIn)
	if err != nil {
		return nil, err
	}

	// Return the response
	return &model.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *authService) Login(req model.LoginRequest) (*model.AuthResponse, error) {
	// Find the user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify the password
	if !hash.CheckPasswordHash(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, s.config.JWTSecret, s.config.JWTExpiresIn)
	if err != nil {
		return nil, err
	}

	// Return the response
	return &model.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// VerifyToken verifies a JWT token and returns the associated user ID
func (s *authService) VerifyToken(token string) (uint, error) {
	claims, err := jwt.VerifyToken(token, s.config.JWTSecret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
