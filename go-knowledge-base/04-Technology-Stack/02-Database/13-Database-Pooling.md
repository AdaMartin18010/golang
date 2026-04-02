# 数据库连接池 (Connection Pooling)

> **分类**: 开源技术堆栈  
> **标签**: #database #pooling #performance

---

## 连接池配置

```go
db, err := sql.Open("postgres", dsn)
if err != nil {
    log.Fatal(err)
}

// 连接池设置
db.SetMaxOpenConns(25)        // 最大打开连接数
db.SetMaxIdleConns(5)         // 最大空闲连接数
db.SetConnMaxLifetime(5 * time.Minute)   // 连接最大生命周期
db.SetConnMaxIdleTime(10 * time.Minute)  // 空闲连接最大时间
```

---

## 监控连接池

```go
type PoolStats struct {
    db *sql.DB
}

func (ps *PoolStats) Collect() map[string]interface{} {
    stats := ps.db.Stats()
    
    return map[string]interface{}{
        "open_connections":    stats.OpenConnections,
        "in_use":              stats.InUse,
        "idle":                stats.Idle,
        "wait_count":          stats.WaitCount,
        "wait_duration_ms":    stats.WaitDuration.Milliseconds(),
        "max_idle_closed":     stats.MaxIdleClosed,
        "max_lifetime_closed": stats.MaxLifetimeClosed,
    }
}

// 导出到 Prometheus
var (
    openConnections = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "db_open_connections",
        Help: "The number of established connections",
    })
    
    inUseConnections = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "db_in_use_connections",
        Help: "The number of connections currently in use",
    })
)

func recordDBMetrics(db *sql.DB) {
    stats := db.Stats()
    openConnections.Set(float64(stats.OpenConnections))
    inUseConnections.Set(float64(stats.InUse))
}
```

---

## 动态调整

```go
type AdaptivePool struct {
    db      *sql.DB
    minConn int
    maxConn int
    target  float64  // 目标使用率
}

func (ap *AdaptivePool) Adjust() {
    stats := ap.db.Stats()
    
    if stats.OpenConnections == 0 {
        return
    }
    
    usage := float64(stats.InUse) / float64(stats.OpenConnections)
    
    // 使用率过高，增加连接
    if usage > ap.target && stats.OpenConnections < ap.maxConn {
        newMax := int(float64(stats.OpenConnections) * 1.2)
        if newMax > ap.maxConn {
            newMax = ap.maxConn
        }
        ap.db.SetMaxOpenConns(newMax)
    }
    
    // 使用率过低，减少连接
    if usage < ap.target*0.5 && stats.OpenConnections > ap.minConn {
        newMax := int(float64(stats.OpenConnections) * 0.8)
        if newMax < ap.minConn {
            newMax = ap.minConn
        }
        ap.db.SetMaxOpenConns(newMax)
    }
}
```

---

## 连接池最佳实践

### 1. 预连接

```go
// 验证连接
func verifyPool(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return db.PingContext(ctx)
}

// 预热连接池
func warmupPool(db *sql.DB, n int) {
    conns := make([]*sql.Conn, 0, n)
    
    for i := 0; i < n; i++ {
        conn, err := db.Conn(context.Background())
        if err != nil {
            break
        }
        conns = append(conns, conn)
    }
    
    // 释放回池中
    for _, conn := range conns {
        conn.Close()
    }
}
```

### 2. 健康检查

```go
type PoolHealthChecker struct {
    db     *sql.DB
    ticker *time.Ticker
}

func (phc *PoolHealthChecker) Start() {
    phc.ticker = time.NewTicker(30 * time.Second)
    
    go func() {
        for range phc.ticker.C {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            
            if err := phc.db.PingContext(ctx); err != nil {
                log.Printf("Pool health check failed: %v", err)
            }
            
            cancel()
        }
    }()
}
```

---

## 连接池问题排查

```go
func DiagnosePool(db *sql.DB) {
    stats := db.Stats()
    
    // 检查等待
    if stats.WaitCount > 100 {
        log.Printf("WARNING: High wait count: %d", stats.WaitCount)
        log.Printf("Consider increasing MaxOpenConns (current: %d)", 
            stats.OpenConnections)
    }
    
    // 检查连接关闭
    if stats.MaxIdleClosed > 100 {
        log.Printf("WARNING: Many idle connections closed: %d", 
            stats.MaxIdleClosed)
        log.Printf("Consider increasing MaxIdleConns or ConnMaxIdleTime")
    }
    
    // 检查生命周期
    if stats.MaxLifetimeClosed > 100 {
        log.Printf("INFO: Many connections closed due to max lifetime: %d",
            stats.MaxLifetimeClosed)
    }
}
```
