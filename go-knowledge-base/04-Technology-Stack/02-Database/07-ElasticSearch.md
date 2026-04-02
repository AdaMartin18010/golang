# ElasticSearch

> **分类**: 开源技术堆栈
> **标签**: #elasticsearch #search #logging

---

## 客户端

```go
import "github.com/elastic/go-elasticsearch/v8"

es, err := elasticsearch.NewClient(elasticsearch.Config{
    Addresses: []string{"http://localhost:9200"},
    Username:  "elastic",
    Password:  "password",
})
if err != nil {
    log.Fatal(err)
}
```

---

## 索引文档

```go
import "bytes"
import "encoding/json"

doc := struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}{
    Title: "Go Tutorial",
    Body:  "Learn Go programming",
}

data, _ := json.Marshal(doc)

res, err := es.Index(
    "articles",
    bytes.NewReader(data),
    es.Index.WithDocumentID("1"),
)
if err != nil {
    log.Fatal(err)
}
defer res.Body.Close()
```

---

## 搜索

```go
query := map[string]interface{}{
    "query": map[string]interface{}{
        "match": map[string]interface{}{
            "title": "Go",
        },
    },
}

var buf bytes.Buffer
json.NewEncoder(&buf).Encode(query)

res, err := es.Search(
    es.Search.WithIndex("articles"),
    es.Search.WithBody(&buf),
)
if err != nil {
    log.Fatal(err)
}
defer res.Body.Close()

var result map[string]interface{}
json.NewDecoder(res.Body).Decode(&result)
```

---

## 聚合查询

```go
aggQuery := map[string]interface{}{
    "aggs": map[string]interface{}{
        "by_category": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "category.keyword",
            },
        },
    },
}
```
