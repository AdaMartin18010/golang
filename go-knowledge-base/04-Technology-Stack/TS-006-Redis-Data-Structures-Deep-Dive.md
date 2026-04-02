# TS-006: Redis 数据结构深度解析 (Redis Data Structures Deep Dive)

> **维度**: Technology Stack
> **级别**: S (16+ KB)
> **标签**: #redis #data-structures #internals #performance
> **权威来源**: [Redis Data Types](https://redis.io/docs/data-types/), [Redis Internals](https://redis.io/docs/reference/internals/)

---

## 底层数据结构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Redis Object System                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  RedisObject (robj)                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  type:  STRING | LIST | HASH | SET | ZSET | STREAM | ...          │    │
│  │  encoding: RAW | INT | EMBSTR | HT | ZIPLIST | INTSET | ...        │    │
│  │  ptr ──► 实际数据结构 (sds, dict, ziplist, skiplist, etc.)         │    │
│  │  refcount: 引用计数                                                 │    │
│  │  lru: 最后访问时间 (用于淘汰)                                        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  Encoding 演进 (Redis 8.2):                                                 │
│  String: INT < 8B → EMBSTR < 44B → RAW                                     │
│  List:   ZIPLIST (小) → QUICKLIST (压缩节点) → LISTPACK (Redis 7+)         │
│  Hash:   ZIPLIST < 512 → HT                                                │
│  Set:    INTSET < 512 → HT                                                 │
│  ZSet:   ZIPLIST < 128 → SKIPLIST + HT                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## String 内部实现

### SDS (Simple Dynamic String)

```c
// sds 结构 (Redis 8.2)
struct __attribute__ ((__packed__)) sdshdr64 {
    uint64_t len;        // 已用长度
    uint64_t alloc;      // 分配总长度
    unsigned char flags; // 类型标记
    char buf[];          // 柔性数组
};

// 优势:
// - O(1) 获取长度
// - 预分配减少内存重分配
// - 二进制安全 (允许 \0)
// - 兼容 C 字符串
```

### 整数优化

```bash
# 小整数使用共享对象 (0-9999)
SET counter 100
# 内部: ptr 直接存储整数，而非字符串

# 大整数
SET big "12345678901234567890"
# 内部: raw sds
```

---

## Hash 内部实现

### 渐进式 Rehash

```
Hash Table 扩容过程:

步骤 1: 创建新表 (2x 大小)
ht[0]: 旧表 (使用中)
ht[1]: 新表 (空)

步骤 2: 渐进式迁移
每次 CRUD 操作迁移 N 个槽位
┌─────────┐      ┌─────────┐
│  ht[0]  │ ──►  │  ht[1]  │
│  [0]    │      │  [0]    │
│  [1]    │ ──►  │  [1]    │  每次操作迁移 1 个槽
│  [2]    │      │  [2]    │
│  ...    │ ──►  │  ...    │
└─────────┘      └─────────┘

步骤 3: 交换完成
ht[0] = ht[1]
ht[1] = nil
```

### 内存优化

```bash
# Hash 编码转换阈值
hash-max-ziplist-entries 512    # 字段数
hash-max-ziplist-value 64       # 值大小 (字节)

# 小 Hash 使用 ziplist (连续内存)
HSET user:1 name "Alice" age 30
# 编码: ziplist (O(n) 查找但内存高效)

# 大 Hash 使用 hashtable
HSET big:hash field1 value1 ... field1000 value1000
# 编码: hashtable (O(1) 查找)
```

---

## ZSet (Sorted Set) 内部实现

### SkipList + HashTable 双结构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        ZSet Structure                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  HashTable (用于 O(1) 成员查找)                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  "member1" ──► {score: 100, skiplist_node_ptr}                 │    │
│  │  "member2" ──► {score: 200, skiplist_node_ptr}                 │    │
│  │  "member3" ──► {score: 150, skiplist_node_ptr}                 │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  SkipList (用于范围查询)                                                  │
│                              ┌──────────────────────┐                   │
│  level 3:    ┌───────────────┼──────────────────────┼───────────┐      │
│  level 2:    ├───────────────┼──────────────────────┤───────┐   │      │
│  level 1:    ├───┐   ┌───────┼───────┐   ┌───────┐  │   ┌───┼───┐   │  │
│  level 0:    50  100  120   150     180 200     250 300 350 400 450   │  │
│              │   │   │       │       │   │       │   │   │   │   │    │  │
│             m1   m2  m4     m3      m6  m5      m8  m7  m9  m10 m11   │  │
│                                                                          │
│  复杂度:                                                                  │
│  ZADD: O(log N)                                                          │
│  ZRANK: O(log N)                                                         │
│  ZRANGE: O(log N + M)  (M = 返回元素数)                                   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 生产优化建议

| 场景 | 推荐结构 | 避免 |
|------|---------|------|
| 热点 Key | 拆分 + 本地缓存 | 单 Key 大 Value |
| 计数器 | String (INCR) | Hash 存储单个计数 |
| 会话存储 | Hash (HGETALL) | 大量小 String |
| 排行榜 | ZSet | 自己维护排序 |
| 时间线 | Stream / List | ZSet 按时间排序 |

---

## 参考文献

1. [Redis Data Types](https://redis.io/docs/data-types/)
2. [Redis Internals](https://redis.io/docs/reference/internals/)
3. [Redis Source Code](https://github.com/redis/redis)
