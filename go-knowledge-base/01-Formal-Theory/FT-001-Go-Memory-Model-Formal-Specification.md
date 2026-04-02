# FT-001: Go 内存模型形式化规范 (Go Memory Model Formal Specification)

> **维度**: Formal Theory
> **级别**: S (25+ KB)
> **标签**: #memory-model #happens-before #formal-verification
> **权威来源**: [The Go Memory Model](https://go.dev/ref/mem), [Hardware Memory Models](https://research.swtch.com/hwmm)

---

## 官方规范解读

Go 内存模型定义了 goroutine 对共享变量的读写可见性保证。这是理解 Go 并发编程的基础。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go Memory Model Hierarchy                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Sequential Consistency                                                     │
│       │                                                                     │
│       ├──► Within a single goroutine: Program order = Execution order       │
│       │                                                                     │
│       └──► Between goroutines: Happens-before relation                      │
│                                                                             │
│  Happens-Before Rules:                                                      │
│  ─────────────────────                                                      │
│  1. Initialization: var declaration → main.main()                           │
│  2. Goroutine creation: go statement → goroutine start                      │
│  3. Channel communication: send → receive (same channel)                    │
│  4. Lock: Unlock → Lock (same mutex)                                        │
│  5. Once: f() return → subsequent f() calls                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Happens-Before 形式化定义

$$
\begin{aligned}
&\text{Happens-Before Relation } (\xrightarrow{hb}): \\
&1. \text{Program Order: } \forall e_1, e_2 \in G: \text{if } i < j \text{ then } e_i \xrightarrow{hb} e_j \\
&2. \text{Channel Send-Receive: } ch <- v \xrightarrow{hb} <-ch \\
&3. \text{Mutex: } mu.Unlock() \xrightarrow{hb} mu.Lock() \\
&4. \text{Once: } once.Do(f)_{return} \xrightarrow{hb} once.Do(f)_{subsequent} \\
\\
&\text{Synchronizes-With } (\xrightarrow{sw}): \\
&\text{If } A \xrightarrow{sw} B \text{ then } A \xrightarrow{hb} B \\
\\
&\text{Visibility: } \\
&\text{If } A \xrightarrow{hb} B \text{ and } A \text{ writes } x, B \text{ reads } x \Rightarrow B \text{ sees } A \text{'s write}
\end{aligned}
$$

---

## 核心代码分析

### Go Runtime 内存屏障实现

```go
// src/runtime/internal/atomic/atomic_amd64.go

//go:nosplit
//go:noinline
func Load(ptr *uint32) uint32 {
 return *ptr
}

// LoadAcquire 具有获取语义的加载
//go:nosplit
//go:noinline
func LoadAcquire(ptr *uint32) uint32 {
 // On amd64, normal loads are already acquire
 return *ptr
}

// StoreRelease 具有释放语义的存储
//go:nosplit
//go:noinline
func StoreRelease(ptr *uint32, val uint32) {
 // On amd64, normal stores are already release
 *ptr = val
}
```

### Channel 的 happens-before 保证

```go
// src/runtime/chan.go

// 无缓冲 channel
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
    // ...

    // 1. 发送方写入数据
    typedmemmove(c.elemtype, qp, ep)

    // 2. 唤醒接收方
    goready(gp, skip+1)

    // 关键：发送方的写入发生在接收方唤醒之前
    // 这建立了 happens-before 关系
}

// 接收方
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
    // ...

    // 1. 接收方被唤醒（由发送方唤醒）

    // 2. 读取数据
    // 根据 happens-before，这里一定能看到发送方写入的数据
    if ep != nil {
        typedmemmove(c.elemtype, ep, qp)
    }
}
```

---

## 常见模式验证

### 模式 1: Channel 同步

```go
var done = make(chan bool)
var msg string

func a() {
    msg = "hello, world"
    done <- true
}

func b() {
    <-done
    println(msg)  // 一定打印 "hello, world"
}
```

**形式化证明:**

- 发送 `done <- true` happens-before 接收 `<-done`
- `msg = "hello..."` 在发送之前（程序顺序）
- 因此 `msg` 的写入 happens-before 读取

### 模式 2: Mutex 同步

```go
var mu sync.Mutex
var shared int

func writer() {
    mu.Lock()
    shared = 42
    mu.Unlock()  // Unlock happens-before Lock
}

func reader() {
    mu.Lock()
    println(shared)  // 一定看到 42（如果 writer 先执行）
    mu.Unlock()
}
```

---

## 反模式与陷阱

### 数据竞争（Data Race）

```go
// 错误：没有 happens-before 关系
var flag bool
var result int

func setup() {
    result = 42
    flag = true  // 普通写入，没有同步
}

func main() {
    go setup()
    for !flag {}  // 普通读取，没有同步
    println(result)  // 可能打印 0！
}
```

**修正:**

```go
var flag atomic.Bool  // 使用 atomic
var result int

func setup() {
    result = 42
    flag.Store(true)
}

func main() {
    go setup()
    for !flag.Load() {}
    println(result)  // 一定打印 42
}
```

---

## 硬件内存模型对比

| 架构 | Store-Load 重排序 | Go 保证 |
|------|------------------|---------|
| x86/amd64 | 不允许 | Strong memory model |
| ARM | 允许 | Go runtime 插入屏障 |
| RISC-V | 允许 | Go runtime 插入屏障 |

```go
// src/runtime/internal/atomic/sys_linux_arm64.s

// 在 ARM64 上需要显式屏障
text ·StoreRelease(SB), NOSPLIT, $0-16
    MOVD ptr+0(FP), R0
    MOVW val+8(FP), R1
    STLRW R1, (R0)  // Store-Release
    RET
```

---

## 参考文献

1. [The Go Memory Model](https://go.dev/ref/mem) - 官方规范
2. [Memory Models](https://research.swtch.com/mm) - Russ Cox 系列文章
3. [Happens-Before](https://en.wikipedia.org/wiki/Happened-before) - Lamport, 1978
4. [Java Memory Model](https://www.cs.umd.edu/~pugh/java/memoryModel/) - 对比参考
