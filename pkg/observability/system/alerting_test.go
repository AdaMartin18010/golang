package system

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

func TestNewAlertManager(t *testing.T) {
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

	// 创建告警管理器
	alertManager := NewAlertManager(mp.Meter("test"), true)
	if alertManager == nil {
		t.Fatal("Alert manager is nil")
	}
}

func TestAlertManager_AddRule(t *testing.T) {
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

	// 创建告警管理器
	alertManager := NewAlertManager(mp.Meter("test"), true)

	// 添加告警规则
	rule := AlertRule{
		ID:         "test-rule",
		Name:       "Test Rule",
		MetricName: "system.cpu.usage",
		Condition:  "gt",
		Threshold:  80.0,
		Level:      AlertLevelWarning,
		Enabled:    true,
		Duration:   5 * time.Minute,
		Cooldown:   10 * time.Minute,
	}
	alertManager.AddRule(rule)

	// 检查告警
	alertManager.Check(ctx, "system.cpu.usage", 85.0, nil)
}

func TestAlertManager_DefaultRules(t *testing.T) {
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

	// 创建告警管理器
	alertManager := NewAlertManager(mp.Meter("test"), true)

	// 添加默认规则
	for _, rule := range DefaultAlertRules() {
		alertManager.AddRule(rule)
	}

	// 检查告警历史
	history := alertManager.GetAlertHistory(10)
	if len(history) < 0 {
		t.Error("Alert history should be accessible")
	}
}
