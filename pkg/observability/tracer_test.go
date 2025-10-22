package observability

import (
	"context"
	"testing"
	"time"
)

func TestTracerBasic(t *testing.T) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()
	span, ctx := tracer.StartSpan(ctx, "test-operation")
	if span == nil {
		t.Fatal("Expected span, got nil")
	}

	if span.Name != "test-operation" {
		t.Errorf("Expected span name 'test-operation', got '%s'", span.Name)
	}

	if span.TraceID == "" {
		t.Error("Expected non-empty TraceID")
	}

	if span.SpanID == "" {
		t.Error("Expected non-empty SpanID")
	}

	// 添加小延迟以确保duration不为0
	time.Sleep(time.Millisecond)

	span.Finish()

	if span.Duration == 0 {
		t.Error("Expected non-zero duration")
	}
}

func TestSpanTags(t *testing.T) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()
	span, _ := tracer.StartSpan(ctx, "test-operation")

	span.SetTag("http.method", "GET")
	span.SetTag("http.url", "/api/test")

	if span.Tags["http.method"] != "GET" {
		t.Errorf("Expected tag 'http.method' to be 'GET', got '%s'", span.Tags["http.method"])
	}

	if span.Tags["http.url"] != "/api/test" {
		t.Errorf("Expected tag 'http.url' to be '/api/test', got '%s'", span.Tags["http.url"])
	}

	span.Finish()
}

func TestSpanLogs(t *testing.T) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()
	span, _ := tracer.StartSpan(ctx, "test-operation")

	span.LogFields(map[string]interface{}{
		"event":   "cache_hit",
		"key":     "user:123",
		"latency": 5,
	})

	if len(span.Logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(span.Logs))
	}

	log := span.Logs[0]
	if log.Fields["event"] != "cache_hit" {
		t.Errorf("Expected event 'cache_hit', got '%v'", log.Fields["event"])
	}

	span.Finish()
}

func TestSpanError(t *testing.T) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()
	span, _ := tracer.StartSpan(ctx, "test-operation")

	// 模拟错误
	err := context.DeadlineExceeded
	span.SetError(err)

	if span.Status.Code != StatusError {
		t.Errorf("Expected status code %d, got %d", StatusError, span.Status.Code)
	}

	if span.Tags["error"] != "true" {
		t.Errorf("Expected error tag to be 'true', got '%s'", span.Tags["error"])
	}

	if span.Status.Message != err.Error() {
		t.Errorf("Expected status message '%s', got '%s'", err.Error(), span.Status.Message)
	}

	span.Finish()
}

func TestNestedSpans(t *testing.T) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()

	// 父span
	parentSpan, ctx := tracer.StartSpan(ctx, "parent-operation")
	if parentSpan == nil {
		t.Fatal("Expected parent span, got nil")
	}

	// 子span
	childSpan, _ := tracer.StartSpan(ctx, "child-operation")
	if childSpan == nil {
		t.Fatal("Expected child span, got nil")
	}

	// 验证继承关系
	if childSpan.TraceID != parentSpan.TraceID {
		t.Errorf("Expected child TraceID to match parent, got parent=%s, child=%s",
			parentSpan.TraceID, childSpan.TraceID)
	}

	if childSpan.ParentID != parentSpan.SpanID {
		t.Errorf("Expected child ParentID to match parent SpanID, got parent=%s, child=%s",
			parentSpan.SpanID, childSpan.ParentID)
	}

	childSpan.Finish()
	parentSpan.Finish()
}

func TestProbabilitySampler(t *testing.T) {
	tests := []struct {
		name        string
		probability float64
	}{
		{"Always sample", 1.0},
		{"Never sample", 0.0},
		{"50% sample", 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := NewProbabilitySampler(tt.probability)

			samples := 0
			total := 1000

			for i := 0; i < total; i++ {
				if sampler.ShouldSample("test") {
					samples++
				}
			}

			// 允许一定误差
			if tt.probability == 1.0 && samples != total {
				t.Errorf("Expected all samples to be taken, got %d/%d", samples, total)
			}

			if tt.probability == 0.0 && samples != 0 {
				t.Errorf("Expected no samples to be taken, got %d/%d", samples, total)
			}

			// 对于0.5，允许40%-60%的范围
			if tt.probability == 0.5 {
				ratio := float64(samples) / float64(total)
				if ratio < 0.4 || ratio > 0.6 {
					t.Errorf("Expected ~50%% samples, got %.2f%% (%d/%d)",
						ratio*100, samples, total)
				}
			}
		})
	}
}

func TestContextPropagation(t *testing.T) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	ctx := context.Background()

	// 创建span
	span, ctx := tracer.StartSpan(ctx, "test-operation")
	if span == nil {
		t.Fatal("Expected span, got nil")
	}

	// 从context获取span
	retrievedSpan := SpanFromContext(ctx)
	if retrievedSpan == nil {
		t.Fatal("Expected to retrieve span from context, got nil")
	}

	if retrievedSpan.SpanID != span.SpanID {
		t.Errorf("Expected SpanID %s, got %s", span.SpanID, retrievedSpan.SpanID)
	}

	span.Finish()
}

func TestGlobalTracer(t *testing.T) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	// 设置全局tracer
	SetGlobalTracer(tracer)

	// 使用全局tracer
	ctx := context.Background()
	span, ctx := StartSpan(ctx, "global-operation")
	if span == nil {
		t.Fatal("Expected span from global tracer, got nil")
	}

	if span.Name != "global-operation" {
		t.Errorf("Expected span name 'global-operation', got '%s'", span.Name)
	}

	span.Finish()
}

func TestInMemoryRecorder(t *testing.T) {
	recorder := NewInMemoryRecorder()

	span := &Span{
		TraceID:   "trace-123",
		SpanID:    "span-123",
		Name:      "test-operation",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Millisecond * 100),
		Duration:  time.Millisecond * 100,
		Tags:      map[string]string{"key": "value"},
	}

	recorder.RecordSpan(span)

	spans := recorder.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("Expected 1 span, got %d", len(spans))
	}

	if spans[0].TraceID != "trace-123" {
		t.Errorf("Expected TraceID 'trace-123', got '%s'", spans[0].TraceID)
	}

	recorder.Clear()

	spans = recorder.GetSpans()
	if len(spans) != 0 {
		t.Errorf("Expected 0 spans after clear, got %d", len(spans))
	}
}

func BenchmarkStartSpan(b *testing.B) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span, _ := tracer.StartSpan(ctx, "test-operation")
		span.Finish()
	}
}

func BenchmarkSpanWithTags(b *testing.B) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span, _ := tracer.StartSpan(ctx, "test-operation")
		span.SetTag("key1", "value1")
		span.SetTag("key2", "value2")
		span.SetTag("key3", "value3")
		span.Finish()
	}
}

func BenchmarkNestedSpans(b *testing.B) {
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("test-service", recorder, sampler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		parentSpan, ctx := tracer.StartSpan(ctx, "parent")
		childSpan, _ := tracer.StartSpan(ctx, "child")
		childSpan.Finish()
		parentSpan.Finish()
	}
}
