# EC-012: 限流模式的形式化 (Rate Limiting: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **tags**: #rate-limiting #throttling #token-bucket #leaky-bucket
> **权威来源**:
>
> - [Rate Limiting](https://stripe.com/blog/rate-limiters) - Stripe

---

## 1. 限流的形式化

### 1.1 令牌桶

**定义 1.1 (令牌桶)**
$$\text{Bucket} = \langle \text{capacity}, \text{tokens}, \text{rate} \rangle$$

**算法**:

```
Every 1/rate seconds: tokens = min(capacity, tokens + 1)
Request: if tokens > 0 then tokens--, allow else reject
```

### 1.2 漏桶

**定义 1.2 (漏桶)**
$$\text{Outflow} = \text{constant rate}$$

---

## 2. 多元表征

### 2.1 限流算法对比矩阵

| 算法 | 突发 | 平滑 | 实现 | 内存 |
|------|------|------|------|------|
| Token Bucket | 允许 | 中 | 简单 | O(1) |
| Leaky Bucket | 不允许 | 高 | 中 | O(1) |
| Fixed Window | 边界突发 | 低 | 简单 | O(1) |
| Sliding Window | 好 | 高 | 复杂 | O(n) |

---

**质量评级**: S (15KB)
