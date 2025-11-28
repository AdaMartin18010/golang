// Package main 展示如何集成日志轮转和 OTLP
//
// 本示例展示：
// 1. 如何配置日志轮转
// 2. 如何将日志与 OTLP 集成
// 3. 如何在不同环境使用不同配置
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/pkg/logger"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/sampling"
	"log/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 根据环境选择日志配置
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var rotationCfg logger.RotationConfig
	switch env {
	case "production":
		rotationCfg = logger.ProductionRotationConfig("logs/app.log")
		log.Println("Using production log configuration")
	case "development":
		rotationCfg = logger.DevelopmentRotationConfig("logs/app.log")
		log.Println("Using development log configuration")
	default:
		rotationCfg = logger.DefaultRotationConfig("logs/app.log")
		log.Println("Using default log configuration")
	}

	// 2. 创建轮转日志记录器
	appLogger, err := logger.NewRotatingLogger(slog.LevelInfo, rotationCfg)
	if err != nil {
		log.Fatalf("Failed to create rotating logger: %v", err)
	}
	slog.SetDefault(appLogger.Logger)
	appLogger.Info("Application started", "env", env)

	// 3. 初始化 OTLP（可选）
	otlpEnabled := os.Getenv("OTLP_ENABLED") == "true"
	if otlpEnabled {
		sampler, err := sampling.NewProbabilisticSampler(0.5)
		if err != nil {
			appLogger.Error("Failed to create sampler", "error", err)
		} else {
			otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
				ServiceName:    "example-service",
				ServiceVersion: "v1.0.0",
				Endpoint:       os.Getenv("OTLP_ENDPOINT"),
				Insecure:       true,
				Sampler:        sampler,
				SampleRate:     0.5,
			})
			if err != nil {
				appLogger.Warn("Failed to initialize OTLP", "error", err)
			} else {
				defer otlpClient.Shutdown(ctx)
				appLogger.Info("OTLP initialized successfully")

				// 使用日志导出器（占位实现）
				logExporter := otlpClient.LogExporter()
				if logExporter != nil {
					appLogger.Info("Log exporter available", "enabled", true)
				}
			}
		}
	}

	// 4. 模拟业务逻辑
	appLogger.Info("Processing requests...")
	for i := 0; i < 10; i++ {
		appLogger.Info("Request processed",
			"request_id", i,
			"timestamp", time.Now().Unix(),
		)
		time.Sleep(100 * time.Millisecond)
	}

	// 5. 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Shutting down...")
}
