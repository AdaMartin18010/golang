# 标准库API

**难度**: 入门 | **预计阅读**: 15分钟

---

## 📖 常用标准库

### fmt - 格式化I/O

```go
// 打印
fmt.Println("Hello")
fmt.Printf("Number: %d\n", 42)

// 格式化字符串
s := fmt.Sprintf("Value: %v", value)

// 扫描输入
var name string
fmt.Scanf("%s", &name)
```

---

### strings - 字符串操作

```go
// 包含
strings.Contains("hello", "ell") // true

// 分割
parts := strings.Split("a,b,c", ",")

// 连接
s := strings.Join([]string{"a", "b"}, ",")

// 替换
s := strings.Replace("hello", "l", "L", -1)

// 大小写
upper := strings.ToUpper("hello")
lower := strings.ToLower("HELLO")

// 修剪
trimmed := strings.TrimSpace("  hello  ")
```

---

### time - 时间处理

```go
// 当前时间
now := time.Now()

// 解析时间
t, _ := time.Parse("2006-01-02", "2025-10-28")

// 格式化
formatted := time.Now().Format("2006-01-02 15:04:05")

// 时间运算
tomorrow := time.Now().Add(24 * time.Hour)
diff := t2.Sub(t1)

// 定时器
timer := time.NewTimer(5 * time.Second)
<-timer.C

// Ticker
ticker := time.NewTicker(time.Second)
for t := range ticker.C {
    fmt.Println(t)
}
```

---

### os - 操作系统接口

```go
// 环境变量
value := os.Getenv("PATH")
os.Setenv("KEY", "value")

// 文件操作
file, _ := os.Create("file.txt")
defer file.Close()
file.WriteString("content")

data, _ := os.ReadFile("file.txt")

// 命令行参数
args := os.Args[1:]
```

---

### io - I/O原语

```go
// 复制
io.Copy(dst, src)

// 读取全部
data, _ := io.ReadAll(reader)

// 管道
reader, writer := io.Pipe()

// 多写
w := io.MultiWriter(os.Stdout, file)
```

---

## 📚 相关资源

- [Standard Library Documentation](https://pkg.go.dev/std)

**下一步**: [03-常用第三方库](./03-常用第三方库.md)

---

**最后更新**: 2025-10-28

