package otlp

import (
	"context"
	"log/slog"
)

// LoggerIntegration 日志集成工具
// 提供将标准库日志与 OTLP 集成的辅助函数
type LoggerIntegration struct {
	exporter LogExporter
	enabled  bool
}

// NewLoggerIntegration 创建日志集成工具
func NewLoggerIntegration(exporter LogExporter) *LoggerIntegration {
	return &LoggerIntegration{
		exporter: exporter,
		enabled:  exporter != nil,
	}
}

// ExportLog 导出日志到 OTLP
// 将 slog.Record 转换为 LogRecord 并导出
func (li *LoggerIntegration) ExportLog(ctx context.Context, record slog.Record, traceID, spanID string) error {
	if !li.enabled || li.exporter == nil {
		return nil
	}

	// 转换日志记录
	logRecord := LogRecord{
		Timestamp:  record.Time.UnixNano(),
		Severity:   record.Level.String(),
		Body:       record.Message,
		Attributes: make(map[string]interface{}),
		TraceID:    traceID,
		SpanID:     spanID,
		Resource:   make(map[string]interface{}),
	}

	// 提取属性
	record.Attrs(func(a slog.Attr) bool {
		logRecord.Attributes[a.Key] = a.Value.Any()
		return true
	})

	// 导出日志
	return li.exporter.Export(ctx, []LogRecord{logRecord})
}

// IsEnabled 检查是否启用
func (li *LoggerIntegration) IsEnabled() bool {
	return li.enabled
}

// Enable 启用集成
func (li *LoggerIntegration) Enable() {
	li.enabled = true
}

// Disable 禁用集成
func (li *LoggerIntegration) Disable() {
	li.enabled = false
}

// CreateSlogHandler 创建集成 OTLP 的 slog Handler
// 返回一个包装的 Handler，自动将日志导出到 OTLP
func CreateSlogHandler(baseHandler slog.Handler, integration *LoggerIntegration) slog.Handler {
	return &otlpHandler{
		Handler:     baseHandler,
		integration: integration,
	}
}

// otlpHandler slog Handler 包装器
type otlpHandler struct {
	slog.Handler
	integration *LoggerIntegration
}

func (h *otlpHandler) Handle(ctx context.Context, record slog.Record) error {
	// 先调用基础 Handler
	if err := h.Handler.Handle(ctx, record); err != nil {
		return err
	}

	// 提取追踪信息
	traceID := extractTraceID(ctx)
	spanID := extractSpanID(ctx)

	// 导出到 OTLP
	if h.integration != nil && h.integration.IsEnabled() {
		_ = h.integration.ExportLog(ctx, record, traceID, spanID)
	}

	return nil
}

// extractTraceID 从 context 中提取 TraceID
func extractTraceID(ctx context.Context) string {
	// TODO: 从 OpenTelemetry context 中提取
	// 当前为占位实现
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// extractSpanID 从 context 中提取 SpanID
func extractSpanID(ctx context.Context) string {
	// TODO: 从 OpenTelemetry context 中提取
	// 当前为占位实现
	if spanID, ok := ctx.Value("span_id").(string); ok {
		return spanID
	}
	return ""
}
