# Go Data Types: A Comprehensive Guide

## Table of Contents

1. [Basic Types](#basic-types)
2. [Composite Types](#composite-types)
3. [Zero Values](#zero-values)
4. [Type Conversions](#type-conversions)
5. [Type Declarations](#type-declarations)
6. [Collections](#collections)
7. [Special Types](#special-types)
8. [Best Practices](#best-practices)

## Basic Types

### Numeric Types

```go
// Integers
var (
    age    int     // platform dependent size (32 or 64 bit)
    count  int8    // -128 to 127
    stock  int16   // -32768 to 32767
    users  int32   // -2^31 to 2^31-1
    points int64   // -2^63 to 2^63-1

    // Unsigned integers
    items  uint    // platform dependent size
    small  uint8   // 0 to 255
    medium uint16  // 0 to 65535
    large  uint32  // 0 to 4294967295
    huge   uint64  // 0 to 18446744073709551615
)

// Floating point
var (
    price    float32    // ±1.18E-38 to ±3.4E38
    distance float64    // ±2.23E-308 to ±1.80E308
)

// Complex numbers
var (
    c64  complex64     // Complex numbers with float32 real and imaginary parts
    c128 complex128    // Complex numbers with float64 real and imaginary parts
)
```

### Boolean and String Types

```go
var (
    isValid bool      // true or false
    name    string    // UTF-8 encoded text
)

// String operations
str := "Hello"
length := len(str)
byte := str[0]        // Get byte at index
substr := str[1:4]    // Substring "ell"
```

## Composite Types

### Arrays

```go
// Fixed-size arrays
var numbers [5]int           // Array of 5 integers
matrix := [2][3]int{         // 2D array
    {1, 2, 3},
    {4, 5, 6},
}

// Array initialization
colors := [3]string{"red", "green", "blue"}
numbers := [...]int{1, 2, 3} // Size determined by initializer
```

### Slices

```go
// Dynamic-size arrays
var slice []int              // Nil slice
numbers := make([]int, 5)    // Slice with length 5
numbers := make([]int, 5, 10) // Slice with length 5, capacity 10

// Slice operations
slice = append(slice, 1)     // Append element
slice = slice[1:4]           // Slice of slice
copy(dest, src)              // Copy slices
```

### Maps

```go
// Key-value pairs
var scores map[string]int            // Nil map
scores = make(map[string]int)        // Initialize map

// Map operations
scores["Alice"] = 95                 // Set value
score, exists := scores["Bob"]       // Check existence
delete(scores, "Alice")              // Remove entry

// Map initialization
config := map[string]string{
    "host": "localhost",
    "port": "8080",
}
```

### Structs

```go
// Custom data types
type Person struct {
    Name    string
    Age     int
    Address struct {
        Street  string
        City    string
    }
}

// Struct initialization
person := Person{
    Name: "John",
    Age:  30,
}

// Anonymous structs
point := struct {
    x, y int
}{10, 20}
```

## Zero Values

```go
var (
    intZero     int     // 0
    floatZero   float64 // 0.0
    boolZero    bool    // false
    stringZero  string  // ""
    pointerZero *int    // nil
    sliceZero   []int   // nil
    mapZero     map[string]int // nil
    structZero  struct{}       // {}
)
```

## Type Conversions

```go
// Basic type conversions
i := 42
f := float64(i)           // int to float64
s := string(65)           // int to string (ASCII)
b := []byte("Hello")      // string to byte slice

// Strconv package for string conversions
import "strconv"

str := strconv.Itoa(123)  // Int to string
num, _ := strconv.Atoi("123") // String to int
```

## Type Declarations

```go
// Type aliases
type ID = string

// New types
type Age int
type Money float64

// Method receiver types
type Counter int

func (c *Counter) Increment() {
    *c++
}
```

## Collections

### Slices Best Practices

```go
// Pre-allocate with make
slice := make([]int, 0, 100)  // Capacity 100

// Append efficiently
func appendSlice(slice []int, items ...int) []int {
    if len(items) == 0 {
        return slice
    }

    // Ensure capacity
    newLen := len(slice) + len(items)
    if newLen > cap(slice) {
        newSlice := make([]int, len(slice), newLen*2)
        copy(newSlice, slice)
        slice = newSlice
    }

    return append(slice, items...)
}
```

### Map Best Practices

```go
// Initialize with expected size
users := make(map[string]User, 100)

// Thread-safe maps
import "sync"

type SafeMap struct {
    sync.RWMutex
    data map[string]interface{}
}

func (m *SafeMap) Get(key string) interface{} {
    m.RLock()
    defer m.RUnlock()
    return m.data[key]
}
```

## Special Types

### Interface

```go
// Interface definition
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Empty interface
var i interface{}
i = 42
i = "hello"
```

### Channel

```go
// Channel types
ch := make(chan int)           // Unbuffered channel
buf := make(chan int, 10)      // Buffered channel
```

## Best Practices

### 1. Type Selection

```go
// Use the right type for the job
type UserID int64      // For IDs
type Amount float64    // For money
type Timestamp int64   // For time
```

### 2. Memory Efficiency

```go
// Use appropriate sizes
var small uint8    // For small numbers (0-255)
var normal int     // For most integers
var big int64     // For very large numbers
```

### 3. Type Safety

```go
// Custom types for type safety
type Temperature float64
type Distance float64

func (t Temperature) ToCelsius() Temperature {
    return (t - 32) * 5 / 9
}
```

### 4. Error Handling

```go
// Custom error types
type ValidationError struct {
    Field string
    Error string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Error)
}
```

### 5. Performance Considerations

```go
// String concatenation
var builder strings.Builder
for i := 0; i < 1000; i++ {
    builder.WriteString("a")
}
result := builder.String()

// Slice operations
// Avoid unnecessary allocations
slice = append(slice[:i], slice[i+1:]...) // Remove element
```

Remember:

- Choose the appropriate type for your data
- Consider memory usage and performance
- Use custom types for better type safety
- Leverage zero values
- Use built-in functions and standard library effectively
- Consider thread safety when necessary
- Use interfaces for abstraction
- Handle type conversions safely

This guide covers the essential aspects of Go's type system. Understanding these concepts is crucial for writing efficient and maintainable Go code.
