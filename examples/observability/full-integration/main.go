package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/pkg/observability"
	"github.com/yourusername/golang/pkg/observability/operational"
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

	// 3. 启动可观测性
	if err := obs.Start(); err != nil {
		log.Fatalf("Failed to start observability: %v", err)
	}
	defer obs.Stop(ctx)

	// 4. 创建运维控制端点
	operationalEndpoints := operational.NewOperationalEndpoints(operational.Config{
		Observability: obs,
		Port:          9090,
		PathPrefix:    "/ops",
		Enabled:       true,
	})

	// 5. 启动运维端点
	if err := operationalEndpoints.Start(); err != nil {
		log.Fatalf("Failed to start operational endpoints: %v", err)
	}
	defer operationalEndpoints.Stop(ctx)

	log.Println("Operational endpoints started on :9090")
	log.Println("Available endpoints:")
	log.Println("  - Health:      http://localhost:9090/ops/health")
	log.Println("  - Readiness:   http://localhost:9090/ops/ready")
	log.Println("  - Liveness:    http://localhost:9090/ops/live")
	log.Println("  - Metrics:     http://localhost:9090/ops/metrics")

	// 6. 创建主 HTTP 服务器
	mux := http.NewServeMux()
	port := 8080

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
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// 7. 启动主服务器
	go func() {
		log.Printf("Server starting on :%d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 8. 等待关闭信号
	log.Println("Application running. Press Ctrl+C to shutdown gracefully...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Application shutdown complete")
}
