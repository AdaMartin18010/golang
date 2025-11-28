# 部署配置

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 目录结构

```
deployments/
├── docker/              # Docker 部署配置
│   ├── Dockerfile      # 多阶段构建 Dockerfile
│   ├── docker-compose.yml  # Docker Compose 配置
│   ├── haproxy/        # HAProxy 配置
│   └── README.md       # Docker 部署说明
│
└── kubernetes/         # Kubernetes 部署配置
    ├── deployment.yaml # Deployment 资源
    ├── service.yaml    # Service 资源
    ├── hpa.yaml        # HorizontalPodAutoscaler
    ├── configmap.yaml  # ConfigMap
    ├── secret.yaml.example  # Secret 示例
    └── README.md       # Kubernetes 部署说明
```

---

## 🚀 快速开始

### Docker 部署

```bash
cd deployments/docker
docker-compose up -d
```

详细说明请参考：[Docker 部署 README](docker/README.md)

### Kubernetes 部署

```bash
cd deployments/kubernetes
kubectl apply -f .
```

详细说明请参考：[Kubernetes 部署 README](kubernetes/README.md)

---

## 📚 相关文档

- [部署架构与策略](../docs/deployment/00-部署架构与策略.md)
- [Docker 部署指南](../docs/deployment/01-Docker部署指南.md)
- [Kubernetes 部署指南](../docs/deployment/02-Kubernetes部署指南.md)
- [部署文档索引](../docs/deployment/README.md)

---

**最后更新**: 2025-01-XX
