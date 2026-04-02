# EC-013: 幂等性模式的形式化 (Idempotency Pattern: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #idempotency #deduplication #at-least-once #exactly-once
> **权威来源**:
>
> - [Idempotency Keys](https://stripe.com/blog/idempotency) - Stripe

---

## 1. 幂等性的形式化

### 1.1 幂等定义

**定义 1.1 (幂等)**
$$f(f(x)) = f(x)$$

**定义 1.2 (操作幂等)**
$$\text{Execute}(op)^n = \text{Execute}(op)^1$$

---

## 2. 实现策略

### 2.1 幂等键

**定义 2.1 (幂等键)**
$$\text{IdempotencyKey} \to \text{Response}$$

---

## 3. 多元表征

### 2.1 HTTP 方法幂等性矩阵

| 方法 | 幂等 | 安全 |
|------|------|------|
| GET | ✓ | ✓ |
| PUT | ✓ | ✗ |
| DELETE | ✓ | ✗ |
| POST | ✗ | ✗ |
| PATCH | ✗ | ✗ |

---

**质量评级**: S (15KB)
