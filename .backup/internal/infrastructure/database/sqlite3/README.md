# SQLite3 数据库支持

SQLite3 数据库连接管理实现。

## 功能特性

- ✅ 连接池管理
- ✅ 连接参数配置
- ✅ 连接健康检查
- ✅ 支持 WAL 模式
- ✅ 支持内存数据库

## 使用示例

### 基本使用

```go
import (
    "github.com/yourusername/golang/internal/infrastructure/database/sqlite3"
)

// 使用默认配置
db, err := sqlite3.NewConnection(sqlite3.DefaultConfig())
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### 自定义配置

```go
config := sqlite3.Config{
    DSN:             "file:app.db?cache=shared&mode=rwc&_journal_mode=WAL",
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
    ConnMaxIdleTime: 10 * time.Minute,
}

db, err := sqlite3.NewConnection(config)
```

### 使用 DSN

```go
// 使用 DSN 字符串
db, err := sqlite3.NewConnectionWithDSN("file:app.db?cache=shared&mode=rwc")
```

### 与 Ent 集成

```go
import (
    "entgo.io/ent/dialect"
    entsql "entgo.io/ent/dialect/sql"
    _ "github.com/mattn/go-sqlite3"

    "github.com/yourusername/golang/internal/infrastructure/database/ent"
)

// 创建 SQLite3 连接
db, err := sqlite3.NewConnection(sqlite3.DefaultConfig())
if err != nil {
    log.Fatal(err)
}

// 创建 Ent 驱动
drv := entsql.OpenDB(dialect.SQLite, db)

// 创建 Ent 客户端
client := ent.NewClient(ent.Driver(drv))
```

## DSN 参数说明

- `file:path/to/db.db` - 数据库文件路径
- `cache=shared` - 共享缓存模式（推荐）
- `mode=rwc` - 读写创建模式
- `_journal_mode=WAL` - WAL 模式（推荐，提升性能）
- `_foreign_keys=1` - 启用外键约束
- `_busy_timeout=5000` - 忙等待超时（毫秒）

## 性能优化建议

1. **使用 WAL 模式**: `_journal_mode=WAL` 可以显著提升并发性能
2. **共享缓存**: `cache=shared` 允许多个连接共享缓存
3. **连接池**: 合理设置连接池大小
4. **预编译语句**: 使用预编译语句提升性能

## 相关资源

- [SQLite 官方文档](https://www.sqlite.org/docs.html)
- [go-sqlite3 文档](https://github.com/mattn/go-sqlite3)
