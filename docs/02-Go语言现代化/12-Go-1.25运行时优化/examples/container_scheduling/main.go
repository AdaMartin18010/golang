package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("===================================")
	fmt.Println("  Go 1.25 容器感知调度验证工具")
	fmt.Println("===================================")

	// 基本运行时信息
	fmt.Printf("Go 版本:        %s\n", runtime.Version())
	fmt.Printf("操作系统:       %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	// CPU 信息
	fmt.Println("--- CPU 信息 ---")
	gomaxprocs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	fmt.Printf("GOMAXPROCS:     %d\n", gomaxprocs)
	fmt.Printf("NumCPU:         %d (物理核心数)\n", numCPU)
	fmt.Println()

	// cgroup 检测（仅 Linux）
	if runtime.GOOS == "linux" {
		fmt.Println("--- cgroup 检测 ---")
		if quota, period := readCgroupCPU(); quota > 0 && period > 0 {
			cpuLimit := float64(quota) / float64(period)
			fmt.Printf("cgroup v1/v2:   检测到\n")
			fmt.Printf("CPU 配额:       %d μs\n", quota)
			fmt.Printf("配额周期:       %d μs\n", period)
			fmt.Printf("CPU 限制:       %.2f 核\n", cpuLimit)
			fmt.Println()

			// 判断容器感知调度是否生效
			fmt.Println("--- 容器感知调度状态 ---")
			if gomaxprocs < numCPU {
				fmt.Println("✅ 容器感知调度已生效")
				fmt.Printf("   GOMAXPROCS (%d) < NumCPU (%d)\n", gomaxprocs, numCPU)
				fmt.Printf("   应用已根据容器 CPU 限制 (%.2f 核) 自动调整\n", cpuLimit)
			} else if gomaxprocs == numCPU {
				fmt.Println("ℹ️  可能未在容器中运行")
				fmt.Printf("   GOMAXPROCS (%d) == NumCPU (%d)\n", gomaxprocs, numCPU)
			} else {
				fmt.Println("⚠️  异常情况")
				fmt.Printf("   GOMAXPROCS (%d) > NumCPU (%d)\n", gomaxprocs, numCPU)
			}
		} else {
			fmt.Println("未检测到 cgroup CPU 限制")
			fmt.Println("可能运行在物理机或未限制的容器中")
		}
	} else {
		fmt.Printf("cgroup 检测仅支持 Linux 系统 (当前: %s)\n", runtime.GOOS)
	}
	fmt.Println()

	// 环境变量检查
	fmt.Println("--- 环境变量 ---")
	if env := os.Getenv("GOMAXPROCS"); env != "" {
		fmt.Printf("GOMAXPROCS env: %s (已手动设置，覆盖自动检测)\n", env)
	} else {
		fmt.Println("GOMAXPROCS env: 未设置 (使用自动检测)")
	}
	fmt.Println()

	// Goroutine 信息
	fmt.Println("--- Goroutine 信息 ---")
	fmt.Printf("当前 Goroutine 数: %d\n", runtime.NumGoroutine())
	fmt.Println()

	// 内存信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("--- 内存信息 ---")
	fmt.Printf("已分配内存:     %d MB\n", m.Alloc/(1024*1024))
	fmt.Printf("系统内存:       %d MB\n", m.Sys/(1024*1024))
	fmt.Printf("堆内存:         %d MB\n", m.HeapAlloc/(1024*1024))
	fmt.Println()

	// 建议
	fmt.Println("--- 建议 ---")
	if runtime.GOOS == "linux" {
		if quota, period := readCgroupCPU(); quota > 0 && period > 0 {
			cpuLimit := float64(quota) / float64(period)
			recommendedGOMAXPROCS := int(cpuLimit)
			if recommendedGOMAXPROCS < 1 {
				recommendedGOMAXPROCS = 1
			}

			if gomaxprocs == recommendedGOMAXPROCS || gomaxprocs == recommendedGOMAXPROCS+1 || gomaxprocs == recommendedGOMAXPROCS-1 {
				fmt.Println("✅ GOMAXPROCS 设置合理，无需调整")
			} else {
				fmt.Printf("⚠️  建议 GOMAXPROCS 设置为 %d（当前 %d）\n", recommendedGOMAXPROCS, gomaxprocs)
				fmt.Printf("   可通过环境变量设置: export GOMAXPROCS=%d\n", recommendedGOMAXPROCS)
			}
		}
	} else {
		fmt.Println("非 Linux 系统，建议手动设置 GOMAXPROCS")
	}

	fmt.Println("\n===================================")
}

// readCgroupCPU 读取 cgroup CPU 配额
func readCgroupCPU() (quota, period int64) {
	// 尝试 cgroup v1
	if q, p := readCgroupV1CPU(); q > 0 && p > 0 {
		return q, p
	}

	// 尝试 cgroup v2
	if q, p := readCgroupV2CPU(); q > 0 && p > 0 {
		return q, p
	}

	return 0, 0
}

// readCgroupV1CPU 读取 cgroup v1 CPU 配额
func readCgroupV1CPU() (quota, period int64) {
	// 读取 cpu.cfs_quota_us
	quotaBytes, err := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	if err != nil {
		return 0, 0
	}
	quota, _ = strconv.ParseInt(strings.TrimSpace(string(quotaBytes)), 10, 64)

	// 读取 cpu.cfs_period_us
	periodBytes, err := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
	if err != nil {
		return 0, 0
	}
	period, _ = strconv.ParseInt(strings.TrimSpace(string(periodBytes)), 10, 64)

	return quota, period
}

// readCgroupV2CPU 读取 cgroup v2 CPU 配额
func readCgroupV2CPU() (quota, period int64) {
	// 读取 cpu.max
	maxBytes, err := os.ReadFile("/sys/fs/cgroup/cpu.max")
	if err != nil {
		return 0, 0
	}

	// 格式: "quota period" (例如: "50000 100000")
	parts := strings.Fields(strings.TrimSpace(string(maxBytes)))
	if len(parts) != 2 {
		return 0, 0
	}

	quota, _ = strconv.ParseInt(parts[0], 10, 64)
	period, _ = strconv.ParseInt(parts[1], 10, 64)

	return quota, period
}
