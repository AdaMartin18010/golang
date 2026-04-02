# TS-DB-012: Database Replication Strategies

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #replication #postgresql #mysql #high-availability #master-slave
> **权威来源**:
>
> - [PostgreSQL Streaming Replication](https://www.postgresql.org/docs/current/warm-standby.html) - PostgreSQL
> - [MySQL Replication](https://dev.mysql.com/doc/refman/8.0/en/replication.html) - MySQL

---

## 1. Replication Architecture

### 1.1 Master-Slave Replication

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Master-Slave Replication                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Master (Primary)                             │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                      Write Operations                          │  │   │
│  │  │  INSERT ──► WAL (Write-Ahead Log) ──► Data Files               │  │   │
│  │  │  UPDATE ──►                                                    │  │   │
│  │  │  DELETE ──►                                                    │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                              │                                      │   │
│  │                              ▼                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                    WAL Archiver / Streamer                     │  │   │
│  │  │  - Continuous archiving to archive directory                   │  │   │
│  │  │  - Streaming replication to standby                            │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └──────────────────────────────┬──────────────────────────────────────┘   │
│                                 │                                            │
│                    ┌────────────┼────────────┐                               │
│                    │            │            │                               │
│                    ▼            ▼            ▼                               │
│  ┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐    │
│  │   Standby 1         │ │   Standby 2         │ │   Standby N         │    │
│  │  (Hot Standby)      │ │  (Hot Standby)      │ │  (Hot Standby)      │    │
│  │                     │ │                     │ │                     │    │
│  │  ┌───────────────┐  │ │  ┌───────────────┐  │ │  ┌───────────────┐  │    │
│  │  │ WAL Receiver  │◄─┘ │  │ WAL Receiver  │◄─┘ │  │ WAL Receiver  │◄─┘    │
│  │  └───────┬───────┘    │  └───────┬───────┘    │  └───────┬───────┘       │
│  │          │             │          │             │          │              │
│  │  ┌───────▼───────┐    │  ┌───────▼───────┐    │  ┌───────▼───────┐       │
│  │  │ WAL Applier   │    │  │ WAL Applier   │    │  │ WAL Applier   │       │
│  │  └───────┬───────┘    │  └───────┬───────┘    │  └───────┬───────┘       │
│  │          │             │          │             │          │              │
│  │  ┌───────▼───────┐    │  ┌───────▼───────┐    │  ┌───────▼───────┐       │
│  │  │  Data Files   │    │  │  Data Files   │    │  │  Data Files   │       │
│  │  └───────────────┘    │  └───────────────┘    │  └───────────────┘       │
│  │                       │                       │                          │
│  │  ┌───────────────┐    │  ┌───────────────┐    │  ┌───────────────┐       │
│  │  │ Read Queries  │    │  │ Read Queries  │    │  │ Read Queries  │       │
│  │  └───────────────┘    │  └───────────────┘    │  └───────────────┘       │
│  └─────────────────────┘ └─────────────────────┘ └─────────────────────┘    │
│                                                                              │
│  Replication Lag:                                                            │
│  - Measured as: current_time - last_replay_timestamp                        │
│  - Acceptable: < 1 second for most applications                             │
│  - High lag: > 10 seconds indicates issues                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 PostgreSQL Streaming Replication

```go
// Go client with read replica support
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

type DBCluster struct {
    primary *sql.DB
    replicas []*sql.DB
    nextReplica int
}

func NewDBCluster(primaryDSN string, replicaDSNs []string) (*DBCluster, error) {
    primary, err := sql.Open("postgres", primaryDSN)
    if err != nil {
        return nil, err
    }

    replicas := make([]*sql.DB, len(replicaDSNs))
    for i, dsn := range replicaDSNs {
        db, err := sql.Open("postgres", dsn)
        if err != nil {
            return nil, err
        }
        replicas[i] = db
    }

    return &DBCluster{
        primary:  primary,
        replicas: replicas,
    }, nil
}

// Write to primary
func (c *DBCluster) Write(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    return c.primary.ExecContext(ctx, query, args...)
}

// Read from replica (round-robin)
func (c *DBCluster) Read(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    if len(c.replicas) == 0 {
        return c.primary.QueryContext(ctx, query, args...)
    }

    // Round-robin replica selection
    idx := c.nextReplica % len(c.replicas)
    c.nextReplica++

    return c.replicas[idx].QueryContext(ctx, query, args...)
}

// Check replication lag
func (c *DBCluster) CheckReplicationLag(ctx context.Context) error {
    for i, replica := range c.replicas {
        var lag sql.NullFloat64
        err := replica.QueryRowContext(ctx, `
            SELECT EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp()))
        `).Scan(&lag)

        if err != nil {
            return fmt.Errorf("replica %d: %w", i, err)
        }

        if lag.Valid && lag.Float64 > 10 {
            log.Printf("WARNING: Replica %d lag is %f seconds", i, lag.Float64)
        }
    }
    return nil
}
```

---

## 2. Replication Modes

### 2.1 Asynchronous vs Synchronous

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Asynchronous vs Synchronous Replication                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Asynchronous Replication:                                                   │
│  ┌─────────┐    Write    ┌─────────┐    WAL    ┌─────────┐                 │
│  │  Client │────────────►│  Master │──────────►│ Standby │                 │
│  └─────────┘             └────┬────┘  (async)  └─────────┘                 │
│                               │                                              │
│                               ▼                                              │
│                          ┌─────────┐                                        │
│                          │ Commit  │  (Immediate)                           │
│                          └─────────┘                                        │
│                                                                              │
│  - Fast commit (no standby wait)                                            │
│  - Risk of data loss if master fails                                        │
│  - Default for PostgreSQL and MySQL                                         │
│                                                                              │
│  Synchronous Replication:                                                    │
│  ┌─────────┐    Write    ┌─────────┐    WAL    ┌─────────┐                 │
│  │  Client │────────────►│  Master │──────────►│ Standby │                 │
│  └─────────┘             └────┬────┘  (sync)   └────┬────┘                 │
│                               │                     │                        │
│                               │◄──── Ack ──────────┘                        │
│                               ▼                                              │
│                          ┌─────────┐                                        │
│                          │ Commit  │  (After standby ack)                   │
│                          └─────────┘                                        │
│                                                                              │
│  - No data loss on failover                                                 │
│  - Higher latency (network round-trip)                                      │
│  - Risk of transaction blocking if standby fails                            │
│                                                                              │
│  PostgreSQL Configuration:                                                   │
│  synchronous_commit = remote_apply  # Wait for standby apply                │
│  synchronous_commit = remote_write  # Wait for standby receive              │
│  synchronous_commit = on            # Wait for local fsync                  │
│  synchronous_commit = off           # No wait (async)                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Failover Handling

```go
// Health checking and failover

type Replica struct {
    db     *sql.DB
    dsn    string
    isHealthy bool
    lag    float64
}

type FailoverManager struct {
    primary  *Replica
    standbys []*Replica
    mu       sync.RWMutex
}

func (fm *FailoverManager) HealthCheck(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fm.checkPrimary(ctx)
            fm.checkStandbys(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (fm *FailoverManager) checkPrimary(ctx context.Context) {
    fm.mu.Lock()
    defer fm.mu.Unlock()

    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    err := fm.primary.db.PingContext(ctx)
    wasHealthy := fm.primary.isHealthy
    fm.primary.isHealthy = (err == nil)

    if wasHealthy && !fm.primary.isHealthy {
        log.Println("Primary is down, initiating failover...")
        fm.performFailover()
    }
}

func (fm *FailoverManager) performFailover() {
    // Find healthiest standby
    var bestStandby *Replica
    for _, standby := range fm.standbys {
        if standby.isHealthy && (bestStandby == nil || standby.lag < bestStandby.lag) {
            bestStandby = standby
        }
    }

    if bestStandby == nil {
        log.Fatal("No healthy standby available for failover")
    }

    log.Printf("Promoting %s to primary", bestStandby.dsn)
    // In production: use pg_ctl promote, repmgr, or Patroni
}

// Connection with automatic failover
func (fm *FailoverManager) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    fm.mu.RLock()
    primary := fm.primary
    standbys := make([]*Replica, len(fm.standbys))
    copy(standbys, fm.standbys)
    fm.mu.RUnlock()

    // Try primary first
    if primary.isHealthy {
        rows, err := primary.db.QueryContext(ctx, query, args...)
        if err == nil {
            return rows, nil
        }
    }

    // Fall back to standbys
    for _, standby := range standbys {
        if standby.isHealthy {
            rows, err := standby.db.QueryContext(ctx, query, args...)
            if err == nil {
                log.Println("Using standby for read")
                return rows, nil
            }
        }
    }

    return nil, errors.New("no healthy database available")
}
```

---

## 4. Checklist

```
Replication Checklist:
□ Asynchronous or synchronous chosen appropriately
□ Replication lag monitored
□ Automatic failover configured (Patroni/repmgr)
□ Read replica load balancing
□ Health checks for all nodes
□ Connection retry logic
□ Quorum requirements for sync replication
□ Archive command configured for WAL
□ Backup from standby configured
□ Replication slots for consistency
```
