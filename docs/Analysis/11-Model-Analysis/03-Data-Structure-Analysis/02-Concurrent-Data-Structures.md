# 并发数据结构分析

## 目录

1. [概述](#概述)
2. [理论基础](#理论基础)
3. [并发队列](#并发队列)
4. [并发映射](#并发映射)
5. [无锁数据结构](#无锁数据结构)
6. [原子操作](#原子操作)
7. [性能分析](#性能分析)
8. [应用场景](#应用场景)
9. [最佳实践](#最佳实践)

## 概述

并发数据结构是设计用于在多线程环境中安全使用的数据结构。在Golang中，由于goroutine的轻量级特性，并发数据结构的设计尤为重要。

### 核心挑战

1. **数据竞争**: 多个goroutine同时访问共享数据
2. **死锁**: 多个锁的循环等待
3. **活锁**: 线程不断重试但无法取得进展
4. **性能开销**: 锁机制带来的性能损失

### 设计原则

- **线程安全**: 保证多线程访问的正确性
- **高性能**: 最小化锁竞争和同步开销
- **可扩展性**: 支持高并发访问
- **公平性**: 避免饥饿现象

## 理论基础

### 1. 并发控制理论

#### 1.1 临界区问题

**定义 1.1 (临界区)** 临界区是访问共享资源的代码段，必须满足：

1. **互斥性**: 同一时刻只能有一个线程在临界区内
2. **进展性**: 如果没有线程在临界区内，那么想要进入临界区的线程应该能够进入
3. **有限等待**: 一个线程等待进入临界区的时间应该是有限的

#### 1.2 同步原语

**定义 1.2 (互斥锁)** 互斥锁是一个二元状态变量，支持两个原子操作：

- **Lock()**: 如果锁是自由的，则获取锁；否则阻塞
- **Unlock()**: 释放锁

**定义 1.3 (读写锁)** 读写锁允许多个读者同时访问，但写者必须独占访问：

- **RLock()**: 获取读锁
- **RUnlock()**: 释放读锁
- **Lock()**: 获取写锁
- **Unlock()**: 释放写锁

### 2. 内存模型

#### 2.1 内存序

**定义 1.4 (内存序)** 内存序定义了内存操作的可见性顺序：

- **Relaxed**: 最弱的内存序，只保证原子性
- **Acquire**: 保证后续操作的可见性
- **Release**: 保证前面操作的可见性
- **Acquire-Release**: 同时保证acquire和release语义
- **Sequentially Consistent**: 最强的内存序，全局顺序

#### 2.2 原子操作

**定义 1.5 (原子操作)** 原子操作是不可分割的操作，要么完全执行，要么完全不执行。

在Golang中，原子操作通过`sync/atomic`包提供：

```go
// 原子操作示例
var value int64

// 原子加载
loaded := atomic.LoadInt64(&value)

// 原子存储
atomic.StoreInt64(&value, 42)

// 原子比较并交换
swapped := atomic.CompareAndSwapInt64(&value, 42, 100)

// 原子加法
newValue := atomic.AddInt64(&value, 10)
```

### 3. 无锁编程

#### 3.1 无锁数据结构定义

**定义 1.6 (无锁数据结构)** 无锁数据结构是不使用互斥锁的并发数据结构，通过原子操作实现同步。

**定义 1.7 (无等待数据结构)** 无等待数据结构是无锁数据结构的一种，保证每个操作都能在有限步数内完成。

#### 3.2 无锁算法的正确性

**定理 1.1 (线性化)** 无锁算法的执行历史是线性化的，即存在一个全局顺序，使得每个操作看起来都在某个瞬间原子地执行。

**证明**: 通过构造线性化点来证明。对于每个操作，选择一个线性化点，使得在该点之前的所有操作都已完成，之后的操作都未开始。

## 并发队列

### 1. 基于锁的并发队列

#### 1.1 互斥锁队列

```go
// MutexQueue 基于互斥锁的并发队列
type MutexQueue[T any] struct {
    elements []T
    head     int
    tail     int
    size     int
    capacity int
    mutex    sync.Mutex
}

// NewMutexQueue 创建新的互斥锁队列
func NewMutexQueue[T any](capacity int) *MutexQueue[T] {
    return &MutexQueue[T]{
        elements: make([]T, capacity),
        head:     0,
        tail:     -1,
        size:     0,
        capacity: capacity,
    }
}

// Enqueue 入队操作
func (mq *MutexQueue[T]) Enqueue(element T) error {
    mq.mutex.Lock()
    defer mq.mutex.Unlock()
    
    if mq.size >= mq.capacity {
        return errors.New("queue is full")
    }
    
    mq.tail = (mq.tail + 1) % mq.capacity
    mq.elements[mq.tail] = element
    mq.size++
    return nil
}

// Dequeue 出队操作
func (mq *MutexQueue[T]) Dequeue() (T, error) {
    mq.mutex.Lock()
    defer mq.mutex.Unlock()
    
    var zero T
    if mq.size == 0 {
        return zero, errors.New("queue is empty")
    }
    
    element := mq.elements[mq.head]
    mq.head = (mq.head + 1) % mq.capacity
    mq.size--
    return element, nil
}

// Size 获取队列大小
func (mq *MutexQueue[T]) Size() int {
    mq.mutex.Lock()
    defer mq.mutex.Unlock()
    return mq.size
}

// IsEmpty 检查是否为空
func (mq *MutexQueue[T]) IsEmpty() bool {
    return mq.Size() == 0
}

// IsFull 检查是否已满
func (mq *MutexQueue[T]) IsFull() bool {
    mq.mutex.Lock()
    defer mq.mutex.Unlock()
    return mq.size >= mq.capacity
}
```

#### 1.2 读写锁队列

```go
// RWLockQueue 基于读写锁的并发队列
type RWLockQueue[T any] struct {
    elements []T
    head     int
    tail     int
    size     int
    capacity int
    mutex    sync.RWMutex
}

// NewRWLockQueue 创建新的读写锁队列
func NewRWLockQueue[T any](capacity int) *RWLockQueue[T] {
    return &RWLockQueue[T]{
        elements: make([]T, capacity),
        head:     0,
        tail:     -1,
        size:     0,
        capacity: capacity,
    }
}

// Enqueue 入队操作（写操作）
func (rwq *RWLockQueue[T]) Enqueue(element T) error {
    rwq.mutex.Lock()
    defer rwq.mutex.Unlock()
    
    if rwq.size >= rwq.capacity {
        return errors.New("queue is full")
    }
    
    rwq.tail = (rwq.tail + 1) % rwq.capacity
    rwq.elements[rwq.tail] = element
    rwq.size++
    return nil
}

// Dequeue 出队操作（写操作）
func (rwq *RWLockQueue[T]) Dequeue() (T, error) {
    rwq.mutex.Lock()
    defer rwq.mutex.Unlock()
    
    var zero T
    if rwq.size == 0 {
        return zero, errors.New("queue is empty")
    }
    
    element := rwq.elements[rwq.head]
    rwq.head = (rwq.head + 1) % rwq.capacity
    rwq.size--
    return element, nil
}

// Front 查看队首元素（读操作）
func (rwq *RWLockQueue[T]) Front() (T, error) {
    rwq.mutex.RLock()
    defer rwq.mutex.RUnlock()
    
    var zero T
    if rwq.size == 0 {
        return zero, errors.New("queue is empty")
    }
    
    return rwq.elements[rwq.head], nil
}

// Size 获取队列大小（读操作）
func (rwq *RWLockQueue[T]) Size() int {
    rwq.mutex.RLock()
    defer rwq.mutex.RUnlock()
    return rwq.size
}

// IsEmpty 检查是否为空（读操作）
func (rwq *RWLockQueue[T]) IsEmpty() bool {
    return rwq.Size() == 0
}
```

### 2. 基于通道的并发队列

#### 2.1 缓冲通道队列

```go
// ChannelQueue 基于通道的并发队列
type ChannelQueue[T any] struct {
    elements chan T
    capacity int
}

// NewChannelQueue 创建新的通道队列
func NewChannelQueue[T any](capacity int) *ChannelQueue[T] {
    return &ChannelQueue[T]{
        elements: make(chan T, capacity),
        capacity: capacity,
    }
}

// Enqueue 入队操作
func (cq *ChannelQueue[T]) Enqueue(element T) error {
    select {
    case cq.elements <- element:
        return nil
    default:
        return errors.New("queue is full")
    }
}

// Dequeue 出队操作
func (cq *ChannelQueue[T]) Dequeue() (T, error) {
    select {
    case element := <-cq.elements:
        return element, nil
    default:
        var zero T
        return zero, errors.New("queue is empty")
    }
}

// TryEnqueue 尝试入队（非阻塞）
func (cq *ChannelQueue[T]) TryEnqueue(element T) bool {
    select {
    case cq.elements <- element:
        return true
    default:
        return false
    }
}

// TryDequeue 尝试出队（非阻塞）
func (cq *ChannelQueue[T]) TryDequeue() (T, bool) {
    select {
    case element := <-cq.elements:
        return element, true
    default:
        var zero T
        return zero, false
    }
}

// Size 获取队列大小
func (cq *ChannelQueue[T]) Size() int {
    return len(cq.elements)
}

// Capacity 获取队列容量
func (cq *ChannelQueue[T]) Capacity() int {
    return cq.capacity
}

// IsEmpty 检查是否为空
func (cq *ChannelQueue[T]) IsEmpty() bool {
    return len(cq.elements) == 0
}

// IsFull 检查是否已满
func (cq *ChannelQueue[T]) IsFull() bool {
    return len(cq.elements) == cq.capacity
}
```

#### 2.2 优先级通道队列

```go
// PriorityChannelQueue 优先级通道队列
type PriorityChannelQueue[T any] struct {
    highPriority chan T
    lowPriority  chan T
    capacity     int
}

// NewPriorityChannelQueue 创建优先级通道队列
func NewPriorityChannelQueue[T any](capacity int) *PriorityChannelQueue[T] {
    return &PriorityChannelQueue[T]{
        highPriority: make(chan T, capacity),
        lowPriority:  make(chan T, capacity),
        capacity:     capacity,
    }
}

// EnqueueHigh 高优先级入队
func (pcq *PriorityChannelQueue[T]) EnqueueHigh(element T) error {
    select {
    case pcq.highPriority <- element:
        return nil
    default:
        return errors.New("high priority queue is full")
    }
}

// EnqueueLow 低优先级入队
func (pcq *PriorityChannelQueue[T]) EnqueueLow(element T) error {
    select {
    case pcq.lowPriority <- element:
        return nil
    default:
        return errors.New("low priority queue is full")
    }
}

// Dequeue 出队（优先处理高优先级）
func (pcq *PriorityChannelQueue[T]) Dequeue() (T, error) {
    // 优先处理高优先级队列
    select {
    case element := <-pcq.highPriority:
        return element, nil
    default:
        // 高优先级队列为空，处理低优先级队列
        select {
        case element := <-pcq.lowPriority:
            return element, nil
        default:
            var zero T
            return zero, errors.New("both queues are empty")
        }
    }
}

// Size 获取总大小
func (pcq *PriorityChannelQueue[T]) Size() int {
    return len(pcq.highPriority) + len(pcq.lowPriority)
}

// HighPrioritySize 获取高优先级队列大小
func (pcq *PriorityChannelQueue[T]) HighPrioritySize() int {
    return len(pcq.highPriority)
}

// LowPrioritySize 获取低优先级队列大小
func (pcq *PriorityChannelQueue[T]) LowPrioritySize() int {
    return len(pcq.lowPriority)
}
```

## 并发映射

### 1. 基于锁的并发映射

#### 1.1 单一锁映射

```go
// MutexMap 基于单一锁的并发映射
type MutexMap[K comparable, V any] struct {
    data  map[K]V
    mutex sync.RWMutex
}

// NewMutexMap 创建新的互斥锁映射
func NewMutexMap[K comparable, V any]() *MutexMap[K, V] {
    return &MutexMap[K, V]{
        data: make(map[K]V),
    }
}

// Set 设置键值对
func (mm *MutexMap[K, V]) Set(key K, value V) {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()
    mm.data[key] = value
}

// Get 获取值
func (mm *MutexMap[K, V]) Get(key K) (V, bool) {
    mm.mutex.RLock()
    defer mm.mutex.RUnlock()
    value, exists := mm.data[key]
    return value, exists
}

// Delete 删除键值对
func (mm *MutexMap[K, V]) Delete(key K) {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()
    delete(mm.data, key)
}

// Has 检查键是否存在
func (mm *MutexMap[K, V]) Has(key K) bool {
    mm.mutex.RLock()
    defer mm.mutex.RUnlock()
    _, exists := mm.data[key]
    return exists
}

// Size 获取映射大小
func (mm *MutexMap[K, V]) Size() int {
    mm.mutex.RLock()
    defer mm.mutex.RUnlock()
    return len(mm.data)
}

// Clear 清空映射
func (mm *MutexMap[K, V]) Clear() {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()
    mm.data = make(map[K]V)
}

// Keys 获取所有键
func (mm *MutexMap[K, V]) Keys() []K {
    mm.mutex.RLock()
    defer mm.mutex.RUnlock()
    
    keys := make([]K, 0, len(mm.data))
    for key := range mm.data {
        keys = append(keys, key)
    }
    return keys
}

// Values 获取所有值
func (mm *MutexMap[K, V]) Values() []V {
    mm.mutex.RLock()
    defer mm.mutex.RUnlock()
    
    values := make([]V, 0, len(mm.data))
    for _, value := range mm.data {
        values = append(values, value)
    }
    return values
}
```

#### 1.2 分段锁映射

```go
// SegmentMap 基于分段锁的并发映射
type SegmentMap[K comparable, V any] struct {
    segments []*segment[K, V]
    size     int
}

type segment[K comparable, V any] struct {
    data  map[K]V
    mutex sync.RWMutex
}

// NewSegmentMap 创建新的分段锁映射
func NewSegmentMap[K comparable, V any](segmentCount int) *SegmentMap[K, V] {
    segments := make([]*segment[K, V], segmentCount)
    for i := 0; i < segmentCount; i++ {
        segments[i] = &segment[K, V]{
            data: make(map[K]V),
        }
    }
    
    return &SegmentMap[K, V]{
        segments: segments,
        size:     segmentCount,
    }
}

// getSegment 获取键对应的段
func (sm *SegmentMap[K, V]) getSegment(key K) *segment[K, V] {
    hash := sm.hash(key)
    return sm.segments[hash%uint32(sm.size)]
}

// hash 计算键的哈希值
func (sm *SegmentMap[K, V]) hash(key K) uint32 {
    // 简单的哈希函数，实际应用中可以使用更复杂的哈希算法
    keyStr := fmt.Sprintf("%v", key)
    h := fnv.New32a()
    h.Write([]byte(keyStr))
    return h.Sum32()
}

// Set 设置键值对
func (sm *SegmentMap[K, V]) Set(key K, value V) {
    segment := sm.getSegment(key)
    segment.mutex.Lock()
    defer segment.mutex.Unlock()
    segment.data[key] = value
}

// Get 获取值
func (sm *SegmentMap[K, V]) Get(key K) (V, bool) {
    segment := sm.getSegment(key)
    segment.mutex.RLock()
    defer segment.mutex.RUnlock()
    value, exists := segment.data[key]
    return value, exists
}

// Delete 删除键值对
func (sm *SegmentMap[K, V]) Delete(key K) {
    segment := sm.getSegment(key)
    segment.mutex.Lock()
    defer segment.mutex.Unlock()
    delete(segment.data, key)
}

// Has 检查键是否存在
func (sm *SegmentMap[K, V]) Has(key K) bool {
    segment := sm.getSegment(key)
    segment.mutex.RLock()
    defer segment.mutex.RUnlock()
    _, exists := segment.data[key]
    return exists
}

// Size 获取映射总大小
func (sm *SegmentMap[K, V]) Size() int {
    total := 0
    for _, segment := range sm.segments {
        segment.mutex.RLock()
        total += len(segment.data)
        segment.mutex.RUnlock()
    }
    return total
}
```

### 2. 基于原子操作的并发映射

#### 2.1 原子指针映射

```go
// AtomicMap 基于原子指针的并发映射
type AtomicMap[K comparable, V any] struct {
    data unsafe.Pointer // *map[K]V
}

// NewAtomicMap 创建新的原子映射
func NewAtomicMap[K comparable, V any]() *AtomicMap[K, V] {
    m := make(map[K]V)
    return &AtomicMap[K, V]{
        data: unsafe.Pointer(&m),
    }
}

// Set 设置键值对
func (am *AtomicMap[K, V]) Set(key K, value V) {
    for {
        oldData := atomic.LoadPointer(&am.data)
        oldMap := *(*map[K]V)(oldData)
        
        // 创建新映射
        newMap := make(map[K]V, len(oldMap)+1)
        for k, v := range oldMap {
            newMap[k] = v
        }
        newMap[key] = value
        
        // 原子交换
        if atomic.CompareAndSwapPointer(&am.data, oldData, unsafe.Pointer(&newMap)) {
            break
        }
    }
}

// Get 获取值
func (am *AtomicMap[K, V]) Get(key K) (V, bool) {
    data := atomic.LoadPointer(&am.data)
    m := *(*map[K]V)(data)
    value, exists := m[key]
    return value, exists
}

// Delete 删除键值对
func (am *AtomicMap[K, V]) Delete(key K) {
    for {
        oldData := atomic.LoadPointer(&am.data)
        oldMap := *(*map[K]V)(oldData)
        
        if _, exists := oldMap[key]; !exists {
            break // 键不存在，无需删除
        }
        
        // 创建新映射
        newMap := make(map[K]V, len(oldMap)-1)
        for k, v := range oldMap {
            if k != key {
                newMap[k] = v
            }
        }
        
        // 原子交换
        if atomic.CompareAndSwapPointer(&am.data, oldData, unsafe.Pointer(&newMap)) {
            break
        }
    }
}
```

## 无锁数据结构

### 1. 无锁栈

#### 1.1 基于CAS的无锁栈

```go
// LockFreeStack 无锁栈
type LockFreeStack[T any] struct {
    head unsafe.Pointer // *node[T]
}

type node[T any] struct {
    value T
    next  unsafe.Pointer // *node[T]
}

// NewLockFreeStack 创建新的无锁栈
func NewLockFreeStack[T any]() *LockFreeStack[T] {
    return &LockFreeStack[T]{}
}

// Push 入栈
func (lfs *LockFreeStack[T]) Push(value T) {
    newNode := &node[T]{value: value}
    
    for {
        oldHead := atomic.LoadPointer(&lfs.head)
        newNode.next = oldHead
        
        if atomic.CompareAndSwapPointer(&lfs.head, oldHead, unsafe.Pointer(newNode)) {
            break
        }
    }
}

// Pop 出栈
func (lfs *LockFreeStack[T]) Pop() (T, bool) {
    for {
        oldHead := atomic.LoadPointer(&lfs.head)
        if oldHead == nil {
            var zero T
            return zero, false
        }
        
        headNode := (*node[T])(oldHead)
        newHead := headNode.next
        
        if atomic.CompareAndSwapPointer(&lfs.head, oldHead, newHead) {
            return headNode.value, true
        }
    }
}

// IsEmpty 检查是否为空
func (lfs *LockFreeStack[T]) IsEmpty() bool {
    return atomic.LoadPointer(&lfs.head) == nil
}
```

#### 1.2 基于ABA问题的解决方案

```go
// ABAFreeStack ABA问题解决方案
type ABAFreeStack[T any] struct {
    head unsafe.Pointer // *taggedNode[T]
}

type taggedNode[T any] struct {
    node *node[T]
    tag  uint64
}

// NewABAFreeStack 创建新的ABA安全栈
func NewABAFreeStack[T any]() *ABAFreeStack[T] {
    return &ABAFreeStack[T]{}
}

// Push 入栈
func (afs *ABAFreeStack[T]) Push(value T) {
    newNode := &node[T]{value: value}
    
    for {
        oldHead := atomic.LoadPointer(&afs.head)
        oldTagged := (*taggedNode[T])(oldHead)
        
        var oldNode *node[T]
        var oldTag uint64
        if oldTagged != nil {
            oldNode = oldTagged.node
            oldTag = oldTagged.tag
        }
        
        newNode.next = unsafe.Pointer(oldNode)
        newTagged := &taggedNode[T]{
            node: newNode,
            tag:  oldTag + 1,
        }
        
        if atomic.CompareAndSwapPointer(&afs.head, oldHead, unsafe.Pointer(newTagged)) {
            break
        }
    }
}

// Pop 出栈
func (afs *ABAFreeStack[T]) Pop() (T, bool) {
    for {
        oldHead := atomic.LoadPointer(&afs.head)
        if oldHead == nil {
            var zero T
            return zero, false
        }
        
        oldTagged := (*taggedNode[T])(oldHead)
        oldNode := oldTagged.node
        oldTag := oldTagged.tag
        
        newHead := oldNode.next
        var newTagged *taggedNode[T]
        if newHead != nil {
            newTagged = &taggedNode[T]{
                node: (*node[T])(newHead),
                tag:  oldTag + 1,
            }
        }
        
        if atomic.CompareAndSwapPointer(&afs.head, oldHead, unsafe.Pointer(newTagged)) {
            return oldNode.value, true
        }
    }
}
```

### 2. 无锁队列

#### 2.1 基于CAS的无锁队列

```go
// LockFreeQueue 无锁队列
type LockFreeQueue[T any] struct {
    head unsafe.Pointer // *node[T]
    tail unsafe.Pointer // *node[T]
}

// NewLockFreeQueue 创建新的无锁队列
func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
    dummy := &node[T]{}
    return &LockFreeQueue[T]{
        head: unsafe.Pointer(dummy),
        tail: unsafe.Pointer(dummy),
    }
}

// Enqueue 入队
func (lfq *LockFreeQueue[T]) Enqueue(value T) {
    newNode := &node[T]{value: value}
    
    for {
        tail := atomic.LoadPointer(&lfq.tail)
        tailNode := (*node[T])(tail)
        
        // 尝试更新尾节点的next指针
        if atomic.CompareAndSwapPointer(&tailNode.next, nil, unsafe.Pointer(newNode)) {
            // 更新尾指针
            atomic.CompareAndSwapPointer(&lfq.tail, tail, unsafe.Pointer(newNode))
            break
        } else {
            // 帮助其他线程更新尾指针
            atomic.CompareAndSwapPointer(&lfq.tail, tail, tailNode.next)
        }
    }
}

// Dequeue 出队
func (lfq *LockFreeQueue[T]) Dequeue() (T, bool) {
    for {
        head := atomic.LoadPointer(&lfq.head)
        tail := atomic.LoadPointer(&lfq.tail)
        headNode := (*node[T])(head)
        next := atomic.LoadPointer(&headNode.next)
        
        if head == tail {
            if next == nil {
                var zero T
                return zero, false
            }
            // 帮助其他线程更新尾指针
            atomic.CompareAndSwapPointer(&lfq.tail, tail, next)
        } else {
            if next == nil {
                continue
            }
            
            nextNode := (*node[T])(next)
            if atomic.CompareAndSwapPointer(&lfq.head, head, next) {
                return nextNode.value, true
            }
        }
    }
}
```

## 原子操作

### 1. 基础原子操作

```go
// AtomicCounter 原子计数器
type AtomicCounter struct {
    value int64
}

// NewAtomicCounter 创建新的原子计数器
func NewAtomicCounter() *AtomicCounter {
    return &AtomicCounter{}
}

// Increment 增加计数
func (ac *AtomicCounter) Increment() int64 {
    return atomic.AddInt64(&ac.value, 1)
}

// Decrement 减少计数
func (ac *AtomicCounter) Decrement() int64 {
    return atomic.AddInt64(&ac.value, -1)
}

// Get 获取当前值
func (ac *AtomicCounter) Get() int64 {
    return atomic.LoadInt64(&ac.value)
}

// Set 设置值
func (ac *AtomicCounter) Set(value int64) {
    atomic.StoreInt64(&ac.value, value)
}

// CompareAndSet 比较并设置
func (ac *AtomicCounter) CompareAndSet(expected, new int64) bool {
    return atomic.CompareAndSwapInt64(&ac.value, expected, new)
}
```

### 2. 原子引用

```go
// AtomicReference 原子引用
type AtomicReference[T any] struct {
    value unsafe.Pointer
}

// NewAtomicReference 创建新的原子引用
func NewAtomicReference[T any]() *AtomicReference[T] {
    return &AtomicReference[T]{}
}

// Set 设置值
func (ar *AtomicReference[T]) Set(value T) {
    atomic.StorePointer(&ar.value, unsafe.Pointer(&value))
}

// Get 获取值
func (ar *AtomicReference[T]) Get() (T, bool) {
    ptr := atomic.LoadPointer(&ar.value)
    if ptr == nil {
        var zero T
        return zero, false
    }
    return *(*T)(ptr), true
}

// CompareAndSet 比较并设置
func (ar *AtomicReference[T]) CompareAndSet(expected, new T) bool {
    var expectedPtr unsafe.Pointer
    if !reflect.ValueOf(expected).IsZero() {
        expectedPtr = unsafe.Pointer(&expected)
    }
    
    var newPtr unsafe.Pointer
    if !reflect.ValueOf(new).IsZero() {
        newPtr = unsafe.Pointer(&new)
    }
    
    return atomic.CompareAndSwapPointer(&ar.value, expectedPtr, newPtr)
}
```

## 性能分析

### 1. 性能对比

#### 1.1 队列性能对比

```go
// 性能基准测试
func BenchmarkMutexQueue(b *testing.B) {
    queue := NewMutexQueue[int](1000)
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            if i%2 == 0 {
                queue.Enqueue(i)
            } else {
                queue.Dequeue()
            }
            i++
        }
    })
}

func BenchmarkChannelQueue(b *testing.B) {
    queue := NewChannelQueue[int](1000)
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            if i%2 == 0 {
                queue.TryEnqueue(i)
            } else {
                queue.TryDequeue()
            }
            i++
        }
    })
}

func BenchmarkLockFreeQueue(b *testing.B) {
    queue := NewLockFreeQueue[int]()
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            if i%2 == 0 {
                queue.Enqueue(i)
            } else {
                queue.Dequeue()
            }
            i++
        }
    })
}
```

#### 1.2 映射性能对比

```go
func BenchmarkMutexMap(b *testing.B) {
    m := NewMutexMap[int, string]()
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            if i%3 == 0 {
                m.Set(i, fmt.Sprintf("value_%d", i))
            } else if i%3 == 1 {
                m.Get(i)
            } else {
                m.Delete(i)
            }
            i++
        }
    })
}

func BenchmarkSegmentMap(b *testing.B) {
    m := NewSegmentMap[int, string](16)
    b.ResetTimer()
    
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            if i%3 == 0 {
                m.Set(i, fmt.Sprintf("value_%d", i))
            } else if i%3 == 1 {
                m.Get(i)
            } else {
                m.Delete(i)
            }
            i++
        }
    })
}
```

### 2. 内存使用分析

#### 2.1 内存开销对比

| 数据结构 | 每个元素开销 | 锁开销 | 总开销 |
|---------|-------------|--------|--------|
| 互斥锁队列 | 1个元素 | 1个锁 | 低 |
| 通道队列 | 1个元素 | 无锁 | 中等 |
| 无锁队列 | 1个元素+指针 | 无锁 | 中等 |
| 互斥锁映射 | 1个键值对 | 1个锁 | 低 |
| 分段锁映射 | 1个键值对 | N个锁 | 中等 |
| 原子映射 | 1个键值对 | 无锁 | 高 |

#### 2.2 缓存性能分析

```go
// 缓存友好的数据结构设计
type CacheFriendlyQueue[T any] struct {
    elements []T
    head     int
    tail     int
    size     int
    capacity int
    mutex    sync.Mutex
    // 填充到64字节边界，避免伪共享
    _ [64 - 8*5]byte
}

// 避免伪共享的设计
type PaddedMutex struct {
    mutex sync.Mutex
    _     [64 - 8]byte // 填充到64字节
}

type CacheFriendlyMap[K comparable, V any] struct {
    segments []*paddedSegment[K, V]
    size     int
}

type paddedSegment[K comparable, V any] struct {
    data  map[K]V
    mutex PaddedMutex
}
```

## 应用场景

### 1. 高并发服务器

```go
// 高并发HTTP服务器
type ConcurrentServer struct {
    requestQueue *ChannelQueue[*http.Request]
    workerPool   *WorkerPool
}

type WorkerPool struct {
    workers int
    queue   *ChannelQueue[func()]
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go func() {
            for task := range wp.queue.elements {
                task()
            }
        }()
    }
}

func (wp *WorkerPool) Submit(task func()) error {
    return wp.queue.Enqueue(task)
}
```

### 2. 缓存系统

```go
// 并发缓存系统
type ConcurrentCache[K comparable, V any] struct {
    data *SegmentMap[K, V]
    ttl  time.Duration
}

type CacheEntry[V any] struct {
    Value      V
    ExpireTime time.Time
}

func (cc *ConcurrentCache[K, V]) Set(key K, value V) {
    entry := CacheEntry[V]{
        Value:      value,
        ExpireTime: time.Now().Add(cc.ttl),
    }
    cc.data.Set(key, entry)
}

func (cc *ConcurrentCache[K, V]) Get(key K) (V, bool) {
    if entry, exists := cc.data.Get(key); exists {
        if time.Now().Before(entry.ExpireTime) {
            return entry.Value, true
        } else {
            cc.data.Delete(key)
        }
    }
    var zero V
    return zero, false
}
```

### 3. 任务调度器

```go
// 并发任务调度器
type TaskScheduler struct {
    highPriorityQueue *PriorityChannelQueue[*Task]
    lowPriorityQueue  *PriorityChannelQueue[*Task]
    workers           int
    stopChan          chan struct{}
}

type Task struct {
    ID       string
    Priority int
    Handler  func() error
}

func (ts *TaskScheduler) Start() {
    for i := 0; i < ts.workers; i++ {
        go ts.worker()
    }
}

func (ts *TaskScheduler) worker() {
    for {
        select {
        case <-ts.stopChan:
            return
        default:
            // 优先处理高优先级任务
            if task, err := ts.highPriorityQueue.Dequeue(); err == nil {
                task.Handler()
            } else if task, err := ts.lowPriorityQueue.Dequeue(); err == nil {
                task.Handler()
            } else {
                time.Sleep(time.Millisecond)
            }
        }
    }
}

func (ts *TaskScheduler) Submit(task *Task) error {
    if task.Priority > 5 {
        return ts.highPriorityQueue.EnqueueHigh(task)
    } else {
        return ts.lowPriorityQueue.EnqueueLow(task)
    }
}
```

## 最佳实践

### 1. 选择合适的数据结构

#### 1.1 根据并发度选择

```go
// 低并发场景 - 使用简单锁
func LowConcurrencyExample() {
    queue := NewMutexQueue[int](100)
    // 适合并发度 < 10
}

// 中等并发场景 - 使用分段锁
func MediumConcurrencyExample() {
    map := NewSegmentMap[int, string](16)
    // 适合并发度 10-100
}

// 高并发场景 - 使用无锁数据结构
func HighConcurrencyExample() {
    queue := NewLockFreeQueue[int]()
    // 适合并发度 > 100
}
```

#### 1.2 根据访问模式选择

```go
// 读多写少 - 使用读写锁
func ReadHeavyExample() {
    cache := &struct {
        data map[string]interface{}
        mutex sync.RWMutex
    }{
        data: make(map[string]interface{}),
    }
}

// 写多读少 - 使用互斥锁
func WriteHeavyExample() {
    queue := NewMutexQueue[int](1000)
}

// 读写均衡 - 使用无锁数据结构
func BalancedExample() {
    stack := NewLockFreeStack[int]()
}
```

### 2. 性能优化技巧

#### 2.1 减少锁竞争

```go
// 使用本地缓存减少锁竞争
type LocalCache[T any] struct {
    local  map[string]T
    global *MutexMap[string, T]
    mutex  sync.Mutex
}

func (lc *LocalCache[T]) Get(key string) (T, bool) {
    // 先检查本地缓存
    lc.mutex.Lock()
    if value, exists := lc.local[key]; exists {
        lc.mutex.Unlock()
        return value, true
    }
    lc.mutex.Unlock()
    
    // 检查全局缓存
    if value, exists := lc.global.Get(key); exists {
        // 更新本地缓存
        lc.mutex.Lock()
        lc.local[key] = value
        lc.mutex.Unlock()
        return value, true
    }
    
    var zero T
    return zero, false
}
```

#### 2.2 批量操作优化

```go
// 批量操作减少锁开销
type BatchQueue[T any] struct {
    elements []T
    mutex    sync.Mutex
    batchSize int
}

func (bq *BatchQueue[T]) EnqueueBatch(elements []T) {
    bq.mutex.Lock()
    defer bq.mutex.Unlock()
    bq.elements = append(bq.elements, elements...)
}

func (bq *BatchQueue[T]) DequeueBatch(count int) []T {
    bq.mutex.Lock()
    defer bq.mutex.Unlock()
    
    if count > len(bq.elements) {
        count = len(bq.elements)
    }
    
    result := make([]T, count)
    copy(result, bq.elements[:count])
    bq.elements = bq.elements[count:]
    return result
}
```

### 3. 错误处理

#### 3.1 超时处理

```go
// 带超时的操作
func (cq *ChannelQueue[T]) EnqueueWithTimeout(element T, timeout time.Duration) error {
    select {
    case cq.elements <- element:
        return nil
    case <-time.After(timeout):
        return errors.New("enqueue timeout")
    }
}

func (cq *ChannelQueue[T]) DequeueWithTimeout(timeout time.Duration) (T, error) {
    select {
    case element := <-cq.elements:
        return element, nil
    case <-time.After(timeout):
        var zero T
        return zero, errors.New("dequeue timeout")
    }
}
```

#### 3.2 优雅关闭

```go
// 优雅关闭的队列
type GracefulQueue[T any] struct {
    elements chan T
    closed   int32
}

func (gq *GracefulQueue[T]) Close() {
    atomic.StoreInt32(&gq.closed, 1)
    close(gq.elements)
}

func (gq *GracefulQueue[T]) Enqueue(element T) error {
    if atomic.LoadInt32(&gq.closed) == 1 {
        return errors.New("queue is closed")
    }
    
    select {
    case gq.elements <- element:
        return nil
    default:
        return errors.New("queue is full")
    }
}

func (gq *GracefulQueue[T]) Dequeue() (T, error) {
    element, ok := <-gq.elements
    if !ok {
        var zero T
        return zero, errors.New("queue is closed")
    }
    return element, nil
}
```

## 总结

并发数据结构是构建高并发系统的核心组件，选择合适的数据结构需要考虑：

1. **并发度**: 根据并发访问的线程数选择合适的数据结构
2. **访问模式**: 根据读写比例选择锁类型
3. **性能要求**: 根据延迟和吞吐量要求选择实现方式
4. **内存限制**: 考虑内存开销和缓存友好性

主要技术包括：
- **锁机制**: 互斥锁、读写锁、分段锁
- **原子操作**: CAS、原子指针、原子引用
- **无锁编程**: 无锁栈、无锁队列、ABA问题解决
- **通道机制**: Golang特有的并发原语

通过合理选择和优化，可以构建高性能、高可靠的并发数据结构，为各种高并发应用提供坚实的基础。 