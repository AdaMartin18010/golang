# Kubernetes

> **分类**: 工程与云原生

---

## Go 客户端

```go
import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

config, _ := rest.InClusterConfig()
clientset, _ := kubernetes.NewForConfig(config)

pods, _ := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
```

---

## 健康检查

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
```

---

## 资源管理

```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```
