# 运维控制功能完整指南

> **版本**: v1.0.0
> **Go版本**: 1.25+

完整的运维控制功能，提供生产环境必需的所有运维控制能力。

## 📋 功能清单

### ✅ 运维端点（12 个）

1. **健康检查** (`/ops/health`, `/ops/healthz`)
   - 检查应用健康状态
   - 返回详细的健康信息（CPU、内存、Goroutine 等）
   - Kubernetes 兼容

2. **就绪检查** (`/ops/ready`, `/ops/readiness`)
   - 检查应用是否准备好接收流量
   - 检查关键依赖是否就绪
   - Kubernetes 就绪探针

3. **存活检查** (`/ops/live`, `/ops/liveness`)
   - 检查应用是否存活
   - Kubernetes 存活探针

4. **指标导出** (`/ops/metrics`)
   - JSON 格式指标导出
   - 包含所有系统指标

5. **Prometheus 指标** (`/ops/metrics/prometheus`)
   - Prometheus 格式指标导出
   - 可直接被 Prometheus 抓取

6. **仪表板数据** (`/ops/dashboard`)
   - 完整的仪表板数据
   - JSON 格式

7. **诊断报告** (`/ops/diagnostics`)
   - 系统诊断信息
   - 问题检测和建议

8. **配置重载** (`/ops/config/reload`)
   - 动态重载配置
   - POST 请求触发

9. **性能分析** (`/ops/debug/pprof/`)
   - CPU 性能分析
   - 内存分析
   - Goroutine 分析
   - 阻塞分析

10. **系统信息** (`/ops/info`)
    - 平台信息
    - Kubernetes 信息（如果在 K8s 中）
    - 主机信息

11. **版本信息** (`/ops/version`)
    - 应用版本
    - 构建时间
    - Git 提交信息

12. **服务发现** (`/ops/services`)
    - 服务列表
    - 服务信息

### ✅ 优雅关闭

- **信号处理** (SIGINT, SIGTERM)
- **超时控制**
- **并发关闭**
- **资源清理**

### ✅ 熔断器

- **三种状态** (关闭、打开、半开)
- **自动恢复**
- **失败计数**
- **超时控制**

### ✅ 重试机制

- **指数退避**
- **最大重试次数**
- **随机抖动**
- **可重试错误判断**

### ✅ 超时控制

- **操作超时**
- **HTTP 超时中间件**
- **可配置超时时间**

### ✅ 服务发现

- **服务注册**
- **服务注销**
- **服务列表**
- **健康检查集成**

## 🚀 快速开始

### 1. 基本使用

```go
import (
    "github.com/yourusername/golang/pkg/observability"
    "github.com/yourusername/golang/pkg/observability/operational"
)

// 创建可观测性集成
obs, _ := observability.NewObservability(observability.Config{
    ServiceName:            "my-service",
    EnableSystemMonitoring: true,
})

// 创建运维端点
endpoints := operational.NewOperationalEndpoints(operational.Config{
    Observability: obs,
    Port:          9090,
    PathPrefix:    "/ops",
    Enabled:       true,
})

// 启动
endpoints.Start()
defer endpoints.Stop(ctx)
```

### 2. 优雅关闭

```go
// 创建关闭管理器
shutdownManager := operational.NewShutdownManager(30 * time.Second)

// 注册关闭函数
shutdownManager.Register(operational.GracefulShutdown("observability", obs.Stop))
shutdownManager.Register(operational.GracefulShutdown("endpoints", endpoints.Stop))

// 等待关闭信号
shutdownManager.WaitForShutdown()
```

### 3. 熔断器

```go
// 创建熔断器
circuitBreaker := operational.NewCircuitBreaker(operational.CircuitBreakerConfig{
    Name:         "external-api",
    MaxFailures:  5,
    ResetTimeout: 60 * time.Second,
})

// 使用熔断器
err := circuitBreaker.Execute(ctx, func() error {
    return callExternalAPI()
})
```

### 4. 重试机制

```go
// 使用重试机制
err := operational.Retry(ctx, operational.DefaultRetryConfig(), func() error {
    return doSomething()
})
```

## 📊 端点使用示例

### 健康检查

```bash
# 健康检查
curl http://localhost:9090/ops/health

# 响应
{
  "status": true,
  "timestamp": "2025-01-XX...",
  "message": "healthy",
  "details": {
    "memory_usage": 45.2,
    "cpu_usage": 12.5,
    "goroutines": 150,
    "gc": 42
  }
}
```

### 就绪检查

```bash
# 就绪检查
curl http://localhost:9090/ops/ready

# 响应
{
  "ready": true,
  "checks": {
    "system_monitor": true,
    "observability": true
  }
}
```

### 指标导出

```bash
# JSON 格式
curl http://localhost:9090/ops/metrics

# Prometheus 格式
curl http://localhost:9090/ops/metrics/prometheus
```

### 性能分析

```bash
# CPU 性能分析（30秒）
curl http://localhost:9090/ops/debug/pprof/profile?seconds=30 > cpu.prof

# 堆内存分析
curl http://localhost:9090/ops/debug/pprof/heap > heap.prof

# Goroutine 分析
curl http://localhost:9090/ops/debug/pprof/goroutine > goroutine.prof
```

## 🎯 Kubernetes 集成

### Deployment 配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-service
spec:
  template:
    spec:
      containers:
      - name: app
        image: my-service:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090  # 运维端点端口
        livenessProbe:
          httpGet:
            path: /ops/live
            port: 9090
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ops/ready
            port: 9090
          initialDelaySeconds: 5
          periodSeconds: 5
        startupProbe:
          httpGet:
            path: /ops/health
            port: 9090
          initialDelaySeconds: 0
          periodSeconds: 5
          failureThreshold: 30
```

## 🔧 最佳实践

### 1. 端口分离

- **应用端口**: 8080（业务请求）
- **运维端口**: 9090（运维端点）

### 2. 安全配置

```go
// 生产环境建议：
// 1. 运维端点仅监听 localhost
// 2. 使用认证中间件
// 3. 限制访问 IP
endpoints := operational.NewOperationalEndpoints(operational.Config{
    Observability: obs,
    Port:          9090,
    PathPrefix:    "/ops",
    Enabled:       true,
})
```

### 3. 优雅关闭

```go
// 在 main 函数中
shutdownManager := operational.NewShutdownManager(30 * time.Second)
shutdownManager.Register(func(ctx context.Context) error {
    return obs.Stop(ctx)
})
shutdownManager.Register(func(ctx context.Context) error {
    return endpoints.Stop(ctx)
})

// 等待关闭信号
if err := shutdownManager.WaitForShutdown(); err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

### 4. 熔断器使用

```go
// 为外部 API 调用添加熔断保护
circuitBreaker := operational.NewCircuitBreaker(operational.CircuitBreakerConfig{
    Name:         "payment-api",
    MaxFailures:  5,
    ResetTimeout: 60 * time.Second,
})

err := circuitBreaker.Execute(ctx, func() error {
    return paymentAPI.Process(ctx, payment)
})
```

## 📚 更多文档

- [运维控制 README](../pkg/observability/operational/README.md)
- [完整使用指南](./OBSERVABILITY-COMPLETE-GUIDE.md)
- [配置集成指南](./CONFIG-INTEGRATION.md)

---

**版本**: v1.0.0
**最后更新**: 2025-01-XX
