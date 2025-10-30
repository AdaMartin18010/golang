# 02-PostgreSQL编程

> Go语言PostgreSQL数据库编程完全指南

---

## 📋 目录

- [02-PostgreSQL编程](#02-postgresql编程)
  - [📋 目录](#-目录)
  - [📚 章节概览](#-章节概览)
  - [1. 环境准备](#1-环境准备)
    - [1.1 安装驱动](#11-安装驱动)
    - [1.2 基本连接](#12-基本连接)
      - [使用pq驱动](#使用pq驱动)
  - [2. pgx驱动使用](#2-pgx驱动使用)
    - [2.1 连接池配置](#21-连接池配置)
      - [pgx连接池架构可视化](#pgx连接池架构可视化)
      - [连接获取与释放流程](#连接获取与释放流程)
    - [2.2 基本查询](#22-基本查询)
  - [3. CRUD操作](#3-crud操作)
    - [3.1 插入数据](#31-插入数据)
    - [3.2 批量插入](#32-批量插入)
    - [3.3 查询数据](#33-查询数据)
    - [3.4 更新和删除](#34-更新和删除)
  - [4. 高级特性](#4-高级特性)
    - [4.1 JSON/JSONB支持](#41-jsonjsonb支持)
    - [4.2 数组类型](#42-数组类型)
    - [4.3 全文搜索](#43-全文搜索)
  - [5. 事务处理](#5-事务处理)
    - [5.1 基本事务](#51-基本事务)
    - [5.2 Savepoint](#52-savepoint)
  - [6. 性能优化](#6-性能优化)
    - [6.1 预处理语句](#61-预处理语句)
    - [6.2 批量操作](#62-批量操作)
  - [💡 最佳实践](#-最佳实践)
    - [1. 驱动选择](#1-驱动选择)
    - [2. 连接管理](#2-连接管理)
    - [3. 性能优化](#3-性能优化)
    - [4. 安全性](#4-安全性)
  - [🔗 相关章节](#-相关章节)

## 📚 章节概览

本章介绍如何在Go语言中使用PostgreSQL数据库，包括pgx驱动的使用、JSON支持、数组类型和全文搜索等PostgreSQL特有功能。

---

## 1. 环境准备

### 1.1 安装驱动

PostgreSQL有两个主流驱动：

```bash
# pq驱动 (database/sql兼容)
go get github.com/lib/pq

# pgx驱动 (推荐，功能更强大)
go get github.com/jackc/pgx/v5
```

### 1.2 基本连接

#### 使用pq驱动

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

func main() {
    // DSN格式
    dsn := "host=localhost port=5432 user=postgres password=secret dbname=testdb sslmode=disable"

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 验证连接
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    fmt.Println("PostgreSQL连接成功!")
}
```

---

## 2. pgx驱动使用

### 2.1 连接池配置

#### pgx连接池架构可视化

```mermaid
graph TB
    subgraph "应用层"
        App1[Goroutine 1]
        App2[Goroutine 2]
        App3[Goroutine 3]
        App4[Goroutine 4]
    end

    subgraph "pgxpool连接池"
        Pool[连接池管理器<br/>MaxConns=25<br/>MinConns=5]

        subgraph "空闲连接"
            Idle1[Conn 1]
            Idle2[Conn 2]
            Idle3[Conn 3]
        end

        subgraph "使用中连接"
            Busy1[Conn 4 - 查询中]
            Busy2[Conn 5 - 事务中]
        end

        WaitQueue[等待队列<br/>Goroutine Queue]
    end

    subgraph "PostgreSQL服务器"
        PG[(PostgreSQL<br/>Database)]
    end

    App1 -->|获取连接| Pool
    App2 -->|获取连接| Pool
    App3 -->|请求连接| WaitQueue
    App4 -->|请求连接| WaitQueue

    Pool -->|分配| Idle1
    Pool -->|分配| Idle2
    Pool -->|连接满，加入队列| WaitQueue

    Busy1 -->|释放| Idle1
    Busy2 -->|释放| Idle2

    Idle1 -.TCP连接.-> PG
    Idle2 -.TCP连接.-> PG
    Idle3 -.TCP连接.-> PG
    Busy1 -.TCP连接.-> PG
    Busy2 -.TCP连接.-> PG

    style Pool fill:#e1ffe1
    style WaitQueue fill:#ffe1e1
    style PG fill:#e1f5ff
    style Idle1 fill:#fff4e1
    style Idle2 fill:#fff4e1
    style Idle3 fill:#fff4e1
    style Busy1 fill:#ffe1e1
    style Busy2 fill:#ffe1e1
```

#### 连接获取与释放流程

```mermaid
sequenceDiagram
    participant App as 应用Goroutine
    participant Pool as pgxpool.Pool
    participant Conn as pgx.Conn
    participant PG as PostgreSQL

    Note over App,PG: 连接获取流程

    App->>Pool: pool.Acquire(ctx)

    alt 有空闲连接
        Pool->>Conn: 分配空闲连接
        Pool-->>App: 返回连接
        Note over Conn: 状态: 使用中
    else 无空闲连接但未达MaxConns
        Pool->>PG: 创建新连接
        PG-->>Pool: 连接建立
        Pool-->>App: 返回新连接
    else 连接池已满
        Note over App: 等待其他连接释放
        Pool->>App: 阻塞等待
    end

    Note over App,PG: 执行查询

    App->>Conn: Query(sql, args)
    Conn->>PG: 发送SQL
    PG-->>Conn: 返回结果集
    Conn-->>App: Rows

    App->>App: 处理结果

    Note over App,PG: 连接释放流程

    App->>Conn: conn.Release()
    Conn->>Pool: 归还连接
    Note over Conn: 状态: 空闲

    opt 有等待的Goroutine
        Pool->>App: 唤醒等待者
    end
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

func initDB() *pgxpool.Pool {
    // DSN格式
    dsn := "postgres://postgres:secret@localhost:5432/testdb?sslmode=disable"

    // 解析配置
    config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        log.Fatal(err)
    }

    // 设置连接池参数
    config.MaxConns = 25  // 最大连接数
    config.MinConns = 5   // 最小连接数（预热）

    // 创建连接池
    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatal(err)
    }

    return pool
}

func main() {
    pool := initDB()
    defer pool.Close()

    // 测试连接
    var greeting string
    err := pool.QueryRow(context.Background(), "SELECT 'Hello PostgreSQL!'").Scan(&greeting)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(greeting)
}
```

### 2.2 基本查询

```go
package main

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
    ID       int
    Username string
    Email    string
    Age      int
}

// 查询单行
func getUserByID(pool *pgxpool.Pool, id int) (*User, error) {
    query := "SELECT id, username, email, age FROM users WHERE id = $1"

    user := &User{}
    err := pool.QueryRow(context.Background(), query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Age,
    )

    if err != nil {
        return nil, err
    }

    return user, nil
}

// 查询多行
func getAllUsers(pool *pgxpool.Pool) ([]User, error) {
    query := "SELECT id, username, email, age FROM users"

    rows, err := pool.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Age)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, rows.Err()
}
```

---

## 3. CRUD操作

### 3.1 插入数据

```go
// 插入并返回ID
func insertUser(pool *pgxpool.Pool, user User) (int, error) {
    query := `
        INSERT INTO users(username, email, age)
        VALUES($1, $2, $3)
        RETURNING id
    `

    var id int
    err := pool.QueryRow(
        context.Background(),
        query,
        user.Username,
        user.Email,
        user.Age,
    ).Scan(&id)

    return id, err
}
```

### 3.2 批量插入

```go
package main

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// 使用CopyFrom批量插入（最高效）
func batchInsertUsers(pool *pgxpool.Pool, users []User) error {
    // 准备数据
    rows := make([][]interface{}, len(users))
    for i, user := range users {
        rows[i] = []interface{}{user.Username, user.Email, user.Age}
    }

    // 使用CopyFrom
    _, err := pool.CopyFrom(
        context.Background(),
        pgx.Identifier{"users"},
        []string{"username", "email", "age"},
        pgx.CopyFromRows(rows),
    )

    return err
}

// 使用Batch批量操作
func batchInsertWithBatch(pool *pgxpool.Pool, users []User) error {
    batch := &pgx.Batch{}

    for _, user := range users {
        batch.Queue(
            "INSERT INTO users(username, email, age) VALUES($1, $2, $3)",
            user.Username,
            user.Email,
            user.Age,
        )
    }

    results := pool.SendBatch(context.Background(), batch)
    defer results.Close()

    // 处理所有结果
    for i := 0; i < len(users); i++ {
        _, err := results.Exec()
        if err != nil {
            return err
        }
    }

    return nil
}
```

### 3.3 查询数据

```go
// 使用命名参数（pgx v5）
func getUsersByAge(pool *pgxpool.Pool, minAge, maxAge int) ([]User, error) {
    query := `
        SELECT id, username, email, age
        FROM users
        WHERE age BETWEEN $1 AND $2
        ORDER BY age
    `

    rows, err := pool.Query(context.Background(), query, minAge, maxAge)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Age)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}
```

### 3.4 更新和删除

```go
// 更新用户
func updateUser(pool *pgxpool.Pool, user User) error {
    query := "UPDATE users SET username=$1, email=$2, age=$3 WHERE id=$4"

    tag, err := pool.Exec(
        context.Background(),
        query,
        user.Username,
        user.Email,
        user.Age,
        user.ID,
    )

    if err != nil {
        return err
    }

    if tag.RowsAffected() == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}

// 删除用户
func deleteUser(pool *pgxpool.Pool, id int) error {
    query := "DELETE FROM users WHERE id=$1"

    tag, err := pool.Exec(context.Background(), query, id)
    if err != nil {
        return err
    }

    if tag.RowsAffected() == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}
```

---

## 4. 高级特性

### 4.1 JSON/JSONB支持

```go
package main

import (
    "context"
    "encoding/json"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
    Bio      string   `json:"bio"`
    Hobbies  []string `json:"hobbies"`
    Settings map[string]interface{} `json:"settings"`
}

type UserWithProfile struct {
    ID       int
    Username string
    Profile  Profile
}

// 插入JSON数据
func insertUserWithProfile(pool *pgxpool.Pool, user UserWithProfile) error {
    query := `
        INSERT INTO users(username, profile)
        VALUES($1, $2)
    `

    profileJSON, err := json.Marshal(user.Profile)
    if err != nil {
        return err
    }

    _, err = pool.Exec(context.Background(), query, user.Username, profileJSON)
    return err
}

// 查询JSON数据
func getUserProfile(pool *pgxpool.Pool, id int) (*UserWithProfile, error) {
    query := "SELECT id, username, profile FROM users WHERE id = $1"

    user := &UserWithProfile{}
    var profileJSON []byte

    err := pool.QueryRow(context.Background(), query, id).Scan(
        &user.ID,
        &user.Username,
        &profileJSON,
    )
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(profileJSON, &user.Profile); err != nil {
        return nil, err
    }

    return user, nil
}

// 使用JSONB查询
func searchUsersByHobby(pool *pgxpool.Pool, hobby string) ([]UserWithProfile, error) {
    query := `
        SELECT id, username, profile
        FROM users
        WHERE profile->'hobbies' ? $1
    `

    rows, err := pool.Query(context.Background(), query, hobby)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []UserWithProfile
    for rows.Next() {
        var user UserWithProfile
        var profileJSON []byte

        err := rows.Scan(&user.ID, &user.Username, &profileJSON)
        if err != nil {
            return nil, err
        }

        json.Unmarshal(profileJSON, &user.Profile)
        users = append(users, user)
    }

    return users, nil
}
```

### 4.2 数组类型

```go
package main

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/lib/pq"
)

// 使用数组类型
func insertTags(pool *pgxpool.Pool, articleID int, tags []string) error {
    query := "UPDATE articles SET tags = $1 WHERE id = $2"

    _, err := pool.Exec(context.Background(), query, tags, articleID)
    return err
}

// 查询包含某个标签的文章
func getArticlesByTag(pool *pgxpool.Pool, tag string) ([]int, error) {
    query := "SELECT id FROM articles WHERE $1 = ANY(tags)"

    rows, err := pool.Query(context.Background(), query, tag)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ids []int
    for rows.Next() {
        var id int
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        ids = append(ids, id)
    }

    return ids, nil
}
```

### 4.3 全文搜索

```go
// 创建全文搜索索引
func createFullTextIndex(pool *pgxpool.Pool) error {
    queries := []string{
        "ALTER TABLE articles ADD COLUMN tsv tsvector",
        "UPDATE articles SET tsv = to_tsvector('english', title || ' ' || content)",
        "CREATE INDEX idx_articles_tsv ON articles USING GIN(tsv)",
    }

    for _, query := range queries {
        if _, err := pool.Exec(context.Background(), query); err != nil {
            return err
        }
    }

    return nil
}

// 全文搜索
func searchArticles(pool *pgxpool.Pool, keyword string) ([]Article, error) {
    query := `
        SELECT id, title, content
        FROM articles
        WHERE tsv @@ to_tsquery('english', $1)
        ORDER BY ts_rank(tsv, to_tsquery('english', $1)) DESC
    `

    rows, err := pool.Query(context.Background(), query, keyword)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var articles []Article
    for rows.Next() {
        var article Article
        err := rows.Scan(&article.ID, &article.Title, &article.Content)
        if err != nil {
            return nil, err
        }
        articles = append(articles, article)
    }

    return articles, nil
}
```

---

## 5. 事务处理

### 5.1 基本事务

```go
package main

import (
    "context"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// 使用事务
func transferMoney(pool *pgxpool.Pool, fromID, toID int, amount float64) error {
    // 开始事务
    tx, err := pool.Begin(context.Background())
    if err != nil {
        return err
    }
    defer tx.Rollback(context.Background())

    // 扣款
    _, err = tx.Exec(
        context.Background(),
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount,
        fromID,
    )
    if err != nil {
        return err
    }

    // 加款
    _, err = tx.Exec(
        context.Background(),
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount,
        toID,
    )
    if err != nil {
        return err
    }

    // 提交事务
    return tx.Commit(context.Background())
}
```

### 5.2 Savepoint

```go
// 使用Savepoint
func complexTransaction(pool *pgxpool.Pool) error {
    tx, err := pool.Begin(context.Background())
    if err != nil {
        return err
    }
    defer tx.Rollback(context.Background())

    // 第一部分操作
    _, err = tx.Exec(context.Background(), "INSERT INTO log(message) VALUES('step 1')")
    if err != nil {
        return err
    }

    // 创建savepoint
    _, err = tx.Exec(context.Background(), "SAVEPOINT sp1")
    if err != nil {
        return err
    }

    // 第二部分操作（可能失败）
    _, err = tx.Exec(context.Background(), "INSERT INTO users(username) VALUES('test')")
    if err != nil {
        // 回滚到savepoint
        tx.Exec(context.Background(), "ROLLBACK TO SAVEPOINT sp1")
    }

    // 继续其他操作
    _, err = tx.Exec(context.Background(), "INSERT INTO log(message) VALUES('step 3')")
    if err != nil {
        return err
    }

    return tx.Commit(context.Background())
}
```

---

## 6. 性能优化

### 6.1 预处理语句

```go
// 使用预处理语句
func batchInsertPrepared(pool *pgxpool.Pool, users []User) error {
    ctx := context.Background()

    // 准备语句
    _, err := pool.Exec(ctx, "PREPARE insert_user AS INSERT INTO users(username, email, age) VALUES($1, $2, $3)")
    if err != nil {
        return err
    }
    defer pool.Exec(ctx, "DEALLOCATE insert_user")

    // 批量执行
    for _, user := range users {
        _, err := pool.Exec(ctx, "EXECUTE insert_user($1, $2, $3)", user.Username, user.Email, user.Age)
        if err != nil {
            return err
        }
    }

    return nil
}
```

### 6.2 批量操作

```go
// 使用UNNEST进行批量更新
func batchUpdate(pool *pgxpool.Pool, updates map[int]string) error {
    if len(updates) == 0 {
        return nil
    }

    ids := make([]int, 0, len(updates))
    usernames := make([]string, 0, len(updates))

    for id, username := range updates {
        ids = append(ids, id)
        usernames = append(usernames, username)
    }

    query := `
        UPDATE users
        SET username = u.username
        FROM UNNEST($1::int[], $2::text[]) AS u(id, username)
        WHERE users.id = u.id
    `

    _, err := pool.Exec(context.Background(), query, ids, usernames)
    return err
}
```

---

## 💡 最佳实践

### 1. 驱动选择

- ✅ 推荐使用pgx驱动（性能更好，功能更强）
- ✅ pgx支持PostgreSQL特有功能
- ✅ pgx提供更好的类型支持

### 2. 连接管理

- ✅ 使用连接池管理连接
- ✅ 合理设置MaxConns和MinConns
- ✅ 使用context控制超时

### 3. 性能优化

- ✅ 使用CopyFrom进行批量插入
- ✅ 利用JSONB索引加速JSON查询
- ✅ 使用GiN索引支持全文搜索
- ✅ 合理使用数组类型

### 4. 安全性

- ✅ 始终使用参数化查询
- ✅ 使用SSL连接
- ✅ 最小权限原则
- ✅ 定期备份数据

---

## 🔗 相关章节

- [01-MySQL编程](01-MySQL编程.md)
- [03-Redis编程](03-Redis编程.md)
- [04-ORM框架-GORM](../01-语言基础/README.md)

---

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3
