# AD-009: 容量规划与扩展策略 (Capacity Planning & Scaling Strategies)

> **维度**: Application Domains
> **级别**: S (16+ KB)
> **标签**: #capacity-planning #scaling #load-testing #resource-planning
> **权威来源**: [The Art of Capacity Planning](https://www.oreilly.com/library/view/the-art-of/9780596518578/), [Google SRE Book](https://sre.google/sre-book/table-of-contents/)

---

## 容量规划模型

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Capacity Planning Framework                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 需求预测                                                                  │
│     ├── 历史数据分析 (时间序列预测)                                            │
│     ├── 业务增长预测                                                          │
│     └── 季节性/事件性波动                                                      │
│                                                                              │
│  2. 容量计算                                                                  │
│     ├── 单实例容量 = RPS/QPS × Latency                                         │
│     ├── 所需实例数 = 总需求 / 单实例容量                                        │
│     └── 冗余系数 = 1 / (1 - 目标利用率)                                        │
│                                                                              │
│  3. 验证测试                                                                  │
│     ├── 负载测试 (Load Testing)                                                │
│     ├── 压力测试 (Stress Testing)                                              │
│     └── 混沌测试 (Chaos Engineering)                                           │
│                                                                              │
│  4. 持续监控                                                                  │
│     ├── 关键指标告警                                                          │
│     ├── 自动扩缩容                                                             │
│     └── 定期容量评审                                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 计算公式

### 基本公式

```
容量需求 = 峰值流量 × (1 + 安全边际)

单实例容量 = (1 / 平均响应时间) × 并发连接数

所需实例数 = ceil(总流量 / 单实例容量 × 冗余系数)

资源利用率 = 实际使用 / 总容量

目标利用率 = 通常 60-70% (保留突发缓冲)

Little's Law:
并发用户数 = 吞吐量 × 平均响应时间
L = λ × W

示例:
- 目标: 支持 10,000 RPS
- 单实例: 100 RPS @ 100ms 延迟
- 所需实例: 10,000 / 100 = 100 实例
- 冗余系数: 1.5
- 最终: 150 实例
```

### 存储容量规划

```
存储需求 = 日增量 × 保留天数 × (1 + 增长系数)

压缩后存储 = 原始大小 × 压缩率

示例:
- 日日志量: 1TB
- 保留期: 30 天
- 压缩率: 0.3
- 原始存储: 30TB
- 实际存储: 30TB × 0.3 = 9TB
- 年增长: 50%
- 1 年后: 9TB × 1.5 = 13.5TB
```

---

## 扩展策略

### 水平扩展 vs 垂直扩展

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Scaling Strategies                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  垂直扩展 (Scale Up)                                                          │
│  ┌─────────┐         ┌─────────┐         ┌─────────┐                       │
│  │   2C    │   ─►    │   8C    │   ─►    │  32C    │                       │
│  │  4GB    │         │  16GB   │         │  128GB  │                       │
│  └─────────┘         └─────────┘         └─────────┘                       │
│       │                  │                  │                              │
│       └──────────────────┴──────────────────┘                              │
│                   限制: 硬件上限，成本指数增长                                │
│                   适用: 数据库，单点服务                                      │
│                                                                              │
│  水平扩展 (Scale Out)                                                         │
│  ┌─────┐ ┌─────┐ ┌─────┐    ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐              │
│  │ 2C  │ │ 2C  │ │ 2C  │    │ 2C  │ │ 2C  │ │ 2C  │ │ 2C  │              │
│  │4GB  │ │4GB  │ │4GB  │ ──►│4GB  │ │4GB  │ │4GB  │ │4GB  │              │
│  └─────┘ └─────┘ └─────┘    └─────┘ └─────┘ └─────┘ └─────┘              │
│     3 节点                    6 节点                                        │
│                                                                              │
│     优势: 线性扩展，容错性好，成本可控                                        │
│     适用: 无状态服务，分布式系统                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 自动扩缩容

```go
package capacity

import (
    "context"
    "math"
    "time"
)

// AutoScaler 自动扩缩容
type AutoScaler struct {
    minReplicas int
    maxReplicas int
    targetCPU   float64 // 目标 CPU 利用率
    targetRPS   float64 // 目标 RPS/实例

    metrics MetricsClient
    scaler Scaler
}

type ScalingDecision struct {
    CurrentReplicas int
    DesiredReplicas int
    Reason          string
    Metrics         ScalingMetrics
}

type ScalingMetrics struct {
    CPUUtilization float64
    MemoryUsage    float64
    RPS            float64
    LatencyP99     time.Duration
}

func (a *AutoScaler) Evaluate(ctx context.Context) (*ScalingDecision, error) {
    // 获取当前指标
    metrics, err := a.metrics.GetMetrics(ctx)
    if err != nil {
        return nil, err
    }

    currentReplicas, err := a.scaler.GetReplicas(ctx)
    if err != nil {
        return nil, err
    }

    // 计算期望副本数
    desiredReplicas := a.calculateDesiredReplicas(metrics, currentReplicas)

    // 边界检查
    if desiredReplicas < a.minReplicas {
        desiredReplicas = a.minReplicas
    }
    if desiredReplicas > a.maxReplicas {
        desiredReplicas = a.maxReplicas
    }

    return &ScalingDecision{
        CurrentReplicas: currentReplicas,
        DesiredReplicas: desiredReplicas,
        Reason:          a.getScalingReason(metrics),
        Metrics:         metrics,
    }, nil
}

func (a *AutoScaler) calculateDesiredReplicas(metrics ScalingMetrics, current int) int {
    // 基于 CPU
    cpuBased := int(math.Ceil(float64(current) * metrics.CPUUtilization / a.targetCPU))

    // 基于 RPS
    rpsBased := int(math.Ceil(metrics.RPS / a.targetRPS))

    // 取最大值
    desired := max(cpuBased, rpsBased)

    // 平滑处理 (避免震荡)
    if abs(desired-current) <= 2 {
        return current // 变化太小，保持稳定
    }

    return desired
}

func (a *AutoScaler) getScalingReason(metrics ScalingMetrics) string {
    if metrics.CPUUtilization > a.targetCPU*1.2 {
        return "high_cpu_utilization"
    }
    if metrics.RPS > a.targetRPS*float64(a.maxReplicas) {
        return "high_request_rate"
    }
    if metrics.CPUUtilization < a.targetCPU*0.5 {
        return "low_cpu_utilization"
    }
    return "stable"
}
```

---

## 容量测试

### 测试类型

| 测试类型 | 目的 | 方法 |
|----------|------|------|
| 负载测试 | 验证正常负载下的表现 | 模拟预期流量 |
| 压力测试 | 找到系统极限 | 逐步增加直到崩溃 |
| 浸泡测试 | 验证长期稳定性 | 持续运行数天 |
| 峰值测试 | 验证突发处理能力 | 突然增加流量 |
| 容量测试 | 确定最大容量 | 找到资源瓶颈 |

### 测试工具

```bash
# k6 - 现代负载测试工具
k6 run --vus 100 --duration 30s script.js

# Vegeta - HTTP 负载测试
echo "GET http://localhost:8080/" | vegeta attack -rate=1000 -duration=60s | vegeta report

# Locust - Python 分布式测试
locust -f locustfile.py --host=http://localhost:8080

# JMeter - 企业级测试
jmeter -n -t test_plan.jmx -l results.jtl
```

---

## 容量规划检查清单

- [ ] 定义 SLA (可用性、延迟、吞吐量)
- [ ] 收集历史流量数据
- [ ] 识别关键路径和瓶颈
- [ ] 计算单实例容量
- [ ] 确定冗余和安全边际
- [ ] 进行负载测试验证
- [ ] 配置监控和告警
- [ ] 制定扩容计划
- [ ] 定期容量评审 (月度/季度)

---

## 参考文献

1. [The Art of Capacity Planning](https://www.oreilly.com/library/view/the-art-of/9780596518578/) - Arun Kejariwal
2. [Google SRE Book](https://sre.google/sre-book/table-of-contents/)
3. [Performance Testing Guidance](https://learn.microsoft.com/en-us/azure/architecture/guide/design-principles/performance-testing)
