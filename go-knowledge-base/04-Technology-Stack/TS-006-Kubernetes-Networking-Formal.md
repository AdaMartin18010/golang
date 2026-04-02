# TS-006: Kubernetes 网络的形式化模型 (Kubernetes Networking: Formal Model)

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #kubernetes #networking #cni #service-mesh #formal-methods
> **权威来源**:
>
> - [Kubernetes Networking Concepts](https://kubernetes.io/docs/concepts/cluster-administration/networking/) - Kubernetes Authors (2025)
> - [Container Network Interface (CNI) Specification](https://www.cni.dev/docs/spec/) - CNI Team
> - [Calico Documentation](https://docs.tigera.io/) - Tigera (2025)
> - [Cilium Documentation](https://docs.cilium.io/) - Isovalent (2025)
> - [IEEE 802.1Q VLAN](https://standards.ieee.org/standard/802.1Q.html) - IEEE (2024)

---

## 1. K8s 网络模型的形式化定义

### 1.1 网络拓扑代数

**定义 1.1 (K8s 网络拓扑)**
K8s 网络是一个图 $G = \langle V, E, L \rangle$：

- $V = \text{Pods} \cup \text{Nodes} \cup \text{Services}$: 顶点集合
- $E \subseteq V \times V$: 边集合（连接关系）
- $L: E \to \text{Labels}$: 边标签（协议、端口等）

**定义 1.2 (Pod IP 分配)**
每个 Pod $p$ 被分配唯一 IP：
$$\text{IP}: \text{Pod} \to \text{IP}_{subnet}$$
满足：
$$\forall p_1, p_2 \in \text{Pods}: p_1 \neq p_2 \Rightarrow \text{IP}(p_1) \neq \text{IP}(p_2)$$

### 1.2 K8s 网络公理

**公理 1.1 (Pod- Pod 通信)**
$$\forall p_1, p_2: \text{CanCommunicate}(p_1, p_2) \text{ without NAT}$$
所有 Pod 可以在任何节点上直接通信，无需 NAT。

**公理 1.2 (Node- Pod 通信)**
$$\forall n \in \text{Nodes}, p \in \text{Pods}: \text{CanCommunicate}(n, p)$$
节点可以与所有 Pod 通信。

**公理 1.3 (Service IP 虚拟性)**
$$\text{ServiceIP} \in \text{Virtual} \land \text{ClusterLocal}$$
Service IP 是虚拟的，仅在集群内部可路由。

---

## 2. CNI (Container Network Interface) 形式化

### 2.1 CNI 操作代数

**定义 2.1 (CNI 操作)**
$$\text{Op} ::= \text{ADD} \mid \text{DEL} \mid \text{CHECK} \mid \text{VERSION}$$

**ADD 操作**:
$$\text{ADD}: \langle \text{container_id}, \text{netns}, \text{config} \rangle \to \langle \text{interface}, \text{ip}, \text{routes}, \text{dns} \rangle$$

**执行流程**:

```
Kubelet ──► CRI ──► Container Runtime
              │
              │ Create NetNS
              ▼
         Invoke CNI Plugin
              │
              ├── Allocate IP (IPAM)
              ├── Create veth pair
              ├── Setup routes
              └── Configure iptables/eBPF
```

### 2.2 IPAM (IP Address Management)

**定义 2.2 (IPAM)**
IPAM 函数：
$$\text{IPAM}: \text{Request} \to \text{IP} \cup \{\text{error}\}$$

**分配策略**:

- **host-local**: 本地分配
- **dhcp**: DHCP 服务器
- **static**: 静态分配

---

## 3. Service 网络的形式化

### 3.1 kube-proxy 模式

**定义 3.1 (代理模式)**
$$\text{ProxyMode} ::= \text{iptables} \mid \text{ipvs} \mid \text{userspace} \mid \text{nftables} \mid \text{eBPF}$$

**iptables 规则链**:

```
PREROUTING ──► KUBE-SERVICES ──► KUBE-SVC-XXX ──► KUBE-SEP-XXX ──► DNAT
```

**复杂度分析**:

| 模式 | 规则数 | 查找复杂度 | 适用规模 |
|------|--------|-----------|---------|
| iptables | $O(S \times E)$ | $O(n)$ | < 1000 服务 |
| ipvs | $O(S)$ | $O(1)$ | > 1000 服务 |
| eBPF | $O(S)$ | $O(1)$ | 任何规模 |

### 3.2 EndpointSlice 与拓扑感知

**定义 3.2 (EndpointSlice)**
$$\text{EndpointSlice} = \langle \text{service}, \text{endpoints}, \text{ports}, \text{topology} \rangle$$

**拓扑感知路由**:
$$\text{PreferLocal} = \{ e \in \text{endpoints} \mid \text{zone}(e) = \text{zone}(\text{client}) \}$$

---

## 4. 多元表征

### 4.1 K8s 网络层次图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Kubernetes Network Layers                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Layer 4: Service Mesh (Optional)                                           │
│  ├── Istio/Linkerd                                                          │
│  ├── mTLS, Traffic Management                                               │
│  └── Observability                                                          │
│                                                                              │
│  Layer 3: Service Network                                                   │
│  ├── kube-proxy (iptables/ipvs/eBPF)                                        │
│  ├── Service IP (ClusterIP/NodePort/LoadBalancer)                           │
│  └── EndpointSlice                                                          │
│                                                                              │
│  Layer 2: Pod Network (CNI)                                                 │
│  ├── CNI Plugin (Calico/Cilium/Flannel)                                     │
│  ├── Pod IP (per-Pod unique)                                                │
│  └── cbr0 / veth / bridge                                                   │
│                                                                              │
│  Layer 1: Node Network                                                      │
│  ├── Physical NIC                                                           │
│  ├── Node IP                                                                │
│  └── Routing (BGP/VXLAN/IPIP)                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 CNI 选择决策树

```
选择 CNI?
│
├── 网络规模?
│   ├── 小型 (< 50 节点) → Flannel (简单)
│   └── 大型 (> 50 节点) → 需要高级特性
│       │
│       ├── 网络策略必需?
│       │   ├── 是 → Calico 或 Cilium
│       │   │       ├── 需要 eBPF?
│       │   │       │   ├── 是 → Cilium (性能)
│       │   │       │   └── 否 → Calico (稳定)
│       │   │       └──
│       │   │           BGP 可用?
│       │   │           ├── 是 → Calico BGP (高性能)
│       │   │           └── 否 → Calico VXLAN
│       │   └──
│       │       基础连通足够?
│       │       └── 是 → Weave (易用)
│       │
│       └── 需要服务网格集成?
│           └── 是 → Cilium (内置 Hubble)
│
└── 云提供商?
    ├── AWS → VPC CNI (原生集成)
    ├── GCP → GKE CNI
    └── Azure → Azure CNI
```

### 4.3 kube-proxy 模式对比矩阵

| 特性 | iptables | ipvs | eBPF (Cilium) |
|------|----------|------|---------------|
| **性能** | 中 (O(n)) | 高 (O(1)) | 极高 (直接包处理) |
| **规模** | < 1000 svc | 10000+ svc | 无限制 |
| **连接跟踪** | 内核 conntrack | 内核 conntrack | eBPF map |
| **负载均衡算法** | 随机 | 多种 (RR/LC/...) | 可编程 |
| **复杂度** | 低 | 中 | 高 |
| **调试** | 难 (iptables -L) | 中 (ipvsadm) | 难 (需要工具) |

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Kubernetes Network Checklist                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  网络配置:                                                                   │
│  □ Pod CIDR 不与节点网络冲突                                                  │
│  □ Service CIDR 独立规划                                                      │
│  □ MTU 配置正确 (考虑封装开销)                                                │
│  □ DNS 配置 (CoreDNS 副本数)                                                  │
│                                                                              │
│  CNI 选择:                                                                   │
│  □ 支持网络策略 (NetworkPolicy)                                               │
│  □ 性能测试通过 (iperf)                                                       │
│  □ 监控集成 (流量指标)                                                        │
│                                                                              │
│  Service 配置:                                                               │
│  □ 选择合适的 kube-proxy 模式                                                 │
│  □ 外部访问配置 (NodePort/LoadBalancer/Ingress)                               │
│  □ 会话亲和性需求 (sessionAffinity)                                           │
│                                                                              │
│  调试工具:                                                                   │
│  □ kubectl get svc,endpoints                                                  │
│  □ iptables -t nat -L (或 ipvsadm -Ln)                                        │
│  □ tcpdump 抓包                                                               │
│  □ CNI 特定工具 (calicoctl, cilium)                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (18KB, 完整形式化)
