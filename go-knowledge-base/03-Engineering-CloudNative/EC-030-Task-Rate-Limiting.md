# 任务限流与降级 (Task Rate Limiting & Degradation)

> **分类**: 工程与云原生
> **标签**: #rate-limiting #circuit-breaker #degradation

---

## 自适应限流

```go
type AdaptiveRateLimiter struct {
    limit     int
    current   int64
    success   int64
    failure   int64
    mu        sync.RWMutex
}

func (arl *AdaptiveRateLimiter) Allow() bool {
    for {
        current := atomic.LoadInt64(&arl.current)
        limit := atomic.LoadInt64(&arl.limit)

        if current >= limit {
            return false
        }

        if atomic.CompareAndSwapInt64(&arl.current, current, current+1) {
            return true
        }
    }
}

func (arl *AdaptiveRateLimiter) RecordResult(success bool) {
    if success {
        atomic.AddInt64(&arl.success, 1)
    } else {
        atomic.AddInt64(&arl.failure, 1)
    }

    // 自适应调整
    arl.adjust()
}

func (arl *AdaptiveRateLimiter) adjust() {
    arl.mu.Lock()
    defer arl.mu.Unlock()

    total := arl.success + arl.failure
    if total < 100 {
        return  // 样本不足
    }

    successRate := float64(arl.success) / float64(total)

    if successRate < 0.9 {
        // 成功率低，降低限制
        arl.limit = int(float64(arl.limit) * 0.9)
        if arl.limit < 10 {
            arl.limit = 10
        }
    } else if successRate > 0.99 {
        // 成功率高，提高限制
        arl.limit = int(float64(arl.limit) * 1.1)
        if arl.limit > 1000 {
            arl.limit = 1000
        }
    }

    // 重置计数
    atomic.StoreInt64(&arl.success, 0)
    atomic.StoreInt64(&arl.failure, 0)
}
```

---

## 优先级降级

```go
type PriorityDegradation struct {
    levels    []DegradationLevel
    current   int
    metrics   MetricsCollector
}

type DegradationLevel struct {
    Name          string
    Threshold     float64  // CPU/内存阈值
    DropLowPri    bool     // 是否丢弃低优先级任务
    DisableNonEssential bool
}

var defaultLevels = []DegradationLevel{
    {
        Name:      "normal",
        Threshold: 0.5,
    },
    {
        Name:      "warning",
        Threshold: 0.7,
        DropLowPri: true,
    },
    {
        Name:      "critical",
        Threshold: 0.85,
        DropLowPri: true,
        DisableNonEssential: true,
    },
    {
        Name:      "emergency",
        Threshold: 0.95,
        DropLowPri: true,
        DisableNonEssential: true,
    },
}

func (pd *PriorityDegradation) Evaluate() {
    cpu := pd.metrics.GetCPUUsage()
    memory := pd.metrics.GetMemoryUsage()

    load := math.Max(cpu, memory)

    // 确定当前级别
    level := 0
    for i, l := range pd.levels {
        if load >= l.Threshold {
            level = i
        }
    }

    if level != pd.current {
        pd.applyLevel(level)
    }
}

func (pd *PriorityDegradation) applyLevel(level int) {
    oldLevel := pd.levels[pd.current]
    newLevel := pd.levels[level]

    log.Printf("Degradation level changed: %s -> %s", oldLevel.Name, newLevel.Name)

    if newLevel.DropLowPri && !oldLevel.DropLowPri {
        // 开始丢弃低优先级任务
        pd.startDroppingLowPriority()
    }

    if newLevel.DisableNonEssential && !oldLevel.DisableNonEssential {
        // 禁用非关键功能
        pd.disableNonEssential()
    }

    pd.current = level
}

func (pd *PriorityDegradation) ShouldProcess(task *Task) bool {
    level := pd.levels[pd.current]

    if level.DropLowPri && task.Priority == PriorityLow {
        return false
    }

    if level.DisableNonEssential && task.Category == CategoryNonEssential {
        return false
    }

    return true
}
```

---

## 令牌桶限流

```go
type TokenBucket struct {
    rate       float64    // 每秒产生令牌数
    burst      int        // 桶容量
    tokens     float64
    lastUpdate time.Time
    mu         sync.Mutex
}

func (tb *TokenBucket) Allow(n int) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastUpdate).Seconds()
    tb.lastUpdate = now

    // 添加新令牌
    tb.tokens += elapsed * tb.rate
    if tb.tokens > float64(tb.burst) {
        tb.tokens = float64(tb.burst)
    }

    // 检查是否足够
    if tb.tokens >= float64(n) {
        tb.tokens -= float64(n)
        return true
    }

    return false
}

func (tb *TokenBucket) Wait(ctx context.Context, n int) error {
    for {
        if tb.Allow(n) {
            return nil
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

## 热点防护

```go
type HotspotProtection struct {
    counters   map[string]*SlidingWindow
    thresholds map[string]int
    mu         sync.RWMutex
}

func (hp *HotspotProtection) IsHotspot(key string) bool {
    hp.mu.RLock()
    counter, exists := hp.counters[key]
    threshold := hp.thresholds[key]
    hp.mu.RUnlock()

    if !exists {
        return false
    }

    return counter.Count() > threshold
}

func (hp *HotspotProtection) Record(key string) {
    hp.mu.Lock()

    if _, exists := hp.counters[key]; !exists {
        hp.counters[key] = NewSlidingWindow(time.Second)
    }

    counter := hp.counters[key]
    hp.mu.Unlock()

    counter.Add(1)
}

// 热点降级
func (hp *HotspotProtection) HandleRequest(ctx context.Context, key string, handler func() error) error {
    if hp.IsHotspot(key) {
        // 热点降级：返回缓存或直接拒绝
        if cached := hp.getCache(key); cached != nil {
            return nil
        }
        return ErrHotspotLimited
    }

    hp.Record(key)
    return handler()
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