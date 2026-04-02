# DRF-SC 保证 (Data-Race-Free → Sequential Consistency)

> **分类**: 形式理论
> **难度**: 专家
> **前置知识**: Happens-Before、顺序一致性

---

## 概述

**DRF-SC** 是 Go 内存模型的核心保证：

> **无数据竞争 (DRF) 的程序表现出顺序一致性 (SC) 的行为。**

这意味着程序员只需要避免数据竞争，就可以像编写单线程程序一样思考并发程序。

---

## 形式化定义

### 顺序一致性 (Sequential Consistency)

**定义** (Lamport, 1979):

```
程序的执行是顺序一致的，如果:
1. 所有处理器的操作按某种全局顺序执行
2. 每个处理器的操作按其程序顺序出现
```

### 数据竞争 (Data Race)

**定义**:

```
两个内存访问构成数据竞争，如果:
1. 访问同一内存位置
2. 至少一个是写操作
3. 它们之间没有 happens-before 关系
```

形式化:

```
Race(e₁, e₂) ⟺
  same_location(e₁, e₂) ∧
  (is_write(e₁) ∨ is_write(e₂)) ∧
  ¬(e₁ →hb e₂) ∧ ¬(e₂ →hb e₁)
```

### DRF-SC 定理

```
Theorem (DRF-SC):
如果一个程序没有数据竞争，那么它的所有执行都是顺序一致的。

形式化:
DRF(program) ⟹ SC(program)

其中:
- DRF(program) = ¬∃e₁, e₂ ∈ program: Race(e₁, e₂)
- SC(program) = 所有执行等价于某顺序执行
```

---

## 证明草图

### 目标

证明: 无数据竞争的程序，所有读操作看到唯一确定的值。

### 步骤 1: 建立偏序

```
对于无数据竞争程序，冲突访问 (同一位置) 有全序:

∀e₁, e₂ 访问同一位置:
  e₁ →hb e₂  ∨  e₂ →hb e₁

否则就构成数据竞争。
```

### 步骤 2: 扩展为全序

```
使用拓扑排序，将 happens-before 偏序扩展为全序:

全序 ≻ 满足:
1. e₁ →hb e₂ ⟹ e₁ ≻ e₂
2. 所有事件可比
```

### 步骤 3: 证明等价性

```
对于任意读操作 r，设 w 是 ≻ 中 r 之前的最后一次写:

1. w 存在 (程序必须初始化)
2. w 是唯一的 (全序性质)
3. r 看到 w 的值 (happens-before)
4. 没有 w' 在 w 和 r 之间 (w 的定义)

因此，读看到确定的值。
```

### 步骤 4: 顺序一致性

```
全序 ≻ 就是顺序一致性的见证:
- 按 ≻ 顺序执行所有操作
- 每个 goroutine 的操作按其程序顺序 (由 →hb 保证)
- 结果与实际执行相同
```

---

## DRF 的边界

### 什么是 DRF 不保证的

DRF-SC **不保证**: 无数据竞争的程序是正确的。

```go
// 无数据竞争，但仍有 bug
var x int
var done bool

func main() {
    go func() {
        x = 42
        done = true
    }()
    for !done {}  // 忙等待
    fmt.Println(x)
}

// 问题: done 的读取没有同步
// 但在某些架构上可能工作
```

### 需要额外的同步

```go
// 正确版本
var x int
var done = make(chan bool)

func main() {
    go func() {
        x = 42
        done <- true
    }()
    <-done  // 同步点
    fmt.Println(x)  // 一定看到 42
}
```

---

## 与其他语言对比

| 语言 | 内存模型 | DRF-SC |
|------|----------|--------|
| Go | Happens-before | ✅ 有 |
| Java | Happens-before | ✅ 有 |
| C++ | Happens-before | ✅ 有 (DRF-SC 或 RWC) |
| JavaScript | Event loop | ✅ 单线程 |
| Rust | Happens-before | ✅ 有 |

---

## 实际意义

### 对程序员的含义

```
只需要:
1. 使用正确的同步 (channel、mutex、atomic)
2. 避免数据竞争 (race detector 帮助)

不需要担心:
- 内存重排序
- 缓存一致性
- 处理器架构差异
```

### 对编译器的含义

```
可以进行的优化:
- 重排序无竞争的操作
- 缓存值到寄存器

不能进行的优化:
- 跨越同步点的重排序
- 引入新的数据竞争
```

---

## 弱于 DRF 的程序

### Racy 程序

有数据竞争的程序，行为**未定义**:

```go
var x int

go func() { x = 1 }()
go func() { x = 2 }()

// x 的最终值不确定
// 可能 1，可能 2，可能损坏
```

### 良性竞争 (Benign Races)

某些数据竞争可能无害，但**极不推荐**:

```go
// 不要这样做
var counter int64

func Inc() {
    counter++  // 非原子操作，数据竞争
}
```

---

## 形式化工具

### 验证 DRF

```
工具: Go Race Detector
命令: go run -race program.go
原理: 基于 happens-before 的动态分析
```

### 模型检测

```
工具: TLA+, SPIN
用途: 验证小规模并发算法的正确性
```

---

## 练习

### 练习 1

分析以下程序是否有数据竞争，是否满足 DRF-SC:

```go
var x int
var mu sync.Mutex

func main() {
    go func() {
        mu.Lock()
        x = 1
        mu.Unlock()
    }()

    mu.Lock()
    fmt.Println(x)
    mu.Unlock()
}
```

### 练习 2

证明: 使用 channel 同步的程序是 DRF 的。

### 练习 3

给出一个程序: 无数据竞争，但有并发 bug。

---

## 参考

- "The Go Memory Model"
- "Memory Models: A Case for Rethinking Parallel Languages and Hardware" (Sarita Adve et al.)
- "Java Memory Model Pragmatics" (Doug Lea)
- POPL 2022: "DRF-SC for GPU"
