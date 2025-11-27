// Package otlp provides OpenTelemetry integration for observability,
// including logging, tracing, and metrics collection.
//
// 设计原则：
// 1. 标准化：遵循 OpenTelemetry 标准，确保与各种后端系统兼容
// 2. 可观测性：提供日志、追踪、指标三大支柱的完整支持
// 3. 性能：最小化对应用性能的影响，使用异步和批处理机制
// 4. 可配置：支持灵活的配置选项，适应不同环境需求
//
// 核心功能：
// - Logger: 结构化日志记录，支持 JSON 格式和上下文集成
// - Tracer: 分布式追踪，支持跨服务调用链追踪
// - Metrics: 指标收集，支持 Counter、Gauge、Histogram 等类型
//
// 使用场景：
// - 微服务架构中的分布式追踪
// - 应用性能监控（APM）
// - 错误追踪和调试
// - 业务指标监控
package otlp

import (
	"context"
	"log/slog"
	"os"
)

// Logger 是 OpenTelemetry 集成的结构化日志记录器。
//
// 设计说明：
// - 基于 Go 标准库的 slog 包实现
// - 支持 JSON 格式输出，便于日志聚合和分析
// - 可集成 OpenTelemetry 追踪上下文（TraceID、SpanID）
// - 支持多级日志（Debug、Info、Warn、Error、Fatal）
//
// 使用示例：
//
//	logger := otlp.NewLogger()
//	logger.Info("Application started", "port", 8080)
//
//	// 带上下文的日志（集成追踪信息）
//	ctx := context.Background()
//	logger.WithContext(ctx).Info("Processing request", "user_id", 123)
//
// 注意事项：
// - 日志级别可通过 HandlerOptions 配置
// - 生产环境建议使用 JSON 格式，便于日志系统解析
// - 敏感信息不应记录在日志中
type Logger struct {
	*slog.Logger
}

// NewLogger 创建一个新的 OpenTelemetry 日志记录器。
//
// 功能说明：
// - 使用 JSON 格式输出到标准输出
// - 默认日志级别为 Info
// - 支持结构化字段（key-value pairs）
//
// 配置选项：
// - Level: 日志级别（Debug/Info/Warn/Error/Fatal）
// - AddSource: 是否添加源代码位置信息
// - ReplaceAttr: 自定义属性替换函数
//
// 返回：
// - *Logger: 配置好的日志记录器实例
//
// 示例：
//
//	logger := otlp.NewLogger()
//	logger.Info("Server starting", "port", 8080, "env", "production")
func NewLogger() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

// WithContext 从上下文中提取追踪信息并返回带有追踪上下文的日志记录器。
//
// 功能说明：
// - 从 context 中提取 OpenTelemetry 追踪信息（TraceID、SpanID）
// - 将追踪信息作为日志字段自动添加
// - 实现日志与追踪的关联，便于问题排查
//
// 参数：
// - ctx: 包含追踪信息的上下文
//
// 返回：
// - *slog.Logger: 带有追踪上下文的日志记录器
//
// 使用场景：
// - HTTP 请求处理中记录请求日志
// - 工作流执行中记录步骤日志
// - 异步任务处理中记录任务日志
//
// 示例：
//
//	ctx := context.Background()
//	// 假设 ctx 中已包含追踪信息
//	logger.WithContext(ctx).Info("Processing request", "method", "GET", "path", "/api/users")
//
// 注意事项：
// - 当前实现为占位符，需要集成 OpenTelemetry 的 context 提取功能
// - 追踪信息会自动添加到日志的 "trace_id" 和 "span_id" 字段
// - 如果 context 中没有追踪信息，则返回原始日志记录器
//
// TODO: 集成 OpenTelemetry 追踪上下文提取功能
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	// 从 context 中提取追踪信息
	// TODO: 集成 OpenTelemetry 追踪上下文
	// 实现思路：
	// 1. 使用 otel.GetTextMapPropagator().Extract() 提取追踪信息
	// 2. 从 span 中获取 TraceID 和 SpanID
	// 3. 使用 logger.With() 添加追踪字段
	// 示例：
	//   span := trace.SpanFromContext(ctx)
	//   if span.SpanContext().IsValid() {
	//       traceID := span.SpanContext().TraceID().String()
	//       spanID := span.SpanContext().SpanID().String()
	//       return l.Logger.With("trace_id", traceID, "span_id", spanID)
	//   }
	return l.Logger
}
