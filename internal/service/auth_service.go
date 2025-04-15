package service

import (
	"context"
	"errors"
	"myapp/internal/model"
	"myapp/internal/pkg/jwt"
	"myapp/internal/repository"
)

var (
	// ErrInvalidCredentials is returned when provided credentials are invalid
	ErrInvalidCredentials = model.ErrInvalidCredentials
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrUnauthorized is returned when a user is not authorized
	ErrUnauthorized = errors.New("unauthorized")
)

// AuthResponse contains the response after authentication
type AuthResponse struct {
	User         model.UserResponse `json:"user"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
}

// RefreshResponse contains the response after refreshing tokens
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(ctx context.Context, email, password string) (*AuthResponse, error)
	Register(ctx context.Context, req *model.RegisterRequest) (*AuthResponse, error)
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*RefreshResponse, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	jwt      *jwt.Service
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository, jwtService *jwt.Service) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwtService,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, model.ErrEmailAlreadyExists
	}

	// Check if username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, model.ErrUsernameAlreadyExists
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

	// Generate tokens
	accessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, model.ErrUserNotFound
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetUserByToken retrieves a user from a token
func (s *authService) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	userID, tokenType, err := s.jwt.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Only access tokens should be used for authentication
	if tokenType != jwt.AccessToken {
		return nil, ErrUnauthorized
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// RefreshTokens refreshes the access and refresh tokens
func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*RefreshResponse, error) {
	userID, tokenType, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Verify it's a refresh token
	if tokenType != jwt.RefreshToken {
		return nil, ErrUnauthorized
	}

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUnauthorized
	}

	// Generate new tokens
	newAccessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
