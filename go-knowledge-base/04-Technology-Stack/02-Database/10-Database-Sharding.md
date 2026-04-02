# 数据库分片 (Database Sharding)

> **分类**: 开源技术堆栈  
> **标签**: #sharding #database #scaling

---

## 分片策略

### 哈希分片

```go
type ShardingStrategy interface {
    GetShard(key string) int
}

type HashSharding struct {
    shardCount int
}

func (h *HashSharding) GetShard(key string) int {
    hash := fnv32(key)
    return int(hash) % h.shardCount
}

func fnv32(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32()
}
```

### 范围分片

```go
type RangeSharding struct {
    ranges []Range
}

type Range struct {
    Min   int64
    Max   int64
    Shard int
}

func (r *RangeSharding) GetShard(key int64) int {
    for _, rng := range r.ranges {
        if key >= rng.Min && key < rng.Max {
            return rng.Shard
        }
    }
    return 0
}
```

---

## 分片管理器

```go
type ShardingManager struct {
    shards  []*sql.DB
    strategy ShardingStrategy
}

func (sm *ShardingManager) GetDB(key string) *sql.DB {
    shardIndex := sm.strategy.GetShard(key)
    return sm.shards[shardIndex]
}

func (sm *ShardingManager) Query(ctx context.Context, key string, query string, args ...interface{}) (*sql.Rows, error) {
    db := sm.GetDB(key)
    return db.QueryContext(ctx, query, args...)
}

func (sm *ShardingManager) QueryAll(ctx context.Context, query string, args ...interface{}) ([]*sql.Rows, error) {
    var results []*sql.Rows
    var errs []error
    
    for _, db := range sm.shards {
        rows, err := db.QueryContext(ctx, query, args...)
        if err != nil {
            errs = append(errs, err)
            continue
        }
        results = append(results, rows)
    }
    
    if len(errs) > 0 {
        return nil, fmt.Errorf("query errors: %v", errs)
    }
    return results, nil
}
```

---

## 全局 ID 生成

```go
// 雪花算法
type Snowflake struct {
    mu        sync.Mutex
    lastStamp int64
    sequence  int64
    workerID  int64
}

func (s *Snowflake) NextID() int64 {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now().UnixMilli()
    if now == s.lastStamp {
        s.sequence = (s.sequence + 1) & 4095
        if s.sequence == 0 {
            for now <= s.lastStamp {
                now = time.Now().UnixMilli()
            }
        }
    } else {
        s.sequence = 0
    }
    
    s.lastStamp = now
    
    return ((now - epoch) << 22) | (s.workerID << 12) | s.sequence
}
```

---

## 跨分片事务

```go
func (sm *ShardingManager) CrossShardTx(ctx context.Context, keys []string, fn func(map[string]*sql.Tx) error) error {
    // 收集涉及的分片
    shardMap := make(map[int]*sql.DB)
    for _, key := range keys {
        shard := sm.strategy.GetShard(key)
        shardMap[shard] = sm.shards[shard]
    }
    
    // 开启所有事务
    txs := make(map[string]*sql.Tx)
    for shard, db := range shardMap {
        tx, err := db.BeginTx(ctx, nil)
        if err != nil {
            // 回滚已开启的事务
            for _, t := range txs {
                t.Rollback()
            }
            return err
        }
        txs[fmt.Sprintf("shard_%d", shard)] = tx
    }
    
    // 执行业务逻辑
    if err := fn(txs); err != nil {
        for _, t := range txs {
            t.Rollback()
        }
        return err
    }
    
    // 提交所有事务
    for _, t := range txs {
        if err := t.Commit(); err != nil {
            // 需要补偿或人工介入
            return err
        }
    }
    
    return nil
}
```
