# 资源限制 (Resource Limits)

> **分类**: 工程与云原生
> **标签**: #resources #cgroups #limits

---

## 内存限制

```go
import "runtime"

// 设置内存限制 (Go 1.19+)
func SetMemoryLimit(limit int64) {
    // 设置软限制
    runtime.SetMemoryLimit(limit)
}

// 监控内存使用
func MonitorMemory() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("Alloc = %v MB\n", m.Alloc/1024/1024)
    fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc/1024/1024)
    fmt.Printf("Sys = %v MB\n", m.Sys/1024/1024)
    fmt.Printf("NumGC = %v\n", m.NumGC)
}
```

---

## CPU 限制

### GOMAXPROCS

```go
import "runtime"

// 设置使用的 CPU 核心数
func init() {
    // 自动检测，但容器环境需要手动设置
    // runtime.GOMAXPROCS(runtime.NumCPU())

    // 容器感知
    cpus := cpuset.CountCPUs()
    runtime.GOMAXPROCS(cpus)
}
```

### 自适应限流

```go
type AdaptiveLimiter struct {
    targetCPU float64  // 目标 CPU 使用率
    current   float64
    rate      int64    // 当前速率
}

func (l *AdaptiveLimiter) Limit(ctx context.Context, fn func() error) error {
    for {
        if l.current < l.targetCPU {
            return fn()
        }

        select {
        case <-time.After(10 * time.Millisecond):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

---

## 文件描述符限制

```go
import "syscall"

func GetFileDescriptorLimit() (uint64, uint64, error) {
    var rLimit syscall.Rlimit
    err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
    return rLimit.Cur, rLimit.Max, err
}

func SetFileDescriptorLimit(max uint64) error {
    var rLimit syscall.Rlimit
    rLimit.Cur = max
    rLimit.Max = max
    return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}
```

---

## 速率限制

```go
type RateLimiter struct {
    rate   int           // 每秒请求数
    burst  int           // 突发容量
    tokens chan struct{} // 令牌桶
}

func NewRateLimiter(rate, burst int) *RateLimiter {
    rl := &RateLimiter{
        rate:   rate,
        burst:  burst,
        tokens: make(chan struct{}, burst),
    }

    // 初始化令牌
    for i := 0; i < burst; i++ {
        rl.tokens <- struct{}{}
    }

    // 持续补充令牌
    go rl.refill()

    return rl
}

func (rl *RateLimiter) refill() {
    ticker := time.NewTicker(time.Second / time.Duration(rl.rate))
    defer ticker.Stop()

    for range ticker.C {
        select {
        case rl.tokens <- struct{}{}:
        default: // 桶满
        }
    }
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}
```

---

## Kubernetes 资源

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

---

## 资源监控

```go
func ResourceMetrics() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return map[string]interface{}{
        "memory_alloc_mb":     m.Alloc / 1024 / 1024,
        "memory_sys_mb":       m.Sys / 1024 / 1024,
        "gc_cycles":           m.NumGC,
        "goroutines":          runtime.NumGoroutine(),
        "cpu_cores":           runtime.GOMAXPROCS(0),
        "cgroups_memory":      getCgroupMemoryLimit(),
    }
}
```
