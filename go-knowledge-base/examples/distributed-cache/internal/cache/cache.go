package cache

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyTooLarge = errors.New("key size exceeds maximum")
	ErrValueTooLarge = errors.New("value size exceeds maximum")
)

// EvictionPolicy defines the cache eviction policy
type EvictionPolicy int

const (
	LRU EvictionPolicy = iota
	LFU
	TTL
)

// Config holds cache configuration
type Config struct {
	MaxSize        int64
	MaxKeySize     int
	MaxValueSize   int
	EvictionPolicy EvictionPolicy
	DefaultTTL     time.Duration
}

// Item represents a cached item
type Item struct {
	Key        string
	Value      []byte
	Expiration int64
	Frequency  int64
	LastAccess int64
	element    *list.Element
}

// IsExpired checks if the item has expired
func (i *Item) IsExpired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

// Cache provides thread-safe caching
type Cache struct {
	mu          sync.RWMutex
	data        map[string]*Item
	lruList     *list.List
	maxSize     int64
	currentSize int64
	config      *Config
	
	// Stats
	hits   int64
	misses int64
}

// New creates a new cache instance
func New(config *Config) *Cache {
	if config.MaxSize <= 0 {
		config.MaxSize = 1024 * 1024 * 1024 // 1GB default
	}
	if config.MaxKeySize <= 0 {
		config.MaxKeySize = 1024 // 1KB
	}
	if config.MaxValueSize <= 0 {
		config.MaxValueSize = 1024 * 1024 // 1MB
	}
	
	return &Cache{
		data:    make(map[string]*Item),
		lruList: list.New(),
		maxSize: config.MaxSize,
		config:  config,
	}
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()
	
	if !exists {
		atomic.AddInt64(&c.misses, 1)
		return nil, ErrKeyNotFound
	}
	
	if item.IsExpired() {
		c.Delete(key)
		atomic.AddInt64(&c.misses, 1)
		return nil, ErrKeyNotFound
	}
	
	// Update access patterns
	atomic.StoreInt64(&item.LastAccess, time.Now().UnixNano())
	atomic.AddInt64(&item.Frequency, 1)
	
	// Update LRU list
	c.mu.Lock()
	c.lruList.MoveToFront(item.element)
	c.mu.Unlock()
	
	atomic.AddInt64(&c.hits, 1)
	return item.Value, nil
}

// Set stores a value in cache
func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	if len(key) > c.config.MaxKeySize {
		return ErrKeyTooLarge
	}
	if len(value) > c.config.MaxValueSize {
		return ErrValueTooLarge
	}
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Check if key already exists
	if existing, ok := c.data[key]; ok {
		c.currentSize -= int64(len(existing.Key) + len(existing.Value))
		existing.Value = value
		existing.Expiration = c.expirationFromTTL(ttl)
		c.currentSize += int64(len(key) + len(value))
		c.lruList.MoveToFront(existing.element)
		return nil
	}
	
	// Evict if necessary
	for c.currentSize+int64(len(key)+len(value)) > c.maxSize {
		if !c.evict() {
			break
		}
	}
	
	// Create new item
	item := &Item{
		Key:        key,
		Value:      value,
		Expiration: c.expirationFromTTL(ttl),
		Frequency:  1,
		LastAccess: time.Now().UnixNano(),
	}
	
	// Add to LRU list
	elem := c.lruList.PushFront(item)
	item.element = elem
	
	c.data[key] = item
	c.currentSize += int64(len(key) + len(value))
	
	return nil
}

// Delete removes a key from cache
func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if item, exists := c.data[key]; exists {
		c.lruList.Remove(item.element)
		c.currentSize -= int64(len(item.Key) + len(item.Value))
		delete(c.data, key)
		return true
	}
	
	return false
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()
	
	if !exists {
		return false
	}
	
	if item.IsExpired() {
		c.Delete(key)
		return false
	}
	
	return true
}

// Flush removes all items from cache
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data = make(map[string]*Item)
	c.lruList = list.New()
	c.currentSize = 0
}

// Stats returns cache statistics
func (c *Cache) Stats() Stats {
	return Stats{
		Size:      c.currentSize,
		MaxSize:   c.maxSize,
		Items:     int64(len(c.data)),
		Hits:      atomic.LoadInt64(&c.hits),
		Misses:    atomic.LoadInt64(&c.misses),
	}
}

// Stats holds cache statistics
type Stats struct {
	Size    int64
	MaxSize int64
	Items   int64
	Hits    int64
	Misses  int64
}

// HitRatio returns the cache hit ratio
func (s Stats) HitRatio() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total)
}

// evict removes the least recently used item
func (c *Cache) evict() bool {
	elem := c.lruList.Back()
	if elem == nil {
		return false
	}
	
	item := elem.Value.(*Item)
	c.lruList.Remove(elem)
	c.currentSize -= int64(len(item.Key) + len(item.Value))
	delete(c.data, item.Key)
	
	return true
}

// expirationFromTTL calculates expiration time from TTL
func (c *Cache) expirationFromTTL(ttl time.Duration) int64 {
	if ttl <= 0 {
		if c.config.DefaultTTL > 0 {
			ttl = c.config.DefaultTTL
		} else {
			return 0 // No expiration
		}
	}
	return time.Now().Add(ttl).UnixNano()
}

// Keys returns all non-expired keys
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	keys := make([]string, 0, len(c.data))
	now := time.Now().UnixNano()
	
	for key, item := range c.data {
		if item.Expiration == 0 || item.Expiration > now {
			keys = append(keys, key)
		}
	}
	
	return keys
}
