package repository

import "myapp/internal/db"

type Repositories struct {
	User UserRepository
	// Add other repos here
}

func NewRepositories(db *db.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		// Initialize other repos
	}
}
