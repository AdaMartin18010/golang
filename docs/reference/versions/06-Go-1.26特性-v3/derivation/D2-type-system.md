# D2: 类型系统

> **层级**: 推导层 (Derivation)
> **地位**: Go 1.26的完整类型规则
> **依赖**: D1

---

## 1. 类型判断

### 1.1 基本形式

```
类型判断: Γ ⊢ e : T

Γ: 类型环境 (变量→类型的映射)
e: 表达式
T: 类型
```

### 1.2 环境操作

```
空环境: ∅
扩展: Γ, x:T
查找: Γ(x) = T 如果x在Γ中
```

---

## 2. Go 1.26特有规则

### 2.1 new表达式规则

```
[new-intro]
Γ ⊢ e : T    T is not nil    T is addressable
───────────────────────────────────────────────
Γ ⊢ new(e) : *T

[new-elim]
Γ ⊢ e : *T
────────────────────────────────
Γ ⊢ *e : T

[new-equiv]
Γ ⊢ e : T
────────────────────────────────
Γ ⊢ new(e) ≡ &e' : *T    (e' fresh copy of e)
```

### 2.2 递归泛型规则

```
[recursive-constraint-def]
Γ, X: Type, X: C[X] ⊢ interface { methods } : Constraint
─────────────────────────────────────────────────────────
Γ ⊢ type C[X C[X]] interface { methods } : Constraint

[recursive-constraint-sat]
Γ ⊢ T : Type
Γ ⊢ T implements C[T]
Γ ⊢ structurally_recursive(C[T])
────────────────────────────────
Γ ⊢ T satisfies C[T C[T]]

[recursive-fun]
Γ, X: Type, X: C[X] ⊢ fun(x X) R { body } : X → R
───────────────────────────────────────────────────
Γ ⊢ func f[X C[X]](x X) R { body } : ∀X:C[X]. X → R
```

---

## 3. 类型安全

### 3.1 保持性 (Preservation)

```
定理 (类型保持):
  如果 Γ ⊢ e : T 且 e → e'，那么 Γ ⊢ e' : T

证明 (对求值规则归纳):
  基本情况: e已经是值，无e'，定理空真成立

  归纳步骤:
    情况[new-step]:
      e → e' 蕴含 new(e) → new(e')
      由归纳假设，Γ ⊢ e' : T
      由[new-intro]，Γ ⊢ new(e') : *T
      ∴ 保持
```

### 3.2 进展性 (Progress)

```
定理 (进展):
  如果 Γ ⊢ e : T，那么e是值或存在e'使得e → e'

证明 (对e的结构归纳):
  基本情况: e是值，定理成立

  归纳步骤:
    情况new(e):
      如果e是值v，则new(v) → p (分配并存储)
      如果e → e'，则new(e) → new(e')
      ∴ 可进展
```

### 3.3 类型安全定理

```
定理 (类型安全):
  如果 Γ ⊢ e : T，那么:
  1. e不会陷入 stuck 状态（除非panic）
  2. e的求值结果具有类型T

推论:
  良类型的Go 1.26程序不会出"类型错误"
```

---

**下一章**: [D3-定理与证明](D3-theorems.md) - 核心定理的严格证明
