package stack

import (
	"sync"
)

// Stack 栈接口
type Stack[T any] interface {
	Push(item T)
	Pop() (T, bool)
	Peek() (T, bool)
	Size() int
	IsEmpty() bool
	Clear()
}

// SimpleStack 简单栈实现
type SimpleStack[T any] struct {
	items []T
	mu    sync.RWMutex
}

// NewSimpleStack 创建简单栈
func NewSimpleStack[T any]() *SimpleStack[T] {
	return &SimpleStack[T]{
		items: make([]T, 0),
	}
}

// Push 入栈
func (s *SimpleStack[T]) Push(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Pop 出栈
func (s *SimpleStack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	lastIndex := len(s.items) - 1
	item := s.items[lastIndex]
	s.items = s.items[:lastIndex]
	return item, true
}

// Peek 查看栈顶元素
func (s *SimpleStack[T]) Peek() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Size 获取栈大小
func (s *SimpleStack[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// IsEmpty 检查栈是否为空
func (s *SimpleStack[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Clear 清空栈
func (s *SimpleStack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = s.items[:0]
}

// ToSlice 转换为切片
func (s *SimpleStack[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result
}

// MaxStack 最大栈（支持O(1)获取最大值）
type MaxStack[T comparable] struct {
	stack    []T
	maxStack []T
	mu       sync.RWMutex
	compare  func(a, b T) bool // 返回true表示a > b
}

// NewMaxStack 创建最大栈
func NewMaxStack[T comparable](compare func(a, b T) bool) *MaxStack[T] {
	return &MaxStack[T]{
		stack:    make([]T, 0),
		maxStack: make([]T, 0),
		compare:  compare,
	}
}

// Push 入栈
func (ms *MaxStack[T]) Push(item T) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.stack = append(ms.stack, item)
	if len(ms.maxStack) == 0 || ms.compare(item, ms.maxStack[len(ms.maxStack)-1]) {
		ms.maxStack = append(ms.maxStack, item)
	} else {
		ms.maxStack = append(ms.maxStack, ms.maxStack[len(ms.maxStack)-1])
	}
}

// Pop 出栈
func (ms *MaxStack[T]) Pop() (T, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if len(ms.stack) == 0 {
		var zero T
		return zero, false
	}
	lastIndex := len(ms.stack) - 1
	item := ms.stack[lastIndex]
	ms.stack = ms.stack[:lastIndex]
	ms.maxStack = ms.maxStack[:lastIndex]
	return item, true
}

// Peek 查看栈顶元素
func (ms *MaxStack[T]) Peek() (T, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if len(ms.stack) == 0 {
		var zero T
		return zero, false
	}
	return ms.stack[len(ms.stack)-1], true
}

// Max 获取最大值
func (ms *MaxStack[T]) Max() (T, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if len(ms.maxStack) == 0 {
		var zero T
		return zero, false
	}
	return ms.maxStack[len(ms.maxStack)-1], true
}

// Size 获取栈大小
func (ms *MaxStack[T]) Size() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.stack)
}

// IsEmpty 检查栈是否为空
func (ms *MaxStack[T]) IsEmpty() bool {
	return ms.Size() == 0
}

// Clear 清空栈
func (ms *MaxStack[T]) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.stack = ms.stack[:0]
	ms.maxStack = ms.maxStack[:0]
}

// MinStack 最小栈（支持O(1)获取最小值）
type MinStack[T comparable] struct {
	stack    []T
	minStack []T
	mu       sync.RWMutex
	compare  func(a, b T) bool // 返回true表示a < b
}

// NewMinStack 创建最小栈
func NewMinStack[T comparable](compare func(a, b T) bool) *MinStack[T] {
	return &MinStack[T]{
		stack:    make([]T, 0),
		minStack: make([]T, 0),
		compare:  compare,
	}
}

// Push 入栈
func (ms *MinStack[T]) Push(item T) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.stack = append(ms.stack, item)
	if len(ms.minStack) == 0 || ms.compare(item, ms.minStack[len(ms.minStack)-1]) {
		ms.minStack = append(ms.minStack, item)
	} else {
		ms.minStack = append(ms.minStack, ms.minStack[len(ms.minStack)-1])
	}
}

// Pop 出栈
func (ms *MinStack[T]) Pop() (T, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if len(ms.stack) == 0 {
		var zero T
		return zero, false
	}
	lastIndex := len(ms.stack) - 1
	item := ms.stack[lastIndex]
	ms.stack = ms.stack[:lastIndex]
	ms.minStack = ms.minStack[:lastIndex]
	return item, true
}

// Peek 查看栈顶元素
func (ms *MinStack[T]) Peek() (T, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if len(ms.stack) == 0 {
		var zero T
		return zero, false
	}
	return ms.stack[len(ms.stack)-1], true
}

// Min 获取最小值
func (ms *MinStack[T]) Min() (T, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if len(ms.minStack) == 0 {
		var zero T
		return zero, false
	}
	return ms.minStack[len(ms.minStack)-1], true
}

// Size 获取栈大小
func (ms *MinStack[T]) Size() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.stack)
}

// IsEmpty 检查栈是否为空
func (ms *MinStack[T]) IsEmpty() bool {
	return ms.Size() == 0
}

// Clear 清空栈
func (ms *MinStack[T]) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.stack = ms.stack[:0]
	ms.minStack = ms.minStack[:0]
}

