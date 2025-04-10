# OOP: Go vs JavaScript - A Comparative Guide

## Table of Contents

- [OOP: Go vs JavaScript - A Comparative Guide](#oop-go-vs-javascript---a-comparative-guide)
  - [Table of Contents](#table-of-contents)
  - [Class Definition](#class-definition)
    - [JavaScript](#javascript)
    - [Go](#go)
  - [Constructor Patterns](#constructor-patterns)
    - [JavaScript](#javascript-1)
    - [Go](#go-1)
  - [Methods and Functions](#methods-and-functions)
    - [JavaScript](#javascript-2)
    - [Go](#go-2)
  - [Inheritance vs Composition](#inheritance-vs-composition)
    - [JavaScript](#javascript-3)
    - [Go](#go-3)
  - [Interfaces and Abstract Classes](#interfaces-and-abstract-classes)
    - [JavaScript](#javascript-4)
    - [Go](#go-4)
  - [Encapsulation](#encapsulation)
    - [JavaScript](#javascript-5)
    - [Go](#go-5)
  - [Error Handling](#error-handling)
    - [JavaScript](#javascript-6)
    - [Go](#go-6)
  - [Working with Properties](#working-with-properties)
    - [JavaScript](#javascript-7)
    - [Go](#go-7)
  - [Best Practices](#best-practices)
    - [JavaScript](#javascript-8)
    - [Go](#go-8)

## Class Definition

### JavaScript

```javascript
class User {
  constructor(firstName, lastName, email) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.email = email;
  }
}
```

### Go

```go
type User struct {
    FirstName string
    LastName  string
    Email     string
}
```

**Key Differences:**

- JavaScript uses the `class` keyword
- Go uses `struct` for object definitions
- No explicit constructor keyword in Go
- Go's struct fields must have type declarations

## Constructor Patterns

### JavaScript

```javascript
class User {
  constructor(firstName, lastName, email) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.email = email;
  }

  // Factory method
  static create(data) {
    return new User(data.firstName, data.lastName, data.email);
  }
}
```

### Go

```go
// Constructor function
func NewUser(firstName, lastName, email string) *User {
    return &User{
        FirstName: firstName,
        LastName:  lastName,
        Email:     email,
    }
}
```

**Key Differences:**

- JavaScript uses the `constructor` method
- Go uses constructor functions by convention (usually prefixed with `New`)
- Go returns pointers to structs
- JavaScript constructors are part of the class definition

## Methods and Functions

### JavaScript

```javascript
class User {
  getFullName() {
    return `${this.firstName} ${this.lastName}`;
  }

  updateEmail(newEmail) {
    this.email = newEmail;
  }
}
```

### Go

```go
// Value receiver
func (u User) GetFullName() string {
    return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// Pointer receiver
func (u *User) UpdateEmail(newEmail string) {
    u.Email = newEmail
}
```

**Key Differences:**

- Go uses explicit receiver syntax
- Go distinguishes between value and pointer receivers
- JavaScript methods are always defined within the class
- Go methods are defined outside the struct

## Inheritance vs Composition

### JavaScript

```javascript
class Entity {
  constructor() {
    this.id = null;
    this.createdAt = new Date();
  }
}

class User extends Entity {
  constructor(firstName, lastName) {
    super();
    this.firstName = firstName;
    this.lastName = lastName;
  }
}
```

### Go

```go
type Entity struct {
    ID        int
    CreatedAt time.Time
}

type User struct {
    Entity            // Embedding
    FirstName string
    LastName  string
}
```

**Key Differences:**

- JavaScript uses classical inheritance with `extends`
- Go uses composition through embedding
- No `super` calls in Go
- Go's approach is more flexible and less hierarchical

## Interfaces and Abstract Classes

### JavaScript

```javascript
// Abstract class
class Repository {
  create() {
    throw new Error("Not implemented");
  }
  find() {
    throw new Error("Not implemented");
  }
}

// Implementation
class UserRepository extends Repository {
  create(user) {
    // Implementation
  }
  find(id) {
    // Implementation
  }
}
```

### Go

```go
type Repository interface {
    Create(interface{}) error
    Find(id int) (interface{}, error)
}

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Create(data interface{}) error {
    // Implementation
    return nil
}

func (r *UserRepository) Find(id int) (interface{}, error) {
    // Implementation
    return nil, nil
}
```

**Key Differences:**

- Go interfaces are implicit
- JavaScript needs explicit inheritance
- Go interfaces are typically smaller
- No abstract classes in Go

## Encapsulation

### JavaScript

```javascript
class User {
  #privateField = "private";

  getPrivateField() {
    return this.#privateField;
  }
}
```

### Go

```go
type User struct {
    publicField  string // Exported
    privateField string // Unexported
}

func (u *User) GetPrivateField() string {
    return u.privateField
}
```

**Key Differences:**

- Go uses capitalization for access control
- JavaScript uses `#` for private fields
- Go's encapsulation is at package level
- JavaScript's private fields are truly private

## Error Handling

### JavaScript

```javascript
class UserService {
  async createUser(userData) {
    try {
      if (!userData.email) {
        throw new Error("Email required");
      }
      return await this.repository.save(userData);
    } catch (error) {
      throw new Error(`Failed to create user: ${error.message}`);
    }
  }
}
```

### Go

```go
type UserService struct {
    repository Repository
}

func (s *UserService) CreateUser(userData *User) error {
    if userData.Email == "" {
        return fmt.Errorf("email required")
    }

    if err := s.repository.Save(userData); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}
```

**Key Differences:**

- Go uses explicit error returns
- JavaScript uses try/catch blocks
- Go errors are values
- JavaScript errors are exceptions

## Working with Properties

### JavaScript

```javascript
class User {
  set name(value) {
    this._name = value.trim();
  }

  get name() {
    return this._name;
  }
}
```

### Go

```go
type User struct {
    name string
}

func (u *User) SetName(name string) {
    u.name = strings.TrimSpace(name)
}

func (u *User) Name() string {
    return u.name
}
```

**Key Differences:**

- JavaScript has built-in getters/setters
- Go uses explicit methods
- Go naming conventions for getters/setters differ
- No direct property accessors in Go

## Best Practices

### JavaScript

```javascript
// Dependency Injection
class UserService {
  constructor(repository, logger) {
    this.repository = repository;
    this.logger = logger;
  }
}

// Method Chaining
class QueryBuilder {
  where(condition) {
    this.condition = condition;
    return this;
  }
}
```

### Go

```go
// Dependency Injection
type UserService struct {
    repository Repository
    logger     Logger
}

func NewUserService(repo Repository, logger Logger) *UserService {
    return &UserService{
        repository: repo,
        logger:    logger,
    }
}

// Method Chaining
type QueryBuilder struct {
    condition string
}

func (q *QueryBuilder) Where(condition string) *QueryBuilder {
    q.condition = condition
    return q
}
```

**Key Differences:**

- Go emphasizes explicit dependencies
- JavaScript often uses class properties
- Go returns pointers for method chaining
- JavaScript returns `this` for chaining

Remember:

- Go promotes composition over inheritance
- JavaScript supports both classical inheritance and composition
- Go's approach is more explicit and verbose
- JavaScript offers more flexibility but can be less predictable
- Go's type system is stricter and provides better compile-time safety
- JavaScript's dynamic nature allows for more flexible but potentially less maintainable code

This comparison shows how both languages achieve similar OOP goals through different approaches, with Go focusing on simplicity and explicitness while JavaScript provides more traditional OOP features with greater flexibility.
