package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

// Logger 日志记录器
type Logger struct {
	*slog.Logger
	level  slog.Level
	writer io.Writer
}

// Config 日志配置
type Config struct {
	Level      slog.Level
	Output     io.Writer
	AddSource  bool
	JSONFormat bool
}

// NewLogger 创建日志记录器
func NewLogger(level slog.Level) *Logger {
	return NewLoggerWithConfig(Config{
		Level:      level,
		Output:     os.Stdout,
		AddSource:  false,
		JSONFormat: true,
	})
}

// NewLoggerWithConfig 使用配置创建日志记录器
func NewLoggerWithConfig(config Config) *Logger {
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	var handler slog.Handler
	if config.JSONFormat {
		handler = slog.NewJSONHandler(config.Output, opts)
	} else {
		handler = slog.NewTextHandler(config.Output, opts)
	}

	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
		level:  config.Level,
		writer: config.Output,
	}
}

// WithContext 添加上下文信息
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	// 从 context 中提取追踪信息
	attrs := []slog.Attr{}

	// 提取 TraceID
	if traceID := getTraceID(ctx); traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}

	// 提取 SpanID
	if spanID := getSpanID(ctx); spanID != "" {
		attrs = append(attrs, slog.String("span_id", spanID))
	}

	// 提取 UserID
	if userID := getUserID(ctx); userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}

	if len(attrs) > 0 {
		return l.Logger.With(attrs)
	}

	return l.Logger
}

// WithFields 添加字段
func (l *Logger) WithFields(fields ...any) *slog.Logger {
	return l.Logger.With(fields...)
}

// WithError 添加错误字段
func (l *Logger) WithError(err error) *slog.Logger {
	if err == nil {
		return l.Logger
	}
	return l.Logger.With("error", err.Error())
}

// Debug 记录调试日志
func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Info 记录信息日志
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Warn 记录警告日志
func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

// Error 记录错误日志
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// LogRequest 记录请求日志
func (l *Logger) LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	logger := l.WithContext(ctx)
	logger.Info("HTTP request",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration", duration,
	)
}

// LogError 记录错误日志
func (l *Logger) LogError(ctx context.Context, err error, msg string, args ...any) {
	allArgs := []any{"error", err.Error()}
	allArgs = append(allArgs, args...)

	logger := l.WithContext(ctx)
	logger.Error(msg, allArgs...)
}

// getTraceID 从 context 中获取 TraceID
func getTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// getSpanID 从 context 中获取 SpanID
func getSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value("span_id").(string); ok {
		return spanID
	}
	return ""
}

// getUserID 从 context 中获取 UserID
func getUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}
