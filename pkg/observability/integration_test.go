package observability

import (
	"context"
	"net"
	"testing"
	"time"
)

// otlpCollectorAvailable 探测本地 OTLP collector（gRPC 4317）是否可达；
// 不可达时跳过依赖外部服务的集成测试
func otlpCollectorAvailable(t *testing.T) {
	t.Helper()
	conn, err := net.DialTimeout("tcp", "localhost:4317", 200*time.Millisecond)
	if err != nil {
		t.Skip("OTLP collector not available")
	}
	_ = conn.Close()
}

func TestNewObservability(t *testing.T) {
	obs, err := NewObservability(Config{
		ServiceName:            "test-service",
		ServiceVersion:         "v1.0.0",
		OTLPEndpoint:           "localhost:4317",
		OTLPInsecure:           true,
		SampleRate:             0.5,
		EnableSystemMonitoring: false, // 禁用系统监控以避免依赖问题
	})
	if err != nil {
		t.Fatalf("Failed to create observability: %v", err)
	}

	if obs == nil {
		t.Fatal("Observability is nil")
	}
}

func TestObservability_StartStop(t *testing.T) {
	otlpCollectorAvailable(t)

	ctx := context.Background()

	obs, err := NewObservability(Config{
		ServiceName:            "test-service",
		ServiceVersion:         "v1.0.0",
		OTLPEndpoint:           "localhost:4317",
		OTLPInsecure:           true,
		SampleRate:             0.5,
		EnableSystemMonitoring: false,
	})
	if err != nil {
		t.Fatalf("Failed to create observability: %v", err)
	}

	// 启动
	if err := obs.Start(); err != nil {
		t.Fatalf("Failed to start observability: %v", err)
	}

	// 等待一段时间
	time.Sleep(100 * time.Millisecond)

	// 停止
	if err := obs.Stop(ctx); err != nil {
		t.Fatalf("Failed to stop observability: %v", err)
	}
}

func TestObservability_TracerMeter(t *testing.T) {
	ctx := context.Background()

	obs, err := NewObservability(Config{
		ServiceName:            "test-service",
		ServiceVersion:         "v1.0.0",
		OTLPEndpoint:           "localhost:4317",
		OTLPInsecure:           true,
		SampleRate:             0.5,
		EnableSystemMonitoring: false,
	})
	if err != nil {
		t.Fatalf("Failed to create observability: %v", err)
	}

	// 获取追踪器
	tracer := obs.Tracer("test-tracer")
	if tracer == nil {
		t.Error("Tracer should not be nil")
	}

	// 获取指标器
	meter := obs.Meter("test-meter")
	if meter == nil {
		t.Error("Meter should not be nil")
	}

	// 使用追踪器
	_, span := tracer.Start(ctx, "test-span")
	if span == nil {
		t.Error("Span should not be nil")
	}
	span.End()
}
