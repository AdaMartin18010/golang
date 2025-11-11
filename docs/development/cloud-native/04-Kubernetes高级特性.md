# Kubernetes高级特性

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.23+ / Kubernetes 1.28+

---

## 📋 目录

- [Kubernetes高级特性](#kubernetes高级特性)
  - [📋 目录](#-目录)
  - [📋 目录](#-目录-1)
  - [4.1 📚 自定义资源(CRD)](#41--自定义资源crd)
  - [4.2 🤖 Operator模式](#42--operator模式)
  - [4.3 📊 高级调度](#43--高级调度)
    - [亲和性与反亲和性](#亲和性与反亲和性)
    - [Taint和Toleration](#taint和toleration)
    - [优先级与抢占](#优先级与抢占)
  - [4.4 💾 存储管理](#44--存储管理)
    - [StatefulSet](#statefulset)
    - [动态存储供应](#动态存储供应)
    - [卷快照](#卷快照)
  - [4.5 🌐 网络高级特性](#45--网络高级特性)
    - [Service Topology](#service-topology)
    - [Endpoint Slices](#endpoint-slices)
  - [4.6 🔧 集群管理](#46--集群管理)
    - [集群备份与恢复](#集群备份与恢复)
    - [集群升级](#集群升级)
  - [4.7 📈 监控与日志](#47--监控与日志)
    - [Prometheus Operator](#prometheus-operator)
    - [日志聚合](#日志聚合)
  - [4.8 🎯 最佳实践](#48--最佳实践)
  - [4.9 ⚠️ 常见问题](#49-️-常见问题)
    - [Q1: CRD与ConfigMap的区别？](#q1-crd与configmap的区别)
    - [Q2: Operator如何实现幂等性？](#q2-operator如何实现幂等性)
    - [Q3: StatefulSet何时使用？](#q3-statefulset何时使用)
  - [4.10 📚 扩展阅读](#410--扩展阅读)
    - [官方文档](#官方文档)
    - [相关文档](#相关文档)

---

## 📋 目录

## 4.1 📚 自定义资源(CRD)

**Custom Resource Definition**: 扩展Kubernetes API，定义自定义资源类型。

**创建CRD**:

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: databases.example.com
spec:
  group: example.com
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              size:
                type: string
                enum: ["small", "medium", "large"]
              version:
                type: string
              replicas:
                type: integer
                minimum: 1
                maximum: 5
          status:
            type: object
            properties:
              state:
                type: string
              message:
                type: string
  scope: Namespaced
  names:
    plural: databases
    singular: database
    kind: Database
    shortNames:
    - db
```

**使用自定义资源**:

```yaml
apiVersion: example.com/v1
kind: Database
metadata:
  name: my-database
spec:
  size: medium
  version: "5.7"
  replicas: 3
```

**Go客户端操作CRD**:

```go
import (
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/client-go/dynamic"
)

type Database struct {
    Spec DatabaseSpec `json:"spec"`
}

type DatabaseSpec struct {
    Size     string `json:"size"`
    Version  string `json:"version"`
    Replicas int    `json:"replicas"`
}

func CreateDatabase(ctx Context.Context, client dynamic.Interface, db *Database) error {
    gvr := schema.GroupVersionResource{
        Group:    "example.com",
        Version:  "v1",
        Resource: "databases",
    }

    unstructuredObj := &unstructured.Unstructured{
        Object: map[string]interface{}{
            "apiVersion": "example.com/v1",
            "kind":       "Database",
            "metadata": map[string]interface{}{
                "name": "my-database",
            },
            "spec": db.Spec,
        },
    }

    _, err := client.Resource(gvr).Namespace("default").Create(ctx, unstructuredObj, metav1.CreateOptions{})
    return err
}
```

## 4.2 🤖 Operator模式

**Operator**: 使用CRD和控制器自动化管理复杂应用。

**使用Operator SDK创建**:

```bash
# 初始化项目
operator-sdk init --domain=example.com --repo=github.com/myorg/database-operator

# 创建API
operator-sdk create api --group=database --version=v1 --kind=Database --resource --controller

# 生成CRD
make generate
make manifests
```

**Controller实现**:

```go
package controllers

import (
    "Context"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"

    databasev1 "github.com/myorg/database-operator/api/v1"
)

type DatabaseReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

func (r *DatabaseReconciler) Reconcile(ctx Context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)

    // 获取Database资源
    database := &databasev1.Database{}
    if err := r.Get(ctx, req.NamespacedName, database); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // 业务逻辑: 创建/更新Deployment
    deployment := r.createDeployment(database)
    if err := r.Create(ctx, deployment); err != nil {
        log.Error(err, "Failed to create Deployment")
        return ctrl.Result{}, err
    }

    // 更新状态
    database.Status.State = "Ready"
    if err := r.Status().Update(ctx, database); err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}

func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&databasev1.Database{}).
        Owns(&appsv1.Deployment{}).
        Complete(r)
}
```

## 4.3 📊 高级调度

### 亲和性与反亲和性

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  template:
    spec:
      affinity:
        # Pod亲和性
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - cache
            topologyKey: kubernetes.io/hostname
        # Pod反亲和性
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - web-app
              topologyKey: topology.kubernetes.io/zone
        # 节点亲和性
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: node-type
                operator: In
                values:
                - high-memory
```

### Taint和Toleration

```yaml
# 给节点添加Taint
kubectl taint nodes node1 key=value:NoSchedule

# Pod容忍Taint
apiVersion: v1
kind: Pod
metadata:
  name: tolerant-pod
spec:
  tolerations:
  - key: "key"
    operator: "Equal"
    value: "value"
    effect: "NoSchedule"
  containers:
  - name: app
    image: nginx
```

### 优先级与抢占

```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: high-priority
value: 1000
globalDefault: false
description: "High priority class"

---
apiVersion: v1
kind: Pod
metadata:
  name: high-priority-pod
spec:
  priorityClassName: high-priority
  containers:
  - name: app
    image: nginx
```

## 4.4 💾 存储管理

### StatefulSet

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        ports:
        - containerPort: 3306
          name: mysql
        volumeMounts:
        - name: data
          mountPath: /var/lib/mysql
        env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-secret
              key: password
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: "fast-ssd"
      resources:
        requests:
          storage: 10Gi
```

### 动态存储供应

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3
  iops: "3000"
  throughput: "125"
  encrypted: "true"
reclaimPolicy: Retain
allowVolumeExpansion: true
volumeBindingMode: WaitForFirstConsumer
```

### 卷快照

```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: mysql-snapshot
spec:
  volumeSnapshotClassName: csi-snapclass
  source:
    persistentVolumeClaimName: mysql-data-mysql-0
```

## 4.5 🌐 网络高级特性

### Service Topology

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: my-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  topologyKeys:
  - "kubernetes.io/hostname"
  - "topology.kubernetes.io/zone"
  - "*"
```

### Endpoint Slices

```yaml
apiVersion: discovery.k8s.io/v1
kind: EndpointSlice
metadata:
  name: my-service-abc
  labels:
    kubernetes.io/service-name: my-service
addressType: IPv4
ports:
- name: http
  protocol: TCP
  port: 80
endpoints:
- addresses:
  - "10.1.2.3"
  conditions:
    ready: true
  hostname: pod-1
  nodeName: node-1
```

## 4.6 🔧 集群管理

### 集群备份与恢复

**etcd备份**:

```bash
# 备份
ETCDCTL_API=3 etcdctl snapshot save /backup/etcd-snapshot.db \
  --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key

# 验证
ETCDCTL_API=3 etcdctl snapshot status /backup/etcd-snapshot.db

# 恢复
ETCDCTL_API=3 etcdctl snapshot restore /backup/etcd-snapshot.db \
  --data-dir=/var/lib/etcd-restore
```

### 集群升级

```bash
# 升级控制平面
kubeadm upgrade plan
kubeadm upgrade apply v1.28.0

# 升级节点
kubeadm upgrade node

# 升级kubelet
apt-mark unhold kubelet kubectl
apt-get update && apt-get install -y kubelet=1.28.0-00 kubectl=1.28.0-00
apt-mark hold kubelet kubectl
systemctl daemon-reload
systemctl restart kubelet
```

## 4.7 📈 监控与日志

### Prometheus Operator

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
spec:
  replicas: 2
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      team: frontend
  resources:
    requests:
      memory: 400Mi
  storage:
    volumeClaimTemplate:
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 50Gi
```

### 日志聚合

**Loki Stack**:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: promtail-config
data:
  promtail.yaml: |
    server:
      http_listen_port: 9080
      grpc_listen_port: 0
    positions:
      filename: /tmp/positions.yaml
    clients:
      - url: http://loki:3100/loki/api/v1/push
    scrape_configs:
    - job_name: kubernetes-pods
      kubernetes_sd_configs:
      - role: pod
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        target_label: app
```

## 4.8 🎯 最佳实践

1. **使用CRD扩展功能**: 而非修改Kubernetes核心
2. **Operator自动化运维**: 减少人工干预
3. **合理调度**: 使用亲和性和污点优化资源分布
4. **StatefulSet管理有状态应用**: 确保数据持久化
5. **动态存储供应**: 简化存储管理
6. **定期备份etcd**: 确保集群可恢复
7. **渐进式升级**: 测试后再生产环境升级
8. **监控关键指标**: API Server、etcd、调度器
9. **日志集中管理**: 使用ELK或Loki Stack
10. **安全加固**: RBAC、Pod Security、网络策略

## 4.9 ⚠️ 常见问题

### Q1: CRD与ConfigMap的区别？

**A**:

- CRD：定义新的API资源类型，有完整的CRUD和validation
- ConfigMap：存储配置数据，无业务逻辑

### Q2: Operator如何实现幂等性？

**A**:

- 使用状态机模式
- 检查资源当前状态
- 只在必要时执行操作

### Q3: StatefulSet何时使用？

**A**:

- 需要稳定的网络标识
- 需要持久化存储
- 需要有序部署和扩缩容

## 4.10 📚 扩展阅读

### 官方文档

- [Kubernetes CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/)

### 相关文档

- [01-Go与容器化基础.md](./01-Go与容器化基础.md)
- [02-Dockerfile最佳实践.md](./02-Dockerfile最佳实践.md)
- [03-Go与Kubernetes入门.md](./03-Go与Kubernetes入门.md)

---

**文档维护者**: Go Documentation Team
**最后更新**: 2025-10-29
**文档状态**: 完成
**适用版本**: Kubernetes 1.27+, Go 1.21+
