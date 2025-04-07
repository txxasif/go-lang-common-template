package service

import (
	"context"
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

// TodoService defines the interface for todo operations
type TodoService interface {
	Create(ctx context.Context, userID uint, req *model.TodoCreateRequest) (*model.Todo, error)
	GetByID(ctx context.Context, id uint) (*model.Todo, error)
	GetByUserID(ctx context.Context, userID uint) ([]*model.Todo, error)
	Update(ctx context.Context, userID uint, id uint, req *model.TodoUpdateRequest) (*model.Todo, error)
	Delete(ctx context.Context, userID uint, id uint) error
}

type todoService struct {
	todoRepo repository.TodoRepository
}

// NewTodoService creates a new TodoService instance
func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoService{todoRepo: todoRepo}
}

func (s *todoService) Create(ctx context.Context, userID uint, req *model.TodoCreateRequest) (*model.Todo, error) {
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) GetByID(ctx context.Context, id uint) (*model.Todo, error) {
	return s.todoRepo.GetByID(ctx, id)
}

func (s *todoService) GetByUserID(ctx context.Context, userID uint) ([]*model.Todo, error) {
	return s.todoRepo.GetByUserID(ctx, userID)
}

func (s *todoService) Update(ctx context.Context, userID uint, id uint, req *model.TodoUpdateRequest) (*model.Todo, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if todo.UserID != userID {
		return nil, errors.New("unauthorized: todo does not belong to user")
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	todo.Completed = req.Completed

	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) Delete(ctx context.Context, userID uint, id uint) error {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.UserID != userID {
		return errors.New("unauthorized: todo does not belong to user")
	}

	return s.todoRepo.Delete(ctx, id)
}
