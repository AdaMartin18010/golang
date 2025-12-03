package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"

	"github.com/yourusername/golang/pkg/observability/system"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("advanced-features-example"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// 2. 创建 TracerProvider
	tp := trace.NewTracerProvider(trace.WithResource(res))
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down TracerProvider: %v", err)
		}
	}()

	// 3. 创建 MeterProvider
	mp := metric.NewMeterProvider(metric.WithResource(res))
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down MeterProvider: %v", err)
		}
	}()

	// 4. 创建系统监控器（启用所有功能）
	systemMonitor, err := system.NewSystemMonitor(system.SystemConfig{
		Meter:             mp.Meter("system"),
		Tracer:            tp.Tracer("system"),
		Enabled:           true,
		CollectInterval:   5 * time.Second,
		EnableDiskMonitor: true,
		EnableLoadMonitor: true,
		EnableAPMMonitor:  true,
		RateLimitConfig: &system.RateLimiterConfig{
			Meter:   mp.Meter("ratelimit"),
			Enabled: true,
			Limit:   100, // 每秒 100 个请求
			Window:  1 * time.Second,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create system monitor: %v", err)
	}

	// 5. 启动监控
	if err := systemMonitor.Start(ctx); err != nil {
		log.Fatalf("Failed to start system monitor: %v", err)
	}
	defer func() {
		if err := systemMonitor.Stop(ctx); err != nil {
			log.Printf("Error stopping system monitor: %v", err)
		}
	}()

	log.Println("Advanced features example started. Collecting metrics...")

	// 6. 演示 APM 监控
	apmMonitor := systemMonitor.GetAPMMonitor()
	if apmMonitor != nil {
		for i := 0; i < 10; i++ {
			start := time.Now()
			// 模拟请求处理
			time.Sleep(50 * time.Millisecond)
			duration := time.Since(start)

			// 记录请求
			apmMonitor.RecordRequest(ctx, duration, 200)
		}
		log.Println("Recorded 10 requests via APM monitor")
	}

	// 7. 演示限流器
	rateLimiter := systemMonitor.GetRateLimiter()
	if rateLimiter != nil {
		allowed := 0
		rejected := 0
		for i := 0; i < 150; i++ {
			if rateLimiter.Allow(ctx) {
				allowed++
			} else {
				rejected++
			}
			time.Sleep(10 * time.Millisecond)
		}
		log.Printf("Rate limiter: allowed=%d, rejected=%d", allowed, rejected)
	}

	// 8. 演示负载监控
	loadMonitor := systemMonitor.GetLoadMonitor()
	if loadMonitor != nil {
		for i := 0; i < 5; i++ {
			loadMonitor.RecordRequest(ctx)
			time.Sleep(100 * time.Millisecond)
		}
		log.Println("Recorded 5 requests via load monitor")
	}

	// 9. 演示指标聚合
	aggregator := systemMonitor.GetAggregator()
	if aggregator != nil {
		aggregator.RecordCounter("custom.counter", 1)
		aggregator.RecordGauge("custom.gauge", 42.0)
		aggregator.RecordHistogram("custom.histogram", 1.5)
		aggregator.RecordHistogram("custom.histogram", 2.0)
		aggregator.RecordHistogram("custom.histogram", 1.8)

		stats := aggregator.GetHistogramStats("custom.histogram")
		log.Printf("Histogram stats: count=%d, mean=%.2f, min=%.2f, max=%.2f",
			stats.Count, stats.Mean, stats.Min, stats.Max)
	}

	// 10. 检查 Kubernetes 信息
	k8sMonitor := systemMonitor.GetKubernetesMonitor()
	if k8sMonitor != nil && k8sMonitor.IsEnabled() {
		info := k8sMonitor.GetInfo()
		log.Printf("Kubernetes info: pod=%s, namespace=%s, node=%s",
			info.PodName, info.PodNamespace, info.NodeName)
	} else {
		log.Println("Not running in Kubernetes")
	}

	// 11. 健康检查
	healthChecker := systemMonitor.GetHealthChecker()
	if healthChecker != nil {
		status := healthChecker.Check(ctx)
		log.Printf("Health status: %+v", status)
	}

	// 模拟运行一段时间
	select {
	case <-time.After(30 * time.Second):
		log.Println("Application running for 30 seconds.")
	case <-ctx.Done():
		log.Println("Context cancelled.")
	}

	fmt.Println("Check your OpenTelemetry Collector for all metrics.")
}
