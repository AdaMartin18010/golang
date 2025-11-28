# Docker 部署配置

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 目录结构

```
deployments/docker/
├── Dockerfile              # 多阶段构建 Dockerfile
├── docker-compose.yml      # Docker Compose 完整配置
├── .dockerignore          # Docker 构建忽略文件
├── haproxy/
│   └── haproxy.cfg        # HAProxy 负载均衡配置
├── otel/
│   └── otel-collector-config.yaml  # OpenTelemetry Collector 配置
├── prometheus/
│   └── prometheus.yml     # Prometheus 配置
└── grafana/
    └── provisioning/      # Grafana 配置
```

---

## 🚀 快速开始

### 1. 构建镜像

```bash
# 从项目根目录构建
docker build -f deployments/docker/Dockerfile -t app:latest .
```

### 2. 使用 Docker Compose 启动

```bash
cd deployments/docker
docker-compose up -d
```

这将启动以下服务：
- **应用服务** (app) - 端口 8080
- **负载均衡器** (haproxy) - 端口 80
- **PostgreSQL 数据库** (db) - 端口 5432
- **Redis 缓存** (cache) - 端口 6379
- **Kafka 消息队列** (kafka) - 端口 9092
- **OpenTelemetry Collector** (otel-collector) - 端口 4317
- **Prometheus** (prometheus) - 端口 9090
- **Grafana** (grafana) - 端口 3000

### 3. 检查服务状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app

# 查看所有日志
docker-compose logs -f
```

### 4. 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

---

## 📝 配置说明

### 环境变量

主要环境变量配置在 `docker-compose.yml` 中：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DATABASE_URL` | 数据库连接字符串 | `postgres://user:password@db:5432/dbname?sslmode=disable` |
| `REDIS_URL` | Redis 连接字符串 | `redis://cache:6379/0` |
| `KAFKA_BROKERS` | Kafka Broker 地址 | `kafka:9092` |
| `OTLP_ENDPOINT` | OTLP 端点地址 | `http://otel-collector:4317` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `PORT` | 服务端口 | `8080` |

### 端口映射

| 服务 | 容器端口 | 主机端口 | 说明 |
|------|---------|---------|------|
| app | 8080 | 8080 | 应用服务 |
| haproxy | 80 | 80 | 负载均衡器 |
| haproxy | 8404 | 8404 | HAProxy 统计页面 |
| prometheus | 9090 | 9090 | Prometheus 监控 |
| grafana | 3000 | 3000 | Grafana 可视化 |

---

## 🔧 高级配置

### 自定义配置

1. **修改环境变量**：编辑 `docker-compose.yml` 中的 `environment` 部分
2. **修改端口映射**：编辑 `docker-compose.yml` 中的 `ports` 部分
3. **修改资源限制**：编辑 `docker-compose.yml` 中的 `deploy.resources` 部分

### 数据持久化

数据卷配置在 `docker-compose.yml` 中：

- `postgres_data` - PostgreSQL 数据
- `redis_data` - Redis 数据
- `prometheus_data` - Prometheus 数据
- `grafana_data` - Grafana 数据

### 网络配置

所有服务都在 `app-network` 网络中，可以通过服务名相互访问。

---

## 📚 相关文档

- [Docker 部署指南](../../docs/deployment/01-Docker部署指南.md)
- [部署架构与策略](../../docs/deployment/00-部署架构与策略.md)
- [Kubernetes 部署指南](../../docs/deployment/02-Kubernetes部署指南.md)

---

## 🐛 故障排查

### 常见问题

1. **端口被占用**
   ```bash
   # 检查端口占用
   netstat -tulpn | grep :8080
   # 或修改 docker-compose.yml 中的端口映射
   ```

2. **服务启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs app
   # 检查服务依赖
   docker-compose ps
   ```

3. **数据库连接失败**
   ```bash
   # 检查数据库服务
   docker-compose ps db
   # 查看数据库日志
   docker-compose logs db
   ```

---

**最后更新**: 2025-01-XX
