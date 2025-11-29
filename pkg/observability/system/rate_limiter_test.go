package system

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

func TestNewRateLimiter(t *testing.T) {
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

	// 创建限流器
	limiter, err := NewRateLimiter(RateLimiterConfig{
		Meter:   mp.Meter("test"),
		Enabled: true,
		Limit:   100,
		Window:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}

	if limiter == nil {
		t.Fatal("Rate limiter is nil")
	}
}

func TestRateLimiter_Allow(t *testing.T) {
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

	// 创建限流器（限制为每秒 5 个请求）
	limiter, err := NewRateLimiter(RateLimiterConfig{
		Meter:   mp.Meter("test"),
		Enabled: true,
		Limit:   5,
		Window:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}

	// 前 5 个请求应该被允许
	allowed := 0
	for i := 0; i < 5; i++ {
		if limiter.Allow(ctx) {
			allowed++
		}
	}
	if allowed != 5 {
		t.Errorf("Expected 5 allowed requests, got %d", allowed)
	}

	// 第 6 个请求应该被拒绝
	if limiter.Allow(ctx) {
		t.Error("6th request should be rejected")
	}
}

func TestRateLimiter_UpdateLimit(t *testing.T) {
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

	// 创建限流器
	limiter, err := NewRateLimiter(RateLimiterConfig{
		Meter:   mp.Meter("test"),
		Enabled: true,
		Limit:   100,
		Window:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create rate limiter: %v", err)
	}

	// 更新限制
	limiter.UpdateLimit(200)
	if limiter.GetLimit() != 200 {
		t.Errorf("Expected limit 200, got %d", limiter.GetLimit())
	}
}
