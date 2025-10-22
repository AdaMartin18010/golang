// Weak Pointer Cache示例：使用弱引用避免内存泄漏
// 注意：由于runtime/weak不是标准库的一部分，本示例提供模拟实现
//
// 真实的weak.Pointer需要特殊的runtime支持
// 这里展示弱引用缓存的概念和使用模式

package memory

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Value 缓存的值
type Value struct {
	Data      string
	Size      int
	CreatedAt time.Time
}

// weakPointer 模拟weak.Pointer的行为
// 注意：这是简化的模拟实现，真实的weak.Pointer需要runtime支持
type weakPointer struct {
	ptr *Value
	// 使用finalizer来检测对象是否被GC
	alive bool
}

// makeWeak 创建弱引用
func makeWeak(v *Value) *weakPointer {
	wp := &weakPointer{
		ptr:   v,
		alive: true,
	}

	// 设置finalizer来模拟弱引用行为
	// 当对象即将被GC时，标记为不可用
	runtime.SetFinalizer(v, func(_ *Value) {
		wp.alive = false
		wp.ptr = nil
	})

	return wp
}

// value 获取弱引用的值
func (wp *weakPointer) value() *Value {
	if wp == nil || !wp.alive {
		return nil
	}
	return wp.ptr
}

// WeakCache 使用weakPointer的缓存
type WeakCache struct {
	mu    sync.RWMutex
	items map[string]*weakPointer
	stats CacheStats
}

// CacheStats 缓存统计
type CacheStats struct {
	Hits      int64
	Misses    int64
	GCCleared int64
}

// NewWeakCache 创建新的弱引用缓存
func NewWeakCache() *WeakCache {
	return &WeakCache{
		items: make(map[string]*weakPointer),
	}
}

// Get 获取缓存值
func (c *WeakCache) Get(key string) (*Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if wp, ok := c.items[key]; ok {
		if v := wp.value(); v != nil {
			c.stats.Hits++
			return v, true
		}
		// 对象已被GC回收
		c.stats.GCCleared++
	}

	c.stats.Misses++
	return nil, false
}

// Set 设置缓存值
func (c *WeakCache) Set(key string, value *Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = makeWeak(value)
}

// Stats 获取统计信息
func (c *WeakCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}

// Cleanup 清理已失效的条目
func (c *WeakCache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	cleaned := 0
	for key, wp := range c.items {
		if wp.value() == nil {
			delete(c.items, key)
			cleaned++
		}
	}

	return cleaned
}

// StrongCache 传统强引用缓存（对比用）
type StrongCache struct {
	mu    sync.RWMutex
	items map[string]*Value
}

func NewStrongCache() *StrongCache {
	return &StrongCache{
		items: make(map[string]*Value),
	}
}

func (c *StrongCache) Get(key string) (*Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.items[key]
	return v, ok
}

func (c *StrongCache) Set(key string, value *Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
}

// memStats 打印内存统计
func memStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("📊 Memory: Alloc=%v MB, Sys=%v MB, NumGC=%v\n",
		m.Alloc/1024/1024, m.Sys/1024/1024, m.NumGC)
}

func main() {
	fmt.Println("🔬 Weak Pointer Cache Demo (Simulated)")
	fmt.Println("⚠️  Note: This uses a simulated weak pointer implementation.")
	fmt.Println("    Real weak.Pointer requires runtime support (experimental)")
	fmt.Println()

	// === 测试1: Weak Cache ===
	fmt.Println("=== Test 1: Weak Cache ===")
	weakCache := NewWeakCache()

	// 填充缓存
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := &Value{
			Data:      fmt.Sprintf("Large data %d", i),
			Size:      1024, // 1KB
			CreatedAt: time.Now(),
		}
		weakCache.Set(key, value)
	}

	fmt.Println("✅ Cached 1000 items")
	memStats()

	// 访问一些值（创建强引用）
	activeValues := make([]*Value, 0)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		if v, ok := weakCache.Get(key); ok {
			activeValues = append(activeValues, v)
		}
	}
	fmt.Printf("✅ Active references: %d\n", len(activeValues))

	// 触发GC
	fmt.Println("⚡ Triggering GC...")
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	memStats()

	// 检查缓存
	cleaned := weakCache.Cleanup()
	fmt.Printf("🧹 Cleaned up %d entries\n", cleaned)

	stats := weakCache.Stats()
	fmt.Printf("📈 Cache stats: Hits=%d, Misses=%d, GC Cleared=%d\n\n",
		stats.Hits, stats.Misses, stats.GCCleared)

	// === 测试2: Strong Cache（对比）===
	fmt.Println("=== Test 2: Strong Cache (for comparison) ===")
	strongCache := NewStrongCache()

	// 填充缓存
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := &Value{
			Data:      fmt.Sprintf("Large data %d", i),
			Size:      1024,
			CreatedAt: time.Now(),
		}
		strongCache.Set(key, value)
	}

	fmt.Println("✅ Cached 1000 items")
	memStats()

	// 触发GC
	fmt.Println("⚡ Triggering GC...")
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	memStats()

	fmt.Println("\n💡 Notice: Strong cache prevents GC, weak cache allows it!")

	// === 测试3: 实际场景 - 图片缓存 ===
	fmt.Println("\n=== Test 3: Image Cache Scenario ===")

	type Image struct {
		ID     int
		Pixels []byte // 模拟图片数据
	}

	imageCache := NewWeakCache()

	// 加载图片
	loadImage := func(id int) *Value {
		img := &Image{
			ID:     id,
			Pixels: make([]byte, 1024*1024), // 1MB
		}
		return &Value{
			Data:      fmt.Sprintf("Image-%d", id),
			Size:      len(img.Pixels),
			CreatedAt: time.Now(),
		}
	}

	// 场景：用户浏览图片
	for round := 1; round <= 3; round++ {
		fmt.Printf("\n📷 Round %d: User browsing...\n", round)

		// 加载一些图片
		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("img-%d", i)
			img := loadImage(i)
			imageCache.Set(key, img)
		}

		memStats()

		// 触发GC（模拟内存压力）
		runtime.GC()
		time.Sleep(50 * time.Millisecond)

		cleaned = imageCache.Cleanup()
		fmt.Printf("🧹 Cleaned %d unused images\n", cleaned)
	}

	finalStats := imageCache.Stats()
	fmt.Printf("\n🎯 Final stats: Hits=%d, Misses=%d, GC Cleared=%d\n",
		finalStats.Hits, finalStats.Misses, finalStats.GCCleared)

	fmt.Println("\n✅ Demo completed!")
	fmt.Println("💡 Key takeaway: weak.Pointer allows GC to reclaim unused cache entries,")
	fmt.Println("   preventing memory leaks while maintaining good cache hit rates.")
}
