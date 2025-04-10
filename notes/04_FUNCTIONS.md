# Go Functions: A Comprehensive Guide

## Table of Contents

1. [Basic Function Syntax](#basic-function-syntax)
2. [Function Types](#function-types)
3. [Advanced Function Patterns](#advanced-function-patterns)
4. [Error Handling](#error-handling)
5. [Function Best Practices](#function-best-practices)
6. [Closures and Anonymous Functions](#closures-and-anonymous-functions)
7. [Method Functions](#method-functions)
8. [Functional Programming Patterns](#functional-programming-patterns)

## Basic Function Syntax

### Simple Functions

```go
// Basic function
func greet(name string) string {
    return "Hello, " + name
}

// Multiple parameters
func add(a, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

### Named Return Values

```go
// Named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // Naked return
}

// With documentation
// Calculate returns the sum and product of two numbers
func Calculate(a, b int) (sum, product int) {
    sum = a + b
    product = a * b
    return
}
```

## Function Types

### Variadic Functions

```go
// Variable number of arguments
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}

// Usage
sum(1, 2, 3)
numbers := []int{1, 2, 3}
sum(numbers...)
```

### First-Class Functions

```go
// Function as type
type Operation func(a, b int) int

// Function as parameter
func process(op Operation, x, y int) int {
    return op(x, y)
}

// Function returning function
func multiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}
```

## Advanced Function Patterns

### Decorator Pattern

```go
// Function decorator
func withLogging(fn func() error) func() error {
    return func() error {
        log.Println("Starting function")
        err := fn()
        log.Println("Function completed")
        return err
    }
}

// Usage
func someOperation() error {
    // Implementation
    return nil
}

decorated := withLogging(someOperation)
```

### Options Pattern

```go
type Server struct {
    host string
    port int
    timeout time.Duration
}

type ServerOption func(*Server)

func WithHost(host string) ServerOption {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(options ...ServerOption) *Server {
    server := &Server{
        host: "localhost",
        port: 8080,
    }

    for _, option := range options {
        option(server)
    }

    return server
}
```

## Error Handling

### Error Wrapping

```go
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // Process file
    return nil
}

// Custom error types
type ValidationError struct {
    Field string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

## Function Best Practices

### Parameter Validation

```go
func processUser(user *User) error {
    // Validate parameters
    if user == nil {
        return errors.New("user cannot be nil")
    }

    if user.Name == "" {
        return &ValidationError{
            Field: "name",
            Message: "cannot be empty",
        }
    }

    return nil
}
```

### Resource Management

```go
func processWithResources() error {
    // Acquire resources
    resource, err := acquireResource()
    if err != nil {
        return fmt.Errorf("failed to acquire resource: %w", err)
    }
    defer resource.Release()

    // Use resource
    return nil
}
```

## Closures and Anonymous Functions

### Basic Closure

```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// Usage
increment := counter()
fmt.Println(increment()) // 1
fmt.Println(increment()) // 2
```

### Goroutine Closure

```go
func processItems(items []string) {
    for _, item := range items {
        // Create new variable in loop
        item := item
        go func() {
            fmt.Println(item)
        }()
    }
}
```

## Method Functions

### Receiver Methods

```go
type Rectangle struct {
    Width, Height float64
}

// Value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```

### Interface Implementation

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}
```

## Functional Programming Patterns

### Map Function

```go
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, item := range slice {
        result[i] = f(item)
    }
    return result
}

// Usage
numbers := []int{1, 2, 3}
squares := Map(numbers, func(x int) int {
    return x * x
})
```

### Reduce Function

```go
func Reduce[T, U any](slice []T, initial U, f func(U, T) U) U {
    result := initial
    for _, item := range slice {
        result = f(result, item)
    }
    return result
}

// Usage
sum := Reduce(numbers, 0, func(acc, x int) int {
    return acc + x
})
```

### Chain Pattern

```go
type Chain struct {
    value int
}

func (c *Chain) Add(x int) *Chain {
    c.value += x
    return c
}

func (c *Chain) Multiply(x int) *Chain {
    c.value *= x
    return c
}

// Usage
result := new(Chain).
    Add(5).
    Multiply(2).
    Add(10).
    value
```

Remember:

- Keep functions focused and small
- Use meaningful parameter and return value names
- Handle errors appropriately
- Use defer for cleanup
- Document public functions
- Use appropriate receiver types for methods
- Consider using functional options for complex configurations
- Leverage closures when appropriate
- Use generics for type-safe generic functions
- Follow consistent error handling patterns
- Use method chaining when it improves readability
- Consider performance implications of your function patterns

This guide covers the essential aspects of functions in Go. Understanding these patterns and best practices is crucial for writing clean, maintainable, and efficient Go code.
