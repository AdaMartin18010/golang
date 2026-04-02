# 对比分析：Redis vs Memcached (Comparison: Redis vs Memcached)

> **分类**: 技术栈对比
> **相关文档**: TS-002

---

## 核心对比

| 特性 | Redis | Memcached |
|------|-------|-----------|
| **数据结构** | 丰富 (String, List, Set, Hash, ZSet, Stream, Vector) | 仅 String |
| **持久化** | RDB + AOF | 不支持 |
| **集群** | Redis Cluster (原生) | 客户端分片 |
| **多线程** | I/O 多线程 (8.2+) | 多线程 |
| **内存效率** | 中等 | 高 |
| **功能** | 发布订阅、Lua、事务 | 简单 KV |
| **适用场景** | 缓存、消息队列、向量搜索 | 纯缓存 |

---

## 数据结构对比

```
Redis:
├── String (二进制安全)
├── List (Linked List / QuickList)
├── Set (Hash Table / IntSet)
├── Sorted Set (Skip List + Hash)
├── Hash (ziplist / hashtable)
├── Stream (Radix Tree)
├── Bitmap
├── HyperLogLog
└── Vector (8.2+)

Memcached:
└── String only
```

---

## 性能对比

| 指标 | Redis | Memcached |
|------|-------|-----------|
| 读 QPS | 1M+ | 1M+ |
| 写 QPS | 800K | 900K |
| 内存/Key | 较高 | 较低 |
| 网络 I/O | 高吞吐 (8.2+) | 高吞吐 |
| 小对象 (<1KB) | 快 | 更快 |
| 大对象 (>100KB) | 更适合 | 一般 |

---

## 选择建议

| 场景 | 推荐 |
|------|------|
| 简单缓存 | Memcached |
| 复杂数据结构 | Redis |
| 需要持久化 | Redis |
| 消息队列 | Redis Streams |
| 分布式锁 | Redis RedLock |
| 向量搜索 | Redis 8.2+ |
| 纯内存效率 | Memcached |

---

## 架构对比

```
Redis Cluster (去中心化)
┌─────┐  ┌─────┐  ┌─────┐
│M1-S1│  │M2-S2│  │M3-S3│
└──┬──┘  └──┬──┘  └──┬──┘
   │        │        │
   └────────┼────────┘
            │
      Gossip Protocol

Memcached (客户端分片)
┌─────┐  ┌─────┐  ┌─────┐
│ S1  │  │ S2  │  │ S3  │
└─────┘  └─────┘  └─────┘
   ↑        ↑        ↑
   └────────┼────────┘
            │
      Client Hash Ring
```

---

## 总结

- **Memcached**: 简单、高效、纯缓存
- **Redis**: 功能丰富、多用途、生态系统完善
