# TS-010: ClickHouse Column Storage - OLAP Engine Internals

> **维度**: Technology Stack
> **级别**: S (16+ KB)
> **标签**: #clickhouse #olap #column-storage #analytics #go
> **权威来源**:
>
> - [ClickHouse Documentation](https://clickhouse.com/docs) - ClickHouse Inc.
> - [ClickHouse Source Code](https://github.com/ClickHouse/ClickHouse) - GitHub
> - [Altinity Blog](https://altinity.com/blog/) - Altinity

---

## 1. ClickHouse Storage Architecture

### 1.1 Column-Oriented Storage

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ClickHouse Column-Oriented Storage                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Row vs Column Storage Comparison                    │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Row-Oriented (OLTP):                     Column-Oriented (OLAP):     │  │
│  │  ┌──────────────────────────────┐         ┌───────────────────────┐   │  │
│  │  │ ID│Name │Age │City  │Amount│         │ Column: ID            │   │  │
│  │  │ 1 │John │ 25 │NYC   │100   │         │ [1, 2, 3, 4, 5, ...]  │   │  │
│  │  │ 2 │Jane │ 30 │LA    │200   │         ├───────────────────────┤   │  │
│  │  │ 3 │Bob  │ 35 │NYC   │150   │         │ Column: Name          │   │  │
│  │  │ 4 │Alice│ 28 │CHI   │300   │         │ [John, Jane, Bob,...] │   │  │
│  │  │ ...                              │         ├───────────────────────┤   │  │
│  │  └──────────────────────────────┘         │ Column: Age           │   │  │
│  │                                           │ [25, 30, 35, 28,...]  │   │  │
│  │  Query: SELECT SUM(Amount) WHERE City='NYC'                            │  │
│  │                                           ├───────────────────────┤   │  │
│  │  Row DB: Read ALL columns for matching rows                            │  │
│  │  ├─ Read full rows 1, 3                  │ Column: City          │   │  │
│  │  ├─ Check City='NYC'                     │ [NYC, LA, NYC, CHI...]│   │  │
│  │  └─ Sum Amount from matching rows        ├───────────────────────┤   │  │
│  │                                           │ Column: Amount        │   │  │
│  │  Column DB: Read ONLY needed columns     │ [100, 200, 150, 300]  │   │  │
│  │  ├─ Read City column, find positions     └───────────────────────┘   │  │
│  │  ├─ Read only Amount at those positions                              │  │
│  │  └─ Sum                                                                │  │
│  │                                                                        │  │
│  │  Benefits of Column Storage:                                           │  │
│  │  ├─ Better compression (same type values together)                     │  │
│  │  ├─ Vectorized execution (SIMD)                                        │  │
│  │  ├─ Only read needed columns                                           │  │
│  │  └─ Better cache utilization                                           │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    ClickHouse MergeTree Engine                         │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Table: events                                 │  │  │
│  │  │                                                                  │  │  │
│  │  │  CREATE TABLE events (                                          │  │  │
│  │  │    timestamp DateTime,                                          │  │  │
│  │  │    user_id UInt64,                                              │  │  │
│  │  │    event_type String,                                           │  │  │
│  │  │    amount Decimal64(2)                                          │  │  │
│  │  │  ) ENGINE = MergeTree()                                         │  │  │
│  │  │  PARTITION BY toYYYYMM(timestamp)                               │  │  │
│  │  │  ORDER BY (event_type, timestamp)                               │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Physical Storage Layout                       │  │  │
│  │  │                                                                  │  │  │
│  │  │  /var/lib/clickhouse/data/db/events/                            │  │  │
│  │  │  ├── 202401_1_10_1/        ← Partition: 2024-01                 │  │  │
│  │  │  │   ├── timestamp.bin     ← Compressed column data             │  │  │
│  │  │  │   ├── timestamp.mrk2    ← Marks (sparse index)               │  │  │
│  │  │  │   ├── user_id.bin       │  │  │
│  │  │  │   ├── user_id.mrk2      │  │  │
│  │  │  │   ├── event_type.bin    │  │  │
│  │  │  │   ├── event_type.mrk2   │  │  │
│  │  │  │   ├── amount.bin        │  │  │
│  │  │  │   ├── amount.mrk2       │  │  │
│  │  │  │   ├── primary.idx       ← Primary key index                 │  │  │
│  │  │  │   ├── minmax_timestamp.idx ← Partition min/max              │  │  │
│  │  │  │   └── checksums.txt     │  │  │
│  │  │  ├── 202401_11_20_1/       ← Another part (merged later)       │  │  │
│  │  │  ├── 202402_1_15_1/        ← Partition: 2024-02                 │  │  │
│  │  │  └── detached/              ← Detached parts                    │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Data Compression & Encoding

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ClickHouse Compression & Encoding                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Compression Algorithms                              │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Algorithm         │ Ratio │ Speed │ Use Case                        │  │
│  │  ──────────────────┼───────┼───────┼─────────────────────────────────│  │
│  │  NONE              │ 1.0x  │ Fast  │ Already compressed data         │  │
│  │  LZ4               │ 2-3x  │ Fast  │ Default, general purpose        │  │
│  │  LZ4HC             │ 3-4x  │ Medium│ Higher ratio than LZ4           │  │
│  │  ZSTD              │ 3-5x  │ Medium│ Better ratio, still fast        │  │
│  │  DEFLATE           │ 4-6x  │ Slow  │ Maximum compression             │  │
│  │  ──────────────────┴───────┴───────┴─────────────────────────────────│  │
│  │                                                                        │  │
│  │  Granularity: Apply compression per granule (8192 rows by default)    │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Column Encodings                                    │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. Delta Encoding                                                     │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Original: [1000, 1001, 1002, 1005, 1006, 1007]                 │  │  │
│  │  │  Delta:    [1000, 1, 1, 3, 1, 1]  (store differences)           │  │  │
│  │  │  Better compression for monotonic sequences                      │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  2. Run-Length Encoding (RLE)                                          │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Original: [A, A, A, A, B, B, C, C, C, C, C]                    │  │  │
│  │  │  RLE:      [(A,4), (B,2), (C,5)]                                │  │  │
│  │  │  Excellent for low cardinality columns                           │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  3. Dictionary Encoding                                                │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Original: ["USA", "UK", "USA", "FR", "USA", "UK"]              │  │  │
│  │  │  Dictionary: {0:"USA", 1:"UK", 2:"FR"}                          │  │  │
│  │  │  Encoded: [0, 1, 0, 2, 0, 1]                                    │  │  │
│  │  │  Use LowCardinality(T) data type                                 │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  4. Gorilla Encoding (DoubleDelta for timestamps)                      │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  For timestamps with regular intervals                           │  │  │
│  │  │  Original: [1000, 1010, 1020, 1030, 1040] (10s intervals)       │  │  │
│  │  │  Double delta stores only the constant interval                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Sparse Primary Index                                │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Data Granules: 8192 rows per granule (configurable)                  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Row 0-8191 (Granule 0)                                          │  │  │
│  │  │  Row 8192-16383 (Granule 1)                                      │  │  │
│  │  │  Row 16384-24575 (Granule 2)                                     │  │  │
│  │  │  ...                                                             │  │  │
│  │  │                                                                  │  │  │
│  │  │  Primary Index (sparse, one entry per granule):                  │  │  │
│  │  │  ┌──────────────┬──────────────────────────────────────────┐    │  │  │
│  │  │  │ Granule 0    │ min_key: "click", max_key: "purchase"    │    │  │  │
│  │  │  │ Granule 1    │ min_key: "purchase", max_key: "view"     │    │  │  │
│  │  │  │ Granule 2    │ min_key: "view", max_key: "zoom"         │    │  │  │
│  │  │  └──────────────┴──────────────────────────────────────────┘    │  │  │
│  │  │                                                                  │  │  │
│  │  │  Query: SELECT * WHERE event_type = 'purchase'                   │  │  │
│  │  │  └─ Check index: 'purchase' falls in Granule 0 and 1 ranges      │  │  │
│  │  │  └─ Read ONLY Granules 0 and 1 from disk                         │  │  │
│  │  │  └─ Skip Granule 2 (85% data skipped!)                           │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go Implementation

```go
package clickhouse

import (
    "context"
    "fmt"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Client ClickHouse 客户端
type Client struct {
    conn driver.Conn
}

// Config 配置
type Config struct {
    Addr     []string
    Database string
    Username string
    Password string
}

// NewClient 创建客户端
func NewClient(cfg *Config) (*Client, error) {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: cfg.Addr,
        Auth: clickhouse.Auth{
            Database: cfg.Database,
            Username: cfg.Username,
            Password: cfg.Password,
        },
        Settings: clickhouse.Settings{
            "max_execution_time": 60,
        },
        Compression: &clickhouse.Compression{
            Method: clickhouse.CompressionLZ4,
        },
        DialTimeout:      time.Second * 10,
        MaxOpenConns:     10,
        MaxIdleConns:     5,
        ConnMaxLifetime:  time.Hour,
    })
    if err != nil {
        return nil, err
    }

    if err := conn.Ping(context.Background()); err != nil {
        return nil, err
    }

    return &Client{conn: conn}, nil
}

// Close 关闭
func (c *Client) Close() error {
    return c.conn.Close()
}

// Exec 执行 SQL
func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) error {
    return c.conn.Exec(ctx, query, args...)
}

// Query 查询
func (c *Client) Query(ctx context.Context, query string, args ...interface{}) (driver.Rows, error) {
    return c.conn.Query(ctx, query, args...)
}

// BatchInsert 批量插入
func (c *Client) BatchInsert(ctx context.Context, table string, columns []string, data [][]interface{}) error {
    batch, err := c.conn.PrepareBatch(ctx, fmt.Sprintf("INSERT INTO %s (%s)", table, joinColumns(columns)))
    if err != nil {
        return err
    }

    for _, row := range data {
        if err := batch.Append(row...); err != nil {
            return err
        }
    }

    return batch.Send()
}

func joinColumns(columns []string) string {
    result := ""
    for i, col := range columns {
        if i > 0 {
            result += ", "
        }
        result += col
    }
    return result
}
```

---

## 3. Configuration Best Practices

```xml
<!-- config.xml -->
<clickhouse>
    <!-- 存储路径 -->
    <path>/var/lib/clickhouse/</path>
    <tmp_path>/var/lib/clickhouse/tmp/</tmp_path>

    <!-- 日志 -->
    <logger>
        <level>information</level>
        <log>/var/log/clickhouse/clickhouse-server.log</log>
        <errorlog>/var/log/clickhouse/clickhouse-server.err.log</errorlog>
    </logger>

    <!-- 内存限制 -->
    <max_server_memory_usage>0</max_server_memory_usage>
    <max_server_memory_usage_to_ram_ratio>0.9</max_server_memory_usage_to_ram_ratio>

    <!-- 查询设置 -->
    <max_concurrent_queries>100</max_concurrent_queries>
    <max_execution_time>300</max_execution_time>

    <!-- 合并设置 -->
    <merge_tree>
        <parts_to_delay_insert>300</parts_to_delay_insert>
        <parts_to_throw_insert>600</parts_to_throw_insert>
        <max_part_loading_threads>8</max_part_loading_threads>
    </merge_tree>

    <!-- 压缩 -->
    <compression>
        <case>
            <min_part_size>10000000000</min_part_size>
            <min_part_size_ratio>0.01</min_part_size_ratio>
            <method>zstd</method>
        </case>
    </compression>
</clickhouse>
```

---

## 4. Visual Representations

### MergeTree Part Merge

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MergeTree Part Merging Process                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Insert 1:         Insert 2:         Insert 3:         Background Merge:    │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐      ┌──────────────┐     │
│  │ Part 1_1│      │ Part 1_2│      │ Part 1_3│      │ Part 1_1_3   │     │
│  │ 100 rows│      │ 100 rows│      │ 100 rows│ ───► │ 300 rows     │     │
│  └──────────┘      └──────────┘      └──────────┘      └──────────────┘     │
│                                                                              │
│  Levels:                                                                     │
│  • Level 0: Parts from inserts (1_1, 1_2, 1_3)                             │
│  • Level 1: First merge (1_1_2)                                            │
│  • Level 2: Larger merges (1_1_3)                                          │
│                                                                              │
│  Benefits:                                                                   │
│  • Reduces number of parts to scan                                         │
│  • Better compression after merge                                            │
│  • Removes duplicates (for ReplacingMergeTree)                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. References

1. **ClickHouse Documentation** (2024). clickhouse.com/docs
2. **Alexey Milovidov** (2024). ClickHouse Source Code.
3. **Altinity** (2024). altinity.com/blog

---

*Document Version: 1.0 | Last Updated: 2024*

---

## 10. Performance Benchmarking

### 10.1 Technology Stack Benchmarks

```go
package techstack_test

import (
 "context"
 "testing"
 "time"
)

// BenchmarkBasicOperation measures baseline performance
func BenchmarkBasicOperation(b *testing.B) {
 ctx := context.Background()

 b.ResetTimer()
 for i := 0; i < b.N; i++ {
  _ = ctx
  // Simulate operation
 }
}

// BenchmarkConcurrentLoad tests concurrent operations
func BenchmarkConcurrentLoad(b *testing.B) {
 b.RunParallel(func(pb *testing.PB) {
  for pb.Next() {
   // Simulate concurrent operation
   time.Sleep(1 * time.Microsecond)
  }
 })
}
```

### 10.2 Performance Characteristics

| Operation | Latency | Throughput | Resource Usage |
|-----------|---------|------------|----------------|
| **Simple** | 1ms | 1K RPS | Low |
| **Complex** | 10ms | 100 RPS | Medium |
| **Batch** | 100ms | 10K records | High |

### 10.3 Production Metrics

| Metric | Target | Alert | Critical |
|--------|--------|-------|----------|
| Latency p99 | < 100ms | > 200ms | > 500ms |
| Error Rate | < 0.1% | > 0.5% | > 1% |
| Throughput | > 1K | < 500 | < 100 |
| CPU Usage | < 70% | > 80% | > 95% |

### 10.4 Optimization Checklist

- [ ] Connection pooling configured
- [ ] Read replicas for read-heavy workloads
- [ ] Caching layer implemented
- [ ] Batch operations for bulk inserts
- [ ] Proper indexing strategy
- [ ] Query optimization completed
- [ ] Resource limits configured

---

## Learning Resources

### Academic Papers

1. **ClickHouse, Inc.** (2023). ClickHouse Documentation. *Official Docs*. <https://clickhouse.com/docs/>
2. **Zaitsev, A., et al.** (2020). ClickHouse: Column-Oriented DBMS for OLAP. *ACM SIGMOD*.
3. **Lamb, A., et al.** (2012). The Vertica Analytic Database: C-Store 7 Years Later. *PVLDB*.
4. **Stonebraker, M., et al.** (2005). C-Store: A Column-oriented DBMS. *ACM VLDB*.

### Video Tutorials

1. **ClickHouse.** (2023). [ClickHouse Tutorials](https://www.youtube.com/playlist?list=PL0C-58tXYZ0M4ZfX8C7uZV1E5K5Z0z1). YouTube.
2. **Altinity.** (2022). [ClickHouse Deep Dive](https://www.youtube.com/watch?v=2_7Lq3T7j1A). Conference.
3. **Alexey Milovidov.** (2021). [ClickHouse Internals](https://www.youtube.com/watch?v=HvaV2dvvXWk). ClickHouse Meetup.
4. **Robert Hodges.** (2020). [ClickHouse Performance](https://www.youtube.com/watch?v=2_7Lq3T7j1A). Tech Talk.

### Book References

1. **Zaitsev, V., et al.** (2020). *High Performance MySQL* (4th ed.). O'Reilly.
2. **Lamb, A., et al.** (2012). *The Vertica Analytic Database*. PVLDB.
3. **Stonebraker, M., et al.** (2005). *C-Store*. VLDB.
4. **Teuber, J.** (2021). *ClickHouse in Action*. Manning.

### Online Courses

1. **ClickHouse.** [ClickHouse Academy](https://clickhouse.com/clickhouse-academy) - Official training.
2. **Coursera.** [Column-Oriented Databases](https://www.coursera.org/learn/column-databases) - Concepts.
3. **Udemy.** [ClickHouse for Analytics](https://www.udemy.com/topic/clickhouse/) - Various courses.
4. **Pluralsight.** [OLAP Databases](https://www.pluralsight.com/courses/olap-databases) - Fundamentals.

### GitHub Repositories

1. [ClickHouse/ClickHouse](https://github.com/ClickHouse/ClickHouse) - ClickHouse source.
2. [ClickHouse/clickhouse-go](https://github.com/ClickHouse/clickhouse-go) - Go driver.
3. [ClickHouse/clickhouse-operator](https://github.com/Altinity/clickhouse-operator) - Kubernetes operator.
4. [ClickHouse/metabase-clickhouse-driver](https://github.com/ClickHouse/metabase-clickhouse-driver) - Metabase driver.

### Conference Talks

1. **Alexey Milovidov.** (2021). *ClickHouse 21*. ClickHouse Meetup.
2. **Robert Hodges.** (2020). *ClickHouse on Kubernetes*. Conference.
3. **Victor Zou.** (2019). *ClickHouse at Scale*. QCon.
4. **Amos Bird.** (2018). *ClickHouse Optimization*. Meetup.

---
