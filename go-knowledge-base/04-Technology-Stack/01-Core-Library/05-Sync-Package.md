# TS-CL-005: Go sync Package - Concurrency Primitives Deep Dive

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #concurrency #mutex #waitgroup #once #pool #syncmap
> **权威来源**:
> - [Go sync Package](https://golang.org/pkg/sync/) - Go standard library
> - [Go Memory Model](https://golang.org/ref/mem) - Memory model
> - [The Go Scheduler](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html) - Ardan Labs

---

## 1. sync.Mutex - Mutual Exclusion Lock

### 1.1 Implementation Details

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       sync.Mutex State Machine                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  States:                                                                     │
│  ┌──────────┐    Lock()     ┌──────────┐    Lock()     ┌──────────┐        │
│  │ Unlocked │──────────────►│  Locked  │──────────────►│  Locked  │        │
│  │  (0)     │               │  (1)     │  (blocked)    │  (N)     │        │
│  └──────────┘               └────┬─────┘               └────┬─────┘        │
│       ▲                          │                          │              │
│       │ Unlock()                 │ Unlock()                 │ Unlock()     │
│       └──────────────────────────┴──────────────────────────┘              │
│                                                                              │
│  Internal Structure:                                                         │
│  type Mutex struct {                                                         │
│      state int32    // 0=unlocked, 1=locked, N=locked with waiters         │
│      sema  uint32   // Semaphore for parking goroutines                     │
│  }                                                                           │
│                                                                              │
│  Fast Path: atomic CAS on state (uncontended case)                          │
│  Slow Path: semaphore-based blocking (contended case)                       │
│                                                                              │
│  Lock Contention:                                                            │
│  - Uncontended: ~10ns (single atomic operation)                             │
│  - Contended: ~100ns-1μs (semaphore operations)                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Usage Patterns

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Basic mutex usage
func basicMutex() {
    var mu sync.Mutex
    counter := 0
    
    for i := 0; i < 1000; i++ {
        go func() {
            mu.Lock()
            defer mu.Unlock()
            counter++
        }()
    }
    
    time.Sleep(time.Second)
    fmt.Println(counter) // 1000
}

// Struct with embedded mutex
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *Counter) Get() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

func (c *Counter) Add(n int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count += n
}

// NEVER copy a mutex
type BadExample struct {
    mu sync.Mutex
    data int
}

func badUsage() {
    b1 := BadExample{data: 1}
    b2 := b1 // COPIES the mutex - WRONG!
    // b1.mu and b2.mu are now the same mutex
    // This causes undefined behavior
}

// Correct way: pass by pointer
func goodUsage() {
    b1 := &BadExample{data: 1}
    b2 := b1 // Both point to same instance - OK
}
```

---

## 2. sync.RWMutex - Read-Write Lock

### 2.1 Implementation Details

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       sync.RWMutex Implementation                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  State: Multiple readers OR single writer                                    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  type RWMutex struct {                                               │   │
│  │      w           Mutex        // Held by writers                     │   │
│  │      writerSem   uint32       // Semaphore for writers               │   │
│  │      readerSem   uint32       // Semaphore for readers               │   │
│  │      readerCount atomic.Int32 // Number of pending readers           │   │
│  │      readerWait  atomic.Int32 // Number of departing readers         │   │
│  │  }                                                                   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Read Lock (RLock):                                                          │
│  1. Increment readerCount                                                    │
│  2. If writer active (readerCount < 0), block on readerSem                  │
│  3. Return immediately otherwise                                             │
│                                                                              │
│  Read Unlock (RUnlock):                                                      │
│  1. Decrement readerCount                                                    │
│  2. If was last reader and writer waiting, signal writer                    │
│                                                                              │
│  Write Lock (Lock):                                                          │
│  1. Acquire w mutex                                                          │
│  2. Announce intent: readerCount -= maxReaders (large negative)             │
│  3. Wait for existing readers to finish                                      │
│                                                                              │
│  Write Unlock (Unlock):                                                      │
│  1. Remove intent: readerCount += maxReaders                                 │
│  2. Signal waiting readers                                                   │
│  3. Release w mutex                                                          │
│                                                                              │
│  Performance:                                                                │
│  - Uncontended read: ~10ns                                                   │
│  - Uncontended write: ~15ns                                                  │
│  - Contended: depends on number of goroutines                               │
│                                                                              │
│  Warning: Write starvation possible with constant read load                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Usage Patterns

```go
// RWMutex for read-heavy workloads
type Cache struct {
    mu    sync.RWMutex
    items map[string]interface{}
    hits  int64
    misses int64
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    val, ok := c.items[key]
    if ok {
        atomic.AddInt64(&c.hits, 1)
    } else {
        atomic.AddInt64(&c.misses, 1)
    }
    return val, ok
}

func (c *Cache) Set(key string, val interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = val
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

func (c *Cache) GetAll() map[string]interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // Return copy to avoid data race
    result := make(map[string]interface{}, len(c.items))
    for k, v := range c.items {
        result[k] = v
    }
    return result
}

// When NOT to use RWMutex:
// 1. Write-heavy workloads (use Mutex instead)
// 2. Short critical sections (Mutex may be faster)
// 3. When write starvation is a concern
```

---

## 3. sync.WaitGroup - Goroutine Synchronization

### 3.1 Implementation

```go
// WaitGroup waits for a collection of goroutines to finish
type WaitGroup struct {
    noCopy noCopy  // Prevents copying
    
    // state: high 32 bits = counter, low 32 bits = waiter count
    state atomic.Uint64
    sema  uint32   // Semaphore for blocking waiters
}

// Add adds delta to the counter
go func(wg *sync.WaitGroup) {
    wg.Add(1)
    defer wg.Done()
    // Do work
}(wg)

// Wait blocks until counter is zero
wg.Wait()
```

### 3.2 Usage Patterns

```go
// Parallel processing with WaitGroup
func processItems(items []Item) {
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()
            process(i)
        }(item)
    }
    
    wg.Wait()
    fmt.Println("All items processed")
}

// Bounded concurrency with semaphore pattern
func processWithLimit(items []Item, maxWorkers int) {
    var wg sync.WaitGroup
    sem := make(chan struct{}, maxWorkers)
    
    for _, item := range items {
        wg.Add(1)
        sem <- struct{}{} // Acquire semaphore
        
        go func(i Item) {
            defer wg.Done()
            defer func() { <-sem }() // Release semaphore
            process(i)
        }(item)
    }
    
    wg.Wait()
}

// Common mistake: Add in goroutine (WRONG)
func badPattern() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        go func() {
            wg.Add(1) // WRONG - race condition!
            defer wg.Done()
            work()
        }()
    }
    wg.Wait()
}

// Correct pattern: Add before starting goroutine
func goodPattern() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1) // CORRECT
        go func() {
            defer wg.Done()
            work()
        }()
    }
    wg.Wait()
}

// Best practice: Pass wg as parameter
func bestPractice(items []Item, wg *sync.WaitGroup) {
    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()
            process(i)
        }(item)
    }
}
```

---

## 4. sync.Once - Guaranteed Single Execution

### 4.1 Implementation

```go
type Once struct {
    done atomic.Uint32
    m    Mutex
}

func (o *Once) Do(f func()) {
    // Fast path
    if o.done.Load() == 0 {
        o.doSlow(f)
    }
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()
    if o.done.Load() == 0 {
        defer o.done.Store(1)
        f()
    }
}
```

### 4.2 Usage Patterns

```go
// Singleton pattern
var (
    instance *Service
    once     sync.Once
)

func GetService() *Service {
    once.Do(func() {
        instance = &Service{
            // Expensive initialization
            connection: createConnection(),
            cache:      make(map[string]interface{}),
        }
    })
    return instance
}

// Go 1.21+: OnceValue, OnceValues for return values
var config = sync.OnceValue(func() Config {
    return loadConfigFromFile("config.yaml")
})

// Usage
func main() {
    cfg := config() // Computed once, cached thereafter
    fmt.Println(cfg.Database.Host)
}

// Once with error handling
var (
    db     *sql.DB
    dbOnce sync.Once
    dbErr  error
)

func GetDB() (*sql.DB, error) {
    dbOnce.Do(func() {
        db, dbErr = sql.Open("postgres", dsn)
        if dbErr != nil {
            return
        }
        dbErr = db.Ping()
    })
    return db, dbErr
}
```

---

## 5. sync.Pool - Object Reuse

### 5.1 Implementation

```go
type Pool struct {
    noCopy noCopy
    
    local     unsafe.Pointer // [P]poolLocal
    localSize uintptr
    
    victim     unsafe.Pointer // [P]poolLocal from previous GC
    victimSize uintptr
    
    New func() interface{}
}

// Per-P poolLocal to reduce contention
type poolLocal struct {
    poolLocalInternal
    pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}

type poolLocalInternal struct {
    private interface{}      // Only used by owner P
    shared  poolChain        // Local P can pushHead/popHead; others popTail
}
```

### 5.2 Usage Patterns

```go
// Buffer pool for reducing GC pressure
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func processData(data []byte) []byte {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf[:cap(buf)]) // Reset slice before returning
    
    // Use buf...
    n := copy(buf, data)
    return buf[:n]
}

// Important considerations:
// 1. Objects may be GC'd at any time
// 2. No guarantee of getting same object back
// 3. Clear sensitive data before Put
// 4. Don't use for connection pooling (use dedicated pool)

// Reset objects properly
type Buffer struct {
    buf []byte
}

func (b *Buffer) Reset() {
    b.buf = b.buf[:0] // Keep capacity, reset length
}

var bufPool = sync.Pool{
    New: func() interface{} {
        return &Buffer{buf: make([]byte, 0, 1024)}
    },
}

func useBuffer() {
    b := bufPool.Get().(*Buffer)
    defer func() {
        b.Reset()
        bufPool.Put(b)
    }()
    
    // Use buffer...
}

// When NOT to use sync.Pool:
// 1. Connection pooling (use dedicated connection pool)
// 2. When you need guaranteed object reuse
// 3. For small objects (allocation may be faster)
// 4. When objects have complex cleanup needs
```

---

## 6. sync.Map - Concurrent Safe Map

### 6.1 Implementation Overview

```go
type Map struct {
    mu Mutex
    read atomic.Pointer[readOnly]
    dirty map[interface{}]*entry
    misses int
}

type readOnly struct {
    m       map[interface{}]*entry
    amended bool // true if dirty contains entries not in read
}

type entry struct {
    p atomic.Pointer[interface{}] // nil=expunged, interface{}=value
}
```

### 6.2 When to Use

```go
// sync.Map is optimized for:
// 1. Entry written once but read many times
// 2. Multiple goroutines read/write disjoint sets of keys

// Good use case: Cache with infrequent writes
var cache sync.Map

func getFromCache(key string) (value interface{}, ok bool) {
    return cache.Load(key)
}

func setCache(key string, value interface{}) {
    cache.Store(key, value)
}

// Compute if not exists (atomic)
func getOrCompute(key string, compute func() interface{}) interface{} {
    val, ok := cache.Load(key)
    if ok {
        return val
    }
    
    newVal := compute()
    actual, loaded := cache.LoadOrStore(key, newVal)
    if loaded {
        return actual // Another goroutine stored first
    }
    return newVal
}

// Range over all entries
cache.Range(func(key, value interface{}) bool {
    fmt.Printf("%s: %v\n", key, value)
    return true // continue iteration
})

// Delete
cache.Delete(key)

// For most cases, use regular map with mutex
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]interface{}
}

func (sm *SafeMap) Get(key string) (interface{}, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    val, ok := sm.m[key]
    return val, ok
}

func (sm *SafeMap) Set(key string, val interface{}) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.m[key] = val
}

// sync.Map vs map+Mutex:
// - sync.Map: Better for many reads, few writes, many keys
// - map+Mutex: Better for write-heavy, type safety
```

---

## 7. Best Practices

```
sync Package Best Practices:

Mutex:
□ Lock for minimum time possible
□ Use defer Unlock for safety
□ Prefer RWMutex for read-heavy workloads
□ Never copy mutexes
□ Consider atomic operations for simple counters

WaitGroup:
□ Call Add before starting goroutine
□ Call Done with defer
□ Don't reuse without waiting
□ Check for negative counter

Pool:
□ Reset objects before Put
□ Don't rely on objects persisting
□ Size pool appropriately
□ Don't pool connections (use dedicated pool)

General:
□ Prefer channels for coordination
□ Use sync primitives for shared state
□ Profile to identify bottlenecks
□ Avoid premature optimization
```

---

## 8. Checklist

```
Concurrency Checklist:
□ No data races (use -race detector)
□ Proper mutex usage
□ RWMutex for read-heavy patterns
□ WaitGroup Add before goroutine
□ Once for lazy initialization
□ Pool for temporary objects
□ Atomic operations for simple cases
□ No busy waiting
□ Proper cleanup on shutdown
```
