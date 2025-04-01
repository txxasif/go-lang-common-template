package service

import (
	"myapp/internal/config"
	"myapp/internal/repository"
)

type Services struct {
	Auth AuthService
	// Add other services here
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth: NewAuthService(repos.User, cfg),
		// Initialize other services
	}
}
