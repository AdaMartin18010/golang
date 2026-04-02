# 整洁代码 (Clean Code)

> **分类**: 工程与云原生

---

## 命名规范

### 变量

```go
// 好: 有意义
userCount := 10
maxRetries := 3

// 不好: 无意义
n := 10
m := 3
```

### 函数

```go
// 好: 动词 + 名词
func CalculateTotal(items []Item) float64
func ValidateEmail(email string) error

// 不好: 模糊
func Do(items []Item) float64
func Check(s string) error
```

### 接口

```go
// 好: 动词 + er
type Reader interface { Read() }
type Writer interface { Write() }

// 不好: 名词
type Read interface {}
```

---

## 函数设计

### 单一职责

```go
// 好: 只做一件事
func ParseUser(data []byte) (*User, error)
func SaveUser(user *User) error

// 不好: 多职责
func ProcessUser(data []byte) (*User, error) {
    // 解析 + 验证 + 保存
}
```

---

## 错误处理

### 错误传播

```go
// 好: 添加上下文
if err != nil {
    return fmt.Errorf("fetch user %d: %w", id, err)
}

// 不好: 吞没错误
if err != nil {
    return nil
}
```

### 早返回

```go
// 好: 减少嵌套
func Process(data []byte) error {
    if len(data) == 0 {
        return errors.New("empty data")
    }

    user, err := Parse(data)
    if err != nil {
        return err
    }

    return Save(user)
}
```

---

## 参考

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
