# TS-044-Database-Sharding-Partitioning

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: 2026 (Horizontal Sharding, Vertical Partitioning, Consistent Hashing)
> **Size**: >20KB

---

## 1. 分片与分区概述

### 1.1 扩展策略对比

| 策略 | 数据分布 | 适用场景 | 复杂度 |
|------|---------|---------|--------|
| 垂直分区 | 按列分离 | 宽表、冷热数据 | 低 |
| 水平分区 | 单表内分区 | 时间序列、地理数据 | 中 |
| 水平分片 | 多实例分布 | 大规模数据、高并发 | 高 |

### 1.2 分片架构

```
┌─────────────────────────────────────────┐
│         Sharded Database Architecture   │
├─────────────────────────────────────────┤
│                                         │
│  Application                            │
│       │                                 │
│       ▼                                 │
│  ┌─────────────┐                        │
│  │ Sharding    │                        │
│  │ Router      │                        │
│  │ (Proxy/     │  - SQL解析             │
│  │  Middleware)│  - 路由计算             │
│  └──────┬──────┘  - 结果聚合             │
│         │                               │
│    ┌────┼────┬─────────┬────────┐      │
│    │    │    │         │        │      │
│    ▼    ▼    ▼         ▼        ▼      │
│  ┌──┐ ┌──┐ ┌──┐     ┌──┐    ┌──┐      │
│  │DB1│ │DB2│ │DB3│ ... │DBn│    │DBn+1│      │
│  └──┘ └──┘ └──┘     └──┘    └──┘      │
│  Shard1 Shard2 Shard3      ShardN      │
│                                         │
│  路由键: user_id % N                     │
│  或: hash(user_id) % N                   │
└─────────────────────────────────────────┘
```

---

## 2. 分片策略

### 2.1 哈希分片

```go
// 一致性哈希分片
type ConsistentHasher struct {
    ring   *consistent.Consistent
    shards map[string]*sql.DB
}

func NewConsistentHasher(shards []Shard) *ConsistentHasher {
    ring := consistent.New()
    shardMap := make(map[string]*sql.DB)

    for _, shard := range shards {
        ring.Add(shard.ID)
        shardMap[shard.ID] = shard.DB
    }

    return &ConsistentHasher{
        ring:   ring,
        shards: shardMap,
    }
}

func (c *ConsistentHasher) GetShard(key string) (*sql.DB, error) {
    shardID, err := c.ring.Get(key)
    if err != nil {
        return nil, err
    }

    return c.shards[shardID], nil
}

// 使用
hasher := NewConsistentHasher(shards)
db, _ := hasher.GetShard("user:12345")
```

### 2.2 范围分片

```go
// 按ID范围分片
type RangeShard struct {
    ID      string
    DB      *sql.DB
    MinID   int64
    MaxID   int64
}

type RangeRouter struct {
    shards []RangeShard
}

func (r *RangeRouter) GetShard(id int64) (*sql.DB, error) {
    for _, shard := range r.shards {
        if id >= shard.MinID && id <= shard.MaxID {
            return shard.DB, nil
        }
    }
    return nil, ErrShardNotFound
}

// 配置示例
shards := []RangeShard{
    {ID: "shard1", MinID: 1, MaxID: 1000000},
    {ID: "shard2", MinID: 1000001, MaxID: 2000000},
    {ID: "shard3", MinID: 2000001, MaxID: 3000000},
}
```

### 2.3 列表分片

```go
// 按地区分片
type ListRouter struct {
    regionMap map[string]*sql.DB
}

func (r *ListRouter) GetShardByRegion(region string) (*sql.DB, error) {
    db, ok := r.regionMap[region]
    if !ok {
        return nil, ErrRegionNotFound
    }
    return db, nil
}

// 配置
regionMap := map[string]*sql.DB{
    "US-East":  db1,
    "US-West":  db2,
    "EU-West":  db3,
    "AP-South": db4,
}
```

---

## 3. Sharding Proxy实现

### 3.1 完整实现

```go
// ShardingProxy 实现SQL路由
type ShardingProxy struct {
    router      ShardRouter
    parser      *sqlparser.Parser
    aggregators map[string]ResultAggregator
}

type ShardRouter interface {
    Route(sql string, params []interface{}) ([]ShardTarget, error)
}

type ShardTarget struct {
    ShardID string
    DB      *sql.DB
}

func (p *ShardingProxy) Execute(ctx context.Context, query string, args ...interface{}) (*Result, error) {
    // 1. 解析SQL
    stmt, err := p.parser.Parse(query)
    if err != nil {
        return nil, err
    }

    // 2. 确定路由
    targets, err := p.router.Route(query, args)
    if err != nil {
        return nil, err
    }

    // 3. 根据SQL类型执行
    switch stmt.(type) {
    case *sqlparser.Select:
        return p.executeSelect(ctx, stmt.(*sqlparser.Select), targets)
    case *sqlparser.Insert:
        return p.executeInsert(ctx, stmt.(*sqlparser.Insert), targets)
    case *sqlparser.Update:
        return p.executeUpdate(ctx, stmt.(*sqlparser.Update), targets)
    case *sqlparser.Delete:
        return p.executeDelete(ctx, stmt.(*sqlparser.Delete), targets)
    default:
        return nil, ErrUnsupportedQuery
    }
}

func (p *ShardingProxy) executeSelect(ctx context.Context, stmt *sqlparser.Select, targets []ShardTarget) (*Result, error) {
    // 单分片查询
    if len(targets) == 1 {
        return p.querySingle(ctx, targets[0], stmt)
    }

    // 跨分片查询 - 并行执行
    results := make([]*Result, len(targets))
    var wg sync.WaitGroup
    errChan := make(chan error, len(targets))

    for i, target := range targets {
        wg.Add(1)
        go func(idx int, t ShardTarget) {
            defer wg.Done()

            result, err := p.querySingle(ctx, t, stmt)
            if err != nil {
                errChan <- err
                return
            }
            results[idx] = result
        }(i, target)
    }

    wg.Wait()
    close(errChan)

    if len(errChan) > 0 {
        return nil, <-errChan
    }

    // 4. 聚合结果
    return p.aggregateResults(results, stmt)
}

func (p *ShardingProxy) aggregateResults(results []*Result, stmt *sqlparser.Select) (*Result, error) {
    // 合并行
    allRows := []Row{}
    for _, r := range results {
        allRows = append(allRows, r.Rows...)
    }

    // 处理ORDER BY
    if len(stmt.OrderBy) > 0 {
        allRows = p.sortRows(allRows, stmt.OrderBy)
    }

    // 处理LIMIT/OFFSET
    if stmt.Limit != nil {
        offset := 0
        limit := len(allRows)

        if stmt.Limit.Offset != nil {
            offset = evalExpr(stmt.Limit.Offset)
        }
        if stmt.Limit.Rowcount != nil {
            limit = evalExpr(stmt.Limit.Rowcount)
        }

        end := offset + limit
        if end > len(allRows) {
            end = len(allRows)
        }
        allRows = allRows[offset:end]
    }

    // 处理聚合函数
    if p.hasAggregation(stmt) {
        return p.computeAggregation(allRows, stmt)
    }

    return &Result{Rows: allRows}, nil
}
```

---

## 4. 分片键选择

### 4.1 选择原则

```
好的分片键特征:
1. 高基数 (Cardinality) - 避免热点
2. 均匀分布 - 数据均衡
3. 访问局部性 - 相关数据在一起
4. 不变性 - 避免重新分片
```

### 4.2 常见选择

| 实体 | 推荐分片键 | 理由 |
|------|-----------|------|
| 用户 | user_id | 高基数，均匀 |
| 订单 | user_id | 按用户查询多 |
| 商品 | category_id + id | 分类查询 |
| 日志 | timestamp | 时间范围查询 |

### 4.3 热点处理

```go
// 热点检测
type HotspotDetector struct {
    counters map[string]* sliding.Window
    threshold int
}

func (h *HotspotDetector) Check(key string) bool {
    counter := h.counters[key]
    count := counter.Increment()

    if count > h.threshold {
        // 热点key，使用二级分片
        return true
    }
    return false
}

// 热点key二次分片
func (p *ShardingProxy) routeWithHotspot(key string) ShardTarget {
    if p.hotspotDetector.Check(key) {
        // 使用hash(key + random)分散到多个shard
        suffix := rand.Intn(10)
        hotKey := fmt.Sprintf("%s#%d", key, suffix)
        return p.router.GetShard(hotKey)
    }
    return p.router.GetShard(key)
}
```

---

## 5. 数据迁移与扩容

### 5.1 在线扩容

```go
// 双写迁移策略
type MigrationService struct {
    oldShards []Shard
    newShards []Shard
    phase     MigrationPhase
}

type MigrationPhase int

const (
    PhaseDualWrite MigrationPhase = iota
    PhaseBackfill
    PhaseVerify
    PhaseCutover
    PhaseCleanup
)

func (m *MigrationService) DualWrite(ctx context.Context, key string, data interface{}) error {
    // 写入旧分片
    if err := m.writeToOld(ctx, key, data); err != nil {
        return err
    }

    // 异步写入新分片
    go func() {
        newShard := m.newRouter.GetShard(key)
        newShard.Write(key, data)
    }()

    return nil
}

func (m *MigrationService) Backfill(ctx context.Context) error {
    // 遍历旧分片数据
    for _, oldShard := range m.oldShards {
        cursor := oldShard.Scan()

        for cursor.Next() {
            key, data := cursor.Get()

            // 计算新分片位置
            newShard := m.newRouter.GetShard(key)

            // 写入新分片
            if err := newShard.Write(key, data); err != nil {
                m.logFailedMigration(key, err)
            }
        }
    }

    return nil
}

func (m *MigrationService) Verify(ctx context.Context) error {
    // 抽样验证数据一致性
    samples := m.sampleKeys(0.01)  // 1%抽样

    for _, key := range samples {
        oldData := m.readFromOld(ctx, key)
        newData := m.readFromNew(ctx, key)

        if !reflect.DeepEqual(oldData, newData) {
            return fmt.Errorf("data mismatch for key %s", key)
        }
    }

    return nil
}
```

---

## 6. 分布式ID生成

### 6.1 Snowflake算法

```go
// 64位ID结构:
// 0|0000000000 0000000000 0000000000 0000000000 0|00000|00000|000000000000
// 符号位|  41位时间戳(毫秒)                        |数据中心|机器ID|  12位序列号

type Snowflake struct {
    mu        sync.Mutex
    epoch     int64
    nodeID    int64
    sequence  int64
    lastTime  int64
}

func NewSnowflake(nodeID int64) *Snowflake {
    return &Snowflake{
        epoch:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
        nodeID:   nodeID,
        sequence: 0,
    }
}

func (s *Snowflake) Generate() int64 {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UnixMilli()

    if now < s.lastTime {
        panic("clock moved backwards")
    }

    if now == s.lastTime {
        s.sequence = (s.sequence + 1) & 0xFFF
        if s.sequence == 0 {
            // 等待下一毫秒
            for now <= s.lastTime {
                now = time.Now().UnixMilli()
            }
        }
    } else {
        s.sequence = 0
    }

    s.lastTime = now

    // 组合ID
    id := ((now - s.epoch) << 22) |
          (s.nodeID << 12) |
          s.sequence

    return id
}

func (s *Snowflake) Parse(id int64) (timestamp time.Time, nodeID, sequence int64) {
    timestamp = time.UnixMilli((id >> 22) + s.epoch)
    nodeID = (id >> 12) & 0x3FF
    sequence = id & 0xFFF
    return
}
```

---

## 7. 最佳实践

### 7.1 查询优化

```sql
-- 跨分片查询避免
-- BAD: 需要查询所有分片
SELECT * FROM orders WHERE status = 'pending';

-- GOOD: 带分片键
SELECT * FROM orders WHERE user_id = 123 AND status = 'pending';

-- 聚合优化
-- BAD: 跨分片聚合
SELECT COUNT(*) FROM orders;

-- GOOD: 预聚合或限制范围
SELECT COUNT(*) FROM orders WHERE user_id = 123;
```

### 7.2 事务处理

```go
// 分片事务 - Saga模式
func (s *Service) CrossShardTransaction(ctx context.Context, userID int64, order Order) error {
    userShard := s.router.GetShard(userID)
    orderShard := s.router.GetShard(order.ID)

    // 两阶段提交简化版

    // Phase 1: 准备
    userTx := userShard.Begin()
    if err := userTx.UpdateBalance(userID, -order.Amount); err != nil {
        userTx.Rollback()
        return err
    }

    orderTx := orderShard.Begin()
    if err := orderTx.CreateOrder(order); err != nil {
        orderTx.Rollback()
        userTx.Rollback()
        return err
    }

    // Phase 2: 提交
    if err := userTx.Commit(); err != nil {
        orderTx.Rollback()
        return err
    }

    if err := orderTx.Commit(); err != nil {
        // 需要补偿机制
        s.compensateUserBalance(userID, order.Amount)
        return err
    }

    return nil
}
```

---

## 8. 参考文献

1. "Designing Data-Intensive Applications" - Martin Kleppmann
2. Vitess Documentation (YouTube Sharding)
3. ShardingSphere Documentation
4. CitusDB Architecture
5. CockroachDB Distribution Layer

---

*Last Updated: 2026-04-03*
