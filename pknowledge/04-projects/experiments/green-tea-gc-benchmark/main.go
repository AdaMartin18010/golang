// Green Tea GC 性能测试
// 个人实验项目

package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	fmt.Println("=== Green Tea GC Benchmark ===")
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Println()

	// 测试 1: 小对象分配
	benchmarkSmallAllocations()

	// 测试 2: GC 压力测试
	benchmarkGCPressure()

	// 测试 3: 内存占用
	measureMemoryUsage()
}

// 测试小对象分配
func benchmarkSmallAllocations() {
	fmt.Println("--- Small Object Allocations ---")

	sizes := []int{1, 64, 128, 512}
	for _, size := range sizes {
		duration := runAllocationBenchmark(size, 1000000)
		fmt.Printf("Size %4d bytes: %v\n", size, duration)
	}
	fmt.Println()
}

func runAllocationBenchmark(size, count int) time.Duration {
	start := time.Now()
	for i := 0; i < count; i++ {
		_ = make([]byte, size)
	}
	return time.Since(start)
}

// GC 压力测试
func benchmarkGCPressure() {
	fmt.Println("--- GC Pressure Test ---")

	var maxPause time.Duration
	var totalPause time.Duration

	for i := 0; i < 10; i++ {
		allocateMemory(100 * 1024 * 1024)

		start := time.Now()
		runtime.GC()
		pause := time.Since(start)

		if pause > maxPause {
			maxPause = pause
		}
		totalPause += pause
	}

	fmt.Printf("Max Pause: %v\n", maxPause)
	fmt.Printf("Avg Pause: %v\n", totalPause/10)
	fmt.Println()
}

func allocateMemory(size int) [][]byte {
	var data [][]byte
	for i := 0; i < size/1024; i++ {
		data = append(data, make([]byte, 1024))
	}
	return data
}

// 内存占用测量
func measureMemoryUsage() {
	fmt.Println("--- Memory Usage ---")

	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	objects := make([][]byte, 10000)
	for i := range objects {
		objects[i] = make([]byte, 1024)
	}

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	objects = nil
	runtime.GC()

	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)

	fmt.Printf("Heap Alloc (initial): %d KB\n", m1.HeapAlloc/1024)
	fmt.Printf("Heap Alloc (after alloc): %d KB\n", m2.HeapAlloc/1024)
	fmt.Printf("Heap Alloc (after free): %d KB\n", m3.HeapAlloc/1024)
	fmt.Println()
}

func init() {
	debug.SetGCPercent(100)
}
