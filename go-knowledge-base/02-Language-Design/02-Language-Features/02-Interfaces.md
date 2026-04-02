# Go 接口 (Interfaces)

> **维度**: 语言设计 (Language Design)
> **分类**: 类型系统核心
> **难度**: 进阶
> **Go 版本**: Go 1.0+ (泛型支持 Go 1.18+)
> **最后更新**: 2026-04-02

---

## 1. 问题陈述 (Problem Statement)

### 1.1 面向对象中的继承问题

传统面向对象语言使用显式继承，面临以下挑战：

| 挑战 | 描述 | Go 的解决方式 |
|------|------|---------------|
| **紧耦合** | 子类与父类强绑定 | 隐式实现，解耦合 |
| **脆弱基类** | 父类修改影响所有子类 | 组合优于继承 |
| **钻石继承** | 多重继承的二义性 | 无继承，只有组合 |
| **模拟类** | 需要模拟不相关的类 | 小接口，按需实现 |

### 1.2 接口设计目标

```
Go 接口设计哲学:
┌─────────────────────────────────────────────────────────┐
│  1. 隐式实现 (Implicit Implementation)                  │
│     → 无需显式声明实现关系                              │
├─────────────────────────────────────────────────────────┤
│  2. 鸭子类型 (Duck Typing)                              │
│     → "如果它走起来像鸭子，叫起来像鸭子，那就是鸭子"    │
├─────────────────────────────────────────────────────────┤
│  3. 小接口原则 (Small Interfaces)                       │
│     → 接口应该小而专注，便于组合                        │
├─────────────────────────────────────────────────────────┤
│  4. 行为抽象 (Behavior Abstraction)                     │
│     → 关注能做什么，而不是是什么                        │
└─────────────────────────────────────────────────────────┘
```

---

## 2. 形式化方法 (Formal Approach)

### 2.1 结构子类型 (Structural Subtyping)

Go 使用结构子类型而非名义子类型：

```
名义子类型 (Nominal):
  class Dog extends Animal  // 显式声明

结构子类型 (Structural):
  type Animal interface { Speak() }
  type Dog struct {}
  func (d Dog) Speak() {}   // 隐式实现 Animal

形式化定义:
  类型 T 实现接口 I，当且仅当 T 的方法集包含 I 的所有方法。

  implements(T, I) ↔ ∀m ∈ I.Methods, m ∈ T.Methods
```

### 2.2 接口值的表示

```
接口值在运行时的内存布局:

interface {
    tab  *itab          // 类型描述符和方法表
    data unsafe.Pointer // 指向实际数据
}

itab 结构:
{
    inter *interfacetype  // 接口类型元数据
    _type *_type          // 具体类型元数据
    hash  uint32          // 类型哈希，用于类型断言
    _     [4]byte         // 填充
    fun   [1]uintptr      // 方法表 (变长数组)
                          // fun[i] 是接口第 i 个方法的实际实现地址
}

空接口 (interface{}) 的特殊表示:
{
    _type *_type          // 类型元数据
    data  unsafe.Pointer  // 指向实际数据
}
```

### 2.3 方法集规则

```
方法集定义:

对于类型 T:
  - 值接收者方法: (t T) Method() → 属于 T 和 *T 的方法集
  - 指针接收者方法: (t *T) Method() → 仅属于 *T 的方法集

类型实现接口的条件:
  ┌──────────────────┬──────────────────┬──────────────────┐
  │   方法接收者     │   值类型 T       │   指针类型 *T    │
  ├──────────────────┼──────────────────┼──────────────────┤
  │   (t T)          │       ✓          │       ✓          │
  │   (t *T)         │       ✗          │       ✓          │
  └──────────────────┴──────────────────┴──────────────────┘

示例:
  type Reader interface { Read([]byte) (int, error) }

  type File struct{}
  func (f File) Read([]byte) (int, error) { ... }
  // File 和 *File 都实现了 Reader

  func (f *File) Write([]byte) (int, error) { ... }
  // 只有 *File 实现了 Writer (假设 Writer 需要 Write 方法)
```

---

## 3. 实现细节 (Implementation)

### 3.1 接口定义与实现

```go
package main

import (
    "fmt"
    "io"
)

// 定义接口
// 小接口：只包含一个方法
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 接口组合
type ReadWriter interface {
    Reader
    Writer
}

// 实现类型
type Buffer struct {
    data []byte
    pos  int
}

// 值接收者实现
func (b Buffer) Read(p []byte) (n int, err error) {
    n = copy(p, b.data[b.pos:])
    b.pos += n
    return n, nil
}

// 指针接收者实现
func (b *Buffer) Write(p []byte) (n int, err error) {
    b.data = append(b.data, p...)
    return len(p), nil
}

func main() {
    // Buffer 未实现 Write，所以不能赋值给 ReadWriter
    // var rw ReadWriter = Buffer{} // 编译错误

    // *Buffer 实现了 Read 和 Write
    var rw ReadWriter = &Buffer{}

    // 但 *Buffer 可以读取 Buffer 的 Read 方法
    // (因为 Read 使用值接收者)
}
```

### 3.2 空接口与类型断言

```go
package main

import "fmt"

// 空接口可以保存任何类型的值
func printAny(v interface{}) {
    fmt.Printf("value: %v, type: %T\n", v, v)
}

// 类型断言
func processValue(v interface{}) {
    // 类型开关
    switch x := v.(type) {
    case int:
        fmt.Printf("Integer: %d\n", x)
    case string:
        fmt.Printf("String: %s\n", x)
    case []byte:
        fmt.Printf("Byte slice with length: %d\n", len(x))
    default:
        fmt.Printf("Unknown type: %T\n", x)
    }
}

// 类型断言 with ok
func getInt(v interface{}) (int, bool) {
    i, ok := v.(int)
    return i, ok
}

// panic 风险
func mustInt(v interface{}) int {
    return v.(int) // 如果 v 不是 int，会 panic
}

func main() {
    printAny(42)           // value: 42, type: int
    printAny("hello")      // value: hello, type: string
    printAny([]int{1,2,3}) // value: [1 2 3], type: []int

    processValue(100)
    processValue("world")
}
```

### 3.3 泛型与接口 (Go 1.18+)

```go
package main

// 泛型约束使用接口
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64
}

// 泛型函数
func Max[T Number](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 近似约束 (~)
type MyInt int  // MyInt 的底层类型是 int

// 使用 ~int 允许 MyInt 满足 Number 约束

// 接口作为约束
type Stringer interface {
    String() string
}

func ToStringSlice[T Stringer](items []T) []string {
    result := make([]string, len(items))
    for i, item := range items {
        result[i] = item.String()
    }
    return result
}

// 类型集合
type Ordered interface {
    Integer | Float | ~string
}
```

### 3.4 常见接口模式

```go
package main

import "io"

// 1. Reader/Writer 模式 (标准库)
// io.Reader, io.Writer, io.Closer, io.ReadWriter

// 2. Stringer 接口 (fmt 包)
type Stringer interface {
    String() string
}

// 3. Error 接口
// type error interface {
//     Error() string
// }

// 4. 比较接口
type Comparable interface {
    Compare(other interface{}) int
}

// 5. 选项模式 (Functional Options)
type ServerOption func(*Server)

type Server struct {
    addr     string
    timeout  int
    maxConns int
}

func WithTimeout(timeout int) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func WithMaxConns(max int) ServerOption {
    return func(s *Server) {
        s.maxConns = max
    }
}

func NewServer(addr string, opts ...ServerOption) *Server {
    s := &Server{
        addr:     addr,
        timeout:  30,
        maxConns: 100,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// 使用
// srv := NewServer(":8080", WithTimeout(60), WithMaxConns(200))
```

---

## 4. 语义分析 (Semantic Analysis)

### 4.1 动态分派语义

```
接口调用是动态分派:

var r Reader = &Buffer{}
n, err := r.Read(p)

编译期:
  - 检查 r 是否实现了 Reader (静态检查)
  - 生成接口调用代码

运行期:
  1. 从 r.tab.fun 找到 Read 方法地址
  2. 将 r.data 作为接收者
  3. 调用实际方法

性能考虑:
  - 接口调用有间接开销 (~2-3ns)
  - 但支持内联优化 (去虚拟化)
```

### 4.2 nil 接口的陷阱

```go
package main

import "fmt"

type MyError struct {
    msg string
}

func (e *MyError) Error() string {
    return e.msg
}

func mayFail() error {
    var e *MyError = nil  // nil 指针
    return e              // 返回 nil 吗？
}

func main() {
    err := mayFail()

    // 陷阱: err != nil !
    // 因为 err 是接口值，包含 (type=*MyError, data=nil)
    // 接口值本身不是 nil

    fmt.Println(err == nil) // false!

    // 正确做法
    var e *MyError = nil
    if e != nil {
        return e
    }
    return nil  // 返回真正的 nil 接口
}
```

---

## 5. 权衡分析 (Trade-offs)

### 5.1 接口 vs 泛型

| 场景 | 接口 | 泛型 | 推荐 |
|------|------|------|------|
| 不同数据，相同算法 | 使用 | 使用 | 泛型 (无运行时开销) |
| 运行时类型未知 | 必须使用 | 无法使用 | 接口 |
| 需要方法集 | 使用 | 配合约束使用 | 两者皆可 |
| 性能敏感 | 有间接开销 | 零开销 | 泛型 |
| 代码简洁性 | 简洁 | 类型参数复杂 | 接口 |

### 5.2 大接口 vs 小接口

```
小接口 (Go 推荐):
  优势:
    ✓ 易于实现 (方法少)
    ✓ 便于组合
    ✓ 职责清晰

  示例: io.Reader, io.Writer, io.Closer

大接口 (反模式):
  劣势:
    ✗ 实现困难
    ✗ 违反 ISP (接口隔离原则)
    ✗ 修改影响面广

  示例: 包含 Read/Write/Seek/Close/Flush 的大接口

折中: io.ReadWriteCloser (组合小接口)
```

---

## 6. 视觉表示 (Visual Representations)

### 6.1 接口值内存布局

```
非空接口值:
┌─────────────────────────────────────────────────────────────┐
│                     Interface Value                         │
├────────────────────────┬────────────────────────────────────┤
│        tab             │               data                 │
│    (*itab, 8 bytes)    │       (unsafe.Pointer, 8 bytes)    │
├────────┬───────────────┼────────────────────────────────────┤
│        │               │                                    │
│        ▼               ▼                                    │
│   ┌──────────────┐   ┌──────────────┐                      │
│   │     itab     │   │   实际数据    │                      │
│   ├──────────────┤   │   (堆分配)    │                      │
│   │ inter        │   └──────────────┘                      │
│   │ _type        │                                          │
│   │ hash         │                                          │
│   │ fun[0]       │ ──► Read 方法地址                        │
│   │ fun[1]       │ ──► Write 方法地址                       │
│   └──────────────┘                                          │
└─────────────────────────────────────────────────────────────┘

空接口值 (interface{}):
┌─────────────────────────────────────────────────────────────┐
│                     Empty Interface                         │
├────────────────────────┬────────────────────────────────────┤
│        _type           │               data                 │
│    (*_type, 8 bytes)   │       (unsafe.Pointer, 8 bytes)    │
├────────────────────────┴────────────────────────────────────┤
│                                                             │
│  _type 指向类型描述符，包含类型大小、对齐、方法集等信息       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 6.2 类型实现关系

```
接口实现关系图:

        ┌──────────────┐
        │  io.Reader   │
        │  Read()      │
        └──────┬───────┘
               │
        ┌──────┴───────┐
        │  io.Writer   │
        │  Write()     │
        └──────┬───────┘
               │
               ▼
        ┌──────────────┐
        │io.ReadWriter │◄────┐
        │ (组合接口)    │     │
        └──────────────┘     │
                             │
    ┌────────────────────────┼────────────────────────┐
    │                        │                        │
    ▼                        ▼                        ▼
┌─────────┐           ┌─────────┐           ┌──────────┐
│  File   │           │ Buffer  │           │  Pipe    │
│         │           │         │           │          │
│ Read()  │           │ Read()  │           │ Read()   │
│ Write() │           │ Write() │           │ Write()  │
└─────────┘           └─────────┘           └──────────┘
```

---

## 7. 最佳实践

### 7.1 接口设计原则

```go
// 1. 小接口原则
// 好：小而专注
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 不好：大而全
type ReadWriteSeekCloserFlusher interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
    Seek(offset int64, whence int) (int64, error)
    Close() error
    Flush() error
}

// 2. 在消费者侧定义接口
// 包 A (消费者)
package storage

type FileReader interface {
    ReadFile(name string) ([]byte, error)
}

func LoadConfig(r FileReader, name string) ([]byte, error) {
    return r.ReadFile(name)
}

// 包 B (生产者)
package os

func ReadFile(name string) ([]byte, error) { ... }

// os 自动满足 storage.FileReader，无需依赖

// 3. 接受接口，返回具体类型
func Process(r io.Reader) error  // 接受接口，灵活

func NewBuffer() *Buffer          // 返回具体类型，明确
```

### 7.2 常见陷阱

```go
// 1. nil 接口陷阱
func returnsNilError() error {
    var p *MyError = nil
    return p  // 不是 nil 接口！
}

// 2. 值接收者 vs 指针接收者
func (t T) Method() {}  // T 和 *T 都有此方法
func (t *T) Method() {} // 只有 *T 有

// 3. 接口嵌套循环
// type A interface { B }
// type B interface { A } // 编译错误：循环嵌入
```

---

## 8. 相关资源

- [LD-016-Interface-Internals.md](./LD-016-Interface-Internals.md) - 接口内部实现
- [LD-007-Go-Reflection-Interface-Internals.md](../LD-007-Go-Reflection-Interface-Internals.md)
- [io 包文档](https://pkg.go.dev/io)

---

*S-Level Quality Document | Generated: 2026-04-02*
