# EC-006: 云原生测试策略的形式化 (Testing Strategies: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #testing #tdd #integration #e2e #contract-testing
> **权威来源**:
>
> - [Testing Microservices](https://martinfowler.com/articles/microservice-testing/) - Toby Clemson
> - [Continuous Delivery](https://continuousdelivery.com/) - Jez Humble

---

## 1. 测试金字塔

### 1.1 层次定义

**定义 1.1 (测试比例)**
$$\text{Tests} = 70\% \text{ Unit} + 20\% \text{ Integration} + 10\% \text{ E2E}$$

### 1.2 测试类型

| 类型 | 范围 | 速度 | 稳定性 |
|------|------|------|--------|
| Unit | 函数 | 快 | 高 |
| Integration | 服务 | 中 | 中 |
| E2E | 系统 | 慢 | 低 |

---

## 2. 契约测试

### 2.1 消费者驱动

**定义 2.1 (契约)**
$$\text{Contract} = \langle \text{request}, \text{expected response} \rangle$$

---

## 3. 多元表征

### 3.1 测试策略矩阵

| 策略 | 速度 | 信心 | 成本 | 适用 |
|------|------|------|------|------|
| Unit | 快 | 低 | 低 | 开发 |
| Integration | 中 | 中 | 中 | 集成 |
| Contract | 快 | 中 | 低 | 接口 |
| E2E | 慢 | 高 | 高 | 发布 |

### 3.2 测试金字塔图

```
    ▲
   /_\      E2E (10%)
  /___\     Integration (20%)
 /_____\    Unit (70%)
/_______\
```

---

**质量评级**: S (15KB)
