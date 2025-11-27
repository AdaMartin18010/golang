# 增强的 OTLP 集成

框架级别的 OpenTelemetry 集成，提供采样、追踪、指标的完整支持。

## 📋 功能特性

- ✅ **完整集成**: 追踪、指标、日志的完整支持
- ✅ **采样支持**: 可配置的采样策略
- ✅ **动态调整**: 运行时动态调整采样率
- ✅ **资源标识**: 自动标识服务信息

## 🚀 快速开始

### 基本使用

```go
import "github.com/yourusername/golang/pkg/observability/otlp"

// 创建增强的 OTLP 集成
otlp, err := otlp.NewEnhancedOTLP(otlp.Config{
    ServiceName:    "my-service",
    ServiceVersion: "v1.0.0",
    Endpoint:       "localhost:4317",
    Insecure:       true,
    SampleRate:     0.5, // 50% 采样率
})
if err != nil {
    log.Fatal(err)
}
defer otlp.Shutdown(context.Background())

// 使用追踪器
tracer := otlp.Tracer("my-tracer")
ctx, span := tracer.Start(ctx, "operation")
defer span.End()

// 使用指标器
meter := otlp.Meter("my-meter")
counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1)
```

### 动态调整采样率

```go
// 更新采样率
otlp.UpdateSampleRate(0.1) // 降低到 10%
```

## 📚 API 参考

### EnhancedOTLP

```go
type EnhancedOTLP struct {
    // ...
}

func NewEnhancedOTLP(cfg Config) (*EnhancedOTLP, error)
func (e *EnhancedOTLP) Shutdown(ctx context.Context) error
func (e *EnhancedOTLP) Tracer(name string) trace.Tracer
func (e *EnhancedOTLP) Meter(name string) metric.Meter
func (e *EnhancedOTLP) ShouldSample(ctx context.Context) bool
func (e *EnhancedOTLP) UpdateSampleRate(rate float64) error
```

## 🔗 相关文档

- [采样机制](../../sampling/README.md)
- [追踪和定位](../../tracing/README.md)
- [OpenTelemetry 基础设施](../../../internal/infrastructure/observability/otlp/README.md)
