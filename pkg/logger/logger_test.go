package logger

import (
	"bytes"
	"context"
	"log/slog"
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

	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	loggerWithCtx := logger.WithContext(ctx)

	if loggerWithCtx == nil {
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
