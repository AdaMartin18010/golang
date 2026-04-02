# TS-DB-008: Vector Databases

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #vector-database #embeddings #similarity-search #pgvector #pinecone
> **权威来源**:
>
> - [pgvector](https://github.com/pgvector/pgvector) - PostgreSQL vector extension
> - [Vector Database Guide](https://www.pinecone.io/learn/vector-database/) - Pinecone

---

## 1. Vector Database Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Vector Database Architecture                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Traditional Database vs Vector Database:                                    │
│                                                                              │
│  Traditional Query:                         Vector Query:                    │
│  SELECT * FROM products                     SELECT * FROM images             │
│  WHERE category = 'electronics'             ORDER BY embedding <->           │
│  AND price < 1000;                          '[0.1, 0.2, ...]' LIMIT 5;      │
│  (Exact match)                              (Similarity search)              │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                       Vector Space                                  │   │
│  │                                                                     │   │
│  │                         ▲                                           │   │
│  │                        /│\                                          │   │
│  │                       / │ \                                         │   │
│  │                      /  │  \                                        │   │
│  │                     /   ●   \     Query vector                      │   │
│  │                    /  /│\    \                                      │   │
│  │                   /  / │ \    \                                     │   │
│  │                  ●  /  │  \    ●   Nearest neighbors                │   │
│  │                v1  /   │   \   v2                                   │   │
│  │                   /    │    \                                       │   │
│  │                  ●     │     ●                                      │   │
│  │                v3      │     v4                                     │   │
│  │                        ●                                            │   │
│  │                       v5                                            │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Key Concepts:                                                               │
│  - Embedding: High-dimensional vector representation (e.g., 384, 768, 1536) │
│  - Distance Metric: Euclidean (L2), Cosine similarity, Dot product          │
│  - Approximate Nearest Neighbor (ANN): Fast similarity search               │
│  - Index Types: HNSW, IVFFlat, IVFPQ                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. pgvector with PostgreSQL

### 2.1 Setup and Configuration

```sql
-- Install pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create table with vector column
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    embedding VECTOR(1536)  -- OpenAI embedding dimension
);

-- Create index for similarity search
-- HNSW index: Good balance of speed and recall
CREATE INDEX ON items USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- IVFFlat index: Good for static datasets
CREATE INDEX ON items USING ivfflat (embedding vector_l2_ops)
WITH (lists = 100);
```

### 2.2 Go Integration

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/pgvector/pgvector-go"
)

const (
    dimension = 1536
)

type Item struct {
    ID          int32
    Name        string
    Description string
    Embedding   []float32
}

func main() {
    ctx := context.Background()

    // Connect to database
    pool, err := pgxpool.New(ctx, "postgres://user:password@localhost/db")
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // Create table
    if err := createTable(ctx, pool); err != nil {
        log.Fatal(err)
    }

    // Insert item with embedding
    embedding := make([]float32, dimension)
    // ... populate embedding from ML model

    item := &Item{
        Name:        "Example Item",
        Description: "This is an example",
        Embedding:   embedding,
    }

    if err := insertItem(ctx, pool, item); err != nil {
        log.Fatal(err)
    }

    // Search similar items
    queryEmbedding := make([]float32, dimension)
    // ... populate from query

    results, err := searchSimilar(ctx, pool, queryEmbedding, 5)
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range results {
        fmt.Printf("ID: %d, Name: %s, Distance: %f\n", r.ID, r.Name, r.Distance)
    }
}

func createTable(ctx context.Context, pool *pgxpool.Pool) error {
    _, err := pool.Exec(ctx, `
        CREATE EXTENSION IF NOT EXISTS vector;

        CREATE TABLE IF NOT EXISTS items (
            id SERIAL PRIMARY KEY,
            name TEXT,
            description TEXT,
            embedding VECTOR(1536)
        );

        CREATE INDEX IF NOT EXISTS items_embedding_idx
        ON items USING ivfflat (embedding vector_cosine_ops)
        WITH (lists = 100);
    `)
    return err
}

func insertItem(ctx context.Context, pool *pgxpool.Pool, item *Item) error {
    _, err := pool.Exec(ctx, `
        INSERT INTO items (name, description, embedding)
        VALUES ($1, $2, $3)
    `, item.Name, item.Description, pgvector.NewVector(item.Embedding))
    return err
}

type SearchResult struct {
    ID       int32
    Name     string
    Distance float64
}

func searchSimilar(ctx context.Context, pool *pgxpool.Pool, embedding []float32, limit int) ([]SearchResult, error) {
    rows, err := pool.Query(ctx, `
        SELECT id, name, 1 - (embedding <=> $1) as similarity
        FROM items
        ORDER BY embedding <=> $1
        LIMIT $2
    `, pgvector.NewVector(embedding), limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []SearchResult
    for rows.Next() {
        var r SearchResult
        if err := rows.Scan(&r.ID, &r.Name, &r.Distance); err != nil {
            return nil, err
        }
        results = append(results, r)
    }

    return results, rows.Err()
}

// Distance operators in pgvector:
// <-> - Euclidean distance (L2)
// <#> - Negative inner product
// <=> - Cosine distance (1 - cosine similarity)
```

---

## 3. Vector Search Algorithms

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Vector Search Algorithms                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Exact Search (KNN)                                                       │
│     - Calculate distance to every vector                                     │
│     - O(N) complexity                                                        │
│     - 100% accuracy                                                          │
│     - Good for: Small datasets (<10K vectors)                                │
│                                                                              │
│  2. IVF (Inverted File Index)                                                │
│     - Cluster vectors into Voronoi cells                                     │
│     - Query only nearest cells                                               │
│     - Fast but lower recall                                                  │
│     - Good for: Medium datasets, static data                                 │
│                                                                              │
│  3. HNSW (Hierarchical Navigable Small World)                                │
│     - Multi-layer graph structure                                            │
│     - Approximate nearest neighbor search                                    │
│     - Excellent speed/recall tradeoff                                        │
│     - Good for: Dynamic datasets, high performance                           │
│                                                                              │
│  4. PQ (Product Quantization)                                                │
│     - Compress vectors for memory efficiency                                 │
│     - Fast distance computation on compressed data                           │
│     - Good for: Billions of vectors, memory constrained                      │
│                                                                              │
│  Performance Comparison (1M vectors, 768 dim):                               │
│  ┌─────────────┬──────────────┬──────────────┬──────────────┐              │
│  │ Algorithm   │ Query Time   │ Recall@10    │ Memory       │              │
│  ├─────────────┼──────────────┼──────────────┼──────────────┤              │
│  │ Exact       │ 100ms        │ 100%         │ 3GB          │              │
│  │ IVF-Flat    │ 5ms          │ 95%          │ 3GB          │              │
│  │ HNSW        │ 1ms          │ 99%          │ 4GB          │              │
│  │ IVF-PQ      │ 2ms          │ 90%          │ 500MB        │              │
│  └─────────────┴──────────────┴──────────────┴──────────────┘              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Use Cases

```
Vector Database Use Cases:

1. Semantic Search
   - Natural language queries
   - Find documents by meaning, not keywords
   - Example: "Find documents about climate change"

2. Recommendation Systems
   - Item-to-item similarity
   - User preference matching
   - Content-based filtering

3. Image Search
   - Visual similarity search
   - Reverse image search
   - Product recommendation by image

4. RAG (Retrieval Augmented Generation)
   - Store document embeddings
   - Retrieve relevant context for LLM
   - Reduce hallucinations

5. Anomaly Detection
   - Find outliers in vector space
   - Fraud detection
   - Quality control
```

---

## 5. Checklist

```
Vector Database Checklist:
□ Embedding model chosen
□ Dimensionality determined
□ Distance metric selected (cosine/euclidean/dot)
□ Index type appropriate for dataset size
□ Batch insertion for initial load
□ Query parameters tuned (ef_search, nprobe)
□ Monitoring for query latency
□ Backup strategy for embeddings
```
