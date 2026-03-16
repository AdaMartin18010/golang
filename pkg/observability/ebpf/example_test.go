package ebpf

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// ExampleCollector 演示如何使用 eBPF 收集器
func ExampleCollector() {
	// 设置 OpenTelemetry
	tracer, meter, shutdown := setupOTel()
	defer shutdown()

	// 创建 eBPF 收集器
	collector, err := NewCollector(Config{
		Tracer:                  tracer,
		Meter:                   meter,
		Enabled:                 true,
		CollectInterval:         5 * time.Second,
		EnableSyscallTracking:   true,
		EnableNetworkMonitoring: true,
	})
	if err != nil {
		log.Printf("Failed to create collector: %v", err)
		return
	}

	// 启动收集器
	if err := collector.Start(); err != nil {
		log.Printf("Failed to start collector: %v", err)
		return
	}
	defer collector.Stop()

	// 运行一段时间收集数据
	ctx := context.Background()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 6; i++ { // 运行 60 秒
		<-ticker.C

		// 打印统计信息
		printStats(ctx, collector)
	}
}

// ExampleSyscallTracer 演示如何使用系统调用追踪器
func ExampleSyscallTracer() {
	// 设置 OpenTelemetry
	tracer, meter, shutdown := setupOTel()
	defer shutdown()

	// 创建系统调用追踪器
	syscallTracer, err := NewSyscallTracer(SyscallTracerConfig{
		Tracer:     tracer,
		Meter:      meter,
		Enabled:    true,
		BufferSize: 64 * 1024,
	})
	if err != nil {
		log.Printf("Failed to create syscall tracer: %v", err)
		return
	}

	// 启动追踪器
	if err := syscallTracer.Start(); err != nil {
		log.Printf("Failed to start syscall tracer: %v", err)
		return
	}
	defer syscallTracer.Stop()

	// 运行 30 秒
	time.Sleep(30 * time.Second)

	// 获取统计
	ctx := context.Background()
	stats, err := syscallTracer.GetSyscallStats(ctx)
	if err != nil {
		log.Printf("Failed to get syscall stats: %v", err)
		return
	}

	fmt.Println("System call statistics:")
	for syscallID, count := range stats {
		fmt.Printf("  Syscall %d: %d calls\n", syscallID, count)
	}
}

// ExampleNetworkTracer 演示如何使用网络追踪器
func ExampleNetworkTracer() {
	// 设置 OpenTelemetry
	tracer, meter, shutdown := setupOTel()
	defer shutdown()

	// 创建网络追踪器
	networkTracer, err := NewNetworkTracer(NetworkTracerConfig{
		Tracer:        tracer,
		Meter:         meter,
		Enabled:       true,
		TrackInbound:  true,
		TrackOutbound: true,
		BufferSize:    64 * 1024,
	})
	if err != nil {
		log.Printf("Failed to create network tracer: %v", err)
		return
	}

	// 启动追踪器
	if err := networkTracer.Start(); err != nil {
		log.Printf("Failed to start network tracer: %v", err)
		return
	}
	defer networkTracer.Stop()

	// 定期打印网络统计
	ctx := context.Background()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 12; i++ {
		<-ticker.C

		// 活跃连接数
		activeConns, err := networkTracer.GetActiveConnections(ctx)
		if err != nil {
			log.Printf("Failed to get active connections: %v", err)
			continue
		}

		// 连接统计
		connStats, err := networkTracer.GetConnectionStats(ctx)
		if err != nil {
			log.Printf("Failed to get connection stats: %v", err)
			continue
		}

		// 连接详情
		details, err := networkTracer.GetConnectionDetails(ctx)
		if err != nil {
			log.Printf("Failed to get connection details: %v", err)
			continue
		}

		fmt.Printf("\n=== Network Stats (%d) ===\n", i+1)
		fmt.Printf("Active connections: %d\n", activeConns)
		fmt.Printf("Connection stats by PID: %+v\n", connStats)
		fmt.Println("Connection details:")
		for _, conn := range details {
			fmt.Printf("  %s:%d -> %s:%d (sent: %d, recv: %d)\n",
				conn.SrcIP, conn.SrcPort,
				conn.DstIP, conn.DstPort,
				conn.BytesSent, conn.BytesRecv)
		}
	}
}

// printStats 打印收集器的统计信息
func printStats(ctx context.Context, collector *Collector) {
	status := collector.GetStatus()
	fmt.Printf("\n=== Collector Status ===\n")
	fmt.Printf("Enabled: %v, Started: %v\n", status.Enabled, status.Started)
	fmt.Printf("Syscall tracking: %v, Network monitoring: %v\n",
		status.SyscallTracerEnabled, status.NetworkTracerEnabled)

	// 系统调用统计
	syscallStats, err := collector.GetSyscallStats(ctx)
	if err != nil {
		log.Printf("Failed to get syscall stats: %v", err)
	} else if len(syscallStats) > 0 {
		fmt.Println("Syscall stats:")
		for id, count := range syscallStats {
			fmt.Printf("  Syscall %d: %d\n", id, count)
		}
	}

	// 网络统计
	activeConns, err := collector.GetActiveConnections(ctx)
	if err != nil {
		log.Printf("Failed to get active connections: %v", err)
	} else {
		fmt.Printf("Active connections: %d\n", activeConns)
	}

	connStats, err := collector.GetConnectionStats(ctx)
	if err != nil {
		log.Printf("Failed to get connection stats: %v", err)
	} else if len(connStats) > 0 {
		fmt.Printf("Connection stats by PID: %+v\n", connStats)
	}
}

// setupOTel 设置 OpenTelemetry
func setupOTel() (trace.Tracer, metric.Meter, func()) {
	ctx := context.Background()

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ebpf-example"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// 创建 trace exporter
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("Failed to create trace exporter: %v", err)
	}

	// 创建 tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
	)

	// 创建 metric exporter
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		log.Fatalf("Failed to create metric exporter: %v", err)
	}

	// 创建 meter provider
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)

	// 设置全局 provider
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)

	// 返回 tracer, meter 和关闭函数
	return tracerProvider.Tracer("ebpf-example"),
		meterProvider.Meter("ebpf-example"),
		func() {
			tracerProvider.Shutdown(ctx)
			meterProvider.Shutdown(ctx)
		}
}
