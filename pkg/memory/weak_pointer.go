// Package memory provides weak pointer cache using Go 1.24+ weak package.
//
// Go 1.24 introduced the weak package (import "weak") which provides
// weak.Pointer[T] for non-intrusive references that do not prevent GC.
// In Go 1.26.2 this is fully stable and production-ready.
package memory

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"weak"
)

// Value 缓存的值
type Value struct {
	Data      string
	Size      int
	CreatedAt time.Time
}

// WeakCache 使用 weak.Pointer 的缓存
// 允许 GC 回收未被外部引用的缓存条目，防止内存泄漏
type WeakCache struct {
	mu    sync.RWMutex
	items map[string]weak.Pointer[Value]
	stats CacheStats
}

// CacheStats 缓存统计
type CacheStats struct {
	Hits      int64
	Misses    int64
	GCEvicted int64
}

// NewWeakCache 创建新的弱引用缓存
func NewWeakCache() *WeakCache {
	return &WeakCache{
		items: make(map[string]weak.Pointer[Value]),
	}
}

// Get 获取缓存值
// 如果值已被 GC 回收，返回 (nil, false)
func (c *WeakCache) Get(key string) (*Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	wp, ok := c.items[key]
	if !ok {
		c.stats.Misses++
		return nil, false
	}

	if v := wp.Value(); v != nil {
		c.stats.Hits++
		return v, true
	}

	// 对象已被 GC 回收
	c.stats.GCEvicted++
	return nil, false
}

// Set 设置缓存值
func (c *WeakCache) Set(key string, value *Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = weak.Make(value)
}

// Stats 获取统计信息（浅拷贝）
func (c *WeakCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}

// Cleanup 清理已失效的条目（已被 GC 回收的弱引用）
// 返回清理的条目数量
func (c *WeakCache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	cleaned := 0
	for key, wp := range c.items {
		if wp.Value() == nil {
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

// NewStrongCache 创建强引用缓存
func NewStrongCache() *StrongCache {
	return &StrongCache{
		items: make(map[string]*Value),
	}
}

// Get 获取缓存值
func (c *StrongCache) Get(key string) (*Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.items[key]
	return v, ok
}

// Set 设置缓存值
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
