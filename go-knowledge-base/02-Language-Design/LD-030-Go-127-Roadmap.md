# LD-030-Go-127-Roadmap

> **Dimension**: 02-Language-Design
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Go 1.27 Roadmap
> **Size**: >20KB

---

## 1. Go 1.27 Overview

Expected: August 2026

## 2. Planned Features

### 2.1 Generic Methods

```go
func (db *DB) QueryOne[T any](query string) (*T, error)
```

### 2.2 Green Tea GC Final

Only GC option, nogreenteagc removed

### 2.3 json/v2

GA preparation

---

## 3. Platform Changes

- macOS 12 support removed
- PowerPC ELFv2

---

## References

1. Go Roadmap
2. Generic Methods Proposal

---

*Last Updated: 2026-04-03*
