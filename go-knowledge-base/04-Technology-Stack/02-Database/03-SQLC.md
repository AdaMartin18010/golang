# SQLC - 类型安全 SQL

> **分类**: 开源技术堆栈

---

## 什么是 SQLC

SQLC 从 SQL 生成类型安全的 Go 代码，无需手动编写样板代码。

```
SQL 查询 → sqlc generate → Go 代码
```

---

## 配置

### sqlc.yaml

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "db"
        out: "db"
```

---

## SQL 查询

### query.sql

```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

---

## 生成代码使用

```go
import "database/sql"
import "github.com/user/project/db"

// 初始化
conn, _ := sql.Open("postgres", dsn)
queries := db.New(conn)

// 使用
user, err := queries.GetUser(ctx, 1)
if err != nil {
    log.Fatal(err)
}
fmt.Println(user.Name)

// 列表
users, err := queries.ListUsers(ctx)
for _, u := range users {
    fmt.Println(u.Name)
}

// 创建
user, err := queries.CreateUser(ctx, db.CreateUserParams{
    Name:  "Alice",
    Email: "alice@example.com",
})
```

---

## 优势

| 优势 | 说明 |
|------|------|
| 类型安全 | 编译时检查 |
| 性能 | 无反射开销 |
| 简单 | SQL 即代码 |
| 无依赖 | 标准库兼容 |
