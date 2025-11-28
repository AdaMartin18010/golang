// Package main 展示如何使用系统监控功能
//
// 本示例展示：
// 1. 如何初始化系统监控
// 2. 如何获取平台信息
// 3. 如何获取资源统计
// 4. 如何集成到 OTLP
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
	log.Println("Initializing OTLP...")
	sampler, err := sampling.NewProbabilisticSampler(0.5)
	if err != nil {
		log.Fatalf("Failed to create sampler: %v", err)
	}

	otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:    "system-monitoring-example",
		ServiceVersion: "v1.0.0",
		Endpoint:       os.Getenv("OTLP_ENDPOINT"),
		Insecure:       true,
		Sampler:        sampler,
		SampleRate:     0.5,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize OTLP: %v", err)
	} else {
		defer otlpClient.Shutdown(ctx)
		log.Println("OTLP initialized successfully")
	}

	// 2. 创建系统监控器
	log.Println("Initializing system monitor...")
	systemMonitor, err := system.NewSystemMonitor(system.SystemConfig{
		Meter:           otlpClient.Meter("system"),
		Enabled:         true,
		CollectInterval: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create system monitor: %v", err)
	}

	// 3. 启动监控
	if err := systemMonitor.Start(); err != nil {
		log.Fatalf("Failed to start system monitor: %v", err)
	}
	defer systemMonitor.Stop()
	log.Println("System monitor started")

	// 4. 显示平台信息
	info := systemMonitor.GetPlatformInfo()
	log.Printf("Platform Information:")
	log.Printf("  OS: %s", info.OS)
	log.Printf("  Arch: %s", info.Arch)
	log.Printf("  Go Version: %s", info.GoVersion)
	log.Printf("  Hostname: %s", info.Hostname)
	log.Printf("  CPUs: %d", info.CPUs)
	log.Printf("  Container ID: %s", info.ContainerID)
	log.Printf("  Container Name: %s", info.ContainerName)
	log.Printf("  Kubernetes Pod: %s", info.KubernetesPod)
	log.Printf("  Kubernetes Node: %s", info.KubernetesNode)
	log.Printf("  Virtualization: %s", info.Virtualization)

	// 5. 环境检测
	log.Printf("\nEnvironment Detection:")
	if systemMonitor.IsContainer() {
		log.Printf("  ✓ Running in container")
		containerID, containerName := systemMonitor.GetPlatformInfo().ContainerID, systemMonitor.GetPlatformInfo().ContainerName
		log.Printf("    Container ID: %s", containerID)
		log.Printf("    Container Name: %s", containerName)
	} else {
		log.Printf("  ✗ Not in container")
	}

	if systemMonitor.IsKubernetes() {
		log.Printf("  ✓ Running in Kubernetes")
		pod, node := systemMonitor.GetPlatformInfo().KubernetesPod, systemMonitor.GetPlatformInfo().KubernetesNode
		log.Printf("    Pod: %s", pod)
		log.Printf("    Node: %s", node)
	} else {
		log.Printf("  ✗ Not in Kubernetes")
	}

	if systemMonitor.IsVirtualized() {
		log.Printf("  ✓ Running in virtualized environment")
		log.Printf("    Type: %s", systemMonitor.GetPlatformInfo().Virtualization)
	} else {
		log.Printf("  ✗ Not virtualized (bare-metal)")
	}

	// 6. 定期显示资源统计
	log.Printf("\nStarting resource monitoring...")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 获取内存统计
				memStats := systemMonitor.GetMemoryStats()
				log.Printf("\nResource Statistics:")
				log.Printf("  Memory Alloc: %d bytes (%.2f MB)", memStats.Alloc, float64(memStats.Alloc)/1024/1024)
				log.Printf("  Memory Total: %d bytes (%.2f MB)", memStats.TotalAlloc, float64(memStats.TotalAlloc)/1024/1024)
				log.Printf("  Memory Sys: %d bytes (%.2f MB)", memStats.Sys, float64(memStats.Sys)/1024/1024)
				log.Printf("  Heap Alloc: %d bytes (%.2f MB)", memStats.HeapAlloc, float64(memStats.HeapAlloc)/1024/1024)
				log.Printf("  GC Count: %d", memStats.NumGC)
				log.Printf("  Goroutines: %d", systemMonitor.GetGoroutineCount())
			}
		}
	}()

	// 7. 等待中断信号
	log.Println("\nSystem monitoring is running. Press Ctrl+C to stop...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down...")
}
