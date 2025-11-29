package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/golang/pkg/observability"
	"github.com/yourusername/golang/pkg/observability/system"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 创建完整的可观测性集成
	obs, err := observability.NewObservability(observability.Config{
		ServiceName:            "my-service",
		ServiceVersion:         "v1.0.0",
		OTLPEndpoint:           "localhost:4317",
		OTLPInsecure:           true,
		SampleRate:             0.5, // 50% 采样率
		MetricInterval:         10 * time.Second,
		TraceBatchTimeout:      5 * time.Second,
		TraceBatchSize:         512,
		EnableSystemMonitoring: true,
		SystemCollectInterval:  5 * time.Second,
		EnableDiskMonitor:      true,
		EnableLoadMonitor:       true,
		EnableAPMMonitor:        true,
		RateLimitConfig: &system.RateLimiterConfig{
			Meter:   nil, // 将在内部设置
			Enabled: true,
			Limit:   100, // 每秒 100 个请求
			Window:  1 * time.Second,
		},
		HealthThresholds: system.DefaultHealthThresholds(),
	})
	if err != nil {
		log.Fatalf("Failed to create observability: %v", err)
	}

	// 2. 启动所有监控
	if err := obs.Start(); err != nil {
		log.Fatalf("Failed to start observability: %v", err)
	}
	defer func() {
		if err := obs.Stop(ctx); err != nil {
			log.Printf("Error stopping observability: %v", err)
		}
	}()

	log.Println("Complete observability integration started!")

	// 3. 使用追踪
	tracer := obs.Tracer("my-service")
	ctx, span := tracer.Start(ctx, "main-operation")
	defer span.End()

	// 4. 使用指标
	meter := obs.Meter("my-service")
	counter, _ := meter.Int64Counter("requests_total")
	counter.Add(ctx, 1)

	// 5. 使用 APM 监控
	apmMonitor := obs.GetAPMMonitor()
	if apmMonitor != nil {
		start := time.Now()
		// 模拟业务操作
		time.Sleep(50 * time.Millisecond)
		duration := time.Since(start)
		apmMonitor.RecordRequest(ctx, duration, 200)
	}

	// 6. 使用限流器
	rateLimiter := obs.GetRateLimiter()
	if rateLimiter != nil {
		for i := 0; i < 10; i++ {
			if rateLimiter.Allow(ctx) {
				log.Printf("Request %d allowed", i+1)
			} else {
				log.Printf("Request %d rejected", i+1)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	// 7. 使用告警系统
	alertManager := obs.GetAlertManager()
	if alertManager != nil {
		// 检查 CPU 使用率
		alertManager.Check(ctx, "system.cpu.usage", 85.0, nil)
	}

	// 8. 使用诊断工具
	diagnostics := obs.GetDiagnostics()
	if diagnostics != nil {
		report, err := diagnostics.GenerateReport(ctx)
		if err == nil {
			log.Printf("Diagnostic report generated: %d issues found", len(report.Issues))
			if len(report.Recommendations) > 0 {
				log.Println("Recommendations:")
				for _, rec := range report.Recommendations {
					log.Printf("  - %s", rec)
				}
			}
		}
	}

	// 9. 使用资源预测
	predictor := obs.GetPredictor()
	if predictor != nil {
		prediction, err := predictor.Predict(ctx, "system.memory.usage", 1*time.Hour)
		if err == nil {
			log.Printf("Memory usage prediction: %.2f (confidence: %.2f, trend: %s)",
				prediction.PredictedValue, prediction.Confidence, prediction.Trend)
		}
	}

	// 10. 导出仪表板数据
	dashboardExporter := obs.GetDashboardExporter()
	if dashboardExporter != nil {
		// JSON 格式
		jsonData, err := dashboardExporter.ExportJSON(ctx)
		if err == nil {
			log.Printf("Dashboard data (JSON): %d bytes", len(jsonData))
		}

		// Prometheus 格式
		promData, err := dashboardExporter.ExportForPrometheus(ctx)
		if err == nil {
			log.Printf("Dashboard data (Prometheus):\n%s", promData)
		}
	}

	// 11. 导出指标
	metricsExporter := obs.GetMetricsExporter()
	if metricsExporter != nil {
		snapshot, err := metricsExporter.Export(ctx)
		if err == nil {
			log.Printf("Metrics snapshot: %d metrics at %s", len(snapshot.Metrics), snapshot.Timestamp)
		}

		jsonMetrics, err := metricsExporter.ExportJSON(ctx)
		if err == nil {
			log.Printf("Metrics JSON: %d bytes", len(jsonMetrics))
		}
	}

	// 12. 检查平台信息
	platformInfo := obs.GetPlatformInfo()
	log.Printf("Platform: %s/%s, Container: %v, Kubernetes: %v, Virtualized: %v",
		platformInfo.OS, platformInfo.Arch,
		obs.IsContainer(), obs.IsKubernetes(), obs.IsVirtualized())

	// 13. 检查 Kubernetes 信息
	if obs.IsKubernetes() {
		k8sInfo := obs.GetKubernetesInfo()
		log.Printf("Kubernetes: Pod=%s, Namespace=%s, Node=%s",
			k8sInfo.PodName, k8sInfo.PodNamespace, k8sInfo.NodeName)
	}

	// 14. 健康检查
	systemMonitor := obs.GetSystemMonitor()
	if systemMonitor != nil {
		healthChecker := systemMonitor.GetHealthChecker()
		if healthChecker != nil {
			status := healthChecker.Check(ctx)
			log.Printf("Health status: Healthy=%v, Degraded=%v, Message=%s",
				status.Healthy, status.Degraded, status.Message)
		}
	}

	// 模拟运行一段时间
	log.Println("Running for 30 seconds...")
	select {
	case <-time.After(30 * time.Second):
		log.Println("Application running completed.")
	case <-ctx.Done():
		log.Println("Context cancelled.")
	}

	fmt.Println("\n✅ Complete observability integration demo completed!")
	fmt.Println("Check your OpenTelemetry Collector for all metrics and traces.")
}
