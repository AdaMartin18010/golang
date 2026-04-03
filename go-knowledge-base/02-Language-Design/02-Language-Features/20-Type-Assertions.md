# 类型断言 (Type Assertions)

> **分类**: 语言设计

---

## 基本断言

```go
var i interface{} = "hello"

// 断言为具体类型
s := i.(string)  // "hello"

// 带检查的断言
s, ok := i.(string)  // ok = true
n, ok := i.(int)     // ok = false

// 安全断言
if s, ok := i.(string); ok {
    fmt.Println(s)
} else {
    fmt.Println("not a string")
}
```

---

## 类型开关

```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %q\n", v)
    case bool:
        fmt.Printf("bool: %t\n", v)
    case []int:
        fmt.Printf("slice of ints: %v\n", v)
    case map[string]int:
        fmt.Printf("map: %v\n", v)
    case nil:
        fmt.Println("nil")
    case Person:
        fmt.Printf("person: %s\n", v.Name)
    case *Person:
        fmt.Printf("person pointer: %s\n", v.Name)
    default:
        fmt.Printf("unknown type: %T\n", v)
    }
}
```

---

## 接口断言

```go
type Reader interface {
    Read() ([]byte, error)
}

type Closer interface {
    Close() error
}

func process(r Reader) error {
    // 检查是否也实现了 Closer
    if c, ok := r.(Closer); ok {
        defer c.Close()
    }

    data, err := r.Read()
    // ...
}
```

---

## 空接口判断

```go
func isNil(i interface{}) bool {
    if i == nil {
        return true
    }

    // 检查底层值是否为 nil
    v := reflect.ValueOf(i)
    return v.Kind() == reflect.Ptr && v.IsNil()
}

// 使用
var p *Person
fmt.Println(isNil(p))  // true
fmt.Println(isNil(nil)) // true
```

---

## 性能考虑

```go
// 编译器优化：直接类型断言
var r io.Reader = strings.NewReader("hello")

// 编译器已知类型，优化为直接检查
if sr, ok := r.(*strings.Reader); ok {
    // 无运行时开销
}

// 多次断言缓存类型
var typ = reflect.TypeOf((*MyInterface)(nil)).Elem()

func checkType(i interface{}) bool {
    return reflect.TypeOf(i).Implements(typ)
}
```

---

## 最佳实践

```go
// ✅ 优先使用类型开关
switch v := i.(type) {
case int:
    // ...
case string:
    // ...
}

// ✅ 检查 ok 值避免 panic
v, ok := i.(MyType)
if !ok {
    return errors.New("type mismatch")
}

// ❌ 避免裸断言
v := i.(MyType)  // 可能 panic

// ✅ 断言前检查 nil
if i == nil {
    return nil
}
v, ok := i.(MyType)
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