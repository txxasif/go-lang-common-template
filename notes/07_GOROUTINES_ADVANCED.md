# Ultimate Guide to Go Concurrency: Context & Goroutines

## Table of Contents

1. [Core Concepts](#core-concepts)
2. [Visual Understanding](#visual-understanding)
3. [Basic Patterns](#basic-patterns)
4. [Advanced Patterns](#advanced-patterns)
5. [Real-World Examples](#real-world-examples)
6. [Debugging & Troubleshooting](#debugging--troubleshooting)
7. [Best Practices & Pitfalls](#best-practices--pitfalls)

## Core Concepts

### What is a Goroutine?

```go
// A goroutine is a lightweight thread managed by the Go runtime
// Visual Representation:
/*
OS Thread:     ────────────────────────►
Goroutines:    ─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─►
                │ │ │ │ │ │ │ │ │ │
                G G G G G G G G G G
*/

// Simple Example
func main() {
    // Main goroutine
    go func() {
        // New goroutine
        fmt.Println("Hello from goroutine!")
    }()
    time.Sleep(time.Second)
}
```

### Understanding Context

```go
// Context is for cancellation, deadlines, and passing values
type ContextExample struct {
    ctx    context.Context
    cancel context.CancelFunc
}

/*
Context Tree:
                   Background
                       │
                  WithCancel
                 /    |     \
           Value   Timeout  Deadline
            |        |         |
         Cancel   Cancel    Cancel
*/

// Different Types of Contexts
func contextTypes() {
    // 1. Root context
    root := context.Background()

    // 2. Cancellable context
    ctx1, cancel1 := context.WithCancel(root)
    defer cancel1()

    // 3. Timeout context
    ctx2, cancel2 := context.WithTimeout(root, 5*time.Second)
    defer cancel2()

    // 4. Deadline context
    ctx3, cancel3 := context.WithDeadline(root, time.Now().Add(1*time.Hour))
    defer cancel3()

    // 5. Value context
    ctx4 := context.WithValue(root, "key", "value")
}
```

## Visual Understanding

### Goroutine Lifecycle

```go
/*
Goroutine Lifecycle:

Creation:   main()──►go func()──►Running
              │
Execution:    │──►Processing──►Completed
              │
Termination:  └──►Done

Memory Usage: [Small Stack]──►[Grows as needed]──►[Released]
*/

func goroutineLifecycle() {
    // Start
    go func() {
        // Initialize
        fmt.Println("Starting")

        // Process
        time.Sleep(time.Second)

        // Cleanup
        fmt.Println("Done")
    }()
}
```

### Channel Communication

```go
/*
Channel Communication:

Producer ──[chan]──► Consumer
    │                   │
    └───────────────────┘
         Synchronization

Buffered Channel:
[Item1][Item2][Item3][ ][ ]
  Full  Full   Full  Empty Empty
*/

func channelExample() {
    // Unbuffered channel
    ch1 := make(chan int)

    // Buffered channel
    ch2 := make(chan int, 5)

    // Producer
    go func() {
        for i := 0; i < 5; i++ {
            ch2 <- i
        }
        close(ch2)
    }()

    // Consumer
    for v := range ch2 {
        fmt.Println(v)
    }
}
```

## Basic Patterns

### Worker Pool with Context

```go
/*
Worker Pool Architecture:

Jobs ──► [Job Queue] ──┬──► Worker1 ──┐
                      ├──► Worker2 ──┼──► [Result Queue] ──► Results
                      └──► Worker3 ──┘
*/

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID int
    Data  string
    Error error
}

func WorkerPool(ctx context.Context, numWorkers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                select {
                case <-ctx.Done():
                    return
                default:
                    results <- processJob(job)
                }
            }
        }(i)
    }

    // Close results when all workers are done
    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

### Rate Limiter with Burst

```go
/*
Rate Limiting:

Requests: R R R R R R R R
Time:     │ │ │ │ │ │ │ │
Bucket:   └─┴─┴─┴─┴─┴─┴─┘
          [Token Bucket Algorithm]
*/

type RateLimiter struct {
    rate     time.Duration
    burst    int
    tokens   chan struct{}
    stopChan chan struct{}
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
    rl := &RateLimiter{
        rate:     rate,
        burst:    burst,
        tokens:   make(chan struct{}, burst),
        stopChan: make(chan struct{}),
    }

    // Token generator
    go func() {
        ticker := time.NewTicker(rate)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                select {
                case rl.tokens <- struct{}{}:
                default:
                }
            case <-rl.stopChan:
                return
            }
        }
    }()

    return rl
}
```

## Advanced Patterns

### Circuit Breaker

```go
/*
Circuit Breaker States:

Closed ──► Half-Open ──► Open
   │           │           │
   └───────────┴───────────┘
*/

type CircuitBreaker struct {
    mu             sync.RWMutex
    failureCount   int
    threshold      int
    resetTimeout   time.Duration
    lastFailure    time.Time
    state          string
}

func (cb *CircuitBreaker) Execute(ctx context.Context, cmd func() error) error {
    if !cb.canExecute() {
        return errors.New("circuit open")
    }

    err := cmd()
    cb.recordResult(err)
    return err
}
```

### Pub/Sub Pattern

```go
/*
Publisher/Subscriber:

Publisher ──► [Topic] ──┬──► Subscriber1
                       ├──► Subscriber2
                       └──► Subscriber3
*/

type PubSub struct {
    mu     sync.RWMutex
    subs   map[string][]chan interface{}
    closed bool
}

func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
    ps.mu.Lock()
    defer ps.mu.Unlock()

    ch := make(chan interface{}, 1)
    ps.subs[topic] = append(ps.subs[topic], ch)
    return ch
}

func (ps *PubSub) Publish(topic string, msg interface{}) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()

    if ps.closed {
        return
    }

    for _, ch := range ps.subs[topic] {
        go func(ch chan interface{}) {
            ch <- msg
        }(ch)
    }
}
```

## Real-World Examples

### HTTP Server with Graceful Shutdown

```go
func main() {
    srv := &http.Server{
        Addr: ":8080",
    }

    // Server run context
    serverCtx, serverStopCtx := context.WithCancel(context.Background())

    // Listen for shutdown signal
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sig

        // Shutdown signal with grace period
        shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

        go func() {
            <-shutdownCtx.Done()
            if shutdownCtx.Err() == context.DeadlineExceeded {
                log.Fatal("graceful shutdown timed out.. forcing exit.")
            }
        }()

        // Trigger graceful shutdown
        err := srv.Shutdown(shutdownCtx)
        if err != nil {
            log.Fatal(err)
        }
        serverStopCtx()
    }()

    // Run the server
    err := srv.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }

    // Wait for server context to be stopped
    <-serverCtx.Done()
}
```

### Database Connection Pool

```go
/*
Connection Pool:

Requests ──► [Available Connections] ──► DB
              │ │ │ │ │
              └─┴─┴─┴─┘
             [Connection Pool]
*/

type DBPool struct {
    mu          sync.Mutex
    connections chan *sql.DB
    maxSize     int
    connStr     string
}

func NewDBPool(maxSize int, connStr string) *DBPool {
    return &DBPool{
        connections: make(chan *sql.DB, maxSize),
        maxSize:     maxSize,
        connStr:     connStr,
    }
}

func (p *DBPool) GetConn(ctx context.Context) (*sql.DB, error) {
    select {
    case conn := <-p.connections:
        return conn, nil
    default:
        return p.createConn()
    }
}
```

## Debugging & Troubleshooting

### Goroutine Leak Detection

```go
/*
Leak Detection:

Before: [G1][G2][G3]
After:  [G1][G2][G3][G4]...
*/

func detectLeaks() {
    initial := runtime.NumGoroutine()

    // Run your code

    time.Sleep(time.Second)
    current := runtime.NumGoroutine()

    if current > initial {
        fmt.Printf("Possible goroutine leak: %d additional goroutines\n",
            current - initial)
    }
}
```

### Context Cancellation Debugging

```go
func debugContext(ctx context.Context) {
    select {
    case <-ctx.Done():
        switch ctx.Err() {
        case context.Canceled:
            log.Println("Context was canceled")
        case context.DeadlineExceeded:
            log.Println("Context deadline exceeded")
        }
    default:
        log.Println("Context still valid")
    }
}
```

## Best Practices & Pitfalls

### Memory Management

```go
// Good Practice
func goodPractice() {
    const maxWorkers = 100
    sem := make(chan struct{}, maxWorkers)

    for work := range workQueue {
        sem <- struct{}{} // Acquire
        go func(w Work) {
            defer func() { <-sem }() // Release
            process(w)
        }(work)
    }
}

// Bad Practice - Potential Memory Leak
func badPractice() {
    for work := range workQueue {
        go func() {
            // work variable is shared across all goroutines
            process(work)
        }()
    }
}
```

### Error Handling

```go
type Result struct {
    Value interface{}
    Err   error
}

func safeGoroutine(ctx context.Context) <-chan Result {
    results := make(chan Result, 1)

    go func() {
        defer func() {
            if r := recover(); r != nil {
                results <- Result{Err: fmt.Errorf("panic: %v", r)}
            }
            close(results)
        }()

        // Do work
        value, err := doWork(ctx)
        results <- Result{Value: value, Err: err}
    }()

    return results
}
```

Remember:

- Always use context for cancellation and timeouts
- Clean up resources properly
- Handle panics in goroutines
- Be careful with shared memory
- Use buffered channels when appropriate
- Implement proper error handling
- Monitor goroutine count
- Use worker pools for controlled concurrency
- Implement graceful shutdown
- Test concurrent code thoroughly

This extended guide provides a comprehensive understanding of Go concurrency patterns with detailed examples and visual representations. The examples cover both basic and advanced use cases, making it suitable for beginners and experienced developers alike.
