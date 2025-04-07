package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"myapp/internal/middleware"
	"myapp/internal/model"
	httputil "myapp/internal/pkg/http"
	"myapp/internal/service"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	todoService service.TodoService
}

// NewTodoHandler creates a new TodoHandler instance
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// Create handles the creation of a new todo
// @Summary Create a new todo
// @Description Create a new todo item for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.TodoCreateRequest true "Todo creation details"
// @Success 201 {object} model.Todo
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos [post]
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req model.TodoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	todo, err := h.todoService.Create(r.Context(), userID, &req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	httputil.JSON(w, http.StatusCreated, todo)
}

// GetByID handles retrieving a todo by ID
// @Summary Get a todo by ID
// @Description Get a specific todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 200 {object} model.Todo
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/{id} [get]
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get todo ID from URL parameters
	todoID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.GetByID(r.Context(), uint(todoID))
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "Failed to get todo")
		return
	}

	// Check if the todo belongs to the user
	if todo.UserID != userID {
		httputil.Error(w, http.StatusForbidden, "Access denied")
		return
	}

	httputil.JSON(w, http.StatusOK, todo)
}

// GetByUserID handles retrieving all todos for a user
// @Summary Get all todos for a user
// @Description Get all todo items for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Todo
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos [get]
func (h *TodoHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	todos, err := h.todoService.GetByUserID(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "Failed to get todos")
		return
	}

	httputil.JSON(w, http.StatusOK, todos)
}

// Update handles updating a todo
// @Summary Update a todo
// @Description Update an existing todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Param request body model.TodoUpdateRequest true "Todo update details"
// @Success 200 {object} model.Todo
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get todo ID from URL parameters
	todoID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req model.TodoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	todo, err := h.todoService.Update(r.Context(), userID, uint(todoID), &req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	httputil.JSON(w, http.StatusOK, todo)
}

// Delete handles deleting a todo
// @Summary Delete a todo
// @Description Delete an existing todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get todo ID from URL parameters
	todoID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	if err := h.todoService.Delete(r.Context(), userID, uint(todoID)); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
