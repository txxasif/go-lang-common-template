package service

import "errors"

var (
	// ErrEmailAlreadyExists is returned when a user with the same email already exists
	ErrEmailAlreadyExists = errors.New("email already exists")
	// ErrUsernameAlreadyExists is returned when a user with the same username already exists
	ErrUsernameAlreadyExists = errors.New("username already exists")
)
