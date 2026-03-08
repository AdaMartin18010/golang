# 完整集成指南

> **版本**: v1.0.0
> **Go版本**: 1.25+

本指南说明如何将可观测性和运维控制功能完整集成到应用中。

## 📋 集成步骤

### 1. 配置集成

在 `configs/config.yaml` 中添加可观测性配置：

```yaml
observability:
  otlp:
    endpoint: "localhost:4317"
    insecure: true
    service_name: "my-service"
    service_version: "v1.0.0"
  system:
    enabled: true
    collect_interval: 5s
    enable_disk_monitor: true
    enable_load_monitor: true
    enable_apm_monitor: true
    rate_limit:
      enabled: true
      limit: 100
      window: 1s
    health_thresholds:
      max_memory_usage: 90.0
      max_cpu_usage: 95.0
      max_goroutines: 10000
    alerts:
      - id: "cpu-high"
        name: "CPU Usage High"
        metric_name: "system.cpu.usage"
        condition: "gt"
        threshold: 80.0
        level: "warning"
        enabled: true
        duration: 5m
        cooldown: 10m
```

### 2. 代码集成

在 `main.go` 中集成：

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/yourusername/golang/internal/config"
    "github.com/yourusername/golang/pkg/observability"
    "github.com/yourusername/golang/pkg/observability/operational"
)

func main() {
    // 1. 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. 创建可观测性集成
    obsConfig := observability.ConfigFromAppConfig(cfg)
    obs, err := observability.NewObservability(obsConfig)
    if err != nil {
        log.Fatalf("Failed to create observability: %v", err)
    }

    // 3. 应用告警规则
    observability.ApplyAlertRules(obs, cfg.Observability.System.Alerts)

    // 4. 启动可观测性
    if err := obs.Start(); err != nil {
        log.Fatalf("Failed to start observability: %v", err)
    }

    // 5. 创建运维控制端点
    operationalEndpoints := operational.NewOperationalEndpoints(operational.Config{
        Observability: obs,
        Port:          9090,
        PathPrefix:    "/ops",
        Enabled:       true,
    })

    // 6. 启动运维端点
    if err := operationalEndpoints.Start(); err != nil {
        log.Fatalf("Failed to start operational endpoints: %v", err)
    }

    // 7. 创建主 HTTP 服务器
    mux := http.NewServeMux()

    // 业务路由（使用追踪和指标）
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ctx, span := obs.Tracer("server").Start(r.Context(), "handler")
        defer span.End()

        meter := obs.Meter("server")
        counter, _ := meter.Int64Counter("requests_total")
        counter.Add(ctx, 1)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello, World!"))
    })

    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
        Handler: mux,
    }

    // 8. 启动主服务器
    go func() {
        log.Printf("Server starting on :%d", cfg.Server.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // 9. 创建优雅关闭管理器
    shutdownManager := operational.NewShutdownManager(30 * time.Second)

    // 注册关闭函数
    shutdownManager.Register(operational.GracefulShutdown("http-server", func(ctx context.Context) error {
        return server.Shutdown(ctx)
    }))
    shutdownManager.Register(operational.GracefulShutdown("observability", func(ctx context.Context) error {
        return obs.Stop(ctx)
    }))
    shutdownManager.Register(operational.GracefulShutdown("operational-endpoints", func(ctx context.Context) error {
        return operationalEndpoints.Stop(ctx)
    }))

    // 10. 等待关闭信号
    log.Println("Application running. Press Ctrl+C to shutdown gracefully...")
    if err := shutdownManager.WaitForShutdown(); err != nil {
        log.Printf("Shutdown error: %v", err)
    }

    log.Println("Application shutdown complete")
}
```

## 🎯 功能使用

### 追踪

```go
// 创建追踪
ctx, span := obs.Tracer("service-name").Start(ctx, "operation-name")
defer span.End()

// 添加属性
span.SetAttributes(
    attribute.String("key", "value"),
)
```

### 指标

```go
// 创建指标
meter := obs.Meter("service-name")
counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1)

gauge, _ := meter.Float64Gauge("memory_usage")
gauge.Record(ctx, 45.2)
```

### 健康检查

```bash
# 健康检查
curl http://localhost:9090/ops/health

# 就绪检查
curl http://localhost:9090/ops/ready

# 存活检查
curl http://localhost:9090/ops/live
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
# CPU 性能分析
go tool pprof http://localhost:9090/ops/debug/pprof/profile?seconds=30

# 堆内存分析
go tool pprof http://localhost:9090/ops/debug/pprof/heap

# Goroutine 分析
go tool pprof http://localhost:9090/ops/debug/pprof/goroutine
```

## 🔧 高级功能

### 熔断器

```go
import "github.com/yourusername/golang/pkg/observability/operational"

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

### 重试机制

```go
// 使用重试机制
err := operational.Retry(ctx, operational.DefaultRetryConfig(), func() error {
    return doSomething()
})
```

### 超时控制

```go
// 为操作添加超时
err := operational.WithTimeout(ctx, 5*time.Second, func(ctx context.Context) error {
    return longRunningOperation(ctx)
})
```

## 📊 Kubernetes 集成

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
        - containerPort: 8080  # 应用端口
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

## 📚 更多文档

- [运维控制完整指南](./OPERATIONAL-CONTROL.md)
- [配置集成指南](./CONFIG-INTEGRATION.md)
- [完整使用指南](./OBSERVABILITY-COMPLETE-GUIDE.md)

---

**版本**: v1.0.0
**最后更新**: 2025-01-XX
