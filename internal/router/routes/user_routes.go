package routes

import (
	"myapp/internal/handler"

	"github.com/go-chi/chi/v5"
)

// SetupUserRoutes sets up user-related routes
func SetupUserRoutes(r chi.Router, h *handler.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/me", h.UserHandler.GetProfile)
	})
}
