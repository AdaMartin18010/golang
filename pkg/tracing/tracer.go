package tracing

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracer 追踪器
// 提供分布式追踪和错误定位能力
type Tracer struct {
	tracer trace.Tracer
}

// NewTracer 创建追踪器
func NewTracer(name string) *Tracer {
	return &Tracer{
		tracer: otel.Tracer(name),
	}
}

// StartSpan 开始一个新的 Span
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// StartSpanWithAttributes 开始一个带属性的 Span
func (t *Tracer) StartSpanWithAttributes(ctx context.Context, name string, attrs map[string]interface{}) (context.Context, trace.Span) {
	spanOpts := []trace.SpanStartOption{
		trace.WithAttributes(convertAttributes(attrs)...),
	}
	return t.tracer.Start(ctx, name, spanOpts...)
}

// RecordError 记录错误
func (t *Tracer) RecordError(span trace.Span, err error, attrs ...attribute.KeyValue) {
	if err == nil {
		return
	}

	span.RecordError(err, trace.WithAttributes(attrs...))
	span.SetStatus(codes.Error, err.Error())
}

// RecordPanic 记录 panic
func (t *Tracer) RecordPanic(span trace.Span, r interface{}) {
	span.RecordError(fmt.Errorf("panic: %v", r), trace.WithAttributes(
		attribute.String("event", "panic"),
		attribute.String("panic.value", fmt.Sprintf("%v", r)),
	))
	span.SetStatus(codes.Error, fmt.Sprintf("panic: %v", r))
}

// AddAttributes 添加属性
func (t *Tracer) AddAttributes(span trace.Span, attrs map[string]interface{}) {
	span.SetAttributes(convertAttributes(attrs)...)
}

// AddStackTrace 添加堆栈跟踪
func (t *Tracer) AddStackTrace(span trace.Span) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	span.SetAttributes(attribute.String("stack.trace", string(buf[:n])))
}

// LocateError 定位错误
// 记录错误的完整上下文信息，包括堆栈跟踪、调用位置等
func (t *Tracer) LocateError(ctx context.Context, err error, attrs map[string]interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	// 记录错误
	t.RecordError(span, err, convertAttributes(attrs)...)

	// 添加堆栈跟踪
	t.AddStackTrace(span)

	// 添加调用位置信息
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		fn := runtime.FuncForPC(pc)
		span.SetAttributes(
			attribute.String("error.location.file", file),
			attribute.Int("error.location.line", line),
			attribute.String("error.location.function", fn.Name()),
		)
	}
}

// convertAttributes 转换属性
func convertAttributes(attrs map[string]interface{}) []attribute.KeyValue {
	result := make([]attribute.KeyValue, 0, len(attrs))
	for k, v := range attrs {
		result = append(result, convertAttribute(k, v))
	}
	return result
}

// convertAttribute 转换单个属性
func convertAttribute(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case int:
		return attribute.Int(key, v)
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	case []string:
		return attribute.StringSlice(key, v)
	case []int:
		return attribute.IntSlice(key, v)
	case []int64:
		return attribute.Int64Slice(key, v)
	case []float64:
		return attribute.Float64Slice(key, v)
	case []bool:
		return attribute.BoolSlice(key, v)
	case time.Time:
		return attribute.String(key, v.Format(time.RFC3339))
	case time.Duration:
		return attribute.String(key, v.String())
	default:
		return attribute.String(key, fmt.Sprintf("%v", v))
	}
}

// Span 包装器，提供便捷方法
type Span struct {
	trace.Span
	tracer *Tracer
}

// NewSpan 创建 Span 包装器
func NewSpan(span trace.Span, tracer *Tracer) *Span {
	return &Span{
		Span:   span,
		tracer: tracer,
	}
}

// RecordError 记录错误
func (s *Span) RecordError(err error, attrs ...attribute.KeyValue) {
	s.tracer.RecordError(s.Span, err, attrs...)
}

// RecordPanic 记录 panic
func (s *Span) RecordPanic(r interface{}) {
	s.tracer.RecordPanic(s.Span, r)
}

// AddAttributes 添加属性
func (s *Span) AddAttributes(attrs map[string]interface{}) {
	s.tracer.AddAttributes(s.Span, attrs)
}

// AddStackTrace 添加堆栈跟踪
func (s *Span) AddStackTrace() {
	s.tracer.AddStackTrace(s.Span)
}

// SetStatus 设置状态
func (s *Span) SetStatus(code codes.Code, description string) {
	s.Span.SetStatus(code, description)
}

// End 结束 Span
func (s *Span) End(opts ...trace.SpanEndOption) {
	s.Span.End(opts...)
}
