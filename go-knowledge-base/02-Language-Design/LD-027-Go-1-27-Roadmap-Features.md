# LD-027-Go-1-27-Roadmap-Features

> **Dimension**: 02-Language-Design
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Go 1.27 Roadmap
> **Size**: >20KB

---

## 1. Go 1.27 Overview

Expected release: August 2026

Key features:

- Generic Methods (accepted)
- Green Tea GC finalized
- json/v2 GA preparation
- Goroutine leak detection default

---

## 2. Generic Methods

### 2.1 Proposal

Author: Robert Griesemer
Status: Accepted December 2025

### 2.2 Syntax

```go
type Container[T any] struct {
    items []T
}

func (c Container[T]) MapTo[U any](fn func(T) U) Container[U] {
    result := make([]U, len(c.items))
    for i, item := range c.items {
        result[i] = fn(item)
    }
    return Container[U]{items: result}
}
```

### 2.3 Database Example

```go
type DB struct{}

func (db *DB) QueryOne[T any](query string, args ...any) (*T, error) {
    // Implementation
}

// Usage
user, err := db.QueryOne[User]("SELECT * FROM users WHERE id = ?", 1)
```

---

## 3. Green Tea GC Finalization

GOEXPERIMENT=nogreenteagc removed
Only GC option in Go 1.27

---

## 4. json/v2 Status

Current: Experimental (GOEXPERIMENT=jsonv2)
Goal: GA in Go 1.27 or 1.28

Features:

- Streaming processing
- Better error messages
- Stricter defaults

---

## 5. References

1. Generic Methods Proposal
2. Go Release Cycle
3. json/v2 Design Doc

---

*Last Updated: 2026-04-03*
