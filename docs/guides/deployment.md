# 部署指南

## 本地开发

### 使用 Docker Compose

```bash
cd deployments/docker
docker-compose up -d
```

这将启动：
- PostgreSQL 数据库
- OpenTelemetry Collector
- Prometheus
- Grafana
- Jaeger

### 运行应用

```bash
go run ./cmd/server
```

## 生产部署

### Docker 部署

#### 构建镜像

```bash
docker build -f deployments/docker/Dockerfile -t golang-app:latest .
```

#### 运行容器

```bash
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=user \
  -e DB_PASSWORD=password \
  -e DB_NAME=golang \
  golang-app:latest
```

### Kubernetes 部署

#### 应用配置

```bash
kubectl apply -f deployments/kubernetes/deployment.yaml
kubectl apply -f deployments/kubernetes/service.yaml
```

#### 检查状态

```bash
kubectl get pods
kubectl get services
```

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

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Jaeger**: http://localhost:16686

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

