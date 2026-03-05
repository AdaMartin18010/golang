# 垃圾回收器形式化模型

> Go GC算法的数学模型与性能分析

---

## 一、GC状态机模型

### 1.1 GC生命周期

```text
GC状态定义:
────────────────────────────────────────

GCState ::=
  | OFF         (GC关闭)
  | MARK        (标记阶段)
  | MARK_TERM   (标记终止)
  | SWEEP       (清扫阶段)

状态转移：
OFF → MARK → MARK_TERM → SWEEP → OFF

触发条件：
- OFF → MARK：堆大小达到阈值
- MARK → MARK_TERM：标记队列空
- MARK_TERM → SWEEP：STW完成
- SWEEP → OFF：清扫完成

定义 1.1 (GC周期):
────────────────────────────────────────
GC周期 C = (T_start, T_end, H_before, H_after, LiveSet)

其中：
├─ T_start：GC开始时间
├─ T_end：GC结束时间
├─ H_before：GC前堆大小
├─ H_after：GC后堆大小
└─ LiveSet：存活对象集合

GC暂停时间：
PauseTime = T_STW_end - T_STW_start

GC吞吐量：
GCThroughput = (AllocatedBytes - FreedBytes) / (T_end - T_start)
```

### 1.2 三色不变式

```text
定义 1.2 (三色抽象):
────────────────────────────────────────

对象颜色：
Color ::= WHITE | GRAY | BLACK

WHITE：未访问，候选垃圾
GRAY：已访问，但引用未完全扫描
BLACK：已完全扫描，存活

三色不变式 (Tri-color Invariant):
────────────────────────────────────────
没有BLACK对象直接引用WHITE对象

形式化：
∀o₁, o₂ ∈ Objects,
Color(o₁) = BLACK ∧ Reference(o₁, o₂)
⟹ Color(o₂) ∈ {GRAY, BLACK}

不变式维护：
当BLACK对象获得WHITE引用时，
必须将WHITE对象变为GRAY或BLACK

写屏障 (Write Barrier):
────────────────────────────────────────
在并发标记期间，对象引用修改时：

原引用：obj.field = old_val
新引用：obj.field = new_val

写屏障操作：
1. 若obj是BLACK且new_val是WHITE
2. 将new_val着色为GRAY
3. 加入标记队列

或采用Dijkstra式写屏障：
1. 若obj是BLACK
2. 将obj重新着色为GRAY
3. 重新扫描obj

定理 1.1 (写屏障正确性):
────────────────────────────────────────
写屏障保证三色不变式在并发标记期间始终成立

证明：
情况1：WHITE → WHITE
  不破坏不变式

情况2：WHITE/GRAY/BLACK → GRAY/BLACK
  目标对象不再是WHITE，安全

情况3：BLACK → WHITE (危险)
  写屏障拦截，将WHITE转为GRAY
  或把BLACK转为GRAY重新扫描
  不变式恢复∎
```

---

## 二、标记算法形式化

### 2.1 并发标记

```text
定义 2.1 (标记工作列表):
────────────────────────────────────────
WorkList：待扫描的GRAY对象队列

标记算法：
while WorkList ≠ ∅:
    obj = WorkList.pop()
    for ref in obj.references:
        if ref.color == WHITE:
            ref.color = GRAY
            WorkList.push(ref)
    obj.color = BLACK

并发修改：
- 多标记线程并发执行
- 需要同步访问WorkList
- 使用无锁队列或工作窃取

定义 2.2 (标记终止):
────────────────────────────────────────
全局标记终止条件：
∀obj ∈ Objects, obj.color ≠ GRAY

即：WorkList为空且所有P的本地工作列表为空

终止检测算法：
1. 停止所有标记线程
2. 检查全局WorkList
3. 检查每个P的本地队列
4. 若全部为空，标记完成

定理 2.1 (标记完备性):
────────────────────────────────────────
若从GC Roots可达的对象，必被标记为BLACK

证明：
- GC Roots初始为GRAY，加入WorkList
- 归纳：若对象o在WorkList中，必被处理
- 处理时将所有WHITE子对象加入WorkList
- 可达性传递，所有可达对象都被处理∎

定理 2.2 (标记正确性):
────────────────────────────────────────
若对象被标记为BLACK，则它从GC Roots可达

证明：
- 只有GRAY对象能变为BLACK
- 只有GC Roots或被BLACK引用的对象变为GRAY
- 归纳可得BLACK对象都可达∎
```

### 2.2 增量与并发

```text
定义 2.3 (增量标记):
────────────────────────────────────────
将标记工作分散到多次小步骤中执行

用户程序与GC交替执行：
User → GC(Δt) → User → GC(Δt) → ...

目标：
每次GC暂停时间 < T_max (通常1ms)

定义 2.4 (债务模型):
────────────────────────────────────────
GC债务：程序分配内存的速率

债务比率 = AllocationRate / GCScanRate

若债务比率 > 1，标记跟不上分配，需要辅助标记或STW

辅助标记 (Mutator Assistance):
────────────────────────────────────────
当分配速度快于标记速度时，
用户程序(Mutator)参与标记工作

形式化：
每分配N字节，执行M字节的标记工作

比例调整：
AssistRatio = (HeapSize - LiveSize) / LiveSize

定理 2.3 (标记完成保证):
────────────────────────────────────────
使用辅助标记，GC必能在堆耗尽前完成标记

证明：
- 设分配速率为A
- 标记速率为S (包括专用标记线程和辅助标记)
- 需要标记的对象数为L (LiveSize)
- 标记完成时间：T = L / S
- 期间分配量：A × T = A × L / S
- 辅助标记保证S > A × L / (HeapSize - L)
- 故分配量 < HeapSize - L，不会溢出∎
```

---

## 三、清扫与内存管理

### 3.1 清扫算法

```text
定义 3.1 (清扫阶段):
────────────────────────────────────────
回收WHITE对象的内存

延迟清扫 (Lazy Sweeping):
────────────────────────────────────────
不在GC结束时立即清扫全部
而是在分配时按需清扫

清扫单位：mspan (内存页组)
每个mspan有清扫状态

分配时检查：
if mspan.swept == false:
    sweep(mspan)
    mspan.swept = true

分配从mspan中查找空闲对象

定义 3.2 (清扫粒度):
────────────────────────────────────────
Granularity = SpanSize / ObjectSize

细粒度：对象小，清扫开销分散
粗粒度：对象大，清扫效率高

定理 3.1 (清扫延迟摊销):
────────────────────────────────────────
延迟清扫将清扫延迟摊销到多次分配中

每次分配的额外开销：O(1)

证明：
- 设总清扫工作为W
- 设分配次数为N
- 每次分配触发W/N的清扫工作
- 常数时间完成∎
```

### 3.2 内存压缩

```text
定义 3.3 (堆碎片):
────────────────────────────────────────
外部碎片：
自由内存总量足够，但没有连续大块

内部碎片：
分配的内存大于实际需要

碎片率：
Fragmentation = 1 - (LargestFreeBlock / TotalFreeSpace)

Go不压缩：
────────────────────────────────────────
Go GC不进行内存压缩(对象不移动)

原因：
1. 简化实现
2. 避免指针更新开销
3. 使用TCMalloc式分配器缓解碎片

但代价：
- 长期运行的程序可能积累碎片
- 内存使用可能高于实际需要
```

---

## 四、GC调优理论

### 4.1 GOGC参数模型

```text
定义 4.1 (GOGC):
────────────────────────────────────────
GOGC = g，表示目标堆大小为存活对象的(g+100)%

目标堆大小：
TargetHeap = LiveSize × (1 + g/100)

触发阈值：
GCThreshold = LiveSize_after_last_GC × (1 + g/100)

示例：
GOGC = 100 (默认)
LiveSize = 100MB
则下次GC在堆达到200MB时触发

定义 4.2 (GC频率):
────────────────────────────────────────
GC间隔：
Interval = AllocationAmount / AllocationRate

其中AllocationAmount = LiveSize × g/100

GC频率：
Frequency = 1 / Interval = AllocationRate / (LiveSize × g/100)

调优公式：
增大GOGC → GC频率降低，但峰值内存增加
减小GOGC → GC频率增加，峰值内存降低

定理 4.1 (GOGC与CPU权衡):
────────────────────────────────────────
设g为GOGC值，则GC CPU占比：

GC_CPU_fraction ≈ (1 / g) × (MarkingCost / AllocationCost)

证明：
- 每次GC处理g × LiveSize的分配量
- 标记成本与LiveSize成正比
- 故每次GC CPU成本与LiveSize相关
- GC频率与1/g成正比∎
```

### 4.2 软内存限制

```text
定义 4.3 (GOMEMLIMIT):
────────────────────────────────────────
SoftMemoryLimit = M

目标：
Runtime会尽量保持内存使用 < M

但不会OOM，只是更激进的GC

实现机制：
1. 监控当前内存使用
2. 若接近限制，增加GC频率
3. 若超过限制，强制GC，甚至返回内存给OS

动态调整：
EffectiveGOGC = min(GOGC, f(MemoryUsage, Limit))

其中f是递减函数

定理 4.2 (内存限制保证):
────────────────────────────────────────
在软内存限制下，程序不会OOM
(除非实际需要超过限制)

但代价：
可能频繁的GC导致CPU使用率升高
```

---

## 五、形式化性能分析

### 5.1 暂停时间分析

```text
定义 5.1 (STW暂停):
────────────────────────────────────────
Go GC有两个STW阶段：
1. 标记开始 (很短，< 100μs)
2. 标记终止 (稍长，< 1ms)

暂停时间组成：
PauseTime = T_stop_world + T_root_scan + T_flush

其中：
├─ T_stop_world：停止所有P
├─ T_root_scan：扫描GC Roots
└─ T_flush：刷新写屏障缓冲区

优化技术：
- 并发标记准备
- 根对象增量扫描
- 屏障缓冲区异步处理

定理 5.1 (亚毫秒GC):
────────────────────────────────────────
Go 1.8+ 保证99%的GC暂停 < 1ms

证明思路：
- STW阶段只包含必须串行的操作
- 所有可并发的工作都在并发阶段完成
- 通过精细的工程优化，STW时间被压缩到亚毫秒级∎
```

### 5.2 吞吐量分析

```text
定义 5.2 (GC吞吐量):
────────────────────────────────────────
ApplicationThroughput =
  UsefulWork / (UsefulWork + GCWork)

Marcov模型：
设：
├─ λ：对象分配速率
├─ μ：对象死亡速率
├─ s：存活对象比例
└─ c：标记速率

稳态条件：
λ × s = μ × (1-s)

GC开销：
Overhead = λ / c

目标：
最小化Overhead，同时保持PauseTime < 阈值
```

---

*本章建立了Go垃圾回收器的形式化模型，包括状态机、标记算法和性能分析。*
