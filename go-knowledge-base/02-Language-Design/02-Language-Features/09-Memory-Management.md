# 内存管理

> **分类**: 语言设计

---

## 内存分配器

### 分级分配

```
Tiny:  < 16B
Small: < 32KB
Large: >= 32KB
```

---

## 栈与堆

```go
// 栈分配 (通常)
x := 42

// 堆分配 (逃逸)
func getX() *int {
    x := 42
    return &x  // 逃逸到堆
}
```

---

## 逃逸分析

```bash
go build -gcflags="-m" main.go
```
