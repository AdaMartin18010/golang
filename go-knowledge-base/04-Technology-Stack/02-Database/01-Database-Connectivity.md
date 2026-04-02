# TS-DB-001: Database Connectivity in Go

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #golang #database #sql #connection-pool #datasource
> **权威来源**:
>
> - [database/sql Package](https://golang.org/pkg/database/sql/) - Go standard library
> - [Go database/sql tutorial](http://go-database-sql.org/) - VividCortex
> - [SQL Injection Prevention](https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html) - OWASP

---

## 1. database/sql Architecture

### 1.1 Package Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    database/sql Architecture                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Application Code                                  │   │
│  │  db.Query()  db.Exec()  db.Prepare()  db.Begin()                   │   │
│  └───────────────────────────┬─────────────────────────────────────────┘   │
│                              │                                              │
│  ┌───────────────────────────▼─────────────────────────────────────────┐   │
│  │                      database/sql (stdlib)                           │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                    DB (Connection Pool)                        │  │   │
│  │  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐              │  │   │
│  │  │  │  Conn 1     │ │  Conn 2     │ │  Conn N     │              │  │   │
│  │  │  │ (Active)    │ │ (Idle)      │ │ (Active)    │              │  │   │
│  │  │  └─────────────┘ └─────────────┘ └─────────────┘              │  │   │
│  │  │                                                                      │  │   │
│  │  │  - MaxOpenConns (default: unlimited)                                │  │   │
│  │  │  - MaxIdleConns (default: 2)                                        │  │   │
│  │  │  - ConnMaxLifetime (default: unlimited)                             │  │   │
│  │  │  - ConnMaxIdleTime (default: unlimited)                             │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                     Tx (Transaction)                          │  │   │
│  │  │  - Bound to a single connection                               │  │   │
│  │  │  - Commit() / Rollback()                                      │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                     Stmt (Prepared Statement)                 │  │   │
│  │  │  - Compiled query cached on connection                        │  │   │
│  │  │  - Bind parameters for safety                                 │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                     Rows (Result Iterator)                    │  │   │
│  │  │  - Streaming result set                                       │  │   │
│  │  │  - Close() to release resources                               │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  └───────────────────────────┬─────────────────────────────────────────┘   │
│                              │                                              │
│  ┌───────────────────────────▼─────────────────────────────────────────┐   │
│  │                    database/driver (interface)                       │   │
│  │  - Driver interface                                                 │   │
│  │  - Conn interface                                                   │   │
│  │  - Stmt interface                                                   │   │
│  │  - Tx interface                                                     │   │
│  └───────────────────────────┬─────────────────────────────────────────┘   │
│                              │                                              │
│  ┌───────────────────────────▼─────────────────────────────────────────┐   │
│  │                    Driver Implementation                             │   │
│  │  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐        │   │
│  │  │  MySQL    │  │ PostgreSQL│  │  SQLite   │  │  Others   │        │   │
│  │  │  driver   │  │  driver   │  │  driver   │  │           │        │   │
│  │  └───────────┘  └───────────┘  └───────────┘  └───────────┘        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Connection Management

### 2.1 Database Connection

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
    _ "github.com/lib/pq"              // PostgreSQL driver
)

func createDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)                 // Maximum open connections
    db.SetMaxIdleConns(5)                  // Maximum idle connections
    db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection
    db.SetConnMaxIdleTime(1 * time.Minute) // Maximum idle time

    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return db, nil
}

// Proper shutdown
func shutdownDB(db *sql.DB) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.Close(); err != nil {
        log.Printf("Error closing database: %v", err)
    }
}
```

### 2.2 Connection Pool Tuning

```
Connection Pool Formula:

MaxOpenConns = (number of cores * 2) + number of disks

For typical web application:
- Small instance (1-2 cores): 10-20 connections
- Medium instance (4 cores): 25-50 connections
- Large instance (8+ cores): 50-100+ connections

MaxIdleConns:
- Usually MaxOpenConns / 4
- Should be >= expected concurrent idle connections

ConnMaxLifetime:
- MySQL wait_timeout - 1 minute buffer
- PostgreSQL: 1 hour or less
- Prevents stale connection issues

ConnMaxIdleTime (Go 1.15+):
- Usually 1-5 minutes
- Closes idle connections sooner
```

---

## 3. Query Execution

### 3.1 Basic Queries

```go
// Query (multiple rows)
func getUsers(ctx context.Context, db *sql.DB) ([]User, error) {
    rows, err := db.QueryContext(ctx,
        "SELECT id, name, email FROM users WHERE active = ?",
        true)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

// QueryRow (single row)
func getUserByID(ctx context.Context, db *sql.DB, id int64) (*User, error) {
    var u User
    err := db.QueryRowContext(ctx,
        "SELECT id, name, email FROM users WHERE id = ?",
        id).Scan(&u.ID, &u.Name, &u.Email)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &u, nil
}

// Exec (INSERT, UPDATE, DELETE)
func createUser(ctx context.Context, db *sql.DB, name, email string) (int64, error) {
    result, err := db.ExecContext(ctx,
        "INSERT INTO users (name, email) VALUES (?, ?)",
        name, email)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}
```

### 3.2 Prepared Statements

```go
// Prepare once, execute many times
func batchInsertUsers(ctx context.Context, db *sql.DB, users []User) error {
    stmt, err := db.PrepareContext(ctx,
        "INSERT INTO users (name, email) VALUES (?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, u := range users {
        _, err := stmt.ExecContext(ctx, u.Name, u.Email)
        if err != nil {
            return err
        }
    }

    return nil
}
```

---

## 4. Transactions

```go
// Transaction with proper rollback on error
func transferMoney(ctx context.Context, db *sql.DB, fromID, toID int64, amount float64) error {
    tx, err := db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
    })
    if err != nil {
        return err
    }
    defer tx.Rollback() // Safe to call even after Commit()

    // Deduct from sender
    result, err := tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance - ? WHERE id = ? AND balance >= ?",
        amount, fromID, amount)
    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return errors.New("insufficient funds or invalid account")
    }

    // Add to receiver
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

## 5. Best Practices

```
Database Connectivity Best Practices:
□ Use connection pooling
□ Set appropriate timeouts
□ Use context for cancellation
□ Always use prepared statements (or parameterized queries)
□ Check sql.ErrNoRows explicitly
□ Close rows after iteration
□ Use transactions for consistency
□ Handle null values properly
□ Don't ignore errors
□ Use struct scanning libraries (sqlx, gorm)
```

---

## 6. Checklist

```
Database Connectivity Checklist:
□ Connection pool configured
□ Ping on startup
□ Context timeouts set
□ Prepared statements used
□ Transactions for multi-step operations
□ Proper error handling
□ Resource cleanup (defer close)
□ Connection limits appropriate
□ Metrics/logging enabled
```
