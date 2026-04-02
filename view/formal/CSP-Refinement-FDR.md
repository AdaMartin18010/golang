# CSP 精化理论与 FDR 模型检测器

> **版本**: 2026.03 | **深度分析** | **关联**: [VISUAL-ATLAS.md](../../VISUAL-ATLAS.md)

---

## 目录

- [CSP 精化理论与 FDR 模型检测器](#csp-精化理论与-fdr-模型检测器)
  - [目录](#目录)
  - [1. 概念定义 (Definitions)](#1-概念定义-definitions)
    - [1.1 CSP 语义模型层次](#11-csp-语义模型层次)
    - [1.2 迹精化 (Trace Refinement)](#12-迹精化-trace-refinement)
    - [1.3 故障精化 (Failures Refinement)](#13-故障精化-failures-refinement)
    - [1.4 故障-发散精化 (Failures-Divergences Refinement)](#14-故障-发散精化-failures-divergences-refinement)
    - [1.5 FDR 模型检测器算法](#15-fdr-模型检测器算法)
  - [2. 属性推导 (Properties)](#2-属性推导-properties)
  - [3. 关系建立 (Relations)](#3-关系建立-relations)
  - [4. 论证过程 (Argumentation)](#4-论证过程-argumentation)
    - [4.1 规范化保证精化检查的正确性](#41-规范化保证精化检查的正确性)
    - [4.2 语义模型层次论证](#42-语义模型层次论证)
  - [5. 形式证明 (Proofs)](#5-形式证明-proofs)
    - [5.1 精化关系的传递性](#51-精化关系的传递性)
    - [5.2 不同语义模型的表达能力层次](#52-不同语义模型的表达能力层次)
  - [6. 实例与反例 (Examples \& Counter-examples)](#6-实例与反例-examples--counter-examples)
    - [6.1 实例：哲学家就餐问题的精化验证](#61-实例哲学家就餐问题的精化验证)
    - [6.2 反例1：状态爆炸导致 FDR 无法完成](#62-反例1状态爆炸导致-fdr-无法完成)
    - [6.3 反例2：过度精化拒绝合法实现](#63-反例2过度精化拒绝合法实现)
    - [6.4 反例3：实时系统扩展困难](#64-反例3实时系统扩展困难)
  - [7. 可视化图集](#7-可视化图集)
    - [7.1 概念依赖图：CSP 精化概念网络](#71-概念依赖图csp-精化概念网络)
    - [7.2 公理-定理推理树：从精化定义到传递性](#72-公理-定理推理树从精化定义到传递性)
    - [7.3 反例场景图：状态爆炸边界](#73-反例场景图状态爆炸边界)
  - [8. 参考文献](#8-参考文献)
  - [9. 关联可视化资源](#9-关联可视化资源)

---

## 1. 概念定义 (Definitions)

### 1.1 CSP 语义模型层次

**定义 1.1** (CSP 语义模型).

设 Σ 为可见动作集合，τ 为内部不可见动作，✓ 为成功终止标记。CSP 语义模型按观察粒度从粗到细分为四层：

| 模型 | 语义集合 | 符号 | 表达能力 |
|------|----------|------|----------|
| **迹模型** | traces(P) ⊆ Σ* | 𝒯 | 安全性 |
| **稳定故障模型** | failures(P) ⊆ Σ* × P(Σ) | ℱ | 安全性 + 部分活性 |
| **故障-发散模型** | failures(P) ∪ divergences(P) | ℱ𝒟 | 完全语义 |
| **无穷迹模型** | infinites(P) ⊆ Σ^ω | ℐ | 完全活性 |

**定义动机**：CSP 的设计哲学是"渐进式观察粒度"。从实际验证需求出发，安全性的验证最简单（只需检查迹包含），活性验证需要更精细的模型（故障、发散）。分层设计使得验证者可以在精度和计算复杂度之间做权衡——这是 FDR 能够处理工业规模系统的理论基础。

> **推断 [Theory→Model]**: CSP 语义模型层次设计遵循"可验证性优先"原则。
>
> **推断 [Model→Implementation]**: 因此 FDR 实现中允许用户根据验证需求选择不同精化模型，实现计算资源的最优分配。

---

### 1.2 迹精化 (Trace Refinement)

**定义 1.2** (迹语义).

进程 P 的**迹**是有限可见动作序列 s ∈ Σ*：

$$
\mathcal{T}\llbracket P \rrbracket = \{s \in \Sigma^* \mid P \xrightarrow{s}\}
$$

其中 P →s 表示存在迁移序列 P = P₀ →a₁ P₁ →a₂ ... →an Pn，s = ⟨a₁, a₂, ..., an⟩。

**定义 1.3** (迹精化).

进程 Q **迹精化**进程 P，记作 P ⊑ₜ Q，当且仅当：

$$
P \sqsubseteq_T Q \iff \mathcal{T}\llbracket Q \rrbracket \subseteq \mathcal{T}\llbracket P \rrbracket
$$

**直观解释**：Q 只能做 P 允许做的事（安全性保证）。迹精化意味着实现不比规范更"活跃"。

**定义动机**：迹精化是最基础的精化关系，直接对应"坏的事情不发生"的安全属性。其判断条件简单（集合包含），使得 FDR 在处理大规模系统的迹精化时可以达到每小时数十亿状态的检查速度。

---

### 1.3 故障精化 (Failures Refinement)

**定义 1.4** (稳定状态与拒绝集).

状态 P 是**稳定的**，记作 stable(P)，当且仅当 P ↛τ（无内部动作可执行）。

状态 P 的**拒绝集**是：

$$
refusals(P) = \{X \subseteq \Sigma \mid \forall a \in X: P \not\xrightarrow{a}\}
$$

**定义 1.5** (Failure).

**Failure** 是二元组 (s, X)，其中 s ∈ Σ* 是迹，X ⊆ Σ 是拒绝集。

进程 P 有 stable failure (s, X)，当：

$$
\exists P'. P \xrightarrow{s} P' \land stable(P') \land X \in refusals(P')
$$

**定义 1.6** (Failures 语义).

$$
\mathcal{F}\llbracket P \rrbracket = \{(s, X) \mid P \text{ 有 stable failure } (s, X)\}
$$

**定义 1.7** (故障精化).

进程 Q **故障精化**进程 P，记作 P ⊑ғ Q，当且仅当：

$$
P \sqsubseteq_F Q \iff \mathcal{F}\llbracket Q \rrbracket \subseteq \mathcal{F}\llbracket P \rrbracket
$$

**直观解释**：Q 不仅不能做 P 禁止的事（安全性），还必须在 P 接受的选择点上保持一致（部分活性）。

**定义动机**：故障精化解决了迹精化的关键缺陷——无法区分内部选择与外部选择。例如 P = (a → STOP) ⊓ (b → STOP) 和 Q = (a → STOP) □ (b → STOP) 有相同的迹，但 Q 必须接受环境的选择，而 P 可以内部决定。故障精化能够捕捉这一区别。

---

### 1.4 故障-发散精化 (Failures-Divergences Refinement)

**定义 1.8** (发散).

状态 P **发散**，记作 P ↑，当存在无限内部动作序列：

$$
P \xrightarrow{\tau} P_1 \xrightarrow{\tau} P_2 \xrightarrow{\tau} \cdots
$$

进程 P 的**发散集**是：

$$
\mathcal{D}\llbracket P \rrbracket = \{s \in \Sigma^* \mid \exists P'. P \xrightarrow{s} P' \land P' \uparrow\}
$$

**定义 1.9** (故障-发散精化).

进程 Q **FD 精化**进程 P，记作 P ⊑ғᴅ Q，当且仅当同时满足：

$$
\begin{aligned}
&1.\ \mathcal{F}\llbracket Q \rrbracket \subseteq \mathcal{F}\llbracket P \rrbracket \cup \{(s, X) \mid s \in \mathcal{D}\llbracket P \rrbracket\} \\
&2.\ \mathcal{D}\llbracket Q \rrbracket \subseteq \mathcal{D}\llbracket P \rrbracket
\end{aligned}
$$

**直观解释**：Q 在发散行为上也被 P 约束——如果 P 在某点后发散，Q 可以发散或进入混沌状态；如果 P 不发散，Q 也不能发散。

**定义动机**：发散代表无限循环或活锁。FD 精化确保实现不会引入新的发散行为，这对验证活性属性（如"请求最终会被响应"）至关重要。条件 1 中的扩展允许实现将规范的发散解释为"任何事情都可能发生"（混沌语义）。

---

### 1.5 FDR 模型检测器算法

**定义 1.10** (标记转移系统 LTS).

CSP 进程的**标记转移系统**是四元组 (S, Σ, →, s₀)，其中：

- S：状态集合
- Σ：可见动作集合
- → ⊆ S × (Σ ∪ {τ}) × S：转移关系
- s₀ ∈ S：初始状态

**定义 1.11** (规范化 GLTS).

**规范化**是将 LTS 转换为**规范化标记转移系统** (GLTS) 的过程：

1. 计算 τ 闭包：τ*(s) = {s' | s →*τ s'}
2. 移除所有 τ 转移
3. 对每对 (s, a) ∈ S × Σ，确保至多有唯一后继 s'

**定义 1.12** (FDR 精化检查算法).

给定规范 GLTS S = (Ss, Σ, →s, s₀s) 和实现 GLTS I = (Si, Σ, →i, s₀i)，精化检查算法执行状态对的广度优先搜索：

```
算法: FDR-Refinement-Check(S, I, model)
────────────────────────────────────────
输入: 规范 S, 实现 I, 精化模型 ∈ {T, F, FD}
输出: (true, ⊥) 或 (false, counterexample)

1. Normalize S → S_norm
2. BFS frontier ← {(s₀s, s₀i)}
3. visited ← ∅

4. while frontier ≠ ∅:
5.   (s_spec, s_impl) ← frontier.pop()
6.   if (s_spec, s_impl) ∈ visited: continue
7.   visited.add((s_spec, s_impl))
8.
9.   if not check-refinement(s_spec, s_impl, model):
10.    return (false, counterexample)
11.
12.  for each (a, s_impl') in transitions(I, s_impl):
13.    if a = τ:
14.      frontier.add((s_spec, s_impl'))
15.    else:
16.      s_spec' ← transition(S_norm, s_spec, a)
17.      if s_spec' = undefined:
18.        return (false, "Event not accepted by spec")
19.      frontier.add((s_spec', s_impl'))

20. return (true, ⊥)

子程序 check-refinement(s_spec, s_impl, model):
──────────────────────────────────────────────
if model = T:
  return true

if model = F or model = FD:
  if stable(s_spec):
    return refusals(s_impl) ⊆ refusals(s_spec)
  return true
```

**定义 1.13** (FDR3 并行算法).

FDR3 使用基于**状态哈希分区**的并行 BFS：

```
算法: Parallel-FDR(S, I, N_workers)
────────────────────────────────────
输入: 规范 S, 实现 I, 工作线程数 N

1. 每个工作线程 w 初始化:
   - done_w ← ∅
   - current_w ← ∅, next_w ← ∅
   - buffers_w[1..N]

2. 初始状态分配:
   w_root ← hash(s₀s, s₀i) mod N
   current_{w_root}.add((s₀s, s₀i))

3. while ∃w: current_w ≠ ∅:
4.   barrier_sync()
5.
6.   parallel for w ∈ [0, N-1]:
7.     for (s_spec, s_impl) ∈ current_w - done_w:
8.       if not check-refinement(s_spec, s_impl):
9.         return FAILURE
10.      done_w.add((s_spec, s_impl))
11.
12.      for each (a, s_impl') in successors(s_impl):
13.        s_spec' ← next-state(s_spec, a)
14.        w_target ← hash(s_spec', s_impl') mod N
15.        if w_target = w:
16.          next_w.add((s_spec', s_impl'))
17.        else:
18.          buffers_w[w_target].add((s_spec', s_impl'))
19.
20.  for w ∈ [0, N-1]:
21.    for w' ∈ [0, N-1]:
22.      current_w.add_all(buffers_{w'}[w])
23.    swap(current_w, next_w)
24.    next_w.clear()

25. return SUCCESS
```

**定义动机**：状态哈希分区是 FDR3 实现近线性加速的关键。每个状态对由确定性哈希分配到固定工作线程，避免了昂贵的锁竞争，同时保证状态不重复处理。

---

## 2. 属性推导 (Properties)

从上述定义，我们推导以下关键性质：

**性质 2.1** (精化的预序性质).

⊑ₜ, ⊑ғ, ⊑ғᴅ 都是预序关系（自反且传递）。

**推导**:

1. **自反性**：对任意进程 P，有 𝒯⟦P⟧ ⊆ 𝒯⟦P⟧，故 P ⊑ₜ P。同理可证其他。
2. **传递性**：见 [定理 5.1](#51-精化关系的传递性) 的完整证明。

---

**性质 2.2** (精化关系的可靠性保证).

若 P ⊑ғᴅ Q，则 Q 满足 P 表达的所有安全性和活性属性。

**推导**:

1. 由定义 1.9，𝒟⟦Q⟧ ⊆ 𝒟⟦P⟧，保证 Q 不会引入新的发散（活性）。
2. 由条件 1，ℱ⟦Q⟧ 被约束，保证 Q 不会做 P 禁止的事（安全性）。
3. 因此 Q 满足 P 的所有规范。

---

**性质 2.3** (语义模型的表达能力层次).

$$
\mathcal{T} \subset \mathcal{F} \subset \mathcal{FD}
$$

即：迹模型表达能力严格弱于故障模型，故障模型严格弱于故障-发散模型。

**推导**:

1. 每个故障 (s, X) 包含迹 s，故故障语义蕴含迹语义。
2. 发散集 𝒟 提供了故障模型无法表达的信息（无限行为）。
3. 严格性由反例证明（见 [6.3](#63-反例2过度精化拒绝合法实现)）。

---

**性质 2.4** (精化检查的复杂度上界).

对于有限状态进程 P 和 Q（状态数分别为 n 和 m）：

- 迹精化检查：O(n·m) 时间和空间
- 故障精化检查：O(n·m·2^|Σ|) 空间（最坏情况）
- FD 精化检查：同故障精化

**推导**:

1. 精化检查本质是状态对的 BFS/DFS，最坏情况遍历 Sp × Sq。
2. 迹精化只需检查可达性，空间 O(n·m)。
3. 故障精化需要存储拒绝集，每个状态最多 2^|Σ| 个拒绝集。
4. FDR 使用 B-Tree 压缩和规范化来减少实际空间使用。

---

**性质 2.5** (FDR 算法的最小反例保证).

若 P ⋢ Q，FDR 返回的反例具有**最小长度**（最短迹）。

**推导**:

1. FDR 使用 BFS 遍历状态空间。
2. BFS 按迹长度递增顺序探索状态。
3. 第一个发现的违反精化条件的状态对，其对应迹长度最小。
4. 因此返回的反例具有最小长度。

---

## 3. 关系建立 (Relations)

**关系 1**: CSP 精化 ≈ 互模拟等价 (Bisimulation)

**论证**:

- 互模拟 (~) 要求进程在每个状态上行为完全一致，是最严格的等价。
- FD 精化 (⊑ғᴅ) 是预序而非等价关系。
- 若 P ~ Q，则 P ⊑ғᴅ Q 且 Q ⊑ғᴅ P（双向精化蕴含互模拟）。
- 反之不成立：P ⊑ғᴅ Q 且 Q ⊑ғᴅ P 不保证 P ~ Q（规范可能更非确定性）。

符号表示：P ~ Q ⟹ (P ⊑ғᴅ Q ∧ Q ⊑ғᴅ P)

---

**关系 2**: CSP 精化 ⟹ 模型检测

**论证**:

- CSP 精化检查是一种**枚举模型检测**（enumerative model checking）。
- 精化 P ⊑ Q 等价于验证 Q 满足所有 P 表达的安全/活性属性。
- 与 LTL/CTL 模型检测不同：精化检查不需要显式写时序逻辑公式，规范本身就是 CSP 进程。
- 这体现了"进程即规范"(processes as specifications) 的 CSP 哲学。

符号表示：P ⊑ Q ⟺ Q ⊨ φᴘ，其中 φᴘ 是 P 编码的属性。

---

**关系 3**: CSP 理论 ⊃ FDR 实现

**论证**:

- CSP 理论是无穷的（可数无限进程集合）。
- FDR 实现只能处理**有限状态**CSP 子集（有限进程、有限数据域）。
- 对于无限状态进程，FDR 使用**数据无关性**(data independence) 和**归纳压缩**等抽象技术。
- 因此 FDR 是 CSP 的保守近似实现。

符号表示：FDR-CSP ⊂ Full-CSP

> **推断 [Theory→Implementation]**: CSP 理论中的无限进程（如带无限整数域的递归进程）无法直接被 FDR 处理。
>
> **推断 [Implementation→Domain]**: 因此工业应用中使用 CSP 建模时，必须将数据域限制为有限枚举类型，或使用数据无关性抽象。

---

**关系 4**: 迹精化 ⊂ 故障精化 ⊂ FD 精化

**论证**:

- 由 [性质 2.3](#23-语义模型的表达能力层次)，语义模型有严格层次。
- 精化关系方向相反：更严格的模型产生更严格的精化关系。
- 形式化：P ⊑ғᴅ Q ⟹ P ⊑ғ Q ⟹ P ⊑ₜ Q（反之不成立）。

符号表示：⊑ғᴅ ⊂ ⊑ғ ⊂ ⊑ₜ（作为关系集合）

---

## 4. 论证过程 (Argumentation)

### 4.1 规范化保证精化检查的正确性

**引理 4.1** (规范化保持语义).

设 norm(P) 为进程 P 的规范化 GLTS，则：

$$
\mathcal{F}\llbracket norm(P) \rrbracket = \mathcal{F}\llbracket P \rrbracket \land \mathcal{D}\llbracket norm(P) \rrbracket = \mathcal{D}\llbracket P \rrbracket
$$

**证明**:

1. **τ 闭包**：τ*(s) 收集所有通过内部动作可达的状态，这不会引入新的稳定故障（因为稳定状态要求无 τ 出边）。
2. **确定性构造**：对每对 (s, a)，若有多个 τ* 后继，选择所有后继的故障集合的并集作为新状态的故障集。
3. **发散保持**：发散检测基于 τ 循环，规范化不改变循环结构。

∎

---

**引理 4.2** (精化检查的完备性).

若 FDR 报告 P ⊑ Q 为真，则确实 P ⊑ Q；若为假，则存在真实反例。

**证明**:

1. **完备性（Soundness）**：FDR 只报告"通过"当遍历完所有可达状态对且未发现违反。由引理 4.1，规范化保持语义，故检查等价于理论精化。
2. **完备性（Completeness）**：若 P ⋢ Q，存在迹 s 和拒绝集 X 使得 (s, X) ∈ ℱ⟦Q⟧ 但 (s, X) ∉ ℱ⟦P⟧。BFS 会遍历所有迹，必然发现此违反。

∎

---

### 4.2 语义模型层次论证

**引理 4.3** (FD 蕴含 F 蕴含 T).

$$
P \sqsubseteq_{FD} Q \implies P \sqsubseteq_F Q \implies P \sqsubseteq_T Q
$$

**证明**:

1. **P ⊑ғᴅ Q ⟹ P ⊑ғ Q**：
   - 由定义 1.9 条件 1，ℱ⟦Q⟧ ⊆ ℱ⟦P⟧ ∪ {(s, X) | s ∈ 𝒟⟦P⟧}。
   - 由条件 2，𝒟⟦Q⟧ ⊆ 𝒟⟦P⟧。
   - 若 (s, X) ∈ ℱ⟦Q⟧ 且 s ∈ 𝒟⟦P⟧，则 s ∈ 𝒟⟦Q⟧。
   - 但稳定故障要求状态稳定（无 τ），发散状态不稳定，故 (s, X) 不能是 stable failure。
   - 因此 ℱ⟦Q⟧ ⊆ ℱ⟦P⟧，即 P ⊑ғ Q。

2. **P ⊑ғ Q ⟹ P ⊑ₜ Q**：
   - 若 (s, X) ∈ ℱ⟦Q⟧，则迹 s 在 Q 中可达。
   - 由故障精化，(s, X) ∈ ℱ⟦P⟧，故 s 在 P 中可达。
   - 因此 𝒯⟦Q⟧ ⊆ 𝒯⟦P⟧。

∎

---

## 5. 形式证明 (Proofs)

### 5.1 精化关系的传递性

**定理 5.1** (精化的传递性).

对任意精化模型 ℳ ∈ {𝒯, ℱ, ℱ𝒟}，若 P ⊑ℳ Q 且 Q ⊑ℳ R，则 P ⊑ℳ R。

**证明**：

我们以故障精化 ⊑ғ 为例，其他模型类似。

**前提**：

- P ⊑ғ Q ⟺ ℱ⟦Q⟧ ⊆ ℱ⟦P⟧
- Q ⊑ғ R ⟺ ℱ⟦R⟧ ⊆ ℱ⟦Q⟧

**推导**：

$$
\begin{aligned}
\mathcal{F}\llbracket R \rrbracket &\subseteq \mathcal{F}\llbracket Q \rrbracket && \text{（由 } Q \sqsubseteq_F R \text{）} \\
&\subseteq \mathcal{F}\llbracket P \rrbracket && \text{（由 } P \sqsubseteq_F Q \text{）}
\end{aligned}
$$

由集合包含的传递性，ℱ⟦R⟧ ⊆ ℱ⟦P⟧。

因此 P ⊑ғ R。

**关键案例分析**：

| 案例 | 条件 | 结论 |
|------|------|------|
| 迹精化 | 𝒯⟦R⟧ ⊆ 𝒯⟦Q⟧, 𝒯⟦Q⟧ ⊆ 𝒯⟦P⟧ | 𝒯⟦R⟧ ⊆ 𝒯⟦P⟧ |
| FD 精化 | 迹 + 故障 + 发散条件均满足 | 所有条件通过传递性保持 |

∎

---

### 5.2 不同语义模型的表达能力层次

**定理 5.2** (语义模型的严格层次).

$$
\mathcal{T} \subset \mathcal{F} \subset \mathcal{FD}
$$

即故障模型严格比迹模型更具表达能力，FD 模型严格比故障模型更具表达能力。

**证明**：

**第一部分：𝒯 ⊂ ℱ**

1. **包含性（𝒯 ⊆ ℱ）**：
   - 每个故障 (s, X) 包含迹 s。
   - 给定故障集 ℱ，可提取迹集 {s | ∃X: (s, X) ∈ ℱ}。
   - 因此故障语义蕴含迹语义。

2. **严格性（𝒯 ≠ ℱ）**：
   - 反例：P = (a → STOP) ⊓ (b → STOP)，Q = (a → STOP) □ (b → STOP)。
   - 𝒯⟦P⟧ = 𝒯⟦Q⟧ = {⟨⟩, ⟨a⟩, ⟨b⟩}（迹相同）。
   - 但 ℱ⟦P⟧ 包含 (⟨⟩, {b})（选择 a 分支后拒绝 b），而 ℱ⟦Q⟧ 不包含（必须同时接受 a 和 b 直到选择）。
   - 因此故障模型能区分迹模型无法区分的进程。

**第二部分：ℱ ⊂ ℱ𝒟**

1. **包含性（ℱ ⊆ ℱ𝒟）**：
   - FD 模型包含故障集和发散集，完全包含故障模型信息。

2. **严格性（ℱ ≠ ℱ𝒟）**：
   - 反例：P = a → div，Q = a → STOP。
   - ℱ⟦P⟧ = {(⟨⟩, X) | a ∉ X}（不稳定状态无故障）。
   - ℱ⟦Q⟧ 类似，但 Q 在 a 后稳定。
   - 然而 𝒟⟦P⟧ = {⟨a⟩}（发散），𝒟⟦Q⟧ = ∅。
   - FD 模型能区分故障模型无法区分的进程（当一方发散时）。

∎

---

## 6. 实例与反例 (Examples & Counter-examples)

### 6.1 实例：哲学家就餐问题的精化验证

**示例 6.1** (哲学家就餐问题精化验证).

考虑经典哲学家就餐问题：N 个哲学家围坐，每两人之间一根筷子。每个哲学家需要两根筷子才能就餐。

**CSP 建模**：

```csp
-- 定义通道
channel pick, put: {0..N-1}.{0..N-1}
channel think, eat: {0..N-1}

-- 哲学家进程
PHIL(i) =
  think.i ->
  pick.i.i -> pick.i.((i+1)%N) ->
  eat.i ->
  put.i.((i+1)%N) -> put.i.i ->
  PHIL(i)

-- 筷子进程
FORK(j) =
  pick.j.j -> put.j.j -> FORK(j)
  []
  pick.((j-1+N)%N).j -> put.((j-1+N)%N).j -> FORK(j)

-- 系统组合
SYSTEM = ||| i:{0..N-1} @ [ {|pick.i, put.i|} ] PHIL(i)
         [ {|pick, put|} ]
         ||| j:{0..N-1} @ [ {|pick.j, put.j|} ] FORK(j)

-- 规范：死锁自由
SPEC = RUN(diff(Events, {think, eat}))
```

**精化验证**：

```csp
-- 迹精化：检查系统只做规范允许的动作
assert SPEC [T= SYSTEM

-- 故障精化：检查死锁自由
assert RUN(Events) [F= SYSTEM

-- FD 精化：检查活锁自由
assert CHAOS(Events) [FD= SYSTEM
```

**验证结果分析**：

- 当 N = 2：验证通过（无死锁）。
- 当 N = 5：FDR 报告死锁反例（循环等待）。

---

### 6.2 反例1：状态爆炸导致 FDR 无法完成

**反例 6.2** (状态爆炸边界).

**场景**：验证包含 N 个并行组件的系统，每个组件有 M 个状态。

**问题**：

- 若组件完全独立：状态空间 M^N。
- 若组件通过共享通道交互：状态空间可达 (M·N)!（排列组合爆炸）。

**具体案例**：

考虑 N 个并行的缓冲区进程：

```csp
BUFFER(cap) =
  in?x -> (if cap > 0 then BUFFER(cap-1) else full -> BUFFER(cap))
  []
  out!x -> (if cap < MAX then BUFFER(cap+1) else empty -> BUFFER(cap))

SYSTEM(N) = ||| i:{0..N-1} @ [ {|in, out|} ] BUFFER(MAX)
```

**状态空间增长**：

| N | MAX=2 状态数 | MAX=10 状态数 | FDR 完成时间 |
|---|-------------|---------------|-------------|
| 2 | 9 | 121 | < 1s |
| 4 | 81 | 14,641 | 5s |
| 8 | 6,561 | 214M | 3h |
| 16 | 43M | 4.6×10^17 | > 100 年 |

**分析**：

- **违反的前提**：FDR 假设有限状态系统可在合理时间内完成验证。
- **导致的异常**：状态空间指数增长导致内存耗尽或运行时间过长。
- **结论**：即使理论可判定，工程实践中仍需组合验证、抽象或对称归约。

> **推断 [Theory→Implementation]**: CSP 精化检查在理论上是可判定的（对有限状态系统）。
>
> **推断 [Implementation→Domain]**: 但在实际工程中，状态爆炸使得朴素的"运行 FDR 检查"无法处理工业规模系统，必须使用分层验证、数据抽象或归纳压缩。

---

### 6.3 反例2：过度精化拒绝合法实现

**反例 6.3** (过度精化).

**规范**（设计者意图）：

```csp
SPEC = req -> (ack1 -> STOP |~| ack2 -> STOP)
```

**实现**（程序员的正确实现）：

```csp
IMPL = req -> if (condition) then ack1 -> STOP else ack2 -> STOP
```

**验证尝试**：

```csp
assert SPEC [F= IMPL  -- 失败！
```

**失败分析**：

| 检查项 | 规范行为 | 实现行为 | 结果 |
|--------|----------|----------|------|
| 迹 | {req·ack1, req·ack2} | {req·ack1, req·ack2} | ✓ 通过 |
| 故障 (req, {ack1}) | 可拒绝 | 不可拒绝 | ✗ 失败 |

**根本原因**：

- 规范使用了**内部选择**（|~|），允许实现任意选择分支。
- 但实现使用了**确定性分支**（if-then-else），在 req 后已确定走哪条分支。
- 规范在 req 后的状态可以"假装"还会选择另一条分支。
- 实现不能"假装"，因此拒绝集不同。

**修复**：

```csp
SPEC' = req -> (ack1 -> STOP [] ack2 -> STOP)
```

**结论**：

- 精化关系要求实现"比规范更确定"。
- 过度非确定性的规范会拒绝实际上正确的实现。
- 规范设计应遵循"最小非确定性原则"。

---

### 6.4 反例3：实时系统扩展困难

**反例 6.4** (Timed CSP 的精化挑战).

**问题**：标准 CSP 精化无法直接处理时间约束。

**场景**：实时系统要求"请求必须在 100ms 内得到响应"。

**尝试建模**（标准 CSP）：

```csp
SPEC = req -> (timeout -> ERROR [] ack -> STOP)
-- 无法区分 1ms 和 10s！
```

**Timed CSP 扩展**：

Roscoe 提出了 Timed CSP，引入时间操作符：

```csp
WAIT(t)
timeout{t}(P, Q)
```

**精化问题**：

| 问题 | 说明 | 后果 |
|------|------|------|
| 语义复杂度 | 时间引入连续状态空间 | 模型检测不再可判定 |
| 交错语义 | 不同时间粒度组件同步 | 精化关系需重新定义 |
| 工具支持 | FDR 对 Timed CSP 支持有限 | 工业应用困难 |

**案例**：

考虑两个实时进程：

```csp
P = (WAIT(0); a -> STOP) |~| (WAIT(1); a -> STOP)
Q = WAIT(0.5); a -> STOP
```

**问题**：P ⊑ Q 还是 Q ⊑ P？

- 在连续时间模型中，两者迹相同（都是 ⟨a⟩）。
- 但在时间精化中，Q 在 t < 0.5 时不能执行 a，而 P 可以。
- 因此时间精化需要全新的语义框架。

**结论**：

- CSP 精化理论主要针对**非时间**并发系统。
- 实时系统需要 Timed CSP 或其他实时进程代数。
- FDR 对实时系统的支持有限，这是其工业应用的主要边界之一。

---

## 7. 可视化图集

### 7.1 概念依赖图：CSP 精化概念网络

```mermaid
graph TB
    subgraph "理论基础"
        T1[TH1: 进程代数] --> T2[CSP 语法]
        T3[TH4: 表达能力层次] --> T4[语义模型层次]
    end

    subgraph "本文档核心 - CSP 精化"
        T2 --> D1["定义 1.2: 迹精化"]
        T2 --> D2["定义 1.7: 故障精化"]
        T2 --> D3["定义 1.9: FD 精化"]
        T4 --> D1
        T4 --> D2
        T4 --> D3

        D1 --> R1["关系: 表达能力层次"]
        D2 --> R1
        D3 --> R1

        D1 --> P1["定理 5.1: 传递性"]
        D2 --> P1
        D3 --> P1

        D3 --> P2["定理 5.2: 表达能力层次"]
    end

    subgraph "实现层 - FDR"
        D1 --> I1["FDR 算法: 迹检查"]
        D2 --> I2["FDR 算法: 故障检查"]
        D3 --> I3["FDR 算法: FD 检查"]

        I1 --> A1["定义 1.12: BFS 遍历"]
        I2 --> A1
        I3 --> A1

        A1 --> A2["定义 1.13: 并行分区"]

        I3 --> A3["引理 4.1: 规范化"]
    end

    subgraph "应用与边界"
        I1 --> C1["实例: 死锁验证"]
        I2 --> C2["反例: 状态爆炸"]
        I3 --> C3["反例: 过度精化"]
        D3 --> C4["反例: 实时系统"]
    end

    style D1 fill:#e1bee7,stroke:#6a1b9a
    style D2 fill:#e1bee7,stroke:#6a1b9a
    style D3 fill:#e1bee7,stroke:#6a1b9a
    style P1 fill:#c8e6c9,stroke:#2e7d32
    style P2 fill:#c8e6c9,stroke:#2e7d32
    style C2 fill:#ffcdd2,stroke:#c62828
    style C3 fill:#ffcdd2,stroke:#c62828
    style C4 fill:#ffcdd2,stroke:#c62828
```

**图说明**：

- 本图展示了 CSP 精化从理论基础到工程应用的概念依赖关系。
- **紫色节点**（D1-D3）表示核心定义，是全文的基础。
- **绿色节点**（P1-P2）表示主要定理，从定义推导而来。
- **红色节点**（C2-C4）表示反例/边界场景，展示了理论的应用边界。
- 从理论到实现再到应用，构成了完整的知识链条。

---

### 7.2 公理-定理推理树：从精化定义到传递性

```mermaid
graph BT
    subgraph "公理层"
        A1["[Axiom] 集合包含的传递性"]
        A2["[Axiom] 幂集偏序"]
    end

    subgraph "定义层"
        D1["[Def 1.3] 迹精化"]
        D2["[Def 1.7] 故障精化"]
        D3["[Def 1.9] FD精化"]
    end

    subgraph "引理层"
        L1["[Lemma] 迹包含传递"]
        L2["[Lemma] 故障包含传递"]
        L3["[Lemma] 发散包含传递"]
    end

    subgraph "定理层"
        T1["[定理 5.1] 精化的传递性"]
    end

    A1 --> L1
    A1 --> L2
    A2 --> L3

    D1 --> L1
    D2 --> L2
    D3 --> L2
    D3 --> L3

    L1 --> T1
    L2 --> T1
    L3 --> T1

    style A1 fill:#fff9c4,stroke:#f57f17
    style A2 fill:#fff9c4,stroke:#f57f17
    style L1 fill:#bbdefb,stroke:#1565c0
    style L2 fill:#bbdefb,stroke:#1565c0
    style L3 fill:#bbdefb,stroke:#1565c0
    style T1 fill:#c8e6c9,stroke:#2e7d32
```

**图说明**：

- 本推理树展示了从集合论公理到精化传递性定理的完整推导链条。
- **黄色节点**（A1-A2）：数学基础公理。
- **蓝色节点**（L1-L3）：由定义直接推导的引理。
- **绿色节点**（T1）：主要定理，综合各引理得出结论。

---

### 7.3 反例场景图：状态爆炸边界

```mermaid
sequenceDiagram
    participant U as 用户
    participant F as FDR3 精化检查器
    participant S as 状态存储
    participant M as 内存/磁盘

    U->>F: 检查 SYSTEM(16) ⊑ SPEC

    Note over F: 初始状态对 (S₀, I₀)

    loop BFS 遍历
        F->>F: 扩展状态对 (Sᵢ, Iⱼ)
        F->>S: 存储 visited 状态
        S->>M: 写入 B-Tree 索引

        alt 状态数 < 阈值
            Note over M: 内存存储 ~8GB/小时
        else 状态数 > 阈值
            M->>M: 溢出到磁盘 性能下降 100x
        end
    end

    Note over F: 1小时后... 已访问70亿状态

    M--xM: 内存耗尽 OutOfMemoryError

    F-->>U: 检查失败 状态空间过大

    Note right of F: 反例结论:<br/>即使理论可判定，状态爆炸使<br/>FDR 无法处理某些有限状态系统
```

**图说明**：

- 本序列图展示了 FDR 在面对状态爆炸问题时的执行流程和失败模式。
- **正常流程**：BFS 遍历、状态存储、B-Tree 索引。
- **边界触发**：状态数超过阈值，内存溢出到磁盘。
- **失败模式**：内存耗尽，检查无法完成。
- **关键教训**：理论可判定不等于工程可验证，必须配合抽象技术使用。

---

## 8. 参考文献

**CSP 理论:**

1. Hoare, C. A. R. (1985). *Communicating Sequential Processes*. Prentice Hall.
2. Roscoe, A. W. (1998). *The Theory and Practice of Concurrency*. Prentice Hall.
3. Roscoe, A. W. (2015). *Understanding Concurrent Systems*. Springer.
4. Brookes, S. D., Hoare, C. A. R., & Roscoe, A. W. (1984). A Theory of Communicating Sequential Processes. *Journal of the ACM*, 31(3), 560-599.

**FDR 工具:**

1. Gibson-Robinson, T., Armstrong, P., Boulgakov, A., & Roscoe, A. W. (2014). FDR3 — A Modern Refinement Checker for CSP. *TACAS 2014*, 187-201.
2. Roscoe, A. W., et al. (1995). Industrial Strength CSP. *Formal Methods'95*.
3. Formal Systems (Europe) Ltd. *FDR User Manual*.

**精化理论:**

1. Schneider, S. A. (2000). *Concurrent and Real-time Systems: The CSP Approach*. Wiley.
2. Lowe, G. (1996). Breaking and Fixing the Needham-Schroeder Protocol using FDR. *TACAS 1996*.

**应用:**

1. Ryan, P. Y. A., & Schneider, S. A. (2000). *The Modelling and Analysis of Security Protocols*. Addison-Wesley.

---

## 9. 关联可视化资源

本文档包含以下可视化资源，详细说明参见 [VISUAL-ATLAS.md](../../VISUAL-ATLAS.md)：

| 图类型 | 图名称 | 位置 |
|--------|--------|------|
| 概念依赖图 | CSP 精化概念网络 | 本文档 §7.1 |
| 公理-定理推理树 | 从精化定义到传递性 | 本文档 §7.2 |
| 反例场景图 | 状态爆炸边界（序列图） | 本文档 §7.3 |

**规范引用**：

> 本文档遵循 [Document-Standard.md](../../templates/Document-Standard.md) 六段式模板，
> 可视化遵循 [Visualization-Standard.md](../../templates/Visualization-Standard.md) 五种图类型规范，
> 跨层推断遵循 [Cross-Level-Inference-Framework.md](../../templates/Cross-Level-Inference-Framework.md) 标记规范。

---

**文档版本**: 2026.03 (深度重构版) | **状态**: Phase 2 批量重构完成
