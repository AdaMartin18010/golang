# 核心API参考

**难度**: 入门 | **预计阅读**: 10分钟

---

## 📋 目录

- [1. 📖 标准库核心包](#1--标准库核心包)
- [2. 📚 相关资源](#2--相关资源)

---

## 1. 📖 标准库核心包

### net/http

```go
// HTTP服务器
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)

// HTTP客户端
resp, err := http.Get("https://example.com")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)

// 自定义请求
req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
req.Header.Set("Content-Type", "application/json")
client := &http.Client{}
resp, _ := client.Do(req)
```

---

### encoding/json

```go
// 序列化
data, _ := json.Marshal(struct)

// 反序列化
var result MyStruct
json.Unmarshal(data, &result)

// 流式编解码
json.NewEncoder(w).Encode(data)
json.NewDecoder(r.Body).Decode(&data)
```

---

### context

```go
// 超时控制
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 值传递
ctx = context.WithValue(ctx, "key", "value")
value := ctx.Value("key")

// 取消信号
ctx, cancel := context.WithCancel(context.Background())
go func() {
    <-ctx.Done()
    // 清理...
}()
cancel()
```

---

### sync

```go
// 互斥锁
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

// 读写锁
var rwmu sync.RWMutex
rwmu.RLock()
defer rwmu.RUnlock()

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // 工作...
}()
wg.Wait()
```

---

## 📚 相关资源

- [Go Standard Library](https://pkg.go.dev/std)

**下一步**: [02-标准库API](./02-标准库API.md)

---

**最后更新**: 2025-10-28

