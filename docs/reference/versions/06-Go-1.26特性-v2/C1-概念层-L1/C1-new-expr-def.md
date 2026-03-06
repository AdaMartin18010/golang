# new表达式扩展 (new Expression Extension)

> **文档层级**: C1-概念层 (Concept Layer L1)
> **文档类型**: 概念定义 (Concept Definition)
> **形式化基础**: [A6-new表达式扩展公理](../C2-原理层-L2/C2-公理系统.md#A6)
> **最后更新**: 2026-03-06

---

## 一、概念定义

### 1.1 形式化定义

```
new表达式 : Go 1.26引入的内置表达式，允许用值初始化new分配的内存

语法形式:
  new(表达式)

类型规则:
  Γ ⊢ e : T    T is value type    T is addressable
  ─────────────────────────────────────────────────
  Γ ⊢ new(e) : *T

语义等价:
  new(v) ≡ &v'  其中v'是v的副本
```

### 1.2 语法演变

| 版本 | 语法 | 说明 |
|------|------|------|
| Go 1.0-1.25 | `new(T)` | 仅支持类型参数，分配零值 |
| **Go 1.26** | `new(expr)` | ✅ 支持表达式，分配并初始化 |

### 1.3 与传统方式对比

```go
// Go 1.25及之前：繁琐的写法
value := 42
ptr := &value           // 需要中间变量
// 或
ptr := &[]int{42}[0]    // 丑陋的hack

// Go 1.26：简洁优雅
ptr := new(42)          // ✨ 一行完成
```

---

## 二、核心性质

### 2.1 语义等价性 (Th1.1)

由[Th1.1](../R-参考层/R-定理索引.md#Th1.1)保证：

```
定理: ∀T: Type, v: T. new(T(v)) ≡ &T(v)

证明概要:
  new(v) = alloc(typeof(v)) ; store(alloc(typeof(v)), v)
  &v     = addressof(v)
  两者都返回指向T(v)副本的地址
  ∴ 语义等价
```

### 2.2 类型安全性 (Th1.3)

```
定理: Γ ⊢ e : T → Γ ⊢ new(e) : *T

推论: new表达式保持类型安全性
```

### 2.3 逃逸行为

```go
// new(expr)的逃逸行为与&相同
// 由编译器的逃逸分析决定

func example() *int {
    return new(42)  // 逃逸到堆
}

func example2() {
    ptr := new(42)  // 可能栈分配（如果编译器能证明不逃逸）
    _ = *ptr
}
```

---

## 三、使用场景

### 3.1 可选字段配置

```go
type Config struct {
    Timeout  *int
    MaxConns *int
}

// Go 1.26 简洁写法
config := Config{
    Timeout:  new(30),
    MaxConns: new(100),
}
```

**应用模式**: [C3-可选字段模式](../C3-实践层-L3/C3-可选字段模式.md)

### 3.2 延迟初始化

```go
type LazyValue struct {
    ptr *expensiveStruct
}

func (l *LazyValue) init() {
    if l.ptr == nil {
        l.ptr = new(expensiveStruct{...})  // ✨ 简洁
    }
}
```

### 3.3 构造者模式

```go
type Builder struct {
    port *int
}

func (b *Builder) Port(p int) *Builder {
    b.port = new(p)  // ✨ 流畅API
    return b
}
```

---

## 四、约束与限制

### 4.1 类型约束

```go
// ✅ 支持：值类型
new(42)
new("hello")
new(struct{ X int }{42})

// ❌ 不支持：非值类型
new(make([]int, 10))     // 错误：make返回非值
new(func() {})           // 错误：函数非值类型上下文
```

### 4.2 可寻址性要求

```go
// ✅ 支持：可寻址值
x := 42
new(x)                   // 变量
new(42)                  // 常量
new(getValue())          // 函数返回值（临时值）

// ⚠️ 注意：某些表达式可能需要额外处理
```

---

## 五、相关概念

### 5.1 概念关系

```mermaid
graph TD
    A[new表达式] --> B[内存分配 A1]
    A --> C[指针语义 A3]
    A --> D[new(expr)公理 A6]

    D --> E[Th1.1 语义等价]
    D --> F[Th1.3 类型安全]

    E --> G[可选字段模式]
    E --> H[延迟初始化模式]
    E --> I[构造者模式]
```

### 5.2 相关文档

- **形式化**: [C2-new-expr-formal](../C2-原理层-L2/C2-new-expr-formal.md)
- **定理**: [Th1.1](../R-参考层/R-定理索引.md#Th1.1), [Th1.3](../R-参考层/R-定理索引.md#Th1.3)
- **模式**: [C3-可选字段模式](../C3-实践层-L3/C3-可选字段模式.md)
- **对比**: `&T{}` vs `new(T())`

---

## 六、历史演进

```
Go 1.0: new(T) 引入 - 分配零值内存
   ↓
Go 1.26: new(expr) 扩展 - 支持表达式初始化
   ↓
未来: 可能的进一步扩展（如new(expr, options...)）
```

---

## 七、最佳实践

| 场景 | 推荐 | 示例 |
|------|------|------|
| 需要初始化值 | ✅ `new(v)` | `new(42)` |
| 只需零值 | ✅ `new(T)` | `new(int)` |
| 取已有变量地址 | ✅ `&x` | `&value` |
| 复合字面量 | ✅ `&T{...}` | `&Config{...}` |

---

**概念分类**: 语言特性 - 表达式
**Go版本**: 1.26+
**依赖公理**: A1, A2, A3, A6
**支持定理**: Th1.1, Th1.3
