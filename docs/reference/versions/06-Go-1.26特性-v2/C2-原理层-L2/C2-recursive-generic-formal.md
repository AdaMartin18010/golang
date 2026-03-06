# 递归泛型形式化理论

> **文档层级**: C2-原理层 (Principle Layer L2)
> **文档类型**: 形式语义 (Formal Semantics)
> **依赖**: [A7-递归约束公理](C2-公理系统.md#A7)
> **最后更新**: 2026-03-06

---

## 一、递归约束的形式定义

### 1.1 语法定义

```
递归约束 C[T C[T]] 的形式文法:

Constraint ::= type Name '[' TypeParam ConstraintRef ']' InterfaceType

ConstraintRef ::= Name '[' TypeParam ']'

InterfaceType ::= interface '{' MethodDecl* '}'
```

### 1.2 类型规则

```
[recursive-constraint-wellformed]
────────────────────────────────
C = μX.F(X)          (不动点表示)
F structurally_recursive
────────────────────────────────
⊢ C wellformed

[recursive-constraint-satisfaction]
────────────────────────────────
Γ ⊢ T : Type
Γ ⊢ T satisfies C[T]
────────────────────────────────
Γ ⊢ RecursiveType valid
```

---

## 二、终止性证明 (Th1.2)

### 2.1 定理陈述

```
Th1.2: ∀C[T C[T]]: Constraint. wellformed(C) → terminates(unfold(C))

定理: 结构良好的递归约束的展开过程保证终止。
```

### 2.2 完整证明

```
证明 (结构归纳法):

定义:
  - C[T C[T]] = μX.F(X)  (约束的不动点表示)
  - unfold(C, n) = Fⁿ(⊥) (n次展开)
  - depth(T) = 类型T的定义深度

引理 1 (结构递归递减):
  如果F是结构递归的，那么对于任何T，
  depth(unfold(C[T], n+1)) < depth(unfold(C[T], n))

  证明:
    结构递归定义要求递归调用时"子结构"规模减小。
    对于类型约束，这意味着:
    - 每次递归引用C[T]，T的复杂度必须降低
    - 例如：树的子节点数量少于父节点

引理 2 (良基性):
  Go类型系统是良基的，不存在无限下降链。
  即：对于任何类型T，depth(T)是有限的。

主证明:
  1. 设C是结构良好的递归约束
  2. 根据引理1，每次展开depth递减
  3. 根据引理2，depth有限且非负
  4. 由良序原理，递减序列必在有限步终止
  5. 因此unfold(C)在有限步内达到不动点

  ∴ terminates(unfold(C))
```

### 2.3 复杂度分析

```
时间复杂度: O(d) 其中d是类型的最大深度
空间复杂度: O(d) 用于存储展开过程中的中间约束

典型值:
  - 树结构: d ≤ 树高
  - 链表: d ≤ 长度
  - AST: d ≤ 语法树深度
```

---

## 三、类型安全证明

### 3.1 保持性 (Preservation)

```
定理: 如果 Γ ⊢ e : T 且 e → e'，那么 Γ ⊢ e' : T

对于递归泛型:
  实例化前: func Walk[T Node[T]](node T)
  实例化后: func Walk(node *TreeNode)

  类型保持因为:
    - 实例化替换类型参数为具体类型
    - 具体类型满足约束
    - 因此函数体类型正确
```

### 3.2 进展性 (Progress)

```
定理: 如果 Γ ⊢ e : T，那么e是值或可以规约

对于递归泛型:
  类型检查通过后，递归泛型函数可以:
    - 正常执行
    - 递归调用（每次递归参数规模减小）
    - 最终达到基础情况
```

---

## 四、不动点理论

### 4.1 最小不动点

```
定义: μX.F(X) 是F的最小不动点

性质:
  μX.F(X) = F(μX.F(X))

对于递归约束:
  C[T C[T]] = C[T C[T C[T]]]  (无限展开)

  实际实现使用有限近似:
    C₀ = ⊥
    C₁ = F(C₀)
    C₂ = F(C₁)
    ...
    Cₙ = F(Cₙ₋₁) 直到满足约束
```

### 4.2 约束求解

```
算法: 递归约束类型检查

输入: 类型T，约束C[T C[T]]
输出: T是否满足C

过程:
  1. 展开C[T C[T]]得到方法集M
  2. 对于M中的每个方法m:
     a. 检查T是否有方法m
     b. 如果m的参数涉及C，递归检查
  3. 如果所有方法都满足，返回true
  4. 如果递归深度超过阈值，报错（防止无限递归）
```

---

## 五、形式化示例

### 5.1 树遍历的形式化

```
定义约束:
  Node[T Node[T]] := interface {
    Children() []T
  }

定义函数:
  Walk[T Node[T]](node T, fn func(T)) : Unit

类型检查 Walk:
  1. T 必须满足 Node[T]
  2. 即 T 必须有 Children() []T 方法
  3. 在函数体内，node.Children() 返回 []T
  4. 对于每个 child ∈ children，Walk(child, fn) 类型正确
  5. 递归调用类型正确因为 child : T

终止性:
  假设树是有限的，每次递归处理子节点
  树深度有限 → 递归终止
```

### 5.2 约束展开示例

```
类型: FileNode

检查 FileNode satisfies Node[FileNode]:

  展开约束:
    Node[FileNode] = interface {
      Children() []FileNode
    }

  检查 FileNode:
    - 有 Children() []*FileNode
    - 注意: []*FileNode 与 []FileNode 不同！

  修正:
    定义 FileNode.Children() []Node 或调整约束
```

---

## 六、与其他语言对比

| 语言 | 递归约束 | 终止性检查 |
|------|----------|------------|
| Rust |  Trait bounds | 递归类型检查 |
| Haskell | 类型类 | 类型系统保证 |
| Scala | F-bounded polymorphism | 编译期检查 |
| Go 1.26 | 递归接口约束 | 结构递归检查 |

---

**相关定理**: [Th1.2](../R-参考层/R-定理索引.md#Th1.2)
**相关公理**: [A7](C2-公理系统.md#A7)
