# 错误处理 (Error Handling)

> **分类**: 语言设计

---

## 设计原则

Go 采用**显式错误返回**而非异常。

```go
func doSomething() error {
    f, err := os.Open("file")
    if err != nil {
        return err
    }
    defer f.Close()

    // 使用 f
    return nil
}
```

---

## Error 接口

```go
type error interface {
    Error() string
}
```

任何实现 `Error()` 方法的类型都是 error。

---

## 错误创建

### 1. errors.New

```go
return errors.New("something went wrong")
```

### 2. fmt.Errorf

```go
return fmt.Errorf("open file: %v", err)

// Go 1.13+: 错误包装
return fmt.Errorf("open file: %w", err)
```

---

## 错误检查

### 标准检查

```go
if err != nil {
    // 处理错误
}
```

### 错误比较 (1.13+)

```go
if errors.Is(err, os.ErrNotExist) {
    // 文件不存在
}
```

### 错误类型断言 (1.13+)

```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // 获取具体类型
    fmt.Println(pathErr.Path)
}
```

---

## 自定义错误

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// 使用
return &ValidationError{Field: "email", Message: "invalid format"}
```

---

## 最佳实践

### 1. 立即处理或返回

```go
// 好: 立即返回
f, err := os.Open("file")
if err != nil {
    return err
}

// 不好: 延迟处理
f, err := os.Open("file")
// ... 其他代码
if err != nil {
    return err
}
```

### 2. 添加上下文

```go
if err != nil {
    return fmt.Errorf("process user %d: %w", userID, err)
}
```

### 3. 不忽略错误

```go
// 不好
f, _ := os.Open("file")

// 好
f, err := os.Open("file")
if err != nil {
    return err
}
```

---

## 优缺点

| 优点 | 缺点 |
|------|------|
| 显式 | 代码冗长 |
| 可预测 | 易遗漏 |
| 无隐藏控制流 | 多层嵌套 |
| 编译时检查 | 无堆栈信息 |
