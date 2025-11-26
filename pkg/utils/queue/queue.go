package queue

import (
	"sync"
)

// Queue 队列接口
type Queue[T any] interface {
	Enqueue(item T)
	Dequeue() (T, bool)
	Peek() (T, bool)
	Size() int
	IsEmpty() bool
	Clear()
}

// SimpleQueue 简单队列实现
type SimpleQueue[T any] struct {
	items []T
	mu    sync.RWMutex
}

// NewSimpleQueue 创建简单队列
func NewSimpleQueue[T any]() *SimpleQueue[T] {
	return &SimpleQueue[T]{
		items: make([]T, 0),
	}
}

// Enqueue 入队
func (q *SimpleQueue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

// Dequeue 出队
func (q *SimpleQueue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek 查看队首元素
func (q *SimpleQueue[T]) Peek() (T, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// Size 获取队列大小
func (q *SimpleQueue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items)
}

// IsEmpty 检查队列是否为空
func (q *SimpleQueue[T]) IsEmpty() bool {
	return q.Size() == 0
}

// Clear 清空队列
func (q *SimpleQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = q.items[:0]
}

// ToSlice 转换为切片
func (q *SimpleQueue[T]) ToSlice() []T {
	q.mu.RLock()
	defer q.mu.RUnlock()
	result := make([]T, len(q.items))
	copy(result, q.items)
	return result
}

// PriorityQueue 优先队列
type PriorityQueue[T comparable] struct {
	items []priorityItem[T]
	mu    sync.RWMutex
}

type priorityItem[T comparable] struct {
	value    T
	priority int
}

// NewPriorityQueue 创建优先队列
func NewPriorityQueue[T comparable]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items: make([]priorityItem[T], 0),
	}
}

// Enqueue 入队（带优先级）
func (pq *PriorityQueue[T]) Enqueue(item T, priority int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	
	newItem := priorityItem[T]{
		value:    item,
		priority: priority,
	}
	
	// 插入排序，优先级高的在前
	inserted := false
	for i, existing := range pq.items {
		if priority > existing.priority {
			pq.items = append(pq.items[:i], append([]priorityItem[T]{newItem}, pq.items[i:]...)...)
			inserted = true
			break
		}
	}
	
	if !inserted {
		pq.items = append(pq.items, newItem)
	}
}

// Dequeue 出队
func (pq *PriorityQueue[T]) Dequeue() (T, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(pq.items) == 0 {
		var zero T
		return zero, false
	}
	item := pq.items[0]
	pq.items = pq.items[1:]
	return item.value, true
}

// Peek 查看队首元素
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	if len(pq.items) == 0 {
		var zero T
		return zero, false
	}
	return pq.items[0].value, true
}

// Size 获取队列大小
func (pq *PriorityQueue[T]) Size() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return len(pq.items)
}

// IsEmpty 检查队列是否为空
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.Size() == 0
}

// Clear 清空队列
func (pq *PriorityQueue[T]) Clear() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.items = pq.items[:0]
}

// CircularQueue 循环队列
type CircularQueue[T any] struct {
	items []T
	front int
	rear  int
	size  int
	cap   int
	mu    sync.RWMutex
}

// NewCircularQueue 创建循环队列
func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	return &CircularQueue[T]{
		items: make([]T, capacity),
		front: 0,
		rear:  0,
		size:  0,
		cap:   capacity,
	}
}

// Enqueue 入队
func (cq *CircularQueue[T]) Enqueue(item T) bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.size == cq.cap {
		return false // 队列已满
	}
	cq.items[cq.rear] = item
	cq.rear = (cq.rear + 1) % cq.cap
	cq.size++
	return true
}

// Dequeue 出队
func (cq *CircularQueue[T]) Dequeue() (T, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if cq.size == 0 {
		var zero T
		return zero, false
	}
	item := cq.items[cq.front]
	cq.front = (cq.front + 1) % cq.cap
	cq.size--
	return item, true
}

// Peek 查看队首元素
func (cq *CircularQueue[T]) Peek() (T, bool) {
	cq.mu.RLock()
	defer cq.mu.RUnlock()
	if cq.size == 0 {
		var zero T
		return zero, false
	}
	return cq.items[cq.front], true
}

// Size 获取队列大小
func (cq *CircularQueue[T]) Size() int {
	cq.mu.RLock()
	defer cq.mu.RUnlock()
	return cq.size
}

// IsEmpty 检查队列是否为空
func (cq *CircularQueue[T]) IsEmpty() bool {
	return cq.Size() == 0
}

// IsFull 检查队列是否已满
func (cq *CircularQueue[T]) IsFull() bool {
	cq.mu.RLock()
	defer cq.mu.RUnlock()
	return cq.size == cq.cap
}

// Clear 清空队列
func (cq *CircularQueue[T]) Clear() {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	cq.front = 0
	cq.rear = 0
	cq.size = 0
}

// Capacity 获取队列容量
func (cq *CircularQueue[T]) Capacity() int {
	return cq.cap
}

