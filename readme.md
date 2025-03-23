# Go API with Chi Router and GORM

A scalable Go API boilerplate using Chi router and GORM with PostgreSQL. This project follows best practices for structure, authentication, and deployment.

## Features

- **Clean Architecture**: Organized using feature-based modular structure
- **Authentication**: JWT-based auth with register and login endpoints
- **Database**: PostgreSQL integration with GORM for ORM
- **Routing**: Chi router with middleware support
- **API Documentation**: Postman collection included
- **Containerization**: Docker and Docker Compose setup

## Project Structure

```
myapp/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── db.go
│   ├── features/
│   │   ├── auth/
│   │   │   ├── handler.go
│   │   │   ├── models.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   └── service.go
│   │   ├── user/
│   │   │   ├── handler.go
│   │   │   ├── models.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   └── service.go
│   │   └── product/
│   │       ├── handler.go
│   │       ├── models.go
│   │       ├── repository.go
│   │       ├── routes.go
│   │       └── service.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── middleware.go
│   ├── model/
│   │   └── common.go
│   └── router/
│       └── router.go
├── pkg/
│   ├── hash/
│   │   └── password.go
│   └── jwt/
│       └── jwt.go
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL
- Docker (optional)

### Local Development

1. Clone the repository

   ```bash
   git clone https://github.com/yourusername/myapp.git
   cd myapp
   ```

2. Install dependencies

   ```bash
   go mod download
   ```

3. Set up environment variables

   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

4. Run the application
   ```bash
   go run cmd/api/main.go
   ```

### Using Docker

1. Build and start the containers

   ```bash
   docker-compose up -d
   ```

2. The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication

- **POST /api/register** - Create a new user account
- **POST /api/login** - Authenticate and get a JWT token

### User Management

- **GET /api/user** - Get current user profile (authenticated)
- **PUT /api/user** - Update user profile (authenticated)
- **GET /api/users** - List all users (authenticated)

### Products

- **GET /api/products** - List all products (public)
- **GET /api/products/{id}** - Get a specific product (public)
- **POST /api/products** - Create a new product (authenticated)
- **PUT /api/products/{id}** - Update a product (authenticated)
- **DELETE /api/products/{id}** - Delete a product (authenticated)

## Testing the API

Import the included Postman collection to test the API endpoints:

1. Open Postman
2. Import the `GoApp API.postman_collection.json` file
3. Set the `base_url` variable to your API URL (default: `http://localhost:8080`)
4. Use the Register endpoint to create a user
5. Use the Login endpoint to obtain a JWT token
6. Test other endpoints with the token automatically included in requests

## Project Design Decisions

### Router Organization

Routes are organized by feature with separate registration for each module, keeping the main.go file clean and focused.

### Authentication

JWT-based authentication with middleware to protect routes:

- Tokens expire after a configurable time
- Password hashing using bcrypt
- Middleware for securing routes

### Database Layer

- Repository pattern with GORM
- Clean separation between database operations and business logic
- Automatic migrations

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
