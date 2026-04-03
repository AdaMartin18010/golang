# TS-041-ClickHouse-25-Features

> **Dimension**: 04-Technology-Stack  
> **Status**: S-Level Academic  
> **Created**: 2026-04-03  
> **Version**: ClickHouse 25.x  
> **Size**: >20KB

---

## 1. ClickHouse 25 Overview

25.8 LTS released August 2025

## 2. Key Features

### 2.1 JSON Type (Production)

```sql
CREATE TABLE events (
    timestamp DateTime64(3),
    event_data JSON
);

INSERT INTO events VALUES 
(now(), '{"type": "click", "target": "button"}');

SELECT event_data.type FROM events;
```

### 2.2 Variant Type

```sql
CREATE TABLE metrics (
    name String,
    value Variant(UInt64, Float64, String)
);
```

### 2.3 Dynamic Tables

```sql
CREATE DYNAMIC TABLE hourly_stats AS
SELECT 
    toStartOfHour(timestamp) as hour,
    count() as cnt
FROM events
GROUP BY hour;
```

## 3. Performance

- Query condition cache
- Lightweight projections
- SIMD optimizations (AVX-512)

## 4. Integration

- OTLP native ingestion
- S3 integration improvements
- Kafka enhancements

---

## References

1. ClickHouse 25 Release Notes
2. ClickHouse Docs

---

*Last Updated: 2026-04-03*
