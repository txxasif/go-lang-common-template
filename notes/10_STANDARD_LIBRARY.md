# Go Standard Library Core Packages: A Comprehensive Guide

This guide covers Go's standard library core packages, essential for building efficient, idiomatic applications.

## Table of Contents

- [Introduction](#introduction)
- [Strings Package](#strings-package)
- [Bytes Package](#bytes-package)
- [Time Package](#time-package)
- [Sync Package](#sync-package)
- [IO Package](#io-package)
- [Additional Core Packages](#additional-core-packages)
- [Best Practices](#best-practices)
- [Advanced Patterns](#advanced-patterns)
- [Testing with the Standard Library](#testing-with-the-standard-library)

## Introduction

Go's standard library is a cornerstone of its ecosystem, offering well-tested, performant packages for common tasks.

### Why the Standard Library?

- **Performance**: Optimized for Go's runtime
- **Reliability**: Thoroughly tested and stable
- **Consistency**: Adheres to Go's design principles
- **No Dependencies**: Reduces external library reliance

## Strings Package

The `strings` package provides tools for manipulating UTF-8 encoded, immutable strings.

### Key Operations

```go
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    s := "Hello, 世界"

    // Basic operations
    fmt.Println("Length (bytes):", len(s))               // 13
    fmt.Println("Rune count:", utf8.RuneCountInString(s)) // 9
    fmt.Println("Contains '世界':", strings.Contains(s, "世界")) // true
    fmt.Println("Split:", strings.Split(s, ","))         // [Hello  世界]
    fmt.Println("Join:", strings.Join([]string{"a", "b"}, "-")) // a-b
    fmt.Println("ToUpper:", strings.ToUpper(s))          // HELLO, 世界
}
```

### String Builder

`strings.Builder` is a mutable buffer for efficient string concatenation.

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    var b strings.Builder
    b.WriteString("Hello")
    b.WriteString(", ")
    b.WriteString("World!")
    fmt.Println(b.String()) // Hello, World!
    fmt.Println("Capacity:", b.Cap()) // Shows allocated capacity
}
```

## Bytes Package

The `bytes` package handles mutable byte slices, complementing the `strings` package.

### Byte Slices

```go
package main

import (
    "bytes"
    "fmt"
)

func main() {
    b := []byte("Hello")

    fmt.Println("Equal to 'Hello':", bytes.Equal(b, []byte("Hello"))) // true
    fmt.Println("Contains 'ell':", bytes.Contains(b, []byte("ell")))  // true
    fmt.Println("Split:", bytes.Split(b, []byte("l")))               // [[H e] [] [o]]
    fmt.Println("TrimSpace:", string(bytes.TrimSpace([]byte("  hi  ")))) // hi
}
```

### Bytes Buffer

`bytes.Buffer` is a resizable buffer for reading and writing bytes.

```go
package main

import (
    "bytes"
    "fmt"
)

func main() {
    var buf bytes.Buffer
    buf.WriteString("Hello")
    buf.WriteByte(',')
    buf.Write([]byte(" World!"))
    fmt.Println(buf.String()) // Hello, World!

    // Read a chunk
    chunk := make([]byte, 5)
    n, _ := buf.Read(chunk)
    fmt.Println(string(chunk[:n])) // Hello
}
```

## Time Package

The `time` package manages time-related operations, including formatting, parsing, and durations.

### Time Handling

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    now := time.Now()
    fmt.Println("Now:", now.Format("2006-01-02 15:04:05"))

    // Custom time
    t := time.Date(2025, time.April, 11, 12, 0, 0, 0, time.UTC)
    fmt.Println("Custom:", t)

    // Operations
    duration := 2 * time.Hour
    fmt.Println("Future:", now.Add(duration))
    fmt.Println("Is after:", t.After(now)) // true

    // Parse
    parsed, _ := time.Parse("2006-01-02", "2025-04-11")
    fmt.Println("Parsed:", parsed)
}
```

### Timers and Tickers

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Timer: one-time event
    timer := time.NewTimer(2 * time.Second)
    defer timer.Stop()
    fmt.Println("Waiting for timer...")
    <-timer.C
    fmt.Println("Timer fired!")

    // Ticker: repeating events
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    done := make(chan bool)
    go func() {
        time.Sleep(2 * time.Second)
        done <- true
    }()
    for {
        select {
        case <-ticker.C:
            fmt.Println("Tick")
        case <-done:
            fmt.Println("Done")
            return
        }
    }
}
```

## Sync Package

The `sync` package provides synchronization primitives for concurrent programming.

### Mutexes

`sync.Mutex` and `sync.RWMutex` ensure mutual exclusion.

```go
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu    sync.RWMutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count
}

func main() {
    c := &SafeCounter{}
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Inc()
        }()
    }
    wg.Wait()
    fmt.Println("Count:", c.Value()) // 100
}
```

### WaitGroup

`sync.WaitGroup` synchronizes goroutines.

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }
    wg.Wait()
    fmt.Println("All workers completed")
}
```

## IO Package

The `io` package defines core I/O interfaces and utilities.

### Core Interfaces

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

func main() {
    r := strings.NewReader("Hello, World!")
    buf := make([]byte, 5)
    n, err := r.Read(buf)
    fmt.Printf("Read %d bytes: %s\n", n, buf[:n]) // Read 5 bytes: Hello
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Utilities

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

func main() {
    // Copy
    src := strings.NewReader("Hello")
    var dst bytes.Buffer
    n, err := io.Copy(&dst, src)
    fmt.Printf("Copied %d bytes: %s\n", n, dst.String()) // Copied 5 bytes: Hello

    // MultiReader
    r1 := strings.NewReader("Hello, ")
    r2 := strings.NewReader("World!")
    multi := io.MultiReader(r1, r2)
    data, _ := io.ReadAll(multi)
    fmt.Println(string(data)) // Hello, World!

    // TeeReader
    var buf bytes.Buffer
    tee := io.TeeReader(strings.NewReader("Hi"), &buf)
    out, _ := io.ReadAll(tee)
    fmt.Println(string(out), buf.String()) // Hi Hi
}
```

## Additional Core Packages

### fmt Package

The `fmt` package handles formatted I/O.

```go
package main

import "fmt"

func main() {
    // Printing
    name := "Alice"
    fmt.Printf("Hello, %s! Score: %d\n", name, 95)
    fmt.Fprintf(&bytes.Buffer{}, "Log: %v", name)

    // Scanning
    var input string
    fmt.Scan(&input)
    fmt.Println("Input:", input)
}
```

### os Package

The `os` package provides platform-independent OS interactions.

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // File operations
    err := os.WriteFile("test.txt", []byte("Hello"), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    data, err := os.ReadFile("test.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(data)) // Hello

    // Environment
    fmt.Println("PATH:", os.Getenv("PATH"))
}
```

## Best Practices

### General Guidelines

- **Check Errors**: Always handle errors explicitly
- **Close Resources**: Use `defer` for cleanup
- **Use Standard Functions**: Prefer library implementations
- **Profile Performance**: Optimize based on benchmarks

### Package-Specific Tips

- **Strings**:

  - Use `strings.Builder` for heavy concatenation
  - Prefer `strings.Contains` over manual loops
  - Handle UTF-8 correctly with `unicode/utf8`

- **Bytes**:

  - Use `bytes.Buffer` for dynamic byte growth
  - Avoid unnecessary conversions to strings

- **Time**:

  - Use `time.Now()` for current time
  - Specify layouts explicitly
  - Stop `Timer` and `Ticker` to prevent leaks

- **Sync**:

  - Use `defer` for mutex unlocks
  - Prefer `RWMutex` for read-heavy workloads
  - Use `Once` for initialization

- **IO**:
  - Buffer I/O operations for efficiency
  - Implement interfaces for custom readers/writers
  - Use utilities like `io.Copy` for simplicity

## Advanced Patterns

### Concurrent String Processing

```go
package main

import (
    "fmt"
    "strings"
    "sync"
)

func processChunk(chunk string) string {
    return strings.ToUpper(chunk)
}

func main() {
    input := "hello,world,go,programming"
    chunks := strings.Split(input, ",")
    results := make([]string, len(chunks))
    var wg sync.WaitGroup

    for i, chunk := range chunks {
        wg.Add(1)
        go func(i int, s string) {
            defer wg.Done()
            results[i] = processChunk(s)
        }(i, chunk)
    }

    wg.Wait()
    fmt.Println(strings.Join(results, ",")) // HELLO,WORLD,GO,PROGRAMMING
}
```

### Time-Based Caching

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Cache struct {
    mu     sync.RWMutex
    data   map[string]string
    expiry map[string]time.Time
    ttl    time.Duration
}

func NewCache(ttl time.Duration) *Cache {
    c := &Cache{
        data:   make(map[string]string),
        expiry: make(map[string]time.Time),
        ttl:    ttl,
    }
    go c.cleanup()
    return c
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
    c.expiry[key] = time.Now().Add(c.ttl)
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    if time.Now().After(c.expiry[key]) {
        return "", false
    }
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        c.mu.Lock()
        for k, expiry := range c.expiry {
            if time.Now().After(expiry) {
                delete(c.data, k)
                delete(c.expiry, k)
            }
        }
        c.mu.Unlock()
    }
}

func main() {
    cache := NewCache(2 * time.Second)
    cache.Set("key", "value")
    if v, ok := cache.Get("key"); ok {
        fmt.Println("Found:", v) // Found: value
    }
    time.Sleep(3 * time.Second)
    if _, ok := cache.Get("key"); !ok {
        fmt.Println("Expired")
    }
}
```

## Testing with the Standard Library

### Testing Strings and Bytes

```go
package main

import (
    "bytes"
    "strings"
    "testing"
)

func TestStringBuilder(t *testing.T) {
    var b strings.Builder
    b.WriteString("test")
    if got := b.String(); got != "test" {
        t.Errorf("Expected 'test', got %q", got)
    }
}

func TestBytesBuffer(t *testing.T) {
    var buf bytes.Buffer
    buf.WriteString("hello")
    if got := buf.String(); got != "hello" {
        t.Errorf("Expected 'hello', got %q", got)
    }
}
```

### Testing Time-Dependent Code

```go
package main

import (
    "testing"
    "time"
)

func TestTimeout(t *testing.T) {
    timer := time.NewTimer(100 * time.Millisecond)
    select {
    case <-timer.C:
        return
    case <-time.After(200 * time.Millisecond):
        t.Error("Timer didn't fire in time")
    }
}
```

_Last Updated: April 11, 2025_
