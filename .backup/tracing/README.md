# 追踪和定位

框架级别的分布式追踪和错误定位能力，提供完整的追踪上下文和错误定位信息。

## 📋 功能特性

- ✅ **分布式追踪**: 基于 OpenTelemetry 的分布式追踪
- ✅ **错误定位**: 自动记录错误的完整上下文（堆栈跟踪、调用位置）
- ✅ **属性记录**: 支持丰富的属性记录
- ✅ **Panic 捕获**: 自动捕获和记录 panic

## 🚀 快速开始

### 基本追踪

```go
import "github.com/yourusername/golang/pkg/tracing"

tracer := tracing.NewTracer("my-service")

// 开始 Span
ctx, span := tracer.StartSpan(ctx, "operation-name")
defer span.End()

// 添加属性
tracer.AddAttributes(span, map[string]interface{}{
    "user.id": 123,
    "operation.type": "create",
})

// 记录错误
if err != nil {
    tracer.RecordError(span, err)
}
```

### 错误定位

```go
// 自动记录错误的完整上下文
tracer.LocateError(ctx, err, map[string]interface{}{
    "user.id": 123,
    "request.id": "req-456",
})
// 这会自动记录：
// - 错误信息
// - 堆栈跟踪
// - 调用位置（文件、行号、函数名）
// - 自定义属性
```

### Panic 捕获

```go
ctx, span := tracer.StartSpan(ctx, "risky-operation")
defer func() {
    if r := recover(); r != nil {
        tracer.RecordPanic(span, r)
        span.End()
        panic(r) // 重新抛出
    }
    span.End()
}()

// 可能 panic 的操作
riskyOperation()
```

## 📚 API 参考

### Tracer

```go
type Tracer struct {
    tracer trace.Tracer
}

func NewTracer(name string) *Tracer
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
func (t *Tracer) RecordError(span trace.Span, err error, attrs ...attribute.KeyValue)
func (t *Tracer) LocateError(ctx context.Context, err error, attrs map[string]interface{})
```

## 🔗 相关文档

- [OpenTelemetry 集成](../observability/README.md)
- [采样机制](../sampling/README.md)
