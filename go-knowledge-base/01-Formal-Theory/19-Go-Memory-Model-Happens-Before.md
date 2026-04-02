# Go 内存模型与 Happens-Before 关系 (Go Memory Model & Happens-Before)

> **分类**: 形式理论
> **标签**: #memory-model #happens-before #synchronization #drf-sc
> **参考**: The Go Memory Model (go.dev/ref/mem), Russ Cox 系列论文

---

## 内存模型基础

### 顺序一致性 (Sequential Consistency)

Go 内存模型基于 **DRF-SC** (Data-Race-Free Sequential Consistency) 原则：

> **定理**: 如果程序没有数据竞争，那么它的执行结果与某种顺序一致的执行相同。

```
┌─────────────────────────────────────────────────────────────────┐
│                    Sequential Consistency                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Program Order:        ──────────────────────────────►          │
│  (每个 goroutine 内)                                              │
│                                                                  │
│  Goroutine 1:    A ──► B ──► C                                  │
│                     \      /                                    │
│                      \    /  同步边                              │
│                       \  /                                       │
│  Goroutine 2:        D ──► E ──► F                              │
│                                                                  │
│  全局顺序: A → B → C → D → E → F (满足所有 happens-before)         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Happens-Before 定义

**定义**: 如果事件 A happens-before 事件 B（记作 A ≺ B），则：

1. A 的内存写入对 B 可见
2. A 和 B 不能重排序（A 必须在 B 之前执行）

```go
// 示例：Happens-Before 关系
var a, b int

func f() {
    a = 1           // 事件 A
    b = 2           // 事件 B
    // A ≺ B (程序顺序)
}

func g() {
    if b == 2 {     // 事件 C (读取 b)
        fmt.Println(a) // 事件 D (读取 a)
        // 注意：没有 happens-before 关系保证 C ≺ D 时 a=1 可见！
    }
}
```

---

## Happens-Before 规则详解

### 1. 程序顺序 (Program Order)

```go
// 在同一个 goroutine 内，程序顺序创建 happens-before 关系
func programOrder() {
    x := 1      // A
    y := 2      // B
    // A ≺ B (程序顺序)
}
```

### 2. Goroutine 启动

```go
// go 语句 happens-before goroutine 开始执行
func goroutineStart() {
    var a string
    a = "hello"           // A
    go func() {
        fmt.Println(a)    // B
        // A ≺ B (go 语句 happens-before goroutine 执行)
    }()
}
```

### 3. Channel 操作

```go
// 规则 1: 发送 happens-before 对应的接收完成
func channelSendRecv() {
    ch := make(chan int)
    var a string

    go func() {
        a = "hello"       // A
        ch <- 0           // B (发送)
    }()

    <-ch                  // C (接收完成)
    fmt.Println(a)        // D
    // A ≺ B ≺ C ≺ D，因此 "hello" 一定可见
}

// 规则 2: 关闭 happens-before 接收到零值
func channelClose() {
    ch := make(chan int)
    var a string

    go func() {
        a = "hello"       // A
        close(ch)         // B (关闭)
    }()

    <-ch                  // C (接收，返回零值)
    fmt.Println(a)        // D
    // A ≺ B ≺ C ≺ D
}

// 规则 3: 无缓冲通道：接收 happens-before 发送完成
func unbufferedChannel() {
    ch := make(chan int)  // 无缓冲通道
    var a string

    go func() {
        a = "hello"       // A
        <-ch              // B (接收)
        fmt.Println(a)    // 此时 a 一定可见
    }()

    ch <- 0               // C (发送完成)
    // B ≺ C
}

// 规则 4: 缓冲通道：第 k 次接收 happens-before 第 k+C 次发送
func bufferedChannel() {
    ch := make(chan int, 10)  // 容量为 10
    // 第 1 次接收 happens-before 第 11 次发送
    // 第 2 次接收 happens-before 第 12 次发送
}
```

### 4. Mutex/RWMutex

```go
// 第 n 次 Unlock happens-before 第 n+m 次 Lock
func mutexHB() {
    var mu sync.Mutex
    var a string

    go func() {
        mu.Lock()
        a = "hello"       // A
        mu.Unlock()       // B
    }()

    mu.Lock()             // C
    fmt.Println(a)        // D
    mu.Unlock()
    // B ≺ C，因此 "hello" 可见
}
```

### 5. Once

```go
// once.Do(f) 中的 f() 完成 happens-before 任何 once.Do 返回
func onceHB() {
    var once sync.Once
    var a string

    go func() {
        once.Do(func() {
            a = "hello"   // A
        })                // B (Do 完成)
    }()

    once.Do(func() {
        // 这个函数不会执行
    })                    // C (Do 返回)
    // B ≺ C，因此 "hello" 可见
}
```

### 6. 原子操作 (Go 1.19+)

```go
// Go 1.19 引入新的 atomic 类型，提供顺序一致原子操作
var counter atomic.Int64

func atomicHB() {
    counter.Add(1)        // A
    // 所有原子操作表现为按某种全局顺序执行
    // 如果另一个 goroutine 看到 counter > 0，则 A happens-before 它
}
```

---

## 数据竞争 (Data Race)

### 定义

**数据竞争**发生在：

- 两个 goroutine 并发访问同一内存位置
- 至少一个是写操作
- 没有 happens-before 关系

```go
// 数据竞争示例
var counter int

func raceExample() {
    go func() { counter++ }()  // 写
    go func() { counter++ }()  // 写
    // 两个写操作并发，没有同步，数据竞争！
}
```

### 竞争检测

```bash
# 启用竞争检测器
go run -race main.go
go test -race ./...
```

```go
// 竞争检测器可以检测的竞争模式
func raceDetectorExample() {
    var x int

    // 写-写竞争
    go func() { x = 1 }()
    go func() { x = 2 }()

    // 读-写竞争
    go func() { fmt.Println(x) }()
    go func() { x = 3 }()
}
```

---

## 内存屏障与编译器重排序

### 编译器重排序

```go
// 编译器可能重排序没有依赖的操作
func compilerReorder() {
    a := 1      // A
    b := 2      // B
    // 编译器可能先执行 B 再执行 A，因为无依赖

    // 但同步操作阻止重排序
    var mu sync.Mutex
    mu.Lock()
    a = 3       // C
    mu.Unlock() // D
    // C 和 D 不能被重排序到 Lock 之前或之后
}
```

### CPU 内存模型

```
┌─────────────────────────────────────────────────────────────────┐
│                    CPU Memory Hierarchy                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │ Core 0   │  │ Core 1   │  │ Core 2   │  │ Core 3   │        │
│  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │        │
│  │ │ L1$  │ │  │ │ L1$  │ │  │ │ L1$  │ │  │ │ L1$  │ │        │
│  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │        │
│  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │  │ ┌──────┐ │        │
│  │ │ L2$  │ │  │ │ L2$  │ │  │ │ L2$  │ │  │ │ L2$  │ │        │
│  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │  │ └──────┘ │        │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │
│       └─────────────┴─────────────┴─────────────┘                │
│                         │                                       │
│                    ┌────┴────┐                                  │
│                    │  L3$    │                                  │
│                    └────┬────┘                                  │
│                         │                                       │
│                    ┌────┴────┐                                  │
│                    │ Memory  │                                  │
│                    └─────────┘                                  │
│                                                                  │
│  问题：Core 0 的写入何时对 Core 1 可见？                          │
│  答案：通过缓存一致性协议 + 内存屏障                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 常见模式与反模式

### 正确模式

```go
// 模式 1: 使用 Channel 传递所有权
func correctPattern1() {
    ch := make(chan *Data)

    go func() {
        data := &Data{Value: 42}  // A
        ch <- data                 // B
    }()

    data := <-ch                  // C
    fmt.Println(data.Value)       // D，安全访问
}

// 模式 2: 使用 Mutex 保护共享状态
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

// 模式 3: 使用 atomic 进行简单计数
var counter atomic.Int64

func increment() {
    counter.Add(1)
}
```

### 反模式

```go
// 反模式 1: 双重检查锁定（无 volatile）
func antiPattern1() {
    var instance *Singleton
    var mu sync.Mutex

    getInstance := func() *Singleton {
        if instance == nil {        // A (无同步读取)
            mu.Lock()
            if instance == nil {
                instance = &Singleton{}  // B
            }
            mu.Unlock()
        }
        return instance            // 可能返回未完全初始化的对象！
    }
    _ = getInstance
}

// 反模式 2: 无同步的忙等待
func antiPattern2() {
    var done bool

    go func() {
        // 做一些工作
        done = true  // A
    }()

    for !done {     // B - 可能永远看不到 done = true！
        // 忙等待
    }
}

// 修正：使用 channel 或 atomic
func fixedPattern() {
    done := make(chan struct{})

    go func() {
        // 做一些工作
        close(done)
    }()

    <-done  // 正确等待
}
```

---

## 形式化验证

### happens-before 图

```
Goroutine 1:    W(x=1) ──► ch<- ──► close(ch)
                         │
                         │ (同步)
                         ▼
Goroutine 2:            <-ch ──► R(x)

推导:
1. W(x=1) ≺ ch<-   (程序顺序)
2. ch<- ≺ <-ch     (channel 规则)
3. <-ch ≺ R(x)     (程序顺序)
4. 因此 W(x=1) ≺ R(x) (传递性)
5. 所以 R(x) 一定看到 1
```

### 循环依赖检测

```go
// 死锁检测：循环等待 happens-before 关系
func deadlockDetection() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        ch1 <- 1      // A
        <-ch2         // B
    }()

    go func() {
        ch2 <- 2      // C
        <-ch1         // D
    }()

    // A ≺ B，C ≺ D
    // 但如果 B 等待 C，D 等待 A，形成循环：
    // B 等待 C，C ≺ D，D 等待 A，A ≺ B
    // 循环依赖导致死锁
}
```
