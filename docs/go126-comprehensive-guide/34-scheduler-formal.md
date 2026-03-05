# Goroutine调度器形式化分析

> Go G-M-P调度器的数学模型与性能分析

---

## 一、调度器状态空间

### 1.1 状态定义

```text
定义 1.1 (调度器状态):
────────────────────────────────────────
调度器状态 S = (G, M, P, Q, B)

其中：
├─ G：Goroutine集合
├─ M：Machine(OS线程)集合
├─ P：Processor集合
├─ Q：运行队列集合
└─ B：阻塞集合

Goroutine状态：
G ∈ Goroutine → State
State ::= Runnable | Running | Waiting | Dead

状态转移函数：
δ : S × Event → S

定义 1.2 (全局状态不变式):
────────────────────────────────────────
对于任何有效状态S，必须满足：

1. 每个Running的G必绑定到一个M：
   ∀g ∈ G, State(g) = Running ⟹ ∃!m ∈ M, Bound(g, m)

2. 每个M至多绑定一个P：
   ∀m ∈ M, |{p ∈ P | Bound(m, p)}| ≤ 1

3. 每个P至多绑定一个M：
   ∀p ∈ P, |{m ∈ M | Bound(p, m)}| ≤ 1

4. Runnable的G必在某个队列中：
   ∀g ∈ G, State(g) = Runnable ⟹
   ∃p ∈ P, g ∈ LocalQueue(p) ∨ g ∈ GlobalQueue

5. Waiting的G不在任何队列：
   ∀g ∈ G, State(g) = Waiting ⟹
   ¬(∃p ∈ P, g ∈ LocalQueue(p)) ∧ g ∉ GlobalQueue
```

### 1.2 状态转移规则

```text
规则 1.1 (创建Goroutine):
────────────────────────────────────────
前置条件：go f() 语句执行

转移：
S = (G, M, P, Q, B)
→
S' = (G ∪ {g'}, M, P, Q ∪ {g'}, B)

其中g'是新创建的goroutine，初始状态为Runnable

新G的放置策略：
1. 尝试放入当前P的本地队列
2. 若本地队列满，放入全局队列
3. 若全局队列满，放入其他P的队列(工作窃取准备)

形式化：
Let p = CurrentP()
If |LocalQueue(p)| < MaxLocalQueueSize Then
    LocalQueue(p) = LocalQueue(p) ∪ {g'}
Else
    GlobalQueue = GlobalQueue ∪ {g'}

规则 1.2 (调度Goroutine):
────────────────────────────────────────
前置条件：M需要获取G执行

转移：
S = (G, M, P, Q, B)
→
S' = (G, M, P, Q', B)

其中Q' = Q \ {g}，State(g)从Runnable变为Running

调度顺序：
1. 检查runnext (M绑定的P的下一个G)
2. 从P的本地队列获取
3. 从全局队列获取 (每61次检查一次)
4. 从其他P窃取 (工作窃取)
5. 从网络轮询器获取

形式化：
SelectG() =
    if runnext ≠ nil then runnext
    else if LocalQueue ≠ ∅ then Dequeue(LocalQueue)
    else if GlobalQueue ≠ ∅ then GlobalSteal()
    else if StealableFromOtherP() then WorkSteal()
    else PollNetwork()

规则 1.3 (阻塞Goroutine):
────────────────────────────────────────
前置条件：G执行阻塞操作(channel、锁、系统调用)

转移：
S = (G, M, P, Q, B)
→
S' = (G, M', P, Q, B ∪ {(g, reason)})

其中State(g)变为Waiting

对于系统调用：
M释放P，进入系统调用状态
M' = M 但状态改变

对于channel操作：
G被加入channel的等待队列
M继续执行其他G

规则 1.4 (唤醒Goroutine):
────────────────────────────────────────
前置条件：阻塞的G可以被唤醒

转移：
S = (G, M, P, Q, B)
→
S' = (G, M, P, Q ∪ {g}, B \
 {(g, _)})

State(g)从Waiting变为Runnable

G被放入某个P的本地队列或全局队列
```

---

## 二、工作窃取算法形式化

### 2.1 窃取策略

```text
定义 2.1 (工作窃取):
────────────────────────────────────────
当P的本地队列为空时，从其他P窃取G

窃取目标选择：
StealTarget(p₀) = argmax_{p ∈ P, p ≠ p₀} |LocalQueue(p)|

窃取数量：
窃取目标P队列的一半：
StealCount(p) = ⌈|LocalQueue(p)| / 2⌉

窃取的原子性：
窃取操作必须是原子的，防止竞争

定义 2.2 (窃取成功率):
────────────────────────────────────────
设N为P的数量，L为平均队列长度

单次窃取成功概率：
P(steal_success) = 1 - (1 - 1/N)^L

当N较大或L较小时，成功率下降

定理 2.1 (负载均衡):
────────────────────────────────────────
工作窃取算法在O(log N)轮后达到近似负载均衡

证明概要：
- 每次窃取减少最大队列与最小队列的差距
- 通过势函数分析，每轮窃取减少总方差
- 经过O(log N)轮，方差降到常数级别∎

定理 2.2 (窃取开销):
────────────────────────────────────────
每次窃取需要O(1)次CAS操作

由于窃取是竞争操作，可能需要重试
预期重试次数为O(1) (低竞争情况下)
```

### 2.2 调度延迟分析

```text
定义 2.3 (调度延迟):
────────────────────────────────────────
Goroutine g的调度延迟：
T_schedule(g) = t_start(g) - t_runnable(g)

其中：
├─ t_runnable(g)：g变为Runnable的时间
└─ t_start(g)：g开始执行的时间

延迟组成：
T_schedule = T_queue + T_steal + T_context_switch

其中：
├─ T_queue：在队列中等待时间
├─ T_steal：工作窃取寻找时间
└─ T_context_switch：上下文切换时间

定理 2.3 (延迟上界):
────────────────────────────────────────
在N个P的系统中，Goroutine的最大调度延迟为：
T_max = O(N × T_quantum)

其中T_quantum是时间片长度

证明：
最坏情况下，Goroutine需要等待：
- 所有其他P上的Goroutine执行一轮
- 通过工作窃取被发现
- 总共O(N)个时间片∎

定义 2.4 (公平性):
────────────────────────────────────────
调度器的公平性度量：
Fairness = min_g (CPU_time(g) / Expected_time(g))

其中Expected_time与Goroutine的优先级和创建时间相关

Go调度器使用先进先出(FIFO)保证公平性
```

---

## 三、性能模型

### 3.1 吞吐量分析

```text
定义 3.1 (系统吞吐量):
────────────────────────────────────────
Throughput = 单位时间内完成的Goroutine数

影响因素：
├─ N：P的数量
├─ T_ctx：上下文切换开销
├─ T_work：平均工作负载
├─ B：阻塞概率
└─ S：同步开销

吞吐量公式：
Throughput = N / (T_work + T_ctx + B × T_block + S)

其中T_block是阻塞等待时间

定理 3.1 (最优P数量):
────────────────────────────────────────
设C为CPU核心数，I为IO密集型比例

最优P数量：
P_optimal = C / (1 - I)

当I = 0 (纯CPU密集型)：P_optimal = C
当I > 0 (有IO)：P_optimal > C

证明：
CPU密集型任务不应超过C，否则导致上下文切换浪费
IO密集型任务可以超过C，因为阻塞时可以让出CPU
```

### 3.2 扩展性分析

```text
定义 3.2 (强扩展性):
────────────────────────────────────────
Speedup(N) = T(1) / T(N)

其中T(N)是N个P时的执行时间

理想线性扩展：Speedup(N) = N

实际扩展受限于：
- Amdahl定律：串行部分限制
- 同步开销：锁、原子操作
- 负载不均衡：工作窃取延迟

定理 3.2 (Amdahl界限):
────────────────────────────────────────
设f为必须串行执行的比例

最大加速比：
Speedup_max(N) = 1 / (f + (1-f)/N)

当N → ∞时，Speedup_max → 1/f

示例：
若f = 5%，则最大加速比为20倍
无论使用多少核心，无法超过20倍加速

定义 3.3 (弱扩展性):
────────────────────────────────────────
Efficiency(N) = T(1) / (N × T(N))

保持每个P的工作负载不变，增加问题规模

Go调度器弱扩展性良好，适合大规模并发
```

---

## 四、形式化验证

### 4.1 安全性证明

```text
定理 4.1 (调度安全性):
────────────────────────────────────────
调度器保证：
1. 不会丢失Goroutine
2. 不会重复执行Goroutine
3. 不会饿死Goroutine

证明：
1. 每个创建的Goroutine都被放入队列
   要么执行，要么在队列中等待
   队列操作是完备的

2. 一个Goroutine只能从队列中取出一次
   取出后变为Running，不会再次入队
   直到变为Runnable才可能重新入队

3. 公平性保证：
   - 本地队列是FIFO
   - 全局队列定期被检查
   - 工作窃取最终会发现所有G
   - 因此任何G最终都会被执行∎

定理 4.2 (互斥正确性):
────────────────────────────────────────
若Goroutine g持有锁m，则没有其他G持有m

证明：
锁m的状态由原子操作维护
Acquire(m)成功当且仅当m未被持有
Release(m)将m状态设为未持有
原子性保证了互斥∎
```

### 4.2 活性证明

```text
定理 4.3 (调度活性):
────────────────────────────────────────
若存在Runnable的Goroutine，则最终有Goroutine被调度

证明：
情况1：存在空闲的M-P对
  直接调度Goroutine

情况2：所有M都在Running
  某个M的Goroutine会完成或阻塞
  释放的M会调度新的G

情况3：需要工作窃取
  空闲的P会尝试窃取
  若窃取成功，调度继续
  若失败，说明系统无工作
  与前提矛盾∎

定理 4.4 (无死锁):
────────────────────────────────────────
在只使用Go同步原语(channel、mutex、select)的程序中，
调度器不会导致死锁，除非程序逻辑本身有死锁

注意：
这并不意味着Go程序不会死锁
程序逻辑错误(如循环等待)仍会导致死锁
但调度器本身不会引入额外的死锁
```

---

*本章建立了Go调度器的形式化模型，包括状态空间、工作窃取算法和性能分析。*
