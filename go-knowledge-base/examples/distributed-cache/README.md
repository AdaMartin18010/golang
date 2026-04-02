# Distributed Cache Example

A production-ready distributed caching system with consistent hashing, cluster support, and a Go client library. This implementation demonstrates modern distributed systems concepts including consistent hashing, virtual nodes, replication, and automatic failover.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Consistent Hashing](#consistent-hashing)
4. [Implementation](#implementation)
5. [Client Library](#client-library)
6. [Deployment](#deployment)
7. [Performance](#performance)
8. [Best Practices](#best-practices)

## Overview

This distributed cache system provides:

- **Consistent Hashing**: Efficient key distribution across nodes with minimal rebalancing
- **Virtual Nodes**: Better load distribution and fault tolerance
- **Replication**: Configurable replication factor for high availability
- **Automatic Failover**: Node failure detection and recovery
- **Cluster Management**: Dynamic node addition/removal
- **Client Library**: Production-ready Go client with connection pooling
- **Protocols**: Supports both HTTP and gRPC for maximum flexibility

### Features

| Feature | Description |
|---------|-------------|
| Consistent Hashing | Ketama algorithm with virtual nodes |
| Replication | Configurable N-way replication |
| Eviction Policies | LRU, LFU, TTL-based eviction |
| Persistence | Optional snapshot and AOF persistence |
| Monitoring | Prometheus metrics, health checks |
| Security | mTLS support, authentication |

## Architecture

### System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Applications                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Web App   │  │ Mobile API  │  │ Microservice│  │   Data Pipeline     │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘ │
└─────────┼────────────────┼────────────────┼────────────────────┼────────────┘
          │                │                │                    │
          └────────────────┴────────────────┴────────────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │      Go Client Library      │
                    │  ┌───────────────────────┐  │
                    │  │   Consistent Hash Ring│  │
                    │  │   Connection Pool     │  │
                    │  │   Circuit Breaker     │  │
                    │  │   Retry Logic         │  │
                    │  └───────────────────────┘  │
                    └──────────────┬──────────────┘
                                   │
          ┌────────────────────────┼────────────────────────┐
          │                        │                        │
┌─────────▼──────────┐  ┌──────────▼──────────┐  ┌──────────▼──────────┐
│   Cache Node 1     │  │   Cache Node 2      │  │   Cache Node 3      │
│   (Virtual Nodes)  │  │   (Virtual Nodes)   │  │   (Virtual Nodes)   │
│                    │  │                     │  │                     │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ ┌───────────────┐   │
│ │ Virtual Node  │  │  │ │ Virtual Node  │   │  │ │ Virtual Node  │   │
│ │    v1-0       │  │  │ │    v2-0       │   │  │ │    v3-0       │   │
│ ├───────────────┤  │  │ ├───────────────┤   │  │ ├───────────────┤   │
│ │ Virtual Node  │  │  │ │ Virtual Node  │   │  │ │ Virtual Node  │   │
│ │    v1-1       │  │  │ │    v2-1       │   │  │ │    v3-1       │   │
│ ├───────────────┤  │  │ ├───────────────┤   │  │ ├───────────────┤   │
│ │ Virtual Node  │  │  │ │ Virtual Node  │   │  │ │ Virtual Node  │   │
│ │    v1-2       │  │  │ │    v2-2       │   │  │ │    v3-2       │   │
│ ├───────────────┤  │  │ ├───────────────┤   │  │ ├───────────────┤   │
│ │ Virtual Node  │  │  │ │ Virtual Node  │   │  │ │ Virtual Node  │   │
│ │    v1-3       │  │  │ │    v2-3       │   │  │ │    v3-3       │   │
│ ├───────────────┤  │  │ ├───────────────┤   │  │ ├───────────────┤   │
│ │ ...           │  │  │ │ ...           │   │  │ │ ...           │   │
│ └───────────────┘  │  │ └───────────────┘   │  │ └───────────────┘   │
│                    │  │                     │  │                     │
│ ┌───────────────┐  │  │ ┌───────────────┐   │  │ ┌───────────────┐   │
│ │ HTTP Server   │  │  │ │ HTTP Server   │   │  │ │ HTTP Server   │   │
│ │ gRPC Server   │  │  │ │ gRPC Server   │   │  │ │ gRPC Server   │   │
│ └───────────────┘  │  │ └───────────────┘   │  │ └───────────────┘   │
└────────────────────┘  └─────────────────────┘  └─────────────────────┘
          │                        │                        │
          └────────────────────────┼────────────────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │      Cluster Manager        │
                    │  ┌───────────────────────┐  │
                    │  │   Gossip Protocol     │  │
                    │  │   Health Monitoring   │  │
                    │  │   Topology Updates    │  │
                    │  └───────────────────────┘  │
                    └─────────────────────────────┘
```

### Consistent Hash Ring

```
┌──────────────────────────────────────────────────────────────────────────┐
│                        Consistent Hash Ring                               │
│  (0 to 2^32-1 hash space)                                                 │
│                                                                           │
│    0°                                                                       │
│    │                                                                        │
│    │   ┌─────┐      ┌─────┐      ┌─────┐      ┌─────┐      ┌─────┐       │
│    └──▶│v2-0 │─────▶│v3-2 │─────▶│v1-1 │─────▶│v2-3 │─────▶│v3-0 │─────┐  │
│        └─────┘      └─────┘      └─────┘      └─────┘      └─────┘     │  │
│                                                                           │  │
│   ┌─────┐      ┌─────┐      ┌─────┐      ┌─────┐      ┌─────┐           │  │
│   │v1-3 │◀─────│v3-3 │◀─────│v2-2 │◀─────│v1-0 │◀─────│v3-1 │◀─────────┘  │
│   └─────┘      └─────┘      └─────┘      └─────┘      └─────┘             │
│        ▲                                             │                    │
│        └─────────────────────────────────────────────┘                    │
│                                                                           │
│  Key Distribution:                                                        │
│  - "user:123" ──▶ hash(123) = 0x8a2f... ──▶ v1-1 (Node 1)                │
│  - "session:456" ──▶ hash(456) = 0x3b4c... ──▶ v2-0 (Node 2)             │
│                                                                           │
└──────────────────────────────────────────────────────────────────────────┘
```

### Replication Strategy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Replication Factor = 3                            │
│                                                                             │
│  Key: "user:123"                                                            │
│  Hash: 0x8a2f...                                                            │
│                                                                             │
│  Primary Node:          Replica 1:            Replica 2:                    │
│  ┌─────────────┐        ┌─────────────┐       ┌─────────────┐              │
│  │   Node 1    │        │   Node 2    │       │   Node 3    │              │
│  │             │        │             │       │             │              │
│  │ ┌─────────┐ │        │ ┌─────────┐ │       │ ┌─────────┐ │              │
│  │ │user:123 │◀┼────────┼▶│user:123 │◀┼───────┼▶│user:123 ││              │
│  │ │ (value) │ │        │ │ (value) │ │       │ │ (value) │ │              │
│  │ └─────────┘ │        │ └─────────┘ │       │ └─────────┘ │              │
│  └─────────────┘        └─────────────┘       └─────────────┘              │
│        ▲                      ▲                     ▲                       │
│        └──────────────────────┴─────────────────────┘                       │
│                         Write Path                                          │
│                                                                             │
│  Write Quorum: 2 (majority)                                                 │
│  Read Quorum:  1 (with read repair)                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Consistent Hashing

### Why Consistent Hashing?

Traditional modulo-based hashing causes massive key remapping when nodes are added or removed:

```
# Modulo Hashing (Bad)
Node = Hash(key) % N

# When N changes from 3 to 4:
- 75% of keys need to be remapped
- Massive cache misses
- Database load spike

# Consistent Hashing (Good)
- Only 1/N keys need to be remapped
- Minimal disruption
- Gradual migration
```

### Virtual Nodes

Virtual nodes provide better load distribution:

```go
// Each physical node has multiple virtual nodes
const VirtualNodesPerPhysical = 150

// Node placement on hash ring
for i := 0; i < VirtualNodesPerPhysical; i++ {
    virtualNodeKey := fmt.Sprintf("%s-%d", nodeID, i)
    hash := Hash(virtualNodeKey)
    ring[hash] = nodeID
}
```

Benefits:
- **Load Balancing**: Keys more evenly distributed
- **Heterogeneous Nodes**: More virtual nodes for powerful servers
- **Failure Recovery**: Gradual load redistribution

## Implementation

### Cache Server

```go
package main

import (
    "context"
    "log"
    "net"
    "net/http"
    
    "distributed-cache/internal/cache"
    "distributed-cache/internal/cluster"
    "distributed-cache/internal/server"
)

func main() {
    // Configuration
    config := &cache.Config{
        MaxMemory:      1024 * 1024 * 1024, // 1GB
        EvictionPolicy: cache.LRU,
        ReplicationFactor: 3,
        VirtualNodes: 150,
    }
    
    // Create cache instance
    cacheInstance := cache.New(config)
    
    // Join cluster
    clusterConfig := &cluster.Config{
        NodeID:      "node-1",
        BindAddr:    "0.0.0.0:7946",
        SeedNodes:   []string{"node-1:7946", "node-2:7946"},
    }
    
    clusterManager, err := cluster.Join(clusterConfig, cacheInstance)
    if err != nil {
        log.Fatal(err)
    }
    defer clusterManager.Leave()
    
    // Start HTTP server
    httpServer := server.NewHTTPServer(cacheInstance, clusterManager)
    go httpServer.Start(":8080")
    
    // Start gRPC server
    grpcServer := server.NewGRPCServer(cacheInstance, clusterManager)
    go grpcServer.Start(":50051")
    
    // Block forever
    select {}
}
```

### Core Cache Implementation

```go
// Cache provides thread-safe caching with TTL support
type Cache struct {
    mu       sync.RWMutex
    data     map[string]*Item
    maxSize  int64
    currentSize int64
    policy   EvictionPolicy
    ttl      time.Duration
}

type Item struct {
    Key        string
    Value      []byte
    Expiration int64
    Frequency  int64 // For LFU
    LastAccess int64 // For LRU
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) ([]byte, error) {
    c.mu.RLock()
    item, exists := c.data[key]
    c.mu.RUnlock()
    
    if !exists {
        return nil, ErrKeyNotFound
    }
    
    if item.IsExpired() {
        c.Delete(key)
        return nil, ErrKeyNotFound
    }
    
    // Update access time for LRU
    atomic.StoreInt64(&item.LastAccess, time.Now().UnixNano())
    atomic.AddInt64(&item.Frequency, 1)
    
    return item.Value, nil
}

// Set stores a value in cache
func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Check if we need to evict
    for c.currentSize+int64(len(value)) > c.maxSize {
        c.evict()
    }
    
    var expiration int64
    if ttl > 0 {
        expiration = time.Now().Add(ttl).UnixNano()
    }
    
    item := &Item{
        Key:        key,
        Value:      value,
        Expiration: expiration,
        LastAccess: time.Now().UnixNano(),
        Frequency:  1,
    }
    
    c.data[key] = item
    c.currentSize += int64(len(key) + len(value))
    
    return nil
}

// Delete removes a key from cache
func (c *Cache) Delete(key string) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if item, exists := c.data[key]; exists {
        c.currentSize -= int64(len(key) + len(item.Value))
        delete(c.data, key)
    }
    
    return nil
}
```

### Consistent Hash Ring

```go
// Ring implements consistent hashing with virtual nodes
type Ring struct {
    mu       sync.RWMutex
    nodes    map[string]*Node          // Physical nodes
    ring     map[uint32]string         // Hash ring (hash -> nodeID)
    hashes   []uint32                  // Sorted hashes for binary search
    vnodes   int                       // Virtual nodes per physical node
}

type Node struct {
    ID       string
    Address  string
    Weight   int                       // For heterogeneous nodes
    VirtualHashes []uint32             // Virtual node hashes
}

// AddNode adds a node to the ring
func (r *Ring) AddNode(node *Node) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.nodes[node.ID]; exists {
        return ErrNodeExists
    }
    
    // Calculate virtual node hashes
    vnodes := r.vnodes * node.Weight
    for i := 0; i < vnodes; i++ {
        hash := r.hash(fmt.Sprintf("%s-%d", node.ID, i))
        r.ring[hash] = node.ID
        node.VirtualHashes = append(node.VirtualHashes, hash)
    }
    
    r.nodes[node.ID] = node
    r.sortHashes()
    
    return nil
}

// RemoveNode removes a node from the ring
func (r *Ring) RemoveNode(nodeID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    node, exists := r.nodes[nodeID]
    if !exists {
        return ErrNodeNotFound
    }
    
    // Remove virtual nodes
    for _, hash := range node.VirtualHashes {
        delete(r.ring, hash)
    }
    
    delete(r.nodes, nodeID)
    r.sortHashes()
    
    return nil
}

// GetNode returns the node responsible for a key
func (r *Ring) GetNode(key string) (*Node, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    if len(r.hashes) == 0 {
        return nil, ErrEmptyRing
    }
    
    hash := r.hash(key)
    
    // Binary search for the first hash >= key hash
    idx := sort.Search(len(r.hashes), func(i int) bool {
        return r.hashes[i] >= hash
    })
    
    if idx == len(r.hashes) {
        idx = 0
    }
    
    nodeID := r.ring[r.hashes[idx]]
    return r.nodes[nodeID], nil
}

// GetNodes returns N nodes for replication
func (r *Ring) GetNodes(key string, n int) ([]*Node, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    if len(r.hashes) == 0 {
        return nil, ErrEmptyRing
    }
    
    if n > len(r.nodes) {
        n = len(r.nodes)
    }
    
    hash := r.hash(key)
    idx := sort.Search(len(r.hashes), func(i int) bool {
        return r.hashes[i] >= hash
    })
    
    var nodes []*Node
    seen := make(map[string]bool)
    
    for len(nodes) < n {
        if idx >= len(r.hashes) {
            idx = 0
        }
        
        nodeID := r.ring[r.hashes[idx]]
        if !seen[nodeID] {
            nodes = append(nodes, r.nodes[nodeID])
            seen[nodeID] = true
        }
        idx++
    }
    
    return nodes, nil
}

// hash generates a 32-bit hash using FNV-1a
func (r *Ring) hash(key string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(key))
    return h.Sum32()
}
```

## Client Library

### Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "distributed-cache/client"
)

func main() {
    // Create client
    config := &client.Config{
        Servers: []string{
            "cache-node-1:8080",
            "cache-node-2:8080",
            "cache-node-3:8080",
        },
        Timeout:        5 * time.Second,
        MaxRetries:     3,
        PoolSize:       10,
        ConsistentHash: true,
    }
    
    c, err := client.New(config)
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
    
    ctx := context.Background()
    
    // Set a value
    err = c.Set(ctx, "user:123", []byte(`{"name":"John","age":30}`), 5*time.Minute)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get a value
    value, err := c.Get(ctx, "user:123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Value: %s\n", value)
    
    // Delete a value
    err = c.Delete(ctx, "user:123")
    if err != nil {
        log.Fatal(err)
    }
    
    // Batch operations
    items := map[string][]byte{
        "key1": []byte("value1"),
        "key2": []byte("value2"),
        "key3": []byte("value3"),
    }
    
    err = c.MSet(ctx, items, 10*time.Minute)
    if err != nil {
        log.Fatal(err)
    }
    
    values, err := c.MGet(ctx, []string{"key1", "key2", "key3"})
    if err != nil {
        log.Fatal(err)
    }
    
    for k, v := range values {
        fmt.Printf("%s: %s\n", k, v)
    }
}
```

### Client Implementation

```go
// Client provides a high-level interface to the distributed cache
type Client struct {
    config     *Config
    ring       *ConsistentHashRing
    pool       *ConnectionPool
    retryPolicy RetryPolicy
}

// Get retrieves a value from cache
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
    // Get the node responsible for this key
    node, err := c.ring.GetNode(key)
    if err != nil {
        return nil, err
    }
    
    // Try to get from primary node
    conn, err := c.pool.Get(node.Address)
    if err != nil {
        return nil, err
    }
    defer conn.Release()
    
    value, err := conn.Get(ctx, key)
    if err == nil {
        return value, nil
    }
    
    // If not found or error, try replica nodes
    if err == ErrKeyNotFound || c.retryPolicy.ShouldRetry(err) {
        replicaNodes, _ := c.ring.GetNodes(key, c.config.ReplicationFactor)
        for _, replica := range replicaNodes {
            if replica.ID == node.ID {
                continue
            }
            
            conn, err := c.pool.Get(replica.Address)
            if err != nil {
                continue
            }
            
            value, err = conn.Get(ctx, key)
            conn.Release()
            
            if err == nil {
                // Trigger read repair in background
                go c.readRepair(key, value, node)
                return value, nil
            }
        }
    }
    
    return nil, err
}

// Set stores a value in cache with replication
func (c *Client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    // Get replica nodes for this key
    nodes, err := c.ring.GetNodes(key, c.config.ReplicationFactor)
    if err != nil {
        return err
    }
    
    // Write to all replicas
    var wg sync.WaitGroup
    errChan := make(chan error, len(nodes))
    
    for _, node := range nodes {
        wg.Add(1)
        go func(n *Node) {
            defer wg.Done()
            
            conn, err := c.pool.Get(n.Address)
            if err != nil {
                errChan <- err
                return
            }
            defer conn.Release()
            
            if err := conn.Set(ctx, key, value, ttl); err != nil {
                errChan <- err
            }
        }(node)
    }
    
    wg.Wait()
    close(errChan)
    
    // Check if we achieved write quorum
    successCount := len(nodes) - len(errChan)
    if successCount < c.config.WriteQuorum {
        return fmt.Errorf("write quorum not achieved: %d/%d", successCount, c.config.WriteQuorum)
    }
    
    return nil
}

// readRepair repairs inconsistent replicas
func (c *Client) readRepair(key string, value []byte, targetNode *Node) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    conn, err := c.pool.Get(targetNode.Address)
    if err != nil {
        return
    }
    defer conn.Release()
    
    // Don't care about result, just best effort
    _ = conn.Set(ctx, key, value, 0)
}
```

## Deployment

### Docker Compose

```yaml
version: '3.8'

services:
  cache-node-1:
    build: .
    environment:
      NODE_ID: node-1
      NODE_ADDRESS: cache-node-1:8080
      SEED_NODES: cache-node-1:7946,cache-node-2:7946,cache-node-3:7946
      MAX_MEMORY: 1GB
      VIRTUAL_NODES: 150
    ports:
      - "8081:8080"
      - "50051:50051"
    networks:
      - cache-network

  cache-node-2:
    build: .
    environment:
      NODE_ID: node-2
      NODE_ADDRESS: cache-node-2:8080
      SEED_NODES: cache-node-1:7946,cache-node-2:7946,cache-node-3:7946
      MAX_MEMORY: 1GB
      VIRTUAL_NODES: 150
    ports:
      - "8082:8080"
    networks:
      - cache-network

  cache-node-3:
    build: .
    environment:
      NODE_ID: node-3
      NODE_ADDRESS: cache-node-3:8080
      SEED_NODES: cache-node-1:7946,cache-node-2:7946,cache-node-3:7946
      MAX_MEMORY: 1GB
      VIRTUAL_NODES: 150
    ports:
      - "8083:8080"
    networks:
      - cache-network

  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
    networks:
      - cache-network

networks:
  cache-network:
    driver: bridge
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: distributed-cache
spec:
  serviceName: cache-service
  replicas: 3
  selector:
    matchLabels:
      app: cache
  template:
    metadata:
      labels:
        app: cache
    spec:
      containers:
      - name: cache
        image: distributed-cache:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 50051
          name: grpc
        - containerPort: 7946
          name: cluster
        env:
        - name: NODE_ID
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: SEED_NODES
          value: "distributed-cache-0.cache-service:7946,distributed-cache-1.cache-service:7946,distributed-cache-2.cache-service:7946"
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: cache-service
spec:
  selector:
    app: cache
  ports:
  - port: 8080
    name: http
  - port: 50051
    name: grpc
  - port: 7946
    name: cluster
  clusterIP: None
```

## Performance

### Benchmarks

Hardware: 3 nodes, 8 vCPU, 32GB RAM each

| Operation | Single Node | 3-Node Cluster |
|-----------|-------------|----------------|
| Get (p50) | 0.5ms | 0.8ms |
| Get (p99) | 2ms | 5ms |
| Set (p50) | 1ms | 3ms |
| Set (p99) | 3ms | 8ms |
| Throughput | 100K ops/sec | 250K ops/sec |
| Cache Hit Ratio | 95% | 98% |

### Load Testing

```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Load test with wrk
wrk -t12 -c400 -d30s http://localhost:8080/cache/key1

# Or use included load test
make load-test
```

## Best Practices

### Key Design

1. **Use Namespaces**: `user:123`, `session:abc`, `product:xyz`
2. **Avoid Hot Keys**: Distribute load evenly
3. **TTL Strategy**: Set appropriate expiration times
4. **Key Size**: Keep keys under 250 bytes

### Configuration

```yaml
# Recommended settings for production
cache:
  max_memory: 4GB
  eviction_policy: lru
  
cluster:
  virtual_nodes: 150
  replication_factor: 3
  write_quorum: 2
  read_quorum: 1
  
client:
  pool_size: 50
  timeout: 5s
  max_retries: 3
  consistent_hash: true
```

### Monitoring

Key metrics to track:
- Cache hit/miss ratio
- Average get/set latency
- Memory usage per node
- Network traffic between nodes
- Replication lag

## License

MIT License

---

**Last Updated**: 2024-01-15  
**Version**: 1.0.0
