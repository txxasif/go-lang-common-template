# Application Documentation: A Complete Guide

## Technical Architecture Overview

### Service Layer Organization

The application follows a clean architecture pattern with well-defined service interfaces and implementations:

```go
// Service Interfaces
type AuthService interface {
    Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error)
    Login(ctx context.Context, email, password string) (*model.AuthResponse, error)
    GetUserByToken(ctx context.Context, token string) (*model.User, error)
}

type TodoService interface {
    Create(ctx context.Context, userID uint, req *model.TodoCreateRequest) (*model.Todo, error)
    GetByID(ctx context.Context, userID, todoID uint) (*model.Todo, error)
    List(ctx context.Context, userID uint) ([]*model.Todo, error)
    Update(ctx context.Context, userID, todoID uint, req *model.TodoUpdateRequest) (*model.Todo, error)
    Delete(ctx context.Context, userID, todoID uint) error
}

// Service Implementations
type authService struct {
    userRepo repository.UserRepository
    jwt      *JWTService
}

type todoService struct {
    todoRepo repository.TodoRepository
}
```

Key points about service organization:

- Each service has a clear interface defining its contract
- Services are responsible for business logic and orchestration
- Dependencies are injected through constructors
- Services coordinate between repositories and handle business rules

### Repository Layer Structure

The repository layer provides data access abstraction:

```go
// Repository Interfaces
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id uint) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id uint) error
}

type TodoRepository interface {
    Create(ctx context.Context, todo *model.Todo) error
    GetByID(ctx context.Context, id uint) (*model.Todo, error)
    List(ctx context.Context, userID uint) ([]*model.Todo, error)
    Update(ctx context.Context, todo *model.Todo) error
    Delete(ctx context.Context, id uint) error
}

// Repository Implementations
type userRepository struct {
    db *gorm.DB
}

type todoRepository struct {
    db *gorm.DB
}
```

Repository layer characteristics:

- Abstracts database operations
- Implements data access patterns
- Handles database-specific logic
- Provides a clean interface for services

### Handler Organization

Handlers are organized by domain and follow a consistent pattern:

```go
// Handler Groups
type Handler struct {
    UserHandler *userHandler
    TodoHandler *todoHandler
}

// Individual Handlers
type userHandler struct {
    authService service.AuthService
}

type todoHandler struct {
    todoService service.TodoService
}

// Handler Methods
func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Request parsing
    // Validation
    // Service call
    // Response formatting
}

func (h *todoHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Authentication check
    // Request parsing
    // Validation
    // Service call
    // Response formatting
}
```

Handler organization principles:

- Each handler focuses on a specific domain
- Handlers delegate business logic to services
- Consistent error handling and response formatting
- Clear separation of concerns

### Route Setup and Middleware

Routes are organized using Chi router with middleware:

```go
func New(h *handler.Handler, authService service.AuthService) http.Handler {
    r := chi.NewRouter()

    // Global Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.CORS)

    // Public Routes
    r.Group(func(r chi.Router) {
        r.Post("/register", h.UserHandler.Register)
        r.Post("/login", h.UserHandler.Login)
    })

    // Protected Routes
    r.Group(func(r chi.Router) {
        r.Use(authmiddleware.Auth(authService))

        // User Routes
        r.Get("/users/me", h.UserHandler.GetProfile)
        r.Put("/users/me", h.UserHandler.UpdateProfile)

        // Todo Routes
        r.Post("/todos", h.TodoHandler.Create)
        r.Get("/todos", h.TodoHandler.List)
        r.Get("/todos/{id}", h.TodoHandler.GetByID)
        r.Put("/todos/{id}", h.TodoHandler.Update)
        r.Delete("/todos/{id}", h.TodoHandler.Delete)
    })

    return r
}
```

Route organization features:

- Clear separation of public and protected routes
- Consistent middleware application
- Logical grouping of related endpoints
- RESTful route naming conventions

### Dependency Injection

The application uses constructor-based dependency injection:

```go
func NewApp() (*App, error) {
    // Load configuration
    cfg := config.Load()

    // Initialize database
    db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    todoRepo := repository.NewTodoRepository(db)

    // Initialize services
    authService := service.NewAuthService(userRepo, jwtService)
    todoService := service.NewTodoService(todoRepo)

    // Initialize handlers
    userHandler := handler.NewUserHandler(authService)
    todoHandler := handler.NewTodoHandler(todoService)
    handlers := handler.New(userHandler, todoHandler)

    // Initialize router
    router := router.New(handlers, authService)

    return &App{
        Config:   cfg,
        Database: db,
        Router:   router,
    }, nil
}
```

Dependency injection benefits:

- Clear dependency relationships
- Easy testing through mocking
- Flexible component replacement
- Explicit initialization order

## Component Interactions Deep Dive

### Repository Layer Implementation Details

```go
// Base Repository Interface
type BaseRepository[T any] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id uint) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id uint) error
}

// User Repository Implementation
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, nil
    }
    return &user, err
}

// Todo Repository Implementation
type todoRepository struct {
    db *gorm.DB
}

func (r *todoRepository) List(ctx context.Context, userID uint) ([]*model.Todo, error) {
    var todos []*model.Todo
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Find(&todos).Error
    return todos, err
}
```

Key aspects of repository implementation:

- Uses GORM for database operations
- Implements context for request-scoped operations
- Handles error cases appropriately
- Provides type-safe operations

### Service Layer Implementation Details

```go
// Auth Service Implementation
type authService struct {
    userRepo repository.UserRepository
    jwt      *JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwt *JWTService) service.AuthService {
    return &authService{
        userRepo: userRepo,
        jwt:      jwt,
    }
}

func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
    // 1. Check for existing user
    existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to check email: %w", err)
    }
    if existingUser != nil {
        return nil, model.ErrEmailAlreadyExists
    }

    // 2. Create new user
    hashedPassword, err := HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    user := &model.User{
        Email:     req.Email,
        Username:  req.Username,
        Password:  hashedPassword,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 3. Generate JWT token
    token, err := s.jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10))
    if err != nil {
        return nil, fmt.Errorf("failed to generate token: %w", err)
    }

    return &model.AuthResponse{
        User:  user.ToResponse(),
        Token: token,
    }, nil
}

// Todo Service Implementation
type todoService struct {
    todoRepo repository.TodoRepository
}

func (s *todoService) Create(ctx context.Context, userID uint, req *model.TodoCreateRequest) (*model.Todo, error) {
    todo := &model.Todo{
        UserID:      userID,
        Title:       req.Title,
        Description: req.Description,
        Completed:   false,
    }

    if err := s.todoRepo.Create(ctx, todo); err != nil {
        return nil, fmt.Errorf("failed to create todo: %w", err)
    }

    return todo, nil
}
```

Service layer characteristics:

- Implements business logic
- Coordinates between repositories
- Handles error wrapping
- Manages transactions when needed

### Handler Layer Implementation Details

```go
// User Handler Implementation
type userHandler struct {
    authService service.AuthService
}

func NewUserHandler(authService service.AuthService) *userHandler {
    return &userHandler{
        authService: authService,
    }
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req model.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.NewServiceError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body").Write(w)
        return
    }

    // 2. Validate request
    if errs := validation.ValidateRegisterRequest(&req); errs.HasErrors() {
        response.NewValidationError(errs.Errors).Write(w)
        return
    }

    // 3. Call service
    resp, err := h.authService.Register(r.Context(), &req)
    if err != nil {
        switch {
        case errors.Is(err, model.ErrEmailAlreadyExists):
            response.NewServiceError(http.StatusConflict, "EMAIL_EXISTS", "Email already exists").Write(w)
        default:
            response.NewServiceError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user").Write(w)
        }
        return
    }

    // 4. Return response
    response.NewSuccess(http.StatusCreated, resp).Write(w)
}

// Todo Handler Implementation
type todoHandler struct {
    todoService service.TodoService
}

func (h *todoHandler) Create(w http.ResponseWriter, r *http.Request) {
    // 1. Get user from context
    userID, err := middleware.GetUserIDFromContext(r)
    if err != nil {
        response.NewServiceError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized").Write(w)
        return
    }

    // 2. Parse request
    var req model.TodoCreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.NewServiceError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body").Write(w)
        return
    }

    // 3. Validate request
    if errs := validation.ValidateTodoCreateRequest(&req); errs.HasErrors() {
        response.NewValidationError(errs.Errors).Write(w)
        return
    }

    // 4. Call service
    todo, err := h.todoService.Create(r.Context(), userID, &req)
    if err != nil {
        response.NewServiceError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create todo").Write(w)
        return
    }

    // 5. Return response
    response.NewSuccess(http.StatusCreated, todo).Write(w)
}
```

Handler layer characteristics:

- Handles HTTP-specific concerns
- Manages request/response lifecycle
- Implements error handling
- Coordinates with services

### Component Interaction Flow

Here's how components interact during a typical request:

1. **HTTP Request Arrives**

   ```go
   // Router matches request to handler
   r.Post("/register", h.UserHandler.Register)
   ```

2. **Handler Processing**

   ```go
   // Handler receives request
   func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
       // Parse and validate request
       // Call service layer
       resp, err := h.authService.Register(r.Context(), &req)
   }
   ```

3. **Service Layer Processing**

   ```go
   // Service implements business logic
   func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
       // Check for existing user
       existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)

       // Create new user
       user := &model.User{...}
       if err := s.userRepo.Create(ctx, user); err != nil {
           return nil, err
       }

       // Generate token
       token, err := s.jwt.GenerateToken(...)
   }
   ```

4. **Repository Layer Processing**

   ```go
   // Repository handles database operations
   func (r *userRepository) Create(ctx context.Context, user *model.User) error {
       return r.db.WithContext(ctx).Create(user).Error
   }
   ```

5. **Response Flow**
   ```go
   // Handler formats response
   response.NewSuccess(http.StatusCreated, resp).Write(w)
   ```

### Error Handling Strategy

The application implements a consistent error handling approach:

```go
// Custom Error Types
var (
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrUserNotFound      = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
)

// Error Response Structure
type ErrorResponse struct {
    Error   string            `json:"error"`
    Message string            `json:"message"`
    Errors  []ValidationError `json:"errors,omitempty"`
}

// Error Handling in Handlers
func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
    resp, err := h.authService.Register(r.Context(), &req)
    if err != nil {
        switch {
        case errors.Is(err, model.ErrEmailAlreadyExists):
            response.NewServiceError(http.StatusConflict, "EMAIL_EXISTS", "Email already exists").Write(w)
        case errors.Is(err, model.ErrInvalidCredentials):
            response.NewServiceError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials").Write(w)
        default:
            response.NewServiceError(http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred").Write(w)
        }
        return
    }
}
```

### Transaction Management

For operations requiring multiple database changes:

```go
func (s *todoService) UpdateWithTags(ctx context.Context, userID, todoID uint, req *model.TodoUpdateRequest) (*model.Todo, error) {
    // Start transaction
    tx := s.db.Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }

    // Defer rollback in case of error
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Update todo
    todo, err := s.todoRepo.Update(ctx, todoID, req)
    if err != nil {
        tx.Rollback()
        return nil, err
    }

    // Update tags
    if err := s.tagRepo.UpdateTodoTags(ctx, todoID, req.Tags); err != nil {
        tx.Rollback()
        return nil, err
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return todo, nil
}
```

This deep dive provides a comprehensive understanding of how the application's components interact and work together to handle requests and manage data. Each layer has specific responsibilities and follows consistent patterns for error handling, validation, and data management.

## Table of Contents

1. [Introduction](#introduction)
2. [Architecture Overview](#architecture-overview)
3. [Understanding the Layers](#understanding-the-layers)
4. [Complete User Flow Example](#complete-user-flow-example)
5. [Authentication Deep Dive](#authentication-deep-dive)
6. [Todo Management System](#todo-management-system)
7. [Error Handling System](#error-handling-system)
8. [Security Implementation](#security-implementation)
9. [Testing and Debugging](#testing-and-debugging)
10. [Best Practices and Patterns](#best-practices-and-patterns)

## Introduction

Welcome to the application documentation! This guide will help you understand how our application works, from handling HTTP requests to managing data in the database. We'll use real-world examples and analogies to make complex concepts easier to understand.

### What is This Application?

Think of this application as a digital workspace where:

- Users can register and log in (like creating an account on any website)
- Users can create and manage their tasks (like a to-do list)
- Everything is secure and protected (like a bank's online system)

## Architecture Overview

### The Big Picture

Imagine our application as a restaurant:

```
Customer (HTTP Request) → Host (Router) → Waiter (Handler) → Chef (Service) → Storage (Repository) → Fridge (Database)
```

Each part has a specific job:

- Host (Router): Greets customers and directs them to the right table
- Waiter (Handler): Takes orders and serves food
- Chef (Service): Prepares the food following recipes
- Storage (Repository): Manages ingredients and supplies
- Fridge (Database): Where all the ingredients are stored

### Layer Responsibilities

1. **Router Layer (The Host)**

   ```go
   // router/router.go
   func New(h *handler.Handler, authService service.AuthService) http.Handler {
       r := chi.NewRouter()

       // Public routes (no authentication needed)
       r.Group(func(r chi.Router) {
           routes.SetupAuthRoutes(r, h)  // Login, Register
       })

       // Protected routes (need authentication)
       r.Group(func(r chi.Router) {
           r.Use(authmiddleware.Auth(authService))  // Check if user is logged in
           routes.SetupUserRoutes(r, h)    // User profile
           routes.SetupTodoRoutes(r, h)    // Todo management
       })
   }
   ```

   - Like a restaurant host who:
     - Greets customers at the door
     - Checks if they have a reservation (authentication)
     - Directs them to the right section (public or private areas)

2. **Handler Layer (The Waiters)**

   ```go
   // handler/user_handler.go
   type UserHandler struct {
       authService service.AuthService
   }

   func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
       // 1. Take the order (parse request)
       var req model.RegisterRequest
       if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
           // Handle invalid order
           return
       }

       // 2. Check if the order makes sense (validate)
       if errs := validation.ValidateRegisterRequest(&req); errs.HasErrors() {
           // Tell customer what's wrong with their order
           return
       }

       // 3. Send order to kitchen (call service)
       resp, err := h.authService.Register(r.Context(), &req)

       // 4. Serve the food (return response)
       response.NewSuccess(http.StatusCreated, resp).Write(w)
   }
   ```

   - Like waiters who:
     - Take orders from customers
     - Check if the order is valid
     - Send orders to the kitchen
     - Serve the prepared food

3. **Service Layer (The Chefs)**

   ```go
   // service/auth_service.go
   type authService struct {
       userRepo repository.UserRepository
       jwt      *JWTService
   }

   func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
       // 1. Check if we have the ingredients (user exists)
       existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
       if existingUser != nil {
           return nil, model.ErrEmailAlreadyExists
       }

       // 2. Prepare the dish (create user)
       user := &model.User{
           Email:     req.Email,
           Username:  req.Username,
           Password:  hashedPassword,
           FirstName: req.FirstName,
           LastName:  req.LastName,
       }

       // 3. Store the prepared dish (save to database)
       if err := s.userRepo.Create(ctx, user); err != nil {
           return nil, err
       }

       // 4. Create a receipt (generate token)
       token, err := s.jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10))

       return &model.AuthResponse{
           User:  user.ToResponse(),
           Token: token,
       }, nil
   }
   ```

   - Like chefs who:
     - Check if they have the right ingredients
     - Follow recipes to prepare food
     - Store prepared dishes
     - Create receipts for orders

4. **Repository Layer (The Storage Manager)**

   ```go
   // repository/user_repository.go
   type userRepository struct {
       db *gorm.DB
   }

   func (r *userRepository) Create(ctx context.Context, user *model.User) error {
       return r.db.WithContext(ctx).Create(user).Error
   }
   ```

   - Like storage managers who:
     - Keep track of all ingredients
     - Store new ingredients
     - Retrieve ingredients when needed
     - Manage the storage system

## Complete User Flow Example

Let's follow a complete user journey from registration to creating a todo:

### 1. User Registration

```http
# Step 1: User sends registration request
POST /register
Content-Type: application/json

{
    "email": "john@example.com",
    "username": "john",
    "password": "SecurePass123!",
    "firstName": "John",
    "lastName": "Doe"
}

# Step 2: Application validates and creates user
# Step 3: User receives response
{
    "success": true,
    "data": {
        "user": {
            "id": 1,
            "email": "john@example.com",
            "username": "john",
            "firstName": "John",
            "lastName": "Doe",
            "createdAt": "2024-04-07T12:00:00Z"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

### 2. User Login

```http
# Step 1: User sends login request
POST /login
Content-Type: application/json

{
    "email": "john@example.com",
    "password": "SecurePass123!"
}

# Step 2: Application validates credentials
# Step 3: User receives new token
{
    "success": true,
    "data": {
        "user": {
            "id": 1,
            "email": "john@example.com",
            "username": "john"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

### 3. Creating a Todo

```http
# Step 1: User sends todo creation request
POST /todos
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
    "title": "Complete project documentation",
    "description": "Write comprehensive documentation for the application"
}

# Step 2: Application creates todo
# Step 3: User receives created todo
{
    "success": true,
    "data": {
        "id": 1,
        "title": "Complete project documentation",
        "description": "Write comprehensive documentation for the application",
        "completed": false,
        "createdAt": "2024-04-07T12:05:00Z"
    }
}
```

## Authentication Deep Dive

### How Authentication Works

1. **Token Generation**

   ```go
   // service/jwt.go
   func (s *JWTService) GenerateToken(userID string) (string, error) {
       claims := jwt.MapClaims{
           "user_id": userID,
           "exp":     time.Now().Add(time.Hour * 24).Unix(),
       }
       return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secretKey)
   }
   ```

   - Like a secure ID card that:
     - Contains user information
     - Has an expiration date
     - Is digitally signed

2. **Token Validation**

   ```go
   // middleware/auth_middleware.go
   func Auth(authService service.AuthService) func(http.Handler) http.Handler {
       return func(next http.Handler) http.Handler {
           return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
               // 1. Check for ID card (token)
               authHeader := r.Header.Get("Authorization")
               if authHeader == "" {
                   httputil.Error(w, http.StatusUnauthorized, "Authorization header required")
                   return
               }

               // 2. Verify ID card (validate token)
               user, err := authService.GetUserByToken(r.Context(), token)
               if err != nil {
                   httputil.Error(w, http.StatusUnauthorized, "Invalid or expired token")
                   return
               }

               // 3. Allow access (add user to context)
               ctx := context.WithValue(r.Context(), userKey, user)
               next.ServeHTTP(w, r.WithContext(ctx))
           })
       }
   }
   ```

   - Like a security guard who:
     - Checks if you have an ID card
     - Verifies if it's valid and not expired
     - Allows you to enter if everything is okay

## Error Handling System

### Types of Errors

1. **Validation Errors**

   ```json
   {
     "error": "Validation failed",
     "errors": [
       {
         "field": "email",
         "message": "must be a valid email address"
       },
       {
         "field": "password",
         "message": "must be at least 8 characters long"
       }
     ]
   }
   ```

   - Like a form checker who:
     - Points out missing information
     - Highlights incorrect formats
     - Provides clear instructions

2. **Business Logic Errors**

   ```json
   {
     "error": "EMAIL_EXISTS",
     "message": "User with this email already exists"
   }
   ```

   - Like a system that:
     - Checks for duplicate entries
     - Verifies business rules
     - Prevents invalid operations

3. **Authentication Errors**
   ```json
   {
     "error": "UNAUTHORIZED",
     "message": "Invalid or expired token"
   }
   ```
   - Like a security system that:
     - Detects invalid access attempts
     - Expires old credentials
     - Protects sensitive information

## Security Implementation

### 1. Password Security

```go
// Password hashing
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// Password verification
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

- Like a secure vault that:
  - Encrypts sensitive information
  - Uses strong encryption methods
  - Never stores plain text

### 2. Token Security

```go
// Token generation with expiration
claims := jwt.MapClaims{
    "user_id": userID,
    "exp":     time.Now().Add(time.Hour * 24).Unix(),
}
```

- Like a temporary access card that:
  - Expires after a set time
  - Contains minimal necessary information
  - Uses strong encryption

## Testing and Debugging

### 1. Unit Testing

```go
func TestUserHandler_Register(t *testing.T) {
    // Setup
    mockAuthService := &MockAuthService{}
    handler := NewUserHandler(mockAuthService)

    // Test cases
    tests := []struct {
        name     string
        request  *model.RegisterRequest
        wantErr  bool
        wantCode int
    }{
        {
            name: "valid registration",
            request: &model.RegisterRequest{
                Email:     "test@example.com",
                Username:  "testuser",
                Password:  "SecurePass123!",
                FirstName: "Test",
                LastName:  "User",
            },
            wantErr:  false,
            wantCode: http.StatusCreated,
        },
    }

    // Run tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### 2. Integration Testing

```go
func TestUserRegistrationFlow(t *testing.T) {
    // Setup test environment
    db := setupTestDB()
    app := setupTestApp(db)

    // Test registration flow
    req := createRegistrationRequest()
    resp := sendRequest(app, req)

    // Verify results
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    assert.NotEmpty(t, resp.Token)
}
```

## Best Practices and Patterns

### 1. Code Organization

- Keep related code together
- Use clear, descriptive names
- Follow consistent patterns

### 2. Error Handling

- Use specific error types
- Provide helpful error messages
- Handle errors at the appropriate level

### 3. Security

- Validate all input
- Use secure defaults
- Follow security best practices

### 4. Testing

- Write tests for all features
- Use both unit and integration tests
- Test error cases and edge conditions

This enhanced documentation provides a comprehensive understanding of the application's architecture and functionality. It uses real-world analogies and practical examples to make complex concepts more accessible. Would you like me to expand on any particular section or add more examples?
