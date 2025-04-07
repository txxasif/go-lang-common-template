package routes

import (
	"myapp/internal/handler"

	"github.com/go-chi/chi/v5"
)

// SetupTodoRoutes sets up all todo-related routes
func SetupTodoRoutes(r chi.Router, h *handler.Handler) {
	r.Route("/todos", func(r chi.Router) {
		r.Post("/", h.TodoHandler.Create)
		r.Get("/", h.TodoHandler.GetByUserID)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.TodoHandler.GetByID)
			r.Put("/", h.TodoHandler.Update)
			r.Delete("/", h.TodoHandler.Delete)
		})
	})
}
