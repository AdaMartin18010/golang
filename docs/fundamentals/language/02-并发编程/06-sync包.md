# sync包与并发安全模式

> **简介**: 详解Go标准库sync包的并发原语，包括Mutex、RWMutex、WaitGroup、Once等
> **版本**: Go 1.23+  
> **难度**: ⭐⭐⭐  
> **标签**: #并发 #sync #锁 #并发安全

<!-- TOC START -->
- [sync包与并发安全模式](#sync包与并发安全模式)
  - [1. 理论基础](#1-理论基础)
  - [2. 典型用法](#2-典型用法)
    - [互斥锁Mutex](#互斥锁mutex)
      - [Mutex状态机可视化](#mutex状态机可视化)
      - [Mutex并发访问时序图](#mutex并发访问时序图)
    - [读写锁RWMutex](#读写锁rwmutex)
      - [RWMutex并发控制可视化](#rwmutex并发控制可视化)
      - [RWMutex读写时序图](#rwmutex读写时序图)
    - [WaitGroup](#waitgroup)
      - [WaitGroup工作流程](#waitgroup工作流程)
      - [WaitGroup时序图](#waitgroup时序图)
    - [Once](#once)
      - [sync.Once单次执行保证](#synconce单次执行保证)
      - [多Goroutine调用Once时序图](#多goroutine调用once时序图)
  - [3. 工程分析与最佳实践](#3-工程分析与最佳实践)
  - [4. 常见陷阱](#4-常见陷阱)
  - [5. 单元测试建议](#5-单元测试建议)
  - [6. 参考文献](#6-参考文献)
<!-- TOC END -->


## 📋 目录


- [1. 理论基础](#1-理论基础)
- [2. 典型用法](#2-典型用法)
  - [互斥锁Mutex](#互斥锁mutex)
    - [Mutex状态机可视化](#mutex状态机可视化)
    - [Mutex并发访问时序图](#mutex并发访问时序图)
  - [读写锁RWMutex](#读写锁rwmutex)
    - [RWMutex并发控制可视化](#rwmutex并发控制可视化)
    - [RWMutex读写时序图](#rwmutex读写时序图)
  - [WaitGroup](#waitgroup)
    - [WaitGroup工作流程](#waitgroup工作流程)
    - [WaitGroup时序图](#waitgroup时序图)
  - [Once](#once)
    - [sync.Once单次执行保证](#synconce单次执行保证)
    - [多Goroutine调用Once时序图](#多goroutine调用once时序图)
- [3. 工程分析与最佳实践](#3-工程分析与最佳实践)
- [4. 常见陷阱](#4-常见陷阱)
- [5. 单元测试建议](#5-单元测试建议)
- [6. 参考文献](#6-参考文献)

## 1. 理论基础

Go的sync包提供了多种并发原语，保障多Goroutine环境下的数据一致性和同步。

- **互斥锁（Mutex）**：保证同一时刻只有一个Goroutine访问临界区。
- **读写锁（RWMutex）**：读操作可并发，写操作独占。
- **等待组（WaitGroup）**：用于等待一组Goroutine完成。
- **Once**：确保某段代码只执行一次。
- **Cond**：条件变量，支持复杂同步。

---

## 2. 典型用法

### 互斥锁Mutex

#### Mutex状态机可视化

```mermaid
stateDiagram-v2
    [*] --> Unlocked: 初始状态
    
    Unlocked --> Locked: Goroutine A<br/>调用 Lock()
    
    state Locked {
        [*] --> Executing: 获得锁
        Executing --> [*]: 执行临界区代码
    }
    
    Locked --> Unlocked: Goroutine A<br/>调用 Unlock()
    
    state Locked_Contention {
        direction LR
        [*] --> WaitingQueue: Goroutine B, C, D<br/>等待获取锁
        WaitingQueue --> GotLock: 锁释放后<br/>按FIFO唤醒
    }
    
    Unlocked --> Locked_Contention: 多个Goroutine<br/>竞争
    Locked_Contention --> Unlocked: 最后一个<br/>释放锁
    
    note right of Locked
        临界区:
        - 只有一个Goroutine执行
        - 其他Goroutine阻塞等待
    end note
    
    note left of Locked_Contention
        锁竞争:
        - 等待队列 (FIFO)
        - 公平模式 vs 非公平模式
    end note
```

#### Mutex并发访问时序图

```mermaid
sequenceDiagram
    participant G1 as Goroutine 1
    participant Mutex as sync.Mutex
    participant G2 as Goroutine 2
    participant G3 as Goroutine 3
    
    Note over G1,G3: 初始状态：Mutex未锁定
    
    G1->>Mutex: Lock()
    Mutex-->>G1: 获得锁 ✓
    Note over G1: 进入临界区
    
    G2->>Mutex: Lock()
    Note over G2: ⏸️ 阻塞等待
    
    G3->>Mutex: Lock()
    Note over G3: ⏸️ 阻塞等待
    
    Note over Mutex: 等待队列: [G2, G3]
    
    G1->>G1: 执行临界区代码
    G1->>Mutex: Unlock()
    Mutex-->>G2: 唤醒 G2 ✓
    Note over G2: 获得锁，进入临界区
    Note over G3: 仍在等待
    
    G2->>G2: 执行临界区代码
    G2->>Mutex: Unlock()
    Mutex-->>G3: 唤醒 G3 ✓
    Note over G3: 获得锁，进入临界区
    
    G3->>G3: 执行临界区代码
    G3->>Mutex: Unlock()
    Note over Mutex: Mutex回到未锁定状态
```

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    counter int
    mu      sync.Mutex
)

func increment(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    mu.Lock()           // 获取锁
    defer mu.Unlock()   // 确保释放锁
    
    // 临界区：修改共享资源
    temp := counter
    time.Sleep(10 * time.Millisecond) // 模拟处理时间
    counter = temp + 1
    fmt.Printf("Goroutine %d: counter = %d\n", id, counter)
}

func main() {
    var wg sync.WaitGroup
    
    // 启动10个并发Goroutine
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go increment(i, &wg)
    }
    
    wg.Wait()
    fmt.Printf("Final counter: %d\n", counter) // 输出: 10
}
```

### 读写锁RWMutex

#### RWMutex并发控制可视化

```mermaid
graph TB
    subgraph "读写锁状态"
        Idle[空闲状态<br/>RWMutex]
        Reading[读模式<br/>多个Reader并发]
        Writing[写模式<br/>独占访问]
    end
    
    subgraph "Reader Goroutines"
        R1[Reader 1]
        R2[Reader 2]
        R3[Reader 3]
    end
    
    subgraph "Writer Goroutines"
        W1[Writer 1]
        W2[Writer 2]
    end
    
    Idle -->|RLock| Reading
    Reading -->|所有RUnlock| Idle
    Idle -->|Lock| Writing
    Writing -->|Unlock| Idle
    
    R1 -->|RLock| Reading
    R2 -->|RLock| Reading
    R3 -->|RLock| Reading
    
    W1 -->|Lock - 等待所有Reader完成| Writing
    W2 -->|Lock - 进入等待队列| Writing
    
    Reading -.阻塞.-> W1
    Writing -.阻塞.-> R1
    Writing -.阻塞.-> W2
    
    style Reading fill:#e1ffe1
    style Writing fill:#ffe1e1
    style Idle fill:#e1f5ff
```

#### RWMutex读写时序图

```mermaid
sequenceDiagram
    participant R1 as Reader 1
    participant R2 as Reader 2
    participant RW as RWMutex
    participant W1 as Writer 1
    participant R3 as Reader 3
    
    Note over R1,R3: 场景：多个读并发，写独占
    
    R1->>RW: RLock()
    RW-->>R1: 获得读锁 ✓
    Note over R1: 读取数据
    
    R2->>RW: RLock()
    RW-->>R2: 获得读锁 ✓ (并发)
    Note over R2: 读取数据
    Note over RW: 读计数: 2
    
    W1->>RW: Lock()
    Note over W1: ⏸️ 阻塞 (等待所有Reader完成)
    
    R3->>RW: RLock()
    Note over R3: ⏸️ 阻塞 (Writer等待中，不允许新Reader)
    
    R1->>RW: RUnlock()
    Note over RW: 读计数: 1
    
    R2->>RW: RUnlock()
    Note over RW: 读计数: 0
    
    RW-->>W1: 获得写锁 ✓
    Note over W1: 独占写入
    
    W1->>RW: Unlock()
    RW-->>R3: 获得读锁 ✓
    Note over R3: 读取数据
    
    R3->>RW: RUnlock()
```

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type SafeMap struct {
    data map[string]int
    rw   sync.RWMutex
}

// 读操作：可并发
func (m *SafeMap) Get(key string) (int, bool) {
    m.rw.RLock()
    defer m.rw.RUnlock()
    
    val, ok := m.data[key]
    return val, ok
}

// 写操作：独占访问
func (m *SafeMap) Set(key string, value int) {
    m.rw.Lock()
    defer m.rw.Unlock()
    
    m.data[key] = value
}

func main() {
    sm := &SafeMap{
        data: make(map[string]int),
    }
    
    // 多个Reader并发读取
    for i := 0; i < 5; i++ {
        go func(id int) {
            for j := 0; j < 3; j++ {
                val, ok := sm.Get("key")
                fmt.Printf("Reader %d: %v, %v\n", id, val, ok)
                time.Sleep(10 * time.Millisecond)
            }
        }(i)
    }
    
    // 单个Writer写入
    go func() {
        for i := 0; i < 3; i++ {
            sm.Set("key", i)
            fmt.Printf("Writer: set key = %d\n", i)
            time.Sleep(50 * time.Millisecond)
        }
    }()
    
    time.Sleep(300 * time.Millisecond)
}
```

### WaitGroup

#### WaitGroup工作流程

```mermaid
flowchart TB
    Start([主Goroutine]) --> Init[var wg sync.WaitGroup<br/>counter = 0]
    Init --> Add1[wg.Add3<br/>counter = 3]
    Add1 --> Launch1[启动 Goroutine 1]
    Add1 --> Launch2[启动 Goroutine 2]
    Add1 --> Launch3[启动 Goroutine 3]
    
    Launch1 --> G1[Goroutine 1执行]
    Launch2 --> G2[Goroutine 2执行]
    Launch3 --> G3[Goroutine 3执行]
    
    G1 --> Done1[wg.Done<br/>counter = 2]
    G2 --> Done2[wg.Done<br/>counter = 1]
    G3 --> Done3[wg.Done<br/>counter = 0]
    
    Add1 --> Wait[wg.Wait<br/>阻塞等待]
    
    Done1 --> Check1{counter == 0?}
    Check1 -->|否| Wait
    
    Done2 --> Check2{counter == 0?}
    Check2 -->|否| Wait
    
    Done3 --> Check3{counter == 0?}
    Check3 -->|是| Unblock[唤醒主Goroutine]
    
    Unblock --> End([所有任务完成])
    
    style Init fill:#e1f5ff
    style Wait fill:#ffe1e1
    style Unblock fill:#e1ffe1
    style End fill:#e1ffe1
```

#### WaitGroup时序图

```mermaid
sequenceDiagram
    participant Main as 主Goroutine
    participant WG as WaitGroup<br/>(counter)
    participant G1 as Goroutine 1
    participant G2 as Goroutine 2
    participant G3 as Goroutine 3
    
    Main->>WG: Add(3)
    Note over WG: counter = 3
    
    Main->>G1: 启动 go func()
    Main->>G2: 启动 go func()
    Main->>G3: 启动 go func()
    
    Main->>WG: Wait()
    Note over Main: ⏸️ 阻塞等待
    
    par 并发执行
        G1->>G1: 执行任务1
        G2->>G2: 执行任务2
        G3->>G3: 执行任务3
    end
    
    G1->>WG: Done()
    Note over WG: counter = 2
    
    G2->>WG: Done()
    Note over WG: counter = 1
    
    G3->>WG: Done()
    Note over WG: counter = 0
    
    WG-->>Main: 唤醒主Goroutine ✓
    Note over Main: 继续执行
```

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // 确保Done被调用
    
    fmt.Printf("Worker %d: 开始工作\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d: 完成工作\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    // 启动5个worker
    for i := 1; i <= 5; i++ {
        wg.Add(1) // 每启动一个Goroutine，计数器+1
        go worker(i, &wg)
    }
    
    fmt.Println("主Goroutine: 等待所有worker完成...")
    wg.Wait() // 阻塞，直到计数器为0
    fmt.Println("主Goroutine: 所有worker已完成！")
}
```

### Once

#### sync.Once单次执行保证

```mermaid
stateDiagram-v2
    [*] --> NotExecuted: 初始状态<br/>done = 0
    
    NotExecuted --> Executing: 第一个Goroutine<br/>调用 Do(f)
    
    state Executing {
        [*] --> RunningFunc: 执行 f()
        RunningFunc --> SetDone: 设置 done = 1
        SetDone --> [*]
    }
    
    Executing --> Executed: 执行完成
    
    NotExecuted --> Blocked: 其他Goroutine<br/>调用 Do(f)
    Blocked --> WaitForCompletion: 等待第一个<br/>Goroutine完成
    WaitForCompletion --> Executed
    
    Executed --> Executed: 后续调用 Do(f)<br/>直接返回，不执行
    
    note right of Executing
        关键特性:
        - 只有第一个调用执行函数
        - 其他调用阻塞等待
        - 执行完成后，done=1
    end note
    
    note left of Executed
        已执行状态:
        - 函数只执行一次
        - 后续调用立即返回
        - 线程安全
    end note
```

#### 多Goroutine调用Once时序图

```mermaid
sequenceDiagram
    participant G1 as Goroutine 1
    participant G2 as Goroutine 2
    participant Once as sync.Once
    participant G3 as Goroutine 3
    participant Init as 初始化函数
    
    Note over G1,Init: 场景：多个Goroutine并发调用once.Do()
    
    G1->>Once: Do(init)
    Once->>Init: 执行 init() ✓
    Note over Init: 初始化操作
    
    G2->>Once: Do(init)
    Note over G2: ⏸️ 阻塞等待
    
    G3->>Once: Do(init)
    Note over G3: ⏸️ 阻塞等待
    
    Init-->>Once: 完成
    Once-->>G1: 返回
    
    Once-->>G2: 返回 (不执行init)
    Once-->>G3: 返回 (不执行init)
    
    Note over Once: done = 1, 后续调用直接返回
    
    G1->>Once: Do(init)
    Once-->>G1: 立即返回 (不执行)
```

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    instance *Singleton
    once     sync.Once
)

type Singleton struct {
    data string
}

// GetInstance 使用sync.Once实现线程安全的单例
func GetInstance() *Singleton {
    once.Do(func() {
        fmt.Println("创建Singleton实例（只执行一次）")
        time.Sleep(100 * time.Millisecond) // 模拟初始化耗时
        instance = &Singleton{data: "singleton instance"}
    })
    return instance
}

func main() {
    var wg sync.WaitGroup
    
    // 10个Goroutine并发调用GetInstance
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            inst := GetInstance()
            fmt.Printf("Goroutine %d: %p - %s\n", id, inst, inst.data)
        }(i)
    }
    
    wg.Wait()
    // 输出：所有Goroutine获得同一个实例（地址相同）
    // "创建Singleton实例"只打印一次
}
```

---

## 3. 工程分析与最佳实践

- 推荐优先使用channel实现同步，sync适合低层并发控制。
- Mutex/RWMutex适合保护共享资源，避免数据竞争。
- WaitGroup适合任务编排，避免忙等。
- Once适合单例、懒加载等场景。
- Cond适合复杂同步，需谨慎使用。
- 尽量缩小锁的粒度，减少锁竞争。

---

## 4. 常见陷阱

- 忘记Unlock会导致死锁。
- 多次Unlock会panic。
- WaitGroup的Add/Done不匹配会导致永久阻塞。
- RWMutex写锁不可重入。

---

## 5. 单元测试建议

- 测试并发场景下的数据一致性与死锁边界。
- 使用-race检测数据竞争。

---

## 6. 参考文献

- Go官方文档：<https://golang.org/pkg/sync/>
- Go Blog: <https://blog.golang.org/share-memory-by-communicating>
- 《Go语言高级编程》

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.25.3+
