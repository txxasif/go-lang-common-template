package repository

import (
	"context"
	"myapp/internal/model"

	"gorm.io/gorm"
)

// TodoRepository defines the interface for todo operations
type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	GetByID(ctx context.Context, id uint) (*model.Todo, error)
	GetByUserID(ctx context.Context, userID uint) ([]*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new TodoRepository instance
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) GetByID(ctx context.Context, id uint) (*model.Todo, error) {
	var todo model.Todo
	if err := r.db.WithContext(ctx).First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.Todo, error) {
	var todos []*model.Todo
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Todo{}, id).Error
}
