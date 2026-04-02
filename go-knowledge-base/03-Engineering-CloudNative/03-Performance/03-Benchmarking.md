# 基准测试

> **分类**: 工程与云原生

---

## 基本用法

```go
func BenchmarkFib(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fib(20)
    }
}
```

---

## 内存分析

```go
func BenchmarkAlloc(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}
```

---

## 比较

```bash
go test -bench=. -benchmem
go test -bench=Benchmark -count=5
```
