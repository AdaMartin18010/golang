# TS-004: Elasticsearch 9.0 内部机制 (Elasticsearch 9.0 Internals)

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #elasticsearch9 #lucene #inverted-index #sharding
> **权威来源**: [Elasticsearch Reference](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html), [Lucene Documentation](https://lucene.apache.org/core/documentation.html)

---

## 架构演进

```
Elasticsearch 7.x (2019)    ES 8.x (2022)              ES 9.0 (2026) ⭐️
      │                          │                          │
      ▼                          ▼                          ▼
┌─────────────┐          ┌───────────────┐          ┌─────────────────┐
│  Lucene 8   │          │  Lucene 9     │          │  Lucene 10      │
│  Type Removal│─────────►│  Security     │─────────►│  AI/ML Native   │
│             │          │  by Default   │          │  Vector Search  │
└─────────────┘          └───────────────┘          │  Semantic Search│
                                                    └─────────────────┘
```

---

## 倒排索引 (Inverted Index)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Inverted Index Structure                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Documents                              Inverted Index                       │
│  ─────────                              ──────────────                       │
│                                                                              │
│  Doc 1: "the quick brown fox"           Term      Doc IDs (Posting List)    │
│  Doc 2: "the lazy dog"                  ────      ─────────────────────     │
│  Doc 3: "quick dog jumps"               the       [1, 2]                    │
│                                         quick     [1, 3]                    │
│                                         brown     [1]                       │
│                                         fox       [1]                       │
│                                         lazy      [2]                       │
│                                         dog       [2, 3]                    │
│                                         jumps     [3]                       │
│                                                                              │
│  查询 "quick dog":                                                          │
│  quick → [1, 3]                                                             │
│  dog   → [2, 3]                                                             │
│  AND   → [3]  (交集)                                                        │
│                                                                              │
│  优化：Skip Lists, Bitsets, Roaring Bitmaps                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 分片与副本

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Elasticsearch Sharding                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Index: logs-2024                                                            │
│  ├── Settings: 3 primary shards, 1 replica                                  │
│  │                                                                           │
│  │  Node 1              Node 2              Node 3                           │
│  │  ┌──────────┐       ┌──────────┐       ┌──────────┐                     │
│  │  │ P0 (5GB) │       │ P1 (5GB) │       │ P2 (5GB) │  Primary Shards    │
│  │  │ R1 (5GB) │◄─────►│ R2 (5GB) │◄─────►│ R0 (5GB) │  Replicas          │
│  │  └──────────┘       └──────────┘       └──────────┘                     │
│  │       │                  │                  │                             │
│  │       └──────────────────┼──────────────────┘                             │
│  │                          │                                                │
│  │                    Routing: hash(_id) % num_shards                        │
│  │                    doc ABC → shard 1                                      │
│  │                                                                           │
│  │  路由公式: shard = hash(routing) % number_of_primary_shards               │
│  │  ⚠️ 创建后不能修改主分片数!                                                │
│  │                                                                           │
│  └──────────────────────────────────────────────────────────────────────────┘
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ES 9.0 新特性

### 1. 原生向量搜索

```json
// ES 9.0: 原生 dense_vector 支持
{
  "mappings": {
    "properties": {
      "text": {
        "type": "text"
      },
      "embedding": {
        "type": "dense_vector",
        "dims": 768,
        "index": true,
        "similarity": "cosine",
        "index_options": {
          "type": "hnsw",
          "m": 16,
          "ef_construction": 100
        }
      }
    }
  }
}

// 语义搜索
{
  "query": {
    "knn": {
      "field": "embedding",
      "query_vector": [0.1, 0.2, ...],
      "k": 10,
      "num_candidates": 100
    }
  }
}
```

### 2. 语义搜索 (Semantic Search)

```json
// ES 9.0: 内置 NLP 模型
{
  "settings": {
    "index": {
      "default_pipeline": "nlp-pipeline"
    }
  },
  "mappings": {
    "properties": {
      "content": {
        "type": "semantic_text",
        "model_id": "sentence-transformers/all-MiniLM-L6-v2"
      }
    }
  }
}

// 自动向量化 + 语义搜索
{
  "query": {
    "semantic": {
      "field": "content",
      "query": "machine learning applications"
    }
  }
}
```

---

## 写入流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Document Indexing Flow                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Client ──IndexRequest──► Coordinating Node                               │
│                                 │                                            │
│  2. Routing (hash(_id) % shards)                                             │
│                                 │                                            │
│                                 ▼                                            │
│  3. Primary Shard ──┬─► In-Memory Buffer (Lucene Segment)                    │
│                     │   │                                                   │
│                     │   ├─► Translog ( durability )                          │
│                     │   │                                                   │
│                     │   └─► Refresh (1s default) ──► FileSystem Cache        │
│                     │                                                       │
│                     └──► Replicate to Replica Shards                         │
│                              │                                              │
│                              ▼                                              │
│  4. Flush (30min or translog size) ──► Disk (Immutable Segment)             │
│                                                                              │
│  5. Merge (background) ──► Merge small segments to large                     │
│                                                                              │
│  近实时 (NRT): 1s refresh 间隔                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 与关系数据库对比

| 特性 | Elasticsearch | PostgreSQL |
|------|---------------|------------|
| 存储 | 倒排索引 | B-Tree |
| 查询 | 全文搜索、向量 | 精确匹配、JOIN |
| Schema | 灵活 (mapping) | 严格 |
| 事务 | 最终一致 | ACID |
| 适用 | 搜索、日志、分析 | 事务、关系 |

---

## 参考文献

1. [Elasticsearch: The Definitive Guide](https://www.elastic.co/guide/en/elasticsearch/guide/current/index.html)
2. [Lucene in Action](https://www.manning.com/books/lucene-in-action-second-edition)
3. [Elasticsearch 9.0 Release Notes](https://www.elastic.co/guide/en/elasticsearch/reference/current/release-notes.html)
