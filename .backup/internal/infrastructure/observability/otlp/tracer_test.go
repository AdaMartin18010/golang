package otlp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

// TestNewTracerProvider 测试 NewTracerProvider 构造函数
func TestNewTracerProvider(t *testing.T) {
	ctx := context.Background()
	endpoint := "localhost:4317"
	insecure := true

	shutdown, err := NewTracerProvider(ctx, endpoint, insecure)

	require.NoError(t, err, "NewTracerProvider 不应返回错误")
	require.NotNil(t, shutdown, "NewTracerProvider 应返回关闭函数")

	// 清理
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = shutdown(shutdownCtx)
}

// TestNewTracerProvider_InvalidEndpoint 测试无效端点
func TestNewTracerProvider_InvalidEndpoint(t *testing.T) {
	ctx := context.Background()

	// 测试无效端点
	shutdown, err := NewTracerProvider(ctx, "invalid-endpoint", true)
	// 行为取决于实现，但不应 panic
	if err == nil && shutdown != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}
}

// TestNewTracerProvider_SecureConnection 测试安全连接
func TestNewTracerProvider_SecureConnection(t *testing.T) {
	ctx := context.Background()

	// 测试 insecure=false (TLS)
	shutdown, err := NewTracerProvider(ctx, "localhost:4317", false)
	// 可能因为没有 TLS 配置而失败，但不应 panic
	if err == nil && shutdown != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}
}

// TestNewTracerProvider_GlobalProvider 测试全局 Provider 设置
func TestNewTracerProvider_GlobalProvider(t *testing.T) {
	ctx := context.Background()

	shutdown, err := NewTracerProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}()

	// 验证全局 TracerProvider 已设置
	tracer := otel.Tracer("test-service")
	assert.NotNil(t, tracer, "全局 Tracer 不应为 nil")
}

// TestNewTracerProvider_Shutdown 测试关闭函数
func TestNewTracerProvider_Shutdown(t *testing.T) {
	ctx := context.Background()

	shutdown, err := NewTracerProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)

	// 测试关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = shutdown(shutdownCtx)
	assert.NoError(t, err, "关闭函数不应返回错误")
}

// TestNewTracerProvider_MultipleShutdown 测试多次关闭
func TestNewTracerProvider_MultipleShutdown(t *testing.T) {
	ctx := context.Background()

	shutdown, err := NewTracerProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 第一次关闭
	err = shutdown(shutdownCtx)
	require.NoError(t, err)

	// 第二次关闭 - 不应 panic
	_ = shutdown(shutdownCtx)
}

// TestNewTracerProvider_WithContextCancellation 测试上下文取消
func TestNewTracerProvider_WithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	shutdown, err := NewTracerProvider(ctx, "localhost:4317", true)
	if err == nil && shutdown != nil {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = shutdown(shutdownCtx)
		}()
	}

	// 取消上下文
	cancel()
}

// TestNewTracerProvider_TracerCreation 测试创建 Tracer
func TestNewTracerProvider_TracerCreation(t *testing.T) {
	ctx := context.Background()

	shutdown, err := NewTracerProvider(ctx, "localhost:4317", true)
	require.NoError(t, err)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}()

	// 创建 Tracer
	tracer := otel.Tracer("test-service")
	require.NotNil(t, tracer)

	// 启动 Span
	ctx, span := tracer.Start(ctx, "test-operation")
	assert.NotNil(t, ctx, "Start 应返回 context")
	assert.NotNil(t, span, "Start 应返回 span")

	defer span.End()
}

// TestNewTracerProvider_EmptyEndpoint 测试空端点
func TestNewTracerProvider_EmptyEndpoint(t *testing.T) {
	ctx := context.Background()

	// 空端点可能导致错误，但不应 panic
	shutdown, err := NewTracerProvider(ctx, "", true)
	if err == nil && shutdown != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}
}
