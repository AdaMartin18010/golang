# 可观测性完整指南

## 📋 概述

本文档提供了完整的可观测性功能使用指南，包括所有功能的详细说明、配置方法和最佳实践。

## 🚀 快速开始

### 1. 基本集成

```go
import (
    "github.com/yourusername/golang/pkg/observability"
    "github.com/yourusername/golang/pkg/observability/system"
)

// 创建可观测性集成
obs, err := observability.NewObservability(observability.Config{
    ServiceName:            "my-service",
    ServiceVersion:         "v1.0.0",
    OTLPEndpoint:           "localhost:4317",
    OTLPInsecure:           true,
    SampleRate:             0.5,
    EnableSystemMonitoring: true,
    SystemCollectInterval:  5 * time.Second,
})
if err != nil {
    log.Fatal(err)
}

// 启动
obs.Start()
defer obs.Stop(ctx)
```

### 2. 完整功能集成

```go
obs, err := observability.NewObservability(observability.Config{
    ServiceName:            "my-service",
    ServiceVersion:         "v1.0.0",
    OTLPEndpoint:           "localhost:4317",
    OTLPInsecure:           true,
    SampleRate:             0.5,
    MetricInterval:         10 * time.Second,
    TraceBatchTimeout:      5 * time.Second,
    TraceBatchSize:         512,
    EnableSystemMonitoring: true,
    SystemCollectInterval:  5 * time.Second,
    EnableDiskMonitor:      true,
    EnableLoadMonitor:       true,
    EnableAPMMonitor:        true,
    RateLimitConfig: &system.RateLimiterConfig{
        Enabled: true,
        Limit:   100,
        Window:  1 * time.Second,
    },
    HealthThresholds: system.DefaultHealthThresholds(),
})
```

## 📊 功能模块详解

### 1. 追踪 (Tracing)

```go
// 获取追踪器
tracer := obs.Tracer("my-service")

// 创建 Span
ctx, span := tracer.Start(ctx, "operation")
defer span.End()

// 添加属性
span.SetAttributes(
    attribute.String("user.id", "123"),
    attribute.Int("request.size", 1024),
)
```

### 2. 指标 (Metrics)

```go
// 获取指标器
meter := obs.Meter("my-service")

// 创建计数器
counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1)

// 创建直方图
histogram, _ := meter.Float64Histogram("request_duration")
histogram.Record(ctx, 0.125)

// 创建 Gauge
gauge, _ := meter.Int64ObservableGauge("active_connections")
```

### 3. 系统监控

```go
systemMonitor := obs.GetSystemMonitor()

// 获取平台信息
platformInfo := obs.GetPlatformInfo()
fmt.Printf("OS: %s, Container: %v\n", platformInfo.OS, obs.IsContainer())

// 健康检查
healthChecker := systemMonitor.GetHealthChecker()
status := healthChecker.Check(ctx)
fmt.Printf("Health: %v\n", status.Healthy)
```

### 4. APM 监控

```go
apmMonitor := obs.GetAPMMonitor()
if apmMonitor != nil {
    start := time.Now()
    // 执行业务逻辑
    duration := time.Since(start)
    apmMonitor.RecordRequest(ctx, duration, 200)
}
```

### 5. 限流器

```go
rateLimiter := obs.GetRateLimiter()
if rateLimiter != nil {
    if rateLimiter.Allow(ctx) {
        // 处理请求
    } else {
        // 请求被限流
    }
}
```

### 6. 告警系统

```go
alertManager := obs.GetAlertManager()
if alertManager != nil {
    // 添加自定义告警规则
    rule := system.AlertRule{
        ID:         "custom-alert",
        Name:       "Custom Alert",
        MetricName: "system.cpu.usage",
        Condition:  "gt",
        Threshold:  80.0,
        Level:      system.AlertLevelWarning,
        Enabled:    true,
        Duration:   5 * time.Minute,
        Cooldown:   10 * time.Minute,
    }
    alertManager.AddRule(rule)

    // 添加告警处理器
    alertManager.AddHandler(myAlertHandler)

    // 检查指标
    alertManager.Check(ctx, "system.cpu.usage", 85.0, nil)
}
```

### 7. 诊断工具

```go
diagnostics := obs.GetDiagnostics()
if diagnostics != nil {
    // 生成诊断报告
    report, err := diagnostics.GenerateReport(ctx)
    if err == nil {
        fmt.Printf("Issues: %d\n", len(report.Issues))
        for _, issue := range report.Issues {
            fmt.Printf("  - %s: %s\n", issue.Level, issue.Description)
        }
    }

    // 导出 JSON
    jsonData, _ := diagnostics.ExportJSON(ctx)
    fmt.Println(string(jsonData))
}
```

### 8. 资源预测

```go
predictor := obs.GetPredictor()
if predictor != nil {
    // 预测内存使用
    prediction, err := predictor.Predict(ctx, "system.memory.usage", 1*time.Hour)
    if err == nil {
        fmt.Printf("Predicted: %.2f (confidence: %.2f, trend: %s)\n",
            prediction.PredictedValue,
            prediction.Confidence,
            prediction.Trend,
        )
    }
}
```

### 9. 指标导出

```go
metricsExporter := obs.GetMetricsExporter()
if metricsExporter != nil {
    // 导出快照
    snapshot, _ := metricsExporter.Export(ctx)
    fmt.Printf("Metrics: %d\n", len(snapshot.Metrics))

    // 导出 JSON
    jsonData, _ := metricsExporter.ExportJSON(ctx)
    fmt.Println(string(jsonData))

    // 查询历史
    history := metricsExporter.GetHistory(10)
    fmt.Printf("History: %d snapshots\n", len(history))
}
```

### 10. 仪表板导出

```go
dashboardExporter := obs.GetDashboardExporter()
if dashboardExporter != nil {
    // JSON 格式
    jsonData, _ := dashboardExporter.ExportJSON(ctx)
    fmt.Println(string(jsonData))

    // Prometheus 格式
    promData, _ := dashboardExporter.ExportForPrometheus(ctx)
    fmt.Println(promData)
}
```

## 🎯 最佳实践

### 1. 采样率配置

```go
// 生产环境：低采样率（1-10%）
SampleRate: 0.01

// 开发环境：高采样率（50-100%）
SampleRate: 0.5
```

### 2. 监控间隔

```go
// 系统监控：5-10 秒
SystemCollectInterval: 5 * time.Second

// 指标导出：10-30 秒
MetricInterval: 10 * time.Second
```

### 3. 告警规则

```go
// CPU 告警：超过 80% 持续 5 分钟
rule := system.AlertRule{
    MetricName: "system.cpu.usage",
    Condition:  "gt",
    Threshold:  80.0,
    Duration:   5 * time.Minute,
    Cooldown:   10 * time.Minute,
}

// 内存告警：超过 90% 立即告警
rule := system.AlertRule{
    MetricName: "system.memory.usage",
    Condition:  "gt",
    Threshold:  90.0,
    Duration:   0, // 立即告警
    Cooldown:   5 * time.Minute,
}
```

### 4. 健康检查

```go
// 设置合理的阈值
thresholds := system.HealthThresholds{
    MaxMemoryUsage: 90.0,
    MaxCPUUsage:    95.0,
    MaxGoroutines:  10000,
}

// 定期检查
healthChecker.CheckPeriodically(ctx, func(status system.HealthStatus) {
    if !status.Healthy {
        // 发送告警
        sendAlert(status)
    }
})
```

## 📚 配置参考

### 完整配置示例

```go
config := observability.Config{
    // OTLP 配置
    ServiceName:       "my-service",
    ServiceVersion:    "v1.0.0",
    OTLPEndpoint:      "localhost:4317",
    OTLPInsecure:     true,
    SampleRate:        0.1, // 10% 采样率
    MetricInterval:    10 * time.Second,
    TraceBatchTimeout: 5 * time.Second,
    TraceBatchSize:    512,

    // 系统监控配置
    EnableSystemMonitoring: true,
    SystemCollectInterval:  5 * time.Second,
    EnableDiskMonitor:     true,
    EnableLoadMonitor:      true,
    EnableAPMMonitor:       true,

    // 限流器配置
    RateLimitConfig: &system.RateLimiterConfig{
        Enabled: true,
        Limit:   100, // 每秒 100 个请求
        Window:  1 * time.Second,
    },

    // 健康检查配置
    HealthThresholds: system.HealthThresholds{
        MaxMemoryUsage: 90.0,
        MaxCPUUsage:    95.0,
        MaxGoroutines:  10000,
    },
}
```

## 🔧 故障排查

### 1. 指标未导出

- 检查 OTLP 端点是否可访问
- 检查采样率是否设置过低
- 检查指标导出间隔是否合理

### 2. 告警未触发

- 检查告警规则是否启用
- 检查冷却时间是否过长
- 检查阈值是否设置合理

### 3. 健康检查失败

- 检查健康阈值是否设置合理
- 检查系统资源是否真的不足
- 查看健康状态详细信息

## 📖 相关文档

- [系统监控 README](../pkg/observability/system/README.md)
- [完整实现报告](./COMPLETE-IMPLEMENTATION-FINAL-REPORT.md)
- [高级功能报告](./ULTIMATE-ADVANCED-FEATURES.md)

---

**版本**: v1.0.0
**最后更新**: 2025-01-XX
