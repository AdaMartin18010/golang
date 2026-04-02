# 负载均衡 (Load Balancing)

> **分类**: 开源技术堆栈
> **标签**: #loadbalancer #proxy #microservices

---

## 客户端负载均衡

### 轮询算法

```go
type RoundRobin struct {
    backends []string
    current  uint32
}

func (r *RoundRobin) Next() string {
    idx := atomic.AddUint32(&r.current, 1)
    return r.backends[idx%uint32(len(r.backends))]
}
```

### 带健康检查

```go
type HealthAwareLB struct {
    backends map[string]*Backend
    healthy  []*Backend
    mu       sync.RWMutex
}

type Backend struct {
    URL     string
    Healthy bool
    client  *http.Client
}

func (lb *HealthAwareLB) CheckHealth() {
    for _, backend := range lb.backends {
        resp, err := backend.client.Get(backend.URL + "/health")
        backend.Healthy = err == nil && resp.StatusCode == 200
    }
    lb.updateHealthyList()
}

func (lb *HealthAwareLB) updateHealthyList() {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    lb.healthy = lb.healthy[:0]
    for _, b := range lb.backends {
        if b.Healthy {
            lb.healthy = append(lb.healthy, b)
        }
    }
}
```

---

## 反向代理

### 基于 httputil

```go
import "net/http/httputil"

func NewReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
    director := func(req *http.Request) {
        target := targets[rand.Intn(len(targets))]
        req.URL.Scheme = target.Scheme
        req.URL.Host = target.Host
        req.URL.Path = target.Path + req.URL.Path
    }

    return &httputil.ReverseProxy{
        Director: director,
        ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
            log.Printf("proxy error: %v", err)
            w.WriteHeader(http.StatusBadGateway)
        },
    }
}
```

---

## gRPC 负载均衡

```go
// 使用 resolver
conn, err := grpc.Dial(
    "service-name",
    grpc.WithDefaultServiceConfig(`{
        "loadBalancingConfig": [{"round_robin": {}}]
    }`),
)
```

---

## 熔断器

```go
import "github.com/sony/gobreaker"

var cb *gobreaker.CircuitBreaker

func init() {
    settings := gobreaker.Settings{
        MaxRequests: 3,
        Interval:    10 * time.Second,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            return counts.Requests >= 3 && failureRatio >= 0.6
        },
        OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
            log.Printf("circuit breaker %s: %s -> %s", name, from, to)
        },
    }

    cb = gobreaker.NewCircuitBreaker(settings)
}

func CallWithCircuitBreaker(ctx context.Context) (interface{}, error) {
    return cb.Execute(func() (interface{}, error) {
        return makeRequest(ctx)
    })
}
```
