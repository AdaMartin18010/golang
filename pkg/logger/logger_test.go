package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(slog.LevelInfo)
	if logger == nil {
		t.Error("Expected logger, got nil")
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Level:      slog.LevelDebug,
		Output:     &buf,
		AddSource:  true,
		JSONFormat: true,
	}

	logger := NewLoggerWithConfig(config)
	if logger == nil {
		t.Error("Expected logger, got nil")
	}

	logger.Debug("test message")
	if buf.Len() == 0 {
		t.Error("Expected log output, got empty")
	}
}

func TestLogger_WithContext(t *testing.T) {
	logger := NewLogger(slog.LevelInfo)

	// 测试传统的 context value 方式
	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	loggerWithCtx := logger.WithContext(ctx)

	if loggerWithCtx == nil {
		t.Error("Expected logger with context, got nil")
	}

	// 测试 OpenTelemetry context（需要实际的 span）
	// 这里只测试不会 panic
	ctx2 := context.Background()
	loggerWithCtx2 := logger.WithContext(ctx2)
	if loggerWithCtx2 == nil {
		t.Error("Expected logger with context, got nil")
	}
}

func TestLogger_WithFields(t *testing.T) {
	logger := NewLogger(slog.LevelInfo)
	loggerWithFields := logger.WithFields("key", "value")

	if loggerWithFields == nil {
		t.Error("Expected logger with fields, got nil")
	}
}

func TestLogger_WithError(t *testing.T) {
	logger := NewLogger(slog.LevelInfo)
	err := NewError("test error")
	loggerWithError := logger.WithError(err)

	if loggerWithError == nil {
		t.Error("Expected logger with error, got nil")
	}

	// Test with nil error
	loggerWithNilError := logger.WithError(nil)
	if loggerWithNilError == nil {
		t.Error("Expected logger, got nil")
	}
}

func TestLogger_LogRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithConfig(Config{
		Level:      slog.LevelInfo,
		Output:     &buf,
		JSONFormat: true,
	})

	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	logger.LogRequest(ctx, "GET", "/users", 200, 100*time.Millisecond)

	if buf.Len() == 0 {
		t.Error("Expected log output, got empty")
	}
}

func TestLogger_LogError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithConfig(Config{
		Level:      slog.LevelError,
		Output:     &buf,
		JSONFormat: true,
	})

	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	err := NewError("test error")
	logger.LogError(ctx, err, "operation failed", "key", "value")

	if buf.Len() == 0 {
		t.Error("Expected log output, got empty")
	}
}

// NewError 创建一个简单的错误用于测试
func NewError(msg string) error {
	return &testError{msg: msg}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestLogger_SampleRate(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Level:      slog.LevelDebug,
		Output:     &buf,
		JSONFormat: true,
		SampleRate: 0.0, // 不采样
	}
	logger := NewLoggerWithConfig(config)

	logger.Info("test message")
	if buf.Len() > 0 {
		t.Error("Expected no log output with 0 sample rate, got output")
	}

	// 测试错误日志不受采样影响
	buf.Reset()
	logger.Error("error message")
	if buf.Len() == 0 {
		t.Error("Expected error log output, got empty")
	}
}

func TestLogger_ServiceNameAndVersion(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		Level:          slog.LevelInfo,
		Output:         &buf,
		JSONFormat:     true,
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
	}
	logger := NewLoggerWithConfig(config)

	logger.Info("test message")
	output := buf.String()
	if !contains(output, "test-service") || !contains(output, "1.0.0") {
		t.Error("Expected service name and version in log output")
	}
}

func TestLogger_SetLevel(t *testing.T) {
	logger := NewLogger(slog.LevelInfo)

	if logger.GetLevel() != slog.LevelInfo {
		t.Error("Expected Info level, got different level")
	}

	logger.SetLevel(slog.LevelWarn)
	if logger.GetLevel() != slog.LevelWarn {
		t.Error("Expected Warn level after SetLevel, got different level")
	}
}

func TestLogger_WithRequestID(t *testing.T) {
	logger := NewLogger(slog.LevelInfo)

	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	loggerWithCtx := logger.WithContext(ctx)

	if loggerWithCtx == nil {
		t.Error("Expected logger with context, got nil")
	}
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
