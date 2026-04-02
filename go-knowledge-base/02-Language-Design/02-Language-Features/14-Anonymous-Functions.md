# 匿名函数与闭包 (Anonymous Functions & Closures)

> **分类**: 语言设计

---

## 匿名函数

```go
// 定义并立即调用
result := func(a, b int) int {
    return a + b
}(1, 2)  // result = 3

// 赋值给变量
add := func(a, b int) int {
    return a + b
}
result := add(3, 4)  // result = 7
```

---

## 闭包

```go
func makeMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor  // 捕获 factor
    }
}

double := makeMultiplier(2)
triple := makeMultiplier(3)

fmt.Println(double(5))  // 10
fmt.Println(triple(5))  // 15
```

---

## 闭包陷阱

### 循环变量问题 (Go < 1.22)

```go
// ❌ 错误
funcs := []func(){}
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)  // 都输出 3！
    })
}

// ✅ 正确or i := 0; i < 3; i++ {
    i := i  // 创建副本
    funcs = append(funcs, func() {
        fmt.Println(i)  // 输出 0, 1, 2
    })
}
```

**Go 1.22+**: 自动修复，每次迭代新变量。

---

## 实际应用

### 装饰器模式

```go
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        log.Printf("%s %s %v", r.Method, r.URL, time.Since(start))
    }
}
```

### 延迟执行

```go
func deferExample() {
    for i := 0; i < 3; i++ {
        defer func(n int) {
            fmt.Println(n)
        }(i)  // 传参求值
    }
    // 输出: 2, 1, 0
}
```
