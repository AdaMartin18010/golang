// Package main 展示如何完整使用可观测性功能
//
// 本示例展示：
// 1. 如何初始化 OTLP（追踪、指标）
// 2. 如何使用日志轮转功能
// 3. 如何集成 eBPF 收集器（框架）
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/pkg/logger"
	"github.com/yourusername/golang/pkg/observability/ebpf"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/sampling"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"log/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 初始化 OTLP（追踪和指标）
	log.Println("Initializing OTLP...")
	sampler, err := sampling.NewProbabilisticSampler(0.5)
	if err != nil {
		log.Fatalf("Failed to create sampler: %v", err)
	}

	otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:       "example-service",
		ServiceVersion:    "v1.0.0",
		Endpoint:          "localhost:4317",
		Insecure:          true,
		Sampler:           sampler,
		SampleRate:        0.5,
		MetricInterval:    10 * time.Second,
		TraceBatchTimeout: 5 * time.Second,
		TraceBatchSize:    512,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize OTLP: %v (continuing without OTLP)", err)
	} else {
		defer otlpClient.Shutdown(ctx)
		log.Println("OTLP initialized successfully")
	}

	// 2. 初始化日志轮转
	log.Println("Initializing rotating logger...")
	rotationCfg := logger.ProductionRotationConfig("logs/app.log")
	appLogger, err := logger.NewRotatingLogger(slog.LevelInfo, rotationCfg)
	if err != nil {
		log.Fatalf("Failed to create rotating logger: %v", err)
	}
	slog.SetDefault(appLogger.Logger)
	appLogger.Info("Application started", "version", "v1.0.0")

	// 3. 使用追踪
	tracer := otlpClient.Tracer("example-service")
	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	// 4. 使用指标
	meter := otlpClient.Meter("example-service")
	counter, err := meter.Int64Counter(
		"requests_total",
		metric.WithDescription("Total number of requests"),
	)
	if err != nil {
		log.Printf("Failed to create counter: %v", err)
	} else {
		counter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("method", "GET"),
			attribute.String("path", "/api/users"),
		))
	}

	// 5. 初始化 eBPF 收集器（框架）
	log.Println("Initializing eBPF collector...")
	ebpfCollector, err := ebpf.NewCollector(ebpf.Config{
		Tracer:                    tracer,
		Meter:                     meter,
		Enabled:                   false, // 设置为 true 需要实际的 eBPF 程序
		CollectInterval:          5 * time.Second,
		EnableSyscallTracking:     true,
		EnableNetworkMonitoring:   true,
		EnablePerformanceProfiling: false,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize eBPF collector: %v", err)
	} else {
		if err := ebpfCollector.Start(); err != nil {
			log.Printf("Warning: Failed to start eBPF collector: %v", err)
		} else {
			defer ebpfCollector.Stop()
			log.Println("eBPF collector initialized (framework only)")
		}
	}

	// 6. 模拟业务逻辑
	log.Println("Running example...")
	runExample(ctx, tracer, meter, appLogger)

	// 7. 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

func runExample(ctx context.Context, tracer interface{}, meter metric.Meter, logger *logger.Logger) {
	// 模拟一些业务操作
	for i := 0; i < 10; i++ {
		logger.Info("Processing request", "request_id", i)
		time.Sleep(100 * time.Millisecond)
	}
}
