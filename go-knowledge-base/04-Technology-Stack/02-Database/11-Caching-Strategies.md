# TS-DB-011: Caching Strategies and Patterns

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #caching #redis #cache-strategy #cache-aside #write-through
> **权威来源**:
>
> - [Redis Caching Strategies](https://redis.io/docs/manual/client-side-caching/) - Redis
> - [Cache Patterns](https://docs.microsoft.com/en-us/azure/architecture/patterns/cache-aside) - Microsoft Azure

---

## 1. Cache Architecture Patterns

### 1.1 Pattern Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Cache Architecture Patterns                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Cache-Aside (Lazy Loading)                                              │
│  ┌─────────────┐         Cache Miss          ┌─────────────┐               │
│  │ Application │◄────────────────────────────│   Cache     │               │
│  └──────┬──────┘                             └─────────────┘               │
│         │                                                                    │
│         │ Read Request                                                       │
│         ▼                                                                    │
│  ┌─────────────┐         Not Found           ┌─────────────┐               │
│  │    Cache    │────────────────────────────►│  Database   │               │
│  └─────────────┘                             └──────┬──────┘               │
│                                                     │                        │
│                                                     │ Write to Cache        │
│                                                     ▼                        │
│                                              ┌─────────────┐               │
│                                              │ Return Data │               │
│                                              └─────────────┘               │
│                                                                              │
│  2. Read-Through                                                              │
│  ┌─────────────┐                             ┌─────────────┐               │
│  │ Application │◄────────────────────────────│    Cache    │               │
│  └─────────────┘    Cache manages loading    │   (Manages  │               │
│                                              │   loading)  │               │
│                                              └──────┬──────┘               │
│                                                     │                        │
│                                                     ▼                        │
│                                              ┌─────────────┐               │
│                                              │  Database   │               │
│                                              └─────────────┘               │
│                                                                              │
│  3. Write-Through                                                             │
│  ┌─────────────┐                             ┌─────────────┐               │
│  │ Application │────────────────────────────►│    Cache    │               │
│  └─────────────┘                             └──────┬──────┘               │
│                                                     │                        │
│                                                     │ Sync Write            │
│                                                     ▼                        │
│                                              ┌─────────────┐               │
│                                              │  Database   │               │
│                                              └─────────────┘               │
│                                                                              │
│  4. Write-Behind (Write-Back)                                                 │
│  ┌─────────────┐                             ┌─────────────┐               │
│  │ Application │────────────────────────────►│    Cache    │               │
│  └─────────────┘                             └──────┬──────┘               │
│                                                     │                        │
│                                                     │ Async Write           │
│                                                     ▼                        │
│                                              ┌─────────────┐               │
│                                              │  Database   │               │
│                                              │  (Async)    │               │
│                                              └─────────────┘               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Cache-Aside Implementation

```go
package cache

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "time"

    "github.com/go-redis/redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type CacheAside struct {
    rdb        *redis.Client
    ttl        time.Duration
    lockTTL    time.Duration
    loadFunc   func(ctx context.Context, key string) (interface{}, error)
}

func NewCacheAside(rdb *redis.Client, ttl time.Duration, loadFunc func(ctx context.Context, key string) (interface{}, error)) *CacheAside {
    return &CacheAside{
        rdb:      rdb,
        ttl:      ttl,
        lockTTL:  10 * time.Second,
        loadFunc: loadFunc,
    }
}

// Get with cache-aside pattern
func (c *CacheAside) Get(ctx context.Context, key string) (interface{}, error) {
    // 1. Try cache first
    val, err := c.rdb.Get(ctx, key).Result()
    if err == nil {
        var result interface{}
        if err := json.Unmarshal([]byte(val), &result); err != nil {
            return nil, err
        }
        return result, nil
    }

    if err != redis.Nil {
        return nil, err // Redis error
    }

    // 2. Cache miss - acquire lock for cache stampede protection
    lockKey := fmt.Sprintf("lock:%s", key)
    locked, err := c.rdb.SetNX(ctx, lockKey, "1", c.lockTTL).Result()
    if err != nil {
        return nil, err
    }

    if !locked {
        // Another goroutine is loading - wait and retry
        time.Sleep(100 * time.Millisecond)
        return c.Get(ctx, key)
    }
    defer c.rdb.Del(ctx, lockKey)

    // 3. Double-check after acquiring lock
    val, err = c.rdb.Get(ctx, key).Result()
    if err == nil {
        var result interface{}
        if err := json.Unmarshal([]byte(val), &result); err != nil {
            return nil, err
        }
        return result, nil
    }

    // 4. Load from database
    data, err := c.loadFunc(ctx, key)
    if err != nil {
        return nil, err
    }

    // 5. Write to cache (async fire-and-forget)
    go c.setCache(key, data)

    return data, nil
}

func (c *CacheAside) setCache(key string, data interface{}) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    jsonData, err := json.Marshal(data)
    if err != nil {
        return
    }

    c.rdb.Set(ctx, key, jsonData, c.ttl)
}

// Invalidate cache
func (c *CacheAside) Invalidate(ctx context.Context, key string) error {
    return c.rdb.Del(ctx, key).Err()
}

// Update cache on write
func (c *CacheAside) Set(ctx context.Context, key string, value interface{}) error {
    jsonData, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return c.rdb.Set(ctx, key, jsonData, c.ttl).Err()
}
```

---

## 3. Cache Eviction Strategies

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Cache Eviction Strategies                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. LRU (Least Recently Used)                                               │
│     - Evict least recently accessed items                                   │
│     - Good for: General purpose caching                                     │
│     - Implementation: Linked hash map or timestamp tracking                 │
│                                                                              │
│  2. LFU (Least Frequently Used)                                             │
│     - Evict least frequently accessed items                                 │
│     - Good for: Workloads with popular items                                │
│     - Implementation: Counter + min-heap                                    │
│                                                                              │
│  3. TTL (Time To Live)                                                      │
│     - Expire items after specified time                                     │
│     - Good for: Session data, temporary data                                │
│     - Implementation: Background cleanup or lazy expiration                 │
│                                                                              │
│  4. Random Replacement                                                      │
│     - Randomly select item to evict                                         │
│     - Good for: Simple implementation, uniform access                       │
│     - Implementation: Random selection                                      │
│                                                                              │
│  5. FIFO (First In First Out)                                               │
│     - Evict oldest items first                                              │
│     - Good for: Streaming data                                              │
│     - Implementation: Queue                                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Cache Problems and Solutions

### 4.1 Cache Stampede

```go
// Solution 1: Singleflight (one-at-a-time loading)
import "golang.org/x/sync/singleflight"

type StampedeProtectedCache struct {
    rdb   *redis.Client
    group singleflight.Group
}

func (c *StampedeProtectedCache) Get(ctx context.Context, key string, load func() (interface{}, error)) (interface{}, error) {
    val, err := c.rdb.Get(ctx, key).Result()
    if err == nil {
        return val, nil
    }

    // Singleflight: only one goroutine loads for this key
    v, err, _ := c.group.Do(key, func() (interface{}, error) {
        return load()
    })

    return v, err
}

// Solution 2: Probabilistic early expiration
func (c *Cache) GetWithEarlyExpire(ctx context.Context, key string) (string, error) {
    val, ttl, err := c.rdb.Get(ctx, key).Result(), c.rdb.TTL(ctx, key).Val(), nil

    if err == redis.Nil {
        return "", ErrCacheMiss
    }
    if err != nil {
        return "", err
    }

    // If TTL is low and random condition met, treat as miss
    if ttl < 5*time.Second && rand.Float64() < 0.1 {
        return "", ErrCacheMiss
    }

    return val, nil
}
```

### 4.2 Cache Penetration

```go
// Solution: Cache empty results with short TTL
func (c *Cache) GetWithNullCache(ctx context.Context, key string) (interface{}, error) {
    val, err := c.rdb.Get(ctx, key).Result()
    if err == nil {
        if val == "__NULL__" {
            return nil, ErrNotFound
        }
        return val, nil
    }

    if err != redis.Nil {
        return nil, err
    }

    // Load from DB
    data, err := c.loadFromDB(key)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            // Cache null to prevent repeated DB queries
            c.rdb.Set(ctx, key, "__NULL__", 1*time.Minute)
            return nil, ErrNotFound
        }
        return nil, err
    }

    c.rdb.Set(ctx, key, data, c.ttl)
    return data, nil
}
```

### 4.3 Cache Breakdown (Hotspot)

```go
// Solution: Local cache + distributed cache
type TwoLevelCache struct {
    local  *ristretto.Cache  // In-memory LRU
    remote *redis.Client
}

func (c *TwoLevelCache) Get(ctx context.Context, key string) (interface{}, error) {
    // 1. Check local cache
    if val, found := c.local.Get(key); found {
        return val, nil
    }

    // 2. Check remote cache
    val, err := c.remote.Get(ctx, key).Result()
    if err == nil {
        // Populate local cache
        c.local.Set(key, val, 1)
        return val, nil
    }

    return nil, ErrCacheMiss
}
```

---

## 5. Multi-Layer Caching

```go
// L1: Local (in-memory) - fastest, smallest
// L2: Redis - distributed, medium
// L3: Database - persistent, slowest

type MultiLayerCache struct {
    l1    *ristretto.Cache  // Local
    l2    *redis.Client     // Redis
    db    *sql.DB           // Database
    ttlL1 time.Duration
    ttlL2 time.Duration
}

func (c *MultiLayerCache) Get(ctx context.Context, key string, loader func(string) (interface{}, error)) (interface{}, error) {
    // L1 cache
    if val, found := c.l1.Get(key); found {
        return val, nil
    }

    // L2 cache
    val, err := c.l2.Get(ctx, key).Result()
    if err == nil {
        // Populate L1
        c.l1.Set(key, val, 0)
        return val, nil
    }

    // Database
    data, err := loader(key)
    if err != nil {
        return nil, err
    }

    // Populate caches
    c.l1.Set(key, data, 0)
    c.l2.Set(ctx, key, data, c.ttlL2)

    return data, nil
}
```

---

## 6. Checklist

```
Caching Strategy Checklist:
□ Appropriate pattern chosen (cache-aside/read-through/write-through)
□ Cache stampede protection implemented
□ Empty result caching (prevent cache penetration)
□ TTL configured appropriately
□ Eviction strategy matches workload
□ Cache invalidation strategy defined
□ Multi-layer cache if needed
□ Monitoring for hit/miss rates
□ Warm-up strategy for cold start
```
