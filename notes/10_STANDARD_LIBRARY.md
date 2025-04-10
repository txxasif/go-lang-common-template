# Go Standard Library Core Packages: A Comprehensive Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Strings Package](#strings-package)
3. [Bytes Package](#bytes-package)
4. [Time Package](#time-package)
5. [Sync Package](#sync-package)
6. [IO Package](#io-package)
7. [Best Practices](#best-practices)
8. [Advanced Patterns](#advanced-patterns)

## Introduction

The Go standard library provides a rich set of packages that are essential for everyday programming. These packages are designed to be efficient, safe, and easy to use. Understanding these core packages is crucial for writing idiomatic Go code.

### Why Use Standard Library?

1. **Performance**: Optimized for Go's runtime
2. **Safety**: Well-tested and production-ready
3. **Consistency**: Follows Go's design principles
4. **Maintainability**: Widely understood and documented

## Strings Package

### Understanding Strings in Go

In Go, strings are immutable sequences of bytes. They are UTF-8 encoded by default, making them suitable for international text.

#### Basic String Operations

```go
// String creation
s := "Hello, 世界"  // UTF-8 encoded string

// Length (in bytes, not characters)
length := len(s)  // Returns 13 (7 for "Hello, " + 6 for "世界")

// Character count
runeCount := utf8.RuneCountInString(s)  // Returns 9
```

#### String Manipulation

```go
// Substring
sub := s[0:5]  // "Hello"

// Contains
contains := strings.Contains(s, "世界")  // true

// Split
parts := strings.Split(s, ",")  // ["Hello", " 世界"]

// Join
joined := strings.Join(parts, "-")  // "Hello- 世界"
```

#### String Builder

The `strings.Builder` is an efficient way to build strings, especially when concatenating multiple strings.

```go
var builder strings.Builder

// Write operations
builder.WriteString("Hello")
builder.WriteString(", ")
builder.WriteString("World!")

// Get result
result := builder.String()  // "Hello, World!"
```

**Why Use String Builder?**

1. **Efficiency**: Minimizes memory allocations
2. **Performance**: Faster than regular string concatenation
3. **Memory**: Reduces garbage collection pressure

## Bytes Package

### Understanding Bytes in Go

The `bytes` package provides functions for manipulating byte slices, which are mutable sequences of bytes.

#### Byte Slice Operations

```go
// Create byte slice
b := []byte("Hello")

// Compare
result := bytes.Compare(b, []byte("Hello"))  // 0 (equal)

// Contains
contains := bytes.Contains(b, []byte("ell"))  // true

// Split
parts := bytes.Split(b, []byte("l"))  // [][]byte{[]byte("He"), []byte("o")}
```

#### Buffer Operations

The `bytes.Buffer` is a variable-sized buffer of bytes with Read and Write methods.

```go
// Create buffer
var buf bytes.Buffer

// Write operations
buf.WriteString("Hello")
buf.WriteByte(',')
buf.Write([]byte(" World!"))

// Read operations
result := buf.String()  // "Hello, World!"
```

**When to Use Bytes Package?**

1. **Performance**: When working with raw bytes
2. **Memory**: When minimizing allocations
3. **IO**: When dealing with binary data

## Time Package

### Understanding Time in Go

The `time` package provides functionality for measuring and displaying time.

#### Time Creation and Formatting

```go
// Current time
now := time.Now()

// Specific time
t := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)

// Formatting
formatted := t.Format("2006-01-02 15:04:05")

// Parsing
parsed, _ := time.Parse("2006-01-02", "2023-01-01")
```

#### Time Operations

```go
// Duration
duration := time.Hour * 2

// Add time
future := now.Add(duration)

// Subtract time
past := now.Add(-duration)

// Time comparison
isAfter := future.After(now)  // true
```

#### Tickers and Timers

```go
// Ticker (repeating)
ticker := time.NewTicker(time.Second)
defer ticker.Stop()

// Timer (one-time)
timer := time.NewTimer(time.Second)
defer timer.Stop()

// Usage
select {
case <-ticker.C:
    // Every second
case <-timer.C:
    // After one second
}
```

## Sync Package

### Understanding Concurrency in Go

The `sync` package provides basic synchronization primitives.

#### Mutex

A Mutex is a mutual exclusion lock.

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

#### WaitGroup

A WaitGroup waits for a collection of goroutines to finish.

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // Do work
    }(i)
}

wg.Wait()  // Wait for all goroutines
```

#### RWMutex

A RWMutex is a reader/writer mutual exclusion lock.

```go
type SafeMap struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (m *SafeMap) Get(key string) interface{} {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.data[key]
}

func (m *SafeMap) Set(key string, value interface{}) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = value
}
```

## IO Package

### Understanding IO in Go

The `io` package provides basic interfaces to I/O primitives.

#### Basic IO Operations

```go
// Reader interface
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer interface
type Writer interface {
    Write(p []byte) (n int, err error)
}

// Closer interface
type Closer interface {
    Close() error
}
```

#### File Operations

```go
// Open file
file, err := os.Open("file.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// Read file
data := make([]byte, 100)
n, err := file.Read(data)

// Write file
err = os.WriteFile("output.txt", data, 0644)
```

#### IO Utilities

```go
// Copy
src := strings.NewReader("Hello")
dst := new(bytes.Buffer)
io.Copy(dst, src)

// MultiReader
r1 := strings.NewReader("Hello")
r2 := strings.NewReader("World")
multi := io.MultiReader(r1, r2)

// TeeReader
var buf bytes.Buffer
tee := io.TeeReader(strings.NewReader("Hello"), &buf)
```

## Best Practices

### String Operations

1. **Use `strings.Builder` for multiple concatenations**

   - More efficient than `+` operator
   - Reduces memory allocations
   - Better performance for large strings

2. **Prefer `strings.Contains` over manual searching**

   - More readable
   - Handles edge cases
   - Optimized implementation

3. **Use `strings.Split` for complex parsing**
   - Handles multiple separators
   - Consistent behavior
   - Memory efficient

### Byte Operations

1. **Use `bytes.Buffer` for byte manipulation**

   - Efficient for growing buffers
   - Provides both read and write operations
   - Memory efficient

2. **Prefer `bytes.Compare` for byte slice comparison**
   - More efficient than manual comparison
   - Handles edge cases
   - Consistent ordering

### Time Operations

1. **Always use `time.Now()` for current time**

   - More accurate than manual time
   - Handles timezone correctly
   - Thread-safe

2. **Use `time.After` for timeouts**
   - Cleaner than manual timers
   - Handles cancellation
   - Memory efficient

### Sync Operations

1. **Always use `defer` with mutex unlocks**

   - Prevents deadlocks
   - Handles panics
   - More maintainable

2. **Use `RWMutex` for read-heavy operations**
   - Better performance
   - Allows concurrent reads
   - Maintains consistency

### IO Operations

1. **Always check errors**

   - Prevents silent failures
   - Better error handling
   - More reliable code

2. **Use `defer` for closing resources**
   - Ensures cleanup
   - Handles panics
   - More maintainable

## Advanced Patterns

### String Pooling

```go
type StringPool struct {
    mu sync.RWMutex
    m  map[string]string
}

func (p *StringPool) Get(s string) string {
    p.mu.RLock()
    if cached, ok := p.m[s]; ok {
        p.mu.RUnlock()
        return cached
    }
    p.mu.RUnlock()

    p.mu.Lock()
    defer p.mu.Unlock()

    if cached, ok := p.m[s]; ok {
        return cached
    }

    p.m[s] = s
    return s
}
```

### Time-based Rate Limiting

```go
type RateLimiter struct {
    mu       sync.Mutex
    rate     time.Duration
    lastTime time.Time
}

func (r *RateLimiter) Wait() {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(r.lastTime)
    if elapsed < r.rate {
        time.Sleep(r.rate - elapsed)
    }
    r.lastTime = time.Now()
}
```

### IO Pipeline

```go
type Pipeline struct {
    stages []func(io.Reader) io.Reader
}

func (p *Pipeline) Process(r io.Reader) io.Reader {
    for _, stage := range p.stages {
        r = stage(r)
    }
    return r
}
```

Remember:

- Use appropriate package for the task
- Handle errors properly
- Clean up resources
- Consider performance implications
- Use standard library functions when available
- Follow Go idioms and patterns
- Document complex operations
- Test edge cases
- Consider concurrent access
- Use appropriate synchronization primitives

This guide covers the fundamental and advanced aspects of Go's core standard library packages. Understanding these packages is crucial for writing efficient and idiomatic Go code.
