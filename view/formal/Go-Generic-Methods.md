# Go泛型方法形式化 (2026提案 #77273)

## 思维导图：泛型方法架构

```
                    泛型方法 (Generic Methods)
                           │
            ┌──────────────┼──────────────┐
            ▼              ▼              ▼
    类型参数位置      接收者泛型       方法级约束
            │              │              │
    ┌───────┴───────┐ ┌────┴────┐  ┌──────┴──────┐
    ▼               ▼ ▼         ▼  ▼             ▼
函数级[T any]    struct{      func (s *Slice[T])
                 data []T     Map[U any](f func(T)U)
                 }             *Slice[U]
            │              │              │
            └──────────────┼──────────────┘
                           ▼
            ┌───────────────────────────────┐
            │   关键约束：独立于接口           │
            │                               │
            │  • 接口方法不能有类型参数        │
            │  • 泛型具体方法不匹配接口        │
            │  • 通过类型推断调用              │
            └───────────────────────────────┘
```

---

## 1. 形式化定义

### 1.1 泛型方法语法

**定义 1.1.1** (泛型方法).
泛型方法 **m = (R, P̄, C̄, T̄, U)**：

- **R**: 接收者类型（可能含类型参数）
- **P̄**: 方法级类型参数
- **C̄**: 方法级约束
- **T̄**: 参数类型
- **U**: 返回类型

**BNF (泛型方法)**:

```
<generic_method>  ::= "func" "(" <receiver> ")" <method_name>
                      "[" <method_type_params> "]"
                      "(" <params> ")" <result> <body>

<receiver>        ::= <ident> <type>  (* 类型可能含类型参数 *)

<method_type_params> ::= <type_param> ("," <type_param>)*

<method_call>     ::= <expr> "." <method_name> <type_args>? "(" <args>? ")"
```

**示例**:

```go
// 接收者含类型参数T，方法引入新类型参数U
func (s *Slice[T]) Map[U any](f func(T) U) *Slice[U] {
    result := make(Slice[U], len(*s))
    for i, v := range *s {
        result[i] = f(v)
    }
    return &result
}

// 调用：类型推断
s := &Slice[int]{1, 2, 3}
s.Map(func(x int) string { return fmt.Sprint(x) })  // U推断为string
```

---

## 2. 类型系统

### 2.1 方法类型规则

**规则 2.1.1** (泛型方法类型).
$$
\frac{\Gamma \vdash recv: T[\bar{A}] \quad T[\bar{A}] \text{ has method } m[\bar{P}]: \forall(\bar{P}).(\bar{\tau}) \to \tau}{\Gamma \vdash recv.m[\bar{B}](\bar{e}) : \tau[\bar{A}/\bar{P}_{recv}, \bar{B}/\bar{P}_{method}]}
$$

### 2.2 类型推断

**定义 2.2.1** (方法调用推断).
方法调用 **recv.m(e₁, e₂, ...)** 的类型推断：

1. 从recv类型确定接收者类型参数
2. 从参数推断方法类型参数
3. 结合约束求解

**示例推导**:

```
s.Map(strconv.Itoa)

已知:
- s: Slice[int]  → T = int
- strconv.Itoa: func(int) string

推导:
- 匹配 func(T) U 与 func(int) string
- T = int (已知)
- U = string

结果: Map[string]
```

---

## 3. 与接口的关系

### 3.1 关键约束

**约束 3.1.1** (接口方法限制).
接口方法声明**不能**有类型参数。

```go
// 非法：接口方法不能有类型参数
type Container[T any] interface {
    Map[U any](func(T) U) Container[U]  // ERROR!
}
```

### 3.2 方法集

**定义 3.2.1** (具体方法 vs 接口方法).

- 具体类型可以有泛型方法
- 接口只能有非泛型方法
- 泛型具体方法**不匹配**接口方法

**示例**:

```go
type Mapper[T any] interface {
    Map(func(T) any) any  // 非泛型
}

func (s *Slice[T]) Map[U any](func(T) U) *Slice[U]  // 泛型

// Slice[T] 不满足 Mapper[T] 因为方法签名不同
```

---

## 4. 语法扩展

### 4.1 语法变更

**BNF (语法调整)**:

```
(* 原语法：类型参数只在操作数名后 *)
<old_call>        ::= <operand_name> "[" <type_args> "]"

(* 新语法：类型参数作为PrimaryExpr的一部分 *)
<primary_expr>    ::= <operand>
                    | <conversion>
                    | <method_expr>
                    | <primary_expr> <selector>
                    | <primary_expr> "[" <index> "]"
                    | <primary_expr> "[" <slice> "]"
                    | <primary_expr> <type_assertion>
                    | <primary_expr> <arguments>
                    | <primary_expr> <type_args>  (* 新增 *)
```

### 4.2 解析策略

类型参数与索引表达式语法相似，解析时：

- **T[int]**：类型实例化（类型参数）
- **expr[int]**：索引表达式
- **expr.Method[type]**：方法类型参数

---

## 5. 应用场景

### 5.1 集合操作

```go
func (s *Slice[T]) Filter(pred func(T) bool) *Slice[T]
func (s *Slice[T]) Reduce[U any](init U, f func(U, T) U) U
func (s *Slice[T]) FlatMap[U any](f func(T) *Slice[U]) *Slice[U]
```

### 5.2 链式调用

```go
s.Map(strings.ToUpper).
  Filter(func(s string) bool { return len(s) > 3 }).
  Reduce(0, func(acc int, s string) int { return acc + len(s) })
```

---

## 6. 元理论

### 6.1 类型安全

**定理 6.1.1** (类型安全).
正确类型检查的泛型方法调用不会导致运行时类型错误。

### 6.2 向后兼容

**定理 6.2.1** (兼容性).
引入泛型方法不改变现有非泛型代码的行为。

---

## 参考文献

1. Go Authors. (2026). *Proposal: Generic Methods for Go (#77273)*.
2. Go Authors. (2026). *Go Specification Update*.

---

*文档版本: 2026-03-29 | 提案状态: Open | 预计版本: Go 1.26+*
