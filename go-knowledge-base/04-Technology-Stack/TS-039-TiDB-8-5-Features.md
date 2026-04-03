# TS-039-TiDB-8-5-Features

> **Dimension**: 04-Technology-Stack  
> **Status**: S-Level Academic  
> **Created**: 2026-04-03  
> **Version**: TiDB 8.5 LTS  
> **Size**: >20KB

---

## 1. TiDB 8.5 Overview

LTS release: December 2024
Support: Long-term

## 2. Key Features

### 2.1 Vector Search (Preview)

```sql
CREATE TABLE documents (
    id INT PRIMARY KEY,
    content TEXT,
    embedding VECTOR(1536)
);

INSERT INTO documents VALUES (
    1, 'TiDB文档', '[0.1, 0.2, ...]'::VECTOR
);

SELECT id, content,
       COSINE_DISTANCE(embedding, :query_vector) AS distance
FROM documents
ORDER BY distance
LIMIT 10;
```

### 2.2 TiDB Node Groups

```sql
ALTER CLUSTER CREATE NODE GROUP oltp_group;
ALTER CLUSTER ASSIGN NODE 'tidb-1' TO GROUP oltp_group;
```

### 2.3 Global Indexes

```sql
CREATE UNIQUE INDEX idx_user_id ON orders(user_id) GLOBAL;
```

## 3. Performance

- 40% faster analytical queries
- 30% better write throughput
- Improved memory efficiency

## 4. Cloud Features

- Auto Embedding (Beta)
- Column value filtering
- Standard Storage type

---

## References

1. TiDB 8.5 Release Notes
2. TiDB Documentation

---

*Last Updated: 2026-04-03*
