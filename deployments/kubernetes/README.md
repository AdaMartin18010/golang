# Kubernetes 部署配置

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 文件说明

| 文件 | 说明 |
|------|------|
| `deployment.yaml` | Deployment 资源定义 |
| `service.yaml` | Service 资源定义 |
| `hpa.yaml` | HorizontalPodAutoscaler 资源定义 |
| `configmap.yaml` | ConfigMap 资源定义 |
| `secret.yaml.example` | Secret 资源示例（需要创建真实的 Secret） |
| `ingress.yaml.example` | Ingress 资源示例（需要根据实际情况修改） |

---

## 🚀 部署步骤

### 1. 创建命名空间（可选）

```bash
kubectl create namespace app
```

### 2. 创建 ConfigMap

```bash
kubectl apply -f configmap.yaml
```

### 3. 创建 Secret

```bash
# 从示例文件创建（需要修改实际值）
kubectl create secret generic db-secret \
  --from-literal=url=postgres://user:password@postgres-service:5432/dbname?sslmode=disable

# 或从文件创建
kubectl create secret generic db-secret \
  --from-file=url=./secret-url.txt
```

### 4. 创建 Deployment

```bash
kubectl apply -f deployment.yaml
```

### 5. 创建 Service

```bash
kubectl apply -f service.yaml
```

### 6. 创建 HPA（可选）

```bash
kubectl apply -f hpa.yaml
```

### 7. 创建 Ingress（可选，需要 Ingress Controller）

```bash
# 修改 ingress.yaml.example 中的域名和配置
# 然后应用配置
kubectl apply -f ingress.yaml
```

### 8. 检查部署状态

```bash
# 查看 Pod 状态
kubectl get pods -l app=app

# 查看 Service
kubectl get svc app-service

# 查看 HPA
kubectl get hpa app-hpa

# 查看 Ingress
kubectl get ingress app-ingress

# 查看日志
kubectl logs -l app=app -f
```

---

## 📚 相关文档

- Kubernetes 部署指南
- 部署架构与策略

---

**最后更新**: 2025-01-XX
