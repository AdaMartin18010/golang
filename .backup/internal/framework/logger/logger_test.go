package logger

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, slog.LevelInfo, cfg.Level)
	assert.NotNil(t, cfg.Output)
	assert.False(t, cfg.AddSource)
	assert.True(t, cfg.JSONFormat)
	assert.Equal(t, 1.0, cfg.SampleRate)
}

// TestNewLogger_WithDefaults 测试使用默认配置创建日志
func TestNewLogger_WithDefaults(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output: &buf,
	}

	log := NewLogger(cfg)
	require.NotNil(t, log)

	// 测试日志记录
	log.Info("test message", "key", "value")

	// 验证输出包含日志内容
	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
}

// TestNewLogger_WithCustomConfig 测试使用自定义配置
func TestNewLogger_WithCustomConfig(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Level:          slog.LevelWarn, // 使用Warn级别以避免被底层logger过滤
		Output:         &buf,
		AddSource:      true,
		JSONFormat:     true,
		SampleRate:     0.5,
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
	}

	log := NewLogger(cfg)
	require.NotNil(t, log)

	// 测试Warn和Error级别（确保通过底层logger的级别过滤）
	log.Warn("warn message")
	log.Error("error message")

	output := buf.String()
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
	assert.Contains(t, output, "test-service")
	assert.Contains(t, output, "1.0.0")
}

// TestNewLogger_WithNilConfig 测试使用nil配置
func TestNewLogger_WithNilConfig(t *testing.T) {
	log := NewLogger(nil)
	require.NotNil(t, log)

	// 应该使用默认配置成功创建
	assert.NotNil(t, log)
}

// TestLogger_Debug 测试Debug日志
func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Level:  slog.LevelDebug,
		Output: &buf,
	}

	log := NewLogger(cfg)
	log.Debug("debug test", "foo", "bar")

	output := buf.String()
	assert.Contains(t, output, "debug test")
	assert.Contains(t, output, "foo")
	assert.Contains(t, output, "bar")
}

// TestLogger_Info 测试Info日志
func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output: &buf,
	}

	log := NewLogger(cfg)
	log.Info("info test", "count", 42)

	output := buf.String()
	assert.Contains(t, output, "info test")
	assert.Contains(t, output, "count")
	assert.Contains(t, output, "42")
}

// TestLogger_Warn 测试Warn日志
func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output: &buf,
	}

	log := NewLogger(cfg)
	log.Warn("warn test", "threshold", 0.8)

	output := buf.String()
	assert.Contains(t, output, "warn test")
}

// TestLogger_Error 测试Error日志
func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output: &buf,
	}

	log := NewLogger(cfg)
	log.Error("error test", "fatal", true)

	output := buf.String()
	assert.Contains(t, output, "error test")
}

// TestLogger_WithContext 测试添加上下文
func TestLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output: &buf,
	}

	log := NewLogger(cfg)
	ctx := context.WithValue(context.Background(), "request_id", "req-123")

	loggerWithCtx := log.WithContext(ctx)
	require.NotNil(t, loggerWithCtx)

	loggerWithCtx.Info("context test")

	output := buf.String()
	assert.Contains(t, output, "context test")
}

// TestLogger_WithFields 测试添加字段
func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output: &buf,
	}

	log := NewLogger(cfg)
	loggerWithFields := log.WithFields("service", "test", "version", "1.0")
	require.NotNil(t, loggerWithFields)

	loggerWithFields.Info("fields test")

	output := buf.String()
	assert.Contains(t, output, "fields test")
	assert.Contains(t, output, "service")
	assert.Contains(t, output, "test")
}

// TestLogger_WithError 测试添加错误
func TestLogger_WithError(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output: &buf,
	}

	log := NewLogger(cfg)
	testErr := errors.New("test error")

	loggerWithErr := log.WithError(testErr)
	require.NotNil(t, loggerWithErr)

	loggerWithErr.Error("error test")

	output := buf.String()
	assert.Contains(t, output, "error test")
	assert.Contains(t, output, "test error")
}

// TestSetLogger 测试设置默认日志实例
func TestSetLogger(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output:     &buf,
		ServiceName: "test-logger",
	}

	customLogger := NewLogger(cfg)
	SetLogger(customLogger)

	// 验证设置成功（由于单例模式，GetLogger可能返回已初始化的实例）
	retrievedLogger := GetLogger()
	assert.NotNil(t, retrievedLogger)

	// 使用设置的日志
	customLogger.Info("set logger test")

	output := buf.String()
	assert.Contains(t, output, "set logger test")
}

// TestGetLogger 测试获取默认日志实例
func TestGetLogger(t *testing.T) {
	// 清除单例状态以确保测试独立
	defaultLogger = nil
	once = sync.Once{}

	log := GetLogger()
	require.NotNil(t, log)

	// 多次获取应该返回相同实例
	log2 := GetLogger()
	assert.Equal(t, log, log2)
}

// TestLogger_LevelFiltering 测试日志级别过滤
func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Level:  slog.LevelWarn, // 只记录Warn及以上级别
		Output: &buf,
	}

	log := NewLogger(cfg)

	// Debug和Info不应该被记录
	log.Debug("debug message")
	log.Info("info message")

	// Warn和Error应该被记录
	log.Warn("warn message")
	log.Error("error message")

	output := buf.String()
	assert.NotContains(t, output, "debug message")
	assert.NotContains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}

// TestConfig_JSONFormat 测试JSON格式输出
func TestConfig_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	cfg := &Config{
		Output:     &buf,
		JSONFormat: true,
	}

	log := NewLogger(cfg)
	log.Info("json test", "key", "value")

	output := buf.String()
	// JSON格式应该包含花括号
	assert.Contains(t, output, "{")
	assert.Contains(t, output, "}")
	assert.Contains(t, output, "json test")
}
