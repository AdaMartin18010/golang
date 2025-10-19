// Mock implementation of ASan example - no CGO required

package main

import (
	"fmt"
	"runtime"
	"time"
)

// MockMemoryBlock 模拟C内存块
type MockMemoryBlock struct {
	size      int
	freed     bool
	data      []byte
	allocTime time.Time
}

var (
	mockAllocations         = make(map[uintptr]*MockMemoryBlock)
	nextMockPtr     uintptr = 0x1000
)

// mockMalloc 模拟C的malloc
func mockMalloc(size int) uintptr {
	nextMockPtr++
	ptr := nextMockPtr

	mockAllocations[ptr] = &MockMemoryBlock{
		size:      size,
		freed:     false,
		data:      make([]byte, size),
		allocTime: time.Now(),
	}

	return ptr
}

// mockFree 模拟C的free
func mockFree(ptr uintptr) {
	if block, ok := mockAllocations[ptr]; ok {
		if block.freed {
			fmt.Printf("❌ [ASan模拟] double-free detected: %#x\n", ptr)
			return
		}
		block.freed = true
		delete(mockAllocations, ptr)
	} else {
		fmt.Printf("❌ [ASan模拟] invalid free: %#x\n", ptr)
	}
}

// demonstrateMemoryLeak 演示内存泄漏
func demonstrateMemoryLeak() {
	fmt.Println("\n=== 1. 内存泄漏检测 (模拟) ===")

	// 分配内存但不释放
	ptr := mockMalloc(1024)
	fmt.Printf("✅ [ASan模拟] 分配内存: %#x (1024 字节)\n", ptr)

	// 故意不调用 mockFree(ptr)
	fmt.Println("⚠️  [ASan模拟] 内存未释放 - 将被检测为泄漏!")

	// 在真实ASan中，这会被检测为内存泄漏
}

// demonstrateBufferOverflow 演示缓冲区溢出
func demonstrateBufferOverflow() {
	fmt.Println("\n=== 2. 缓冲区溢出检测 (模拟) ===")

	ptr := mockMalloc(10)
	block := mockAllocations[ptr]

	fmt.Printf("✅ [ASan模拟] 分配内存: %#x (10 字节)\n", ptr)

	// 正常访问
	for i := 0; i < 10; i++ {
		block.data[i] = byte(i)
	}
	fmt.Println("✅ [ASan模拟] 正常访问: block[0-9] = OK")

	// 越界访问
	fmt.Println("⚠️  [ASan模拟] 尝试越界访问: block[10]")
	if len(block.data) <= 10 {
		fmt.Println("❌ [ASan模拟] buffer-overflow detected: 超出分配大小!")
	}

	mockFree(ptr)
	fmt.Printf("✅ [ASan模拟] 释放内存: %#x\n", ptr)
}

// demonstrateUseAfterFree 演示Use-After-Free
func demonstrateUseAfterFree() {
	fmt.Println("\n=== 3. Use-After-Free检测 (模拟) ===")

	ptr := mockMalloc(1024)
	fmt.Printf("✅ [ASan模拟] 分配内存: %#x (1024 字节)\n", ptr)

	mockFree(ptr)
	fmt.Printf("✅ [ASan模拟] 释放内存: %#x\n", ptr)

	// 尝试使用已释放的内存
	fmt.Println("⚠️  [ASan模拟] 尝试访问已释放的内存...")
	if block, ok := mockAllocations[ptr]; !ok || block.freed {
		fmt.Println("❌ [ASan模拟] use-after-free detected: 访问已释放的内存!")
	}
}

// demonstrateDoubleFree 演示Double Free
func demonstrateDoubleFree() {
	fmt.Println("\n=== 4. Double Free检测 (模拟) ===")

	ptr := mockMalloc(512)
	fmt.Printf("✅ [ASan模拟] 分配内存: %#x (512 字节)\n", ptr)

	mockFree(ptr)
	fmt.Printf("✅ [ASan模拟] 第一次释放: %#x\n", ptr)

	// 尝试再次释放
	fmt.Println("⚠️  [ASan模拟] 尝试第二次释放...")
	mockFree(ptr)
}

// checkMemoryLeaks 检查内存泄漏
func checkMemoryLeaks() {
	fmt.Println("\n=== 内存泄漏报告 ===")

	if len(mockAllocations) == 0 {
		fmt.Println("✅ [ASan模拟] 没有检测到内存泄漏")
		return
	}

	fmt.Printf("❌ [ASan模拟] 检测到 %d 个内存泄漏:\n", len(mockAllocations))
	for ptr, block := range mockAllocations {
		if !block.freed {
			fmt.Printf("  - 地址: %#x, 大小: %d 字节, 分配时间: %s\n",
				ptr, block.size, block.allocTime.Format("15:04:05"))
		}
	}
}

// printMemoryStats 打印内存统计
func printMemoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("\n=== Go运行时内存统计 ===")
	fmt.Printf("分配的堆内存: %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("总分配内存:   %.2f MB\n", float64(m.TotalAlloc)/1024/1024)
	fmt.Printf("系统内存:     %.2f MB\n", float64(m.Sys)/1024/1024)
	fmt.Printf("GC次数:       %d\n", m.NumGC)
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║  Go 1.23+ AddressSanitizer (ASan) 模拟演示       ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")

	fmt.Println("\n⚠️  注意: 这是一个模拟版本，用于演示ASan的概念")
	fmt.Println("   真实的ASan需要CGO和C编译器支持")
	fmt.Println("   安装C编译器后，可以使用 'go build -asan' 启用真实ASan")

	// 演示各种内存安全问题
	demonstrateMemoryLeak()
	demonstrateBufferOverflow()
	demonstrateUseAfterFree()
	demonstrateDoubleFree()

	// 检查泄漏
	runtime.GC() // 触发垃圾回收
	checkMemoryLeaks()

	// 打印内存统计
	printMemoryStats()

	fmt.Println("\n╔══════════════════════════════════════════════════╗")
	fmt.Println("║  ASan最佳实践建议                                ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println("1. 在CI/CD中集成ASan测试")
	fmt.Println("2. 定期运行ASan扫描")
	fmt.Println("3. 立即修复发现的问题")
	fmt.Println("4. 结合race detector使用")
	fmt.Println("5. 仅在开发/测试环境启用（性能开销大）")

	fmt.Println("\n提示: 要使用真实ASan，请参考README.md安装C编译器")
}
