# 📚 Go 1.25.3 快速参考手册 - 2025

**版本**: Go 1.25.3  
**更新日期**: 2025-10-28  
**类型**: 快速参考  
**用途**: 日常开发速查

---

## 📋 目录

- [1. 基础语法速查](#1-基础语法速查)
- [2. 类型系统速查](#2-类型系统速查)
- [3. 并发编程速查](#3-并发编程速查)
- [4. 常用标准库](#4-常用标准库)
- [5. 错误处理模式](#5-错误处理模式)
- [6. 性能优化技巧](#6-性能优化技巧)
- [7. 常见陷阱](#7-常见陷阱)
- [8. 最佳实践](#8-最佳实践)

---

## 1. 基础语法速查

### 变量声明

```go
// var声明
var x int                    // 零值: 0
var y int = 10              // 初始化
var z = 20                  // 类型推导

// 短变量声明
a := 30                     // 只能在函数内

// 批量声明
var (
    name string
    age  int
    addr string
)

// 多变量
x, y := 1, 2
x, y = y, x                 // 交换
```

### 常量

```go
const Pi = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)

// iota枚举
const (
    Monday = iota    // 0
    Tuesday          // 1
    Wednesday        // 2
)
```

### 控制流

```go
// if
if x > 0 {
    // ...
}

if err := doSomething(); err != nil {
    return err
}

// switch
switch x {
case 1:
    // ...
case 2, 3:
    // ...
default:
    // ...
}

// 无表达式switch
switch {
case x > 0:
    // ...
case x < 0:
    // ...
}

// for
for i := 0; i < 10; i++ {
    // ...
}

for condition {
    // while风格
}

for {
    // 无限循环
}

// range
for i, v := range slice {
    // ...
}
for k, v := range map {
    // ...
}
for i, r := range "string" {
    // 按rune遍历
}
```

---

## 2. 类型系统速查

### 基本类型

```go
// 布尔
var b bool // false

// 整数
var i8  int8   // -128 to 127
var i16 int16
var i32 int32
var i64 int64
var i   int    // 平台相关

var u8  uint8  // 0 to 255
var u16 uint16
var u32 uint32
var u64 uint64
var u   uint

// 浮点
var f32 float32
var f64 float64

// 复数
var c64  complex64
var c128 complex128

// 字符串
var s string // ""

// 字节/字符
var b byte // uint8
var r rune // int32
```

### 复合类型

```go
// 数组
var arr [5]int
arr := [3]int{1, 2, 3}
arr := [...]int{1, 2, 3, 4}

// 切片
var s []int
s := []int{1, 2, 3}
s := make([]int, 5)      // len=5, cap=5
s := make([]int, 5, 10)  // len=5, cap=10

// 映射
var m map[string]int
m := make(map[string]int)
m := map[string]int{"a": 1, "b": 2}

// 结构体
type Person struct {
    Name string
    Age  int
}
p := Person{"Alice", 30}
p := Person{Name: "Bob", Age: 25}
```

### 指针、函数、通道

```go
// 指针
var p *int
x := 42
p = &x
*p = 100

// 函数
func add(a, b int) int {
    return a + b
}
f := func(x int) int { return x * 2 }

// 通道
ch := make(chan int)       // 无缓冲
ch := make(chan int, 10)   // 缓冲

// 发送/接收
ch <- 42
v := <-ch
v, ok := <-ch

close(ch)
```

### 接口

```go
// 定义
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 空接口
var any interface{}
var any any  // Go 1.18+

// 类型断言
s := i.(string)
s, ok := i.(string)

// 类型switch
switch v := i.(type) {
case int:
    // v是int
case string:
    // v是string
}
```

### 泛型（Go 1.18+）

```go
// 泛型函数
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 泛型类型
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// 使用
result := Min(3, 5)
stack := Stack[int]{}
stack.Push(42)
```

---

## 3. 并发编程速查

### Goroutine

```go
// 启动
go func() {
    // ...
}()

go myFunction()

// 传参（避免闭包陷阱）
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}
```

### Channel

```go
// 创建
ch := make(chan int)       // 无缓冲
ch := make(chan int, 10)   // 缓冲

// 操作
ch <- value        // 发送
value := <-ch      // 接收
value, ok := <-ch  // 检查关闭
close(ch)          // 关闭

// 单向channel
func send(ch chan<- int) {}    // 只写
func recv(ch <-chan int) {}     // 只读

// select
select {
case v := <-ch1:
    // ...
case ch2 <- value:
    // ...
case <-time.After(1 * time.Second):
    // 超时
default:
    // 非阻塞
}
```

### 同步原语

```go
// Mutex
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

// RWMutex
var rwmu sync.RWMutex
rwmu.RLock()    // 读锁
rwmu.RUnlock()
rwmu.Lock()     // 写锁
rwmu.Unlock()

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // ...
}()
wg.Wait()

// Once
var once sync.Once
once.Do(func() {
    // 只执行一次
})

// atomic
import "sync/atomic"
var counter int64
atomic.AddInt64(&counter, 1)
value := atomic.LoadInt64(&counter)
atomic.StoreInt64(&counter, 100)
```

### Context

```go
// 创建
ctx := context.Background()
ctx := context.TODO()

// WithCancel
ctx, cancel := context.WithCancel(parent)
defer cancel()

// WithTimeout
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// WithDeadline
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)
defer cancel()

// WithValue
ctx := context.WithValue(parent, key, value)
value := ctx.Value(key)

// 使用
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    return result
}
```

---

## 4. 常用标准库

### fmt - 格式化

```go
// 打印
fmt.Println("Hello")
fmt.Printf("x=%d\n", x)
fmt.Sprintf("x=%d", x)

// 格式化动词
%v   // 默认格式
%+v  // 带字段名的结构体
%#v  // Go语法表示
%T   // 类型
%t   // bool
%d   // 十进制整数
%f   // 浮点数
%s   // 字符串
%p   // 指针
```

### strings - 字符串

```go
strings.Contains(s, substr)
strings.HasPrefix(s, prefix)
strings.HasSuffix(s, suffix)
strings.Index(s, substr)
strings.Join([]string{"a", "b"}, ",")
strings.Split(s, ",")
strings.Replace(s, old, new, n)
strings.ToLower(s)
strings.ToUpper(s)
strings.TrimSpace(s)
```

### strconv - 转换

```go
// 字符串 → 数字
i, err := strconv.Atoi("42")
i, err := strconv.ParseInt("42", 10, 64)
f, err := strconv.ParseFloat("3.14", 64)

// 数字 → 字符串
s := strconv.Itoa(42)
s := strconv.FormatInt(42, 10)
s := strconv.FormatFloat(3.14, 'f', 2, 64)
```

### time - 时间

```go
// 当前时间
now := time.Now()

// 时间操作
later := now.Add(1 * time.Hour)
diff := t2.Sub(t1)

// 格式化（Magic: 2006-01-02 15:04:05）
s := now.Format("2006-01-02 15:04:05")
t, err := time.Parse("2006-01-02", "2025-10-28")

// 定时器
timer := time.NewTimer(5 * time.Second)
<-timer.C

ticker := time.NewTicker(1 * time.Second)
for t := range ticker.C {
    // 每秒执行
}
ticker.Stop()

// Sleep
time.Sleep(1 * time.Second)
```

### io - 输入输出

```go
// Reader
n, err := reader.Read(buf)

// Writer
n, err := writer.Write(data)

// Copy
n, err := io.Copy(dst, src)

// 常用工具
io.ReadAll(r)
io.WriteString(w, s)
```

### os - 操作系统

```go
// 文件操作
f, err := os.Open("file.txt")
f, err := os.Create("file.txt")
defer f.Close()

data, err := os.ReadFile("file.txt")
err := os.WriteFile("file.txt", data, 0644)

// 目录
err := os.Mkdir("dir", 0755)
err := os.MkdirAll("dir/subdir", 0755)
err := os.Remove("file.txt")
err := os.RemoveAll("dir")

// 环境变量
value := os.Getenv("PATH")
os.Setenv("KEY", "value")

// 参数
args := os.Args  // [程序名, arg1, arg2, ...]
```

### net/http - HTTP

```go
// 客户端
resp, err := http.Get("https://example.com")
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)

client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Get(url)

// 服务端
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello")
}
```

### encoding/json - JSON

```go
// 编码
data, err := json.Marshal(v)
data, err := json.MarshalIndent(v, "", "  ")

// 解码
err := json.Unmarshal(data, &v)

// 结构体标签
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"-"` // 忽略
}
```

---

## 5. 错误处理模式

### 基本模式

```go
// 返回错误
func DoSomething() error {
    if err := step1(); err != nil {
        return err
    }
    return nil
}

// 包装错误（Go 1.13+）
import "fmt"
return fmt.Errorf("step1 failed: %w", err)

// 判断错误
import "errors"
if errors.Is(err, os.ErrNotExist) {
    // ...
}

// 错误类型断言
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println(pathErr.Path)
}
```

### 自定义错误

```go
// 简单错误
var ErrNotFound = errors.New("not found")

// 错误类型
type ValidationError struct {
    Field string
    Value interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid %s: %v", e.Field, e.Value)
}
```

### defer处理

```go
func process() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    // 可能panic的代码
    return nil
}
```

---

## 6. 性能优化技巧

### 切片预分配

```go
// ❌ 低效
var slice []int
for i := 0; i < 1000; i++ {
    slice = append(slice, i)
}

// ✅ 高效
slice := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    slice = append(slice, i)
}
```

### 避免不必要的分配

```go
// ❌ 每次分配
for i := 0; i < n; i++ {
    buf := make([]byte, 1024)
    // 使用buf
}

// ✅ 复用
buf := make([]byte, 1024)
for i := 0; i < n; i++ {
    // 复用buf
}
```

### 字符串拼接

```go
// ❌ 低效（大量拼接）
s := ""
for i := 0; i < 1000; i++ {
    s += "x"
}

// ✅ 高效
var b strings.Builder
for i := 0; i < 1000; i++ {
    b.WriteString("x")
}
s := b.String()
```

### sync.Pool复用对象

```go
var pool = sync.Pool{
    New: func() interface{} {
        return new(MyObject)
    },
}

obj := pool.Get().(*MyObject)
// 使用obj
pool.Put(obj)
```

### 避免接口装箱

```go
// ❌ 装箱开销
func process(values []interface{}) {
    for _, v := range values {
        // ...
    }
}

// ✅ 泛型避免装箱（Go 1.18+）
func process[T any](values []T) {
    for _, v := range values {
        // ...
    }
}
```

---

## 7. 常见陷阱

### 1. 闭包循环变量

```go
// ❌ 错误
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // 可能全输出5
    }()
}

// ✅ 正确
for i := 0; i < 5; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}
```

### 2. range复制

```go
// ❌ 修改副本
for _, item := range items {
    item.Value = 100  // 无效！
}

// ✅ 使用索引
for i := range items {
    items[i].Value = 100
}
```

### 3. 切片陷阱

```go
// 共享底层数组
a := []int{1, 2, 3, 4, 5}
b := a[1:3]  // [2 3]
b[0] = 20
// a变成[1 20 3 4 5]

// ✅ 复制避免
b := make([]int, 2)
copy(b, a[1:3])
```

### 4. defer执行时机

```go
// ❌ 立即求值
x := 10
defer fmt.Println(x)  // 打印10
x = 20

// ✅ 延迟求值
defer func() {
    fmt.Println(x)  // 打印20
}()
```

### 5. map并发不安全

```go
// ❌ 数据竞争
m := make(map[string]int)
go func() { m["a"] = 1 }()
go func() { m["b"] = 2 }()

// ✅ 使用锁
var mu sync.RWMutex
go func() {
    mu.Lock()
    m["a"] = 1
    mu.Unlock()
}()

// ✅ 或使用sync.Map
var m sync.Map
m.Store("a", 1)
```

---

## 8. 最佳实践

### 命名规范

```go
// 包名：小写，单词
package httputil

// 导出：大写开头
type User struct {}
func GetUser() {}

// 未导出：小写开头
var counter int
func processData() {}

// 接口：-er后缀
type Reader interface {}
type Writer interface {}

// 缩写：全大写或全小写
var userID int  // ✅
var userId int  // ❌

// 常量：驼峰式
const MaxRetries = 3
```

### 错误处理

```go
// ✅ 立即处理
result, err := doSomething()
if err != nil {
    return err
}

// ✅ 包装错误
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// ❌ 忽略错误
doSomething()  // 不要这样
```

### 函数设计

```go
// ✅ 参数少
func New(name string, age int) *User

// ❌ 参数多，考虑使用options或config
func New(name, email, phone string, age, score int, active bool)

// ✅ Options模式
type Options struct {
    Timeout time.Duration
    Retries int
}

func New(opts Options) *Client
```

### 并发安全

```go
// ✅ 使用channel通信
ch := make(chan int)
go producer(ch)
consumer(ch)

// ✅ 使用锁保护共享数据
var (
    data map[string]int
    mu   sync.RWMutex
)

// ✅ 使用Context控制生命周期
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            doWork()
        }
    }
}
```

### 资源管理

```go
// ✅ 使用defer确保释放
func process() error {
    f, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer f.Close()
    
    // 处理文件
    return nil
}
```

---

## 📚 相关文档

- [Go 1.25.3完整知识体系总览](./00-Go-1.25.3完整知识体系总览-2025.md)
- [核心机制完整解析](./fundamentals/language/00-Go-1.25.3核心机制完整解析/)
- [常见问题解答](./reference/guides/03-常见问题.md)

---

**更新日期**: 2025-10-28  
**版本**: v1.0  
**维护**: Go形式化理论体系项目组

---

> **快速参考，日常必备** 📖

