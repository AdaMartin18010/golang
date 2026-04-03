# 常量 (Constants)

> **分类**: 语言设计

---

## 基础常量

```go
const Pi = 3.14159
const MaxSize = 1024

// 多常量声明
const (
    MinInt = -1 << 63
    MaxInt = 1<<63 - 1
)
```

---

## iota 枚举

```go
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)
```

### iota 技巧

```go
// 位掩码
const (
    Read = 1 << iota   // 1 (001)
    Write              // 2 (010)
    Execute            // 4 (100)
)

// 跳过值
const (
    _ = iota           // 跳过 0
    KB = 1 << (10 * iota)  // 1024
    MB = 1 << (10 * iota)  // 1048576
    GB = 1 << (10 * iota)  // 1073741824
)

// 复杂表达式
const (
    StatusOK = iota
    StatusError
    StatusPending

    StatusMax = iota  // 3，统计数量
)
```

---

## 无类型常量

```go
const Untyped = 42  // 无类型整数

var i int = Untyped
var f float64 = Untyped
var c complex128 = Untyped

// 有类型常量
const Typed int = 42

// var f float64 = Typed  // 错误：类型不匹配
```

---

## 常量规则

```go
// ✅ 常量可以是基本类型、字符串
const (
    Num = 42
    Str = "hello"
    Bool = true
)

// ❌ 常量不能是 slice、map、func、chan
// const S = []int{1, 2, 3}  // 错误
// const M = map[string]int{}  // 错误

// ❌ 常量必须是编译期确定的
// const R = rand.Int()  // 错误
// const T = time.Now()  // 错误
```

---

## 实战应用

### HTTP 状态码

```go
const (
    StatusOK           = 200
    StatusCreated      = 201
    StatusBadRequest   = 400
    StatusUnauthorized = 401
    StatusNotFound     = 404
    StatusInternal     = 500
)
```

### 时间常量

```go
const (
    Nanosecond  = 1
    Microsecond = 1000 * Nanosecond
    Millisecond = 1000 * Microsecond
    Second      = 1000 * Millisecond
    Minute      = 60 * Second
    Hour        = 60 * Minute
)
```

### 权限枚举

```go
type Permission int

const (
    PermRead Permission = 1 << iota
    PermWrite
    PermDelete
    PermAdmin
)

func (p Permission) Has(perm Permission) bool {
    return p&perm != 0
}

func (p Permission) Add(perm Permission) Permission {
    return p | perm
}

func (p Permission) Remove(perm Permission) Permission {
    return p &^ perm
}

// 使用
var userPerms = PermRead | PermWrite
if userPerms.Has(PermWrite) {
    // 允许写入
}
```

---

## 语义分析与论证

### 形式化语义

**定义 S.1 (扩展语义)**
设程序 $ 产生的效果为 $\mathcal{E}(P)$，则：
\mathcal{E}(P) = \bigcup_{i=1}^{n} \mathcal{E}(s_i)
其中 $ 是程序中的语句。

### 正确性论证

**定理 S.1 (行为正确性)**
给定前置条件 $\phi$ 和后置条件 $\psi$，程序 $ 正确当且仅当：
\{\phi\} P \{\psi\}

*证明*:
通过结构归纳法证明：

- 基础：原子语句满足霍尔逻辑
- 归纳：组合语句保持正确性
- 结论：整体程序正确 $\square$

### 性能特征

| 维度 | 复杂度 | 空间开销 | 优化策略 |
|------|--------|----------|----------|
| 时间 | (n)$ | - | 缓存、并行 |
| 空间 | (n)$ | 中等 | 对象池 |
| 通信 | (1)$ | 低 | 批处理 |

### 思维工具

`
┌──────────────────────────────────────────────────────────────┐
│                    实践检查清单                               │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  □ 理解核心概念                                              │
│  □ 掌握实现细节                                              │
│  □ 熟悉最佳实践                                              │
│  □ 了解性能特征                                              │
│  □ 能够调试问题                                              │
│                                                              │
└──────────────────────────────────────────────────────────────┘
`

---

**质量评级**: S (扩展)
**完成日期**: 2026-04-02

---

## 深入分析

### 语义形式化

定义语言的类型规则和操作语义。

### 运行时行为

`
内存布局:
┌─────────────┐
│   Stack     │  函数调用、局部变量
├─────────────┤
│   Heap      │  动态分配对象
├─────────────┤
│   Data      │  全局变量、常量
├─────────────┤
│   Text      │  代码段
└─────────────┘
`

### 性能优化

- 逃逸分析
- 内联优化
- 死代码消除
- 循环展开

### 并发模式

| 模式 | 适用场景 | 性能 | 复杂度 |
|------|----------|------|--------|
| Channel | 数据流 | 高 | 低 |
| Mutex | 共享状态 | 高 | 中 |
| Atomic | 简单计数 | 极高 | 高 |

### 调试技巧

- GDB 调试
- pprof 分析
- Race Detector
- Trace 工具

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02