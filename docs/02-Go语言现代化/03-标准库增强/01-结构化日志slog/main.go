package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

// DemonstrateSlog 包含了对 slog 包不同特性的演示。
func DemonstrateSlog() {
	BasicUsage()
	LogLevels()
	ContextualLogging()
	CustomizingHandler()
	GroupingAttributes()
}

// --- 1. 基础用法：TextHandler vs JSONHandler ---
func BasicUsage() {
	slog.Info("--- 1. Basic Usage ---")

	// 使用默认的 TextHandler
	slog.Info("Default logger (TextHandler)", "user", "Alice", "age", 30)

	// 创建并设置为默认的 JSONHandler
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(jsonHandler))
	slog.Info("Switched to JSONHandler", "user", "Bob", "permissions", []string{"read", "write"})

	// 切换回默认的 TextHandler 以便后续演示
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	slog.Info("--- End of Basic Usage ---\n")
}

// --- 2. 日志级别控制 ---
func LogLevels() {
	slog.Info("--- 2. Log Levels ---")

	// 创建一个只显示 Info 级别及以上日志的 Handler
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	logger.Debug("This debug message will be ignored.")
	logger.Info("This is an info message.")
	logger.Warn("This is a warning message.")
	logger.Error("This is an error message.")
	slog.Info("--- End of Log Levels ---\n")
}

// --- 3. 上下文日志：使用 With() ---
func ContextualLogging() {
	slog.Info("--- 3. Contextual Logging with With() ---")

	// 模拟一个 HTTP 请求的处理过程
	requestID := "req-xyz-789"
	logger := slog.With("request_id", requestID, "component", "APIHandler")

	logger.Info("Request processing started.")

	// 在处理过程中添加更多上下文
	userLogger := logger.With("user_id", "user-456")
	userLogger.Info("User authenticated.")
	userLogger.Warn("User has exceeded their API quota.")

	logger.Info("Request processing finished.")
	slog.Info("--- End of Contextual Logging ---\n")
}

// --- 4. 自定义 Handler 选项 ---
func CustomizingHandler() {
	slog.Info("--- 4. Customizing Handler Options ---")

	opts := &slog.HandlerOptions{
		// AddSource: true, // 包含源文件名和行号
		Level: slog.LevelDebug,
		// 自定义时间戳格式
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format(time.RFC3339Nano))
			}
			return a
		},
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsonHandler)

	logger.Debug("A debug message with custom timestamp format.")
	slog.Info("--- End of Customizing Handler ---\n")
}

// --- 5. 分组属性：使用 slog.Group ---
func GroupingAttributes() {
	slog.Info("--- 5. Grouping Attributes ---")

	// 使用默认的 TextHandler，以便清晰地看到分组效果
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	reqGroup := slog.Group("request",
		slog.String("method", "GET"),
		slog.String("url", "/api/v1/users"),
		slog.Int("status", 200),
	)

	resGroup := slog.Group("response",
		slog.Int("size", 1024),
		slog.Duration("duration", 50*time.Millisecond),
	)

	slog.Info("HTTP request processed", reqGroup, resGroup)
	slog.Info("--- End of Grouping Attributes ---\n")
}

// WithContextLogger 演示了如何在 context.Context 中传递 logger
func WithContextLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func main() {
	DemonstrateSlog()
}
