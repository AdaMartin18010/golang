package counter

import (
	"sync"
	"sync/atomic"
	"time"
)

// Counter 计数器接口
type Counter interface {
	Increment() int64
	Decrement() int64
	Add(delta int64) int64
	Get() int64
	Reset()
}

// SimpleCounter 简单计数器
type SimpleCounter struct {
	value int64
}

// NewSimpleCounter 创建简单计数器
func NewSimpleCounter() *SimpleCounter {
	return &SimpleCounter{}
}

// Increment 增加1
func (c *SimpleCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement 减少1
func (c *SimpleCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Add 增加指定值
func (c *SimpleCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Get 获取当前值
func (c *SimpleCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset 重置计数器
func (c *SimpleCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

// Set 设置值
func (c *SimpleCounter) Set(value int64) {
	atomic.StoreInt64(&c.value, value)
}

// MaxCounter 最大计数器（只增不减）
type MaxCounter struct {
	value int64
}

// NewMaxCounter 创建最大计数器
func NewMaxCounter() *MaxCounter {
	return &MaxCounter{}
}

// Increment 增加1
func (c *MaxCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Add 增加指定值
func (c *MaxCounter) Add(delta int64) int64 {
	if delta <= 0 {
		return c.Get()
	}
	return atomic.AddInt64(&c.value, delta)
}

// Get 获取当前值
func (c *MaxCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset 重置计数器
func (c *MaxCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

// MinCounter 最小计数器（只减不增）
type MinCounter struct {
	value int64
}

// NewMinCounter 创建最小计数器
func NewMinCounter(initial int64) *MinCounter {
	return &MinCounter{
		value: initial,
	}
}

// Decrement 减少1
func (c *MinCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Subtract 减少指定值
func (c *MinCounter) Subtract(delta int64) int64 {
	if delta <= 0 {
		return c.Get()
	}
	return atomic.AddInt64(&c.value, -delta)
}

// Get 获取当前值
func (c *MinCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset 重置计数器
func (c *MinCounter) Reset(initial int64) {
	atomic.StoreInt64(&c.value, initial)
}

// RateCounter 速率计数器
type RateCounter struct {
	counts    []int64
	window    time.Duration
	interval  time.Duration
	mu        sync.RWMutex
	lastIndex int
	lastTime  time.Time
}

// NewRateCounter 创建速率计数器
func NewRateCounter(window, interval time.Duration) *RateCounter {
	if interval <= 0 {
		interval = time.Second
	}
	if window <= 0 {
		window = time.Minute
	}
	
	numBuckets := int(window / interval)
	if numBuckets <= 0 {
		numBuckets = 1
	}
	
	return &RateCounter{
		counts:   make([]int64, numBuckets),
		window:   window,
		interval: interval,
		lastTime: time.Now(),
	}
}

// Increment 增加1
func (rc *RateCounter) Increment() {
	rc.Add(1)
}

// Add 增加指定值
func (rc *RateCounter) Add(delta int64) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rc.lastTime)
	
	// 更新过期的桶
	steps := int(elapsed / rc.interval)
	if steps > 0 {
		for i := 0; i < steps && i < len(rc.counts); i++ {
			rc.lastIndex = (rc.lastIndex + 1) % len(rc.counts)
			rc.counts[rc.lastIndex] = 0
		}
		rc.lastTime = now
	}
	
	// 增加当前桶的值
	rc.counts[rc.lastIndex] += delta
}

// Get 获取当前速率（每秒）
func (rc *RateCounter) Get() float64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	
	var total int64
	for _, count := range rc.counts {
		total += count
	}
	
	return float64(total) / rc.window.Seconds()
}

// Reset 重置计数器
func (rc *RateCounter) Reset() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	
	for i := range rc.counts {
		rc.counts[i] = 0
	}
	rc.lastIndex = 0
	rc.lastTime = time.Now()
}

// SlidingWindowCounter 滑动窗口计数器
type SlidingWindowCounter struct {
	buckets   []int64
	window    time.Duration
	interval  time.Duration
	mu        sync.RWMutex
	lastIndex int
	lastTime  time.Time
}

// NewSlidingWindowCounter 创建滑动窗口计数器
func NewSlidingWindowCounter(window, interval time.Duration) *SlidingWindowCounter {
	if interval <= 0 {
		interval = time.Second
	}
	if window <= 0 {
		window = time.Minute
	}
	
	numBuckets := int(window / interval)
	if numBuckets <= 0 {
		numBuckets = 1
	}
	
	return &SlidingWindowCounter{
		buckets:  make([]int64, numBuckets),
		window:   window,
		interval: interval,
		lastTime: time.Now(),
	}
}

// Increment 增加1
func (swc *SlidingWindowCounter) Increment() {
	swc.Add(1)
}

// Add 增加指定值
func (swc *SlidingWindowCounter) Add(delta int64) {
	swc.mu.Lock()
	defer swc.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(swc.lastTime)
	
	// 更新过期的桶
	steps := int(elapsed / swc.interval)
	if steps > 0 {
		for i := 0; i < steps && i < len(swc.buckets); i++ {
			swc.lastIndex = (swc.lastIndex + 1) % len(swc.buckets)
			swc.buckets[swc.lastIndex] = 0
		}
		swc.lastTime = now
	}
	
	// 增加当前桶的值
	swc.buckets[swc.lastIndex] += delta
}

// Get 获取窗口内的总数
func (swc *SlidingWindowCounter) Get() int64 {
	swc.mu.RLock()
	defer swc.mu.RUnlock()
	
	var total int64
	for _, count := range swc.buckets {
		total += count
	}
	return total
}

// Reset 重置计数器
func (swc *SlidingWindowCounter) Reset() {
	swc.mu.Lock()
	defer swc.mu.Unlock()
	
	for i := range swc.buckets {
		swc.buckets[i] = 0
	}
	swc.lastIndex = 0
	swc.lastTime = time.Now()
}

// MultiCounter 多键计数器
type MultiCounter struct {
	counters map[string]*SimpleCounter
	mu       sync.RWMutex
}

// NewMultiCounter 创建多键计数器
func NewMultiCounter() *MultiCounter {
	return &MultiCounter{
		counters: make(map[string]*SimpleCounter),
	}
}

// Increment 增加指定键的计数
func (mc *MultiCounter) Increment(key string) int64 {
	mc.mu.Lock()
	counter, ok := mc.counters[key]
	if !ok {
		counter = NewSimpleCounter()
		mc.counters[key] = counter
	}
	mc.mu.Unlock()
	
	return counter.Increment()
}

// Decrement 减少指定键的计数
func (mc *MultiCounter) Decrement(key string) int64 {
	mc.mu.Lock()
	counter, ok := mc.counters[key]
	if !ok {
		counter = NewSimpleCounter()
		mc.counters[key] = counter
	}
	mc.mu.Unlock()
	
	return counter.Decrement()
}

// Add 增加指定键的计数
func (mc *MultiCounter) Add(key string, delta int64) int64 {
	mc.mu.Lock()
	counter, ok := mc.counters[key]
	if !ok {
		counter = NewSimpleCounter()
		mc.counters[key] = counter
	}
	mc.mu.Unlock()
	
	return counter.Add(delta)
}

// Get 获取指定键的计数
func (mc *MultiCounter) Get(key string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	counter, ok := mc.counters[key]
	if !ok {
		return 0
	}
	return counter.Get()
}

// GetAll 获取所有计数
func (mc *MultiCounter) GetAll() map[string]int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	result := make(map[string]int64)
	for key, counter := range mc.counters {
		result[key] = counter.Get()
	}
	return result
}

// Reset 重置指定键的计数
func (mc *MultiCounter) Reset(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	if counter, ok := mc.counters[key]; ok {
		counter.Reset()
	}
}

// ResetAll 重置所有计数
func (mc *MultiCounter) ResetAll() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	for _, counter := range mc.counters {
		counter.Reset()
	}
}

// Keys 获取所有键
func (mc *MultiCounter) Keys() []string {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	keys := make([]string, 0, len(mc.counters))
	for key := range mc.counters {
		keys = append(keys, key)
	}
	return keys
}

