# LD-029-Go-126-Features

> **Dimension**: 02-Language-Design
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Go 1.26
> **Size**: >20KB

---

## 1. Go 1.26 Release

Release date: February 10, 2026

## 2. Key Features

### 2.1 Green Tea GC (Default)

- 10-40% GC overhead reduction
- Page-based scanning
- Enabled by default

### 2.2 new() Expression

```go
ptr := new(int64(300))
```

### 2.3 Recursive Type Constraints

```go
type Adder[A Adder[A]] interface {
    Add(A) A
}
```

### 2.4 crypto/hpke

RFC 9180 implementation

### 2.5 errors.AsType

```go
if netErr, ok := errors.AsType[*NetworkError](err); ok {
    // Use netErr
}
```

---

## 3. Performance

| Operation | Improvement |
|-----------|-------------|
| GC | 10-40% |
| cgo | 30% |
| Small alloc | 30% |
| ReadAll | 2x |

---

## References

1. Go 1.26 Release Notes
2. Green Tea GC Design Doc

---

*Last Updated: 2026-04-03*
