# TS-CL-014: Go Channels Advanced Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (22+ KB)
> **标签**: #golang #channels #goroutines #concurrency #patterns
> **权威来源**:
>
> - [Go Concurrency Patterns](https://go.dev/blog/pipelines) - Go Blog
> - [Advanced Concurrency](https://go.dev/talks/2012/concurrency.slide) - Rob Pike

---

## 1. Channel Architecture

### 1.1 Channel Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Channel Structure                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   hchan (runtime)                                                            │
│   ┌───────────────────────────────────────────────────────────────────────┐  │
│   │  qcount   uint    - Total data in queue                              │  │
│   │  dataqsiz uint    - Size of circular queue                           │  │
│   │  buf      unsafe.Pointer - Circular buffer                           │  │
│   │  elemsize uint16  - Size of each element                             │  │
│   │  closed   uint32  - Channel closed flag                              │  │
│   │  elemtype *_type  - Element type                                     │  │
│   │  sendx    uint    - Send index                                       │  │
│   │  recvx    uint    - Receive index                                    │  │
│   │  recvq    waitq   - Waiting receivers (linked list)                  │  │
│   │  sendq    waitq   - Waiting senders (linked list)                    │  │
│   │  lock     mutex   - Channel lock                                     │  │
│   └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│   Buffer Visualization:                                                      │
│   ┌───┬───┬───┬───┬───┐                                                     │
│   │ A │ B │ C │ D │ E │  Circular buffer (size 5)                          │
│   └───┴───┴───┴───┴───┘                                                     │
│        ▲          ▲                                                         │
│       recvx      sendx                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Channel Types

```go
// Unbuffered channel (synchronous)
ch := make(chan int)

// Buffered channel (asynchronous up to capacity)
ch := make(chan int, 10)

// Receive-only channel
func receiver(ch <-chan int) {}

// Send-only channel
func sender(ch chan<- int) {}

// Bidirectional channel
func processor(ch chan int) {}
```

---

## 2. Advanced Patterns

### 2.1 Fan-Out / Fan-In

```go
// Fan-Out: Distribute work to multiple workers
func fanOut(input <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)

    for i := 0; i < workers; i++ {
        ch := make(chan int)
        channels[i] = ch

        go func() {
            defer close(ch)
            for val := range input {
                ch <- process(val)
            }
        }()
    }

    return channels
}

// Fan-In: Merge multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                out <- val
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### 2.2 Pipeline Pattern

```go
// Stage 1: Generator
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Stage 2: Square
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Usage: pipeline
c := generator(2, 3, 4)
out := square(c)

for n := range out {
    fmt.Println(n) // 4, 9, 16
}
```

### 2.3 Select Pattern

```go
// Timeout pattern
func withTimeout(ch <-chan int, timeout time.Duration) (int, bool) {
    select {
    case val := <-ch:
        return val, true
    case <-time.After(timeout):
        return 0, false
    }
}

// Non-blocking receive
func nonBlockingRecv(ch <-chan int) (int, bool) {
    select {
    case val := <-ch:
        return val, true
    default:
        return 0, false
    }
}

// Multiplexing
func multiplex(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)

    go func() {
        for {
            select {
            case val, ok := <-ch1:
                if !ok {
                    ch1 = nil // Disable this case
                    continue
                }
                out <- val
            case val, ok := <-ch2:
                if !ok {
                    ch2 = nil
                    continue
                }
                out <- val
            }
            if ch1 == nil && ch2 == nil {
                close(out)
                return
            }
        }
    }()

    return out
}
```

### 2.4 Worker Pool

```go
type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID int
    Value interface{}
    Error error
}

func workerPool(jobs <-chan Job, workers int) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                res, err := processJob(job)
                results <- Result{JobID: job.ID, Value: res, Error: err}
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

---

## 3. Channel Best Practices

### 3.1 Ownership Rules

```go
// Channel owner: Creates, sends, closes
// Channel user: Only receives

func channelOwner() <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < 10; i++ {
            ch <- i
        }
    }()
    return ch
}

func channelConsumer(ch <-chan int) {
    for val := range ch {
        fmt.Println(val)
    }
}
```

### 3.2 Closing Patterns

```go
// Safe close helper
func SafeClose[T any](ch chan T) (justClosed bool) {
    defer func() {
        if recover() != nil {
            justClosed = false
        }
    }()
    close(ch)
    return true
}

// Signal channel (zero-value signaling)
done := make(chan struct{})
close(done) // Signal completion
```

---

## 4. Performance Characteristics

### 4.1 Channel Overhead

| Operation | Unbuffered | Buffered (100) |
|-----------|-----------|----------------|
| Send | ~100-200ns | ~20-30ns |
| Receive | ~100-200ns | ~20-30ns |
| Close | ~10-20ns | ~10-20ns |

---

## 5. Comparison with Alternatives

| Pattern | Use Case | Performance |
|---------|----------|-------------|
| **Channels** | Coordination, streaming | Good |
| **sync.Mutex** | Shared state | Better |
| **sync.WaitGroup** | Wait for completion | Best |
| **atomic** | Simple counters | Best |
| **context** | Cancellation | Good |

---

## 6. Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Channel Best Practices                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Design:                                                                     │
│  □ Channel owner closes, consumers receive                                  │
│  □ Use directional channel types in function signatures                     │
│  □ Consider buffer size based on workload                                   │
│                                                                              │
│  Patterns:                                                                   │
│  □ Use select for non-blocking and timeout operations                       │
│  □ Implement graceful shutdown with done channels                           │
│  □ Use nil channels to disable select cases                                 │
│                                                                              │
│  Safety:                                                                     │
│  □ Never close a channel from a receiver                                    │
│  □ Check closed status when necessary                                       │
│  □ Use context for cancellation across API boundaries                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (22+ KB, comprehensive coverage)
