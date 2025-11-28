// Package main 展示如何从配置文件驱动日志和 OTLP 初始化
//
// 本示例展示：
// 1. 如何从配置文件加载日志配置
// 2. 如何从配置文件加载 OTLP 配置
// 3. 如何统一管理所有可观测性配置
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/pkg/logger"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/sampling"
	"log/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 从配置创建日志记录器
	appLogger, err := logger.CreateLoggerFromConfig(
		cfg.Logging.Level,
		cfg.Logging.Format,
		cfg.Logging.Output,
		cfg.Logging.OutputPath,
		logger.RotationConfig{
			Filename:   cfg.Logging.OutputPath,
			MaxSize:    cfg.Logging.Rotation.MaxSize,
			MaxBackups: cfg.Logging.Rotation.MaxBackups,
			MaxAge:     cfg.Logging.Rotation.MaxAge,
			Compress:   cfg.Logging.Rotation.Compress,
			LocalTime:  true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	slog.SetDefault(appLogger.Logger)
	appLogger.Info("Application started")

	// 3. 从配置初始化 OTLP
	if cfg.OTLP.Endpoint != "" {
		sampler, err := sampling.NewProbabilisticSampler(0.5)
		if err != nil {
			appLogger.Warn("Failed to create sampler", "error", err)
		} else {
			otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
				ServiceName:       cfg.OTLP.ServiceName,
				ServiceVersion:    cfg.OTLP.ServiceVersion,
				Endpoint:          cfg.OTLP.Endpoint,
				Insecure:          cfg.OTLP.Insecure,
				Sampler:           sampler,
				SampleRate:        0.5,
				MetricInterval:    10 * time.Second,
				TraceBatchTimeout: 5 * time.Second,
				TraceBatchSize:    512,
			})
			if err != nil {
				appLogger.Warn("Failed to initialize OTLP", "error", err)
			} else {
				defer otlpClient.Shutdown(ctx)
				appLogger.Info("OTLP initialized successfully",
					"endpoint", cfg.OTLP.Endpoint,
					"service_name", cfg.OTLP.ServiceName,
				)

				// 使用追踪和指标
				tracer := otlpClient.Tracer("example-service")
				meter := otlpClient.Meter("example-service")

				// 记录一些指标
				counter, _ := meter.Int64Counter("requests_total")
				counter.Add(ctx, 1)

				// 创建追踪
				ctx, span := tracer.Start(ctx, "example-operation")
				appLogger.Info("Operation started", "trace_id", "example")
				time.Sleep(100 * time.Millisecond)
				span.End()
			}
		}
	}

	// 4. 模拟业务逻辑
	appLogger.Info("Processing requests...")
	for i := 0; i < 5; i++ {
		appLogger.Info("Request processed",
			"request_id", i,
			"timestamp", time.Now().Unix(),
		)
		time.Sleep(200 * time.Millisecond)
	}

	// 5. 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down...")
}
