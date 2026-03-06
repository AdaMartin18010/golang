# F3: 核心概念定义

> **层级**: 基础层 (Foundation)
> **地位**: Go 1.26语言特性的严格数学定义
> **依赖**: [F1-元语言](F1-metalanguage.md), [F2-公理系统](F2-axioms.md)

---

## C1: new表达式 (new Expression)

### 数学定义

```
┌─────────────────────────────────────────────────────────────┐
│ 定义 C1.1 (new表达式语法)                                     │
│                                                              │
│   new_expr ::= new(expression)                              │
│                                                              │
│   其中 expression 是任意可求值的Go表达式                       │
└─────────────────────────────────────────────────────────────┘
```

### 语义定义

```
定义 C1.2 (new表达式语义)

给定:
  e: Expression, 类型为T

语义:
  ⟦new(e)⟧ρ =
    let v = ⟦e⟧ρ in                    (求值e)
    let T = typeof(v) in               (获取类型)
    let (p, σ₁) = alloc(σ, T) in       [A1] 分配内存
    let σ₂ = store(σ₁, p, v) in        [A2] 存储值
    (p, σ₂)                            (返回指针和新状态)

类型:
  Γ ⊢ e : T
  ─────────────────
  Γ ⊢ new(e) : *T
```

### 核心性质

```
定理 C1.1 (初始化保证)
  ∀e: Expression.
    let p = new(e) in
    *p = e

证明:
  由A6定义，new(e)执行store(p, e)
  由A2，store后load得到存储值
  ∴ *p = e

定理 C1.2 (内存分配)
  new(e)返回的指针指向新分配的内存

证明:
  由A6，new调用alloc
  由A1，alloc返回新内存地址
  ∴ 内存是新的
```

### 与相关概念的关系

```
new(e) 与 &e' 的关系:

  new(e):  分配内存 → 复制e → 返回地址
  &e':     创建e的副本e' → 返回e'的地址

  关键区别:
    - new(e)的内存位置由运行时决定（栈或堆）
    - &e'的内存位置由变量作用域决定

  语义等价:
    在"获得可修改的e的副本的指针"这个语义上，两者等价
```

---

## C2: 递归泛型约束 (Recursive Generic Constraint)

### 数学定义

```
┌─────────────────────────────────────────────────────────────┐
│ 定义 C2.1 (递归约束语法)                                      │
│                                                              │
│   recursive_constraint ::=                                   │
│     type Name[T Name[T]] interface { method_set }           │
│                                                              │
│ 其中约束定义中Name引用了自身                                   │
└─────────────────────────────────────────────────────────────┘
```

### 不动点表示

```
定义 C2.2 (递归约束的不动点表示)

递归约束 C[T C[T]] 可以表示为不动点:

  C = μX.F(X)

其中:
  μ 是最小不动点算子
  F 是约束构造函数

展开:
  C = F(C) = F(F(C)) = F(F(F(C))) = ...

结构递归要求:
  F 必须是结构递归的（每次应用减小参数规模）
```

### 类型满足关系

```
定义 C2.3 (约束满足)

T satisfies C[T C[T]] ↔
  ∀m ∈ methods(C). T implements m
  ∧
  ∀recursive_ref ∈ C. terminates(unfold(C, T))

其中:
  - T实现C的所有方法
  - 约束展开在有限步内终止 [A7]
```

### 终止性条件

```
定义 C2.4 (结构递归)

约束C关于类型参数T是结构递归的，如果:

  depth(unfold(C, T, n+1)) < depth(unfold(C, T, n))

  其中depth是类型定义的深度度量

直观:
  每次递归引用都指向"更小"的类型
  例如：树的子节点 < 父节点
```

### 典型实例

```
实例 C2.1 (树节点)

定义:
  type Node[T Node[T]] interface {
    Children() []T
  }

分析:
  方法Children()返回[]T
  T是Node[T]的实现类型
  子节点是父节点的子结构
  ∴ 结构递归，展开终止

满足的类型:
  type FileNode struct {
    Children []*FileNode
  }
  func (f *FileNode) Children() []*FileNode { return f.Children }
```

---

## C3: GreenTeaGC

### 数学模型

```
┌─────────────────────────────────────────────────────────────┐
│ 定义 C3.1 (GC状态机)                                          │
│                                                              │
│   State ::= IDLE | MARK | SWEEP                             │
│                                                              │
│   转换:                                                     │
│     IDLE --trigger--> MARK --complete--> SWEEP --complete--> IDLE│
│     MARK --write-barrier--> MARK                            │
└─────────────────────────────────────────────────────────────┘
```

### 并发标记算法

```
定义 C3.2 (三色标记)

对象颜色:
  WHITE: 未访问（潜在垃圾）
  GREY:  已访问，字段未处理
  BLACK: 已访问，字段已处理

算法:
  1. 初始化: 所有对象WHITE，根对象GREY
  2. 标记: 取GREY对象，处理其字段，标记为BLACK
  3. 写屏障: 修改指针时，将新目标标记为GREY
  4. 完成: 无GREY对象时，WHITE对象即为垃圾

正确性:
  由A8保证不会丢失存活对象
```

### 性能模型

```
定义 C3.3 (停顿时间分布)

设P为停顿时间随机变量:
  P(P < 1ms) = 0.99      [A8]
  P(P < 5ms) = 0.999
  E[P] < 0.5ms

增量标记预算:
  budget = target_pause_time × work_rate

  确保每个增量片的工作量可控
```

---

## C4: HPKE (Hybrid Public Key Encryption)

### 数学定义

```
┌─────────────────────────────────────────────────────────────┐
│ 定义 C4.1 (HPKE方案)                                          │
│                                                              │
│   HPKE = (KEM, KDF, AEAD)                                   │
│                                                              │
│ 其中:                                                       │
│   KEM:  Key Encapsulation Mechanism                         │
│   KDF:  Key Derivation Function                             │
│   AEAD: Authenticated Encryption with Associated Data       │
└─────────────────────────────────────────────────────────────┘
```

### 算法组件

```
定义 C4.2 (密钥封装)

Encap(pk):
  输入: 接收方公钥pk
  输出: (enc, shared_secret)

  性质:
    1. enc可公开传输
    2. 只有拥有对应私钥sk的接收方能解封装
    3. shared_secret均匀分布

Decap(sk, enc):
  输入: 接收方私钥sk, 封装enc
  输出: shared_secret

  正确性: Decap(sk, Encap(pk).enc) = Encap(pk).shared_secret
```

### 加密流程

```
定义 C4.3 (HPKE加密)

Seal(pk, pt, info):
  1. (enc, ss) = KEM.Encap(pk)
  2. key = KDF(ss, info)
  3. ct = AEAD.Seal(key, pt)
  4. return (enc, ct)

Open(sk, enc, ct, info):
  1. ss = KEM.Decap(sk, enc)
  2. key = KDF(ss, info)
  3. pt = AEAD.Open(key, ct)
  4. return pt
```

---

## 概念依赖图

```
A1-A3 (内存基础)
      ↓
   C1 (new)
      ↓
  ┌───┴───┐
  ↓       ↓
A6    应用模式

A5 (泛型)
      ↓
   C2 (递归泛型)
      ↓
  ┌───┴───┐
  ↓       ↓
A7    树遍历

A8 (GC)
      ↓
   C3 (GreenTeaGC)
      ↓
   低延迟优化
```

---

**下一层**: [D1-形式语义](../derivation/D1-formal-semantics.md) - 基于核心概念建立形式语义
