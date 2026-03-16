// Package logger 提供框架级别的统一日志系统
//
// 设计原理：
// 1. 统一的日志接口，框架内所有组件使用相同的日志系统
// 2. 基于 slog（Go 1.21+ 标准库），提供结构化日志
// 3. 支持 OpenTelemetry 集成，自动添加追踪信息
// 4. 支持日志级别、采样率、输出格式等配置
// 5. 线程安全，支持并发写入
//
// 架构位置：
// - Framework Layer (internal/framework/logger/)
// - 被所有层使用（Domain, Application, Infrastructure, Interfaces）
//
// 使用场景：
// 1. 框架内部日志记录
// 2. 应用日志记录（通过依赖注入）
// 3. 调试和问题排查
// 4. 生产环境监控
//
// 示例：
//
//	// 创建全局日志实例
//	logger := framework.NewLogger(framework.LoggerConfig{
//	    Level:      slog.LevelInfo,
//	    JSONFormat: true,
//	    ServiceName: "my-service",
//	})
//
//	// 使用日志
//	logger.Info("Application started", "port", 8080)
//
//	// 带上下文的日志
//	ctx := context.WithValue(context.Background(), "request_id", "req-123")
//	logger.WithContext(ctx).Info("Processing request")
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/yourusername/golang/pkg/logger"
)

// Logger 框架日志接口
//
// 功能说明：
// - 提供统一的日志记录接口
// - 支持结构化日志（key-value pairs）
// - 支持日志级别控制
// - 支持上下文集成（追踪信息）
type Logger interface {
	// Debug 记录调试日志
	Debug(msg string, args ...any)
	// Info 记录信息日志
	Info(msg string, args ...any)
	// Warn 记录警告日志
	Warn(msg string, args ...any)
	// Error 记录错误日志
	Error(msg string, args ...any)
	// WithContext 添加上下文信息（追踪ID、请求ID等）
	WithContext(ctx context.Context) *slog.Logger
	// WithFields 添加字段
	WithFields(fields ...any) *slog.Logger
	// WithError 添加错误字段
	WithError(err error) *slog.Logger
}

// Config 日志配置
type Config struct {
	// Level 日志级别（Debug/Info/Warn/Error）
	Level slog.Level
	// Output 日志输出（默认 os.Stdout）
	Output io.Writer
	// AddSource 是否添加源代码位置
	AddSource bool
	// JSONFormat 是否使用 JSON 格式（默认 true）
	JSONFormat bool
	// SampleRate 采样率 (0.0-1.0)，1.0 表示记录所有日志
	SampleRate float64
	// ServiceName 服务名称，会添加到所有日志中
	ServiceName string
	// ServiceVersion 服务版本，会添加到所有日志中
	ServiceVersion string
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Level:      slog.LevelInfo,
		Output:     os.Stdout,
		AddSource:  false,
		JSONFormat: true,
		SampleRate: 1.0,
	}
}

// frameworkLogger 框架日志实现
type frameworkLogger struct {
	*logger.Logger
	mu sync.RWMutex
}

var (
	// defaultLogger 默认日志实例（全局单例）
	defaultLogger Logger
	once          sync.Once
)

// NewLogger 创建新的日志实例
//
// 功能说明：
// - 根据配置创建日志实例
// - 如果未指定配置，使用默认配置
//
// 参数：
// - config: 日志配置（可选，nil 时使用默认配置）
//
// 返回：
// - Logger: 日志实例
func NewLogger(config *Config) Logger {
	cfg := DefaultConfig()
	if config != nil {
		if config.Level != 0 {
			cfg.Level = config.Level
		}
		if config.Output != nil {
			cfg.Output = config.Output
		}
		cfg.AddSource = config.AddSource
		if config.JSONFormat {
			cfg.JSONFormat = true
		}
		if config.SampleRate > 0 {
			cfg.SampleRate = config.SampleRate
		}
		cfg.ServiceName = config.ServiceName
		cfg.ServiceVersion = config.ServiceVersion
	}

	baseLogger := logger.NewLoggerWithConfig(logger.Config{
		Level:          cfg.Level,
		Output:         cfg.Output,
		AddSource:      cfg.AddSource,
		JSONFormat:     cfg.JSONFormat,
		SampleRate:     cfg.SampleRate,
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
	})

	return &frameworkLogger{
		Logger: baseLogger,
	}
}

// GetLogger 获取默认日志实例（单例模式）
//
// 功能说明：
// - 首次调用时创建默认日志实例
// - 后续调用返回同一个实例
//
// 返回：
// - Logger: 默认日志实例
func GetLogger() Logger {
	once.Do(func() {
		defaultLogger = NewLogger(nil)
	})
	return defaultLogger
}

// SetLogger 设置默认日志实例
//
// 功能说明：
// - 允许替换默认日志实例
// - 用于测试或自定义日志配置
//
// 参数：
// - l: 新的日志实例
func SetLogger(l Logger) {
	defaultLogger = l
}

// Debug 记录调试日志
func (l *frameworkLogger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Info 记录信息日志
func (l *frameworkLogger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Warn 记录警告日志
func (l *frameworkLogger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

// Error 记录错误日志
func (l *frameworkLogger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// WithContext 添加上下文信息
func (l *frameworkLogger) WithContext(ctx context.Context) *slog.Logger {
	return l.Logger.WithContext(ctx)
}

// WithFields 添加字段
func (l *frameworkLogger) WithFields(fields ...any) *slog.Logger {
	return l.Logger.WithFields(fields...)
}

// WithError 添加错误字段
func (l *frameworkLogger) WithError(err error) *slog.Logger {
	return l.Logger.WithError(err)
}
