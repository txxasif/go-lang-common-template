package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"myapp/internal/handler"
	authmiddleware "myapp/internal/middleware"
	"myapp/internal/router/routes"
	"myapp/internal/service"
)

// New creates a new router with all routes configured
func New(h *handler.Handler, authService service.AuthService) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	// Swagger UI routes (public)
	r.Group(func(r chi.Router) {
		routes.SetupSwaggerRoutes(r)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		routes.SetupAuthRoutes(r, h)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Auth middleware
		r.Use(authmiddleware.Auth(authService))

		// User routes
		routes.SetupUserRoutes(r, h)

		// Todo routes
		routes.SetupTodoRoutes(r, h)
	})

	return r
}
