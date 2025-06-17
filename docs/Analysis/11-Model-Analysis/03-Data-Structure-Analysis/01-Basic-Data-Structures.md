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

### 1. 数学定义

**定义 1.1 (数组)** 数组是一个有序的元素序列，定义为：

$$A = [a_0, a_1, a_2, ..., a_{n-1}]$$

其中 $a_i$ 是第 $i$ 个元素，$n$ 是数组的长度。

**定义 1.2 (数组操作)** 数组支持以下基本操作：

- **访问**: $Access(A, i) = a_i$，时间复杂度 $O(1)$
- **修改**: $Update(A, i, x) = [a_0, ..., a_{i-1}, x, a_{i+1}, ..., a_{n-1}]$，时间复杂度 $O(1)$
- **插入**: $Insert(A, i, x)$，时间复杂度 $O(n)$
- **删除**: $Delete(A, i)$，时间复杂度 $O(n)$

### 2. Golang实现

#### 2.1 基础数组实现

```go
// Array 泛型数组实现
type Array[T any] struct {
    elements []T
    length   int
    capacity int
}

// NewArray 创建新数组
func NewArray[T any](capacity int) *Array[T] {
    return &Array[T]{
        elements: make([]T, capacity),
        length:   0,
        capacity: capacity,
    }
}

// Access 访问元素
func (a *Array[T]) Access(index int) (T, error) {
    if index < 0 || index >= a.length {
        var zero T
        return zero, fmt.Errorf("index %d out of bounds", index)
    }
    return a.elements[index], nil
}

// Update 更新元素
func (a *Array[T]) Update(index int, value T) error {
    if index < 0 || index >= a.length {
        return fmt.Errorf("index %d out of bounds", index)
    }
    a.elements[index] = value
    return nil
}

// Insert 插入元素
func (a *Array[T]) Insert(index int, value T) error {
    if index < 0 || index > a.length {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    // 检查容量
    if a.length >= a.capacity {
        a.resize(a.capacity * 2)
    }
    
    // 移动元素
    for i := a.length; i > index; i-- {
        a.elements[i] = a.elements[i-1]
    }
    
    a.elements[index] = value
    a.length++
    return nil
}

// Delete 删除元素
func (a *Array[T]) Delete(index int) error {
    if index < 0 || index >= a.length {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    // 移动元素
    for i := index; i < a.length-1; i++ {
        a.elements[i] = a.elements[i+1]
    }
    
    a.length--
    
    // 缩容
    if a.length < a.capacity/4 && a.capacity > 10 {
        a.resize(a.capacity / 2)
    }
    
    return nil
}

// resize 调整数组大小
func (a *Array[T]) resize(newCapacity int) {
    newElements := make([]T, newCapacity)
    copy(newElements, a.elements[:a.length])
    a.elements = newElements
    a.capacity = newCapacity
}

// Length 获取数组长度
func (a *Array[T]) Length() int {
    return a.length
}

// Capacity 获取数组容量
func (a *Array[T]) Capacity() int {
    return a.capacity
}

// IsEmpty 检查是否为空
func (a *Array[T]) IsEmpty() bool {
    return a.length == 0
}
```

#### 2.2 动态数组实现

```go
// DynamicArray 动态数组实现
type DynamicArray[T any] struct {
    elements []T
    size     int
}

// NewDynamicArray 创建动态数组
func NewDynamicArray[T any]() *DynamicArray[T] {
    return &DynamicArray[T]{
        elements: make([]T, 0),
        size:     0,
    }
}

// Append 追加元素
func (da *DynamicArray[T]) Append(value T) {
    da.elements = append(da.elements, value)
    da.size++
}

// Get 获取元素
func (da *DynamicArray[T]) Get(index int) (T, error) {
    if index < 0 || index >= da.size {
        var zero T
        return zero, fmt.Errorf("index %d out of bounds", index)
    }
    return da.elements[index], nil
}

// Set 设置元素
func (da *DynamicArray[T]) Set(index int, value T) error {
    if index < 0 || index >= da.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    da.elements[index] = value
    return nil
}

// Remove 移除元素
func (da *DynamicArray[T]) Remove(index int) error {
    if index < 0 || index >= da.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    da.elements = append(da.elements[:index], da.elements[index+1:]...)
    da.size--
    return nil
}

// Size 获取大小
func (da *DynamicArray[T]) Size() int {
    return da.size
}

// Clear 清空数组
func (da *DynamicArray[T]) Clear() {
    da.elements = make([]T, 0)
    da.size = 0
}
```

### 3. 性能分析

#### 3.1 时间复杂度

| 操作 | 平均情况 | 最坏情况 | 说明 |
|------|---------|---------|------|
| 访问 | O(1) | O(1) | 直接索引访问 |
| 修改 | O(1) | O(1) | 直接索引修改 |
| 插入 | O(n) | O(n) | 需要移动元素 |
| 删除 | O(n) | O(n) | 需要移动元素 |
| 追加 | O(1) | O(n) | 动态扩容 |

#### 3.2 空间复杂度

- **静态数组**: O(n) - 固定大小
- **动态数组**: O(n) - 可扩容，但存在空间浪费

#### 3.3 内存布局

```go
// 数组内存布局分析
type ArrayLayout struct {
    // 数组头部 (24 bytes on 64-bit)
    ptr      *byte  // 8 bytes - 指向数据
    len      int    // 8 bytes - 长度
    cap      int    // 8 bytes - 容量
    
    // 数据部分 (连续内存)
    data     []byte // n * sizeof(T) bytes
}

// 内存对齐示例
type AlignedArray struct {
    a bool   // 1 byte + 7 bytes padding
    b int64  // 8 bytes
    c int32  // 4 bytes + 4 bytes padding
}
```

## 链表 (Linked List)

### 1. 数学定义

**定义 2.1 (链表)** 链表是一个由节点组成的序列，每个节点包含数据和指向下一个节点的指针：

$$L = n_0 \rightarrow n_1 \rightarrow n_2 \rightarrow ... \rightarrow n_{k-1} \rightarrow nil$$

其中 $n_i = (data_i, next_i)$，$next_i$ 指向 $n_{i+1}$。

**定义 2.2 (链表操作)** 链表支持以下基本操作：

- **访问**: $Access(L, i)$，时间复杂度 $O(i)$
- **插入**: $Insert(L, i, x)$，时间复杂度 $O(i)$
- **删除**: $Delete(L, i)$，时间复杂度 $O(i)$
- **搜索**: $Search(L, x)$，时间复杂度 $O(n)$

### 2. Golang实现

#### 2.1 单链表实现

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
    mutex sync.RWMutex
}

// NewLinkedList 创建新链表
func NewLinkedList[T any]() *LinkedList[T] {
    return &LinkedList[T]{}
}

// InsertAt 在指定位置插入
func (ll *LinkedList[T]) InsertAt(index int, data T) error {
    ll.mutex.Lock()
    defer ll.mutex.Unlock()
    
    if index < 0 || index > ll.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    newNode := &Node[T]{Data: data}
    
    if index == 0 {
        newNode.Next = ll.head
        ll.head = newNode
    } else {
        current := ll.head
        for i := 0; i < index-1; i++ {
            current = current.Next
        }
        newNode.Next = current.Next
        current.Next = newNode
    }
    
    ll.size++
    return nil
}

// DeleteAt 删除指定位置的元素
func (ll *LinkedList[T]) DeleteAt(index int) error {
    ll.mutex.Lock()
    defer ll.mutex.Unlock()
    
    if index < 0 || index >= ll.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    if index == 0 {
        ll.head = ll.head.Next
    } else {
        current := ll.head
        for i := 0; i < index-1; i++ {
            current = current.Next
        }
        current.Next = current.Next.Next
    }
    
    ll.size--
    return nil
}

// Get 获取指定位置的元素
func (ll *LinkedList[T]) Get(index int) (T, error) {
    ll.mutex.RLock()
    defer ll.mutex.RUnlock()
    
    if index < 0 || index >= ll.size {
        var zero T
        return zero, fmt.Errorf("index %d out of bounds", index)
    }
    
    current := ll.head
    for i := 0; i < index; i++ {
        current = current.Next
    }
    
    return current.Data, nil
}

// Set 设置指定位置的元素
func (ll *LinkedList[T]) Set(index int, data T) error {
    ll.mutex.Lock()
    defer ll.mutex.Unlock()
    
    if index < 0 || index >= ll.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    current := ll.head
    for i := 0; i < index; i++ {
        current = current.Next
    }
    
    current.Data = data
    return nil
}

// Size 获取链表大小
func (ll *LinkedList[T]) Size() int {
    ll.mutex.RLock()
    defer ll.mutex.RUnlock()
    return ll.size
}

// IsEmpty 检查是否为空
func (ll *LinkedList[T]) IsEmpty() bool {
    return ll.Size() == 0
}

// Clear 清空链表
func (ll *LinkedList[T]) Clear() {
    ll.mutex.Lock()
    defer ll.mutex.Unlock()
    ll.head = nil
    ll.size = 0
}
```

#### 2.2 双向链表实现

```go
// DNode 双向链表节点
type DNode[T any] struct {
    Data T
    Prev *DNode[T]
    Next *DNode[T]
}

// DoublyLinkedList 双向链表
type DoublyLinkedList[T any] struct {
    head *DNode[T]
    tail *DNode[T]
    size int
    mutex sync.RWMutex
}

// NewDoublyLinkedList 创建双向链表
func NewDoublyLinkedList[T any]() *DoublyLinkedList[T] {
    return &DoublyLinkedList[T]{}
}

// InsertAt 在指定位置插入
func (dll *DoublyLinkedList[T]) InsertAt(index int, data T) error {
    dll.mutex.Lock()
    defer dll.mutex.Unlock()
    
    if index < 0 || index > dll.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    newNode := &DNode[T]{Data: data}
    
    if index == 0 {
        // 插入到头部
        newNode.Next = dll.head
        if dll.head != nil {
            dll.head.Prev = newNode
        }
        dll.head = newNode
        if dll.tail == nil {
            dll.tail = newNode
        }
    } else if index == dll.size {
        // 插入到尾部
        newNode.Prev = dll.tail
        if dll.tail != nil {
            dll.tail.Next = newNode
        }
        dll.tail = newNode
        if dll.head == nil {
            dll.head = newNode
        }
    } else {
        // 插入到中间
        current := dll.head
        for i := 0; i < index; i++ {
            current = current.Next
        }
        
        newNode.Prev = current.Prev
        newNode.Next = current
        current.Prev.Next = newNode
        current.Prev = newNode
    }
    
    dll.size++
    return nil
}

// DeleteAt 删除指定位置的元素
func (dll *DoublyLinkedList[T]) DeleteAt(index int) error {
    dll.mutex.Lock()
    defer dll.mutex.Unlock()
    
    if index < 0 || index >= dll.size {
        return fmt.Errorf("index %d out of bounds", index)
    }
    
    if index == 0 {
        // 删除头部
        dll.head = dll.head.Next
        if dll.head != nil {
            dll.head.Prev = nil
        } else {
            dll.tail = nil
        }
    } else if index == dll.size-1 {
        // 删除尾部
        dll.tail = dll.tail.Prev
        if dll.tail != nil {
            dll.tail.Next = nil
        } else {
            dll.head = nil
        }
    } else {
        // 删除中间节点
        current := dll.head
        for i := 0; i < index; i++ {
            current = current.Next
        }
        
        current.Prev.Next = current.Next
        current.Next.Prev = current.Prev
    }
    
    dll.size--
    return nil
}

// Get 获取指定位置的元素
func (dll *DoublyLinkedList[T]) Get(index int) (T, error) {
    dll.mutex.RLock()
    defer dll.mutex.RUnlock()
    
    if index < 0 || index >= dll.size {
        var zero T
        return zero, fmt.Errorf("index %d out of bounds", index)
    }
    
    // 优化：从头部或尾部开始遍历
    if index < dll.size/2 {
        current := dll.head
        for i := 0; i < index; i++ {
            current = current.Next
        }
        return current.Data, nil
    } else {
        current := dll.tail
        for i := dll.size - 1; i > index; i-- {
            current = current.Prev
        }
        return current.Data, nil
    }
}
```

### 3. 性能分析

#### 3.1 时间复杂度

| 操作 | 单链表 | 双向链表 | 说明 |
|------|--------|---------|------|
| 访问 | O(n) | O(n) | 需要遍历 |
| 插入头部 | O(1) | O(1) | 直接操作 |
| 插入尾部 | O(n) | O(1) | 双向链表有尾指针 |
| 插入中间 | O(n) | O(n) | 需要遍历 |
| 删除头部 | O(1) | O(1) | 直接操作 |
| 删除尾部 | O(n) | O(1) | 双向链表有尾指针 |
| 删除中间 | O(n) | O(n) | 需要遍历 |

#### 3.2 空间复杂度

- **单链表**: O(n) - 每个节点需要一个指针
- **双向链表**: O(n) - 每个节点需要两个指针

## 栈 (Stack)

### 1. 数学定义

**定义 3.1 (栈)** 栈是一个后进先出(LIFO)的抽象数据类型，定义为：

$$Stack = (D, O, A)$$

其中：

- $D = \{s | s \text{ 是元素序列}\}$
- $O = \{push, pop, top, empty, size\}$
- $A$ 包含以下公理：

$$
\begin{align}
empty(new()) &= true \\
empty(push(s, x)) &= false \\
top(push(s, x)) &= x \\
pop(push(s, x)) &= s \\
size(new()) &= 0 \\
size(push(s, x)) &= size(s) + 1
\end{align}
$$

### 2. Golang实现

#### 2.1 基于数组的栈

```go
// ArrayStack 基于数组的栈实现
type ArrayStack[T any] struct {
    elements []T
    top      int
    mutex    sync.RWMutex
}

// NewArrayStack 创建新栈
func NewArrayStack[T any](capacity int) *ArrayStack[T] {
    return &ArrayStack[T]{
        elements: make([]T, capacity),
        top:      -1,
    }
}

// Push 入栈
func (as *ArrayStack[T]) Push(element T) error {
    as.mutex.Lock()
    defer as.mutex.Unlock()

    if as.top >= len(as.elements)-1 {
        // 扩容
        newElements := make([]T, len(as.elements)*2)
        copy(newElements, as.elements)
        as.elements = newElements
    }

    as.top++
    as.elements[as.top] = element
    return nil
}

// Pop 出栈
func (as *ArrayStack[T]) Pop() (T, error) {
    as.mutex.Lock()
    defer as.mutex.Unlock()

    var zero T
    if as.IsEmpty() {
        return zero, errors.New("stack is empty")
    }

    element := as.elements[as.top]
    as.top--
    return element, nil
}

// Top 查看栈顶元素
func (as *ArrayStack[T]) Top() (T, error) {
    as.mutex.RLock()
    defer as.mutex.RUnlock()

    var zero T
    if as.IsEmpty() {
        return zero, errors.New("stack is empty")
    }

    return as.elements[as.top], nil
}

// IsEmpty 检查是否为空
func (as *ArrayStack[T]) IsEmpty() bool {
    return as.top == -1
}

// Size 获取栈大小
func (as *ArrayStack[T]) Size() int {
    as.mutex.RLock()
    defer as.mutex.RUnlock()
    return as.top + 1
}

// Clear 清空栈
func (as *ArrayStack[T]) Clear() {
    as.mutex.Lock()
    defer as.mutex.Unlock()
    as.top = -1
}
```

#### 2.2 基于链表的栈

```go
// LinkedStack 基于链表的栈实现
type LinkedStack[T any] struct {
    head *Node[T]
    size int
    mutex sync.RWMutex
}

// NewLinkedStack 创建新栈
func NewLinkedStack[T any]() *LinkedStack[T] {
    return &LinkedStack[T]{}
}

// Push 入栈
func (ls *LinkedStack[T]) Push(element T) {
    ls.mutex.Lock()
    defer ls.mutex.Unlock()

    newNode := &Node[T]{
        Data: element,
        Next: ls.head,
    }
    ls.head = newNode
    ls.size++
}

// Pop 出栈
func (ls *LinkedStack[T]) Pop() (T, error) {
    ls.mutex.Lock()
    defer ls.mutex.Unlock()

    var zero T
    if ls.IsEmpty() {
        return zero, errors.New("stack is empty")
    }

    element := ls.head.Data
    ls.head = ls.head.Next
    ls.size--
    return element, nil
}

// Top 查看栈顶元素
func (ls *LinkedStack[T]) Top() (T, error) {
    ls.mutex.RLock()
    defer ls.mutex.RUnlock()

    var zero T
    if ls.IsEmpty() {
        return zero, errors.New("stack is empty")
    }

    return ls.head.Data, nil
}

// IsEmpty 检查是否为空
func (ls *LinkedStack[T]) IsEmpty() bool {
    return ls.head == nil
}

// Size 获取栈大小
func (ls *LinkedStack[T]) Size() int {
    ls.mutex.RLock()
    defer ls.mutex.RUnlock()
    return ls.size
}
```

### 3. 应用场景

#### 3.1 函数调用栈

```go
// 函数调用栈模拟
type CallFrame struct {
    FunctionName string
    Parameters   []interface{}
    ReturnValue  interface{}
    LocalVars    map[string]interface{}
}

type CallStack struct {
    frames []*CallFrame
}

func (cs *CallStack) Push(frame *CallFrame) {
    cs.frames = append(cs.frames, frame)
}

func (cs *CallStack) Pop() *CallFrame {
    if len(cs.frames) == 0 {
        return nil
    }
    frame := cs.frames[len(cs.frames)-1]
    cs.frames = cs.frames[:len(cs.frames)-1]
    return frame
}

func (cs *CallStack) Top() *CallFrame {
    if len(cs.frames) == 0 {
        return nil
    }
    return cs.frames[len(cs.frames)-1]
}
```

#### 3.2 表达式求值

```go
// 中缀表达式求值
func EvaluateInfixExpression(expression string) (float64, error) {
    stack := NewArrayStack[float64](100)
    operatorStack := NewArrayStack[string](100)

    tokens := strings.Fields(expression)

    for _, token := range tokens {
        switch token {
        case "+", "-", "*", "/":
            for !operatorStack.IsEmpty() {
                top, _ := operatorStack.Top()
                if precedence(top) >= precedence(token) {
                    op, _ := operatorStack.Pop()
                    b, _ := stack.Pop()
                    a, _ := stack.Pop()
                    result := applyOperator(a, b, op)
                    stack.Push(result)
                } else {
                    break
                }
            }
            operatorStack.Push(token)
        case "(":
            operatorStack.Push(token)
        case ")":
            for !operatorStack.IsEmpty() {
                op, _ := operatorStack.Top()
                if op == "(" {
                    operatorStack.Pop()
                    break
                }
                operatorStack.Pop()
                b, _ := stack.Pop()
                a, _ := stack.Pop()
                result := applyOperator(a, b, op)
                stack.Push(result)
            }
        default:
            if num, err := strconv.ParseFloat(token, 64); err == nil {
                stack.Push(num)
            }
        }
    }

    for !operatorStack.IsEmpty() {
        op, _ := operatorStack.Pop()
        b, _ := stack.Pop()
        a, _ := stack.Pop()
        result := applyOperator(a, b, op)
        stack.Push(result)
    }

    result, _ := stack.Pop()
    return result, nil
}

func precedence(operator string) int {
    switch operator {
    case "+", "-":
        return 1
    case "*", "/":
        return 2
    default:
        return 0
    }
}

func applyOperator(a, b float64, operator string) float64 {
    switch operator {
    case "+":
        return a + b
    case "-":
        return a - b
    case "*":
        return a * b
    case "/":
        return a / b
    default:
        return 0
    }
}
```

## 队列 (Queue)

### 1. 数学定义

**定义 4.1 (队列)** 队列是一个先进先出(FIFO)的抽象数据类型，定义为：

$$Queue = (D, O, A)$$

其中：

- $D = \{q | q \text{ 是元素序列}\}$
- $O = \{enqueue, dequeue, front, empty, size\}$
- $A$ 包含以下公理：

$$
\begin{align}
empty(new()) &= true \\
empty(enqueue(q, x)) &= false \\
front(enqueue(q, x)) &= \begin{cases}
x & \text{if } empty(q) \\
front(q) & \text{otherwise}
\end{cases} \\
dequeue(enqueue(q, x)) &= \begin{cases}
new() & \text{if } empty(q) \\
enqueue(dequeue(q), x) & \text{otherwise}
\end{cases}
\end{align}
$$

### 2. Golang实现

#### 2.1 基于数组的队列

```go
// ArrayQueue 基于数组的队列实现
type ArrayQueue[T any] struct {
    elements []T
    front    int
    rear     int
    size     int
    capacity int
    mutex    sync.RWMutex
}

// NewArrayQueue 创建新队列
func NewArrayQueue[T any](capacity int) *ArrayQueue[T] {
    return &ArrayQueue[T]{
        elements: make([]T, capacity),
        front:    0,
        rear:     -1,
        size:     0,
        capacity: capacity,
    }
}

// Enqueue 入队
func (aq *ArrayQueue[T]) Enqueue(element T) error {
    aq.mutex.Lock()
    defer aq.mutex.Unlock()

    if aq.IsFull() {
        return errors.New("queue is full")
    }

    aq.rear = (aq.rear + 1) % aq.capacity
    aq.elements[aq.rear] = element
    aq.size++
    return nil
}

// Dequeue 出队
func (aq *ArrayQueue[T]) Dequeue() (T, error) {
    aq.mutex.Lock()
    defer aq.mutex.Unlock()

    var zero T
    if aq.IsEmpty() {
        return zero, errors.New("queue is empty")
    }

    element := aq.elements[aq.front]
    aq.front = (aq.front + 1) % aq.capacity
    aq.size--
    return element, nil
}

// Front 查看队首元素
func (aq *ArrayQueue[T]) Front() (T, error) {
    aq.mutex.RLock()
    defer aq.mutex.RUnlock()

    var zero T
    if aq.IsEmpty() {
        return zero, errors.New("queue is empty")
    }

    return aq.elements[aq.front], nil
}

// IsEmpty 检查是否为空
func (aq *ArrayQueue[T]) IsEmpty() bool {
    return aq.size == 0
}

// IsFull 检查是否已满
func (aq *ArrayQueue[T]) IsFull() bool {
    return aq.size == aq.capacity
}

// Size 获取队列大小
func (aq *ArrayQueue[T]) Size() int {
    aq.mutex.RLock()
    defer aq.mutex.RUnlock()
    return aq.size
}
```

#### 2.2 基于链表的队列

```go
// LinkedQueue 基于链表的队列实现
type LinkedQueue[T any] struct {
    head *Node[T]
    tail *Node[T]
    size int
    mutex sync.RWMutex
}

// NewLinkedQueue 创建新队列
func NewLinkedQueue[T any]() *LinkedQueue[T] {
    return &LinkedQueue[T]{}
}

// Enqueue 入队
func (lq *LinkedQueue[T]) Enqueue(element T) {
    lq.mutex.Lock()
    defer lq.mutex.Unlock()

    newNode := &Node[T]{Data: element}

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
    lq.mutex.Lock()
    defer lq.mutex.Unlock()

    var zero T
    if lq.IsEmpty() {
        return zero, errors.New("queue is empty")
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
    lq.mutex.RLock()
    defer lq.mutex.RUnlock()

    var zero T
    if lq.IsEmpty() {
        return zero, errors.New("queue is empty")
    }

    return lq.head.Data, nil
}

// IsEmpty 检查是否为空
func (lq *LinkedQueue[T]) IsEmpty() bool {
    return lq.head == nil
}

// Size 获取队列大小
func (lq *LinkedQueue[T]) Size() int {
    lq.mutex.RLock()
    defer lq.mutex.RUnlock()
    return lq.size
}
```

### 3. 应用场景

#### 3.1 任务调度

```go
// 任务调度器
type Task struct {
    ID       string
    Priority int
    Data     interface{}
    Handler  func(interface{}) error
}

type TaskScheduler struct {
    queue *ArrayQueue[*Task]
    mutex sync.RWMutex
}

func NewTaskScheduler(capacity int) *TaskScheduler {
    return &TaskScheduler{
        queue: NewArrayQueue[*Task](capacity),
    }
}

func (ts *TaskScheduler) AddTask(task *Task) error {
    return ts.queue.Enqueue(task)
}

func (ts *TaskScheduler) ProcessNextTask() error {
    task, err := ts.queue.Dequeue()
    if err != nil {
        return err
    }

    return task.Handler(task.Data)
}

func (ts *TaskScheduler) GetPendingTaskCount() int {
    return ts.queue.Size()
}
```

#### 3.2 消息队列

```go
// 消息队列
type Message struct {
    ID      string
    Topic   string
    Payload []byte
    Time    time.Time
}

type MessageQueue struct {
    messages *LinkedQueue[*Message]
    mutex    sync.RWMutex
}

func NewMessageQueue() *MessageQueue {
    return &MessageQueue{
        messages: NewLinkedQueue[*Message](),
    }
}

func (mq *MessageQueue) Publish(message *Message) {
    mq.messages.Enqueue(message)
}

func (mq *MessageQueue) Consume() (*Message, error) {
    return mq.messages.Dequeue()
}

func (mq *MessageQueue) GetMessageCount() int {
    return mq.messages.Size()
}
```

## 性能对比

### 1. 时间复杂度对比

| 操作 | 数组 | 链表 | 栈(数组) | 栈(链表) | 队列(数组) | 队列(链表) |
|------|------|------|----------|----------|------------|------------|
| 访问 | O(1) | O(n) | O(1) | O(1) | O(n) | O(n) |
| 插入头部 | O(n) | O(1) | O(1) | O(1) | O(n) | O(1) |
| 插入尾部 | O(1) | O(n) | O(1) | O(1) | O(1) | O(1) |
| 删除头部 | O(n) | O(1) | O(1) | O(1) | O(1) | O(1) |
| 删除尾部 | O(1) | O(n) | O(1) | O(1) | O(n) | O(1) |

### 2. 空间复杂度对比

| 数据结构 | 空间复杂度 | 额外开销 | 说明 |
|---------|-----------|---------|------|
| 数组 | O(n) | 低 | 连续存储，无指针开销 |
| 链表 | O(n) | 高 | 每个节点需要指针 |
| 栈(数组) | O(n) | 低 | 动态扩容 |
| 栈(链表) | O(n) | 高 | 每个节点需要指针 |
| 队列(数组) | O(n) | 低 | 循环数组 |
| 队列(链表) | O(n) | 高 | 每个节点需要指针 |

### 3. 缓存性能对比

```go
// 缓存性能测试
func BenchmarkArrayAccess(b *testing.B) {
    arr := make([]int, 1000)
    for i := 0; i < 1000; i++ {
        arr[i] = i
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = arr[i%1000]
    }
}

func BenchmarkLinkedListAccess(b *testing.B) {
    list := NewLinkedList[int]()
    for i := 0; i < 1000; i++ {
        list.InsertAt(i, i)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        list.Get(i % 1000)
    }
}
```

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

### 1. 选择原则

#### 1.1 根据访问模式选择

```go
// 随机访问频繁 - 选择数组
func RandomAccessExample() {
    data := make([]int, 1000)
    for i := 0; i < 1000; i++ {
        data[i] = i
    }

    // 随机访问
    for i := 0; i < 1000; i++ {
        idx := rand.Intn(1000)
        _ = data[idx]
    }
}

// 频繁插入删除 - 选择链表
func FrequentInsertDeleteExample() {
    list := NewLinkedList[int]()

    // 频繁插入删除
    for i := 0; i < 1000; i++ {
        list.InsertAt(0, i)
        list.DeleteAt(0)
    }
}
```

#### 1.2 根据内存要求选择

```go
// 内存敏感 - 选择数组
func MemorySensitiveExample() {
    // 数组：连续内存，缓存友好
    data := make([]int, 1000000)
    for i := 0; i < 1000000; i++ {
        data[i] = i
    }
}

// 内存不敏感 - 选择链表
func MemoryInsensitiveExample() {
    // 链表：分散内存，但动态增长
    list := NewLinkedList[int]()
    for i := 0; i < 1000000; i++ {
        list.InsertAt(i, i)
    }
}
```

### 2. 性能优化

#### 2.1 预分配内存

```go
// 预分配内存避免频繁扩容
func PreallocateExample() {
    // 预分配容量
    stack := NewArrayStack[int](1000)
    for i := 0; i < 1000; i++ {
        stack.Push(i)
    }
}
```

#### 2.2 使用对象池

```go
// 使用对象池减少GC压力
var nodePool = sync.Pool{
    New: func() interface{} {
        return &Node[int]{}
    },
}

func GetNode() *Node[int] {
    return nodePool.Get().(*Node[int])
}

func PutNode(node *Node[int]) {
    node.Data = 0
    node.Next = nil
    nodePool.Put(node)
}
```

### 3. 并发安全

#### 3.1 使用适当的锁

```go
// 读写锁优化
type OptimizedStack[T any] struct {
    elements []T
    top      int
    mutex    sync.RWMutex
}

func (os *OptimizedStack[T]) Top() (T, error) {
    os.mutex.RLock() // 读锁
    defer os.mutex.RUnlock()

    var zero T
    if os.top == -1 {
        return zero, errors.New("stack is empty")
    }
    return os.elements[os.top], nil
}
```

#### 3.2 无锁数据结构

```go
// 无锁栈实现
type LockFreeStack[T any] struct {
    head unsafe.Pointer
}

type node[T any] struct {
    value T
    next  unsafe.Pointer
}

func (lfs *LockFreeStack[T]) Push(value T) {
    newNode := &node[T]{value: value}
    for {
        oldHead := lfs.head
        newNode.next = oldHead
        if atomic.CompareAndSwapPointer(&lfs.head, oldHead, unsafe.Pointer(newNode)) {
            break
        }
    }
}

func (lfs *LockFreeStack[T]) Pop() (T, bool) {
    for {
        oldHead := lfs.head
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
```

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
