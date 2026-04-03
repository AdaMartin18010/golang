# FT-036-Kubernetes-135-Formal-Model

> **Dimension**: 01-Formal-Theory
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: Kubernetes 1.35 "Timbernetes" (Released: December 17, 2025)
> **Size**: >20KB

---

## 1. Kubernetes 1.35 概览

### 1.1 发布信息

- **发布名称**: Timbernetes (The World Tree Release)
- **发布日期**: 2025年12月17日
- **发布负责人**: Drew Hagen
- **支持周期**: 约至2026年12月

### 1.2 发布统计

| 类别 | 数量 |
|------|------|
| 总增强 | 60 |
| Stable (GA) | 17 |
| Beta | 19 |
| Alpha | 22 |

**主题灵感**: 北欧神话世界树Yggdrasil——根深(稳定基础)、干强(核心特性)、枝茂(新能力)

---

## 2. 核心GA特性

### 2.1 In-Place Pod Resize (KEP-1287)

**状态**: GA ( graduated to Stable)

**形式化定义**:

```
设 Pod P = (C₁, C₂, ..., Cₙ) 为容器集合
设资源规范 R = (CPU, Memory)

传统重调度:
  resize(P, R') → 终止P → 创建P' → 启动C'

In-Place Resize:
  resize(P, R') → 调整cgroups → P保持运行

约束条件:
  ∀c ∈ C: c.state = Running → resize_possible(c)
  QoS(P) = QoS(P')  // QoS类不能改变
```

**使用**:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: resizable-pod
spec:
  containers:
  - name: app
    image: myapp:latest
    resources:
      requests:
        cpu: "500m"
        memory: "512Mi"
      limits:
        cpu: "1000m"
        memory: "1Gi"
```

```bash
# 热调整CPU
kubectl patch pod resizable-pod --subresource resize -p '{
  "spec": {
    "containers": [{
      "name": "app",
      "resources": {
        "requests": {"cpu": "700m"},
        "limits": {"cpu": "700m"}
      }
    }]
  }
}'
```

**限制**:

- 不能改变QoS类
- Init容器和临时容器不可调整
- Windows Pods不支持
- 静态CPU/内存管理策略的Pod不可调整

### 2.2 侧链容器增强

**形式化模型**:

```
Pod生命周期:
  Init Phase → Main Phase

传统Init容器:
  I₁ → I₂ → ... → Iₙ (顺序执行，全部完成后启动主容器)

侧链容器 (restartPolicy: Always):
  I_sidecar ∥ C_main  // 并行运行

容器重启规则 (KEP-4960, Beta):
  restartPolicyRules:
    - action: RestartContainer | RestartAllContainers
      exitCodes:
        operator: In | NotIn
        values: [0, 1, ...]
```

```yaml
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: watcher
    image: watcher:1.0
    restartPolicy: Always  # 标记为侧链
  containers:
  - name: main-app
    image: myapp:1.0
  restartPolicyRules:  # Beta
  - action: RestartAllContainers
    exitCodes:
      operator: In
      values: [88]  # watcher退出时重启整个Pod
```

### 2.3 Pod级别资源 (KEP-5419, Alpha)

```yaml
apiVersion: v1
kind: Pod
spec:
  resources:  # Pod级别资源
    requests:
      pod.kubernetes.io/memory: "64Gi"
      pod.kubernetes.io/cpu: "8"
  containers:
  - name: trainer
    image: ml-trainer:latest
    # 不指定资源限制，使用Pod级别
```

---

## 3. AI/ML工作负载增强

### 3.1 Gang调度 (KEP-4671, Alpha)

**问题**: 传统Kubernetes逐个调度Pod，导致分布式训练作业的"部分分配"死锁

**形式化定义**:

```
设 PodGroup G = {P₁, P₂, ..., Pₙ}
设资源需求 R = Σ resources(Pᵢ)

传统调度:
  ∀Pᵢ ∈ G: schedule(Pᵢ) 顺序执行
  可能结果: 部分Pod调度，部分Pending → 死锁

Gang调度:
  atomic_schedule(G):
    if available_resources ≥ R:
      schedule_all(G)
    else:
      queue_all(G)  // 全部等待
```

```yaml
apiVersion: scheduling.x-k8s.io/v1alpha1
kind: PodGroup
metadata:
  name: training-group
spec:
  scheduleTimeoutSeconds: 300
  minMember: 4  # 4个Pod必须一起调度
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    scheduling.x-k8s.io/pod-group: training-group
spec:
  containers:
  - name: trainer
    resources:
      limits:
        nvidia.com/gpu: 2
```

### 3.2 动态资源分配 (DRA)

**状态**: Feature gate locked to enabled-by-default

| 特性 | KEP | 状态 | 描述 |
|------|-----|------|------|
| Device Binding Conditions | #5007 | Beta | 硬件就绪前等待 |
| Partitionable Devices | #4815 | Alpha | 动态分区(MIG切片) |
| Prioritized Alternatives | #4816 | Beta | 备用设备请求 |
| Device Taints/Tolerations | #5055 | Alpha | 设备级调度约束 |
| Consumable Capacity | #5075 | Alpha | 保证资源共享 |

### 3.3 OCI镜像卷源 (KEP-4639, Beta)

```yaml
apiVersion: v1
kind: Pod
spec:
  volumes:
  - name: model-volume
    ociImage:
      reference: registry.example.com/models/v1:latest
      pullPolicy: IfNotPresent
  containers:
  - name: inference
    image: inference-server:latest
    volumeMounts:
    - name: model-volume
      mountPath: /models
```

**优势**:

- 模型与应用解耦
- 减小镜像体积
- 消除init container复杂性

---

## 4. 安全增强

### 4.1 用户命名空间 (Beta, 默认启用)

**KEP**: 127

```yaml
apiVersion: v1
kind: Pod
spec:
  hostUsers: false  // 启用用户命名空间
  containers:
  - name: app
    image: myapp:latest
```

**安全模型**:

```
容器UID 0 (root) → 映射到主机UID 100000+ (无特权)

效果:
- 容器逃逸后无root权限
- 显著降低容器逃逸风险
```

### 4.2 Pod证书 (KEP-4317, Beta)

**形式化定义**:

```
Pod身份 = X.509证书
签发者 = Kubelet (代表集群CA)
生命周期 = 自动管理 (创建、轮换、撤销)

优势:
- 原生工作负载身份
- 无需外部cert-manager
- 自动轮换
- mTLS简化
```

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    image: myapp:latest
    volumeMounts:
    - name: cert-volume
      mountPath: /certs
  volumes:
  - name: cert-volume
    projected:
      sources:
      - serviceAccountToken:
          path: token
      - certificate:
          path: tls.crt
      - certificate:
          path: tls.key
```

### 4.3 约束模拟 (KEP-5284, Alpha)

```yaml
# 限制模拟权限
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: constrained-impersonator
rules:
- apiGroups: [""]
  resources: ["users"]
  verbs: ["impersonate"]
  # 只能模拟特定用户
  resourceNames: ["service-account-1"]
```

---

## 5. 网络增强

### 5.1 PreferSameNode流量分发 (KEP-3015, GA)

```yaml
apiVersion: v1
kind: Service
spec:
  trafficDistribution:
    config:
      distributionPolicy: PreferSameNode  // 优先同节点
      // PreferSameZone 重命名为更清晰
```

**路由策略对比**:

| 策略 | 描述 |
|------|------|
| PreferSameNode | 优先同节点Pod (最低延迟) |
| PreferSameZone | 优先同可用区Pod |
| ClusterDefault | 集群默认策略 |

### 5.2 WebSockets替换SPDY (KEP-4006, Beta)

```
旧协议: SPDY (已弃用)
新协议: WebSockets

影响命令:
- kubectl exec
- kubectl attach
- kubectl port-forward
- kubectl cp
```

---

## 6. 存储增强

### 6.1 设备绑定条件 (KEP-5007, Beta)

```yaml
apiVersion: resource.k8s.io/v1alpha3
kind: ResourceClaim
spec:
  devices:
    config:
    - conditions:
      - type: Ready  # 等待设备就绪
        status: "True"
```

### 6.2 可变CSI节点可分配 (KEP-4876, Beta)

动态卷附着容量管理。

---

## 7. 调度器增强

### 7.1 机会性批处理调度 (KEP-5598, Beta)

```yaml
apiVersion: batch/v1
kind: Job
spec:
  parallelism: 1000
  completions: 1000
  schedulingPolicy:
    batchScheduling: Enabled  # 批量调度
```

**效果**: 显著降低大规模部署的调度器开销。

### 7.2 VPA集成In-Place Resize

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
spec:
  updatePolicy:
    updateMode: "InPlaceOrRecreate"  # Beta
    # 先尝试原地调整，失败则驱逐重建
```

---

## 8. 弃用和移除

### 8.1 cgroup v1支持移除 (KEP-5573)

**影响**: Kubelet在cgroup v1节点上默认启动失败

**迁移**:

```bash
# 检查cgroup版本
stat -fc %T /sys/fs/cgroup/
# cgroup2fs = v2 (OK)
# tmpfs = v1 (需要迁移)
```

### 8.2 Ingress NGINX进入维护模式

- 完全退役: 2026年3月
- **迁移路径**: Gateway API

### 8.3 kube-proxy ipvs模式弃用

**KEP**: 5495
**替代**: nftables

---

## 9. 形式化验证模型

### 9.1 In-Place Resize正确性

```tla
--------------------------- MODULE PodResize ----------------------------
EXTENDS Integers, Sequences, FiniteSets

CONSTANTS Pods, Resources

VARIABLES podStates, resourceSpecs

Init ==
  podStates = [p ∈ Pods |-> "Running"]
  resourceSpecs = [p ∈ Pods |-> defaultResources]

Resize(p, newSpec) ==
  ∧ podStates[p] = "Running"
  ∧ CanResizeInPlace(p, newSpec)
  ∧ resourceSpecs' = [resourceSpecs EXCEPT ![p] = newSpec]
  ∧ podStates' = [podStates EXCEPT ![p] = "Resizing"]

CompleteResize(p) ==
  ∧ podStates[p] = "Resizing"
  ∧ podStates' = [podStates EXCEPT ![p] = "Running"]
  ∧ UNCHANGED resourceSpecs

Safety ==
  ∀p ∈ Pods:
    podStates[p] = "Running" ⇒
      ActualResources(p) ≤ resourceSpecs[p].limits

Liveness ==
  ∀p ∈ Pods, spec ∈ Resources:
    ◇(resizeRequested(p, spec) ⇒ ◇(resourceSpecs[p] = spec))
============================================================================
```

---

## 10. 升级指南

### 10.1 升级到1.35

```bash
# 1. 检查弃用API
kubectl get --raw=/apis | jq '.groups[].versions[].groupVersion'

# 2. 验证cgroup v2
stat -fc %T /sys/fs/cgroup/

# 3. 滚动升级
kubeadm upgrade plan
kubeadm upgrade apply v1.35.0

# 4. 验证
kubectl version
kubectl get nodes
```

### 10.2 关键变更清单

| 变更 | 行动 |
|------|------|
| cgroup v1移除 | 迁移到v2 |
| Ingress NGINX退役 | 迁移到Gateway API |
| 旧API移除 | 更新manifests |

---

## 11. 参考文献

1. Kubernetes 1.35 Release Notes
2. KEP-1287: In-Place Pod Resize
3. KEP-4671: Gang Scheduling
4. KEP-4317: Pod Certificates
5. KEP-127: User Namespaces
6. Timbernetes Release Blog

---

*Last Updated: 2026-04-03*
