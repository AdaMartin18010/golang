# 1. 🗄️ SQLite 深度解析

> **简介**: 本文档详细阐述了 SQLite 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 🗄️ SQLite 深度解析](#1-️-sqlite-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 数据库连接](#131-数据库连接)
    - [1.3.2 基础操作](#132-基础操作)
    - [1.3.3 事务处理](#133-事务处理)
    - [1.3.4 性能优化](#134-性能优化)
    - [1.3.5 并发控制](#135-并发控制)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 使用场景最佳实践](#141-使用场景最佳实践)
    - [1.4.2 性能优化最佳实践](#142-性能优化最佳实践)
    - [1.4.3 并发控制最佳实践](#143-并发控制最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**SQLite 是什么？**

SQLite 是一个轻量级的、嵌入式的、无服务器的 SQL 数据库引擎。

**核心特性**:

- ✅ **零配置**: 无需服务器，无需配置
- ✅ **轻量级**: 库文件小，资源占用低
- ✅ **文件数据库**: 数据库存储在单个文件中
- ✅ **ACID 事务**: 支持完整的 ACID 事务
- ✅ **跨平台**: 支持多种操作系统和架构

---

## 1.2 选型论证

**为什么选择 SQLite？**

**论证矩阵**:

| 评估维度 | 权重 | SQLite | PostgreSQL | MySQL | 说明 |
|---------|------|--------|-----------|-------|------|
| **轻量级** | 30% | 10 | 5 | 5 | SQLite 最轻量 |
| **易用性** | 25% | 10 | 7 | 7 | SQLite 零配置 |
| **性能** | 20% | 8 | 10 | 9 | SQLite 性能优秀 |
| **并发支持** | 15% | 6 | 10 | 10 | SQLite 并发较弱 |
| **功能完整性** | 10% | 7 | 10 | 10 | SQLite 功能完整 |
| **加权总分** | - | **8.50** | 7.75 | 7.60 | SQLite 得分最高（轻量级场景） |

**核心优势**:

1. **轻量级（权重 30%）**:
   - 库文件小，资源占用低
   - 适合嵌入式应用和移动应用
   - 无需独立的数据库服务器

2. **易用性（权重 25%）**:
   - 零配置，开箱即用
   - 数据库存储在单个文件中
   - 部署简单，无需维护

3. **性能（权重 20%）**:
   - 对于单用户或低并发场景性能优秀
   - 读写速度快
   - 适合中小型应用

**为什么不选择其他数据库？**

1. **PostgreSQL**:
   - ✅ 功能强大，并发支持好
   - ❌ 需要独立的数据库服务器
   - ❌ 配置和维护复杂
   - ❌ 不适合嵌入式场景

2. **MySQL**:
   - ✅ 功能丰富，生态成熟
   - ❌ 需要独立的数据库服务器
   - ❌ 配置和维护复杂
   - ❌ 不适合嵌入式场景

**适用场景**:

- ✅ 嵌入式应用
- ✅ 移动应用
- ✅ 小型 Web 应用
- ✅ 开发和测试环境
- ✅ 单用户应用
- ✅ 数据分析和报表

**不适用场景**:

- ❌ 高并发 Web 应用
- ❌ 多用户同时写入
- ❌ 需要复杂网络访问
- ❌ 大规模数据存储

---

## 1.3 实际应用

### 1.3.1 数据库连接

**使用 go-sqlite3 连接**:

```go
// internal/infrastructure/database/sqlite/client.go
package sqlite

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type Client struct {
    db *sql.DB
}

func NewClient(dbPath string) (*Client, error) {
    // 连接字符串
    dsn := dbPath + "?_journal_mode=WAL&_foreign_keys=1"

    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, err
    }

    // 配置连接池
    db.SetMaxOpenConns(1)  // SQLite 建议单连接
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0)

    // 测试连接
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return &Client{db: db}, nil
}

func (c *Client) Close() error {
    return c.db.Close()
}

func (c *Client) DB() *sql.DB {
    return c.db
}
```

**使用 Ent ORM 连接**:

```go
// 使用 Ent ORM 连接 SQLite
import (
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/sql"
    _ "github.com/mattn/go-sqlite3"
)

func NewEntClient(dbPath string) (*ent.Client, error) {
    drv, err := sql.Open(dialect.SQLite, dbPath+"?_fk=1")
    if err != nil {
        return nil, err
    }

    client := ent.NewClient(ent.Driver(drv))
    return client, nil
}
```

### 1.3.2 基础操作

**创建表**:

```go
// 创建表
func (c *Client) CreateTable(ctx context.Context) error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )
    `

    _, err := c.db.ExecContext(ctx, query)
    return err
}
```

**插入数据**:

```go
// 插入数据
func (c *Client) CreateUser(ctx context.Context, email, name string) (int64, error) {
    query := `INSERT INTO users (email, name) VALUES (?, ?)`

    result, err := c.db.ExecContext(ctx, query, email, name)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}
```

**查询数据**:

```go
// 查询数据
type User struct {
    ID        int64
    Email     string
    Name      string
    CreatedAt time.Time
}

func (c *Client) GetUser(ctx context.Context, id int64) (*User, error) {
    query := `SELECT id, email, name, created_at FROM users WHERE id = ?`

    var user User
    err := c.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.CreatedAt,
    )
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

### 1.3.3 事务处理

**基本事务**:

```go
// 事务处理
func (c *Client) CreateUserWithProfile(ctx context.Context, email, name string) error {
    tx, err := c.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 插入用户
    result, err := tx.ExecContext(ctx,
        "INSERT INTO users (email, name) VALUES (?, ?)",
        email, name,
    )
    if err != nil {
        return err
    }

    userID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // 插入用户配置
    _, err = tx.ExecContext(ctx,
        "INSERT INTO user_profiles (user_id, settings) VALUES (?, ?)",
        userID, "{}",
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

### 1.3.4 性能优化

**WAL 模式**:

```go
// 启用 WAL 模式提高并发性能
func (c *Client) EnableWAL() error {
    _, err := c.db.Exec("PRAGMA journal_mode=WAL")
    return err
}
```

**优化设置**:

```go
// 性能优化设置
func (c *Client) Optimize() error {
    optimizations := []string{
        "PRAGMA journal_mode=WAL",           // WAL 模式
        "PRAGMA synchronous=NORMAL",         // 同步模式
        "PRAGMA cache_size=-64000",          // 缓存大小（64MB）
        "PRAGMA foreign_keys=ON",            // 外键约束
        "PRAGMA temp_store=MEMORY",          // 临时存储
    }

    for _, pragma := range optimizations {
        if _, err := c.db.Exec(pragma); err != nil {
            return err
        }
    }

    return nil
}
```

### 1.3.5 并发控制

**使用 WAL 模式支持并发读取**:

```go
// WAL 模式支持多读一写
func (c *Client) InitWithWAL(dbPath string) error {
    dsn := dbPath + "?_journal_mode=WAL&_foreign_keys=1&_busy_timeout=5000"

    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return err
    }

    c.db = db
    return nil
}
```

**使用文件锁控制并发写入**:

```go
// 使用文件锁控制并发写入
import (
    "os"
    "syscall"
)

func (c *Client) LockForWrite() error {
    file, err := os.OpenFile(c.dbPath+".lock", os.O_CREATE|os.O_EXCL, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // 获取排他锁
    return syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
}
```

---

## 1.4 最佳实践

### 1.4.1 使用场景最佳实践

**为什么需要合理选择使用场景？**

合理选择使用场景可以充分发挥 SQLite 的优势，避免其局限性。

**使用场景选择原则**:

1. **适合场景**:
   - 嵌入式应用和移动应用
   - 小型 Web 应用（低并发）
   - 开发和测试环境
   - 单用户应用
   - 数据分析和报表

2. **不适合场景**:
   - 高并发 Web 应用
   - 多用户同时写入
   - 需要复杂网络访问
   - 大规模数据存储

**实际应用示例**:

```go
// 使用场景判断
func ShouldUseSQLite(concurrentUsers int, dataSize int64) bool {
    // 并发用户数少于 100
    if concurrentUsers > 100 {
        return false
    }

    // 数据大小少于 100GB
    if dataSize > 100*1024*1024*1024 {
        return false
    }

    // 主要是读取操作
    // 写入操作较少

    return true
}
```

**最佳实践要点**:

1. **并发控制**: SQLite 适合低并发场景，高并发应使用 PostgreSQL 或 MySQL
2. **数据大小**: 适合中小型数据，大规模数据应使用其他数据库
3. **网络访问**: 适合本地访问，网络访问应使用客户端-服务器数据库

### 1.4.2 性能优化最佳实践

**为什么需要性能优化？**

合理的性能优化可以提高 SQLite 的读写性能，减少资源占用。

**性能优化原则**:

1. **启用 WAL 模式**: 提高并发读取性能
2. **调整同步模式**: 平衡性能和可靠性
3. **设置缓存大小**: 提高查询性能
4. **使用索引**: 加速查询
5. **批量操作**: 减少事务开销

**实际应用示例**:

```go
// 性能优化最佳实践
func (c *Client) OptimizeForPerformance() error {
    // 1. 启用 WAL 模式
    if _, err := c.db.Exec("PRAGMA journal_mode=WAL"); err != nil {
        return err
    }

    // 2. 设置同步模式（NORMAL 平衡性能和可靠性）
    if _, err := c.db.Exec("PRAGMA synchronous=NORMAL"); err != nil {
        return err
    }

    // 3. 设置缓存大小（64MB）
    if _, err := c.db.Exec("PRAGMA cache_size=-64000"); err != nil {
        return err
    }

    // 4. 启用外键约束
    if _, err := c.db.Exec("PRAGMA foreign_keys=ON"); err != nil {
        return err
    }

    // 5. 临时存储使用内存
    if _, err := c.db.Exec("PRAGMA temp_store=MEMORY"); err != nil {
        return err
    }

    return nil
}
```

**最佳实践要点**:

1. **WAL 模式**: 启用 WAL 模式可以提高并发读取性能
2. **同步模式**: 根据场景选择合适的同步模式（FULL/NORMAL/OFF）
3. **缓存大小**: 根据可用内存设置合适的缓存大小
4. **索引优化**: 为常用查询字段创建索引
5. **批量操作**: 使用事务批量操作减少开销

### 1.4.3 并发控制最佳实践

**为什么需要并发控制？**

SQLite 的并发能力有限，需要合理控制并发访问。

**并发控制原则**:

1. **WAL 模式**: 支持多读一写
2. **连接池**: 使用单连接或小连接池
3. **文件锁**: 使用文件锁控制并发写入
4. **超时设置**: 设置合理的超时时间
5. **重试机制**: 实现重试机制处理锁冲突

**实际应用示例**:

```go
// 并发控制最佳实践
func (c *Client) InitWithConcurrencyControl(dbPath string) error {
    // 使用 WAL 模式和超时设置
    dsn := dbPath + "?_journal_mode=WAL&_foreign_keys=1&_busy_timeout=5000"

    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return err
    }

    // SQLite 建议使用单连接或小连接池
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0)

    c.db = db
    return nil
}

// 带重试的写入操作
func (c *Client) WriteWithRetry(ctx context.Context, query string, args ...interface{}) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        _, err := c.db.ExecContext(ctx, query, args...)
        if err == nil {
            return nil
        }

        // 检查是否是锁冲突
        if strings.Contains(err.Error(), "database is locked") {
            if i < maxRetries-1 {
                time.Sleep(time.Millisecond * time.Duration(100*(i+1)))
                continue
            }
        }

        return err
    }

    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

**最佳实践要点**:

1. **WAL 模式**: 启用 WAL 模式支持多读一写
2. **连接池**: 使用单连接或小连接池（SQLite 建议）
3. **超时设置**: 设置合理的 busy_timeout 处理锁冲突
4. **重试机制**: 实现重试机制处理临时锁冲突
5. **读写分离**: 尽量分离读写操作，减少锁竞争

---

## 📚 扩展阅读

- [SQLite 官方文档](https://www.sqlite.org/docs.html)
- [go-sqlite3 官方文档](https://github.com/mattn/go-sqlite3)
- [Ent ORM SQLite 支持](https://entgo.io/docs/dialects/#sqlite)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 SQLite 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
