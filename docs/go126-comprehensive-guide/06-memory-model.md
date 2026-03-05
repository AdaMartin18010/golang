# Go内存模型

> Happens-before关系与同步语义的形式化分析

---

## 一、内存模型基础

### 1.1 Happens-before关系

```
形式化定义:
────────────────────────────────────────
事件a happens-before事件b (记作 a ≺ b):
如果a在程序顺序中先于b，且内存效应可见

Happens-before性质:
├─ 传递性: a ≺ b ∧ b ≺ c ⟹ a ≺ c
├─ 非自反: ¬(a ≺ a)
├─ 偏序: 并非所有事件对都可比较
└─ 同步: a ≺ b 保证a的写对b可见

内存操作视角:
写操作W(x) = v → 将值v写入x
读操作R(x) → 从x读取值

可见性保证:
W(x) = v ≺ R(x) ⟹ R(x)可能读到v
```

### 1.2 程序顺序与执行

```
程序顺序:
────────────────────────────────────────
单个goroutine内的语句按源代码顺序执行

但编译器和CPU可能重排:
├─ 不违反单goroutine语义的指令重排
├─ 写缓存延迟
└─ 读预取优化

同步原语建立跨goroutine顺序:
channel、锁、atomic操作建立happens-before边
```

---

## 二、同步原语语义

### 2.1 Channel同步

```
Channel Happens-before规则:
────────────────────────────────────────

无缓冲Channel:
send(ch, v) ≺ recv(ch, v)
发送在接收完成前发生

有缓冲Channel:
第k个send ≺ 第k个recv (当缓冲满时)
close(ch) ≺ recv返回零值+false

示例:
var c = make(chan int, 10)
var a string

func f() {
    a = "hello, world"
    c <- 0  // send
}

func main() {
    go f()
    <-c      // recv
    print(a)  // 保证看到"hello, world"
}

证明:
c <- 0 在 <-c 之前发生
a = "..." 在 c <- 0 之前 (程序顺序)
∴ a = "..." ≺ print(a)
```

### 2.2 Mutex同步

```
Mutex Happens-before规则:
────────────────────────────────────────
Unlock(m) ≺ Lock(m) 在相同mutex上

示例:
var mu sync.Mutex
var shared int

func writer() {
    mu.Lock()
    shared = 42
    mu.Unlock()  // 释放前写可见
}

func reader() {
    mu.Lock()    // 获取后看到writer的写
    print(shared)
    mu.Unlock()
}

RWMutex扩展:
RLock ≺ RUnlock
Lock ≺ Unlock
RUnlock ≺ Lock (当有等待的写者)
```

### 2.3 WaitGroup同步

```
WaitGroup语义:
────────────────────────────────────────
Add(n) ≺ Wait() 当计数器>0
Done() ≺ Wait() 当计数器归零

使用模式:
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func() {
        defer wg.Done()
        process(item)
    }()
}
wg.Wait()  // 等待所有goroutine完成

正确性:
每个Done()在Wait返回前发生
process()在Done()前发生 (程序顺序)
∴ process() ≺ Wait之后代码
```

---

## 三、Atomic操作

### 3.1 Atomic语义

```
Atomic保证:
────────────────────────────────────────
├─ 原子性: 操作不可中断
├─ 顺序一致性: 所有goroutine看到相同操作顺序
└─ 无数据竞争: atomic变量访问不构成竞争

Happens-before:
atomic.StoreUint64(&x, v) ≺ atomic.LoadUint64(&x)

注意:
atomic操作不提供对非atomic变量的顺序保证
需要配合atomic变量作为同步信号

示例:
var done uint32
var shared string

func writer() {
    shared = "ready"
    atomic.StoreUint32(&done, 1)
}

func reader() {
    for atomic.LoadUint32(&done) == 0 {
        // spin
    }
    print(shared)  // 保证看到"ready"
}
```

### 3.2 Atomic vs Mutex

```
选择指南:
────────────────────────────────────────
Atomic适用:
├─ 计数器、标志位
├─ 简单的状态机
└─ 高频访问的低争用场景

Mutex适用:
├─ 保护临界区
├─ 复杂不变式
├─ 条件等待
└─ 多变量一致性

性能对比:
├─ Atomic: 通常1-2ns
├─ Mutex: 通常10-50ns (无竞争)
└─ Mutex: 竞争时内核上下文切换
```

---

## 四、Once与Pool

### 4.1 Once语义

```
sync.Once保证:
────────────────────────────────────────
once.Do(f):
├─ f执行且仅执行一次
├─ 所有Do返回时f已完成
└─ 多个Do并发调用只有一个执行f，其余等待

Happens-before:
f()中的操作 ≺ 任何Do返回

应用: 单例初始化
var instance *Resource
var once sync.Once

func GetInstance() *Resource {
    once.Do(func() {
        instance = createResource()
    })
    return instance
}

初始化 happens-before 返回
```

### 4.2 Pool语义

```
sync.Pool语义:
────────────────────────────────────────
├─ Get(): 从池中获取对象
├─ Put(x): 将对象放回池中
└─ GC时池可能被清空

无显式同步保证:
Put(x) ≺ Get() 返回x 不保证

使用场景:
├─ 临时对象复用，减少GC压力
├─ 对象创建成本高
└─ 对象生命周期短暂

type Buffer struct {
    pool sync.Pool
}

func (b *Buffer) Get() []byte {
    v := b.pool.Get()
    if v == nil {
        return make([]byte, 0, 1024)
    }
    return v.([]byte)[:0]
}

func (b *Buffer) Put(buf []byte) {
    b.pool.Put(buf)
}
```

---

## 五、数据竞争

### 5.1 数据竞争定义

```
数据竞争:
────────────────────────────────────────
两个goroutine并发访问同一内存位置，且:
├─ 至少一个是写操作
└─ 无happens-before关系

形式化:
W(x) ∥ R(x): 竞争
W(x) ∥ W(x): 竞争

后果 (未定义行为):
├─ 读取脏值
├─ 状态不一致
├─ 编译器优化失效
└─ 程序崩溃

检测:
go test -race
编译时插入检测代码
```

### 5.2 避免竞争

```
同步策略:
────────────────────────────────────────

1. 不共享:
   每个goroutine有独立数据副本

2. 通过通信共享:
   使用channel传递数据所有权

3. 通过同步保护:
   使用mutex、atomic保护共享状态

Go并发哲学:
"不要通过共享内存来通信，而要通过通信来共享内存"

示例 (通信替代共享):
// 不良: 共享计数器
var counter int
var mu sync.Mutex

// 良好: channel传递
func counterWorker(in <-chan int, out chan<- int) {
    count := 0
    for range in {
        count++
    }
    out <- count
}
```

---

## 六、Go 1.26增强

### 6.1 内存模型更新

```
Go 1.26改进:
────────────────────────────────────────
├─ 更精确的happens-before追踪
├─ race detector精度提升
├─ 减少误报
└─ 增强weak memory ordering支持

新工具支持:
├─ runtime.SetRaceCallback
└─ 细粒度竞争分析
```

### 6.2 泄漏检测集成

```
Goroutine泄漏与内存:
────────────────────────────────────────
泄漏goroutine持有内存不释放
Go 1.26: runtime.SetGoroutineLeakCallback

检测模式:
├─ 无退出的goroutine
├─ 阻塞在channel/锁
└─ 资源累积
```

---

*本章形式化分析了Go内存模型，为编写正确的并发程序提供理论基础。*
