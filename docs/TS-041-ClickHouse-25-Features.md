# TS-041: ClickHouse 25 Columnar Storage - S-Level Technical Reference

**Version:** ClickHouse 25.1  
**Status:** S-Level (Expert/Architectural)  
**Last Updated:** 2026-04-03  
**Classification:** OLAP / Columnar Storage / Vectorized Query Execution

---

## 1. Executive Summary

ClickHouse 25 introduces significant enhancements to its columnar storage engine, including the new MergeTree engine optimizations, advanced compression algorithms, and improved vectorized query execution. This document provides comprehensive technical analysis of ClickHouse's storage architecture, data organization, and query processing mechanisms.

---

## 2. Storage Architecture

### 2.1 Table Engine Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ClickHouse MergeTree Engine Architecture                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Table Structure                                     │  │
│  │                                                                        │  │
│  │  CREATE TABLE events (                                                │  │
│  │      timestamp DateTime64(3),                                         │  │
│  │      user_id UInt64,                                                  │  │
│  │      event_type LowCardinality(String),                               │  │
│  │      properties String CODEC(ZSTD(3)),                                │  │
│  │      value Float64 CODEC(Gorilla, LZ4)                                │  │
│  │  ) ENGINE = MergeTree()                                               │  │
│  │  PARTITION BY toYYYYMMDD(timestamp)                                   │  │
│  │  ORDER BY (user_id, timestamp)                                        │  │
│  │  PRIMARY KEY (user_id)                                                │  │
│  │  TTL timestamp + INTERVAL 90 DAY                                      │  │
│  │  SETTINGS index_granularity = 8192;                                   │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Physical Storage Layout                             │  │
│  │                                                                        │  │
│  │  /var/lib/clickhouse/data/database/table/                             │  │
│  │  │                                                                      │  │
│  │  ├── 20240101_1_1_0/          ← Partition directory                    │  │
│  │  │   │                                                                   │  │
│  │  │   ├── checksums.txt         ← SHA256 checksums for all files         │  │
│  │  │   ├── columns.txt           ← Column names and types                 │  │
│  │  │   ├── count.txt             ← Row count                              │  │
│  │  │   ├── primary.idx          ← Primary index (sparse)                  │  │
│  │  │   ├── minmax_timestamp.idx ← Partition column index                 │  │
│  │  │   │                                                                   │  │
│  │  │   ├── timestamp.bin         ← Column data (compressed)               │  │
│  │  │   ├── timestamp.mrk2        ← Mark files (offsets to granules)       │  │
│  │  │   │                                                                   │  │
│  │  │   ├── user_id.bin            ← Column data                           │  │
│  │  │   ├── user_id.mrk2           ← Mark files                            │  │
│  │  │   │                                                                   │  │
│  │  │   ├── event_type.bin         ← Dictionary-encoded column             │  │
│  │  │   ├── event_type.mrk2        ← Mark files                            │  │
│  │  │   │                                                                   │  │
│  │  │   ├── properties.bin         ← ZSTD compressed data                  │  │
│  │  │   ├── properties.mrk2        ← Mark files                            │  │
│  │  │   │                                                                   │  │
│  │  │   ├── value.bin              ← Gorilla + LZ4 compressed              │  │
│  │  │   ├── value.mrk2             ← Mark files                            │  │
│  │  │   │                                                                   │  │
│  │  │   ├── skp_idx_idx_name.idx   ← Skipping index (optional)             │  │
│  │  │   └── part_columns.txt       ← Columns in this part                  │  │
│  │  │                                                                       │  │
│  │  ├── 20240101_2_2_0/          ← Another part in same partition          │  │
│  │  ├── 20240102_3_3_0/          ← Different partition                     │  │
│  │  └── detached/                 ← Detached parts for recovery            │  │
│  │                                                                          │  │
│  │  Part naming: {partition}_{min_block}_{max_block}_{level}               │  │
│  │  • partition: Value of partition key                                    │  │
│  │  • min_block: Minimum block number in this part                         │  │
│  │  • max_block: Maximum block number in this part                         │  │
│  │  • level: Merge depth (0 = inserted, >0 = merged)                       │  │
│  │                                                                          │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 MergeTree Data Organization

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ClickHouse MergeTree Data Organization                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Granule Structure (index_granularity = 8192):                               │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  Table Data:                                                           │   │
│  │  ┌─────────┬─────────┬─────────┬─────────┬─────────┬─────────┐       │   │
│  │  │ Granule │ Granule │ Granule │ Granule │ Granule │ Granule │       │   │
│  │  │   0     │   1     │   2     │   3     │   4     │   5     │       │   │
│  │  │8192 rows│8192 rows│8192 rows│8192 rows│8192 rows│8192 rows│       │   │
│  │  │[0-8191] │[8192-  │[16384- │[24576- │[32768- │[40960- │       │   │
│  │  │         │ 16383] │ 24575] │ 32767] │ 40959] │ 49151] │       │   │
│  │  └────┬────┴────┬────┴────┬────┴────┬────┴────┬────┴────┬────┘       │   │
│  │       │         │         │         │         │         │            │   │
│  │       ▼         ▼         ▼         ▼         ▼         ▼            │   │
│  │  ┌────────────────────────────────────────────────────────────────┐  │   │
│  │  │                    Primary Index (sparse)                       │  │   │
│  │  │                                                                  │  │   │
│  │  │  ┌───────────────────────────────────────────────────────────┐ │  │   │
│  │  │  │ Mark │ Granule Start │ Primary Key Value (first in granule)│ │  │   │
│  │  │  ├──────┼───────────────┼─────────────────────────────────────┤ │  │   │
│  │  │  │  0   │      0        │ user_id=100, timestamp=2024-01-... │ │  │   │
│  │  │  │  1   │     8192      │ user_id=150, timestamp=2024-01-... │ │  │   │
│  │  │  │  2   │    16384      │ user_id=200, timestamp=2024-01-... │ │  │   │
│  │  │  │  3   │    24576      │ user_id=250, timestamp=2024-01-... │ │  │   │
│  │  │  │  4   │    32768      │ user_id=300, timestamp=2024-01-... │ │  │   │
│  │  │  │  5   │    40960      │ user_id=350, timestamp=2024-01-... │ │  │   │
│  │  │  └──────┴───────────────┴─────────────────────────────────────┘ │  │   │
│  │  │                                                                  │  │   │
│  │  │  Index size: 1 mark per 8192 rows (configurable)                │  │   │
│  │  │  For 1 billion rows: ~122K marks, ~4MB index size               │  │   │
│  │  │                                                                  │  │   │
│  │  └────────────────────────────────────────────────────────────────┘  │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Column File Format:                                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  timestamp.bin (compressed):                                           │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │   │
│  │  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │   │   │
│  │  │  │ Compressed  │ │ Compressed  │ │ Compressed  │               │   │   │
│  │  │  │  Block 0    │ │  Block 1    │ │  Block 2    │ ...           │   │   │
│  │  │  │ (64KB-1MB) │ │ (64KB-1MB) │ │ (64KB-1MB) │               │   │   │
│  │  │  │             │ │             │ │             │               │   │   │
│  │  │  │ 8192 rows  │ │ 8192 rows  │ │ 8192 rows  │               │   │   │
│  │  │  │ Delta,     │ │ Delta,     │ │ Delta,     │               │   │   │
│  │  │  │ LZ4        │ │ LZ4        │ │ LZ4        │               │   │   │
│  │  │  └─────────────┘ └─────────────┘ └─────────────┘               │   │   │
│  │  └─────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                        │   │
│  │  timestamp.mrk2 (mark file):                                           │   │
│  │  ┌─────────────────────────────────────────────────────────────────┐   │   │
│  │  │  Offset in    │ Offset in    │ Rows in      │                   │   │   │
│  │  │  .bin file   │ decompressed │ block        │                   │   │   │
│  │  ├───────────────┼──────────────┼──────────────┤                   │   │   │
│  │  │  0           │  0           │  8192        │  ← Granule 0      │   │   │
│  │  │  45000       │  0           │  8192        │  ← Granule 1      │   │   │
│  │  │  89000       │  0           │  8192        │  ← Granule 2      │   │   │
│  │  │  ...         │  ...         │  ...         │                   │   │   │
│  │  └─────────────────────────────────────────────────────────────────┘   │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Compression Algorithms

### 3.1 Codec Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ClickHouse Compression Codecs                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  General-Purpose Compression:                                                │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  Codec          │ Ratio │ Speed (GB/s) │ Use Case                     │   │
│  ├─────────────────┼───────┼──────────────┼──────────────────────────────┤   │
│  │  LZ4            │ 2-3x  │  0.8 / 3.5   │ General, low CPU             │   │
│  │  LZ4HC          │ 3-4x  │  0.1 / 3.5   │ Better ratio, same decomp    │   │
│  │  ZSTD           │ 4-6x  │  0.3 / 1.0   │ Best ratio, moderate CPU     │   │
│  │  ZSTD(level)    │ 4-8x  │ variable     │ Configurable compression     │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Specialized Columnar Codecs:                                                │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │  Codec          │ Target Type    │ Ratio │ Description                │   │
│  ├─────────────────┼────────────────┼───────┼────────────────────────────┤   │
│  │  None           │ Any            │ 1x    │ No compression             │   │
│  │  Delta          │ Integer, Date  │ 2-5x  │ Delta encoding             │   │
│  │  DeltaDelta     │ Timestamp      │ 3-8x  │ Second-order delta         │   │
│  │  DoubleDelta    │ Timestamp      │ 4-10x │ Third-order delta          │   │
│  │  Gorilla        │ Float          │ 2-4x  │ XOR-based float compression│   │
│  │  FPC            │ Float          │ 2-4x  │ Floating-point compression │   │
│  │  T64            │ Integer        │ 2-4x  │ Trim to 64/32/16/8 bits    │   │
│  │  LowCardinality │ String, Enum   │ 5-50x │ Dictionary encoding        │   │
│  │  SET            │ Low cardinality│ 3-10x │ Bitmap encoding            │   │
│  │  Encrypted      │ Any            │ 1x    │ AES-256-GCM encryption     │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Codec Chaining:                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  CODEC(Delta, ZSTD)        ← Timestamp column                        │   │
│  │  ┌─────────┐    ┌─────────┐    ┌─────────┐                            │   │
│  │  │ Raw     │───▶│ Delta   │───▶│ ZSTD    │──▶ Disk                    │   │
│  │  │ Values  │    │ Encoded │    │ Compressed                              │   │
│  │  └─────────┘    └─────────┘    └─────────┘                            │   │
│  │                                                                        │   │
│  │  CODEC(Gorilla, LZ4)       ← Float column                            │   │
│  │  ┌─────────┐    ┌─────────┐    ┌─────────┐                            │   │
│  │  │ Float   │───▶│ Gorilla │───▶│ LZ4     │──▶ Disk                    │   │
│  │  │ Values  │    │ XOR     │    │ Compressed                              │   │
│  │  └─────────┘    └─────────┘    └─────────┘                            │   │
│  │                                                                        │   │
│  │  CODEC(LowCardinality)     ← String enum                             │   │
│  │  ┌─────────┐    ┌─────────┐    ┌─────────┐                            │   │
│  │  │ String  │───▶│ Dict    │───▶│ LZ4     │──▶ Disk                    │   │
│  │  │ Values  │    │ Indices │    │ Compressed                              │   │
│  │  └─────────┘    └─────────┘    └─────────┘                            │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Gorilla Float Compression

```
ALGORITHM GorillaCompress(values):
    INPUT:  values - Array of float64 values
    OUTPUT: compressed bitstream
    
    1. // First value: store verbatim
       prev_bits ← float_to_bits(values[0])
       write_bits(prev_bits, 64)
    
    2. prev_xor ← 0
       prev_leading_zeros ← 0
       prev_trailing_zeros ← 0
    
    3. FOR each value in values[1:]:
       a. curr_bits ← float_to_bits(value)
       b. xor ← prev_bits XOR curr_bits
       
       c. IF xor == 0:
              // Same value, write single 0 bit
              write_bit(0)
          ELSE:
              write_bit(1)
              
              leading_zeros ← count_leading_zeros(xor)
              trailing_zeros ← count_trailing_zeros(xor)
              meaningful_bits ← 64 - leading_zeros - trailing_zeros
              
              IF leading_zeros >= prev_leading_zeros AND 
                 trailing_zeros >= prev_trailing_zeros:
                  // Use previous block
                  write_bit(0)
                  write_bits(xor >> trailing_zeros, 
                            64 - prev_leading_zeros - prev_trailing_zeros)
              ELSE:
                  // New block description
                  write_bit(1)
                  write_bits(leading_zeros, 6)   // 0-63 leading zeros
                  write_bits(meaningful_bits, 6) // 1-64 meaningful bits
                  write_bits(xor >> trailing_zeros, meaningful_bits)
                  
                  prev_leading_zeros ← leading_zeros
                  prev_trailing_zeros ← trailing_zeros
       
       d. prev_bits ← curr_bits
    
    4. RETURN bitstream

ALGORITHM GorillaDecompress(bitstream, count):
    1. // First value
       prev_bits ← read_bits(64)
       result ← [bits_to_float(prev_bits)]
    
    2. prev_leading_zeros ← 0
       prev_trailing_zeros ← 0
    
    3. FOR i ← 1 TO count - 1:
       IF read_bit() == 0:
          // Same value
          result.append(result[-1])
       ELSE:
          IF read_bit() == 0:
              // Same block
              meaningful_bits ← 64 - prev_leading_zeros - prev_trailing_zeros
              xor ← read_bits(meaningful_bits) << prev_trailing_zeros
          ELSE:
              // New block
              leading_zeros ← read_bits(6)
              meaningful_bits ← read_bits(6)
              xor ← read_bits(meaningful_bits) << 
                    (64 - leading_zeros - meaningful_bits)
              prev_leading_zeros ← leading_zeros
              prev_trailing_zeros ← 64 - leading_zeros - meaningful_bits
          
          curr_bits ← prev_bits XOR xor
          result.append(bits_to_float(curr_bits))
          prev_bits ← curr_bits
    
    4. RETURN result
```

---

## 4. Vectorized Query Execution

### 4.1 Query Execution Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ClickHouse Vectorized Query Execution                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Query:                                                                      │
│  SELECT user_id, avg(value)                                                 │
│  FROM events                                                                │
│  WHERE timestamp >= '2024-01-01' AND event_type = 'purchase'                │
│  GROUP BY user_id                                                           │
│  ORDER BY avg(value) DESC                                                   │
│  LIMIT 100;                                                                 │
│                                                                              │
│  Execution Pipeline:                                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐   │
│  │                                                                        │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │   Read      │  Read granules from disk                             │   │
│  │  │   From Disk │  Uses primary index for filtering                    │   │
│  │  │             │  Skips irrelevant granules                           │   │
│  │  └──────┬──────┘                                                       │   │
│  │         │ Decompress blocks (SIMD-accelerated)                         │   │
│  │         ▼                                                              │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │  Decompress │  LZ4/ZSTD decompression                               │   │
│  │  │             │  Uses hardware acceleration when available            │   │
│  │  └──────┬──────┘                                                       │   │
│  │         │ Decode special codecs (Delta, Gorilla, etc.)                 │   │
│  │         ▼                                                              │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │    Filter   │  WHERE timestamp >= '2024-01-01'                       │   │
│  │  │             │  Vectorized comparison using SIMD                     │   │
│  │  │  SIMD SSE/  │  Produces selection vector (bitmap)                   │   │
│  │  │  AVX2/AVX512│  No branching - data-parallel execution               │   │
│  │  └──────┬──────┘                                                       │   │
│  │         │                                                              │   │
│  │         ▼                                                              │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │    Filter   │  WHERE event_type = 'purchase'                        │   │
│  │  │             │  Dictionary lookup for LowCardinality                 │   │
│  │  └──────┬──────┘                                                       │   │
│  │         │ Combine filters (AND)                                        │   │
│  │         ▼                                                              │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │   Project   │  Select user_id, value columns                        │   │
│  │  │             │  Gather from filtered rows                            │   │
│  │  └──────┬──────┘                                                       │   │
│  │         │                                                              │   │
│  │         ▼                                                              │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │  Aggregate  │  GROUP BY user_id                                     │   │
│  │  │             │  Hash table: user_id -> {count, sum}                  │   │
│  │  │  HashAggregate                                                       │   │
│  │  │             │  Parallel aggregation with merge                        │   │
│  │  └──────┬──────┘                                                       │   │
│  │         │ Finalize: sum / count = avg                                  │   │
│  │         ▼                                                              │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │    Sort     │  ORDER BY avg(value) DESC                             │   │
│  │  │             │  Partial sort if large result set                     │   │
│  │  │  Top-N Sort │  Priority queue for LIMIT optimization                │   │
│  │  └──────┬──────┘                                                       │   │
│  │         │                                                              │   │
│  │         ▼                                                              │   │
│  │  ┌─────────────┐                                                       │   │
│  │  │    Limit    │  LIMIT 100                                            │   │
│  │  │             │                                                       │   │
│  │  └─────────────┘                                                       │   │
│  │                                                                        │   │
│  └───────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 SIMD Vectorized Operations

```cpp
// ClickHouse SIMD-accelerated column operations
// ClickHouse/src/Columns/ColumnsCommon.cpp

#include <immintrin.h>

// Vectorized filter using AVX2
size_t filterColumnAVX2(const UInt8* data, size_t size, 
                        IColumn::Filter& filter) {
    const __m256i zero = _mm256_setzero_si256();
    size_t result_size = 0;
    
    // Process 32 bytes at a time
    size_t i = 0;
    for (; i + 32 <= size; i += 32) {
        __m256i bytes = _mm256_loadu_si256(
            reinterpret_cast<const __m256i*>(data + i));
        
        // Compare with zero (non-zero = true)
        __m256i cmp = _mm256_cmpeq_epi8(bytes, zero);
        int mask = _mm256_movemask_epi8(cmp);
        
        // Invert mask: 0 -> 1 (true), 0xFF -> 0 (false)
        mask = ~mask;
        
        // Store 32-bit mask
        *reinterpret_cast<UInt32*>(&filter[result_size]) = mask;
        result_size += _mm_popcnt_u32(mask);
    }
    
    // Handle remaining elements
    for (; i < size; ++i) {
        filter[result_size] = data[i] != 0;
        result_size += data[i] != 0;
    }
    
    return result_size;
}

// Vectorized comparison for WHERE clauses
template<typename T>
void compareColumnConstAVX2(const T* data, size_t size, 
                            T value, UInt8* result,
                            ComparisonFunc comp) {
    const size_t vec_size = 32 / sizeof(T);  // 4 for UInt64, 8 for UInt32
    
    // Broadcast value to all lanes
    __m256i vec_value;
    if constexpr (sizeof(T) == 8) {
        vec_value = _mm256_set1_epi64x(value);
    } else if constexpr (sizeof(T) == 4) {
        vec_value = _mm256_set1_epi32(value);
    }
    
    size_t i = 0;
    for (; i + vec_size <= size; i += vec_size) {
        __m256i vec_data = _mm256_loadu_si256(
            reinterpret_cast<const __m256i*>(data + i));
        
        __m256i cmp_result;
        if constexpr (sizeof(T) == 8) {
            cmp_result = _mm256_cmpgt_epi64(vec_data, vec_value);
        } else {
            cmp_result = _mm256_cmpgt_epi32(vec_data, vec_value);
        }
        
        // Store comparison results
        _mm256_storeu_si256(reinterpret_cast<__m256i*>(result + i), 
                           cmp_result);
    }
    
    // Scalar fallback for remaining
    for (; i < size; ++i) {
        result[i] = comp(data[i], value) ? 0xFF : 0;
    }
}
```

---

## 5. Merge Process

### 5.1 Background Merge Algorithm

```
ALGORITHM MergeParts(parts_to_merge):
    INPUT:  parts_to_merge - List of parts to merge
    OUTPUT: merged_part
    
    1. // Create output part directory
       merged_part ← create_new_part()
       merged_part.min_block ← min(part.min_block for part in parts_to_merge)
       merged_part.max_block ← max(part.max_block for part in parts_to_merge)
       merged_part.level ← max(part.level for part in parts_to_merge) + 1
    
    2. FOR each column in table schema:
       a. // Open readers for all parts
          readers ← []
          FOR part in parts_to_merge:
              reader ← ColumnReader(part.column_file)
              readers.append(reader)
       
       b. // Create merger (k-way merge for sorted data)
          merger ← MergeSorter(readers, sort_key)
       
       c. // Write merged column
          writer ← ColumnWriter(merged_part.column_file)
          
          WHILE merger.has_next():
              granule ← merger.next_granule()
              writer.write_granule(granule)
       
       d. // Write mark file
          writer.write_marks(merged_part.mark_file)
    
    3. // Merge primary index
       merged_index ← merge_indexes([p.primary_idx for p in parts_to_merge])
       write_index(merged_part.primary_idx, merged_index)
    
    4. // Calculate and write checksums
       checksums ← calculate_checksums(merged_part)
       write_checksums(merged_part.checksums_file, checksums)
    
    5. // Atomically replace old parts with merged part
       transaction ← begin_transaction()
       FOR part in parts_to_merge:
           transaction.remove_part(part)
       transaction.add_part(merged_part)
       transaction.commit()
    
    6. RETURN merged_part

// Merge selection algorithm
ALGORITHM SelectPartsToMerge(parts, max_parts_to_merge):
    // Score-based merge selection
    candidates ← []
    
    FOR partition in group_parts_by_partition(parts):
        partition_parts ← sorted(partition, key=lambda p: p.level)
        
        // Score based on:
        // 1. Size ratio (prefer similar sized parts)
        // 2. Age (prefer older parts)
        // 3. Number of parts at same level
        
        FOR level, level_parts in group_by_level(partition_parts):
            IF length(level_parts) >= 2:
                score ← calculate_merge_score(level_parts)
                candidates.append((score, level_parts))
    
    // Select best candidates
    candidates.sort_by_score(descending)
    RETURN candidates[0:max_parts_to_merge]
```

---

## 6. Performance Benchmarks

### 6.1 Compression Ratios

| Data Type | Raw Size | LZ4 | ZSTD(1) | ZSTD(9) | Best Codec | Compression |
|-----------|----------|-----|---------|---------|------------|-------------|
| Timestamp | 8 bytes | 3.2x | 5.1x | 7.8x | Delta+ZSTD | 12.5x |
| UInt64 ID | 8 bytes | 2.1x | 3.5x | 4.2x | T64+ZSTD | 6.8x |
| Float64 | 8 bytes | 2.8x | 4.5x | 5.2x | Gorilla+LZ4 | 8.2x |
| String (URL) | 65 bytes avg | 4.5x | 8.2x | 11.5x | LowCardinality | 25x |
| JSON | 500 bytes avg | 6.2x | 12.5x | 18.2x | ZSTD(9) | 18.2x |

### 6.2 Query Performance

| Query Type | Rows/sec | Data Scanned/sec | Latency (10M rows) |
|------------|----------|------------------|-------------------|
| Full Scan | 2.5B | 20GB | 0.5s |
| Filtered Scan (1%) | 1.2B | 10GB | 0.8s |
| Aggregation | 800M | 6.4GB | 1.2s |
| Group By (cardinality 1M) | 200M | 1.6GB | 5s |
| Join (hash) | 50M | 400MB | 20s |

---

## 7. References

1. **ClickHouse Documentation**
   - URL: https://clickhouse.com/docs

2. **MergeTree Engine Guide**
   - URL: https://clickhouse.com/docs/en/engines/table-engines/mergetree-family/mergetree

3. **Data Compression in ClickHouse**
   - URL: https://clickhouse.com/docs/en/sql-reference/statements/create/table#column-compression-codec

4. **Gorilla: A Fast, Scalable, In-Memory Time Series Database**
   - URL: https://www.vldb.org/pvldb/vol8/p1816-teller.pdf

5. **ClickHouse Source Code**
   - URL: https://github.com/ClickHouse/ClickHouse

---

*Document generated for S-Level technical reference.*
