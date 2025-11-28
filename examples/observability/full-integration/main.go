// Package main 展示完整的可观测性集成
//
// 本示例展示：
// 1. 如何集成 OTLP、系统监控、日志轮转
// 2. 如何获取平台信息
// 3. 如何监控系统资源
// 4. 如何在容器和 Kubernetes 环境中使用
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/pkg/logger"
	"github.com/yourusername/golang/pkg/observability"
	"github.com/yourusername/golang/pkg/sampling"
	"log/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 初始化日志轮转
	log.Println("Initializing logger...")
	rotationCfg := logger.ProductionRotationConfig("logs/app.log")
	appLogger, err := logger.NewRotatingLogger(slog.LevelInfo, rotationCfg)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	slog.SetDefault(appLogger.Logger)
	appLogger.Info("Application starting...")

	// 2. 创建采样器
	sampler, err := sampling.NewProbabilisticSampler(0.5)
	if err != nil {
		appLogger.Error("Failed to create sampler", "error", err)
		os.Exit(1)
	}

	// 3. 初始化完整的可观测性集成
	appLogger.Info("Initializing observability...")
	obs, err := observability.NewObservability(observability.Config{
		ServiceName:            "full-integration-example",
		ServiceVersion:         "v1.0.0",
		OTLPEndpoint:           getEnv("OTLP_ENDPOINT", "localhost:4317"),
		OTLPInsecure:           true,
		SampleRate:             0.5,
		MetricInterval:         10 * time.Second,
		TraceBatchTimeout:      5 * time.Second,
		TraceBatchSize:         512,
		EnableSystemMonitoring: true,
		SystemCollectInterval:  5 * time.Second,
	})
	if err != nil {
		appLogger.Error("Failed to initialize observability", "error", err)
		os.Exit(1)
	}

	// 4. 启动可观测性
	if err := obs.Start(); err != nil {
		appLogger.Error("Failed to start observability", "error", err)
		os.Exit(1)
	}
	defer obs.Stop(ctx)

	appLogger.Info("Observability initialized successfully")

	// 5. 显示平台信息
	platformInfo := obs.GetPlatformInfo()
	appLogger.Info("Platform Information",
		"os", platformInfo.OS,
		"arch", platformInfo.Arch,
		"go_version", platformInfo.GoVersion,
		"hostname", platformInfo.Hostname,
		"cpus", platformInfo.CPUs,
		"container_id", platformInfo.ContainerID,
		"container_name", platformInfo.ContainerName,
		"k8s_pod", platformInfo.KubernetesPod,
		"k8s_node", platformInfo.KubernetesNode,
		"virtualization", platformInfo.Virtualization,
	)

	// 6. 环境检测
	appLogger.Info("Environment Detection",
		"is_container", obs.IsContainer(),
		"is_kubernetes", obs.IsKubernetes(),
		"is_virtualized", obs.IsVirtualized(),
	)

	// 7. 使用追踪
	tracer := obs.Tracer("example-service")
	ctx, span := tracer.Start(ctx, "main-operation")
	defer span.End()

	appLogger.Info("Main operation started")

	// 8. 使用指标
	meter := obs.Meter("example-service")
	counter, _ := meter.Int64Counter("requests_total")
	counter.Add(ctx, 1)

	// 9. 定期显示系统资源
	systemMonitor := obs.GetSystemMonitor()
	if systemMonitor != nil {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					memStats := systemMonitor.GetMemoryStats()
					appLogger.Info("System Resources",
						"memory_alloc_mb", float64(memStats.Alloc)/1024/1024,
						"memory_total_mb", float64(memStats.TotalAlloc)/1024/1024,
						"gc_count", memStats.NumGC,
						"goroutines", systemMonitor.GetGoroutineCount(),
					)
				}
			}
		}()
	}

	// 10. 模拟业务逻辑
	appLogger.Info("Processing requests...")
	for i := 0; i < 10; i++ {
		ctx, span := tracer.Start(ctx, "process-request")
		appLogger.Info("Request processed",
			"request_id", i,
			"timestamp", time.Now().Unix(),
		)
		counter.Add(ctx, 1)
		time.Sleep(500 * time.Millisecond)
		span.End()
	}

	// 11. 等待中断信号
	appLogger.Info("Application is running. Press Ctrl+C to stop...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down...")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
