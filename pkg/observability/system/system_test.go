package system

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

func TestNewSystemMonitor(t *testing.T) {
	ctx := context.Background()

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	// 创建 MeterProvider
	mp := metric.NewMeterProvider(metric.WithResource(res))

	// 创建系统监控器
	monitor, err := NewSystemMonitor(SystemConfig{
		Meter:           mp.Meter("test"),
		Enabled:         true,
		CollectInterval: 1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create system monitor: %v", err)
	}

	if monitor == nil {
		t.Fatal("System monitor is nil")
	}
}

func TestSystemMonitor_StartStop(t *testing.T) {
	ctx := context.Background()

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	// 创建 MeterProvider
	mp := metric.NewMeterProvider(metric.WithResource(res))

	// 创建系统监控器
	monitor, err := NewSystemMonitor(SystemConfig{
		Meter:           mp.Meter("test"),
		Enabled:         true,
		CollectInterval: 1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create system monitor: %v", err)
	}

	// 启动监控
	if err := monitor.Start(); err != nil {
		t.Fatalf("Failed to start monitor: %v", err)
	}

	// 等待一段时间
	time.Sleep(100 * time.Millisecond)

	// 停止监控
	if err := monitor.Stop(); err != nil {
		t.Fatalf("Failed to stop monitor: %v", err)
	}
}

func TestSystemMonitor_GetPlatformInfo(t *testing.T) {
	ctx := context.Background()

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	// 创建 MeterProvider
	mp := metric.NewMeterProvider(metric.WithResource(res))

	// 创建系统监控器
	monitor, err := NewSystemMonitor(SystemConfig{
		Meter:           mp.Meter("test"),
		Enabled:         true,
		CollectInterval: 1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create system monitor: %v", err)
	}

	// 获取平台信息
	info := monitor.GetPlatformInfo()
	if info.OS == "" {
		t.Error("Platform OS should not be empty")
	}
}

func TestSystemMonitor_HealthCheck(t *testing.T) {
	ctx := context.Background()

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	// 创建 MeterProvider
	mp := metric.NewMeterProvider(metric.WithResource(res))

	// 创建系统监控器
	monitor, err := NewSystemMonitor(SystemConfig{
		Meter:           mp.Meter("test"),
		Enabled:         true,
		CollectInterval: 1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create system monitor: %v", err)
	}

	// 启动监控
	if err := monitor.Start(); err != nil {
		t.Fatalf("Failed to start monitor: %v", err)
	}
	defer monitor.Stop()

	// 执行健康检查
	status := monitor.CheckHealth(ctx)
	if status.Timestamp.IsZero() {
		t.Error("Health status timestamp should not be zero")
	}
}
