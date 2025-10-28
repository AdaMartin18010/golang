# Kubernetes部署Go应用

> **简介**: Kubernetes部署Go应用完整指南，包括Deployment、Service、配置管理和最佳实践

> **版本**: Go 1.25.3, Kubernetes 1.28+  
> **难度**: ⭐⭐⭐⭐  
> **标签**: #Kubernetes #K8s #部署 #云原生

---

## 📚 目录

1. [基础概念](#基础概念)
2. [Deployment部署](#deployment部署)
3. [Service服务](#service服务)
4. [配置管理](#配置管理)
5. [健康检查](#健康检查)
6. [最佳实践](#最佳实践)

---

## 1. 基础概念

### Kubernetes核心对象

- **Pod**: 最小部署单元，包含一个或多个容器
- **Deployment**: 管理Pod的副本和更新
- **Service**: 为Pod提供稳定的网络访问
- **ConfigMap**: 配置数据
- **Secret**: 敏感数据
- **Ingress**: HTTP/HTTPS路由

---

## 2. Deployment部署

### 基本Deployment

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
        - name: ENV
          value: "production"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

**部署**:
```bash
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl get pods
```

---

### 滚动更新

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # 最多额外创建1个Pod
      maxUnavailable: 0  # 更新时保持可用
  template:
    # ... (同上)
```

**更新应用**:
```bash
# 更新镜像
kubectl set image deployment/myapp myapp=myapp:2.0.0

# 查看更新状态
kubectl rollout status deployment/myapp

# 回滚
kubectl rollout undo deployment/myapp

# 查看历史
kubectl rollout history deployment/myapp
```

---

### 水平自动扩缩容（HPA）

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
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

**应用**:
```bash
kubectl apply -f hpa.yaml
kubectl get hpa
```

---

## 3. Service服务

### ClusterIP（内部访问）

```yaml
# service-clusterip.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: ClusterIP
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

---

### NodePort（外部访问）

```yaml
# service-nodeport.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-nodeport
spec:
  type: NodePort
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
    nodePort: 30080  # 30000-32767
```

---

### LoadBalancer（云负载均衡）

```yaml
# service-lb.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-lb
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

### Ingress（HTTP路由）

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
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
  tls:
  - hosts:
    - myapp.example.com
    secretName: myapp-tls
```

---

## 4. 配置管理

### ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
data:
  app.yaml: |
    server:
      port: 8080
      host: 0.0.0.0
    database:
      host: postgres
      port: 5432
  LOG_LEVEL: "info"
  ENV: "production"
```

**在Deployment中使用**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:1.0.0
        # 环境变量
        env:
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: myapp-config
              key: LOG_LEVEL
        # 挂载文件
        volumeMounts:
        - name: config
          mountPath: /etc/config
      volumes:
      - name: config
        configMap:
          name: myapp-config
```

---

### Secret

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secret
type: Opaque
data:
  # base64编码
  DB_PASSWORD: cGFzc3dvcmQxMjM=
  API_KEY: YXBpa2V5MTIz
```

**创建Secret**:
```bash
# 从字面值创建
kubectl create secret generic myapp-secret \
  --from-literal=DB_PASSWORD=password123 \
  --from-literal=API_KEY=apikey123

# 从文件创建
kubectl create secret generic myapp-secret \
  --from-file=./credentials.txt
```

**使用Secret**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:1.0.0
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: myapp-secret
              key: DB_PASSWORD
```

---

## 5. 健康检查

### Liveness Probe（存活探针）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:1.0.0
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3
```

---

### Readiness Probe（就绪探针）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:1.0.0
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          successThreshold: 1
          failureThreshold: 3
```

**Go应用中实现**:
```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
    // 检查依赖服务是否就绪
    if !checkDatabase() {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Ready"))
}

func main() {
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/ready", readyHandler)
    http.ListenAndServe(":8080", nil)
}
```

---

### Startup Probe（启动探针）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:1.0.0
        startupProbe:
          httpGet:
            path: /health
            port: 8080
          failureThreshold: 30
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          periodSeconds: 10
```

---

## 6. 最佳实践

### 1. 完整部署配置

```yaml
# complete-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
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
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
        version: v1.0.0
    spec:
      # 优雅关闭
      terminationGracePeriodSeconds: 30
      
      # 安全上下文
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      
      containers:
      - name: myapp
        image: myapp:1.0.0
        imagePullPolicy: IfNotPresent
        
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        
        # 环境变量
        env:
        - name: ENV
          value: production
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: myapp-secret
              key: DB_PASSWORD
        
        # 资源限制
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        
        # 健康检查
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
          initialDelaySeconds: 10
          periodSeconds: 5
        
        # 挂载配置
        volumeMounts:
        - name: config
          mountPath: /etc/config
          readOnly: true
      
      volumes:
      - name: config
        configMap:
          name: myapp-config
```

---

### 2. 使用Kustomize

```
k8s/
├── base/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── kustomization.yaml
├── overlays/
│   ├── dev/
│   │   ├── kustomization.yaml
│   │   └── replica-patch.yaml
│   └── prod/
│       ├── kustomization.yaml
│       └── replica-patch.yaml
```

**base/kustomization.yaml**:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- deployment.yaml
- service.yaml

commonLabels:
  app: myapp
```

**overlays/prod/kustomization.yaml**:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- ../../base

replicas:
- name: myapp
  count: 5

images:
- name: myapp
  newTag: v1.0.0
```

**部署**:
```bash
# 开发环境
kubectl apply -k k8s/overlays/dev

# 生产环境
kubectl apply -k k8s/overlays/prod
```

---

### 3. 使用Helm

```bash
# 创建Helm chart
helm create myapp

# 目录结构
myapp/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   └── _helpers.tpl
```

**values.yaml**:
```yaml
replicaCount: 3

image:
  repository: myapp
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  hosts:
    - host: myapp.example.com
      paths:
        - path: /
          pathType: Prefix

resources:
  limits:
    cpu: 500m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

**部署**:
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

### 4. 多环境配置

```bash
# values-dev.yaml
replicaCount: 1
image:
  tag: dev

# values-prod.yaml
replicaCount: 5
image:
  tag: 1.0.0
resources:
  limits:
    cpu: 1000m
    memory: 512Mi
```

**部署不同环境**:
```bash
# 开发环境
helm install myapp ./myapp -f values-dev.yaml

# 生产环境
helm install myapp ./myapp -f values-prod.yaml
```

---

### 5. 监控和日志

```yaml
# ServiceMonitor for Prometheus
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: myapp
spec:
  selector:
    matchLabels:
      app: myapp
  endpoints:
  - port: metrics
    interval: 30s
```

**日志收集**:
```yaml
# 使用fluentd收集日志
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/myapp*.log
      pos_file /var/log/fluentd-myapp.pos
      tag kubernetes.*
      format json
    </source>
```

---

## 🎯 常用命令

```bash
# 创建资源
kubectl apply -f deployment.yaml

# 查看资源
kubectl get deployments
kubectl get pods
kubectl get services
kubectl describe pod <pod-name>

# 日志
kubectl logs <pod-name>
kubectl logs -f <pod-name>  # 实时日志

# 执行命令
kubectl exec -it <pod-name> -- /bin/sh

# 端口转发
kubectl port-forward <pod-name> 8080:8080

# 扩缩容
kubectl scale deployment myapp --replicas=5

# 删除资源
kubectl delete -f deployment.yaml
kubectl delete deployment myapp
```

---

## 🔗 相关资源

- [Docker部署](./02-Docker部署.md)
- [CI/CD流程](./04-CI-CD流程.md)
- [监控与日志](./05-监控与日志.md)

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3  
**Kubernetes版本**: 1.28+
