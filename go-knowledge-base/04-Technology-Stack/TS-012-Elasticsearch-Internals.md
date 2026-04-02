# TS-012: Elasticsearch 内部机制 (Elasticsearch Internals)

> **维度**: Technology Stack
> **级别**: S (17+ KB)
> **标签**: #elasticsearch #search-engine #lucene #inverted-index
> **权威来源**: [Elasticsearch Guide](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html), [Lucene](https://lucene.apache.org/core/documentation.html)
> **版本**: Elasticsearch 9.0+

---

## 架构概述

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Elasticsearch Cluster                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      Cluster: es-prod                               │    │
│  │                                                                      │    │
│  │  Master-Eligible Nodes               Data Nodes                     │    │
│  │  ┌─────────────┐  ┌─────────────┐    ┌─────────────┐  ┌────────────┐│    │
│  │  │  master-1   │  │  master-2   │    │  data-1     │  │  data-2    ││    │
│  │  │  (Active)   │  │  (Standby)  │    │  Hot Tier   │  │  Warm Tier ││    │
│  │  └─────────────┘  └─────────────┘    └─────────────┘  └────────────┘│    │
│  │                                                                      │    │
│  │  角色:                                                               │    │
│  │  - master: 集群管理、索引创建、节点发现                                  │    │
│  │  - data: 存储数据、执行搜索                                             │    │
│  │  - ingest: 文档预处理                                                   │    │
│  │  - coordinating: 请求路由、聚合                                         │    │
│  │  - remote_cluster_client: 跨集群搜索                                    │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  分片分配:                                                                   │
│  Index: logs-2026.04.01                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Shard-0 (P) │  │ Shard-1 (P) │  │ Shard-2 (P) │  │ Shard-3 (P) │  Primary │
│  │ Shard-2 (R) │  │ Shard-3 (R) │  │ Shard-0 (R) │  │ Shard-1 (R) │  Replica │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘         │
│     data-1          data-2          data-1          data-2                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 倒排索引 (Inverted Index)

### 索引结构

```
文档集合:
┌─────────────────────────────────────────────────────────────────────────────┐
│ DocID │ Content                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│   1   │ "The quick brown fox"                                               │
│   2   │ "Jumped over the lazy dog"                                          │
│   3   │ "The quick dog"                                                     │
└─────────────────────────────────────────────────────────────────────────────┘

倒排索引:
┌─────────────────────────────────────────────────────────────────────────────┐
│ Term      │ DocIDs (Posting List)           │ TermFreq │ Positions          │
├─────────────────────────────────────────────────────────────────────────────┤
│ the       │ [1, 2, 3]                       │ 1, 1, 1  │ [1, 3, 1]          │
│ quick     │ [1, 3]                          │ 1, 1     │ [2, 2]             │
│ brown     │ [1]                             │ 1        │ [3]                │
│ fox       │ [1]                             │ 1        │ [4]                │
│ jumped    │ [2]                             │ 1        │ [1]                │
│ over      │ [2]                             │ 1        │ [2]                │
│ lazy      │ [2]                             │ 1        │ [4]                │
│ dog       │ [2, 3]                          │ 1, 1     │ [5, 3]             │
└─────────────────────────────────────────────────────────────────────────────┘

附加数据结构:
- Term Dictionary: 词项排序列表 (FST - Finite State Transducer)
- Postings List: 文档 ID 列表 (FOR/RoaringBitmap 压缩)
- Term Frequency: 词频
- Position: 词项位置 (短语查询)
- Offset: 字符偏移 (高亮)
```

### 分析链 (Analysis Chain)

```
输入文本: "The Quick Brown Foxes!"
    │
    ▼
Character Filters (字符过滤器)
    │ - HTML Strip
    │ - Mapping ("$" → "dollar")
    ▼
"The Quick Brown Foxes"
    │
    ▼
Tokenizer (分词器)
    │ - Standard: 按单词分割
    ▼
["The", "Quick", "Brown", "Foxes"]
    │
    ▼
Token Filters (词项过滤器)
    │ - Lowercase: ["the", "quick", "brown", "foxes"]
    │ - Stop: 移除停用词 ["quick", "brown", "foxes"] (可选)
    │ - Stemmer: ["quick", "brown", "fox"]
    ▼
输出词项: ["the", "quick", "brown", "fox"]
```

---

## Lucene 段 (Segments)

### 段结构

```
Index (索引):
├── _0.cfe          # 复合文件入口
├── _0.cfs          # 复合文件 (存储域数据)
├── _0.si           # 段信息
├── _0_1.liv        # 存活文档 (删除标记)
├── _1.cfe          # 另一个段
├── _1.cfs
├── _1.si
└── segments_2      # 段列表 (提交点)

段内文件:
- .tim: Term 索引
- .tip: Term 字典
- .doc: Postings (文档 ID)
- .pos: Positions
- .pay: Payloads
- .dvd: DocValues
- .dvm: DocValues 元数据

段合并 (Segment Merge):
- 后台线程持续执行
- 策略: TieredMergePolicy (默认)
- 目标: 减少段数量，优化查询性能
```

---

## 搜索执行流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Query Execution                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 查询解析 (Query Parsing)                                                  │
│     "title:elasticsearch AND body:search"                                    │
│              │                                                               │
│              ▼                                                               │
│     BooleanQuery                                                             │
│     ├── TermQuery (title:elasticsearch)                                      │
│     └── TermQuery (body:search)                                              │
│                                                                              │
│  2. 分片路由 (Shard Routing)                                                  │
│     - 计算路由: shard = hash(routing) % num_shards                          │
│     - 广播到所有相关分片 (DFS Query Then Fetch / Query Then Fetch)            │
│                                                                              │
│  3. 分片级别搜索                                                              │
│     ┌─────────────────────────────────────────────────────────────────┐      │
│     │ Shard 0                                                        │      │
│     │  - 加载 Term Dictionary (FST)                                  │      │
│     │  - 定位 Postings List                                          │      │
│     │  - 计算相关性评分 (TF/IDF, BM25)                                │      │
│     │  - 返回 Top N 结果                                             │      │
│     └─────────────────────────────────────────────────────────────────┘      │
│                                                                              │
│  4. 结果聚合 (Reduce Phase)                                                   │
│     - 收集所有分片结果                                                        │
│     - 全局排序                                                                │
│     - 聚合 (Aggregations)                                                     │
│                                                                              │
│  5. Fetch Phase (如果需要文档内容)                                             │
│     - 根据 ID 获取完整文档                                                    │
│     - 高亮处理                                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### BM25 评分算法

```go
package es

// BM25 评分公式
// Score(D,Q) = Σ IDF(q) * (f(q,D) * (k1 + 1)) / (f(q,D) + k1 * (1 - b + b * |D|/avgDL))
//
// IDF(q) = log(1 + (N - n(q) + 0.5) / (n(q) + 0.5))
//
// 其中:
// - N: 文档总数
// - n(q): 包含词项 q 的文档数
// - f(q,D): 词项 q 在文档 D 中的频率
// - |D|: 文档 D 的长度
// - avgDL: 平均文档长度
// - k1: 词频饱和度参数 (默认 1.2)
// - b: 长度归一化参数 (默认 0.75)

// BM25Score 计算 BM25 分数
func BM25Score(tf, docLen, avgDocLen float64, docFreq, totalDocs int, k1, b float64) float64 {
    // IDF
    nq := float64(docFreq)
    N := float64(totalDocs)
    idf := math.Log(1 + (N - nq + 0.5) / (nq + 0.5))

    // TF 部分
    tfPart := (tf * (k1 + 1)) / (tf + k1 * (1 - b + b * docLen / avgDocLen))

    return idf * tfPart
}
```

---

## Go 客户端示例

```go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "log"

    "github.com/elastic/go-elasticsearch/v9"
    "github.com/elastic/go-elasticsearch/v9/esapi"
)

type Document struct {
    Title   string   `json:"title"`
    Content string   `json:"content"`
    Tags    []string `json:"tags"`
}

func main() {
    // 创建客户端
    cfg := elasticsearch.Config{
        Addresses: []string{"http://localhost:9200"},
        Username:  "elastic",
        Password:  "password",
    }

    es, err := elasticsearch.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 索引文档
    doc := Document{
        Title:   "Elasticsearch Guide",
        Content: "Learn Elasticsearch internals...",
        Tags:    []string{"search", "database"},
    }

    data, _ := json.Marshal(doc)

    req := esapi.IndexRequest{
        Index:      "articles",
        DocumentID: "1",
        Body:       bytes.NewReader(data),
        Refresh:    "true",
    }

    res, err := req.Do(context.Background(), es)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    // 搜索
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "multi_match": map[string]interface{}{
                "query":  "elasticsearch",
                "fields": []string{"title^3", "content"},
            },
        },
        "highlight": map[string]interface{}{
            "fields": map[string]interface{}{
                "content": map[string]interface{}{},
            },
        },
    }

    var buf bytes.Buffer
    json.NewEncoder(&buf).Encode(query)

    searchRes, err := es.Search(
        es.Search.WithContext(context.Background()),
        es.Search.WithIndex("articles"),
        es.Search.WithBody(&buf),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer searchRes.Body.Close()

    log.Println(searchRes)
}
```

---

## 性能优化

| 优化项 | 建议 |
|--------|------|
| 分片数 | 单分片 20-50GB，避免过多小分片 |
| 副本数 | 生产环境至少 1 个副本 |
| 刷新间隔 | 日志场景 `refresh_interval: 30s` |
| 批量写入 | 使用 `_bulk` API，批次 5-15MB |
| 映射优化 | 禁用 `_all` 字段，使用 `keyword` 聚合 |
| 冻结索引 | 历史数据使用冻结索引 |

---

## 参考文献

1. [Elasticsearch Reference](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
2. [Lucene Documentation](https://lucene.apache.org/core/documentation.html)
3. [BM25: The Next Generation of Lucene Relevance](https://www.elastic.co/blog/practical-bm25-part-2-the-bm25-algorithm-and-its-variables)
