package memory

import (
	"sync"
	"testing"
)

// =============================================================================
// 内存池测试
// =============================================================================

// TestGenericPool 测试通用对象池
func TestGenericPool(t *testing.T) {
	// 创建字符串池
	pool := NewGenericPool(
		func() interface{} {
			return new(string)
		},
		func(obj interface{}) {
			s := obj.(*string)
			*s = ""
		},
		100,
	)

	// 测试Get
	obj1 := pool.Get().(*string)
	if obj1 == nil {
		t.Fatal("Get returned nil")
	}

	// 测试Put
	*obj1 = "test"
	pool.Put(obj1)

	// 再次Get应该获取到相同对象
	obj2 := pool.Get().(*string)
	if *obj2 != "" {
		t.Errorf("Object not reset, got %s", *obj2)
	}

	// 测试统计
	stats := pool.Stats()
	if stats.Gets != 2 {
		t.Errorf("Expected 2 gets, got %d", stats.Gets)
	}
	if stats.Puts != 1 {
		t.Errorf("Expected 1 put, got %d", stats.Puts)
	}
}

// TestGenericPoolHitRate 测试命中率
func TestGenericPoolHitRate(t *testing.T) {
	pool := NewGenericPool(
		func() interface{} {
			return new(int)
		},
		nil,
		10,
	)

	// 先Put一些对象
	for i := 0; i < 5; i++ {
		pool.Put(new(int))
	}

	// 然后Get
	for i := 0; i < 10; i++ {
		pool.Get()
	}

	hitRate := pool.HitRate()
	if hitRate < 0 || hitRate > 100 {
		t.Errorf("Invalid hit rate: %.2f%%", hitRate)
	}
}

// TestGenericPoolMaxSize 测试最大大小限制
func TestGenericPoolMaxSize(t *testing.T) {
	maxSize := 5
	pool := NewGenericPool(
		func() interface{} {
			return new(int)
		},
		nil,
		maxSize,
	)

	// Put超过maxSize的对象
	for i := 0; i < 10; i++ {
		pool.Put(new(int))
	}

	stats := pool.Stats()
	if stats.Size > maxSize {
		t.Errorf("Pool size %d exceeds max size %d", stats.Size, maxSize)
	}
}

// TestGenericPoolClear 测试清空池
func TestGenericPoolClear(t *testing.T) {
	pool := NewGenericPool(
		func() interface{} {
			return new(int)
		},
		nil,
		10,
	)

	// Put一些对象
	for i := 0; i < 5; i++ {
		pool.Put(new(int))
	}

	// 清空
	pool.Clear()

	stats := pool.Stats()
	if stats.Size != 0 {
		t.Errorf("Pool not cleared, size: %d", stats.Size)
	}
}

// TestBytePool 测试字节池
func TestBytePool(t *testing.T) {
	pool := NewBytePool([]int{256, 1024, 4096})

	// 获取256字节
	buf1 := pool.Get(200)
	if len(*buf1) != 200 || cap(*buf1) != 256 {
		t.Errorf("Unexpected buffer size: len=%d, cap=%d", len(*buf1), cap(*buf1))
	}

	// 归还
	pool.Put(buf1)

	// 再次获取应该复用
	buf2 := pool.Get(256)
	if cap(*buf2) != 256 {
		t.Errorf("Buffer not reused, cap=%d", cap(*buf2))
	}
}

// TestBytePoolLargeSize 测试大尺寸分配
func TestBytePoolLargeSize(t *testing.T) {
	pool := NewBytePool([]int{256, 1024})

	// 请求超过最大池大小的缓冲区
	buf := pool.Get(5000)
	if len(*buf) != 5000 {
		t.Errorf("Expected buffer length 5000, got %d", len(*buf))
	}

	// 归还不应该panic
	pool.Put(buf)
}

// TestPoolManager 测试池管理器
func TestPoolManager(t *testing.T) {
	manager := NewPoolManager()

	// 注册池
	pool1 := NewGenericPool(func() interface{} { return new(int) }, nil, 10)
	pool2 := NewGenericPool(func() interface{} { return new(string) }, nil, 10)

	manager.Register("int_pool", pool1)
	manager.Register("string_pool", pool2)

	// 获取池
	if p, ok := manager.Get("int_pool"); !ok || p == nil {
		t.Error("Failed to get int_pool")
	}

	if p, ok := manager.Get("string_pool"); !ok || p == nil {
		t.Error("Failed to get string_pool")
	}

	// 获取统计
	stats := manager.AllStats()
	if len(stats) != 2 {
		t.Errorf("Expected 2 pools, got %d", len(stats))
	}
}

// TestDefaultBytePool 测试默认字节池
func TestDefaultBytePool(t *testing.T) {
	// 使用默认池
	buf := GetBytes(1024)
	if len(*buf) != 1024 {
		t.Errorf("Expected buffer length 1024, got %d", len(*buf))
	}

	PutBytes(buf)
}

// =============================================================================
// 并发测试
// =============================================================================

// TestGenericPoolConcurrent 并发测试
func TestGenericPoolConcurrent(t *testing.T) {
	pool := NewGenericPool(
		func() interface{} {
			return new(int)
		},
		func(obj interface{}) {
			i := obj.(*int)
			*i = 0
		},
		1000,
	)

	var wg sync.WaitGroup
	concurrency := 10
	iterations := 1000

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				obj := pool.Get().(*int)
				*obj = j
				pool.Put(obj)
			}
		}()
	}

	wg.Wait()

	stats := pool.Stats()
	t.Logf("Concurrent stats: Gets=%d, Puts=%d, HitRate=%.2f%%",
		stats.Gets, stats.Puts, pool.HitRate())
}

// TestBytePoolConcurrent 字节池并发测试
func TestBytePoolConcurrent(t *testing.T) {
	pool := NewBytePool([]int{256, 1024, 4096})

	var wg sync.WaitGroup
	concurrency := 10
	iterations := 1000

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				size := 256 + (j%3)*768 // 256, 1024, 4096
				buf := pool.Get(size)
				// 使用缓冲区
				(*buf)[0] = byte(j)
				pool.Put(buf)
			}
		}()
	}

	wg.Wait()
}

// =============================================================================
// 基准测试
// =============================================================================

// BenchmarkGenericPool 对象池基准测试
func BenchmarkGenericPool(b *testing.B) {
	pool := NewGenericPool(
		func() interface{} {
			return make([]byte, 1024)
		},
		func(obj interface{}) {
			// Reset
		},
		100,
	)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := pool.Get()
			pool.Put(obj)
		}
	})
}

// BenchmarkBytePool 字节池基准测试
func BenchmarkBytePool(b *testing.B) {
	pool := NewBytePool([]int{256, 1024, 4096, 65536})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := pool.Get(1024)
			pool.Put(buf)
		}
	})
}

// BenchmarkDefaultBytePool 默认字节池基准测试
func BenchmarkDefaultBytePool(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetBytes(1024)
			PutBytes(buf)
		}
	})
}

// BenchmarkDirectAlloc 直接分配基准测试（对比）
func BenchmarkDirectAlloc(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := make([]byte, 1024)
			_ = buf
		}
	})
}
