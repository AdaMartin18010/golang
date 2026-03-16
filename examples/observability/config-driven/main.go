package main

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/golang/pkg/observability"
	"github.com/yourusername/golang/pkg/observability/otlp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 创建可观测性配置
	obsConfig := observability.Config{
		ServiceName:    "example-service",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		OTLP: otlp.Config{
			Endpoint: "localhost:4317",
			Insecure: true,
		},
		Metrics: observability.MetricsConfig{
			Enabled:      true,
			ExportPeriod: 30 * time.Second,
		},
		Tracing: observability.TracingConfig{
			Enabled:         true,
			SampleRate:      1.0,
			ExportBatchSize: 100,
		},
	}

	// 2. 创建可观测性集成
	obs, err := observability.NewObservability(obsConfig)
	if err != nil {
		log.Fatalf("Failed to create observability: %v", err)
	}

	// 3. 启动
	if err := obs.Start(); err != nil {
		log.Fatalf("Failed to start observability: %v", err)
	}
	defer func() {
		if err := obs.Stop(ctx); err != nil {
			log.Printf("Error stopping observability: %v", err)
		}
	}()

	log.Println("Observability started!")

	// 4. 使用追踪和指标
	tracer := obs.Tracer("my-service")
	ctx, span := tracer.Start(ctx, "example-operation")
	defer span.End()

	// 记录一些指标
	counter := obs.Counter("example.counter", "Example counter")
	counter.Add(ctx, 1)

	log.Println("Operation traced and metrics recorded")

	// 保持运行一段时间
	time.Sleep(5 * time.Second)
}
