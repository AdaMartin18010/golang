# 向量数据库 (Vector Database)

> **分类**: 开源技术堆栈
> **标签**: #vector #embedding #ai #milvus #pinecone

---

## Milvus

```go
import "github.com/milvus-io/milvus-sdk-go/v2/client"

ctx := context.Background()

milvusClient, err := client.NewGrpcClient(ctx, "localhost:19530")
if err != nil {
    log.Fatal(err)
}
defer milvusClient.Close()
```

### 创建 Collection

```go
schema := &entity.Schema{
    CollectionName: "document_vectors",
    Fields: []*entity.Field{
        {
            Name:       "id",
            DataType:   entity.FieldTypeInt64,
            PrimaryKey: true,
            AutoID:     true,
        },
        {
            Name:     "embedding",
            DataType: entity.FieldTypeFloatVector,
            TypeParams: map[string]string{
                "dim": "1536",
            },
        },
        {
            Name:     "content",
            DataType: entity.FieldTypeVarChar,
            TypeParams: map[string]string{
                "max_length": "65535",
            },
        },
    },
}

err = milvusClient.CreateCollection(ctx, schema, 2)
```

### 向量搜索

```go
// 准备查询向量
queryVector := []float32{0.1, 0.2, 0.3, /* ... 1536 dims */}

sp, _ := entity.NewIndexHNSWSearchParam(100)

searchResult, err := milvusClient.Search(
    ctx,
    "document_vectors",
    nil,
    "",
    []string{"content"},
    []entity.Vector{entity.FloatVector(queryVector)},
    "embedding",
    entity.L2,
    10,
    sp,
)
```

---

## PGVector (PostgreSQL)

```go
// 创建扩展
_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

// 创建表
_, err = db.Exec(`
    CREATE TABLE documents (
        id SERIAL PRIMARY KEY,
        content TEXT,
        embedding vector(1536)
    )
`)

// 插入向量
_, err = db.Exec(
    "INSERT INTO documents (content, embedding) VALUES ($1, $2)",
    "content text",
    pgvector.NewVector(embedding),
)

// 相似度搜索
rows, err := db.Query(`
    SELECT content, embedding <=> $1 as distance
    FROM documents
    ORDER BY embedding <=> $1
    LIMIT 10
`, pgvector.NewVector(queryVec))
```

---

## 使用场景

| 场景 | 说明 |
|------|------|
| RAG | 检索增强生成 |
| 语义搜索 | 基于含义的文档搜索 |
| 推荐系统 | 相似内容推荐 |
| 去重 | 相似文档检测 |
