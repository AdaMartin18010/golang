# TS-DB-013: Database Connection Pooling

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #database #connection-pool #performance #golang #sql
> **权威来源**:
>
> - [database/sql Connection Pool](https://go.dev/doc/database/manage-connections) - Go team
> - [PostgreSQL Connection Pooling](https://www.postgresql.org/docs/current/runtime-config-connection.html) - PostgreSQL

---

## 1. Connection Pool Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Database Connection Pool Architecture                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Application                                                                  │
│     │                                                                         │
│     ▼                                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Connection Pool (sql.DB)                          │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                    Idle Connection Pool                        │  │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐          │  │   │
│  │  │  │  Conn 1 │  │  Conn 2 │  │  Conn 3 │  │  ...    │          │  │   │
│  │  │  │ (Idle)  │  │ (Idle)  │  │ (Idle)  │  │         │          │  │   │
│  │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘          │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                   Active Connections                           │  │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐          │  │   │
│  │  │  │  Conn A │  │  Conn B │  │  Conn C │  │  Conn D │          │  │   │
│  │  │  │ (In Tx) │  │ (Query) │  │ (Query) │  │ (In Tx) │          │  │   │
│  │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘          │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  │  Pool Configuration:                                                 │   │
│  │  - MaxOpenConns: Maximum open connections (default: unlimited)      │   │
│  │  - MaxIdleConns: Maximum idle connections (default: 2)              │   │
│  │  - ConnMaxLifetime: Maximum lifetime of a connection                │   │
│  │  - ConnMaxIdleTime: Maximum idle time before close (Go 1.15+)       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                               │
│                              │                                               │
│  ┌───────────────────────────┼─────────────────────────────────────────┐   │
│  │                    Database Server                                   │   │
│  │  ┌────────────────────────┴─────────────────────────────────────┐   │   │
│  │  │                     Connection Slots                           │   │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐          │   │   │
│  │  │  │ Conn 1  │  │ Conn 2  │  │ Conn 3  │  │  ...    │          │   │   │
│  │  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘          │   │   │
│  │  │  max_connections (PostgreSQL default: 100)                    │   │   │
│  │  └───────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Connection Pool Configuration

### 2.1 Go sql.DB Pool Settings

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "runtime"
    "time"

    _ "github.com/lib/pq"
)

func configurePool(db *sql.DB) {
    // Maximum number of open connections to the database
    // Formula: (core_count * 2) + effective_spindle_count
    // For SSD/cloud: core_count * 2
    // Default: 0 (unlimited - dangerous!)
    db.SetMaxOpenConns(25)

    // Maximum number of connections in the idle connection pool
    // Should be <= MaxOpenConns
    // Default: 2 (often too low for production)
    db.SetMaxIdleConns(10)

    // Maximum amount of time a connection may be reused
    // Prevents issues with stale connections, connection leaks
    // Should be less than database's wait_timeout
    // MySQL: typically 1 hour, PostgreSQL: typically 8 hours
    db.SetConnMaxLifetime(30 * time.Minute)

    // Maximum amount of time a connection can be idle (Go 1.15+)
    // Closes idle connections sooner, reduces memory usage
    db.SetConnMaxIdleTime(10 * time.Minute)
}

// Calculate optimal pool size based on environment
func calculatePoolSize() int {
    cores := runtime.NumCPU()
    // Conservative formula for web applications
    // Each request typically takes ~10-50ms of DB time
    // With concurrent requests, we need enough connections
    return cores * 4
}

func main() {
    dsn := "postgres://user:password@localhost/dbname?sslmode=disable"

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Configure pool
    configurePool(db)

    // Verify with ping
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        log.Fatal(err)
    }

    // Print pool stats
    printStats(db)
}

func printStats(db *sql.DB) {
    stats := db.Stats()
    fmt.Printf("Open Connections: %d\n", stats.OpenConnections)
    fmt.Printf("In Use: %d\n", stats.InUse)
    fmt.Printf("Idle: %d\n", stats.Idle)
    fmt.Printf("Wait Count: %d\n", stats.WaitCount)
    fmt.Printf("Wait Duration: %v\n", stats.WaitDuration)
    fmt.Printf("Max Idle Closed: %d\n", stats.MaxIdleClosed)
    fmt.Printf("Max Lifetime Closed: %d\n", stats.MaxLifetimeClosed)
}
```

### 2.2 Pool Size Guidelines

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Pool Size Recommendations                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Small Application (1-2 cores):                                             │
│  - MaxOpenConns: 5-10                                                       │
│  - MaxIdleConns: 2-5                                                        │
│                                                                              │
│  Medium Application (4-8 cores):                                            │
│  - MaxOpenConns: 20-40                                                      │
│  - MaxIdleConns: 5-10                                                       │
│                                                                              │
│  Large Application (16+ cores):                                             │
│  - MaxOpenConns: 50-100                                                     │
│  - MaxIdleConns: 10-25                                                      │
│                                                                              │
│  Important Constraints:                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Database max_connections                                           │   │
│  │  ├── PostgreSQL: default 100                                        │   │
│  │  ├── MySQL: depends on RAM, typically 150-1000                      │   │
│  │  └── Consider: App instances × MaxOpenConns < max_connections       │   │
│  │                                                                     │   │
│  │  Example:                                                           │   │
│  │  - 5 app instances                                                  │   │
│  │  - 20 MaxOpenConns each                                             │   │
│  │  - Total: 100 connections                                           │   │
│  │  - PostgreSQL default max_connections: 100 ✓                        │   │
│  │  - Leave room for admin connections (usually 3 reserved)            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Connection Lifecycle

```go
// Connection lifecycle management

type ManagedDB struct {
    *sql.DB
    maxRetries int
    retryDelay  time.Duration
}

func (mdb *ManagedDB) QueryWithRetry(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    var rows *sql.Rows
    var err error

    for i := 0; i < mdb.maxRetries; i++ {
        rows, err = mdb.QueryContext(ctx, query, args...)
        if err == nil {
            return rows, nil
        }

        // Check if it's a retryable error
        if !isRetryableError(err) {
            return nil, err
        }

        time.Sleep(mdb.retryDelay * time.Duration(i+1))
    }

    return nil, err
}

func isRetryableError(err error) bool {
    if err == nil {
        return false
    }

    errStr := err.Error()
    retryableErrors := []string{
        "connection refused",
        "connection reset",
        "broken pipe",
        "deadlock",
        "lock wait timeout",
    }

    for _, retryable := range retryableErrors {
        if contains(errStr, retryable) {
            return true
        }
    }

    return false
}
```

---

## 4. Monitoring Pool Health

```go
// Health check and monitoring

type PoolMonitor struct {
    db     *sql.DB
    ticker *time.Ticker
    done   chan struct{}
}

func NewPoolMonitor(db *sql.DB, interval time.Duration) *PoolMonitor {
    return &PoolMonitor{
        db:     db,
        ticker: time.NewTicker(interval),
        done:   make(chan struct{}),
    }
}

func (pm *PoolMonitor) Start() {
    go func() {
        for {
            select {
            case <-pm.ticker.C:
                pm.checkHealth()
            case <-pm.done:
                return
            }
        }
    }()
}

func (pm *PoolMonitor) checkHealth() {
    stats := pm.db.Stats()

    // Alert if wait count is high (indicates pool exhaustion)
    if stats.WaitCount > 100 {
        log.Printf("WARNING: High connection wait count: %d", stats.WaitCount)
    }

    // Alert if idle connections are being closed frequently
    if stats.MaxIdleClosed > 10 {
        log.Printf("WARNING: High idle connection closure: %d", stats.MaxIdleClosed)
    }

    // Check if we're near max connections
    if stats.InUse > int(float64(stats.OpenConnections)*0.8) {
        log.Printf("WARNING: High connection usage: %d/%d",
            stats.InUse, stats.OpenConnections)
    }
}

func (pm *PoolMonitor) Stop() {
    close(pm.done)
    pm.ticker.Stop()
}
```

---

## 5. Pool Exhaustion Handling

```go
// Graceful handling of pool exhaustion

func queryWithTimeout(db *sql.DB, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // This will fail fast if pool is exhausted and timeout is reached
    rows, err := db.QueryContext(ctx, "SELECT * FROM users")
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("query timeout (pool may be exhausted): %w", err)
        }
        return err
    }
    defer rows.Close()

    return nil
}

// Circuit breaker for database operations
type DBCircuitBreaker struct {
    db           *sql.DB
    failureCount int
    threshold    int
    lastFailure  time.Time
    cooldown     time.Duration
    mu           sync.RWMutex
}

func (cb *DBCircuitBreaker) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    if cb.isOpen() {
        return nil, errors.New("database circuit breaker open")
    }

    rows, err := cb.db.QueryContext(ctx, query, args...)
    if err != nil {
        cb.recordFailure()
        return nil, err
    }

    cb.recordSuccess()
    return rows, nil
}

func (cb *DBCircuitBreaker) isOpen() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    if cb.failureCount < cb.threshold {
        return false
    }

    return time.Since(cb.lastFailure) < cb.cooldown
}

func (cb *DBCircuitBreaker) recordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.failureCount++
    cb.lastFailure = time.Now()
}

func (cb *DBCircuitBreaker) recordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.failureCount > 0 {
        cb.failureCount--
    }
}
```

---

## 6. Checklist

```
Connection Pool Checklist:
□ MaxOpenConns configured appropriately
□ MaxIdleConns configured (not default 2)
□ ConnMaxLifetime set (prevent stale connections)
□ ConnMaxIdleTime set (Go 1.15+)
□ Pool size accounts for multiple app instances
□ Monitor pool statistics
□ Handle pool exhaustion gracefully
□ Set query timeouts
□ Use context for cancellation
□ Test pool behavior under load
```
