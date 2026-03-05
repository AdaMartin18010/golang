# 第一章：Go 1.26 新特性概览

> 发布日期：2026年2月
> 兼容性：Go 1 兼容性承诺保证

---

## 1.1 语言层面的重大改进

### 1.1.1 `new` 函数支持表达式（核心特性）

Go 1.26 最显著的语法改进是 `new` 内置函数现在接受表达式作为参数，而不仅仅是类型。

#### 语法对比

**Go 1.25 及之前:**

```go
// 必须先声明变量，再取地址
x := int64(300)
ptr := &x

// 或者使用辅助函数
func intPtr(v int) *int { return &v }
age := intPtr(25)
```

**Go 1.26:**

```go
// 直接传入表达式
ptr := new(int64(300))

// 更实用的例子：可选字段初始化
type Config struct {
    Timeout *time.Duration `json:"timeout,omitempty"`
    Retries *int           `json:"retries,omitempty"`
}

func NewConfig() Config {
    return Config{
        // 简洁地设置可选字段
        Timeout: new(30 * time.Second),
        Retries: new(3),
    }
}
```

#### 实际应用场景

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

// API 请求中的可选字段
type CreateUserRequest struct {
    Name     string     `json:"name"`
    Age      *int       `json:"age,omitempty"`       // 可选
    Birthday *time.Time `json:"birthday,omitempty"`  // 可选
    Metadata *Metadata  `json:"metadata,omitempty"`  // 可选
}

type Metadata struct {
    Source  string `json:"source"`
    Version string `json:"version"`
}

func createUserRequest(name string, age int, birthday time.Time) ([]byte, error) {
    req := CreateUserRequest{
        Name: name,
        // Go 1.26: 一行代码设置可选字段
        Age:      new(age),
        Birthday: new(birthday),
        Metadata: new(Metadata{
            Source:  "web",
            Version: "1.0",
        }),
    }
    return json.Marshal(req)
}

func main() {
    data, _ := createUserRequest("Alice", 30, time.Date(1994, 1, 1, 0, 0, 0, 0, time.UTC))
    fmt.Println(string(data))
    // 输出: {"name":"Alice","age":30,"birthday":"1994-01-01T00:00:00Z","metadata":{"source":"web","version":"1.0"}}
}
```

**语义规则:**

```text
new(T)  // T 是类型 → 返回 *T，值为零值
new(E)  // E 是表达式 → 返回 *T（E 的类型），值为 E 的求值结果
```

### 1.1.2 泛型递归类型约束（核心特性）

Go 1.26 解除了泛型类型不能在其类型参数列表中引用自身的限制。

#### 概念解释

```go
// Go 1.25: 编译错误 - invalid recursive type
// type T[P T[P]] struct{}

// Go 1.26: 合法！
type Ordered[T Ordered[T]] interface {
    Less(T) bool
}

// 应用：通用比较器
type Tree[T Ordered[T]] struct {
    root *Node[T]
}

type Node[T Ordered[T]] struct {
    value T
    left  *Node[T]
    right *Node[T]
}

func (t *Tree[T]) Insert(v T) {
    // 可以使用 Less 方法比较
    if t.root == nil {
        t.root = &Node[T]{value: v}
        return
    }
    t.root.insert(v)
}

func (n *Node[T]) insert(v T) {
    if v.Less(n.value) {
        if n.left == nil {
            n.left = &Node[T]{value: v}
        } else {
            n.left.insert(v)
        }
    } else {
        if n.right == nil {
            n.right = &Node[T]{value: v}
        } else {
            n.right.insert(v)
        }
    }
}
```

#### 完整示例：自引用约束

```go
package main

import (
    "fmt"
    "net/netip"
)

// Adder 约束：类型必须能与同类型进行加法运算
type Adder[A Adder[A]] interface {
    Add(A) A
}

// 数值类型实现 Adder
func (n IntNumber) Add(other IntNumber) IntNumber {
    return n + other
}

type IntNumber int

// 泛型算法：对任何满足 Adder 约束的类型生效
func Sum[T Adder[T]](items []T) T {
    var sum T
    for _, item := range items {
        sum = sum.Add(item)
    }
    return sum
}

// Ordered 约束：类型必须可比较
type Ordered[T Ordered[T]] interface {
    Less(T) bool
}

// netip.Addr 天然满足 Ordered 约束，因为它有 Less 方法

type Container[T Ordered[T]] struct {
    items []T
}

func (c *Container[T]) Add(item T) {
    c.items = append(c.items, item)
}

func (c *Container[T]) FindMin() T {
    if len(c.items) == 0 {
        var zero T
        return zero
    }
    min := c.items[0]
    for _, item := range c.items[1:] {
        if item.Less(min) {
            min = item
        }
    }
    return min
}

func main() {
    // 使用 IntNumber
    numbers := []IntNumber{1, 2, 3, 4, 5}
    fmt.Println("Sum:", Sum(numbers)) // Sum: 15

    // 使用 netip.Addr（它有 Less 方法）
    addrs := Container[netip.Addr]{}
    addrs.Add(netip.MustParseAddr("192.168.1.1"))
    addrs.Add(netip.MustParseAddr("10.0.0.1"))
    addrs.Add(netip.MustParseAddr("172.16.0.1"))
    fmt.Println("Min IP:", addrs.FindMin()) // Min IP: 10.0.0.1
}
```

---

## 1.2 性能提升

### 1.2.1 Green Tea 垃圾回收器（默认启用）

```text
┌─────────────────────────────────────────────────────────────────┐
│                    Green Tea GC 优势                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  性能提升                    适用场景                             │
│  ─────────                   ─────────                          │
│  • 10-40% GC 开销降低        • 大量小对象分配                     │
│  • 更好的局部性               • 高并发应用                        │
│  • 更好的 CPU 可扩展性        • 内存密集型服务                     │
│  • 向量化扫描 (Ice Lake/Zen4+)                                   │
│                                                                 │
│  禁用方式（不推荐）：                                              │
│  GOEXPERIMENT=nogreenteagc go build                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**性能基准测试:**

```go
package main

import (
    "runtime"
    "testing"
)

// 模拟大量小对象分配
func BenchmarkSmallObjectAllocation(b *testing.B) {
    type Small struct {
        a, b, c int64
        d       string
    }

    b.ReportAllocs()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            s := &Small{a: 1, b: 2, c: 3, d: "test"}
            _ = s
        }
    })
}

// 查看 GC 统计
func PrintGCStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    println("GC runs:", m.NumGC)
    println("Total GC pause (ns):", m.PauseTotalNs)
    println("Heap alloc:", m.HeapAlloc)
    println("Heap objects:", m.HeapObjects)
}
```

### 1.2.2 cgo 开销降低 30%

```go
// cgo 调用示例
package main

/*
#include <stdio.h>
#include <stdlib.h>

void hello(const char* name) {
    printf("Hello, %s!\n", name);
}

int add(int a, int b) {
    return a + b;
}
*/
import "C"
import "unsafe"

func main() {
    // 调用 C 函数 - Go 1.26 更快！
    name := C.CString("Go 1.26")
    defer C.free(unsafe.Pointer(name))

    C.hello(name)

    result := C.add(40, 2)
    println("Result:", result)
}
```

### 1.2.3 栈上切片分配优化

Go 1.26 编译器可以在更多情况下将切片的底层数组分配在栈上。

```go
package main

// 这个函数在 Go 1.26 中，data 的底层数组可能分配在栈上
func processSmallData() int {
    data := make([]int, 10)  // 小切片，可能栈分配
    sum := 0
    for i := range data {
        data[i] = i * 2
        sum += data[i]
    }
    return sum
}

// 大切片仍然堆分配
func processLargeData() int {
    data := make([]int, 10000)  // 大切片，堆分配
    sum := 0
    for i := range data {
        data[i] = i
        sum += data[i]
    }
    return sum
}
```

**逃逸分析检查:**

```bash
# 查看逃逸分析结果
go build -gcflags="-m" main.go

# 禁用新优化（用于调试）
go build -gcflags="-d=variablemakehash=n" main.go
```

---

## 1.3 工具链改进

### 1.3.1 完全重写的 `go fix`

Go 1.26 的 `go fix` 命令基于 Go 分析框架完全重写，引入了 **modernizers** 概念。

```bash
# 使用 go fix 自动升级代码
go fix ./...

# 查看可用的 modernizers
go fix -list
```

**自动修复示例:**

```go
// 修复前 (Go 1.20 风格)
for i := 0; i < len(items); i++ {
    item := items[i]
    process(item)
}

// go fix 自动转换为 (Go 1.22+ 风格)
for i, item := range items {
    process(item)
}
```

### 1.3.2 内联分析器

```go
// 使用 //go:fix inline 指令标记可内联函数

//go:fix inline
func helper(x int) int {
    return x * 2
}

func main() {
    // go fix 会自动内联所有对 helper 的调用
    result := helper(5)  // 变成: result := 5 * 2
}
```

---

## 1.4 运行时增强

### 1.4.1 Goroutine 泄漏检测（实验性）

```go
package main

import (
    "context"
    "time"
)

// 这是一个典型的 goroutine 泄漏模式
func leakyProcess(ctx context.Context) {
    ch := make(chan result)

    // 启动 goroutine
    go func() {
        res := doWork()
        ch <- res  // 如果接收者提前返回，这个发送会永远阻塞
    }()

    select {
    case res := <-ch:
        useResult(res)
    case <-ctx.Done():
        return  // 泄漏！goroutine 仍然在等待发送
    }
}

// 修复：使用缓冲通道或确保接收
type result struct{}

func doWork() result { return result{} }
func useResult(r result) {}

func fixedProcess(ctx context.Context) {
    ch := make(chan result, 1)  // 缓冲通道

    go func() {
        res := doWork()
        select {
        case ch <- res:
        case <-ctx.Done():
            // 优雅退出
        }
    }()

    select {
    case res := <-ch:
        useResult(res)
    case <-ctx.Done():
        return
    }
}
```

**启用泄漏检测:**

```bash
GOEXPERIMENT=goroutineleakprofile go run main.go

# 查看泄漏 profile
go tool pprof http://localhost:6060/debug/pprof/goroutineleak
```

### 1.4.2 新的调度器指标

```go
package main

import (
    "fmt"
    "runtime/metrics"
)

func main() {
    // 读取新的 goroutine 状态指标
    samples := []metrics.Sample{
        {Name: "/sched/goroutines:goroutines"},
        {Name: "/sched/goroutines/runnable:goroutines"},
        {Name: "/sched/goroutines/running:goroutines"},
        {Name: "/sched/goroutines/waiting:goroutines"},
        {Name: "/sched/goroutines/not-in-go:goroutines"},
        {Name: "/sched/goroutines-created:goroutines"},
        {Name: "/sched/threads:threads"},
    }

    metrics.Read(samples)

    for _, s := range samples {
        fmt.Printf("%s: %d\n", s.Name, s.Value.Uint64())
    }
}
```

---

## 1.5 标准库更新

### 1.5.1 `errors.AsType` - 泛型错误检查

```go
package main

import (
    "errors"
    "fmt"
    "net"
)

func main() {
    err := &net.DNSError{Name: "example.com"}

    // 旧方式（Go 1.13+）
    var target *net.DNSError
    if errors.As(err, &target) {
        fmt.Println("DNS error:", target.Name)
    }

    // 新方式（Go 1.26）- 类型安全、更快、更简洁
    if dnsErr, ok := errors.AsType[*net.DNSError](err); ok {
        fmt.Println("DNS error:", dnsErr.Name)
    }

    // 多类型检查更清晰
    if connErr, ok := errors.AsType[*net.OpError](err); ok {
        fmt.Println("Network op failed:", connErr.Op)
    } else if dnsErr, ok := errors.AsType[*net.DNSError](err); ok {
        fmt.Println("DNS resolution failed:", dnsErr.Name)
    }
}
```

**性能对比:**

```text
BenchmarkAs-8        12606744    95.62 ns/op    40 B/op    2 allocs/op
BenchmarkAsType-8    37961869    30.26 ns/op    24 B/op    1 allocs/op
```

### 1.5.2 `crypto/hpke` - 混合公钥加密

```go
package main

import (
    "fmt"
    "crypto/rand"

    "golang.org/x/crypto/hpke"
)

func main() {
    // 配置: ML-KEM + X25519, HKDF-SHA256, AES-256-GCM
    kem, kdf, aead := hpke.MLKEM768X25519(), hpke.HKDFSHA256(), hpke.AES256GCM()

    // 接收方生成密钥对
    recipientPriv, err := kem.GenerateKey()
    if err != nil {
        panic(err)
    }
    recipientPub := recipientPriv.PublicKey()

    // 发送方加密消息
    message := []byte("secret message")
    ciphertext, err := hpke.Seal(recipientPub, kdf, aead, []byte("context"), message)
    if err != nil {
        panic(err)
    }

    // 接收方解密消息
    plaintext, err := hpke.Open(recipientPriv, kdf, aead, []byte("context"), ciphertext)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Decrypted: %s\n", plaintext)
}
```

### 1.5.3 `signal.NotifyContext` 改进

```go
package main

import (
    "context"
    "fmt"
    "os"
    "syscall"
)

func main() {
    // Go 1.26: context 取消时包含信号信息
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    <-ctx.Done()

    // Go 1.25: 只能看到 "context canceled"
    // Go 1.26: 可以看到具体是哪个信号
    fmt.Println("Error:", ctx.Err())
    fmt.Println("Cause:", context.Cause(ctx))  // 输出: signal: interrupt
}
```

---

## 1.6 实验性特性

### 1.6.1 SIMD 支持 (`simd/archsimd`)

```go
// +build goexperiment.simd

package main

import (
    "fmt"
    "golang.org/x/exp/simd/archsimd"
)

func main() {
    // 128-bit 向量操作
    a := archsimd.Int8x16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
    b := archsimd.Int8x16{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

    c := a.Add(b)
    fmt.Println(c) // [17 17 17 17 17 17 17 17 17 17 17 17 17 17 17 17]
}
```

### 1.6.2 安全内存擦除 (`runtime/secret`)

```go
// +build goexperiment.runtimesecret

package main

import (
    "crypto/ecdh"
    "crypto/rand"
    "golang.org/x/exp/runtime/secret"
)

func deriveSessionKey(peerPub *ecdh.PublicKey) (*ecdh.PublicKey, []byte, error) {
    var pubKey *ecdh.PublicKey
    var sessionKey []byte
    var err error

    secret.Do(func() {
        // 在这个块内分配的敏感数据会在结束时自动擦除
        priv, e := ecdh.P256().GenerateKey(rand.Reader)
        if e != nil {
            err = e
            return
        }

        shared, e := priv.ECDH(peerPub)
        if e != nil {
            err = e
            return
        }

        // 派生会话密钥
        sessionKey = hkdfDerive(shared)
        pubKey = priv.PublicKey()
    })

    return pubKey, sessionKey, err
}

func hkdfDerive(secret []byte) []byte {
    // HKDF 实现...
    return secret[:32] // 简化示例
}
```

---

## 1.7 迁移指南

### 从 Go 1.25 迁移到 1.26

```bash
# 1. 升级 Go 版本
go mod edit -go=1.26

# 2. 运行 go fix 自动更新代码
go fix ./...

# 3. 运行测试
go test ./...

# 4. 检查弃用警告
go vet ./...
```

### 兼容性注意事项

| 特性 | 兼容性影响 | 建议操作 |
|------|-----------|----------|
| Green Tea GC | 无 | 享受性能提升 |
| `new(表达式)` | 完全向后兼容 | 逐步采用新语法 |
| 递归类型约束 | 需要 Go 1.26 | 检查 go.mod 版本 |
| `errors.AsType` | 需要 Go 1.26 | 条件编译或版本检查 |

---

## 1.8 小结

Go 1.26 带来了多项重大改进：

1. **语言特性**: `new(表达式)` 和递归类型约束让代码更简洁表达力更强
2. **性能**: Green Tea GC 和 cgo 优化显著提升运行时性能
3. **工具**: 重写的 `go fix` 使代码现代化更加容易
4. **运行时**: Goroutine 泄漏检测帮助发现并发问题
5. **标准库**: `errors.AsType` 和 `crypto/hpke` 等新增功能
