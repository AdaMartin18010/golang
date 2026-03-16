package lru

import (
	"container/list"
	"sync"
)

// LRUCache LRU缓存
type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*list.Element
	list     *list.List
	mu       sync.RWMutex
}

// entry 缓存条目
type entry[K comparable, V any] struct {
	key   K
	value V
}

// NewLRUCache 创建LRU缓存
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	if capacity <= 0 {
		capacity = 100 // 默认容量
	}
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
	}
}

// Get 获取值
func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	
	if elem, ok := lru.cache[key]; ok {
		// 移动到链表头部
		lru.list.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	
	var zero V
	return zero, false
}

// Put 设置值
func (lru *LRUCache[K, V]) Put(key K, value V) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	
	if elem, ok := lru.cache[key]; ok {
		// 更新值并移动到头部
		elem.Value.(*entry[K, V]).value = value
		lru.list.MoveToFront(elem)
		return
	}
	
	// 检查容量
	if lru.list.Len() >= lru.capacity {
		// 删除尾部元素
		back := lru.list.Back()
		if back != nil {
			lru.list.Remove(back)
			delete(lru.cache, back.Value.(*entry[K, V]).key)
		}
	}
	
	// 添加到头部
	elem := lru.list.PushFront(&entry[K, V]{key: key, value: value})
	lru.cache[key] = elem
}

// Delete 删除键
func (lru *LRUCache[K, V]) Delete(key K) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	
	if elem, ok := lru.cache[key]; ok {
		lru.list.Remove(elem)
		delete(lru.cache, key)
		return true
	}
	return false
}

// Contains 检查键是否存在
func (lru *LRUCache[K, V]) Contains(key K) bool {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	_, ok := lru.cache[key]
	return ok
}

// Size 获取缓存大小
func (lru *LRUCache[K, V]) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	return lru.list.Len()
}

// Capacity 获取缓存容量
func (lru *LRUCache[K, V]) Capacity() int {
	return lru.capacity
}

// Clear 清空缓存
func (lru *LRUCache[K, V]) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.cache = make(map[K]*list.Element)
	lru.list = list.New()
}

// Keys 获取所有键
func (lru *LRUCache[K, V]) Keys() []K {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	
	keys := make([]K, 0, lru.list.Len())
	for elem := lru.list.Front(); elem != nil; elem = elem.Next() {
		keys = append(keys, elem.Value.(*entry[K, V]).key)
	}
	return keys
}

// Values 获取所有值
func (lru *LRUCache[K, V]) Values() []V {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	
	values := make([]V, 0, lru.list.Len())
	for elem := lru.list.Front(); elem != nil; elem = elem.Next() {
		values = append(values, elem.Value.(*entry[K, V]).value)
	}
	return values
}

// Peek 查看值（不更新访问顺序）
func (lru *LRUCache[K, V]) Peek(key K) (V, bool) {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	
	if elem, ok := lru.cache[key]; ok {
		return elem.Value.(*entry[K, V]).value, true
	}
	
	var zero V
	return zero, false
}

// GetOldest 获取最旧的键值对
func (lru *LRUCache[K, V]) GetOldest() (K, V, bool) {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	
	if lru.list.Len() == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	
	back := lru.list.Back()
	entry := back.Value.(*entry[K, V])
	return entry.key, entry.value, true
}

// GetNewest 获取最新的键值对
func (lru *LRUCache[K, V]) GetNewest() (K, V, bool) {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	
	if lru.list.Len() == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	
	front := lru.list.Front()
	entry := front.Value.(*entry[K, V])
	return entry.key, entry.value, true
}

// Resize 调整容量
func (lru *LRUCache[K, V]) Resize(newCapacity int) {
	if newCapacity <= 0 {
		return
	}
	
	lru.mu.Lock()
	defer lru.mu.Unlock()
	
	lru.capacity = newCapacity
	
	// 如果当前大小超过新容量，删除多余的条目
	for lru.list.Len() > newCapacity {
		back := lru.list.Back()
		if back != nil {
			lru.list.Remove(back)
			delete(lru.cache, back.Value.(*entry[K, V]).key)
		}
	}
}

