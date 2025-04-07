# System Patterns

## Architecture Overview

The application follows a clean architecture pattern with clear separation of concerns:

```
Application Structure
├── internal/
│   ├── bootstrap/     # Application initialization
│   ├── config/       # Configuration management
│   ├── db/          # Database setup
│   ├── handler/     # HTTP handlers
│   ├── middleware/  # HTTP middleware
│   ├── model/       # Data models
│   ├── pkg/         # Internal packages
│   ├── repository/  # Data access layer
│   ├── router/      # Route definitions
│   ├── service/     # Business logic
│   └── validation/  # Input validation
```

## Design Patterns

### Repository Pattern

- Abstracts data persistence operations
- Enables swappable data sources
- Centralizes data access logic
- Located in `internal/repository/`

### Service Layer Pattern

- Contains business logic
- Orchestrates data flow
- Handles complex operations
- Located in `internal/service/`

### Validation Pattern

- Input validation layer
- Request validation
- Data sanitization
- Located in `internal/validation/`

### Middleware Pattern

- Authentication middleware
- Logging middleware
- CORS middleware
- Recovery middleware
- Located in `internal/middleware/`

## Component Relationships

### Request Flow

1. HTTP Request → Router
2. Router → Middleware Chain
3. Middleware → Handler
4. Handler → Service
5. Service → Repository
6. Repository → Database

### Validation Flow

1. Request received
2. Validation layer processes
3. Handler processes valid request
4. Service executes business logic
5. Repository handles data access

## System Boundaries

### External Interfaces

- REST API endpoints
- Database connections
- Environment configuration
- External service integrations (when added)

### Internal Boundaries

- Clear separation between layers
- Service layer for business logic
- Repository layer for data access
- Validation layer for input handling

## Error Handling

- Consistent error types
- Proper HTTP status codes
- Detailed error messages
- Error wrapping when appropriate
- Validation error handling

## Code Organization

### Layer Structure

Each layer follows a consistent pattern:

```go
// handler.go
type Handler struct {
    service Service
}

// service.go
type Service struct {
    repo Repository
}

// repository.go
type Repository struct {
    db *gorm.DB
}

// validation.go
type Validator struct {
    // validation rules
}
```

### Common Patterns

- Constructor injection
- Interface-based design
- Error wrapping
- Consistent naming conventions
- Standard CRUD operations
- Input validation

## Testing Strategy

- Unit tests for business logic
- Integration tests for repositories
- API tests for endpoints
- Validation tests
- Mock interfaces for dependencies
