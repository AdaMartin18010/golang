# 快速开始指南

## 🚀 5 分钟快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/yourusername/golang/pkg/observability"
)

func main() {
    ctx := context.Background()

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

    // 使用追踪
    tracer := obs.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "operation")
    defer span.End()

    // 使用指标
    meter := obs.Meter("my-service")
    counter, _ := meter.Int64Counter("requests_total")
    counter.Add(ctx, 1)

    log.Println("Observability started!")
}
```

### 2. 完整功能使用

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/yourusername/golang/pkg/observability"
    "github.com/yourusername/golang/pkg/observability/system"
)

func main() {
    ctx := context.Background()

    // 创建完整的可观测性集成
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
    if err != nil {
        log.Fatal(err)
    }

    // 启动
    obs.Start()
    defer obs.Stop(ctx)

    // 使用所有功能
    useAllFeatures(ctx, obs)

    log.Println("Complete observability integration running!")
}

func useAllFeatures(ctx context.Context, obs *observability.Observability) {
    // 1. 追踪
    tracer := obs.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "operation")
    defer span.End()

    // 2. 指标
    meter := obs.Meter("my-service")
    counter, _ := meter.Int64Counter("requests_total")
    counter.Add(ctx, 1)

    // 3. APM 监控
    apmMonitor := obs.GetAPMMonitor()
    if apmMonitor != nil {
        start := time.Now()
        // 执行业务逻辑
        time.Sleep(50 * time.Millisecond)
        apmMonitor.RecordRequest(ctx, time.Since(start), 200)
    }

    // 4. 限流器
    rateLimiter := obs.GetRateLimiter()
    if rateLimiter != nil && rateLimiter.Allow(ctx) {
        // 处理请求
    }

    // 5. 告警
    alertManager := obs.GetAlertManager()
    if alertManager != nil {
        alertManager.Check(ctx, "system.cpu.usage", 85.0, nil)
    }

    // 6. 诊断
    diagnostics := obs.GetDiagnostics()
    if diagnostics != nil {
        report, _ := diagnostics.GenerateReport(ctx)
        log.Printf("Diagnostic report: %d issues", len(report.Issues))
    }

    // 7. 预测
    predictor := obs.GetPredictor()
    if predictor != nil {
        prediction, _ := predictor.Predict(ctx, "system.memory.usage", 1*time.Hour)
        log.Printf("Predicted memory: %.2f", prediction.PredictedValue)
    }

    // 8. 仪表板
    dashboardExporter := obs.GetDashboardExporter()
    if dashboardExporter != nil {
        jsonData, _ := dashboardExporter.ExportJSON(ctx)
        log.Printf("Dashboard data: %d bytes", len(jsonData))
    }
}
```

### 3. 配置文件使用

```yaml
# configs/observability.yaml
observability:
  otlp:
    service_name: "my-service"
    endpoint: "localhost:4317"
    insecure: true
    sample_rate: 0.5
  system:
    enabled: true
    collect_interval: 5s
```

```go
// 从配置文件加载
import "github.com/spf13/viper"

viper.SetConfigFile("configs/observability.yaml")
viper.ReadInConfig()

obs, _ := observability.NewObservability(observability.Config{
    ServiceName: viper.GetString("observability.otlp.service_name"),
    // ...
})
```

## 📚 更多示例

查看 `examples/observability/` 目录了解更多示例：

- `complete-integration/main.go` - 完整集成示例
- `system-monitoring/main.go` - 系统监控示例
- `health-check/main.go` - 健康检查示例
- `advanced-features/main.go` - 高级功能示例

## 🔧 配置选项

### OTLP 配置

- `ServiceName` - 服务名称
- `ServiceVersion` - 服务版本
- `OTLPEndpoint` - OTLP 端点地址
- `OTLPInsecure` - 是否使用不安全连接
- `SampleRate` - 采样率（0.0-1.0）
- `MetricInterval` - 指标导出间隔
- `TraceBatchTimeout` - 追踪批处理超时
- `TraceBatchSize` - 追踪批处理大小

### 系统监控配置

- `EnableSystemMonitoring` - 是否启用系统监控
- `SystemCollectInterval` - 系统监控收集间隔
- `EnableDiskMonitor` - 是否启用磁盘监控
- `EnableLoadMonitor` - 是否启用负载监控
- `EnableAPMMonitor` - 是否启用 APM 监控
- `RateLimitConfig` - 限流器配置
- `HealthThresholds` - 健康检查阈值

## 🎯 下一步

1. 查看 [完整使用指南](./OBSERVABILITY-COMPLETE-GUIDE.md)
2. 查看 [功能总览](./OBSERVABILITY-FEATURES-SUMMARY.md)
3. 查看 [系统监控 README](../pkg/observability/system/README.md)

---

**版本**: v1.0.0
