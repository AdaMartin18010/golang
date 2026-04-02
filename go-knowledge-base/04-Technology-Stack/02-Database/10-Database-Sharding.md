# TS-DB-010: Database Sharding Strategies

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #sharding #partitioning #scalability #database #distributed
> **权威来源**:
>
> - [Database Sharding](https://docs.microsoft.com/en-us/azure/architecture/patterns/sharding) - Microsoft Azure
> - [PostgreSQL Partitioning](https://www.postgresql.org/docs/current/ddl-partitioning.html) - PostgreSQL

---

## 1. Sharding Architecture

### 1.1 Horizontal Partitioning

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Database Sharding Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Before Sharding:                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Single Database                                 │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                    Users Table (100M rows)                     │  │   │
│  │  │  ID 1-100,000,000                                              │  │   │
│  │  │  CPU: 100%    Memory: 95%    Disk: 90%                         │  │   │
│  │  │  Query time: 5s+                                               │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  After Sharding (by user_id % 4):                                            │
│  ┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐   │
│  │    Shard 0          │ │    Shard 1          │ │    Shard 2          │   │
│  │  ┌───────────────┐  │ │  ┌───────────────┐  │ │  ┌───────────────┐  │   │
│  │  │ Users (ID%4=0)│  │ │  │ Users (ID%4=1)│  │ │  │ Users (ID%4=2)│  │   │
│  │  │ 25M rows      │  │ │  │ 25M rows      │  │ │  │ 25M rows      │  │   │
│  │  │ CPU: 30%      │  │ │  │ CPU: 28%      │  │ │  │ CPU: 32%      │  │   │
│  │  └───────────────┘  │ │  └───────────────┘  │ │  └───────────────┘  │   │
│  └─────────────────────┘ └─────────────────────┘ └─────────────────────┘   │
│                                                                              │
│  ┌─────────────────────┐                                                    │
│  │    Shard 3          │                                                    │
│  │  ┌───────────────┐  │                                                    │
│  │  │ Users (ID%4=3)│  │                                                    │
│  │  │ 25M rows      │  │                                                    │
│  │  │ CPU: 29%      │  │                                                    │
│  │  └───────────────┘  │                                                    │
│  └─────────────────────┘                                                    │
│                                                                              │
│  Query time: <100ms on each shard                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Sharding Strategies

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Sharding Strategies                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Hash Sharding                                                            │
│     shard = hash(key) % num_shards                                          │
│                                                                              │
│     Pros: Even distribution                                                  │
│     Cons: Range queries across shards, rebalancing on scale                 │
│                                                                              │
│  2. Range Sharding                                                           │
│     Shard 1: ID 1-1,000,000                                                 │
│     Shard 2: ID 1,000,001-2,000,000                                         │
│     Shard 3: ID 2,000,001-3,000,000                                         │
│                                                                              │
│     Pros: Efficient range queries                                            │
│     Cons: Hot spots (latest data in last shard)                             │
│                                                                              │
│  3. List Sharding                                                            │
│     Shard 1: US, CA                                                         │
│     Shard 2: EU countries                                                   │
│     Shard 3: APAC countries                                                 │
│                                                                              │
│     Pros: Data locality, compliance                                          │
│     Cons: Uneven distribution                                                │
│                                                                              │
│  4. Composite Sharding                                                       │
│     Primary: region (list)                                                   │
│     Secondary: user_id % N (hash)                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Sharding Implementation

```go
package sharding

import (
    "context"
    "database/sql"
    "fmt"
    "hash/fnv"
)

// ShardManager manages database shards
type ShardManager struct {
    shards []*sql.DB
    config *ShardingConfig
}

type ShardingConfig struct {
    ShardCount int
    Strategy   ShardingStrategy
}

type ShardingStrategy int

const (
    HashStrategy ShardingStrategy = iota
    RangeStrategy
    ListStrategy
)

// NewShardManager creates a new shard manager
func NewShardManager(shards []*sql.DB, config *ShardingConfig) *ShardManager {
    return &ShardManager{
        shards: shards,
        config: config,
    }
}

// GetShard returns the appropriate shard for a key
func (sm *ShardManager) GetShard(key string) (*sql.DB, error) {
    shardIndex := sm.calculateShard(key)
    if shardIndex < 0 || shardIndex >= len(sm.shards) {
        return nil, fmt.Errorf("invalid shard index: %d", shardIndex)
    }
    return sm.shards[shardIndex], nil
}

func (sm *ShardManager) calculateShard(key string) int {
    switch sm.config.Strategy {
    case HashStrategy:
        return sm.hashShard(key)
    case RangeStrategy:
        return sm.rangeShard(key)
    default:
        return sm.hashShard(key)
    }
}

func (sm *ShardManager) hashShard(key string) int {
    h := fnv.New32a()
    h.Write([]byte(key))
    return int(h.Sum32()) % sm.config.ShardCount
}

func (sm *ShardManager) rangeShard(key string) int {
    // Parse key as integer for range sharding
    var id int
    fmt.Sscanf(key, "%d", &id)

    // Assuming even distribution
    shardSize := 1000000 // 1M per shard
    return id / shardSize
}

// Query executes a query on the appropriate shard
func (sm *ShardManager) Query(ctx context.Context, shardKey string, query string, args ...interface{}) (*sql.Rows, error) {
    shard, err := sm.GetShard(shardKey)
    if err != nil {
        return nil, err
    }
    return shard.QueryContext(ctx, query, args...)
}

// Execute executes a write on the appropriate shard
func (sm *ShardManager) Execute(ctx context.Context, shardKey string, query string, args ...interface{}) (sql.Result, error) {
    shard, err := sm.GetShard(shardKey)
    if err != nil {
        return nil, err
    }
    return shard.ExecContext(ctx, query, args...)
}

// QueryAll executes query on all shards and merges results
func (sm *ShardManager) QueryAll(ctx context.Context, query string, args ...interface{}) ([]*sql.Rows, error) {
    results := make([]*sql.Rows, 0, len(sm.shards))

    for _, shard := range sm.shards {
        rows, err := shard.QueryContext(ctx, query, args...)
        if err != nil {
            // Close already opened results
            for _, r := range results {
                r.Close()
            }
            return nil, err
        }
        results = append(results, rows)
    }

    return results, nil
}
```

---

## 3. Cross-Shard Operations

```go
// CrossShardTransaction coordinates transactions across shards
type CrossShardTransaction struct {
    shards map[int]*sql.Tx
    mu     sync.Mutex
}

func (cst *CrossShardTransaction) Begin(shards []*sql.DB) error {
    cst.shards = make(map[int]*sql.Tx)

    for i, shard := range shards {
        tx, err := shard.Begin()
        if err != nil {
            // Rollback already started transactions
            cst.Rollback()
            return err
        }
        cst.shards[i] = tx
    }

    return nil
}

func (cst *CrossShardTransaction) Commit() error {
    // Two-phase commit would be needed for true ACID
    // Simplified: commit all, hope for the best

    var lastErr error
    for _, tx := range cst.shards {
        if err := tx.Commit(); err != nil {
            lastErr = err
        }
    }

    return lastErr
}

func (cst *CrossShardTransaction) Rollback() {
    for _, tx := range cst.shards {
        tx.Rollback()
    }
}
```

---

## 4. Checklist

```
Sharding Checklist:
□ Sharding strategy chosen appropriately
□ Shard key selected (high cardinality)
□ Cross-shard operations minimized
□ Rebalancing strategy defined
□ Query routing layer implemented
□ Monitoring for shard balance
□ Backup strategy per shard
□ Failover per shard
□ Application code handles shard failures
```
