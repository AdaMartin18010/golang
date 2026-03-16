package pool

import (
	"strings"
	"sync"
)

// Pool 对象池接口
type Pool[T any] interface {
	Get() T
	Put(item T)
	Clear()
	Size() int
}

// SimplePool 简单对象池
type SimplePool[T any] struct {
	pool    []T
	newFunc func() T
	mu      sync.Mutex
}

// NewSimplePool 创建简单对象池
func NewSimplePool[T any](newFunc func() T) *SimplePool[T] {
	return &SimplePool[T]{
		pool:    make([]T, 0),
		newFunc: newFunc,
	}
}

// Get 获取对象
func (p *SimplePool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if len(p.pool) > 0 {
		item := p.pool[len(p.pool)-1]
		p.pool = p.pool[:len(p.pool)-1]
		return item
	}
	
	if p.newFunc != nil {
		return p.newFunc()
	}
	
	var zero T
	return zero
}

// Put 归还对象
func (p *SimplePool[T]) Put(item T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool = append(p.pool, item)
}

// Clear 清空对象池
func (p *SimplePool[T]) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool = p.pool[:0]
}

// Size 获取对象池大小
func (p *SimplePool[T]) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.pool)
}

// BoundedPool 有界对象池
type BoundedPool[T any] struct {
	pool    chan T
	newFunc func() T
	maxSize int
}

// NewBoundedPool 创建有界对象池
func NewBoundedPool[T any](maxSize int, newFunc func() T) *BoundedPool[T] {
	if maxSize <= 0 {
		maxSize = 100 // 默认容量
	}
	return &BoundedPool[T]{
		pool:    make(chan T, maxSize),
		newFunc: newFunc,
		maxSize: maxSize,
	}
}

// Get 获取对象
func (p *BoundedPool[T]) Get() T {
	select {
	case item := <-p.pool:
		return item
	default:
		if p.newFunc != nil {
			return p.newFunc()
		}
		var zero T
		return zero
	}
}

// Put 归还对象
func (p *BoundedPool[T]) Put(item T) {
	select {
	case p.pool <- item:
	default:
		// 池已满，丢弃对象
	}
}

// Clear 清空对象池
func (p *BoundedPool[T]) Clear() {
	for {
		select {
		case <-p.pool:
		default:
			return
		}
	}
}

// Size 获取对象池大小
func (p *BoundedPool[T]) Size() int {
	return len(p.pool)
}

// Capacity 获取对象池容量
func (p *BoundedPool[T]) Capacity() int {
	return p.maxSize
}

// BufferPool 缓冲区池
type BufferPool struct {
	pool *sync.Pool
}

// NewBufferPool 创建缓冲区池
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024)
			},
		},
	}
}

// Get 获取缓冲区
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

// Put 归还缓冲区
func (bp *BufferPool) Put(buf []byte) {
	// 只归还合理大小的缓冲区
	if cap(buf) <= 64*1024 {
		buf = buf[:0]
		bp.pool.Put(buf)
	}
}

// StringBuilderPool 字符串构建器池
type StringBuilderPool struct {
	pool *sync.Pool
}

// NewStringBuilderPool 创建字符串构建器池
func NewStringBuilderPool() *StringBuilderPool {
	return &StringBuilderPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

// Get 获取字符串构建器
func (sbp *StringBuilderPool) Get() *strings.Builder {
	return sbp.pool.Get().(*strings.Builder)
}

// Put 归还字符串构建器
func (sbp *StringBuilderPool) Put(sb *strings.Builder) {
	sb.Reset()
	sbp.pool.Put(sb)
}

