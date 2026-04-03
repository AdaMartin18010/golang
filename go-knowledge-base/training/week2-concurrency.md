# Week 2: Go Concurrency Deep Dive

## Module Overview

**Duration:** 40 hours (5 days)  
**Prerequisites:** Week 1 completion (Go Fundamentals)  
**Learning Goal:** Master Go's concurrency model and implement robust concurrent systems

---

## Learning Objectives

By the end of this week, you will be able to:

1. **CSP Fundamentals**
   - Explain Communicating Sequential Processes (CSP) theory
   - Contrast Go's concurrency model with thread-based models
   - Understand the happens-before relationship

2. **Goroutine Mastery**
   - Manage goroutine lifecycles properly
   - Understand the Go scheduler (GMP model)
   - Avoid goroutine leaks

3. **Channel Patterns**
   - Implement common channel patterns (pipelines, fan-out, fan-in)
   - Use select effectively for multiple channel operations
   - Handle channel closing and cancellation

4. **Synchronization**
   - Use sync package primitives appropriately
   - Implement lock-free algorithms where beneficial
   - Apply atomic operations correctly

5. **Concurrency Debugging**
   - Detect and fix race conditions
   - Use the race detector effectively
   - Profile concurrent applications
   - Debug deadlocks

---

## Reading Assignments

### Required Reading (Complete by Day 3)

1. **[CSP Theory](../01-Formal-Theory/03-Concurrency-Models/01-CSP-Theory.md)**
   - Study: Hoare's CSP, process algebra
   - Understand: Synchronous communication
   - Learn: Channel-based communication benefits

2. **[Go Concurrency Semantics](../01-Formal-Theory/03-Concurrency-Models/02-Go-Concurrency-Semantics.md)**
   - Master: Happens-before relationships
   - Learn: Memory visibility guarantees
   - Study: Channel synchronization semantics

3. **[GMP Scheduler Deep Dive](../02-Language-Design/FT-002-GMP-Scheduler-Deep-Dive.md)**
   - Understand: Goroutines, Machine threads, Processors
   - Learn: Work stealing algorithm
   - Study: Scheduler optimizations

4. **[Go Channels Deep Dive](../02-Language-Design/LD-004-Go-Channels-Formal.md)**
   - Master: Buffered vs unbuffered channels
   - Learn: Channel implementation internals
   - Study: Select statement mechanics

5. **[Sync Package Internals](../02-Language-Design/30-Go-sync-Package-Internals.md)**
   - Understand: Mutex, RWMutex implementation
   - Learn: WaitGroup and Once patterns
   - Study: Pool and Map usage

### Supplementary Reading (Complete by Day 5)

6. **[Go Memory Model](../01-Formal-Theory/04-Memory-Models/01-Happens-Before.md)**
   - Learn: Data races vs race conditions
   - Understand: Safe publication patterns
   - Study: Atomic visibility guarantees

7. **[Lock-Free Programming](../03-Engineering-CloudNative/03-Performance/06-Lock-Free-Programming.md)**
   - Understand: Compare-and-swap (CAS)
   - Learn: Atomic pointer operations
   - Study: When to use lock-free approaches

8. **[Race Detection Guide](../03-Engineering-CloudNative/03-Performance/04-Race-Detection.md)**
   - Master: Using the race detector
   - Learn: Common race patterns
   - Study: False positives and negatives

---

## Hands-on Exercises

### Day 1: Goroutine Fundamentals

#### Exercise 1.1: Goroutine Lifecycle (2 hours)

Understand goroutine creation and management:

```go
package concurrency

import (
    "context"
    "fmt"
    "runtime"
    "sync"
    "time"
)

// GoroutineTracker tracks active goroutines
type GoroutineTracker struct {
    wg     sync.WaitGroup
    ctx    context.Context
    cancel context.CancelFunc
}

func NewGoroutineTracker() *GoroutineTracker {
    ctx, cancel := context.WithCancel(context.Background())
    return &GoroutineTracker{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (gt *GoroutineTracker) Start(fn func(context.Context)) {
    gt.wg.Add(1)
    go func() {
        defer gt.wg.Done()
        fn(gt.ctx)
    }()
}

func (gt *GoroutineTracker) Stop() {
    gt.cancel()
    gt.wg.Wait()
}

// Exercise 1.1.1: Proper goroutine cleanup
func DemonstrateGoroutineLeak() {
    // BAD: Goroutine leak - never terminates
    ch := make(chan int)
    go func() {
        for val := range ch {
            fmt.Println(val)
        }
    }()
    
    // Channel never closed, goroutine leaked
    ch <- 1
    // ... never close ch
}

func DemonstrateProperCleanup() {
    // GOOD: Using context for cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    ch := make(chan int)
    go func() {
        for {
            select {
            case val := <-ch:
                fmt.Println(val)
            case <-ctx.Done():
                return // Proper cleanup
            }
        }
    }()
    
    ch <- 1
    cancel() // Signal goroutine to stop
}
```

**Tasks:**
1. Create a goroutine pool with size limit
2. Implement graceful shutdown for multiple goroutines
3. Profile goroutine count under load
4. Detect and fix goroutine leaks

**Deliverable:** Goroutine management package with leak detection

#### Exercise 1.2: Scheduler Observation (2 hours)

Understand GOMAXPROCS and scheduling:

```go
package concurrency

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func DemonstrateScheduling() {
    // Observe scheduler behavior
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
    fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
    
    // Create CPU-intensive work
    var wg sync.WaitGroup
    workers := runtime.GOMAXPROCS(0) * 2
    
    start := time.Now()
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // CPU-intensive calculation
            var result int
            for j := 0; j < 1e9; j++ {
                result += j
            }
            fmt.Printf("Worker %d completed, result: %d\n", id, result)
        }(i)
    }
    
    wg.Wait()
    fmt.Printf("Total time: %v\n", time.Since(start))
}

// Demonstrate preemption
func DemonstratePreemption() {
    runtime.GOMAXPROCS(1) // Force single thread
    
    var wg sync.WaitGroup
    wg.Add(2)
    
    // Goroutine 1: CPU-intensive
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            fmt.Println("Worker 1:", i)
            // Simulate CPU work
            for j := 0; j < 1e8; j++ {}
        }
    }()
    
    // Goroutine 2: Frequent yields
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            fmt.Println("Worker 2:", i)
            time.Sleep(1 * time.Millisecond) // Yield
        }
    }()
    
    wg.Wait()
}
```

**Deliverable:** Scheduler behavior analysis report

---

### Day 2: Channel Patterns

#### Exercise 2.1: Pipeline Pattern (3 hours)

Build a data processing pipeline:

```go
package pipeline

import (
    "context"
    "fmt"
    "runtime"
)

// Stage represents a pipeline stage
type Stage func(context.Context, <-chan int) <-chan int

// Generator creates the source channel
func Generator(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Filter removes numbers not satisfying predicate
func Filter(predicate func(int) bool) Stage {
    return func(ctx context.Context, in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                if predicate(n) {
                    select {
                    case out <- n:
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }()
        return out
    }
}

// Transform applies function to each element
func Transform(fn func(int) int) Stage {
    return func(ctx context.Context, in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                select {
                case out <- fn(n):
                case <-ctx.Done():
                    return
                }
            }
        }()
        return out
    }
}

// Reduce aggregates values
func Reduce(fn func(int, int) int, initial int) Stage {
    return func(ctx context.Context, in <-chan int) <-chan int {
        out := make(chan int, 1)
        go func() {
            defer close(out)
            result := initial
            for n := range in {
                result = fn(result, n)
            }
            select {
            case out <- result:
            case <-ctx.Done():
            }
        }()
        return out
    }
}

// Compose chains stages together
func Compose(stages ...Stage) Stage {
    return func(ctx context.Context, in <-chan int) <-chan int {
        out := in
        for _, stage := range stages {
            out = stage(ctx, out)
        }
        return out
    }
}

// Usage example
func ExamplePipeline() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Build pipeline: Generator -> Filter(even) -> Square -> Sum
    pipeline := Compose(
        Filter(func(n int) bool { return n%2 == 0 }),
        Transform(func(n int) int { return n * n }),
        Reduce(func(acc, n int) int { return acc + n }, 0),
    )
    
    source := Generator(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    result := pipeline(ctx, source)
    
    fmt.Println("Result:", <-result) // 220 (4 + 16 + 36 + 64 + 100)
}
```

**Tasks:**
1. Add error handling to pipeline stages
2. Implement parallel map stage (fan-out/fan-in)
3. Add metrics collection to each stage
4. Create a generic version using generics

**Deliverable:** Production-ready pipeline package

#### Exercise 2.2: Fan-Out / Fan-In (3 hours)

Implement parallel processing patterns:

```go
package concurrency

import (
    "context"
    "sync"
)

// FanOut distributes work to multiple workers
func FanOut(ctx context.Context, input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    
    for i := 0; i < workers; i++ {
        outputs[i] = processWorker(ctx, input, i)
    }
    
    return outputs
}

func processWorker(ctx context.Context, input <-chan int, id int) <-chan int {
    output := make(chan int)
    go func() {
        defer close(output)
        for val := range input {
            // Process value
            result := process(val, id)
            select {
            case output <- result:
            case <-ctx.Done():
                return
            }
        }
    }()
    return output
}

func process(val, workerID int) int {
    // Simulate CPU-intensive work
    return val * val + workerID
}

// FanIn merges multiple channels into one
func FanIn(ctx context.Context, channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    multiplexed := make(chan int)
    
    // Multiplex function
    multiplex := func(c <-chan int) {
        defer wg.Done()
        for val := range c {
            select {
            case multiplexed <- val:
            case <-ctx.Done():
                return
            }
        }
    }
    
    // Start goroutine for each input channel
    wg.Add(len(channels))
    for _, ch := range channels {
        go multiplex(ch)
    }
    
    // Close output when all input channels close
    go func() {
        wg.Wait()
        close(multiplexed)
    }()
    
    return multiplexed
}

// OrderedFanIn preserves input order
func OrderedFanIn(ctx context.Context, channels ...<-chan int) <-chan int {
    // Implementation that preserves order using sequence numbers
    // ... (more complex implementation)
    return nil
}

// Example usage
func ParallelProcessingExample() {
    ctx := context.Background()
    
    // Create input
    input := make(chan int, 100)
    go func() {
        defer close(input)
        for i := 1; i <= 100; i++ {
            input <- i
        }
    }()
    
    // Distribute to 4 workers
    workers := FanOut(ctx, input, 4)
    
    // Merge results
    results := FanIn(ctx, workers...)
    
    // Collect results
    var sum int
    for val := range results {
        sum += val
    }
}
```

**Deliverable:** Parallel processing utilities with benchmarks

---

### Day 3: Select and Synchronization

#### Exercise 3.1: Select Mastery (3 hours)

Master the select statement:

```go
package concurrency

import (
    "context"
    "fmt"
    "time"
)

// TimeoutPattern demonstrates timeout with select
func TimeoutPattern() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()
    
    select {
    case result := <-ch:
        fmt.Println("Received:", result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}

// NonBlockingSelect demonstrates non-blocking operations
func NonBlockingSelect() {
    ch := make(chan int)
    
    // Non-blocking send
    select {
    case ch <- 1:
        fmt.Println("Sent value")
    default:
        fmt.Println("Channel full, skipping")
    }
    
    // Non-blocking receive
    select {
    case val := <-ch:
        fmt.Println("Received:", val)
    default:
        fmt.Println("No value available")
    }
}

// PrioritySelect handles multiple channels with priority
func PrioritySelect(ctx context.Context) {
    highPriority := make(chan string)
    lowPriority := make(chan string)
    
    for {
        select {
        case msg := <-highPriority:
            fmt.Println("High priority:", msg)
            // Process immediately
            
        case msg := <-lowPriority:
            // Check high priority again
            select {
            case highMsg := <-highPriority:
                fmt.Println("High priority:", highMsg)
                // Return low priority message
                go func() { lowPriority <- msg }()
            default:
                fmt.Println("Low priority:", msg)
            }
            
        case <-ctx.Done():
            return
        }
    }
}

// RandomSelect with default for load balancing
func RandomSelect(inputs []<-chan int) <-chan int {
    output := make(chan int)
    
    go func() {
        defer close(output)
        cases := make([]reflect.SelectCase, len(inputs))
        for i, ch := range inputs {
            cases[i] = reflect.SelectCase{
                Dir:  reflect.SelectRecv,
                Chan: reflect.ValueOf(ch),
            }
        }
        
        for len(cases) > 0 {
            chosen, value, ok := reflect.Select(cases)
            if !ok {
                // Channel closed, remove it
                cases = append(cases[:chosen], cases[chosen+1:]...)
                continue
            }
            output <- value.Interface().(int)
        }
    }()
    
    return output
}
```

**Tasks:**
1. Implement rate limiting using select and time.Tick
2. Create a heartbeat mechanism
3. Build a connection pool with select
4. Implement circuit breaker with select

**Deliverable:** Select pattern library with 10+ patterns

#### Exercise 3.2: Sync Primitives (2 hours)

Master synchronization primitives:

```go
package concurrency

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// WaitGroupExample demonstrates proper WaitGroup usage
func WaitGroupExample() {
    var wg sync.WaitGroup
    urls := []string{"url1", "url2", "url3"}
    
    for _, url := range urls {
        wg.Add(1) // Add BEFORE starting goroutine
        go func(u string) {
            defer wg.Done()
            fetchURL(u)
        }(url)
    }
    
    wg.Wait() // Wait for all to complete
}

func fetchURL(url string) {
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Fetched:", url)
}

// OnceExample demonstrates one-time initialization
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{data: "initialized"}
        fmt.Println("Singleton initialized")
    })
    return instance
}

// RWMutexExample demonstrates read-write lock
type SafeCache struct {
    mu    sync.RWMutex
    data  map[string]string
}

func NewSafeCache() *SafeCache {
    return &SafeCache{data: make(map[string]string)}
}

func (c *SafeCache) Get(key string) (string, bool) {
    c.mu.RLock()         // Multiple readers allowed
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

func (c *SafeCache) Set(key, value string) {
    c.mu.Lock()          // Exclusive access
    defer c.mu.Unlock()
    c.data[key] = value
}

// PoolExample demonstrates sync.Pool
func PoolExample() {
    var bufferPool = sync.Pool{
        New: func() interface{} {
            return make([]byte, 1024)
        },
    }
    
    // Get from pool
    buf := bufferPool.Get().([]byte)
    
    // Use buffer
    copy(buf, "data")
    
    // Return to pool
    bufferPool.Put(buf)
}

// AtomicOperations demonstrates atomic primitives
func AtomicOperations() {
    var counter int64
    var flag int32
    
    // Atomic increment
    atomic.AddInt64(&counter, 1)
    
    // Compare and swap
    if atomic.CompareAndSwapInt32(&flag, 0, 1) {
        fmt.Println("Acquired lock")
    }
    
    // Atomic load/store
    val := atomic.LoadInt64(&counter)
    atomic.StoreInt64(&counter, val+1)
}

// CondExample demonstrates condition variables
func CondExample() {
    var mu sync.Mutex
    cond := sync.NewCond(&mu)
    ready := false
    
    // Waiter
    go func() {
        mu.Lock()
        for !ready {
            cond.Wait() // Releases lock and waits
        }
        fmt.Println("Proceeding!")
        mu.Unlock()
    }()
    
    time.Sleep(100 * time.Millisecond)
    
    // Signaler
    mu.Lock()
    ready = true
    cond.Signal() // Wake one waiter
    mu.Unlock()
}
```

**Deliverable:** Synchronization primitives reference guide

---

### Day 4: Concurrency Patterns

#### Exercise 4.1: Worker Pool (3 hours)

Build a robust worker pool:

```go
package workerpool

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "sync/atomic"
)

var (
    ErrPoolClosed = errors.New("worker pool is closed")
    ErrQueueFull  = errors.New("task queue is full")
)

// Task represents work to be done
type Task func() error

// Pool manages a pool of workers
type Pool struct {
    workers   int
    queueSize int
    
    tasks    chan Task
    wg       sync.WaitGroup
    ctx      context.Context
    cancel   context.CancelFunc
    
    active   int64
    pending  int64
    processed int64
    errors   int64
}

// New creates a new worker pool
func New(workers, queueSize int) *Pool {
    ctx, cancel := context.WithCancel(context.Background())
    
    p := &Pool{
        workers:   workers,
        queueSize: queueSize,
        tasks:     make(chan Task, queueSize),
        ctx:       ctx,
        cancel:    cancel,
    }
    
    p.start()
    return p
}

func (p *Pool) start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *Pool) worker(id int) {
    defer p.wg.Done()
    
    for {
        select {
        case task, ok := <-p.tasks:
            if !ok {
                return
            }
            
            atomic.AddInt64(&p.active, 1)
            atomic.AddInt64(&p.pending, -1)
            
            if err := task(); err != nil {
                atomic.AddInt64(&p.errors, 1)
            }
            
            atomic.AddInt64(&p.active, -1)
            atomic.AddInt64(&p.processed, 1)
            
        case <-p.ctx.Done():
            return
        }
    }
}

// Submit adds a task to the pool
func (p *Pool) Submit(task Task) error {
    select {
    case p.tasks <- task:
        atomic.AddInt64(&p.pending, 1)
        return nil
    case <-p.ctx.Done():
        return ErrPoolClosed
    default:
        return ErrQueueFull
    }
}

// SubmitBlocking adds a task, blocking if queue is full
func (p *Pool) SubmitBlocking(task Task) error {
    select {
    case p.tasks <- task:
        atomic.AddInt64(&p.pending, 1)
        return nil
    case <-p.ctx.Done():
        return ErrPoolClosed
    }
}

// Stats returns pool statistics
func (p *Pool) Stats() PoolStats {
    return PoolStats{
        Workers:   p.workers,
        Active:    int(atomic.LoadInt64(&p.active)),
        Pending:   int(atomic.LoadInt64(&p.pending)),
        Processed: atomic.LoadInt64(&p.processed),
        Errors:    atomic.LoadInt64(&p.errors),
    }
}

type PoolStats struct {
    Workers   int
    Active    int
    Pending   int
    Processed int64
    Errors    int64
}

// Stop gracefully shuts down the pool
func (p *Pool) Stop() {
    p.cancel()
    close(p.tasks)
    p.wg.Wait()
}

// StopWithContext waits for context or graceful shutdown
func (p *Pool) StopWithContext(ctx context.Context) error {
    done := make(chan struct{})
    go func() {
        p.Stop()
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

**Extensions:**
1. Dynamic worker scaling
2. Priority task queue
3. Task timeout handling
4. Result aggregation

**Deliverable:** Production-ready worker pool with full test suite

#### Exercise 4.2: Rate Limiter (2 hours)

Implement token bucket rate limiter:

```go
package ratelimit

import (
    "context"
    "sync"
    "time"
)

// TokenBucket implements token bucket algorithm
type TokenBucket struct {
    capacity   int
    tokens     int
    refillRate time.Duration
    mu         sync.Mutex
    lastRefill time.Time
}

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
    return &TokenBucket{
        capacity:   capacity,
        tokens:     capacity,
        refillRate: refillRate,
        lastRefill: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    
    tb.refill()
    
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    return false
}

func (tb *TokenBucket) Wait(ctx context.Context) error {
    for {
        tb.mu.Lock()
        tb.refill()
        
        if tb.tokens > 0 {
            tb.tokens--
            tb.mu.Unlock()
            return nil
        }
        
        // Calculate wait time
        waitTime := tb.refillRate / time.Duration(tb.capacity)
        tb.mu.Unlock()
        
        select {
        case <-time.After(waitTime):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func (tb *TokenBucket) refill() {
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill)
    tokensToAdd := int(elapsed / tb.refillRate)
    
    if tokensToAdd > 0 {
        tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
        tb.lastRefill = now
    }
}

// RateLimiter for multiple keys
type RateLimiter struct {
    buckets   map[string]*TokenBucket
    mu        sync.RWMutex
    capacity  int
    rate      time.Duration
}

func NewRateLimiter(capacity int, rate time.Duration) *RateLimiter {
    return &RateLimiter{
        buckets:  make(map[string]*TokenBucket),
        capacity: capacity,
        rate:     rate,
    }
}

func (rl *RateLimiter) Allow(key string) bool {
    rl.mu.RLock()
    bucket, exists := rl.buckets[key]
    rl.mu.RUnlock()
    
    if !exists {
        rl.mu.Lock()
        bucket = NewTokenBucket(rl.capacity, rl.rate)
        rl.buckets[key] = bucket
        rl.mu.Unlock()
    }
    
    return bucket.Allow()
}
```

**Deliverable:** Rate limiting package with sliding window variant

---

### Day 5: Debugging and Best Practices

#### Exercise 5.1: Race Detection (2 hours)

Learn to find and fix race conditions:

```go
package race

import (
    "sync"
    "testing"
)

// Detectable race condition
func DetectableRace() {
    var counter int
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // Race: unsynchronized access
        }()
    }
    
    wg.Wait()
}

// Fixed version
func FixedRace() {
    var counter int
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    
    wg.Wait()
}

// Using atomic
func AtomicFix() {
    var counter int64
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)
        }()
    }
    
    wg.Wait()
}

// Hidden race: slice access
func HiddenSliceRace() {
    results := make([]int, 100)
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            results[i] = i * i // Race: shared loop variable
        }()
    }
    
    wg.Wait()
}

// Fixed: pass loop variable
func FixedSliceRace() {
    results := make([]int, 100)
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            results[idx] = idx * idx
        }(i) // Pass i as argument
    }
    
    wg.Wait()
}

// Channel race
func ChannelRace() {
    ch := make(chan int)
    
    go func() {
        ch <- 1
    }()
    
    go func() {
        ch <- 2 // May panic if first send closes channel
    }()
    
    <-ch
}
```

**Test with race detector:**
```bash
go test -race ./...
```

**Deliverable:** Race detection guide with 10 examples

#### Exercise 5.2: Deadlock Detection (2 hours)

Learn to prevent and fix deadlocks:

```go
package deadlock

import (
    "sync"
    "time"
)

// Classic deadlock
func ClassicDeadlock() {
    var mu1, mu2 sync.Mutex
    
    go func() {
        mu1.Lock()
        time.Sleep(10 * time.Millisecond)
        mu2.Lock() // Waits forever
        // ...
        mu2.Unlock()
        mu1.Unlock()
    }()
    
    mu2.Lock()
    time.Sleep(10 * time.Millisecond)
    mu1.Lock() // Waits forever
    // ...
    mu1.Unlock()
    mu2.Unlock()
}

// Prevention: consistent lock ordering
func PreventedDeadlock() {
    var mu1, mu2 sync.Mutex
    
    go func() {
        mu1.Lock()
        mu2.Lock()
        // ...
        mu2.Unlock()
        mu1.Unlock()
    }()
    
    mu1.Lock()
    mu2.Lock()
    // ...
    mu2.Unlock()
    mu1.Unlock()
}

// Channel deadlock
func ChannelDeadlock() {
    ch := make(chan int)
    ch <- 1 // Deadlock: no receiver
    <-ch
}

// Fixed: buffered channel or separate goroutine
func FixedChannelDeadlock() {
    ch := make(chan int, 1) // Buffered
    ch <- 1
    <-ch
}

// WaitGroup deadlock
func WaitGroupDeadlock() {
    var wg sync.WaitGroup
    wg.Add(1)
    wg.Wait() // Deadlock: no Done() called
}
```

**Deliverable:** Deadlock prevention checklist and examples

---

## Code Review Checklist for Concurrent Code

### Correctness

- [ ] All goroutines have a termination condition
- [ ] Channels are closed by sender only
- [ ] Context cancellation is checked in long-running goroutines
- [ ] No data races (verified with -race flag)
- [ ] No potential deadlocks from lock ordering
- [ ] No goroutine leaks (all goroutines can exit)

### Channel Usage

- [ ] Buffered channels used when appropriate
- [ ] Channel direction specified in function signatures
- [ ] Select has default case when non-blocking needed
- [ ] Context cancellation in all select statements
- [ ] Channel closing logic is correct

### Synchronization

- [ ] Mutexes are not copied
- [ ] Lock ordering is consistent across codebase
- [ ] RWMutex used when read-heavy
- [ ] Atomic operations used for simple counters
- [ ] sync.Once used for one-time initialization

### Performance

- [ ] No unnecessary goroutine creation
- [ ] Worker pools used for bounded concurrency
- [ ] Channel buffering tuned for workload
- [ ] Avoids lock contention through sharding
- [ ] No busy waiting (use channels or sync.Cond)

### Testing

- [ ] Tests run with -race flag in CI
- [ ] Concurrent tests use stress testing
- [ ] Race conditions are reproducibly tested
- [ ] Timeout mechanisms tested
- [ ] Resource cleanup verified under load

---

## Assessment Criteria

### Knowledge Assessment (30%)

**Quiz Topics:**
1. CSP theory and Go's implementation
2. GMP scheduler internals
3. Channel semantics and happens-before
4. Sync primitive selection
5. Race condition patterns

**Passing Score:** 80%

### Coding Challenge (50%)

**Problem:** Build a concurrent download manager

**Requirements:**
- Download multiple files concurrently
- Rate limiting per domain
- Progress reporting
- Pause/resume capability
- Bandwidth throttling
- Race-free implementation

**Evaluation Criteria:**
- Correctness: 30%
- Concurrency patterns: 25%
- Performance: 20%
- Code quality: 15%
- Race safety: 10%

### System Design (20%)

**Problem:** Design a concurrent event processor

**Requirements:**
- Handle 10K events/second
- At-least-once delivery
- Ordered processing per partition
- Horizontal scalability

---

## Resources and References

### Official Resources
- [Go Memory Model](https://golang.org/ref/mem)
- [Concurrency in Go](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/)
- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide)

### Recommended Reading
- "Concurrency in Go" by Katherine Cox-Buday
- "The Go Programming Language" Chapters 8-9
- "100 Go Mistakes" - Concurrency chapters

### Tools
- Race detector: `go test -race`
- Deadlock detection: go-deadlock library
- Visualization: gotrace
- Profiling: go tool pprof

---

## Week 2 Completion Checklist

- [ ] All CSP theory understood
- [ ] Channel patterns mastered (10+ patterns)
- [ ] Worker pool implemented and tested
- [ ] Rate limiter working correctly
- [ ] Race detector used on all code
- [ ] No goroutine leaks in any exercise
- [ ] Week 2 assessment passed (80%+)

---

*Next: [Week 3: Cloud-Native Patterns](week3-cloudnative.md)*
