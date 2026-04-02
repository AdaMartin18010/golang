# 性能优化 (Optimization)

> **分类**: 工程与云原生

---

## 内存优化

### 预分配

```go
// 好: 预分配
result := make([]int, 0, len(data))
for _, v := range data {
    result = append(result, v*2)
}

// 不好: 重复分配
var result []int
for _, v := range data {
    result = append(result, v*2)
}
```

### sync.Pool

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func process() {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    // 使用 buf
}
```

---

## 并发优化

### 限制 goroutine

```go
semaphore := make(chan struct{}, 10)

for _, item := range items {
    semaphore <- struct{}{}
    go func(i Item) {
        defer func() { <-semaphore }()
        process(i)
    }(item)
}
```

---

## 编译优化

### PGO (Go 1.20+)

```bash
# 1. 收集性能分析
go test -cpuprofile=cpu.pprof

# 2. 使用 PGO 构建
go build -pgo=cpu.pprof
```

**效果**: 通常提升 2-4% 性能
