package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/pkg/observability/ebpf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	log.Println("🚀 eBPF 监控示例")
	log.Println("使用 Cilium eBPF 库进行系统级监控")
	log.Println("")

	// 检查是否为 Linux
	if os.Getenv("GOOS") != "linux" {
		log.Println("⚠️  eBPF 需要 Linux 环境")
		log.Println("当前系统不支持，仅展示框架使用")
	}

	// 初始化 OpenTelemetry
	ctx := context.Background()

	// 创建 stdout exporter (用于演示)
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal("Failed to create exporter:", err)
	}

	// 创建 tracer provider
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ebpf-monitoring-example"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)

	// 获取 tracer 和 meter
	tracer := otel.Tracer("ebpf-example")
	meter := otel.Meter("ebpf-example")

	// 创建 eBPF 收集器
	log.Println("📊 创建 eBPF 收集器...")
	collector, err := ebpf.NewCollector(ebpf.Config{
		Tracer:                  tracer,
		Meter:                   meter,
		Enabled:                 true,
		EnableSyscallTracking:   true,
		EnableNetworkMonitoring: true,
		CollectInterval:         5 * time.Second,
	})
	if err != nil {
		log.Fatal("Failed to create eBPF collector:", err)
	}

	// 启动收集
	log.Println("▶️  启动 eBPF 监控...")
	if err := collector.Start(); err != nil {
		log.Printf("⚠️  启动 eBPF 监控失败: %v", err)
		log.Println("这是正常的，eBPF 需要：")
		log.Println("  1. Linux 环境")
		log.Println("  2. Root 权限或 CAP_BPF capability")
		log.Println("  3. 编译的 eBPF 程序 (make generate-ebpf)")
		log.Println("")
		log.Println("继续运行以展示框架集成...")
	} else {
		log.Println("✅ eBPF 监控已启动")
		log.Println("")
		log.Println("监控功能：")
		log.Println("  ✅ 系统调用追踪 (sys_enter/sys_exit)")
		log.Println("  ✅ TCP 连接监控 (connect/accept/close)")
		log.Println("  ✅ 网络流量统计 (bytes sent/recv)")
		log.Println("  ✅ 延迟测量 (syscall/connection latency)")
	}
	defer collector.Stop()

	// 模拟一些工作负载
	log.Println("")
	log.Println("🔄 模拟工作负载...")
	go simulateWorkload(ctx, tracer)

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	log.Println("")
	log.Println("📡 eBPF 监控运行中...")
	log.Println("按 Ctrl+C 停止")
	log.Println("")

	<-sigCh
	log.Println("")
	log.Println("🛑 收到停止信号，正在清理...")
}

// simulateWorkload 模拟工作负载
func simulateWorkload(ctx context.Context, tracer trace.Tracer) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 创建一些 span 来展示追踪
			ctx, span := tracer.Start(ctx, "simulated-work")
			time.Sleep(100 * time.Millisecond)

			// 模拟一些系统调用
			_ = os.Getpid()
			_ = os.Getuid()

			span.End()
		}
	}
}
