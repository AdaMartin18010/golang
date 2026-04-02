# TS-CL-014: Advanced Go Channels Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #channels #concurrency #patterns #select
> **权威来源**:
>
> - [Go Concurrency Patterns](https://golang.org/doc/codewalk/sharemem/) - Go team
> - [Channels are First-Class](https://dave.cheney.net/2014/03/19/channel-axioms) - Dave Cheney

---

## 1. Channel Internals

### 1.1 Channel Structure

```go
// runtime channel structure
type hchan struct {
    qcount   uint           // Total data in queue
    dataqsiz uint           // Size of circular queue
    buf      unsafe.Pointer // Points to circular queue
    elemsize uint16         // Element size
    closed   uint32         // Is channel closed?
    elemtype *_type         // Element type
    sendx    uint           // Send index
    recvx    uint           // Receive index
    recvq    waitq          // List of recv waiters
    sendq    waitq          // List of send waiters
    lock     mutex          // Protects all fields
}

type waitq struct {
    first *sudog
    last  *sudog
}
```

### 1.2 Channel Operations

```
Buffered Channel (capacity N):
┌─────────────────────────────────────────────────────────────────┐
│                      Circular Buffer                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐       │
│  │   [0]    │   [1]    │   [2]    │   ...    │  [N-1]   │       │
│  │    A     │    B     │          │          │          │       │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘       │
│       ▲                             │                            │
│       │                             │                            │
│    recvx                          sendx                          │
│                                                                  │
│  qcount = 2 (A and B in buffer)                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

Send Operation:
1. If qcount < dataqsiz: write to buf[sendx], increment sendx, increment qcount
2. If qcount == dataqsiz: block on sendq until space available

Receive Operation:
1. If qcount > 0: read from buf[recvx], increment recvx, decrement qcount
2. If qcount == 0: block on recvq until data available
```

---

## 2. Channel Patterns

### 2.1 Pipeline Pattern

```go
// Generator produces numbers
func Generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Square squares numbers
func Square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Merge combines multiple channels
func Merge(cs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(cs))

    for _, c := range cs {
        go func(ch <-chan int) {
            defer wg.Done()
            for n := range ch {
                out <- n
            }
        }(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

// Usage
func main() {
    c1 := Generator(1, 2, 3)
    c2 := Generator(4, 5, 6)

    merged := Merge(c1, c2)
    squared := Square(merged)

    for n := range squared {
        fmt.Println(n)
    }
}
```

### 2.2 Worker Pool Pattern

```go
func WorkerPool(jobs <-chan Job, results chan<- Result, workerCount int) {
    var wg sync.WaitGroup
    wg.Add(workerCount)

    for i := 0; i < workerCount; i++ {
        go func(id int) {
            defer wg.Done()
            for job := range jobs {
                result := processJob(job)
                results <- result
            }
        }(i)
    }

    go func() {
        wg.Wait()
        close(results)
    }()
}

// Usage with bounded concurrency
func main() {
    jobs := make(chan Job, 100)
    results := make(chan Result, 100)

    // Start workers
    WorkerPool(jobs, results, 10)

    // Send jobs
    go func() {
        for i := 0; i < 1000; i++ {
            jobs <- Job{ID: i}
        }
        close(jobs)
    }()

    // Collect results
    for result := range results {
        fmt.Println(result)
    }
}
```

### 2.3 Fan-Out, Fan-In Pattern

```go
// Fan-out: Distribute work to multiple workers
func FanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)

    for i := 0; i < workers; i++ {
        out := make(chan int)
        outputs[i] = out

        go func() {
            defer close(out)
            for val := range input {
                out <- process(val)
            }
        }()
    }

    return outputs
}

// Fan-in: Combine multiple channels into one
func FanIn(inputs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(len(inputs))

    for _, in := range inputs {
        go func(ch <-chan int) {
            defer wg.Done()
            for val := range ch {
                out <- val
            }
        }(in)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### 2.4 Select Pattern

```go
// Timeout pattern
func WithTimeout(ch <-chan Result, timeout time.Duration) (Result, error) {
    select {
    case result := <-ch:
        return result, nil
    case <-time.After(timeout):
        return Result{}, errors.New("timeout")
    }
}

// Non-blocking send
func NonBlockingSend(ch chan<- int, value int) bool {
    select {
    case ch <- value:
        return true
    default:
        return false
    }
}

// Non-blocking receive
func NonBlockingReceive(ch <-chan int) (int, bool) {
    select {
    case value := <-ch:
        return value, true
    default:
        return 0, false
    }
}

// Graceful shutdown
func GracefulShutdown(ctx context.Context, work <-chan Job) {
    for {
        select {
        case job := <-work:
            processJob(job)
        case <-ctx.Done():
            // Drain remaining work
            for job := range work {
                processJob(job)
            }
            return
        }
    }
}
```

### 2.5 Quit Channel Pattern

```go
func Worker(jobs <-chan Job, quit <-chan struct{}, results chan<- Result) {
    for {
        select {
        case job := <-jobs:
            results <- processJob(job)
        case <-quit:
            // Clean shutdown
            return
        }
    }
}

// Alternative with done channel
func WorkerWithDone(jobs <-chan Job, done chan<- struct{}) {
    defer close(done)
    for job := range jobs {
        processJob(job)
    }
}
```

---

## 3. Common Mistakes

```go
// Mistake 1: Sending on closed channel
go func() {
    ch <- value // Panic if ch is closed
}()
close(ch) // DON'T

// Fix: Use select with done channel

// Mistake 2: Closing channel multiple times
close(ch)
close(ch) // Panic!

// Fix: Ensure only one close, or use sync.Once

// Mistake 3: Not checking closed channel
for {
    val := <-ch // Returns zero value forever after close
}

// Fix: Use range or check second return value
for val := range ch { // Automatically exits on close
}

// or
val, ok := <-ch
if !ok {
    // Channel closed
}

// Mistake 4: Unnecessary buffering
ch := make(chan int, 1000) // Too large

// Mistake 5: Goroutine leak from blocked sends
func Leak() {
    ch := make(chan int)
    go func() {
        ch <- 1 // Blocks forever if no receiver
    }()
    // Forgot to receive!
}
```

---

## 4. Checklist

```
Channel Patterns Checklist:
□ Use channels for goroutine communication
□ Prefer unbuffered channels for synchronization
□ Close channels from sender side only
□ Check closed status with ok idiom
□ Use select for multiple channel operations
□ Handle timeouts with time.After
□ Prevent goroutine leaks
□ Use context for cancellation
□ Document channel ownership (who closes)
```
