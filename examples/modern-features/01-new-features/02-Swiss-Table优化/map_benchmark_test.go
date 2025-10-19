package performance_analysis

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

/*
   ======================================================================
   Swiss Table vs. Old Map - 性能与内存基准测试
   ======================================================================

   🎯 测试目的:
   本文件旨在通过一系列基准测试，量化对比 Go 1.22+ 引入的 Swiss Table Map
   与旧版（Go 1.21 及更早）的 Chained Hashing Map 在性能和内存占用上的差异。

   ⚙️ 如何运行测试:
   1. 确保你安装了两个版本的 Go (例如，1.21.x 和 1.22.x 或更高版本)。
      可以使用 `gvm` 或 `go` 命令的多版本管理功能。
   2. 使用 Go 1.21.x 编译并运行测试:
      $ go1.21.x test -bench=. -benchmem -count=3 > old_map_results.txt
   3. 使用 Go 1.22.x (或更高) 编译并运行测试:
      $ go1.22.x test -bench=. -benchmem -count=3 > swiss_table_results.txt
   4. 对比 `old_map_results.txt` 和 `swiss_table_results.txt` 文件中的数据。
      - `ns/op`: 每次操作耗费的纳秒数 (越低越好)。
      - `B/op`: 每次操作分配的字节数 (越低越好)。
      - `allocs/op`: 每次操作的内存分配次数 (越低越好)。

   🔍 预期结果:
   - Swiss Table 在大多数读写和迭代场景下应该有更低的 `ns/op`。
   - Swiss Table 在所有场景下应该有显著更低的 `B/op` 和 `allocs/op`，
     尤其是在内存分析测试中。
*/

const mapSize = 1_000_000 // 定义测试用的map大小

var keys []int // 预先生成一组用于测试的键

// setupKeys 函数负责生成和打乱测试用的键，以确保测试的随机性。
func setupKeys() {
	if len(keys) == mapSize {
		return
	}
	keys = make([]int, mapSize)
	for i := 0; i < mapSize; i++ {
		keys[i] = i
	}
	rand.New(rand.NewSource(42)).Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
}

// --- 性能基准测试 ---

// BenchmarkWrite_Dense 测试密集写入性能。
func BenchmarkWrite_Dense(b *testing.B) {
	setupKeys()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[int]int, mapSize)
		for _, k := range keys {
			m[k] = k
		}
	}
}

// BenchmarkRead_Dense 测试密集读取性能。
func BenchmarkRead_Dense(b *testing.B) {
	setupKeys()
	m := make(map[int]int, mapSize)
	for _, k := range keys {
		m[k] = k
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 随机读取 mapSize 次
		for j := 0; j < mapSize; j++ {
			_ = m[keys[j]]
		}
	}
}

// BenchmarkRead_HitMiss_50_50 测试50%命中率的读取性能。
func BenchmarkRead_HitMiss_50_50(b *testing.B) {
	setupKeys()
	m := make(map[int]int, mapSize/2) // Map中只有一半的键
	for i := 0; i < mapSize/2; i++ {
		m[keys[i]] = keys[i]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range keys {
			_ = m[k] // 一半会命中，一半会错过
		}
	}
}

// BenchmarkWriteDelete_Sparse 测试稀疏写入和删除的性能。
func BenchmarkWriteDelete_Sparse(b *testing.B) {
	setupKeys()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[int]int)
		// 模拟一个动态增删的过程
		for j := 0; j < mapSize; j++ {
			m[keys[j]] = j
			if j%10 == 0 { // 每10次写入，删除一个元素
				delete(m, keys[j/10])
			}
		}
	}
}

// BenchmarkIteration 测试遍历性能。
func BenchmarkIteration(b *testing.B) {
	setupKeys()
	m := make(map[int]int, mapSize)
	for _, k := range keys {
		m[k] = k
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum int
		for k, v := range m {
			sum += k + v
		}
	}
}

// --- 内存占用分析 ---

// captureMemStats 是一个辅助函数，用于获取当前goroutine的内存分配统计信息。
func captureMemStats(b *testing.B) uint64 {
	b.StopTimer()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	b.StartTimer()
	return ms.Alloc
}

// TestMemoryProfile 分析并打印不同大小map的内存占用。
// 注意：这是一个测试（Test）而不是基准（Benchmark），因为它只运行一次并打印结果。
func TestMemoryProfile(t *testing.T) {
	sizes := []int{100, 1_000, 10_000, 100_000, 1_000_000}

	fmt.Printf("\n--- Memory Usage Profile ---\n")
	fmt.Printf("%-15s | %-15s\n", "Map Size", "Memory (Bytes)")
	fmt.Println("---------------------------------")

	for _, size := range sizes {
		runtime.GC() // 在每次测试前运行垃圾回收，以获得更准确的读数
		startMem := captureMemStats(&testing.B{})

		m := make(map[int]int, size)
		for i := 0; i < size; i++ {
			m[i] = i
		}

		endMem := captureMemStats(&testing.B{})

		// 打印map本身占用的内存，并让它在作用域结束后被回收
		fmt.Printf("%-15d | %-15d\n", size, endMem-startMem)
		runtime.KeepAlive(m)
	}
	fmt.Println("---------------------------------")
}
