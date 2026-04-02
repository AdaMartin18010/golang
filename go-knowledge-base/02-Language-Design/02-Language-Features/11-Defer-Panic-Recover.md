# Defer, Panic, Recover

> **分类**: 语言设计

---

## Defer

### 基本用法

```go
func readFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()  // 函数返回时执行

    // 处理文件
    return nil
}
```

### 多个 Defer (LIFO)

```go
func multipleDefer() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    // 输出: 3 2 1
}
```

### Defer 参数求值

```go
func deferArgs() {
    i := 0
    defer fmt.Println(i)  // 输出 0，defer 注册时求值
    i++
}
```

---

## Panic

### 触发 Panic

```go
func mayPanic() {
    panic("something went wrong")
}

func caller() {
    mayPanic()  // 传播 panic
    fmt.Println("不会执行到这里")
}
```

### Panic 场景

```go
// 除零
var a = 1 / 0  // panic

// 空指针
type T struct{}
var p *T
p.Method()  // panic

// 数组越界
arr := []int{1, 2, 3}
_ = arr[10]  // panic
```

---

## Recover

### 捕获 Panic

```go
func safeCall() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered: %v\n", r)
        }
    }()

    mayPanic()
}
```

### HTTP Server 中的 Recover

```go
func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v\n%s", err, debug.Stack())
                c.AbortWithStatusJSON(500, gin.H{
                    "error": "internal server error",
                })
            }
        }()
        c.Next()
    }
}
```

---

## 最佳实践

```go
// ✅ 用 defer 释放资源
func good() {
    f, _ := os.Open("file")
    defer f.Close()
}

// ✅ recover 只在 defer 中有效
func goodRecover() {
    defer func() {
        recover()  // 正确
    }()
}

// ❌ recover 不在 defer 中无效
func badRecover() {
    recover()  // 无效
}

// ✅ 不要滥用 panic
func badUsage() {
    if user == nil {
        panic("user is nil")  // ❌ 应该返回 error
    }
}

// ✅ 应该这样
func goodUsage() (*User, error) {
    if user == nil {
        return nil, errors.New("user is nil")
    }
    return user, nil
}
```
