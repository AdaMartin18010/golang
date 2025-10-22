package observability

import (
	"context"
	"fmt"
	"time"
)

// =============================================================================
// 示例用法
// =============================================================================

// ExampleTracing 追踪示例
func ExampleTracing() {
	// 创建追踪器
	recorder := NewInMemoryRecorder()
	sampler := NewProbabilitySampler(1.0) // 100%采样
	tracer := NewTracer("my-service", recorder, sampler)

	// 设置为全局追踪器
	SetGlobalTracer(tracer)

	// 开始一个span
	ctx := context.Background()
	span, ctx := StartSpan(ctx, "handle-request")
	defer span.Finish()

	// 添加标签
	span.SetTag("http.method", "GET")
	span.SetTag("http.url", "/api/users")

	// 模拟一些工作
	processRequest(ctx)

	// 记录日志
	span.LogFields(map[string]interface{}{
		"event": "request_processed",
		"count": 42,
	})

	fmt.Println("Tracing example completed")
}

func processRequest(ctx context.Context) {
	// 创建子span
	span, _ := StartSpan(ctx, "database-query")
	defer span.Finish()

	span.SetTag("db.type", "postgres")
	span.SetTag("db.statement", "SELECT * FROM users")

	// 模拟数据库查询
	time.Sleep(10 * time.Millisecond)
}

// ExampleMetrics 指标示例
func ExampleMetrics() {
	// 创建和注册指标
	requestCounter := RegisterCounter(
		"http_requests_total",
		"Total HTTP requests",
		map[string]string{"endpoint": "/api/users"},
	)

	requestDuration := RegisterHistogram(
		"http_request_duration_seconds",
		"HTTP request latency",
		nil,
		map[string]string{"endpoint": "/api/users"},
	)

	activeConnections := RegisterGauge(
		"active_connections",
		"Number of active connections",
		nil,
	)

	// 使用指标
	activeConnections.Inc()
	defer activeConnections.Dec()

	start := time.Now()

	// 模拟请求处理
	time.Sleep(50 * time.Millisecond)

	duration := time.Since(start).Seconds()
	requestDuration.Observe(duration)
	requestCounter.Inc()

	// 导出指标（Prometheus格式）
	metrics := ExportMetrics()
	fmt.Printf("Metrics exported: %d bytes\n", len(metrics))
}

// ExampleLogging 日志示例
func ExampleLogging() {
	// 创建日志记录器
	logger := NewLogger(InfoLevel, nil)

	// 添加指标钩子
	logger.AddHook(NewMetricsHook())

	// 基本日志
	logger.Info("Service started")

	// 带字段的日志
	logger.WithField("user_id", "123").
		WithField("action", "login").
		Info("User logged in")

	// 带多个字段
	logger.WithFields(map[string]interface{}{
		"endpoint": "/api/users",
		"method":   "GET",
		"status":   200,
		"duration": 45.5,
	}).Info("Request completed")

	// 格式化日志
	logger.Infof("Processed %d requests in %v", 100, time.Second*5)

	fmt.Println("Logging example completed")
}

// ExampleIntegration 集成示例
func ExampleIntegration() {
	// 设置追踪
	recorder := NewInMemoryRecorder()
	sampler := &AlwaysSampler{}
	tracer := NewTracer("my-service", recorder, sampler)
	SetGlobalTracer(tracer)

	// 设置日志
	logger := NewLogger(InfoLevel, nil)
	logger.AddHook(NewMetricsHook())
	SetDefaultLogger(logger)

	// 设置指标
	requestCounter := RegisterCounter(
		"requests_total",
		"Total requests",
		nil,
	)

	// 处理请求
	ctx := context.Background()
	handleRequest(ctx, requestCounter, logger)

	fmt.Println("Integration example completed")
}

func handleRequest(ctx context.Context, counter *Counter, logger *Logger) {
	// 开始追踪
	span, ctx := StartSpan(ctx, "handle-request")
	defer span.Finish()

	// 使用带追踪信息的日志
	logger.WithContext(ctx).Info("Request started")

	// 模拟处理
	time.Sleep(20 * time.Millisecond)

	// 更新指标
	counter.Inc()

	// 记录完成
	logger.WithContext(ctx).Info("Request completed")
	span.SetStatus(StatusOK, "Success")
}

// RunAllExamples 运行所有示例
func RunAllExamples() {
	fmt.Println("=== Tracing Example ===")
	ExampleTracing()

	fmt.Println("\n=== Metrics Example ===")
	ExampleMetrics()

	fmt.Println("\n=== Logging Example ===")
	ExampleLogging()

	fmt.Println("\n=== Integration Example ===")
	ExampleIntegration()
}
