# 通用数据库抽象层

框架级别的通用数据库抽象，支持多种数据库驱动（PostgreSQL、SQLite3、MySQL），提供统一的接口。

## 📋 功能特性

- ✅ **统一接口**: 提供统一的数据库操作接口
- ✅ **多驱动支持**: 支持 PostgreSQL、SQLite3、MySQL
- ✅ **连接池管理**: 自动管理连接池
- ✅ **事务支持**: 完整的事务支持
- ✅ **上下文支持**: 所有操作支持 Context
- ✅ **统计信息**: 提供连接池统计信息

## 🚀 快速开始

### 基本使用

```go
import "github.com/yourusername/golang/pkg/database"

// 创建 PostgreSQL 连接
db, err := database.NewDatabase(database.Config{
    Driver:       database.DriverPostgreSQL,
    DSN:          "postgres://user:password@localhost/dbname?sslmode=disable",
    MaxOpenConns: 25,
    MaxIdleConns: 5,
})
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// 执行查询
rows, err := db.Query(ctx, "SELECT id, name FROM users WHERE id = $1", 1)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

// 执行更新
result, err := db.Exec(ctx, "UPDATE users SET name = $1 WHERE id = $2", "New Name", 1)
if err != nil {
    log.Fatal(err)
}
```

### 使用事务

```go
// 开始事务
tx, err := db.Begin(ctx)
if err != nil {
    log.Fatal(err)
}
defer tx.Rollback()

// 在事务中执行操作
_, err = tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "John")
if err != nil {
    return err
}

// 提交事务
if err := tx.Commit(); err != nil {
    return err
}
```

### 切换数据库

```go
// PostgreSQL
db, _ := database.NewDatabase(database.Config{
    Driver: database.DriverPostgreSQL,
    DSN:    "postgres://...",
})

// SQLite3
db, _ := database.NewDatabase(database.Config{
    Driver: database.DriverSQLite3,
    DSN:    "file:app.db?cache=shared&mode=rwc",
})

// MySQL
db, _ := database.NewDatabase(database.Config{
    Driver: database.DriverMySQL,
    DSN:    "user:password@tcp(localhost:3306)/dbname",
})
```

## 📚 API 参考

### Database 接口

```go
type Database interface {
    Driver() Driver
    DB() *sql.DB
    Ping(ctx context.Context) error
    Close() error
    Stats() sql.DBStats
    Begin(ctx context.Context) (Transaction, error)
    Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
    Prepare(ctx context.Context, query string) (*sql.Stmt, error)
}
```

### Transaction 接口

```go
type Transaction interface {
    Commit() error
    Rollback() error
    Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
    Prepare(ctx context.Context, query string) (*sql.Stmt, error)
}
```

## 🔧 配置选项

```go
type Config struct {
    Driver          Driver        // 数据库驱动类型
    DSN             string        // 数据源名称
    MaxOpenConns    int           // 最大打开连接数
    MaxIdleConns    int           // 最大空闲连接数
    ConnMaxLifetime time.Duration // 连接最大生存时间
    ConnMaxIdleTime time.Duration // 连接最大空闲时间
    PingTimeout     time.Duration // Ping 超时时间
}
```

## 🎯 最佳实践

1. **使用 Context**: 所有操作都应该使用 Context 以支持超时和取消
2. **连接池配置**: 根据应用负载合理配置连接池参数
3. **事务管理**: 使用 defer 确保事务正确回滚
4. **错误处理**: 始终检查并处理错误
5. **资源清理**: 使用 defer 确保关闭连接和结果集

## 🔗 相关文档

- [数据库基础设施说明](../../internal/infrastructure/README.md#数据库实现)
