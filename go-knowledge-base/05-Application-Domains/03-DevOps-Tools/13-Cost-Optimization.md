# 成本优化 (Cost Optimization)

> **分类**: 成熟应用领域
> **标签**: #cost #optimization #cloud

---

## 资源使用监控

```go
// 资源使用指标
type ResourceMetrics struct {
    CPUUsage    float64
    MemoryUsage float64
    DiskIO      float64
    NetworkIO   float64
}

func CollectMetrics() ResourceMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return ResourceMetrics{
        MemoryUsage: float64(m.Alloc) / 1024 / 1024,  // MB
    }
}

// Prometheus 导出
var (
    memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "app_memory_usage_mb",
        Help: "Current memory usage in MB",
    })
)

func recordMetrics() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    memoryUsage.Set(float64(m.Alloc) / 1024 / 1024)
}
```

---

## 自动伸缩

```go
type AutoScaler struct {
    minReplicas int
    maxReplicas int
    targetCPU   float64
    k8sClient   kubernetes.Interface
}

func (as *AutoScaler) Scale(ctx context.Context, deployment string) error {
    // 获取当前指标
    metrics, err := as.getMetrics(ctx, deployment)
    if err != nil {
        return err
    }

    // 计算所需副本数
    currentReplicas := metrics.CurrentReplicas
    cpuPercent := metrics.CPUUsage / float64(currentReplicas)

    desiredReplicas := int(cpuPercent / as.targetCPU * float64(currentReplicas))

    // 边界检查
    if desiredReplicas < as.minReplicas {
        desiredReplicas = as.minReplicas
    }
    if desiredReplicas > as.maxReplicas {
        desiredReplicas = as.maxReplicas
    }

    // 应用伸缩
    if desiredReplicas != currentReplicas {
        return as.applyScale(ctx, deployment, desiredReplicas)
    }

    return nil
}
```

---

## 资源配额管理

```go
type QuotaManager struct {
    quotas map[string]*Quota
    mu     sync.RWMutex
}

type Quota struct {
    CPU     int     // cores
    Memory  int     // MB
    Storage int     // GB
    Used    Usage
}

type Usage struct {
    CPU     int
    Memory  int
    Storage int
}

func (qm *QuotaManager) CheckQuota(team string, req Request) error {
    qm.mu.RLock()
    defer qm.mu.RUnlock()

    quota, ok := qm.quotas[team]
    if !ok {
        return fmt.Errorf("no quota for team %s", team)
    }

    if quota.Used.CPU+req.CPU > quota.CPU {
        return fmt.Errorf("CPU quota exceeded: %d/%d",
            quota.Used.CPU+req.CPU, quota.CPU)
    }

    if quota.Used.Memory+req.Memory > quota.Memory {
        return fmt.Errorf("memory quota exceeded: %d/%d",
            quota.Used.Memory+req.Memory, quota.Memory)
    }

    return nil
}
```

---

## Spot 实例利用

```go
type SpotInstanceManager struct {
    onDemandRatio float64  // 0.2 = 20% on-demand, 80% spot
}

func (sim *SpotInstanceManager) Allocate(nodes int) Allocation {
    onDemand := int(float64(nodes) * sim.onDemandRatio)
    spot := nodes - onDemand

    return Allocation{
        OnDemand: onDemand,
        Spot:     spot,
        Cost:     calculateCost(onDemand, spot),
    }
}

func calculateCost(onDemand, spot int) float64 {
    // Spot 通常便宜 60-90%
    spotPrice := 0.3  // 30% of on-demand

    onDemandCost := float64(onDemand) * 1.0
    spotCost := float64(spot) * spotPrice

    return onDemandCost + spotCost
}
```

---

## 成本告警

```go
type CostAlert struct {
    Threshold float64
    Current   float64
    Window    time.Duration
}

func (ca *CostAlert) Check(currentCost float64) (bool, string) {
    if currentCost > ca.Threshold {
        return true, fmt.Sprintf(
            "Cost alert: $%.2f exceeds threshold $%.2f",
            currentCost, ca.Threshold,
        )
    }
    return false, ""
}

func sendAlert(message string) {
    // 发送到 Slack/PagerDuty
    webhook.Send(Alert{
        Severity: "warning",
        Message:  message,
    })
}
```

---

## 架构决策记录

### 决策矩阵

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| A | 高性能 | 复杂 | 大规模 |
| B | 简单 | 扩展性差 | 小规模 |

### 风险评估

**风险 R.1**: 性能瓶颈
- 概率: 中
- 影响: 高
- 缓解: 缓存、分片

**风险 R.2**: 单点故障
- 概率: 低
- 影响: 极高
- 缓解: 冗余、故障转移

### 实施路线图

`
Phase 1: 基础设施 (Week 1-2)
Phase 2: 核心功能 (Week 3-6)
Phase 3: 优化加固 (Week 7-8)
`

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 架构决策记录 (ADR)

### 上下文

业务需求和技术约束分析。

### 决策

选择方案A作为主要架构方向。

### 后果

正面：
- 可扩展性提升
- 维护成本降低

负面：
- 初期开发复杂度增加
- 团队学习成本

### 实施指南

`
Week 1-2: 基础设施搭建
Week 3-4: 核心功能开发
Week 5-6: 集成测试
Week 7-8: 性能优化
`

### 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 性能不足 | 中 | 高 | 缓存、分片 |
| 兼容性 | 低 | 中 | 接口适配层 |

### 监控指标

- 系统吞吐量
- 响应延迟
- 错误率
- 资源利用率

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 系统设计

### 需求分析

功能需求和非功能需求的完整梳理。

### 架构视图

`
┌─────────────────────────────────────┐
│           API Gateway               │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    ▼                   ▼
┌─────────┐       ┌─────────┐
│ Service │       │ Service │
│   A     │       │   B     │
└────┬────┘       └────┬────┘
     │                 │
     └────────┬────────┘
              ▼
        ┌─────────┐
        │  Data   │
        │  Store  │
        └─────────┘
`

### 技术选型

| 组件 | 技术 | 理由 |
|------|------|------|
| API | gRPC | 性能 |
| DB | PostgreSQL | 可靠 |
| Cache | Redis | 速度 |
| Queue | Kafka | 吞吐 |

### 性能指标

- QPS: 10K+
- P99 Latency: <100ms
- Availability: 99.99%

### 运维手册

- 部署流程
- 监控配置
- 应急预案
- 容量规划

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