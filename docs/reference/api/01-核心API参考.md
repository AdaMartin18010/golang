# 核心API参考

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

## 📋 目录

- [核心API参考](#核心api参考)
  - [📋 目录](#-目录)
  - [1. 📖 标准库核心包](#1--标准库核心包)
    - [net/http](#nethttp)
    - [encoding/json](#encodingjson)
    - [Context](#context)
    - [sync](#sync)
    - [fmt](#fmt)
    - [io](#io)
    - [time](#time)
    - [os](#os)
    - [strings](#strings)
    - [strconv](#strconv)
    - [errors](#errors)
    - [log](#log)
  - [📚 相关资源](#-相关资源)
  - [🔗 导航](#-导航)

---

---

## 1. 📖 标准库核心包

### net/http

```go
// HTTP服务器
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)

// HTTP客户端
resp, err := http.Get("https://example.com")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)

// 自定义请求
req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
req.Header.Set("Content-Type", "application/json")
client := &http.Client{}
resp, _ := client.Do(req)
```

---

### encoding/json

```go
// 序列化
data, _ := json.Marshal(struct)

// 反序列化
var result MyStruct
json.Unmarshal(data, &result)

// 流式编解码
json.NewEncoder(w).Encode(data)
json.NewDecoder(r.Body).Decode(&data)
```

---

### Context

```go
// 超时控制
ctx, cancel := Context.WithTimeout(Context.Background(), 5*time.Second)
defer cancel()

// 值传递
ctx = Context.WithValue(ctx, "key", "value")
value := ctx.Value("key")

// 取消信号
ctx, cancel := Context.WithCancel(Context.Background())
go func() {
    <-ctx.Done()
    // 清理...
}()
cancel()
```

---

### sync

```go
// Mutex
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

// 读写锁
var rwmu sync.RWMutex
rwmu.RLock()
defer rwmu.RUnlock()

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // 工作...
}()
wg.Wait()

// Once - 确保只执行一次
var once sync.Once
once.Do(func() {
    // 初始化代码，只执行一次
})

// Map - 并发安全的map
var m sync.Map
m.Store("key", "value")
if v, ok := m.Load("key"); ok {
    fmt.Println(v)
}
```

---

### fmt

```go
// 格式化输出
fmt.Printf("Name: %s, Age: %d\n", name, age)
fmt.Sprintf("Formatted: %v", value)

// 打印到Writer
fmt.Fprintf(w, "Output: %s\n", data)

// Scan输入
var input string
fmt.Scanf("%s", &input)

// 常用动词
// %v  默认格式
// %+v 带字段名
// %#v Go语法表示
// %T  类型
// %t  布尔值
// %d  十进制整数
// %f  浮点数
// %s  字符串
// %p  指针
```

---

### io

```go
// 复制
io.Copy(dst, src)

// 限制读取
r := io.LimitReader(reader, 1024)

// 多个Reader
r := io.MultiReader(r1, r2, r3)

// 读取全部
data, _ := io.ReadAll(reader)

// 管道
pr, pw := io.Pipe()
go func() {
    pw.Write([]byte("data"))
    pw.Close()
}()
io.Copy(os.Stdout, pr)

// TeeReader - 同时读取和写入
r := io.TeeReader(reader, writer)
```

---

### time

```go
// 当前时间
now := time.Now()

// 解析
t, _ := time.Parse("2006-01-02", "2025-10-28")

// 格式化
s := t.Format("2006-01-02 15:04:05")

// 时间运算
future := now.Add(24 * time.Hour)
past := now.Add(-1 * time.Hour)

// 定时器
timer := time.NewTimer(5 * time.Second)
<-timer.C

// Ticker
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for range ticker.C {
    // 每秒执行
}

// Sleep
time.Sleep(2 * time.Second)
```

---

### os

```go
// 文件操作
f, _ := os.Open("file.txt")
defer f.Close()

f, _ := os.Create("output.txt")
f.WriteString("content")

// 环境变量
val := os.Getenv("PATH")
os.Setenv("MY_VAR", "value")

// 命令行参数
args := os.Args

// 目录
os.Mkdir("dir", 0755)
os.MkdirAll("path/to/dir", 0755)
os.Remove("file")
os.RemoveAll("dir")

// 文件信息
info, _ := os.Stat("file.txt")
size := info.Size()
isDir := info.IsDir()
```

---

### strings

```go
// 连接
s := strings.Join([]string{"a", "b", "c"}, ",")

// 分割
parts := strings.Split("a,b,c", ",")

// 判断
strings.Contains("hello", "ll")     // true
strings.HasPrefix("hello", "he")    // true
strings.HasSuffix("hello", "lo")    // true

// 替换
s := strings.Replace("hello", "l", "L", -1)
s := strings.ReplaceAll("hello", "l", "L")

// 大小写
strings.ToUpper("hello")  // HELLO
strings.ToLower("HELLO")  // hello

// 修剪
strings.TrimSpace("  hello  ")  // "hello"
strings.Trim("xxhelloxx", "x")  // "hello"
```

---

### strconv

```go
// 字符串转数字
i, _ := strconv.Atoi("42")
i64, _ := strconv.ParseInt("42", 10, 64)
f, _ := strconv.ParseFloat("3.14", 64)
b, _ := strconv.ParseBool("true")

// 数字转字符串
s := strconv.Itoa(42)
s := strconv.FormatInt(42, 10)
s := strconv.FormatFloat(3.14, 'f', 2, 64)
s := strconv.FormatBool(true)

// Quote
s := strconv.Quote("hello\nworld")  // "\"hello\\nworld\""
```

---

### errors

```go
// 创建错误
err := errors.New("something went wrong")

// 包装错误 (Go 1.13+)
err := fmt.Errorf("failed to process: %w", originalErr)

// 判断错误
if errors.Is(err, os.ErrNotExist) {
    // 处理文件不存在
}

// 类型断言
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Path:", pathErr.Path)
}

// 连接多个错误 (Go 1.20+)
err := errors.Join(err1, err2, err3)
```

---

### log

```go
// 基本日志
log.Println("Info message")
log.Printf("User: %s", username)

// 致命错误
log.Fatal("Fatal error")  // 输出后os.Exit(1)

// Panic
log.Panic("Panic message")  // 输出后panic

// 自定义Logger
logger := log.New(os.Stdout, "PREFIX: ", log.Ldate|log.Ltime)
logger.Println("Custom log")

// 标志
log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

// 输出目标
log.SetOutput(file)
```

---

## 📚 相关资源

- [Go Standard Library](https://pkg.go.dev/std)
- [Go 1.25.3 Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)

---

## 🔗 导航
