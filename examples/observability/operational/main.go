package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/golang/pkg/observability"
	"github.com/yourusername/golang/pkg/observability/operational"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 创建可观测性集成
	obs, err := observability.NewObservability(observability.Config{
		ServiceName:            "operational-demo",
		ServiceVersion:         "1.0.0",
		OTLPEndpoint:           "localhost:4317",
		OTLPInsecure:           true,
		EnableSystemMonitoring: true,
		SystemCollectInterval:  5 * time.Second,
		EnableDiskMonitor:      true,
		EnableLoadMonitor:      true,
		EnableAPMMonitor:       true,
	})
	if err != nil {
		log.Fatalf("Failed to create observability: %v", err)
	}

	// 2. 启动可观测性
	if err := obs.Start(); err != nil {
		log.Fatalf("Failed to start observability: %v", err)
	}

	// 3. 创建运维控制端点
	endpoints := operational.NewOperationalEndpoints(operational.Config{
		Observability: obs,
		Port:          9090,
		PathPrefix:    "/ops",
		Enabled:       true,
	})

	// 4. 启动运维端点
	if err := endpoints.Start(); err != nil {
		log.Fatalf("Failed to start operational endpoints: %v", err)
	}

	log.Println("Operational endpoints started on :9090")
	log.Println("Available endpoints:")
	log.Println("  - Health:     http://localhost:9090/ops/health")
	log.Println("  - Readiness:  http://localhost:9090/ops/ready")
	log.Println("  - Liveness:   http://localhost:9090/ops/live")
	log.Println("  - Metrics:    http://localhost:9090/ops/metrics")
	log.Println("  - Prometheus: http://localhost:9090/ops/metrics/prometheus")
	log.Println("  - Dashboard:  http://localhost:9090/ops/dashboard")
	log.Println("  - Diagnostics: http://localhost:9090/ops/diagnostics")
	log.Println("  - Info:       http://localhost:9090/ops/info")
	log.Println("  - Version:    http://localhost:9090/ops/version")
	log.Println("  - Pprof:      http://localhost:9090/ops/debug/pprof/")

	// 5. 创建优雅关闭管理器
	shutdownManager := operational.NewShutdownManager(30 * time.Second)

	// 注册关闭函数
	shutdownManager.Register(operational.GracefulShutdown("observability", func(ctx context.Context) error {
		return obs.Stop(ctx)
	}))
	shutdownManager.Register(operational.GracefulShutdown("operational-endpoints", func(ctx context.Context) error {
		return endpoints.Stop(ctx)
	}))

	// 6. 演示熔断器
	circuitBreaker := operational.NewCircuitBreaker(operational.CircuitBreakerConfig{
		Name:         "external-api",
		MaxFailures:  5,
		ResetTimeout: 60 * time.Second,
	})

	// 模拟一些操作
	go func() {
		for i := 0; i < 10; i++ {
			err := circuitBreaker.Execute(ctx, func() error {
				// 模拟外部 API 调用
				time.Sleep(100 * time.Millisecond)
				if i%3 == 0 {
					return fmt.Errorf("simulated error")
				}
				return nil
			})
			if err != nil {
				log.Printf("Circuit breaker error: %v (state: %s)", err, circuitBreaker.GetState())
			} else {
				log.Printf("Operation succeeded (state: %s)", circuitBreaker.GetState())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 7. 等待关闭信号
	log.Println("Application running. Press Ctrl+C to shutdown gracefully...")
	if err := shutdownManager.WaitForShutdown(); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Application shutdown complete")
}
