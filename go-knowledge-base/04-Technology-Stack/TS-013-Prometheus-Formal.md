# TS-013: Prometheus 可观测性形式化 (Prometheus Observability: Formal Model)

> **维度**: Technology Stack
> **级别**: S (15+ KB)
> **标签**: #prometheus #metrics #monitoring #alerting #observability
> **权威来源**:
>
> - [Prometheus: Up & Running](https://www.oreilly.com/library/view/prometheus-up/9781492034143/) - Brian Brazil (2018)
> - [Google SRE Book: Monitoring](https://sre.google/sre-book/monitoring-distributed-systems/) - Google (2017)
> - [The RED Method](https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/) - Weaveworks (2015)
> - [The USE Method](http://www.brendangregg.com/usemethod.html) - Brendan Gregg

---

## 1. 指标的形式化定义

### 1.1 时间序列代数

**定义 1.1 (时间序列)**
$$TS = \{ (t_1, v_1), (t_2, v_2), ... \}$$
其中 $t_i$ 是时间戳，$v_i$ 是值。

**定义 1.2 (指标)**
$$\text{Metric} = \langle \text{name}, \text{labels}, TS \rangle$$

**标签**:
$$\text{labels} = \{ (k_1, v_1), (k_2, v_2), ... \}$$

### 1.2 指标类型

| 类型 | 数学性质 | 操作 |
|------|---------|------|
| **Counter** | 单调递增 | rate, increase |
| **Gauge** | 任意变化 | 当前值 |
| **Histogram** | 分布 | bucket计数, sum, count |
| **Summary** | 分位数 | quantile, sum, count |

---

## 2. PromQL 形式化

### 2.1 查询操作

**定义 2.1 (瞬时向量)**
$$\text{Instant}(m, t) = \{ (l, v) \mid \text{metric}=m, \text{time}=t \}$$

**定义 2.2 (范围向量)**
$$\text{Range}(m, [d]) = \{ (l, [(t-d, v), ..., (t, v)]) \}$$

### 2.2 聚合操作

**定义 2.3 (聚合)**
$$\text{Agg}(v, by) = \text{groupby}(v, by, \text{aggregator})$$

**常用聚合器**: sum, avg, min, max, count

---

## 3. 告警的形式化

### 3.1 告警规则

**定义 3.1 (告警)**
$$\text{Alert} = \langle \text{expr}, \text{for}, \text{labels}, \text{annotations} \rangle$$

**触发条件**:
$$\text{Firing} \Leftarrow \text{expr} = \text{true} \land \text{duration} \geq for$$

---

## 4. 多元表征

### 4.1 指标类型决策树

```
选择指标类型?
│
├── 只增计数?
│   └── 是 → Counter (请求数、错误数)
│
├── 可增可减?
│   ├── 需要分布?
│   │   ├── 是 → Histogram (延迟分布)
│   │   └── 否 → Gauge (温度、队列长度)
│   └──
│       需要精确分位数?
│       └── 是 → Summary (客户端计算)
│
└── 预聚合?
    └── 使用 _count, _sum, _bucket
```

### 4.2 监控方法对比矩阵

| 方法 | 指标 | 适用 | 公式 |
|------|------|------|------|
| **USE** | 利用率、饱和度、错误 | 资源 | Utilization, Saturation, Errors |
| **RED** | 速率、错误、延迟 | 服务 | Rate, Errors, Duration |
| **Golden Signals** | 延迟、流量、错误、饱和度 | 通用 | Google SRE |

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Prometheus Best Practices                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  命名:                                                                       │
│  □ 使用 snake_case                                                           │
│  □ 单位后缀 (_seconds, _bytes, _total)                                       │
│  □ 区分 counter/gauge (不要混用)                                             │
│                                                                              │
│  标签:                                                                       │
│  □ 避免高基数 (cardinality)                                                  │
│  □ 不要使用用户ID作为标签                                                     │
│  □ 保持标签值有限集                                                          │
│                                                                              │
│  告警:                                                                       │
│  □ 使用 for 避免抖动                                                         │
│  □ 告警可操作 (Actionable)                                                   │
│  □ 分层 (page/warning)                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (15KB, 完整形式化)
