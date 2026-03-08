# 附录 A：快速参考卡片

> Go 1.26 快速参考 - 随时查阅的手册

---

## A.1 语法速查

### 基础语法

```go
// 变量声明
var x int = 10           // 完整声明
x := 10                  // 短声明（函数内）
const Pi = 3.14          // 常量

// 数据类型
int, int8, int16, int32, int64    // 整数
uint, uintptr                      // 无符号整数
float32, float64                   // 浮点
complex64, complex128              // 复数
string                             // 字符串
bool                               // 布尔
byte // alias for uint8
rune // alias for int32

// 复合类型
[5]int          // 数组
[]int           // 切片
map[K]V         // 映射
chan T          // 通道
func(T) R       // 函数
struct{...}     // 结构体
interface{...}  // 接口
```

### 控制结构

```go
// 条件
if x > 0 { ... } else { ... }
if err := do(); err != nil { ... }

// 循环
for i := 0; i < n; i++ { ... }     // C-style
for condition { ... }              // while
for { ... }                        // 无限
for i, v := range slice { ... }    // range
for i := range n { ... }           // Go 1.22+ 整数 range

// switch
switch x {
case 1: ...
case 2, 3: ...
default: ...
}

switch {  // 无表达式
 case x > 0: ...
 case x < 0: ...
}

switch x := i.(type) {  // 类型 switch
 case int: ...
 case string: ...
}

// defer/panic/recover
defer cleanup()           // LIFO 执行
panic("error")            // 运行时错误
recover()                 // 捕获 panic
```

---

## A.2 并发速查

### Goroutine 和 Channel

```go
// 启动 goroutine
go function()
go func() { ... }()

// Channel 操作
ch := make(chan int)      // 无缓冲
ch := make(chan int, 10)  // 有缓冲

ch <- v                   // 发送
v := <-ch                 // 接收
v, ok := <-ch             // 检查关闭
close(ch)                 // 关闭

// select
select {
case v := <-ch1:
    // 处理
 case ch2 <- v:
    // 发送
 case <-time.After(5 * time.Second):
    // 超时
 default:
    // 默认
}

// 同步原语
var mu sync.Mutex         // 互斥锁
mu.Lock(); defer mu.Unlock()

var rw sync.RWMutex       // 读写锁
rw.RLock(); defer rw.RUnlock()

var wg sync.WaitGroup     // 等待组
wg.Add(1); go func() { defer wg.Done() }()
wg.Wait()

var once sync.Once        // 单次执行
once.Do(func() { ... })

// Context
ctx := context.Background()
ctx, cancel := context.WithCancel(ctx)
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
```

---

## A.3 标准库速查

### 常用包

```go
// fmt - 格式化
fmt.Printf("%s %d %v %+v %#v %T %%", ...)
fmt.Sprintf("...", ...)
fmt.Fprintf(w, "...", ...)

// strings
strings.HasPrefix(s, prefix)
strings.HasSuffix(s, suffix)
strings.Contains(s, substr)
strings.Index(s, substr)
strings.Split(s, sep)
strings.Join(slice, sep)
strings.ToLower(s)
strings.ToUpper(s)
strings.TrimSpace(s)
strings.ReplaceAll(s, old, new)

// strconv
strconv.Atoi(s)
strconv.Itoa(i)
strconv.ParseFloat(s, 64)
strconv.FormatFloat(f, 'f', -1, 64)
strconv.ParseBool(s)
strconv.FormatBool(b)

// os
os.Open(filename)
os.Create(filename)
os.Stat(filename)
os.Remove(filename)
os.Getenv(key)
os.Setenv(key, value)
os.Exit(code)
os.Stdin / os.Stdout / os.Stderr

// io
io.Copy(dst, src)
io.ReadAll(r)
io.WriteString(w, s)

// bufio
scanner := bufio.NewScanner(r)
scanner.Scan()
line := scanner.Text()

writer := bufio.NewWriter(w)
writer.WriteString(s)
writer.Flush()

// time
time.Now()
time.Since(t)
time.Until(t)
time.Sleep(d)
time.Parse(layout, value)
t.Format(layout)  // "2006-01-02 15:04:05"

// encoding/json
json.Marshal(v)
json.Unmarshal(data, &v)
json.NewEncoder(w).Encode(v)
json.NewDecoder(r).Decode(&v)

// net/http
http.Get(url)
http.Post(url, contentType, body)
http.HandleFunc(pattern, handler)
http.ListenAndServe(addr, handler)

// log
log.Println(...)
log.Printf("...", ...)
log.Fatal(...)
log.Panic(...)

// testing
func TestXxx(t *testing.T) {
    t.Error(args)      // 继续
    t.Fatal(args)      // 停止
    t.Log(args)
}

func BenchmarkXxx(b *testing.B) {
    for b.Loop() { ... }  // Go 1.24+
}
```

---

## A.4 泛型速查

```go
// 预定义约束
import "golang.org/x/exp/constraints"

constraints.Integer       // ~int | ~int8 | ...
constraints.Signed        // 有符号整数
constraints.Unsigned      // 无符号整数
constraints.Float         // ~float32 | ~float64
constraints.Complex       // 复数
constraints.Ordered       // 可比较排序

// 类型约束定义
type Number interface {
    constraints.Integer | constraints.Float
}

type Stringer interface {
    String() string
}

// 泛型函数
func Min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}

// 泛型类型
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) { ... }
func (s *Stack[T]) Pop() (T, bool) { ... }

// 类型推断
Min[int](1, 2)     // 显式
Min(1, 2)          // 推断

// 类型集合
type IntOrString interface {
    ~int | ~string  // 近似类型
}

// Go 1.26: 递归约束
type Ordered[T Ordered[T]] interface {
    Less(T) bool
}
```

---

## A.5 错误处理速查

```go
// 错误创建
errors.New("message")
fmt.Errorf("message: %w", err)  // 包装

// 错误检查
if err != nil { ... }

errors.Is(err, target)        // 错误链包含
errors.As(err, &target)       // 类型断言

// Go 1.26: 泛型错误检查
if target, ok := errors.AsType[*MyError](err); ok {
    // target 是 *MyError 类型
}

// 自定义错误
type MyError struct {
    Op   string
    Err  error
}

func (e *MyError) Error() string {
    return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *MyError) Unwrap() error {
    return e.Err
}
```

---

## A.6 项目结构模板

```
myproject/
├── api/                    # API 定义 (protobuf/OpenAPI)
├── cmd/                    # 可执行程序入口
│   ├── server/
│   │   └── main.go
│   └── cli/
│       └── main.go
├── internal/               # 内部实现
│   ├── domain/            # 领域模型
│   ├── service/           # 业务逻辑
│   ├── repository/        # 数据访问
│   └── infrastructure/    # 基础设施
├── pkg/                    # 可导入的包
├── configs/               # 配置文件
├── scripts/               # 脚本
├── test/                  # 测试数据/集成测试
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

---

## A.7 Makefile 模板

```makefile
.PHONY: all build test clean run lint docker

APP_NAME=myapp
GO=go
GOFLAGS=-v

all: lint test build

build:
 $(GO) build $(GOFLAGS) -o bin/$(APP_NAME) ./cmd/server

test:
 $(GO) test $(GOFLAGS) -race -cover ./...

lint:
 golangci-lint run

clean:
 rm -rf bin/
 rm -f coverage.out

run: build
 ./bin/$(APP_NAME)

dev:
 air

docker:
 docker build -t $(APP_NAME):latest .

generate:
 $(GO) generate ./...

fmt:
 $(GO) fmt ./...

vet:
 $(GO) vet ./...
```

---

## A.8 常用命令

```bash
# 开发
go run main.go                    # 运行
go build                          # 构建
go install                        # 安装到 $GOPATH/bin
go test ./...                     # 运行测试
go test -v -race ./...            # 详细 + 竞态检测
go test -bench=. ./...            # 基准测试
go test -cover ./...              # 覆盖率

# 依赖管理
go mod init module_name           # 初始化模块
go mod tidy                       # 整理依赖
go mod download                   # 下载依赖
go mod vendor                     # 创建 vendor
go list -m all                    # 列出依赖
go get package@version            # 添加/更新依赖

# 代码工具
go fmt ./...                      # 格式化
go vet ./...                      # 静态分析
go fix ./...                      # 升级代码 (Go 1.26)

go generate ./...                 # 代码生成
godoc -http=:6060                 # 文档服务器

# 调试
go tool pprof cpu.prof            # CPU 分析
go tool pprof heap.prof           # 内存分析
go tool trace trace.out           # 追踪分析
go tool objdump binary            # 反汇编

# 交叉编译
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
GOOS=darwin GOARCH=arm64 go build
```

---

## A.9 Go 1.26 新特性速查

```go
// 1. new 支持表达式
ptr := new(int(42))           // *int = 42
ptr := new(compute())         // *T = compute() 结果

// 2. 递归类型约束
type Adder[A Adder[A]] interface {
    Add(A) A
}

// 3. errors.AsType
if e, ok := errors.AsType[*MyError](err); ok {
    // e 是 *MyError 类型
}

// 4. Green Tea GC (默认启用)
// GOEXPERIMENT=nogreenteagc 禁用

// 5. Goroutine 泄漏检测
// GOEXPERIMENT=goroutineleakprofile
import _ "net/http/pprof"
// 访问 /debug/pprof/goroutineleak
```

---

*打印此页作为日常开发参考。*
