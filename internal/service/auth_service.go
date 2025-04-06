package service

import (
	"context"
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"strconv"
)

var (
	// ErrInvalidCredentials is returned when provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserAlreadyExists is returned when a user with the given email or username already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotFound is returned when a user doesn't exist
	ErrUserNotFound = errors.New("user not found")
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(ctx context.Context, email, password string) (*model.AuthResponse, error)
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error)
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	jwt      *JWTService
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      NewJWTService(jwtSecret),
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Create new user
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, email, password string) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// GetUser retrieves a user by ID
func (s *authService) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	userIDStr, err := s.jwt.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	user, err := s.userRepo.GetByID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
