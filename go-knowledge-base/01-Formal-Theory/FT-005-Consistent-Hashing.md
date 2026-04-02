# FT-005: 一致性哈希算法 (Consistent Hashing)

> **维度**: Formal Theory
> **级别**: S (15+ KB)
> **标签**: #consistent-hashing #distributed-cache #load-balancing #virtual-nodes
> **权威来源**: [Consistent Hashing and Random Trees](https://www.akamai.com/us/en/multimedia/documents/technical-publication/consistent-hashing-and-random-trees-distributed-caching-protocols-for-relieving-hot-spots-on-the-world-wide-web-technical-publication.pdf), [Dynamo Paper](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf)

---

## 问题背景

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Traditional Hashing Problem                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  N = 4 nodes: node0, node1, node2, node3                                    │
│                                                                             │
│  Traditional: hash(key) % N                                                 │
│                                                                             │
│  Problem: When N changes (add/remove node), almost all keys remap!          │
│                                                                             │
│  Example:                                                                   │
│  hash("user:100") = 1234                                                    │
│  N=4: 1234 % 4 = 2 → node2                                                  │
│  N=5: 1234 % 5 = 4 → node4  (REMAPPED!)                                     │
│                                                                             │
│  Cache invalidation rate: ~75% when adding 1 node to 4                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 一致性哈希原理

### 核心思想

将节点和数据都映射到同一个哈希环上，数据由顺时针方向的第一个节点负责。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Consistent Hash Ring                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         hash=0°                                             │
│                            │                                                │
│                            ▼                                                │
│                   ┌─────────────────┐                                       │
│          hash=270°│                 │hash=90°                               │
│         ◄─────────┤   HASH SPACE    ├─────────►                             │
│                   │   [0, 2³²-1]    │                                       │
│                   │                 │                                       │
│                   └─────────────────┘                                       │
│                         ▲                                                   │
│                         │                                                   │
│                       hash=180°                                             │
│                                                                             │
│  Nodes on ring:                                                             │
│  • node-A at hash=45°                                                       │
│  • node-B at hash=135°                                                      │
│  • node-C at hash=225°                                                      │
│  • node-D at hash=315°                                                      │
│                                                                             │
│  Key placement:                                                             │
│  • key "user:1" → hash=50° → clockwise → node-B (135°)                      │
│  • key "user:2" → hash=200° → clockwise → node-C (225°)                     │
│                                                                             │
│  Add node-E at hash=100°:                                                   │
│  • Only keys between node-A (45°) and node-E (100°) move to node-E          │
│  • ~25% of keys affected (vs 75% in traditional)                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 形式化定义

$$
\begin{aligned}
&\text{Hash Ring: } H = [0, 2^{32}-1] \text{ (circular)} \\
&\text{Nodes: } N = \{n_1, n_2, ..., n_k\}, \text{ each } n_i \mapsto h(n_i) \in H \\
&\text{Keys: } K = \{k_1, k_2, ..., k_m\}, \text{ each } k_j \mapsto h(k_j) \in H \\
\\
&\text{Assignment Function: } \\
&assign(k) = \arg\min_{n \in N} \{ d(h(k), h(n)) \} \\
&\text{where } d(a, b) = (b - a) \mod 2^{32} \\
\\
&\text{Node Removal Impact: } \\
&\text{affected keys} = \{ k \in K : assign_{old}(k) = n_{removed} \} \\
&\text{keys remapped} = \text{affected keys} \\
&\text{remap ratio} = \frac{|\text{affected keys}|}{|K|} \approx \frac{1}{|N|}
\end{aligned}
$$

---

## Go 实现

```go
package consistenthash

import (
 "hash/crc32"
 "sort"
 "strconv"
 "sync"
)

// HashFunc 哈希函数类型
type HashFunc func(data []byte) uint32

// Map 一致性哈希环
type Map struct {
 hash     HashFunc       // 哈希函数
 replicas int            // 虚拟节点倍数
 ring     []int          // 排序后的哈希环
 nodes    map[int]string // 哈希值 -> 真实节点

 mu sync.RWMutex
}

// New 创建一致性哈希
func New(replicas int, fn HashFunc) *Map {
 m := &Map{
  replicas: replicas,
  hash:     fn,
  nodes:    make(map[int]string),
 }
 if m.hash == nil {
  m.hash = crc32.ChecksumIEEE
 }
 return m
}

// Add 添加节点
func (m *Map) Add(nodes ...string) {
 m.mu.Lock()
 defer m.mu.Unlock()

 for _, node := range nodes {
  // 创建虚拟节点
  for i := 0; i < m.replicas; i++ {
   hash := int(m.hash([]byte(strconv.Itoa(i) + node)))
   m.ring = append(m.ring, hash)
   m.nodes[hash] = node
  }
 }

 // 排序哈希环
 sort.Ints(m.ring)
}

// Remove 移除节点
func (m *Map) Remove(node string) {
 m.mu.Lock()
 defer m.mu.Unlock()

 // 移除该节点的所有虚拟节点
 for i := 0; i < m.replicas; i++ {
  hash := int(m.hash([]byte(strconv.Itoa(i) + node)))
  delete(m.nodes, hash)
 }

 // 重建哈希环
 m.ring = m.ring[:0]
 for hash := range m.nodes {
  m.ring = append(m.ring, hash)
 }
 sort.Ints(m.ring)
}

// Get 获取 key 对应的节点
func (m *Map) Get(key string) string {
 m.mu.RLock()
 defer m.mu.RUnlock()

 if len(m.ring) == 0 {
  return ""
 }

 hash := int(m.hash([]byte(key)))

 // 二分查找第一个 >= hash 的位置
 idx := sort.Search(len(m.ring), func(i int) bool {
  return m.ring[i] >= hash
 })

 // 如果超出范围，回到第一个节点
 if idx == len(m.ring) {
  idx = 0
 }

 return m.nodes[m.ring[idx]]
}

// GetN 获取 key 对应的 N 个不同节点（用于副本）
func (m *Map) GetN(key string, n int) []string {
 m.mu.RLock()
 defer m.mu.RUnlock()

 if len(m.ring) == 0 || n <= 0 {
  return nil
 }

 if n > len(m.ring) {
  n = len(m.ring)
 }

 hash := int(m.hash([]byte(key)))
 idx := sort.Search(len(m.ring), func(i int) bool {
  return m.ring[i] >= hash
 })

 result := make([]string, 0, n)
 seen := make(map[string]struct{})

 for len(result) < n {
  if idx >= len(m.ring) {
   idx = 0
  }

  node := m.nodes[m.ring[idx]]
  if _, ok := seen[node]; !ok {
   result = append(result, node)
   seen[node] = struct{}{}
  }

  idx++
 }

 return result
}
```

---

## 虚拟节点

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Virtual Nodes                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Problem: Without virtual nodes, uneven distribution causes hot spots       │
│                                                                              │
│  Physical Nodes          Virtual Nodes (replicas=3)                          │
│  ─────────────           ─────────────────────────                           │
│                                                                              │
│  node-A ────►  hash(node-A:0) = 100                                          │
│               hash(node-A:1) = 500                                           │
│               hash(node-A:2) = 900                                           │
│                                                                              │
│  node-B ────►  hash(node-B:0) = 200                                          │
│               hash(node-B:1) = 600                                           │
│               hash(node-B:2) = 1000                                          │
│                                                                              │
│  node-C ────►  hash(node-C:0) = 300                                          │
│               hash(node-C:1) = 700                                           │
│               hash(node-C:2) = 1100                                          │
│                                                                              │
│  Benefits:                                                                   │
│  1. Better load distribution (standard deviation ↓)                         │
│  2. Smoother migration when adding/removing nodes                           │
│  3. Multiple points for replication                                         │
│                                                                              │
│  Optimal replicas: 100-200 for production                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 性能分析

| 操作 | 时间复杂度 | 说明 |
|------|-----------|------|
| Add | O(replicas × log(replicas × nodes)) | 排序开销 |
| Remove | O(replicas × nodes) | 重建环 |
| Get | O(log(replicas × nodes)) | 二分查找 |
| GetN | O(N × log(replicas × nodes)) | N 个不同节点 |

---

## 实际应用

### 1. Redis Cluster

```go
// Redis 集群使用一致性哈希的变体（哈希槽）
const HashSlots = 16384

func slot(key string) int {
 start := strings.Index(key, "{")
 if start != -1 {
  end := strings.Index(key[start:], "}")
  if end != -1 && end != 1 {
   key = key[start+1 : start+end]
  }
 }
 return int(crc16(key) % HashSlots)
}
```

### 2. Nginx 负载均衡

```nginx
upstream backend {
    consistent_hash $request_uri;
    server backend1.example.com;
    server backend2.example.com;
    server backend3.example.com;
}
```

### 3. DynamoDB 分区

Dynamo 使用一致性哈希配合虚拟节点，每个物理节点负责多个虚拟节点。

---

## 参考文献

1. [Consistent Hashing and Random Trees](https://www.akamai.com/us/en/multimedia/documents/technical-publication/consistent-hashing-and-random-trees-distributed-caching-protocols-for-relieving-hot-spots-on-the-world-wide-web-technical-publication.pdf) - Karger et al., 1997
2. [Dynamo: Amazon's Highly Available Key-value Store](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf) - DeCandia et al., 2007
3. [Consistent Hashing in System Design](https://systemdesignprimer.com/consistent-hashing/) - System Design Primer
