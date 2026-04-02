# AD-009: 容量规划的形式化理论与实践 (Capacity Planning: Formal Theory & Practice)

> **维度**: Application Domains
> **级别**: S (16+ KB)
> **tags**: #capacity-planning #scaling #load-forecasting #performance #sre
> **权威来源**:
>
> - [The Art of Capacity Planning](https://www.oreilly.com/library/view/the-art-of/9780596518578/) - John Allspaw
> - [Site Reliability Engineering](https://sre.google/sre-book/table-of-contents/) - Google
> - [Capacity Planning for Web Operations](https://www.usenix.org/legacy/publications/login/2005-12/pdfs/allspaw.pdf) - USENIX
> - [Forecasting: Principles and Practice](https://otexts.com/fpp3/) - Hyndman & Athanasopoulos

---

## 1. 形式化基础

### 1.1 容量规划定义

**定义 1.1 (容量)**
容量是系统在给定服务质量 (QoS) 约束下处理工作负载的能力。

**定义 1.2 (容量利用率)**
$$U = \frac{\text{实际负载}}{\text{容量}} \times 100\%$$

**定义 1.3 (容量需求)**
$$C_{required} = \frac{L_{peak}}{U_{target}} \times SF$$

其中：

- $L_{peak}$: 峰值负载
- $U_{target}$: 目标利用率 (通常 60-70%)
- $SF$: 安全系数

### 1.2 容量规划定理

**定理 1.1 (利用率与延迟关系)**
根据排队论，当利用率 $U \to 1$ 时，平均延迟 $W \to \infty$。

*证明* (基于 M/M/1 队列):
$$W = \frac{1}{\mu - \lambda} = \frac{1}{\mu(1 - U)}$$
当 $U \to 1$，分母 $\to 0$，故 $W \to \infty$。

$\square$

**公理 1.1 (容量安全边际)**
生产系统应保持至少 30% 的容量余量以应对突发流量。

---

## 2. 容量规划模型

### 2.1 Little's Law

**定理 2.1 (Little's Law)**
对于稳定系统：
$$L = \lambda \cdot W$$

其中：

- $L$: 系统中平均请求数
- $\lambda$: 平均到达率 (请求/秒)
- $W$: 平均停留时间 (秒)

**应用**: 若已知 QPS 和目标响应时间，可计算所需并发处理能力。

### 2.2 扩展公式

**定义 2.1 (水平扩展)**
$$C_{total} = n \cdot C_{unit} \cdot \eta$$

其中 $\eta$ 是扩展效率因子 (通常 0.8-0.95，考虑协调开销)。

**定义 2.2 (垂直扩展)**
$$C_{new} = C_{base} \cdot f(scale\_factor)$$

垂直扩展通常非线性，受 Amdahl 定律约束。

### 2.3 容量规划对比矩阵

| 维度 | 水平扩展 | 垂直扩展 |
|------|----------|----------|
| **上限** | 理论无上限 | 硬件限制 |
| **成本模型** | 线性 | 超线性 |
| **故障隔离** | 好 | 差 |
| **复杂度** | 高 | 低 |
| **适用场景** | 无状态服务 | 数据库、单体应用 |

---

## 3. 负载预测模型

### 3.1 预测方法

**定义 3.1 (线性预测)**
$$L(t) = L_0 + r \cdot t$$

其中 $r$ 是增长率。

**定义 3.2 (指数增长)**
$$L(t) = L_0 \cdot (1 + r)^t$$

**定义 3.3 (季节性调整)**
$$L(t) = T(t) \cdot S(t) \cdot I(t)$$

其中 $T$ 是趋势，$S$ 是季节性因子，$I$ 是随机因子。

### 3.2 预测精度评估

| 指标 | 公式 | 说明 |
|------|------|------|
| MAE | $\frac{1}{n}\sum|y_i - \hat{y}_i|$ | 平均绝对误差 |
| RMSE | $\sqrt{\frac{1}{n}\sum(y_i - \hat{y}_i)^2}$ | 均方根误差 |
| MAPE | $\frac{100\%}{n}\sum|\frac{y_i - \hat{y}_i}{y_i}|$ | 平均绝对百分比误差 |

---

## 4. 多元表征

### 4.1 容量规划决策树

```
需要规划系统容量?
│
├── 了解当前负载?
│   ├── 否 → 部署监控和指标收集
│   │       ├── APM 工具 (Datadog, New Relic)
│   │       ├── 基础设施监控 (Prometheus)
│   │       └── 日志分析 (ELK)
│   │
│   └── 是 → 收集历史数据
│           ├── QPS/吞吐量
│           ├── 响应时间 (p50, p95, p99)
│           ├── 错误率
│           └── 资源使用率 (CPU, 内存, IO, 网络)
│
├── 识别瓶颈资源?
│   ├── CPU 限制? → 计算密集型优化或扩展
│   ├── 内存限制? → 缓存优化或增加内存
│   ├── IO 限制? → 使用 SSD, 数据库优化
│   ├── 网络限制? → 带宽升级或 CDN
│   └── 并发限制? → 连接池优化或水平扩展
│
├── 预测增长?
│   ├── 业务增长率?
│   │   └── 历史数据 → 线性/指数回归
│   ├── 季节性波动?
│   │   └── 使用季节性分解
│   └── 营销事件?
│       └── 准备突发容量 (3-5x)
│
└── 制定容量计划
    ├── 短期 (1-3个月)
    │   └── 根据当前趋势线性扩展
    ├── 中期 (3-12个月)
    │   └── 考虑架构优化
    └── 长期 (1年+)
        └── 可能的技术栈升级
```

### 4.2 容量规划检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Capacity Planning Checklist                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  数据收集阶段                                                                │
│  ═══════════════════════════════════════════════════════════════════════     │
│  □ 收集至少 30 天历史数据                                                    │
│  □ 识别每日/每周/季节性模式                                                  │
│  □ 记录峰值负载和持续时间                                                    │
│  □ 测量各资源组件使用率                                                      │
│  □ 建立基线性能指标                                                          │
│                                                                              │
│  分析阶段                                                                    │
│  ═══════════════════════════════════════════════════════════════════════     │
│  □ 识别瓶颈组件                                                              │
│  □ 计算单位容量 (每实例 QPS)                                                 │
│  □ 测量扩展效率                                                              │
│  □ 评估单点故障风险                                                          │
│                                                                              │
│  预测阶段                                                                    │
│  ═══════════════════════════════════════════════════════════════════════     │
│  □ 应用增长模型 (线性/指数)                                                  │
│  □ 考虑产品发布影响                                                          │
│  □ 预留突发容量 (30-50%)                                                     │
│  □ 设置利用率警戒线 (70%)                                                    │
│                                                                              │
│  实施阶段                                                                    │
│  ═══════════════════════════════════════════════════════════════════════     │
│  □ 优先水平扩展                                                              │
│  □ 实施自动扩展策略                                                          │
│  □ 建立容量审查流程                                                          │
│  □ 定期回顾和更新计划                                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 容量规划工具对比

| 工具 | 类型 | 强项 | 适用场景 |
|------|------|------|----------|
| **Excel/Google Sheets** | 电子表格 | 简单、灵活 | 小规模、初创公司 |
| **Datadog** | APM | 全栈监控 | 企业级云原生 |
| **Prometheus + Grafana** | 开源监控 | 可扩展、免费 | 自建监控体系 |
| **AWS Cost Explorer** | 云成本 | 与 AWS 集成 | AWS 环境 |
| **Turbo360** | 成本优化 | 多云支持 | 多云环境 |
| **PlanForCloud** | 预测工具 | 场景建模 | 长期规划 |

---

## 5. 负载测试与验证

### 5.1 负载测试类型

| 类型 | 目标 | 持续时间 | 负载模式 |
|------|------|----------|----------|
| **负载测试** | 验证容量 | 持续 | 目标负载 |
| **压力测试** | 找到极限 | 递增 | 直到失败 |
| **浸泡测试** | 发现内存泄漏 | 长时间 | 稳定负载 |
| **峰值测试** | 验证突发处理 | 短时 | 突然高峰 |
| **断点测试** | 确定故障点 | 递增 | 逐步增加 |

### 5.2 负载测试工具

```go
// 使用 Go 编写负载测试示例
package loadtest

import (
    "context"
    "sync"
    "testing"
    "time"
)

// LoadTestSpec 定义负载测试规格
type LoadTestSpec struct {
    Duration       time.Duration
    Concurrency    int
    RampUpTime     time.Duration
    TargetRPS      int
    RequestFunc    func(ctx context.Context) error
}

// RunLoadTest 执行负载测试
func RunLoadTest(spec LoadTestSpec) *LoadTestResult {
    ctx, cancel := context.WithTimeout(context.Background(), spec.Duration)
    defer cancel()

    result := &LoadTestResult{
        StartTime: time.Now(),
    }

    var wg sync.WaitGroup
    sem := make(chan struct{}, spec.Concurrency)

    ticker := time.NewTicker(time.Second / time.Duration(spec.TargetRPS))
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            wg.Wait()
            result.EndTime = time.Now()
            return result
        case <-ticker.C:
            sem <- struct{}{}
            wg.Add(1)

            go func() {
                defer wg.Done()
                defer func() { <-sem }()

                start := time.Now()
                err := spec.RequestFunc(ctx)
                latency := time.Since(start)

                result.Record(latency, err)
            }()
        }
    }
}

type LoadTestResult struct {
    StartTime    time.Time
    EndTime      time.Time
    TotalReqs    int64
    SuccessReqs  int64
    FailedReqs   int64
    Latencies    []time.Duration
}

func (r *LoadTestResult) Record(latency time.Duration, err error) {
    r.TotalReqs++
    if err != nil {
        r.FailedReqs++
    } else {
        r.SuccessReqs++
    }
    r.Latencies = append(r.Latencies, latency)
}

func (r *LoadTestResult) SuccessRate() float64 {
    if r.TotalReqs == 0 {
        return 0
    }
    return float64(r.SuccessReqs) / float64(r.TotalReqs) * 100
}
```

---

## 6. 云原生容量规划

### 6.1 Kubernetes 容量管理

```yaml
# HorizontalPodAutoscaler 示例
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-service
  minReplicas: 3
  maxReplicas: 50
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

### 6.2 成本优化策略

| 策略 | 实施复杂度 | 节省潜力 | 风险 |
|------|------------|----------|------|
| **预留实例** | 低 | 30-60% | 长期承诺 |
| **Spot 实例** | 中 | 60-90% | 可能中断 |
| **自动伸缩** | 中 | 20-40% | 配置不当 |
| **右规模** | 中 | 10-30% | 性能影响 |
| **无服务器** | 高 | 变量 | 冷启动 |

---

## 7. 关系网络

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Capacity Planning Context                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  相关领域                                                                    │
│  ├── Performance Engineering                                                │
│  ├── Site Reliability Engineering (SRE)                                     │
│  ├── DevOps / Platform Engineering                                          │
│  ├── Cloud Cost Management                                                  │
│  └── Business Continuity Planning                                           │
│                                                                              │
│  理论基础                                                                    │
│  ├── Queuing Theory (Little's Law)                                          │
│  ├── Time Series Analysis                                                   │
│  ├── Statistical Forecasting                                                │
│  └── Control Theory                                                         │
│                                                                              │
│  实践框架                                                                    │
│  ├── Google SRE Book - Capacity Planning                                    │
│  ├── AWS Well-Architected - Cost Optimization                               │
│  ├── Azure Advisor                                                          │
│  └── FinOps Foundation                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. 思维工具总结

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Capacity Planning Toolkit                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  核心公式                                                                    │
│  ═══════════════════════════════════════════════════════════════════════     │
│  • 容量需求 = 峰值负载 / 目标利用率 × 安全系数                               │
│  • Little's Law: L = λ × W                                                  │
│  • 扩展后容量 = n × 单位容量 × 效率因子                                      │
│                                                                              │
│  关键指标                                                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ 利用率: 目标 60-70%，警戒线 80%                                     │    │
│  │ 响应时间: p95 < 200ms, p99 < 500ms                                  │    │
│  │ 错误率: < 0.1%                                                      │    │
│  │ 并发数: 根据 Little's Law 计算                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  扩展原则                                                                    │
│  ═══════════════════════════════════════════════════════════════════════     │
│  1. 先水平，后垂直                                                          │
│  2. 无状态服务优先水平扩展                                                  │
│  3. 有状态服务考虑分片                                                      │
│  4. 预留 30-50% 突发容量                                                    │
│  5. 自动化伸缩策略                                                          │
│                                                                              │
│  常见误区:                                                                   │
│  ❌ 仅关注 CPU/内存，忽略 IO/网络                                            │
│  ❌ 忽视数据库连接池限制                                                      │
│  ❌ 未考虑第三方服务容量                                                      │
│  ❌ 过度配置造成浪费                                                          │
│  ❌ 缺乏定期容量审查                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16KB)
**完成日期**: 2026-04-02
