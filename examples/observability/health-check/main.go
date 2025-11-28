// Package main 展示如何使用健康检查功能
//
// 本示例展示：
// 1. 如何配置健康检查
// 2. 如何定期执行健康检查
// 3. 如何根据健康状态做出决策
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/observability/system"
	"github.com/yourusername/golang/pkg/sampling"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 初始化 OTLP
	sampler, _ := sampling.NewProbabilisticSampler(0.5)
	otlpClient, _ := otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:    "health-check-example",
		ServiceVersion: "v1.0.0",
		Endpoint:       os.Getenv("OTLP_ENDPOINT"),
		Insecure:       true,
		Sampler:        sampler,
	})
	defer otlpClient.Shutdown(ctx)

	// 2. 创建系统监控器（带健康检查）
	systemMonitor, err := system.NewSystemMonitor(system.SystemConfig{
		Meter:            otlpClient.Meter("system"),
		Enabled:          true,
		CollectInterval:  5 * time.Second,
		EnableDiskMonitor: true,
		HealthThresholds: system.HealthThresholds{
			MaxMemoryUsage: 85.0,  // 85% 内存使用率
			MaxCPUUsage:    90.0,   // 90% CPU 使用率
			MaxGoroutines:  5000,   // 5000 个 Goroutine
			MinGCInterval:  500 * time.Millisecond,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create system monitor: %v", err)
	}

	// 3. 启动监控
	if err := systemMonitor.Start(); err != nil {
		log.Fatalf("Failed to start system monitor: %v", err)
	}
	defer systemMonitor.Stop()

	log.Println("System monitor started with health checking")

	// 4. 定期执行健康检查
	healthChecker := systemMonitor.GetHealthChecker()
	healthChecker.CheckPeriodically(ctx, func(status system.HealthStatus) {
		if status.Healthy {
			log.Printf("✓ Health check passed: %s", status.Message)
		} else {
			log.Printf("✗ Health check failed: %s", status.Message)
			log.Printf("  Memory: %.2f%%, CPU: %.2f%%, Goroutines: %d, GC: %d",
				status.MemoryUsage, status.CPUUsage, status.Goroutines, status.GC)
		}
	})

	// 5. 模拟一些负载
	go func() {
		for i := 0; i < 100; i++ {
			go func() {
				time.Sleep(10 * time.Second)
			}()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 6. 等待中断信号
	log.Println("Health checking is running. Press Ctrl+C to stop...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
