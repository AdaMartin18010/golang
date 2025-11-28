# 生产环境最佳实践

## 概述

本文档提供在生产环境中使用可观测性功能的最佳实践和建议。

## 1. 系统监控配置

### 1.1 监控间隔

```go
// 生产环境推荐配置
systemMonitor, err := system.NewSystemMonitor(system.SystemConfig{
    Meter:            otlpClient.Meter("system"),
    Enabled:          true,
    CollectInterval:  10 * time.Second, // 生产环境使用 10 秒间隔
    EnableDiskMonitor: true,
    HealthThresholds: system.HealthThresholds{
        MaxMemoryUsage: 85.0,  // 85% 内存使用率阈值
        MaxCPUUsage:    90.0,   // 90% CPU 使用率阈值
        MaxGoroutines:  5000,   // 5000 个 Goroutine 阈值
    },
})
```

### 1.2 健康检查配置

```go
// 配置健康检查阈值
healthThresholds := system.HealthThresholds{
    MaxMemoryUsage: 85.0,  // 根据实际内存配置调整
    MaxCPUUsage:    90.0,  // 根据 CPU 核心数调整
    MaxGoroutines:  5000,   // 根据应用特性调整
    MinGCInterval:  1 * time.Second,
}

// 定期执行健康检查
healthChecker := systemMonitor.GetHealthChecker()
healthChecker.CheckPeriodically(ctx, func(status system.HealthStatus) {
    if !status.Healthy {
        // 发送告警
        alerting.SendAlert("system_unhealthy", status.Message)
    }
})
```

## 2. OTLP 配置

### 2.1 生产环境配置

```go
otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
    ServiceName:       "production-service",
    ServiceVersion:    getVersion(), // 从构建信息获取
    Endpoint:          os.Getenv("OTLP_ENDPOINT"),
    Insecure:          false, // 生产环境使用 TLS
    SampleRate:        0.1,   // 生产环境使用 10% 采样率
    MetricInterval:    30 * time.Second, // 30 秒导出间隔
    TraceBatchTimeout: 10 * time.Second, // 10 秒批处理超时
    TraceBatchSize:    1024,  // 更大的批处理大小
})
```

### 2.2 TLS 配置

```go
// 生产环境必须使用 TLS
opts := []otlptracegrpc.Option{
    otlptracegrpc.WithEndpoint(cfg.Endpoint),
    otlptracegrpc.WithTLSCredentials(credentials.NewTLS(&tls.Config{
        ServerName: "otel-collector.example.com",
    })),
}
```

### 2.3 采样策略

```go
// 根据请求类型调整采样率
sampler := sampling.NewAdaptiveSampler(sampling.AdaptiveConfig{
    BaseRate:           0.1,  // 基础采样率 10%
    ErrorRate:          1.0,  // 错误请求 100% 采样
    SlowRequestRate:    0.5,  // 慢请求 50% 采样
    SlowRequestThreshold: 1 * time.Second,
})
```

## 3. 日志配置

### 3.1 生产环境日志配置

```go
// 使用生产环境配置
rotationCfg := logger.ProductionRotationConfig("logs/app.log")
// 或自定义配置
customCfg := logger.RotationConfig{
    Filename:   "/var/log/app/app.log",
    MaxSize:    500,  // 500MB
    MaxBackups: 30,   // 保留 30 个备份
    MaxAge:     90,   // 保留 90 天
    Compress:   true, // 压缩旧日志
    LocalTime:  true,
}

logger, err := logger.NewRotatingLogger(slog.LevelInfo, customCfg)
```

### 3.2 日志级别

```go
// 生产环境推荐使用 Info 级别
// 开发环境可以使用 Debug 级别
level := slog.LevelInfo
if os.Getenv("ENV") == "development" {
    level = slog.LevelDebug
}
```

### 3.3 敏感信息过滤

```go
// 创建带敏感信息过滤的 Logger
type SensitiveLogger struct {
    *logger.Logger
}

func (l *SensitiveLogger) Info(msg string, args ...any) {
    // 过滤敏感信息
    filteredArgs := filterSensitiveInfo(args)
    l.Logger.Info(msg, filteredArgs...)
}
```

## 4. 容器环境配置

### 4.1 Docker 环境

```go
// 检测容器环境
if systemMonitor.IsContainer() {
    info := systemMonitor.GetPlatformInfo()
    log.Printf("Running in container: %s", info.ContainerID)

    // 容器环境特殊配置
    // 1. 使用更短的日志保留时间
    // 2. 使用更小的批处理大小
    // 3. 使用更频繁的健康检查
}
```

### 4.2 Kubernetes 环境

```go
// 检测 Kubernetes 环境
if systemMonitor.IsKubernetes() {
    info := systemMonitor.GetPlatformInfo()
    log.Printf("Running in Kubernetes: Pod=%s, Node=%s",
        info.KubernetesPod, info.KubernetesNode)

    // Kubernetes 环境特殊配置
    // 1. 使用 Pod 名作为服务名
    // 2. 添加 Kubernetes 标签到指标
    // 3. 使用更短的超时时间
}
```

## 5. 错误处理和重试

### 5.1 监控错误处理

```go
// 使用重试机制
retryConfig := system.DefaultRetryConfig()
retryConfig.MaxRetries = 5
retryConfig.InitialDelay = 2 * time.Second

err := system.Retry(retryConfig, func() error {
    return systemMonitor.Start()
})
if err != nil {
    log.Fatalf("Failed to start system monitor after retries: %v", err)
}
```

### 5.2 OTLP 连接错误处理

```go
// 优雅降级：OTLP 连接失败时继续运行
otlpClient, err := otlp.NewEnhancedOTLP(cfg)
if err != nil {
    log.Printf("Warning: Failed to initialize OTLP: %v (continuing without OTLP)", err)
    // 继续运行，但不导出指标和追踪
} else {
    defer otlpClient.Shutdown(ctx)
}
```

## 6. 性能优化

### 6.1 指标导出间隔

```go
// 根据指标重要性调整导出间隔
// 关键指标：10 秒
// 一般指标：30 秒
// 低频指标：60 秒
```

### 6.2 批处理大小

```go
// 根据网络条件调整批处理大小
// 高带宽：1024
// 低带宽：256
// 不稳定网络：128
```

### 6.3 采样率

```go
// 根据流量调整采样率
// 高流量：0.01 (1%)
// 中流量：0.1 (10%)
// 低流量：0.5 (50%)
```

## 7. 安全考虑

### 7.1 敏感信息

```go
// 不要在日志中记录敏感信息
// ❌ 错误示例
logger.Info("User login", "password", password)

// ✅ 正确示例
logger.Info("User login", "user_id", userID)
```

### 7.2 TLS 配置

```go
// 生产环境必须使用 TLS
// 配置证书验证
tlsConfig := &tls.Config{
    ServerName:         "otel-collector.example.com",
    InsecureSkipVerify: false, // 生产环境必须验证
}
```

### 7.3 访问控制

```go
// 限制监控数据的访问
// 1. 使用网络隔离
// 2. 使用认证和授权
// 3. 加密传输
```

## 8. 监控告警

### 8.1 健康检查告警

```go
healthChecker.CheckPeriodically(ctx, func(status system.HealthStatus) {
    if !status.Healthy {
        // 发送告警
        sendAlert(Alert{
            Level:   "warning",
            Message: status.Message,
            Metrics: map[string]float64{
                "memory_usage": status.MemoryUsage,
                "cpu_usage":    status.CPUUsage,
                "goroutines":   float64(status.Goroutines),
            },
        })
    }
})
```

### 8.2 资源使用告警

```go
// 监控资源使用并告警
if memStats.Alloc > threshold {
    sendAlert(Alert{
        Level:   "critical",
        Message: "Memory usage exceeds threshold",
    })
}
```

## 9. 部署建议

### 9.1 资源限制

```yaml
# Kubernetes 资源限制
resources:
  limits:
    cpu: "2"
    memory: "4Gi"
  requests:
    cpu: "1"
    memory: "2Gi"
```

### 9.2 健康检查端点

```go
// 提供健康检查 HTTP 端点
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    status := systemMonitor.CheckHealth(r.Context())
    if status.Healthy {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    json.NewEncoder(w).Encode(status)
})
```

### 9.3 优雅关闭

```go
// 优雅关闭所有监控
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

<-sigChan

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// 停止监控
systemMonitor.Stop()

// 关闭 OTLP（确保数据导出完成）
otlpClient.Shutdown(ctx)
```

## 10. 故障排查

### 10.1 常见问题

1. **OTLP 连接失败**
   - 检查网络连接
   - 检查防火墙设置
   - 检查 TLS 配置

2. **指标不导出**
   - 检查 Meter 是否正确初始化
   - 检查指标是否注册
   - 检查导出间隔设置

3. **内存泄漏**
   - 检查 Goroutine 数量
   - 检查内存使用趋势
   - 检查 GC 频率

### 10.2 调试模式

```go
// 启用调试模式
if os.Getenv("DEBUG") == "true" {
    // 使用更详细的日志
    logger.SetLevel(slog.LevelDebug)

    // 使用更短的收集间隔
    collectInterval = 1 * time.Second

    // 使用更高的采样率
    sampleRate = 1.0
}
```

## 11. 性能基准

### 11.1 监控开销

- CPU 开销：< 2%
- 内存开销：< 50MB
- 网络开销：< 1MB/s（取决于采样率）

### 11.2 优化建议

1. 使用合理的采样率
2. 使用批处理减少网络请求
3. 使用异步导出避免阻塞
4. 定期清理旧数据

## 📚 相关文档

- [使用指南](./usage-guide.md)
- [系统监控实现](./system-monitoring-implementation.md)
- [OTLP 集成](../pkg/observability/otlp/README.md)
