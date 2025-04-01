.PHONY: help dev prod down build test lint clean docker-dev docker-prod docker-down docker-clean

# Default target
.DEFAULT_GOAL := help

# Project variables
PROJECT_NAME=myapp
DOCKER_COMPOSE=docker compose

# Help target
help: ## Display available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development environment setup
dev: ## Start development environment
	@echo "Starting development environment..."
	cp -n .env.example .env 2>/dev/null || true
	$(DOCKER_COMPOSE) up -d
	@echo "Development environment is running."
	@echo "API is available at http://localhost:8080"
	@echo "PostgreSQL is available at localhost:5432"
	@echo "PgAdmin is available at http://localhost:5050"

# Production environment setup
prod: ## Start production environment
	@echo "Starting production environment..."
	cp -n .env.example .env 2>/dev/null || true
	$(DOCKER_COMPOSE) --profile prod up -d api-prod postgres
	@echo "Production environment is running."
	@echo "API is available at http://localhost:8080"

# Stop environment
down: ## Stop development environment
	@echo "Stopping environment..."
	$(DOCKER_COMPOSE) down
	@echo "Environment stopped."

# Docker scripts
docker-dev: ## Start development environment using script
	./scripts/docker-start.sh dev

docker-prod: ## Start production environment using script
	./scripts/docker-start.sh prod

docker-down: ## Stop containers using script
	./scripts/docker-start.sh down

docker-clean: ## Clean Docker environment using script
	./scripts/docker-start.sh clean

# Build containers
build: ## Build Docker containers
	@echo "Building containers..."
	$(DOCKER_COMPOSE) build
	@echo "Build complete."

# Rebuild containers
rebuild: ## Rebuild Docker containers from scratch
	@echo "Rebuilding containers..."
	$(DOCKER_COMPOSE) build --no-cache
	@echo "Rebuild complete."

# Run tests
test: ## Run tests
	$(DOCKER_COMPOSE) exec api go test ./... -v

# Run linter
lint: ## Run linter
	$(DOCKER_COMPOSE) exec api go vet ./...
	$(DOCKER_COMPOSE) exec api golangci-lint run

# Clean up
clean: down ## Clean up development environment
	@echo "Cleaning up..."
	$(DOCKER_COMPOSE) down -v
	rm -rf tmp/
	@echo "Clean complete."

# Show logs
logs: ## Show logs from containers
	$(DOCKER_COMPOSE) logs -f

# Enter API container shell
shell: ## Open shell in API container
	$(DOCKER_COMPOSE) exec api sh

# Migrate database
migrate: ## Run database migrations
	$(DOCKER_COMPOSE) exec api go run cmd/api/main.go migrate

# Show container status
ps: ## Show container status
	$(DOCKER_COMPOSE) ps

# Restart containers
restart: ## Restart containers
	$(DOCKER_COMPOSE) restart

# Pull latest images
pull: ## Pull latest images
	$(DOCKER_COMPOSE) pull