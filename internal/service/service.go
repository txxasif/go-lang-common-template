package service

import (
	"time"
)

// Services holds all services
type Services struct {
	Auth AuthService
	Todo TodoService
}

// JWTConfig contains the configuration for JWT tokens
type JWTConfig struct {
	Secret           string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}
