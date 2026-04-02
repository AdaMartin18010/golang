# TS-004: Elasticsearch Query DSL - Internals & Go Implementation

> **维度**: Technology Stack
> **级别**: S (18+ KB)
> **标签**: #elasticsearch #search #lucene #query-dsl #go
> **权威来源**:
>
> - [Elasticsearch Documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html) - Elastic
> - [Lucene in Action](https://lucene.apache.org/core/) - Apache Lucene
> - [Elasticsearch: The Definitive Guide](https://www.elastic.co/guide/en/elasticsearch/guide/current/index.html)

---

## 1. Elasticsearch Internal Architecture

### 1.1 Cluster & Node Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Elasticsearch Cluster Architecture                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        Elasticsearch Cluster                           │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │  │
│  │  │                      Master-Eligible Nodes                       │ │  │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                         │ │  │
│  │  │  │ Master  │  │ Master  │  │ Master  │  voting_config_only      │ │  │
│  │  │  │ (Active)│  │ (Standby)│  │ (Standby)│                        │ │  │
│  │  │  │ node.master: true    │  │ node.data: false                   │ │  │
│  │  │  └─────────┘  └─────────┘  └─────────┘                         │ │  │
│  │  │                                                                  │ │  │
│  │  │  Responsibilities:                                             │ │  │
│  │  │  - Cluster state management                                    │ │  │
│  │  │  - Index/shard allocation                                      │ │  │
│  │  │  - Node membership                                             │ │  │
│  │  └─────────────────────────────────────────────────────────────────┘ │  │
│  │                                                                      │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │  │
│  │  │                        Data Nodes                                │ │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │ │  │
│  │  │  │ Data Node 1 │  │ Data Node 2 │  │ Data Node 3 │             │ │  │
│  │  │  │             │  │             │  │             │             │ │  │
│  │  │  │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌─────────┐ │             │ │  │
│  │  │  │ │Shard P0 │ │  │ │Shard P1 │ │  │ │Shard P0R│ │             │ │  │
│  │  │  │ │(Primary)│ │  │ │(Primary)│ │  │ │(Replica)│ │             │ │  │
│  │  │  │ ├─────────┤ │  │ ├─────────┤ │  │ ├─────────┤ │             │ │  │
│  │  │  │ │Shard P1R│ │  │ │Shard P0R│ │  │ │Shard P1R│ │             │ │  │
│  │  │  │ │(Replica)│ │  │ │(Replica)│ │  │ │(Replica)│ │             │ │  │
│  │  │  │ └─────────┘ │  │ └─────────┘ │  │ └─────────┘ │             │ │  │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘             │ │  │
│  │  └─────────────────────────────────────────────────────────────────┘ │  │
│  │                                                                      │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │  │
│  │  │                   Coordinating/Client Nodes                      │ │  │
│  │  │  ┌─────────────┐  ┌─────────────┐                              │ │  │
│  │  │  │ Client 1    │  │ Client 2    │                              │ │  │
│  │  │  │ (No data)   │  │ (No data)   │                              │ │  │
│  │  │  │ Query agg   │  │ Query agg   │                              │ │  │
│  │  │  │ Routing     │  │ Routing     │                              │ │  │
│  │  │  └─────────────┘  └─────────────┘                              │ │  │
│  │  └─────────────────────────────────────────────────────────────────┘ │  │
│  │                                                                      │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐ │  │
│  │  │                    Ingest Nodes (Optional)                       │ │  │
│  │  │  - Document preprocessing (pipelines)                            │ │  │
│  │  │  - ETL operations before indexing                                │ │  │
│  │  └─────────────────────────────────────────────────────────────────┘ │  │
│  │                                                                      │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Index & Shard Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Index & Shard Internal Structure                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Index: "products"                                                           │
│  Settings: 3 primary shards, 1 replica per shard                             │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Shard Distribution                                  │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Node 1              Node 2              Node 3                        │  │
│  │  ┌─────────┐        ┌─────────┐        ┌─────────┐                    │  │
│  │  │ P0      │        │ P1      │        │ P2      │                    │  │
│  │  │ Primary │        │ Primary │        │ Primary │                    │  │
│  │  ├─────────┤        ├─────────┤        ├─────────┤                    │  │
│  │  │ R1      │        │ R2      │        │ R0      │                    │  │
│  │  │ Replica │        │ Replica │        │ Replica │                    │  │
│  │  └─────────┘        └─────────┘        └─────────┘                    │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Shard Internal Structure (Lucene Index)             │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Shard = Lucene Index                                                  │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                    Segment Structure                             │  │  │
│  │  │                                                                  │  │  │
│  │  │  Segment 0 (Immutable)      Segment 1 (Immutable)               │  │  │
│  │  │  ┌───────────────────┐      ┌───────────────────┐               │  │  │
│  │  │  │ .tim (Terms Dict) │      │ .tim (Terms Dict) │               │  │  │
│  │  │  │ .tip (Terms Index)│      │ .tip (Terms Index)│               │  │  │
│  │  │  │ .doc (Postings)   │      │ .doc (Postings)   │               │  │  │
│  │  │  │ .pos (Positions)  │      │ .pos (Positions)  │               │  │  │
│  │  │  │ .pay (Payloads)   │      │ .pay (Payloads)   │               │  │  │
│  │  │  │ .dvd (DocValues)  │      │ .dvd (DocValues)  │               │  │  │
│  │  │  │ .dvm (DV Metadata)│      │ .dvm (DV Metadata)│               │  │  │
│  │  │  │ .fdx (Field Index)│      │ .fdx (Field Index)│               │  │  │
│  │  │  │ .fdt (Field Data) │      │ .fdt (Field Data) │               │  │  │
│  │  │  │ .nvd (Norms)      │      │ .nvd (Norms)      │               │  │  │
│  │  │  │ .nvm (Norms Meta) │      │ .nvm (Norms Meta) │               │  │  │
│  │  │  │ .si (Segment Info)│      │ .si (Segment Info)│               │  │  │
│  │  │  └───────────────────┘      └───────────────────┘               │  │  │
│  │  │                                                                  │  │  │
│  │  │  Translog (Write-Ahead Log)                                      │  │  │
│  │  │  ┌─────────────────────────────────────────────────────────────┐ │  │  │
│  │  │  │ Op 1: INDEX doc_001                                         │ │  │  │
│  │  │  │ Op 2: INDEX doc_002                                         │ │  │  │
│  │  │  │ Op 3: DELETE doc_001                                        │ │  │  │
│  │  │  │ Op 4: INDEX doc_003                                         │ │  │  │
│  │  │  │ ...                                                         │ │  │  │
│  │  │  └─────────────────────────────────────────────────────────────┘ │  │  │
│  │  │                                                                  │  │  │
│  │  │  ┌─────────────────────────────────────────────────────────────┐ │  │  │
│  │  │  │                    Merge Process                             │ │  │  │
│  │  │  │                                                              │ │  │  │
│  │  │  │   Segment 0 + Segment 1 ──────► Merged Segment 2            │ │  │  │
│  │  │  │   (删除标记的文档被物理删除)                                  │ │  │  │
│  │  │  └─────────────────────────────────────────────────────────────┘ │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Inverted Index Structure                            │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Documents:                                                            │  │
│  │  Doc 1: "The quick brown fox"                                          │  │
│  │  Doc 2: "The quick blue fox"                                           │  │
│  │  Doc 3: "The slow brown dog"                                           │  │
│  │                                                                        │  │
│  │  Inverted Index (Term Dictionary → Postings List):                     │  │
│  │  ┌─────────┬─────────────────────────────────────────────────────┐    │  │
│  │  │ Term    │ Postings (Doc ID, Freq, Position, Offset)         │    │  │
│  │  ├─────────┼─────────────────────────────────────────────────────┤    │  │
│  │  │ the     │ (1,1,[0],[0]) (2,1,[0],[0]) (3,1,[0],[0])        │    │  │
│  │  │ quick   │ (1,1,[1],[4]) (2,1,[1],[4])                      │    │  │
│  │  │ brown   │ (1,1,[2],[10]) (3,1,[2],[9])                     │    │  │
│  │  │ fox     │ (1,1,[3],[16]) (2,1,[3],[15])                    │    │  │
│  │  │ blue    │ (2,1,[2],[10])                                   │    │  │
│  │  │ slow    │ (3,1,[1],[4])                                    │    │  │
│  │  │ dog     │ (3,1,[3],[15])                                   │    │  │
│  │  └─────────┴─────────────────────────────────────────────────────┘    │  │
│  │                                                                        │  │
│  │  FST (Finite State Transducer) for Term Dictionary:                    │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │                                                                  │  │  │
│  │  │      (b)──►(l)──►(u)──►(e)──►[postings ptr]                    │  │  │
│  │  │     /                                                           │  │  │
│  │  │  (r)──►(o)──►(w)──►(n)──►[postings ptr]                        │  │  │
│  │  │  /                                                              │  │  │
│  │  │ (b)                                                            │  │  │
│  │  │  \                                                             │  │  │
│  │  │   ...                                                          │  │  │
│  │  │                                                                  │  │  │
│  │  │  Memory-efficient term lookup (shared prefixes)                │  │  │
│  │  │                                                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.3 Query Execution Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Query Execution Flow                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Query Phase                                         │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Step 1: Client sends query                                            │  │
│  │  ┌─────────┐      ┌─────────────┐                                      │  │
│  │  │ Client  │─────►│ Coordinator │                                      │  │
│  │  └─────────┘      │ Node        │                                      │  │
│  │                   └──────┬──────┘                                      │  │
│  │                          │                                            │  │
│  │  Step 2: Query Parsing & Rewrite                                       │  │
│  │  ┌────────────────────────────────────────────────────────────┐       │  │
│  │  │ - Parse JSON Query DSL                                     │       │  │
│  │  │ - Rewrite multi-term queries (e.g., expand wildcards)      │       │  │
│  │  │ - Apply alias resolution                                   │       │  │
│  │  │ - Validate query against mapping                           │       │  │
│  │  └────────────────────────────────────────────────────────────┘       │  │
│  │                          │                                            │  │
│  │  Step 3: Identify Target Shards                                        │  │
│  │  ┌────────────────────────────────────────────────────────────┐       │  │
│  │  │ _routing provided?                                         │       │  │
│  │  │   ├── YES ──► Single shard lookup                          │       │  │
│  │  │   └── NO  ──► Broadcast to all shards                      │       │  │
│  │  │ Can shard filtering (date range) reduce set?               │       │  │
│  │  │   ├── YES ──► Filter shards                                │       │  │
│  │  │   └── NO  ──► All shards                                   │       │  │
│  │  └────────────────────────────────────────────────────────────┘       │  │
│  │                          │                                            │  │
│  │  Step 4: Scatter to Shards                                             │  │
│  │                   │                                                    │  │
│  │          ┌────────┼────────┬────────┐                                 │  │
│  │          ▼        ▼        ▼        ▼                                 │  │
│  │       ┌────┐  ┌────┐  ┌────┐  ┌────┐                                 │  │
│  │       │ P0 │  │ P1 │  │ P2 │  │ P0R│  (Target shards)               │  │
│  │       └────┘  └────┘  └────┘  └────┘                                 │  │
│  │          │       │       │       │                                    │  │
│  │          │       │       │       │                                    │  │
│  │  Step 5: Lucene Query Execution (per shard)                            │  │
│  │  ┌────────────────────────────────────────────────────────────┐       │  │
│  │  │ - Build Lucene Query from DSL                              │       │  │
│  │  │ - Search index segments (skipping deleted docs)            │       │  │
│  │  │ - Collect top N results (priority queue)                   │       │  │
│  │  │ - Compute scores (if needed)                               │       │  │
│  │  │ - Return: doc IDs + scores + sort values                   │       │  │
│  │  └────────────────────────────────────────────────────────────┘       │  │
│  │          │       │       │       │                                    │  │
│  │          │       │       │       │                                    │  │
│  │  Step 6: Gather on Coordinator                                         │  │
│  │          └───────┴───────┴───────┘                                    │  │
│  │                          │                                            │  │
│  │                          ▼                                            │  │
│  │  ┌────────────────────────────────────────────────────────────┐       │  │
│  │  │ - Merge shard results into global top N                    │       │  │
│  │  │ - Resolve document _source from _id                        │       │  │
│  │  │ - Apply global aggregations                                │       │  │
│  │  │ - Build final response                                     │       │  │
│  │  └────────────────────────────────────────────────────────────┘       │  │
│  │                          │                                            │  │
│  │                          ▼                                            │  │
│  │  ┌────────────────────────────────────────────────────────────┐       │  │
│  │  │                      Response                              │       │  │
│  │  │ {                                                          │       │  │
│  │  │   "took": 15,                                            │       │  │
│  │  │   "timed_out": false,                                    │       │  │
│  │  │   "hits": {                                              │       │  │
│  │  │     "total": { "value": 1000 },                          │       │  │
│  │  │     "hits": [ ... ]                                      │       │  │
│  │  │   },                                                     │       │  │
│  │  │   "aggregations": { ... }                                │       │  │
│  │  │ }                                                          │       │  │
│  │  └────────────────────────────────────────────────────────────┘       │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Fetch Phase (if needed)                             │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Query Phase 只返回 doc IDs 和排序值                                    │  │
│  │  Fetch Phase 获取实际的 _source 字段                                    │  │
│  │                                                                        │  │
│  │  ┌─────────────┐         ┌─────────────┐                              │  │
│  │  │ Doc ID: 1   │────────►│ Shard P0    │──► Get stored fields       │  │
│  │  │ Doc ID: 5   │────────►│ Shard P1    │──► Get stored fields       │  │
│  │  │ Doc ID: 3   │────────►│ Shard P2    │──► Get stored fields       │  │
│  │  └─────────────┘         └─────────────┘                              │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Query DSL Deep Dive

### 2.1 Query vs Filter Context

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Query vs Filter Context                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Query Context                                       │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │  Purpose: "How well does this document match?"                         │  │
│  │  Result:  Returns _score for relevance ranking                         │  │
│  │                                                                        │  │
│  │  Example:                                                              │  │
│  │  {                                                                     │  │
│  │    "query": {                                                          │  │
│  │      "match": {                                                        │  │
│  │        "title": "quick brown fox"                                      │  │
│  │      }                                                                 │  │
│  │    }                                                                   │  │
│  │  }                                                                     │  │
│  │                                                                        │  │
│  │  Scoring (BM25):                                                       │  │
│  │  score(q,d) = Σ IDF(q_i) * (f(q_i,d) * (k1 + 1)) /                     │  │
│  │               (f(q_i,d) + k1 * (1 - b + b * |d|/avgdl))                │  │
│  │                                                                        │  │
│  │  Where:                                                                │  │
│  │  - IDF = log(1 + (N - n(q) + 0.5) / (n(q) + 0.5))                     │  │
│  │  - N = total documents                                                 │  │
│  │  - n(q) = documents containing term                                    │  │
│  │  - f(q,d) = term frequency in doc                                      │  │
│  │  - |d| = document length                                               │  │
│  │  - avgdl = average document length                                     │  │
│  │  - k1, b = tuning parameters (default 1.2, 0.75)                       │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Filter Context                                      │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │  Purpose: "Does this document match?" (Yes/No)                         │  │
│  │  Result:  No scoring, binary inclusion                                 │  │
│  │  Optimization: Caching of filter bitsets                               │  │
│  │                                                                        │  │
│  │  Example:                                                              │  │
│  │  {                                                                     │  │
│  │    "query": {                                                          │  │
│  │      "bool": {                                                         │  │
│  │        "must": [                                                       │  │
│  │          { "match": { "title": "search" } }                            │  │
│  │        ],                                                              │  │
│  │        "filter": [                                                     │  │
│  │          { "term": { "status": "published" } },                        │  │
│  │          { "range": { "created_at": { "gte": "2024-01-01" } } }        │  │
│  │        ]                                                               │  │
│  │      }                                                                 │  │
│  │    }                                                                   │  │
│  │  }                                                                     │  │
│  │                                                                        │  │
│  │  Filter Bitset Caching:                                                │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Query 1: status:published ──► [1,0,1,1,0,1,0,1...] (bitset)   │  │  │
│  │  │  (cached for reuse)                                             │  │  │
│  │  │                                                                  │  │  │
│  │  │  Query 2: status:published AND created_at > X                   │  │  │
│  │  │  - Use cached bitset for status                                 │  │  │
│  │  │  - AND with range filter result                                 │  │  │
│  │  │  - Result: [1,0,0,1,0,1,0,0...]                                 │  │  │
│  │  │                                                                  │  │  │
│  │  │  Bitwise operations are extremely fast!                         │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Full-Text Queries

```go
package elasticsearch

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"

    "github.com/elastic/go-elasticsearch/v8"
    "github.com/elastic/go-elasticsearch/v8/esapi"
)

// MatchQuery 全文匹配查询
func MatchQuery(client *elasticsearch.Client, index, field, query string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                field: map[string]interface{}{
                    "query":     query,
                    "operator":  "or",       // or / and
                    "fuzziness": "AUTO",     // AUTO / 0 / 1 / 2
                    "prefix_length": 0,
                    "max_expansions": 50,
                    "minimum_should_match": "75%",
                },
            },
        },
        "highlight": map[string]interface{}{
            "fields": map[string]interface{}{
                field: map[string]interface{}{},
            },
        },
    }

    return executeSearch(client, index, q)
}

// MultiMatchQuery 多字段匹配
func MultiMatchQuery(client *elasticsearch.Client, index, query string, fields []string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "multi_match": map[string]interface{}{
                "query":  query,
                "fields": fields,           // e.g., ["title^3", "content"]
                "type":   "best_fields",    // best_fields / most_fields / cross_fields / phrase / phrase_prefix
                "tie_breaker": 0.3,
            },
        },
    }

    return executeSearch(client, index, q)
}

// MatchPhraseQuery 短语匹配
func MatchPhraseQuery(client *elasticsearch.Client, index, field, phrase string, slop int) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "match_phrase": map[string]interface{}{
                field: map[string]interface{}{
                    "query": phrase,
                    "slop":  slop,  // 允许的词语间隔
                },
            },
        },
    }

    return executeSearch(client, index, q)
}

// QueryStringQuery 查询字符串 (Lucene 语法)
func QueryStringQuery(client *elasticsearch.Client, index, query string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "query_string": map[string]interface{}{
                "query":            query,  // e.g., "title:(quick OR brown) AND fox"
                "default_field":    "content",
                "default_operator": "AND",
                "allow_leading_wildcard": false,
                "analyze_wildcard": true,
                "time_zone": "+08:00",
            },
        },
    }

    return executeSearch(client, index, q)
}

// SimpleQueryStringQuery 简化查询字符串 (容错)
func SimpleQueryStringQuery(client *elasticsearch.Client, index, query string, fields []string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "simple_query_string": map[string]interface{}{
                "query":  query,  // Supports + | - " * ()
                "fields": fields,
                "default_operator": "AND",
                "flags": "OR|AND|NOT|PREFIX|PHRASE|PRECEDENCE",  // 启用操作符
            },
        },
    }

    return executeSearch(client, index, q)
}
```

### 2.3 Term-Level Queries

```go
// TermQuery 精确值查询
func TermQuery(client *elasticsearch.Client, index, field string, value interface{}) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "term": map[string]interface{}{
                field: value,
            },
        },
    }

    return executeSearch(client, index, q)
}

// TermsQuery 多值精确匹配
func TermsQuery(client *elasticsearch.Client, index, field string, values []interface{}) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "terms": map[string]interface{}{
                field: values,
            },
        },
    }

    return executeSearch(client, index, q)
}

// RangeQuery 范围查询
func RangeQuery(client *elasticsearch.Client, index, field string, gte, lte interface{}) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "range": map[string]interface{}{
                field: map[string]interface{}{
                    "gte": gte,
                    "lte": lte,
                    // "gt": gt,
                    // "lt": lt,
                    // "format": "yyyy-MM-dd",
                    // "time_zone": "+08:00",
                },
            },
        },
    }

    return executeSearch(client, index, q)
}

// ExistsQuery 字段存在查询
func ExistsQuery(client *elasticsearch.Client, index, field string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "exists": map[string]interface{}{
                "field": field,
            },
        },
    }

    return executeSearch(client, index, q)
}

// PrefixQuery 前缀查询
func PrefixQuery(client *elasticsearch.Client, index, field, prefix string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "prefix": map[string]interface{}{
                field: prefix,
            },
        },
    }

    return executeSearch(client, index, q)
}

// WildcardQuery 通配符查询 (性能较低)
func WildcardQuery(client *elasticsearch.Client, index, field, pattern string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "wildcard": map[string]interface{}{
                field: pattern,  // e.g., "ki*y" or "ki?y"
            },
        },
    }

    return executeSearch(client, index, q)
}

// RegexpQuery 正则表达式查询
func RegexpQuery(client *elasticsearch.Client, index, field, pattern string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "regexp": map[string]interface{}{
                field: map[string]interface{}{
                    "value":         pattern,
                    "flags":         "ALL",  // INTERSECTION|COMPLEMENT|EMPTY
                    "max_determinized_states": 10000,
                },
            },
        },
    }

    return executeSearch(client, index, q)
}

// IdsQuery ID 批量查询
func IdsQuery(client *elasticsearch.Client, index string, ids []string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "ids": map[string]interface{}{
                "values": ids,
            },
        },
    }

    return executeSearch(client, index, q)
}
```

### 2.4 Compound Queries

```go
// BoolQuery 布尔组合查询
type BoolQuery struct {
    Must               []map[string]interface{}
    Filter             []map[string]interface{}
    Should             []map[string]interface{}
    MustNot            []map[string]interface{}
    MinimumShouldMatch interface{}
}

// Build 构建查询
func (bq *BoolQuery) Build() map[string]interface{} {
    boolClause := map[string]interface{}{}

    if len(bq.Must) > 0 {
        boolClause["must"] = bq.Must
    }
    if len(bq.Filter) > 0 {
        boolClause["filter"] = bq.Filter
    }
    if len(bq.Should) > 0 {
        boolClause["should"] = bq.Should
        if bq.MinimumShouldMatch != nil {
            boolClause["minimum_should_match"] = bq.MinimumShouldMatch
        }
    }
    if len(bq.MustNot) > 0 {
        boolClause["must_not"] = bq.MustNot
    }

    return map[string]interface{}{
        "query": map[string]interface{}{
            "bool": boolClause,
        },
    }
}

// ComplexSearch 复杂搜索示例
func ComplexSearch(client *elasticsearch.Client, index string) (*SearchResponse, error) {
    bq := &BoolQuery{
        Must: []map[string]interface{}{
            {
                "multi_match": map[string]interface{}{
                    "query":  "database performance",
                    "fields": []string{"title^3", "content"},
                },
            },
        },
        Filter: []map[string]interface{}{
            {
                "term": map[string]interface{}{"status": "published"},
            },
            {
                "range": map[string]interface{}{
                    "created_at": map[string]interface{}{
                        "gte": "now-7d/d",
                        "lte": "now/d",
                    },
                },
            },
            {
                "exists": map[string]interface{}{"field": "author"},
            },
        },
        MustNot: []map[string]interface{}{
            {
                "term": map[string]interface{}{"category": "spam"},
            },
        },
        Should: []map[string]interface{}{
            {
                "term": map[string]interface{}{"featured": true},
            },
            {
                "range": map[string]interface{}{
                    "rating": map[string]interface{}{"gte": 4.5},
                },
            },
        },
        MinimumShouldMatch: 1,
    }

    return executeSearch(client, index, bq.Build())
}

// FunctionScoreQuery 函数评分查询
func FunctionScoreQuery(client *elasticsearch.Client, index string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "function_score": map[string]interface{}{
                "query": map[string]interface{}{
                    "match": map[string]interface{}{
                        "title": "elasticsearch",
                    },
                },
                "functions": []map[string]interface{}{
                    {
                        "filter": map[string]interface{}{
                            "term": map[string]interface{}{"featured": true},
                        },
                        "weight": 2.0,
                    },
                    {
                        "field_value_factor": map[string]interface{}{
                            "field":    "popularity",
                            "factor":   1.2,
                            "modifier": "log1p",
                            "missing":  1.0,
                        },
                    },
                    {
                        "gauss": map[string]interface{}{
                            "created_at": map[string]interface{}{
                                "origin":        "now",
                                "scale":         "7d",
                                "offset":        "1d",
                                "decay":         0.5,
                            },
                        },
                    },
                },
                "score_mode":   "sum",    // multiply / sum / avg / first / max / min
                "boost_mode":   "multiply", // multiply / replace / sum / avg / max / min
                "max_boost":    100,
                "min_score":    1,
            },
        },
    }

    return executeSearch(client, index, q)
}

// NestedQuery 嵌套对象查询
func NestedQuery(client *elasticsearch.Client, index, path string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "query": map[string]interface{}{
            "nested": map[string]interface{}{
                "path": "comments",
                "query": map[string]interface{}{
                    "bool": map[string]interface{}{
                        "must": []map[string]interface{}{
                            {
                                "match": map[string]interface{}{
                                    "comments.content": "elasticsearch",
                                },
                            },
                            {
                                "range": map[string]interface{}{
                                    "comments.rating": map[string]interface{}{"gte": 4},
                                },
                            },
                        },
                    },
                },
                "inner_hits": map[string]interface{}{},  // 返回匹配的嵌套文档
            },
        },
    }

    return executeSearch(client, index, q)
}
```

### 2.5 Aggregations

```go
// TermsAggregation 分桶聚合
func TermsAggregation(client *elasticsearch.Client, index, field string, size int) (*SearchResponse, error) {
    q := map[string]interface{}{
        "size": 0,  // 不返回文档，只返回聚合结果
        "aggs": map[string]interface{}{
            "by_category": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": field,
                    "size":  size,
                    "order": map[string]interface{}{
                        "_count": "desc",
                    },
                    "min_doc_count": 1,
                },
            },
        },
    }

    return executeSearch(client, index, q)
}

// DateHistogramAggregation 日期直方图聚合
func DateHistogramAggregation(client *elasticsearch.Client, index, field string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "size": 0,
        "aggs": map[string]interface{}{
            "sales_over_time": map[string]interface{}{
                "date_histogram": map[string]interface{}{
                    "field":            field,
                    "calendar_interval": "month",  // minute / hour / day / week / month / quarter / year
                    "format":           "yyyy-MM-dd",
                    "min_doc_count":    0,
                    "extended_bounds": map[string]interface{}{
                        "min": "2024-01-01",
                        "max": "2024-12-31",
                    },
                },
            },
        },
    }

    return executeSearch(client, index, q)
}

// MultiMetricsAggregation 多指标聚合
func MultiMetricsAggregation(client *elasticsearch.Client, index string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "size": 0,
        "aggs": map[string]interface{}{
            "by_category": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": "category",
                },
                "aggs": map[string]interface{}{  // 子聚合
                    "avg_price": map[string]interface{}{
                        "avg": map[string]interface{}{
                            "field": "price",
                        },
                    },
                    "max_price": map[string]interface{}{
                        "max": map[string]interface{}{
                            "field": "price",
                        },
                    },
                    "min_price": map[string]interface{}{
                        "min": map[string]interface{}{
                            "field": "price",
                        },
                    },
                    "stats_price": map[string]interface{}{
                        "stats": map[string]interface{}{
                            "field": "price",
                        },
                    },
                    "percentiles_price": map[string]interface{}{
                        "percentiles": map[string]interface{}{
                            "field":    "price",
                            "percents": []float64{25, 50, 75, 95, 99},
                        },
                    },
                    "cardinality_products": map[string]interface{}{
                        "cardinality": map[string]interface{}{  // 去重计数
                            "field": "product_id",
                        },
                    },
                },
            },
        },
    }

    return executeSearch(client, index, q)
}

// PipelineAggregation 管道聚合
func PipelineAggregation(client *elasticsearch.Client, index string) (*SearchResponse, error) {
    q := map[string]interface{}{
        "size": 0,
        "aggs": map[string]interface{}{
            "sales_per_month": map[string]interface{}{
                "date_histogram": map[string]interface{}{
                    "field":            "sale_date",
                    "calendar_interval": "month",
                },
                "aggs": map[string]interface{}{
                    "total_sales": map[string]interface{}{
                        "sum": map[string]interface{}{
                            "field": "amount",
                        },
                    },
                    "sales_derivative": map[string]interface{}{  // 环比变化
                        "derivative": map[string]interface{}{
                            "buckets_path": "total_sales",
                        },
                    },
                    "moving_avg": map[string]interface{}{  // 移动平均
                        "moving_avg": map[string]interface{}{
                            "buckets_path": "total_sales",
                            "window":       3,
                            "model":        "linear",
                        },
                    },
                },
            },
        },
    }

    return executeSearch(client, index, q)
}
```

### 2.6 Go Client Implementation

```go
package elasticsearch

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "strings"
    "time"

    "github.com/elastic/go-elasticsearch/v8"
    "github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client ES 客户端封装
type Client struct {
    es      *elasticsearch.Client
    indices *IndexManager
}

// Config 客户端配置
type Config struct {
    Addresses []string
    Username  string
    Password  string
    APIKey    string
    CloudID   string

    // 重试配置
    MaxRetries    int
    RetryOnStatus []int

    // 连接配置
    TransportMaxIdleConns        int
    TransportMaxIdleConnsPerHost int
    TransportIdleConnTimeout     time.Duration
}

// NewClient 创建客户端
func NewClient(cfg *Config) (*Client, error) {
    esCfg := elasticsearch.Config{
        Addresses: cfg.Addresses,
        Username:  cfg.Username,
        Password:  cfg.Password,
        APIKey:    cfg.APIKey,
        CloudID:   cfg.CloudID,
    }

    if cfg.MaxRetries > 0 {
        esCfg.MaxRetries = cfg.MaxRetries
    }
    if len(cfg.RetryOnStatus) > 0 {
        esCfg.RetryOnStatus = cfg.RetryOnStatus
    }

    es, err := elasticsearch.NewClient(esCfg)
    if err != nil {
        return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
    }

    c := &Client{es: es}
    c.indices = &IndexManager{client: c}

    return c, nil
}

// Ping 健康检查
func (c *Client) Ping(ctx context.Context) error {
    res, err := c.es.Ping(c.es.Ping.WithContext(ctx))
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("ping failed: %s", res.String())
    }
    return nil
}

// SearchResponse 搜索响应
type SearchResponse struct {
    Took     int64                  `json:"took"`
    TimedOut bool                   `json:"timed_out"`
    Shards   ShardInfo              `json:"_shards"`
    Hits     HitsInfo               `json:"hits"`
    Aggregations map[string]interface{} `json:"aggregations,omitempty"`
}

type ShardInfo struct {
    Total      int `json:"total"`
    Successful int `json:"successful"`
    Skipped    int `json:"skipped"`
    Failed     int `json:"failed"`
}

type HitsInfo struct {
    Total    TotalValue    `json:"total"`
    MaxScore float64       `json:"max_score"`
    Hits     []HitDocument `json:"hits"`
}

type TotalValue struct {
    Value    int64  `json:"value"`
    Relation string `json:"relation"`
}

type HitDocument struct {
    Index   string                 `json:"_index"`
    Type    string                 `json:"_type"`
    ID      string                 `json:"_id"`
    Score   float64                `json:"_score"`
    Source  json.RawMessage        `json:"_source"`
    Highlight map[string][]string  `json:"highlight,omitempty"`
}

// Search 执行搜索
func (c *Client) Search(ctx context.Context, index string, query map[string]interface{}) (*SearchResponse, error) {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(query); err != nil {
        return nil, err
    }

    res, err := c.es.Search(
        c.es.Search.WithContext(ctx),
        c.es.Search.WithIndex(index),
        c.es.Search.WithBody(&buf),
        c.es.Search.WithTrackTotalHits(true),
        c.es.Search.WithPretty(),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.IsError() {
        return nil, fmt.Errorf("search error: %s", res.String())
    }

    var sr SearchResponse
    if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
        return nil, err
    }

    return &sr, nil
}

// IndexDocument 索引文档
func (c *Client) IndexDocument(ctx context.Context, index, id string, document interface{}) error {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(document); err != nil {
        return err
    }

    req := esapi.IndexRequest{
        Index:      index,
        DocumentID: id,
        Body:       &buf,
        Refresh:    "true",
    }

    res, err := req.Do(ctx, c.es)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("index error: %s", res.String())
    }

    return nil
}

// BulkIndex 批量索引
func (c *Client) BulkIndex(ctx context.Context, index string, documents map[string]interface{}) error {
    var buf bytes.Buffer

    for id, doc := range documents {
        // Action metadata
        meta := map[string]interface{}{
            "index": map[string]interface{}{
                "_index": index,
                "_id":    id,
            },
        }
        if err := json.NewEncoder(&buf).Encode(meta); err != nil {
            return err
        }

        // Document source
        if err := json.NewEncoder(&buf).Encode(doc); err != nil {
            return err
        }
    }

    res, err := c.es.Bulk(
        c.es.Bulk.WithContext(ctx),
        c.es.Bulk.WithBody(&buf),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("bulk error: %s", res.String())
    }

    return nil
}

// DeleteDocument 删除文档
func (c *Client) DeleteDocument(ctx context.Context, index, id string) error {
    req := esapi.DeleteRequest{
        Index:      index,
        DocumentID: id,
    }

    res, err := req.Do(ctx, c.es)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("delete error: %s", res.String())
    }

    return nil
}

// UpdateByQuery 批量更新
func (c *Client) UpdateByQuery(ctx context.Context, index string, query, script map[string]interface{}) error {
    body := map[string]interface{}{
        "query":  query,
        "script": script,
    }

    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(body); err != nil {
        return err
    }

    res, err := c.es.UpdateByQuery(
        []string{index},
        c.es.UpdateByQuery.WithContext(ctx),
        c.es.UpdateByQuery.WithBody(&buf),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("update by query error: %s", res.String())
    }

    return nil
}

// ScrollSearch 滚动搜索 (大数据集)
func (c *Client) ScrollSearch(ctx context.Context, index string, query map[string]interface{}, scrollTime time.Duration) (<-chan []HitDocument, <-chan error) {
    results := make(chan []HitDocument)
    errCh := make(chan error, 1)

    go func() {
        defer close(results)
        defer close(errCh)

        var buf bytes.Buffer
        if err := json.NewEncoder(&buf).Encode(query); err != nil {
            errCh <- err
            return
        }

        // Initial search
        res, err := c.es.Search(
            c.es.Search.WithContext(ctx),
            c.es.Search.WithIndex(index),
            c.es.Search.WithBody(&buf),
            c.es.Search.WithScroll(scrollTime),
            c.es.Search.WithSize(1000),
        )
        if err != nil {
            errCh <- err
            return
        }

        var sr SearchResponse
        if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
            res.Body.Close()
            errCh <- err
            return
        }
        res.Body.Close()

        scrollID := extractScrollID(&sr)

        for len(sr.Hits.Hits) > 0 {
            select {
            case results <- sr.Hits.Hits:
            case <-ctx.Done():
                c.clearScroll(ctx, scrollID)
                return
            }

            // Get next batch
            res, err = c.es.Scroll(
                c.es.Scroll.WithContext(ctx),
                c.es.Scroll.WithScrollID(scrollID),
                c.es.Scroll.WithScroll(scrollTime),
            )
            if err != nil {
                errCh <- err
                return
            }

            if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
                res.Body.Close()
                errCh <- err
                return
            }
            res.Body.Close()

            scrollID = extractScrollID(&sr)
        }

        c.clearScroll(ctx, scrollID)
    }()

    return results, errCh
}

func (c *Client) clearScroll(ctx context.Context, scrollID string) {
    if scrollID != "" {
        c.es.ClearScroll(
            c.es.ClearScroll.WithContext(ctx),
            c.es.ClearScroll.WithScrollID(scrollID),
        )
    }
}

func extractScrollID(sr *SearchResponse) string {
    // Extract from _scroll_id field in response
    return ""
}

// IndexManager 索引管理
type IndexManager struct {
    client *Client
}

// CreateIndex 创建索引
func (im *IndexManager) CreateIndex(ctx context.Context, index string, settings, mappings map[string]interface{}) error {
    body := map[string]interface{}{
        "settings":  settings,
        "mappings":  mappings,
    }

    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(body); err != nil {
        return err
    }

    res, err := im.client.es.Indices.Create(
        index,
        im.client.es.Indices.Create.WithContext(ctx),
        im.client.es.Indices.Create.WithBody(&buf),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("create index error: %s", res.String())
    }

    return nil
}

// DeleteIndex 删除索引
func (im *IndexManager) DeleteIndex(ctx context.Context, indices ...string) error {
    res, err := im.client.es.Indices.Delete(
        indices,
        im.client.es.Indices.Delete.WithContext(ctx),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("delete index error: %s", res.String())
    }

    return nil
}

// PutMapping 更新映射
func (im *IndexManager) PutMapping(ctx context.Context, index string, mapping map[string]interface{}) error {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
        return err
    }

    res, err := im.client.es.Indices.PutMapping(
        []string{index},
        im.client.es.Indices.PutMapping.WithContext(ctx),
        im.client.es.Indices.PutMapping.WithBody(&buf),
    )
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("put mapping error: %s", res.String())
    }

    return nil
}

// Indices 返回索引管理器
func (c *Client) Indices() *IndexManager {
    return c.indices
}

// 辅助函数
func executeSearch(client *elasticsearch.Client, index string, query map[string]interface{}) (*SearchResponse, error) {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(query); err != nil {
        return nil, err
    }

    res, err := client.Search(
        client.Search.WithIndex(index),
        client.Search.WithBody(&buf),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.IsError() {
        body, _ := io.ReadAll(res.Body)
        return nil, fmt.Errorf("search error: %s", string(body))
    }

    var sr SearchResponse
    if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
        return nil, err
    }

    return &sr, nil
}
```

---

## 3. Configuration Best Practices

### 3.1 Index Settings

```json
{
  "settings": {
    "number_of_shards": 5,
    "number_of_replicas": 1,

    "refresh_interval": "5s",
    "translog": {
      "durability": "async",
      "sync_interval": "5s"
    },

    "index": {
      "store": {
        "preload": ["nvd", "dvd"]
      },
      "codec": "best_compression",
      "max_result_window": 10000,
      "max_terms_count": 65536
    },

    "analysis": {
      "analyzer": {
        "custom_ik": {
          "type": "custom",
          "tokenizer": "ik_max_word",
          "filter": ["lowercase", "synonym_filter"]
        }
      },
      "filter": {
        "synonym_filter": {
          "type": "synonym_graph",
          "synonyms_path": "analysis/synonyms.txt",
          "updateable": true
        }
      }
    }
  },
  "mappings": {
    "dynamic": "strict",
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "custom_ik",
        "search_analyzer": "ik_smart",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "created_at": {
        "type": "date",
        "format": "strict_date_optional_time||epoch_millis"
      },
      "price": {
        "type": "scaled_float",
        "scaling_factor": 100
      },
      "tags": {
        "type": "keyword"
      },
      "location": {
        "type": "geo_point"
      }
    }
  }
}
```

---

## 4. Performance Tuning

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Elasticsearch Performance Tuning                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Indexing Optimization                               │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. Bulk API (批量写入)                                                │  │
│  │     - 批量大小: 5-15 MB 每批                                          │  │
│  │     - 使用 _bulk endpoint                                             │  │
│  │                                                                        │  │
│  │  2. Refresh Strategy                                                   │  │
│  │     - 大批量导入时: refresh_interval = -1 (禁用)                      │  │
│  │     - 实时搜索场景: refresh_interval = 1s                             │  │
│  │                                                                        │  │
│  │  3. Translog Durability                                                │  │
│  │     - 高吞吐: durability = async                                      │  │
│  │     - 高可靠: durability = request                                    │  │
│  │                                                                        │  │
│  │  4. Mapping Optimization                                               │  │
│  │     - 禁用 _all 字段 (ES 6.0+ 已默认禁用)                             │  │
│  │     - 禁用 dynamic mapping (index.mapping.total_fields.limit)         │  │
│  │     - 使用 keyword 替代 text 用于聚合/排序                            │  │
│  │                                                                        │  │
│  │  5. Segment Merge Tuning                                               │  │
│  │     - index.merge.policy.max_merge_at_once: 10                        │  │
│  │     - index.merge.policy.segments_per_tier: 10                        │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Search Optimization                                 │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. Query Caching                                                      │  │
│  │     - Filter context 自动缓存                                         │  │
│  │     - 使用 constant_score 包装 filter                                 │  │
│  │                                                                        │  │
│  │  2. Avoid Deep Paging                                                  │  │
│  │     - 使用 search_after 替代 from/size                                │  │
│  │     - 使用 scroll API 导出大数据                                      │  │
│  │                                                                        │  │
│  │  3. Routing                                                            │  │
│  │     - 使用 _routing 将相关文档路由到同一分片                          │  │
│  │     - 减少跨分片查询                                                  │  │
│  │                                                                        │  │
│  │  4. Eager Global Ordinals                                              │  │
│  │     - 高基数字段聚合时使用                                            │  │
│  │     - 在 refresh 时预构建                                             │  │
│  │                                                                        │  │
│  │  5. Doc Values vs Field Data                                          │  │
│  │     - 优先使用 doc_values (磁盘)                                      │  │
│  │     - fielddata = true 仅用于 text 字段聚合                           │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Cluster Sizing                                      │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  内存建议:                                                             │  │
│  │  - JVM Heap: 不超过 32GB (压缩指针优化)                               │  │
│  │  - OS Cache: 剩余全部内存 (用于文件缓存)                              │  │
│  │  - 比例: 1:24 (Heap:Total RAM)                                        │  │
│  │                                                                        │  │
│  │  分片建议:                                                             │  │
│  │  - 每节点分片数 < 20 * heap(GB)                                       │  │
│  │  - 分片大小: 20-50 GB                                                 │  │
│  │  - 避免过多小分片 (merge overhead)                                    │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Visual Representations

### 5.1 Inverted Index with FST

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Inverted Index & FST Structure                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Term Dictionary (FST)                               │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │                    ┌─►[u]─►[c]─►[k]─►[y]────► Postings Ptr            │  │
│  │                   /                                                   │  │
│  │  [r]─►[u]─►[n]─►[n]─►[i]─►[n]─►[g]────► Postings Ptr                │  │
│  │                   \                                                   │  │
│  │                    └─►[s]─►[t]────────────────► Postings Ptr          │  │
│  │                                                                        │  │
│  │  [s]─►[l]─►[o]─►[w]──────────────────────────► Postings Ptr           │  │
│  │                                                                        │  │
│  │  [t]─►[h]─►[e]───────────────────────────────► Postings Ptr           │  │
│  │                                                                        │  │
│  │  FST Properties:                                                       │  │
│  │  - Shared prefixes compressed (run = r-u-n prefix)                   │  │
│  │  - Sorted order traversal                                              │  │
│  │  - Can seek to arbitrary term                                          │  │
│  │  - Can enumerate terms with given prefix                               │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Postings List Structure                             │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Term: "fox"                                                           │  │
│  │                                                                        │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Doc ID │ Term Freq │ Positions        │ Offsets               │  │  │
│  │  ├─────────────────────────────────────────────────────────────────┤  │  │
│  │  │   1     │     1     │ [3]              │ [(16, 19)]            │  │  │
│  │  │   2     │     1     │ [3]              │ [(15, 18)]            │  │  │
│  │  │   4     │     2     │ [1, 5]           │ [(4, 7), (22, 25)]    │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Encoding: Delta + Variable Byte                                       │  │
│  │  Doc IDs: 1, 2, 4 ──► 1, 1, 2 (deltas) ──► VB encoded bytes          │  │
│  │                                                                        │  │
│  │  Skip List (for faster conjunctions):                                  │  │
│  │  Level 2: ───────► 4                                                   │  │
│  │  Level 1: ─────► 2 ──────► 4                                           │  │
│  │  Level 0: 1 ──► 2 ──► 4                                                │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Query Execution Plan

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Query Execution Plan Visualization                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Query: (title:elasticsearch OR content:elasticsearch) AND status:published │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Query Tree                                          │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │                           AND (MUST)                                   │  │
│  │                          /            \                                │  │
│  │                         /              \                               │  │
│  │                       OR               TERM                           │  │
│  │                     /    \            status:"published"               │  │
│  │                   /        \                                           │  │
│  │               TERM        TERM                                         │  │
│  │         title:"es"    content:"es"                                    │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Execution Flow                                      │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  Step 1: Filter Cache Check                                            │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  status:"published" in cache?                                     │  │  │
│  │  │  ├── YES ──► Use cached bitset [1,0,1,1,0,1,1,0...]             │  │  │
│  │  │  └── NO  ──► Execute term query, cache result                   │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Step 2: Query Execution                                               │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  title:"elasticsearch"                                            │  │  │
│  │  │  └──► Union postings: [1, 5, 10, 15, 23, 45...]                  │  │  │
│  │  │                                                                   │  │  │
│  │  │  content:"elasticsearch"                                          │  │  │
│  │  │  └──► Union postings: [3, 5, 12, 15, 23, 48...]                  │  │  │
│  │  │                                                                   │  │  │
│  │  │  OR (title OR content)                                            │  │  │
│  │  │  └──► Merge with skip lists: [1, 3, 5, 10, 12, 15, 23, 45, 48...]│  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Step 3: Filter Application (Bitwise AND)                              │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Query Result:  [1, 3, 5, 10, 12, 15, 23, 45, 48...]            │  │  │
│  │  │  Cached Filter: [1, 0, 1, 1, 0, 1, 1, 0, 1...]  (published)    │  │  │
│  │  │  ─────────────────────────────────────────────────              │  │  │
│  │  │  AND Result:    [1, 0, 5, 10, 0, 15, 23, 0, 48...]              │  │  │
│  │  │                                                                   │  │  │
│  │  │  Final: [1, 5, 10, 15, 23, 48...]                                │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Step 4: Scoring (if needed)                                           │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  For each doc in result:                                          │  │  │
│  │  │  score = BM25(title) + BM25(content)                              │  │  │
│  │  │                                                                   │  │  │
│  │  │  Title match weight: 3.0 (field boost)                           │  │  │
│  │  │  Content match weight: 1.0                                        │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  │  Step 5: Top N Collection                                              │  │
│  │  ┌─────────────────────────────────────────────────────────────────┐  │  │
│  │  │  Min-Heap (size = from + size)                                    │  │  │
│  │  │                                                                   │  │  │
│  │  │  Insert each scored doc into heap                                 │  │  │
│  │  │  If heap full, compare with min, replace if higher                │  │  │
│  │  │                                                                   │  │  │
│  │  │  Final: Top N documents by score                                  │  │  │
│  │  └─────────────────────────────────────────────────────────────────┘  │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.3 Cluster State Management

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Elasticsearch Cluster State Management                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    Cluster State Structure                             │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  {                                                                     │  │
│  │    "cluster_uuid": "abc123...",                                        │  │
│  │    "version": 1523,           ◄─── 每次变更递增                        │  │
│  │                                                                        │  │
│  │    "nodes": {               ◄─── 集群节点信息                          │  │
│  │      "node_1": { name: "es-1", attributes: {...} },                   │  │
│  │      "node_2": { name: "es-2", attributes: {...} },                   │  │
│  │      "node_3": { name: "es-3", attributes: {...} }                    │  │
│  │    },                                                                  │  │
│  │                                                                        │  │
│  │    "metadata": {            ◄─── 索引元数据                            │  │
│  │      "indices": {                                                      │  │
│  │        "products": {                                                   │  │
│  │          "state": "open",                                               │  │
│  │          "settings": {...},                                             │  │
│  │          "mappings": {...},                                             │  │
│  │          "aliases": [...]                                               │  │
│  │        }                                                               │  │
│  │      }                                                                 │  │
│  │    },                                                                  │  │
│  │                                                                        │  │
│  │    "routing_table": {       ◄─── 分片路由信息                          │  │
│  │      "indices": {                                                      │  │
│  │        "products": {                                                   │  │
│  │          "shards": {                                                   │  │
│  │            "0": [                                                      │  │
│  │              { shard: 0, primary: true, node: "node_1", state: "STARTED" },│  │
│  │              { shard: 0, primary: false, node: "node_2", state: "STARTED" }│  │
│  │            ],                                                          │  │
│  │            "1": [...],                                                 │  │
│  │            "2": [...]                                                  │  │
│  │          }                                                             │  │
│  │        }                                                               │  │
│  │      }                                                                 │  │
│  │    },                                                                  │  │
│  │                                                                        │  │
│  │    "routing_nodes": {       ◄─── 节点视角的路由                        │  │
│  │      "node_1": [ { shard: "products", shard: 0, primary: true }, ... ],│  │
│  │      "node_2": [ ... ],                                                │  │
│  │      "unassigned": [ ... ]                                             │  │
│  │    }                                                                   │  │
│  │  }                                                                     │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    State Update Flow                                   │  │
│  ├───────────────────────────────────────────────────────────────────────┤  │
│  │                                                                        │  │
│  │  1. Create Index Request                                               │  │
│  │     Client ──► Master Node                                             │  │
│  │                                                                        │  │
│  │  2. Master Updates State                                               │  │
│  │     - Add index to metadata                                            │  │
│  │     - Create routing entries                                           │  │
│  │     - Increment version                                                │  │
│  │                                                                        │  │
│  │  3. Publish to All Nodes                                               │  │
│  │     Master ──► Node 1, Node 2, Node 3                                  │  │
│  │                                                                        │  │
│  │  4. Nodes Acknowledge                                                  │  │
│  │     Node 1, Node 2, Node 3 ──► Master                                  │  │
│  │                                                                        │  │
│  │  5. Cluster Health Update                                              │  │
│  │     - Wait for shard allocation                                        │  │
│  │     - Update health status                                             │  │
│  │                                                                        │  │
│  │  6. Response to Client                                                 │  │
│  │     Master ──► Client (acknowledged)                                   │  │
│  │                                                                        │  │
│  │  Failure Handling:                                                     │  │
│  │  - If node doesn't ack, master retries                                 │  │
│  │  - If master fails, new master elected                                 │  │
│  │  - State has version, stale updates rejected                           │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. References

1. **Gormley, C., & Tong, Z.** (2015). Elasticsearch: The Definitive Guide. O'Reilly Media.
2. **McCandless, M., Hatcher, E., & Gospodnetic, O.** (2010). Lucene in Action, Second Edition. Manning Publications.
3. **Elastic Documentation** (2024). elasticsearch.org
4. **Baeza-Yates, R., & Ribeiro-Neto, B.** (2011). Modern Information Retrieval. Addison-Wesley.

---

*Document Version: 1.0 | Last Updated: 2024*
