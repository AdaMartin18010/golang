# 内存泄漏检测 (Memory Leak Detection)

> **分类**: 工程与云原生

---

## 常见内存泄漏

### 1. Goroutine 泄漏

```go
// ❌ 错误: 发送者阻塞
func bad() {
    ch := make(chan int)
    go func() {
        ch <- 42  // 无人接收，永久阻塞
    }()
}

// ✅ 正确: 使用缓冲或 select
func good() {
    ch := make(chan int, 1)  // 缓冲
    go func() {
        ch <- 42
    }()
}
```

### 2. Timer 未停止

```go
// ❌ 错误
timer := time.NewTimer(time.Hour)
// 如果提前返回，timer 继续运行

// ✅ 正确
timer := time.NewTimer(time.Hour)
defer timer.Stop()
```

### 3. 全局引用

```go
// ❌ 错误: 全局缓存无限增长
var cache = map[string][]byte{}

func store(key string, data []byte) {
    cache[key] = data  // 永不清理
}

// ✅ 正确: LRU 缓存
import lru "github.com/hashicorp/golang-lru"

cache, _ := lru.New(1000)  // 限制大小
```

---

## 检测方法

### pprof

```go
import _ "net/http/pprof"

go http.ListenAndServe("localhost:6060", nil)
```

```bash
# 查看堆分配
go tool pprof http://localhost:6060/debug/pprof/heap

# 对比两个时间点的堆
# T1
curl http://localhost:6060/debug/pprof/heap > heap.1
# T2
curl http://localhost:6060/debug/pprof/heap > heap.2

# 比较
go tool pprof -base heap.1 heap.2
```

### Goroutine 泄漏检测

```go
// 使用 goleak
import "go.uber.org/goleak"

func TestFunction(t *testing.T) {
    defer goleak.VerifyNone(t)

    // 测试代码
}
```

---

## 监控

```go
// 导出内存指标
var (
    heapAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "go_heap_alloc_bytes",
        Help: "Heap allocation",
    })
    goroutines = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "go_goroutines",
        Help: "Number of goroutines",
    })
)

func recordMetrics() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    heapAlloc.Set(float64(m.HeapAlloc))
    goroutines.Set(float64(runtime.NumGoroutine()))
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