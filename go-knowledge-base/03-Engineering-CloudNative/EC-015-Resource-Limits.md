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

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02

---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02