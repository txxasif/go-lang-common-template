package model

import "errors"

var (
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when a user already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrEmailAlreadyExists is returned when a user with the same email already exists
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrUsernameAlreadyExists is returned when a user with the same username already exists
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token expired")

	// ErrUnauthorized is returned when a user is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrValidation is returned when validation fails
	ErrValidation = errors.New("validation failed")
)
