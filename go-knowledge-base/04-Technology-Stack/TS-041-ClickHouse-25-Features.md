# TS-041-ClickHouse-25-Features

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: ClickHouse 25.x (25.1-25.8)
> **Size**: >20KB

---

## 1. ClickHouse 25.x 概览

### 1.1 版本时间线

| 版本 | 发布日期 | 类型 | 关键特性 |
|------|---------|------|---------|
| 25.1 | 2025-01 | 主要 | JSON类型、Query Condition Cache |
| 25.2 | 2025-02 | 主要 | Variant类型、Streaming |
| 25.3 | 2025-03 | 主要 | ACL、动态表 |
| 25.4 | 2025-04 | 主要 | 多集群共享 |
| 25.5 | 2025-05 | 主要 | OTLP Ingestion |
| 25.6 | 2025-06 | 主要 | 轻量物化投影 |
| 25.7 | 2025-07 | 主要 | Bzip3压缩 |
| 25.8 | 2025-08 | LTS | 长期支持版本 |

### 1.2 LTS版本

- **25.8 LTS**: 长期支持，企业推荐
- 向后移植修复到LTS版本
- 稳定性优于新功能

---

## 2. JSON数据类型 (Production Ready)

### 2.1 概述

ClickHouse 25.1+ 中 JSON 类型达到生产可用状态。

**特点**:

- 二进制列式存储
- 嵌套元素自动推断类型
- 可选schema定义
- 支持缺失字段

### 2.2 使用方式

```sql
-- 使用JSON类型
CREATE TABLE events (
    timestamp DateTime64(3),
    user_id UInt64,
    event_data JSON,  -- 动态JSON

    INDEX idx_event_type event_data.type TYPE bloom_filter GRANULARITY 3
) ENGINE = MergeTree()
ORDER BY timestamp;

-- 插入JSON数据
INSERT INTO events VALUES
(now(), 1, '{"type": "click", "target": "button", "value": 42}'),
(now(), 2, '{"type": "scroll", "depth": 1000}');

-- 查询JSON字段
SELECT
    event_data.type,
    count() as event_count
FROM events
WHERE event_data.type = 'click'
GROUP BY event_data.type;
```

### 2.3 带Schema的JSON

```sql
-- 定义JSON schema
CREATE TABLE structured_events (
    timestamp DateTime64(3),
    user_id UInt64,
    event_data JSON(
        type String,
        target String,
        value Nullable(Int64),
        depth Nullable(Int64)
    )
) ENGINE = MergeTree()
ORDER BY timestamp;
```

### 2.4 向后兼容性

```sql
-- 25.1之前的Object('json') 重命名为 JSON 类型
-- 支持从旧格式无缝迁移
```

---

## 3. Variant数据类型

### 3.1 概述

Variant类型允许列存储多种数据类型。

**应用场景**:

- 属性键值对 (不同值类型)
- 指标存储
- 半结构化数据

### 3.2 使用示例

```sql
-- 创建Variant列
CREATE TABLE metrics (
    timestamp DateTime,
    name String,
    value Variant(UInt64, Float64, String)
) ENGINE = MergeTree()
ORDER BY timestamp;

-- 插入不同类型数据
INSERT INTO metrics VALUES
(now(), 'cpu_usage', 45.5),      -- Float64
(now(), 'memory_bytes', 8589934592),  -- UInt64
(now(), 'status', 'healthy');     -- String

-- 查询时类型推断
SELECT
    name,
    value.Float64 as cpu_usage
FROM metrics
WHERE name = 'cpu_usage';
```

---

## 4. 动态表 (25.3+)

### 4.1 概述

动态表是ClickHouse的流处理抽象。

**特性**:

- 物化视图 + 流处理的结合
- 增量处理
- 时间窗口支持

### 4.2 创建动态表

```sql
-- 创建动态表 (增量物化视图)
CREATE DYNAMIC TABLE hourly_stats
AS SELECT
    toStartOfHour(timestamp) as hour,
    count() as event_count,
    uniqExact(user_id) as unique_users
FROM events
GROUP BY hour;

-- 时间窗口聚合
CREATE DYNAMIC TABLE moving_avg
AS SELECT
    timestamp,
    avg(value) OVER (ORDER BY timestamp RANGE BETWEEN INTERVAL '1' HOUR PRECEDING AND CURRENT ROW) as avg_value
FROM metrics;
```

---

## 5. ACL和安全增强 (25.3+)

### 5.1 行级安全 (RLS)

```sql
-- 启用RLS
ALTER TABLE sensitive_data ENABLE ROW LEVEL SECURITY;

-- 创建策略
CREATE POLICY tenant_isolation ON sensitive_data
    FOR SELECT
    USING tenant_id = currentUser();

-- 基于属性的访问控制
CREATE POLICY region_access ON sales_data
    FOR ALL
    USING region IN (SELECT allowed_regions FROM user_permissions WHERE user = currentUser());
```

### 5.2 列级加密

```sql
-- 列级加密
CREATE TABLE encrypted_data (
    id UInt64,
    sensitive_column String ENCRYPTED WITH KEY 'my_key',
    public_column String
) ENGINE = MergeTree()
ORDER BY id;
```

---

## 6. 流处理增强

### 6.1 时间窗口

```sql
-- 滚动窗口
SELECT
    tumbleStart(timestamp, INTERVAL '5' MINUTE) as window_start,
    count() as event_count
FROM stream_data
GROUP BY tumbleStart(timestamp, INTERVAL '5' MINUTE);

-- 滑动窗口
SELECT
    hopStart(timestamp, INTERVAL '1' MINUTE, INTERVAL '5' MINUTE) as window_start,
    avg(value) as avg_value
FROM sensor_data
GROUP BY hopStart(timestamp, INTERVAL '1' MINUTE, INTERVAL '5' MINUTE);

-- 会话窗口
SELECT
    sessionStart(timestamp, INTERVAL '30' MINUTE) as session_start,
    count() as events_in_session
FROM user_activity
GROUP BY sessionStart(timestamp, INTERVAL '30' MINUTE);
```

### 6.2 Watermark处理

```sql
-- 设置Watermark延迟容忍
CREATE STREAM sensor_stream (
    timestamp DateTime64(3) WATERMARK timestamp - INTERVAL '10' SECOND,
    sensor_id String,
    value Float64
);
```

---

## 7. OTLP Ingestion (25.5+)

### 7.1 OpenTelemetry原生支持

```sql
-- 自动创建OTLP表
SYSTEM CREATE OTLP TABLES;

-- 生成端点URL
SELECT * FROM system.build_options
WHERE name = 'OTLP_ENDPOINT';

-- OTLP gRPC端点: http://<host>:9363/v1/traces
-- OTLP HTTP端点: http://<host>:9363/v1/metrics
```

### 7.2 配置

```xml
<!-- config.xml -->
<opentelemetry>
    <enabled>true</enabled>
    <http_port>9363</http_port>
    <grpc_port>9364</grpc_port>
    <metrics_table>otel_metrics</metrics_table>
    <logs_table>otel_logs</logs_table>
    <traces_table>otel_traces</traces_table>
</opentelemetry>
```

---

## 8. 查询性能优化

### 8.1 Query Condition Cache

```sql
-- 查询条件缓存加速重复过滤
SET use_query_condition_cache = 1;

-- 自动缓存常用过滤条件
-- 显著加速仪表板类查询
```

### 8.2 轻量物化投影 (25.6+)

```sql
-- 轻量投影 (不存储数据，仅元数据)
CREATE TABLE events (
    timestamp DateTime,
    user_id UInt64,
    event_type String,
    value Float64,
    PROJECTION event_summary
    (
        SELECT
            toStartOfHour(timestamp),
            event_type,
            count(),
            avg(value)
        GROUP BY toStartOfHour(timestamp), event_type
    )
) ENGINE = MergeTree()
ORDER BY timestamp;

-- 轻量投影特点:
-- - 实时计算
-- - 不占用额外存储
-- - 适合频繁变化的聚合
```

### 8.3 投影增强

```sql
-- 正常投影 (存储预聚合数据)
CREATE PROJECTION monthly_sales
ON TABLE transactions
(
    SELECT
        toYYYYMM(transaction_date) as month,
        product_category,
        sum(amount),
        count()
    GROUP BY month, product_category
);

-- 查询自动使用投影
SELECT
    toYYYYMM(transaction_date),
    sum(amount)
FROM transactions
GROUP BY toYYYYMM(transaction_date);
-- 自动路由到投影数据
```

---

## 9. 压缩优化

### 9.1 Bzip3压缩 (25.7+)

```sql
-- 使用Bzip3压缩
CREATE TABLE compressed_data (
    id UInt64,
    content String CODEC(Bzip3)
) ENGINE = MergeTree()
ORDER BY id
SETTINGS compression = 'bzip3';

-- Bzip3特点:
-- - 高压缩比
-- - 适合冷数据
-- - 解压缩较慢
```

### 9.2 压缩算法对比

| 算法 | 压缩比 | 速度 | 适用场景 |
|------|--------|------|---------|
| LZ4 | 低 | 极快 | 热数据 |
| ZSTD | 中 | 快 | 通用 |
| Bzip3 | 高 | 慢 | 冷数据/归档 |

---

## 10. 多集群共享 (25.4+)

### 10.1 跨集群查询

```sql
-- 查询远程集群
SELECT * FROM remote('cluster_b', default.events)
WHERE timestamp > now() - INTERVAL 1 HOUR;

-- 分布式JOIN
SELECT
    a.user_id,
    a.event_count,
    b.user_profile
FROM clusterA.events a
JOIN clusterB.users b ON a.user_id = b.id;
```

### 10.2 集群联邦

```xml
<remote_servers>
    <federation>
        <shard>
            <replica>
                <host>cluster-a-node-1</host>
                <port>9000</port>
            </replica>
        </shard>
        <shard>
            <replica>
                <host>cluster-b-node-1</host>
                <port>9000</port>
            </replica>
        </shard>
    </federation>
</remote_servers>
```

---

## 11. 可观测性增强

### 11.1 EXPLAIN PIPELINE

```sql
-- 查看查询执行管道
EXPLAIN PIPELINE
SELECT count()
FROM events
WHERE timestamp > now() - INTERVAL 1 DAY;

-- 输出管道阶段和算子
```

### 11.2 查询画像

```sql
-- 详细查询分析
EXPLAIN ANALYZE
SELECT
    user_id,
    count() as event_count
FROM events
GROUP BY user_id
ORDER BY event_count DESC
LIMIT 10;

-- 显示:
-- - 实际执行时间
-- - 行数统计
-- - 内存使用
-- - 算子成本
```

### 11.3 系统表增强

```sql
-- 查询缓存统计
SELECT
    query,
    read_rows,
    result_bytes,
    query_cache_usage
FROM system.query_log
WHERE event_time > now() - INTERVAL 1 HOUR;

-- 查看物化视图状态
SELECT
    database,
    table,
    last_successful_update,
    last_error
FROM system.materialized_views;
```

---

## 12. 最佳实践

### 12.1 JSON数据建模

```sql
-- 推荐: 常用字段提取为物理列
CREATE TABLE events_optimized (
    timestamp DateTime64(3),
    user_id UInt64,
    event_type LowCardinality(String),  -- 提取为列
    event_data JSON,  -- 动态内容保留

    INDEX idx_type event_type TYPE bloom_filter,
    INDEX idx_user user_id TYPE minmax
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (user_id, timestamp);
```

### 12.2 流处理模式

```sql
-- 实时聚合 + 历史查询结合
CREATE TABLE raw_events (...) ENGINE = MergeTree();

CREATE MATERIALIZED VIEW events_5min_agg
ENGINE = SummingMergeTree()
AS SELECT
    toStartOfFiveMinute(timestamp) as bucket,
    event_type,
    count() as event_count
FROM raw_events
GROUP BY bucket, event_type;

-- 查询时合并
SELECT
    bucket,
    sum(event_count)
FROM events_5min_agg
WHERE bucket > now() - INTERVAL 1 DAY
GROUP BY bucket;
```

---

## 13. 升级指南

### 13.1 升级到25.x

```bash
# 1. 备份元数据
clickhouse-backup create metadata_backup

# 2. 滚动升级 (逐个副本)
systemctl stop clickhouse-server
# 升级包
systemctl start clickhouse-server

# 3. 验证
SELECT version();
```

### 13.2 从Object('json')迁移

```sql
-- 自动迁移
ALTER TABLE old_table
MODIFY COLUMN json_column JSON;

-- 或使用MATERIALIZE
ALTER TABLE old_table
MATERIALIZE COLUMN json_column;
```

---

## 14. 参考文献

1. ClickHouse 25.x Release Notes
2. ClickHouse JSON Documentation
3. ClickHouse Streaming Guide
4. ClickHouse Performance Tuning
5. ClickHouse Security Best Practices

---

*Last Updated: 2026-04-03*
