# Self-Assessment Quiz Questions

## Week 1: Go Fundamentals

### Multiple Choice Questions

**Q1. What is the zero value of a slice in Go?**

A. `nil`
B. Empty slice `[]`
C. Panic
D. Compiler error

<details>
<summary>Answer</summary>
A. `nil` - The zero value of a slice is nil, not an empty slice. A nil slice has no underlying array.
</details>

---

**Q2. Which of the following correctly describes Go's interface satisfaction?**

A. Explicit declaration required
B. Implicit through method matching
C. Inheritance-based
D. Tag-based annotation

<details>
<summary>Answer</summary>
B. Implicit through method matching - Go uses structural typing (duck typing), not nominal typing. A type implements an interface by implementing its methods.
</details>

---

**Q3. What happens when you send to a closed channel?**

A. Returns false
B. Blocks forever
C. Panics
D. Returns zero value

<details>
<summary>Answer</summary>
C. Panics - Sending on a closed channel causes a panic. Receiving from a closed channel returns the zero value immediately.
</details>

---

**Q4. Which is the correct way to handle errors in Go?**

A. `try-catch` blocks
B. Multiple return values with error as last
C. Global error variable
D. Exceptions

<details>
<summary>Answer</summary>
B. Multiple return values with error as last - Go uses explicit error returns, not exceptions. Functions return `(result, error)` with error as the last value.
</details>

---

**Q5. What does `defer` do in Go?**

A. Delays function execution to end of goroutine
B. Delays function execution to end of enclosing function
C. Runs function asynchronously
D. Creates a callback

<details>
<summary>Answer</summary>
B. Delays function execution to end of enclosing function - Deferred calls are executed in LIFO order when the function returns, regardless of how it returns (normal or panic).
</details>

---

**Q6. What is the difference between `make` and `new` in Go?**

A. No difference
B. `make` allocates, `new` initializes
C. `new` returns a pointer, `make` returns a value
D. `make` is for slices/maps/channels, `new` is for any type

<details>
<summary>Answer</summary>
D. `make` is for slices/maps/channels, `new` is for any type - `new(T)` allocates zeroed memory and returns `*T`. `make` initializes slices, maps, and channels with specified capacity.
</details>

---

**Q7. What is the zero value of a map in Go?**

A. Empty map
B. `nil`
C. Panic on creation
D. Zero-length map

<details>
<summary>Answer</summary>
B. `nil` - Like slices, maps are nil by default. You cannot write to a nil map without panicking; you must use `make` to initialize.
</details>

---

**Q8. Which statement about Go strings is correct?**

A. Strings are mutable
B. Strings are UTF-8 encoded
C. Strings are null-terminated
D. String indexing returns runes

<details>
<summary>Answer</summary>
B. Strings are UTF-8 encoded - Go strings are immutable byte slices that typically contain UTF-8 text. Indexing returns bytes, not runes.
</details>

---

**Q9. What happens when you append to a slice that shares an underlying array?**

A. Always creates a new array
B. May overwrite other slices if capacity exceeded
C. Panics
D. Returns error

<details>
<summary>Answer</summary>
B. May overwrite other slices if capacity exceeded - If append doesn't require reallocation (within capacity), it modifies the shared underlying array, affecting other slices.
</details>

---

**Q10. What is the purpose of the `init()` function?**

A. Initialize variables
B. Called automatically before main()
C. Constructor for structs
D. Initialize imports

<details>
<summary>Answer</summary>
B. Called automatically before main() - `init()` functions run automatically in package initialization order before `main()` starts. Multiple `init()` functions per package are allowed.
</details>

---

### Short Answer Questions

**Q11. Explain the difference between value receivers and pointer receivers.**

<details>
<summary>Sample Answer</summary>

Value receivers operate on a copy of the value, while pointer receivers operate on the original value. Key differences:

**Value Receivers:**

- Work on a copy (safe from modification)
- Safer for concurrent access
- Better for small, immutable types
- Cannot modify the original

**Pointer Receivers:**

- Work on the original value
- Can modify the receiver
- Required for large structs (avoid copying)
- Must use consistently for all methods of a type

Example:

```go
type Counter struct { count int }

// Value receiver - cannot modify original
func (c Counter) Increment() { c.count++ }

// Pointer receiver - can modify original
func (c *Counter) Increment() { c.count++ }
```

Best practice: Use pointer receivers when the method needs to modify the receiver or when the struct is large (>100 bytes).
</details>

---

**Q12. When should you use a pointer vs a value?**

<details>
<summary>Sample Answer</summary>

**Use Pointers When:**

1. The function needs to modify the argument
2. The struct is large (avoid copying overhead)
3. You need to distinguish between "no value" and "zero value"
4. For consistency when some methods use pointer receivers
5. When implementing interface methods that require pointers

**Use Values When:**

1. The data is small and simple (primitives, small structs)
2. You want to ensure the original isn't modified
3. For immutability guarantees
4. When working with map keys (must be comparable)

**Common Types and Conventions:**

- `time.Time`: Value (immutable)
- `sync.Mutex`: Value (but don't copy!)
- `bytes.Buffer`: Pointer (mutable)
- Large structs (>100 bytes): Pointer
- Small structs (<64 bytes): Either

</details>

---

**Q13. Explain Go's memory model and the happens-before relationship.**

<details>
<summary>Sample Answer</summary>

Go's memory model defines when memory writes are guaranteed to be visible to reads:

**Happens-Before Rules:**

1. **Within a goroutine:** Reads/writes are in program order
2. **Channel send/receive:** A send happens-before the corresponding receive completes
3. **Channel close:** Closing happens-before a receive that returns zero value
4. **Mutex:** Unlock happens-before subsequent Lock
5. **Once:** The function call happens-before any return from Once
6. **WaitGroup:** Wait happens-before Wait returns

**Example - Safe Publication:**

```go
var msg string
var done = make(chan bool)

func setup() {
    msg = "hello, world"  // Happens-before
    done <- true          // the receive
}

func main() {
    go setup()
    <-done                // Happens-before
    println(msg)          // this read (guaranteed to see "hello")
}
```

Without proper synchronization (like the channel), there's no guarantee the write to `msg` is visible.
</details>

---

### Code Analysis

**Q14. What is wrong with this code?**

```go
func process(items []string) {
    for _, item := range items {
        go func() {
            fmt.Println(item)
        }()
    }
}
```

<details>
<summary>Answer</summary>

The goroutine captures the loop variable `item` by reference, not by value. All goroutines will see the last value of `item`.

**Fix:** Pass `item` as an argument:

```go
func process(items []string) {
    for _, item := range items {
        go func(i string) {
            fmt.Println(i)
        }(item)  // Pass current value
    }
}
```

Or use a local variable:

```go
func process(items []string) {
    for _, item := range items {
        item := item  // Create new variable
        go func() {
            fmt.Println(item)
        }()
    }
}
```

</details>

---

**Q15. Identify the race condition:**

```go
var counter int
var wg sync.WaitGroup

for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        counter++
    }()
}
wg.Wait()
```

<details>
<summary>Answer</summary>

Multiple goroutines access `counter` without synchronization. The `++` operation is not atomic (read-modify-write), causing a data race.

**Fix with sync.Mutex:**

```go
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
```

**Fix with atomic:**

```go
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
```

</details>

---

**Q16. What does this code print?**

```go
func main() {
    s := make([]int, 0, 10)
    s = append(s, 1, 2, 3)

    b := s[:5]
    b[3] = 4
    b[4] = 5

    fmt.Println(len(s), cap(s), s)
    fmt.Println(len(b), cap(b), b)
}
```

<details>
<summary>Answer</summary>

```
3 10 [1 2 3]
5 10 [1 2 3 4 5]
```

Both slices share the same underlying array. `s` has length 3, but `b` extends to length 5 using the available capacity. Modifying `b` affects the underlying array that `s` also references.

If we then append to `s`:

```go
s = append(s, 6)
// s is now [1 2 3 6], b is [1 2 3 6 5]
// Index 3 was overwritten!
```

</details>

---

## Week 2: Concurrency

**Q17. What is the output of this program?**

```go
func main() {
    ch := make(chan int)
    go func() { ch <- 1 }()
    go func() { ch <- 2 }()
    fmt.Println(<-ch)
}
```

A. Always 1
B. Always 2
C. Either 1 or 2
D. Deadlock

<details>
<summary>Answer</summary>
C. Either 1 or 2 - The order of goroutine execution is non-deterministic. Whichever send completes first will be received. The other send will block until someone else receives from the channel.
</details>

---

**Q18. Explain the difference between buffered and unbuffered channels.**

<details>
<summary>Sample Answer</summary>

**Unbuffered Channels (make(chan T)):**

- Synchronous communication
- Send blocks until receive is ready
- Receive blocks until send is ready
- Guarantees synchronization point
- Use for signaling, coordination

**Buffered Channels (make(chan T, n)):**

- Asynchronous up to capacity n
- Send blocks only when buffer is full
- Receive blocks only when buffer is empty
- Decouples sender and receiver timing
- Use for work queues, batching

**Example:**

```go
// Unbuffered - synchronize
ch := make(chan bool)
go func() {
    doWork()
    ch <- true  // Blocks until main receives
}()
<-ch  // Waits for goroutine to finish

// Buffered - decouple
ch := make(chan int, 10)
for i := 0; i < 10; i++ {
    ch <- i  // Doesn't block (buffer has room)
}
```

</details>

---

**Q19. What does `select` do when multiple cases are ready?**

<details>
<summary>Sample Answer</summary>

When multiple cases in a `select` are ready, one is chosen pseudo-randomly. This prevents starvation and ensures fairness.

```go
select {
case v1 := <-ch1:
    // Might be chosen if both ready
case v2 := <-ch2:
    // Might be chosen if both ready
}
```

**To prioritize channels:**

```go
// Check high priority first
select {
case v := <-highPriority:
    return v
default:
}

// Then check low priority
select {
case v := <-highPriority:
    return v
case v := <-lowPriority:
    return v
}
```

Or use a for loop with default:

```go
for {
    select {
    case v := <-highPriority:
        return v
    case v := <-lowPriority:
        // Check high again
        select {
        case v2 := <-highPriority:
            // Return low priority item
            go func() { lowPriority <- v }()
            return v2
        default:
            return v
        }
    }
}
```

</details>

---

**Q20. Explain context cancellation propagation.**

<details>
<summary>Sample Answer</summary>

When a parent context is cancelled, all derived child contexts are automatically cancelled. This creates a tree of cancellation:

```go
// Root context
ctx := context.Background()

// Add timeout
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// Add values
ctx = context.WithValue(ctx, "requestID", "123")

// Child context inherits parent's cancellation
childCtx, childCancel := context.WithCancel(ctx)
defer childCancel()

// When parent times out, childCtx.Done() is also closed
```

**Best Practices:**

1. Always accept context as first parameter: `func(ctx context.Context, ...)`
2. Pass context through call chain, don't store in structs
3. Check `ctx.Done()` in long-running operations
4. Respect cancellation in loops and I/O
5. Add timeout for external calls

**Checking cancellation:**

```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()  // Return cancelled error
    case work := <-workCh:
        process(work)
    }
}
```

</details>

---

## Week 3: Cloud-Native

**Q21. What is the purpose of a circuit breaker?**

A. To encrypt traffic
B. To prevent cascading failures
C. To load balance requests
D. To cache responses

<details>
<summary>Answer</summary>
B. To prevent cascading failures - Circuit breakers stop requests to failing services, allowing them to recover and preventing the failure from spreading to the rest of the system.
</details>

---

**Q22. What are the three pillars of observability?**

<details>
<summary>Sample Answer</summary>

1. **Metrics**: Numeric measurements over time
   - Counters (request count, error count)
   - Gauges (current connections, memory usage)
   - Histograms (latency distributions)
   - Used for: alerting, dashboards, trend analysis

2. **Logs**: Timestamped event records
   - Structured (JSON) preferred over unstructured
   - Different levels (DEBUG, INFO, WARN, ERROR)
   - Correlation IDs for request tracing
   - Used for: debugging, audit trails, forensics

3. **Traces**: End-to-end request flow
   - Distributed across services
   - Spans represent operations
   - Parent-child relationships
   - Used for: performance analysis, bottleneck identification

Together they provide comprehensive system visibility. Metrics tell you THAT something is wrong, logs tell you WHAT happened, traces tell you WHERE the problem is.
</details>

---

**Q23. When should you use gRPC vs REST?**

<details>
<summary>Sample Answer</summary>

**Use gRPC When:**

- Internal service-to-service communication
- High performance requirements (binary Protocol Buffers)
- Need strong typing and code generation
- Streaming (client, server, or bidirectional)
- Both services use Go or have good gRPC support
- You control both client and server

**Use REST When:**

- External/public APIs
- Browser-based clients
- Need simple debugging with curl/browser
- Caching with HTTP proxies/CDN
- Wide client compatibility needed
- Human-readable responses preferred

**Hybrid Approach:**

- REST externally, gRPC internally
- GraphQL for complex queries
- WebSocket for real-time features

</details>

---

**Q24. What is graceful shutdown?**

<details>
<summary>Sample Answer</summary>

Graceful shutdown ensures in-flight requests complete before the application exits:

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: handler}

    // Start server in goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}
```

**Key Steps:**

1. Stop accepting new connections
2. Wait for active requests to complete
3. Close database connections
4. Flush logs and metrics
5. Exit cleanly

</details>

---

## Week 4: System Design

**Q25. Explain the CAP theorem.**

<details>
<summary>Sample Answer</summary>

CAP theorem states that in a distributed system with data replication, you can only guarantee two of three properties:

**Consistency (C):** All nodes see the same data at the same time. Every read receives the most recent write.

**Availability (A):** Every request receives a response (success or failure), without guarantee it contains the most recent write.

**Partition Tolerance (P):** The system continues to operate despite arbitrary message loss or failure of part of the system (network partition).

**Key Insight:** In distributed systems, network partitions are inevitable, so the real choice is between CP (Consistency + Partition tolerance) and AP (Availability + Partition tolerance).

**Examples:**

- **CP:** PostgreSQL, HBase, MongoDB (configured)
- **AP:** Cassandra, DynamoDB, Couchbase
- **CA:** Single-node databases (not truly distributed)

**PACELC Theorem Extension:** If there is a Partition, choose between Availability and Consistency; Else, choose between Latency and Consistency.
</details>

---

**Q26. What is the difference between strong consistency and eventual consistency?**

<details>
<summary>Sample Answer</summary>

**Strong Consistency:**

- After a write, all subsequent reads will see the new value
- Reads return the most recent write or an error
- Simpler application logic
- Higher latency (synchronous replication)
- Lower availability during partitions

**Eventual Consistency:**

- After a write, reads may not see the new value immediately
- All replicas will eventually converge to the same value
- Applications must handle stale data
- Lower latency (asynchronous replication)
- Higher availability during partitions

**When to Use:**

- **Strong:** Financial transactions, inventory management, user authentication
- **Eventual:** Social media likes, comments, analytics, recommendations

**Other Consistency Models:**

- **Causal Consistency:** Preserves happens-before relationships
- **Read-Your-Writes:** Your writes are always visible to you
- **Monotonic Reads:** Successive reads see non-decreasing values

</details>

---

**Q27. Design a rate limiter for an API gateway.**

<details>
<summary>Sample Answer</summary>

**Token Bucket Algorithm:**

- Bucket has fixed capacity C
- Tokens added at rate R per second
- Each request consumes 1 token
- If no tokens available, request is rejected
- Allows burst up to bucket capacity

**Implementation:**

```go
type TokenBucket struct {
    capacity int
    tokens   float64
    rate     float64
    lastTime time.Time
    mu       sync.Mutex
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastTime).Seconds()
    tb.tokens = min(tb.capacity, tb.tokens + elapsed*tb.rate)
    tb.lastTime = now

    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }
    return false
}
```

**Distributed Rate Limiting:**

- Use Redis with Lua scripts for atomic operations
- Or use centralized rate limiter service
- Handle clock skew between nodes

**Rate Limit Headers:**

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1640995200
Retry-After: 60
```

**Alternative: Sliding Window**

- More accurate than token bucket
- Tracks actual timestamps
- Higher memory usage

</details>

---

## Advanced Topics

**Q28. How does Raft achieve consensus?**

<details>
<summary>Sample Answer</summary>

Raft uses a strong leader approach with three sub-problems:

**1. Leader Election:**

- Nodes start as Followers
- If no heartbeat received within timeout, becomes Candidate
- Candidate increments term, votes for itself, requests votes from peers
- Wins if majority votes received
- Split vote handled with randomized timeouts

**2. Log Replication:**

- Leader accepts client writes
- Appends to its log, then replicates to followers
- Entry committed when majority have replicated
- Followers apply committed entries to state machine

**3. Safety:**

- Only one leader per term
- Leader has all committed entries
- Log matching property ensures consistency
- Committed entries never overwritten

**State Machine:**

```
Follower -> [timeout] -> Candidate -> [wins election] -> Leader
    ^                         | [loses election]          |
    |                         v                           |
    +-------------------- Follower <----------------------+
```

**Compared to Paxos:**

- Raft is easier to understand and implement
- Separates leader election, log replication, safety
- Stronger leader model simplifies decisions

</details>

---

**Q29. What are CRDTs and when are they useful?**

<details>
<summary>Sample Answer</summary>

**Conflict-free Replicated Data Types** are data structures that can be replicated across nodes and merged without coordination.

**Required Properties:**

1. **Commutative:** A op B = B op A (order doesn't matter)
2. **Associative:** (A op B) op C = A op (B op C) (grouping doesn't matter)
3. **Idempotent:** A op A = A (duplicate operations have no effect)

**Common CRDT Types:**

**G-Counter (Grow-only Counter):**

```go
type GCounter map[string]int

func (c GCounter) Increment(node string) {
    c[node]++
}

func (c GCounter) Merge(other GCounter) {
    for node, val := range other {
        if val > c[node] {
            c[node] = val
        }
    }
}

func (c GCounter) Value() int {
    sum := 0
    for _, v := range c {
        sum += v
    }
    return sum
}
```

**PN-Counter:** Supports increment and decrement (two G-Counters)

**LWW-Register:** Last-write-wins based on timestamps

**OR-Set:** Add-wins or remove-wins set

**Use Cases:**

- Collaborative editing (Google Docs)
- Shopping carts (can add/remove while offline)
- Distributed counters
- Conflict resolution in disconnected systems
- Multi-region databases

</details>

---

**Q30. Explain event sourcing and its benefits.**

<details>
<summary>Sample Answer</summary>

**Event Sourcing** stores state changes as a sequence of events rather than just current state.

**Core Pattern:**

```
Command -> Validation -> Event -> Event Store -> Projection -> Read Model
```

**Benefits:**

1. **Complete Audit Trail:**
   - Every change recorded with timestamp and reason
   - Compliance and forensic capabilities

2. **Temporal Queries:**
   - Reconstruct state at any point in time
   - "What did the system look like yesterday?"
   - "How did we get to this state?"

3. **Debugging and Replay:**
   - Replay events to reproduce bugs
   - Test new logic against historical events
   - Copy production events to staging

4. **Natural Event-Driven:**
   - Fits well with microservices
   - Easy to add new read models
   - Cross-service communication via events

5. **Supports CQRS:**
   - Separate read and write models
   - Optimize each for their purpose

**Challenges:**

- Event schema evolution (versioning)
- Event store size (use snapshots)
- Learning curve
- Eventual consistency of projections

**Snapshot Strategy:**

- Take periodic snapshots of aggregate state
- Replay only events after snapshot
- Reduces replay time for large event histories

</details>

---

## Practice Coding Challenges

### Challenge 1: Concurrent Cache

Implement a thread-safe LRU cache with TTL support.

```go
type LRUCache struct {
    // Implementation required
}

func NewLRUCache(capacity int, ttl time.Duration) *LRUCache
func (c *LRUCache) Get(key string) (value interface{}, ok bool)
func (c *LRUCache) Set(key string, value interface{})
```

**Requirements:**

- O(1) Get and Set operations
- Automatic eviction of expired items
- Thread-safe for concurrent access
- Memory efficient

---

### Challenge 2: Rate Limiter

Implement a token bucket rate limiter.

```go
type RateLimiter struct {
    // Implementation required
}

func NewRateLimiter(rate int, burst int) *RateLimiter
func (rl *RateLimiter) Allow() bool
func (rl *RateLimiter) Wait(ctx context.Context) error
```

**Requirements:**

- Support burst traffic
- Thread-safe
- Context-aware waiting
- Efficient memory usage

---

### Challenge 3: Worker Pool

Implement a dynamic worker pool that scales based on queue depth.

```go
type WorkerPool struct {
    // Implementation required
}

func NewWorkerPool(minWorkers, maxWorkers int) *WorkerPool
func (wp *WorkerPool) Submit(task func()) error
func (wp *WorkerPool) Stop()
```

**Requirements:**

- Scale up when queue grows
- Scale down when idle
- Graceful shutdown
- Handle panics in tasks

---

### Challenge 4: Circuit Breaker

Implement a circuit breaker with half-open state.

```go
type CircuitBreaker struct {
    // Implementation required
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker
func (cb *CircuitBreaker) Execute(fn func() error) error
func (cb *CircuitBreaker) State() State
```

**Requirements:**

- Three states: Closed, Open, HalfOpen
- Automatic recovery testing in HalfOpen
- Configurable failure threshold and timeout
- Thread-safe

---

### Challenge 5: Retry with Backoff

Implement exponential backoff with jitter.

```go
type RetryConfig struct {
    MaxRetries   int
    BaseDelay    time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}

func Retry(ctx context.Context, config RetryConfig, fn func() error) error
```

**Requirements:**

- Exponential backoff
- Random jitter to avoid thundering herd
- Respect context cancellation
- Return last error on exhaustion

---

## Scoring Guide

| Score | Level | Recommendation |
|-------|-------|----------------|
| 0-30% | Beginner | Review Week 1 materials thoroughly |
| 30-50% | Novice | Review specific weak areas in Week 1-2 |
| 50-70% | Intermediate | Ready for Week 3 content, review gaps |
| 70-85% | Advanced | Ready for Week 4+ content |
| 85-100% | Expert | Consider advanced track, mentor others |

---

## Answer Key Summary

| Q | Answer | Topic |
|---|--------|-------|
| 1 | A | Slice zero value |
| 2 | B | Interface satisfaction |
| 3 | C | Closed channel send |
| 4 | B | Error handling |
| 5 | B | Defer behavior |
| 6 | D | make vs new |
| 7 | B | Map zero value |
| 8 | B | String encoding |
| 9 | B | Slice sharing |
| 10 | B | init() function |
| 11 | Essay | Value vs pointer receivers |
| 12 | Essay | Pointer usage guidelines |
| 13 | Essay | Memory model |
| 14 | Essay | Loop variable capture |
| 15 | Essay | Race condition |
| 16 | Essay | Slice capacity |
| 17 | C | Channel non-determinism |
| 18 | Essay | Buffered vs unbuffered |
| 19 | Essay | Select behavior |
| 20 | Essay | Context propagation |
| 21 | B | Circuit breaker purpose |
| 22 | Essay | Observability pillars |
| 23 | Essay | gRPC vs REST |
| 24 | Essay | Graceful shutdown |
| 25 | Essay | CAP theorem |
| 26 | Essay | Consistency models |
| 27 | Essay | Rate limiter design |
| 28 | Essay | Raft consensus |
| 29 | Essay | CRDTs |
| 30 | Essay | Event sourcing |

---

*Use these questions for self-assessment before each week's assessment.*
