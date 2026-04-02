# net/http 详解

> **分类**: 开源技术堆栈

---

## HTTP 服务器

### 基本用法

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello %s", r.URL.Path)
})

http.ListenAndServe(":8080", nil)
```

### Handler 接口

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

---

## HTTP 客户端

```go
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

body, _ := io.ReadAll(resp.Body)
```

### 带超时

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
resp, err := client.Get(url)
```

---

## 中间件模式

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL, time.Since(start))
    })
}
```
