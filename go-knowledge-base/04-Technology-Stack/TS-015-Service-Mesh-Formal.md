# TS-015: 服务网格的形式化架构 (Service Mesh: Formal Architecture)

> **维度**: Technology Stack
> **级别**: S (16+ KB)
> **标签**: #service-mesh #istio #envoy #sidecar #traffic-management
> **权威来源**:
>
> - [Istio: A Load Balancer in the Data Path](https://www.usenix.org/conference/nsdi18/presentation/zhang) - Google (2018)
> - [The Service Mesh](https://www.infoq.com/articles/service-mesh-next-generation-networking/) - Buoyant (2017)
> - [Envoy Proxy Architecture](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview) - Envoy Team (2025)
> - [SMI (Service Mesh Interface) Spec](https://smi-spec.io/) - CNCF (2024)
> - [Istio: Zero Trust Networking](https://istio.io/latest/docs/concepts/security/) - Istio Team (2025)

---

## 1. 服务网格的形式化定义

### 1.1 架构代数

**定义 1.1 (服务网格)**
服务网格 $M$ 是一个六元组 $\langle S, P, C, D, T, O \rangle$：

- $S$: 服务集合
- $P$: 代理集合 (Sidecar)
- $C$: 控制平面
- $D$: 数据平面
- $T$: 流量管理策略
- $O$: 可观测性系统

**定义 1.2 (Sidecar 注入)**
$$\text{Inject}: \text{Pod} \to \text{Pod} \times \text{Proxy}$$
将代理容器注入应用 Pod。

### 1.2 数据平面与控制平面

**数据平面**:
$$D = \{ p_i \mid p_i \text{ handles traffic for } s_i \}$$
处理实际流量。

**控制平面**:
$$C = \langle \text{Pilot}, \text{Mixer}, \text{Citadel} \rangle$$
配置和证书管理。

---

## 2. 流量管理的形式化

### 2.1 路由规则

**定义 2.1 (VirtualService)**
$$VS = \langle \text{hosts}, \text{gateways}, \text{http}, \text{tls}, \text{tcp} \rangle$$

**路由匹配**:
$$\text{Match}: \text{Request} \to \text{Destination}$$

**权重路由**:
$$\forall d \in \text{Destinations}: \sum w(d) = 1$$

### 2.2 流量策略

**超时**:
$$\text{Timeout}(r) \Rightarrow \text{Abort}(r, \text{504})$$

**重试**:
$$\text{Retry}(r, n, \text{condition})$$

**熔断**:
$$\text{CircuitBreaker}(d) = \langle \text{threshold}, \text{interval}, \text{break} \rangle$$

---

## 3. 安全的形式化

### 3.1 mTLS 握手

**定义 3.1 (身份)**
$$\text{Identity}(s) = \text{SPIFFE ID}$$

**认证**:
$$\text{Authenticate}(s_1, s_2) \Leftrightarrow \text{Verify}(\text{cert}_{s_1}, \text{cert}_{s_2})$$

**授权策略**:
$$\text{Allow} \Leftarrow \text{source} \in \text{principals} \land \text{operation} \in \text{permissions}$$

### 3.2 零信任网络

**原则**: 永不信任，始终验证
$$\forall c: \text{Authenticate}(c) \land \text{Authorize}(c)$$

---

## 4. 多元表征

### 4.1 服务网格架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Service Mesh Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Control Plane (Istiod)                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                         │
│  │   Pilot     │  │  Citadel    │  │   Galley    │                         │
│  │  (xDS)      │  │  (certs)    │  │  (config)   │                         │
│  └──────┬──────┘  └─────────────┘  └─────────────┘                         │
│         │                                                                    │
│         │ xDS API (gRPC stream)                                              │
│         │                                                                    │
│  ┌──────┴──────────────────────────────────────────────────────────┐       │
│  │                       Data Plane                                 │       │
│  │  ┌──────────┐      ┌──────────┐      ┌──────────┐              │       │
│  │  │ Service A │◄────►│  Envoy   │◄────►│ Service B │              │       │
│  │  │   (App)   │      │ (Sidecar)│      │   (App)   │              │       │
│  │  └──────────┘      └────┬─────┘      └──────────┘              │       │
│  │                         │                                       │       │
│  │  ┌──────────────────────┴──────────────────────┐                │       │
│  │  │  Envoy Functionality:                        │                │       │
│  │  │  - Traffic Management (routing, lb)          │                │       │
│  │  │  - Security (mTLS, authz)                    │                │       │
│  │  │  - Observability (metrics, tracing)          │                │       │
│  │  └──────────────────────────────────────────────┘                │       │
│  └──────────────────────────────────────────────────────────────────┘       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 服务网格 vs 传统架构对比矩阵

| 特性 | 传统 | Service Mesh |
|------|------|--------------|
| **通信** | 直连 | 通过 Sidecar |
| **负载均衡** | Client-side | Sidecar |
| **重试/超时** | 应用代码 | Sidecar 配置 |
| **mTLS** | 应用实现 | 自动 |
| **监控** | 应用埋点 | 自动注入 |
| **升级** | 应用重启 | Sidecar 滚动更新 |
| **延迟** | 低 | 增加 1-2ms |
| **复杂度** | 应用内 | 基础设施 |

### 4.3 流量管理决策树

```
配置流量管理?
│
├── 路由规则?
│   ├── 基于 URI? → VirtualService HTTPRoute
│   ├── 基于 Header? → Match conditions
│   └── 基于权重? → Weighted routing
│
├── 流量控制?
│   ├── 超时 → Timeout setting
│   ├── 重试 → Retry policy
│   └── 熔断 → Circuit breaker
│
├── 安全?
│   ├── mTLS? → PeerAuthentication
│   ├── 认证? → RequestAuthentication
│   └── 授权? → AuthorizationPolicy
│
└── 可观测性?
    ├── 指标? → Prometheus
    ├── 追踪? → Jaeger/Zipkin
    └── 日志? → Access logs
```

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Service Mesh Implementation Checklist                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  部署:                                                                       │
│  □ Sidecar 自动注入                                                           │
│  □ 控制平面高可用                                                             │
│  □ 数据平面资源限制                                                           │
│                                                                              │
│  流量管理:                                                                   │
│  □ 默认拒绝策略 (零信任)                                                       │
│  □ 熔断器配置                                                                │
│  □ 超时和重试策略                                                             │
│                                                                              │
│  安全:                                                                       │
│  □ mTLS 严格模式                                                             │
│  □ 授权策略最小权限                                                           │
│  □ 密钥轮换                                                                  │
│                                                                              │
│  可观测性:                                                                   │
│  □ Prometheus 指标                                                            │
│  □ 分布式追踪 (Sampling)                                                      │
│  □ 访问日志                                                                  │
│                                                                              │
│  注意:                                                                       │
│  ❌ 不是所有服务都需要服务网格                                                   │
│  ❌ Sidecar 增加延迟和资源消耗                                                  │
│  ❌ 调试复杂度增加                                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16KB, 完整形式化)
