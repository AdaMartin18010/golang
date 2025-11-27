# Docker 部署配置

本目录包含完整的 Docker Compose 配置，用于本地开发和测试环境。

## 📋 服务列表

### 应用服务
- **app**: 主应用服务 (端口: 8080)
- **temporal-worker**: Temporal 工作流工作器

### 数据库服务
- **db**: PostgreSQL 16 (端口: 5432)
- **redis**: Redis 7.2 (端口: 6379)
- **mongodb**: MongoDB 7.0 (端口: 27017)
- **temporal-db**: Temporal 专用 PostgreSQL 数据库

### 工作流服务
- **temporal**: Temporal 工作流引擎 (端口: 7233)
- **temporal-ui**: Temporal Web UI (端口: 8088)

### 可观测性服务
- **otel-collector**: OpenTelemetry 收集器
  - OTLP gRPC: 4317
  - OTLP HTTP: 4318
  - Health Check: 13133
- **prometheus**: 指标收集 (端口: 9090)
- **grafana**: 可视化面板 (端口: 3000)
- **jaeger**: 分布式追踪 (端口: 16686)

### 监控导出器
- **redis-exporter**: Redis 指标导出器 (端口: 9121)
- **mongodb-exporter**: MongoDB 指标导出器 (端口: 9216)
- **postgres-exporter**: PostgreSQL 指标导出器 (端口: 9187)

## 🚀 快速开始

### 启动所有服务

```bash
cd deployments/docker
docker-compose up -d
```

### 查看服务状态

```bash
docker-compose ps
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f redis
```

### 停止服务

```bash
docker-compose down
```

### 停止并删除数据卷

```bash
docker-compose down -v
```

## 🔐 默认凭据

### PostgreSQL
- **主机**: db
- **端口**: 5432
- **用户**: user
- **密码**: password
- **数据库**: mydb

### Redis
- **主机**: redis
- **端口**: 6379
- **密码**: redispassword
- **数据库**: 0

### MongoDB
- **主机**: mongodb
- **端口**: 27017
- **用户**: admin
- **密码**: mongopassword
- **数据库**: mydb
- **认证数据库**: admin

### Grafana
- **URL**: http://localhost:3000
- **用户**: admin
- **密码**: admin

## 🌐 服务访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 应用 | http://localhost:8080 | 主应用服务 |
| Temporal UI | http://localhost:8088 | 工作流管理界面 |
| Grafana | http://localhost:3000 | 监控面板 |
| Prometheus | http://localhost:9090 | 指标查询 |
| Jaeger UI | http://localhost:16686 | 分布式追踪 |

## 📊 监控配置

### Prometheus 监控目标

- OpenTelemetry Collector: `otel-collector:8888`
- Redis: `redis-exporter:9121`
- MongoDB: `mongodb-exporter:9216`
- PostgreSQL: `postgres-exporter:9187`
- 应用: `app:8080`

### Grafana 数据源

Grafana 已自动配置 Prometheus 数据源，无需手动添加。

## 🔧 环境变量

应用服务支持以下环境变量：

### 数据库配置
- `DATABASE_HOST`: PostgreSQL 主机
- `DATABASE_PORT`: PostgreSQL 端口
- `DATABASE_USER`: PostgreSQL 用户
- `DATABASE_PASSWORD`: PostgreSQL 密码
- `DATABASE_DBNAME`: PostgreSQL 数据库名
- `DATABASE_SSLMODE`: SSL 模式

### Redis 配置
- `REDIS_HOST`: Redis 主机
- `REDIS_PORT`: Redis 端口
- `REDIS_PASSWORD`: Redis 密码
- `REDIS_DB`: Redis 数据库编号

### MongoDB 配置
- `MONGODB_URI`: MongoDB 连接字符串

### 可观测性配置
- `OTEL_EXPORTER_OTLP_ENDPOINT`: OpenTelemetry 导出端点

### Temporal 配置
- `TEMPORAL_ADDRESS`: Temporal 服务地址
- `TEMPORAL_TASK_QUEUE`: Temporal 任务队列名称

## 📝 配置文件

- `docker-compose.yml`: Docker Compose 主配置文件
- `Dockerfile`: 应用镜像构建文件
- `otel-collector-config.yaml`: OpenTelemetry Collector 配置
- `prometheus.yaml`: Prometheus 监控配置
- `grafana-datasources.yaml`: Grafana 数据源配置
- `grafana-dashboards.yaml`: Grafana 仪表板配置

## 🛠️ 开发建议

1. **数据持久化**: 所有数据库数据都存储在 Docker 卷中，重启不会丢失数据
2. **健康检查**: 所有服务都配置了健康检查，确保服务正常启动
3. **依赖管理**: 使用 `depends_on` 和健康检查条件确保服务启动顺序
4. **网络隔离**: 所有服务在 `app-network` 网络中，可以相互访问

## 🔍 故障排查

### 服务无法启动

```bash
# 检查服务日志
docker-compose logs <service-name>

# 检查服务健康状态
docker-compose ps
```

### 连接问题

确保服务在同一个网络中，并且依赖的服务已健康启动。

### 端口冲突

如果端口已被占用，可以在 `docker-compose.yml` 中修改端口映射。

## 📚 相关文档

- [Docker Compose 文档](https://docs.docker.com/compose/)
- [OpenTelemetry 文档](https://opentelemetry.io/docs/)
- [Temporal 文档](https://docs.temporal.io/)
- [Prometheus 文档](https://prometheus.io/docs/)
