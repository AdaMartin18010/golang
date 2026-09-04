# LD-027-Go-1-27-Roadmap-Features

> **Dimension**: 02-Language-Design
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Go 1.27 Roadmap
> **Size**: >20KB

> **维度**: Language Design
> **级别**: B (1 KB)
> **标签**: #ld
> **Go 版本**: 1.27+
> **Bloom 层级**: L2   <!-- L2 理解层：roadmap 与提案状态跟踪 -->
> **前置概念**: [LD-026 Go 1.26 特性](LD-026-Go-126-New-Features.md) · [03-Evolution/06 提案流程](03-Evolution/06-Proposal-Process.md) | **后置概念**: [LD-037 泛型方法](LD-037-Go-1.27-Generic-Methods.md) · [03-Evolution/07](03-Evolution/07-Go126-to-Go127.md)
> **定理链**: Roadmap 提案 → Accepted 状态 → 语言/库/工具链变更清单 → 升级决策
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
