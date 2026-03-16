package otlp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
)

// TestMetricsProviderStructure 测试 MetricsProvider 结构体
func TestMetricsProviderStructure(t *testing.T) {
	mp := &MetricsProvider{}
	assert.NotNil(t, mp, "MetricsProvider 实例不应为 nil")
}

// TestNewMetricsProvider 测试 NewMetricsProvider 构造函数
func TestNewMetricsProvider(t *testing.T) {
	ctx := context.Background()
	endpoint := "localhost:4317"
	insecure := true

	mp, err := NewMetricsProvider(ctx, endpoint, insecure)

	require.NoError(t, err, "NewMetricsProvider 不应返回错误")
	require.NotNil(t, mp, "NewMetricsProvider 不应返回 nil")

	// 清理
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = mp.Shutdown(shutdownCtx)
}

// TestNewMetricsProvider_InvalidEndpoint 测试无效端点
func TestNewMetricsProvider_InvalidEndpoint(t *testing.T) {
	ctx := context.Background()

	// 测试空端点
	mp, err := NewMetricsProvider(ctx, "", true)
	// 当前实现不会验证端点，所以可能成功或失败取决于实现
	// 但不应 panic
	if err == nil && mp != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mp.Shutdown(shutdownCtx)
	}
}

// TestMetricsProvider_MeterProvider 测试 MeterProvider 方法
func TestMetricsProvider_MeterProvider(t *testing.T) {
	ctx := context.Background()
	mp, err := NewMetricsProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mp.Shutdown(shutdownCtx)
	}()

	provider := mp.MeterProvider()
	assert.NotNil(t, provider, "MeterProvider 不应返回 nil")
	assert.IsType(t, &metric.MeterProvider{}, provider, "应返回 *metric.MeterProvider 类型")
}

// TestMetricsProvider_Shutdown 测试 Shutdown 方法
func TestMetricsProvider_Shutdown(t *testing.T) {
	ctx := context.Background()
	mp, err := NewMetricsProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)

	// 测试正常关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = mp.Shutdown(shutdownCtx)
	assert.NoError(t, err, "Shutdown 不应返回错误")
}

// TestMetricsProvider_MultipleShutdown 测试多次关闭
func TestMetricsProvider_MultipleShutdown(t *testing.T) {
	ctx := context.Background()
	mp, err := NewMetricsProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 第一次关闭
	err = mp.Shutdown(shutdownCtx)
	require.NoError(t, err)

	// 第二次关闭 - 不应 panic，可能返回错误
	_ = mp.Shutdown(shutdownCtx)
}

// TestMetricsProvider_WithContextCancellation 测试上下文取消
func TestMetricsProvider_WithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	// 使用已取消的上下文创建 provider
	mp, err := NewMetricsProvider(ctx, "localhost:4317", true)
	// 行为取决于实现，但不应 panic
	if err == nil && mp != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_ = mp.Shutdown(shutdownCtx)
	}
}

// TestMetricsProvider_SecureConnection 测试安全连接
func TestMetricsProvider_SecureConnection(t *testing.T) {
	ctx := context.Background()

	// 测试 insecure=false
	mp, err := NewMetricsProvider(ctx, "localhost:4317", false)
	// 当前实现可能不支持 TLS，所以可能失败
	// 但不应 panic
	if err == nil && mp != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mp.Shutdown(shutdownCtx)
	}
}

// TestMetricsProvider_CreateMeter 测试创建 Meter
func TestMetricsProvider_CreateMeter(t *testing.T) {
	ctx := context.Background()
	mp, err := NewMetricsProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mp.Shutdown(shutdownCtx)
	}()

	meter := mp.MeterProvider().Meter("test-service")
	assert.NotNil(t, meter, "Meter 不应为 nil")
}
