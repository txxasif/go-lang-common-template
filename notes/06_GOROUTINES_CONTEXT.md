# Go Concurrency: Context and Goroutines Guide

## Table of Contents

- [Go Concurrency: Context and Goroutines Guide](#go-concurrency-context-and-goroutines-guide)
  - [Table of Contents](#table-of-contents)
  - [Goroutines Basics](#goroutines-basics)
    - [Simple Goroutine](#simple-goroutine)
  - [Context Fundamentals](#context-fundamentals)
    - [Context Types](#context-types)
  - [Patterns and Best Practices](#patterns-and-best-practices)
    - [Worker Pool Pattern](#worker-pool-pattern)
    - [Pipeline Pattern](#pipeline-pattern)
  - [Advanced Patterns](#advanced-patterns)
    - [Fan-Out Fan-In Pattern](#fan-out-fan-in-pattern)
    - [Rate Limiting Pattern](#rate-limiting-pattern)
  - [Error Handling](#error-handling)
    - [Error Propagation](#error-propagation)
  - [Performance Considerations](#performance-considerations)
    - [Goroutine Pool](#goroutine-pool)

## Goroutines Basics

### Simple Goroutine

```go
func main() {
    // Start a goroutine
    go func() {
        fmt.Println("Running in goroutine")
    }()

    // Wait for goroutine to finish
    time.Sleep(time.Second)
}
```

Visual Representation:

```
Main Thread:  ────────────────────────►
                    │
Goroutine:          ├──────►
                    │
Time: ─────────────────────────────►
```

## Context Fundamentals

### Context Types

```go
// Create root context
ctx := context.Background()

// With cancellation
ctx, cancel := context.WithCancel(ctx)
defer cancel()

// With timeout
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// With deadline
ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
defer cancel()

// With values
ctx = context.WithValue(ctx, "key", "value")
```

Visual Representation:

```
Context Hierarchy:

Background
    │
    ├── WithCancel
    │     └── WithValue
    │
    ├── WithTimeout
    │     └── WithValue
    │
    └── WithDeadline
          └── WithValue
```

## Patterns and Best Practices

### Worker Pool Pattern

```go
func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
    for {
        select {
        case <-ctx.Done():
            return
        case job, ok := <-jobs:
            if !ok {
                return
            }
            results <- process(job)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    const numWorkers = 3
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start workers
    for i := 0; i < numWorkers; i++ {
        go worker(ctx, i, jobs, results)
    }

    // Send jobs
    for i := 0; i < 100; i++ {
        jobs <- i
    }
    close(jobs)
}
```

Visual Representation:

```
Jobs Channel:    [Job1]─[Job2]─[Job3]─►
                    │      │      │
Workers:        Worker1 Worker2 Worker3
                    │      │      │
Results Channel: [Res1]─[Res2]─[Res3]─►
```

### Pipeline Pattern

```go
func generator(ctx context.Context) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 0; ; i++ {
            select {
            case <-ctx.Done():
                return
            case out <- i:
            }
        }
    }()
    return out
}

func multiplier(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for num := range in {
            select {
            case <-ctx.Done():
                return
            case out <- num * 2:
            }
        }
    }()
    return out
}
```

Visual Representation:

```
Pipeline Flow:

Generator ──► Multiplier ──► Consumer
   [1,2,3] ──► [2,4,6] ──► Process
      │           │           │
      └───────────┴───────────┘
          Context Control
```

## Advanced Patterns

### Fan-Out Fan-In Pattern

```go
func fanOut(ctx context.Context, source <-chan int, n int) []<-chan int {
    channels := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        channels[i] = processor(ctx, source)
    }
    return channels
}

func fanIn(ctx context.Context, channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    multiplexed := make(chan int)

    multiplex := func(c <-chan int) {
        defer wg.Done()
        for i := range c {
            select {
            case <-ctx.Done():
                return
            case multiplexed <- i:
            }
        }
    }

    wg.Add(len(channels))
    for _, c := range channels {
        go multiplex(c)
    }

    go func() {
        wg.Wait()
        close(multiplexed)
    }()

    return multiplexed
}
```

Visual Representation:

```
Fan-Out/Fan-In:

Source ─┬─► Processor1 ─┐
        ├─► Processor2 ─┼─► Merged Output
        └─► Processor3 ─┘
```

### Rate Limiting Pattern

```go
func rateLimitedWorker(ctx context.Context, rate time.Duration) {
    ticker := time.NewTicker(rate)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // Do work
        }
    }
}
```

Visual Representation:

```
Rate Limiting:

Time: ─────┬─────┬─────┬─────►
           │     │     │
Work:      W     W     W
           │     │     │
Rate:      R     R     R
```

## Error Handling

### Error Propagation

```go
type Result struct {
    Value int
    Err   error
}

func processWithError(ctx context.Context) <-chan Result {
    results := make(chan Result)
    go func() {
        defer close(results)
        for {
            select {
            case <-ctx.Done():
                results <- Result{Err: ctx.Err()}
                return
            default:
                // Process and send result
                if err := process(); err != nil {
                    results <- Result{Err: err}
                    return
                }
            }
        }
    }()
    return results
}
```

Visual Representation:

```
Error Handling Flow:

Operation ──► Error Check ──┬─► Success Path
                           │
                           └─► Error Path
                                   │
                                   v
                             Error Handling
```

## Performance Considerations

### Goroutine Pool

```go
type Pool struct {
    work    chan func()
    sem     chan struct{}
    timeout time.Duration
}

func NewPool(size int, timeout time.Duration) *Pool {
    return &Pool{
        work:    make(chan func()),
        sem:     make(chan struct{}, size),
        timeout: timeout,
    }
}

func (p *Pool) Submit(task func()) error {
    select {
    case p.sem <- struct{}{}:
        go func() {
            defer func() { <-p.sem }()
            task()
        }()
        return nil
    case <-time.After(p.timeout):
        return errors.New("pool overloaded")
    }
}
```

Visual Representation:

```
Pool Management:

Tasks ──► Pool Queue [│││││] ──┬─► Worker1
                              ├─► Worker2
                              └─► Worker3

Semaphore: [□□□] (Limited Slots)
```

Remember:

- Always use context for cancellation
- Properly clean up resources
- Handle errors appropriately
- Consider rate limiting
- Use appropriate channel patterns
- Monitor goroutine lifecycles
- Implement proper timeout mechanisms
- Use worker pools for controlled concurrency
- Implement graceful shutdown
- Consider memory and CPU implications

This guide provides a comprehensive overview of Go concurrency patterns using context and goroutines. The visual representations help understand the flow and relationships between different components in concurrent programs.
