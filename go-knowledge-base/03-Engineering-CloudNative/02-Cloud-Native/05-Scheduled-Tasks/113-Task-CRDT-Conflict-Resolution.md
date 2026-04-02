# CRDT 冲突解决实现 (CRDT Conflict Resolution Implementation)

> **分类**: 工程与云原生
> **标签**: #crdt #conflict-free #eventual-consistency #distributed
> **参考**: Shapiro et al. "A comprehensive study of Convergent and Commutative Replicated Data Types"

---

## CRDT 理论基础

```
强一致性 (CP)               最终一致性 (AP) + CRDT
      │                             │
      ▼                             ▼
┌─────────────┐              ┌─────────────┐
│  Consensus  │              │   Merge     │
│  (Paxos/    │              │   Function  │
│   Raft)     │              │  (单调性保证) │
└─────────────┘              └─────────────┘
      │                             │
   高延迟                          低延迟
   高可用性损失                      始终可用
   需要协调                         无协调
```

---

## CRDT 数学定义

$$
\begin{aligned}
&\text{State-based CRDT (CvRDT):} \\
&S: \text{状态空间} \\
&\sqcup: S \times S \rightarrow S \text{ (合并函数)} \\
&\forall a, b \in S: a \sqcup b = b \sqcup a \text{ (交换律)} \\
&\forall a, b, c \in S: (a \sqcup b) \sqcup c = a \sqcup (b \sqcup c) \text{ (结合律)} \\
&\forall a \in S: a \sqcup a = a \text{ (幂等律)} \\
\\
&\text{Operation-based CRDT (CmRDT):} \\
&\forall o_1, o_2 \in \text{Operations}: \\
&\quad \text{if } \text{source}(o_1) \parallel \text{source}(o_2) \Rightarrow o_1 \circ o_2 = o_2 \circ o_1
\end{aligned}
$$
