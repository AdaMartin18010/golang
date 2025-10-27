# Kubernetes部署

**难度**: 高级 | **预计阅读**: 25分钟 | **前置知识**: Kubernetes基础、Docker

---

## 📋 目录

- [1. 📖 概念介绍](#1--概念介绍)
- [2. 🎯 核心资源](#2--核心资源)
  - [2.1 Deployment](#21-deployment)
  - [2.2 Service](#22-service)
  - [2.3 ConfigMap](#23-configmap)
  - [2.4 Secret](#24-secret)
  - [2.5 Ingress](#25-ingress)
- [3. 💡 最佳实践](#3--最佳实践)
- [4. 🔧 Helm部署](#4--helm部署)
- [5. 📚 相关资源](#5--相关资源)

---

## 1. 📖 概念介绍

Kubernetes（K8s）是容器编排平台，用于自动化部署、扩展和管理容器化应用。Go应用非常适合K8s部署。

---

## 2. 🎯 核心资源

### 2.1 Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

部署：
```bash
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl get pods
```

---

### 2.2 Service

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: LoadBalancer
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

---

### 2.3 ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
data:
  app.conf: |
    log_level=info
    max_connections=100
  feature.flags: |
    feature_x=true
    feature_y=false
```

使用：
```yaml
# 在Deployment中引用
spec:
  containers:
  - name: myapp
    envFrom:
    - configMapRef:
        name: myapp-config
```

---

### 2.4 Secret

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
stringData:
  url: "postgres://user:pass@db:5432/myapp"
  password: "secretpassword"
```

创建：
```bash
kubectl create secret generic db-secret \
  --from-literal=url='postgres://...' \
  --from-literal=password='secret'
```

---

### 2.5 Ingress

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: myapp-service
            port:
              number: 80
```

---

## 3. 💡 最佳实践

### 3.1 资源限制

```yaml
resources:
  requests:  # 调度所需最小资源
    memory: "64Mi"
    cpu: "250m"
  limits:    # 最大可用资源
    memory: "128Mi"
    cpu: "500m"
```

### 3.2 健康检查

```yaml
livenessProbe:   # 存活探针
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  
readinessProbe:  # 就绪探针
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
```

### 3.3 滚动更新

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1        # 最多多出1个Pod
    maxUnavailable: 0   # 最多0个不可用
```

### 3.4 水平扩展

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## 4. 🔧 Helm部署

### 4.1 Chart结构

```
myapp/
├── Chart.yaml
├── values.yaml
└── templates/
    ├── deployment.yaml
    ├── service.yaml
    └── ingress.yaml
```

### 4.2 Chart.yaml

```yaml
apiVersion: v2
name: myapp
version: 1.0.0
appVersion: "1.0.0"
```

### 4.3 values.yaml

```yaml
replicaCount: 3

image:
  repository: myapp
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  type: LoadBalancer
  port: 80
```

### 4.4 使用Helm

```bash
# 安装
helm install myapp ./myapp

# 升级
helm upgrade myapp ./myapp

# 回滚
helm rollback myapp

# 卸载
helm uninstall myapp
```

---

## 5. 📚 相关资源

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Helm](https://helm.sh/)

**下一步**: [04-CI-CD流程](./04-CI-CD流程.md)

---

**最后更新**: 2025-10-28
