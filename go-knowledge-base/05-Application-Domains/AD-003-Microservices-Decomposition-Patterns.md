# AD-003: 微服务拆分与边界划分 (Microservices Decomposition & Boundary Patterns)

> **维度**: Application Domains
> **级别**: S (18+ KB)
> **标签**: #microservices #decomposition #bounded-context #service-boundary
> **权威来源**: [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman

---

## 拆分策略

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Decomposition Strategies                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Monolith                                                                  │
│     │                                                                        │
│     ├──► Decompose by Business Capability                                   │
│     │      ├── Order Service                                                │
│     │      ├── Inventory Service                                            │
│     │      └── Shipping Service                                             │
│     │                                                                        │
│     ├──► Decompose by Subdomain (DDD)                                       │
│     │      ├── Core Domain (Competitive advantage)                          │
│     │      ├── Supporting Subdomain                                         │
│     │      └── Generic Subdomain                                            │
│     │                                                                        │
│     └──► Decompose by Entity (Anti-pattern!)                                │
│            ├── User Service (CRUD only)                                     │
│            └── Order Service (CRUD only)                                    │
│            ⚠️ 避免贫血领域模型                                                │
│                                                                              │
│  推荐组合: Business Capability + DDD Subdomain                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 边界划分模式

### 1. 按业务能力 (Business Capability)

```
E-commerce System
├── Sales Capability
│   ├── Pricing
│   ├── Promotion
│   └── Discount
│   └── Service: pricing-service
│
├── Fulfillment Capability
│   ├── Inventory
│   ├── Warehouse
│   └── Shipping
│   └── Service: fulfillment-service
│
└── Customer Capability
    ├── Profile
    ├── Preference
    └── Support
    └── Service: customer-service
```

### 2. 按子域 (Subdomain)

```
Core Domain (最复杂，内部开发)
├── Order Management
├── Payment Processing
└── Risk Assessment

Supporting Subdomain (可外包)
├── Notification
├── Document Generation
└── Reporting

Generic Subdomain (使用现成方案)
├── Authentication (Auth0, Keycloak)
├── Logging (ELK)
└── Monitoring (Prometheus)
```

---

## 服务粒度决策

### 粒度评估维度

| 维度 | 过小 | 合适 | 过大 |
|------|------|------|------|
| 团队规模 | < 2人 | 5-9人 | > 15人 |
| 代码量 | < 1000行 | 1-10万行 | > 50万行 |
| 变更频率 | 极少 | 独立演进 | 耦合发布 |
| 数据库 | 共享表 | 独立 Schema | 多服务共享 |

### 反模式：纳米服务 (Nanoservices)

```
反模式：服务过小
├── User-Service (只有 CRUD)
├── User-Profile-Service
├── User-Preference-Service
└── User-Auth-Service

问题：
- 网络开销 > 业务逻辑
- 部署复杂
- 难以测试

正解：
└── User-Service (内聚的功能集合)
    ├── Profile
    ├── Preference
    └── Auth
```

---

## 数据库边界

### 每个服务一个数据库

```
Order Service ──► Order DB
Inventory Service ──► Inventory DB
Shipping Service ──► Shipping DB

跨服务查询：
- 避免直接 JOIN
- 使用 API 或 事件
- CQRS 模式

数据一致性：
- 服务内：ACID
- 跨服务：Saga / Eventual Consistency
```

---

## 拆分检查清单

- [ ] 是否可以独立部署？
- [ ] 是否有独立的数据库？
- [ ] 团队是否可以独立开发？
- [ ] 失败是否隔离？
- [ ] 技术栈是否可以不同？
- [ ] 是否可以独立扩展？

---

## 参考文献

1. [Building Microservices](https://samnewman.io/books/building_microservices/) - Sam Newman
2. [Monolith to Microservices](https://samnewman.io/books/monolith-to-microservices/) - Sam Newman
3. [Domain-Driven Design](https://www.domainlanguage.com/ddd/) - Eric Evans
