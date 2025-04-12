# Ultimate Guide to Go Concurrency: Context & Goroutines

This guide provides a comprehensive dive into Go's concurrency model, focusing on goroutines and context. It covers core concepts, practical patterns, real-world applications, and best practices, with clear examples and visual aids to enhance understanding.

## Table of Contents

- [Ultimate Guide to Go Concurrency: Context \& Goroutines](#ultimate-guide-to-go-concurrency-context--goroutines)
  - [Table of Contents](#table-of-contents)
  - [Core Concepts](#core-concepts)
    - [Goroutines](#goroutines)
    - [Context](#context)
  - [Visual Understanding](#visual-understanding)
    - [Goroutine Lifecycle](#goroutine-lifecycle)
    - [Channel Mechanics](#channel-mechanics)
  - [Basic Patterns](#basic-patterns)
    - [Worker Pool](#worker-pool)
    - [Rate Limiter](#rate-limiter)
  - [Advanced Patterns](#advanced-patterns)
    - [Circuit Breaker](#circuit-breaker)
    - [Pub/Sub System](#pubsub-system)
  - [Real-World Examples](#real-world-examples)
    - [HTTP Server with Graceful Shutdown](#http-server-with-graceful-shutdown)
    - [Concurrent Data Pipeline](#concurrent-data-pipeline)
  - [Debugging \& Troubleshooting](#debugging--troubleshooting)
    - [Goroutine Leaks](#goroutine-leaks)
    - [Context Debugging](#context-debugging)
  - [Best Practices \& Pitfalls](#best-practices--pitfalls)
    - [Resource Management](#resource-management)
    - [Error Handling](#error-handling)
    - [Testing Concurrent Code](#testing-concurrent-code)
  - [Key Takeaways](#key-takeaways)

---

## Core Concepts

### Goroutines

A **goroutine** is a lightweight thread managed by the Go runtime, not the OS. It allows concurrent execution with minimal overhead, making it ideal for scalable applications.

**Key Points:**

- Stack size starts small (2 KB) and grows dynamically.
- Scheduled by Go's runtime, not the OS kernel.
- Thousands or millions can run efficiently.

**Example:**

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        fmt.Println("Hello from goroutine!")
    }()
    time.Sleep(time.Second) // Wait for goroutine to finish
}
```

**Visual:**

```
OS Thread:     ────────────────────────►
Goroutines:    ─┬─┬─┬─┬─┬─┬─┬─┬─┬─┬─►
                G1 G2 G3 G4 G5 G6 G7 G8 G9 G10
```

**বাংলায় ব্যাখ্যা:**  
গোরুটিন হলো গোলাং-এর একটি হালকা থ্রেড, যা অপারেটিং সিস্টেমের পরিবর্তে গো রানটাইম দ্বারা নিয়ন্ত্রিত হয়। এটি অনেক কম মেমোরি ব্যবহার করে এবং হাজার হাজার গোরুটিন একসাথে চলতে পারে। এটি কল্পনা করো একটা ব্যস্ত রেস্টুরেন্টের মতো, যেখানে অনেক ওয়েটার (গোরুটিন) একসাথে কাজ করছে, কিন্তু তাদের নিয়ন্ত্রণ করছে রেস্টুরেন্ট ম্যানেজার (গো রানটাইম)।

**বাংলায় উদাহরণ:**  
ধরো, তুমি একটা দোকানে কাজ করছো। তুমি একা সব ক্রেতার অর্ডার নিতে পারছো না। তাই তুমি দুইজন সহকারী নিয়োগ করলে, যারা একই সাথে ক্রেতাদের সেবা দিচ্ছে। এখানে প্রতিটি সহকারী হলো গোরুটিন।

**কোড উদাহরণ:**

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        fmt.Println("সহকারী ১: ক্রেতার অর্ডার নিচ্ছি!")
    }()
    go func() {
        fmt.Println("সহকারী ২: ক্রেতার অর্ডার নিচ্ছি!")
    }()
    time.Sleep(time.Second) // সহকারীদের কাজ শেষ করার জন্য অপেক্ষা
    fmt.Println("প্রধান দোকানদার: সবাই কাজ শেষ করেছে!")
}
```

**আউটপুট:**

```
সহকারী ১: ক্রেতার অর্ডার নিচ্ছি!
সহকারী ২: ক্রেতার অর্ডার নিচ্ছি!
প্রধান দোকানদার: সবাই কাজ শেষ করেছে!
```

**মনে রাখার টিপস:**  
গোরুটিন মানে একসাথে অনেক কাজ করা, যেমন একটা দোকানে অনেক সহকারী একসাথে ক্রেতার সেবা দিচ্ছে। কিন্তু তাদের শেষ করার জন্য `time.Sleep` বা চ্যানেল দিয়ে অপেক্ষা করতে হয়।

### Context

The `context` package provides a way to manage cancellation, deadlines, and request-scoped values across goroutines.

**Key Features:**

- **Cancellation**: Signal goroutines to stop gracefully.
- **Deadlines/Timeouts**: Enforce time-bound operations.
- **Values**: Pass request-scoped data (use sparingly).

**Example:**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx := context.Background()

    // Cancellable context
    ctxCancel, cancel := context.WithCancel(ctx)
    defer cancel()

    // Timeout context
    ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 2*time.Second)
    defer cancelTimeout()

    // Value context
    ctxValue := context.WithValue(ctx, "userID", "12345")

    fmt.Println("Context setup complete")
}
```

**Visual:**

```
Context Tree:
                   Background
                       │
                  WithCancel
                 /    |     \
           Value   Timeout  Deadline
            |        |         |
         Cancel   Cancel    Cancel
```

**বাংলায় ব্যাখ্যা:**  
কনটেক্সট হলো গোলাং-এর একটি টুল, যা গোরুটিনগুলোর কাজ নিয়ন্ত্রণ করে। এটি দিয়ে তুমি গোরুটিনকে বলতে পারো কখন থামতে হবে, কতক্ষণ কাজ করতে হবে, বা কিছু তথ্য শেয়ার করতে পারো। এটি কল্পনা করো একটা টাইমারের মতো, যা কর্মীদের বলে, "এত সময়ের মধ্যে কাজ শেষ করো, নয়তো থেমে যাও।"

**বাংলায় উদাহরণ:**  
ধরো, তুমি একটা দোকানে একজন কর্মীকে বললে, "৫ মিনিটের মধ্যে ১০টা প্যাকেট তৈরি করো।" যদি সময় শেষ হয়ে যায়, তুমি তাকে থামিয়ে দিবে। এটাই কনটেক্সটের কাজ।

**কোড উদাহরণ:**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    go func() {
        fmt.Println("কর্মী: প্যাকেট তৈরি শুরু করছি...")
        select {
        case <-time.After(3 * time.Second): // ৩ সেকেন্ড সময় লাগবে
            fmt.Println("কর্মী: প্যাকেট তৈরি শেষ!")
        case <-ctx.Done():
            fmt.Println("কর্মী: সময় শেষ, থামছি!", ctx.Err())
        }
    }()

    time.Sleep(3 * time.Second) // অপেক্ষা
}
```

**আউটপুট:**

```
কর্মী: প্যাকেট তৈরি শুরু করছি...
কর্মী: সময় শেষ, থামছি! context deadline exceeded
```

**মনে রাখার টিপস:**  
কনটেক্সট মানে একটা নিয়ন্ত্রক, যেমন দোকানের ম্যানেজার যে কর্মীদের সময়সীমা দিয়ে কাজ করায়। `WithTimeout` বা `WithCancel` দিয়ে তুমি গোরুটিনের কাজ বন্ধ করতে পারো।

---

## Visual Understanding

### Goroutine Lifecycle

Goroutines follow a clear lifecycle from creation to termination.

**Visual:**

```
Creation:   main()──►go func()──►Running
              │
Execution:    │──►Processing──►Completed
              │
Termination:  └──►Done

Stack: [2KB]──►[Grows as needed]──►[Released]
```

**Example:**

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        fmt.Println("Goroutine started")
        time.Sleep(time.Second)
        fmt.Println("Goroutine done")
    }()
    time.Sleep(2 * time.Second)
}
```

**বাংলায় ব্যাখ্যা:**  
গোরুটিনের একটা জীবনচক্র আছে: শুরু, কাজ করা, এবং শেষ। এটি কল্পনা করো একটা কর্মীর মতো, যে দোকানে এসে কাজ শুরু করে, ক্রেতার অর্ডার নেয়, এবং কাজ শেষ করে চলে যায়।

**বাংলায় উদাহরণ:**

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        fmt.Println("কর্মী: দোকানে এসেছি, কাজ শুরু!")
        time.Sleep(time.Second)
        fmt.Println("কর্মী: কাজ শেষ, চললাম!")
    }()
    fmt.Println("ম্যানেজার: কর্মীকে কাজে পাঠিয়েছি")
    time.Sleep(2 * time.Second)
}
```

**আউটপুট:**

```
ম্যানেজার: কর্মীকে কাজে পাঠিয়েছি
কর্মী: দোকানে এসেছি, কাজ শুরু!
কর্মী: কাজ শেষ, চললাম!
```

**মনে রাখার টিপস:**  
গোরুটিন মানে কর্মী, যে শুরু থেকে শেষ পর্যন্ত কাজ করে। তবে তুমি তাকে অপেক্ষা না করলে সে কাজ না শেষ করেই চলে যেতে পারে। তাই `time.Sleep` বা চ্যানেল দরকার।

### Channel Mechanics

Channels enable safe communication and synchronization between goroutines.

**Types:**

- **Unbuffered**: Blocks until sender and receiver are ready.
- **Buffered**: Allows sending up to capacity without blocking.

**Visual:**

```
Unbuffered:
Producer ──[chan]──► Consumer
    │                   │
    └──────Sync───────┘

Buffered:
[Item1][Item2][ ][ ]
  Full  Full  Empty Empty
```

**Example:**

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan int, 2) // Buffered channel

    go func() {
        for i := 1; i <= 2; i++ {
            ch <- i
        }
        close(ch)
    }()

    for v := range ch {
        fmt.Println(v)
    }
}
```

**বাংলায় ব্যাখ্যা:**  
চ্যানেল হলো দুই গোরুটিনের মধ্যে পাইপ, যার মাধ্যমে তারা ডেটা শেয়ার করে। আনবাফারড চ্যানেল মানে হাতে হাতে ডেলিভারি, আর বাফারড চ্যানেল মানে মেলবক্সে চিঠি রাখা।

**বাংলায় উদাহরণ (আনবাফারড চ্যানেল):**

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string) // আনবাফারড চ্যানেল

    go func() {
        time.Sleep(1 * time.Second)
        ch <- "পণ্য ডেলিভারি সম্পন্ন!"
    }()

    fmt.Println("অপেক্ষা করছি পণ্যের জন্য...")
    msg := <-ch
    fmt.Println("বার্তা পেয়েছি:", msg)
}
```

**আউটপুট:**

```
অপেক্ষা করছি পণ্যের জন্য...
বার্তা পেয়েছি: পণ্য ডেলিভারি সম্পন্ন!
```

**বাংলায় উদাহরণ (বাফারড চ্যানেল):**

```go
package main

import (
    "fmt"
)

func main() {
    ch := make(chan string, 2) // বাফারড চ্যানেল, ক্যাপাসিটি ২

    ch <- "পণ্য ১"
    ch <- "পণ্য ২"
    fmt.Println("দুটো পণ্য পাঠানো হয়েছে")

    fmt.Println("প্রথম পণ্য গ্রহণ:", <-ch)
    fmt.Println("দ্বিতীয় পণ্য গ্রহণ:", <-ch)
}
```

**আউটপুট:**

```
দুটো পণ্য পাঠানো হয়েছে
প্রথম পণ্য গ্রহণ: পণ্য ১
দ্বিতীয় পণ্য গ্রহণ: পণ্য ২
```

**মনে রাখার টিপস:**

- আনবাফারড চ্যানেল মানে দোকানে সরাসরি পণ্য হস্তান্তর।
- বাফারড চ্যানেল মানে গুদামে পণ্য জমা রাখা, যেখানে সীমিত জায়গা আছে।

---

## Basic Patterns

### Worker Pool

A worker pool distributes tasks across a fixed number of goroutines, controlled by a context.

**Architecture:**

```
Jobs ──► [Job Queue] ──┬──► Worker1 ──┐
                      ├──► Worker2 ──┼──► [Results]
                      └──► Worker3 ──┘
```

**Example:**

```go
package main

import (
    "context"
    "fmt"
    "sync"
)

type Job struct {
    ID int
}

type Result struct {
    JobID int
    Value int
}

func workerPool(ctx context.Context, jobs <-chan Job, results chan<- Result) {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case job, ok := <-jobs:
                    if !ok {
                        return
                    }
                    results <- Result{JobID: job.ID, Value: job.ID * 2}
                }
            }
        }(i)
    }
    go func() {
        wg.Wait()
        close(results)
    }()
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    jobs := make(chan Job, 5)
    results := make(chan Result, 5)

    go workerPool(ctx, jobs, results)

    for i := 1; i <= 5; i++ {
        jobs <- Job{ID: i}
    }
    close(jobs)

    for r := range results {
        fmt.Printf("Job %d: %d\n", r.JobID, r.Value)
    }
}
```

**বাংলায় ব্যাখ্যা:**  
ওয়ার্কার পুল হলো এমন একটি পদ্ধতি, যেখানে কিছু নির্দিষ্ট সংখ্যক কর্মী (গোরুটিন) একসাথে কাজ ভাগ করে নেয়। এটি কল্পনা করো একটা কারখানার মতো, যেখানে ৩ জন কর্মী একটা লাইনে কাজ করে এবং প্রত্যেকে কাজ নিয়ে ফলাফল পাঠায়।

**বাংলায় উদাহরণ:**  
ধরো, একটা ফলের দোকানে ৩ জন কর্মী আছে। তারা ফল ধুয়ে, কেটে, এবং প্যাক করে। প্রত্যেকে একটা কাজ নেয়, শেষ করে ফলাফল (প্যাকেট) পাঠায়।

**কোড উদাহরণ:**

```go
package main

import (
    "context"
    "fmt"
    "sync"
)

type Job struct {
    ID int
}

type Result struct {
    JobID int
    Value string
}

func workerPool(ctx context.Context, jobs <-chan Job, results chan<- Result) {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    fmt.Printf("কর্মী %d: কাজ বন্ধ!\n", id)
                    return
                case job, ok := <-jobs:
                    if !ok {
                        fmt.Printf("কর্মী %d: কাজ শেষ!\n", id)
                        return
                    }
                    results <- Result{JobID: job.ID, Value: fmt.Sprintf("প্যাকেট %d তৈরি", job.ID)}
                }
            }
        }(i)
    }
    go func() {
        wg.Wait()
        close(results)
    }()
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    jobs := make(chan Job, 5)
    results := make(chan Result, 5)

    go workerPool(ctx, jobs, results)

    for i := 1; i <= 5; i++ {
        jobs <- Job{ID: i}
    }
    close(jobs)

    for r := range results {
        fmt.Printf("ফলাফল: কাজ %d -> %s\n", r.JobID, r.Value)
    }
}
```

**আউটপুট:**

```
ফলাফল: কাজ 1 -> প্যাকেট 1 তৈরি
ফলাফল: কাজ 2 -> প্যাকেট 2 তৈরি
ফলাফল: কাজ 3 -> প্যাকেট 3 তৈরি
ফলাফল: কাজ 4 -> প্যাকেট 4 তৈরি
ফলাফল: কাজ 5 -> প্যাকেট 5 তৈরি
কর্মী 0: কাজ শেষ!
কর্মী 1: কাজ শেষ!
কর্মী 2: কাজ শেষ!
```

**মনে রাখার টিপস:**  
ওয়ার্কার পুল মানে কারখানার লাইন, যেখানে সীমিত কর্মী কাজ ভাগ করে নেয়। কনটেক্সট দিয়ে তুমি তাদের থামাতে পারো।

### Rate Limiter

Rate limiting controls the frequency of operations, useful for APIs or resource-intensive tasks.

**Visual:**

```
Requests: R R R R R
Time:     │ │ │ │ │
Tokens:   [T][T][ ]
```

**Example:**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type RateLimiter struct {
    tokens chan struct{}
    stop   chan struct{}
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, burst),
        stop:   make(chan struct{}),
    }
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
            case <-rl.stop:
                return
            }
        }
    }()
    return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-rl.tokens:
        return nil
    }
}

func (rl *RateLimiter) Stop() {
    close(rl.stop)
}

func main() {
    rl := NewRateLimiter(500*time.Millisecond, 2)
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    defer rl.Stop()

    for i := 1; i <= 5; i++ {
        if err := rl.Wait(ctx); err != nil {
            fmt.Println("Context timeout:", err)
            return
        }
        fmt.Printf("Request %d processed at %v\n", i, time.Now())
    }
}
```

**বাংলায় ব্যাখ্যা:**  
রেট লিমিটার কাজের গতি নিয়ন্ত্রণ করে। এটি কল্পনা করো একটা টাইমারের মতো, যেখানে প্রতি সেকেন্ডে শুধু একটা গাড়ি যেতে পারে।

**বাংলায় উদাহরণ:**  
ধরো, তুমি একটা দোকানে আছো, যেখানে প্রতি ৫ সেকেন্ডে একজন ক্রেতাকে সেবা দেওয়া যায়। যদি অনেক ক্রেতা আসে, তবে তাদের অপেক্ষা করতে হবে।

**কোড উদাহরণ:**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type RateLimiter struct {
    tokens chan struct{}
    stop   chan struct{}
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, burst),
        stop:   make(chan struct{}),
    }
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
            case <-rl.stop:
                return
            }
        }
    }()
    return rl
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-rl.tokens:
        return nil
    }
}

func (rl *RateLimiter) Stop() {
    close(rl.stop)
}

func main() {
    rl := NewRateLimiter(500*time.Millisecond, 2)
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    defer rl.Stop()

    for i := 1; i <= 5; i++ {
        if err := rl.Wait(ctx); err != nil {
            fmt.Println("সময় শেষ:", err)
            return
        }
        fmt.Printf("ক্রেতা %d কে সেবা দেওয়া হলো: %v\n", i, time.Now())
    }
}
```

**আউটপুট (উদাহরণ):**

```
ক্রেতা 1 কে সেবা দেওয়া হলো: 2025-04-11 12:00:00
ক্রেতা 2 কে সেবা দেওয়া হলো: 2025-04-11 12:00:00
ক্রেতা 3 কে সেবা দেওয়া হলো: 2025-04-11 12:00:00.5
ক্রেতা 4 কে সেবা দেওয়া হলো: 2025-04-11 12:00:01
সময় শেষ: context deadline exceeded
```

**মনে রাখার টিপস:**  
রেট লিমিটার মানে সিগন্যাল, যা কাজের গতি ধীর করে। এটি ব্যবহার করো যখন তুমি সার্ভার বা রিসোর্সের উপর চাপ কমাতে চাও।

---

## Advanced Patterns

### Circuit Breaker

A circuit breaker prevents cascading failures by halting operations when errors exceed a threshold.

**States:**

```
Closed ──► Half-Open ──► Open
   │           │           │
   └───────────┴───────────┘
```

**Example:**

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

type CircuitBreaker struct {
    mu          sync.RWMutex
    state       string
    failures    int
    threshold   int
    resetAfter  time.Duration
    lastFailure time.Time
}

func NewCircuitBreaker(threshold int, resetAfter time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:      "Closed",
        threshold:  threshold,
        resetAfter: resetAfter,
    }
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    cb.mu.RLock()
    if cb.state == "Open" && time.Since(cb.lastFailure) < cb.resetAfter {
        cb.mu.RUnlock()
        return errors.New("circuit breaker open")
    }
    cb.mu.RUnlock()

    cb.mu.Lock()
    if cb.state == "Open" && time.Since(cb.lastFailure) >= cb.resetAfter {
        cb.state = "Half-Open"
    }
    cb.mu.Unlock()

    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.threshold {
            cb.state = "Open"
        }
        return err
    }

    if cb.state == "Half-Open" {
        cb.state = "Closed"
        cb.failures = 0
    }
    return nil
}

func main() {
    cb := NewCircuitBreaker(2, 5*time.Second)
    ctx := context.Background()

    failingFn := func() error {
        return errors.New("service failed")
    }

    for i := 1; i <= 5; i++ {
        err := cb.Execute(ctx, failingFn)
        fmt.Printf("Attempt %d: %v\n", i, err)
        time.Sleep(1 * time.Second)
    }
}
```

**বাংলায় ব্যাখ্যা:**  
সার্কিট ব্রেকার একটা সুরক্ষা ব্যবস্থা, যা বারবার ব্যর্থতা হলে কাজ বন্ধ করে দেয়। এটি কল্পনা করো একটা বৈদ্যুতিক সার্কিটের মতো, যেখানে অনেক সমস্যা হলে ফিউজ বন্ধ হয়ে যায়।

**বাংলায় উদাহরণ:**  
ধরো, তুমি একটা দোকানে অনলাইন অর্ডার পাঠাচ্ছ। যদি সার্ভার ২ বার ব্যর্থ হয়, তবে তুমি কিছুক্ষণের জন্য অর্ডার বন্ধ করে দেবে।

**কোড উদাহরণ:**

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

type CircuitBreaker struct {
    mu          sync.RWMutex
    state       string
    failures    int
    threshold   int
    resetAfter  time.Duration
    lastFailure time.Time
}

func NewCircuitBreaker(threshold int, resetAfter time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:      "Closed",
        threshold:  threshold,
        resetAfter: resetAfter,
    }
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    cb.mu.RLock()
    if cb.state == "Open" && time.Since(cb.lastFailure) < cb.resetAfter {
        cb.mu.RUnlock()
        return errors.New("সার্কিট ব্রেকার খোলা")
    }
    cb.mu.RUnlock()

    cb.mu.Lock()
    if cb.state == "Open" && time.Since(cb.lastFailure) >= cb.resetAfter {
        cb.state = "Half-Open"
    }
    cb.mu.Unlock()

    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.threshold {
            cb.state = "Open"
        }
        return err
    }

    if cb.state == "Half-Open" {
        cb.state = "Closed"
        cb.failures = 0
    }
    return nil
}

func main() {
    cb := NewCircuitBreaker(2, 5*time.Second)
    ctx := context.Background()

    failingFn := func() error {
        return errors.New("সার্ভার ব্যর্থ")
    }

    for i := 1; i <= 5; i++ {
        err := cb.Execute(ctx, failingFn)
        fmt.Printf("চেষ্টা %d: %v\n", i, err)
        time.Sleep(1 * time.Second)
    }
}
```

**আউটপুট:**

```
চেষ্টা 1: সার্ভার ব্যর্থ
চেষ্টা 2: সার্ভার ব্যর্থ
চেষ্টা 3: সার্কিট ব্রেকার খোলা
চেষ্টা 4: সার্কিট ব্রেকার খোলা
চেষ্টা 5: সার্কিট ব্রেকার খোলা
```

**মনে রাখার টিপস:**  
সার্কিট ব্রেকার মানে ফিউজ, যা ব্যর্থতা থেকে সিস্টেমকে বাঁচায়। এটি ব্যবহার করো যখন তুমি সার্ভারের উপর বারবার চাপ দিতে চাও না।

### Pub/Sub System

A publish/subscribe system allows multiple subscribers to receive messages from a topic.

**Architecture:**

```
Publisher ──► [Topic] ──┬──► Sub1
                       ├──► Sub2
                       └──► Sub3
```

**Example:**

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type PubSub struct {
    mu   sync.RWMutex
    subs map[string][]chan string
}

func NewPubSub() *PubSub {
    return &PubSub{
        subs: make(map[string][]chan string),
    }
}

func (ps *PubSub) Subscribe(topic string) <-chan string {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    ch := make(chan string, 1)
    ps.subs[topic] = append(ps.subs[topic], ch)
    return ch
}

func (ps *PubSub) Publish(topic, msg string) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()
    for _, ch := range ps.subs[topic] {
        ch <- msg
    }
}

func (ps *PubSub) Close() {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    for _, subs := range ps.subs {
        for _, ch := range subs {
            close(ch)
        }
    }
    ps.subs = nil
}

func main() {
    ps := NewPubSub()
    topic := "news"

    sub1 := ps.Subscribe(topic)
    sub2 := ps.Subscribe(topic)

    go func() {
        for msg := range sub1 {
            fmt.Printf("Sub1 received: %s\n", msg)
        }
    }()
    go func() {
        for msg := range sub2 {
            fmt.Printf("Sub2 received: %s\n", msg)
        }
    }()

    ps.Publish(topic, "Breaking News!")
    time.Sleep(100 * time.Millisecond)
    ps.Close()
}
```

**বাংলায় ব্যাখ্যা:**  
পাব/সাব সিস্টেমে একজন প্রকাশক (পাবলিশার) বার্তা পাঠায়, এবং অনেক গ্রাহক (সাবস্ক্রাইবার) তা গ্রহণ করে। এটি কল্পনা করো একটা রেডিও স্টেশনের মতো, যেখানে একজন ঘোষক বার্তা দেয়, এবং অনেক শ্রোতা তা শোনে।

**বাংলায় উদাহরণ:**  
ধরো, একটা নিউজ চ্যানেলে খবর প্রকাশিত হয়। যারা সাবস্ক্রাইব করেছে (যেমন, ফোন অ্যাপে নোটিফিকেশন চালু করেছে), তারা খবর পায়।

**কোড উদাহরণ:**

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type PubSub struct {
    mu   sync.RWMutex
    subs map[string][]chan string
}

func NewPubSub() *PubSub {
    return &PubSub{
        subs: make(map[string][]chan string),
    }
}

func (ps *PubSub) Subscribe(topic string) <-chan string {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    ch := make(chan string, 1)
    ps.subs[topic] = append(ps.subs[topic], ch)
    return ch
}

func (ps *PubSub) Publish(topic, msg string) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()
    for _, ch := range ps.subs[topic] {
        ch <- msg
    }
}

func (ps *PubSub) Close() {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    for _, subs := range ps.subs {
        for _, ch := range subs {
            close(ch)
        }
    }
    ps.subs = nil
}

func main() {
    ps := NewPubSub()
    topic := "খবর"

    sub1 := ps.Subscribe(topic)
    sub2 := ps.Subscribe(topic)

    go func() {
        for msg := range sub1 {
            fmt.Printf("গ্রাহক ১ পেয়েছে: %s\n", msg)
        }
    }()
    go func() {
        for msg := range sub2 {
            fmt.Printf("গ্রাহক ২ পেয়েছে: %s\n", msg)
        }
    }()

    ps.Publish(topic, "জরুরি খবর: বৃষ্টি হবে!")
    time.Sleep(100 * time.Millisecond)
    ps.Close()
}
```

**আউটপুট:**

```
গ্রাহক ১ পেয়েছে: জরুরি খবর: বৃষ্টি হবে!
গ্রাহক ২ পেয়েছে: জরুরি খবর: বৃষ্টি হবে!
```

**মনে রাখার টিপস:**  
পাব/সাব মানে রেডিও ব্রডকাস্ট। একজন পাঠায়, অনেকে পায়। চ্যানেল ব্যবহার করে নিরাপদে বার্তা পাঠাও।

---

## Real-World Examples

### HTTP Server with Graceful Shutdown

An HTTP server that shuts down gracefully, ensuring in-flight requests complete.

**Example:**

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    srv := &http.Server{
        Addr: ":8080",
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            time.Sleep(2 * time.Second) // Simulate work
            w.Write([]byte("Hello, World!"))
        }),
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig

        shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
        defer shutdownCancel()

        if err := srv.Shutdown(shutdownCtx); err != nil {
            log.Printf("Shutdown error: %v", err)
        }
        cancel()
    }()

    log.Println("Server starting on :8080")
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Server error: %v", err)
    }

    <-ctx.Done()
    log.Println("Server stopped")
}
```

**বাংলায় ব্যাখ্যা:**  
এইচটিটিপি সার্ভার গ্রেসফুল শাটডাউন দিয়ে বন্ধ হয়, যাতে চলমান কাজগুলো শেষ হয়। এটি কল্পনা করো একটা দোকানের মতো, যেখানে দোকান বন্ধের আগে সব ক্রেতার অর্ডার শেষ করা হয়।

**বাংলায় উদাহরণ:**  
ধরো, তুমি একটা অনলাইন দোকান চালাচ্ছ। যখন তুমি দোকান বন্ধ করতে চাও, তখন প্রথমে চলমান অর্ডারগুলো ডেলিভারি করে তারপর সার্ভার বন্ধ করো।

**কোড উদাহরণ:**

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    srv := &http.Server{
        Addr: ":8080",
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            time.Sleep(2 * time.Second) // কাজ প্রক্রিয়াকরণের সময়
            w.Write([]byte("স্বাগতম, বিশ্ব!"))
        }),
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig

        shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
        defer shutdownCancel()

        if err := srv.Shutdown(shutdownCtx); err != nil {
            log.Printf("শাটডাউন ত্রুটি: %v", err)
        }
        cancel()
    }()

    log.Println("সার্ভার শুরু হচ্ছে :8080 পোর্টে")
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("সার্ভার ত্রুটি: %v", err)
    }

    <-ctx.Done()
    log.Println("সার্ভার বন্ধ হয়েছে")
}
```

**মনে রাখার টিপস:**  
গ্রেসফুল শাটডাউন মানে দোকান বন্ধের আগে সব ক্রেতার কাজ শেষ করা। কনটেক্সট এবং সিগন্যাল দিয়ে তুমি এটি নিয়ন্ত্রণ করতে পারো।

### Concurrent Data Pipeline

A pipeline processes data in stages concurrently, such as reading, transforming, and writing.

**Example:**

```go
package main

import (
    "context"
    "fmt"
    "sync"
)

func generator(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case <-ctx.Done():
                return
            case out <- n:
            }
        }
    }()
    return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case <-ctx.Done():
                return
            case out <- n * n:
            }
        }
    }()
    return out
}

func printer(ctx context.Context, in <-chan int) {
    for n := range in {
        select {
        case <-ctx.Done():
            return
        default:
            fmt.Println(n)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    gen := generator(ctx, 1, 2, 3, 4)
    squared := square(ctx, gen)

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        printer(ctx, squared)
    }()

    wg.Wait()
}
```

**বাংলায় ব্যাখ্যা:**  
ডেটা পাইপলাইন হলো এমন একটি প্রক্রিয়া, যেখানে ডেটা ধাপে ধাপে প্রক্রিয়া করা হয়। এটি কল্পনা করো একটা কারখানার মতো, যেখানে ফল ধোয়া, কাটা, এবং প্যাক করা হয়।

**বাংলায় উদাহরণ:**  
ধরো, তুমি ফল প্রক্রিয়া করছো। প্রথমে ফল সংগ্রহ করো, তারপর ধুয়ে, কেটে, এবং প্যাক করো। প্রতিটি ধাপ একটা গোরুটিন দিয়ে করা হয়।

**কোড উদাহরণ:**

```go
package main

import (
    "context"
    "fmt"
    "sync"
)

func generator(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case <-ctx.Done():
                return
            case out <- n:
            }
        }
    }()
    return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case <-ctx.Done():
                return
            case out <- n * n:
            }
        }
    }()
    return out
}

func printer(ctx context.Context, in <-chan int) {
    for n := range in {
        select {
        case <-ctx.Done():
            return
        default:
            fmt.Println("ফলাফল:", n)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    gen := generator(ctx, 1, 2, 3, 4)
    squared := square(ctx, gen)

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        printer(ctx, squared)
    }()

    wg.Wait()
}
```

**আউটপুট:**

```
ফলাফল: 1
ফলাফল: 4
ফলাফল: 9
ফলাফল: 16
```

**মনে রাখার টিপস:**  
পাইপলাইন মানে কারখানার লাইন। প্রতিটি ধাপ একটা গোরুটিন, এবং চ্যানেল দিয়ে ডেটা পাঠানো হয়।

---

## Debugging & Troubleshooting

### Goroutine Leaks

Goroutine leaks occur when goroutines are created but never terminate, consuming resources.

**Detection Example:**

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    before := runtime.NumGoroutine()
    fmt.Printf("Goroutines before: %d\n", before)

    // Simulate leak
    go func() {
        <-make(chan struct{}) // Never closes
    }()

    time.Sleep(time.Second)
    after := runtime.NumGoroutine()
    fmt.Printf("Goroutines after: %d\n", after)

    if after > before {
        fmt.Printf("Leak detected: %d extra goroutines\n", after-before)
    }
}
```

**Prevention Tips:**

- Always close channels when done.
- Use context cancellation to terminate goroutines.
- Ensure all paths in a goroutine lead to exit.

### Context Debugging

Debugging context issues involves checking for cancellation or timeout errors.

**Example:**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func debugContext(ctx context.Context) {
    select {
    case <-ctx.Done():
        switch ctx.Err() {
        case context.Canceled:
            fmt.Println("Context canceled")
        case context.DeadlineExceeded:
            fmt.Println("Context deadline exceeded")
        default:
            fmt.Println("Unknown context error")
        }
    default:
        fmt.Println("Context active")
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    debugContext(ctx)
    time.Sleep(2 * time.Second)
    debugContext(ctx)
}
```

---

## Best Practices & Pitfalls

### Resource Management

**Good Practice:**

```go
package main

import (
    "fmt"
)

type Work struct {
    ID int
}

func process(w Work) {
    fmt.Printf("Processing work %d\n", w.ID)
}

func main() {
    workQueue := make(chan Work, 10)
    const maxWorkers = 3
    sem := make(chan struct{}, maxWorkers)

    go func() {
        for i := 1; i <= 5; i++ {
            workQueue <- Work{ID: i}
        }
        close(workQueue)
    }()

    for w := range workQueue {
        sem <- struct{}{}
        go func(work Work) {
            defer func() { <-sem }()
            process(work)
        }(w)
    }

    // Wait for all workers to finish
    for i := 0; i < maxWorkers; i++ {
        sem <- struct{}{}
    }
}
```

**Pitfall (Leak):**

```go
// Bad: Captures loop variable
for w := range workQueue {
    go func() {
        process(w) // Uses shared w, causing race
    }()
}
```

### Error Handling

Wrap goroutines with panic recovery and proper error propagation.

**Example:**

```go
package main

import (
    "context"
    "errors"
    "fmt"
)

type Result struct {
    Value int
    Err   error
}

func safeGoroutine(ctx context.Context, id int) <-chan Result {
    ch := make(chan Result, 1)
    go func() {
        defer func() {
            if r := recover(); r != nil {
                ch <- Result{Err: fmt.Errorf("panic: %v", r)}
            }
            close(ch)
        }()
        select {
        case <-ctx.Done():
            ch <- Result{Err: ctx.Err()}
        default:
            if id < 0 {
                ch <- Result{Err: errors.New("negative ID")}
                return
            }
            ch <- Result{Value: id * 2}
        }
    }()
    return ch
}

func main() {
    ctx := context.Background()
    results := safeGoroutine(ctx, 5)
    r := <-results
    fmt.Printf("Result: %v, Error: %v\n", r.Value, r.Err)
}
```

### Testing Concurrent Code

Testing concurrent code requires handling race conditions and ensuring deterministic outcomes.

**Example:**

```go
package main

import (
    "sync"
    "testing"
)

func TestConcurrentCounter(t *testing.T) {
    var counter int
    var mu sync.Mutex
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }

    wg.Wait()
    if counter != 100 {
        t.Errorf("Expected counter to be 100, got %d", counter)
    }
}
```

**Tips:**

- Use `-race` flag to detect data races (`go test -race`).
- Use `sync.WaitGroup` for synchronization.
- Test timeouts and cancellations with `context`.

---

## Key Takeaways

- **Goroutines**: Use for lightweight concurrency, but monitor for leaks.
- **Context**: Essential for cancellation, timeouts, and graceful shutdown.
- **Channels**: Prefer for synchronization over shared memory.
- **Patterns**: Worker pools, rate limiters, and circuit breakers solve common problems.
- **Debugging**: Use runtime metrics and context checks to diagnose issues.
- **Best Practices**:
  - Always defer resource cleanup.
  - Handle errors and panics explicitly.
  - Test concurrency thoroughly with race detection.

_Last Updated: April 11, 2025_

```

```
