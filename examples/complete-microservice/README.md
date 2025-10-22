# 🎯 完整微服务示例

> **版本**: v2.0.0  
> **类型**: 完整应用示例  
> **难度**: 中级

这是一个完整的微服务应用示例，展示了如何集成和使用所有核心模块。

---

## 📋 目录

- [🎯 完整微服务示例](#-完整微服务示例)
  - [📋 目录](#-目录)
  - [✨ 功能特性](#-功能特性)
    - [核心模块集成](#核心模块集成)
    - [应用特性](#应用特性)
  - [🏗️ 架构设计](#️-架构设计)
    - [系统架构](#系统架构)
    - [请求处理流程](#请求处理流程)
  - [🚀 快速开始](#-快速开始)
    - [前置要求](#前置要求)
    - [安装依赖](#安装依赖)
    - [运行应用](#运行应用)
    - [测试应用](#测试应用)
  - [📚 API文档](#-api文档)
    - [GET /health](#get-health)
    - [POST /api/process](#post-apiprocess)
    - [GET /metrics](#get-metrics)
  - [📊 性能指标](#-性能指标)
    - [基准测试](#基准测试)
    - [预期性能](#预期性能)
    - [资源优化](#资源优化)
  - [🎓 最佳实践](#-最佳实践)
    - [1. 可观测性集成](#1-可观测性集成)
    - [2. 内存管理](#2-内存管理)
    - [3. 并发控制](#3-并发控制)
    - [4. 优雅关闭](#4-优雅关闭)
    - [5. 错误处理](#5-错误处理)
  - [🔧 配置](#-配置)
    - [环境变量](#环境变量)
    - [配置文件](#配置文件)
  - [🧪 测试](#-测试)
    - [单元测试](#单元测试)
    - [集成测试](#集成测试)
    - [压力测试](#压力测试)
  - [📈 监控](#-监控)
    - [Prometheus集成](#prometheus集成)
    - [关键指标](#关键指标)
  - [🚀 部署](#-部署)
    - [Docker](#docker)
    - [Kubernetes](#kubernetes)
  - [📚 扩展阅读](#-扩展阅读)
  - [💡 提示](#-提示)

---

## ✨ 功能特性

### 核心模块集成

- ✅ **AI-Agent** - 智能代理系统
- ✅ **Concurrency** - 并发模式（Worker Pool, Rate Limiter）
- ✅ **HTTP/3** - 现代化HTTP服务器
- ✅ **Memory** - 内存管理（对象池）
- ✅ **Observability** - 完整可观测性（Tracing, Metrics, Logging）

### 应用特性

- ✅ 优雅启动和关闭
- ✅ 健康检查端点
- ✅ 指标导出（Prometheus格式）
- ✅ 分布式追踪
- ✅ 结构化日志
- ✅ 请求限流
- ✅ 对象池优化
- ✅ 并发处理

---

## 🏗️ 架构设计

### 系统架构

```text
┌─────────────────────────────────────────────┐
│           Microservice Application          │
├─────────────────────────────────────────────┤
│                                             │
│  ┌────────────┐  ┌────────────┐           │
│  │   HTTP     │  │  AI-Agent  │           │
│  │  Server    │  │   System   │           │
│  └────────────┘  └────────────┘           │
│         │               │                   │
│  ┌──────▼───────────────▼────┐            │
│  │    Request Handler         │            │
│  └────────────┬────────────────┘            │
│               │                             │
│  ┌────────────▼────────────┐               │
│  │  Concurrency Patterns   │               │
│  │  • Worker Pool          │               │
│  │  • Rate Limiter         │               │
│  └────────────┬────────────┘               │
│               │                             │
│  ┌────────────▼────────────┐               │
│  │   Memory Management     │               │
│  │   • Object Pool         │               │
│  └────────────┬────────────┘               │
│               │                             │
│  ┌────────────▼────────────┐               │
│  │    Observability        │               │
│  │  • Tracing              │               │
│  │  • Metrics              │               │
│  │  • Logging              │               │
│  └─────────────────────────┘               │
│                                             │
└─────────────────────────────────────────────┘
```

### 请求处理流程

```text
1. HTTP Request
   ↓
2. Tracing (Start Span)
   ↓
3. Rate Limiting Check
   ↓
4. Get Request from Pool
   ↓
5. Worker Pool Processing
   ↓
6. Metrics Recording
   ↓
7. Structured Logging
   ↓
8. Return Request to Pool
   ↓
9. HTTP Response
```

---

## 🚀 快速开始

### 前置要求

- Go 1.25.3+
- 端口 8080 可用

### 安装依赖

```bash
cd examples/complete-microservice
go mod tidy
```

### 运行应用

```bash
go run main.go
```

输出示例：

```text
2025-10-22T10:00:00+08:00 INFO Starting microservice application...
2025-10-22T10:00:00+08:00 INFO Server listening addr=:8080
2025-10-22T10:00:00+08:00 INFO Microservice started successfully
```

### 测试应用

**健康检查**:

```bash
curl http://localhost:8080/health
```

**处理请求**:

```bash
curl -X POST http://localhost:8080/api/process
```

**查看指标**:

```bash
curl http://localhost:8080/metrics
```

---

## 📚 API文档

### GET /health

健康检查端点。

**响应**:

```json
{
  "status": "healthy",
  "version": "v2.0.0"
}
```

### POST /api/process

业务处理端点，展示完整的功能集成。

**特性**:

- 分布式追踪
- 请求限流（100 req/s）
- 对象池管理
- 并发处理（Worker Pool）
- 指标记录
- 结构化日志

**响应**:

```json
{
  "request_id": "req-1729566000123456",
  "processed": 10,
  "duration": 0.123
}
```

### GET /metrics

Prometheus格式的指标导出。

**指标**:

- `requests_total` - 总请求数
- `request_duration_seconds` - 请求处理时长

---

## 📊 性能指标

### 基准测试

```bash
# 运行基准测试
go test -bench=. -benchmem

# 压力测试
hey -n 10000 -c 100 http://localhost:8080/api/process
```

### 预期性能

| 指标 | 数值 |
|------|------|
| 吞吐量 | 1000+ req/s |
| 平均延迟 | <100ms |
| P99延迟 | <500ms |
| 内存占用 | <100MB |
| CPU使用 | <50% (单核) |

### 资源优化

- **对象池**: 减少 80% 内存分配
- **Worker Pool**: 高效并发处理
- **Rate Limiter**: 防止过载
- **Graceful Shutdown**: 零停机时间

---

## 🎓 最佳实践

### 1. 可观测性集成

```go
// 为每个请求创建追踪Span
span, ctx := observability.StartSpan(ctx, "operation-name")
defer span.Finish()

// 使用上下文日志
observability.WithContext(ctx).Info("Processing", "key", "value")

// 记录关键指标
counter.Inc()
histogram.Observe(duration)
```

### 2. 内存管理

```go
// 使用对象池减少GC压力
req := pool.Get()
defer pool.Put(req)

// 处理请求
process(req)
```

### 3. 并发控制

```go
// 使用Worker Pool进行并发处理
jobs := make(chan Job, 100)
results := patterns.WorkerPool(ctx, workerCount, jobs)

// 使用Rate Limiter限流
if !limiter.Allow() {
    return ErrRateLimitExceeded
}
```

### 4. 优雅关闭

```go
// 监听系统信号
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

<-sigChan

// 带超时的优雅关闭
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

server.Shutdown(ctx)
```

### 5. 错误处理

```go
// 使用上下文记录错误
if err != nil {
    observability.WithContext(ctx).Error("Operation failed", 
        "error", err,
        "request_id", reqID)
    return err
}
```

---

## 🔧 配置

### 环境变量

```bash
# 服务器地址
export SERVER_ADDR=:8080

# Worker Pool大小
export WORKER_POOL_SIZE=5

# 日志级别 (DEBUG, INFO, WARN, ERROR)
export LOG_LEVEL=INFO

# Rate Limiter（请求/秒）
export RATE_LIMIT=100
```

### 配置文件

创建 `config.yaml`:

```yaml
server:
  addr: ":8080"
  read_timeout: 30s
  write_timeout: 30s

concurrency:
  worker_pool_size: 5
  rate_limit: 100

memory:
  pool_size: 1000

observability:
  log_level: "INFO"
  metrics_enabled: true
  tracing_enabled: true
```

---

## 🧪 测试

### 单元测试

```bash
go test -v ./...
```

### 集成测试

```bash
go test -v -tags=integration ./...
```

### 压力测试

```bash
# 使用hey工具
hey -n 10000 -c 100 -m POST http://localhost:8080/api/process

# 使用wrk工具
wrk -t10 -c100 -d30s http://localhost:8080/api/process
```

---

## 📈 监控

### Prometheus集成

1. 添加到 `prometheus.yml`:

    ```yaml
    scrape_configs:
    - job_name: 'microservice'
        static_configs:
        - targets: ['localhost:8080']
        metrics_path: '/metrics'
    ```

2. 重启Prometheus

3. 访问 Grafana 查看仪表盘

### 关键指标

- `requests_total` - 请求总数
- `request_duration_seconds` - 请求时长分布
- `go_goroutines` - Goroutine数量
- `go_memstats_alloc_bytes` - 内存分配

---

## 🚀 部署

### Docker

创建 `Dockerfile`:

```dockerfile
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o microservice main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/microservice .
EXPOSE 8080
CMD ["./microservice"]
```

构建和运行:

```bash
docker build -t microservice:v2.0.0 .
docker run -p 8080:8080 microservice:v2.0.0
```

### Kubernetes

创建 `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: microservice
spec:
  replicas: 3
  selector:
    matchLabels:
      app: microservice
  template:
    metadata:
      labels:
        app: microservice
    spec:
      containers:
      - name: microservice
        image: microservice:v2.0.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

---

## 📚 扩展阅读

- [完整文档](../../docs/README.md)
- [API文档](../../API_DOCUMENTATION.md)
- [性能优化指南](../../docs/07-性能优化/README.md)
- [微服务架构](../../docs/05-微服务架构/README.md)

---

## 💡 提示

- 在生产环境中使用环境变量或配置文件
- 启用TLS/SSL加密
- 实施认证和授权
- 添加请求ID追踪
- 配置日志轮转
- 设置监控告警

---

**示例愉快！** 🎉

如有问题，欢迎提Issue或参与Discussions。
