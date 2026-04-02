# 字符串处理 (String Handling)

> **分类**: 语言设计

---

## 字符串基础

```go
// Go 字符串是不可变字节序列
s := "Hello, 世界"

// 索引得到字节，不是字符
b := s[0]  // 72 (H)

// 长度是字节数
len(s)  // 13 (不是 9)

// 转换为 rune 切片处理字符
runes := []rune(s)
char := runes[7]  // '世'
```

---

## strings 包

### 常用函数

```go
import "strings"

// 包含
strings.Contains("hello", "ll")     // true
strings.HasPrefix("hello", "he")    // true
strings.HasSuffix("hello", "lo")    // true

// 查找
strings.Index("hello", "ll")        // 2
strings.LastIndex("hello", "l")     // 3

// 替换
strings.Replace("hello", "l", "L", 1)   // "heLlo"
strings.ReplaceAll("hello", "l", "L") // "heLLo"

// 分割与连接
parts := strings.Split("a,b,c", ",")  // ["a", "b", "c"]
joined := strings.Join(parts, "-")     // "a-b-c"

// 修剪
strings.TrimSpace("  hello  ")           // "hello"
strings.Trim("xxhelloxx", "x")          // "hello"
strings.TrimPrefix("hello", "he")       // "llo"
```

---

## Builder 高效拼接

```go
import "strings"

// ✅ 高效（推荐）
func buildString(items []string) string {
    var b strings.Builder
    b.Grow(100)  // 预分配容量

    for _, item := range items {
        b.WriteString(item)
        b.WriteByte(',')
    }

    return b.String()
}

// ❌ 低效
func badBuild(items []string) string {
    result := ""
    for _, item := range items {
        result += item + ","  // 每次分配新内存
    }
    return result
}
```

---

## 格式化

```go
import "fmt"

// Sprintf
name := "Alice"
age := 30
s := fmt.Sprintf("Name: %s, Age: %d", name, age)

// 常用动词
%d  // 十进制整数
%f  // 浮点数
%s  // 字符串
%v  // 默认格式
%+v // 带字段名
%#v // Go 语法格式
%T  // 类型
%%  // 百分号

// 宽度与精度
fmt.Sprintf("%10s", "hi")     // "        hi"
fmt.Sprintf("%-10s", "hi")    // "hi        "
fmt.Sprintf("%.2f", 3.14159)  // "3.14"
```

---

## Unicode 处理

```go
import "unicode"
import "unicode/utf8"

// 字符分类
unicode.IsDigit('1')     // true
unicode.IsLetter('a')    // true
unicode.IsSpace(' ')     // true
unicode.IsChinese('世')  // true

// UTF-8 解码
s := "Hello, 世界"
for len(s) > 0 {
    r, size := utf8.DecodeRuneInString(s)
    fmt.Printf("%c ", r)
    s = s[size:]
}

// 统计字符数
utf8.RuneCountInString("Hello, 世界")  // 9

// 验证 UTF-8
utf8.ValidString("hello")  // true
```

---

## 性能对比

| 操作 | 方法 | 时间复杂度 |
|------|------|-----------|
| 拼接少量 | + | O(n) |
| 拼接大量 | strings.Builder | O(n) |
| 分割 | strings.Split | O(n) |
| 查找 | strings.Index | O(n) |
| 替换 | strings.Replace | O(n*m) |

---

## 最佳实践

```go
// 1. 需要修改时用 Builder
var b strings.Builder

// 2. 预分配容量
b.Grow(expectedSize)

// 3. 比较时统一大小写
strings.EqualFold("Hello", "hello")  // true

// 4. 大量重复操作预编译正则
var re = regexp.MustCompile(`\d+`)

// 5. 字符串与字节切片转换零拷贝
b := []byte(s)  // 分配新内存
// 无法零拷贝，Go 字符串不可变
```
