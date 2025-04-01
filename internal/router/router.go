package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"myapp/internal/handler"
	appMiddleware "myapp/internal/middleware"
	"myapp/internal/service"
)

type Router struct {
	chi.Router
}

func New(h *handler.Handler, services *service.Services) *Router {
	r := &Router{
		Router: chi.NewRouter(),
	}

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Auth routes
		setupAuthRoutes(r, h, services.Auth)

		// User routes
		setupUserRoutes(r, h, services.Auth)

		// Other routes can be added in their own setup functions
		// setupProductRoutes(r, h, services)
		// setupOrderRoutes(r, h, services)
	})

	return r
}

// setupAuthRoutes configures all authentication-related routes
func setupAuthRoutes(r chi.Router, h *handler.Handler, authService service.AuthService) {
	r.Group(func(r chi.Router) {
		r.Post("/login", h.Auth.Login)
		r.Post("/register", h.Auth.Register)
	})
}

// setupUserRoutes configures all user-related routes
func setupUserRoutes(r chi.Router, h *handler.Handler, authService service.AuthService) {
	r.Group(func(r chi.Router) {
		// Protected routes
		r.Use(appMiddleware.Auth(authService))
		r.Get("/user", h.Auth.GetUser)

		// Additional user routes would go here
		// r.Put("/user", h.User.UpdateUser)
		// r.Delete("/user", h.User.DeleteUser)
	})
}

// Additional route setup functions would be defined here
// func setupProductRoutes(r chi.Router, h *handler.Handler, authService service.AuthService) {...}
// func setupOrderRoutes(r chi.Router, h *handler.Handler, authService service.AuthService) {...}
