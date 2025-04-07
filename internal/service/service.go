package service

import (
	"myapp/internal/repository"
)

// Services holds all services
type Services struct {
	Auth AuthService
	Todo TodoService
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, jwtSecret string) *Services {
	return &Services{
		Auth: NewAuthService(repos.User, jwtSecret),
		Todo: NewTodoService(repos.Todo),
	}
}
