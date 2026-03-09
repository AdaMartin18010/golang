# Docker 部署指南

> **版本**: v1.0
> **日期**: 2025-01-XX

---

## 📋 目录

- [1. Dockerfile 编写](#1-dockerfile-编写)
- [2. Docker Compose 配置](#2-docker-compose-配置)
- [3. HAProxy 配置](#3-haproxy-配置)
- [4. 多阶段构建](#4-多阶段构建)
- [5. 最佳实践](#5-最佳实践)

---

## 1. Dockerfile 编写

### 1.1 基础 Dockerfile

```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/server

# 运行阶段
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/app .

# 更改文件所有者
RUN chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./app"]
```

### 1.2 优化后的 Dockerfile

```dockerfile
# 多阶段构建 - 依赖阶段
FROM golang:1.21-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# 多阶段构建 - 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app ./cmd/server

# 多阶段构建 - 运行阶段
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/app /app
ENV TZ=Asia/Shanghai
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app", "-health-check"]
CMD ["/app"]
```

---

## 2. Docker Compose 配置

### 2.1 完整配置示例

```yaml
version: '3.8'

services:
  # 应用服务
  app:
    build:
      context: .
      dockerfile: Dockerfile
    image: app:latest
    container_name: app
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/dbname?sslmode=disable
      - REDIS_URL=redis://cache:6379/0
      - KAFKA_BROKERS=kafka:9092
      - OTLP_ENDPOINT=http://otel-collector:4317
      - LOG_LEVEL=info
      - PORT=8080
    depends_on:
      db:
        condition: service_healthy
      cache:
        condition: service_healthy
      kafka:
        condition: service_started
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # 负载均衡器
  haproxy:
    image: haproxy:2.8-alpine
    container_name: haproxy
    ports:
      - "80:80"
      - "8404:8404"  # Stats
    volumes:
      - ./deployments/docker/haproxy/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg:ro
    depends_on:
      - app
    networks:
      - app-network
    restart: unless-stopped

  # PostgreSQL 数据库
  db:
    image: postgres:15-alpine
    container_name: db
    environment:
      - POSTGRES_DB=dbname
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - PGDATA=/var/lib/postgresql/data/pgdata
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d dbname"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 缓存
  cache:
    image: redis:7-alpine
    container_name: cache
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Kafka 消息队列
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - app-network
    restart: unless-stopped

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    container_name: kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    networks:
      - app-network
    restart: unless-stopped

  # OpenTelemetry Collector
  otel-collector:
    image: otel/opentelemetry-collector:latest
    container_name: otel-collector
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./deployments/docker/otel/otel-collector-config.yaml:/etc/otel-collector-config.yaml:ro
    ports:
      - "4317:4317"   # OTLP gRPC receiver
      - "4318:4318"   # OTLP HTTP receiver
      - "8888:8888"   # Prometheus metrics
    networks:
      - app-network
    restart: unless-stopped

  # Prometheus
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    volumes:
      - ./deployments/docker/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
    ports:
      - "9090:9090"
    networks:
      - app-network
    restart: unless-stopped

  # Grafana
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    volumes:
      - grafana_data:/var/lib/grafana
      - ./deployments/docker/grafana/provisioning:/etc/grafana/provisioning:ro
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    networks:
      - app-network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:

networks:
  app-network:
    driver: bridge
```

### 2.2 开发环境配置

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - go-mod-cache:/go/pkg/mod
    environment:
      - GO_ENV=development
      - AIR_WORKSPACE=/app
    command: air -c .air.toml

volumes:
  go-mod-cache:
```

---

## 3. HAProxy 配置

### 3.1 完整配置

```cfg
global
    daemon
    maxconn 4096
    log stdout local0
    chroot /var/lib/haproxy
    stats socket /run/haproxy/admin.sock mode 660 level admin
    stats timeout 30s
    user haproxy
    group haproxy

defaults
    mode http
    log global
    option httplog
    option dontlognull
    option forwardfor
    option http-server-close
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms
    timeout http-request 10s
    timeout http-keep-alive 10s

# 统计页面
frontend stats
    bind *:8404
    stats enable
    stats uri /stats
    stats refresh 30s
    stats admin if TRUE

# HTTP 前端
frontend http_front
    bind *:80
    mode http

    # 健康检查路径直接返回
    acl is_health_check path_beg /health
    http-request return status 200 content-type "text/plain" string "OK" if is_health_check

    # 默认后端
    default_backend http_back

# HTTP 后端
backend http_back
    mode http
    balance roundrobin

    # 健康检查
    option httpchk GET /health
    http-check expect status 200

    # 服务器配置
    server app1 app:8080 check inter 5s fall 3 rise 2
    server app2 app:8080 check inter 5s fall 3 rise 2
    server app3 app:8080 check inter 5s fall 3 rise 2

    # 连接限制
    stick-table type ip size 100k expire 30s store http_req_rate(10s)
    http-request track-sc0 src
    http-request deny deny_status 429 if { sc_http_req_rate(0) gt 100 }
```

### 3.2 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `maxconn` | 最大连接数 | 4096 |
| `timeout connect` | 连接超时 | 5000ms |
| `timeout client` | 客户端超时 | 50000ms |
| `timeout server` | 服务器超时 | 50000ms |
| `balance` | 负载均衡算法 | roundrobin |
| `inter` | 健康检查间隔 | 5s |
| `fall` | 失败次数 | 3 |
| `rise` | 成功次数 | 2 |

---

## 4. 多阶段构建

### 4.1 构建阶段优化

```dockerfile
# 阶段 1: 依赖下载
FROM golang:1.21-alpine AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# 阶段 2: 代码构建
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s' \
    -a -installsuffix cgo \
    -o app ./cmd/server

# 阶段 3: 运行环境
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
```

### 4.2 使用 Scratch 镜像

```dockerfile
# 最小化镜像（使用 scratch）
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/app /app
ENV TZ=Asia/Shanghai
EXPOSE 8080
CMD ["/app"]
```

---

## 5. 最佳实践

### 5.1 镜像优化

1. **使用多阶段构建**：减少镜像大小
2. **使用 Alpine 或 Scratch**：最小化基础镜像
3. **合并 RUN 命令**：减少镜像层数
4. **使用 .dockerignore**：排除不必要的文件
5. **缓存依赖**：优化构建速度

### 5.2 安全最佳实践

1. **使用非 root 用户**：降低安全风险
2. **扫描镜像漏洞**：使用 Trivy 等工具
3. **最小权限原则**：只暴露必要的端口
4. **使用 Secret**：管理敏感信息
5. **定期更新基础镜像**：修复安全漏洞

### 5.3 性能最佳实践

1. **健康检查**：确保服务可用性
2. **资源限制**：防止资源耗尽
3. **日志管理**：使用日志驱动
4. **网络优化**：使用自定义网络
5. **卷管理**：持久化数据

---

**最后更新**: 2025-01-XX
