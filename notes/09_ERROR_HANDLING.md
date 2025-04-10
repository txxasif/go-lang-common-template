# Go Error Handling: A Comprehensive Guide

## Table of Contents

1. [Error Interface](#error-interface)
2. [Basic Error Handling](#basic-error-handling)
3. [Custom Errors](#custom-errors)
4. [Error Wrapping](#error-wrapping)
5. [Panic and Recover](#panic-and-recover)
6. [Error Handling Patterns](#error-handling-patterns)
7. [Best Practices](#best-practices)
8. [Advanced Concepts](#advanced-concepts)

## Error Interface

### Understanding the Error Interface

In Go, errors are values that implement the built-in `error` interface:

```go
type error interface {
    Error() string
}
```

This simple interface means any type that has an `Error() string` method can be used as an error.

### Why This Design?

1. **Simplicity**: The interface is simple and clear
2. **Flexibility**: Any type can be an error
3. **Composability**: Errors can be composed and wrapped
4. **No Exceptions**: Go uses explicit error handling instead of exceptions

## Basic Error Handling

### The Standard Pattern

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

### Key Concepts:

1. **Multiple Return Values**: Functions return both result and error
2. **Nil Error**: No error is represented by `nil`
3. **Error Checking**: Always check errors immediately
4. **Error Propagation**: Return errors up the call stack

### Common Error Creation

```go
// Using errors.New
err := errors.New("something went wrong")

// Using fmt.Errorf
err := fmt.Errorf("invalid value: %v", value)
```

## Custom Errors

### Creating Custom Error Types

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

### When to Use Custom Errors

1. When you need to carry additional information
2. When you need to distinguish between different error types
3. When you need to implement specific error handling logic

### Error Type Assertion

```go
func handleError(err error) {
    if ve, ok := err.(*ValidationError); ok {
        // Handle validation error
        fmt.Printf("Validation error in field %s: %s\n", ve.Field, ve.Message)
    } else {
        // Handle other errors
        fmt.Println("Other error:", err)
    }
}
```

## Error Wrapping

### Understanding Error Wrapping

Error wrapping allows you to add context to errors while preserving the original error.

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
```

### Unwrapping Errors

```go
func handleError(err error) {
    if err == nil {
        return
    }

    // Unwrap the error
    if unwrapped := errors.Unwrap(err); unwrapped != nil {
        fmt.Println("Original error:", unwrapped)
    }

    // Check if error is of specific type
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File does not exist")
    }
}
```

## Panic and Recover

### Understanding Panic

Panic is used for unrecoverable errors that should stop program execution.

```go
func mustDivide(a, b float64) float64 {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}
```

### Using Recover

Recover is used to catch panics and handle them gracefully.

```go
func safeDivide(a, b float64) (result float64, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic occurred: %v", r)
        }
    }()

    result = mustDivide(a, b)
    return
}
```

### When to Use Panic

1. Only for truly unrecoverable errors
2. When continuing execution would cause more problems
3. In initialization code that must succeed
4. Never in library code that others will use

## Error Handling Patterns

### 1. Error Wrapping Pattern

```go
func processData(data []byte) error {
    if err := validateData(data); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    if err := storeData(data); err != nil {
        return fmt.Errorf("storage failed: %w", err)
    }

    return nil
}
```

### 2. Error Chain Pattern

```go
type ErrorChain struct {
    err error
    msg string
}

func (e *ErrorChain) Error() string {
    return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *ErrorChain) Unwrap() error {
    return e.err
}
```

### 3. Error Aggregation Pattern

```go
type MultiError struct {
    errors []error
}

func (e *MultiError) Error() string {
    var msgs []string
    for _, err := range e.errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (e *MultiError) Add(err error) {
    if err != nil {
        e.errors = append(e.errors, err)
    }
}
```

## Best Practices

### 1. Always Handle Errors

```go
// Bad
file, _ := os.Open("file.txt")

// Good
file, err := os.Open("file.txt")
if err != nil {
    return err
}
```

### 2. Add Context to Errors

```go
// Bad
return err

// Good
return fmt.Errorf("failed to process user %s: %w", userID, err)
```

### 3. Use Custom Errors When Appropriate

```go
// Bad
return fmt.Errorf("validation failed")

// Good
return &ValidationError{
    Field: "email",
    Message: "invalid format",
}
```

### 4. Don't Ignore Errors

```go
// Bad
_ = file.Close()

// Good
if err := file.Close(); err != nil {
    log.Printf("failed to close file: %v", err)
}
```

## Advanced Concepts

### 1. Error Inspection

```go
func inspectError(err error) {
    for err != nil {
        fmt.Printf("Error: %v\n", err)
        err = errors.Unwrap(err)
    }
}
```

### 2. Error Classification

```go
type ErrorKind int

const (
    KindInvalid ErrorKind = iota
    KindNotFound
    KindPermission
)

type ClassifiedError struct {
    kind ErrorKind
    err  error
}

func (e *ClassifiedError) Error() string {
    return e.err.Error()
}

func (e *ClassifiedError) Kind() ErrorKind {
    return e.kind
}
```

### 3. Error Recovery Strategy

```go
type RecoveryStrategy func(error) error

func withRecovery(strategy RecoveryStrategy, f func() error) error {
    err := f()
    if err == nil {
        return nil
    }
    return strategy(err)
}
```

### 4. Error Metrics

```go
type ErrorMetrics struct {
    Counts map[string]int
    mu     sync.Mutex
}

func (m *ErrorMetrics) Record(err error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.Counts[err.Error()]++
}
```

Remember:

- Errors are values, not exceptions
- Always check errors immediately
- Add context to errors when propagating
- Use custom errors for specific error types
- Use panic only for truly unrecoverable errors
- Document error conditions in function signatures
- Consider error handling in concurrent code
- Use error wrapping to preserve error context
- Implement proper error cleanup
- Consider error metrics and monitoring

This guide covers the fundamental and advanced aspects of error handling in Go. Understanding these concepts is crucial for writing robust and maintainable Go code.
