package bootstrap

import (
	"myapp/internal/config"
	"myapp/internal/db"
	"myapp/internal/handler"
	"myapp/internal/repository"
	"myapp/internal/router"
	"myapp/internal/service"
)

type App struct {
	Config   *config.Config
	Database *db.DB
	Router   *router.Router
	Handler  *handler.Handler
	Services *service.Services
	Repos    *repository.Repositories
}

func NewApp() (*App, error) {
	// Load config
	cfg := config.New()

	// Setup DB
	database, err := db.Connect(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repos
	repos := repository.NewRepositories(database)

	// Initialize services
	services := service.NewServices(repos, cfg)

	// Initialize handlers
	h := handler.New(
		handler.WithAuthHandler(services.Auth),
	)

	// Setup router
	r := router.New(h, services)

	return &App{
		Config:   cfg,
		Database: database,
		Router:   r,
		Handler:  h,
		Services: services,
		Repos:    repos,
	}, nil
}
