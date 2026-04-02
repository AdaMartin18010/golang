# AD-010: System Design Interview Preparation

> **Dimension**: Application Domains
> **Level**: S (20+ KB)
> **Tags**: #system-design #interview #scalability #reliability #distributed-systems

---

## 1. Interview Framework

### 1.1 RASCAL Framework

| Letter | Component | Time | Key Questions |
|--------|-----------|------|---------------|
| R | Requirements | 5 min | Functional? Non-functional? |
| A | Architecture | 10 min | High-level design? |
| S | Scale | 5 min | Capacity planning? |
| C | Components | 15 min | Database? Cache? Queue? |
| A | Algorithms | 10 min | Sharding? Consistency? |
| L | Logistics | 5 min | Monitoring? Deployment? |

### 1.2 Requirements Gathering

**Functional Requirements**:

- Core features
- API endpoints
- User interactions

**Non-Functional Requirements**:

- Availability: 99.9%, 99.99%, 99.999%
- Latency: P50, P95, P99 targets
- Throughput: QPS/RPS requirements
- Consistency: Strong vs Eventual
- Durability: Data retention

---

## 2. Capacity Planning

### 2.1 Back-of-Envelope Calculations

| Metric | Calculation | Example |
|--------|-------------|---------|
| Daily Active Users | Given | 100M DAU |
| Requests per User | Estimated | 10 req/day |
| Total Daily Requests | DAU × RPU | 1B requests/day |
| Peak QPS | Daily × 2 / 86400 | ~23K QPS |
| Write QPS | Read/Write ratio | 2.3K writes |
| Read QPS | 90% of total | 20.7K reads |
| Storage per day | Write QPS × size × time | 2TB/day |
| 5-year storage | Daily × 365 × 5 × replica | 11PB |

### 2.2 Resource Estimation

```
Server Requirements:
- 1 server: 64GB RAM, 16 cores
- QPS capacity per server: ~1K QPS
- Servers needed: 23 (with 2x redundancy) = 46 servers

Database Requirements:
- Write QPS: 2.3K
- MySQL can handle: ~1K writes/second
- Shards needed: 4-8

Cache Requirements:
- Cache hit ratio target: 90%
- Working set size: 100GB
- Redis memory needed: 128GB (with overhead)
- Redis nodes: 4 (32GB each)
```

---

## 3. System Components

### 3.1 Load Balancer

```go
package design

// Types of Load Balancers
const (
    Layer4 LoadBalancerType = iota // Transport layer
    Layer7                         // Application layer
)

// Algorithms
type LBAlgorithm int

const (
    RoundRobin LBAlgorithm = iota
    LeastConnections
    IPHash
    WeightedRoundRobin
)

// Health Check
type HealthChecker struct {
    interval    time.Duration
    timeout     time.Duration
    threshold   int
    unhealthy   map[string]int
}

func (h *HealthChecker) Check(endpoint string) bool {
    resp, err := http.Get(endpoint + "/health")
    if err != nil || resp.StatusCode != 200 {
        h.unhealthy[endpoint]++
        return h.unhealthy[endpoint] < h.threshold
    }
    h.unhealthy[endpoint] = 0
    return true
}
```

### 3.2 Database Selection

| Use Case | Database | Reason |
|----------|----------|--------|
| Transactions | PostgreSQL | ACID compliance |
| High write throughput | Cassandra | LSM trees |
| Flexible schema | MongoDB | Document model |
| Caching | Redis | In-memory |
| Search | Elasticsearch | Inverted index |
| Graph data | Neo4j | Graph traversal |
| Time series | InfluxDB | Time-optimized |

### 3.3 Caching Strategies

| Strategy | Use Case | Pros | Cons |
|----------|----------|------|------|
| Cache-aside | Read-heavy | Simple | Stale data |
| Write-through | Consistency needed | No stale data | Write latency |
| Write-behind | Write-heavy | Fast writes | Data loss risk |
| Read-through | Complex backend | Transparent | Cache miss penalty |

### 3.4 Message Queues

| Feature | Kafka | RabbitMQ | SQS | Pulsar |
|---------|-------|----------|-----|--------|
| Throughput | Very High | High | Medium | Very High |
| Persistence | Yes | Optional | Yes | Yes |
| Ordering | Partition | Queue | Best effort | Partition |
| Replay | Yes | No | No | Yes |
| Delayed messages | No | Yes | Yes | Yes |

---

## 4. Scalability Patterns

### 4.1 Horizontal Scaling

```
┌─────────────────────────────────────────────────────────────────┐
│                      Horizontal Scaling                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────┐     ┌──────────────┐     ┌─────────────────────┐  │
│  │  CDN    │────>│ Load Balancer │────>│  App Servers        │  │
│  └─────────┘     └──────────────┘     │  ┌─────┐ ┌─────┐   │  │
│                                        │  │ S1  │ │ S2  │   │  │
│                                        │  └─────┘ └─────┘   │  │
│                                        │  ┌─────┐ ┌─────┐   │  │
│                                        │  │ S3  │ │ S4  │   │  │
│                                        │  └─────┘ └─────┘   │  │
│                                        └─────────────────────┘  │
│                                                                  │
│  Stateless Application Servers                                   │
│  - Share nothing                                                 │
│  - Session in Redis/Cookie                                       │
│  - Easy to add/remove instances                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Database Sharding

```go
package design

// Sharding Strategies

type Sharder interface {
    GetShard(key string) int
}

// Hash-based sharding
type HashSharder struct {
    numShards int
}

func (h *HashSharder) GetShard(key string) int {
    hash := fnv32(key)
    return int(hash % uint32(h.numShards))
}

// Range-based sharding
type RangeSharder struct {
    ranges []Range
}

type Range struct {
    Min   int
    Max   int
    Shard int
}

func (r *RangeSharder) GetShard(key string) int {
    id := parseInt(key)
    for _, rng := range r.ranges {
        if id >= rng.Min && id < rng.Max {
            return rng.Shard
        }
    }
    return 0
}

// Directory-based sharding
type DirectorySharder struct {
    lookup map[string]int
}

func (d *DirectorySharder) GetShard(key string) int {
    return d.lookup[key]
}
```

### 4.3 Consistent Hashing

```go
package design

import (
    "hash/crc32"
    "sort"
)

type ConsistentHash struct {
    replicas int
    ring     []uint32
    nodes    map[uint32]string
}

func NewConsistentHash(replicas int) *ConsistentHash {
    return &ConsistentHash{
        replicas: replicas,
        nodes:    make(map[uint32]string),
    }
}

func (ch *ConsistentHash) Add(node string) {
    for i := 0; i < ch.replicas; i++ {
        hash := ch.hash(node + string(rune(i)))
        ch.nodes[hash] = node
        ch.ring = append(ch.ring, hash)
    }
    sort.Slice(ch.ring, func(i, j int) bool {
        return ch.ring[i] < ch.ring[j]
    })
}

func (ch *ConsistentHash) Get(key string) string {
    if len(ch.ring) == 0 {
        return ""
    }

    hash := ch.hash(key)
    idx := sort.Search(len(ch.ring), func(i int) bool {
        return ch.ring[i] >= hash
    })

    if idx == len(ch.ring) {
        idx = 0
    }

    return ch.nodes[ch.ring[idx]]
}

func (ch *ConsistentHash) hash(key string) uint32 {
    return crc32.ChecksumIEEE([]byte(key))
}
```

---

## 5. Reliability Patterns

### 5.1 Replication

| Strategy | Consistency | Availability | Use Case |
|----------|-------------|--------------|----------|
| Single Leader | Strong | Medium | Read-heavy |
| Multi-Leader | Eventual | High | Multi-region |
| Leaderless | Eventual | Very High | High write |

### 5.2 CAP Theorem

```
┌─────────────────────────────────────────────────────────────────┐
│                         CAP Theorem                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│                    ┌──────────────┐                             │
│                   /    Consistency  \                            │
│                  /         (C)        \                           │
│                 /                       \                          │
│                /                         \                         │
│       ┌───────┴───────┐             ┌───────┴───────┐            │
│       │  CP Systems   │             │  AP Systems   │            │
│       │               │             │               │            │
│       │ • PostgreSQL  │             │ • Cassandra   │            │
│       │ • MongoDB     │             │ • DynamoDB    │            │
│       │ • Redis       │             │ • CouchDB     │            │
│       └───────────────┘             └───────────────┘            │
│                                                                  │
│                    \                       /                         │
│                     \    Availability    /                          │
│                      \      (A)         /                           │
│                       \               /                            │
│                        \           /                             │
│                         \       /                              │
│                          \   /                               │
│                           \ /                                │
│                      Partition Tolerance                       │
│                            (P)                                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.3 Circuit Breaker

```go
package design

type CircuitBreaker struct {
    state          State
    failureCount   int
    successCount   int
    failureThreshold int
    successThreshold int
    timeout        time.Duration
}

func (cb *CircuitBreaker) Call(f func() error) error {
    if cb.state == StateOpen {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
            cb.failureCount = 0
            cb.successCount = 0
        } else {
            return ErrCircuitOpen
        }
    }

    err := f()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) recordResult(err error) {
    if err != nil {
        cb.failureCount++
        if cb.failureCount >= cb.failureThreshold {
            cb.state = StateOpen
            cb.lastFailure = time.Now()
        }
    } else {
        cb.successCount++
        if cb.state == StateHalfOpen && cb.successCount >= cb.successThreshold {
            cb.state = StateClosed
        }
    }
}
```

---

## 6. Common System Designs

### 6.1 URL Shortener

```
Requirements:
- Shorten long URLs
- Redirect to original URL
- Custom aliases (optional)
- Analytics (optional)

Scale:
- 100M new URLs/month
- 10:1 read:write ratio
- 100B redirects/month

Design:
┌─────────┐     ┌──────────┐     ┌──────────┐
│  Client │────>│  Load    │────>│  App     │
└─────────┘     │ Balancer │     │ Servers  │
                └──────────┘     └────┬─────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
                    v                 v                 v
              ┌──────────┐     ┌──────────┐     ┌──────────┐
              │  Cache   │     │  DB      │     │ Key Gen  │
              │ (Redis)  │     │(MySQL)   │     │ (ZooKeeper│
              └──────────┘     └──────────┘     └──────────┘
```

### 6.2 News Feed System

```
Requirements:
- Post text/image/video
- View news feed
- Follow/unfollow users

Push vs Pull:
- Normal user (1000 followers): Push model
- Celebrity (1M followers): Pull model

Design:
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   User      │────>│  Post       │────>│  Fanout     │
│   Action    │     │  Service    │     │  Service    │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                         ┌─────────────────────┼──────┐
                         │                     │      │
                         v                     v      v
                   ┌──────────┐          ┌────────┐ ┌────────┐
                   │  Redis   │          │User A  │ │User B  │
                   │ Timeline │          │Feed    │ │Feed    │
                   └──────────┘          └────────┘ └────────┘
```

---

## 7. Interview Tips

### 7.1 Common Mistakes

| Mistake | Solution |
|---------|----------|
| Jumping to solutions | Ask clarifying questions first |
| Ignoring non-functional | Always discuss scale |
| Single point of failure | Design for redundancy |
| Premature optimization | Start simple, then scale |
| Not discussing trade-offs | Compare alternatives |

### 7.2 Time Management

| Phase | Duration | Activities |
|-------|----------|------------|
| Understanding | 2-3 min | Clarify requirements |
| Estimation | 5 min | Back-of-envelope math |
| High-level | 10-15 min | Draw architecture |
| Deep dive | 15-20 min | Data model, algorithms |
| Trade-offs | 5-10 min | Discuss alternatives |

---

## References

1. Designing Data-Intensive Applications - Martin Kleppmann
2. System Design Interview - Alex Xu
3. Designing Distributed Systems - Brendan Burns
4. Scalability Rules - Abbott & Fisher

---

**Quality Rating**: S (20+ KB)
**Last Updated**: 2026-04-02
