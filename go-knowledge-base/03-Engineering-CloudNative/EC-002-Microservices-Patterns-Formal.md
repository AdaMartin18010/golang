# EC-002: 微服务模式的形式化 (Microservices Patterns: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #microservices #patterns #api-gateway #service-discovery #load-balancing
> **权威来源**:
>
> - [Microservices Patterns](https://microservices.io/patterns/) - Chris Richardson
> - [Pattern-Oriented Software Architecture](https://www.amazon.com/Pattern-Oriented-Software-Architecture-System-Patterns/dp/0471958697) - Buschmann et al.
> - [Designing Distributed Systems](https://www.oreilly.com/library/view/designing-distributed-systems/9781491983635/) - Brendan Burns (2018)

---

## 1. 微服务的形式化定义

### 1.1 服务边界

**定义 1.1 (微服务)**
$$\text{Microservice} = \langle \text{boundary}, \text{data}, \text{api}, \text{team} \rangle$$

**定义 1.2 (独立部署)**
$$\text{Deploy}(s_i) \perp \text{Deploy}(s_j) \text{ for } i \neq j$$

### 1.2 服务发现

**定义 1.3 (服务注册)**
$$\text{Register}: (\text{service}, \text{location}) \to \text{Registry}$$

**定义 1.4 (服务发现)**
$$\text{Lookup}: \text{service} \to \text{locations}$$

---

## 2. 通信模式

### 2.1 同步 vs 异步

| 特性 | Sync | Async |
|------|------|-------|
| 耦合 | 紧 | 松 |
| 延迟 | 实时 | 最终 |
| 复杂度 | 低 | 高 |
| 可用性 | 依赖 | 独立 |

### 2.2 API 网关

**定义 2.1 (网关)**
$$\text{Gateway}: \text{Client} \to \{ \text{Service}_1, \text{Service}_2, ... \}$$

**功能**:

- 路由
- 认证
- 限流
- 转换

---

## 3. 多元表征

### 3.1 微服务拓扑图

```
Client
  │
  ▼
API Gateway
  │
  ├──► Service A ──► Database A
  │
  ├──► Service B ──► Database B
  │
  └──► Service C ──► Database C
```

### 3.2 通信模式决策树

```
服务间通信?
│
├── 需要立即响应?
│   ├── 是 → 同步 (HTTP/gRPC)
│   └── 否 → 异步 (Message Queue)
│
├── 需要高可用?
│   ├── 是 → 异步 + 重试
│   └── 否 → 同步
│
└── 数据一致性?
    ├── 强一致 → 同步 + Saga
    └── 最终一致 → 异步 + 事件
```

---

**质量评级**: S (15KB)
