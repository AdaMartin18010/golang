# 快速开始指南

## 🚀 一键启动

```bash
# 1. 进入目录
cd deployments/docker

# 2. 启动所有服务
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 查看日志
docker-compose logs -f
```

## 📊 访问服务

| 服务 | 地址 | 说明 |
|------|------|------|
| 应用 | http://localhost:8080 | 主应用服务 |
| Temporal UI | http://localhost:8088 | 工作流管理 |
| Grafana | http://localhost:3000 | 监控面板 (admin/admin) |
| Prometheus | http://localhost:9090 | 指标查询 |
| Jaeger UI | http://localhost:16686 | 分布式追踪 |

## 🗄️ 数据库连接

### PostgreSQL 主节点（写操作）
```bash
psql -h localhost -p 5432 -U user -d mydb
# 密码: password
```

### PostgreSQL 备节点（只读）
```bash
psql -h localhost -p 5433 -U user -d mydb
# 密码: password
```

### Redis
```bash
redis-cli -h localhost -p 6379 -a redispassword
```

## 🔍 常用命令

### 检查复制状态
```bash
./postgresql/check-replication.sh
```

### 检查所有服务
```bash
./scripts/check-services.sh
```

### 数据库备份
```bash
# 备份数据库（自动压缩，保留7天）
./postgresql/backup.sh

# 恢复数据库
./postgresql/restore.sh ./backups/mydb_20240101_120000.sql.gz
```

### 数据库维护
```bash
# 执行维护（VACUUM、ANALYZE）
./postgresql/maintenance.sh

# 性能监控
./postgresql/performance.sh
```

### 查看特定服务日志
```bash
docker-compose logs -f db          # PostgreSQL 主节点
docker-compose logs -f db-replica  # PostgreSQL 备节点
docker-compose logs -f app         # 应用服务
docker-compose logs -f redis       # Redis
```

### 重启服务
```bash
docker-compose restart <service-name>
```

### 停止所有服务
```bash
docker-compose down
```

### 停止并删除数据
```bash
docker-compose down -v
```

## 🐛 故障排查

### 服务无法启动
```bash
# 查看服务日志
docker-compose logs <service-name>

# 检查服务状态
docker-compose ps
```

### PostgreSQL 复制问题
```bash
# 检查主节点复制状态
docker-compose exec db psql -U user -d mydb -c "SELECT * FROM pg_stat_replication;"

# 检查备节点状态
docker-compose exec db-replica psql -U user -d mydb -c "SELECT pg_is_in_recovery();"
```

### 端口冲突
如果端口被占用，可以修改 `docker-compose.yml` 中的端口映射。

## 📚 更多信息

详细文档请参考 [README.md](./README.md)
