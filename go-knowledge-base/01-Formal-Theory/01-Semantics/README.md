# 形式语义学 (Formal Semantics)

> **目标**: 精确理解 Go 语言的数学含义

---

## 三种语义学方法

### 1. 操作语义 (Operational Semantics)

描述程序如何执行。

**小步语义** (Small-Step):

```
(e, σ) → (e', σ')

表示表达式 e 在状态 σ 下一步执行到 e'，状态变为 σ'
```

**大步语义** (Big-Step):

```
(e, σ) ⇓ (v, σ')

表示表达式 e 在状态 σ 下求值为 v，最终状态为 σ'
```

### 2. 指称语义 (Denotational Semantics)

将程序映射到数学对象。

```
〚e〛 : Environment → Value

表达式的含义是从环境到值的函数
```

### 3. 公理语义 (Axiomatic Semantics)

使用逻辑断言描述程序行为。

```
{P} C {Q}

如果在执行命令 C 前置条件 P 成立，
则执行后后置条件 Q 成立
```

---

## Featherweight Go

Go 语言的极简形式化子集。

### 语法

```
t ::=                         // 类型
  | t_I                       // 接口类型
  | t_S                       // 结构体类型

e ::=                       // 表达式
  | x                         // 变量
  | e.f                       // 字段选择
  | e.m(e)                    // 方法调用
  | t_S{e, ..., e}            // 结构体字面量
```

### 类型规则

```
Γ ⊢ e : t_S    fields(t_S) = ... f:t_f ...
────────────────────────────────────────────  (T-Field)
Γ ⊢ e.f : t_f

Γ ⊢ e_0 : t_0    methods(t_0) ∋ m(x:t_1)t_2    Γ ⊢ e_1 : t_1
─────────────────────────────────────────────────────────────  (T-Call)
Γ ⊢ e_0.m(e_1) : t_2
```

### 操作语义

**结构体字段选择**:

```
──────────────────────────────────  (R-Field)
t_S{v_1, ..., v_n}.f_i → v_i
```

**方法调用**:

```
body(m, t_S) = (x, e)
──────────────────────────────────────────────  (R-Call)
t_S{v...}.m(v') → e[x := v', this := t_S{v...}]
```

---

## 类型安全

### 保持性 (Preservation)

```
Theorem: If Γ ⊢ e : T and e → e', then Γ ⊢ e' : T

证明概要:
- 对推导规则进行归纳
- 每种归约保持类型
```

### 进展性 (Progress)

```
Theorem: If Γ ⊢ e : T, then either
  1. e is a value, or
  2. ∃e': e → e'

证明概要:
- 对类型推导进行归纳
- 每种类型形式要么已是值，要么可归约
```

---

## Go 扩展

### 实际 Go 与 FG 的差异

| 特性 | FG | 实际 Go |
|------|-----|---------|
| 赋值 | 无 | 有 |
| 循环 | 无 | 有 |
| 并发 | 无 | 有 (Goroutine/Channel) |
| 接口 | 简单 | 完整 |

---

## 参考

- Featherweight Go (OOPSLA 2020)
- Types and Programming Languages (Pierce)
