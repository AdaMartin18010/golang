# TS-DB-004: Redis Internals and Go Integration

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #redis #cache #data-structures #performance #go-redis
> **权威来源**:
>
> - [Redis Documentation](https://redis.io/documentation) - Redis Labs
> - [Redis Internals](https://redis.io/topics/internals) - Implementation details
> - [go-redis Documentation](https://redis.uptrace.dev/) - Go client
> - [Redis Cluster Specification](https://redis.io/topics/cluster-spec) - Distributed mode

---

## 1. Redis Architecture Overview

### 1.1 Single-Threaded Event Loop

```
┌─────────────────────────────────────────────────────────────────┐
│                      Redis Server Architecture                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐     ┌──────────────────────────────────────┐  │
│  │   Clients    │────►│          Event Loop (Single Thread)   │  │
│  └──────────────┘     ├──────────────────────────────────────┤  │
│                       │                                      │  │
│                       │  ┌─────────┐    ┌─────────────────┐  │  │
│                       │  │  AE (   │    │  Command Table  │  │  │
│                       │  │ epoll/  │───►│  (Hash Table)   │  │  │
│                       │  │ kqueue) │    └────────┬────────┘  │  │
│                       │  └────┬────┘             │           │  │
│                       │       │                  ▼           │  │
│                       │       │         ┌─────────────────┐  │  │
│                       │       │         │  Data Structures │  │  │
│                       │       │         │  (SDS, Dict,    │  │  │
│                       │       │         │   Ziplist, etc) │  │  │
│                       │       │         └────────┬────────┘  │  │
│                       │       │                  │           │  │
│                       │       └──────────────────┘           │  │
│                       │                  │                    │  │
│                       │                  ▼                    │  │
│                       │         ┌─────────────────┐          │  │
│                       │         │   Persistence   │          │  │
│                       │         │  (AOF/RDB)      │          │  │
│                       │         └─────────────────┘          │  │
│                       │                                      │  │
│  ┌──────────────┐     │  ┌────────────────────────────────┐  │  │
│  │  Background  │◄────┘  │  BIO Threads (IO intensive)   │  │  │
│  │  Save/IO     │        │  - AOF fsync                   │  │  │
│  └──────────────┘        │  - RDB save                     │  │  │
│                          │  - Lazy free                    │  │  │
│                          └────────────────────────────────┘  │  │
│                                                              │  │
└──────────────────────────────────────────────────────────────┘  │
```

**Event Loop Pseudocode:**

```c
void aeMain(aeEventLoop *eventLoop) {
    while (!eventLoop->stop) {
        // 1. Check ready sockets (epoll_wait)
        ready_sockets = aeApiPoll(eventLoop, timeout);

        // 2. Process file events (client requests)
        for (socket in ready_sockets) {
            read_query_from_client(socket);
            process_command(client);
            add_reply_to_client(client);
        }

        // 3. Process time events (expired keys, etc.)
        process_time_events(eventLoop);

        // 4. Write pending replies
        write_to_clients();
    }
}
```

### 1.2 Memory Efficiency

Redis uses specialized data structures for memory optimization:

```
┌─────────────────────────────────────────────────────────────────┐
│                 Redis Data Structure Selection                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  String:                                                         │
│  - Integer (int64_t) when possible                              │
│  - Embedded string (<= 44 bytes): SDS embedded in redisObject   │
│  - Raw string: SDS + separate allocation                        │
│                                                                  │
│  Hash:                                                           │
│  - ziplist (entries < 512, entry < 64 bytes)                    │
│  - hash table (dict) when larger                                │
│                                                                  │
│  List:                                                           │
│  - quicklist (ziplist + linked list of ziplists)               │
│  - linked list (rarely used now)                                │
│                                                                  │
│  Set:                                                            │
│  - intset (all integers, < 512 entries)                         │
│  - hash table otherwise                                         │
│                                                                  │
│  Sorted Set:                                                     │
│  - ziplist (entries < 128, member < 64 bytes)                   │
│  - skip list + hash table otherwise                             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Core Data Structures Internals

### 2.1 Simple Dynamic String (SDS)

```c
// SDS Type 5 (not used, embedded in flags)
// SDS Type 8 (for strings < 256 bytes)
struct __attribute__ ((__packed__)) sdshdr8 {
    uint8_t len;        // Current string length
    uint8_t alloc;      // Total allocated bytes (excluding header)
    unsigned char flags; // Type identifier (3 lsb)
    char buf[];         // Actual string data (null-terminated)
};

// SDS Type 16, 32, 64 follow similar pattern
```

**Benefits over C strings:**

- O(1) length lookup
- No buffer overflow (pre-allocated)
- Binary safe (can contain null bytes)
- Amortized O(1) append
- Reduced memory reallocation via over-allocation

### 2.2 Skip List (Sorted Sets)

```c
typedef struct zskiplistNode {
    sds ele;                    // Element value
    double score;               // Sort score
    struct zskiplistNode *backward;  // Back pointer for reverse iteration
    struct zskiplistLevel {
        struct zskiplistNode *forward;
        unsigned long span;     // Distance to next node
    } level[];                  // Variable level array
} zskiplistNode;

typedef struct zskiplist {
    struct zskiplistNode *header, *tail;
    unsigned long length;
    int level;                  // Current max level
} zskiplist;
```

**Skip List Properties:**

- Average O(log n) search/insert/delete
- Deterministic level via random(): P(level=k) = 0.25^k
- Max level: 32 (supports 2^32 elements)
- Also maintains hash table for O(1) ZSCORE

### 2.3 QuickList (Lists)

```c
typedef struct quicklistNode {
    struct quicklistNode *prev;
    struct quicklistNode *next;
    unsigned char *zl;          // ziplist or NULL
    unsigned int sz;            // ziplist size in bytes
    unsigned int count : 16;    // count of items in ziplist
    unsigned int encoding : 2;  // RAW==1 or LZF==2
    unsigned int container : 2; // NONE==1 or ZIPLIST==2
    unsigned int recompress : 1;
    unsigned int attempted_compress : 1;
    unsigned int extra : 10;
} quicklistNode;

typedef struct quicklist {
    quicklistNode *head;
    quicklistNode *tail;
    unsigned long count;        // Total count of all entries
    unsigned long len;          // Number of nodes
    int fill : 16;              // fill factor for individual nodes
    unsigned int compress : 16; // depth of end nodes not to compress
} quicklist;
```

---

## 3. Persistence Mechanisms

### 3.1 RDB (Redis Database Backup)

```
RDB Snapshot Process:
┌─────────────────────────────────────────────────────────────┐
│                                                              │
│  1. Fork child process (copy-on-write)                      │
│     Parent ───fork()───► Child                              │
│       │                    │                                │
│       │ continues          │ writes RDB to disk             │
│       │ serving clients    │                                │
│       │                    │                                │
│  2. Copy-on-write:                                         │
│     - Parent and child share memory pages                   │
│     - Modified pages are copied before write                │
│     - Memory overhead: proportional to write rate           │
│                                                              │
│  3. Atomic replacement:                                     │
│     - Child writes to temp.rdb                              │
│     - Rename to dump.rdb when complete                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 AOF (Append-Only File)

```
AOF Rewrite Process (BGREWRITEAOF):
┌─────────────────────────────────────────────────────────────┐
│                                                              │
│  1. Fork child process                                      │
│                                                              │
│  2. Child creates new AOF from current dataset              │
│     (more compact than original AOF)                        │
│                                                              │
│  3. Parent accumulates new writes in:                       │
│     - Regular AOF buffer                                    │
│     - AOF rewrite buffer                                    │
│                                                              │
│  4. When child finishes:                                    │
│     - Parent sends rewrite buffer to child                  │
│     - Child appends to new AOF                              │
│     - Atomic rename                                         │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**AOF Sync Policies:**

- `appendfsync always`: Every command (safest, slowest)
- `appendfsync everysec`: Every second (default, good balance)
- `appendfsync no`: OS decides (fastest, least safe)

---

## 4. Go Integration

### 4.1 go-redis Client

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// Connection pooling configuration
func createRedisClient() *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        Password:     "",
        DB:           0,

        // Connection pool settings
        PoolSize:     100,              // Maximum connections
        MinIdleConns: 10,               // Minimum idle connections
        MaxConnAge:   time.Hour,        // Maximum connection lifetime
        PoolTimeout:  30 * time.Second, // Timeout for getting conn from pool
        IdleTimeout:  10 * time.Minute, // Idle connection timeout

        // Timeouts
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })

    return rdb
}

// Basic operations with context
func basicOperations(ctx context.Context, rdb *redis.Client) error {
    // String operations
    err := rdb.Set(ctx, "key", "value", time.Hour).Err()
    if err != nil {
        return err
    }

    val, err := rdb.Get(ctx, "key").Result()
    if err == redis.Nil {
        fmt.Println("key does not exist")
    } else if err != nil {
        return err
    }
    fmt.Println("key", val)

    // Hash operations
    rdb.HSet(ctx, "user:1000", map[string]interface{}{
        "name": "John",
        "age": 30,
    })

    // Pipeline for batch operations
    pipe := rdb.Pipeline()
    pipe.Set(ctx, "key1", "value1", 0)
    pipe.Set(ctx, "key2", "value2", 0)
    pipe.Get(ctx, "key1")
    pipe.Get(ctx, "key2")

    cmds, err := pipe.Exec(ctx)
    if err != nil {
        return err
    }

    for _, cmd := range cmds {
        fmt.Println(cmd.Name(), cmd.Args())
    }

    return nil
}
```

### 4.2 Redis Cluster Support

```go
// Cluster client
func createClusterClient() *redis.ClusterClient {
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "localhost:7000",
            "localhost:7001",
            "localhost:7002",
            "localhost:7003",
            "localhost:7004",
            "localhost:7005",
        },

        // Pool settings per node
        PoolSize:     100,
        MinIdleConns: 10,

        // Cluster-specific
        MaxRedirects:   8,      // Maximum redirects for MOVED/ASK
        ReadOnly:       false,  // Allow reads from replicas
        RouteByLatency: false, // Route to nearest node
        RouteRandomly:  false, // Random node selection
    })

    return rdb
}

// Handling MOVED/ASK redirects automatically
func clusterOperations(ctx context.Context, rdb *redis.ClusterClient) {
    // Client automatically handles slot redirections
    rdb.Set(ctx, "key", "value", 0)

    // Cross-slot operations require hash tags
    // {user1000}.name and {user1000}.email go to same slot
    rdb.Set(ctx, "{user1000}.name", "John", 0)
    rdb.Set(ctx, "{user1000}.email", "john@example.com", 0)
}
```

### 4.3 Caching Patterns

```go
// Cache-Aside Pattern with go-redis

type CacheAside struct {
    rdb  *redis.Client
    load func(ctx context.Context, key string) (string, error)
    ttl  time.Duration
}

func (c *CacheAside) Get(ctx context.Context, key string) (string, error) {
    // 1. Try cache first
    val, err := c.rdb.Get(ctx, key).Result()
    if err == nil {
        return val, nil // Cache hit
    }
    if err != redis.Nil {
        return "", err // Redis error
    }

    // 2. Cache miss - load from source
    val, err = c.load(ctx, key)
    if err != nil {
        return "", err
    }

    // 3. Store in cache (async fire-and-forcel)
    go c.rdb.Set(ctx, key, val, c.ttl)

    return val, nil
}

// Write-Through Pattern
func (c *CacheAside) Set(ctx context.Context, key, value string) error {
    // 1. Write to database first
    if err := c.writeToDB(key, value); err != nil {
        return err
    }

    // 2. Update cache
    return c.rdb.Set(ctx, key, value, c.ttl).Err()
}

// Cache warming pattern
func (c *CacheAside) WarmCache(ctx context.Context, keys []string) {
    pipe := c.rdb.Pipeline()

    for _, key := range keys {
        // Skip if already cached
        pipe.Exists(ctx, key)
    }

    exists, _ := pipe.Exec(ctx)

    for i, cmd := range exists {
        if cmd.(*redis.IntCmd).Val() == 0 {
            // Not cached, load it
            go func(k string) {
                c.Get(ctx, k) // This will load and cache
            }(keys[i])
        }
    }
}
```

### 4.4 Distributed Lock with Redlock

```go
// Redlock algorithm implementation

type Redlock struct {
    clients []*redis.Client
    quorum  int
}

func NewRedlock(addrs []string) *Redlock {
    clients := make([]*redis.Client, len(addrs))
    for i, addr := range addrs {
        clients[i] = redis.NewClient(&redis.Options{Addr: addr})
    }
    return &Redlock{
        clients: clients,
        quorum:  len(addrs)/2 + 1,
    }
}

func (r *Redlock) Lock(ctx context.Context, resource string, ttl time.Duration) (*Lock, error) {
    token := generateUniqueToken()

    var acquired int
    start := time.Now()

    for _, client := range r.clients {
        ok, err := r.acquire(ctx, client, resource, token, ttl)
        if err == nil && ok {
            acquired++
        }
    }

    elapsed := time.Since(start)
    validity := ttl - elapsed - driftFactor

    if acquired >= r.quorum && validity > 0 {
        return &Lock{
            resource: resource,
            token:    token,
            validity: validity,
        }, nil
    }

    // Failed to acquire quorum, release all
    r.releaseAll(ctx, resource, token)
    return nil, errors.New("failed to acquire lock")
}

func (r *Redlock) acquire(ctx context.Context, client *redis.Client,
    resource, token string, ttl time.Duration) (bool, error) {

    // NX = only if not exists, PX = millisecond expiry
    return client.SetNX(ctx, resource, token, ttl).Result()
}

const releaseScript = `
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
`

func (r *Redlock) releaseAll(ctx context.Context, resource, token string) {
    for _, client := range r.clients {
        client.Eval(ctx, releaseScript, []string{resource}, token)
    }
}
```

---

## 5. Performance Tuning

### 5.1 Memory Optimization

```go
// Use hash tags to group related keys
// Instead of: user:1000:name, user:1000:email
// Use: user:1000 {user:1000}.name {user:1000}.email

// Or use hash for related fields
rdb.HSet(ctx, "user:1000", map[string]interface{}{
    "name": "John",
    "email": "john@example.com",
})

// Enable compression for large values
rdb.Set(ctx, "compressed", compress(data), ttl)

// Set appropriate maxmemory policy
// maxmemory-policy allkeys-lru (for cache)
// maxmemory-policy volatile-lru (for session store)
```

### 5.2 Pipeline Optimization

```go
// Batch operations with pipelines
func batchSet(ctx context.Context, rdb *redis.Client, items map[string]string) error {
    pipe := rdb.Pipeline()

    for key, value := range items {
        pipe.Set(ctx, key, value, 0)
    }

    _, err := pipe.Exec(ctx)
    return err
}

// For very large batches, chunk them
func chunkedBatch(ctx context.Context, rdb *redis.Client, items map[string]string, chunkSize int) error {
    keys := make([]string, 0, len(items))
    for k := range items {
        keys = append(keys, k)
    }

    for i := 0; i < len(keys); i += chunkSize {
        end := i + chunkSize
        if end > len(keys) {
            end = len(keys)
        }

        pipe := rdb.Pipeline()
        for _, key := range keys[i:end] {
            pipe.Set(ctx, key, items[key], 0)
        }
        if _, err := pipe.Exec(ctx); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 6. Monitoring and Health Checks

```go
// Health check
func healthCheck(ctx context.Context, rdb *redis.Client) error {
    return rdb.Ping(ctx).Err()
}

// Pool statistics
func logPoolStats(rdb *redis.Client) {
    stats := rdb.PoolStats()
    log.Printf("Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d",
        stats.Hits, stats.Misses, stats.Timeouts,
        stats.TotalConns, stats.IdleConns, stats.StaleConns)
}

// Slow log monitoring
func getSlowLogs(ctx context.Context, rdb *redis.Client) {
    logs, err := rdb.SlowLogGet(ctx, 10).Result()
    if err != nil {
        log.Println("Failed to get slow logs:", err)
        return
    }

    for _, log := range logs {
        fmt.Printf("ID=%d Duration=%s Command=%v\n",
            log.ID, log.Duration, log.Args)
    }
}
```

---

## 7. Checklist

```
Redis Configuration Checklist:
□ Set appropriate maxmemory
□ Configure maxmemory-policy
□ Enable persistence (RDB + AOF)
□ Set client output buffer limits
□ Configure timeout for idle connections
□ Enable slow log monitoring

Go Client Checklist:
□ Configure connection pool size
□ Set appropriate timeouts
□ Use context for cancellation
□ Handle redis.Nil for missing keys
□ Use pipelines for batch operations
□ Implement circuit breaker for resilience
```
