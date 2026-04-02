# FT-008: 概率数据结构 (Probabilistic Data Structures)

> **维度**: Formal Theory
> **级别**: S (16+ KB)
> **标签**: #bloom-filter #hyperloglog #count-min-sketch #probabilistic
> **权威来源**: [Bloom Filters by Example](https://llimllib.github.io/bloomfilter-tutorial/), [Redis HyperLogLog](https://redis.io/docs/data-types/probabilistic/hyperloglog/)

---

## 核心原理

概率数据结构牺牲绝对准确性换取空间效率，适用于海量数据场景。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Probabilistic Data Structures                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Bloom Filter│  │  HyperLogLog│  │Count-Min    │  │   Cuckoo    │         │
│  │             │  │             │  │  Sketch     │  │   Filter    │         │
│  │ Membership  │  │Cardinality  │  │ Frequency   │  │ Membership  │         │
│  │ Query       │  │ Estimation  │  │ Estimation  │  │ + Delete    │         │
│  │             │  │             │  │             │  │             │         │
│  │ False +     │  │ ~2% error   │  │ Overestimate│  │ False +     │         │
│  │ No False -  │  │ 12KB/10^9   │  │ Error bound │  │ Can Delete  │         │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Bloom Filter

### 原理

```
位数组 + k 个哈希函数

Insert "hello":
  h1("hello") % m = 3  → set bit[3]
  h2("hello") % m = 7  → set bit[7]
  h3("hello") % m = 11 → set bit[11]

Query "hello":
  h1("hello") % m = 3  → bit[3] = 1 ✓
  h2("hello") % m = 7  → bit[7] = 1 ✓
  h3("hello") % m = 11 → bit[11] = 1 ✓
  Result: MAYBE IN SET

Query "world":
  h1("world") % m = 3  → bit[3] = 1 ✓
  h2("world") % m = 7  → bit[7] = 1 ✓
  h3("world") % m = 15 → bit[15] = 0 ✗
  Result: DEFINITELY NOT IN SET
```

### Go 实现

```go
package bloom

import (
    "hash/fnv"
    "math"
)

// Filter Bloom过滤器
type Filter struct {
    bits []bool
    k    int // 哈希函数数量
    m    int // 位数组大小
    n    int // 已插入元素数
}

// New 创建Bloom过滤器
// expectedElements: 预期元素数量
// falsePositiveRate: 目标误判率
func New(expectedElements int, falsePositiveRate float64) *Filter {
    // 最优位数组大小
    m := int(math.Ceil(-float64(expectedElements) * math.Log(falsePositiveRate) / (math.Log(2) * math.Log(2))))
    // 最优哈希函数数量
    k := int(math.Round(float64(m) / float64(expectedElements) * math.Log(2)))

    return &Filter{
        bits: make([]bool, m),
        k:    k,
        m:    m,
    }
}

// hash 计算多个哈希值
func (f *Filter) hash(data []byte) []int {
    h1 := fnv.New64a()
    h1.Write(data)
    hash1 := h1.Sum64()

    h2 := fnv.New64()
    h2.Write(data)
    hash2 := h2.Sum64()

    results := make([]int, f.k)
    for i := 0; i < f.k; i++ {
        // 使用双重哈希生成 k 个哈希值
        hash := (hash1 + uint64(i)*hash2) % uint64(f.m)
        results[i] = int(hash)
    }
    return results
}

// Add 添加元素
func (f *Filter) Add(data []byte) {
    for _, pos := range f.hash(data) {
        f.bits[pos] = true
    }
    f.n++
}

// Contains 查询元素可能存在
func (f *Filter) Contains(data []byte) bool {
    for _, pos := range f.hash(data) {
        if !f.bits[pos] {
            return false // 肯定不存在
        }
    }
    return true // 可能存在
}

// FalsePositiveRate 计算当前误判率
func (f *Filter) FalsePositiveRate() float64 {
    return math.Pow(1-math.Exp(-float64(f.k*f.n)/float64(f.m)), float64(f.k))
}
```

---

## HyperLogLog

### 原理

基于概率统计估计基数（唯一元素数量），误差约 0.81%。

```go
package main

import (
    "github.com/axiomhq/hyperloglog"
)

func main() {
    h := hyperloglog.New14() // 2^14 个寄存器

    // 添加元素
    h.Insert([]byte("user1"))
    h.Insert([]byte("user2"))
    h.Insert([]byte("user1")) // 重复

    // 估计基数
    count := h.Estimate()
    fmt.Printf("Estimated unique users: %d\n", count) // ~2
}
```

### Redis 实现

```bash
# 添加元素
PFADD visitors user1 user2 user3

# 估计基数
PFCOUNT visitors

# 合并多个HyperLogLog
PFMERGE total visitors1 visitors2
```

---

## Count-Min Sketch

### 应用

```
流量统计：
- 统计每个IP的请求次数（允许高估）
- 空间复杂度 O((1/ε) * log(1/δ))

热门查询：
- 实时统计搜索热词
- 流式处理场景
```

---

## 应用场景

| 场景 | 数据结构 | 收益 |
|------|---------|------|
| 缓存穿透防护 | Bloom Filter | 1GB → 1MB |
| UV 统计 | HyperLogLog | 内存减少 1000x |
| 频率估计 | Count-Min Sketch | 实时流处理 |
| 去重 | Cuckoo Filter | 支持删除操作 |

---

## 参考文献

1. [Bloom Filters by Example](https://llimllib.github.io/bloomfilter-tutorial/)
2. [HyperLogLog: the analysis of a near-optimal cardinality estimation algorithm](http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf)
3. [Redis HyperLogLog](https://redis.io/docs/data-types/probabilistic/hyperloglog/)
