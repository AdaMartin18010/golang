package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger 日志记录器
type Logger struct {
	*slog.Logger
}

// NewLogger 创建日志记录器
func NewLogger(level slog.Level) *Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

// WithContext 添加上下文信息
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	// 从 context 中提取追踪信息
	// TODO: 集成 OpenTelemetry 追踪上下文
	return l.Logger
}

// WithFields 添加字段
func (l *Logger) WithFields(fields ...any) *slog.Logger {
	return l.Logger.With(fields...)
}
