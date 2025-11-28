package logger

import (
	"fmt"
	"io"
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

// RotationConfig 日志轮转配置
type RotationConfig struct {
	// Filename 日志文件路径
	Filename string
	// MaxSize 单个日志文件的最大大小（MB），超过此大小会轮转
	MaxSize int
	// MaxBackups 保留的旧日志文件数量
	MaxBackups int
	// MaxAge 保留旧日志文件的天数
	MaxAge int
	// Compress 是否压缩轮转后的旧日志文件
	Compress bool
	// LocalTime 是否使用本地时间（而非 UTC）
	LocalTime bool
}

// NewRotatingWriter 创建支持轮转的日志写入器
// 基于 lumberjack 实现，支持：
// - 按大小轮转（MaxSize）
// - 按时间轮转（MaxAge）
// - 自动压缩旧日志（Compress）
// - 限制保留文件数量（MaxBackups）
//
// 如果配置无效，会返回错误
func NewRotatingWriter(cfg RotationConfig) (io.Writer, error) {
	if err := ValidateRotationConfig(cfg); err != nil {
		return nil, err
	}
	return &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  cfg.LocalTime,
	}, nil
}

// NewRotatingWriterOrPanic 创建支持轮转的日志写入器，如果失败则 panic
func NewRotatingWriterOrPanic(cfg RotationConfig) io.Writer {
	writer, err := NewRotatingWriter(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create rotating writer: %v", err))
	}
	return writer
}

// NewRotatingLogger 创建支持轮转的日志记录器
func NewRotatingLogger(level slog.Level, cfg RotationConfig) (*Logger, error) {
	writer, err := NewRotatingWriter(cfg)
	if err != nil {
		return nil, err
	}
	return NewLoggerWithConfig(Config{
		Level:      level,
		Output:     writer,
		AddSource:  false,
		JSONFormat: true,
		SampleRate: 1.0,
	}), nil
}

// NewRotatingLoggerOrPanic 创建支持轮转的日志记录器，如果失败则 panic
func NewRotatingLoggerOrPanic(level slog.Level, cfg RotationConfig) *Logger {
	logger, err := NewRotatingLogger(level, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create rotating logger: %v", err))
	}
	return logger
}

// DefaultRotationConfig 返回默认的轮转配置
func DefaultRotationConfig(filename string) RotationConfig {
	return RotationConfig{
		Filename:   filename,
		MaxSize:    100,  // 100MB
		MaxBackups: 10,   // 保留10个备份
		MaxAge:     30,   // 保留30天
		Compress:   true, // 压缩旧日志
		LocalTime:  true, // 使用本地时间
	}
}

// ProductionRotationConfig 返回生产环境的轮转配置
// 更严格的配置，适合生产环境使用
func ProductionRotationConfig(filename string) RotationConfig {
	return RotationConfig{
		Filename:   filename,
		MaxSize:    500,  // 500MB
		MaxBackups: 20,   // 保留20个备份
		MaxAge:     90,   // 保留90天
		Compress:   true, // 压缩旧日志
		LocalTime:  true, // 使用本地时间
	}
}

// DevelopmentRotationConfig 返回开发环境的轮转配置
// 更宽松的配置，适合开发环境使用
func DevelopmentRotationConfig(filename string) RotationConfig {
	return RotationConfig{
		Filename:   filename,
		MaxSize:    50,   // 50MB
		MaxBackups: 5,    // 保留5个备份
		MaxAge:     7,    // 保留7天
		Compress:   false, // 开发环境不压缩
		LocalTime:  true,  // 使用本地时间
	}
}

// ValidateRotationConfig 验证轮转配置的有效性
func ValidateRotationConfig(cfg RotationConfig) error {
	if cfg.Filename == "" {
		return fmt.Errorf("filename is required")
	}
	if cfg.MaxSize <= 0 {
		return fmt.Errorf("max_size must be greater than 0")
	}
	if cfg.MaxBackups < 0 {
		return fmt.Errorf("max_backups must be non-negative")
	}
	if cfg.MaxAge < 0 {
		return fmt.Errorf("max_age must be non-negative")
	}
	return nil
}
