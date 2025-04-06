package handler

import (
	"myapp/internal/service"
)

// Handler contains all HTTP handlers
type Handler struct {
	UserHandler *UserHandler
	TodoHandler *TodoHandler
}

// New creates a new Handler instance
func New(authService service.AuthService, todoService service.TodoService) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(authService),
		TodoHandler: NewTodoHandler(todoService),
	}
}
