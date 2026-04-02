# 接口 (Interfaces)

> **分类**: 语言设计

---

## 定义

接口定义行为契约：

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

---

## 隐式实现

```go
type File struct{}

func (f File) Read(p []byte) (n int, err error) {
    // 实现
    return 0, nil
}

// File 自动实现 Reader，无需声明
var r Reader = File{}
```

---

## 接口值

```go
// 接口值 = (类型, 值)
var r Reader = File{fd: 1}

// 运行时: (File, File{fd: 1})
```

### 结构

```
接口值 {
    tab: *itab    // 类型和方法表
    data: unsafe.Pointer  // 实际数据
}
```

---

## 空接口

```go
// 空接口接受任何类型
var x interface{} = 42
var y interface{} = "hello"
```

---

## 接口组合

```go
type ReadWriter interface {
    Reader
    Writer
}
```

---

## 最佳实践

### 小接口原则

```go
// 好: 小接口
type Reader interface { Read() }
type Writer interface { Write() }
type Closer interface { Close() }

// 不好: 大接口
type BigInterface interface {
    Read()
    Write()
    Close()
    Seek()
    // ...
}
```

---

## 鸭子类型

> "If it looks like a duck and quacks like a duck, it's a duck."

Go 接口实现结构子类型（鸭子类型）：

- 不声明实现关系
- 由编译器检查方法集
- 运行时动态分派
