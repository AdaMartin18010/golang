# AD-008: Performance Optimization Patterns

> **Dimension**: Application Domains
> **Level**: S (16+ KB)
> **Tags**: #performance #optimization #caching #profiling #benchmarking

---

## 1. Performance Fundamentals

### 1.1 Performance Metrics

| Metric | Definition | Target |
|--------|-----------|--------|
| Latency | Time to process a request | P99 < 100ms |
| Throughput | Requests per second | > 10K RPS |
| Error Rate | Failed requests percentage | < 0.1% |
| Resource Usage | CPU/Memory utilization | < 80% |

### 1.2 Latency Numbers

| Operation | Time |
|-----------|------|
| L1 cache reference | 0.5 ns |
| Branch mispredict | 5 ns |
| L2 cache reference | 7 ns |
| Mutex lock/unlock | 100 ns |
| Main memory reference | 100 ns |
| Compress 1K bytes | 10,000 ns |
| Send 2K bytes over network | 200,000 ns |
| SSD random read | 150,000 ns |
| Read 1 MB from SSD | 1,000,000 ns |
| Round trip within datacenter | 500,000 ns |
| Read 1 MB from memory | 10,000,000 ns |
| Round trip between datacenters | 50,000,000 ns |

---

## 2. Profiling and Benchmarking

### 2.1 CPU Profiling

```go
package profiling

import (
    "os"
    "runtime/pprof"
)

func StartCPUProfile(filename string) (*os.File, error) {
    f, err := os.Create(filename)
    if err != nil {
        return nil, err
    }

    if err := pprof.StartCPUProfile(f); err != nil {
        f.Close()
        return nil, err
    }

    return f, nil
}

func StopCPUProfile(f *os.File) {
    pprof.StopCPUProfile()
    f.Close()
}
```

### 2.2 Memory Profiling

```go
package profiling

import (
    "os"
    "runtime"
    "runtime/pprof"
)

func WriteHeapProfile(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    runtime.GC() // Run GC for accurate stats
    return pprof.WriteHeapProfile(f)
}
```

### 2.3 Benchmarking

```go
package main

import (
    "testing"
    "time"
)

// Basic benchmark
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FunctionToBenchmark()
    }
}

// Benchmark with memory allocation
func BenchmarkWithAlloc(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        FunctionWithAlloc()
    }
}

// Benchmark with setup
func BenchmarkWithSetup(b *testing.B) {
    data := setupData()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        ProcessData(data)
    }
}

// Parallel benchmark
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            FunctionToBenchmark()
        }
    })
}
```

---

## 3. Caching Strategies

### 3.1 In-Memory Cache

```go
package cache

import (
    "context"
    "sync"
    "time"
)

type Cache struct {
    data map[string]*item
    mu   sync.RWMutex
    ttl  time.Duration
}

type item struct {
    value      interface{}
    expiration time.Time
}

func NewCache(ttl time.Duration) *Cache {
    c := &Cache{
        data: make(map[string]*item),
        ttl:  ttl,
    }
    go c.cleanup()
    return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, exists := c.data[key]
    if !exists || time.Now().After(item.expiration) {
        return nil, false
    }

    return item.value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = &item{
        value:      value,
        expiration: time.Now().Add(c.ttl),
    }
}

func (c *Cache) cleanup() {
    ticker := time.NewTicker(c.ttl)
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.data {
            if now.After(item.expiration) {
                delete(c.data, key)
            }
        }
        c.mu.Unlock()
    }
}
```

### 3.2 Redis Cache

```go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })

    return &RedisCache{
        client: client,
        ttl:    ttl,
    }
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }

    return json.Unmarshal(data, dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}
```

### 3.3 Cache-Aside Pattern

```go
package cache

import (
    "context"
    "fmt"
)

type CacheAside struct {
    cache   Cache
    loader  func(ctx context.Context, key string) (interface{}, error)
}

func (c *CacheAside) Get(ctx context.Context, key string) (interface{}, error) {
    // Try cache first
    if value, found := c.cache.Get(key); found {
        return value, nil
    }

    // Cache miss - load from source
    value, err := c.loader(ctx, key)
    if err != nil {
        return nil, fmt.Errorf("load failed: %w", err)
    }

    // Store in cache
    c.cache.Set(key, value)

    return value, nil
}
```

---

## 4. Database Optimization

### 4.1 Connection Pooling

```go
package db

import (
    "database/sql"
    "time"
    _ "github.com/lib/pq"
)

func NewDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(1 * time.Minute)

    return db, nil
}
```

### 4.2 Query Optimization

```go
package db

import (
    "context"
    "database/sql"
)

// Prepared statement
type UserRepository struct {
    db    *sql.DB
    getByIDStmt *sql.Stmt
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
    getByIDStmt, err := db.Prepare("SELECT id, name, email FROM users WHERE id = $1")
    if err != nil {
        return nil, err
    }

    return &UserRepository{
        db: db,
        getByIDStmt: getByIDStmt,
    }, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    var user User
    err := r.getByIDStmt.QueryRowContext(ctx, id).Scan(
        &user.ID, &user.Name, &user.Email,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// Batch insert
func (r *UserRepository) BatchInsert(ctx context.Context, users []User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx,
        "INSERT INTO users (name, email) VALUES ($1, $2)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, user := range users {
        if _, err := stmt.ExecContext(ctx, user.Name, user.Email); err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

---

## 5. Concurrency Optimization

### 5.1 Worker Pool

```go
package worker

import (
    "context"
    "sync"
)

type Pool struct {
    workers int
    jobs    chan Job
    wg      sync.WaitGroup
}

type Job func() error

func NewPool(workers int) *Pool {
    return &Pool{
        workers: workers,
        jobs:    make(chan Job, workers * 2),
    }
}

func (p *Pool) Start(ctx context.Context) {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(ctx)
    }
}

func (p *Pool) worker(ctx context.Context) {
    defer p.wg.Done()

    for {
        select {
        case <-ctx.Done():
            return
        case job := <-p.jobs:
            job()
        }
    }
}

func (p *Pool) Submit(job Job) {
    p.jobs <- job
}

func (p *Pool) Stop() {
    close(p.jobs)
    p.wg.Wait()
}
```

### 5.2 Fan-Out Fan-In

```go
package concurrent

import (
    "context"
    "sync"
)

func FanOut(ctx context.Context, inputs []Input, workers int) <-chan Result {
    inputCh := make(chan Input)
    resultCh := make(chan Result)

    // Producer
    go func() {
        defer close(inputCh)
        for _, input := range inputs {
            select {
            case <-ctx.Done():
                return
            case inputCh <- input:
            }
        }
    }()

    // Workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for input := range inputCh {
                result := Process(input)
                select {
                case <-ctx.Done():
                    return
                case resultCh <- result:
                }
            }
        }()
    }

    // Closer
    go func() {
        wg.Wait()
        close(resultCh)
    }()

    return resultCh
}
```

---

## 6. Memory Optimization

### 6.1 Object Pool

```go
package pool

import (
    "sync"
)

type BufferPool struct {
    pool sync.Pool
    size int
}

func NewBufferPool(size int) *BufferPool {
    return &BufferPool{
        size: size,
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, size)
            },
        },
    }
}

func (p *BufferPool) Get() []byte {
    return p.pool.Get().([]byte)
}

func (p *BufferPool) Put(buf []byte) {
    if cap(buf) >= p.size {
        p.pool.Put(buf[:p.size])
    }
}
```

### 6.2 Memory-Efficient Structures

```go
package structs

import "sync/atomic"

// Cache-friendly struct layout
type UserStats struct {
    // Frequently accessed together
    ID        uint64 // 8 bytes
    LoginTime int64  // 8 bytes

    // Frequently accessed together
    Posts    atomic.Uint64 // 8 bytes
    Comments atomic.Uint64 // 8 bytes
    Likes    atomic.Uint64 // 8 bytes
}
```

---

## 7. Network Optimization

### 7.1 HTTP Client Optimization

```go
package http

import (
    "net"
    "net/http"
    "time"
)

func OptimizedClient() *http.Client {
    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   5 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 5 * time.Second,
    }

    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
}
```

---

## 8. Performance Checklist

- [ ] Profile before optimizing
- [ ] Use appropriate data structures
- [ ] Implement caching where beneficial
- [ ] Optimize database queries
- [ ] Use connection pooling
- [ ] Minimize allocations
- [ ] Leverage concurrency appropriately
- [ ] Monitor performance metrics
- [ ] Set SLOs and alerts
- [ ] Regular performance testing

---

## References

1. Designing Data-Intensive Applications - Martin Kleppmann
2. Systems Performance - Brendan Gregg
3. Go Performance Patterns
4. Database Internals - Alex Petrov

---

**Quality Rating**: S (16+ KB)
**Last Updated**: 2026-04-02
