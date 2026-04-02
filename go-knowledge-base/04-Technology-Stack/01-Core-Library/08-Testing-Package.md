# testing 包

> **分类**: 开源技术堆栈

---

## 测试函数

```go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

---

## 辅助函数

```go
func TestComplex(t *testing.T) {
    t.Helper()
    t.Parallel()
    t.Skip("skip this")
    t.Fatal("fatal error")
}
```

---

## 覆盖率

```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```
