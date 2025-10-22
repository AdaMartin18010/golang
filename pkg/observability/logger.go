package observability

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// =============================================================================
// 结构化日志 - Structured Logging
// =============================================================================

// LogLevel 日志级别
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry 日志条目
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
}

// Logger 日志记录器
type Logger struct {
	level      LogLevel
	output     io.Writer
	hooks      []LogHook
	fields     map[string]interface{}
	mu         sync.RWMutex
	slogLogger *slog.Logger
}

// LogHook 日志钩子
type LogHook interface {
	Fire(*LogEntry) error
}

// NewLogger 创建日志记录器
func NewLogger(level LogLevel, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	// 创建slog.Logger
	var slogLevel slog.Level
	switch level {
	case DebugLevel:
		slogLevel = slog.LevelDebug
	case InfoLevel:
		slogLevel = slog.LevelInfo
	case WarnLevel:
		slogLevel = slog.LevelWarn
	case ErrorLevel:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slogLevel,
	})

	return &Logger{
		level:      level,
		output:     output,
		hooks:      make([]LogHook, 0),
		fields:     make(map[string]interface{}),
		slogLogger: slog.New(handler),
	}
}

// WithField 添加字段
func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:      l.level,
		output:     l.output,
		hooks:      l.hooks,
		fields:     make(map[string]interface{}),
		slogLogger: l.slogLogger,
	}

	// 复制现有字段
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// 添加新字段
	newLogger.fields[key] = value

	return newLogger
}

// WithFields 添加多个字段
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &Logger{
		level:      l.level,
		output:     l.output,
		hooks:      l.hooks,
		fields:     make(map[string]interface{}),
		slogLogger: l.slogLogger,
	}

	// 复制现有字段
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// 添加新字段
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// WithContext 从context添加追踪信息
func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := SpanFromContext(ctx)
	if span == nil {
		return l
	}

	return l.WithFields(map[string]interface{}{
		"trace_id": span.TraceID,
		"span_id":  span.SpanID,
	})
}

// AddHook 添加钩子
func (l *Logger) AddHook(hook LogHook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, hook)
}

// log 记录日志
func (l *Logger) log(level LogLevel, message string) {
	if level < l.level {
		return
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    make(map[string]interface{}),
	}

	// 复制字段
	for k, v := range l.fields {
		entry.Fields[k] = v
		// 提取trace_id和span_id
		if k == "trace_id" {
			if traceID, ok := v.(string); ok {
				entry.TraceID = traceID
			}
		}
		if k == "span_id" {
			if spanID, ok := v.(string); ok {
				entry.SpanID = spanID
			}
		}
	}

	// 调用hooks
	for _, hook := range l.hooks {
		if err := hook.Fire(entry); err != nil {
			fmt.Fprintf(os.Stderr, "Hook error: %v\n", err)
		}
	}

	// 使用slog输出
	attrs := make([]slog.Attr, 0, len(entry.Fields))
	for k, v := range entry.Fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	var slogLevel slog.Level
	switch level {
	case DebugLevel:
		slogLevel = slog.LevelDebug
	case InfoLevel:
		slogLevel = slog.LevelInfo
	case WarnLevel:
		slogLevel = slog.LevelWarn
	case ErrorLevel:
		slogLevel = slog.LevelError
	case FatalLevel:
		slogLevel = slog.LevelError + 1
	}

	l.slogLogger.LogAttrs(context.Background(), slogLevel, message, attrs...)
}

// Debug 记录调试日志
func (l *Logger) Debug(message string) {
	l.log(DebugLevel, message)
}

// Info 记录信息日志
func (l *Logger) Info(message string) {
	l.log(InfoLevel, message)
}

// Warn 记录警告日志
func (l *Logger) Warn(message string) {
	l.log(WarnLevel, message)
}

// Error 记录错误日志
func (l *Logger) Error(message string) {
	l.log(ErrorLevel, message)
}

// Fatal 记录致命日志并退出
func (l *Logger) Fatal(message string) {
	l.log(FatalLevel, message)
	os.Exit(1)
}

// Debugf 格式化调试日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Infof 格式化信息日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warnf 格式化警告日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Fatalf 格式化致命日志
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// =============================================================================
// 预设钩子
// =============================================================================

// MetricsHook 指标钩子（记录日志数量）
type MetricsHook struct {
	counters map[LogLevel]*Counter
	mu       sync.RWMutex
}

// NewMetricsHook 创建指标钩子
func NewMetricsHook() *MetricsHook {
	hook := &MetricsHook{
		counters: make(map[LogLevel]*Counter),
	}

	// 为每个级别创建计数器
	for _, level := range []LogLevel{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel} {
		counter := RegisterCounter(
			fmt.Sprintf("log_entries_total_%s", level.String()),
			fmt.Sprintf("Total number of %s log entries", level.String()),
			map[string]string{"level": level.String()},
		)
		hook.counters[level] = counter
	}

	return hook
}

func (h *MetricsHook) Fire(entry *LogEntry) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if counter, ok := h.counters[entry.Level]; ok {
		counter.Inc()
	}

	return nil
}

// FileHook 文件钩子（写入文件）
type FileHook struct {
	file     *os.File
	mu       sync.Mutex
	minLevel LogLevel
}

// NewFileHook 创建文件钩子
func NewFileHook(filename string, minLevel LogLevel) (*FileHook, error) {
	// 使用更安全的文件权限 0600 (只有所有者可读写)
	// #nosec G304 - filename由调用者控制，应在调用处验证
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	return &FileHook{
		file:     file,
		minLevel: minLevel,
	}, nil
}

func (h *FileHook) Fire(entry *LogEntry) error {
	if entry.Level < h.minLevel {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	line := fmt.Sprintf("[%s] %s: %s",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level.String(),
		entry.Message,
	)

	if len(entry.Fields) > 0 {
		line += fmt.Sprintf(" %v", entry.Fields)
	}

	line += "\n"

	_, err := h.file.WriteString(line)
	return err
}

func (h *FileHook) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.file.Close()
}

// =============================================================================
// 全局日志记录器
// =============================================================================

var (
	defaultLogger *Logger
	loggerMu      sync.RWMutex
)

func init() {
	// 初始化默认日志记录器
	defaultLogger = NewLogger(InfoLevel, os.Stdout)

	// 添加指标钩子
	defaultLogger.AddHook(NewMetricsHook())
}

// SetDefaultLogger 设置默认日志记录器
func SetDefaultLogger(logger *Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	defaultLogger = logger
}

// GetDefaultLogger 获取默认日志记录器
func GetDefaultLogger() *Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	return defaultLogger
}

// Debug 使用默认日志记录器记录调试日志
func Debug(message string) {
	GetDefaultLogger().Debug(message)
}

// Info 使用默认日志记录器记录信息日志
func Info(message string) {
	GetDefaultLogger().Info(message)
}

// Warn 使用默认日志记录器记录警告日志
func Warn(message string) {
	GetDefaultLogger().Warn(message)
}

// Error 使用默认日志记录器记录错误日志
func Error(message string) {
	GetDefaultLogger().Error(message)
}

// Fatal 使用默认日志记录器记录致命日志
func Fatal(message string) {
	GetDefaultLogger().Fatal(message)
}

// WithField 使用默认日志记录器添加字段
func WithField(key string, value interface{}) *Logger {
	return GetDefaultLogger().WithField(key, value)
}

// WithFields 使用默认日志记录器添加多个字段
func WithFields(fields map[string]interface{}) *Logger {
	return GetDefaultLogger().WithFields(fields)
}

// WithContext 使用默认日志记录器从context添加追踪信息
func WithContext(ctx context.Context) *Logger {
	return GetDefaultLogger().WithContext(ctx)
}
