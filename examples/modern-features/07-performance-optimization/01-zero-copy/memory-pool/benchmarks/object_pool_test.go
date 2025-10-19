package benchmarks

import (
	"sync"
	"testing"

	memorypool "performance-optimization-examples/01-zero-copy/memory-pool"
)

// 测试数据结构
type TestObject struct {
	ID   int
	Name string
	Data []byte
}

func (t *TestObject) Reset() {
	t.ID = 0
	t.Name = ""
	t.Data = t.Data[:0]
}

// BenchmarkObjectPool 对象池性能基准测试
func BenchmarkObjectPool(b *testing.B) {
	// 创建对象池
	pool := memorypool.NewObjectPool(1000, func() interface{} {
		return &TestObject{
			Data: make([]byte, 0, 1024),
		}
	}, func(obj interface{}) {
		if t, ok := obj.(*TestObject); ok {
			t.Reset()
		}
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 获取对象
			obj := pool.Get().(*TestObject)

			// 模拟使用对象
			obj.ID = 1
			obj.Name = "test"
			obj.Data = append(obj.Data, "hello"...)

			// 归还对象
			pool.Put(obj)
		}
	})
}

// BenchmarkDirectAllocation 直接分配性能基准测试
func BenchmarkDirectAllocation(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 直接创建对象
			obj := &TestObject{
				ID:   1,
				Name: "test",
				Data: make([]byte, 0, 1024),
			}

			// 模拟使用对象
			obj.Data = append(obj.Data, "hello"...)

			// 对象会被GC回收
		}
	})
}

// BenchmarkBufferPool 缓冲区池性能基准测试
func BenchmarkBufferPool(b *testing.B) {
	pool := memorypool.NewBufferPool(1000, 1024)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 获取缓冲区
			buf := pool.GetBuffer()

			// 模拟使用缓冲区
			buf = append(buf, "hello world"...)

			// 归还缓冲区
			pool.PutBuffer(buf)
		}
	})
}

// BenchmarkDirectBuffer 直接分配缓冲区性能基准测试
func BenchmarkDirectBuffer(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 直接创建缓冲区
			buf := make([]byte, 0, 1024)

			// 模拟使用缓冲区
			buf = append(buf, "hello world"...)

			// 缓冲区会被GC回收
		}
	})
}

// BenchmarkRequestPool HTTP请求池性能基准测试
func BenchmarkRequestPool(b *testing.B) {
	pool := memorypool.NewRequestPool(1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 获取请求对象
			req := pool.GetRequest()

			// 模拟使用请求对象
			req.Method = "GET"
			req.URL = "/api/users"
			req.Headers["Content-Type"] = "application/json"
			req.Body = append(req.Body, "{}"...)

			// 归还请求对象
			pool.PutRequest(req)
		}
	})
}

// BenchmarkDirectRequest 直接分配请求对象性能基准测试
func BenchmarkDirectRequest(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 直接创建请求对象
			req := &memorypool.HTTPRequest{
				Method:  "GET",
				URL:     "/api/users",
				Headers: make(map[string]string),
				Body:    make([]byte, 0),
			}

			// 模拟使用请求对象
			req.Headers["Content-Type"] = "application/json"
			req.Body = append(req.Body, "{}"...)

			// 请求对象会被GC回收
		}
	})
}

// BenchmarkConcurrentAccess 并发访问性能基准测试
func BenchmarkConcurrentAccess(b *testing.B) {
	pool := memorypool.NewObjectPool(100, func() interface{} {
		return &TestObject{
			Data: make([]byte, 0, 1024),
		}
	}, func(obj interface{}) {
		if t, ok := obj.(*TestObject); ok {
			t.Reset()
		}
	})

	var wg sync.WaitGroup
	concurrency := 100
	iterations := b.N / concurrency

	b.ResetTimer()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				obj := pool.Get().(*TestObject)
				obj.ID = j
				obj.Name = "test"
				obj.Data = append(obj.Data, "hello"...)
				pool.Put(obj)
			}
		}()
	}
	wg.Wait()
}

// BenchmarkPoolStats 池统计信息性能基准测试
func BenchmarkPoolStats(b *testing.B) {
	pool := memorypool.NewObjectPool(1000, func() interface{} {
		return &TestObject{
			Data: make([]byte, 0, 1024),
		}
	}, func(obj interface{}) {
		if t, ok := obj.(*TestObject); ok {
			t.Reset()
		}
	})

	// 预热池
	for i := 0; i < 100; i++ {
		obj := pool.Get()
		pool.Put(obj)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stats := pool.Stats()
		_ = stats.HitRate
		_ = stats.Created
		_ = stats.Reused
	}
}

// TestPoolStats 测试池统计信息
func TestPoolStats(t *testing.T) {
	pool := memorypool.NewObjectPool(10, func() interface{} {
		return &TestObject{
			Data: make([]byte, 0, 1024),
		}
	}, func(obj interface{}) {
		if t, ok := obj.(*TestObject); ok {
			t.Reset()
		}
	})

	// 获取一些对象
	obj1 := pool.Get()
	obj2 := pool.Get()
	_ = pool.Get() // obj3 - 获取但不归还，用于测试

	// 归还一些对象
	pool.Put(obj1)
	pool.Put(obj2)

	// 检查统计信息
	stats := pool.Stats()
	if stats.Created != 3 {
		t.Errorf("Expected 3 created objects, got %d", stats.Created)
	}
	if stats.Reused != 2 {
		t.Errorf("Expected 2 reused objects, got %d", stats.Reused)
	}
	if stats.PoolSize != 2 {
		t.Errorf("Expected pool size 2, got %d", stats.PoolSize)
	}
}

// TestPoolConcurrency 测试池的并发安全性
func TestPoolConcurrency(t *testing.T) {
	pool := memorypool.NewObjectPool(100, func() interface{} {
		return &TestObject{
			Data: make([]byte, 0, 1024),
		}
	}, func(obj interface{}) {
		if t, ok := obj.(*TestObject); ok {
			t.Reset()
		}
	})

	var wg sync.WaitGroup
	concurrency := 100
	iterations := 1000

	// 启动多个goroutine并发访问池
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				obj := pool.Get().(*TestObject)
				obj.ID = id
				obj.Name = "test"
				obj.Data = append(obj.Data, "hello"...)
				pool.Put(obj)
			}
		}(i)
	}

	wg.Wait()

	// 检查统计信息
	stats := pool.Stats()
	t.Logf("Pool stats: Created=%d, Reused=%d, HitRate=%.2f%%",
		stats.Created, stats.Reused, stats.HitRate)
}

// BenchmarkMemoryUsage 内存使用基准测试
func BenchmarkMemoryUsage(b *testing.B) {
	pool := memorypool.NewObjectPool(1000, func() interface{} {
		return &TestObject{
			Data: make([]byte, 0, 1024),
		}
	}, func(obj interface{}) {
		if t, ok := obj.(*TestObject); ok {
			t.Reset()
		}
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := pool.Get().(*TestObject)
			obj.Data = append(obj.Data, make([]byte, 512)...)
			pool.Put(obj)
		}
	})
}

// BenchmarkPoolVsDirect 池 vs 直接分配对比基准测试
func BenchmarkPoolVsDirect(b *testing.B) {
	b.Run("Pool", func(b *testing.B) {
		pool := memorypool.NewObjectPool(1000, func() interface{} {
			return &TestObject{
				Data: make([]byte, 0, 1024),
			}
		}, func(obj interface{}) {
			if t, ok := obj.(*TestObject); ok {
				t.Reset()
			}
		})

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				obj := pool.Get().(*TestObject)
				obj.ID = 1
				obj.Name = "test"
				obj.Data = append(obj.Data, "hello"...)
				pool.Put(obj)
			}
		})
	})

	b.Run("Direct", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				obj := &TestObject{
					ID:   1,
					Name: "test",
					Data: make([]byte, 0, 1024),
				}
				obj.Data = append(obj.Data, "hello"...)
			}
		})
	})
}
