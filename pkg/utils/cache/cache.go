package cache

import (
	"context"
	"sync"
	"time"
)

// Item 缓存项
type Item struct {
	Value      interface{}
	Expiration time.Time
}

// IsExpired 检查是否过期
func (i *Item) IsExpired() bool {
	return !i.Expiration.IsZero() && time.Now().After(i.Expiration)
}

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Size() int
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*Item
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]*Item),
	}
	go c.startCleanup()
	return c
}

// Get 获取缓存
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if item.IsExpired() {
		delete(c.items, key)
		return nil, false
	}

	return item.Value, true
}

// Set 设置缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.items[key] = &Item{
		Value:      value,
		Expiration: expiration,
	}

	return nil
}

// Delete 删除缓存
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Clear 清空缓存
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*Item)
	return nil
}

// Size 获取缓存大小
func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// startCleanup 启动清理协程
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup 清理过期项
func (c *MemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if !item.Expiration.IsZero() && now.After(item.Expiration) {
			delete(c.items, key)
		}
	}
}

// GetOrSet 获取或设置缓存
func GetOrSet(
	ctx context.Context,
	cache Cache,
	key string,
	ttl time.Duration,
	fn func() (interface{}, error),
) (interface{}, error) {
	// 尝试获取
	if value, ok := cache.Get(ctx, key); ok {
		return value, nil
	}

	// 执行函数
	value, err := fn()
	if err != nil {
		return nil, err
	}

	// 设置缓存
	if err := cache.Set(ctx, key, value, ttl); err != nil {
		return nil, err
	}

	return value, nil
}
