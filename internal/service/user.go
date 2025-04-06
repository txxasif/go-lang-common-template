package service

import (
	"context"
	"myapp/internal/model"
	"myapp/internal/repository"
)

// UserService defines the interface for user operations
type UserService interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

// userService implements the UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, user *model.User) error {
	return s.userRepo.Create(ctx, user)
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail retrieves a user by email
func (s *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, user *model.User) error {
	return s.userRepo.Update(ctx, user)
}

// Delete deletes a user by ID
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
