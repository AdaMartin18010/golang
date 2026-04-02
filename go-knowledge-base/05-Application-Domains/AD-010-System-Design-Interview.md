# AD-010: 系统设计面试指南 (System Design Interview Guide)

> **维度**: Application Domains  
> **级别**: S (18+ KB)  
> **标签**: #system-design #interview #architecture #scalability  
> **权威来源**: [Designing Data-Intensive Applications](https://dataintensive.net/), [System Design Primer](https://github.com/donnemartin/system-design-primer)  

---

## 面试框架 (RASCAL)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      System Design Interview Framework                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  R - Requirements (需求)                                                      │
│     ├── 功能性需求: 系统应该做什么                                              │
│     └── 非功能性需求: 性能、可用性、扩展性                                        │
│                                                                              │
│  A - Architecture (架构)                                                      │
│     ├── 高层设计: API、数据流、组件                                             │
│     └── 技术选型: 数据库、缓存、队列                                            │
│                                                                              │
│  S - Scale (扩展)                                                             │
│     ├── 容量估算: QPS、存储、带宽                                               │
│     └── 扩展策略: 水平/垂直扩展、分片                                           │
│                                                                              │
│  C - Components (组件)                                                        │
│     ├── 详细设计: 每个组件的职责                                                │
│     └── 交互设计: 组件间通信方式                                                │
│                                                                              │
│  A - Algorithms (算法)                                                        │
│     ├── 核心算法: 推荐、排序、搜索                                              │
│     └── 优化策略: 缓存、预计算、批处理                                          │
│                                                                              │
│  L - Logistics (落地)                                                         │
│     ├── 监控告警: 可观测性                                                      │
│     └── 部署运维: CI/CD、容灾                                                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 经典系统设计案例

### 1. 短链接服务 (URL Shortener)

```
需求:
- 生成短链接
- 重定向到原链接
- 支持自定义别名
- 访问统计

流量估算:
- 写入: 1000 URLs/s
- 读取: 10M redirects/s (100:1 读写比)
- 存储: 100M URLs/year × 500B = 50GB/year

设计:
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────►│   API       │────►│   Redis     │
└─────────────┘     │   Gateway   │     │  (Cache)    │
                    └──────┬──────┘     └──────┬──────┘
                           │                     │
                           ▼                     │
                    ┌─────────────┐              │
                    │   App       │              │
                    │   Server    │              │
                    └──────┬──────┘              │
                           │                     │
                           ▼                     │
                    ┌─────────────┐              │
                    │   MySQL     │──────────────┘
                    │  (Master)   │
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │   MySQL     │
                    │  (Replica)  │
                    └─────────────┘

核心算法:
- Base62 编码 (a-zA-Z0-9)
- 自增 ID → 短码
- 冲突解决: 预生成、布隆过滤器
```

### 2. 分布式消息队列 (如 Kafka)

```
需求:
- 高吞吐消息发布/订阅
- 消息持久化
- 分区有序
- 消费者组

架构:
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Producer   │────►│   Broker    │────►│  Consumer   │
│   (App)     │     │   Cluster   │     │   (App)     │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
        ┌─────────┐  ┌─────────┐  ┌─────────┐
        │Partition│  │Partition│  │Partition│
        │   0     │  │   1     │  │   2     │
        │ (Leader)│  │ (Leader)│  │ (Leader)│
        │ (R1,R2) │  │ (R2,R0) │  │ (R0,R1) │
        └─────────┘  └─────────┘  └─────────┘

关键设计:
- 分区: 水平扩展、并行消费
- 副本: ISR (In-Sync Replicas)
- 存储: 日志结构、顺序写、零拷贝
- 协调: ZooKeeper/KRaft
```

### 3. 社交媒体 Feed (如 Twitter)

```
需求:
- 发帖
- 关注/取消关注
- 查看 Timeline (关注者的帖子)
- 点赞、评论

两种模型:

1. 拉模型 (Pull/Fan-out on read)
   - 发帖: 只写入自己的帖子表
   - 读 Timeline: 查询所有关注者的帖子，合并排序
   - 适合: 读少写多、关注数少

2. 推模型 (Push/Fan-out on write)
   - 发帖: 写入自己的帖子 + 推送到所有粉丝的 Timeline
   - 读 Timeline: 直接读取自己的 Timeline 表
   - 适合: 读多写少、关注数少

混合方案:
- 普通用户: Push 模型
- 大 V (>1M 粉丝): Pull 模型

架构:
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────►│   API       │────►│   Redis     │
└─────────────┘     │   Gateway   │     │  (Timeline  │
                    └──────┬──────┘     │   Cache)    │
                           │             └─────────────┘
                           ▼                    │
                    ┌─────────────┐              │
                    │   Post      │              │
                    │   Service   │              │
                    └──────┬──────┘              │
                           │                     │
              ┌────────────┴────────────┐        │
              ▼                         ▼        │
        ┌─────────────┐           ┌─────────────┐│
        │  Post DB    │           │  Fan-out    ││
        │  (Sharding) │           │  Service    │└┘
        └─────────────┘           └──────┬──────┘
                                          │
                                          ▼
                                    ┌─────────────┐
                                    │  Message    │
                                    │  Queue      │
                                    └─────────────┘
```

---

## 容量估算速查

```
常用数字:
- 1M DAU: 日活用户
- 10 RPS/用户: 平均请求率
- 100KB/请求: 平均请求大小
- 100ms: 目标延迟
- 99.9%: 目标可用性

计算示例:
DAU = 10M
请求/用户/天 = 100
峰值倍数 = 3

QPS = 10M × 100 / 86400 × 3 ≈ 35K
带宽 = 35K × 100KB = 3.5GB/s
存储/天 = 10M × 100 × 100KB = 100TB

数据库:
- MySQL: ~1K QPS/实例
- Redis: ~100K QPS/实例
- Cassandra: ~10K QPS/节点

需要实例数:
- MySQL: 35K / 1K = 35 主库 + 70 从库
- Redis: 35K / 100K = 1 集群
```

---

## 技术选型对比

| 场景 | 方案 A | 方案 B | 选择因素 |
|------|--------|--------|---------|
| 缓存 | Redis | Memcached | 数据结构 vs 性能 |
| 数据库 | PostgreSQL | MySQL | 复杂查询 vs 生态 |
| 搜索 | Elasticsearch | Solr | 易用性 vs 定制 |
| 队列 | Kafka | RabbitMQ | 吞吐 vs 路由 |
| 协调 | etcd | ZooKeeper | K8s 原生 vs 成熟 |
| 网关 | Envoy | Nginx | 动态配置 vs 性能 |

---

## 面试要点

### 应该做 ✅
- 先澄清需求
- 做容量估算
- 讨论权衡 (Trade-offs)
- 从简单开始逐步扩展
- 提及监控和运维

### 避免 ❌
- 立即深入细节
- 忽视非功能性需求
- 只有一种方案
- 忽视故障场景
- 过度设计

---

## 参考文献

1. [Designing Data-Intensive Applications](https://dataintensive.net/) - Martin Kleppmann
2. [System Design Primer](https://github.com/donnemartin/system-design-primer)
3. [System Design Interview](https://www.amazon.com/System-Design-Interview-insiders-Second/dp/B08CMF2CQF)
