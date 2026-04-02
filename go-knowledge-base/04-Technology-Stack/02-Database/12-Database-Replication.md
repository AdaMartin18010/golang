# 数据库复制 (Database Replication)

> **分类**: 开源技术堆栈  
> **标签**: #replication #database #scaling

---

## 读写分离

```go
type DBCluster struct {
    master *sql.DB
    slaves []*sql.DB
    counter uint64
}

func (c *DBCluster) Master() *sql.DB {
    return c.master
}

func (c *DBCluster) Slave() *sql.DB {
    // 轮询选择从库
    if len(c.slaves) == 0 {
        return c.master
    }
    
    idx := atomic.AddUint64(&c.counter, 1) % uint64(len(c.slaves))
    return c.slaves[idx]
}

// 写操作
func (c *DBCluster) Write(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    return c.Master().ExecContext(ctx, query, args...)
}

// 读操作
func (c *DBCluster) Read(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    return c.Slave().QueryContext(ctx, query, args...)
}
```

---

## 延迟检测

```go
type ReplicationLagMonitor struct {
    slaves []*sql.DB
}

func (m *ReplicationLagMonitor) CheckLag(ctx context.Context) (map[string]time.Duration, error) {
    lagMap := make(map[string]time.Duration)
    
    for i, slave := range m.slaves {
        var secondsBehind float64
        err := slave.QueryRowContext(ctx, 
            "SHOW SLAVE STATUS").Scan(&secondsBehind)
        
        if err != nil {
            lagMap[fmt.Sprintf("slave-%d", i)] = -1
        } else {
            lagMap[fmt.Sprintf("slave-%d", i)] = 
                time.Duration(secondsBehind) * time.Second
        }
    }
    
    return lagMap, nil
}

func (m *ReplicationLagMonitor) GetHealthySlave(ctx context.Context, maxLag time.Duration) (*sql.DB, error) {
    for _, slave := range m.slaves {
        var lag float64
        err := slave.QueryRowContext(ctx, 
            "SHOW SLAVE STATUS").Scan(&lag)
        
        if err == nil && time.Duration(lag)*time.Second <= maxLag {
            return slave, nil
        }
    }
    
    return nil, errors.New("no healthy slave available")
}
```

---

## 复制模式

### 异步复制

```go
// 主库写入后立即返回
func AsyncWrite(ctx context.Context, db *sql.DB, data interface{}) error {
    _, err := db.ExecContext(ctx, 
        "INSERT INTO data (value) VALUES (?)", data)
    return err
}
```

### 半同步复制

```go
// 等待至少一个从库确认
func SemiSyncWrite(ctx context.Context, master *sql.DB, data interface{}) error {
    tx, err := master.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 启用半同步
    _, err = tx.Exec("SET rpl_semi_sync_master_wait_for_slave_count = 1")
    if err != nil {
        return err
    }
    
    _, err = tx.Exec("INSERT INTO data (value) VALUES (?)", data)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
```

---

## 故障转移

```go
type FailoverManager struct {
    master  *sql.DB
    slaves  []*sql.DB
    current *sql.DB
}

func (fm *FailoverManager) HealthCheck(ctx context.Context) error {
    // 检查主库
    if err := fm.master.PingContext(ctx); err != nil {
        // 主库故障，选举新主库
        newMaster, err := fm.electNewMaster(ctx)
        if err != nil {
            return err
        }
        fm.promoteSlave(newMaster)
    }
    
    return nil
}

func (fm *FailoverManager) electNewMaster(ctx context.Context) (*sql.DB, error) {
    // 选择延迟最小的从库
    var bestSlave *sql.DB
    var minLag time.Duration = time.Hour
    
    for _, slave := range fm.slaves {
        var lag float64
        err := slave.QueryRowContext(ctx, 
            "SHOW SLAVE STATUS").Scan(&lag)
        
        if err == nil {
            slaveLag := time.Duration(lag) * time.Second
            if slaveLag < minLag {
                minLag = slaveLag
                bestSlave = slave
            }
        }
    }
    
    if bestSlave == nil {
        return nil, errors.New("no slave available for promotion")
    }
    
    return bestSlave, nil
}
```
