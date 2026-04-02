# 竞态检测 (Race Detection)

> **分类**: 工程与云原生

---

## 启用 Race Detector

```bash
go run -race main.go
go test -race ./...
```

---

## 常见数据竞争

### 1. 读写竞争

```go
// ❌ 错误
var counter int

go func() {
    counter++  // 写
}()

go func() {
    fmt.Println(counter)  // 读 - 数据竞争!
}()
```

### 2. 切片竞争

```go
// ❌ 错误
s := make([]int, 10)

go func() {
    s[0] = 1
}()

go func() {
    s[0] = 2  // 数据竞争!
}()
```

### 3. Map 竞争

```go
// ❌ 错误
m := make(map[string]int)

go func() {
    m["key"] = 1
}()

go func() {
    _ = m["key"]  // 可能 panic!
}()
```

---

## 解决方案

### Mutex

```go
var (
    mu      sync.Mutex
    counter int
)

go func() {
    mu.Lock()
    counter++
    mu.Unlock()
}()

go func() {
    mu.Lock()
    fmt.Println(counter)
    mu.Unlock()
}()
```

### Atomic

```go
import "sync/atomic"

var counter int64

go func() {
    atomic.AddInt64(&counter, 1)
}()

go func() {
    val := atomic.LoadInt64(&counter)
    fmt.Println(val)
}()
```

### Channel

```go
counter := make(chan int, 1)
counter <- 0

go func() {
    n := <-counter
    n++
    counter <- n
}()
```

---

## CI 集成

```yaml
# GitHub Actions
- name: Race Test
  run: go test -race ./...
```

---

## 性能影响

Race detector 会使程序慢 10-20 倍，仅用于测试环境。
