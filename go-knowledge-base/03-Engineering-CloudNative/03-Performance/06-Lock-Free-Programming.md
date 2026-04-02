# 无锁编程 (Lock-Free Programming)

> **分类**: 工程与云原生
> **标签**: #lock-free #atomic #performance

---

## Atomic 操作

### 基本类型

```go
import "sync/atomic"

var counter int64

// 增加
atomic.AddInt64(&counter, 1)

// 读取
value := atomic.LoadInt64(&counter)

// 写入
atomic.StoreInt64(&counter, 100)

// CAS (Compare-And-Swap)
swapped := atomic.CompareAndSwapInt64(&counter, 100, 200)
```

---

## 无锁队列

```go
type Node struct {
    value interface{}
    next  atomic.Pointer[Node]
}

type LockFreeQueue struct {
    head atomic.Pointer[Node]
    tail atomic.Pointer[Node]
}

func NewLockFreeQueue() *LockFreeQueue {
    dummy := &Node{}
    q := &LockFreeQueue{}
    q.head.Store(dummy)
    q.tail.Store(dummy)
    return q
}

func (q *LockFreeQueue) Enqueue(value interface{}) {
    newNode := &Node{value: value}

    for {
        tail := q.tail.Load()
        next := tail.next.Load()

        if tail == q.tail.Load() {
            if next == nil {
                if tail.next.CompareAndSwap(next, newNode) {
                    q.tail.CompareAndSwap(tail, newNode)
                    return
                }
            } else {
                q.tail.CompareAndSwap(tail, next)
            }
        }
    }
}

func (q *LockFreeQueue) Dequeue() (interface{}, bool) {
    for {
        head := q.head.Load()
        tail := q.tail.Load()
        next := head.next.Load()

        if head == q.head.Load() {
            if head == tail {
                if next == nil {
                    return nil, false
                }
                q.tail.CompareAndSwap(tail, next)
            } else {
                value := next.value
                if q.head.CompareAndSwap(head, next) {
                    return value, true
                }
            }
        }
    }
}
```

---

## 无锁栈

```go
type LockFreeStack struct {
    top atomic.Pointer[Node]
}

func (s *LockFreeStack) Push(value interface{}) {
    newNode := &Node{value: value}

    for {
        oldTop := s.top.Load()
        newNode.next.Store(oldTop)

        if s.top.CompareAndSwap(oldTop, newNode) {
            return
        }
    }
}

func (s *LockFreeStack) Pop() (interface{}, bool) {
    for {
        oldTop := s.top.Load()
        if oldTop == nil {
            return nil, false
        }

        newTop := oldTop.next.Load()
        if s.top.CompareAndSwap(oldTop, newTop) {
            return oldTop.value, true
        }
    }
}
```

---

## RCU (Read-Copy-Update)

```go
type RCU struct {
    data atomic.Pointer[Config]
}

func (r *RCU) Load() *Config {
    return r.data.Load()
}

func (r *RCU) Store(newConfig *Config) {
    oldConfig := r.data.Load()
    r.data.Store(newConfig)

    // 延迟释放旧配置
    go func() {
        time.Sleep(time.Second)  // 等待所有读者
        // 释放 oldConfig
    }()
}
```

---

## 性能对比

| 实现 | 吞吐量 | 延迟 | 适用场景 |
|------|--------|------|----------|
| Mutex | 中等 | 高 | 低竞争 |
| RWMutex | 读高写低 | 中等 | 读多写少 |
| Atomic | 极高 | 极低 | 简单计数器 |
| Lock-Free | 高 | 低 | 高竞争队列 |

---

## 注意事项

1. **ABA 问题**: 使用指针标签或 Hazard Pointers
2. **内存序**: 理解 atomic 的内存序语义
3. **复杂度**: 无锁代码难以验证，谨慎使用
