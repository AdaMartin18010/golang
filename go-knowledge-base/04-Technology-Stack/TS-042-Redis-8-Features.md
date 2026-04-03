# TS-042-Redis-8-Features

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level Academic
> **Created**: 2026-04-03
> **Version**: Redis 8.6
> **Size**: >20KB

---

## 1. Redis 8 Overview

8.0: May 2025
8.6: February 2026

## 2. New Data Structures

### 2.1 Vector Set

```bash
VADD my_vectors vec1 [0.1, 0.2, 0.3]
VSIM my_vectors VEC [0.15, 0.25, 0.35] K 10
```

### 2.2 JSON

```bash
JSON.SET user:1 $ '{"name": "Alice", "age": 30}'
JSON.GET user:1 $.name
```

### 2.3 Time Series

```bash
TS.CREATE sensor TEMPERATURE
TS.ADD sensor * 23.5
TS.RANGE sensor - + AGGREGATION avg 60000
```

## 3. Performance

- 87% faster command execution
- 112% throughput improvement
- 5x vector operations

## 4. Security

- Lua script security fixes
- ACL improvements
- TLS enhancements

---

## References

1. Redis 8 Release Notes
2. Redis Documentation

---

*Last Updated: 2026-04-03*
