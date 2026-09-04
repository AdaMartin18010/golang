# 部署脚本

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 脚本列表

| 脚本 | 说明 | 用法 |
|------|------|------|
| `docker-build.sh` | 构建 Docker 镜像 | `./scripts/deploy/docker-build.sh` |
| `docker-push.sh` | 推送 Docker 镜像 | `./scripts/deploy/docker-push.sh` |
| `k8s-deploy.sh` | 部署到 Kubernetes | `./scripts/deploy/k8s-deploy.sh` |
| `k8s-delete.sh` | 删除 Kubernetes 部署 | `./scripts/deploy/k8s-delete.sh` |

---

## 🚀 使用方法

### Docker 构建

```bash
# 使用默认配置
./scripts/deploy/docker-build.sh

# 自定义配置
IMAGE_NAME=myapp IMAGE_TAG=v1.0.0 ./scripts/deploy/docker-build.sh
```

### Docker 推送

```bash
# 设置环境变量
export REGISTRY_USER=your-username
export IMAGE_NAME=app
export IMAGE_TAG=latest

# 推送镜像
./scripts/deploy/docker-push.sh
```

### Kubernetes 部署

```bash
# 使用默认命名空间
./scripts/deploy/k8s-deploy.sh

# 指定命名空间
NAMESPACE=production ./scripts/deploy/k8s-deploy.sh
```

### Kubernetes 删除

```bash
# 删除部署
./scripts/deploy/k8s-delete.sh

# 指定命名空间
NAMESPACE=production ./scripts/deploy/k8s-delete.sh
```

---

## 📚 相关文档

- Docker 部署指南
- Kubernetes 部署指南
- 部署架构与策略

---

**最后更新**: 2025-01-XX
