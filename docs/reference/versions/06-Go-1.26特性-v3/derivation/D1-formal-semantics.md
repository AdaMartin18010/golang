# D1: 形式语义

> **层级**: 推导层 (Derivation)
> **地位**: Go 1.26特性的形式语义
> **依赖**: [F1](../foundation/F1-metalanguage.md), [F2](../foundation/F2-axioms.md), [F3](../foundation/F3-core-concepts.md)

---

## 1. 操作语义

### 1.1 表达式求值

```
小步语义: e → e'

基本规则:

[Const]
─────────────
n → n          (常量已经是值)

[Add-Left]
e₁ → e₁'
─────────────────
e₁ + e₂ → e₁' + e₂

[Add-Right]
e₁ ⇓ v₁   e₂ → e₂'
─────────────────────
e₁ + e₂ → v₁ + e₂'

[Add]
n₁ + n₂ = n
─────────────────────
n₁ + n₂ → n
```

### 1.2 new表达式语义

```
[new-step]
e → e'
─────────────────────
new(e) → new(e')

[new-value]
───────────────────────────────────────────────── [A6]
new(v) → let p = alloc(typeof(v)) in store(p, v); p
```

### 1.3 内存操作语义

```
[Alloc]
───────────────────────────────────────────────── [A1]
(σ, alloc(T)) → (σ', p)
其中:
  p ∉ dom(σ)
  σ' = σ[p ↦ uninit(T)]

[Store]
p ∈ dom(σ)
───────────────────────────────────────────────── [A2]
(σ, store(p, v)) → (σ[p ↦ v], ())

[Load]
p ∈ dom(σ)   σ(p) = v
─────────────────────────────────────────────────
(σ, load(p)) → (σ, v)
```

---

## 2. 类型语义

### 2.1 new表达式的类型规则

```
[new-typing]
Γ ⊢ e : T    T is addressable
────────────────────────────────
Γ ⊢ new(e) : *T

[new-deref]
Γ ⊢ e : *T
────────────────────────────────
Γ ⊢ *e : T

[new-addr]
Γ ⊢ e : T    T is not pointer
────────────────────────────────
Γ ⊢ &e : *T
```

### 2.2 泛型类型规则

```
[Generic-Fun]
Γ, X: Type ⊢ fun(x X) R { body } : X → R
─────────────────────────────────────────────
Γ ⊢ func f[X any](x X) R { body } : ∀X. X → R

[Generic-Inst]
Γ ⊢ f : ∀X. F(X)    Γ ⊢ T : Type
────────────────────────────────
Γ ⊢ f[T] : F(T)

[Recursive-Constraint]
Γ ⊢ T : Type    T satisfies C[T]
────────────────────────────────
Γ ⊢ T valid under C[T C[T]]
```

---

## 3. 指称语义

### 3.1 语义函数定义

```
⟦·⟧ : Expression → Environment → State → (Value, State)

基本表达式:
⟦n⟧ρσ = (n, σ)                    (常量)
⟦x⟧ρσ = (ρ(x), σ)                 (变量)
⟦e₁ + e₂⟧ρσ =
  let (v₁, σ₁) = ⟦e₁⟧ρσ in
  let (v₂, σ₂) = ⟦e₂⟧ρσ₁ in
  (v₁ + v₂, σ₂)                   (加法)
```

### 3.2 new表达式的指称语义

```
⟦new(e)⟧ρσ =
  let (v, σ₁) = ⟦e⟧ρσ in           (求值e)
  let T = typeof(v) in              (获取类型)
  let p = fresh_address() in        (新地址)
  let σ₂ = σ₁[p ↦ v] in             (存储值)
  (p, σ₂)                           (返回指针)

等价于:
⟦new(e)⟧ρσ = ⟦&e'⟧ρσ
  其中e'是e的副本
```

---

## 4. 语义等价性

### 4.1 上下文等价

```
定义 (上下文等价):

e₁ ≅ e₂ 当且仅当 ∀C. C[e₁] ⇓ v ↔ C[e₂] ⇓ v

其中C是求值上下文
```

### 4.2 new表达式的等价性

```
定理 D1.1 (new语义等价)
  new(v) ≅ &v'

  其中v'是v的副本

证明:
  对任意上下文C:

  C[new(v)] 的求值:
    1. 分配新内存p
    2. 存储v的副本到p
    3. 在C中使用p

  C[&v'] 的求值:
    1. 创建v的副本v'
    2. 获取v'的地址p
    3. 在C中使用p

  两种情况C的行为相同（都获得指向v副本的指针）
  ∴ new(v) ≅ &v'
```

---

**下一章**: [D2-类型系统](D2-type-system.md) - 完整的类型规则系统
