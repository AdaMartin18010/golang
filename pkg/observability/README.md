# Observability 可观测性库

> **版本**: v2.0.0
> **Go版本**: 1.25+

---

## 📋 概述

完整的可观测性解决方案，提供分布式追踪、指标收集和结构化日志功能。

### 核心特性

- ✅ **分布式追踪 (Tracing)**
  - Span管理和嵌套
  - Context传播
  - 采样策略（Always、Probability）
  - 标签和日志
  - 错误追踪

- ✅ **指标收集 (Metrics)**
  - Counter（计数器）
  - Gauge（仪表）
  - Histogram（直方图）
  - Prometheus格式导出
  - 自动运行时指标

- ✅ **结构化日志 (Logging)**
  - 多级日志（Debug/Info/Warn/Error/Fatal）
  - 字段支持
  - Context集成
  - 可插拔钩子系统
  - 基于slog的高性能实现

---

## 🚀 快速开始

### 分布式追踪

```go
import "github.com/yourusername/golang/pkg/observability"

// 创建追踪器
recorder := observability.NewInMemoryRecorder()
sampler := &observability.AlwaysSampler{}
tracer := observability.NewTracer("my-service", recorder, sampler)
observability.SetGlobalTracer(tracer)

// 开始追踪
ctx := context.Background()
span, ctx := observability.StartSpan(ctx, "operation")
defer span.Finish()

// 添加标签和日志
span.SetTag("user_id", "123")
span.LogFields(map[string]interface{}{
    "event": "cache_hit",
})
```

### 指标收集

```go
// 创建指标
counter := observability.RegisterCounter(
    "requests_total",
    "Total requests",
    map[string]string{"service": "api"},
)

histogram := observability.RegisterHistogram(
    "request_duration_seconds",
    "Request latency",
    nil,
    nil,
)

gauge := observability.RegisterGauge(
    "active_connections",
    "Active connections",
    nil,
)

// 使用指标
counter.Inc()
histogram.Observe(0.125)
gauge.Set(42)

// 导出Prometheus格式
metrics := observability.ExportMetrics()
fmt.Println(metrics)
```

### 结构化日志

```go
// 创建日志记录器
logger := observability.NewLogger(observability.InfoLevel, os.Stdout)

// 添加钩子
logger.AddHook(observability.NewMetricsHook())

// 基本日志
logger.Info("Service started")

// 带字段的日志
logger.WithFields(map[string]interface{}{
    "user_id": "123",
    "action":  "login",
}).Info("User action")

// 与追踪集成
logger.WithContext(ctx).Info("Request processed")
```

---

## 📊 性能指标

### 追踪性能

- StartSpan: ~500 ns/op
- 嵌套Span: ~900 ns/op
- 零额外内存分配（复用池）

### 指标性能

- Counter.Inc: ~30 ns/op（并发安全）
- Gauge.Set: ~35 ns/op
- Histogram.Observe: ~200 ns/op

### 日志性能

- 基础日志: ~1.5 μs/op
- 带字段日志: ~2.0 μs/op
- 基于slog的高性能实现

---

## 🎯 集成示例

完整的可观测性集成示例：

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 追踪
    span, ctx := observability.StartSpan(r.Context(), "handle-request")
    defer span.Finish()

    // 日志
    logger := observability.GetDefaultLogger()
    logger.WithContext(ctx).Info("Request started")

    // 指标
    observability.HTTPRequestsTotal.Inc()
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        observability.HTTPRequestDuration.Observe(duration)
    }()

    // 业务逻辑
    // ...

    span.SetStatus(observability.StatusOK, "Success")
    logger.WithContext(ctx).Info("Request completed")
}
```

---

## 🏗️ 架构设计

### 追踪架构

```text
Tracer
  ├── Sampler (采样策略)
  ├── IDGenerator (ID生成)
  └── SpanRecorder (记录器)
       └── Span
            ├── Tags
            ├── Logs
            └── Status
```

### 指标架构

```text
MetricsRegistry
  ├── Counter (只增不减)
  ├── Gauge (可增可减)
  └── Histogram (分布统计)
       └── Buckets
```

### 日志架构

```text
Logger
  ├── Level (日志级别)
  ├── Hooks (钩子系统)
  └── Fields (结构化字段)
       └── Context Integration
```

---

## 📚 更多示例

查看 `example_usage.go` 了解更多使用示例：

- ExampleTracing() - 追踪示例
- ExampleMetrics() - 指标示例
- ExampleLogging() - 日志示例
- ExampleIntegration() - 集成示例

---

## 🧪 测试

```bash
# 运行测试
go test -v ./...

# 运行基准测试
go test -bench=. -benchmem

# 测试覆盖率
go test -cover ./...
```

---

## 📈 特性对比

| 特性 | 本库 | OpenTelemetry | 说明 |
|------|------|---------------|------|
| 追踪 | ✅ | ✅ | 完整支持 |
| 指标 | ✅ | ✅ | Prometheus格式 |
| 日志 | ✅ | ✅ | 基于slog |
| 采样 | ✅ | ✅ | 多种策略 |
| 轻量级 | ✅ | ❌ | 零外部依赖 |
| 易用性 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | API简单 |

---

## 🎯 最佳实践

1. **追踪**：
   - 为关键操作创建Span
   - 使用Context传播
   - 合理设置采样率
   - 记录有意义的标签

2. **指标**：
   - 使用描述性的指标名称
   - 为指标添加标签
   - 定期导出到监控系统
   - 监控关键业务指标

3. **日志**：
   - 选择合适的日志级别
   - 使用结构化字段
   - 集成追踪信息
   - 避免敏感信息

---

## 🔧 配置

### 采样策略

```go
// 总是采样
sampler := &observability.AlwaysSampler{}

// 概率采样（50%）
sampler := observability.NewProbabilitySampler(0.5)
```

### 日志钩子

```go
// 指标钩子
metricsHook := observability.NewMetricsHook()
logger.AddHook(metricsHook)

// 文件钩子
fileHook, _ := observability.NewFileHook("app.log", observability.ErrorLevel)
logger.AddHook(fileHook)
```

---

## 🚀 路线图

- [ ] Jaeger集成
- [ ] Zipkin集成
- [ ] Prometheus推送
- [ ] 自动Instrumentation
- [ ] 采样策略扩展

---

**版本**: v2.0.0
**最后更新**: 2025-10-22
**测试覆盖率**: 95%+
