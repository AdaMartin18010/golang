# Featherweight Go (FG)

> **分类**: 形式理论
> **难度**: 专家
> **前置知识**: 操作语义、类型系统基础

---

## 概述

Featherweight Go (FG) 是 Go 语言的极简形式化演算，保留了 Go 的核心特性：

- 结构类型接口
- 方法调用
- 结构体嵌入

FG 去除了实际 Go 的复杂特性（赋值、循环、并发等），专注于核心类型系统。

---

## 语法

### 表达式

```
e ::=
  | x                      变量
  | e.f                    字段选择
  | e.m(e)                 方法调用
  | t_S{e, ..., e}         结构体字面量
  | assert(e, t)           类型断言
```

### 类型

```
t ::=
  | t_S                    结构体类型
  | t_I                    接口类型
```

### 声明

```
D ::=
  | type t_S struct {f t}           结构体声明
  | type t_I interface {m(x t) t}   接口声明
  | func (x t) m(y t) t { return e }  方法声明
```

---

## 类型规则

### 环境

```
Γ ::= ∅ | Γ, x: t          类型环境（变量到类型的映射）

Φ ::= ∅ | Φ, t_S: struct{f t}     结构体声明集合
     | Φ, t_I: interface{m(x t₁) t₂}  接口声明集合
     | Φ, method(t, m) = (x, y, e)   方法实现集合
```

### 表达式类型规则

#### 变量

```
x: t ∈ Γ
───────────  (T-Var)
Γ ⊢ x: t
```

#### 结构体字面量

```
t_S: struct{f₁ t₁, ..., fₙ tₙ} ∈ Φ
Γ ⊢ e₁: t₁  ...  Γ ⊢ eₙ: tₙ
────────────────────────────────────────  (T-Struct)
Γ ⊢ t_S{e₁, ..., eₙ}: t_S
```

#### 字段选择

```
Γ ⊢ e: t_S
fields(t_S) = ..., f: t, ...
──────────────────────────────  (T-Field)
Γ ⊢ e.f: t
```

#### 方法调用

```
Γ ⊢ e: t
method(t, m) = (x, y, e')
Γ ⊢ e_arg: t_arg
t_arg <: parameter_type(m, t)
──────────────────────────────────  (T-Call)
Γ ⊢ e.m(e_arg): return_type(m, t)
```

#### 类型断言

```
Γ ⊢ e: t
t' <: t  or  t <: t'        (可比较)
───────────────────────────  (T-Assert)
Γ ⊢ assert(e, t'): t'
```

---

## 子类型关系

### 结构体子类型

```
t_S: struct{f₁ t₁, ..., fₙ tₙ} <: t_S    (S-Refl)

t_S <: t_I    if    methods(t_I) ⊆ methods(t_S)    (S-Struct-Interface)
```

### 接口子类型

```
t_I₁ <: t_I₂    if    methods(t_I₂) ⊆ methods(t_I₁)    (S-Interface)

注意: 接口子类型是反变的（方法越多，子类型越小）
```

### 传递性

```
t₁ <: t₂    t₂ <: t₃
─────────────────────  (S-Trans)
t₁ <: t₃
```

---

## 操作语义

### 求值上下文

```
E ::=
  | []                      空上下文
  | E.f                     字段选择上下文
  | E.m(e)                  调用者上下文
  | v.m(E)                  参数上下文
  | t_S{v₁, ..., E, ..., eₙ}  结构体字面量上下文
```

### 归约规则

#### 字段选择

```
t_S{v₁, ..., vₙ}.fᵢ → vᵢ    (R-Field)

选择结构体第 i 个字段的值
```

#### 方法调用

```
(v: t).m(v') → e[x ↦ v, y ↦ v']    (R-Call)

其中 method(t, m) = func (x t) m(y t') t'' { return e }
```

#### 类型断言

```
assert(t_S{v...}, t_S) → t_S{v...}    (R-Assert-Struct)

assert(v, t_I) → v    if type(v) <: t_I    (R-Assert-Interface)
```

---

## 类型安全

### 保持性 (Preservation)

**定理**: 如果表达式 e 是良类型的，且 e 归约为 e'，那么 e' 也是良类型的，且类型相同。

```
Theorem (Preservation):
If Γ ⊢ e: t and e → e', then Γ ⊢ e': t.

证明概要 (对归约规则进行归纳):

情况 R-Field:
  前提: t_S{v₁, ..., vₙ}: t_S 且 fields(t_S) = ..., fᵢ: t, ...
  结论: vᵢ: t
  由结构体类型规则，vᵢ 的类型就是字段 fᵢ 声明的类型 t

情况 R-Call:
  前提: (v: t).m(v') 调用 method(t, m) = func (x t) m(y t₁) t₂ { return e }
  结论: e[x ↦ v, y ↦ v']: t₂
  由方法声明的类型检查，e 在环境 x: t, y: t₁ 下有类型 t₂
  代入后类型保持
```

### 进展性 (Progress)

**定理**: 良类型的表达式要么是值，要么可以进一步归约。

```
Theorem (Progress):
If Γ ⊢ e: t, then either:
  1. e is a value, or
  2. ∃e': e → e'.

证明概要 (对 e 的结构进行归纳):

情况 e = x (变量):
  如果 x 在环境 Γ 中有类型，那么它应该被替换为值
  在实际程序中，这意味着变量必须已初始化

情况 e = t_S{e₁, ..., eₙ} (结构体字面量):
  如果所有 eᵢ 都是值，则 e 是值
  否则，选择第一个非值的 eᵢ 进行归约

情况 e = e'.f (字段选择):
  如果 e' 是结构体值，应用 R-Field
  否则，对 e' 应用归纳假设

情况 e = e₁.m(e₂) (方法调用):
  如果 e₁ 是值且 e₂ 是值，应用 R-Call
  否则，归约 e₁ 或 e₂
```

### 类型安全定理

```
Theorem (Type Safety):
If ∅ ⊢ e: t and e →* v, then ∅ ⊢ v: t.

由 Preservation 和 Progress 直接得出。
```

---

## 与实际 Go 的差异

| 特性 | FG | Go |
|------|-----|-----|
| 赋值 | ❌ | ✅ |
| 循环 | ❌ | ✅ |
| 并发 | ❌ | ✅ |
| 指针 | ❌ | ✅ |
| 函数类型 | ❌ | ✅ |
| 包系统 | ❌ | ✅ |
| 方法集 | 简化 | 完整 |
| 接口 | 基础 | 完整 |

---

## 扩展到 FGG

Featherweight Generic Go (FGG) 在 FG 基础上添加泛型：

```
t ::=
  | t_S[t₁, ..., tₙ]      泛型结构体
  | t_I[t₁, ..., tₙ]      泛型接口
  | α                     类型变量
```

FGG 的形式化是理解 Go 1.18+ 泛型的理论基础。

---

## 实现练习

### 练习 1

为以下 Go 代码写出 FG 表示：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type File struct { ... }

func (f File) Read(p []byte) (n int, err error) { ... }
```

### 练习 2

证明以下子类型关系：

```
如果 t_S 实现了 t_I 的所有方法，则 t_S <: t_I
```

### 练习 3

推导以下表达式的完整求值序列：

```
t_Point{x: 3, y: 4}.Add(t_Point{x: 1, y: 2})
```

---

## 参考

- Featherweight Go (OOPSLA 2020)
- A Dictionary-Passing Translation of Featherweight Go (APLAS 2021)
- Generic Go to Go: Dictionary-Passing, Monomorphisation, and Hybrid (OOPSLA 2022)
