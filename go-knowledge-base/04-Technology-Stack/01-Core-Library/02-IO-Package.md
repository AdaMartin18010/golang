# io 包详解

> **分类**: 开源技术堆栈

---

## 核心接口

### Reader

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

### Writer

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

### ReadWriter

```go
type ReadWriter interface {
    Reader
    Writer
}
```

---

## 常用函数

```go
// 复制
io.Copy(dst, src)
io.CopyN(dst, src, 100)

// 读取全部
io.ReadAll(reader)

// 写入字符串
io.WriteString(writer, "hello")

// 多路读取
io.MultiReader(r1, r2, r3)
```

---

## 实现示例

```go
// 自定义 Reader
type MyReader struct {
    data []byte
    pos  int
}

func (r *MyReader) Read(p []byte) (n int, err error) {
    if r.pos >= len(r.data) {
        return 0, io.EOF
    }
    n = copy(p, r.data[r.pos:])
    r.pos += n
    return n, nil
}
```
