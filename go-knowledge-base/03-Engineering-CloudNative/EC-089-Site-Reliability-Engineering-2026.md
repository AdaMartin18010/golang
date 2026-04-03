# EC-089-Site-Reliability-Engineering-2026

> **Dimension**: 03-Engineering-CloudNative  
> **Status**: S-Level  
> **Created**: 2026-04-03  
> **Version**: SRE 2026 (SLI/SLO/SLA, Error Budgets, Chaos Engineering)  
> **Size**: >20KB 

---

## 1. SRE核心概念

### 1.1 SRE vs DevOps

| 维度 | SRE | DevOps |
|------|-----|--------|
| 起源 | Google 2003 | 社区运动 2009 |
| 重点 | 可靠性工程 | 协作文化 |
| 度量 | SLI/SLO/SLA | 部署频率、MTTR |
| 团队 | 专职SRE团队 | 跨职能团队 |
| 工具 | 标准化平台 | 多样化选择 |

### 1.2 错误预算哲学

```
┌─────────────────────────────────────────┐
│         Error Budget Philosophy         │
├─────────────────────────────────────────┤
│                                         │
│  100%可靠性是错误的目标                  │
│                                         │
│  原因:                                  │
│  1. 成本指数增长                         │
│  2. 抑制创新                             │
│  3. 用户无法感知差异(99.9% vs 99.99%)   │
│                                         │
│  目标可靠性 = 1 - 错误预算               │
│                                         │
└─────────────────────────────────────────┘
```

---

## 2. SLI/SLO/SLA

### 2.1 定义

| 术语 | 定义 | 示例 |
|------|------|------|
| **SLI** | 服务水平指标 | 延迟、可用性、错误率 |
| **SLO** | 服务水平目标 | 99.9%可用性 |
| **SLA** | 服务水平协议 | 未达标赔偿条款 |

### 2.2 SLI选择

```go
// 常见SLI类型
type SLIType string

const (
    Availability SLIType = "availability"
    Latency      SLIType = "latency"
    ErrorRate    SLIType = "error_rate"
    Throughput   SLIType = "throughput"
    Saturation   SLIType = "saturation"
    Durability   SLIType = "durability"
)

// SLI定义
type SLI struct {
    Name        string
    Type        SLIType
    Description string
    Query       string  // PromQL/LogQL
    GoodEvents  string  // 良好事件计数
    TotalEvents string  // 总事件计数
}

// 示例: HTTP可用性
var HTTPAvailabilitySLI = SLI{
    Name:        "http_availability",
    Type:        Availability,
    Description: "HTTP请求成功率",
    Query: `sum(rate(http_requests_total{status!~"5.."}[5m])) 
            / 
            sum(rate(http_requests_total[5m]))`,
}

// 示例: P99延迟
var P99LatencySLI = SLI{
    Name:        "p99_latency",
    Type:        Latency,
    Description: "P99响应延迟",
    Query: `histogram_quantile(0.99, 
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le))`,
}
```

### 2.3 SLO定义

```yaml
# slo.yaml
apiVersion: openslo/v1
kind: SLO
metadata:
  name: payment-service-availability
spec:
  service: payment-service
  description: 支付服务可用性目标
  
  indicator:
    metricSource:
      type: prometheus
      spec:
        query: |
          sum(rate(http_requests_total{service="payment",status!~"5.."}[5m]))
          /
          sum(rate(http_requests_total{service="payment"}[5m]))
  
  timeWindow:
    - duration: 30d
      isRolling: true
  
  budgetingMethod: Occurrences
  
  objectives:
    - target: 0.999  # 99.9%
      displayName: 可用性目标
  
  alertPolicies:
    - name: fast-burn
      condition:
        kind: Burnrate
        spec:
          burnRateCondition:
            lookbackWindow: 1h
            alertAfter: 2%
```

### 2.4 错误预算计算

```go
// 错误预算计算
func CalculateErrorBudget(slo float64, window time.Duration) (budget float64) {
    // 错误预算 = (1 - SLO) * 时间窗口内的总事件
    
    // 例如: 99.9% SLO, 30天窗口
    // 错误预算 = 0.1% * 30天的请求数
    
    errorRate := 1 - slo  // 0.001 for 99.9%
    
    // 假设1000 RPS
    rps := 1000.0
    totalRequests := rps * window.Seconds()
    
    errorBudget := errorRate * totalRequests
    
    return errorBudget  // 30天约259,200个可接受错误
}

// 错误预算消耗速度
func BurnRate(alerts int, window time.Duration) float64 {
    // 消耗速度 = 实际错误率 / 错误预算率
    
    actualErrorRate := float64(alerts) / window.Hours()
    budgetErrorRate := CalculateErrorBudget(0.999, window) / window.Hours()
    
    return actualErrorRate / budgetErrorRate
}

// 多窗口警报
const (
    FastBurnRate   = 14.4  // 2%预算在1小时内消耗
    SlowBurnRate   = 2     // 5%预算在6小时内消耗
    VerySlowBurnRate = 1   // 10%预算在3天内消耗
)
```

---

## 3. 监控和可观测性

### 3.1 黄金信号

```go
// 四个黄金信号
type GoldenSignals struct {
    Latency   float64  // 服务处理请求时间
    Traffic   float64  // 请求量
    Errors    float64  // 错误率
    Saturation float64 // 资源利用率
}

// 实现示例
func RecordGoldenSignals(ctx context.Context, handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // 包装ResponseWriter以捕获状态码
        rw := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        handler.ServeHTTP(rw, r)
        
        duration := time.Since(start)
        
        // 记录指标
        metrics.HTTPDuration.Observe(duration.Seconds())
        metrics.HTTPRequestsTotal.WithLabelValues(
            r.Method, 
            strconv.Itoa(rw.statusCode),
        ).Inc()
        
        if rw.statusCode >= 500 {
            metrics.HTTPErrorsTotal.Inc()
        }
    })
}
```

### 3.2 RED方法

| 指标 | 描述 | PromQL |
|------|------|--------|
| Rate | 每秒请求数 | `rate(http_requests_total[5m])` |
| Errors | 错误率 | `rate(http_requests_total{status=~"5.."}[5m])` |
| Duration | 请求持续时间 | `histogram_quantile(0.99, ...)` |

### 3.3 USE方法

| 指标 | 描述 | 适用资源 |
|------|------|---------|
| Utilization | 资源利用率 | CPU, 内存, 磁盘 |
| Saturation | 饱和度 | 队列长度, 等待时间 |
| Errors | 错误计数 | 硬件错误, 网络丢包 |

---

## 4. 混沌工程

### 4.1 原则

```
1. 建立稳态假设
2. 引入现实世界的混乱
3. 在稳态中引入混乱
4. 否定假设
```

### 4.2 实验类型

```yaml
# 网络延迟
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-delay
spec:
  action: delay
  mode: one
  selector:
    labelSelectors:
      app: payment-service
  delay:
    latency: "100ms"
    correlation: "100"
    jitter: "0ms"
  duration: "5m"

---
# Pod故障
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: pod-kill
spec:
  action: pod-kill
  mode: one
  selector:
    labelSelectors:
      app: order-service
  duration: "30s"

---
# CPU压力
apiVersion: chaos-mesh.org/v1alpha1
kind: StressChaos
metadata:
  name: cpu-stress
spec:
  mode: all
  selector:
    labelSelectors:
      app: api-gateway
  stressors:
    cpu:
      workers: 4
      load: 80
  duration: "10m"
```

### 4.3 自动化解复

```go
// 自动检测和缓解
type ChaosDetection struct {
    metricsClient prometheus.Client
    alertManager  *alertmanager.Client
}

func (cd *ChaosDetection) DetectAnomaly() (*Anomaly, error) {
    // 检测异常模式
    query := `(
        sum(rate(http_requests_total{status=~"5.."}[5m]))
        /
        sum(rate(http_requests_total[5m]))
    ) > 0.01`  // 错误率 > 1%
    
    result, err := cd.metricsClient.Query(query)
    if err != nil {
        return nil, err
    }
    
    if result > threshold {
        return &Anomaly{
            Type:      HighErrorRate,
            Severity:  Critical,
            Timestamp: time.Now(),
        }, nil
    }
    
    return nil, nil
}

func (cd *ChaosDetection) Mitigate(anomaly *Anomaly) error {
    switch anomaly.Type {
    case HighErrorRate:
        // 自动降级
        return cd.triggerCircuitBreaker()
    case HighLatency:
        // 自动扩容
        return cd.scaleUp()
    case ResourceExhaustion:
        // 流量限制
        return cd.enableRateLimit()
    }
    return nil
}
```

---

## 5. 容量规划

### 5.1 容量公式

```
所需容量 = 峰值流量 × 安全系数 / 单机容量

安全系数 = 1.3 ~ 2.0 (通常)

成本优化:
  按需实例 + 预留实例 + Spot实例组合
```

### 5.2 预测模型

```python
import pandas as pd
from prophet import Prophet

# 容量预测
def forecast_capacity(history_data, horizon_days=30):
    df = pd.DataFrame(history_data)
    df.columns = ['ds', 'y']  # Prophet要求列名
    
    model = Prophet(
        yearly_seasonality=True,
        weekly_seasonality=True,
        daily_seasonality=True,
        changepoint_prior_scale=0.05
    )
    
    model.fit(df)
    
    future = model.make_future_dataframe(periods=horizon_days)
    forecast = model.predict(future)
    
    return forecast[['ds', 'yhat', 'yhat_lower', 'yhat_upper']]

# 使用预测结果规划容量
forecast = forecast_capacity(traffic_history)
peak_forecast = forecast['yhat'].max()
required_capacity = peak_forecast * 1.5  # 50% buffer
```

---

## 6. 事件响应

### 6.1 事件严重性分级

| 级别 | 名称 | 响应时间 | 示例 |
|------|------|---------|------|
| P0 | 严重 | 5分钟 | 服务完全中断 |
| P1 | 高 | 30分钟 | 核心功能降级 |
| P2 | 中 | 2小时 | 非核心功能异常 |
| P3 | 低 | 1工作日 | 轻微性能问题 |
| P4 | 极低 | 排期处理 | 优化建议 |

### 6.2 事件响应流程

```
检测 → 分类 → 响应 → 缓解 → 复盘
  │       │       │       │       │
  ▼       ▼       ▼       ▼       ▼
警报   SEV等级  值班响应  修复     事后分析
自动化  影响评估  升级     验证     改进项
```

### 6.3 事后分析模板

```markdown
# 事后分析: [事件标题]

## 摘要
- 日期: YYYY-MM-DD
- 持续时间: HH:MM
- 严重性: P0/P1/P2
- 影响: [用户/服务/数据]

## 时间线
- HH:MM - 检测
- HH:MM - 响应开始
- HH:MM - 缓解
- HH:MM - 恢复

## 根因
[5 Whys分析]

## 教训
- 做得好的:
- 做得差的:

## 改进项
- [ ] 短期修复
- [ ] 长期改进
- [ ] 监控增强
```

---

## 7. 平台工程

### 7.1 IDP (内部开发平台)

```
┌─────────────────────────────────────────┐
│         Internal Developer Platform     │
├─────────────────────────────────────────┤
│                                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ 开发者门户 │ │ 部署平台  │ │ 运维平台  │ │
│  │ Backstage │ │ ArgoCD   │ │ Datadog  │ │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ │
│       │            │            │       │
│       └────────────┼────────────┘       │
│                    │                    │
│              平台API层                  │
│                    │                    │
│       ┌────────────┴────────────┐       │
│       ▼                         ▼       │
│  Kubernetes                Cloud        │
│  Clusters                  Resources    │
│                                         │
└─────────────────────────────────────────┘
```

### 7.2 自助服务能力

```yaml
# 服务模板
apiVersion: platform.company.io/v1
kind: ServiceTemplate
metadata:
  name: microservice-standard
spec:
  template:
    deployment:
      replicas: 3
      resources:
        requests:
          cpu: 100m
          memory: 128Mi
        limits:
          cpu: 1000m
          memory: 512Mi
    monitoring:
      slo:
        availability: 99.9
        latency_p99: 200ms
      alerts:
        - high_error_rate
        - high_latency
    security:
      scan: true
      policies:
        - psp-restricted
```

---

## 8. 最佳实践

### 8.1 SRE团队拓扑

| 模式 | 描述 | 适用 |
|------|------|------|
| 嵌入式 | SRE嵌入产品团队 | 早期阶段 |
| 咨询式 | SRE提供建议 | 中等规模 |
| 平台式 | SRE管理平台 | 大规模 |

### 8.2 关键指标

| 指标 | 目标 | 度量 |
|------|------|------|
| MTTD | < 5分钟 | 平均检测时间 |
| MTTR | < 1小时 | 平均恢复时间 |
| MTBF | > 30天 | 平均故障间隔 |
| 变更成功率 | > 95% | 无事故部署比例 |

---

## 9. 参考文献

1. "Site Reliability Engineering" - Google
2. "The Site Reliability Workbook" - Google
3. "Chaos Engineering" - Netflix
4. "Building Secure & Reliable Systems" - Google
5. "Platform Engineering" - Team Topologies

---

*Last Updated: 2026-04-03*
