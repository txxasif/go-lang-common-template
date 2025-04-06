package model

import (
	"time"
)

// Todo represents a todo item in the system
type Todo struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	User        User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TodoCreateRequest represents the request body for creating a todo
type TodoCreateRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// TodoUpdateRequest represents the request body for updating a todo
type TodoUpdateRequest struct {
	Title       string `json:"title" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
	Completed   bool   `json:"completed"`
}
