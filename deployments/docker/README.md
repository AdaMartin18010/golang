# Docker 部署配置

本目录包含完整的 Docker Compose 配置，用于本地开发和测试环境。

> 💡 **快速开始**: 查看 [QUICKSTART.md](./QUICKSTART.md) 获取快速启动指南

## 📋 服务列表

### 应用服务

- **app**: 主应用服务 (端口: 8080)
- **temporal-worker**: Temporal 工作流工作器

### 数据库服务

- **db**: PostgreSQL 16 主节点 (端口: 5432)
- **db-replica**: PostgreSQL 16 备节点 - 异步复制 (端口: 5433)
- **redis**: Redis 7.2 (端口: 6379)
- **temporal-db**: Temporal 专用 PostgreSQL 数据库

> **注意**: PostgreSQL 支持 JSON/JSONB 数据类型，可以作为文档数据库使用，无需单独的 MongoDB。

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
- **postgres-exporter**: PostgreSQL 主节点指标导出器 (端口: 9187)
- **postgres-replica-exporter**: PostgreSQL 备节点指标导出器 (端口: 9188)

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

### PostgreSQL 主节点

- **主机**: db
- **端口**: 5432
- **用户**: user
- **密码**: password
- **数据库**: mydb
- **复制用户**: replicator
- **复制密码**: replicatorpassword

### PostgreSQL 备节点（只读）

- **主机**: db-replica
- **端口**: 5433
- **用户**: user
- **密码**: password
- **数据库**: mydb
- **复制模式**: 异步流复制

### Redis

- **主机**: redis
- **端口**: 6379
- **密码**: redispassword
- **数据库**: 0

### Grafana

- **URL**: <http://localhost:3000>
- **用户**: admin
- **密码**: admin

## 🌐 服务访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 应用 | <http://localhost:8080> | 主应用服务 |
| Temporal UI | <http://localhost:8088> | 工作流管理界面 |
| Grafana | <http://localhost:3000> | 监控面板 |
| Prometheus | <http://localhost:9090> | 指标查询 |
| Jaeger UI | <http://localhost:16686> | 分布式追踪 |

## 📊 监控配置

### Prometheus 监控目标

- OpenTelemetry Collector: `otel-collector:8888`
- Redis: `redis-exporter:9121`
- PostgreSQL 主节点: `postgres-exporter:9187`
- PostgreSQL 备节点: `postgres-replica-exporter:9187`
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

- `docker-compose.yml`: Docker Compose 主配置文件（主备模式）
- `docker-compose.cluster.yml`: PostgreSQL 集群配置（Patroni + HAProxy）
- `Dockerfile`: 应用镜像构建文件
- `env.example`: 环境变量示例文件
- `otel-collector-config.yaml`: OpenTelemetry Collector 配置
- `prometheus.yaml`: Prometheus 监控配置
- `grafana-datasources.yaml`: Grafana 数据源配置
- `grafana-dashboards.yaml`: Grafana 仪表板配置

### PostgreSQL 配置目录 (`postgresql/`)

- `postgresql.conf`: 主节点配置
- `postgresql-replica.conf`: 备节点配置
- `pg_hba.conf`: 客户端认证配置
- `init-replication.sh`: 复制用户初始化脚本
- `setup-replica.sh`: 备节点初始化脚本
- `check-replication.sh`: 复制状态检查脚本
- `setup-replication-slot.sh`: 复制槽创建脚本
- `backup.sh`: 数据库备份脚本
- `restore.sh`: 数据库恢复脚本
- `maintenance.sh`: 数据库维护脚本（VACUUM、ANALYZE）
- `performance.sh`: 性能监控脚本

### HAProxy 配置目录 (`haproxy/`)

- `haproxy.cfg`: 负载均衡和故障转移配置（用于集群模式）

### 管理脚本目录 (`scripts/`)

- `check-services.sh`: 服务健康状态检查脚本

## 🗄️ PostgreSQL 高可用配置

### 主备复制模式（默认）

当前配置使用 PostgreSQL 流复制实现主备架构：

- **主节点 (db)**: 处理所有写操作
- **备节点 (db-replica)**: 异步复制，用于读操作和故障转移

**特点**:

- 异步流复制，性能影响小
- 自动故障检测
- 支持热备（只读查询）

### 集群模式（可选）

使用 `docker-compose.cluster.yml` 启动 PostgreSQL 集群：

```bash
docker-compose -f docker-compose.yml -f docker-compose.cluster.yml up -d
```

**集群架构**:

- **Patroni**: 高可用管理器，自动故障转移
- **etcd**: 分布式配置存储
- **HAProxy**: 负载均衡和读写分离
- **多节点**: 支持多个 PostgreSQL 节点

**集群特性**:

- 自动主从切换
- 读写分离
- 多节点负载均衡
- 零停机故障转移

### PostgreSQL 文档数据库功能

PostgreSQL 原生支持 JSON/JSONB 数据类型，可以作为文档数据库使用：

```sql
-- 创建包含 JSON 字段的表
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 插入 JSON 文档
INSERT INTO documents (data) VALUES
('{"name": "John", "age": 30, "tags": ["developer", "golang"]}');

-- 查询 JSON 字段
SELECT * FROM documents WHERE data->>'name' = 'John';
SELECT * FROM documents WHERE data @> '{"age": 30}';
SELECT * FROM documents WHERE data->'tags' ? 'golang';

-- 创建 GIN 索引加速查询
CREATE INDEX idx_data_gin ON documents USING GIN (data);
```

## 🛠️ 管理脚本

### 检查服务状态

```bash
# 检查所有服务的健康状态
./scripts/check-services.sh

# 或使用 Docker Compose
docker-compose ps
```

### 检查 PostgreSQL 复制状态

```bash
# 检查主备复制状态和延迟
./postgresql/check-replication.sh

# 或在容器内执行
docker-compose exec db psql -U user -d mydb -c "SELECT * FROM pg_stat_replication;"
docker-compose exec db-replica psql -U user -d mydb -c "SELECT pg_is_in_recovery(), pg_last_wal_receive_lsn(), pg_last_wal_replay_lsn();"
```

### 创建复制槽

```bash
# 在主节点上创建复制槽（用于逻辑复制）
./postgresql/setup-replication-slot.sh
```

### 数据库备份和恢复

```bash
# 备份数据库
./postgresql/backup.sh

# 恢复数据库
./postgresql/restore.sh ./backups/mydb_20240101_120000.sql.gz
```

备份文件会自动压缩并保存在 `./backups/` 目录，默认保留最近 7 天的备份。

### 数据库维护

```bash
# 执行数据库维护（VACUUM、ANALYZE 等）
./postgresql/maintenance.sh

# 性能监控
./postgresql/performance.sh
```

### 查看服务日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f db
docker-compose logs -f db-replica
docker-compose logs -f app
```

## 🛠️ 开发建议

1. **数据持久化**: 所有数据库数据都存储在 Docker 卷中，重启不会丢失数据
2. **健康检查**: 所有服务都配置了健康检查，确保服务正常启动
3. **依赖管理**: 使用 `depends_on` 和健康检查条件确保服务启动顺序
4. **网络隔离**: 所有服务在 `app-network` 网络中，可以相互访问
5. **环境变量**: 可以复制 `env.example` 为 `.env` 来自定义配置

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
