package memory

import (
	"sync"
	"sync/atomic"
	"time"
)

// =============================================================================
// 内存池增强 - Enhanced Memory Pool
// =============================================================================

// PoolStats 内存池统计信息
type PoolStats struct {
	Gets      uint64 // 获取次数
	Puts      uint64 // 归还次数
	Hits      uint64 // 命中次数（从池中获取）
	Misses    uint64 // 未命中次数（新建对象）
	Size      int    // 当前池大小
	MaxSize   int    // 最大池大小
	ObjectAge uint64 // 对象平均年龄（毫秒）
}

// ObjectPool 对象池接口
type ObjectPool interface {
	Get() interface{}
	Put(interface{})
	Stats() PoolStats
	Clear()
	Size() int
}

// =============================================================================
// 通用对象池 - Generic Object Pool
// =============================================================================

// GenericPool 通用对象池（支持统计）
type GenericPool struct {
	pool    sync.Pool
	new     func() interface{}
	reset   func(interface{})
	stats   *poolStatsTracker
	maxSize int32
	size    int32
}

// poolStatsTracker 统计跟踪器
type poolStatsTracker struct {
	gets      uint64
	puts      uint64
	hits      uint64
	misses    uint64
	totalAge  uint64
	ageCount  uint64
	startTime time.Time
}

// NewGenericPool 创建通用对象池
func NewGenericPool(new func() interface{}, reset func(interface{}), maxSize int) *GenericPool {
	p := &GenericPool{
		new:     new,
		reset:   reset,
		maxSize: int32(maxSize),
		stats: &poolStatsTracker{
			startTime: time.Now(),
		},
	}

	p.pool.New = func() interface{} {
		atomic.AddUint64(&p.stats.misses, 1)
		return p.new()
	}

	return p
}

// Get 获取对象
func (p *GenericPool) Get() interface{} {
	atomic.AddUint64(&p.stats.gets, 1)

	obj := p.pool.Get()

	// 如果从池中获取到对象，记录命中
	if obj != nil {
		atomic.AddUint64(&p.stats.hits, 1)
	}

	return obj
}

// Put 归还对象
func (p *GenericPool) Put(obj interface{}) {
	if obj == nil {
		return
	}

	// 检查池大小限制
	currentSize := atomic.LoadInt32(&p.size)
	if p.maxSize > 0 && currentSize >= p.maxSize {
		// 池已满，丢弃对象
		return
	}

	// 重置对象状态
	if p.reset != nil {
		p.reset(obj)
	}

	atomic.AddUint64(&p.stats.puts, 1)
	atomic.AddInt32(&p.size, 1)

	p.pool.Put(obj)
}

// Stats 获取统计信息
func (p *GenericPool) Stats() PoolStats {
	gets := atomic.LoadUint64(&p.stats.gets)
	puts := atomic.LoadUint64(&p.stats.puts)
	hits := atomic.LoadUint64(&p.stats.hits)
	misses := atomic.LoadUint64(&p.stats.misses)

	var avgAge uint64
	if p.stats.ageCount > 0 {
		avgAge = p.stats.totalAge / p.stats.ageCount
	}

	return PoolStats{
		Gets:      gets,
		Puts:      puts,
		Hits:      hits,
		Misses:    misses,
		Size:      int(atomic.LoadInt32(&p.size)),
		MaxSize:   int(p.maxSize),
		ObjectAge: avgAge,
	}
}

// HitRate 命中率
func (p *GenericPool) HitRate() float64 {
	stats := p.Stats()
	if stats.Gets == 0 {
		return 0
	}
	return float64(stats.Hits) / float64(stats.Gets) * 100
}

// Clear 清空池
func (p *GenericPool) Clear() {
	atomic.StoreInt32(&p.size, 0)
	p.pool = sync.Pool{New: func() interface{} {
		atomic.AddUint64(&p.stats.misses, 1)
		return p.new()
	}}
}

// Size 当前池大小
func (p *GenericPool) Size() int {
	return int(atomic.LoadInt32(&p.size))
}

// =============================================================================
// 字节池 - Byte Pool (优化版)
// =============================================================================

// BytePool 字节切片池
type BytePool struct {
	pools []*sync.Pool
	sizes []int
}

// NewBytePool 创建字节池
// sizes: 不同大小的缓冲区，例如 []int{256, 1024, 4096, 65536}
func NewBytePool(sizes []int) *BytePool {
	bp := &BytePool{
		pools: make([]*sync.Pool, len(sizes)),
		sizes: sizes,
	}

	for i, size := range sizes {
		bufSize := size
		bp.pools[i] = &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, bufSize)
				return &buf
			},
		}
	}

	return bp
}

// Get 获取指定大小的字节切片
func (bp *BytePool) Get(size int) *[]byte {
	// 找到最合适的池
	for i, poolSize := range bp.sizes {
		if size <= poolSize {
			buf := bp.pools[i].Get().(*[]byte)
			// 重置切片长度
			*buf = (*buf)[:size]
			return buf
		}
	}

	// 如果没有合适的池，直接分配
	buf := make([]byte, size)
	return &buf
}

// Put 归还字节切片
func (bp *BytePool) Put(buf *[]byte) {
	if buf == nil || *buf == nil {
		return
	}

	capacity := cap(*buf)

	// 找到对应的池
	for i, poolSize := range bp.sizes {
		if capacity == poolSize {
			// 重置切片，避免内存泄漏
			*buf = (*buf)[:0]
			bp.pools[i].Put(buf)
			return
		}
	}

	// 不属于任何池，让GC回收
}

// =============================================================================
// 对象池管理器 - Pool Manager
// =============================================================================

// PoolManager 对象池管理器
type PoolManager struct {
	pools map[string]ObjectPool
	mu    sync.RWMutex
}

// NewPoolManager 创建池管理器
func NewPoolManager() *PoolManager {
	return &PoolManager{
		pools: make(map[string]ObjectPool),
	}
}

// Register 注册池
func (pm *PoolManager) Register(name string, pool ObjectPool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.pools[name] = pool
}

// Get 获取池
func (pm *PoolManager) Get(name string) (ObjectPool, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	pool, ok := pm.pools[name]
	return pool, ok
}

// AllStats 获取所有池的统计
func (pm *PoolManager) AllStats() map[string]PoolStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	stats := make(map[string]PoolStats)
	for name, pool := range pm.pools {
		stats[name] = pool.Stats()
	}
	return stats
}

// ClearAll 清空所有池
func (pm *PoolManager) ClearAll() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, pool := range pm.pools {
		pool.Clear()
	}
}

// =============================================================================
// 预定义池 - Predefined Pools
// =============================================================================

var (
	// DefaultBytePool 默认字节池
	DefaultBytePool = NewBytePool([]int{
		256,    // 256B
		1024,   // 1KB
		4096,   // 4KB
		16384,  // 16KB
		65536,  // 64KB
		262144, // 256KB
	})

	// DefaultPoolManager 默认池管理器
	DefaultPoolManager = NewPoolManager()
)

// GetBytes 从默认池获取字节切片
func GetBytes(size int) *[]byte {
	return DefaultBytePool.Get(size)
}

// PutBytes 归还字节切片到默认池
func PutBytes(buf *[]byte) {
	DefaultBytePool.Put(buf)
}
