package repository

import (
	"gorm.io/gorm"
)

// Repositories holds all repository interfaces
type Repositories struct {
	User UserRepository
	Todo TodoRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		Todo: NewTodoRepository(db),
	}
}
