# TS-042-Redis-8-Features

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: Redis 8.0-8.6 (2025-2026)
> **Size**: >20KB

---

## 1. Redis 8 概览

### 1.1 发布历史

| 版本 | 发布日期 | 关键特性 |
|------|---------|---------|
| 8.0 | 2025年5月 | 9种新数据结构集成 |
| 8.2 | 2025年末 | 流增强、位图操作 |
| 8.4 | 2026年初 | 原子槽迁移、SIMD优化 |
| 8.6 | 2026年2月 | 流幂等性、热键检测 |

### 1.2 核心改进

| 指标 | 改进 |
|------|------|
| 命令执行 | 最高87%更快 |
| 吞吐量 | 最高112%提升 |
| 复制 | 18%更快 |
| 内存(复制) | 35%减少 |
| 向量操作 | 5x吞吐量 |

---

## 2. Redis 8.0: 9种新数据结构

### 2.1 集成概述

Redis 8将原Redis Stack模块集成到核心：

```
Redis 8 Core Data Structures:
├── 基础: String, Hash, List, Set, Sorted Set
├── 模块集成:
│   ├── RediSearch → 内置查询引擎
│   ├── RedisJSON → JSON类型
│   ├── RedisTimeSeries → 时间序列
│   └── 概率数据结构:
│       ├── Bloom Filter
│       ├── Cuckoo Filter
│       ├── Count-Min Sketch
│       ├── Top-K
│       └── T-Digest
└── 新增:
    └── Vector Set (Beta)
```

### 2.2 向量集 (Vector Set)

```bash
# 添加向量
VADD my_vectors vec1 [0.1, 0.2, 0.3, ...]
VADD my_vectors vec2 [0.4, 0.5, 0.6, ...]

# 相似度搜索
VSIM my_vectors VEC [0.15, 0.25, 0.35, ...] K 10

# 带过滤的搜索
VSIM my_vectors VEC [...] K 10 FILTER "category = 'electronics'"
```

**应用场景**:

- 语义搜索
- 推荐系统
- 相似度匹配

### 2.3 JSON数据类型

```bash
# 设置JSON
JSON.SET user:1 $ '{"name": "Alice", "age": 30, "tags": ["dev", "go"]}'

# 路径查询
JSON.GET user:1 $.name
# "Alice"

# 条件查询
JSON.GET user:1 '$.tags[?(@ == "go")]'
# ["go"]

# 数组操作
JSON.ARRAPPEND user:1 $.tags '"rust"'
JSON.GET user:1 $.tags
# ["dev", "go", "rust"]

# 增量更新
JSON.NUMINCRBY user:1 $.age 1
```

### 2.4 时间序列

```bash
# 创建时间序列
TS.CREATE sensor:temperature RETENTION 86400000 LABELS location room1

# 添加样本
TS.ADD sensor:temperature * 23.5
TS.ADD sensor:temperature * 24.1

# 范围查询
TS.RANGE sensor:temperature - + AGGREGATION avg 60000
# 每分钟平均值

# 多序列聚合
TS.MRANGE - + FILTER location=room1 AGGREGATION max 3600000
```

### 2.5 概率数据结构

```bash
# Bloom Filter - 成员检测
BF.ADD users:filter user123
BF.EXISTS users:filter user123
# 1
BF.EXISTS users:filter user999
# 0 (可能有误判)

# Cuckoo Filter - 可删除的Bloom Filter
CF.ADD items:filter item1
CF.DEL items:filter item1

# Count-Min Sketch - 频率估计
CMS.ADD events:counters event1
CMS.ADD events:counters event1
CMS.QUERY events:counters event1
# 2

# Top-K - 热门项
TOPK.ADD trending:keywords "redis" "go" "database"
TOPK.LIST trending:keywords
# ["redis", "go", "database"]

# T-Digest - 分位数估计
TDIGEST.ADD latencies:api 12.5 15.3 11.2 18.7
TDIGEST.QUANTILE latencies:api 0.99
# 18.5
```

---

## 3. Redis 8.2: 流增强

### 3.1 新命令

```bash
# XDELEX: 删除流条目并处理消费者组
XDELEX mystream 1526569495631-0
# 自动处理待处理条目(PEL)

# XACKDEL: 确认并删除
XACKDEL mystream mygroup 1526569495631-0
# 原子性确认+删除

# CLUSTER SLOT-STATS: 每槽使用指标
CLUSTER SLOT-STATS
# 返回每个槽的键数量、内存使用
```

### 3.2 位图操作增强

```bash
# 新的BITOP操作符
BITOP DIFF result key1 key2      # 差集
BITOP DIFF1 result key1 key2     # 单向差集
BITOP ANDOR result key1 key2     # AND后OR
BITOP ONE result key1 key2       # 仅保留一位
```

---

## 4. Redis 8.4: 性能优化

### 4.1 SIMD优化

**AVX2/AVX512/ARM Neon** 加速：

| 操作 | 加速 |
|------|------|
| BITCOUNT | SIMD popcount |
| HyperLogLog | 向量化 |
| 向量操作 (VADD/VSIM) | 并行计算 |

### 4.2 原子槽迁移

```bash
# 新的CLUSTER MIGRATION命令
CLUSTER MIGRATION <slot> TO <node-id>

# 特点:
# - 原子性迁移
# - 无数据丢失
# - 最小中断
```

### 4.3 混合搜索 (FT.HYBRID)

```bash
# 向量 + 全文混合搜索
FT.HYBRID idx "laptop"
  VECTOR query_vec [...]
  RANKING RRF
  LIMIT 0 20

# RRF: Reciprocal Rank Fusion
# 结合多种排名方法
```

---

## 5. Redis 8.6: 流幂等性

### 5.1 XADD幂等性

```bash
# 最多一次交付
XADD mystream IDMPAUTO * field value
# IDMPAUTO: 自动生成幂等ID

# 显式幂等ID
XADD mystream IDMP <id> field value
# 相同ID不会重复添加

# 应用场景:
# - 网络重试安全
# - 恰好一次处理
```

### 5.2 热键检测

```bash
# 检测访问频率最高的键
HOTKEYS
# 返回: key1, key2, ... (按访问频率排序)

# 配置热键阈值
CONFIG SET hotkeys-threshold 1000
```

### 5.3 新淘汰策略

```bash
# 最近修改淘汰 (LRM)
CONFIG SET maxmemory-policy volatile-lrm
# 优先淘汰最近修改的键

CONFIG SET maxmemory-policy allkeys-lrm
# 所有键，按最近修改时间
```

---

## 6. 安全更新

### 6.1 Redis 8.0.4 (2025年10月)

| CVE | 描述 |
|-----|------|
| CVE-2025-49844 | Lua脚本RCE |
| CVE-2025-46817 | Lua整数溢出 |
| CVE-2025-46818 | 跨用户Lua执行 |
| CVE-2025-46819 | Lua越界读取 |

### 6.2 安全建议

```bash
# 升级到最新补丁版本
redis-server --version
# 确保 >= 8.0.4

# 配置安全选项
CONFIG SET requirepass strong_password
CONFIG SET maxclients 10000
```

---

## 7. ACL变更

### 7.1 新类别

```bash
# 模块命令现在包含在标准类别中
ACL SETUSER appuser +@read
# 现在允许: GET, FT.SEARCH, JSON.GET, etc.

ACL SETUSER readonly +@all -@write
# 现在阻止: SET, JSON.SET, FT.ADD, etc.

# 新类别
@search    # 搜索命令
@json      # JSON命令
@timeseries # 时间序列
@bloom     # Bloom Filter
@cuckoo    # Cuckoo Filter
@cms       # Count-Min Sketch
@topk      # Top-K
@tdigest   # T-Digest
```

---

## 8. 配置优化

### 8.1 I/O线程

```bash
# Redis 8.0+ I/O线程
io-threads 8
io-threads-do-reads yes

# 效果: 最高112%吞吐量提升
```

### 8.2 内存优化

```bash
# Redis 8.6内存减少
hash-max-listpack-entries 512
zset-max-listpack-entries 128

# 效果:
# - Hash: 最高16.7%减少
# - ZSet: 最高30.5%减少
```

---

## 9. 部署架构

### 9.1 高可用架构

```
┌─────────────────────────────────────────┐
│         Redis 8 Cluster                 │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐ │
│  │ Master1 │◄──►│ Master2 │◄──►│ Master3 │ │
│  │ Slave1a │    │ Slave2a │    │ Slave3a │ │
│  │ Slave1b │    │ Slave2b │    │ Slave3b │ │
│  └─────────┘    └─────────┘    └─────────┘ │
│       │              │              │      │
│       └──────────────┼──────────────┘      │
│                      │                     │
│              Sentinel (HA)                 │
│              or Redis Cluster              │
└─────────────────────────────────────────┘
```

### 9.2 模块迁移

```bash
# 从Redis Stack迁移到Redis 8
# 无需更改代码，API兼容

# 原有代码:
# MODULE LOAD /path/to/redisearch.so

# Redis 8:
# 无需加载，内置支持
FT.CREATE idx ...
JSON.SET key ...
TS.ADD key ...
```

---

## 10. 最佳实践

### 10.1 数据结构选择

| 场景 | 推荐数据结构 |
|------|-------------|
| 缓存 | String/Hash |
| 会话 | Hash with HEXPIRE |
| 排行榜 | Sorted Set |
| 地理位置 | Geo + Search |
| 语义搜索 | Vector Set |
| 时间序列 | TimeSeries |
| 布隆过滤 | Bloom Filter |
| 实时统计 | Count-Min Sketch |
| JSON文档 | JSON |
| 流处理 | Streams |

### 10.2 性能调优

```bash
# 1. 启用I/O线程
io-threads 8

# 2. 调整内存策略
maxmemory-policy allkeys-lru

# 3. 持久化配置
save 900 1
save 300 10
save 60 10000

# 4. 客户端输出缓冲区
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit replica 256mb 64mb 60
```

---

## 11. 升级指南

### 11.1 从Redis 7.x升级

```bash
# 1. 备份数据
redis-cli BGSAVE

# 2. 测试兼容性
# Redis 8向后兼容Redis 7协议

# 3. 滚动升级 (主从架构)
# 先升级从节点
redis-cli -h slave1 DEBUG RELOAD

# 4. 故障转移到升级后的从节点
redis-cli -h master1 CLUSTER FAILOVER

# 5. 升级原主节点
```

### 11.2 配置迁移

```bash
# 检查废弃配置
redis-server --test-memory

# 更新redis.conf
# 添加新的ACL类别配置
```

---

## 12. 参考文献

1. Redis 8.0 Release Notes
2. Redis 8.6 Performance Blog
3. Redis Vector Set Documentation
4. Redis JSON Commands
5. Redis TimeSeries Guide

---

*Last Updated: 2026-04-03*
