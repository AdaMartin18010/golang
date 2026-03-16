// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// cache_test.go 包含缓存功能的单元测试
package abac

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDefaultCacheConfig 测试默认缓存配置
func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	assert.Equal(t, 10000, config.MaxSize)
	assert.Equal(t, 5*time.Minute, config.DefaultTTL)
}

// TestNewCache_Defaults 测试创建缓存的默认值
func TestNewCache_Defaults(t *testing.T) {
	// 零值配置应该获得默认值
	cache := NewCache(CacheConfig{})
	assert.NotNil(t, cache)

	stats := cache.Stats()
	assert.Equal(t, 10000, stats.MaxSize)
}

// TestLRUCache_SetExistingKey 测试更新已存在的键
func TestLRUCache_SetExistingKey(t *testing.T) {
	cache := NewCache(CacheConfig{
		MaxSize:    10,
		DefaultTTL: 1 * time.Hour,
	})

	cache.Set("key1", "value1", 0)
	cache.Set("key1", "value2", 0)

	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value2", value)
}

// TestLRUCache_Clear 测试缓存清空
func TestLRUCache_Clear(t *testing.T) {
	cache := NewCache(CacheConfig{
		MaxSize:    10,
		DefaultTTL: 1 * time.Hour,
	})

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)
	cache.Clear()

	_, found := cache.Get("key1")
	assert.False(t, found)

	_, found = cache.Get("key2")
	assert.False(t, found)

	stats := cache.Stats()
	assert.Equal(t, 0, stats.Size)
}

// TestLRUCache_Stats_HitRate 测试缓存命中率
func TestLRUCache_Stats_HitRate(t *testing.T) {
	cache := NewCache(CacheConfig{
		MaxSize:    10,
		DefaultTTL: 1 * time.Hour,
	})

	// 未访问时命中率应为 0
	stats := cache.Stats()
	assert.Equal(t, float64(0), stats.HitRate)

	// 添加项
	cache.Set("key1", "value1", 0)

	// 命中
	cache.Get("key1")
	cache.Get("key1")

	// 未命中
	cache.Get("non-existing")

	stats = cache.Stats()
	assert.InDelta(t, 0.67, stats.HitRate, 0.01)
}

// TestNewEvaluationCache 测试创建评估缓存
func TestNewEvaluationCache(t *testing.T) {
	config := CacheConfig{
		MaxSize:    100,
		DefaultTTL: 5 * time.Minute,
	}

	cache := NewEvaluationCache(config)
	assert.NotNil(t, cache)
	assert.NotNil(t, cache.cache)
}

// TestEvaluationCache_NonExistingResult 测试获取不存在的评估结果
func TestEvaluationCache_NonExistingResult(t *testing.T) {
	cache := NewEvaluationCache(DefaultCacheConfig())

	req := Request{
		Subject:  Subject{ID: "user-1"},
		Resource: Resource{Type: "document", ID: "doc-1"},
		Action:   Action{Name: "read"},
	}

	_, found := cache.GetEvaluationResult("policy-1", req)
	assert.False(t, found)
}

// TestEvaluationCache_WrongType 测试获取错误类型的评估结果
func TestEvaluationCache_WrongType(t *testing.T) {
	cache := NewEvaluationCache(DefaultCacheConfig())

	req := Request{
		Subject:  Subject{ID: "user-1"},
		Resource: Resource{Type: "document", ID: "doc-1"},
		Action:   Action{Name: "read"},
	}

	// 直接存储错误类型的值
	key := cache.generateKey("policy-1", req)
	cache.cache.Set(key, "wrong type", 0)

	_, found := cache.GetEvaluationResult("policy-1", req)
	assert.False(t, found)
}

// TestEvaluationCache_Stats 测试评估缓存统计
func TestEvaluationCache_Stats(t *testing.T) {
	cache := NewEvaluationCache(CacheConfig{
		MaxSize:    10,
		DefaultTTL: 1 * time.Hour,
	})

	req := Request{
		Subject:  Subject{ID: "user-1"},
		Resource: Resource{Type: "document", ID: "doc-1"},
		Action:   Action{Name: "read"},
	}

	result := &EvaluationResult{Allowed: true, Decision: Allow}
	cache.SetEvaluationResult("policy-1", req, result, 0)

	stats := cache.Stats()
	assert.Equal(t, 1, stats.Size)
}

// TestNewCacheMiddleware 测试创建缓存中间件
func TestNewCacheMiddleware(t *testing.T) {
	engine := NewEngine()
	config := CacheConfig{
		MaxSize:    100,
		DefaultTTL: 5 * time.Minute,
	}

	middleware := NewCacheMiddleware(engine, config)
	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.engine)
	assert.NotNil(t, middleware.cache)
	assert.True(t, middleware.enabled)
}

// TestCacheMiddleware_Stats 测试缓存中间件统计
func TestCacheMiddleware_Stats(t *testing.T) {
	engine := NewEngine()
	middleware := NewCacheMiddleware(engine, DefaultCacheConfig())

	stats := middleware.Stats()
	assert.Equal(t, 0, stats.Size)
}

// TestCacheMiddleware_Evaluate_NoMatchingPolicy 测试没有匹配策略时的评估
func TestCacheMiddleware_Evaluate_NoMatchingPolicy(t *testing.T) {
	engine := NewEngine()
	middleware := NewCacheMiddleware(engine, DefaultCacheConfig())

	ctx := context.Background()
	req := createTestRequest()

	// 没有策略，应该返回默认拒绝
	result := middleware.Evaluate(ctx, req)
	assert.False(t, result.Allowed)
	assert.False(t, result.Cached)
}

// TestCacheItem_Expired 测试过期缓存项
func TestCacheItem_Expired(t *testing.T) {
	item := &cacheItem{
		key:       "key1",
		value:     "value1",
		expiresAt: time.Now().Add(-1 * time.Second),
	}

	assert.True(t, item.isExpired())
}

// TestCacheItem_NotExpired 测试未过期缓存项
func TestCacheItem_NotExpired(t *testing.T) {
	item := &cacheItem{
		key:       "key1",
		value:     "value1",
		expiresAt: time.Now().Add(1 * time.Hour),
	}

	assert.False(t, item.isExpired())
}
