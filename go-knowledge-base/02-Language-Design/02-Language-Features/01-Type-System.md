# 类型系统 (Type System)

> **分类**: 语言设计

---

## 核心特性

### 1. 静态类型

```go
var x int = 42      // 编译时确定类型
y := "hello"        // 类型推断
```

### 2. 结构子类型

```go
type Reader interface { Read() }

type File struct{}
func (f File) Read() {}  // 自动实现 Reader
```

### 3. 接口类型

```go
// 隐式实现
var r Reader = File{}
```

---

## 类型层级

```
interface{}  (空接口，所有类型的父类型)
    │
    ├── 基本类型: int, float64, string, bool
    │
    ├── 复合类型: struct, array, slice, map
    │
    ├── 函数类型: func
    │
    └── 接口类型: io.Reader, error
```

---

## 类型安全

Go 的类型系统保证：

- 编译时捕获类型错误
- 无隐式类型转换
- 无空指针（nil 检查）

---

## 类型推断

```go
// 显式
var x int = 42

// 推断
x := 42              // int
y := make([]int, 0)  // []int

// 泛型推断
func Max[T ~int](a, b T) T
Max(1, 2)            // T 推断为 int
```

---

## 类型断言与反射

```go
// 类型断言
r := getReader()     // interface{}
f, ok := r.(File)    // 断言为 File

// 类型开关
switch v := r.(type) {
case File:
    // v 是 File
case Socket:
    // v 是 Socket
}
```
