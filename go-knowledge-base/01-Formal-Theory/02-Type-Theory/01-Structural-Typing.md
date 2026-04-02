# 结构类型系统 (Structural Typing)

> **分类**: 形式理论
> **难度**: 进阶
> **前置知识**: 类型系统基础

---

## 概述

Go 使用**结构类型系统** (Structural Typing)，与 Java/C++ 的**名义类型系统** (Nominal Typing) 相对。

**核心原则**: 类型的兼容性由结构决定，而非名称。

---

## 名义类型 vs 结构类型

### 名义类型 (Nominal)

```java
// Java: 名义类型
interface Writer {
    void Write(String data);
}

class File {
    void Write(String data) { ... }
}

// File 不实现 Writer 接口，即使方法签名相同
// 必须显式声明: class File implements Writer
```

### 结构类型 (Structural)

```go
// Go: 结构类型
type Writer interface {
    Write(data string) (n int, err error)
}

type File struct { }

func (f File) Write(data string) (n int, err error) { ... }

// File 自动实现 Writer 接口，无需显式声明
```

---

## 形式化定义

### 结构子类型

```
t₁ <: t₂    ⟺    methods(t₂) ⊆ methods(t₁)

类型 t₁ 是 t₂ 的子类型，当且仅当 t₁ 实现了 t₂ 的所有方法
```

### 方法集计算

```
methods(t) = { m | m 是类型 t 的方法 }

对于结构体 t_S:
  methods(t_S) = 显式声明的方法 + 嵌入类型的方法

对于接口 t_I:
  methods(t_I) = 接口声明的方法
```

---

## 类型规则

### 接口实现

```
t_S: struct { ... }
methods(t_S) ⊇ { m₁, m₂, ..., mₙ }
t_I: interface { m₁; m₂; ...; mₙ }
────────────────────────────────────  (T-Impl)
t_S implements t_I

即: t_S <: t_I
```

### 赋值兼容性

```
Γ ⊢ e: t₁    t₁ <: t₂
──────────────────────  (T-Assign)
Γ ⊢ e: t₂
```

### 接口赋值

```
Γ ⊢ e: t_S    t_S <: t_I
──────────────────────────  (T-Interface)
Γ ⊢ e: t_I
```

---

## 方法集计算详解

### 简单类型

```go
type Point struct {
    X, Y float64
}

func (p Point) Add(q Point) Point { ... }
func (p Point) String() string { ... }
```

方法集：

```
methods(Point) = { Add(Point) Point, String() string }
```

### 嵌入类型

```go
type ColoredPoint struct {
    Point           // 嵌入
    Color string
}
```

方法集：

```
methods(ColoredPoint) = {
    Add(Point) Point,    // 从 Point 提升
    String() string,      // 从 Point 提升
    Color() string        // 如果定义了 Color 方法
}
```

**提升规则**:

```
if t_S embeds t' and m ∈ methods(t'):
  then m ∈ methods(t_S)
  (with receiver adjusted)
```

---

## 类型相等

### 结构相等

```
t₁ = t₂    ⟺    structure(t₁) = structure(t₂)
```

### 名义相等 vs 结构相等

| 特性 | 名义相等 | 结构相等 |
|------|----------|----------|
| 判断依据 | 类型名称 | 类型结构 |
| 示例 | `type A int` ≠ `type B int` | `struct{ X int }` = `struct{ X int }` |
| 兼容性 | 需显式转换 | 自动兼容 |

Go 使用**结构相等**用于接口实现，**名义相等**用于类型定义。

---

## 与接口的关系

### 隐式实现

```go
// 定义接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 定义类型（不声明实现）
type File struct { }
func (f File) Read(p []byte) (n int, err error) { ... }

// 自动满足
var r Reader = File{}  // ✅ 合法
```

形式化：

```
methods(File) ⊇ methods(Reader)
⇒ File <: Reader
```

### 接口组合

```go
type ReadWriter interface {
    Reader      // 嵌入接口
    Writer      // 嵌入接口
}
```

形式化：

```
methods(ReadWriter) = methods(Reader) ∪ methods(Writer)

t <: ReadWriter  ⟺  t <: Reader ∧ t <: Writer
```

---

## 类型推断

### 函数参数

```go
func Process(r Reader) { ... }

Process(File{})  // 自动推断 File <: Reader
```

### 泛型中的结构类型

```go
func Process[T Reader](x T) { ... }

Process(File{})  // T 推断为 File，满足 File <: Reader
```

---

## 优缺点

### 优点

1. **解耦**: 类型和接口独立定义
2. **灵活**: 后期添加接口实现
3. **简洁**: 无需显式声明
4. **鸭子类型**: "如果它走起来像鸭子，叫起来像鸭子，那它就是鸭子"

### 缺点

1. **隐式性**: 实现关系不明确
2. **意外实现**: 可能意外满足接口
3. **文档**: 需要工具辅助找出实现关系
4. **兼容性**: 修改方法可能影响多个接口

---

## 形式化性质

### 自反性

```
t <: t    (所有类型都是自身的子类型)
```

### 传递性

```
t₁ <: t₂    t₂ <: t₃
─────────────────────
t₁ <: t₃
```

### 反对称性（接口）

```
t_I₁ <: t_I₂    t_I₂ <: t_I₁
────────────────────────────
t_I₁ = t_I₂    (方法集相同)
```

---

## 实际影响

### 向后兼容

```go
// 旧代码
type OldInterface interface {
    MethodA()
}

// 新代码添加方法
type NewInterface interface {
    MethodA()
    MethodB()  // 新增
}

// 旧实现自动兼容（只要实现了新方法）
```

### Mock 测试

```go
// 易于创建 Mock
type MockReader struct { }
func (m MockReader) Read(p []byte) (n int, err error) { ... }

// 自动满足 Reader 接口
```

---

## 参考

- Go Language Specification: Interface types
- "Type Systems" (Pierce), Chapter on Subtyping
- "Featherweight Go" (OOPSLA 2020)
