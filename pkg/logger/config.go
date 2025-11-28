package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// ConfigFromFile 从配置文件创建日志配置
// 将 internal/config 的 LoggingConfig 转换为 logger.Config
func ConfigFromFile(cfg interface{}) (Config, error) {
	// 这里需要根据实际的配置结构进行转换
	// 当前为占位实现
	return Config{
		Level:      slog.LevelInfo,
		Output:     os.Stdout,
		AddSource:  false,
		JSONFormat: true,
		SampleRate: 1.0,
	}, nil
}

// CreateLoggerFromConfig 从配置创建日志记录器
// 支持从配置文件或环境变量创建
func CreateLoggerFromConfig(
	level string,
	format string,
	output string,
	outputPath string,
	rotationCfg RotationConfig,
) (*Logger, error) {
	// 解析日志级别
	logLevel, err := parseLogLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// 确定输出目标
	var writer io.Writer
	if output == "file" {
		if outputPath == "" {
			return nil, fmt.Errorf("output_path is required when output is file")
		}

		// 如果配置了轮转，使用轮转写入器
		if rotationCfg.Filename != "" || outputPath != "" {
			if rotationCfg.Filename == "" {
				rotationCfg.Filename = outputPath
			}
			// 如果配置为空，使用默认值
			if rotationCfg.MaxSize == 0 {
				rotationCfg = DefaultRotationConfig(outputPath)
			}

			var err error
			writer, err = NewRotatingWriter(rotationCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to create rotating writer: %w", err)
			}
		} else {
			// 不使用轮转，直接创建文件
			file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to open log file: %w", err)
			}
			writer = file
		}
	} else {
		// 输出到标准输出
		writer = os.Stdout
	}

	// 确定格式
	jsonFormat := strings.ToLower(format) == "json"

	// 创建日志记录器
	return NewLoggerWithConfig(Config{
		Level:      logLevel,
		Output:     writer,
		AddSource:  false,
		JSONFormat: jsonFormat,
		SampleRate: 1.0,
	}), nil
}

// parseLogLevel 解析日志级别字符串
func parseLogLevel(level string) (slog.Level, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN", "WARNING":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}

// CreateLoggerFromFileConfig 从文件配置创建日志记录器
// 这是一个便捷函数，用于从 internal/config.LoggingConfig 创建日志记录器
func CreateLoggerFromFileConfig(
	level string,
	format string,
	output string,
	outputPath string,
	rotationCfg interface{}, // internal/config.RotationConfig
) (*Logger, error) {
	// 转换轮转配置
	var loggerRotationCfg RotationConfig
	if rotationCfg != nil {
		// 这里需要根据实际的配置结构进行转换
		// 当前为占位实现
		loggerRotationCfg = DefaultRotationConfig(outputPath)
	}

	return CreateLoggerFromConfig(level, format, output, outputPath, loggerRotationCfg)
}
