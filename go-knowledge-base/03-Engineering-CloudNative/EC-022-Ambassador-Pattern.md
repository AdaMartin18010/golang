# EC-022: Ambassador Pattern

## Problem Formalization

### The Remote Service Access Challenge

When applications need to connect to remote services (databases, message queues, caches, external APIs), they face challenges with connection management, retry logic, circuit breaking, and service discovery. Embedding this complexity directly in application code creates tight coupling and operational difficulties.

#### Problem Statement

Given:
- Application A needing to access remote service R
- Network characteristics N = {latency, packet loss, partition probability}
- Service R has properties: {replicas, health states, load distribution}

Find a connectivity component C such that:
```
Minimize: ConnectionManagementComplexity(A)
Maximize: Reliability(A → R)
Subject to:
  - A doesn't need to know R's topology
  - Connection pooling is optimized
  - Failover is transparent
  - Observability is comprehensive
```

### Ambassador vs Sidecar

While similar, Ambassador and Sidecar serve different purposes:

| Aspect | Ambassador | Sidecar |
|--------|------------|---------|
| **Primary Role** | Proxy remote connections | Add cross-cutting concerns |
| **Location** | Between app and external service | Alongside app (inbound+outbound) |
| **Scope** | Usually one external service type | All traffic |
| **Knowledge** | Deep service protocol knowledge | Generic HTTP/gRPC handling |

```
Sidecar Pattern (Traffic Management):
┌─────────────────────────────────────────────────────────────┐
│  Pod                                                        │
│  ┌──────────────┐         ┌──────────────┐                 │
│  │   Service    │◄───────►│   Sidecar    │◄────► Other     │
│  │              │         │ (Envoy)      │        Services  │
│  └──────────────┘         └──────────────┘                 │
│                                                             │
│  Manages ALL traffic (inbound and outbound)                 │
└─────────────────────────────────────────────────────────────┘

Ambassador Pattern (Service Proxy):
┌─────────────────────────────────────────────────────────────┐
│  Pod                                                        │
│  ┌──────────────┐         ┌──────────────┐                 │
│  │   Service    │◄───────►│  Ambassador  │◄────► Redis    │
│  │              │         │ (Twemproxy)  │        Cluster  │
│  └──────────────┘         └──────────────┘                 │
│                            ┌──────────────┐                │
│                            │  Ambassador  │◄────► Postgres │
│                            │   (PgBouncer)│        Primary  │
│                            └──────────────┘                │
│                                                             │
│  One ambassador per external service type                   │
│  Deep protocol optimization                                 │
└─────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Ambassador Pattern Types

#### 1. Database Connection Ambassador

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Database Ambassador (PgBouncer)                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Application Perspective:                                               │
│  ┌──────────────┐                                                       │
│  │   Service    │  "Connect to localhost:6432"                          │
│  │              │  (Thinks it's direct connection)                       │
│  └──────┬───────┘                                                       │
│         │                                                               │
│         │ localhost:6432                                                │
│         ▼                                                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Ambassador (PgBouncer)                                          │   │
│  │                                                                  │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  Connection Pool                                            │ │   │
│  │  │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐            │ │   │
│  │  │  │ Conn │ │ Conn │ │ Conn │ │ Conn │ │ Conn │ (max 100)    │ │   │
│  │  │  │  1   │ │  2   │ │  3   │ │  4   │ │  5   │              │ │   │
│  │  │  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘            │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                                                                  │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  Features                                                   │ │   │
│  │  │  • Transaction pooling (mode=transaction)                   │ │   │
│  │  │  • Session pooling (mode=session)                         │ │   │
│  │  │  • Statement pooling (mode=statement)                     │ │   │
│  │  │  • Query timeout enforcement                                │ │   │
│  │  │  • Client authentication                                    │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│         │                                                               │
│         │ Actual Database Connections                                   │
│         ▼                                                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    PostgreSQL Cluster                            │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │   Primary    │◄───────►│   Replica    │                     │   │
│  │  │  (Write)     │  Stream │   (Read)     │                     │   │
│  │  │  :5432       │  Repl   │   :5432      │                     │   │
│  │  └──────────────┘         └──────────────┘                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2. Cache Ambassador

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      Cache Ambassador (Twemproxy)                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Application: "SET user:123 {data}"                                     │
│       │                                                                 │
│       │ localhost:22121 (Memcached protocol)                            │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Twemproxy (Nutcracker)                                          │   │
│  │                                                                  │   │
│  │  Consistent Hashing Ring:                                        │   │
│  │  ┌───────────────────────────────────────────────────────────┐   │   │
│  │  │  Key: "user:123"                                          │   │   │
│  │  │  Hash: 0x8A3F...                                          │   │   │
│  │  │  Maps to: redis-node-2                                    │   │   │
│  │  └───────────────────────────────────────────────────────────┘   │   │
│  │                                                                  │   │
│  │  Features:                                                       │   │
│  │  • Protocol parsing (Redis/Memcached)                          │   │
│  │  • Request routing via consistent hashing                      │   │
│  │  • Connection pooling per backend                              │   │
│  │  • Sharding without client complexity                          │   │
│  │  • Failure handling (eject bad nodes)                          │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                 │
│       ▼                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────┐   │
│  │ redis-node-1 │  │ redis-node-2 │  │ redis-node-3 │  │ ...      │   │
│  │ :6379        │  │ :6379        │  │ :6379        │  │          │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 3. Message Queue Ambassador

```
┌─────────────────────────────────────────────────────────────────────────┐
│                   Message Queue Ambassador Pattern                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Application publishes to localhost                                     │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Ambassador (Custom/Envoy with Kafka protocol support)           │   │
│  │                                                                  │   │
│  │  Responsibilities:                                               │   │
│  │  ┌───────────────────────────────────────────────────────────┐   │   │
│  │  │  • Buffer messages during broker unavailability           │   │   │
│  │  │  • Batch small messages for efficiency                    │   │   │
│  │  │  • Handle producer retries with backoff                   │   │   │
│  │  │  • Schema validation (Avro/Protobuf)                      │   │   │
│  │  │  • Partition assignment strategy                          │   │   │
│  │  │  • Idempotency handling                                   │   │   │
│  │  └───────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Kafka Cluster                                 │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │   │
│  │  │ Broker 1 │  │ Broker 2 │  │ Broker 3 │  │ Broker N │        │   │
│  │  │ :9092    │  │ :9092    │  │ :9092    │  │ :9092    │        │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Redis Ambassador Implementation

```go
// cmd/ambassador/redis/main.go
package main

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "hash/crc32"
    "net"
    "sort"
    "sync"
    "sync/atomic"
    "time"

    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
)

// Config holds ambassador configuration
type Config struct {
    ListenAddr       string
    BackendServers   []string
    PoolSize         int
    MaxRetries       int
    RetryBaseDelay   time.Duration
    HealthCheckInterval time.Duration
    Timeout          time.Duration
}

// RedisAmbassador acts as a proxy between applications and Redis cluster
type RedisAmbassador struct {
    config     *Config
    logger     *zap.Logger
    
    // Consistent hashing ring
    ring       *HashRing
    
    // Connection pools per backend
    pools      map[string]*redis.Client
    poolsMu    sync.RWMutex
    
    // Health status
    healthy    map[string]*atomic.Bool
    
    // Statistics
    stats      *Stats
    
    // Listener
    listener   net.Listener
    ctx        context.Context
    cancel     context.CancelFunc
}

// HashRing implements consistent hashing
type HashRing struct {
    replicas int
    ring     map[uint32]string
    nodes    map[string]struct{}
    sortedKeys []uint32
    mu       sync.RWMutex
}

func NewHashRing(replicas int) *HashRing {
    return &HashRing{
        replicas: replicas,
        ring:     make(map[uint32]string),
        nodes:    make(map[string]struct{}),
    }
}

func (hr *HashRing) AddNode(node string) {
    hr.mu.Lock()
    defer hr.mu.Unlock()
    
    if _, exists := hr.nodes[node]; exists {
        return
    }
    
    hr.nodes[node] = struct{}{}
    
    // Add virtual nodes
    for i := 0; i < hr.replicas; i++ {
        hash := hr.hash(fmt.Sprintf("%s:%d", node, i))
        hr.ring[hash] = node
    }
    
    hr.updateSortedKeys()
}

func (hr *HashRing) RemoveNode(node string) {
    hr.mu.Lock()
    defer hr.mu.Unlock()
    
    if _, exists := hr.nodes[node]; !exists {
        return
    }
    
    delete(hr.nodes, node)
    
    // Remove virtual nodes
    for i := 0; i < hr.replicas; i++ {
        hash := hr.hash(fmt.Sprintf("%s:%d", node, i))
        delete(hr.ring, hash)
    }
    
    hr.updateSortedKeys()
}

func (hr *HashRing) GetNode(key string) string {
    hr.mu.RLock()
    defer hr.mu.RUnlock()
    
    if len(hr.ring) == 0 {
        return ""
    }
    
    hash := hr.hash(key)
    
    // Binary search for first node >= hash
    idx := sort.Search(len(hr.sortedKeys), func(i int) bool {
        return hr.sortedKeys[i] >= hash
    })
    
    if idx == len(hr.sortedKeys) {
        idx = 0
    }
    
    return hr.ring[hr.sortedKeys[idx]]
}

func (hr *HashRing) hash(key string) uint32 {
    return crc32.ChecksumIEEE([]byte(key))
}

func (hr *HashRing) updateSortedKeys() {
    hr.sortedKeys = make([]uint32, 0, len(hr.ring))
    for k := range hr.ring {
        hr.sortedKeys = append(hr.sortedKeys, k)
    }
    sort.Slice(hr.sortedKeys, func(i, j int) bool {
        return hr.sortedKeys[i] < hr.sortedKeys[j]
    })
}

// Stats tracks ambassador statistics
type Stats struct {
    TotalCommands   uint64
    ErrorCommands   uint64
    CacheHits       uint64
    CacheMisses     uint64
    BackendFailures uint64
}

func NewRedisAmbassador(cfg *Config, logger *zap.Logger) (*RedisAmbassador, error) {
    ctx, cancel := context.WithCancel(context.Background())
    
    ra := &RedisAmbassador{
        config:  cfg,
        logger:  logger,
        ring:    NewHashRing(150), // 150 virtual nodes per physical node
        pools:   make(map[string]*redis.Client),
        healthy: make(map[string]*atomic.Bool),
        stats:   &Stats{},
        ctx:     ctx,
        cancel:  cancel,
    }
    
    // Initialize connections to backends
    for _, server := range cfg.BackendServers {
        ra.AddBackend(server)
    }
    
    // Start health checks
    go ra.healthCheckLoop()
    
    return ra, nil
}

func (ra *RedisAmbassador) AddBackend(server string) {
    ra.poolsMu.Lock()
    defer ra.poolsMu.Unlock()
    
    // Create Redis client
    client := redis.NewClient(&redis.Options{
        Addr:         server,
        PoolSize:     ra.config.PoolSize,
        MinIdleConns: ra.config.PoolSize / 4,
        MaxRetries:   ra.config.MaxRetries,
        DialTimeout:  ra.config.Timeout,
        ReadTimeout:  ra.config.Timeout,
        WriteTimeout: ra.config.Timeout,
    })
    
    ra.pools[server] = client
    ra.healthy[server] = &atomic.Bool{}
    ra.healthy[server].Store(true)
    
    // Add to hash ring
    ra.ring.AddNode(server)
    
    ra.logger.Info("added backend", zap.String("server", server))
}

func (ra *RedisAmbassador) RemoveBackend(server string) {
    ra.poolsMu.Lock()
    defer ra.poolsMu.Unlock()
    
    if client, ok := ra.pools[server]; ok {
        client.Close()
        delete(ra.pools, server)
        delete(ra.healthy, server)
        ra.ring.RemoveNode(server)
    }
}

func (ra *RedisAmbassador) healthCheckLoop() {
    ticker := time.NewTicker(ra.config.HealthCheckInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            ra.checkBackends()
        case <-ra.ctx.Done():
            return
        }
    }
}

func (ra *RedisAmbassador) checkBackends() {
    ra.poolsMu.RLock()
    servers := make([]string, 0, len(ra.pools))
    for server := range ra.pools {
        servers = append(servers, server)
    }
    ra.poolsMu.RUnlock()
    
    for _, server := range servers {
        go func(s string) {
            ra.poolsMu.RLock()
            client := ra.pools[s]
            ra.poolsMu.RUnlock()
            
            ctx, cancel := context.WithTimeout(ra.ctx, 5*time.Second)
            defer cancel()
            
            err := client.Ping(ctx).Err()
            wasHealthy := ra.healthy[s].Load()
            isHealthy := err == nil
            
            if wasHealthy != isHealthy {
                ra.healthy[s].Store(isHealthy)
                if isHealthy {
                    ra.logger.Info("backend became healthy", zap.String("server", s))
                    ra.ring.AddNode(s)
                } else {
                    ra.logger.Warn("backend became unhealthy", 
                        zap.String("server", s),
                        zap.Error(err))
                    ra.ring.RemoveNode(s)
                }
            }
        }(server)
    }
}

func (ra *RedisAmbassador) getClient(key string) (*redis.Client, bool) {
    // Get node from hash ring
    node := ra.ring.GetNode(key)
    if node == "" {
        return nil, false
    }
    
    ra.poolsMu.RLock()
    defer ra.poolsMu.RUnlock()
    
    client, ok := ra.pools[node]
    return client, ok
}

func (ra *RedisAmbassador) getHealthyClient() (*redis.Client, bool) {
    ra.poolsMu.RLock()
    defer ra.poolsMu.RUnlock()
    
    for server, client := range ra.pools {
        if ra.healthy[server].Load() {
            return client, true
        }
    }
    
    return nil, false
}

// ExecuteCommand executes a Redis command on the appropriate backend
func (ra *RedisAmbassador) ExecuteCommand(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
    atomic.AddUint64(&ra.stats.TotalCommands, 1)
    
    // Determine which backend to use based on command and key
    var client *redis.Client
    var err error
    
    switch cmd {
    case "GET", "SET", "DEL", "EXPIRE", "TTL":
        // Key-based commands - use hash ring
        if len(args) > 0 {
            key, ok := args[0].(string)
            if ok {
                client, _ = ra.getClient(key)
            }
        }
        
    case "MGET", "MSET":
        // Multi-key commands - complex routing or error
        return nil, fmt.Errorf("multi-key commands not supported, use pipeline")
        
    case "PING", "INFO", "CONFIG":
        // Admin commands - any healthy backend
        client, _ = ra.getHealthyClient()
        
    default:
        // Default: try hash ring with first arg as key
        if len(args) > 0 {
            if key, ok := args[0].(string); ok {
                client, _ = ra.getClient(key)
            }
        }
    }
    
    if client == nil {
        atomic.AddUint64(&ra.stats.ErrorCommands, 1)
        return nil, fmt.Errorf("no available backend")
    }
    
    // Execute command
    result, err := executeRedisCommand(ctx, client, cmd, args...)
    if err != nil {
        atomic.AddUint64(&ra.stats.ErrorCommands, 1)
        return nil, err
    }
    
    return result, nil
}

func executeRedisCommand(ctx context.Context, client *redis.Client, cmd string, args ...interface{}) (interface{}, error) {
    // Convert to go-redis command
    switch cmd {
    case "GET":
        if len(args) < 1 {
            return nil, fmt.Errorf("GET requires key")
        }
        return client.Get(ctx, args[0].(string)).Result()
        
    case "SET":
        if len(args) < 2 {
            return nil, fmt.Errorf("SET requires key and value")
        }
        return client.Set(ctx, args[0].(string), args[1], 0).Result()
        
    case "DEL":
        if len(args) < 1 {
            return nil, fmt.Errorf("DEL requires key(s)")
        }
        keys := make([]string, len(args))
        for i, arg := range args {
            keys[i] = arg.(string)
        }
        return client.Del(ctx, keys...).Result()
        
    case "EXPIRE":
        if len(args) < 2 {
            return nil, fmt.Errorf("EXPIRE requires key and seconds")
        }
        seconds := args[1].(time.Duration)
        return client.Expire(ctx, args[0].(string), seconds).Result()
        
    case "PING":
        return client.Ping(ctx).Result()
        
    default:
        // Generic command
        return client.Do(ctx, append([]interface{}{cmd}, args...)...).Result()
    }
}

// Run starts the ambassador server
func (ra *RedisAmbassador) Run() error {
    listener, err := net.Listen("tcp", ra.config.ListenAddr)
    if err != nil {
        return fmt.Errorf("failed to listen: %w", err)
    }
    ra.listener = listener
    
    ra.logger.Info("ambassador listening",
        zap.String("addr", ra.config.ListenAddr))
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            if ra.ctx.Err() != nil {
                return nil
            }
            ra.logger.Error("accept error", zap.Error(err))
            continue
        }
        
        go ra.handleConnection(conn)
    }
}

func (ra *RedisAmbassador) handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // Handle Redis protocol
    // This is a simplified version - real implementation would use
    // a proper Redis protocol parser
    
    buf := make([]byte, 4096)
    for {
        n, err := conn.Read(buf)
        if err != nil {
            return
        }
        
        // Parse command (simplified)
        cmd, args, err := parseRedisCommand(buf[:n])
        if err != nil {
            conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
            continue
        }
        
        // Execute
        result, err := ra.ExecuteCommand(ra.ctx, cmd, args...)
        if err != nil {
            conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
            continue
        }
        
        // Format response
        response := formatRedisResponse(result)
        conn.Write([]byte(response))
    }
}

func (ra *RedisAmbassador) Shutdown() error {
    ra.cancel()
    
    if ra.listener != nil {
        ra.listener.Close()
    }
    
    ra.poolsMu.Lock()
    defer ra.poolsMu.Unlock()
    
    for _, client := range ra.pools {
        client.Close()
    }
    
    return nil
}

// Simplified protocol parsing - real implementation would be more robust
func parseRedisCommand(data []byte) (string, []interface{}, error) {
    // Simplified parsing - assume "COMMAND arg1 arg2 ..." format
    parts := strings.Fields(string(data))
    if len(parts) == 0 {
        return "", nil, fmt.Errorf("empty command")
    }
    
    cmd := strings.ToUpper(parts[0])
    args := make([]interface{}, len(parts)-1)
    for i, part := range parts[1:] {
        args[i] = part
    }
    
    return cmd, args, nil
}

func formatRedisResponse(result interface{}) string {
    switch v := result.(type) {
    case string:
        return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
    case int64:
        return fmt.Sprintf(":%d\r\n", v)
    case nil:
        return "$-1\r\n"
    default:
        return fmt.Sprintf("+%v\r\n", v)
    }
}

func main() {
    logger, _ := zap.NewProduction()
    
    cfg := &Config{
        ListenAddr:          ":6379",
        BackendServers:      []string{"redis-1:6379", "redis-2:6379", "redis-3:6379"},
        PoolSize:            100,
        MaxRetries:          3,
        RetryBaseDelay:      100 * time.Millisecond,
        HealthCheckInterval: 10 * time.Second,
        Timeout:             5 * time.Second,
    }
    
    ambassador, err := NewRedisAmbassador(cfg, logger)
    if err != nil {
        logger.Fatal("failed to create ambassador", zap.Error(err))
    }
    
    if err := ambassador.Run(); err != nil {
        logger.Fatal("ambassador error", zap.Error(err))
    }
}
```

## Trade-off Analysis

### Ambassador Deployment Options

| Option | Latency | Resource Overhead | Complexity | Use Case |
|--------|---------|-------------------|------------|----------|
| **Per-Pod Ambassador** | ~0.1ms | Medium (1 per pod) | Low | Production standard |
| **Per-Node Ambassador** | ~0.5ms | Low (1 per node) | Medium | Resource-constrained |
| **Shared Cluster Service** | 1-5ms | Very Low | High | Development only |
| **Embedded Library** | 0ms | None | High | Specialized cases |

### When to Use Ambassador vs Sidecar

```
Use AMBASSADOR when:
┌─────────────────────────────────────────────────────────────────────────┐
│  • Connecting to stateful external services (DB, cache, queue)          │
│  • Need deep protocol knowledge (SQL parsing, Redis commands)           │
│  • Connection pooling is critical                                       │
│  • Service discovery for external resources                             │
│  • Different teams own app vs infrastructure                            │
└─────────────────────────────────────────────────────────────────────────┘

Use SIDECAR when:
┌─────────────────────────────────────────────────────────────────────────┐
│  • Need uniform traffic management (all inbound/outbound)               │
│  • mTLS for all service-to-service communication                        │
│  • Observability across all traffic types                               │
│  • Rate limiting at service edge                                        │
│  • Same team owns app and proxy configuration                           │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Ambassador Testing

```go
// test/ambassador/integration_test.go
func TestAmbassadorSharding(t *testing.T) {
    // Setup test Redis instances
    redis1 := startTestRedis(t)
    redis2 := startTestRedis(t)
    
    cfg := &Config{
        ListenAddr:     ":16379",
        BackendServers: []string{redis1.Addr(), redis2.Addr()},
        PoolSize:       10,
    }
    
    ambassador, _ := NewRedisAmbassador(cfg, zap.NewNop())
    go ambassador.Run()
    defer ambassador.Shutdown()
    
    time.Sleep(100 * time.Millisecond)
    
    // Test that keys are distributed
    client := redis.NewClient(&redis.Options{Addr: "localhost:16379"})
    
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("key:%d", i)
        client.Set(context.Background(), key, "value", 0)
    }
    
    // Check distribution across backends
    count1 := redis1.DBSize(context.Background()).Val()
    count2 := redis2.DBSize(context.Background()).Val()
    
    // Both should have keys (distribution check)
    assert.Greater(t, count1, int64(0))
    assert.Greater(t, count2, int64(0))
    
    // Total should be 1000
    assert.Equal(t, int64(1000), count1+count2)
}

func TestAmbassadorFailover(t *testing.T) {
    redis1 := startTestRedis(t)
    redis2 := startTestRedis(t)
    
    cfg := &Config{
        ListenAddr:          ":26379",
        BackendServers:      []string{redis1.Addr(), redis2.Addr()},
        HealthCheckInterval: 100 * time.Millisecond,
    }
    
    ambassador, _ := NewRedisAmbassador(cfg, zap.NewNop())
    go ambassador.Run()
    defer ambassador.Shutdown()
    
    time.Sleep(200 * time.Millisecond)
    
    client := redis.NewClient(&redis.Options{Addr: "localhost:26379"})
    
    // Write a key
    client.Set(context.Background(), "testkey", "value", 0)
    
    // Kill first Redis
    redis1.Shutdown()
    
    // Wait for health check
    time.Sleep(200 * time.Millisecond)
    
    // Read should still work (from second Redis or failover)
    val, err := client.Get(context.Background(), "testkey").Result()
    
    // If consistent hashing puts key on dead node, we get miss
    // If it fails over, we get value
    // Both are valid behaviors depending on requirements
    assert.True(t, err == nil || err == redis.Nil)
}
```

## Summary

The Ambassador Pattern provides:

1. **Simplified Client Configuration**: Apps connect to localhost
2. **Service Discovery**: Automatic backend detection
3. **Resilience**: Health checks and failover
4. **Protocol Optimization**: Connection pooling, batching
5. **Observability**: Centralized metrics for external calls

Key considerations:
- Adds network hop (minimal with localhost)
- Additional component to manage
- Protocol-specific implementation needed
- Consistent hashing complexity for sharding
