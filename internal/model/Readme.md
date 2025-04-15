# GORM Model Guide

This guide provides everything you need to know about writing models using GORM in Go applications. It covers from basic to advanced concepts, with practical examples and best practices.

## Table of Contents

- [GORM Model Guide](#gorm-model-guide)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Basic Setup](#basic-setup)
    - [Configuration](#configuration)
  - [Basic Model Structure](#basic-model-structure)
  - [Primary Keys](#primary-keys)
    - [Auto-incrementing Primary Key](#auto-incrementing-primary-key)
    - [Custom Primary Key](#custom-primary-key)
    - [Composite Primary Key](#composite-primary-key)
  - [Relationships](#relationships)
    - [One-to-One](#one-to-one)
    - [One-to-Many](#one-to-many)
    - [Many-to-Many](#many-to-many)
  - [Query Building](#query-building)
    - [Basic Queries](#basic-queries)
    - [Advanced Queries](#advanced-queries)
    - [Joins](#joins)
  - [Transactions](#transactions)
    - [Basic Transaction](#basic-transaction)
    - [Transaction with Context](#transaction-with-context)
    - [Nested Transactions](#nested-transactions)
  - [Hooks](#hooks)
    - [Available Hooks](#available-hooks)
    - [Example Hooks](#example-hooks)
  - [Scopes](#scopes)
    - [Defining Scopes](#defining-scopes)
    - [Using Scopes](#using-scopes)
  - [Raw SQL](#raw-sql)
    - [Executing Raw SQL](#executing-raw-sql)
    - [Named Arguments](#named-arguments)
  - [Performance Optimization](#performance-optimization)
    - [Indexes](#indexes)
    - [Prepared Statements](#prepared-statements)
    - [Batch Operations](#batch-operations)
  - [Testing](#testing)
    - [Setup Test Database](#setup-test-database)
    - [Test Examples](#test-examples)
  - [Common Pitfalls](#common-pitfalls)
  - [Example Models](#example-models)
    - [Complete User Model](#complete-user-model)
    - [Complete Todo Model](#complete-todo-model)
  - [Request Handling and Validation](#request-handling-and-validation)
    - [Request Flow](#request-flow)
    - [Request Types](#request-types)
      - [1. Registration Request](#1-registration-request)
      - [2. Login Request](#2-login-request)
      - [3. Todo Request](#3-todo-request)
    - [Validation Process](#validation-process)
      - [1. Request Parsing](#1-request-parsing)
      - [2. Validation Rules](#2-validation-rules)
      - [3. Error Handling](#3-error-handling)
      - [4. Validation Middleware](#4-validation-middleware)
    - [Complete Request Flow Example](#complete-request-flow-example)
    - [Error Responses](#error-responses)
      - [1. Validation Errors](#1-validation-errors)
      - [2. Business Logic Errors](#2-business-logic-errors)
      - [3. Internal Server Error](#3-internal-server-error)
    - [Best Practices for Request Handling](#best-practices-for-request-handling)
    - [Testing Request Validation](#testing-request-validation)

## Getting Started

### Installation

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres  # or mysql, sqlite, etc.
```

### Basic Setup

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

func main() {
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
}
```

### Configuration

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // Enable logging
    NamingStrategy: schema.NamingStrategy{
        TablePrefix: "t_",   // table name prefix
        SingularTable: true, // use singular table name
    },
    PrepareStmt: true, // Enable prepared statements
})
```

## Basic Model Structure

Every model should follow this basic structure:

```go
type ModelName struct {
    // Primary key
    ID        uint           `gorm:"primaryKey" json:"id"`

    // Required fields
    FieldName string         `gorm:"not null" json:"field_name"`

    // Optional fields
    OptionalField string     `json:"optional_field"`

    // Timestamps
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## Primary Keys

### Auto-incrementing Primary Key

```go
type User struct {
    ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
    // ... other fields
}
```

### Custom Primary Key

```go
type User struct {
    UserID string `gorm:"primaryKey;type:uuid" json:"user_id"`
    // ... other fields
}
```

### Composite Primary Key

```go
type UserRole struct {
    UserID uint `gorm:"primaryKey" json:"user_id"`
    RoleID uint `gorm:"primaryKey" json:"role_id"`
    // ... other fields
}
```

## Relationships

### One-to-One

```go
type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Profile   Profile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"profile"`
}

type Profile struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    UserID uint   `json:"user_id"`
    User   User   `gorm:"foreignKey:UserID" json:"user"`
}
```

### One-to-Many

One-to-Many relationships are very common in database designs. This relationship exists when a single record in one table is associated with multiple records in another table. Examples include:
- A user having multiple posts
- A department having multiple employees
- A product having multiple reviews

#### Basic Implementation

```go
type User struct {
    ID    uint    `gorm:"primaryKey" json:"id"`
    Name  string  `json:"name"`
    Posts []Post  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"posts"`
}

type Post struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    UserID    uint      `json:"user_id"`
    User      User      `gorm:"foreignKey:UserID" json:"user"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### Tag Options Explained

- `foreignKey`: Specifies which field is used as the foreign key (default is the owner's type name + primary key)
- `references`: Specifies which field the foreign key references (default is the primary key)
- `constraint`: Defines referential actions (CASCADE, SET NULL, RESTRICT, etc.)

#### Complete Example with Operations

```go
// Define models
type Author struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    Name      string     `json:"name"`
    Email     string     `gorm:"uniqueIndex" json:"email"`
    Books     []Book     `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"books"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}

type Book struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    PublishedAt time.Time `json:"published_at"`
    AuthorID    *uint     `json:"author_id"` // Nullable foreign key
    Author      Author    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Creating records with association
func CreateAuthorWithBooks(db *gorm.DB) error {
    author := Author{
        Name:  "John Doe",
        Email: "john@example.com",
        Books: []Book{
            {Title: "GORM Basics", Description: "Introduction to GORM", PublishedAt: time.Now()},
            {Title: "Advanced GORM", Description: "Deep dive into GORM", PublishedAt: time.Now()},
        },
    }
    
    return db.Create(&author).Error
}

// Querying with preload
func GetAuthorWithBooks(db *gorm.DB, authorID uint) (*Author, error) {
    var author Author
    if err := db.Preload("Books").First(&author, authorID).Error; err != nil {
        return nil, err
    }
    return &author, nil
}

// Adding a book to an existing author
func AddBookToAuthor(db *gorm.DB, authorID uint, book Book) error {
    return db.Model(&Author{ID: authorID}).Association("Books").Append(&book)
}

// Count books for an author
func CountAuthorBooks(db *gorm.DB, authorID uint) (int64, error) {
    var count int64
    err := db.Model(&Book{}).Where("author_id = ?", authorID).Count(&count).Error
    return count, err
}

// Find all books for an author without loading the author
func GetAuthorBooks(db *gorm.DB, authorID uint) ([]Book, error) {
    var books []Book
    err := db.Where("author_id = ?", authorID).Find(&books).Error
    return books, err
}
```

#### Self-Referential One-to-Many

A common use case is hierarchical data like comments with replies:

```go
type Comment struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Content   string    `json:"content"`
    ParentID  *uint     `json:"parent_id"`
    Parent    *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### Best Practices

1. **Always handle the error returned by GORM**:
   ```go
   if err := db.Preload("Books").First(&author, authorID).Error; err != nil {
       return nil, err
   }
   ```

2. **Use Preload judiciously** to avoid N+1 queries:
   ```go
   // Good - Single query with preload
   db.Preload("Books").Find(&authors)
   
   // Bad - N+1 queries
   db.Find(&authors)
   for _, author := range authors {
       db.Model(&author).Association("Books").Find(&author.Books)
   }
   ```

3. **Consider nullable foreign keys** when appropriate:
   ```go
   AuthorID *uint `json:"author_id"` // Can be NULL
   ```

4. **Use association mode for complex operations**:
   ```go
   // Replace all books
   db.Model(&author).Association("Books").Replace(&newBooks)
   
   // Delete specific association
   db.Model(&author).Association("Books").Delete(&bookToDelete)
   
   // Clear all associations
   db.Model(&author).Association("Books").Clear()
   ```

5. **Set appropriate ON DELETE constraints**:
   ```go
   // When parent is deleted, set children to NULL
   `gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL"`
   
   // When parent is deleted, delete all children
   `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
   ```

6. **For large collections, consider using limits and pagination**:
   ```go
   db.Preload("Books", func(db *gorm.DB) *gorm.DB {
       return db.Order("published_at DESC").Limit(10)
   }).Find(&authors)
   ```

### Many-to-Many

```go
type User struct {
    ID    uint      `gorm:"primaryKey" json:"id"`
    Roles []Role    `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles"`
}

type Role struct {
    ID    uint     `gorm:"primaryKey" json:"id"`
    Users []User   `gorm:"many2many:user_roles;joinForeignKey:RoleID;joinReferences:UserID" json:"users"`
}
```

## Query Building

### Basic Queries

```go
// Find by primary key
db.First(&user, 1)

// Find by conditions
db.Where("name = ?", "jinzhu").First(&user)
db.Where(&User{Name: "jinzhu"}).First(&user)

// Find all records
db.Find(&users)

// Find with conditions
db.Where("name <> ?", "jinzhu").Find(&users)
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
```

### Advanced Queries

```go
// Select specific fields
db.Select("name", "age").Find(&users)

// Order by
db.Order("age desc, name").Find(&users)

// Limit and Offset
db.Limit(10).Offset(5).Find(&users)

// Group by
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Find(&results)

// Having
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("sum(age) > ?", 100).Find(&results)
```

### Joins

```go
// Inner Join
db.Joins("Profile").Find(&users)

// Left Join
db.Joins("LEFT JOIN profiles ON profiles.user_id = users.id").Find(&users)

// Multiple Joins
db.Joins("Profile").Joins("Company").Find(&users)
```

## Transactions

### Basic Transaction

```go
tx := db.Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
return tx.Commit().Error
```

### Transaction with Context

```go
tx := db.WithContext(ctx).Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
return tx.Commit().Error
```

### Nested Transactions

```go
tx := db.Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    return nil
}); err != nil {
    tx.Rollback()
    return err
}

return tx.Commit().Error
```

## Hooks

### Available Hooks

```go
// Creating
func (u *User) BeforeCreate(tx *gorm.DB) error
func (u *User) AfterCreate(tx *gorm.DB) error

// Updating
func (u *User) BeforeUpdate(tx *gorm.DB) error
func (u *User) AfterUpdate(tx *gorm.DB) error

// Deleting
func (u *User) BeforeDelete(tx *gorm.DB) error
func (u *User) AfterDelete(tx *gorm.DB) error

// Finding
func (u *User) AfterFind(tx *gorm.DB) error
```

### Example Hooks

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Hash password
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email
    return nil
}
```

## Scopes

### Defining Scopes

```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func OlderThan(age int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age > ?", age)
    }
}
```

### Using Scopes

```go
db.Scopes(ActiveUsers, OlderThan(18)).Find(&users)
```

## Raw SQL

### Executing Raw SQL

```go
db.Exec("UPDATE users SET name = ? WHERE id = ?", "jinzhu", 1)

var result Result
db.Raw("SELECT name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)
```

### Named Arguments

```go
db.Raw("SELECT * FROM users WHERE name = @name OR age = @age",
    map[string]interface{}{"name": "jinzhu", "age": 18}).Find(&users)
```

## Performance Optimization

### Indexes

```go
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"index:idx_name"`
    Email     string `gorm:"uniqueIndex"`
    Age       int    `gorm:"index:idx_age"`
}
```

### Prepared Statements

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    PrepareStmt: true,
})
```

### Batch Operations

```go
// Batch Insert
var users = []User{{Name: "jinzhu_1"}, {Name: "jinzhu_2"}}
db.Create(&users)

// Batch Update
db.Model(&User{}).Where("role = ?", "admin").Updates(User{Name: "jinzhu"})
```

## Testing

### Setup Test Database

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatal(err)
    }

    // Migrate the schema
    db.AutoMigrate(&User{})

    return db
}
```

### Test Examples

```go
func TestUserCreate(t *testing.T) {
    db := setupTestDB(t)

    user := User{Name: "jinzhu"}
    if err := db.Create(&user).Error; err != nil {
        t.Fatal(err)
    }

    var found User
    if err := db.First(&found, user.ID).Error; err != nil {
        t.Fatal(err)
    }

    if found.Name != user.Name {
        t.Errorf("Expected name %s, got %s", user.Name, found.Name)
    }
}
```

## Common Pitfalls

1. **N+1 Query Problem**

   ```go
   // Bad
   var users []User
   db.Find(&users)
   for _, user := range users {
       db.Model(&user).Association("Posts").Find(&user.Posts)
   }

   // Good
   var users []User
   db.Preload("Posts").Find(&users)
   ```

2. **Missing Error Handling**

   ```go
   // Bad
   db.Create(&user)

   // Good
   if err := db.Create(&user).Error; err != nil {
       return err
   }
   ```

3. **Incorrect Transaction Usage**

   ```go
   // Bad
   tx := db.Begin()
   tx.Create(&user)
   tx.Commit()

   // Good
   tx := db.Begin()
   if err := tx.Create(&user).Error; err != nil {
       tx.Rollback()
       return err
   }
   return tx.Commit().Error
   ```

## Example Models

### Complete User Model

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;not null" json:"email"`
    Username  string         `gorm:"uniqueIndex;not null" json:"username"`
    Password  string         `gorm:"not null" json:"-"`
    FirstName string         `json:"first_name"`
    LastName  string         `json:"last_name"`
    Profile   Profile        `gorm:"foreignKey:UserID" json:"profile"`
    Posts     []Post         `gorm:"foreignKey:UserID" json:"posts"`
    Roles     []Role         `gorm:"many2many:user_roles;" json:"roles"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email
    return nil
}
```

### Complete Todo Model

```go
type Todo struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Title       string         `gorm:"not null;index" json:"title"`
    Description string         `json:"description"`
    Completed   bool           `gorm:"default:false" json:"completed"`
    UserID      uint           `gorm:"not null;index" json:"user_id"`
    User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Tags        []Tag          `gorm:"many2many:todo_tags;" json:"tags"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Todo) BeforeCreate(tx *gorm.DB) error {
    if t.Title == "" {
        return errors.New("title is required")
    }
    return nil
}

func (t *Todo) AfterUpdate(tx *gorm.DB) error {
    if t.Completed {
        // Send notification
    }
    return nil
}
```

## Request Handling and Validation

### Request Flow

1. **HTTP Request** → **Handler** → **Service** → **Repository** → **Database**
2. Each layer has specific responsibilities:
   - Handler: Request parsing, response formatting
   - Service: Business logic, validation
   - Repository: Database operations
   - Model: Data structure, validation rules

### Request Types

#### 1. Registration Request

```go
type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email,max=100"`
    Username  string `json:"username" validate:"required,username"`
    Password  string `json:"password" validate:"required,password"`
    FirstName string `json:"first_name" validate:"required,name"`
    LastName  string `json:"last_name" validate:"required,name"`
}

// Validation Rules:
// - Email: Required, valid email format, max 100 chars
// - Username: Required, 3-20 chars, alphanumeric + underscore
// - Password: Required, min 8 chars, contains uppercase, lowercase, number, special char
// - FirstName/LastName: Required, 2-50 chars, letters, spaces, hyphens, apostrophes
```

#### 2. Login Request

```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email,max=100"`
    Password string `json:"password" validate:"required,password"`
}

// Validation Rules:
// - Email: Required, valid email format, max 100 chars
// - Password: Required, min 8 chars
```

#### 3. Todo Request

```go
type TodoCreateRequest struct {
    Title       string `json:"title" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"max=500"`
}

type TodoUpdateRequest struct {
    Title       string `json:"title" validate:"omitempty,min=3,max=100"`
    Description string `json:"description" validate:"omitempty,max=500"`
    Completed   bool   `json:"completed"`
}

// Validation Rules:
// - Title: Required for create, 3-100 chars
// - Description: Optional, max 500 chars
// - Completed: Boolean, no validation needed
```

### Validation Process

#### 1. Request Parsing

```go
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httputil.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    // ... validation and processing
}
```

#### 2. Validation Rules

```go
// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    return len(password) >= 8 &&
        regexp.MustCompile(`[A-Z]`).MatchString(password) &&
        regexp.MustCompile(`[a-z]`).MatchString(password) &&
        regexp.MustCompile(`[0-9]`).MatchString(password) &&
        regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    return len(username) >= 3 && len(username) <= 20 &&
        regexp.MustCompile(`^[a-zA-Z]`).MatchString(username) &&
        regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username)
}

func validateName(fl validator.FieldLevel) bool {
    name := fl.Field().String()
    return len(name) >= 2 && len(name) <= 50 &&
        regexp.MustCompile(`^[a-zA-Z\s\-']+$`).MatchString(name) &&
        !regexp.MustCompile(`[\s\-']{2,}`).MatchString(name) &&
        !regexp.MustCompile(`^[\s\-']|[\s\-']$`).MatchString(name)
}
```

#### 3. Error Handling

```go
type ValidationError struct {
    Field      string `json:"field"`
    Message    string `json:"message"`
    StatusCode int    `json:"-"`
}

type ValidationErrors struct {
    Errors     []ValidationError `json:"errors"`
    StatusCode int              `json:"-"`
}

func (ve ValidationErrors) Error() string {
    var messages []string
    for _, err := range ve.Errors {
        messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
    }
    return strings.Join(messages, ", ")
}
```

#### 4. Validation Middleware

```go
func ValidateRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var req interface{}
        // Determine request type based on route
        switch r.URL.Path {
        case "/register":
            req = &RegisterRequest{}
        case "/login":
            req = &LoginRequest{}
        // ... other cases
        }

        if err := json.NewDecoder(r.Body).Decode(req); err != nil {
            httputil.Error(w, http.StatusBadRequest, "Invalid request body")
            return
        }

        if err := validation.ValidateStruct(req); err != nil {
            if ve, ok := err.(validation.ValidationErrors); ok {
                httputil.JSON(w, ve.StatusCode, ve)
                return
            }
            httputil.Error(w, http.StatusBadRequest, "Invalid request data")
            return
        }

        // Store validated request in context
        ctx := context.WithValue(r.Context(), "validated_request", req)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Complete Request Flow Example

```go
// 1. Handler
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Get validated request from context
    req, ok := r.Context().Value("validated_request").(*RegisterRequest)
    if !ok {
        httputil.Error(w, http.StatusBadRequest, "Invalid request")
        return
    }

    // Call service
    authResponse, err := h.authService.Register(r.Context(), req)
    if err != nil {
        switch err {
        case service.ErrUserAlreadyExists:
            httputil.Error(w, http.StatusConflict, "User already exists")
        default:
            httputil.Error(w, http.StatusInternalServerError, "Internal server error")
        }
        return
    }

    // Return response
    httputil.JSON(w, http.StatusCreated, authResponse)
}

// 2. Service
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
    // Check if user exists
    existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, ErrUserAlreadyExists
    }

    // Create user
    user := &model.User{
        Email:     req.Email,
        Username:  req.Username,
        Password:  req.Password, // Will be hashed in BeforeCreate hook
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Generate token
    token, err := s.jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10))
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        User:  user.ToResponse(),
        Token: token,
    }, nil
}

// 3. Repository
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

### Error Responses

#### 1. Validation Errors

```json
{
  "errors": [
    {
      "field": "email",
      "message": "email must be a valid email address"
    },
    {
      "field": "password",
      "message": "password must be at least 8 characters long and contain uppercase, lowercase, number, and special character"
    }
  ]
}
```

#### 2. Business Logic Errors

```json
{
  "error": "User already exists"
}
```

#### 3. Internal Server Error

```json
{
  "error": "Internal server error"
}
```

### Best Practices for Request Handling

1. **Input Validation**

   - Validate at the earliest point possible
   - Use custom validation functions for complex rules
   - Provide clear error messages
   - Return appropriate HTTP status codes

2. **Error Handling**

   - Use custom error types
   - Handle errors at appropriate levels
   - Log errors with context
   - Return user-friendly error messages

3. **Security**

   - Sanitize all input
   - Use parameterized queries
   - Implement rate limiting
   - Validate content types

4. **Performance**

   - Use prepared statements
   - Implement connection pooling
   - Cache when appropriate
   - Use appropriate indexes

5. **Testing**
   - Test all validation rules
   - Test error scenarios
   - Test edge cases
   - Use mock repositories

### Testing Request Validation

```go
func TestRegisterValidation(t *testing.T) {
    tests := []struct {
        name    string
        req     RegisterRequest
        wantErr bool
    }{
        {
            name: "valid request",
            req: RegisterRequest{
                Email:     "test@example.com",
                Username:  "testuser",
                Password:  "Test123!@#",
                FirstName: "Test",
                LastName:  "User",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            req: RegisterRequest{
                Email:     "invalid-email",
                Username:  "testuser",
                Password:  "Test123!@#",
                FirstName: "Test",
                LastName:  "User",
            },
            wantErr: true,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validation.ValidateStruct(&tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

This comprehensive request handling and validation system ensures:

1. Data integrity
2. Security
3. Performance
4. Maintainability
5. Testability
