# database/sql基础

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---
## 📋 目录

- [database/sql基础](#databasesql基础)
  - [1. database/sql简介](#1-databasesql简介)
- [MySQL](#mysql)
- [PostgreSQL](#postgresql)
- [SQLite](#sqlite)
  - [2. 连接数据库](#2-连接数据库)
  - [3. 查询数据](#3-查询数据)
  - [4. 插入和更新](#4-插入和更新)
  - [5. 事务处理](#5-事务处理)
  - [6. 预处理语句](#6-预处理语句)
  - [7. 最佳实践](#7-最佳实践)
  - [🔗 相关资源](#相关资源)

---

## 1. database/sql简介

### 什么是database/sql

**database/sql** 是Go的标准数据库接口：

- 提供统一的数据库API
- 支持多种数据库驱动
- 连接池管理
- 预处理语句
- 事务支持

### 安装数据库驱动

```bash
# MySQL
go get -u github.com/go-sql-driver/mysql

# PostgreSQL
go get -u github.com/lib/pq

# SQLite
go get -u github.com/mattn/go-sqlite3
```

---

## 2. 连接数据库

### 打开连接

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // 连接字符串
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname"

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 验证连接
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to database")
}
```

---

### 连接池配置

```go
func setupDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    // 设置最大打开连接数
    db.SetMaxOpenConns(25)

    // 设置最大空闲连接数
    db.SetMaxIdleConns(5)

    // 设置连接最大存活时间
    db.SetConnMaxLifetime(5 * time.Minute)

    // 设置连接最大空闲时间
    db.SetConnMaxIdleTime(5 * time.Minute)

    // 验证连接
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
```

---

## 3. 查询数据

### 查询单行

```go
type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

func getUser(db *sql.DB, id int) (*User, error) {
    var user User

    query := "SELECT id, name, email, age FROM users WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Age)

    if err == sql.ErrNoRows {
        return nil, nil  // 没找到
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

---

### 查询多行

```go
func listUsers(db *sql.DB) ([]*User, error) {
    query := "SELECT id, name, email, age FROM users"
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // 重要：关闭rows

    var users []*User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    // 检查迭代错误
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}
```

---

### 带参数查询

```go
func searchUsers(db *sql.DB, name string, minAge int) ([]*User, error) {
    query := `
        SELECT id, name, email, age
        FROM users
        WHERE name LIKE ? AND age >= ?
    `

    rows, err := db.Query(query, "%"+name+"%", minAge)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    return users, rows.Err()
}
```

---

## 4. 插入和更新

### 插入数据

```go
func createUser(db *sql.DB, user *User) (int64, error) {
    query := `
        INSERT INTO users (name, email, age)
        VALUES (?, ?, ?)
    `

    result, err := db.Exec(query, user.Name, user.Email, user.Age)
    if err != nil {
        return 0, err
    }

    // 获取插入的ID
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return id, nil
}
```

---

### 更新数据

```go
func updateUser(db *sql.DB, user *User) error {
    query := `
        UPDATE users
        SET name = ?, email = ?, age = ?
        WHERE id = ?
    `

    result, err := db.Exec(query, user.Name, user.Email, user.Age, user.ID)
    if err != nil {
        return err
    }

    // 检查影响的行数
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("no rows affected")
    }

    return nil
}
```

---

### 删除数据

```go
func deleteUser(db *sql.DB, id int) error {
    query := "DELETE FROM users WHERE id = ?"

    result, err := db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("user not found")
    }

    return nil
}
```

---

## 5. 事务处理

### 基本事务

```go
func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    // 开始事务
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    // defer处理提交/回滚
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        } else if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    // 扣款
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
    if err != nil {
        return err
    }

    // 入账
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
    if err != nil {
        return err
    }

    return nil
}
```

---

### 使用Context的事务

```go
func transferMoneyWithContext(ctx Context.Context, db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 扣款
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance - ? WHERE id = ?",
        amount, fromID)
    if err != nil {
        return err
    }

    // 入账
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance + ? WHERE id = ?",
        amount, toID)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

---

## 6. 预处理语句

### 创建预处理语句

```go
func batchInsert(db *sql.DB, users []*User) error {
    // 准备语句
    stmt, err := db.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // 批量执行
    for _, user := range users {
        _, err := stmt.Exec(user.Name, user.Email, user.Age)
        if err != nil {
            return err
        }
    }

    return nil
}
```

---

### 查询预处理语句

```go
func getUsersByAge(db *sql.DB, ages []int) ([]*User, error) {
    stmt, err := db.Prepare("SELECT id, name, email, age FROM users WHERE age = ?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    var users []*User
    for _, age := range ages {
        rows, err := stmt.Query(age)
        if err != nil {
            return nil, err
        }

        for rows.Next() {
            var user User
            if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
                rows.Close()
                return nil, err
            }
            users = append(users, &user)
        }
        rows.Close()
    }

    return users, nil
}
```

---

## 7. 最佳实践

### 1. 使用Context

```go
// ✅ 推荐
func getUser(ctx Context.Context, db *sql.DB, id int) (*User, error) {
    var user User
    err := db.QueryRowContext(ctx,
        "SELECT id, name, email, age FROM users WHERE id = ?",
        id).Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    return &user, err
}
```

---

### 2. 总是关闭Rows

```go
// ✅ 推荐
rows, err := db.Query("SELECT * FROM users")
if err != nil {
    return err
}
defer rows.Close()  // 确保关闭

// ❌ 不推荐：忘记关闭会导致连接泄漏
```

---

### 3. 检查Scan错误

```go
// ✅ 推荐
for rows.Next() {
    if err := rows.Scan(&user.ID, &user.Name); err != nil {
        return nil, err
    }
}
if err := rows.Err(); err != nil {
    return nil, err
}
```

---

### 4. 使用sql.NullXxx处理NULL

```go
type User struct {
    ID    int
    Name  string
    Email sql.NullString  // 可能为NULL
    Age   sql.NullInt64   // 可能为NULL
}

func scanUser(rows *sql.Rows) (*User, error) {
    var user User
    err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    if err != nil {
        return nil, err
    }

    // 检查NULL
    if !user.Email.Valid {
        // Email为NULL
    }

    return &user, nil
}
```

---

### 5. 连接池管理

```go
// ✅ 推荐：合理配置连接池
db.SetMaxOpenConns(25)      // 最大打开连接数
db.SetMaxIdleConns(5)       // 最大空闲连接数
db.SetConnMaxLifetime(5 * time.Minute)   // 连接最大存活时间
db.SetConnMaxIdleTime(5 * time.Minute)   // 连接最大空闲时间
```

---

## 🔗 相关资源
