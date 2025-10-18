package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ===== 容器感知调度验证 =====

// TestContainerAwareScheduling 测试容器感知调度是否生效
func TestContainerAwareScheduling(t *testing.T) {
	gomaxprocs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()

	t.Logf("GOMAXPROCS: %d", gomaxprocs)
	t.Logf("NumCPU (物理核心数): %d", numCPU)

	// 检查是否在容器中
	if quota, period := readCgroupCPU(); quota > 0 && period > 0 {
		cpuLimit := float64(quota) / float64(period)
		t.Logf("cgroup CPU 限制: %.2f 核", cpuLimit)

		// 验证 GOMAXPROCS 是否正确设置
		expectedGOMAXPROCS := int(cpuLimit)
		if expectedGOMAXPROCS == 0 {
			expectedGOMAXPROCS = 1
		}

		// 允许 ±1 的误差（向上/向下取整）
		if gomaxprocs >= expectedGOMAXPROCS-1 && gomaxprocs <= expectedGOMAXPROCS+1 {
			t.Logf("✅ 容器感知调度已生效 (GOMAXPROCS=%d, CPU limit=%.2f)", gomaxprocs, cpuLimit)
		} else {
			t.Logf("⚠️  GOMAXPROCS(%d) 与 CPU 限制(%.2f)不匹配", gomaxprocs, cpuLimit)
		}
	} else {
		t.Logf("ℹ️  未检测到 cgroup CPU 限制（可能运行在物理机上）")
	}

	// 检查环境变量覆盖
	if env := os.Getenv("GOMAXPROCS"); env != "" {
		t.Logf("⚠️  GOMAXPROCS 环境变量已设置: %s (覆盖自动检测)", env)
	}
}

// TestCgroupDetection 测试 cgroup 检测功能
func TestCgroupDetection(t *testing.T) {
	quota, period := readCgroupCPU()

	t.Logf("cgroup 检测结果:")
	t.Logf("  cpu.cfs_quota_us: %d", quota)
	t.Logf("  cpu.cfs_period_us: %d", period)

	if quota > 0 && period > 0 {
		cpuLimit := float64(quota) / float64(period)
		t.Logf("  有效 CPU 核心数: %.2f", cpuLimit)
		t.Logf("  ✅ cgroup CPU 限制检测成功")
	} else {
		t.Logf("  ℹ️  未检测到 cgroup CPU 限制")
	}
}

// ===== 性能基准测试 =====

// cpuIntensiveTask 模拟 CPU 密集型任务
func cpuIntensiveTask(iterations int) int {
	sum := 0
	for i := 0; i < iterations; i++ {
		sum += i * i
	}
	return sum
}

// BenchmarkCorrectGOMAXPROCS 正确的 GOMAXPROCS 设置（容器感知）
func BenchmarkCorrectGOMAXPROCS(b *testing.B) {
	// 使用自动检测的 GOMAXPROCS
	gomaxprocs := runtime.GOMAXPROCS(0)
	b.Logf("使用 GOMAXPROCS=%d (自动检测)", gomaxprocs)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cpuIntensiveTask(10000)
		}
	})
}

// BenchmarkOversubscribedGOMAXPROCS 过度订阅的 GOMAXPROCS（模拟 Go 1.24 问题）
func BenchmarkOversubscribedGOMAXPROCS(b *testing.B) {
	// 模拟 Go 1.24 行为：使用物理 CPU 核心数
	oldGOMAXPROCS := runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(oldGOMAXPROCS)

	b.Logf("使用 GOMAXPROCS=%d (物理核心数)", runtime.NumCPU())

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cpuIntensiveTask(10000)
		}
	})
}

// BenchmarkSchedulingOverhead 调度开销对比
func BenchmarkSchedulingOverhead(b *testing.B) {
	tests := []struct {
		name       string
		gomaxprocs int
	}{
		{"GOMAXPROCS=1", 1},
		{"GOMAXPROCS=2", 2},
		{"GOMAXPROCS=4", 4},
		{"GOMAXPROCS=8", 8},
		{"GOMAXPROCS=16", 16},
		{"GOMAXPROCS=32", 32},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			oldGOMAXPROCS := runtime.GOMAXPROCS(tt.gomaxprocs)
			defer runtime.GOMAXPROCS(oldGOMAXPROCS)

			var wg sync.WaitGroup
			for i := 0; i < b.N; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					cpuIntensiveTask(1000)
				}()
			}
			wg.Wait()
		})
	}
}

// BenchmarkConcurrentLoad 并发负载测试
func BenchmarkConcurrentLoad(b *testing.B) {
	b.Run("Auto-Detected", func(b *testing.B) {
		// 使用自动检测的 GOMAXPROCS
		b.SetParallelism(runtime.GOMAXPROCS(0))
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := 0; i < 100; i++ {
					cpuIntensiveTask(100)
				}
			}
		})
	})

	b.Run("NumCPU", func(b *testing.B) {
		// 使用物理 CPU 核心数
		oldGOMAXPROCS := runtime.GOMAXPROCS(runtime.NumCPU())
		defer runtime.GOMAXPROCS(oldGOMAXPROCS)

		b.SetParallelism(runtime.NumCPU())
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				for i := 0; i < 100; i++ {
					cpuIntensiveTask(100)
				}
			}
		})
	})
}

// ===== 压力测试 =====

// TestHighConcurrency 高并发压力测试
func TestHighConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试")
	}

	const (
		duration   = 5 * time.Second
		goroutines = 1000
	)

	gomaxprocs := runtime.GOMAXPROCS(0)
	t.Logf("开始高并发压力测试 (GOMAXPROCS=%d)", gomaxprocs)

	var (
		wg        sync.WaitGroup
		done      = make(chan struct{})
		ops       atomic.Uint64
		startTime = time.Now()
	)

	// 启动多个 goroutine
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					cpuIntensiveTask(1000)
					ops.Add(1)
				}
			}
		}()
	}

	// 运行指定时间
	time.Sleep(duration)
	close(done)
	wg.Wait()

	elapsed := time.Since(startTime)
	totalOps := ops.Load()
	opsPerSec := float64(totalOps) / elapsed.Seconds()

	t.Logf("压力测试结果:")
	t.Logf("  运行时间: %v", elapsed)
	t.Logf("  总操作数: %d", totalOps)
	t.Logf("  操作/秒: %.0f", opsPerSec)
	t.Logf("  Goroutines: %d", goroutines)
	t.Logf("  GOMAXPROCS: %d", gomaxprocs)
}

// ===== 辅助函数 =====
// (readCgroupCPU, readCgroupV1CPU, readCgroupV2CPU 已在 main.go 中定义)

// ===== 诊断工具 =====

// TestRuntimeDiagnostics 运行时诊断信息
func TestRuntimeDiagnostics(t *testing.T) {
	t.Logf("=== Go 运行时诊断信息 ===")

	// 基本信息
	t.Logf("Go 版本: %s", runtime.Version())
	t.Logf("GOOS: %s", runtime.GOOS)
	t.Logf("GOARCH: %s", runtime.GOARCH)

	// CPU 信息
	t.Logf("\n=== CPU 信息 ===")
	t.Logf("GOMAXPROCS: %d", runtime.GOMAXPROCS(0))
	t.Logf("NumCPU: %d", runtime.NumCPU())

	// cgroup 信息（仅 Linux）
	if runtime.GOOS == "linux" {
		t.Logf("\n=== cgroup 信息 ===")
		if quota, period := readCgroupCPU(); quota > 0 && period > 0 {
			cpuLimit := float64(quota) / float64(period)
			t.Logf("cpu.cfs_quota_us: %d", quota)
			t.Logf("cpu.cfs_period_us: %d", period)
			t.Logf("CPU 限制: %.2f 核", cpuLimit)

			if runtime.GOMAXPROCS(0) < runtime.NumCPU() {
				t.Logf("✅ 检测到容器 CPU 限制，已自动调整 GOMAXPROCS")
			}
		} else {
			t.Logf("未检测到 cgroup CPU 限制")
		}
	}

	// 环境变量
	t.Logf("\n=== 环境变量 ===")
	if env := os.Getenv("GOMAXPROCS"); env != "" {
		t.Logf("GOMAXPROCS env: %s", env)
	} else {
		t.Logf("GOMAXPROCS env: (未设置)")
	}

	// Goroutine 信息
	t.Logf("\n=== Goroutine 信息 ===")
	t.Logf("当前 Goroutine 数: %d", runtime.NumGoroutine())

	// 内存信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	t.Logf("\n=== 内存信息 ===")
	t.Logf("Alloc: %d MB", m.Alloc/(1024*1024))
	t.Logf("Sys: %d MB", m.Sys/(1024*1024))
}

// ===== 对比测试 =====

// ComparePerformance 性能对比测试
// 使用方法：
//
//	# 默认（自动检测）
//	go test -bench=BenchmarkCorrectGOMAXPROCS -benchmem > auto.txt
//
//	# 模拟 Go 1.24（使用 NumCPU）
//	go test -bench=BenchmarkOversubscribedGOMAXPROCS -benchmem > oversubscribed.txt
//
//	# 对比结果
//	benchstat auto.txt oversubscribed.txt
func ComparePerformance() {
	// 此函数为说明性注释
}

// ===== 示例函数 =====

// Example_containerAwareScheduling 容器感知调度示例
func Example_containerAwareScheduling() {
	gomaxprocs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()

	fmt.Printf("GOMAXPROCS: %d\n", gomaxprocs)
	fmt.Printf("NumCPU: %d\n", numCPU)

	if gomaxprocs < numCPU {
		fmt.Println("容器感知调度已生效")
	} else {
		fmt.Println("运行在物理机或未限制容器中")
	}
}
