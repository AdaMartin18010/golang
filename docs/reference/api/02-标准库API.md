# 标准库API

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 📖 常用标准库](#1.-常用标准库)
  - [fmt - 格式化I/O](#fmt-格式化io)
  - [strings - 字符串操作](#strings-字符串操作)
  - [time - 时间处理](#time-时间处理)
  - [os - 操作系统接口](#os-操作系统接口)
  - [io - I/O原语](#io-io原语)
  - [bufio - 缓冲I/O](#bufio-缓冲io)
  - [path/filepath - 文件路径](#pathfilepath-文件路径)
  - [regexp - 正则表达式](#regexp-正则表达式)
  - [math - 数学函数](#math-数学函数)
  - [math/rand - 随机数](#mathrand-随机数)
  - [sort - 排序](#sort-排序)
  - [flag - 命令行参数](#flag-命令行参数)
  - [database/sql - 数据库](#databasesql-数据库)
  - [crypto - 加密](#crypto-加密)
  - [compress - 压缩](#compress-压缩)
  - [archive/tar - TAR归档](#archivetar-tar归档)
- [📚 相关资源](#相关资源)
- [🔗 导航](#导航)

## 1. 📖 常用标准库

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

// TeeReader
r := io.TeeReader(reader, writer)

// LimitReader
r := io.LimitReader(reader, 1024)
```

---

### bufio - 缓冲I/O

```go
// Reader
reader := bufio.NewReader(file)
line, _ := reader.ReadString('\n')
bytes, _ := reader.ReadBytes('\n')

// Scanner - 按行读取
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}

// Writer
writer := bufio.NewWriter(file)
writer.WriteString("data")
writer.Flush()  // 刷新缓冲区
```

---

### path/filepath - 文件路径

```go
// 路径操作
abs, _ := filepath.Abs("./file.txt")
dir := filepath.Dir("/path/to/file.txt")    // /path/to
base := filepath.Base("/path/to/file.txt")  // file.txt
ext := filepath.Ext("file.txt")             // .txt

// 连接路径
path := filepath.Join("dir", "subdir", "file.txt")

// 遍历目录
filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
        fmt.Println(path)
    }
    return nil
})

// 匹配模式
matched, _ := filepath.Match("*.go", "main.go")  // true
```

---

### regexp - 正则表达式

```go
// 匹配
matched, _ := regexp.MatchString(`\d+`, "age: 25")  // true

// 编译
re := regexp.MustCompile(`\d+`)

// 查找
found := re.FindString("age: 25")  // "25"
all := re.FindAllString("1 2 3", -1)  // ["1", "2", "3"]

// 替换
result := re.ReplaceAllString("age: 25", "XX")  // "age: XX"

// 分组
re := regexp.MustCompile(`(\d+)-(\d+)-(\d+)`)
matches := re.FindStringSubmatch("2025-10-28")
// matches[0]: "2025-10-28"
// matches[1]: "2025"
// matches[2]: "10"
// matches[3]: "28"
```

---

### math - 数学函数

```go
// 基本运算
math.Abs(-10)           // 10
math.Ceil(3.14)         // 4
math.Floor(3.14)        // 3
math.Round(3.5)         // 4
math.Max(10, 20)        // 20
math.Min(10, 20)        // 10

// 指数和对数
math.Pow(2, 3)          // 8
math.Sqrt(16)           // 4
math.Log(10)            // 2.302585...
math.Exp(1)             // 2.718281... (e)

// 三角函数
math.Sin(math.Pi / 2)   // 1
math.Cos(0)             // 1
math.Tan(math.Pi / 4)   // 1
```

---

### math/rand - 随机数

```go
// 设置种子 (Go 1.20+自动)
rand.Seed(time.Now().UnixNano())

// 生成随机数
n := rand.Int()                  // 随机整数
n := rand.Intn(100)              // 0-99
f := rand.Float64()              // 0.0-1.0

// 随机字符串
letters := []rune("abcdefghijklmnopqrstuvwxyz")
b := make([]rune, 10)
for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
}
s := string(b)

// crypto/rand - 密码学安全随机数
import "crypto/rand"
bytes := make([]byte, 32)
rand.Read(bytes)
```

---

### sort - 排序

```go
// 整数排序
nums := []int{3, 1, 4, 1, 5, 9}
sort.Ints(nums)

// 字符串排序
strs := []string{"c", "a", "b"}
sort.Strings(strs)

// 自定义排序
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})

// 二分查找
index := sort.SearchInts(sortedNums, target)

// 检查是否已排序
isSorted := sort.IntsAreSorted(nums)
```

---

### flag - 命令行参数

```go
// 定义标志
var (
    host = flag.String("host", "localhost", "Server host")
    port = flag.Int("port", 8080, "Server port")
    debug = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
    flag.Parse()
    
    fmt.Printf("Host: %s\n", *host)
    fmt.Printf("Port: %d\n", *port)
    fmt.Printf("Debug: %v\n", *debug)
    
    // 剩余参数
    args := flag.Args()
}

// 使用: go run main.go -host=example.com -port=9000 -debug
```

---

### database/sql - 数据库

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// 连接
db, _ := sql.Open("mysql", "user:pass@tcp(localhost:3306)/db")
defer db.Close()

// 查询
rows, _ := db.Query("SELECT id, name FROM users WHERE age > ?", 18)
defer rows.Close()
for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
}

// 插入
result, _ := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 25)
id, _ := result.LastInsertId()

// 事务
tx, _ := db.Begin()
tx.Exec("UPDATE ...")
tx.Commit()  // 或 tx.Rollback()
```

---

### crypto - 加密

```go
import (
    "crypto/md5"
    "crypto/sha256"
    "crypto/sha512"
)

// MD5
h := md5.Sum([]byte("data"))
fmt.Printf("%x", h)

// SHA256
h := sha256.Sum256([]byte("data"))
fmt.Printf("%x", h)

// SHA512
h := sha512.Sum512([]byte("data"))
fmt.Printf("%x", h)

// HMAC
import "crypto/hmac"
mac := hmac.New(sha256.New, []byte("key"))
mac.Write([]byte("data"))
signature := mac.Sum(nil)
```

---

### compress - 压缩

```go
import (
    "compress/gzip"
    "compress/zlib"
)

// Gzip压缩
var buf bytes.Buffer
w := gzip.NewWriter(&buf)
w.Write([]byte("data"))
w.Close()

// Gzip解压
r, _ := gzip.NewReader(&buf)
data, _ := io.ReadAll(r)
r.Close()

// Zlib压缩/解压 (类似)
```

---

### archive/tar - TAR归档

```go
// 创建tar
tw := tar.NewWriter(file)
hdr := &tar.Header{
    Name: "file.txt",
    Mode: 0644,
    Size: int64(len(data)),
}
tw.WriteHeader(hdr)
tw.Write(data)
tw.Close()

// 读取tar
tr := tar.NewReader(file)
for {
    hdr, err := tr.Next()
    if err == io.EOF {
        break
    }
    data, _ := io.ReadAll(tr)
}
```

---

## 📚 相关资源

- [Standard Library Documentation](https://pkg.go.dev/std)
- [Go 1.25.3 Release Notes](https://go.dev/doc/go1.25)
- [Effective Go](https://go.dev/doc/effective_go)

---

## 🔗 导航

- **上一页**: [01-核心API参考](./01-核心API参考.md)
- **下一页**: [03-常用第三方库](./03-常用第三方库.md)
- **相关**: [README](./README.md)

---

**最后更新**: 2025-10-29  
**Go版本**: 1.25.3
