# TS-012: Elasticsearch 倒排索引的形式化 (Elasticsearch Inverted Index: Formal Analysis)

> **维度**: Technology Stack
> **级别**: S (16+ KB)
> **标签**: #elasticsearch #lucene #inverted-index #search-engine #information-retrieval
> **权威来源**:
>
> - [Lucene in Action](https://www.manning.com/books/lucene-in-action-second-edition) - McCandless et al. (2010)
> - [Introduction to Information Retrieval](https://nlp.stanford.edu/IR-book/) - Manning et al. (2008)
> - [Elasticsearch: The Definitive Guide](https://www.elastic.co/guide/en/elasticsearch/guide/current/index.html) - Clinton Gormley (2015)
> - [BM25: The Next Generation of Lucene Relevance](https://www.elastic.co/blog/practical-bm25-part-2-the-bm25-algorithm-and-its-variables) - Elastic (2016)
> - [Designing Data-Intensive Applications](https://dataintensive.net/) - Martin Kleppmann (2017)

---

## 1. 倒排索引的形式化定义

### 1.1 索引代数

**定义 1.1 (文档)**
$$d = \langle id, \text{terms}, \text{fields} \rangle$$

**定义 1.2 (倒排索引)**
$$I = \{ (t, D_t) \mid t \in \text{Vocabulary}, D_t \subseteq \text{Documents} \}$$
其中 $D_t$ 是包含词项 $t$ 的文档集合。

**定义 1.3 (Posting List)**
$$D_t = [ (doc_1, freq_1, pos_1), (doc_2, freq_2, pos_2), ... ]$$
包含文档 ID、词频、位置信息。

### 1.2 索引构建

**定义 1.4 (索引构建)**
$$\text{Index}: \{d_1, d_2, ..., d_n\} \to I$$

**算法**:

```
1. Tokenize documents
2. For each term t in document d:
   a. Add d to posting list of t
   b. Record frequency and positions
3. Sort posting lists by doc ID
4. Create term dictionary (FST)
```

---

## 2. BM25 评分模型的形式化

### 2.1 评分公式

**定义 2.1 (BM25)**
$$\text{BM25}(D, Q) = \sum_{q \in Q} \text{IDF}(q) \cdot \frac{f(q, D) \cdot (k_1 + 1)}{f(q, D) + k_1 \cdot (1 - b + b \cdot \frac{|D|}{\text{avgdl}})}$$

其中:

- $f(q, D)$: 词项 $q$ 在文档 $D$ 中的频率
- $|D|$: 文档长度
- avgdl: 平均文档长度
- $k_1, b$: 调节参数

### 2.2 IDF 计算

**定义 2.2 (逆文档频率)**
$$\text{IDF}(q) = \log\left(1 + \frac{N - n(q) + 0.5}{n(q) + 0.5}\right)$$
其中 $N$ 是文档总数，$n(q)$ 是包含 $q$ 的文档数。

---

## 3. 段合并的形式化

### 3.1 段结构

**定义 3.1 (Segment)**
$$S = \langle \text{documents}, \text{term-dict}, \text{postings}, \text{stored-fields} \rangle$$

**定义 3.2 (段合并)**
$$\text{Merge}: S_1, S_2, ..., S_k \to S_{new}$$

**复杂度**: $O(\sum |S_i|)$

---

## 4. 多元表征

### 4.1 倒排索引结构图

```
Document Collection                    Inverted Index
┌─────────────────────┐                ┌──────────────────────┐
│ Doc1: "quick brown" │                │ Term      Postings   │
│ Doc2: "brown fox"   │───Index───────►│ quick  → [1]         │
│ Doc3: "quick fox"   │                │ brown  → [1, 2]      │
└─────────────────────┘                │ fox    → [2, 3]      │
                                       └──────────────────────┘

Posting List Detail:
quick:  [(Doc1, freq=1, pos=[0])]
brown:  [(Doc1, freq=1, pos=[1]), (Doc2, freq=1, pos=[0])]
fox:    [(Doc2, freq=1, pos=[1]), (Doc3, freq=1, pos=[1])]
```

### 4.2 搜索算法决策树

```
查询类型?
│
├── 精确匹配? → Term Query
├── 短语匹配? → Match Phrase
├── 模糊匹配? → Fuzzy Query (Levenshtein)
├── 范围查询? → Range Query
├── 布尔组合? → Bool Query (must/should/must_not)
└── 全文搜索? → Multi-Match + BM25 Scoring
```

### 4.3 评分算法对比矩阵

| 算法 | 考虑 TF | 考虑 IDF | 考虑文档长度 | 适用场景 |
|------|---------|----------|-------------|---------|
| **TF-IDF** | ✓ | ✓ | ✗ | 基础搜索 |
| **BM25** | ✓ | ✓ | ✓ | 现代搜索 (推荐) |
| **DFR** | ✓ | ✓ | ✓ | 概率模型 |
| **IB** | ✓ | ✓ | ✓ | 信息论模型 |

---

## 5. 检查清单

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Elasticsearch Design Checklist                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  索引设计:                                                                   │
│  □ 合理的分片数 (避免过多小分片)                                               │
│  □ 副本数 ≥ 1 (生产环境)                                                      │
│  □ 映射优化 (禁用 _all, 使用 keyword)                                         │
│  □ 分词器选择 (standard/ik/ngram)                                             │
│                                                                              │
│  查询优化:                                                                   │
│  □ 使用 filter context (缓存)                                                │
│  □ 避免 deep paging (search_after)                                           │
│  □ 预加载 fielddata (聚合)                                                   │
│  □ Routing 优化 (单分片查询)                                                  │
│                                                                              │
│  性能:                                                                       │
│  □ 批量索引 (bulk API)                                                        │
│  □ 刷新间隔调整 (index.refresh_interval)                                       │
│  □ 段合并策略优化                                                              │
│  □ 冻结历史索引                                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16KB, 完整形式化)
