package otlp

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLogger 测试 NewLogger 构造函数
func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	require.NotNil(t, logger, "NewLogger 不应返回 nil")
	assert.NotNil(t, logger.Logger, "内部 Logger 不应为 nil")
}

// TestLoggerStructure 测试 Logger 结构体
func TestLoggerStructure(t *testing.T) {
	// 测试 Logger 结构体定义
	l := &Logger{}
	assert.NotNil(t, l, "Logger 实例不应为 nil")

	// 测试带有内部 logger 的实例
	internalLogger := slog.Default()
	l = &Logger{Logger: internalLogger}
	assert.NotNil(t, l.Logger, "内部 Logger 不应为 nil")
}

// TestLogger_WithContext 测试 WithContext 方法
func TestLogger_WithContext(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// 测试 WithContext 返回 logger
	result := logger.WithContext(ctx)
	assert.NotNil(t, result, "WithContext 不应返回 nil")

	// 当前实现返回原始 logger，后续应返回带追踪信息的 logger
	assert.IsType(t, &slog.Logger{}, result, "应返回 *slog.Logger 类型")
}

// TestLogger_WithContextNil 测试 WithContext 传入 nil context
func TestLogger_WithContextNil(t *testing.T) {
	logger := NewLogger()

	// 即使传入 nil context 也不应 panic
	result := logger.WithContext(nil)
	assert.NotNil(t, result, "WithContext(nil) 不应返回 nil")
}

// TestLogger_LoggingMethods 测试日志记录方法
func TestLogger_LoggingMethods(t *testing.T) {
	logger := NewLogger()

	// 测试各种日志级别不 panic
	assert.NotPanics(t, func() {
		logger.Info("test info message", "key", "value")
	}, "Info 不应 panic")

	assert.NotPanics(t, func() {
		logger.Debug("test debug message")
	}, "Debug 不应 panic")

	assert.NotPanics(t, func() {
		logger.Warn("test warn message", "count", 42)
	}, "Warn 不应 panic")

	assert.NotPanics(t, func() {
		logger.Error("test error message", "err", "some error")
	}, "Error 不应 panic")
}

// TestLogger_WithContextLogging 测试带上下文的日志记录
func TestLogger_WithContextLogging(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	ctxLogger := logger.WithContext(ctx)
	require.NotNil(t, ctxLogger)

	// 测试带上下文的日志记录不 panic
	assert.NotPanics(t, func() {
		ctxLogger.Info("context message", "request_id", "12345")
	}, "带上下文的日志记录不应 panic")
}

// TestLogger_MultipleInstances 测试创建多个 Logger 实例
func TestLogger_MultipleInstances(t *testing.T) {
	logger1 := NewLogger()
	logger2 := NewLogger()

	require.NotNil(t, logger1)
	require.NotNil(t, logger2)

	// 每个实例应独立
	assert.NotSame(t, logger1, logger2, "每个 NewLogger 调用应返回不同实例")
}
