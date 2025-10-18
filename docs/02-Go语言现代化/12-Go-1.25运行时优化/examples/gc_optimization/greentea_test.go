package gc_optimization

import (
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

// SmallObject 代表小对象（72字节）
type SmallObject struct {
	ID   int64
	Data [64]byte
}

// MediumObject 代表中等对象（1KB）
type MediumObject struct {
	ID    int64
	Data  [1016]byte
	Valid bool
}

// ===== 基准测试：小对象分配 =====

// BenchmarkSmallObjectAllocation 测试小对象密集分配场景
func BenchmarkSmallObjectAllocation(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		runtime.GC()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			objects := make([]*SmallObject, 10000)
			for j := range objects {
				objects[j] = &SmallObject{ID: int64(j)}
			}
			runtime.KeepAlive(objects)
		}
	})

	b.Run("Concurrent", func(b *testing.B) {
		runtime.GC()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				objects := make([]*SmallObject, 1000)
				for j := range objects {
					objects[j] = &SmallObject{ID: int64(j)}
				}
				runtime.KeepAlive(objects)
			}
		})
	})
}

// BenchmarkMixedSizeAllocation 测试混合大小对象分配
func BenchmarkMixedSizeAllocation(b *testing.B) {
	runtime.GC()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 90% 小对象，10% 中等对象
		for j := 0; j < 900; j++ {
			obj := &SmallObject{ID: int64(j)}
			runtime.KeepAlive(obj)
		}
		for j := 0; j < 100; j++ {
			obj := &MediumObject{ID: int64(j), Valid: true}
			runtime.KeepAlive(obj)
		}
	}
}

// ===== 基准测试：GC 性能 =====

// BenchmarkGCPause 测试 GC 暂停时间
func BenchmarkGCPause(b *testing.B) {
	// 预分配一些内存
	data := make([]*SmallObject, 100000)
	for i := range data {
		data[i] = &SmallObject{ID: int64(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		runtime.GC()
		duration := time.Since(start)

		b.ReportMetric(float64(duration.Microseconds()), "μs/gc")
	}
}

// BenchmarkGCOverhead 测试 GC 开销占比
func BenchmarkGCOverhead(b *testing.B) {
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	startGC := stats.NumGC
	startTime := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟实际工作负载
		objects := make([]*SmallObject, 5000)
		for j := range objects {
			objects[j] = &SmallObject{ID: int64(j)}
		}

		// 一些计算工作
		sum := int64(0)
		for _, obj := range objects {
			sum += obj.ID
		}
		runtime.KeepAlive(sum)
		runtime.KeepAlive(objects)
	}

	totalTime := time.Since(startTime)
	debug.ReadGCStats(&stats)
	gcCount := stats.NumGC - startGC

	b.ReportMetric(float64(gcCount), "gc-cycles")
	b.ReportMetric(float64(stats.PauseTotal)/float64(totalTime), "gc-overhead")
}

// ===== 基准测试：高并发场景 =====

// BenchmarkHighConcurrency 测试高并发分配场景
func BenchmarkHighConcurrency(b *testing.B) {
	numWorkers := runtime.GOMAXPROCS(0)

	b.Run("Workers-"+string(rune(numWorkers)), func(b *testing.B) {
		var wg sync.WaitGroup
		work := make(chan int, b.N)

		// 启动工作器
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range work {
					objects := make([]*SmallObject, 1000)
					for j := range objects {
						objects[j] = &SmallObject{ID: int64(j)}
					}
					runtime.KeepAlive(objects)
				}
			}()
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			work <- i
		}
		close(work)
		wg.Wait()
	})
}

// ===== 基准测试：长时间运行 =====

// BenchmarkLongRunning 测试长时间运行场景
func BenchmarkLongRunning(b *testing.B) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	startAlloc := memStats.TotalAlloc
	startGC := memStats.NumGC

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟持续的小对象分配
		for j := 0; j < 1000; j++ {
			obj := &SmallObject{ID: int64(j)}
			// 一些简单计算
			obj.Data[0] = byte(j % 256)
			runtime.KeepAlive(obj)
		}
	}

	runtime.ReadMemStats(&memStats)
	totalAlloc := memStats.TotalAlloc - startAlloc
	totalGC := memStats.NumGC - startGC

	b.ReportMetric(float64(totalAlloc)/float64(b.N), "bytes/op")
	b.ReportMetric(float64(totalGC), "gc-count")
	b.ReportMetric(float64(memStats.HeapAlloc)/(1024*1024), "heap-MB")
}

// ===== 基准测试：内存碎片 =====

// BenchmarkMemoryFragmentation 测试内存碎片情况
func BenchmarkMemoryFragmentation(b *testing.B) {
	var memStats runtime.MemStats

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 混合分配和释放，模拟碎片产生
		objects := make([]*SmallObject, 10000)
		for j := range objects {
			objects[j] = &SmallObject{ID: int64(j)}
		}

		// 释放一半对象（奇数索引）
		for j := 1; j < len(objects); j += 2 {
			objects[j] = nil
		}

		runtime.GC()
		runtime.ReadMemStats(&memStats)
	}

	b.StopTimer()

	// 计算碎片率
	runtime.ReadMemStats(&memStats)
	heapInUse := memStats.HeapInuse
	heapAlloc := memStats.HeapAlloc
	fragmentation := float64(heapInUse-heapAlloc) / float64(heapInUse) * 100

	b.ReportMetric(fragmentation, "fragmentation-%")
}

// ===== 功能测试：GC 统计 =====

// TestGCStats 测试 GC 统计信息收集
func TestGCStats(t *testing.T) {
	var stats debug.GCStats
	debug.ReadGCStats(&stats)

	t.Logf("GC 统计信息:")
	t.Logf("  NumGC: %d", stats.NumGC)
	t.Logf("  PauseTotal: %v", stats.PauseTotal)

	if len(stats.Pause) > 0 {
		t.Logf("  LastGC: %v", stats.LastGC)
		t.Logf("  LastPause: %v", stats.Pause[0])

		// 计算平均暂停时间
		var totalPause time.Duration
		for _, p := range stats.Pause {
			totalPause += p
		}
		avgPause := totalPause / time.Duration(len(stats.Pause))
		t.Logf("  AvgPause: %v", avgPause)
	}
}

// TestMemoryStats 测试内存统计信息
func TestMemoryStats(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	t.Logf("内存统计信息:")
	t.Logf("  Alloc: %d MB", m.Alloc/(1024*1024))
	t.Logf("  TotalAlloc: %d MB", m.TotalAlloc/(1024*1024))
	t.Logf("  Sys: %d MB", m.Sys/(1024*1024))
	t.Logf("  HeapAlloc: %d MB", m.HeapAlloc/(1024*1024))
	t.Logf("  HeapInuse: %d MB", m.HeapInuse/(1024*1024))
	t.Logf("  HeapObjects: %d", m.HeapObjects)

	// 计算平均对象大小
	if m.Mallocs > 0 {
		avgObjSize := m.TotalAlloc / m.Mallocs
		t.Logf("  AvgObjectSize: %d bytes", avgObjSize)

		// 如果平均对象大小 < 256 bytes，greentea GC 应该有优势
		if avgObjSize < 256 {
			t.Logf("  ✅ 小对象为主，greentea GC 应该有优势")
		} else {
			t.Logf("  ⚠️ 对象较大，greentea GC 优势可能不明显")
		}
	}
}

// ===== 压力测试 =====

// TestStressGC 压力测试 GC 稳定性
func TestStressGC(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试")
	}

	const (
		duration   = 5 * time.Second
		goroutines = 100
	)

	var wg sync.WaitGroup
	done := make(chan struct{})

	// 启动多个 goroutine 持续分配
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					objects := make([]*SmallObject, 1000)
					for j := range objects {
						objects[j] = &SmallObject{ID: int64(j)}
					}
					runtime.KeepAlive(objects)
				}
			}
		}()
	}

	// 运行指定时间
	time.Sleep(duration)
	close(done)
	wg.Wait()

	// 收集 GC 统计
	var stats debug.GCStats
	debug.ReadGCStats(&stats)

	t.Logf("压力测试结果 (%v):", duration)
	t.Logf("  GC 次数: %d", stats.NumGC)
	t.Logf("  GC 总暂停: %v", stats.PauseTotal)
	t.Logf("  平均 GC 间隔: %v", duration/time.Duration(stats.NumGC))
}

// ===== 对比测试助手 =====

// CompareGCPerformance 对比默认 GC 和 greentea GC 性能
// 使用方法:
//
//	go test -bench=BenchmarkSmallObjectAllocation > default.txt
//	GOEXPERIMENT=greentea go test -bench=BenchmarkSmallObjectAllocation > greentea.txt
//	benchcmp default.txt greentea.txt
func CompareGCPerformance() {
	// 此函数为说明性注释，实际对比需要使用 benchcmp 或 benchstat 工具
}
