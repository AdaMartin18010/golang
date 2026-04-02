# CSP 理论 (Communicating Sequential Processes)

> **分类**: 形式理论
> **难度**: 专家
> **前置知识**: 进程代数基础

---

## 概述

CSP 是由 Tony Hoare (1978) 提出的并发计算形式化语言，是 Go 并发模型的理论基础。

**核心理念**: Do not communicate by sharing memory; share memory by communicating.

---

## 基本语法

### 进程 (Process)

```
P ::=
  | STOP               死锁/终止
  | SKIP               成功终止
  | a → P              事件 a 后执行 P (前缀)
  | P ⊓ Q              内部选择
  | P □ Q              外部选择
  | P ||| Q            交错并行
  | P [| A |] Q        同步并行 (在集合 A 上同步)
  | P \ A              隐藏 (限制事件集合 A)
```

### 事件 (Event)

```
a ::=
  | c!v                在通道 c 上发送值 v
  | c?x                在通道 c 上接收值到 x
  | τ                  内部事件 (不可见)
```

---

## 语义

### 迹语义 (Trace Semantics)

进程的**迹** (trace) 是可能发生的事件序列。

```
traces(P) = { 所有可能的有限事件序列 }
```

### 操作语义

#### 前缀

```
─────────────── (Prefix)
a → P --a--> P
```

执行事件 a 后成为进程 P。

#### 外部选择

```
P --a--> P'
──────────────── (Ext-Left)
P □ Q --a--> P'

Q --a--> Q'
──────────────── (Ext-Right)
P □ Q --a--> Q'
```

环境选择执行 P 或 Q 中的可用事件。

#### 内部选择

```
─────── (Int-Left)
P ⊓ Q --> P

─────── (Int-Right)
P ⊓ Q --> Q
```

进程内部决定执行 P 或 Q（非确定性）。

#### 并行组合

```
P --a--> P'    a ∉ A
────────────────────── (Par-Left)
P [| A |] Q --a--> P' [| A |] Q

P --a--> P'    Q --a--> Q'    a ∈ A
──────────────────────────────────── (Par-Sync)
P [| A |] Q --a--> P' [| A |] Q'
```

在集合 A 中的事件必须同步执行。

---

## Go 与 CSP 的对应

| CSP | Go | 说明 |
|-----|-----|------|
| `P ||| Q` | `go P(); Q()` | 并行执行 |
| `c!v → P` | `ch <- v; ...` | Channel 发送 |
| `c?x → P` | `x := <-ch; ...` | Channel 接收 |
| `P □ Q` | `select { case ... }` | 选择 |
| `SKIP` | `return` | 成功终止 |

### Channel 的 CSP 表示

```
// Go: ch := make(chan int, 2)
Channel(c, 2) =
  c?x → Channel(c, 2) ⊓
  (if n < 2 then c!v → Channel(c, n+1) else STOP)
```

---

## 精化 (Refinement)

### 迹精化

```
P ⊑T Q  ⟺  traces(Q) ⊆ traces(P)

Q 是 P 的迹精化，如果 Q 的所有迹都是 P 的迹
```

### 失败精化

```
P ⊑F Q  ⟺  failures(Q) ⊆ failures(P)

失败 = (迹, 拒绝集)
```

### Go 正确性

```
Spec ⊑ Impl

实现 Impl 满足规范 Spec
```

---

## 死锁自由

### 死锁定义

```
进程 P 死锁，如果:
- P 不是 SKIP
- P 不能执行任何事件
```

### 死锁自由证明

```
进程 P 是死锁自由的，如果:
∀s ∈ traces(P), P/s ≠ STOP

其中 P/s 是 P 执行迹 s 后的剩余进程
```

### Go 应用

```go
// 死锁示例
ch := make(chan int)
<-ch  // 死锁: 没有发送者

// CSP: c?x → STOP 没有匹配的发送
```

---

## 与 π-演算关系

| 特性 | CSP | π-演算 |
|------|-----|--------|
| 通道 | 全局命名 | 一等值 (可传递) |
| 移动性 | 无 | 有 (通道可发送) |
| 适用 | 系统描述 | 动态拓扑 |
| Go 对应 | 基础并发 | 通道传递 |

Go 结合了 CSP 的简洁和 π-演算的移动性 (通道可作为值传递)。

---

## 形式化性质

### 交换律

```
P ||| Q = Q ||| P
P □ Q = Q □ P
```

### 结合律

```
(P ||| Q) ||| R = P ||| (Q ||| R)
```

### 分配律

```
P □ (Q ⊓ R) = (P □ Q) ⊓ (P □ R)
```

---

## 参考

- "Communicating Sequential Processes" (Hoare, 1978)
- "The Theory and Practice of Concurrency" (Roscoe)
- University of Padua: Languages for Concurrency and Distribution
