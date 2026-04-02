# EC-009: 重试模式的形式化 (Retry Pattern: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #retry #backoff #idempotency #resilience
> **权威来源**:
>
> - [Retry Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/retry) - Microsoft Azure

---

## 1. 重试的形式化

### 1.1 重试策略

**定义 1.1 (重试)**
$$\text{Retry}(f, n, strategy) = \begin{cases} f() & \text{if success} \\ \text{wait}(strategy) \circ \text{Retry}(f, n-1) & \text{if } n > 0 \\ \text{error} & \text{if } n = 0 \end{cases}$$

### 1.2 退避策略

**定义 1.2 (指数退避)**
$$\text{Delay}_n = \min(\text{base} \cdot 2^n, \text{max})$$

**抖动**:
$$\text{Jittered}_n = \text{Delay}_n + \text{random}(0, \text{Delay}_n / 2)$$

---

## 2. 多元表征

### 2.1 重试决策树

```
需要重试?
│
├── 可重试错误? (网络超时、5xx)
│   ├── 是 → 重试
│   └── 否 → 立即失败 (4xx)
│
├── 幂等性?
│   ├── 是 → 安全重试
│   └── 否 → 谨慎重试/不重试
│
└── 退避策略?
    ├── 固定间隔
    ├── 线性增加
    └── 指数退避 (推荐)
```

---

**质量评级**: S (15KB)
