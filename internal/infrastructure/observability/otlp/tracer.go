package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

// NewTracerProvider 创建并配置 OpenTelemetry 追踪提供者（TracerProvider）。
//
// 功能说明：
// - 创建 OTLP gRPC 导出器，用于将追踪数据发送到 OpenTelemetry Collector
// - 配置服务资源信息（服务名、版本等）
// - 设置全局 TracerProvider 和文本映射传播器
// - 返回关闭函数，用于优雅关闭追踪提供者
//
// 设计原则：
// 1. 标准化：使用 OTLP（OpenTelemetry Protocol）标准协议
// 2. 批处理：使用 Batcher 批量发送追踪数据，提高性能
// 3. 上下文传播：支持 TraceContext 和 Baggage 传播，实现跨服务追踪
// 4. 资源标识：通过资源属性标识服务，便于在追踪系统中识别
//
// 参数：
// - ctx: 上下文，用于创建导出器和资源
// - endpoint: OpenTelemetry Collector 的 gRPC 端点地址
//   示例：localhost:4317（gRPC）或 localhost:4318（HTTP）
// - insecure: 是否使用不安全的连接（不使用 TLS）
//   - true: 使用 grpc.WithInsecure()，适用于开发环境
//   - false: 使用 TLS 加密，适用于生产环境
//
// 返回：
// - func(context.Context) error: 关闭函数，用于优雅关闭 TracerProvider
//   应在应用程序退出时调用，确保所有追踪数据都已发送
// - error: 如果创建失败，返回错误信息
//
// 使用示例：
//
//	ctx := context.Background()
//	shutdown, err := otlp.NewTracerProvider(ctx, "localhost:4317", true)
//	if err != nil {
//	    log.Fatal("Failed to create tracer provider:", err)
//	}
//	defer func() {
//	    if err := shutdown(ctx); err != nil {
//	        log.Printf("Error shutting down tracer provider: %v", err)
//	    }
//	}()
//
//	// 使用全局 TracerProvider 创建 Span
//	tracer := otel.Tracer("my-service")
//	ctx, span := tracer.Start(ctx, "operation-name")
//	defer span.End()
//
// 配置说明：
// - ServiceName: 服务名称，用于在追踪系统中标识服务
// - ServiceVersion: 服务版本，用于区分不同版本的服务实例
// - 其他资源属性可通过 resource.WithAttributes() 添加
//
// 追踪数据流程：
// 1. 应用程序创建 Span
// 2. Span 数据被收集到 TracerProvider
// 3. Batcher 批量处理 Span 数据
// 4. OTLP Exporter 通过 gRPC 发送到 Collector
// 5. Collector 将数据转发到后端系统（Jaeger、Zipkin、Tempo 等）
//
// 注意事项：
// - 确保 OpenTelemetry Collector 已启动并监听指定端点
// - 生产环境应使用 TLS 加密连接（insecure=false）
// - 关闭函数应在应用程序退出时调用，确保数据不丢失
// - 如果 endpoint 为空，可以考虑使用 NoOp TracerProvider
//
// 相关组件：
// - OpenTelemetry Collector: 接收、处理和导出追踪数据
// - Jaeger/Zipkin/Tempo: 追踪数据存储和可视化后端
// - TraceContext Propagation: W3C Trace Context 标准，用于跨服务追踪
// - Baggage: 在追踪上下文中传递额外的键值对数据
func NewTracerProvider(ctx context.Context, endpoint string, insecure bool) (func(context.Context) error, error) {
	// 创建资源（Resource）
	// 资源用于标识产生追踪数据的服务或实体
	// 资源属性包括服务名、版本、部署环境等信息
	// 这些信息会在所有 Span 中自动包含，便于在追踪系统中识别和过滤
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("golang-service"),      // 服务名称
			semconv.ServiceVersion("1.0.0"),            // 服务版本
			// 可以添加更多资源属性，例如：
			// semconv.DeploymentEnvironment("production"),
			// semconv.HostName("server-01"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建 OTLP gRPC 导出器
	// OTLP（OpenTelemetry Protocol）是 OpenTelemetry 的标准协议
	// 支持 gRPC 和 HTTP 两种传输方式
	// 导出器负责将追踪数据发送到 OpenTelemetry Collector
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint), // Collector 端点地址
		// 其他可选配置：
		// otlptracegrpc.WithTimeout(30*time.Second),  // 超时时间
		// otlptracegrpc.WithHeaders(map[string]string{"api-key": "xxx"}), // 自定义头部
	}
	if insecure {
		// 开发环境：使用不安全的连接（不使用 TLS）
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	// 生产环境应使用 TLS：
	// opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials))

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// 创建追踪提供者（TracerProvider）
	// TracerProvider 是追踪系统的核心组件，负责：
	// 1. 管理 Tracer 实例
	// 2. 配置采样策略
	// 3. 处理 Span 的创建和完成
	// 4. 将 Span 数据发送到导出器
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // 使用批处理器，提高性能
		// 批处理器会收集多个 Span，然后批量发送，减少网络开销
		// 可以配置批处理参数：
		// sdktrace.WithBatchTimeout(5*time.Second),  // 批处理超时
		// sdktrace.WithMaxExportBatchSize(512),      // 最大批处理大小
		sdktrace.WithResource(res), // 关联资源信息
		// 其他可选配置：
		// sdktrace.WithSampler(sdktrace.AlwaysSample()), // 采样策略
		// sdktrace.WithIDGenerator(...),                 // ID 生成器
	)

	// 设置为全局 TracerProvider
	// 这样应用程序的其他部分可以通过 otel.Tracer() 获取 Tracer
	otel.SetTracerProvider(tp)

	// 设置文本映射传播器（Text Map Propagator）
	// 传播器负责在跨服务调用时传递追踪上下文
	// TraceContext: W3C Trace Context 标准，传递 TraceID 和 SpanID
	// Baggage: 在追踪上下文中传递额外的键值对数据
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // W3C Trace Context 传播
		propagation.Baggage{},      // Baggage 传播
	))

	// 返回关闭函数
	// 关闭函数会：
	// 1. 停止接受新的 Span
	// 2. 等待所有正在处理的 Span 完成
	// 3. 刷新所有待发送的追踪数据
	// 4. 关闭导出器连接
	return tp.Shutdown, nil
}
