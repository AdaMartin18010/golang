# TS-039-TiDB-8-5-Features

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: TiDB 8.5 LTS (2024-2025)
> **Size**: >20KB

---

## 1. TiDB 8.5 概览

### 1.1 发布信息

- **版本类型**: LTS (Long-Term Support)
- **发布时间**: 2024年12月
- **支持周期**: 长期支持版本
- **主要主题**: AI/ML、SaaS扩展性、多云

### 1.2 版本时间线

| 版本 | 发布日期 | 类型 |
|------|---------|------|
| 8.5.0 | 2024-12-19 | LTS |
| 8.5.1 | 2025-01-17 | Patch |
| 8.5.2 | 2025-06-12 | Patch |
| 8.5.3 | 2025-08-14 | Patch |
| 8.5.4 | 2025-11-27 | Patch |
| 8.5.5 | 2026-01-15 | Patch |

---

## 2. 向量搜索 (Preview)

### 2.1 概述

TiDB 8.5引入原生向量搜索支持，实现AI和语义搜索。

**优势**:

- 无需专用向量数据库
- 统一SQL接口
- 分布式向量存储

### 2.2 基本用法

```sql
-- 创建向量表
CREATE TABLE documents (
    id INT PRIMARY KEY,
    content TEXT,
    embedding VECTOR(1536)
);

-- 插入向量数据
INSERT INTO documents VALUES (
    1,
    'TiDB文档',
    '[0.1, 0.2, 0.3, ...]'::VECTOR
);

-- 向量相似度查询
SELECT id, content,
       COSINE_DISTANCE(embedding, :query_vector) AS distance
FROM documents
ORDER BY distance
LIMIT 10;
```

### 2.3 与全文搜索结合

```sql
-- 混合搜索: 向量 + 全文
SELECT * FROM products
WHERE MATCH(name, description) AGAINST('laptop')
ORDER BY COSINE_DISTANCE(embedding, :vec)
LIMIT 20;
```

---

## 3. TiDB Node Groups (GA)

### 3.1 概述

细粒度计算资源隔离，多租户/多工作负载场景优化。

**架构**:

```
┌─────────────────────────────────────────┐
│           TiDB Cluster                  │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────────┐    ┌─────────────┐    │
│  │ Node Group A│    │ Node Group B│    │
│  │ (OLTP)      │    │ (Analytics) │    │
│  │             │    │             │    │
│  │ TiDB × 3    │    │ TiDB × 2    │    │
│  │ TiKV × 6    │    │ TiKV × 4    │    │
│  └─────────────┘    └─────────────┘    │
│                                         │
│  - 资源隔离                             │
│  - 独立扩缩容                           │
│  - 统一管理                             │
└─────────────────────────────────────────┘
```

### 3.2 创建Node Group

```sql
-- 创建节点组
ALTER CLUSTER CREATE NODE GROUP oltp_group;
ALTER CLUSTER CREATE NODE GROUP analytics_group;

-- 分配节点
ALTER CLUSTER ASSIGN NODE 'tidb-1' TO GROUP oltp_group;
ALTER CLUSTER ASSIGN NODE 'tikv-1' TO GROUP oltp_group;

-- 放置规则
CREATE PLACEMENT POLICY oltp_policy
    CONSTRAINTS = '+group=oltp_group';

CREATE TABLE transactions (...) PLACEMENT POLICY oltp_policy;
```

### 3.3 应用场景

| 场景 | 配置 |
|------|------|
| 多租户SaaS | 每租户独立Node Group |
| 混合负载 | OLTP和OLAP分离 |
| 关键业务 | 独立资源保证SLA |

---

## 4. 全局索引 (Global Indexes)

### 4.1 概述

分区表上的全局索引，支持跨分区唯一约束。

### 4.2 使用示例

```sql
-- 创建分区表
CREATE TABLE orders (
    id BIGINT,
    user_id BIGINT,
    order_date DATE,
    PRIMARY KEY (id, order_date)
) PARTITION BY RANGE (YEAR(order_date)) (
    PARTITION p2023 VALUES LESS THAN (2024),
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026)
);

-- 创建全局索引
CREATE UNIQUE INDEX idx_user_id ON orders(user_id) GLOBAL;

-- 全局索引特点:
-- - 跨分区唯一
-- - 独立存储
-- - 点查优化
```

### 4.3 性能优势

- 点查无需扫描所有分区
- 跨分区唯一约束
- 减少查询延迟

---

## 5. Active PD Followers

### 5.1 概述

提升可用性，减少控制平面负载。

**架构改进**:

```
传统模式:
Client → PD Leader (所有请求)

Active Follower模式:
Client → PD Leader (写请求)
Client → PD Follower (读请求)
```

### 5.2 启用方式

```sql
-- 启用Active PD Follower
SET GLOBAL pd_enable_follower_handle_region = ON;

-- 或使用TiDB配置
[pd-client]
pd_enable_follower_handle_region = true
```

### 5.3 收益

- 减少PD Leader负载
- 提高Region信息获取性能
- 更好的可扩展性

---

## 6. TiProxy增强

### 6.1 智能路由

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────►│   TiProxy   │────►│ TiDB Nodes  │
└─────────────┘     │             │     │             │
                    │ - 负载均衡   │     │ - 健康检查   │
                    │ - 故障转移   │     │ - 读写分离   │
                    │ - 连接池     │     │             │
                    └─────────────┘     └─────────────┘
```

### 6.2 区域感知路由

```yaml
# TiProxy配置
proxy:
  backends:
    - host: tidb-region1-1
      region: us-west
    - host: tidb-region1-2
      region: us-west
    - host: tidb-region2-1
      region: us-east

  route-policy: region-aware  # 优先同区域
```

---

## 7. 计划缓存增强

### 7.1 实例级执行计划缓存

```sql
-- 启用实例级计划缓存
SET GLOBAL tidb_enable_instance_plan_cache = ON;

-- 缓存统计
SHOW STATUS LIKE '%plan_cache%';
```

### 7.2 非预处理语句缓存

```sql
-- 缓存非预处理语句的执行计划
SET GLOBAL tidb_enable_non_prepared_plan_cache = ON;

-- 减少重复查询编译开销
```

---

## 8. 工作负载控制

### 8.1 资源组管理

```sql
-- 创建资源组
CREATE RESOURCE GROUP oltp_group
    RU_PER_SEC = 10000
    PRIORITY = HIGH
    BURSTABLE = YES;

-- 绑定用户到资源组
ALTER USER 'app_user' RESOURCE GROUP oltp_group;

-- 查询级别资源限制
SELECT /*+ RESOURCE_GROUP(analytics_group) */ * FROM large_table;
```

### 8.2 Runaway Queries管理

```sql
-- 自动识别和管理失控查询
SET GLOBAL tidb_runaway_query_limit = 100;

-- 监控失控查询
SELECT * FROM information_schema.runaway_queries;
```

---

## 9. 并行查询增强

### 9.1 优化器改进

- 更好的并行度选择
- 动态并行度调整
- 减少数据传输

### 9.2 使用提示

```sql
-- 强制并行度
SELECT /*+ MAX_EXECUTION_TIME(10000) PARALLEL(4) */
    COUNT(*), AVG(amount)
FROM large_orders;
```

---

## 10. TiDB Cloud新功能

### 10.1 Auto Embedding (Beta)

```
功能: 自动将文本转换为向量
支持提供商:
- Amazon Titan
- OpenAI
- Cohere
- Gemini
- Jina AI
- Hugging Face
- NVIDIA NIM
```

### 10.2 列值过滤 (Changefeed)

```sql
-- 变更数据流中过滤列值
CREATE CHANGEFEED FOR TABLE users
INTO 'kafka://...'
WITH
    COLUMN_FILTER = 'id, name, email',  -- 仅同步指定列
    ROW_FILTER = 'region = "US"';        -- 仅同步符合条件的行
```

### 10.3 Standard Storage类型

**AWS专用**:

- Raft日志独立磁盘资源
- 减少I/O争用
- 提升读写性能
- 更具竞争力的价格

---

## 11. 监控和可观测性

### 11.1 Index Advisor

```sql
-- 获取索引建议
RECOMMEND INDEX FOR SELECT * FROM orders
WHERE user_id = 123 AND status = 'pending';

-- 应用建议索引
CREATE INDEX idx_orders_user_status
ON orders(user_id, status);
```

### 11.2 Schema Cache

- 减少Schema信息获取延迟
- 提升DDL操作性能
- 更好的并发处理

---

## 12. 升级建议

### 12.1 从8.1升级到8.5

```bash
# 1. 检查兼容性
tiup cluster check <cluster-name> --cluster-version v8.5.0

# 2. 备份数据
tiup br backup full --pd <pd-endpoint> --storage <s3-path>

# 3. 滚动升级
tiup cluster upgrade <cluster-name> v8.5.0
```

### 12.2 关键变更

| 特性 | 变更 |
|------|------|
| Index Insight | 移除，改用Index Advisor |
| TiFlash | Pipeline Model默认启用 |
| Runtime Filter | 默认启用 |

---

## 13. 最佳实践

### 13.1 SaaS多租户

```sql
-- 使用Node Group隔离租户
-- 使用Placement Rules控制数据分布

CREATE PLACEMENT POLICY tenant_a_policy
    CONSTRAINTS = '+zone=zone1,+disk=ssd'
    FOLLOWER_CONSTRAINTS = '+zone=zone2'
    LEADER_CONSTRAINTS = '+zone=zone1';

CREATE TABLE tenant_a_data (...) PLACEMENT POLICY tenant_a_policy;
```

### 13.2 AI/ML工作负载

```sql
-- 向量搜索 + 传统查询结合
SELECT p.*, COSINE_DISTANCE(embedding, :vec) as dist
FROM products p
WHERE category = 'electronics'
  AND price BETWEEN 100 AND 500
ORDER BY dist
LIMIT 10;
```

---

## 14. 参考文献

1. TiDB 8.5 Release Notes
2. TiDB Cloud Documentation
3. TiDB Best Practices
4. TiDB Vector Search Preview
5. TiDB Node Groups Guide

---

*Last Updated: 2026-04-03*
