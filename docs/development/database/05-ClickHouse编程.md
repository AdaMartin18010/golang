# ClickHouse编程 - Go语言实战指南

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---
## 📋 目录

- [ClickHouse编程 - Go语言实战指南](#clickhouse编程-go语言实战指南)
  - [ClickHouse概述](#clickhouse概述)
  - [安装与配置](#安装与配置)
- [ClickHouse官方Go驱动](#clickhouse官方go驱动)
- [标准database/sql接口](#标准databasesql接口)
- [Docker方式（推荐）](#docker方式推荐)
- [测试连接](#测试连接)
  - [连接管理](#连接管理)
  - [表设计与操作](#表设计与操作)
  - [数据插入](#数据插入)
  - [数据查询](#数据查询)
  - [物化视图](#物化视图)
  - [分区与分片](#分区与分片)
  - [实时数据分析](#实时数据分析)
  - [性能优化](#性能优化)
  - [监控与运维](#监控与运维)
  - [最佳实践](#最佳实践)
  - [总结](#总结)

---

## ClickHouse概述

### 特点

ClickHouse是一个**OLAP**（在线分析处理）数据库：

- **列式存储**: 适合分析查询，压缩率高
- **高性能**: 单服务器每秒处理数亿行
- **SQL支持**: 完整的SQL方言
- **实时插入**: 支持实时数据流
- **分布式**: 自动数据分片和副本
- **丰富的数据类型**: 数组、嵌套、JSON等

### 适用场景

```text
✅ 日志分析
✅ 用户行为分析
✅ 实时监控
✅ 时序数据存储
✅ 商业智能(BI)
✅ 数据仓库
```

### OLTP vs OLAP

```text
┌──────────────┬────────────┬──────────────┐
│   特性       │   OLTP     │    OLAP      │
├──────────────┼────────────┼──────────────┤
│ 数据量       │ GB级       │ TB-PB级      │
│ 查询类型     │ 点查询     │ 聚合查询     │
│ 写入模式     │ 随机写     │ 批量写       │
│ 典型数据库   │ MySQL/PG   │ ClickHouse   │
└──────────────┴────────────┴──────────────┘
```

---

## 安装与配置

### 安装驱动

```bash
# ClickHouse官方Go驱动
go get github.com/ClickHouse/clickhouse-go/v2

# 标准database/sql接口
go get github.com/ClickHouse/clickhouse-go/v2/lib/driver
```

### 本地安装ClickHouse

```bash
# Docker方式（推荐）
docker run -d \
  --name clickhouse-server \
  --ulimit nofile=262144:262144 \
  -p 9000:9000 \
  -p 8123:8123 \
  clickhouse/clickhouse-server

# 测试连接
curl http://localhost:8123/ping
```

---

## 连接管理

### 原生协议连接

```go
package database

import (
    "Context"
    "crypto/tls"
    "fmt"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// ClickHouse 数据库实例
type ClickHouse struct {
    Conn driver.Conn
}

// NewClickHouse 创建ClickHouse连接（原生协议）
func NewClickHouse(addr string) (*ClickHouse, error) {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{addr}, // 支持多个地址用于负载均衡
        Auth: clickhouse.Auth{
            Database: "default",
            Username: "default",
            Password: "",
        },
        Settings: clickhouse.Settings{
            "max_execution_time": 60, // 最大执行时间(秒)
        },
        DialTimeout:      5 * time.Second,
        MaxOpenConns:     10,
        MaxIdleConns:     5,
        ConnMaxLifetime:  time.Hour,
        ConnOpenStrategy: clickhouse.ConnOpenInOrder,

        // TLS配置（可选）
        TLS: &tls.Config{
            InsecureSkipVerify: true,
        },

        // 调试模式
        Debug: false,
    })

    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }

    // 测试连接
    ctx := Context.Background()
    if err := conn.Ping(ctx); err != nil {
        return nil, fmt.Errorf("ping failed: %w", err)
    }

    fmt.Println("✅ ClickHouse连接成功")

    return &ClickHouse{Conn: conn}, nil
}

// Close 关闭连接
func (c *ClickHouse) Close() error {
    return c.Conn.Close()
}

// Exec 执行SQL
func (c *ClickHouse) Exec(ctx Context.Context, query string, args ...interface{}) error {
    return c.Conn.Exec(ctx, query, args...)
}

// Query 查询
func (c *ClickHouse) Query(ctx Context.Context, dest interface{}, query string, args ...interface{}) error {
    rows, err := c.Conn.Query(ctx, query, args...)
    if err != nil {
        return err
    }
    defer rows.Close()

    return rows.ScanStruct(dest)
}
```

### HTTP协议连接

```go
import (
    "database/sql"
    _ "github.com/ClickHouse/clickhouse-go/v2"
)

// NewClickHouseHTTP 创建ClickHouse连接（HTTP协议）
func NewClickHouseHTTP(dsn string) (*sql.DB, error) {
    // DSN格式: clickhouse://username:password@host:port/database?param=value
    db, err := sql.Open("clickhouse", dsn)
    if err != nil {
        return nil, err
    }

    // 设置连接池
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)

    // 测试连接
    if err := db.Ping(); err != nil {
        return nil, err
    }

    fmt.Println("✅ ClickHouse HTTP连接成功")

    return db, nil
}
```

---

## 表设计与操作

### 表引擎类型

ClickHouse支持多种表引擎：

- **MergeTree**: 最常用，支持索引和分区
- **ReplacingMergeTree**: 去重
- **SummingMergeTree**: 自动求和聚合
- **AggregatingMergeTree**: 聚合函数结果
- **CollapsingMergeTree**: 状态折叠
- **Distributed**: 分布式表

### 创建表

```go
package repository

import (
    "Context"
    "fmt"

    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// EventRepository 事件仓储
type EventRepository struct {
    conn driver.Conn
}

// NewEventRepository 创建事件仓储
func NewEventRepository(conn driver.Conn) *EventRepository {
    return &EventRepository{conn: conn}
}

// CreateTable 创建事件表
func (r *EventRepository) CreateTable(ctx Context.Context) error {
    query := `
    CREATE TABLE IF NOT EXISTS events (
        event_id String,
        user_id UInt64,
        event_type String,
        event_time DateTime,
        properties String,  -- JSON格式
        country String,
        city String,
        device String,
        os String,
        browser String,
        session_id String,
        created_at DateTime DEFAULT now()
    )
    ENGINE = MergeTree()
    PARTITION BY toYYYYMM(event_time)    -- 按月分区
    ORDER BY (event_type, user_id, event_time)  -- 排序键（主键）
    TTL event_time + INTERVAL 90 DAY     -- 数据保留90天
    SETTINGS index_granularity = 8192
    `

    err := r.conn.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("create table failed: %w", err)
    }

    fmt.Println("✅ 表创建成功")

    return nil
}

// CreateDistributedTable 创建分布式表
func (r *EventRepository) CreateDistributedTable(ctx Context.Context, cluster string) error {
    query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS events_distributed AS events
    ENGINE = Distributed(%s, default, events, rand())
    `, cluster)

    return r.conn.Exec(ctx, query)
}

// DropTable 删除表
func (r *EventRepository) DropTable(ctx Context.Context) error {
    return r.conn.Exec(ctx, "DROP TABLE IF EXISTS events")
}

// OptimizeTable 优化表（手动触发合并）
func (r *EventRepository) OptimizeTable(ctx Context.Context) error {
    return r.conn.Exec(ctx, "OPTIMIZE TABLE events FINAL")
}
```

---

## 数据插入

### 批量插入

```go
package models

import "time"

// Event 事件模型
type Event struct {
    EventID    string    `ch:"event_id"`
    UserID     uint64    `ch:"user_id"`
    EventType  string    `ch:"event_type"`
    EventTime  time.Time `ch:"event_time"`
    Properties string    `ch:"properties"`
    Country    string    `ch:"country"`
    City       string    `ch:"city"`
    Device     string    `ch:"device"`
    OS         string    `ch:"os"`
    Browser    string    `ch:"browser"`
    SessionID  string    `ch:"session_id"`
    CreatedAt  time.Time `ch:"created_at"`
}
```

```go
// Insert 插入单条事件
func (r *EventRepository) Insert(ctx Context.Context, event *Event) error {
    query := `
    INSERT INTO events (
        event_id, user_id, event_type, event_time, properties,
        country, city, device, os, browser, session_id
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

    err := r.conn.Exec(ctx, query,
        event.EventID, event.UserID, event.EventType, event.EventTime,
        event.Properties, event.Country, event.City, event.Device,
        event.OS, event.Browser, event.SessionID,
    )

    return err
}

// BatchInsert 批量插入（推荐）
func (r *EventRepository) BatchInsert(ctx Context.Context, events []*Event) error {
    // 使用Batch接口，性能更好
    batch, err := r.conn.PrepareBatch(ctx, "INSERT INTO events")
    if err != nil {
        return err
    }

    for _, event := range events {
        err := batch.Append(
            event.EventID, event.UserID, event.EventType, event.EventTime,
            event.Properties, event.Country, event.City, event.Device,
            event.OS, event.Browser, event.SessionID,
        )
        if err != nil {
            return err
        }
    }

    // 发送批次
    if err := batch.Send(); err != nil {
        return fmt.Errorf("batch send failed: %w", err)
    }

    fmt.Printf("✅ 批量插入成功，插入了 %d 条记录\n", len(events))

    return nil
}

// AsyncInsert 异步插入（ClickHouse 22.3+）
func (r *EventRepository) AsyncInsert(ctx Context.Context, events []*Event) error {
    // 异步插入会在服务器端批量处理
    query := `
    INSERT INTO events (
        event_id, user_id, event_type, event_time, properties,
        country, city, device, os, browser, session_id
    ) VALUES
    `

    // 使用 SETTINGS async_insert=1
    ctx = clickhouse.Context(ctx, clickhouse.WithSettings(clickhouse.Settings{
        "async_insert": 1,
        "wait_for_async_insert": 0, // 不等待插入完成
    }))

    batch, err := r.conn.PrepareBatch(ctx, query)
    if err != nil {
        return err
    }

    for _, event := range events {
        batch.Append(
            event.EventID, event.UserID, event.EventType, event.EventTime,
            event.Properties, event.Country, event.City, event.Device,
            event.OS, event.Browser, event.SessionID,
        )
    }

    return batch.Send()
}
```

---

## 数据查询

### 基本查询

```go
// FindByUserID 根据用户ID查询事件
func (r *EventRepository) FindByUserID(ctx Context.Context, userID uint64, limit int) ([]*Event, error) {
    query := `
    SELECT event_id, user_id, event_type, event_time, properties,
           country, city, device, os, browser, session_id, created_at
    FROM events
    WHERE user_id = ?
    ORDER BY event_time DESC
    LIMIT ?
    `

    rows, err := r.conn.Query(ctx, query, userID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []*Event
    for rows.Next() {
        var event Event
        if err := rows.Scan(
            &event.EventID, &event.UserID, &event.EventType, &event.EventTime,
            &event.Properties, &event.Country, &event.City, &event.Device,
            &event.OS, &event.Browser, &event.SessionID, &event.CreatedAt,
        ); err != nil {
            return nil, err
        }
        events = append(events, &event)
    }

    return events, rows.Err()
}

// CountByEventType 按事件类型统计
func (r *EventRepository) CountByEventType(ctx Context.Context, startTime, endTime time.Time) (map[string]uint64, error) {
    query := `
    SELECT event_type, count() as count
    FROM events
    WHERE event_time >= ? AND event_time < ?
    GROUP BY event_type
    ORDER BY count DESC
    `

    rows, err := r.conn.Query(ctx, query, startTime, endTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]uint64)
    for rows.Next() {
        var eventType string
        var count uint64
        if err := rows.Scan(&eventType, &count); err != nil {
            return nil, err
        }
        result[eventType] = count
    }

    return result, rows.Err()
}
```

### 聚合查询

```go
// GetDailyStatistics 获取每日统计
func (r *EventRepository) GetDailyStatistics(ctx Context.Context, days int) ([]DailyStats, error) {
    query := `
    SELECT
        toDate(event_time) as date,
        count() as total_events,
        uniq(user_id) as unique_users,
        uniq(session_id) as sessions,
        countIf(event_type = 'page_view') as page_views,
        countIf(event_type = 'click') as clicks,
        avg(toUInt32(JSONExtractString(properties, 'duration'))) as avg_duration
    FROM events
    WHERE event_time >= now() - INTERVAL ? DAY
    GROUP BY date
    ORDER BY date DESC
    `

    rows, err := r.conn.Query(ctx, query, days)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []DailyStats
    for rows.Next() {
        var stat DailyStats
        if err := rows.Scan(
            &stat.Date, &stat.TotalEvents, &stat.UniqueUsers,
            &stat.Sessions, &stat.PageViews, &stat.Clicks, &stat.AvgDuration,
        ); err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }

    return stats, rows.Err()
}

// DailyStats 每日统计
type DailyStats struct {
    Date        time.Time
    TotalEvents uint64
    UniqueUsers uint64
    Sessions    uint64
    PageViews   uint64
    Clicks      uint64
    AvgDuration float64
}

// GetTopCountries 获取Top国家
func (r *EventRepository) GetTopCountries(ctx Context.Context, limit int) ([]CountryStats, error) {
    query := `
    SELECT
        country,
        count() as events,
        uniq(user_id) as users,
        uniq(session_id) as sessions
    FROM events
    WHERE event_time >= now() - INTERVAL 7 DAY
    GROUP BY country
    ORDER BY events DESC
    LIMIT ?
    `

    rows, err := r.conn.Query(ctx, query, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []CountryStats
    for rows.Next() {
        var stat CountryStats
        if err := rows.Scan(&stat.Country, &stat.Events, &stat.Users, &stat.Sessions); err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }

    return stats, rows.Err()
}

// CountryStats 国家统计
type CountryStats struct {
    Country  string
    Events   uint64
    Users    uint64
    Sessions uint64
}
```

### 漏斗分析

```go
// FunnelAnalysis 漏斗分析
func (r *EventRepository) FunnelAnalysis(ctx Context.Context, steps []string, window int) (*FunnelResult, error) {
    // 使用windowFunnel函数
    query := `
    SELECT
        windowFunnel(?, 'strict')(event_time, ` + buildFunnelConditions(steps) + `) as level,
        count() as users
    FROM events
    WHERE event_time >= now() - INTERVAL 7 DAY
    GROUP BY level
    ORDER BY level
    `

    rows, err := r.conn.Query(ctx, query, window)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := &FunnelResult{Steps: steps}
    for rows.Next() {
        var level int
        var users uint64
        if err := rows.Scan(&level, &users); err != nil {
            return nil, err
        }
        result.Levels = append(result.Levels, FunnelLevel{Level: level, Users: users})
    }

    return result, rows.Err()
}

// buildFunnelConditions 构建漏斗条件
func buildFunnelConditions(steps []string) string {
    var conditions []string
    for _, step := range steps {
        conditions = append(conditions, fmt.Sprintf("event_type = '%s'", step))
    }
    return strings.Join(conditions, ", ")
}

// FunnelResult 漏斗结果
type FunnelResult struct {
    Steps  []string
    Levels []FunnelLevel
}

// FunnelLevel 漏斗层级
type FunnelLevel struct {
    Level int
    Users uint64
}
```

### 留存分析

```go
// RetentionAnalysis 留存分析
func (r *EventRepository) RetentionAnalysis(ctx Context.Context, days int) ([][]float64, error) {
    query := `
    SELECT
        retention
    FROM
    (
        SELECT
            user_id,
            groupArray(toDate(event_time)) as dates
        FROM events
        WHERE event_time >= today() - INTERVAL ? DAY
        GROUP BY user_id
    )
    ARRAY JOIN retention(dates, today() - INTERVAL ? DAY, toDate(today()), 1) as retention
    `

    rows, err := r.conn.Query(ctx, query, days, days)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var retentionData [][]float64
    // 处理结果
    // ...

    return retentionData, nil
}
```

---

## 物化视图

### 创建物化视图

```go
// CreateMaterializedView 创建物化视图
func (r *EventRepository) CreateMaterializedView(ctx Context.Context) error {
    // 创建目标表
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS event_hourly_stats (
        event_date Date,
        event_hour UInt8,
        event_type String,
        country String,
        total_events UInt64,
        unique_users UInt64
    )
    ENGINE = SummingMergeTree()
    PARTITION BY toYYYYMM(event_date)
    ORDER BY (event_date, event_hour, event_type, country)
    `

    if err := r.conn.Exec(ctx, createTableQuery); err != nil {
        return err
    }

    // 创建物化视图
    createMVQuery := `
    CREATE MATERIALIZED VIEW IF NOT EXISTS event_hourly_stats_mv
    TO event_hourly_stats
    AS SELECT
        toDate(event_time) as event_date,
        toHour(event_time) as event_hour,
        event_type,
        country,
        count() as total_events,
        uniqState(user_id) as unique_users
    FROM events
    GROUP BY event_date, event_hour, event_type, country
    `

    if err := r.conn.Exec(ctx, createMVQuery); err != nil {
        return err
    }

    fmt.Println("✅ 物化视图创建成功")

    return nil
}

// QueryMaterializedView 查询物化视图
func (r *EventRepository) QueryMaterializedView(ctx Context.Context, date time.Time) ([]HourlyStats, error) {
    query := `
    SELECT
        event_hour,
        event_type,
        country,
        sum(total_events) as events,
        uniqMerge(unique_users) as users
    FROM event_hourly_stats
    WHERE event_date = ?
    GROUP BY event_hour, event_type, country
    ORDER BY event_hour, events DESC
    `

    rows, err := r.conn.Query(ctx, query, date)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []HourlyStats
    for rows.Next() {
        var stat HourlyStats
        if err := rows.Scan(&stat.Hour, &stat.EventType, &stat.Country, &stat.Events, &stat.Users); err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }

    return stats, rows.Err()
}

// HourlyStats 小时统计
type HourlyStats struct {
    Hour      uint8
    EventType string
    Country   string
    Events    uint64
    Users     uint64
}
```

---

## 分区与分片

### 分区管理

```go
// ListPartitions 列出所有分区
func (r *EventRepository) ListPartitions(ctx Context.Context) ([]PartitionInfo, error) {
    query := `
    SELECT
        partition,
        rows,
        bytes_on_disk,
        modification_time
    FROM system.parts
    WHERE table = 'events' AND active = 1
    ORDER BY partition
    `

    rows, err := r.conn.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var partitions []PartitionInfo
    for rows.Next() {
        var p PartitionInfo
        if err := rows.Scan(&p.Partition, &p.Rows, &p.BytesOnDisk, &p.ModificationTime); err != nil {
            return nil, err
        }
        partitions = append(partitions, p)
    }

    return partitions, rows.Err()
}

// PartitionInfo 分区信息
type PartitionInfo struct {
    Partition        string
    Rows             uint64
    BytesOnDisk      uint64
    ModificationTime time.Time
}

// DropPartition 删除分区
func (r *EventRepository) DropPartition(ctx Context.Context, partition string) error {
    query := fmt.Sprintf("ALTER TABLE events DROP PARTITION '%s'", partition)
    return r.conn.Exec(ctx, query)
}

// DetachPartition 分离分区（不删除数据）
func (r *EventRepository) DetachPartition(ctx Context.Context, partition string) error {
    query := fmt.Sprintf("ALTER TABLE events DETACH PARTITION '%s'", partition)
    return r.conn.Exec(ctx, query)
}

// AttachPartition 附加分区
func (r *EventRepository) AttachPartition(ctx Context.Context, partition string) error {
    query := fmt.Sprintf("ALTER TABLE events ATTACH PARTITION '%s'", partition)
    return r.conn.Exec(ctx, query)
}
```

---

## 实时数据分析

### 实时看板

```go
package analytics

import (
    "Context"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Dashboard 实时看板
type Dashboard struct {
    conn driver.Conn
}

// NewDashboard 创建看板
func NewDashboard(conn driver.Conn) *Dashboard {
    return &Dashboard{conn: conn}
}

// GetRealtimeMetrics 获取实时指标
func (d *Dashboard) GetRealtimeMetrics(ctx Context.Context) (*RealtimeMetrics, error) {
    query := `
    SELECT
        -- 最近1分钟
        countIf(event_time >= now() - INTERVAL 1 MINUTE) as events_1m,
        uniqIf(user_id, event_time >= now() - INTERVAL 1 MINUTE) as users_1m,

        -- 最近5分钟
        countIf(event_time >= now() - INTERVAL 5 MINUTE) as events_5m,
        uniqIf(user_id, event_time >= now() - INTERVAL 5 MINUTE) as users_5m,

        -- 最近1小时
        countIf(event_time >= now() - INTERVAL 1 HOUR) as events_1h,
        uniqIf(user_id, event_time >= now() - INTERVAL 1 HOUR) as users_1h,

        -- 今日
        countIf(toDate(event_time) = today()) as events_today,
        uniqIf(user_id, toDate(event_time) = today()) as users_today
    FROM events
    WHERE event_time >= today()
    `

    row := d.conn.QueryRow(ctx, query)

    var metrics RealtimeMetrics
    if err := row.Scan(
        &metrics.Events1m, &metrics.Users1m,
        &metrics.Events5m, &metrics.Users5m,
        &metrics.Events1h, &metrics.Users1h,
        &metrics.EventsToday, &metrics.UsersToday,
    ); err != nil {
        return nil, err
    }

    return &metrics, nil
}

// RealtimeMetrics 实时指标
type RealtimeMetrics struct {
    Events1m     uint64
    Users1m      uint64
    Events5m     uint64
    Users5m      uint64
    Events1h     uint64
    Users1h      uint64
    EventsToday  uint64
    UsersToday   uint64
}

// GetActiveUsers 获取活跃用户
func (d *Dashboard) GetActiveUsers(ctx Context.Context, minutes int) ([]ActiveUser, error) {
    query := `
    SELECT
        user_id,
        count() as event_count,
        max(event_time) as last_event_time
    FROM events
    WHERE event_time >= now() - INTERVAL ? MINUTE
    GROUP BY user_id
    ORDER BY event_count DESC
    LIMIT 100
    `

    rows, err := d.conn.Query(ctx, query, minutes)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []ActiveUser
    for rows.Next() {
        var user ActiveUser
        if err := rows.Scan(&user.UserID, &user.EventCount, &user.LastEventTime); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, rows.Err()
}

// ActiveUser 活跃用户
type ActiveUser struct {
    UserID        uint64
    EventCount    uint64
    LastEventTime time.Time
}
```

---

## 性能优化

### 查询优化

```go
// 1. 使用PREWHERE替代WHERE（更快的数据过滤）
query := `
SELECT count()
FROM events
PREWHERE event_type = 'click'  -- 先在列存储层过滤
WHERE country = 'US'
`

// 2. 避免SELECT *，只查询需要的列
query := `
SELECT event_id, event_type, event_time
FROM events
WHERE user_id = ?
`

// 3. 使用FINAL慎重（会降低性能）
query := `
SELECT *
FROM events
-- FINAL  -- 避免使用，除非必需
WHERE user_id = ?
`

// 4. 合理使用LIMIT
query := `
SELECT *
FROM events
ORDER BY event_time DESC
LIMIT 1000  -- 限制返回行数
`

// 5. 使用采样查询（大数据集）
query := `
SELECT event_type, count()
FROM events
SAMPLE 0.1  -- 采样10%的数据
WHERE event_time >= today()
GROUP BY event_type
`
```

### 索引优化

```go
// 创建跳数索引（Skip Index）
func (r *EventRepository) CreateSkipIndex(ctx Context.Context) error {
    queries := []string{
        // MinMax索引（适合范围查询）
        `ALTER TABLE events ADD INDEX idx_user_id_minmax user_id TYPE minmax GRANULARITY 4`,

        // Set索引（适合IN查询）
        `ALTER TABLE events ADD INDEX idx_country_set country TYPE set(100) GRANULARITY 4`,

        // Bloom Filter索引（适合等值查询）
        `ALTER TABLE events ADD INDEX idx_event_id_bloom event_id TYPE bloom_filter GRANULARITY 4`,
    }

    for _, query := range queries {
        if err := r.conn.Exec(ctx, query); err != nil {
            return err
        }
    }

    fmt.Println("✅ 跳数索引创建成功")

    return nil
}
```

### 批量写入优化

```go
// BatchWriter 批量写入器
type BatchWriter struct {
    repo      *EventRepository
    batchSize int
    buffer    []*Event
    mu        sync.Mutex
    ticker    *time.Ticker
    done      Channel struct{}
}

// NewBatchWriter 创建批量写入器
func NewBatchWriter(repo *EventRepository, batchSize int, flushInterval time.Duration) *BatchWriter {
    bw := &BatchWriter{
        repo:      repo,
        batchSize: batchSize,
        buffer:    make([]*Event, 0, batchSize),
        ticker:    time.NewTicker(flushInterval),
        done:      make(Channel struct{}),
    }

    // 启动定时刷新
    go bw.autoFlush()

    return bw
}

// Write 写入事件
func (bw *BatchWriter) Write(event *Event) error {
    bw.mu.Lock()
    defer bw.mu.Unlock()

    bw.buffer = append(bw.buffer, event)

    // 达到批次大小，立即刷新
    if len(bw.buffer) >= bw.batchSize {
        return bw.flush()
    }

    return nil
}

// flush 刷新缓冲区
func (bw *BatchWriter) flush() error {
    if len(bw.buffer) == 0 {
        return nil
    }

    ctx, cancel := Context.WithTimeout(Context.Background(), 10*time.Second)
    defer cancel()

    err := bw.repo.BatchInsert(ctx, bw.buffer)
    if err != nil {
        return err
    }

    bw.buffer = bw.buffer[:0] // 清空缓冲区

    return nil
}

// autoFlush 自动刷新
func (bw *BatchWriter) autoFlush() {
    for {
        select {
        case <-bw.ticker.C:
            bw.mu.Lock()
            bw.flush()
            bw.mu.Unlock()

        case <-bw.done:
            return
        }
    }
}

// Close 关闭写入器
func (bw *BatchWriter) Close() error {
    close(bw.done)
    bw.ticker.Stop()

    bw.mu.Lock()
    defer bw.mu.Unlock()

    return bw.flush()
}
```

---

## 监控与运维

### 系统监控

```go
// GetSystemMetrics 获取系统指标
func (d *Dashboard) GetSystemMetrics(ctx Context.Context) (*SystemMetrics, error) {
    query := `
    SELECT
        uptime() as uptime,
        -- 查询统计
        (SELECT count() FROM system.query_log WHERE event_time >= now() - INTERVAL 1 MINUTE) as queries_1m,
        -- 慢查询
        (SELECT count() FROM system.query_log WHERE query_duration_ms > 1000 AND event_time >= now() - INTERVAL 1 HOUR) as slow_queries_1h,
        -- 当前连接数
        (SELECT count() FROM system.processes) as current_connections,
        -- 表大小
        (SELECT sum(bytes) FROM system.parts WHERE table = 'events' AND active = 1) as table_size_bytes,
        -- 行数
        (SELECT sum(rows) FROM system.parts WHERE table = 'events' AND active = 1) as table_rows
    `

    row := d.conn.QueryRow(ctx, query)

    var metrics SystemMetrics
    if err := row.Scan(
        &metrics.Uptime,
        &metrics.Queries1m,
        &metrics.SlowQueries1h,
        &metrics.CurrentConnections,
        &metrics.TableSizeBytes,
        &metrics.TableRows,
    ); err != nil {
        return nil, err
    }

    return &metrics, nil
}

// SystemMetrics 系统指标
type SystemMetrics struct {
    Uptime             uint64
    Queries1m          uint64
    SlowQueries1h      uint64
    CurrentConnections uint64
    TableSizeBytes     uint64
    TableRows          uint64
}
```

---

## 最佳实践

### 1. 表设计原则

- ✅ 使用合适的表引擎（MergeTree家族）
- ✅ 选择合适的ORDER BY（影响查询性能）
- ✅ 合理设置分区（通常按时间分区）
- ✅ 使用TTL自动清理旧数据
- ✅ 避免过多的列（影响压缩率）

### 2. 查询优化

- ✅ 在ORDER BY的第一列上过滤效果最好
- ✅ 使用PREWHERE进行列裁剪
- ✅ 避免SELECT *
- ✅ 合理使用LIMIT
- ✅ 对大数据集使用SAMPLE

### 3. 写入优化

- ✅ 批量插入（至少1000行）
- ✅ 使用异步插入
- ✅ 避免频繁的小批次写入
- ✅ 合理设置merge参数

### 4. 监控告警

```go
// 监控指标
- 查询延迟 (P50, P95, P99)
- 慢查询数量
- 插入速率 (rows/s)
- 磁盘使用率
- 副本延迟
- Merge速度
```

---

## 总结
