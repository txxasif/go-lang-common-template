# Technical Context

## Technology Stack

- **Language**: Go 1.20+
- **Router**: Chi
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose
- **Validation**: Custom validation layer

## Development Environment

- **Version Control**: Git
- **API Testing**: Postman
- **Hot Reload**: Air
- **Environment Variables**: .env file
- **Build Tool**: Make
- **Code Organization**: Clean Architecture

## Project Dependencies

- Chi router for HTTP routing and middleware
- GORM for database operations
- JWT-go for authentication
- bcrypt for password hashing
- godotenv for environment configuration
- Custom validation layer for input handling

## Architecture Components

### Core Layers

1. Handlers (HTTP Layer)
2. Services (Business Logic)
3. Repositories (Data Access)
4. Models (Data Structures)
5. Validation (Input Handling)

### Key Directories

- `cmd/`: Application entry points
- `internal/`: Private application code
  - `bootstrap/`: Application initialization
  - `config/`: Configuration management
  - `db/`: Database setup
  - `handler/`: HTTP handlers
  - `middleware/`: HTTP middleware
  - `model/`: Data models
  - `pkg/`: Internal packages
  - `repository/`: Data access layer
  - `router/`: Route definitions
  - `service/`: Business logic
  - `validation/`: Input validation
- `scripts/`: Utility scripts
- `docker/`: Container configurations

## Development Guidelines

1. Follow clean architecture principles
2. Use dependency injection
3. Implement proper error handling
4. Write comprehensive tests
5. Document all public APIs
6. Validate all inputs
7. Maintain clear layer separation

## Build and Deployment

### Local Development

```bash
make run # Run locally
make test # Run tests
```

### Docker Deployment

```bash
docker-compose up -d # Development
docker-compose -f docker-compose.prod.yml up -d # Production
```

## Security Considerations

- JWT token expiration
- Password hashing with bcrypt
- Environment variable protection
- CORS configuration
- Input validation
- Rate limiting implementation

## Performance Optimization

- Database connection pooling
- Query optimization
- Proper indexing
- Caching strategies (when implemented)
- Validation efficiency

## Monitoring and Logging

- Structured logging
- Error tracking
- Performance metrics
- Health checks
- Validation metrics
