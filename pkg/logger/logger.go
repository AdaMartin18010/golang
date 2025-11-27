package logger

import (
	"context"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Logger 日志记录器
type Logger struct {
	*slog.Logger
	level       slog.Level
	writer      io.Writer
	sampleRate  float64
	serviceName string
	serviceVer  string
	mu          sync.RWMutex
}

// Config 日志配置
type Config struct {
	Level      slog.Level
	Output     io.Writer
	AddSource  bool
	JSONFormat bool
	// SampleRate 采样率 (0.0-1.0)，1.0 表示记录所有日志，0.5 表示记录 50%
	SampleRate float64
	// ServiceName 服务名称，会添加到所有日志中
	ServiceName string
	// ServiceVersion 服务版本，会添加到所有日志中
	ServiceVersion string
}

// NewLogger 创建日志记录器
func NewLogger(level slog.Level) *Logger {
	return NewLoggerWithConfig(Config{
		Level:      level,
		Output:     os.Stdout,
		AddSource:  false,
		JSONFormat: true,
		SampleRate: 1.0, // 默认记录所有日志
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

	// 如果配置了服务名称或版本，创建带默认字段的 logger
	var baseLogger *slog.Logger
	if config.ServiceName != "" || config.ServiceVersion != "" {
		attrs := []slog.Attr{}
		if config.ServiceName != "" {
			attrs = append(attrs, slog.String("service.name", config.ServiceName))
		}
		if config.ServiceVersion != "" {
			attrs = append(attrs, slog.String("service.version", config.ServiceVersion))
		}
		baseLogger = slog.New(handler).With(attrs)
	} else {
		baseLogger = slog.New(handler)
	}

	sampleRate := config.SampleRate
	if sampleRate <= 0 {
		sampleRate = 1.0
	} else if sampleRate > 1.0 {
		sampleRate = 1.0
	}

	return &Logger{
		Logger:      baseLogger,
		level:       config.Level,
		writer:      config.Output,
		sampleRate:  sampleRate,
		serviceName: config.ServiceName,
		serviceVer:  config.ServiceVersion,
	}
}

// WithContext 添加上下文信息
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	// 从 context 中提取追踪信息
	attrs := []slog.Attr{}

	// 优先从 OpenTelemetry trace context 中提取
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		spanCtx := span.SpanContext()
		attrs = append(attrs, slog.String("trace_id", spanCtx.TraceID().String()))
		attrs = append(attrs, slog.String("span_id", spanCtx.SpanID().String()))
		if spanCtx.TraceFlags().IsSampled() {
			attrs = append(attrs, slog.Bool("sampled", true))
		}
	} else {
		// 回退到从 context value 中提取（兼容旧代码）
		if traceID := getTraceID(ctx); traceID != "" {
			attrs = append(attrs, slog.String("trace_id", traceID))
		}
		if spanID := getSpanID(ctx); spanID != "" {
			attrs = append(attrs, slog.String("span_id", spanID))
		}
	}

	// 提取 UserID
	if userID := getUserID(ctx); userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}

	// 提取 RequestID
	if requestID := getRequestID(ctx); requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
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

// shouldSample 判断是否应该采样
func (l *Logger) shouldSample() bool {
	if l.sampleRate >= 1.0 {
		return true
	}
	if l.sampleRate <= 0 {
		return false
	}
	return rand.Float64() < l.sampleRate
}

// Debug 记录调试日志
func (l *Logger) Debug(msg string, args ...any) {
	if !l.shouldSample() {
		return
	}
	l.Logger.Debug(msg, args...)
}

// Info 记录信息日志
func (l *Logger) Info(msg string, args ...any) {
	if !l.shouldSample() {
		return
	}
	l.Logger.Info(msg, args...)
}

// Warn 记录警告日志
func (l *Logger) Warn(msg string, args ...any) {
	if !l.shouldSample() {
		return
	}
	l.Logger.Warn(msg, args...)
}

// Error 记录错误日志（错误日志总是记录，不受采样影响）
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// SetLevel 动态设置日志级别
func (l *Logger) SetLevel(level slog.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
	// 注意：slog.Logger 的级别是 Handler 级别的，这里只是记录
	// 实际使用时需要重新创建 Handler
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() slog.Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
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

// getRequestID 从 context 中获取 RequestID
func getRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}
