package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"myapp/internal/handler"
	authmiddleware "myapp/internal/middleware"
	"myapp/internal/service"
)

// Router wraps chi.Router to add custom functionality
type Router struct {
	chi.Router
}

// New creates a new router with all routes configured
func New(h *handler.Handler, authService service.AuthService) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/register", h.UserHandler.Register)
		r.Post("/login", h.UserHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Auth middleware
		r.Use(authmiddleware.Auth(authService))

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/me", h.UserHandler.GetProfile)
		})

		// Todo routes
		r.Route("/todos", func(r chi.Router) {
			r.Post("/", h.TodoHandler.Create)
			r.Get("/", h.TodoHandler.GetByUserID)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.TodoHandler.GetByID)
				r.Put("/", h.TodoHandler.Update)
				r.Delete("/", h.TodoHandler.Delete)
			})
		})
	})

	return r
}

// RegisterRouteWithMiddleware registers a route with the given middleware
func (r *Router) RegisterRouteWithMiddleware(method, path string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	// Apply middlewares in reverse order
	h := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](http.HandlerFunc(h)).ServeHTTP
	}

	// Register route
	r.Method(method, path, http.HandlerFunc(h))
}

// RegisterRouteWithValidation registers a route with validation middleware
func (r *Router) RegisterRouteWithValidation(method, path string, handler http.HandlerFunc, model interface{}, middlewares ...func(http.Handler) http.Handler) {
	// Add validation middleware
	middlewares = append(middlewares, ValidateRequest(model))

	// Register route with middleware
	r.RegisterRouteWithMiddleware(method, path, handler, middlewares...)
}

// ValidateRequest returns a middleware that validates the request body against the given model
func ValidateRequest(model interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement request validation
			next.ServeHTTP(w, r)
		})
	}
}

// Additional route setup functions would be defined here
// func setupProductRoutes(r chi.Router, h *handler.Handler, authService service.AuthService) {...}
// func setupOrderRoutes(r chi.Router, h *handler.Handler, authService service.AuthService) {...}
