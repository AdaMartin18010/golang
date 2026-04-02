# 数据库连接 (Database Connectivity)

> **分类**: 开源技术堆栈

---

## database/sql

### 连接数据库

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

db, err := sql.Open("mysql", "user:password@/dbname")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// 验证连接
if err := db.Ping(); err != nil {
    log.Fatal(err)
}
```

---

## 基本操作

### 查询

```go
rows, err := db.Query("SELECT id, name FROM users WHERE age > ?", 18)
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        return err
    }
    fmt.Println(id, name)
}
```

### 执行

```go
result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 30)
if err != nil {
    return err
}

lastID, _ := result.LastInsertId()
rowsAffected, _ := result.RowsAffected()
```

---

## 连接池

```go
db.SetMaxOpenConns(25)      // 最大连接数
db.SetMaxIdleConns(5)       // 最大空闲连接
db.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生命周期
```

---

## 事务

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", 100, 1)
if err != nil {
    return err
}

_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", 100, 2)
if err != nil {
    return err
}

if err := tx.Commit(); err != nil {
    return err
}
```
