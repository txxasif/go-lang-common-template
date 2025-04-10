# Go Loops and Iteration: A Complete Guide

## Table of Contents

1. [For Loop Variations](#for-loop-variations)
2. [While Loop Pattern](#while-loop-pattern)
3. [Infinite Loops](#infinite-loops)
4. [Control Statements](#control-statements)
5. [Common Patterns](#common-patterns)
6. [Best Practices](#best-practices)
7. [Advanced Patterns](#advanced-patterns)

## For Loop Variations

### Basic For Loop

```go
// Traditional C-style for loop
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// Multiple variables in for loop
for i, j := 0, 10; i < j; i, j = i+1, j-1 {
    fmt.Printf("i: %d, j: %d\n", i, j)
}
```

## While Loop Pattern

### Basic While Loop

```go
// Go's while loop equivalent
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}

// With condition function
func shouldContinue() bool {
    // Some condition logic
    return true
}

for shouldContinue() {
    // Loop body
}
```

### Conditional While Loop

```go
// Processing with condition
func processData() bool {
    // Return false when done
    return true
}

for ok := true; ok; ok = processData() {
    // Process data
}
```

## Infinite Loops

### Basic Infinite Loop

```go
// Simple infinite loop
for {
    // Will run forever unless broken
    if someCondition {
        break
    }
}

// With ticker
ticker := time.NewTicker(time.Second)
for {
    select {
    case <-ticker.C:
        // Do something every second
    }
}
```

### Controlled Infinite Loop

```go
// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

for {
    select {
    case <-ctx.Done():
        return
    default:
        // Do work
    }
}
```

## Control Statements

### Break and Continue

```go
// Break with label
OuterLoop:
    for i := 0; i < 5; i++ {
        for j := 0; j < 5; j++ {
            if i*j > 10 {
                break OuterLoop
            }
        }
    }

// Continue with condition
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue // Skip even numbers
    }
    fmt.Println(i)
}
```

## Common Patterns

### Retry Pattern

```go
// Retry with maximum attempts
const maxRetries = 3
for attempts := 0; attempts < maxRetries; attempts++ {
    err := someOperation()
    if err == nil {
        break
    }
    time.Sleep(time.Second * time.Duration(attempts))
}

// Exponential backoff
func withExponentialBackoff(operation func() error) error {
    backoff := time.Second
    maxBackoff := time.Minute

    for {
        err := operation()
        if err == nil {
            return nil
        }

        if backoff > maxBackoff {
            return fmt.Errorf("max retries exceeded: %w", err)
        }

        time.Sleep(backoff)
        backoff *= 2
    }
}
```

### Processing Pattern

```go
// Process until empty
for !queue.IsEmpty() {
    item := queue.Dequeue()
    process(item)
}

// Process with timeout
func processWithTimeout(timeout time.Duration) error {
    timer := time.NewTimer(timeout)
    defer timer.Stop()

    for {
        select {
        case <-timer.C:
            return errors.New("timeout")
        default:
            if done := processNext(); done {
                return nil
            }
        }
    }
}
```

## Best Practices

### Resource Management

```go
// Proper cleanup in loops
for {
    resource, err := acquireResource()
    if err != nil {
        continue
    }
    defer resource.Release()

    if err := useResource(resource); err != nil {
        break
    }
}

// With context
func processWithContext(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Process
        }
    }
}
```

### Error Handling

```go
// Accumulate errors
var errors []error
for _, item := range items {
    if err := process(item); err != nil {
        errors = append(errors, err)
    }
}

// Continue on certain errors
for {
    err := someOperation()
    if err != nil {
        if errors.Is(err, ErrRetryable) {
            continue
        }
        return err
    }
    break
}
```

## Advanced Patterns

### Worker Pool Pattern

```go
func workerPool(numWorkers int, jobs <-chan Job, results chan<- Result) {
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
}
```

### Pipeline Pattern

```go
func pipeline(input <-chan int) <-chan int {
    output := make(chan int)
    go func() {
        defer close(output)
        for val := range input {
            // Process and send
            output <- transform(val)
        }
    }()
    return output
}
```

### Rate Limited Loop

```go
func rateLimitedLoop(rate time.Duration) {
    ticker := time.NewTicker(rate)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Do rate-limited work
        }
    }
}
```

### Event Processing Loop

```go
type EventLoop struct {
    events chan Event
    stop   chan struct{}
}

func (l *EventLoop) Run() {
    for {
        select {
        case event := <-l.events:
            l.handleEvent(event)
        case <-l.stop:
            return
        }
    }
}
```

Remember:

- Use appropriate loop patterns for your use case
- Always consider termination conditions
- Handle resources properly
- Implement proper error handling
- Consider using contexts for cancellation
- Use channels for concurrent operations
- Implement proper rate limiting when needed
- Clean up resources in loops
- Use labels sparingly and only when necessary
- Consider performance implications of your loop patterns

This guide covers the essential aspects of loops in Go, focusing on while-loop patterns and other iteration techniques. Understanding these patterns is crucial for writing efficient and maintainable Go code.
