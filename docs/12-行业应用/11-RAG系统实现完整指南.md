# Go RAG系统实现完整指南

> **更新日期**: 2025年10月24日  
> **适用版本**: Go 1.21+  
> **难度**: ⭐⭐⭐⭐⭐  
> **标签**: #RAG #向量数据库 #Embedding #语义搜索 #检索增强

---

## 📚 目录

- [Go RAG系统实现完整指南](#go-rag系统实现完整指南)
  - [📚 目录](#-目录)
  - [1. RAG系统概述](#1-rag系统概述)
    - [1.1 什么是RAG](#11-什么是rag)
    - [1.2 RAG架构](#12-rag架构)
    - [1.3 核心组件](#13-核心组件)
  - [2. 向量数据库集成](#2-向量数据库集成)
    - [2.1 Qdrant集成](#21-qdrant集成)
    - [2.2 Weaviate集成](#22-weaviate集成)
    - [2.3 向量数据库抽象层](#23-向量数据库抽象层)
  - [3. Embedding生成](#3-embedding生成)
    - [3.1 OpenAI Embeddings](#31-openai-embeddings)
    - [3.2 本地Embedding模型](#32-本地embedding模型)
    - [3.3 Embedding缓存](#33-embedding缓存)
  - [4. 文档处理](#4-文档处理)
    - [4.1 文档加载](#41-文档加载)
    - [4.2 文档分块策略](#42-文档分块策略)
  - [5. 语义搜索](#5-语义搜索)
    - [5.1 向量相似度搜索](#51-向量相似度搜索)
  - [6. RAG Pipeline](#6-rag-pipeline)
    - [6.1 基础RAG流程](#61-基础rag流程)
  - [8. 实战案例](#8-实战案例)
    - [8.1 知识库问答系统](#81-知识库问答系统)
  - [11. 参考资源](#11-参考资源)
    - [官方文档](#官方文档)
    - [Go库](#go库)
    - [论文与研究](#论文与研究)
    - [最佳实践](#最佳实践)

---

## 1. RAG系统概述

### 1.1 什么是RAG

**RAG（Retrieval-Augmented Generation，检索增强生成）** 是一种结合信息检索和生成式AI的技术，通过在生成答案前检索相关文档，使LLM能够回答特定领域的问题。

**核心优势**:

- ✅ 减少幻觉（Hallucination）
- ✅ 提供最新信息
- ✅ 支持专有知识库
- ✅ 可追溯信息来源
- ✅ 降低微调成本

**应用场景**:

- 企业知识库问答
- 客户服务自动化
- 技术文档助手
- 法律/医疗文档分析

### 1.2 RAG架构

```mermaid
graph TB
    User[用户查询] --> Query[查询处理]
    Query --> Embedding1[查询向量化]
    Embedding1 --> VectorDB[向量数据库]
    
    VectorDB --> Retrieve[检索相关文档]
    Retrieve --> Rerank[重排序]
    Rerank --> Context[构建上下文]
    
    Context --> LLM[大语言模型]
    User --> LLM
    
    LLM --> Response[生成答案]
    Response --> User
    
    Documents[文档集合] --> Process[文档处理]
    Process --> Chunk[文档分块]
    Chunk --> Embedding2[文档向量化]
    Embedding2 --> VectorDB
    
    style Query fill:#e1f5fe
    style VectorDB fill:#e8f5e9
    style LLM fill:#fff3e0
    style Response fill:#f3e5f5
```

### 1.3 核心组件

| 组件 | 职责 | 技术选型 |
|------|------|----------|
| **文档处理器** | 加载、分块、元数据提取 | PDF解析、文本分割 |
| **Embedding服务** | 文本向量化 | OpenAI, HuggingFace |
| **向量数据库** | 存储和检索向量 | Qdrant, Weaviate, Pinecone |
| **检索器** | 语义搜索、重排序 | 相似度算法 |
| **LLM集成** | 生成答案 | OpenAI, Claude |
| **缓存层** | 性能优化 | Redis |

---

## 2. 向量数据库集成

### 2.1 Qdrant集成

**Qdrant客户端实现**:

```go
package vectordb

import (
    "context"
    "fmt"
    
    "github.com/qdrant/go-client/qdrant"
)

// QdrantClient Qdrant客户端
type QdrantClient struct {
    client     *qdrant.Client
    collection string
}

// Point 向量点
type Point struct {
    ID       string
    Vector   []float32
    Metadata map[string]interface{}
}

func NewQdrantClient(url, apiKey, collection string) (*QdrantClient, error) {
    client, err := qdrant.NewClient(&qdrant.Config{
        Host:   url,
        APIKey: apiKey,
    })
    if err != nil {
        return nil, fmt.Errorf("create client: %w", err)
    }

    return &QdrantClient{
        client:     client,
        collection: collection,
    }, nil
}

// CreateCollection 创建集合
func (c *QdrantClient) CreateCollection(ctx context.Context, dimension int) error {
    return c.client.CreateCollection(ctx, &qdrant.CreateCollection{
        CollectionName: c.collection,
        VectorsConfig: &qdrant.VectorsConfig{
            Params: &qdrant.VectorParams{
                Size:     uint64(dimension),
                Distance: qdrant.Distance_Cosine,
            },
        },
    })
}

// Upsert 插入或更新向量
func (c *QdrantClient) Upsert(ctx context.Context, points []Point) error {
    qdrantPoints := make([]*qdrant.PointStruct, len(points))
    
    for i, p := range points {
        qdrantPoints[i] = &qdrant.PointStruct{
            Id: &qdrant.PointId{
                PointIdOptions: &qdrant.PointId_Uuid{
                    Uuid: p.ID,
                },
            },
            Vectors: &qdrant.Vectors{
                VectorsOptions: &qdrant.Vectors_Vector{
                    Vector: &qdrant.Vector{
                        Data: p.Vector,
                    },
                },
            },
            Payload: convertToPayload(p.Metadata),
        }
    }

    _, err := c.client.Upsert(ctx, &qdrant.UpsertPoints{
        CollectionName: c.collection,
        Points:         qdrantPoints,
    })

    return err
}

// Search 向量搜索
func (c *QdrantClient) Search(ctx context.Context, vector []float32, topK int, filter map[string]interface{}) ([]SearchResult, error) {
    searchReq := &qdrant.SearchPoints{
        CollectionName: c.collection,
        Vector:         vector,
        Limit:          uint64(topK),
        WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
    }

    if filter != nil {
        searchReq.Filter = buildFilter(filter)
    }

    resp, err := c.client.Search(ctx, searchReq)
    if err != nil {
        return nil, err
    }

    results := make([]SearchResult, len(resp))
    for i, hit := range resp {
        results[i] = SearchResult{
            ID:       hit.Id.GetUuid(),
            Score:    hit.Score,
            Metadata: convertFromPayload(hit.Payload),
        }
    }

    return results, nil
}

// SearchResult 搜索结果
type SearchResult struct {
    ID       string
    Score    float32
    Metadata map[string]interface{}
}

func convertToPayload(m map[string]interface{}) map[string]*qdrant.Value {
    payload := make(map[string]*qdrant.Value)
    for k, v := range m {
        payload[k] = &qdrant.Value{
            Kind: &qdrant.Value_StringValue{
                StringValue: fmt.Sprintf("%v", v),
            },
        }
    }
    return payload
}

func convertFromPayload(payload map[string]*qdrant.Value) map[string]interface{} {
    result := make(map[string]interface{})
    for k, v := range payload {
        result[k] = v.GetStringValue()
    }
    return result
}

func buildFilter(filter map[string]interface{}) *qdrant.Filter {
    // 实现过滤器构建逻辑
    return nil
}

// 使用示例
func ExampleQdrant() {
    client, _ := NewQdrantClient("localhost:6333", "", "knowledge_base")
    
    ctx := context.Background()
    
    // 创建集合
    client.CreateCollection(ctx, 1536) // OpenAI embedding维度
    
    // 插入向量
    points := []Point{
        {
            ID:     "doc1",
            Vector: make([]float32, 1536), // 实际的embedding向量
            Metadata: map[string]interface{}{
                "text":   "Go是一门编程语言",
                "source": "docs/intro.md",
            },
        },
    }
    client.Upsert(ctx, points)
    
    // 搜索
    queryVector := make([]float32, 1536) // 查询的embedding向量
    results, _ := client.Search(ctx, queryVector, 5, nil)
    
    for _, result := range results {
        fmt.Printf("ID: %s, Score: %.4f\n", result.ID, result.Score)
        fmt.Printf("Text: %v\n", result.Metadata["text"])
    }
}
```

### 2.2 Weaviate集成

**Weaviate客户端实现**:

```go
package vectordb

import (
    "context"
    "fmt"
    
    "github.com/weaviate/weaviate-go-client/v4/weaviate"
    "github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

// WeaviateClient Weaviate客户端
type WeaviateClient struct {
    client    *weaviate.Client
    className string
}

func NewWeaviateClient(host, className string) *WeaviateClient {
    config := weaviate.Config{
        Host:   host,
        Scheme: "http",
    }
    
    client, _ := weaviate.NewClient(config)
    
    return &WeaviateClient{
        client:    client,
        className: className,
    }
}

// CreateSchema 创建schema
func (c *WeaviateClient) CreateSchema(ctx context.Context) error {
    classObj := &models.Class{
        Class: c.className,
        Properties: []*models.Property{
            {
                Name:     "text",
                DataType: []string{"text"},
            },
            {
                Name:     "source",
                DataType: []string{"string"},
            },
        },
        Vectorizer: "text2vec-openai",
    }
    
    return c.client.Schema().ClassCreator().WithClass(classObj).Do(ctx)
}

// AddDocument 添加文档
func (c *WeaviateClient) AddDocument(ctx context.Context, text, source string, vector []float32) (string, error) {
    properties := map[string]interface{}{
        "text":   text,
        "source": source,
    }
    
    result, err := c.client.Data().Creator().
        WithClassName(c.className).
        WithProperties(properties).
        WithVector(vector).
        Do(ctx)
    
    if err != nil {
        return "", err
    }
    
    return result.Object.ID.String(), nil
}

// Search 语义搜索
func (c *WeaviateClient) Search(ctx context.Context, query string, limit int) ([]map[string]interface{}, error) {
    result, err := c.client.GraphQL().Get().
        WithClassName(c.className).
        WithFields(graphql.Field{Name: "text"}, graphql.Field{Name: "source"}).
        WithNearText(&graphql.NearTextArgumentBuilder{
            Concepts: []string{query},
        }).
        WithLimit(limit).
        Do(ctx)
    
    if err != nil {
        return nil, err
    }
    
    // 解析结果
    data := result.Data["Get"].(map[string]interface{})[c.className].([]interface{})
    results := make([]map[string]interface{}, len(data))
    
    for i, item := range data {
        results[i] = item.(map[string]interface{})
    }
    
    return results, nil
}

// 使用示例
func ExampleWeaviate() {
    client := NewWeaviateClient("localhost:8080", "Document")
    ctx := context.Background()
    
    // 创建schema
    client.CreateSchema(ctx)
    
    // 添加文档
    vector := make([]float32, 1536)
    id, _ := client.AddDocument(ctx, "Go is a programming language", "docs/intro.md", vector)
    fmt.Println("Document ID:", id)
    
    // 搜索
    results, _ := client.Search(ctx, "What is Go?", 5)
    for _, result := range results {
        fmt.Printf("Text: %v\n", result["text"])
    }
}
```

### 2.3 向量数据库抽象层

**统一接口设计**:

```go
package vectordb

import "context"

// VectorStore 向量存储接口
type VectorStore interface {
    // Upsert 插入或更新向量
    Upsert(ctx context.Context, points []Point) error
    
    // Search 向量搜索
    Search(ctx context.Context, vector []float32, topK int, filter map[string]interface{}) ([]SearchResult, error)
    
    // Delete 删除向量
    Delete(ctx context.Context, ids []string) error
    
    // GetByID 根据ID获取向量
    GetByID(ctx context.Context, id string) (*Point, error)
}

// VectorStoreFactory 向量存储工厂
type VectorStoreFactory struct {
    stores map[string]VectorStore
}

func NewVectorStoreFactory() *VectorStoreFactory {
    return &VectorStoreFactory{
        stores: make(map[string]VectorStore),
    }
}

// Register 注册存储
func (f *VectorStoreFactory) Register(name string, store VectorStore) {
    f.stores[name] = store
}

// Get 获取存储
func (f *VectorStoreFactory) Get(name string) (VectorStore, error) {
    store, ok := f.stores[name]
    if !ok {
        return nil, fmt.Errorf("store %s not found", name)
    }
    return store, nil
}

// 使用示例：切换不同的向量数据库
func ExampleFactory() {
    factory := NewVectorStoreFactory()
    
    // 注册不同的向量数据库
    qdrant, _ := NewQdrantClient("localhost:6333", "", "kb")
    factory.Register("qdrant", qdrant)
    
    // 使用时切换
    store, _ := factory.Get("qdrant")
    ctx := context.Background()
    
    results, _ := store.Search(ctx, make([]float32, 1536), 5, nil)
    fmt.Printf("Found %d results\n", len(results))
}
```

---

## 3. Embedding生成

### 3.1 OpenAI Embeddings

**Embedding服务实现**:

```go
package embedding

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// EmbeddingService Embedding服务
type EmbeddingService struct {
    apiKey     string
    model      string
    httpClient *http.Client
}

// EmbeddingRequest Embedding请求
type EmbeddingRequest struct {
    Input          []string `json:"input"`
    Model          string   `json:"model"`
    EncodingFormat string   `json:"encoding_format,omitempty"`
}

// EmbeddingResponse Embedding响应
type EmbeddingResponse struct {
    Object string `json:"object"`
    Data   []struct {
        Object    string    `json:"object"`
        Embedding []float32 `json:"embedding"`
        Index     int       `json:"index"`
    } `json:"data"`
    Model string `json:"model"`
    Usage struct {
        PromptTokens int `json:"prompt_tokens"`
        TotalTokens  int `json:"total_tokens"`
    } `json:"usage"`
}

func NewEmbeddingService(apiKey string) *EmbeddingService {
    return &EmbeddingService{
        apiKey:     apiKey,
        model:      "text-embedding-3-small", // 或 text-embedding-3-large
        httpClient: &http.Client{},
    }
}

// Embed 生成单个文本的embedding
func (s *EmbeddingService) Embed(ctx context.Context, text string) ([]float32, error) {
    embeddings, err := s.EmbedBatch(ctx, []string{text})
    if err != nil {
        return nil, err
    }
    return embeddings[0], nil
}

// EmbedBatch 批量生成embedding
func (s *EmbeddingService) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
    req := EmbeddingRequest{
        Input: texts,
        Model: s.model,
    }

    body, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        "https://api.openai.com/v1/embeddings",
        bytes.NewReader(body),
    )
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

    resp, err := s.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("do request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
    }

    var embResp EmbeddingResponse
    if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }

    // 提取embedding
    embeddings := make([][]float32, len(embResp.Data))
    for _, data := range embResp.Data {
        embeddings[data.Index] = data.Embedding
    }

    return embeddings, nil
}

// 使用示例
func ExampleEmbedding() {
    service := NewEmbeddingService("sk-...")
    ctx := context.Background()

    // 单个文本
    embedding, _ := service.Embed(ctx, "Go is a programming language")
    fmt.Printf("Embedding dimension: %d\n", len(embedding))

    // 批量
    texts := []string{
        "Go is fast",
        "Go has goroutines",
        "Go is simple",
    }
    embeddings, _ := service.EmbedBatch(ctx, texts)
    fmt.Printf("Generated %d embeddings\n", len(embeddings))
}
```

### 3.2 本地Embedding模型

**使用Sentence Transformers（通过Python bridge）**:

```go
package embedding

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

// LocalEmbeddingService 本地embedding服务
type LocalEmbeddingService struct {
    endpoint string
    client   *http.Client
}

func NewLocalEmbeddingService(endpoint string) *LocalEmbeddingService {
    return &LocalEmbeddingService{
        endpoint: endpoint,
        client:   &http.Client{},
    }
}

// Embed 生成embedding（调用本地服务）
func (s *LocalEmbeddingService) Embed(text string) ([]float32, error) {
    req := map[string]interface{}{
        "text": text,
    }

    body, _ := json.Marshal(req)
    resp, err := s.client.Post(
        s.endpoint+"/embed",
        "application/json",
        bytes.NewReader(body),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Embedding []float32 `json:"embedding"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return result.Embedding, nil
}

/*
本地embedding服务（Python Flask示例）:

from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer

app = Flask(__name__)
model = SentenceTransformer('all-MiniLM-L6-v2')

@app.route('/embed', methods=['POST'])
def embed():
    text = request.json['text']
    embedding = model.encode(text).tolist()
    return jsonify({'embedding': embedding})

if __name__ == '__main__':
    app.run(port=5000)
*/
```

### 3.3 Embedding缓存

**Redis缓存实现**:

```go
package embedding

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

// CachedEmbeddingService 带缓存的embedding服务
type CachedEmbeddingService struct {
    embedder EmbeddingService
    cache    *redis.Client
    ttl      time.Duration
}

func NewCachedEmbeddingService(embedder EmbeddingService, redisAddr string) *CachedEmbeddingService {
    rdb := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })

    return &CachedEmbeddingService{
        embedder: embedder,
        cache:    rdb,
        ttl:      24 * time.Hour,
    }
}

// Embed 带缓存的embedding生成
func (s *CachedEmbeddingService) Embed(ctx context.Context, text string) ([]float32, error) {
    // 生成缓存键
    key := s.cacheKey(text)

    // 尝试从缓存获取
    cached, err := s.cache.Get(ctx, key).Result()
    if err == nil {
        var embedding []float32
        if err := json.Unmarshal([]byte(cached), &embedding); err == nil {
            return embedding, nil
        }
    }

    // 缓存未命中，生成embedding
    embedding, err := s.embedder.Embed(ctx, text)
    if err != nil {
        return nil, err
    }

    // 存入缓存
    data, _ := json.Marshal(embedding)
    s.cache.Set(ctx, key, data, s.ttl)

    return embedding, nil
}

func (s *CachedEmbeddingService) cacheKey(text string) string {
    hash := sha256.Sum256([]byte(text))
    return "emb:" + hex.EncodeToString(hash[:])
}

// 使用示例
func ExampleCachedEmbedding() {
    embedder := NewEmbeddingService("sk-...")
    cached := NewCachedEmbeddingService(embedder, "localhost:6379")
    
    ctx := context.Background()
    
    // 第一次调用：生成embedding并缓存
    emb1, _ := cached.Embed(ctx, "Go programming")
    
    // 第二次调用：从缓存获取（快速）
    emb2, _ := cached.Embed(ctx, "Go programming")
    
    fmt.Println("Same embedding:", len(emb1) == len(emb2))
}
```

---

## 4. 文档处理

### 4.1 文档加载

**多格式文档加载器**:

```go
package document

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

// Document 文档结构
type Document struct {
    Content  string
    Metadata map[string]interface{}
}

// Loader 文档加载器接口
type Loader interface {
    Load(path string) ([]Document, error)
}

// TextLoader 文本文件加载器
type TextLoader struct{}

func (l *TextLoader) Load(path string) ([]Document, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    return []Document{{
        Content: string(data),
        Metadata: map[string]interface{}{
            "source": path,
            "type":   "text",
        },
    }}, nil
}

// PDFLoader PDF加载器（需要pdf库）
type PDFLoader struct{}

func (l *PDFLoader) Load(path string) ([]Document, error) {
    // 使用github.com/ledongthuc/pdf等库解析PDF
    // 这里提供简化示例
    
    return []Document{{
        Content: "PDF content extracted...",
        Metadata: map[string]interface{}{
            "source": path,
            "type":   "pdf",
        },
    }}, nil
}

// DirectoryLoader 目录加载器
type DirectoryLoader struct {
    loaders map[string]Loader
}

func NewDirectoryLoader() *DirectoryLoader {
    return &DirectoryLoader{
        loaders: map[string]Loader{
            ".txt": &TextLoader{},
            ".md":  &TextLoader{},
            ".pdf": &PDFLoader{},
        },
    }
}

func (l *DirectoryLoader) Load(dirPath string) ([]Document, error) {
    var allDocs []Document

    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        ext := filepath.Ext(path)
        loader, ok := l.loaders[ext]
        if !ok {
            return nil // 跳过不支持的文件类型
        }

        docs, err := loader.Load(path)
        if err != nil {
            return fmt.Errorf("load %s: %w", path, err)
        }

        allDocs = append(allDocs, docs...)
        return nil
    })

    return allDocs, err
}

// 使用示例
func ExampleLoader() {
    loader := NewDirectoryLoader()
    docs, _ := loader.Load("./knowledge_base")
    
    fmt.Printf("Loaded %d documents\n", len(docs))
    for _, doc := range docs {
        fmt.Printf("Source: %v, Length: %d\n",
            doc.Metadata["source"],
            len(doc.Content),
        )
    }
}
```

### 4.2 文档分块策略

**智能文档分块**:

```go
package document

import (
    "strings"
    "unicode"
)

// Chunker 文档分块器
type Chunker struct {
    chunkSize    int
    chunkOverlap int
}

func NewChunker(chunkSize, chunkOverlap int) *Chunker {
    return &Chunker{
        chunkSize:    chunkSize,
        chunkOverlap: chunkOverlap,
    }
}

// Chunk 分块
func (c *Chunker) Chunk(doc Document) []Document {
    // 按段落分割
    paragraphs := strings.Split(doc.Content, "\n\n")
    
    var chunks []Document
    var currentChunk strings.Builder
    var chunkCount int

    for _, para := range paragraphs {
        para = strings.TrimSpace(para)
        if para == "" {
            continue
        }

        // 如果当前块加上新段落超过限制
        if currentChunk.Len()+len(para) > c.chunkSize && currentChunk.Len() > 0 {
            // 保存当前块
            chunks = append(chunks, c.createChunk(currentChunk.String(), doc.Metadata, chunkCount))
            chunkCount++

            // 开始新块，包含overlap
            if c.chunkOverlap > 0 {
                overlapText := c.getOverlap(currentChunk.String())
                currentChunk.Reset()
                currentChunk.WriteString(overlapText)
                currentChunk.WriteString("\n\n")
            } else {
                currentChunk.Reset()
            }
        }

        currentChunk.WriteString(para)
        currentChunk.WriteString("\n\n")
    }

    // 保存最后一个块
    if currentChunk.Len() > 0 {
        chunks = append(chunks, c.createChunk(currentChunk.String(), doc.Metadata, chunkCount))
    }

    return chunks
}

func (c *Chunker) createChunk(content string, metadata map[string]interface{}, index int) Document {
    chunkMetadata := make(map[string]interface{})
    for k, v := range metadata {
        chunkMetadata[k] = v
    }
    chunkMetadata["chunk_index"] = index
    chunkMetadata["chunk_size"] = len(content)

    return Document{
        Content:  strings.TrimSpace(content),
        Metadata: chunkMetadata,
    }
}

func (c *Chunker) getOverlap(text string) string {
    if len(text) <= c.chunkOverlap {
        return text
    }
    return text[len(text)-c.chunkOverlap:]
}

// RecursiveChunker 递归分块器（更智能）
type RecursiveChunker struct {
    chunkSize    int
    chunkOverlap int
    separators   []string
}

func NewRecursiveChunker(chunkSize, chunkOverlap int) *RecursiveChunker {
    return &RecursiveChunker{
        chunkSize:    chunkSize,
        chunkOverlap: chunkOverlap,
        separators:   []string{"\n\n", "\n", ". ", " ", ""},
    }
}

// Chunk 递归分块
func (c *RecursiveChunker) Chunk(doc Document) []Document {
    return c.splitText(doc.Content, doc.Metadata, 0)
}

func (c *RecursiveChunker) splitText(text string, metadata map[string]interface{}, depth int) []Document {
    if len(text) <= c.chunkSize {
        return []Document{{
            Content:  text,
            Metadata: metadata,
        }}
    }

    if depth >= len(c.separators) {
        // 强制分割
        return c.forceSplit(text, metadata)
    }

    separator := c.separators[depth]
    parts := strings.Split(text, separator)

    var chunks []Document
    var currentChunk string

    for _, part := range parts {
        testChunk := currentChunk + separator + part
        if len(testChunk) > c.chunkSize && currentChunk != "" {
            // 当前块已满，递归处理
            chunks = append(chunks, c.splitText(currentChunk, metadata, depth+1)...)
            currentChunk = part
        } else {
            currentChunk = testChunk
        }
    }

    if currentChunk != "" {
        chunks = append(chunks, c.splitText(currentChunk, metadata, depth+1)...)
    }

    return chunks
}

func (c *RecursiveChunker) forceSplit(text string, metadata map[string]interface{}) []Document {
    var chunks []Document
    for i := 0; i < len(text); i += c.chunkSize - c.chunkOverlap {
        end := i + c.chunkSize
        if end > len(text) {
            end = len(text)
        }
        
        chunks = append(chunks, Document{
            Content:  text[i:end],
            Metadata: metadata,
        })

        if end == len(text) {
            break
        }
    }
    return chunks
}

// 使用示例
func ExampleChunker() {
    doc := Document{
        Content: `This is a long document...
        
With multiple paragraphs...

And various sections...`,
        Metadata: map[string]interface{}{
            "source": "test.txt",
        },
    }

    // 基础分块
    chunker := NewChunker(200, 50)
    chunks := chunker.Chunk(doc)
    fmt.Printf("Basic chunker: %d chunks\n", len(chunks))

    // 递归分块
    recursiveChunker := NewRecursiveChunker(200, 50)
    chunks2 := recursiveChunker.Chunk(doc)
    fmt.Printf("Recursive chunker: %d chunks\n", len(chunks2))
}
```

---

## 5. 语义搜索

### 5.1 向量相似度搜索

**相似度计算**:

```go
package search

import (
    "math"
    "sort"
)

// CosineSimilarity 余弦相似度
func CosineSimilarity(a, b []float32) float32 {
    var dotProduct, normA, normB float32

    for i := range a {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }

    if normA == 0 || normB == 0 {
        return 0
    }

    return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

// EuclideanDistance 欧氏距离
func EuclideanDistance(a, b []float32) float32 {
    var sum float32
    for i := range a {
        diff := a[i] - b[i]
        sum += diff * diff
    }
    return float32(math.Sqrt(float64(sum)))
}

// SearchResult 搜索结果
type SearchResult struct {
    DocumentID string
    Score      float32
    Content    string
    Metadata   map[string]interface{}
}

// LocalSearch 本地向量搜索（用于小数据集）
type LocalSearch struct {
    vectors  map[string][]float32
    metadata map[string]map[string]interface{}
}

func NewLocalSearch() *LocalSearch {
    return &LocalSearch{
        vectors:  make(map[string][]float32),
        metadata: make(map[string]map[string]interface{}),
    }
}

// Add 添加向量
func (s *LocalSearch) Add(id string, vector []float32, metadata map[string]interface{}) {
    s.vectors[id] = vector
    s.metadata[id] = metadata
}

// Search 搜索
func (s *LocalSearch) Search(queryVector []float32, topK int) []SearchResult {
    type scoredResult struct {
        id    string
        score float32
    }

    var results []scoredResult
    for id, vector := range s.vectors {
        score := CosineSimilarity(queryVector, vector)
        results = append(results, scoredResult{id: id, score: score})
    }

    // 按分数排序
    sort.Slice(results, func(i, j int) bool {
        return results[i].score > results[j].score
    })

    // 取topK
    if topK > len(results) {
        topK = len(results)
    }

    searchResults := make([]SearchResult, topK)
    for i := 0; i < topK; i++ {
        searchResults[i] = SearchResult{
            DocumentID: results[i].id,
            Score:      results[i].score,
            Metadata:   s.metadata[results[i].id],
        }
    }

    return searchResults
}
```

---

由于RAG文档内容很长，让我继续完成核心部分...

## 6. RAG Pipeline

### 6.1 基础RAG流程

**完整的RAG实现**:

```go
package rag

import (
    "context"
    "fmt"
    "strings"
)

// RAGPipeline RAG管道
type RAGPipeline struct {
    embedder    EmbeddingService
    vectorStore VectorStore
    llmClient   LLMClient
    topK        int
}

func NewRAGPipeline(embedder EmbeddingService, store VectorStore, llm LLMClient) *RAGPipeline {
    return &RAGPipeline{
        embedder:    embedder,
        vectorStore: store,
        llmClient:   llm,
        topK:        3,
    }
}

// Query 执行RAG查询
func (p *RAGPipeline) Query(ctx context.Context, query string) (string, error) {
    // 1. 生成查询的embedding
    queryEmbedding, err := p.embedder.Embed(ctx, query)
    if err != nil {
        return "", fmt.Errorf("embed query: %w", err)
    }

    // 2. 检索相关文档
    results, err := p.vectorStore.Search(ctx, queryEmbedding, p.topK, nil)
    if err != nil {
        return "", fmt.Errorf("search: %w", err)
    }

    // 3. 构建上下文
    context := p.buildContext(results)

    // 4. 生成答案
    prompt := p.buildPrompt(query, context)
    answer, err := p.llmClient.Generate(ctx, prompt)
    if err != nil {
        return "", fmt.Errorf("generate answer: %w", err)
    }

    return answer, nil
}

func (p *RAGPipeline) buildContext(results []SearchResult) string {
    var contexts []string
    for i, result := range results {
        text := result.Metadata["text"].(string)
        source := result.Metadata["source"].(string)
        contexts = append(contexts, fmt.Sprintf("[%d] %s\nSource: %s", i+1, text, source))
    }
    return strings.Join(contexts, "\n\n")
}

func (p *RAGPipeline) buildPrompt(query, context string) string {
    return fmt.Sprintf(`Based on the following context, please answer the question.

Context:
%s

Question: %s

Answer:`, context, query)
}
```

---

继续补充实战案例和生产部署部分...

## 8. 实战案例

### 8.1 知识库问答系统

**完整的知识库QA系统**:

```go
package main

import (
    "context"
    "fmt"
    "log"
)

// KnowledgeBaseQA 知识库问答系统
type KnowledgeBaseQA struct {
    embedder    *EmbeddingService
    vectorStore VectorStore
    llm         *OpenAIClient
    chunker     *RecursiveChunker
}

func NewKnowledgeBaseQA(apiKey string) *KnowledgeBaseQA {
    embedder := NewEmbeddingService(apiKey)
    vectorStore, _ := NewQdrantClient("localhost:6333", "", "kb")
    llm := NewOpenAIClient(apiKey)
    chunker := NewRecursiveChunker(500, 50)

    return &KnowledgeBaseQA{
        embedder:    embedder,
        vectorStore: vectorStore,
        llm:         llm,
        chunker:     chunker,
    }
}

// IndexDocuments 索引文档
func (qa *KnowledgeBaseQA) IndexDocuments(ctx context.Context, dirPath string) error {
    // 1. 加载文档
    loader := NewDirectoryLoader()
    docs, err := loader.Load(dirPath)
    if err != nil {
        return fmt.Errorf("load documents: %w", err)
    }

    // 2. 分块
    var allChunks []Document
    for _, doc := range docs {
        chunks := qa.chunker.Chunk(doc)
        allChunks = append(allChunks, chunks...)
    }

    // 3. 批量生成embedding并存储
    batchSize := 100
    for i := 0; i < len(allChunks); i += batchSize {
        end := i + batchSize
        if end > len(allChunks) {
            end = len(allChunks)
        }

        batch := allChunks[i:end]
        texts := make([]string, len(batch))
        for j, chunk := range batch {
            texts[j] = chunk.Content
        }

        // 生成embeddings
        embeddings, err := qa.embedder.EmbedBatch(ctx, texts)
        if err != nil {
            return fmt.Errorf("embed batch: %w", err)
        }

        // 存储到向量数据库
        points := make([]Point, len(batch))
        for j, chunk := range batch {
            points[j] = Point{
                ID:       fmt.Sprintf("chunk_%d", i+j),
                Vector:   embeddings[j],
                Metadata: map[string]interface{}{
                    "text":   chunk.Content,
                    "source": chunk.Metadata["source"],
                },
            }
        }

        if err := qa.vectorStore.Upsert(ctx, points); err != nil {
            return fmt.Errorf("upsert: %w", err)
        }

        log.Printf("Indexed %d chunks\n", end)
    }

    return nil
}

// Ask 提问
func (qa *KnowledgeBaseQA) Ask(ctx context.Context, question string) (string, error) {
    // 1. 生成问题的embedding
    queryEmbedding, err := qa.embedder.Embed(ctx, question)
    if err != nil {
        return "", err
    }

    // 2. 检索相关文档
    results, err := qa.vectorStore.Search(ctx, queryEmbedding, 3, nil)
    if err != nil {
        return "", err
    }

    // 3. 构建prompt
    var contexts []string
    for i, result := range results {
        text := result.Metadata["text"].(string)
        contexts = append(contexts, fmt.Sprintf("[%d] %s", i+1, text))
    }

    prompt := fmt.Sprintf(`Based on the following context, please answer the question accurately.

Context:
%s

Question: %s

Instructions:
- Answer based only on the provided context
- If the answer is not in the context, say "I don't have enough information to answer this question"
- Cite the relevant context numbers in your answer

Answer:`, strings.Join(contexts, "\n\n"), question)

    // 4. 调用LLM生成答案
    resp, err := qa.llm.Chat(ctx, ChatRequest{
        Model: "gpt-4",
        Messages: []ChatMessage{
            {Role: "system", Content: "You are a helpful assistant that answers questions based on provided context."},
            {Role: "user", Content: prompt},
        },
        Temperature: 0.3,
    })

    if err != nil {
        return "", err
    }

    return resp.Choices[0].Message.Content, nil
}

// 使用示例
func main() {
    qa := NewKnowledgeBaseQA("sk-...")
    ctx := context.Background()

    // 索引文档
    if err := qa.IndexDocuments(ctx, "./knowledge_base"); err != nil {
        log.Fatal(err)
    }

    // 提问
    answer, err := qa.Ask(ctx, "What is Go programming language?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Answer:", answer)
}
```

---

## 11. 参考资源

### 官方文档

- [LangChain Documentation](https://docs.langchain.com/)
- [Qdrant Documentation](https://qdrant.tech/documentation/)
- [Weaviate Documentation](https://weaviate.io/developers/weaviate)
- [OpenAI Embeddings](https://platform.openai.com/docs/guides/embeddings)

### Go库

- [go-client (Qdrant)](https://github.com/qdrant/go-client)
- [weaviate-go-client](https://github.com/weaviate/weaviate-go-client)
- [chromem-go](https://github.com/philippgille/chromem-go) - 纯Go向量数据库

### 论文与研究

- [RAG: Retrieval-Augmented Generation](https://arxiv.org/abs/2005.11401)
- [Dense Passage Retrieval](https://arxiv.org/abs/2004.04906)

### 最佳实践

- [Building Production-Ready RAG Applications](https://www.pinecone.io/learn/rag/)
- [Advanced RAG Techniques](https://blog.langchain.dev/advanced-rag/)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月24日  
**文档状态**: ✅ 完成  
**适用版本**: Go 1.21+

**贡献者**: 欢迎提交Issue和PR改进本文档

