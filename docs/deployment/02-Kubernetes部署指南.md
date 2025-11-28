# Kubernetes 部署指南

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 目录

- [1. Kubernetes 资源定义](#1-kubernetes-资源定义)
- [2. 部署配置](#2-部署配置)
- [3. 服务发现](#3-服务发现)
- [4. 自动扩展](#4-自动扩展)
- [5. 配置管理](#5-配置管理)
- [6. 最佳实践](#6-最佳实践)

---

## 1. Kubernetes 资源定义

### 1.1 Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
  namespace: default
  labels:
    app: app
    version: v1.0.0
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: app
  template:
    metadata:
      labels:
        app: app
        version: v1.0.0
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: app-service-account
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
      - name: app
        image: app:latest
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: redis-url
        - name: LOG_LEVEL
          value: "info"
        - name: PORT
          value: "8080"
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        startupProbe:
          httpGet:
            path: /startup
            port: 8080
          initialDelaySeconds: 0
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 30
        volumeMounts:
        - name: config
          mountPath: /etc/app/config
          readOnly: true
      volumes:
      - name: config
        configMap:
          name: app-config
```

### 1.2 Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: app-service
  namespace: default
  labels:
    app: app
spec:
  type: ClusterIP
  selector:
    app: app
  ports:
  - name: http
    port: 80
    targetPort: 8080
    protocol: TCP
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10800
```

### 1.3 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  namespace: default
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  tls:
  - hosts:
    - app.example.com
    secretName: app-tls
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
```

---

## 2. 部署配置

### 2.1 ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
data:
  redis-url: "redis://redis-service:6379/0"
  kafka-brokers: "kafka-service:9092"
  otlp-endpoint: "http://otel-collector:4317"
  log-level: "info"
  port: "8080"
```

### 2.2 Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
  namespace: default
type: Opaque
stringData:
  url: "postgres://user:password@postgres-service:5432/dbname?sslmode=disable"
```

### 2.3 ServiceAccount

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-service-account
  namespace: default
```

---

## 3. 服务发现

### 3.1 DNS 服务发现

Kubernetes 自动为 Service 创建 DNS 记录：

```text
<service-name>.<namespace>.svc.cluster.local
```

示例：

```text
app-service.default.svc.cluster.local
```

### 3.2 环境变量服务发现

Kubernetes 自动注入环境变量：

```bash
APP_SERVICE_HOST=10.0.0.1
APP_SERVICE_PORT=80
APP_SERVICE_PORT_80_TCP=tcp://10.0.0.1:80
```

---

## 4. 自动扩展

### 4.1 HorizontalPodAutoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: app-hpa
  namespace: default
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: app
  minReplicas: 3
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
        periodSeconds: 15
      - type: Pods
        value: 2
        periodSeconds: 15
      selectPolicy: Max
```

### 4.2 VerticalPodAutoscaler

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: app-vpa
  namespace: default
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: app
  updatePolicy:
    updateMode: "Auto"
  resourcePolicy:
    containerPolicies:
    - containerName: app
      minAllowed:
        cpu: 100m
        memory: 128Mi
      maxAllowed:
        cpu: 2
        memory: 2Gi
```

---

## 5. 配置管理

### 5.1 使用 ConfigMap

```yaml
# 从文件创建 ConfigMap
kubectl create configmap app-config \
  --from-file=config.yaml=./config/config.yaml \
  --from-literal=log-level=info

# 在 Pod 中使用
envFrom:
- configMapRef:
    name: app-config
```

### 5.2 使用 Secret

```yaml
# 从文件创建 Secret
kubectl create secret generic db-secret \
  --from-literal=url=postgres://user:pass@host:5432/db

# 在 Pod 中使用
env:
- name: DATABASE_URL
  valueFrom:
    secretKeyRef:
      name: db-secret
      key: url
```

### 5.3 使用 Helm

```yaml
# values.yaml
replicaCount: 3
image:
  repository: app
  tag: latest
  pullPolicy: IfNotPresent
service:
  type: ClusterIP
  port: 80
ingress:
  enabled: true
  hosts:
    - host: app.example.com
      paths: ["/"]
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

---

## 6. 最佳实践

### 6.1 资源管理

1. **设置资源请求和限制**：确保资源分配
2. **使用 QoS 类**：Guaranteed > Burstable > BestEffort
3. **监控资源使用**：使用 Prometheus 和 Grafana
4. **配置 HPA**：自动扩展
5. **使用 VPA**：自动调整资源

### 6.2 健康检查

1. **配置 Liveness Probe**：检测容器是否存活
2. **配置 Readiness Probe**：检测容器是否就绪
3. **配置 Startup Probe**：检测容器是否启动完成
4. **合理设置超时时间**：避免误判
5. **使用 HTTP 健康检查**：更准确

### 6.3 安全最佳实践

1. **使用 ServiceAccount**：最小权限原则
2. **使用 SecurityContext**：非 root 用户运行
3. **使用 NetworkPolicy**：网络隔离
4. **使用 PodSecurityPolicy**：Pod 安全策略
5. **扫描镜像漏洞**：使用 Trivy 等工具

### 6.4 部署策略

1. **使用 RollingUpdate**：零停机部署
2. **配置 maxSurge 和 maxUnavailable**：控制更新速度
3. **使用蓝绿部署**：快速回滚
4. **使用金丝雀部署**：灰度发布
5. **配置资源版本**：版本管理

---

**最后更新**: 2025-01-XX
