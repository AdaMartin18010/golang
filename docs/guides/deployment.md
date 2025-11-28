# 部署指南

> **详细文档**：请参考 [部署文档索引](../deployment/README.md)

---

## 📋 目录

- [本地开发](#本地开发)
- [Docker 部署](#docker-部署)
- [Kubernetes 部署](#kubernetes-部署)
- [环境变量](#环境变量)
- [健康检查](#健康检查)
- [监控和日志](#监控和日志)

---

## 本地开发

### 使用 Docker Compose

```bash
cd deployments/docker
docker-compose up -d
```

这将启动：

- **应用服务** (app) - 端口 8080
- **负载均衡器** (haproxy) - 端口 80
- **PostgreSQL 数据库** (db) - 端口 5432
- **Redis 缓存** (cache) - 端口 6379
- **Kafka 消息队列** (kafka) - 端口 9092
- **OpenTelemetry Collector** (otel-collector) - 端口 4317
- **Prometheus** (prometheus) - 端口 9090
- **Grafana** (grafana) - 端口 3000

### 运行应用

```bash
# 方式 1: 直接运行
go run ./cmd/server

# 方式 2: 使用热重载（推荐）
make run-dev
```

---

## Docker 部署

### 构建镜像

```bash
# 从项目根目录构建
docker build -f deployments/docker/Dockerfile -t app:latest .
```

### 使用 Docker Compose 启动

```bash
cd deployments/docker
docker-compose up -d
```

### 检查服务状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app
```

### 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

**详细说明**：请参考 [Docker 部署指南](../deployment/01-Docker部署指南.md)

---

## Kubernetes 部署

### 前置要求

- Kubernetes 1.25+
- kubectl
- 已配置的 kubeconfig

### 部署步骤

#### 1. 创建 ConfigMap

```bash
kubectl apply -f deployments/kubernetes/configmap.yaml
```

#### 2. 创建 Secret

```bash
# 从示例文件创建（需要修改实际值）
kubectl create secret generic db-secret \
  --from-literal=url=postgres://user:password@postgres-service:5432/dbname?sslmode=disable
```

#### 3. 创建 Deployment

```bash
kubectl apply -f deployments/kubernetes/deployment.yaml
```

#### 4. 创建 Service

```bash
kubectl apply -f deployments/kubernetes/service.yaml
```

#### 5. 创建 HPA（可选）

```bash
kubectl apply -f deployments/kubernetes/hpa.yaml
```

#### 6. 检查状态

```bash
# 查看 Pod 状态
kubectl get pods -l app=app

# 查看 Service
kubectl get svc app-service

# 查看 HPA
kubectl get hpa app-hpa

# 查看日志
kubectl logs -l app=app -f
```

**详细说明**：请参考 [Kubernetes 部署指南](../deployment/02-Kubernetes部署指南.md)

## 环境变量

### 服务器配置

- `SERVER_HOST` - 服务器地址（默认: 0.0.0.0）
- `SERVER_PORT` - 服务器端口（默认: 8080）

### 数据库配置

- `DB_HOST` - 数据库主机（默认: localhost）
- `DB_PORT` - 数据库端口（默认: 5432）
- `DB_USER` - 数据库用户
- `DB_PASSWORD` - 数据库密码
- `DB_NAME` - 数据库名称
- `DB_SSLMODE` - SSL 模式（默认: disable）

### 可观测性配置

- `OTLP_ENDPOINT` - OTLP 端点（默认: localhost:4317）
- `OTLP_INSECURE` - 是否使用不安全连接（默认: true）

## 健康检查

应用提供健康检查端点：

```bash
curl http://localhost:8080/health
```

## 监控和日志

### 查看日志

```bash
# Docker
docker logs <container-id>

# Kubernetes
kubectl logs <pod-name>
```

### 访问监控面板

- **Grafana**: <http://localhost:3000> (admin/admin)
- **Prometheus**: <http://localhost:9090>
- **Jaeger**: <http://localhost:16686>

## 数据库迁移

### 使用 Ent 迁移

```bash
# 生成迁移
go run -mod=mod entgo.io/ent/cmd/ent migrate new <migration-name>

# 应用迁移
go run -mod=mod entgo.io/ent/cmd/ent migrate apply
```

### 使用 SQL 迁移

```bash
# 应用迁移
psql -h localhost -U user -d golang -f migrations/postgres/001_create_users.up.sql
```

## 性能优化

### 数据库连接池

在 `configs/config.yaml` 中配置：

```yaml
database:
  max_conns: 25
```

### 服务器超时

```yaml
server:
  read_timeout: 5s
  write_timeout: 10s
  idle_timeout: 120s
```

## 安全建议

1. **使用环境变量存储敏感信息**
2. **启用 HTTPS**
3. **配置 CORS 策略**
4. **使用数据库 SSL 连接**
5. **定期更新依赖**

## 故障排除

### 数据库连接失败

- 检查数据库是否运行
- 验证连接配置
- 检查网络连接

### 端口冲突

- 修改 `configs/config.yaml` 中的端口配置
- 或使用环境变量覆盖

### 内存不足

- 调整 Docker/Kubernetes 资源限制
- 优化应用配置
