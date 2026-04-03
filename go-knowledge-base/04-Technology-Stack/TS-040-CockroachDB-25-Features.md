# TS-040-CockroachDB-25-Features

> **Dimension**: 04-Technology-Stack  
> **Status**: S-Level Academic  
> **Created**: 2026-04-03  
> **Version**: CockroachDB 25.2  
> **Size**: >20KB

---

## 1. CockroachDB 25 Overview

25.2 release: May 2025

## 2. Key Features

### 2.1 C-SPANN Vector Index

```sql
CREATE TABLE documents (
    id UUID DEFAULT gen_random_uuid(),
    embedding VECTOR(1536),
    PRIMARY KEY (id)
);

CREATE INDEX idx_embedding ON documents USING cspann (embedding);

SELECT id, l2_distance(embedding, $1) 
FROM documents
ORDER BY embedding <-> $1
LIMIT 10;
```

### 2.2 Buffered Writes

30% faster baseline throughput

### 2.3 Leader Leases

Unified read/write authority

## 3. Performance

| Metric | 25.1 | 25.2 | Improvement |
|--------|------|------|-------------|
| SQL latency | ~3ms | ~1.32ms | 56% |
| tpmC | - | 88.1K | +41% |
| Bulk import | - | 4x | faster |

## 4. Security

- Row-level security
- Configurable TLS ciphers
- Physical cluster replication GA

---

## References

1. CockroachDB 25.2 Release Notes
2. CockroachDB Docs

---

*Last Updated: 2026-04-03*
