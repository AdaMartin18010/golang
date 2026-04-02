# 数据库 (Database)

本目录包含 Go 数据库技术栈。

## 文档列表

| 文档 | 内容 |
|------|------|
| [01-Database-Connectivity.md](01-Database-Connectivity.md) | 标准库 database/sql |
| [02-ORM-GORM.md](02-ORM-GORM.md) | GORM |
| [03-SQLC.md](03-SQLC.md) | SQLC |
| [04-Redis.md](04-Redis.md) | Redis |
| [05-MongoDB.md](05-MongoDB.md) | MongoDB |
| [06-ClickHouse.md](06-ClickHouse.md) | ClickHouse |

## 技术选择

| 场景 | 推荐 |
|------|------|
| SQL | database/sql + sqlc |
| ORM | GORM |
| 缓存 | go-redis |
| 文档 | mongo-driver |
| 分析 | clickhouse-go |
