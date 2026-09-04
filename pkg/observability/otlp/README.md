# 增强的 OTLP 集成

框架级别的 OpenTelemetry 集成，提供采样、追踪、指标的完整支持。

## 📋 功能特性

- ✅ **完整集成**: 追踪、指标的完整支持
- ✅ **采样支持**: 可配置的采样策略
- ✅ **动态调整**: 运行时动态调整采样率
- ✅ **资源标识**: 自动标识服务信息
- ✅ **批处理优化**: 可配置的批处理大小和超时
- ✅ **指标导出**: 可配置的指标导出间隔
- ⚠️ **日志导出**: 等待 OpenTelemetry 官方发布

## 🚀 快速开始

### 基本使用

```go
import "github.com/yourusername/golang/pkg/observability/otlp"

// 创建增强的 OTLP 集成
otlp, err := otlp.NewEnhancedOTLP(otlp.Config{
    ServiceName:       "my-service",
    ServiceVersion:    "v1.0.0",
    Endpoint:          "localhost:4317",
    Insecure:          true,
    SampleRate:        0.5,              // 50% 采样率
    MetricInterval:    10 * time.Second, // 指标导出间隔
    TraceBatchTimeout: 5 * time.Second,  // 追踪批处理超时
    TraceBatchSize:    512,               // 追踪批处理大小
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

// Config 配置选项
type Config struct {
    ServiceName       string          // 服务名称
    ServiceVersion    string          // 服务版本
    Endpoint          string          // OTLP 端点地址
    Insecure          bool            // 是否使用不安全连接
    Sampler           sampling.Sampler // 采样器
    SampleRate        float64         // 采样率（0.0-1.0）
    MetricInterval    time.Duration   // 指标导出间隔（默认：10秒）
    TraceBatchTimeout time.Duration   // 追踪批处理超时（默认：5秒）
    TraceBatchSize     int             // 追踪批处理大小（默认：512）
}
```

## 🔗 相关文档

- [采样机制](../../sampling/README.md)
- 追踪和定位
- OpenTelemetry 基础设施
