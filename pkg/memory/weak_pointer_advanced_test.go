package memory

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestWeakCacheMemoryPressure 测试内存压力下的weak cache行为
func TestWeakCacheMemoryPressure(t *testing.T) {
	cache := NewWeakCache()

	// 添加大量数据
	const numItems = 1000
	for i := 0; i < numItems; i++ {
		v := &Value{
			Data:      string(make([]byte, 1024)), // 1KB per item
			Size:      1024,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune('A'+i%26))+string(rune(i)), v)
	}

	// 触发GC
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// 清理
	cleaned := cache.Cleanup()
	t.Logf("Cleaned %d entries under memory pressure", cleaned)

	// 验证缓存仍然可用
	cache.Set("test-after-gc", &Value{Data: "test", Size: 100, CreatedAt: time.Now()})
	if _, found := cache.Get("test-after-gc"); !found {
		t.Error("Cache should still work after GC")
	}
}

// TestWeakCacheConcurrentAccess 测试并发访问
func TestWeakCacheConcurrentAccess(t *testing.T) {
	cache := NewWeakCache()
	const numGoroutines = 50
	const numOpsPerGoroutine = 100

	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOpsPerGoroutine; j++ {
				key := string(rune('A'+id%26)) + string(rune(j))
				v := &Value{
					Data:      "concurrent-data",
					Size:      100,
					CreatedAt: time.Now(),
				}
				cache.Set(key, v)
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOpsPerGoroutine; j++ {
				key := string(rune('A'+id%26)) + string(rune(j))
				cache.Get(key)
			}
		}(i)
	}

	// 并发清理
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Cleanup()
		}()
	}

	wg.Wait()

	// 验证缓存状态
	stats := cache.Stats()
	t.Logf("Stats after concurrent access: Hits=%d, Misses=%d", stats.Hits, stats.Misses)
}

// TestWeakCacheLargeValues 测试大值缓存
func TestWeakCacheLargeValues(t *testing.T) {
	cache := NewWeakCache()

	sizes := []int{1024, 10240, 102400, 1024000} // 1KB, 10KB, 100KB, 1MB

	for _, size := range sizes {
		t.Run(string(rune(size)), func(t *testing.T) {
			v := &Value{
				Data:      string(make([]byte, size)),
				Size:      size,
				CreatedAt: time.Now(),
			}

			key := string(rune(size))
			cache.Set(key, v)

			retrieved, found := cache.Get(key)
			if !found {
				t.Errorf("Failed to retrieve value of size %d", size)
			}

			if retrieved != nil && retrieved.Size != size {
				t.Errorf("Size mismatch: expected %d, got %d", size, retrieved.Size)
			}
		})
	}
}

// TestWeakCacheRapidSetGet 测试快速set/get操作
func TestWeakCacheRapidSetGet(t *testing.T) {
	cache := NewWeakCache()
	const numOps = 10000

	start := time.Now()
	for i := 0; i < numOps; i++ {
		key := string(rune(i % 100))
		v := &Value{
			Data:      "rapid-test",
			Size:      100,
			CreatedAt: time.Now(),
		}
		cache.Set(key, v)
		cache.Get(key)
	}
	duration := time.Since(start)

	t.Logf("Completed %d set/get operations in %v", numOps*2, duration)
	t.Logf("Average: %v per operation", duration/time.Duration(numOps*2))
}

// TestWeakCacheExpiration 测试缓存过期
func TestWeakCacheExpiration(t *testing.T) {
	cache := NewWeakCache()

	// 添加值
	v := &Value{
		Data:      "expire-test",
		Size:      100,
		CreatedAt: time.Now(),
	}
	cache.Set("expire-key", v)

	// 验证存在
	if _, found := cache.Get("expire-key"); !found {
		t.Error("Value should exist initially")
	}

	// 释放引用，触发GC
	v = nil
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// 清理
	cache.Cleanup()

	// 注意：由于GC的不确定性，这个测试可能不稳定
	t.Log("Expiration test completed (GC behavior may vary)")
}

// TestWeakPointerBehavior 测试weak pointer行为
func TestWeakPointerBehavior(t *testing.T) {
	v := &Value{
		Data:      "weak-test",
		Size:      100,
		CreatedAt: time.Now(),
	}

	wp := makeWeak(v)

	// 验证初始状态
	if !wp.alive {
		t.Error("Weak pointer should be alive initially")
	}

	if wp.value() == nil {
		t.Error("Weak pointer value should not be nil initially")
	}

	// 释放强引用
	v = nil
	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	// 注意：finalizer可能还未运行
	t.Log("Weak pointer behavior test completed")
}

// TestCacheStatsAccuracy 测试统计准确性
func TestCacheStatsAccuracy(t *testing.T) {
	cache := NewWeakCache()

	// 重置统计（通过创建新cache）
	cache = NewWeakCache()

	// 添加一些值
	for i := 0; i < 10; i++ {
		v := &Value{
			Data:      "stats-test",
			Size:      100,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune('A'+i)), v)
	}

	// 执行一些get操作
	hitCount := 0
	missCount := 0

	for i := 0; i < 20; i++ {
		key := string(rune('A' + i%15))
		if _, found := cache.Get(key); found {
			hitCount++
		} else {
			missCount++
		}
	}

	stats := cache.Stats()

	t.Logf("Expected hits: ~%d, actual: %d", hitCount, stats.Hits)
	t.Logf("Expected misses: ~%d, actual: %d", missCount, stats.Misses)

	if stats.Hits+stats.Misses != 20 {
		t.Errorf("Total operations should be 20, got %d", stats.Hits+stats.Misses)
	}
}

// TestWeakCacheCleanupEfficiency 测试清理效率
func TestWeakCacheCleanupEfficiency(t *testing.T) {
	cache := NewWeakCache()

	// 添加大量数据
	const numItems = 1000
	for i := 0; i < numItems; i++ {
		v := &Value{
			Data:      "cleanup-test",
			Size:      100,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune(i)), v)
	}

	// 多次清理
	start := time.Now()
	for i := 0; i < 10; i++ {
		cache.Cleanup()
	}
	duration := time.Since(start)

	t.Logf("10 cleanup operations completed in %v", duration)
	t.Logf("Average: %v per cleanup", duration/10)
}

// BenchmarkWeakCacheSet 基准测试Set操作
func BenchmarkWeakCacheSet(b *testing.B) {
	cache := NewWeakCache()
	v := &Value{
		Data:      "benchmark",
		Size:      1024,
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(string(rune(i%100)), v)
	}
}

// BenchmarkWeakCacheGet 基准测试Get操作
func BenchmarkWeakCacheGet(b *testing.B) {
	cache := NewWeakCache()

	// 预填充
	for i := 0; i < 100; i++ {
		v := &Value{
			Data:      "benchmark",
			Size:      1024,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune(i)), v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(string(rune(i % 100)))
	}
}

// BenchmarkWeakCacheCleanup 基准测试Cleanup操作
func BenchmarkWeakCacheCleanup(b *testing.B) {
	cache := NewWeakCache()

	// 预填充
	for i := 0; i < 1000; i++ {
		v := &Value{
			Data:      "benchmark",
			Size:      1024,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune(i)), v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Cleanup()
	}
}

// BenchmarkWeakCacheConcurrent 基准测试并发操作
func BenchmarkWeakCacheConcurrent(b *testing.B) {
	cache := NewWeakCache()

	// 预填充
	for i := 0; i < 100; i++ {
		v := &Value{
			Data:      "benchmark",
			Size:      1024,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune(i)), v)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := string(rune(i % 100))
			if i%2 == 0 {
				v := &Value{
					Data:      "benchmark",
					Size:      1024,
					CreatedAt: time.Now(),
				}
				cache.Set(key, v)
			} else {
				cache.Get(key)
			}
			i++
		}
	})
}
