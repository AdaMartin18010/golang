package memory_allocator

import (
	"fmt"
	"runtime"
	"testing"
)

// ===== Swiss Tables Map 基准测试 =====

// BenchmarkMapLarge 测试大规模 map 性能
func BenchmarkMapLarge(b *testing.B) {
	sizes := []int{1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			m := make(map[int]int, size)
			for i := 0; i < size; i++ {
				m[i] = i
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = m[i%size]
			}
		})
	}
}

// BenchmarkMapInsert 测试 map 插入性能
func BenchmarkMapInsert(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				m := make(map[int]int, size)
				for j := 0; j < size; j++ {
					m[j] = j
				}
				runtime.KeepAlive(m)
			}
		})
	}
}

// BenchmarkMapIteration 测试 map 遍历性能
func BenchmarkMapIteration(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		m := make(map[int]int, size)
		for i := 0; i < size; i++ {
			m[i] = i
		}

		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				sum := 0
				for _, v := range m {
					sum += v
				}
				runtime.KeepAlive(sum)
			}
		})
	}
}

// BenchmarkMapStringKey 测试字符串键 map 性能
func BenchmarkMapStringKey(b *testing.B) {
	const size = 10000
	m := make(map[string]int, size)
	keys := make([]string, size)

	for i := 0; i < size; i++ {
		key := fmt.Sprintf("key_%d", i)
		keys[i] = key
		m[key] = i
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = m[keys[i%size]]
	}
}

// ===== 小对象分配基准测试 =====

type SmallObject struct {
	ID   int64
	Data [64]byte
}

// BenchmarkSmallObjectAllocation 测试小对象分配
func BenchmarkSmallObjectAllocation(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			obj := &SmallObject{ID: int64(i)}
			runtime.KeepAlive(obj)
		}
	})

	b.Run("Batch", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			objects := make([]*SmallObject, 100)
			for j := range objects {
				objects[j] = &SmallObject{ID: int64(j)}
			}
			runtime.KeepAlive(objects)
		}
	})
}

// BenchmarkSliceAllocation 测试切片分配
func BenchmarkSliceAllocation(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				slice := make([]int, size)
				runtime.KeepAlive(slice)
			}
		})
	}
}

// ===== 内存分配模式对比 =====

// BenchmarkAllocationPatterns 对比不同分配模式
func BenchmarkAllocationPatterns(b *testing.B) {
	b.Run("NewEachTime", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			for j := 0; j < 100; j++ {
				obj := &SmallObject{ID: int64(j)}
				runtime.KeepAlive(obj)
			}
		}
	})

	b.Run("ReuseSlice", func(b *testing.B) {
		objects := make([]*SmallObject, 100)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range objects {
				if objects[j] == nil {
					objects[j] = &SmallObject{}
				}
				objects[j].ID = int64(j)
			}
			runtime.KeepAlive(objects)
		}
	})

	b.Run("PreAllocate", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			objects := make([]SmallObject, 100)
			for j := range objects {
				objects[j].ID = int64(j)
			}
			runtime.KeepAlive(objects)
		}
	})
}

// ===== GC 压力测试 =====

// BenchmarkGCPressure 测试 GC 压力
func BenchmarkGCPressure(b *testing.B) {
	b.Run("HighAllocation", func(b *testing.B) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		startGC := memStats.NumGC

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// 大量临时对象
			for j := 0; j < 1000; j++ {
				obj := &SmallObject{ID: int64(j)}
				runtime.KeepAlive(obj)
			}
		}

		b.StopTimer()
		runtime.ReadMemStats(&memStats)
		b.ReportMetric(float64(memStats.NumGC-startGC), "gc-count")
	})

	b.Run("LowAllocation", func(b *testing.B) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		startGC := memStats.NumGC

		// 预分配，减少GC压力
		pool := make([]SmallObject, 1000)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := range pool {
				pool[j].ID = int64(j)
			}
			runtime.KeepAlive(pool)
		}

		b.StopTimer()
		runtime.ReadMemStats(&memStats)
		b.ReportMetric(float64(memStats.NumGC-startGC), "gc-count")
	})
}

// ===== 内存占用测试 =====

// BenchmarkMemoryUsage 测试内存占用
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("MapVsSlice", func(b *testing.B) {
		const size = 10000

		b.Run("Map", func(b *testing.B) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			before := m.HeapAlloc

			data := make(map[int]int, size)
			for i := 0; i < size; i++ {
				data[i] = i
			}

			runtime.ReadMemStats(&m)
			after := m.HeapAlloc

			b.ReportMetric(float64(after-before)/(1024*1024), "MB")
			runtime.KeepAlive(data)
		})

		b.Run("Slice", func(b *testing.B) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			before := m.HeapAlloc

			data := make([]int, size)
			for i := range data {
				data[i] = i
			}

			runtime.ReadMemStats(&m)
			after := m.HeapAlloc

			b.ReportMetric(float64(after-before)/(1024*1024), "MB")
			runtime.KeepAlive(data)
		})
	})
}

// ===== 实际场景模拟 =====

type User struct {
	ID       int64
	Name     string
	Email    string
	Metadata map[string]interface{}
}

// BenchmarkRealWorldScenario 模拟实际应用场景
func BenchmarkRealWorldScenario(b *testing.B) {
	b.Run("CacheSimulation", func(b *testing.B) {
		cache := make(map[int64]*User, 10000)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// 模拟缓存操作
			userID := int64(i % 10000)

			// 查找
			if _, ok := cache[userID]; !ok {
				// 未命中，创建新用户
				cache[userID] = &User{
					ID:       userID,
					Name:     fmt.Sprintf("User%d", userID),
					Email:    fmt.Sprintf("user%d@example.com", userID),
					Metadata: make(map[string]interface{}),
				}
			}

			// 访问
			_ = cache[userID]
		}
	})

	b.Run("DataProcessing", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			// 模拟数据处理：读取-处理-聚合
			data := make(map[string]int, 1000)

			// 读取和处理
			for j := 0; j < 1000; j++ {
				key := fmt.Sprintf("metric_%d", j%100)
				data[key]++
			}

			// 聚合
			sum := 0
			for _, v := range data {
				sum += v
			}

			runtime.KeepAlive(sum)
		}
	})
}

// ===== 对比测试助手 =====

// CompareVersions 对比 Go 1.24 和 1.25 性能
// 使用方法:
//
//	# Go 1.24
//	go test -bench=. -benchmem > go1.24.txt
//
//	# Go 1.25
//	go test -bench=. -benchmem > go1.25.txt
//
//	# 对比
//	benchstat go1.24.txt go1.25.txt
func CompareVersions() {
	// 此函数为说明性注释
}

// ===== 功能测试 =====

// TestMapBasic 测试 map 基本功能
func TestMapBasic(t *testing.T) {
	m := make(map[int]int, 1000)

	// 插入
	for i := 0; i < 1000; i++ {
		m[i] = i * 2
	}

	// 查找
	if v, ok := m[500]; !ok || v != 1000 {
		t.Errorf("期望 m[500]=1000, 实际 %d", v)
	}

	// 删除
	delete(m, 500)
	if _, ok := m[500]; ok {
		t.Error("删除失败")
	}

	// 大小
	if len(m) != 999 {
		t.Errorf("期望大小 999, 实际 %d", len(m))
	}
}

// TestMemoryStats 测试内存统计
func TestMemoryStats(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	t.Logf("内存统计:")
	t.Logf("  Alloc: %d MB", m.Alloc/(1024*1024))
	t.Logf("  TotalAlloc: %d MB", m.TotalAlloc/(1024*1024))
	t.Logf("  Sys: %d MB", m.Sys/(1024*1024))
	t.Logf("  NumGC: %d", m.NumGC)
	t.Logf("  HeapObjects: %d", m.HeapObjects)
}
