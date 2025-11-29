package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/pkg/observability"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 从配置文件加载配置
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 从应用配置创建可观测性配置
	obsConfig := observability.ConfigFromAppConfig(appConfig)

	// 3. 创建可观测性集成
	obs, err := observability.NewObservability(obsConfig)
	if err != nil {
		log.Fatalf("Failed to create observability: %v", err)
	}

	// 4. 应用告警规则
	observability.ApplyAlertRules(obs, appConfig.Observability.System.Alerts)

	// 5. 启动
	if err := obs.Start(); err != nil {
		log.Fatalf("Failed to start observability: %v", err)
	}
	defer func() {
		if err := obs.Stop(ctx); err != nil {
			log.Printf("Error stopping observability: %v", err)
		}
	}()

	log.Println("Config-driven observability started!")

	// 6. 使用追踪和指标
	tracer := obs.Tracer("my-service")
	ctx, span := tracer.Start(ctx, "config-driven-operation")
	defer span.End()

	meter := obs.Meter("my-service")
	counter, _ := meter.Int64Counter("requests_total")
	counter.Add(ctx, 1)

	// 模拟运行
	select {
	case <-time.After(30 * time.Second):
		log.Println("Application running for 30 seconds.")
	case <-ctx.Done():
		log.Println("Context cancelled.")
	}

	fmt.Println("Check your OpenTelemetry Collector for metrics and traces.")
}
