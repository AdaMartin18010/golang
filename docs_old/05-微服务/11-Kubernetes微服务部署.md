# 11. ☸️ Kubernetes微服务部署

> 📚 **简介**：本文档深入探讨在Kubernetes平台上部署和管理Go微服务的完整流程，涵盖容器化、资源编排、服务发现、配置管理、健康检查、自动扩展和持续部署等核心主题。通过本文，读者将掌握在Kubernetes环境中构建生产级微服务的实践技能。

<!-- TOC START -->
- [11. ☸️ Kubernetes微服务部署](#11-️-kubernetes微服务部署)
  - [11.1 📚 Kubernetes基础](#111--kubernetes基础)
    - [核心概念](#核心概念)
    - [架构组件](#架构组件)
  - [11.2 🐳 容器化Go应用](#112--容器化go应用)
    - [Dockerfile最佳实践](#dockerfile最佳实践)
    - [多阶段构建](#多阶段构建)
    - [镜像优化](#镜像优化)
  - [11.3 📋 部署配置](#113--部署配置)
    - [Deployment](#deployment)
    - [Service](#service)
    - [Ingress](#ingress)
  - [11.4 ⚙️ 配置管理](#114-️-配置管理)
    - [ConfigMap](#configmap)
    - [Secret](#secret)
    - [环境变量注入](#环境变量注入)
  - [11.5 💾 存储管理](#115--存储管理)
    - [Volume](#volume)
    - [PersistentVolume](#persistentvolume)
  - [11.6 🔍 健康检查](#116--健康检查)
    - [Liveness Probe](#liveness-probe)
    - [Readiness Probe](#readiness-probe)
    - [Startup Probe](#startup-probe)
  - [11.7 📊 资源管理](#117--资源管理)
    - [资源请求与限制](#资源请求与限制)
    - [QoS类别](#qos类别)
  - [11.8 🚀 自动扩展](#118--自动扩展)
    - [HPA水平扩展](#hpa水平扩展)
    - [VPA垂直扩展](#vpa垂直扩展)
  - [11.9 🔄 滚动更新](#119--滚动更新)
    - [更新策略](#更新策略)
    - [回滚操作](#回滚操作)
  - [11.10 📈 监控与日志](#1110--监控与日志)
    - [Prometheus监控](#prometheus监控)
    - [日志收集](#日志收集)
  - [11.11 🎯 最佳实践](#1111--最佳实践)
  - [11.12 ⚠️ 常见问题](#1112-️-常见问题)
    - [Q1: Pod一直处于Pending状态？](#q1-pod一直处于pending状态)
    - [Q2: 如何调试CrashLoopBackOff？](#q2-如何调试crashloopbackoff)
    - [Q3: 如何优雅关闭应用？](#q3-如何优雅关闭应用)
    - [Q4: 如何实现蓝绿部署？](#q4-如何实现蓝绿部署)
  - [11.13 📚 扩展阅读](#1113--扩展阅读)
    - [官方文档](#官方文档)
    - [相关文档](#相关文档)
    - [工具推荐](#工具推荐)
<!-- TOC END -->

## 11.1 📚 Kubernetes基础

### 核心概念

**Pod**: Kubernetes最小部署单元，包含一个或多个容器。

**ReplicaSet**: 确保指定数量的Pod副本运行。

**Deployment**: 管理ReplicaSet，提供声明式更新。

**Service**: 为Pod提供稳定的网络访问入口。

**Namespace**: 资源隔离和多租户支持。

### 架构组件

```text
Master节点:
├── API Server    # 集群管理的统一入口
├── Scheduler     # 负责Pod调度
├── Controller    # 维护集群状态
└── etcd          # 集群数据存储

Worker节点:
├── Kubelet       # 节点代理，管理Pod
├── Kube-proxy    # 网络代理
└── Container Runtime  # 容器运行时（Docker/containerd）
```

## 11.2 🐳 容器化Go应用

### Dockerfile最佳实践

```dockerfile
# 基础镜像选择
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖（利用Docker缓存层）
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 最终镜像
FROM alpine:latest

# 安装ca证书（HTTPS请求需要）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
```

### 多阶段构建

**优势**:

- 减小镜像体积（去除编译工具）
- 提高安全性（减少攻击面）
- 加快部署速度

**示例**:

```dockerfile
# 阶段1: 构建
FROM golang:1.21 AS builder
WORKDIR /build
COPY . .
RUN go build -ldflags="-s -w" -o app .

# 阶段2: 运行
FROM gcr.io/distroless/base-debian11
COPY --from=builder /build/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
```

### 镜像优化

**体积对比**:

| 基础镜像 | 大小 | 特点 |
|----------|------|------|
| golang:1.21 | ~800MB | 完整开发环境 |
| alpine | ~5MB | 最小Linux发行版 |
| distroless | ~20MB | 无shell，高安全性 |
| scratch | 0MB | 空镜像，仅二进制 |

**推荐组合**:

```dockerfile
# 构建: golang:1.21-alpine
# 运行: alpine 或 distroless
```

## 11.3 📋 部署配置

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: production
  labels:
    app: user-service
    version: v1.0
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
        version: v1.0
    spec:
      containers:
      - name: user-service
        image: myregistry/user-service:v1.0
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        env:
        - name: ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: database.host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database.password
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: production
spec:
  type: ClusterIP
  selector:
    app: user-service
  ports:
  - name: http
    protocol: TCP
    port: 80
    targetPort: 8080
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800
```

**Service类型**:

| 类型 | 用途 | 访问方式 |
|------|------|----------|
| ClusterIP | 集群内部访问 | 内部IP |
| NodePort | 外部访问（测试） | NodeIP:Port |
| LoadBalancer | 外部访问（生产） | 云厂商LB |
| ExternalName | DNS映射 | CNAME |

### Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  namespace: production
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - api.example.com
    secretName: api-tls
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /users
        pathType: Prefix
        backend:
          service:
            name: user-service
            port:
              number: 80
      - path: /orders
        pathType: Prefix
        backend:
          service:
            name: order-service
            port:
              number: 80
```

## 11.4 ⚙️ 配置管理

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: production
data:
  app.properties: |
    server.port=8080
    log.level=info
  database.host: "mysql.default.svc.cluster.local"
  database.port: "3306"
  redis.host: "redis.default.svc.cluster.local"
```

**使用方式**:

```yaml
# 1. 环境变量
env:
- name: DB_HOST
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: database.host

# 2. 文件挂载
volumes:
- name: config
  configMap:
    name: app-config
volumeMounts:
- name: config
  mountPath: /etc/config
  readOnly: true
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: production
type: Opaque
data:
  database.password: cGFzc3dvcmQxMjM=  # base64编码
  api.key: YXBpa2V5MTIzNDU2  # base64编码
```

**创建Secret**:

```bash
# 从文件创建
kubectl create secret generic app-secrets \
  --from-file=./secret.txt \
  --namespace=production

# 从字面值创建
kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password=secret123 \
  --namespace=production

# 从环境文件创建
kubectl create secret generic app-secrets \
  --from-env-file=./secret.env \
  --namespace=production
```

### 环境变量注入

```go
package main

import (
    "os"
    "log"
)

type Config struct {
    Port       string
    DBHost     string
    DBPassword string
    RedisHost  string
}

func LoadConfig() *Config {
    return &Config{
        Port:       getEnv("PORT", "8080"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPassword: os.Getenv("DB_PASSWORD"), // 必需的
        RedisHost:  getEnv("REDIS_HOST", "localhost"),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func main() {
    config := LoadConfig()
    
    if config.DBPassword == "" {
        log.Fatal("DB_PASSWORD must be set")
    }
    
    log.Printf("Starting server on port %s", config.Port)
    // ...
}
```

## 11.5 💾 存储管理

### Volume

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-pod
spec:
  containers:
  - name: app
    image: myapp:latest
    volumeMounts:
    - name: cache
      mountPath: /app/cache
    - name: logs
      mountPath: /var/log/app
  volumes:
  - name: cache
    emptyDir: {}
  - name: logs
    hostPath:
      path: /var/log/pods
      type: DirectoryOrCreate
```

### PersistentVolume

```yaml
# PV定义
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-data
spec:
  capacity:
    storage: 10Gi
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: slow
  nfs:
    server: nfs-server.default.svc.cluster.local
    path: /data

---
# PVC申请
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: data-claim
  namespace: production
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: slow

---
# 使用PVC
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
      - name: app
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: data-claim
```

## 11.6 🔍 健康检查

### Liveness Probe

**存活探针**: 检测容器是否存活，失败则重启。

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
    httpHeaders:
    - name: Custom-Header
      value: Awesome
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  successThreshold: 1
  failureThreshold: 3
```

**Go实现**:

```go
func healthHandler(c *gin.Context) {
    // 检查关键依赖
    if err := checkDatabase(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "error":  err.Error(),
        })
        return
    }
    
    if err := checkRedis(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "error":  err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
    })
}
```

### Readiness Probe

**就绪探针**: 检测容器是否准备好接收流量，失败则移出Service。

```yaml
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
```

**Go实现**:

```go
var isReady atomic.Bool

func readyHandler(c *gin.Context) {
    if !isReady.Load() {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "not ready",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "ready",
    })
}

func initializeApp() {
    // 初始化数据库
    if err := setupDatabase(); err != nil {
        log.Fatal(err)
    }
    
    // 初始化缓存
    if err := setupCache(); err != nil {
        log.Fatal(err)
    }
    
    // 标记为就绪
    isReady.Store(true)
}
```

### Startup Probe

**启动探针**: 检测容器应用是否启动完成，适用于慢启动应用。

```yaml
startupProbe:
  httpGet:
    path: /startup
    port: 8080
  initialDelaySeconds: 0
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 30  # 最多等待300秒
```

## 11.7 📊 资源管理

### 资源请求与限制

```yaml
resources:
  requests:  # 最小保证资源
    cpu: "100m"      # 0.1核
    memory: "128Mi"  # 128MB
  limits:    # 最大使用资源
    cpu: "500m"      # 0.5核
    memory: "512Mi"  # 512MB
```

**CPU单位**:

- `1` = 1核心
- `100m` = 0.1核心（100毫核）

**内存单位**:

- `128Mi` = 128 MiB（1024^2）
- `1Gi` = 1 GiB

### QoS类别

| QoS类别 | 条件 | 驱逐优先级 |
|---------|------|-----------|
| Guaranteed | requests = limits | 最低 |
| Burstable | requests < limits | 中等 |
| BestEffort | 无requests/limits | 最高 |

## 11.8 🚀 自动扩展

### HPA水平扩展

**基于CPU扩展**:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: user-service-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 30
      - type: Pods
        value: 4
        periodSeconds: 30
      selectPolicy: Max
```

**基于自定义指标**:

```yaml
metrics:
- type: Pods
  pods:
    metric:
      name: http_requests_per_second
    target:
      type: AverageValue
      averageValue: "1000"
```

### VPA垂直扩展

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: user-service-vpa
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  updatePolicy:
    updateMode: "Auto"
  resourcePolicy:
    containerPolicies:
    - containerName: user-service
      minAllowed:
        cpu: 100m
        memory: 128Mi
      maxAllowed:
        cpu: 2
        memory: 2Gi
```

## 11.9 🔄 滚动更新

### 更新策略

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # 最多超出期望副本数1个
      maxUnavailable: 0  # 最多0个不可用
```

**更新流程**:

```bash
# 1. 更新镜像
kubectl set image deployment/user-service \
  user-service=myregistry/user-service:v1.1 \
  --namespace=production

# 2. 查看滚动状态
kubectl rollout status deployment/user-service \
  --namespace=production

# 3. 查看历史版本
kubectl rollout history deployment/user-service \
  --namespace=production
```

### 回滚操作

```bash
# 回滚到上一版本
kubectl rollout undo deployment/user-service \
  --namespace=production

# 回滚到指定版本
kubectl rollout undo deployment/user-service \
  --to-revision=2 \
  --namespace=production

# 暂停滚动
kubectl rollout pause deployment/user-service \
  --namespace=production

# 恢复滚动
kubectl rollout resume deployment/user-service \
  --namespace=production
```

## 11.10 📈 监控与日志

### Prometheus监控

**ServiceMonitor**:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: user-service-monitor
  namespace: production
spec:
  selector:
    matchLabels:
      app: user-service
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
```

**Go应用暴露指标**:

```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // 业务路由
    r := gin.Default()
    r.GET("/health", healthHandler)
    r.GET("/api/users", getUsersHandler)
    
    // 指标端点
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))
    
    r.Run(":8080")
}
```

### 日志收集

**结构化日志**:

```go
import "github.com/sirupsen/logrus"

func main() {
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{})
    log.SetOutput(os.Stdout)
    
    log.WithFields(logrus.Fields{
        "service":  "user-service",
        "version":  "v1.0",
        "pod_name": os.Getenv("POD_NAME"),
    }).Info("Service started")
}
```

**Fluent Bit配置**:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluent-bit-config
data:
  fluent-bit.conf: |
    [INPUT]
        Name              tail
        Path              /var/log/containers/*production*.log
        Parser            json
        Tag               kube.*
    
    [OUTPUT]
        Name              es
        Match             *
        Host              elasticsearch
        Port              9200
        Index             k8s-logs
```

## 11.11 🎯 最佳实践

1. **使用多阶段构建**: 减小镜像体积，提高安全性
2. **设置资源限制**: 防止资源滥用，提高稳定性
3. **配置健康检查**: 自动重启失败容器，提高可用性
4. **使用ConfigMap/Secret**: 分离配置和代码
5. **启用HPA**: 自动应对负载变化
6. **实施滚动更新**: 零停机部署
7. **配置就绪探针**: 避免将流量发送到未就绪Pod
8. **使用命名空间**: 资源隔离和权限管理
9. **标签和注解**: 便于资源管理和查询
10. **监控和日志**: 及时发现和排查问题

## 11.12 ⚠️ 常见问题

### Q1: Pod一直处于Pending状态？

**A**: 可能原因：

- 资源不足（CPU/内存）
- 无可用节点
- PVC绑定失败
- 镜像拉取失败

**排查**:

```bash
kubectl describe pod <pod-name>
kubectl get events --sort-by='.lastTimestamp'
```

### Q2: 如何调试CrashLoopBackOff？

**A**:

```bash
# 查看日志
kubectl logs <pod-name> --previous

# 进入容器
kubectl exec -it <pod-name> -- /bin/sh

# 查看事件
kubectl describe pod <pod-name>
```

### Q3: 如何优雅关闭应用？

**A**:

```go
func main() {
    srv := &http.Server{Addr: ":8080"}
    
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // 监听信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // 优雅关闭（5秒超时）
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exited")
}
```

### Q4: 如何实现蓝绿部署？

**A**: 使用两个Deployment和Service标签切换：

```bash
# 部署green版本
kubectl apply -f deployment-green.yaml

# 切换Service指向green
kubectl patch service user-service -p '{"spec":{"selector":{"version":"green"}}}'

# 删除blue版本
kubectl delete deployment user-service-blue
```

## 11.13 📚 扩展阅读

### 官方文档

- [Kubernetes文档](https://kubernetes.io/docs/)
- [kubectl命令参考](https://kubernetes.io/docs/reference/kubectl/)
- [Kubernetes API](https://kubernetes.io/docs/reference/kubernetes-api/)

### 相关文档

- [10-高性能微服务架构.md](./10-高性能微服务架构.md)
- [12-Service Mesh集成.md](./12-Service-Mesh集成.md)
- [../06-云原生/01-容器化部署.md](../06-云原生/01-容器化部署.md)

### 工具推荐

- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [k9s](https://k9scli.io/) - Kubernetes CLI工具
- [Lens](https://k8slens.dev/) - Kubernetes IDE
- [Helm](https://helm.sh/) - Kubernetes包管理器

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Kubernetes 1.27+, Go 1.21+
