package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/pkg/observability"
	"github.com/yourusername/golang/pkg/observability/operational"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 加载配置
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

	// 5. 启动可观测性
	if err := obs.Start(); err != nil {
		log.Fatalf("Failed to start observability: %v", err)
	}

	// 6. 创建运维控制端点
	operationalEndpoints := operational.NewOperationalEndpoints(operational.Config{
		Observability: obs,
		Port:          9090,
		PathPrefix:    "/ops",
		Enabled:       true,
	})

	// 7. 启动运维端点
	if err := operationalEndpoints.Start(); err != nil {
		log.Fatalf("Failed to start operational endpoints: %v", err)
	}

	log.Println("Operational endpoints started on :9090")
	log.Println("Available endpoints:")
	log.Println("  - Health:      http://localhost:9090/ops/health")
	log.Println("  - Readiness:   http://localhost:9090/ops/ready")
	log.Println("  - Liveness:    http://localhost:9090/ops/live")
	log.Println("  - Metrics:     http://localhost:9090/ops/metrics")
	log.Println("  - Prometheus:  http://localhost:9090/ops/metrics/prometheus")
	log.Println("  - Dashboard:   http://localhost:9090/ops/dashboard")
	log.Println("  - Diagnostics: http://localhost:9090/ops/diagnostics")
	log.Println("  - Info:        http://localhost:9090/ops/info")
	log.Println("  - Version:     http://localhost:9090/ops/version")
	log.Println("  - Pprof:       http://localhost:9090/ops/debug/pprof/")

	// 8. 创建主 HTTP 服务器
	mux := http.NewServeMux()

	// 使用运维中间件
	mux.HandleFunc("/", operational.RecoveryMiddleware(
		operational.SecurityHeadersMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 使用追踪
				ctx, span := obs.Tracer("server").Start(r.Context(), "handler")
				defer span.End()

				// 使用指标
				meter := obs.Meter("server")
				counter, _ := meter.Int64Counter("requests_total")
				counter.Add(ctx, 1)

				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "Hello, World! Time: %s\n", time.Now().Format(time.RFC3339))
			}),
		),
	).ServeHTTP)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler: mux,
	}

	// 9. 启动主服务器
	go func() {
		log.Printf("Server starting on :%d", appConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 10. 创建优雅关闭管理器
	shutdownManager := operational.NewShutdownManager(30 * time.Second)

	// 注册关闭函数
	shutdownManager.Register(operational.GracefulShutdown("http-server", func(ctx context.Context) error {
		return server.Shutdown(ctx)
	}))
	shutdownManager.Register(operational.GracefulShutdown("observability", func(ctx context.Context) error {
		return obs.Stop(ctx)
	}))
	shutdownManager.Register(operational.GracefulShutdown("operational-endpoints", func(ctx context.Context) error {
		return operationalEndpoints.Stop(ctx)
	}))

	// 11. 演示熔断器
	circuitBreaker := operational.NewCircuitBreaker(operational.CircuitBreakerConfig{
		Name:         "external-api",
		MaxFailures:  5,
		ResetTimeout: 60 * time.Second,
	})

	// 模拟一些操作
	go func() {
		for i := 0; i < 10; i++ {
			err := circuitBreaker.Execute(ctx, func() error {
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

	// 12. 等待关闭信号
	log.Println("Application running. Press Ctrl+C to shutdown gracefully...")
	if err := shutdownManager.WaitForShutdown(); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Application shutdown complete")
}
