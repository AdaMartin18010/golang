// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件实现了策略缓存，用于提升 ABAC 评估性能。
//
// 缓存机制：
//   - 内存缓存：使用 LRU 算法管理缓存项
//   - TTL 支持：每个缓存项可以设置过期时间
//   - 并发安全：支持多线程并发访问
//   - 缓存统计：提供命中率等统计数据
//
// 缓存策略：
//   - 评估结果缓存：缓存相同请求的评估结果
//   - 策略缓存：缓存解析后的策略规则
//
// 使用示例：
//
//	// 创建缓存
//	cache := abac.NewCache(abac.CacheConfig{
//	    MaxSize:    10000,
//	    DefaultTTL: 5 * time.Minute,
//	})
//
//	// 设置缓存值
//	cache.Set("policy:1", policy, 10*time.Minute)
//
//	// 获取缓存值
//	if value, found := cache.Get("policy:1"); found {
//	    policy := value.(*abac.Policy)
//	}
//
//	// 获取统计信息
//	stats := cache.Stats()
//	fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
package abac

import (
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	MaxSize    int           // 最大缓存项数
	DefaultTTL time.Duration // 默认过期时间
}

// DefaultCacheConfig 返回默认缓存配置
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		MaxSize:    10000,
		DefaultTTL: 5 * time.Minute,
	}
}

// Cache 是 ABAC 策略缓存接口
type Cache interface {
	// Get 获取缓存值
	//
	// 参数：
	//   - key: 缓存键
	//
	// 返回：
	//   - value: 缓存值
	//   - found: 是否找到
	Get(key string) (interface{}, bool)

	// Set 设置缓存值
	//
	// 参数：
	//   - key: 缓存键
	//   - value: 缓存值
	//   - ttl: 过期时间（0 使用默认过期时间）
	Set(key string, value interface{}, ttl time.Duration)

	// Delete 删除缓存项
	//
	// 参数：
	//   - key: 缓存键
	Delete(key string)

	// Clear 清空所有缓存
	Clear()

	// Stats 获取缓存统计信息
	Stats() CacheStats
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits       uint64  `json:"hits"`        // 命中次数
	Misses     uint64  `json:"misses"`      // 未命中次数
	Size       int     `json:"size"`        // 当前缓存项数
	MaxSize    int     `json:"max_size"`    // 最大缓存项数
	HitRate    float64 `json:"hit_rate"`    // 命中率
	Evictions  uint64  `json:"evictions"`   // 淘汰次数
}

// cacheItem 缓存项
type cacheItem struct {
	key       string
	value     interface{}
	expiresAt time.Time
	element   *list.Element // LRU 链表元素
}

// isExpired 检查缓存项是否过期
func (item *cacheItem) isExpired() bool {
	return time.Now().After(item.expiresAt)
}

// LRUCache 实现 LRU 淘汰策略的缓存
type LRUCache struct {
	config    CacheConfig
	items     map[string]*cacheItem
	lruList   *list.List // 用于 LRU 淘汰
	mu        sync.RWMutex
	hits      uint64
	misses    uint64
	evictions uint64
}

// NewCache 创建新的缓存实例
//
// 参数：
//   - config: 缓存配置
//
// 返回：
//   - Cache: 缓存接口
func NewCache(config CacheConfig) Cache {
	if config.MaxSize <= 0 {
		config.MaxSize = 10000
	}
	if config.DefaultTTL <= 0 {
		config.DefaultTTL = 5 * time.Minute
	}

	return &LRUCache{
		config:  config,
		items:   make(map[string]*cacheItem),
		lruList: list.New(),
	}
}

// Get 获取缓存值
//
// 如果缓存项不存在或已过期，返回 found=false
//
// 参数：
//   - key: 缓存键
//
// 返回：
//   - value: 缓存值
//   - found: 是否找到有效缓存
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		c.incrementMisses()
		return nil, false
	}

	// 检查是否过期
	if item.isExpired() {
		c.mu.Lock()
		c.deleteItem(item)
		c.mu.Unlock()
		c.incrementMisses()
		return nil, false
	}

	// 更新 LRU 位置（移到链表尾部表示最近使用）
	c.mu.Lock()
	c.lruList.MoveToBack(item.element)
	c.mu.Unlock()

	c.incrementHits()
	return item.value, true
}

// Set 设置缓存值
//
// 如果缓存已满，会淘汰最久未使用的项
//
// 参数：
//   - key: 缓存键
//   - value: 缓存值
//   - ttl: 过期时间（0 使用默认过期时间）
func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) {
	if ttl <= 0 {
		ttl = c.config.DefaultTTL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果键已存在，更新值和过期时间
	if item, exists := c.items[key]; exists {
		item.value = value
		item.expiresAt = time.Now().Add(ttl)
		c.lruList.MoveToBack(item.element)
		return
	}

	// 检查是否需要淘汰
	if len(c.items) >= c.config.MaxSize {
		c.evictLRU()
	}

	// 创建新缓存项
	item := &cacheItem{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	// 添加到 LRU 链表尾部
	item.element = c.lruList.PushBack(item)
	c.items[key] = item
}

// Delete 删除缓存项
//
// 参数：
//   - key: 缓存键
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, exists := c.items[key]; exists {
		c.deleteItem(item)
	}
}

// Clear 清空所有缓存
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	c.lruList.Init()
}

// Stats 获取缓存统计信息
//
// 返回：
//   - CacheStats: 缓存统计信息
func (c *LRUCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return CacheStats{
		Hits:      c.hits,
		Misses:    c.misses,
		Size:      len(c.items),
		MaxSize:   c.config.MaxSize,
		HitRate:   hitRate,
		Evictions: c.evictions,
	}
}

// deleteItem 删除缓存项（必须在持有写锁时调用）
func (c *LRUCache) deleteItem(item *cacheItem) {
	c.lruList.Remove(item.element)
	delete(c.items, item.key)
}

// evictLRU 淘汰最久未使用的缓存项（必须在持有写锁时调用）
func (c *LRUCache) evictLRU() {
	elem := c.lruList.Front()
	if elem == nil {
		return
	}

	item := elem.Value.(*cacheItem)
	c.deleteItem(item)
	c.evictions++
}

// incrementHits 增加命中计数
func (c *LRUCache) incrementHits() {
	c.mu.Lock()
	c.hits++
	c.mu.Unlock()
}

// incrementMisses 增加未命中计数
func (c *LRUCache) incrementMisses() {
	c.mu.Lock()
	c.misses++
	c.mu.Unlock()
}

// EvaluationCache 评估结果缓存
//
// 专门用于缓存 ABAC 评估结果
type EvaluationCache struct {
	cache  Cache
	config CacheConfig
}

// EvaluationKey 用于缓存评估结果的键
type EvaluationKey struct {
	PolicyID string `json:"policy_id"`
	Subject  string `json:"subject"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// NewEvaluationCache 创建评估结果缓存
//
// 参数：
//   - config: 缓存配置
//
// 返回：
//   - *EvaluationCache: 评估缓存实例
func NewEvaluationCache(config CacheConfig) *EvaluationCache {
	return &EvaluationCache{
		cache:  NewCache(config),
		config: config,
	}
}

// GetEvaluationResult 获取评估结果缓存
//
// 参数：
//   - policyID: 策略ID
//   - req: 访问请求
//
// 返回：
//   - result: 评估结果
//   - found: 是否找到缓存
func (ec *EvaluationCache) GetEvaluationResult(policyID string, req Request) (*EvaluationResult, bool) {
	key := ec.generateKey(policyID, req)
	value, found := ec.cache.Get(key)
	if !found {
		return nil, false
	}

	result, ok := value.(*EvaluationResult)
	if !ok {
		return nil, false
	}

	return result, true
}

// SetEvaluationResult 设置评估结果缓存
//
// 参数：
//   - policyID: 策略ID
//   - req: 访问请求
//   - result: 评估结果
//   - ttl: 过期时间
func (ec *EvaluationCache) SetEvaluationResult(policyID string, req Request, result *EvaluationResult, ttl time.Duration) {
	key := ec.generateKey(policyID, req)
	ec.cache.Set(key, result, ttl)
}

// generateKey 生成缓存键
func (ec *EvaluationCache) generateKey(policyID string, req Request) string {
	key := EvaluationKey{
		PolicyID: policyID,
		Subject:  req.Subject.ID,
		Resource: fmt.Sprintf("%s:%s", req.Resource.Type, req.Resource.ID),
		Action:   req.Action.Name,
	}

	data, _ := json.Marshal(key)
	hash := sha256.Sum256(data)
	return "eval:" + hex.EncodeToString(hash[:8])
}

// Stats 获取缓存统计
//
// 返回：
//   - CacheStats: 缓存统计信息
func (ec *EvaluationCache) Stats() CacheStats {
	return ec.cache.Stats()
}

// Clear 清空缓存
func (ec *EvaluationCache) Clear() {
	ec.cache.Clear()
}

// EvaluationResult 评估结果
type EvaluationResult struct {
	Allowed    bool     `json:"allowed"`
	Decision   Effect   `json:"decision"`
	Reason     string   `json:"reason"`
	EvaluatedAt int64   `json:"evaluated_at"`
}

// CacheMiddleware 缓存中间件
//
// 包装 Engine，添加评估结果缓存功能
type CacheMiddleware struct {
	engine *Engine
	cache  *EvaluationCache
	enabled bool
}

// NewCacheMiddleware 创建缓存中间件
//
// 参数：
//   - engine: ABAC 引擎
//   - config: 缓存配置
//
// 返回：
//   - *CacheMiddleware: 缓存中间件
func NewCacheMiddleware(engine *Engine, config CacheConfig) *CacheMiddleware {
	return &CacheMiddleware{
		engine:  engine,
		cache:   NewEvaluationCache(config),
		enabled: true,
	}
}

// SetEnabled 启用/禁用缓存
//
// 参数：
//   - enabled: 是否启用
func (cm *CacheMiddleware) SetEnabled(enabled bool) {
	cm.enabled = enabled
}

// Evaluate 执行评估（带缓存）
//
// 参数：
//   - ctx: 上下文
//   - req: 访问请求
//
// 返回：
//   - Result: 评估结果
func (cm *CacheMiddleware) Evaluate(ctx context.Context, req Request) Result {
	if !cm.enabled {
		return cm.engine.Evaluate(ctx, req)
	}

	// 尝试从缓存获取
	policies := cm.engine.ListPolicies()
	for _, policy := range policies {
		if !policy.Enabled {
			continue
		}

		cached, found := cm.cache.GetEvaluationResult(policy.ID, req)
		if found {
			return Result{
				Allowed:    cached.Allowed,
				Decision:   cached.Decision,
				Reason:     cached.Reason + " (cached)",
				Cached:     true,
			}
		}
	}

	// 执行实际评估
	result := cm.engine.Evaluate(ctx, req)

	// 缓存结果
	if result.MatchedPolicy != nil {
		cm.cache.SetEvaluationResult(
			result.MatchedPolicy.ID,
			req,
			&EvaluationResult{
				Allowed:     result.Allowed,
				Decision:    result.Decision,
				Reason:      result.Reason,
				EvaluatedAt: time.Now().Unix(),
			},
			cm.cache.config.DefaultTTL,
		)
	}

	return result
}

// InvalidateCache 使缓存失效
//
// 通常在策略变更时调用
func (cm *CacheMiddleware) InvalidateCache() {
	cm.cache.Clear()
}

// Stats 获取缓存统计
//
// 返回：
//   - CacheStats: 缓存统计信息
func (cm *CacheMiddleware) Stats() CacheStats {
	return cm.cache.Stats()
}
