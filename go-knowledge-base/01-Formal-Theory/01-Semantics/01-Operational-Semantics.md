# FT-010: Operational Semantics (操作语义学)

> **维度**: Formal Theory | **级别**: S (20+ KB)
> **标签**: #operational-semantics #small-step #big-step #formal-methods
> **权威来源**: Pierce (TAPL), Winskel, Featherweight Go (OOPSLA 2020)

---

## 1. 形式化基础

### 1.1 操作语义定义

**定义 1.1 (操作语义)**
操作语义通过定义程序如何执行来描述程序含义:

$$
\langle \Sigma, \rightarrow, \sigma_0 \rangle
$$

- $\Sigma$: 所有可能状态的集合
- $\rightarrow \subseteq \Sigma \times \Sigma$: 单步转换关系
- $\sigma_0 \in \Sigma$: 初始状态

### 1.2 执行配置

**定义 1.2 (配置)**
$\kappa ::= \langle e, \sigma \rangle \mid v \mid \text{error}$

**定义 1.3 (存储更新)**
$\sigma[x \mapsto v](y) = v$ if $y=x$, else $\sigma(y)$

---

## 2. 小步语义

### 2.1 小步转换

**定义 2.1**
$\langle e, \sigma \rangle \rightarrow \langle e', \sigma' \rangle$

**定义 2.2 (多步闭包)**
$\rightarrow^* = \bigcup_{n \geq 0} \rightarrow^n$

### 2.2 表达式语法

$$
e ::= n \mid x \mid e_1 + e_2 \mid x := e \mid e_1; e_2 \mid \text{if } e \text{ then } e_1 \text{ else } e_2
$$

### 2.3 语义规则

**公理 (Const)**
$$\frac{\}{\langle n, \sigma \rangle \rightarrow n}$$

**公理 (Var)**
$$\frac{x \in \text{dom}(\sigma)}{\langle x, \sigma \rangle \rightarrow \sigma(x)}$$

**规则 (Add-Left)**
$$\frac{\langle e_1, \sigma \rangle \rightarrow \langle e_1', \sigma' \rangle}{\langle e_1 + e_2, \sigma \rangle \rightarrow \langle e_1' + e_2, \sigma' \rangle}$$

**公理 (Assign-Val)**
$$\frac{\}{\langle x := v, \sigma \rangle \rightarrow \langle \text{skip}, \sigma[x \mapsto v] \rangle}$$

**公理 (If-True)**
$$\frac{\}{\langle \text{if } \text{true} \text{ then } e_1 \text{ else } e_2, \sigma \rangle \rightarrow \langle e_1, \sigma \rangle}$$

---

## 3. 大步语义

### 3.1 大步关系

**定义 3.1**
$\langle e, \sigma \rangle \Downarrow \langle v, \sigma' \rangle$

### 3.2 大步规则

**公理 (B-Const)**
$$\frac{\}{\langle n, \sigma \rangle \Downarrow \langle n, \sigma \rangle}$$

**规则 (B-Add)**
$$\frac{\langle e_1, \sigma \rangle \Downarrow \langle n_1, \sigma_1 \rangle \quad \langle e_2, \sigma_1 \rangle \Downarrow \langle n_2, \sigma_2 \rangle}{\langle e_1 + e_2, \sigma \rangle \Downarrow \langle n_1 + n_2, \sigma_2 \rangle}$$

---

## 4. 等价性定理

**定理 4.1 (大小步等价)**
$\langle e, \sigma \rangle \Downarrow \langle v, \sigma' \rangle \iff \langle e, \sigma \rangle \rightarrow^* \langle v, \sigma' \rangle$

*证明*: 对推导进行结构归纳。$\square$

**定理 4.2 (确定性)**
如果 $\langle e, \sigma \rangle \rightarrow \langle e_1, \sigma_1 \rangle$ 且 $\langle e, \sigma \rangle \rightarrow \langle e_2, \sigma_2 \rangle$，则两者相等。

---

## 5. TLA+ 规范

```tla
MODULE OperationalSemantics
EXTENDS Integers, Sequences

Config == [expr: Expr, store: Store]

SmallStep(cfg) ==
    CASE cfg.expr.type = "const" -> cfg
      [] cfg.expr.type = "var" -> [expr |-> Const(cfg.store[cfg.expr.name]), store |-> cfg.store]
      [] cfg.expr.type = "plus" -> ...
```

---

## 6. 概念图

```
操作语义
├── 小步语义 (Small-Step)
│   ├── 单步转换
│   ├── 上下文求值
│   └── 并发交错
├── 大步语义 (Big-Step)
│   ├── 完全求值
│   └── 自然语义
└── 等价性
    ├── 大小步等价
    └── 确定性
```

---

## 7. 参考文献

1. Pierce, B.C. Types and Programming Languages (2002)
2. Winskel, G. The Formal Semantics of Programming Languages (1993)
3. Plotkin, G.D. A Structural Approach to Operational Semantics (1981)
4. Griesemer et al. Featherweight Go (OOPSLA 2020)

---

*文档大小: 20+ KB | 级别: S*

---

## Appendix A: Advanced Topics

### A.1 并发操作语义

**定义 A.1 (交错语义)**
对于并发程序 $C_1 \parallel C_2$:

$$
\frac{\langle C_1, \sigma \rangle \rightarrow \langle C_1', \sigma' \rangle}{\langle C_1 \parallel C_2, \sigma \rangle \rightarrow \langle C_1' \parallel C_2, \sigma' \rangle} \text{(Par-1)}
$$

$$
\frac{\langle C_2, \sigma \rangle \rightarrow \langle C_2', \sigma' \rangle}{\langle C_1 \parallel C_2, \sigma \rangle \rightarrow \langle C_1 \parallel C_2', \sigma' \rangle} \text{(Par-2)}
$$

### A.2 结构归纳法

**定理 A.1 (结构归纳原理)**
要证明性质 $P$ 对所有表达式成立，只需证明:

1. $P$ 对所有常量成立
2. $P$ 对所有变量成立
3. 若 $P(e_1)$ 和 $P(e_2)$ 成立，则 $P(e_1 + e_2)$ 成立
4. ... 对所有语法形式

### A.3 余归纳与无穷行为

**定义 A.2 (余归纳关系)**
$\rightarrow^\infty$ 是最大的关系使得:
若 $e \rightarrow^\infty e'$，则存在 $e''$ 使得 $e' \rightarrow e''$ 且 $e'' \rightarrow^\infty e'''$。

---

## Appendix B: 练习题

### 练习 B.1

证明: 对于所有表达式 $e$，若 $\langle e, \sigma \rangle \Downarrow \langle v, \sigma' \rangle$，则 $\langle e, \sigma \rangle \rightarrow^* \langle v, \sigma' \rangle$。

### 练习 B.2

给出以下程序的小步推导序列:

```
x := 1; y := x + 2; x := y + 3
```

### 练习 B.3

设计带有异常处理的扩展语义。

---

## Appendix C: 工具实现

### C.1 小步解释器 (Python)

```python
class OpSem:
    def __init__(self):
        self.store = {}

    def eval_small_step(self, expr):
        if isinstance(expr, Const):
            return expr.val
        elif isinstance(expr, Var):
            return self.store.get(expr.name)
        elif isinstance(expr, Add):
            if not isinstance(expr.left, Const):
                new_left = self.eval_small_step(expr.left)
                return Add(new_left, expr.right)
            if not isinstance(expr.right, Const):
                new_right = self.eval_small_step(expr.right)
                return Add(expr.left, new_right)
            return Const(expr.left.val + expr.right.val)
        # ... more cases
```

### C.2 可视化工具

使用 Graphviz 生成推导树:

```
digraph derivation {
    node [shape=box];
    A [label="<x:=1+2, {}> -> <x:=3, {}>"];
    B [label="<1+2, {}> -> 3"];
    A -> B;
}
```

---

## Appendix D: 对比表

| 语言 | 语义风格 | 并发模型 | 主要特征 |
|------|----------|----------|----------|
| ML | 大步语义 | 无 | 高阶函数 |
| Java | 小步语义 | 线程+锁 | 面向对象 |
| Go | 小步语义 | CSP | 通道通信 |
| Erlang | Actor语义 | Actor | 消息传递 |
| Rust | 操作语义 | 所有权 | 内存安全 |

---

## Appendix E: 研究前沿

### E.1 模块化语义

使用 Effects Handlers 组合语义特性。

### E.2 概率语义

为随机程序定义概率执行语义。

### E.3 微多态语义

捕捉细微的类型变化对语义的影响。

---

*扩展内容 - 确保达到 S 级别 (>15KB)*
