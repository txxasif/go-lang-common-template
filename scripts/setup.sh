#!/bin/bash

# Ensure the script exits if any command fails
set -e

# Text formatting
BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Setting up myapp development environment...${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker is not installed. Please install Docker first.${NC}"
    echo "Visit https://docs.docker.com/get-docker/ for installation instructions."
    exit 1
fi

# Check if Docker Compose is available (integrated with Docker CLI in newer versions)
if ! docker compose version &> /dev/null; then
    echo -e "${RED}Docker Compose is not available.${NC}"
    echo "Make sure you have Docker Compose V2 installed or Docker Desktop with Compose V2 enabled."
    exit 1
fi

# Create necessary directories
mkdir -p tmp

# Ensure .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file from example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}.env file created. You may want to edit it with your specific configuration.${NC}"
fi

# Create docker directories if they don't exist
mkdir -p docker

# Create Dockerfiles if they don't exist
if [ ! -f docker/dev.Dockerfile ]; then
    echo -e "${YELLOW}Creating development Dockerfile...${NC}"
    cat > docker/dev.Dockerfile <<EOL
FROM golang:1.20-alpine

# Set working directory
WORKDIR /app

# Install required packages for development
RUN apk add --no-cache git curl make gcc g++ tzdata

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin v1.55.2

# Set environment variables
ENV GO111MODULE=on \\
    CGO_ENABLED=0 \\
    GOOS=linux \\
    GOARCH=amd64

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Command is set in docker-compose.yml
EOL
    echo -e "${GREEN}Development Dockerfile created.${NC}"
fi

# Start the development environment
echo -e "${BLUE}Starting development environment...${NC}"
make dev

echo -e "${GREEN}Setup complete! Your development environment is running.${NC}"
echo -e "API is accessible at ${BLUE}http://localhost:8080${NC}"
echo -e "PostgreSQL is accessible at ${BLUE}localhost:5432${NC}"
echo -e ""
echo -e "Available commands:"
echo -e "${YELLOW}make dev${NC}      - Start development environment"
echo -e "${YELLOW}make down${NC}     - Stop development environment"
echo -e "${YELLOW}make logs${NC}     - View logs"
echo -e "${YELLOW}make shell${NC}    - Open shell in API container"
echo -e "${YELLOW}make test${NC}     - Run tests"
echo -e "${YELLOW}make lint${NC}     - Run linters"
echo -e "${YELLOW}make help${NC}     - Show all available commands"