# MyApp

A Go-based REST API application with user authentication and todo management.

## Features

- User registration and authentication
- JWT-based authentication
- Todo management (CRUD operations)
- Swagger API documentation
- Hot reloading for development

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/myapp.git
cd myapp
```

2. Install dependencies:

```bash
go mod download
```

3. Set up the database:

```bash
# Create a PostgreSQL database
createdb myapp

# Run migrations
go run cmd/migrate/main.go
```

4. Configure the application:

```bash
# Copy the example config file
cp config.example.yaml config.yaml

# Edit config.yaml with your settings
```

## Development

### Running the Application

1. Start the development server with hot reloading:

```bash
# Using air (recommended for development)
air

# Or using go run
go run cmd/api/main.go
```

2. Access the API:

```bash
# API base URL
http://localhost:8080

# Swagger documentation
http://localhost:8080/swagger/
```

### API Documentation

The API is documented using Swagger/OpenAPI. You can access the documentation at:

```
http://localhost:8080/swagger/
```

To update the Swagger documentation after making changes to the API:

```bash
swag init -g cmd/api/main.go -o docs
```

### Project Structure

```
myapp/
├── cmd/              # Application entry points
│   ├── api/          # Main API server
│   └── migrate/      # Database migrations
├── internal/         # Private application code
│   ├── config/       # Configuration
│   ├── handler/      # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   ├── model/        # Data models
│   ├── pkg/          # Shared packages
│   ├── repository/   # Data access layer
│   ├── router/       # HTTP router
│   └── service/      # Business logic
├── docs/             # Swagger documentation
├── migrations/       # Database migration files
└── scripts/          # Utility scripts
```

### Adding New API Endpoints

1. Create a new handler function in the appropriate handler file
2. Add Swagger annotations to document the endpoint:

```go
// @Summary Short description
// @Description Longer description
// @Tags tag-name
// @Accept json
// @Produce json
// @Param parameter-name parameter-type required "Description"
// @Success status-code {object} response-type
// @Failure status-code {object} error-type
// @Router /path [method]
```

3. Update the router to include the new endpoint
4. Generate updated Swagger documentation:

```bash
swag init -g cmd/api/main.go -o docs
```

## Testing

Run the test suite:

```bash
go test ./...
```

## Deployment

1. Build the application:

```bash
go build -o myapp cmd/api/main.go
```

2. Run the application:

```bash
./myapp
```

## License
