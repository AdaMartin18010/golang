# regexp 包

> **分类**: 开源技术堆栈

---

## 基本使用

```go
import "regexp"

// 编译正则
re := regexp.MustCompile(`\b\w+@\w+\.\w+\b`)

// 匹配
matched := re.MatchString("contact@example.com")

// 查找
email := re.FindString("Email me at john@example.com")

// 提取所有
emails := re.FindAllString(text, -1)
```

---

## 常用方法

| 方法 | 说明 |
|------|------|
| `MatchString` | 是否匹配 |
| `FindString` | 查找第一个 |
| `FindAllString` | 查找所有 |
| `FindStringSubmatch` | 提取子匹配 |
| `ReplaceAllString` | 替换 |
| `Split` | 分割 |

---

## 子匹配提取

```go
re := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)

matches := re.FindStringSubmatch("john@example.com")
// matches[0] = "john@example.com"
// matches[1] = "john"
// matches[2] = "example"
// matches[3] = "com"

// 命名捕获
re2 := regexp.MustCompile(`(?P<user>\w+)@(?P<host>\w+)\.(?P<domain>\w+)`)
```

---

## 替换

```go
re := regexp.MustCompile(`\s+`)

// 替换为单个空格
result := re.ReplaceAllString("hello   world", " ")
// "hello world"

// 使用函数
re.ReplaceAllStringFunc(text, func(s string) string {
    return strings.ToUpper(s)
})
```

---

## 预编译优化

```go
// ✅ 全局预编译
var emailRegex = regexp.MustCompile(`^[\w.-]+@[\w.-]+\.\w+$`)

func ValidateEmail(email string) bool {
    return emailRegex.MatchString(email)
}

// ❌ 不要每次编译
func BadValidate(email string) bool {
    re := regexp.MustCompile(`^[\w.-]+@[\w.-]+\.\w+$`)  // 浪费！
    return re.MatchString(email)
}
```
