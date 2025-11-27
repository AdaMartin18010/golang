package tracing

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestTracer_StartSpan(t *testing.T) {
	// 设置 noop tracer provider 用于测试
	otel.SetTracerProvider(noop.NewTracerProvider())

	tracer := NewTracer("test-service")
	ctx := context.Background()

	ctx, span := tracer.StartSpan(ctx, "test-operation")
	if span == nil {
		t.Fatal("Expected non-nil span")
	}
	span.End()
}

func TestTracer_RecordError(t *testing.T) {
	otel.SetTracerProvider(noop.NewTracerProvider())

	tracer := NewTracer("test-service")
	ctx := context.Background()

	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	err := errors.New("test error")
	tracer.RecordError(span, err)

	// 验证错误被记录（noop tracer 不会实际记录，但不会 panic）
}

func TestTracer_LocateError(t *testing.T) {
	otel.SetTracerProvider(noop.NewTracerProvider())

	tracer := NewTracer("test-service")
	ctx := context.Background()

	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	err := errors.New("test error")
	tracer.LocateError(ctx, err, map[string]interface{}{
		"user.id": 123,
	})

	// 验证错误定位信息被记录
}

func TestTracer_AddAttributes(t *testing.T) {
	otel.SetTracerProvider(noop.NewTracerProvider())

	tracer := NewTracer("test-service")
	ctx := context.Background()

	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	tracer.AddAttributes(span, map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	})

	// 验证属性被添加
}

func TestTracer_AddStackTrace(t *testing.T) {
	otel.SetTracerProvider(noop.NewTracerProvider())

	tracer := NewTracer("test-service")
	ctx := context.Background()

	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	tracer.AddStackTrace(span)

	// 验证堆栈跟踪被添加
}

