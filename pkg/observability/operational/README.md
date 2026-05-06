# 运维控制功能

> **版本**: v1.0.0
> **Go版本**: 1.26+

完整的运维控制功能，提供健康检查、就绪检查、指标导出、优雅关闭、性能分析等运维必需的功能。

## 📋 功能特性

### 运维端点 ✅
- ✅ **健康检查** (`/ops/health`, `/ops/healthz`)
- ✅ **就绪检查** (`/ops/ready`, `/ops/readiness`)
- ✅ **存活检查** (`/ops/live`, `/ops/liveness`)
- ✅ **指标导出** (`/ops/metrics`)
- ✅ **Prometheus 指标** (`/ops/metrics/prometheus`)
- ✅ **仪表板数据** (`/ops/dashboard`)
- ✅ **诊断报告** (`/ops/diagnostics`)
- ✅ **配置重载** (`/ops/config/reload`)
- ✅ **性能分析** (`/ops/debug/pprof/`)
- ✅ **系统信息** (`/ops/info`)
- ✅ **版本信息** (`/ops/version`)
- ✅ **服务发现** (`/ops/services`)

### 优雅关闭 ✅
- ✅ **优雅关闭管理器**
- ✅ **信号处理** (SIGINT, SIGTERM)
- ✅ **超时控制**
- ✅ **并发关闭**

### 熔断器 ✅
- ✅ **三种状态** (关闭、打开、半开)
- ✅ **自动恢复**
- ✅ **失败计数**
- ✅ **超时控制**

### 重试机制 ✅
- ✅ **指数退避**
- ✅ **最大重试次数**
- ✅ **随机抖动**
- ✅ **可重试错误判断**

### 超时控制 ✅
- ✅ **操作超时**
- ✅ **HTTP 超时中间件**
- ✅ **可配置超时时间**

### 服务发现 ✅
- ✅ **服务注册**
- ✅ **服务注销**
- ✅ **服务列表**
- ✅ **健康检查集成**

## 🚀 快速开始

### 1. 创建运维端点

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
import "github.com/yourusername/golang/pkg/observability/operational"

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

// 使用熔断器执行操作
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

## 📊 端点说明

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
curl http://localhost:9090/ops/debug/pprof/profile?seconds=30

# 堆内存分析
curl http://localhost:9090/ops/debug/pprof/heap

# Goroutine 分析
curl http://localhost:9090/ops/debug/pprof/goroutine
```

## 🎯 最佳实践

### 1. Kubernetes 集成

```yaml
# Kubernetes Deployment
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
```

### 2. 优雅关闭

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

### 3. 熔断器使用

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

- [完整使用指南](../docs/OBSERVABILITY-COMPLETE-GUIDE.md)
- [配置集成指南](../docs/CONFIG-INTEGRATION.md)

---

**版本**: v1.0.0
**最后更新**: 2025-01-XX
