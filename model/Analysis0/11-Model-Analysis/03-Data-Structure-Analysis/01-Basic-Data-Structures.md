# 基础数据结构分析

## 目录

- [基础数据结构分析](#基础数据结构分析)
  - [目录](#目录)
  - [概述](#概述)
    - [核心特性](#核心特性)
  - [数组 (Array)](#数组-array)
    - [形式化定义](#形式化定义)
    - [Golang实现](#golang实现)
    - [性能分析](#性能分析)
  - [链表 (Linked List)](#链表-linked-list)
    - [形式化定义1](#形式化定义1)
    - [Golang实现1](#golang实现1)
    - [性能分析1](#性能分析1)
  - [栈 (Stack)](#栈-stack)
    - [形式化定义2](#形式化定义2)
    - [Golang实现2](#golang实现2)
    - [性能分析2](#性能分析2)
  - [队列 (Queue)](#队列-queue)
    - [形式化定义3](#形式化定义3)
    - [Golang实现3](#golang实现3)
    - [性能分析3](#性能分析3)
  - [双端队列 (Deque)](#双端队列-deque)
    - [形式化定义4](#形式化定义4)
    - [Golang实现4](#golang实现4)
    - [性能分析4](#性能分析4)
  - [性能对比分析](#性能对比分析)
    - [时间复杂度对比](#时间复杂度对比)
    - [空间复杂度对比](#空间复杂度对比)
    - [应用场景分析](#应用场景分析)
      - [数组适用场景](#数组适用场景)
      - [链表适用场景](#链表适用场景)
      - [栈适用场景](#栈适用场景)
      - [队列适用场景](#队列适用场景)
      - [双端队列适用场景](#双端队列适用场景)
  - [最佳实践](#最佳实践)
    - [1. 选择原则](#1-选择原则)
      - [1.1 性能优先](#11-性能优先)
      - [1.2 内存优先](#12-内存优先)
      - [1.3 实现复杂度](#13-实现复杂度)
    - [2. 优化策略](#2-优化策略)
      - [2.1 内存优化](#21-内存优化)
      - [2.2 并发优化](#22-并发优化)
      - [2.3 缓存优化](#23-缓存优化)
    - [3. 错误处理](#3-错误处理)
      - [3.1 边界检查](#31-边界检查)
      - [3.2 空值检查](#32-空值检查)
    - [4. 测试策略](#4-测试策略)
      - [4.1 单元测试](#41-单元测试)
      - [4.2 性能测试](#42-性能测试)
      - [4.3 并发测试](#43-并发测试)
  - [总结](#总结)

## 概述

基础数据结构是计算机科学的核心基础，为高级数据结构和算法提供基础支撑。本文档对五种基础数据结构进行深入分析，包含形式化定义、Golang实现和性能分析。

### 核心特性

- **形式化定义**: 严格的数学描述和证明
- **Golang实现**: 完整的代码示例和测试
- **性能分析**: 时间复杂度和空间复杂度分析
- **最佳实践**: 实际应用中的最佳实践总结

## 数组 (Array)

### 形式化定义

**定义 2.1** (数组)
数组是一个四元组 $\mathcal{A} = (E, I, V, A_f)$，其中：

- $E$ 是元素类型集合
- $I = \{0, 1, \ldots, n-1\}$ 是索引集合
- $V: I \rightarrow E$ 是值函数
- $A_f: I \rightarrow E$ 是访问函数，$A_f(i) = V(i)$

**性质 2.1** (数组连续性)
对于数组 $\mathcal{A}$，任意两个相邻索引 $i, i+1 \in I$，其对应的内存地址满足：
$\text{addr}(i+1) = \text{addr}(i) + \text{sizeof}(E)$

**性质 2.2** (数组随机访问)
对于任意索引 $i \in I$，访问操作的时间复杂度为 $O(1)$：
$T_{\text{access}}(i) = O(1)$

### Golang实现

```go
// Array 通用数组实现
type Array[T any] struct {
    data []T
    size int
}

// NewArray 创建新数组
func NewArray[T any](capacity int) *Array[T] {
    return &Array[T]{
        data: make([]T, capacity),
        size: 0,
    }
}

// At 获取指定位置的元素
func (a *Array[T]) At(index int) (T, error) {
    if index < 0 || index >= a.size {
        var zero T
        return zero, fmt.Errorf("index %d out of bounds", index)
    }
    return a.data[index], nil
}

// Set 设置指定位置的元素
func (a *Array[T]) Set(index int, value T) error {
    if index < 0 || index >= a.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    a.data[index] = value
    return nil
}

// Append 追加元素
func (a *Array[T]) Append(value T) {
    if a.size >= len(a.data) {
        // 扩容
        newData := make([]T, len(a.data)*2)
        copy(newData, a.data)
        a.data = newData
    }
    a.data[a.size] = value
    a.size++
}

// Size 获取数组大小
func (a *Array[T]) Size() int {
    return a.size
}

// Capacity 获取数组容量
func (a *Array[T]) Capacity() int {
    return len(a.data)
}

// Clear 清空数组
func (a *Array[T]) Clear() {
    a.size = 0
    // 清空数据以避免内存泄漏
    for i := range a.data {
        var zero T
        a.data[i] = zero
    }
}
```

### 性能分析

| 操作 | 时间复杂度 | 空间复杂度 | 说明 |
|------|------------|------------|------|
| 访问 | $O(1)$ | $O(1)$ | 直接内存访问 |
| 搜索 | $O(n)$ | $O(1)$ | 线性搜索 |
| 插入 | $O(n)$ | $O(1)$ | 需要移动元素 |
| 删除 | $O(n)$ | $O(1)$ | 需要移动元素 |
| 追加 | 均摊 $O(1)$ | $O(1)$ | 动态扩容 |

## 链表 (Linked List)

### 形式化定义1

**定义 2.2** (链表节点)
链表节点是一个二元组 $\mathcal{N} = (data, next)$，其中：

- $data$ 是节点数据
- $next$ 是指向下一个节点的指针

**定义 2.3** (链表)
链表是一个三元组 $\mathcal{L} = (N, H, T)$，其中：

- $N$ 是节点集合
- $H \in N$ 是头节点
- $T \in N$ 是尾节点

**性质 2.3** (链表连接性)
对于链表 $\mathcal{L}$ 中的任意节点 $n_i$，存在路径：
$H \rightarrow n_1 \rightarrow n_2 \rightarrow \ldots \rightarrow n_i \rightarrow \ldots \rightarrow T$

### Golang实现1

```go
// Node 链表节点
type Node[T any] struct {
    Data T
    Next *Node[T]
}

// LinkedList 链表实现
type LinkedList[T any] struct {
    head *Node[T]
    tail *Node[T]
    size int
}

// NewLinkedList 创建新链表
func NewLinkedList[T any]() *LinkedList[T] {
    return &LinkedList[T]{
        head: nil,
        tail: nil,
        size: 0,
    }
}

// InsertAtHead 在头部插入
func (l *LinkedList[T]) InsertAtHead(data T) {
    newNode := &Node[T]{
        Data: data,
        Next: l.head,
    }
    
    if l.head == nil {
        l.tail = newNode
    }
    l.head = newNode
    l.size++
}

// InsertAtTail 在尾部插入
func (l *LinkedList[T]) InsertAtTail(data T) {
    newNode := &Node[T]{
        Data: data,
        Next: nil,
    }
    
    if l.tail == nil {
        l.head = newNode
    } else {
        l.tail.Next = newNode
    }
    l.tail = newNode
    l.size++
}

// InsertAt 在指定位置插入
func (l *LinkedList[T]) InsertAt(index int, data T) error {
    if index < 0 || index > l.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    if index == 0 {
        l.InsertAtHead(data)
        return nil
    }
    
    if index == l.size {
        l.InsertAtTail(data)
        return nil
    }
    
    // 找到插入位置的前一个节点
    current := l.head
    for i := 0; i < index-1; i++ {
        current = current.Next
    }
    
    newNode := &Node[T]{
        Data: data,
        Next: current.Next,
    }
    current.Next = newNode
    l.size++
    
    return nil
}

// DeleteAt 删除指定位置的元素
func (l *LinkedList[T]) DeleteAt(index int) error {
    if index < 0 || index >= l.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    if index == 0 {
        l.head = l.head.Next
        if l.head == nil {
            l.tail = nil
        }
        l.size--
        return nil
    }
    
    // 找到删除位置的前一个节点
    current := l.head
    for i := 0; i < index-1; i++ {
        current = current.Next
    }
    
    current.Next = current.Next.Next
    if current.Next == nil {
        l.tail = current
    }
    l.size--
    
    return nil
}

// Get 获取指定位置的元素
func (l *LinkedList[T]) Get(index int) (T, error) {
    if index < 0 || index >= l.size {
        var zero T
        return zero, fmt.Errorf("index %d out of bounds", index)
    }
    
    current := l.head
    for i := 0; i < index; i++ {
        current = current.Next
    }
    
    return current.Data, nil
}

// Size 获取链表大小
func (l *LinkedList[T]) Size() int {
    return l.size
}

// IsEmpty 判断是否为空
func (l *LinkedList[T]) IsEmpty() bool {
    return l.size == 0
}

// Clear 清空链表
func (l *LinkedList[T]) Clear() {
    l.head = nil
    l.tail = nil
    l.size = 0
}
```

### 性能分析1

| 操作 | 时间复杂度 | 空间复杂度 | 说明 |
|------|------------|------------|------|
| 访问 | $O(n)$ | $O(1)$ | 需要遍历 |
| 搜索 | $O(n)$ | $O(1)$ | 线性搜索 |
| 头部插入 | $O(1)$ | $O(1)$ | 直接插入 |
| 尾部插入 | $O(1)$ | $O(1)$ | 有尾指针 |
| 中间插入 | $O(n)$ | $O(1)$ | 需要遍历 |
| 头部删除 | $O(1)$ | $O(1)$ | 直接删除 |
| 尾部删除 | $O(n)$ | $O(1)$ | 需要遍历 |
| 中间删除 | $O(n)$ | $O(1)$ | 需要遍历 |

## 栈 (Stack)

### 形式化定义2

**定义 2.4** (栈)
栈是一个三元组 $\mathcal{S} = (E, O_s, T_s)$，其中：

- $E$ 是元素集合
- $O_s = \{\text{push}, \text{pop}, \text{peek}, \text{isEmpty}\}$ 是栈操作集合
- $T_s$ 是栈顶指针

**性质 2.4** (LIFO性质)
栈遵循后进先出 (LIFO) 原则：
$\forall e_1, e_2 \in E: \text{push}(e_1) \circ \text{push}(e_2) \circ \text{pop}() = e_2$

**性质 2.5** (栈操作复杂度)
对于栈 $\mathcal{S}$，基本操作的时间复杂度为：

- $\text{push}(e): O(1)$
- $\text{pop}(): O(1)$
- $\text{peek}(): O(1)$
- $\text{isEmpty}(): O(1)$

### Golang实现2

```go
// Stack 栈实现
type Stack[T any] struct {
    data []T
    top  int
}

// NewStack 创建新栈
func NewStack[T any](capacity int) *Stack[T] {
    return &Stack[T]{
        data: make([]T, capacity),
        top:  -1,
    }
}

// Push 入栈
func (s *Stack[T]) Push(element T) {
    s.top++
    if s.top >= len(s.data) {
        // 扩容
        newData := make([]T, len(s.data)*2)
        copy(newData, s.data)
        s.data = newData
    }
    s.data[s.top] = element
}

// Pop 出栈
func (s *Stack[T]) Pop() (T, error) {
    if s.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("stack is empty")
    }
    
    element := s.data[s.top]
    s.top--
    return element, nil
}

// Peek 查看栈顶元素
func (s *Stack[T]) Peek() (T, error) {
    if s.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("stack is empty")
    }
    
    return s.data[s.top], nil
}

// IsEmpty 判断栈是否为空
func (s *Stack[T]) IsEmpty() bool {
    return s.top == -1
}

// Size 获取栈大小
func (s *Stack[T]) Size() int {
    return s.top + 1
}

// Capacity 获取栈容量
func (s *Stack[T]) Capacity() int {
    return len(s.data)
}

// Clear 清空栈
func (s *Stack[T]) Clear() {
    s.top = -1
    // 清空数据以避免内存泄漏
    for i := range s.data {
        var zero T
        s.data[i] = zero
    }
}
```

### 性能分析2

| 操作 | 时间复杂度 | 空间复杂度 | 说明 |
|------|------------|------------|------|
| Push | 均摊 $O(1)$ | $O(1)$ | 动态扩容 |
| Pop | $O(1)$ | $O(1)$ | 直接操作 |
| Peek | $O(1)$ | $O(1)$ | 只读操作 |
| IsEmpty | $O(1)$ | $O(1)$ | 状态检查 |

## 队列 (Queue)

### 形式化定义3

**定义 2.5** (队列)
队列是一个四元组 $\mathcal{Q} = (E, O_q, F, R)$，其中：

- $E$ 是元素集合
- $O_q = \{\text{enqueue}, \text{dequeue}, \text{front}, \text{isEmpty}\}$ 是队列操作集合
- $F$ 是队首指针
- $R$ 是队尾指针

**性质 2.6** (FIFO性质)
队列遵循先进先出 (FIFO) 原则：
$\forall e_1, e_2 \in E: \text{enqueue}(e_1) \circ \text{enqueue}(e_2) \circ \text{dequeue}() = e_1$

**性质 2.7** (队列操作复杂度)
对于队列 $\mathcal{Q}$，基本操作的时间复杂度为：

- $\text{enqueue}(e): O(1)$
- $\text{dequeue}(): O(1)$
- $\text{front}(): O(1)$
- $\text{isEmpty}(): O(1)$

### Golang实现3

```go
// Queue 队列实现
type Queue[T any] struct {
    data []T
    head int
    tail int
    size int
}

// NewQueue 创建新队列
func NewQueue[T any](capacity int) *Queue[T] {
    return &Queue[T]{
        data: make([]T, capacity),
        head: 0,
        tail: 0,
        size: 0,
    }
}

// Enqueue 入队
func (q *Queue[T]) Enqueue(element T) {
    if q.size >= len(q.data) {
        // 扩容
        newData := make([]T, len(q.data)*2)
        // 重新排列元素
        for i := 0; i < q.size; i++ {
            newData[i] = q.data[(q.head+i)%len(q.data)]
        }
        q.data = newData
        q.head = 0
        q.tail = q.size
    }
    
    q.data[q.tail] = element
    q.tail = (q.tail + 1) % len(q.data)
    q.size++
}

// Dequeue 出队
func (q *Queue[T]) Dequeue() (T, error) {
    if q.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("queue is empty")
    }
    
    element := q.data[q.head]
    q.head = (q.head + 1) % len(q.data)
    q.size--
    return element, nil
}

// Front 查看队首元素
func (q *Queue[T]) Front() (T, error) {
    if q.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("queue is empty")
    }
    
    return q.data[q.head], nil
}

// IsEmpty 判断队列是否为空
func (q *Queue[T]) IsEmpty() bool {
    return q.size == 0
}

// Size 获取队列大小
func (q *Queue[T]) Size() int {
    return q.size
}

// Capacity 获取队列容量
func (q *Queue[T]) Capacity() int {
    return len(q.data)
}

// Clear 清空队列
func (q *Queue[T]) Clear() {
    q.head = 0
    q.tail = 0
    q.size = 0
    // 清空数据以避免内存泄漏
    for i := range q.data {
        var zero T
        q.data[i] = zero
    }
}
```

### 性能分析3

| 操作 | 时间复杂度 | 空间复杂度 | 说明 |
|------|------------|------------|------|
| Enqueue | 均摊 $O(1)$ | $O(1)$ | 动态扩容 |
| Dequeue | $O(1)$ | $O(1)$ | 直接操作 |
| Front | $O(1)$ | $O(1)$ | 只读操作 |
| IsEmpty | $O(1)$ | $O(1)$ | 状态检查 |

## 双端队列 (Deque)

### 形式化定义4

**定义 2.6** (双端队列)
双端队列是一个五元组 $\mathcal{D} = (E, O_d, F, R, C)$，其中：

- $E$ 是元素集合
- $O_d = \{\text{pushFront}, \text{pushBack}, \text{popFront}, \text{popBack}, \text{front}, \text{back}, \text{isEmpty}\}$ 是双端队列操作集合
- $F$ 是前端指针
- $R$ 是后端指针
- $C$ 是容量

**性质 2.8** (双端操作性质)
双端队列支持两端操作：
$\forall e \in E: \text{pushFront}(e) \circ \text{popBack}() = e$
$\forall e \in E: \text{pushBack}(e) \circ \text{popFront}() = e$

### Golang实现4

```go
// Deque 双端队列实现
type Deque[T any] struct {
    data []T
    head int
    tail int
    size int
}

// NewDeque 创建新双端队列
func NewDeque[T any](capacity int) *Deque[T] {
    return &Deque[T]{
        data: make([]T, capacity),
        head: 0,
        tail: 0,
        size: 0,
    }
}

// PushFront 前端入队
func (d *Deque[T]) PushFront(element T) {
    if d.size >= len(d.data) {
        // 扩容
        newData := make([]T, len(d.data)*2)
        // 重新排列元素
        for i := 0; i < d.size; i++ {
            newData[i+1] = d.data[(d.head+i)%len(d.data)]
        }
        d.data = newData
        d.head = 1
        d.tail = d.size + 1
    } else {
        d.head = (d.head - 1 + len(d.data)) % len(d.data)
    }
    
    d.data[d.head] = element
    d.size++
}

// PushBack 后端入队
func (d *Deque[T]) PushBack(element T) {
    if d.size >= len(d.data) {
        // 扩容
        newData := make([]T, len(d.data)*2)
        // 重新排列元素
        for i := 0; i < d.size; i++ {
            newData[i] = d.data[(d.head+i)%len(d.data)]
        }
        d.data = newData
        d.head = 0
        d.tail = d.size
    }
    
    d.data[d.tail] = element
    d.tail = (d.tail + 1) % len(d.data)
    d.size++
}

// PopFront 前端出队
func (d *Deque[T]) PopFront() (T, error) {
    if d.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("deque is empty")
    }
    
    element := d.data[d.head]
    d.head = (d.head + 1) % len(d.data)
    d.size--
    return element, nil
}

// PopBack 后端出队
func (d *Deque[T]) PopBack() (T, error) {
    if d.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("deque is empty")
    }
    
    d.tail = (d.tail - 1 + len(d.data)) % len(d.data)
    element := d.data[d.tail]
    d.size--
    return element, nil
}

// Front 查看前端元素
func (d *Deque[T]) Front() (T, error) {
    if d.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("deque is empty")
    }
    
    return d.data[d.head], nil
}

// Back 查看后端元素
func (d *Deque[T]) Back() (T, error) {
    if d.IsEmpty() {
        var zero T
        return zero, fmt.Errorf("deque is empty")
    }
    
    index := (d.tail - 1 + len(d.data)) % len(d.data)
    return d.data[index], nil
}

// IsEmpty 判断是否为空
func (d *Deque[T]) IsEmpty() bool {
    return d.size == 0
}

// Size 获取大小
func (d *Deque[T]) Size() int {
    return d.size
}

// Capacity 获取容量
func (d *Deque[T]) Capacity() int {
    return len(d.data)
}

// Clear 清空
func (d *Deque[T]) Clear() {
    d.head = 0
    d.tail = 0
    d.size = 0
    // 清空数据以避免内存泄漏
    for i := range d.data {
        var zero T
        d.data[i] = zero
    }
}
```

### 性能分析4

| 操作 | 时间复杂度 | 空间复杂度 | 说明 |
|------|------------|------------|------|
| PushFront | 均摊 $O(1)$ | $O(1)$ | 动态扩容 |
| PushBack | 均摊 $O(1)$ | $O(1)$ | 动态扩容 |
| PopFront | $O(1)$ | $O(1)$ | 直接操作 |
| PopBack | $O(1)$ | $O(1)$ | 直接操作 |
| Front | $O(1)$ | $O(1)$ | 只读操作 |
| Back | $O(1)$ | $O(1)$ | 只读操作 |

## 性能对比分析

### 时间复杂度对比

| 数据结构 | 访问 | 搜索 | 插入 | 删除 | 特殊操作 |
|----------|------|------|------|------|----------|
| 数组 | $O(1)$ | $O(n)$ | $O(n)$ | $O(n)$ | 随机访问 |
| 链表 | $O(n)$ | $O(n)$ | $O(1)$ | $O(1)$ | 动态大小 |
| 栈 | $O(1)$ | $O(n)$ | $O(1)$ | $O(1)$ | LIFO |
| 队列 | $O(1)$ | $O(n)$ | $O(1)$ | $O(1)$ | FIFO |
| 双端队列 | $O(1)$ | $O(n)$ | $O(1)$ | $O(1)$ | 双端操作 |

### 空间复杂度对比

| 数据结构 | 基础空间 | 额外空间 | 内存局部性 |
|----------|----------|----------|------------|
| 数组 | $O(n)$ | $O(1)$ | 优秀 |
| 链表 | $O(n)$ | $O(n)$ | 较差 |
| 栈 | $O(n)$ | $O(1)$ | 优秀 |
| 队列 | $O(n)$ | $O(1)$ | 优秀 |
| 双端队列 | $O(n)$ | $O(1)$ | 优秀 |

### 应用场景分析

#### 数组适用场景

- **随机访问**: 需要频繁随机访问元素
- **内存效率**: 对内存使用要求严格
- **缓存友好**: 需要良好的缓存性能
- **固定大小**: 数据大小相对固定

#### 链表适用场景

- **动态大小**: 数据大小变化频繁
- **频繁插入删除**: 在中间位置频繁操作
- **内存分散**: 可以接受内存碎片化
- **实现简单**: 需要快速实现原型

#### 栈适用场景

- **函数调用**: 函数调用栈管理
- **表达式求值**: 后缀表达式计算
- **括号匹配**: 语法分析中的括号检查
- **深度优先搜索**: 图的DFS算法

#### 队列适用场景

- **任务调度**: 操作系统任务队列
- **广度优先搜索**: 图的BFS算法
- **消息队列**: 异步消息处理
- **缓冲区管理**: 生产者-消费者模式

#### 双端队列适用场景

- **滑动窗口**: 算法中的滑动窗口
- **单调队列**: 单调递增/递减队列
- **双向遍历**: 需要双向遍历的场景
- **缓存实现**: LRU缓存的双端操作

## 最佳实践

### 1. 选择原则

#### 1.1 性能优先

- **随机访问**: 选择数组
- **频繁插入删除**: 选择链表
- **LIFO操作**: 选择栈
- **FIFO操作**: 选择队列
- **双端操作**: 选择双端队列

#### 1.2 内存优先

- **内存紧张**: 选择数组
- **内存充足**: 可以选择链表
- **缓存敏感**: 选择连续存储结构

#### 1.3 实现复杂度

- **快速原型**: 选择简单实现
- **生产环境**: 选择成熟实现
- **性能要求**: 选择优化实现

### 2. 优化策略

#### 2.1 内存优化

```go
// 使用对象池减少GC压力
var nodePool = sync.Pool{
    New: func() interface{} {
        return &Node{}
    },
}

func getNode() *Node {
    return nodePool.Get().(*Node)
}

func putNode(node *Node) {
    nodePool.Put(node)
}
```

#### 2.2 并发优化

```go
// 使用原子操作避免锁
type AtomicStack[T any] struct {
    head unsafe.Pointer
}

func (s *AtomicStack[T]) Push(element T) {
    newNode := &Node[T]{Data: element}
    for {
        oldHead := s.head
        newNode.Next = (*Node[T])(oldHead)
        if atomic.CompareAndSwapPointer(&s.head, oldHead, unsafe.Pointer(newNode)) {
            break
        }
    }
}
```

#### 2.3 缓存优化

```go
// 使用内存对齐提高缓存性能
type CacheFriendlyNode[T any] struct {
    Data T
    Next *CacheFriendlyNode[T]
    _    [64 - unsafe.Sizeof(T) - unsafe.Sizeof(unsafe.Pointer(nil))]byte // 填充到64字节
}
```

### 3. 错误处理

#### 3.1 边界检查

```go
func (a *Array[T]) boundsCheck(index int) error {
    if index < 0 || index >= a.size {
        return fmt.Errorf("index %d out of bounds [0, %d)", index, a.size)
    }
    return nil
}
```

#### 3.2 空值检查

```go
func (l *LinkedList[T]) nullCheck() error {
    if l.head == nil {
        return fmt.Errorf("linked list is empty")
    }
    return nil
}
```

### 4. 测试策略

#### 4.1 单元测试

```go
func TestArrayOperations(t *testing.T) {
    arr := NewArray[int](10)
    
    // 测试插入
    arr.Append(1)
    arr.Append(2)
    arr.Append(3)
    
    if arr.Size() != 3 {
        t.Errorf("Expected size 3, got %d", arr.Size())
    }
    
    // 测试访问
    if val, err := arr.At(1); err != nil || val != 2 {
        t.Errorf("Expected value 2 at index 1, got %d", val)
    }
}
```

#### 4.2 性能测试

```go
func BenchmarkArrayAppend(b *testing.B) {
    arr := NewArray[int](1000)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        arr.Append(i)
    }
}

func BenchmarkLinkedListAppend(b *testing.B) {
    list := NewLinkedList[int]()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        list.InsertAtTail(i)
    }
}
```

#### 4.3 并发测试

```go
func TestConcurrentStack(t *testing.T) {
    stack := NewStack[int](1000)
    var wg sync.WaitGroup
    
    // 并发推入
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            stack.Push(id)
        }(i)
    }
    
    wg.Wait()
    
    // 验证结果
    if stack.Size() != 1000 {
        t.Errorf("Expected size 1000, got %d", stack.Size())
    }
}
```

---

## 总结

本文档对五种基础数据结构进行了深入分析，包括：

1. **形式化定义**: 严格的数学描述和性质证明
2. **Golang实现**: 完整的代码示例和错误处理
3. **性能分析**: 详细的时间复杂度和空间复杂度分析
4. **应用场景**: 各种数据结构的适用场景分析
5. **最佳实践**: 实际应用中的优化策略和测试方法

这些基础数据结构为高级数据结构和算法提供了重要支撑，在实际应用中需要根据具体需求选择合适的数据结构。

---

**最后更新**: 2024-12-19  
**当前状态**: ✅ 基础数据结构分析完成  
**下一步**: 并发数据结构分析
