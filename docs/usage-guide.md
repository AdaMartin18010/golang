# 可观测性功能使用指南

## 概述

本文档提供完整的使用指南，包括 OTLP、日志轮转和 eBPF 功能的使用方法。

## 1. OTLP 集成

### 1.1 基本使用

```go
import (
    "github.com/yourusername/golang/pkg/observability/otlp"
    "github.com/yourusername/golang/pkg/sampling"
)

// 创建采样器
sampler, err := sampling.NewProbabilisticSampler(0.5) // 50% 采样率
if err != nil {
    log.Fatal(err)
}

// 初始化 OTLP
otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
    ServiceName:       "my-service",
    ServiceVersion:    "v1.0.0",
    Endpoint:          "localhost:4317",
    Insecure:          true,
    Sampler:           sampler,
    SampleRate:        0.5,
    MetricInterval:    10 * time.Second,  // 指标导出间隔
    TraceBatchTimeout: 5 * time.Second,   // 追踪批处理超时
    TraceBatchSize:    512,                // 追踪批处理大小
})
if err != nil {
    log.Fatal(err)
}
defer otlpClient.Shutdown(context.Background())
```

### 1.2 使用追踪

```go
tracer := otlpClient.Tracer("my-service")
ctx, span := tracer.Start(ctx, "operation")
defer span.End()

// 添加属性
span.SetAttributes(
    attribute.String("user.id", "123"),
    attribute.Int("request.size", 1024),
)
```

### 1.3 使用指标

```go
meter := otlpClient.Meter("my-service")

// 创建计数器
counter, _ := meter.Int64Counter(
    "requests_total",
    metric.WithDescription("Total number of requests"),
)

// 记录指标
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("method", "GET"),
    attribute.String("path", "/api/users"),
))

// 创建直方图
histogram, _ := meter.Float64Histogram(
    "request_duration_seconds",
    metric.WithDescription("Request duration in seconds"),
)
histogram.Record(ctx, 0.123)
```

### 1.4 动态调整采样率

```go
// 更新采样率
err := otlpClient.UpdateSampleRate(0.1) // 降低到 10%
if err != nil {
    log.Printf("Failed to update sample rate: %v", err)
}
```

## 2. 日志轮转

### 2.1 基本使用

```go
import "github.com/yourusername/golang/pkg/logger"
import "log/slog"

// 使用默认配置
cfg := logger.DefaultRotationConfig("logs/app.log")
logger, err := logger.NewRotatingLogger(slog.LevelInfo, cfg)
if err != nil {
    log.Fatal(err)
}

// 使用日志
logger.Info("Application started")
```

### 2.2 使用预定义配置

```go
// 生产环境配置
prodCfg := logger.ProductionRotationConfig("logs/app.log")
logger, _ := logger.NewRotatingLogger(slog.LevelInfo, prodCfg)

// 开发环境配置
devCfg := logger.DevelopmentRotationConfig("logs/app.log")
logger, _ := logger.NewRotatingLogger(slog.LevelDebug, devCfg)
```

### 2.3 自定义配置

```go
customCfg := logger.RotationConfig{
    Filename:   "logs/app.log",
    MaxSize:    200,  // 200MB
    MaxBackups: 20,   // 保留20个备份
    MaxAge:     60,   // 保留60天
    Compress:   true, // 压缩旧日志
    LocalTime:  true, // 使用本地时间
}

// 验证配置
if err := logger.ValidateRotationConfig(customCfg); err != nil {
    log.Fatal(err)
}

logger, err := logger.NewRotatingLogger(slog.LevelInfo, customCfg)
if err != nil {
    log.Fatal(err)
}
```

### 2.4 配置文件方式

在 `configs/config.yaml` 中配置：

```yaml
logging:
  level: "info"
  format: "json"
  output: "file"
  output_path: "logs/app.log"
  rotation:
    max_size: 100      # 单个日志文件最大大小（MB）
    max_backups: 10    # 保留的旧日志文件数量
    max_age: 30        # 保留旧日志文件的天数
    compress: true     # 是否压缩轮转后的旧日志文件
```

## 3. eBPF 收集器

### 3.1 基本使用（框架）

```go
import "github.com/yourusername/golang/pkg/observability/ebpf"

// 初始化 eBPF 收集器
collector, err := ebpf.NewCollector(ebpf.Config{
    Tracer:                    tracer,
    Meter:                     meter,
    Enabled:                   true,
    CollectInterval:           5 * time.Second,
    EnableSyscallTracking:     true,
    EnableNetworkMonitoring:   true,
    EnablePerformanceProfiling: false,
})
if err != nil {
    log.Fatal(err)
}

// 启动收集器
if err := collector.Start(); err != nil {
    log.Fatal(err)
}
defer collector.Stop()
```

### 3.2 收集指标

```go
// 收集系统调用指标
if err := collector.CollectSyscallMetrics(ctx); err != nil {
    log.Printf("Failed to collect syscall metrics: %v", err)
}

// 收集网络指标
if err := collector.CollectNetworkMetrics(ctx); err != nil {
    log.Printf("Failed to collect network metrics: %v", err)
}
```

### 3.3 记录追踪

```go
// 记录系统调用追踪
collector.RecordSyscallTrace(ctx, "open", 12345, 10*time.Millisecond)
```

### 3.4 注意事项

⚠️ **重要**：当前的 eBPF 实现仅为框架接口，实际功能需要：

1. 添加 `github.com/cilium/ebpf` 依赖
2. 编写 eBPF C 程序（`.bpf.c` 文件）
3. 编译 eBPF 程序
4. 实现程序加载和数据收集逻辑

参考文档：
- [eBPF 深度解析](../architecture/tech-stack/observability/ebpf.md)
- [cilium/ebpf 官方文档](https://github.com/cilium/ebpf)

## 4. 完整集成示例

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yourusername/golang/pkg/logger"
    "github.com/yourusername/golang/pkg/observability/ebpf"
    "github.com/yourusername/golang/pkg/observability/otlp"
    "github.com/yourusername/golang/pkg/sampling"
    "log/slog"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 1. 初始化 OTLP
    sampler, _ := sampling.NewProbabilisticSampler(0.5)
    otlpClient, _ := otlp.NewEnhancedOTLP(otlp.Config{
        ServiceName:    "my-service",
        ServiceVersion: "v1.0.0",
        Endpoint:       "localhost:4317",
        Insecure:       true,
        Sampler:        sampler,
    })
    defer otlpClient.Shutdown(ctx)

    // 2. 初始化日志轮转
    cfg := logger.ProductionRotationConfig("logs/app.log")
    appLogger, _ := logger.NewRotatingLogger(slog.LevelInfo, cfg)
    slog.SetDefault(appLogger.Logger)

    // 3. 初始化 eBPF（可选）
    tracer := otlpClient.Tracer("my-service")
    meter := otlpClient.Meter("my-service")
    ebpfCollector, _ := ebpf.NewCollector(ebpf.Config{
        Tracer:                  tracer,
        Meter:                   meter,
        Enabled:                 false, // 需要实际的 eBPF 程序
        EnableSyscallTracking:   true,
        EnableNetworkMonitoring: true,
    })
    defer ebpfCollector.Stop()

    // 4. 业务逻辑
    appLogger.Info("Application started")

    // 等待中断
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
}
```

## 5. 最佳实践

### 5.1 OTLP

- 生产环境使用 TLS 连接
- 合理设置采样率（建议 0.1-0.5）
- 监控指标导出间隔，避免过于频繁
- 使用批处理减少网络开销

### 5.2 日志轮转

- 生产环境使用 `ProductionRotationConfig`
- 开发环境使用 `DevelopmentRotationConfig`
- 根据磁盘空间调整 `MaxSize` 和 `MaxBackups`
- 启用压缩以节省空间

### 5.3 eBPF

- 仅在需要系统级监控时启用
- 需要 root 权限（或 CAP_BPF 能力）
- 仅在 Linux 系统上可用
- 考虑使用现成的工具（如 Pixie、Datadog Agent）

## 6. 故障排查

### OTLP 连接失败

```go
// 检查端点是否可达
// 检查防火墙设置
// 检查 TLS 配置
```

### 日志轮转不工作

```go
// 检查文件权限
// 检查磁盘空间
// 验证配置是否正确
if err := logger.ValidateRotationConfig(cfg); err != nil {
    log.Printf("Invalid config: %v", err)
}
```

### eBPF 无法启动

```go
// 检查是否在 Linux 系统上
// 检查是否有 root 权限
// 检查内核版本（需要 4.18+）
// 检查是否加载了 eBPF 程序
```
