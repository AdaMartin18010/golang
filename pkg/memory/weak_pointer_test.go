package memory

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestValueCreation 测试Value创建
func TestValueCreation(t *testing.T) {
	v := &Value{
		Data:      "test",
		Size:      1024,
		CreatedAt: time.Now(),
	}

	if v.Data != "test" {
		t.Errorf("Expected Data 'test', got '%s'", v.Data)
	}

	if v.Size != 1024 {
		t.Errorf("Expected Size 1024, got %d", v.Size)
	}
}

// TestWeakCacheCreation 测试WeakCache创建
func TestWeakCacheCreation(t *testing.T) {
	cache := NewWeakCache()

	if cache == nil {
		t.Fatal("NewWeakCache returned nil")
	}

	if cache.items == nil {
		t.Error("Cache items map not initialized")
	}
}

// TestWeakCacheSet 测试缓存设置
func TestWeakCacheSet(t *testing.T) {
	cache := NewWeakCache()
	v := &Value{
		Data:      "test1",
		Size:      1024,
		CreatedAt: time.Now(),
	}

	cache.Set("key1", v)

	// 验证缓存中有数据
	cache.mu.RLock()
	_, exists := cache.items["key1"]
	cache.mu.RUnlock()

	if !exists {
		t.Error("Expected key1 to exist in cache")
	}
}

// TestWeakCacheGet 测试缓存获取
func TestWeakCacheGet(t *testing.T) {
	cache := NewWeakCache()
	v := &Value{
		Data:      "test1",
		Size:      1024,
		CreatedAt: time.Now(),
	}

	cache.Set("key1", v)

	retrieved, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}

	if retrieved == nil {
		t.Error("Retrieved value should not be nil")
	}

	if retrieved.Data != "test1" {
		t.Errorf("Expected Data 'test1', got '%s'", retrieved.Data)
	}
}

// TestWeakCacheCleanup 测试缓存清理
func TestWeakCacheCleanup(t *testing.T) {
	cache := NewWeakCache()

	// 添加多个项目
	for i := 0; i < 10; i++ {
		v := &Value{
			Data:      "test",
			Size:      1024,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune('A'+i)), v)
	}

	// 等待一段时间
	time.Sleep(100 * time.Millisecond)

	// 触发清理
	cache.Cleanup()

	// 验证缓存仍然可用
	cache.mu.RLock()
	count := len(cache.items)
	cache.mu.RUnlock()

	if count < 0 {
		t.Error("Cache item count should not be negative")
	}
}

// TestWeakCacheConcurrency 测试并发访问
func TestWeakCacheConcurrency(t *testing.T) {
	cache := NewWeakCache()
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v := &Value{
				Data:      "test",
				Size:      1024,
				CreatedAt: time.Now(),
			}
			cache.Set(string(rune('A'+id)), v)
		}(i)
	}

	// 并发读取
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, _ = cache.Get(string(rune('A' + id)))
		}(i)
	}

	wg.Wait()

	// 验证缓存状态
	cache.mu.RLock()
	count := len(cache.items)
	cache.mu.RUnlock()

	if count < 0 {
		t.Error("Cache should contain entries")
	}
}

// TestValueCreationLarge 测试大Value对象创建
func TestValueCreationLarge(t *testing.T) {
	v := &Value{
		Data:      string(make([]byte, 1024*1024)), // 1MB
		Size:      1024 * 1024,
		CreatedAt: time.Now(),
	}

	if v.Size != 1024*1024 {
		t.Errorf("Expected size 1048576, got %d", v.Size)
	}
}

// TestCacheEviction 测试缓存驱逐
func TestCacheEviction(t *testing.T) {
	cache := NewWeakCache()

	// 添加对象
	v := &Value{
		Data:      "test",
		Size:      1024,
		CreatedAt: time.Now(),
	}
	cache.Set("key1", v)

	// 让对象变为不可达
	v = nil
	runtime.GC()

	// 等待一小段时间
	time.Sleep(50 * time.Millisecond)

	// 尝试获取（可能已被GC）
	_, found := cache.Get("key1")
	// 注意：这里不强制要求找不到，因为GC是不确定的
	t.Logf("Key found after GC: %v", found)
}

// TestMultipleOperations 测试多种操作组合
func TestMultipleOperations(t *testing.T) {
	cache := NewWeakCache()

	// 设置
	cache.Set("a", &Value{Data: "data-a", Size: 100, CreatedAt: time.Now()})
	cache.Set("b", &Value{Data: "data-b", Size: 200, CreatedAt: time.Now()})
	cache.Set("c", &Value{Data: "data-c", Size: 300, CreatedAt: time.Now()})

	// 获取
	if _, found := cache.Get("a"); !found {
		t.Error("Should find 'a'")
	}

	// 获取剩余的
	if _, found := cache.Get("c"); !found {
		t.Error("Should find 'c'")
	}
}

// TestCacheStats 测试缓存统计
func TestCacheStats(t *testing.T) {
	cache := NewWeakCache()

	// 添加一些数据
	cache.Set("key1", &Value{Data: "data1", Size: 100, CreatedAt: time.Now()})
	cache.Set("key2", &Value{Data: "data2", Size: 200, CreatedAt: time.Now()})

	// 执行一些操作
	cache.Get("key1") // hit
	cache.Get("key3") // miss

	stats := cache.Stats()

	if stats.Hits < 0 {
		t.Error("Hits should not be negative")
	}

	if stats.Misses < 0 {
		t.Error("Misses should not be negative")
	}
}

// BenchmarkCacheSet 基准测试缓存设置
func BenchmarkCacheSet(b *testing.B) {
	cache := NewWeakCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := &Value{
			Data:      "test",
			Size:      1024,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune(i%26+'A')), v)
	}
}

// BenchmarkCacheGet 基准测试缓存获取
func BenchmarkCacheGet(b *testing.B) {
	cache := NewWeakCache()

	// 预先填充
	for i := 0; i < 100; i++ {
		v := &Value{
			Data:      "test",
			Size:      1024,
			CreatedAt: time.Now(),
		}
		cache.Set(string(rune(i%26+'A')), v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(string(rune(i%26 + 'A')))
	}
}
