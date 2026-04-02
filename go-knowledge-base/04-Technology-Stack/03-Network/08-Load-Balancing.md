# TS-NET-008: Load Balancing Strategies

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #load-balancing #ha-proxy #nginx #round-robin #least-connections
> **权威来源**:
>
> - [Load Balancing Algorithms](https://www.nginx.com/resources/glossary/load-balancing/) - NGINX
> - [HAProxy Documentation](http://cbonte.github.io/haproxy-dconv/) - HAProxy

---

## 1. Load Balancer Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Load Balancer Architecture                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        Clients                                       │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐                │   │
│  │  │ Client 1│  │ Client 2│  │ Client 3│  │ Client N│                │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘                │   │
│  │       └─────────────┴─────────────┴─────────────┘                   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Load Balancer (L4/L7)                            │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                    Algorithm Selection                         │  │   │
│  │  │  - Round Robin                                               │  │   │
│  │  │  - Least Connections                                         │  │   │
│  │  │  - IP Hash                                                   │  │   │
│  │  │  - Weighted Round Robin                                      │  │   │
│  │  │  - Least Response Time                                       │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                   Health Checking                              │  │   │
│  │  │  - TCP check                                                   │  │   │
│  │  │  - HTTP check                                                  │  │   │
│  │  │  - Custom check                                                │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └───────────────────────────────┬─────────────────────────────────────┘   │
│                                  │                                           │
│         ┌────────────────────────┼────────────────────────┐                 │
│         │                        │                        │                 │
│         ▼                        ▼                        ▼                 │
│  ┌─────────────┐          ┌─────────────┐          ┌─────────────┐        │
│  │  Backend 1  │          │  Backend 2  │          │  Backend N  │        │
│  │  (Active)   │          │  (Active)   │          │  (Active)   │        │
│  │  Weight: 5  │          │  Weight: 3  │          │  Weight: 2  │        │
│  └─────────────┘          └─────────────┘          └─────────────┘        │
│                                                                              │
│  Session Persistence:                                                        │
│  - Sticky sessions (cookie-based)                                           │
│  - IP hashing                                                               │
│  - Shared session store (Redis)                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Load Balancing Algorithms

```go
package loadbalancer

import (
    "hash/fnv"
    "sync"
    "sync/atomic"
)

// Backend represents a backend server
type Backend struct {
    Address     string
    Weight      int
    Connections int64
    Healthy     bool
}

// LoadBalancer interface
type LoadBalancer interface {
    SelectBackends() *Backend
    HealthCheck()
}

// RoundRobin load balancer
type RoundRobin struct {
    backends []*Backend
    current  uint64
}

func NewRoundRobin(backends []*Backend) *RoundRobin {
    return &RoundRobin{backends: backends}
}

func (rr *RoundRobin) SelectBackend() *Backend {
    healthy := rr.getHealthy()
    if len(healthy) == 0 {
        return nil
    }

    next := atomic.AddUint64(&rr.current, 1) % uint64(len(healthy))
    return healthy[next]
}

func (rr *RoundRobin) getHealthy() []*Backend {
    var healthy []*Backend
    for _, b := range rr.backends {
        if b.Healthy {
            healthy = append(healthy, b)
        }
    }
    return healthy
}

// LeastConnections load balancer
type LeastConnections struct {
    backends []*Backend
    mu       sync.RWMutex
}

func NewLeastConnections(backends []*Backend) *LeastConnections {
    return &LeastConnections{backends: backends}
}

func (lc *LeastConnections) SelectBackend() *Backend {
    lc.mu.RLock()
    defer lc.mu.RUnlock()

    var selected *Backend
    var minConn int64 = -1

    for _, b := range lc.backends {
        if !b.Healthy {
            continue
        }

        if minConn == -1 || b.Connections < minConn {
            minConn = b.Connections
            selected = b
        }
    }

    return selected
}

func (lc *LeastConnections) IncrementConnections(backend *Backend) {
    atomic.AddInt64(&backend.Connections, 1)
}

func (lc *LeastConnections) DecrementConnections(backend *Backend) {
    atomic.AddInt64(&backend.Connections, -1)
}

// IPHash load balancer
type IPHash struct {
    backends []*Backend
}

func NewIPHash(backends []*Backend) *IPHash {
    return &IPHash{backends: backends}
}

func (ih *IPHash) SelectBackend(clientIP string) *Backend {
    healthy := ih.getHealthy()
    if len(healthy) == 0 {
        return nil
    }

    h := fnv.New32a()
    h.Write([]byte(clientIP))
    index := h.Sum32() % uint32(len(healthy))

    return healthy[index]
}

func (ih *IPHash) getHealthy() []*Backend {
    var healthy []*Backend
    for _, b := range ih.backends {
        if b.Healthy {
            healthy = append(healthy, b)
        }
    }
    return healthy
}

// WeightedRoundRobin load balancer
type WeightedRoundRobin struct {
    backends []*Backend
    weights  []int
    current  int
    cw       int
    mu       sync.Mutex
}

func NewWeightedRoundRobin(backends []*Backend) *WeightedRoundRobin {
    wrr := &WeightedRoundRobin{
        backends: backends,
        weights:  make([]int, len(backends)),
    }

    gcd := 0
    for _, b := range backends {
        gcd = greatestCommonDivisor(gcd, b.Weight)
    }

    for i, b := range backends {
        wrr.weights[i] = b.Weight / gcd
    }

    return wrr
}

func (wrr *WeightedRoundRobin) SelectBackend() *Backend {
    wrr.mu.Lock()
    defer wrr.mu.Unlock()

    healthy := wrr.getHealthy()
    if len(healthy) == 0 {
        return nil
    }

    for {
        wrr.current = (wrr.current + 1) % len(healthy)
        if wrr.current == 0 {
            wrr.cw--
            if wrr.cw <= 0 {
                wrr.cw = wrr.maxWeight()
            }
        }

        if wrr.weights[wrr.current] >= wrr.cw {
            return healthy[wrr.current]
        }
    }
}

func (wrr *WeightedRoundRobin) maxWeight() int {
    max := 0
    for _, w := range wrr.weights {
        if w > max {
            max = w
        }
    }
    return max
}

func (wrr *WeightedRoundRobin) getHealthy() []*Backend {
    var healthy []*Backend
    for _, b := range wrr.backends {
        if b.Healthy {
            healthy = append(healthy, b)
        }
    }
    return healthy
}

func greatestCommonDivisor(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}
```

---

## 3. Health Checking

```go
// Health checker
type HealthChecker struct {
    backends  []*Backend
    interval  time.Duration
    timeout   time.Duration
    checkFunc func(*Backend) bool
}

func NewHealthChecker(backends []*Backend, interval, timeout time.Duration) *HealthChecker {
    return &HealthChecker{
        backends: backends,
        interval: interval,
        timeout:  timeout,
    }
}

func (hc *HealthChecker) Start() {
    ticker := time.NewTicker(hc.interval)
    go func() {
        for range ticker.C {
            hc.checkAll()
        }
    }()
}

func (hc *HealthChecker) checkAll() {
    for _, backend := range hc.backends {
        go func(b *Backend) {
            healthy := hc.checkBackend(b)
            b.Healthy = healthy
        }(backend)
    }
}

func (hc *HealthChecker) checkBackend(backend *Backend) bool {
    client := &http.Client{
        Timeout: hc.timeout,
    }

    resp, err := client.Get("http://" + backend.Address + "/health")
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode == http.StatusOK
}
```

---

## 4. Checklist

```
Load Balancing Checklist:
□ Algorithm chosen appropriately
□ Health checks configured
□ Sticky sessions if needed
□ SSL termination configured
□ Backend weights set
□ Monitoring for backend health
□ Failover handling
□ Graceful backend shutdown
```
