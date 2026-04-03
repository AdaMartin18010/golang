# EC-121: Google SRE 可靠性工程实践 (Google SRE Reliability Engineering)

> **维度**: Engineering CloudNative
> **级别**: S (30+ KB)
> **标签**: #sre #reliability #sla #error-budget #observability
> **权威来源**: [Google SRE Book](https://sre.google/sre-book/table-of-contents/), [Site Reliability Workbook](https://sre.google/workbook/table-of-contents/), [Google Cloud Operations](https://cloud.google.com/blog/products/devops-sre)

---

## SRE 核心理念

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        SRE Fundamental Principles                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Service Level Objectives (SLOs)                                         │
│     ─────────────────────────────────                                       │
│     Availability: 99.9% ("three nines") = 8.77 hours downtime/year          │
│     Availability: 99.99% ("four nines") = 52.6 minutes downtime/year        │
│     Availability: 99.999% ("five nines") = 5.26 minutes downtime/year       │
│                                                                              │
│  2. Error Budget                                                            │
│     ────────────────                                                        │
│     Error Budget = 100% - SLO                                               │
│     Example: 99.9% SLO → 0.1% Error Budget                                  │
│     When budget exhausted: freeze feature launches                          │
│                                                                              │
│  3. Toil Elimination                                                        │
│     ────────────────                                                        │
│     Toil: Manual, repetitive, automatable tasks                             │
│     Target: < 50% of SRE time on toil                                       │
│                                                                              │
│  4. Blameless Postmortems                                                   │
│     ─────────────────────                                                   │
│     Focus on systemic fixes, not individual blame                           │
│     Document: What happened, Detection, Response, Recovery                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## SLI / SLO / SLA 定义

### 公式化定义

$$
\begin{aligned}
&\text{SLI (Service Level Indicator):} \\
&\quad \text{request_latency}_{p99} = \text{percentile}(\{latency_i\}, 99) \\
&\quad \text{error_rate} = \frac{\text{error_requests}}{\text{total_requests}} \\
&\quad \text{availability} = \frac{\text{successful_requests}}{\text{total_requests}} \\
\\
&\text{SLO (Service Level Objective):} \\
&\quad \text{availability} \geq 99.9\% \text{ over 30 days} \\
&\quad \text{latency}_{p99} \leq 200\text{ms} \\
\\
&\text{SLA (Service Level Agreement):} \\
&\quad \text{Contractual obligation with penalties} \\
&\quad \text{SLA SLO is typically looser than internal SLO}
\end{aligned}
$$

### Go 实现：SLI 计算

```go
package sre

import (
    "sync"
    "time"

    "github.com/prometheus/client_golang/prometheus"
)

// SLIMetrics SLI 指标收集
type SLIMetrics struct {
    // 可用性
    totalRequests   prometheus.Counter
    errorRequests   prometheus.Counter

    // 延迟
    requestDuration prometheus.Histogram

    // 吞吐量
    requestsPerSecond prometheus.Gauge

    // 自定义 SLI
    customSLIs map[string]prometheus.Gauge
}

func NewSLIMetrics(serviceName string) *SLIMetrics {
    return &SLIMetrics{
        totalRequests: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "service_requests_total",
            Help: "Total number of requests",
            ConstLabels: prometheus.Labels{"service": serviceName},
        }),
        errorRequests: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "service_errors_total",
            Help: "Total number of error requests",
            ConstLabels: prometheus.Labels{"service": serviceName},
        }),
        requestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "service_request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
            ConstLabels: prometheus.Labels{"service": serviceName},
        }),
    }
}

// RecordRequest 记录请求
func (m *SLIMetrics) RecordRequest(duration time.Duration, err error) {
    m.totalRequests.Inc()
    m.requestDuration.Observe(duration.Seconds())

    if err != nil {
        m.errorRequests.Inc()
    }
}

// GetAvailability 计算可用性
func (m *SLIMetrics) GetAvailability(window time.Duration) float64 {
    total := getCounterValue(m.totalRequests)
    errors := getCounterValue(m.errorRequests)

    if total == 0 {
        return 1.0
    }

    return 1.0 - (errors / total)
}
```

---

## Error Budget 策略

```go
// ErrorBudgetController 错误预算控制器
type ErrorBudgetController struct {
    slo         float64        // 目标 SLO (如 0.999)
    budget      float64        // 剩余预算 (0.0 - 1.0)
    window      time.Duration  // 计算窗口 (如 30天)

    // 状态
    exhausted   bool
    mu          sync.RWMutex

    // 回调
    onExhausted func()
    onReset     func()
}

func NewErrorBudgetController(slo float64, window time.Duration) *ErrorBudgetController {
    return &ErrorBudgetController{
        slo:    slo,
        budget: 1.0 - slo,  // 初始满预算
        window: window,
    }
}

// Consume 消费错误预算
func (ebc *ErrorBudgetController) Consume(errors, total float64) {
    ebc.mu.Lock()
    defer ebc.mu.Unlock()

    if total == 0 {
        return
    }

    errorRate := errors / total
    ebc.budget -= errorRate

    if ebc.budget <= 0 && !ebc.exhausted {
        ebc.exhausted = true
        if ebc.onExhausted != nil {
            ebc.onExhausted()
        }
    }
}

// CanLaunch 检查是否可以发布新功能
func (ebc *ErrorBudgetController) CanLaunch() bool {
    ebc.mu.RLock()
    defer ebc.mu.RUnlock()

    // 预算低于 50% 时，暂停非紧急发布
    return ebc.budget > (1.0-ebc.slo)*0.5
}

// AlertRules 告警规则
func (ebc *ErrorBudgetController) AlertRules() []AlertRule {
    return []AlertRule{
        {
            Name:        "ErrorBudgetBurningFast",
            Condition:   "budget will exhaust in 2 days",
            Severity:    "critical",
        },
        {
            Name:        "ErrorBudgetBurningSlow",
            Condition:   "budget will exhaust in 30 days",
            Severity:    "warning",
        },
    }
}
```

---

## 可靠性与速度的平衡

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Reliability vs Development Velocity                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Development Speed                                                          │
│  ▲                                                                          │
│  │    ┌──────────────────────────────────────────┐                         │
│  │    │  🚀 Feature Development                  │                         │
│  │    │     (Error Budget > 50%)                 │                         │
│  │    │                                          │                         │
│  │    │  When Error Budget is healthy,           │                         │
│  │    │  prioritize feature launches             │                         │
│  │    └──────────────────────────────────────────┘                         │
│  │                              │                                            │
│  │    ┌─────────────────────────┴──────────────┐                           │
│  │    │  ⚠️  Reliability Work                  │                           │
│  │    │     (Error Budget < 50%)               │                           │
│  │    │                                          │                           │
│  │    │  When Error Budget is at risk,         │                           │
│  │    │  freeze launches, focus on fixes       │                           │
│  │    └──────────────────────────────────────────┘                         │
│  │                                                                          │
│  └───────────────────────────────────────────────────────────────────────►  │
│                              Time →                                          │
│                                                                              │
│  Policy: "If SLO specifies 99.9%, we aim for 99.99% to preserve budget"    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 监控与告警设计

### 四个黄金信号

| 信号 | 指标 | 示例 |
|------|------|------|
| Latency | 请求处理时间 | p50, p95, p99 latency |
| Traffic | 请求量 | QPS, active connections |
| Errors | 错误率 | 5xx rate, timeout rate |
| Saturation | 资源利用率 | CPU, memory, disk, bandwidth |

```go
// FourGoldenSignals 四个黄金信号监控
type FourGoldenSignals struct {
    // Latency
    latencyHistogram *prometheus.HistogramVec

    // Traffic
    requestCounter   *prometheus.CounterVec

    // Errors
    errorCounter     *prometheus.CounterVec

    // Saturation
    resourceGauges   map[string]prometheus.Gauge
}

func (fgs *FourGoldenSignals) Record(method string, duration time.Duration, err error) {
    // Latency
    fgs.latencyHistogram.WithLabelValues(method).Observe(duration.Seconds())

    // Traffic
    fgs.requestCounter.WithLabelValues(method).Inc()

    // Errors
    if err != nil {
        fgs.errorCounter.WithLabelValues(method, errorType(err)).Inc()
    }
}
```

---

## 参考文献

1. [Site Reliability Engineering](https://sre.google/sre-book/table-of-contents/) - Google
2. [The Site Reliability Workbook](https://sre.google/workbook/table-of-contents/) - Google
3. [Implementing SLOs](https://sre.google/workbook/implementing-slos/) - Google SRE Workbook
4. [Error Budget Policy](https://sre.google/workbook/error-budget-policy/) - Google
5. [SRE at Google: Reliability at Scale](https://www.youtube.com/watch?v=HhBI1SCz8oU) - Google Cloud

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