# TS-003: Redis 数据结构的代数与复杂度 (Redis Data Structures: Algebra & Complexity)

> **维度**: Technology Stack
> **级别**: S (20+ KB)
> **标签**: #redis #data-structures #complexity-analysis #internals #algorithms
> **权威来源**:
>
> - [Redis Documentation: Internals](https://redis.io/docs/reference/internals/) - Redis Ltd (2025)
> - [Redis Source Code](https://github.com/redis/redis) - GitHub
> - [Skip Lists: A Probabilistic Alternative to Balanced Trees](https://dl.acm.org/doi/10.1145/78973.78977) - Pugh (1990)
> - [The Art of Computer Programming, Vol 3](https://www-cs-faculty.stanford.edu/~knuth/taocp.html) - Knuth (Sorting & Searching)
> - [SipHash: A Fast Short-Input PRF](https://131002.net/siphash/) - Aumasson & Bernstein (2012)

---

## 1. Redis 对象系统的代数结构

### 1.1 对象类型代数

**定义 1.1 (Redis 对象)**
Redis 对象 $o$ 是一个五元组 $\langle \text{type}, \text{encoding}, \text{ptr}, \text{refcount}, \text{lru} \rangle$：

- $type \in \{\text{STRING}, \text{LIST}, \text{HASH}, \text{SET}, \text{ZSET}, ...\}$: 逻辑类型
- $encoding \in \{\text{RAW}, \text{INT}, \text{HT}, \text{ZIPLIST}, ...\}$: 物理编码
- $ptr$: 指向数据的指针
- $refcount \in \mathbb{N}$: 引用计数
- $lru \in \mathbb{N}$: 最后访问时间

**定义 1.2 (编码转换函数)**
$$\text{encode}: \text{Type} \times \text{Data} \to \text{Encoding}$$
根据数据特征选择最优编码。

**示例转换规则**:

```
String:
  len ≤ 20 && is_integer → INT
  len ≤ 44 → EMBSTR (embedded string)
  else → RAW

List:
  len < 512 && size < 64B → ZIPLIST
  else → QUICKLIST (ziplist + linked list)

Hash:
  len < 512 && size < 64B → ZIPLIST
  else → HT (hashtable)

Set:
  len < 512 && integers → INTSET
  else → HT

ZSet:
  len < 128 → ZIPLIST
  else → SKIPLIST + HT
```

### 1.2 引用计数与内存管理

**定义 1.3 (对象生命周期)**
$$\text{lifecycle}(o) = \text{Create} \to \text{Share}^* \to \text{Release}$$

**引用计数操作**:

- `INCRREFCOUNT(o)`: refcount++
- `DECRREFCOUNT(o)`: refcount--; if 0 then `free(o)`

**共享对象池**:
$$\text{SharedPool} = \{ o_i \mid 0 \leq i \leq 9999 \}$$
整数 0-9999 预创建共享对象，减少内存分配。

---

## 2. 核心数据结构的复杂度分析

### 2.1 String: SDS (Simple Dynamic String)

**定义 2.1 (SDS)**
$$\text{SDS} = \langle \text{len}, \text{alloc}, \text{buf} \rangle$$

- $len$: 已用长度
- $alloc$: 分配总长度
- $buf[len+1]$: 字符数组（兼容 C 字符串）

**操作复杂度**:

| 操作 | 复杂度 | 说明 |
|------|--------|------|
| `sdsnew` | $O(n)$ | 复制字符串 |
| `sdslen` | $O(1)$ | 直接返回 len |
| `sdscat` | $O(n)$ 或 $O(1)$ | 有空间则直接追加，否则 realloc |
| `sdsrange` | $O(n)$ | 截取子串 |

**空间预分配策略**:
$$\text{new_alloc} = \begin{cases} \text{len} + \text{addlen} & \text{if } < 1MB \\ \text{len} + \text{addlen} + 1MB & \text{if } \geq 1MB \end{cases}$$
小于 1MB 时翻倍，大于 1MB 时加 1MB。

### 2.2 Hash: 渐进式 Rehash

**定义 2.2 (字典)**
$$\text{Dict} = \langle \text{ht}[0], \text{ht}[1], \text{rehashidx} \rangle$$

- `ht[0]`: 主哈希表
- `ht[1]`: 渐进式 rehash 时使用
- `rehashidx`: 当前 rehash 进度

**渐进式 Rehash 算法**:

```
While rehashing:
  1. Each CRUD op on ht[0] migrates rehashidx bucket to ht[1]
  2. rehashidx++
  3. When ht[0] empty:
     - Free ht[0]
     - ht[0] = ht[1]
     - ht[1] = NULL
     - rehashidx = -1
```

**复杂度保证**:

- 单次操作: $O(1)$ 摊还
- 分摊 rehash 成本，避免停顿

### 2.3 Sorted Set: SkipList + HashTable

**定义 2.3 (SkipList)**
概率性平衡数据结构，期望层数 $E[levels] = O(\log n)$。

**节点结构**:

```c
typeof struct zskiplistNode {
    sds ele;                    // 成员
    double score;               // 分数
    struct zskiplistNode *backward;  // 后退指针
    struct zskiplistLevel {     // 层数组
        struct zskiplistNode *forward;
        unsigned int span;      // 跨度
    } level[];
} zskiplistNode;
```

**操作复杂度**:

| 操作 | SkipList | HashTable | 组合 |
|------|----------|-----------|------|
| `ZADD` | $O(\log n)$ | $O(1)$ | $O(\log n)$ |
| `ZSCORE` | $O(n)$ | $O(1)$ | $O(1)$ (查 HT) |
| `ZRANGE` | $O(\log n + m)$ | - | $O(\log n + m)$ |
| `ZREM` | $O(\log n)$ | $O(1)$ | $O(\log n)$ |

**双结构优势**:

- SkipList: 有序遍历、范围查询
- HashTable: O(1) 成员查找

### 2.4 HyperLogLog: 基数估计

**定义 2.4 (HLL)**
概率数据结构，使用 $O(\log \log n)$ 空间估计基数。

**算法**:

1. 哈希每个元素到 $p$ 位二进制
2. 前 $p$ 位决定桶索引
3. 后 $b-p$ 位找最高 0 的位置 $\rho$
4. 更新桶: $M[j] = \max(M[j], \rho)$

**估计公式**:
$$E = \alpha_m \cdot m^2 \cdot \left( \sum_{j=1}^{m} 2^{-M[j]} \right)^{-1}$$

**误差**: 标准误差 $\approx 1.04/\sqrt{m}$

---

## 3. 内存优化策略

### 3.1 ziplist 压缩列表

**结构**:

```
<zlbytes> <zltail> <zllen> <entry> <entry> ... <entry> <zlend>
```

**Entry 编码**:

| 编码 | 长度 | 说明 |
|------|------|------|
| `00xxxxxx` | 1 byte | 6 bit string len + data |
| `01xxxxxx xxxxxxxx` | 2 bytes | 14 bit string len + data |
| `10xxxxxx` + 4 bytes | 5 bytes | 32 bit string len + data |
| `11000000` | 1 byte | int16_t |
| `11010000` | 1 byte | int32_t |
| `11100000` | 1 byte | int64_t |
| `11110000` | 1 byte | 24 bit signed int |
| `11111110` | 1 byte | 8 bit signed int |

**连锁更新问题**:
当 entry 大小变化需要改变编码长度时，可能触发连锁更新。
**复杂度**: $O(n^2)$ 最坏情况，但 $O(n)$ 摊还。

### 3.2 内存回收策略

**定义 3.1 (驱逐策略)**
$$\text{EvictionPolicy} \in \{\text{LRU}, \text{LFU}, \text{RANDOM}, \text{TTL}, ...\}$$

**LRU 近似**:
Redis 使用 24-bit 字段记录访问时间，精度为秒级。

**抽样 LRU**:

- 默认抽样 5 个 key
- 驱逐最老的
- $O(1)$ 复杂度，近似 LRU

---

## 4. 多元表征

### 4.1 Redis 数据结构层次图

```
Redis Data Structures
├── String (SDS)
│   ├── INT (long long)
│   ├── EMBSTR (len ≤ 44, embedded)
│   └── RAW (len > 44, separate allocation)
│
├── List
│   ├── ZIPLIST (small, compact)
│   └── QUICKLIST (ziplist nodes + doubly linked)
│       └── 配置: list-max-ziplist-size, list-compress-depth
│
├── Hash
│   ├── ZIPLIST (small, field-value pairs)
│   └── HT (hashtable with progressive rehash)
│       └── 负载因子: <0.1 shrink, >1 expand, >5 force rehash
│
├── Set
│   ├── INTSET (integers only, sorted)
│   └── HT (general, values=NULL)
│
├── Sorted Set (ZSet)
│   ├── ZIPLIST (small)
│   └── SKIPLIST + HT (large)
│       ├── SkipList: sorted by score, O(log n) ops
│       └── HashTable: member → score, O(1) lookup
│
├── Bitmap
│   └── RAW string (bit operations)
│
├── HyperLogLog
│   ├── SPARSE (sparse representation)
│   └── DENSE (dense representation, 12KB fixed)
│
├── Stream
│   └── RAX tree (radix tree) + listpacks
│       └── Consumer groups with PEL (Pending Entries List)
│
└── Geospatial
    └── Sorted Set (geohash as score)
```

### 4.2 编码选择决策树

```
选择 Redis 编码?
│
├── String
│   ├── 整数且 0-9999? → 使用共享对象池
│   ├── 整数且可表示为 long? → INT encoding
│   ├── 长度 ≤ 44? → EMBSTR (内存连续，1次分配)
│   └── 长度 > 44? → RAW (分离 sds 结构)
│
├── List
│   ├── 长度 < 512 且元素大小 < 64B? → ZIPLIST
│   └── 否则 → QUICKLIST
│       └── 是否压缩?
│           ├── 是 → LZF 压缩节点
│           └── 否 → 原始数据
│
├── Hash
│   ├── 字段数 < 512 且大小 < 64B? → ZIPLIST
│   └── 否则 → HT
│       └── 渐进式 rehash 进行中?
│           ├── 是 → 操作同时迁移 bucket
│           └── 否 → 正常操作
│
├── Set
│   ├── 全是整数且数量 < 512? → INTSET
│   └── 否则 → HT (value 为 NULL)
│
└── ZSet
    ├── 元素数 < 128 且大小 < 64B? → ZIPLIST
    └── 否则 → SKIPLIST + HT
        ├── 插入/删除: 更新 skiplist 和 hashtable
        └── 查找 score: O(1) via hashtable
        └── 范围查询: O(log n + m) via skiplist
```

### 4.3 数据结构复杂度对比矩阵

| 结构 | 访问 | 插入 | 删除 | 搜索 | 范围 | 空间 | 适用场景 |
|------|------|------|------|------|------|------|---------|
| **String** | $O(1)$ | - | - | - | - | $O(n)$ | 缓存、计数 |
| **List** | $O(n)$ | $O(1)$* | $O(1)$* | $O(n)$ | 支持 | $O(n)$ | 队列、流 |
| **Hash** | $O(1)$ | $O(1)$ | $O(1)$ | $O(1)$ | 不支持 | $O(n)$ | 对象存储 |
| **Set** | - | $O(1)$ | $O(1)$ | $O(1)$ | 不支持 | $O(n)$ | 去重、交集 |
| **ZSet** | - | $O(\log n)$ | $O(\log n)$ | $O(1)$ | $O(\log n+m)$ | $O(n)$ | 排行榜 |
| **Bitmap** | $O(1)$ | - | - | $O(n/w)$** | 支持 | $O(n/w)$ | 在线状态 |
| **HLL** | - | $O(1)$ | - | - | - | $O(12KB)$ | UV 统计 |

*头部/尾部操作 $O(1)$，中间 $O(n)$
**按字长 $w$ 计算

---

## 5. 算法与实现细节

### 5.1 SipHash 哈希算法

Redis 使用 SipHash-1-2 防止哈希洪水攻击。

**特性**:

- 密钥化哈希（keyed hash）
- 抗碰撞攻击
- 性能好（短输入优化）

### 5.2 整数集合 (intset)

**升级策略**:
当插入更大类型整数时，整个集合升级。

- int16 → int32 → int64
- $O(n)$ 升级成本，但摊还 $O(1)$

---

## 6. 参考文献

1. **Pugh, W. (1990)**. Skip Lists: A Probabilistic Alternative to Balanced Trees. *CACM*.
2. **Knuth, D. E. (1998)**. The Art of Computer Programming, Vol 3. *Addison-Wesley*.
3. **Aumasson, J. P., & Bernstein, D. J. (2012)**. SipHash: A Fast Short-Input PRF. *ACS*.
4. **Redis Authors (2025)**. Redis Documentation. *Redis Ltd*.
5. **Flajolet, P., et al. (2007)**. HyperLogLog: The Analysis of a Near-Optimal Cardinality Estimation Algorithm. *AOFA*.

---

## 7. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Redis Design Checklist                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  数据结构选择:                                                               │
│  □ String: 简单值、计数器 (原子 INCR)                                         │
│  □ List: 队列 (LPUSH/BRPOP)、时间线                                           │
│  □ Hash: 对象存储 (字段级操作)                                                 │
│  □ Set: 去重、交集、并集、差集                                                 │
│  □ ZSet: 排行榜、范围查询                                                      │
│  □ Bitmap: 大规模布尔标志 (在线状态)                                          │
│  □ HLL: UV 统计 (误差可接受时)                                                │
│  □ Stream: 消息队列、事件溯源                                                  │
│                                                                              │
│  内存优化:                                                                   │
│  □ 使用合适编码 (利用 ziplist/intset 节省内存)                                │
│  □ 配置 maxmemory 和淘汰策略                                                   │
│  □ 大 key 拆分 (避免阻塞)                                                     │
│  □ 使用压缩 (list-compress-depth, hash-max-ziplist-entries)                  │
│                                                                              │
│  性能优化:                                                                   │
│  □ 避免 O(n) 操作 (KEYS, FLUSHALL)                                            │
│  □ 使用 SCAN 替代 KEYS                                                        │
│  □ Pipeline 批量操作                                                          │
│  □ 合理使用 Pipeline 和事务                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
