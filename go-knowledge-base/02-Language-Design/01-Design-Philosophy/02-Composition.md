# 组合优于继承 (Composition)

> **分类**: 语言设计

---

## 核心理念

**组合** (has-a) 优于 **继承** (is-a)

---

## Go 的组合方式

### 1. 结构体嵌入

```go
type Reader struct { }
func (r Reader) Read() { }

type Writer struct { }
func (w Writer) Write() { }

// 组合
type ReadWriter struct {
    Reader      // 嵌入
    Writer      // 嵌入
}

// 自动获得 Read() 和 Write() 方法
```

### 2. 接口组合

```go
type Reader interface {
    Read()
}

type Writer interface {
    Write()
}

// 接口组合
type ReadWriter interface {
    Reader
    Writer
}
```

---

## 对比继承

| 特性 | 继承 | 组合 |
|------|------|------|
| 耦合度 | 高 | 低 |
| 灵活性 | 低 | 高 |
| 复用性 | 有限 | 高 |
| 理解难度 | 高（隐含关系） | 低（显式关系） |

---

## 最佳实践

```go
// 推荐: 小接口组合
type Reader interface { Read() }
type Writer interface { Write() }
type Closer interface { Close() }

type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

---

## 设计影响

Go 标准库大量使用组合：

- `io.Reader` + `io.Writer` = `io.ReadWriter`
- 可测试性高
- 可扩展性强
