# TS-007: Kubernetes 网络深度解析 (Kubernetes Networking Deep Dive)

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #kubernetes #networking #cni #service-mesh #iptables
> **权威来源**: [K8s Networking](https://kubernetes.io/docs/concepts/services-networking/), [CNI Spec](https://www.cni.dev/)
> **K8s 版本**: 1.34+

---

## K8s 网络模型

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Kubernetes Network Model                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  K8s 网络原则:                                                               │
│  1. 每个 Pod 有独立的 IP (Pod IP)                                            │
│  2. 所有 Pod 可以在任何节点上互相通信 (无需 NAT)                              │
│  3. 所有节点可以与所有 Pod 通信                                               │
│  4. Service IP 是虚拟的，仅在集群内部可路由                                    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         Node-1                                      │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │    │
│  │  │   Pod-A     │  │   Pod-B     │  │   Pod-C     │                 │    │
│  │  │ 10.244.1.2  │  │ 10.244.1.3  │  │ 10.244.1.4  │                 │    │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                 │    │
│  │         │                │                │                        │    │
│  │         └────────────────┴────────────────┘                        │    │
│  │                          │                                         │    │
│  │                    ┌─────┴─────┐                                   │    │
│  │                    │  cbr0    │  (网桥/虚拟接口)                     │    │
│  │                    │ 10.244.1.1/24 │                                │    │
│  │                    └─────┬─────┘                                   │    │
│  │                          │                                         │    │
│  │  ┌───────────────────────┴───────────────────────┐                │    │
│  │  │              eth0 (Node IP)                   │                │    │
│  │  │              192.168.1.10                     │                │    │
│  │  └───────────────────────┬───────────────────────┘                │    │
│  │                          │                                         │    │
│  └──────────────────────────┼─────────────────────────────────────────┘    │
│                             │                                               │
│  ┌──────────────────────────┼─────────────────────────────────────────┐    │
│  │                         Node-2                                      │    │
│  │  ┌─────────────┐         │          ┌─────────────┐                │    │
│  │  │   Pod-D     │◄────────┴─────────►│   Pod-E     │                │    │
│  │  │ 10.244.2.2  │   Direct Routing   │ 10.244.2.3  │                │    │
│  │  └─────────────┘   (Overlay/VPC)   └─────────────┘                │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## CNI (Container Network Interface)

### CNI 工作流程

```
Pod 创建流程:

kubelet ──► CRI (containerd) ──► 创建 Pause 容器
    │
    │ 调用 CNI 插件
    ▼
┌─────────────────────────────────────────────────────────────┐
│                      CNI Plugin                             │
│  1. 分配 IP (从 IPAM)                                        │
│  2. 创建 veth pair (eth0 <-> vethxxx)                       │
│  3. 配置网桥 (cbr0/cni0)                                     │
│  4. 设置路由表                                               │
│  5. 配置 iptables/eBPF 规则                                  │
└─────────────────────────────────────────────────────────────┘
```

### 主流 CNI 对比

| CNI | 模式 | 特点 | 适用场景 |
|-----|------|------|---------|
| Flannel | VXLAN/Host-GW | 简单、轻量 | 小型集群 |
| Calico | BGP/eBPF | 性能高、策略丰富 | 生产环境 |
| Cilium | eBPF | 安全、可观测 | 云原生安全 |
| Weave | VXLAN | 易用、加密 | 多租户 |

### Calico BGP 模式

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Calico BGP Mode                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  每个节点运行 BIRD BGP Daemon，与路由反射器或全互联                         │
│                                                                              │
│  Node-1 (192.168.1.10)                        Node-2 (192.168.1.11)         │
│  ┌─────────────────────────┐                  ┌─────────────────────────┐   │
│  │  BIRD                   │◄────BGP Peering──►│  BIRD                   │   │
│  │  10.244.1.0/24          │                  │  10.244.2.0/24          │   │
│  └─────────────────────────┘                  └─────────────────────────┘   │
│                                                                              │
│  路由表:                                                                     │
│  Node-1: 10.244.2.0/24 via 192.168.1.11                                    │
│  Node-2: 10.244.1.0/24 via 192.168.1.10                                    │
│                                                                              │
│  优势: 无 Overlay 开销，性能接近裸机                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Service 网络

### kube-proxy 模式演进

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      kube-proxy Modes                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. userspace (已废弃)                                                      │
│     - 早期实现，性能差                                                        │
│                                                                              │
│  2. iptables (默认)                                                          │
│     - NAT 规则转发                                                           │
│     - 规则数量 O(Services × Endpoints)                                       │
│     - 大集群性能问题                                                          │
│                                                                              │
│  3. ipvs (推荐，大型集群)                                                     │
│     - 内核负载均衡                                                           │
│     - 哈希表查找 O(1)                                                        │
│     - 支持多种调度算法 (rr, lc, dh, sh, sed, nq)                              │
│                                                                              │
│  4. nftables (K8s 1.33+ 实验性)                                              │
│     - iptables 替代                                                          │
│     - 更好性能和可维护性                                                       │
│                                                                              │
│  5. eBPF (Cilium)                                                            │
│     - 绕过 kube-proxy                                                        │
│     - 直接 Socket 负载均衡                                                   │
│     - 最佳性能                                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### iptables 规则示例

```bash
# 查看 Service NAT 规则
iptables -t nat -L KUBE-SERVICES -n | grep my-service

# 典型规则链
KUBE-SERVICES
  └── KUBE-SVC-XXX (匹配 ClusterIP)
        └── KUBE-SEP-XXX (Endpoints)
              └── DNAT 到 Pod IP

# 示例: Service 10.96.0.1:80 -> Pods 10.244.1.2:8080, 10.244.1.3:8080
-A KUBE-SERVICES -d 10.96.0.1/32 -p tcp -m tcp --dport 80 -j KUBE-SVC-XXX
-A KUBE-SVC-XXX -m statistic --mode random --probability 0.5 -j KUBE-SEP-1
-A KUBE-SVC-XXX -j KUBE-SEP-2
-A KUBE-SEP-1 -p tcp -j DNAT --to-destination 10.244.1.2:8080
-A KUBE-SEP-2 -p tcp -j DNAT --to-destination 10.244.1.3:8080
```

---

## DNS 与服务发现

### CoreDNS 架构

```
Pod ──► CoreDNS ──► etcd (K8s API)
         │
         ├──► Forward . to /etc/resolv.conf
         ├──► Cache
         └──► Prometheus metrics

# Pod DNS 配置
nameserver 10.96.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
options ndots:5

# 完全限定域名 (FQDN)
my-service.my-namespace.svc.cluster.local
```

---

## 网络策略 (Network Policy)

```yaml
# 只允许 frontend 访问 backend 的 8080 端口
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-policy
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - protocol: TCP
          port: 8080
```

---

## 参考文献

1. [Kubernetes Networking](https://kubernetes.io/docs/concepts/services-networking/)
2. [CNI Specification](https://www.cni.dev/docs/spec/)
3. [Calico Documentation](https://docs.tigera.io/)
4. [Cilium Documentation](https://docs.cilium.io/)
