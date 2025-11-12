# Go语法基础

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---
## 📋 目录

- [Go语法基础](#go语法基础)
  - [📚 文档列表](#文档列表)
  - [🚀 快速示例](#快速示例)
  - [📖 系统文档](#系统文档)

---

## 📚 文档列表

1. **[基本语法](./01-基本语法.md)** ⭐⭐⭐⭐⭐
   - 变量声明、常量、零值
   - 作用域、可见性

2. **[数据类型](./02-数据类型.md)** ⭐⭐⭐⭐⭐
   - 基本类型：int, float, string, bool
   - 复合类型：array, slice, map, struct
   - 引用类型：pointer, Channel

3. **[运算符](./03-运算符.md)** ⭐⭐⭐⭐
   - 算术、关系、逻辑、位运算

4. **[流程控制](./04-流程控制.md)** ⭐⭐⭐⭐⭐
   - if/else, switch, for
   - goto, break, continue

5. **[函数](./05-函数.md)** ⭐⭐⭐⭐⭐
   - 函数定义、参数、返回值
   - 闭包、递归、defer

6. **[错误处理](./06-错误处理.md)** ⭐⭐⭐⭐⭐
   - error接口、panic/recover

7. **[包管理](./07-包管理.md)** ⭐⭐⭐⭐
   - package、import、init

---

## 🚀 快速示例

### 变量声明

```go
var x int = 10
y := 20
const Pi = 3.14
```

### 流程控制

```go
if x > 10 {
    fmt.Println("x > 10")
}

for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

### 函数

```go
func add(a, b int) int {
    return a + b
}

func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

---

## 📖 系统文档
