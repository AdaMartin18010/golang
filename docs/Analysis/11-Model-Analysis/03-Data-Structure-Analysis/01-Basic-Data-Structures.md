# 基础数据结构分析

## 目录

1. [概述](#概述)
2. [数组 (Array)](#数组-array)
3. [链表 (Linked List)](#链表-linked-list)
4. [栈 (Stack)](#栈-stack)
5. [队列 (Queue)](#队列-queue)
6. [性能对比](#性能对比)
7. [应用场景](#应用场景)
8. [最佳实践](#最佳实践)

## 概述

基础数据结构是构建复杂系统的基础，它们提供了最基本的数据组织和操作方式。
在Golang中，这些数据结构通常作为更复杂数据结构的基础组件。

### 核心特征

- **简单性**: 概念清晰，实现简单
- **高效性**: 特定操作具有最优性能
- **通用性**: 广泛应用于各种场景
- **可组合性**: 可以组合构建复杂结构

## 数组 (Array)

### 1.1 形式化定义

#### 1.1.1 数学定义

数组可以形式化定义为四元组：

$$\mathcal{A} = \langle \mathcal{E}, \mathcal{I}, \mathcal{V}, \mathcal{A}_f \rangle$$

其中：
- $\mathcal{E}$：元素类型集合
- $\mathcal{I}$：索引集合 $\{0, 1, \ldots, n-1\}$
- $\mathcal{V}$：值函数 $\mathcal{V}: \mathcal{I} \rightarrow \mathcal{E} \cup \{\text{nil}\}$
- $\mathcal{A}_f$：访问函数 $\mathcal{A}_f(i) = \mathcal{V}(i)$

#### 1.1.2 操作语义

**访问操作**：
$$\text{Access}(\mathcal{A}, i) = \mathcal{A}_f(i)$$

**更新操作**：
$$\text{Update}(\mathcal{A}, i, e) = \mathcal{A}' \text{ where } \mathcal{V}'(j) = \begin{cases} 
e & \text{if } j = i \\
\mathcal{V}(j) & \text{otherwise}
\end{cases}$$

### 1.2 Golang实现

#### 1.2.1 基础数组

```go
// Array 基础数组实现
type Array[T any] struct {
    elements []T
    size     int
}

// NewArray 创建新数组
func NewArray[T any](capacity int) *Array[T] {
    return &Array[T]{
        elements: make([]T, capacity),
        size:     0,
    }
}

// Get 获取元素
func (a *Array[T]) Get(index int) (T, error) {
    if index < 0 || index >= a.size {
        var zero T
        return zero, fmt.Errorf("index %d out of bounds", index)
    }
    return a.elements[index], nil
}

// Set 设置元素
func (a *Array[T]) Set(index int, element T) error {
    if index < 0 || index >= a.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    a.elements[index] = element
    return nil
}

// Append 追加元素
func (a *Array[T]) Append(element T) {
    if a.size >= len(a.elements) {
        // 扩容
        newElements := make([]T, len(a.elements)*2)
        copy(newElements, a.elements)
        a.elements = newElements
    }
    a.elements[a.size] = element
    a.size++
}

// Size 获取大小
func (a *Array[T]) Size() int {
    return a.size
}
```

#### 1.2.2 动态数组

```go
// DynamicArray 动态数组实现
type DynamicArray[T any] struct {
    elements []T
    size     int
    capacity int
}

// NewDynamicArray 创建动态数组
func NewDynamicArray[T any](initialCapacity int) *DynamicArray[T] {
    if initialCapacity <= 0 {
        initialCapacity = 10
    }
    return &DynamicArray[T]{
        elements: make([]T, initialCapacity),
        size:     0,
        capacity: initialCapacity,
    }
}

// grow 扩容
func (da *DynamicArray[T]) grow() {
    newCapacity := da.capacity * 2
    newElements := make([]T, newCapacity)
    copy(newElements, da.elements)
    da.elements = newElements
    da.capacity = newCapacity
}

// Add 添加元素
func (da *DynamicArray[T]) Add(element T) {
    if da.size >= da.capacity {
        da.grow()
    }
    da.elements[da.size] = element
    da.size++
}

// Remove 删除元素
func (da *DynamicArray[T]) Remove(index int) error {
    if index < 0 || index >= da.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    // 移动元素
    for i := index; i < da.size-1; i++ {
        da.elements[i] = da.elements[i+1]
    }
    da.size--
    
    // 缩容检查
    if da.size < da.capacity/4 && da.capacity > 10 {
        da.shrink()
    }
    
    return nil
}

// shrink 缩容
func (da *DynamicArray[T]) shrink() {
    newCapacity := da.capacity / 2
    newElements := make([]T, newCapacity)
    copy(newElements, da.elements[:da.size])
    da.elements = newElements
    da.capacity = newCapacity
}
```

### 1.3 性能分析

#### 1.3.1 时间复杂度

| 操作 | 平均情况 | 最坏情况 | 最好情况 |
|------|----------|----------|----------|
| 访问 | $O(1)$ | $O(1)$ | $O(1)$ |
| 搜索 | $O(n)$ | $O(n)$ | $O(1)$ |
| 插入 | $O(1)$ | $O(n)$ | $O(1)$ |
| 删除 | $O(n)$ | $O(n)$ | $O(1)$ |

#### 1.3.2 空间复杂度

- **静态数组**：$O(n)$
- **动态数组**：$O(n)$ (摊销)

## 链表 (Linked List)

### 2.1 形式化定义

#### 2.1.1 数学定义

链表可以形式化定义为：

$$\mathcal{L} = \langle \mathcal{N}, \mathcal{E}, \mathcal{P}, \mathcal{H}, \mathcal{T} \rangle$$

其中：
- $\mathcal{N}$：节点集合
- $\mathcal{E}$：元素集合
- $\mathcal{P}$：指针函数 $\mathcal{P}: \mathcal{N} \rightarrow \mathcal{N} \cup \{\text{nil}\}$
- $\mathcal{H}$：头节点 $\mathcal{H} \in \mathcal{N}$
- $\mathcal{T}$：尾节点 $\mathcal{T} \in \mathcal{N}$

#### 2.1.2 节点结构

节点可以定义为：

$$\mathcal{N}_i = \langle e_i, p_i \rangle$$

其中 $e_i \in \mathcal{E}$ 是元素，$p_i \in \mathcal{N} \cup \{\text{nil}\}$ 是指针。

### 2.2 Golang实现

#### 2.2.1 单链表

```go
// Node 链表节点
type Node[T any] struct {
    Data T
    Next *Node[T]
}

// LinkedList 单链表
type LinkedList[T any] struct {
    head *Node[T]
    size int
}

// NewLinkedList 创建新链表
func NewLinkedList[T any]() *LinkedList[T] {
    return &LinkedList[T]{
        head: nil,
        size: 0,
    }
}

// InsertAtHead 在头部插入
func (ll *LinkedList[T]) InsertAtHead(data T) {
    newNode := &Node[T]{
        Data: data,
        Next: ll.head,
    }
    ll.head = newNode
    ll.size++
}

// InsertAtTail 在尾部插入
func (ll *LinkedList[T]) InsertAtTail(data T) {
    newNode := &Node[T]{
        Data: data,
        Next: nil,
    }
    
    if ll.head == nil {
        ll.head = newNode
    } else {
        current := ll.head
        for current.Next != nil {
            current = current.Next
        }
        current.Next = newNode
    }
    ll.size++
}

// Delete 删除指定元素
func (ll *LinkedList[T]) Delete(data T) bool {
    if ll.head == nil {
        return false
    }
    
    if ll.head.Data == data {
        ll.head = ll.head.Next
        ll.size--
        return true
    }
    
    current := ll.head
    for current.Next != nil {
        if current.Next.Data == data {
            current.Next = current.Next.Next
            ll.size--
            return true
        }
        current = current.Next
    }
    
    return false
}

// Search 搜索元素
func (ll *LinkedList[T]) Search(data T) bool {
    current := ll.head
    for current != nil {
        if current.Data == data {
            return true
        }
        current = current.Next
    }
    return false
}
```

#### 2.2.2 双向链表

```go
// DoublyNode 双向链表节点
type DoublyNode[T any] struct {
    Data T
    Prev *DoublyNode[T]
    Next *DoublyNode[T]
}

// DoublyLinkedList 双向链表
type DoublyLinkedList[T any] struct {
    head *DoublyNode[T]
    tail *DoublyNode[T]
    size int
}

// NewDoublyLinkedList 创建双向链表
func NewDoublyLinkedList[T any]() *DoublyLinkedList[T] {
    return &DoublyLinkedList[T]{
        head: nil,
        tail: nil,
        size: 0,
    }
}

// InsertAtHead 在头部插入
func (dll *DoublyLinkedList[T]) InsertAtHead(data T) {
    newNode := &DoublyNode[T]{
        Data: data,
        Prev: nil,
        Next: dll.head,
    }
    
    if dll.head != nil {
        dll.head.Prev = newNode
    } else {
        dll.tail = newNode
    }
    
    dll.head = newNode
    dll.size++
}

// InsertAtTail 在尾部插入
func (dll *DoublyLinkedList[T]) InsertAtTail(data T) {
    newNode := &DoublyNode[T]{
        Data: data,
        Prev: dll.tail,
        Next: nil,
    }
    
    if dll.tail != nil {
        dll.tail.Next = newNode
    } else {
        dll.head = newNode
    }
    
    dll.tail = newNode
    dll.size++
}

// Delete 删除指定元素
func (dll *DoublyLinkedList[T]) Delete(data T) bool {
    current := dll.head
    
    for current != nil {
        if current.Data == data {
            if current.Prev != nil {
                current.Prev.Next = current.Next
            } else {
                dll.head = current.Next
            }
            
            if current.Next != nil {
                current.Next.Prev = current.Prev
            } else {
                dll.tail = current.Prev
            }
            
            dll.size--
            return true
        }
        current = current.Next
    }
    
    return false
}
```

### 2.3 性能分析

#### 2.3.1 时间复杂度

| 操作 | 单链表 | 双向链表 |
|------|--------|----------|
| 头部插入 | $O(1)$ | $O(1)$ |
| 尾部插入 | $O(n)$ | $O(1)$ |
| 头部删除 | $O(1)$ | $O(1)$ |
| 尾部删除 | $O(n)$ | $O(1)$ |
| 搜索 | $O(n)$ | $O(n)$ |

#### 2.3.2 空间复杂度

- **单链表**：$O(n)$
- **双向链表**：$O(n)$ (每个节点多一个指针)

## 栈 (Stack)

### 3.1 形式化定义

#### 3.1.1 数学定义

栈可以形式化定义为：

$$\mathcal{S} = \langle \mathcal{E}, \mathcal{O}_s, \mathcal{T}_s, \mathcal{I}_s \rangle$$

其中：
- $\mathcal{E}$：元素集合
- $\mathcal{O}_s$：栈操作集合 $\{\text{push}, \text{pop}, \text{peek}, \text{isEmpty}\}$
- $\mathcal{T}_s$：栈顶指针
- $\mathcal{I}_s$：栈内容 $\mathcal{I}_s: \mathbb{N} \rightarrow \mathcal{E}$

#### 3.1.2 LIFO性质

栈遵循后进先出 (LIFO) 原则：

$$\forall e_1, e_2 \in \mathcal{E}: \text{push}(e_1) \circ \text{push}(e_2) \circ \text{pop}() = e_2$$

### 3.2 Golang实现

#### 3.2.1 基于数组的栈

```go
// Stack 栈实现
type Stack[T any] struct {
    elements []T
    top      int
    capacity int
}

// NewStack 创建新栈
func NewStack[T any](capacity int) *Stack[T] {
    return &Stack[T]{
        elements: make([]T, capacity),
        top:      -1,
        capacity: capacity,
    }
}

// Push 入栈
func (s *Stack[T]) Push(element T) error {
    if s.top >= s.capacity-1 {
        return fmt.Errorf("stack overflow")
    }
    s.top++
    s.elements[s.top] = element
    return nil
}

// Pop 出栈
func (s *Stack[T]) Pop() (T, error) {
    var zero T
    if s.IsEmpty() {
        return zero, fmt.Errorf("stack underflow")
    }
    element := s.elements[s.top]
    s.top--
    return element, nil
}

// Peek 查看栈顶元素
func (s *Stack[T]) Peek() (T, error) {
    var zero T
    if s.IsEmpty() {
        return zero, fmt.Errorf("stack is empty")
    }
    return s.elements[s.top], nil
}

// IsEmpty 检查栈是否为空
func (s *Stack[T]) IsEmpty() bool {
    return s.top == -1
}

// Size 获取栈大小
func (s *Stack[T]) Size() int {
    return s.top + 1
}
```

#### 3.2.2 基于链表的栈

```go
// LinkedStack 基于链表的栈
type LinkedStack[T any] struct {
    top  *Node[T]
    size int
}

// NewLinkedStack 创建链表栈
func NewLinkedStack[T any]() *LinkedStack[T] {
    return &LinkedStack[T]{
        top:  nil,
        size: 0,
    }
}

// Push 入栈
func (ls *LinkedStack[T]) Push(element T) {
    newNode := &Node[T]{
        Data: element,
        Next: ls.top,
    }
    ls.top = newNode
    ls.size++
}

// Pop 出栈
func (ls *LinkedStack[T]) Pop() (T, error) {
    var zero T
    if ls.IsEmpty() {
        return zero, fmt.Errorf("stack underflow")
    }
    
    element := ls.top.Data
    ls.top = ls.top.Next
    ls.size--
    return element, nil
}

// Peek 查看栈顶元素
func (ls *LinkedStack[T]) Peek() (T, error) {
    var zero T
    if ls.IsEmpty() {
        return zero, fmt.Errorf("stack is empty")
    }
    return ls.top.Data, nil
}

// IsEmpty 检查栈是否为空
func (ls *LinkedStack[T]) IsEmpty() bool {
    return ls.top == nil
}
```

### 3.3 应用场景

#### 3.3.1 函数调用栈

```go
// CallStack 函数调用栈模拟
type CallFrame struct {
    FunctionName string
    Parameters   []interface{}
    ReturnValue  interface{}
}

type CallStack struct {
    frames *Stack[*CallFrame]
}

func NewCallStack() *CallStack {
    return &CallStack{
        frames: NewStack[*CallFrame](1000),
    }
}

func (cs *CallStack) PushFrame(frame *CallFrame) {
    cs.frames.Push(frame)
}

func (cs *CallStack) PopFrame() (*CallFrame, error) {
    return cs.frames.Pop()
}
```

#### 3.3.2 表达式求值

```go
// ExpressionEvaluator 表达式求值器
type ExpressionEvaluator struct {
    operandStack  *Stack[float64]
    operatorStack *Stack[string]
}

func NewExpressionEvaluator() *ExpressionEvaluator {
    return &ExpressionEvaluator{
        operandStack:  NewStack[float64](100),
        operatorStack: NewStack[string](100),
    }
}

func (ee *ExpressionEvaluator) Evaluate(expression string) (float64, error) {
    // 实现中缀表达式求值
    // 使用两个栈：操作数栈和运算符栈
    return 0, nil
}
```

## 队列 (Queue)

### 4.1 形式化定义

#### 4.1.1 数学定义

队列可以形式化定义为：

$$\mathcal{Q} = \langle \mathcal{E}, \mathcal{O}_q, \mathcal{F}, \mathcal{R}, \mathcal{I}_q \rangle$$

其中：
- $\mathcal{E}$：元素集合
- $\mathcal{O}_q$：队列操作集合 $\{\text{enqueue}, \text{dequeue}, \text{front}, \text{isEmpty}\}$
- $\mathcal{F}$：队首指针
- $\mathcal{R}$：队尾指针
- $\mathcal{I}_q$：队列内容 $\mathcal{I}_q: \mathbb{N} \rightarrow \mathcal{E}$

#### 4.1.2 FIFO性质

队列遵循先进先出 (FIFO) 原则：

$$\forall e_1, e_2 \in \mathcal{E}: \text{enqueue}(e_1) \circ \text{enqueue}(e_2) \circ \text{dequeue}() = e_1$$

### 4.2 Golang实现

#### 4.2.1 基于数组的队列

```go
// Queue 队列实现
type Queue[T any] struct {
    elements []T
    front    int
    rear     int
    size     int
    capacity int
}

// NewQueue 创建新队列
func NewQueue[T any](capacity int) *Queue[T] {
    return &Queue[T]{
        elements: make([]T, capacity),
        front:    0,
        rear:     -1,
        size:     0,
        capacity: capacity,
    }
}

// Enqueue 入队
func (q *Queue[T]) Enqueue(element T) error {
    if q.IsFull() {
        return fmt.Errorf("queue is full")
    }
    
    q.rear = (q.rear + 1) % q.capacity
    q.elements[q.rear] = element
    q.size++
    return nil
}

// Dequeue 出队
func (q *Queue[T]) Dequeue() (T, error) {
    var zero T
    if q.IsEmpty() {
        return zero, fmt.Errorf("queue is empty")
    }
    
    element := q.elements[q.front]
    q.front = (q.front + 1) % q.capacity
    q.size--
    return element, nil
}

// Front 查看队首元素
func (q *Queue[T]) Front() (T, error) {
    var zero T
    if q.IsEmpty() {
        return zero, fmt.Errorf("queue is empty")
    }
    return q.elements[q.front], nil
}

// IsEmpty 检查队列是否为空
func (q *Queue[T]) IsEmpty() bool {
    return q.size == 0
}

// IsFull 检查队列是否已满
func (q *Queue[T]) IsFull() bool {
    return q.size == q.capacity
}

// Size 获取队列大小
func (q *Queue[T]) Size() int {
    return q.size
}
```

#### 4.2.2 基于链表的队列

```go
// LinkedQueue 基于链表的队列
type LinkedQueue[T any] struct {
    head *Node[T]
    tail *Node[T]
    size int
}

// NewLinkedQueue 创建链表队列
func NewLinkedQueue[T any]() *LinkedQueue[T] {
    return &LinkedQueue[T]{
        head: nil,
        tail: nil,
        size: 0,
    }
}

// Enqueue 入队
func (lq *LinkedQueue[T]) Enqueue(element T) {
    newNode := &Node[T]{
        Data: element,
        Next: nil,
    }
    
    if lq.IsEmpty() {
        lq.head = newNode
        lq.tail = newNode
    } else {
        lq.tail.Next = newNode
        lq.tail = newNode
    }
    lq.size++
}

// Dequeue 出队
func (lq *LinkedQueue[T]) Dequeue() (T, error) {
    var zero T
    if lq.IsEmpty() {
        return zero, fmt.Errorf("queue is empty")
    }
    
    element := lq.head.Data
    lq.head = lq.head.Next
    
    if lq.head == nil {
        lq.tail = nil
    }
    
    lq.size--
    return element, nil
}

// Front 查看队首元素
func (lq *LinkedQueue[T]) Front() (T, error) {
    var zero T
    if lq.IsEmpty() {
        return zero, fmt.Errorf("queue is empty")
    }
    return lq.head.Data, nil
}

// IsEmpty 检查队列是否为空
func (lq *LinkedQueue[T]) IsEmpty() bool {
    return lq.head == nil
}
```

### 4.3 应用场景

#### 4.3.1 任务队列

```go
// Task 任务定义
type Task struct {
    ID       string
    Priority int
    Handler  func() error
}

// TaskQueue 任务队列
type TaskQueue struct {
    queue *Queue[*Task]
}

func NewTaskQueue() *TaskQueue {
    return &TaskQueue{
        queue: NewQueue[*Task](1000),
    }
}

func (tq *TaskQueue) AddTask(task *Task) error {
    return tq.queue.Enqueue(task)
}

func (tq *TaskQueue) ProcessNextTask() error {
    task, err := tq.queue.Dequeue()
    if err != nil {
        return err
    }
    return task.Handler()
}
```

#### 4.3.2 消息队列

```go
// Message 消息定义
type Message struct {
    ID      string
    Content interface{}
    Time    time.Time
}

// MessageQueue 消息队列
type MessageQueue struct {
    queue *Queue[*Message]
}

func NewMessageQueue() *MessageQueue {
    return &MessageQueue{
        queue: NewQueue[*Message](10000),
    }
}

func (mq *MessageQueue) SendMessage(msg *Message) error {
    return mq.queue.Enqueue(msg)
}

func (mq *MessageQueue) ReceiveMessage() (*Message, error) {
    return mq.queue.Dequeue()
}
```

## 性能对比

### 1. 时间复杂度对比

| 操作 | 数组 | 链表 | 栈 | 队列 |
|------|------|------|----|----|
| 访问 | $O(1)$ | $O(n)$ | $O(1)$ | $O(1)$ |
| 插入 | $O(n)$ | $O(1)$ | $O(1)$ | $O(1)$ |
| 删除 | $O(n)$ | $O(1)$ | $O(1)$ | $O(1)$ |
| 搜索 | $O(n)$ | $O(n)$ | $O(n)$ | $O(n)$ |

### 2. 空间复杂度对比

| 数据结构 | 空间复杂度 | 额外开销 |
|----------|------------|----------|
| 数组 | $O(n)$ | 无 |
| 链表 | $O(n)$ | 指针开销 |
| 栈 | $O(n)$ | 无 |
| 队列 | $O(n)$ | 无 |

### 3. 缓存性能对比

| 数据结构 | 缓存友好性 | 局部性 |
|----------|------------|--------|
| 数组 | 高 | 好 |
| 链表 | 低 | 差 |
| 栈 | 高 | 好 |
| 队列 | 中等 | 中等 |

## 应用场景

### 1. 数组应用场景

- **数值计算**: 矩阵运算、向量计算
- **图像处理**: 像素数据存储
- **音频处理**: 音频采样数据
- **缓存实现**: 固定大小的缓存

### 2. 链表应用场景

- **内存管理**: 空闲内存块链表
- **文件系统**: 文件分配表
- **哈希表**: 冲突解决
- **LRU缓存**: 最近最少使用缓存

### 3. 栈应用场景

- **函数调用**: 程序执行栈
- **表达式求值**: 中缀转后缀
- **括号匹配**: 语法检查
- **深度优先搜索**: 图遍历

### 4. 队列应用场景

- **任务调度**: 进程调度
- **消息队列**: 异步通信
- **广度优先搜索**: 图遍历
- **缓冲区**: 数据流处理

## 最佳实践

### 1. 选择指南

1. **数组**：适用于随机访问频繁的场景
2. **链表**：适用于频繁插入删除的场景
3. **栈**：适用于后进先出的场景
4. **队列**：适用于先进先出的场景

### 2. 性能优化

1. **预分配容量**：减少动态扩容开销
2. **批量操作**：减少函数调用开销
3. **内存对齐**：提高缓存性能
4. **对象池**：减少GC压力

### 3. 错误处理

1. **边界检查**：防止越界访问
2. **空值检查**：防止空指针异常
3. **容量检查**：防止溢出
4. **类型安全**：使用泛型保证类型安全

## 总结

基础数据结构是构建复杂系统的基础，每种数据结构都有其特定的应用场景和性能特征：

1. **数组**: 适合随机访问，内存连续，缓存友好
2. **链表**: 适合频繁插入删除，动态增长，内存分散
3. **栈**: 后进先出，适合函数调用、表达式求值
4. **队列**: 先进先出，适合任务调度、消息队列

选择合适的数据结构需要考虑：

- 访问模式（随机 vs 顺序）
- 操作频率（插入、删除、访问）
- 内存要求（连续 vs 分散）
- 并发需求（单线程 vs 多线程）

通过合理选择和优化，可以构建高效、可靠的基础数据结构，为更复杂的系统提供坚实的基础。
