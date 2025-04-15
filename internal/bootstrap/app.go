package bootstrap

import (
	"fmt"
	"net/http"

	"myapp/internal/config"
	"myapp/internal/handler"
	"myapp/internal/pkg/jwt"
	"myapp/internal/repository"
	"myapp/internal/router"
	"myapp/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// App holds all application dependencies
type App struct {
	Config   *config.Config
	Database *gorm.DB
	Router   http.Handler
}

// NewApp creates a new App instance
func NewApp(cfg *config.Config) (*App, error) {
	// Setup database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize JWT service with config
	jwtService := jwt.NewServiceWithConfig(jwt.ServiceConfig{
		SecretKey:     cfg.JWT.Secret,
		AccessExpiry:  cfg.JWT.AccessExpiresIn,
		RefreshExpiry: cfg.JWT.RefreshExpiresIn,
	})

	// Initialize services
	authService := service.NewAuthService(repos.User, jwtService)
	todoService := service.NewTodoService(repos.Todo)

	// Initialize handlers
	h := &handler.Handler{
		UserHandler: handler.NewUserHandler(authService),
		TodoHandler: handler.NewTodoHandler(todoService),
	}

	// Setup router
	r := router.New(h, authService)

	return &App{
		Config:   cfg,
		Database: db,
		Router:   r,
	}, nil
}
