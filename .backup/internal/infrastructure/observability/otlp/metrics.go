package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

// MetricsProvider 是 OpenTelemetry 指标提供者的封装。
//
// 功能说明：
// - 管理 MeterProvider 实例，用于创建和管理指标（Metrics）
// - 支持多种指标类型：Counter、Gauge、Histogram、Summary
// - 将指标数据导出到 OpenTelemetry Collector
//
// 设计原则：
// 1. 标准化：遵循 OpenTelemetry Metrics 标准
// 2. 性能：使用周期性读取器（PeriodicReader）批量导出指标
// 3. 资源标识：通过资源属性标识服务，便于指标聚合和过滤
//
// 使用场景：
// - 业务指标监控（请求数、响应时间、错误率等）
// - 系统指标监控（CPU、内存、网络等）
// - 自定义业务指标（订单数、用户数、交易金额等）
//
// 指标类型说明：
// - Counter: 累加型指标，只能增加（如：请求总数、错误总数）
// - Gauge: 瞬时值指标，可以增减（如：当前连接数、队列长度）
// - Histogram: 分布型指标，记录值的分布（如：响应时间分布）
// - Summary: 摘要型指标，提供分位数统计（如：P50、P95、P99 延迟）
//
// 示例：
//
//	// 创建指标提供者
//	mp, err := otlp.NewMetricsProvider(ctx, "localhost:4317", true)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer mp.Shutdown(ctx)
//
//	// 获取 Meter 并创建指标
//	meter := mp.MeterProvider().Meter("my-service")
//	counter, _ := meter.Int64Counter("requests_total", metric.WithDescription("Total requests"))
//	counter.Add(ctx, 1)
type MetricsProvider struct {
	provider *metric.MeterProvider
}

// NewMetricsProvider 创建并配置 OpenTelemetry 指标提供者。
//
// 功能说明：
// - 创建服务资源信息（服务名、版本等）
// - 配置指标导出器（当前为占位符实现）
// - 返回配置好的 MetricsProvider 实例
//
// 参数：
//   - ctx: 上下文，用于创建导出器和资源
//   - endpoint: OpenTelemetry Collector 的 gRPC 端点地址
//     示例：localhost:4317（gRPC）或 localhost:4318（HTTP）
//   - insecure: 是否使用不安全的连接（不使用 TLS）
//   - true: 使用不安全的连接，适用于开发环境
//   - false: 使用 TLS 加密，适用于生产环境
//
// 返回：
// - *MetricsProvider: 配置好的指标提供者实例
// - error: 如果创建失败，返回错误信息
//
// 当前状态：
// - 基础结构已实现
// - OTLP metrics 导出器待实现（TODO）
// - 当前使用空的 MeterProvider，指标数据不会导出
//
// 完整实现示例（待实现）：
//
//	import (
//	    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
//	)
//
//	// 创建 OTLP gRPC 导出器
//	opts := []otlpmetricgrpc.Option{
//	    otlpmetricgrpc.WithEndpoint(endpoint),
//	}
//	if insecure {
//	    opts = append(opts, otlpmetricgrpc.WithInsecure())
//	}
//	exporter, err := otlpmetricgrpc.New(ctx, opts...)
//	if err != nil {
//	    return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
//	}
//
//	// 创建周期性读取器
//	reader := metric.NewPeriodicReader(exporter,
//	    metric.WithInterval(30*time.Second), // 每 30 秒导出一次
//	)
//
//	// 创建指标提供者
//	mp := metric.NewMeterProvider(
//	    metric.WithResource(res),
//	    metric.WithReader(reader),
//	)
//
// 注意事项：
// - 确保 OpenTelemetry Collector 已启动并监听指定端点
// - 生产环境应使用 TLS 加密连接（insecure=false）
// - 指标导出是异步的，使用周期性读取器批量导出
// - 关闭函数应在应用程序退出时调用，确保数据不丢失
//
// TODO: 实现 OTLP metrics 导出器
func NewMetricsProvider(ctx context.Context, endpoint string, insecure bool) (*MetricsProvider, error) {
	// 创建资源（Resource）
	// 资源用于标识产生指标数据的服务或实体
	// 资源属性包括服务名、版本、部署环境等信息
	// 这些信息会在所有指标中自动包含，便于在监控系统中识别和过滤
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("golang-service"), // 服务名称
			semconv.ServiceVersion("1.0.0"),       // 服务版本
			// 可以添加更多资源属性，例如：
			// semconv.DeploymentEnvironment("production"),
			// semconv.HostName("server-01"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// TODO: 实现 OTLP metrics 导出器
	// 暂时使用空的 reader，后续添加 OTLP exporter
	// 实现步骤：
	// 1. 导入 otlpmetricgrpc 包
	// 2. 创建 OTLP gRPC 导出器（类似 tracer.go 中的实现）
	// 3. 创建周期性读取器（PeriodicReader），配置导出间隔
	// 4. 将读取器添加到 MeterProvider
	//
	// 示例代码：
	//   exporter, err := otlpmetricgrpc.New(ctx, opts...)
	//   reader := metric.NewPeriodicReader(exporter, metric.WithInterval(30*time.Second))
	//   mp := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(reader))

	// 创建指标提供者（MeterProvider）
	// MeterProvider 是指标系统的核心组件，负责：
	// 1. 管理 Meter 实例
	// 2. 处理指标的创建和更新
	// 3. 将指标数据发送到导出器
	mp := metric.NewMeterProvider(
		metric.WithResource(res), // 关联资源信息
		// metric.WithReader(metric.NewPeriodicReader(exporter)), // 待实现：添加导出器
		// 其他可选配置：
		// metric.WithView(metric.NewView(...)), // 指标视图配置
	)

	return &MetricsProvider{
		provider: mp,
	}, nil
}

// MeterProvider 返回底层的 MeterProvider 实例。
//
// 功能说明：
// - 提供对底层 OpenTelemetry MeterProvider 的访问
// - 用于创建 Meter 和指标
//
// 返回：
// - *metric.MeterProvider: 底层的指标提供者实例
//
// 使用示例：
//
//	meter := mp.MeterProvider().Meter("my-service")
//	counter, _ := meter.Int64Counter("requests_total")
//	counter.Add(ctx, 1)
func (mp *MetricsProvider) MeterProvider() *metric.MeterProvider {
	return mp.provider
}

// Shutdown 优雅关闭指标提供者。
//
// 功能说明：
// - 停止接受新的指标更新
// - 等待所有正在处理的指标数据完成
// - 刷新所有待导出的指标数据
// - 关闭导出器连接
//
// 参数：
// - ctx: 上下文，用于控制关闭超时
//
// 返回：
// - error: 如果关闭过程中出现错误，返回错误信息
//
// 使用示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := mp.Shutdown(ctx); err != nil {
//	    log.Printf("Error shutting down metrics provider: %v", err)
//	}
//
// 注意事项：
// - 应在应用程序退出时调用，确保所有指标数据都已导出
// - 建议设置超时上下文，避免关闭过程无限等待
// - 关闭后不应再使用该 MetricsProvider
func (mp *MetricsProvider) Shutdown(ctx context.Context) error {
	return mp.provider.Shutdown(ctx)
}
