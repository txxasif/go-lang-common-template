# System Patterns

## Architecture Overview

The application follows a clean architecture pattern with clear separation of concerns:

```
Feature Module
├── handler.go (HTTP Layer)
├── service.go (Business Logic)
├── repository.go (Data Access)
├── models.go (Data Structures)
└── routes.go (Route Definitions)
```

## Design Patterns

### Repository Pattern

- Abstracts data persistence operations
- Enables swappable data sources
- Centralizes data access logic
- Example implementation in `internal/features/*/repository.go`

### Dependency Injection

- Services receive dependencies through constructors
- Facilitates testing and modularity
- Reduces tight coupling between components
- Used consistently across feature modules

### Middleware Pattern

- Authentication middleware
- Logging middleware
- CORS middleware
- Recovery middleware
- Located in `internal/middleware/`

### Service Layer Pattern

- Contains business logic
- Orchestrates data flow
- Handles complex operations
- Independent of HTTP and database concerns

## Component Relationships

### Request Flow

1. HTTP Request → Router
2. Router → Middleware Chain
3. Middleware → Handler
4. Handler → Service
5. Service → Repository
6. Repository → Database

### Authentication Flow

1. Client sends credentials
2. Auth service validates
3. JWT token generated
4. Token included in subsequent requests
5. Auth middleware validates token

## System Boundaries

### External Interfaces

- REST API endpoints
- Database connections
- Environment configuration
- External service integrations (when added)

### Internal Boundaries

- Feature modules are self-contained
- Shared utilities in `pkg/`
- Configuration in `internal/config/`
- Database setup in `internal/db/`

## Error Handling

- Consistent error types
- Proper HTTP status codes
- Detailed error messages
- Error wrapping when appropriate

## Code Organization

### Feature Module Structure

Each feature follows the same structure:

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
```

### Common Patterns

- Constructor injection
- Interface-based design
- Error wrapping
- Consistent naming conventions
- Standard CRUD operations

## Testing Strategy

- Unit tests for business logic
- Integration tests for repositories
- API tests for endpoints
- Mock interfaces for dependencies
