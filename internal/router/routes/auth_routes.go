package routes

import (
	"myapp/internal/handler"

	"github.com/go-chi/chi/v5"
)

// SetupAuthRoutes sets up authentication-related routes
func SetupAuthRoutes(r chi.Router, h *handler.Handler) {
	r.Post("/register", h.UserHandler.Register)
	r.Post("/login", h.UserHandler.Login)
}
