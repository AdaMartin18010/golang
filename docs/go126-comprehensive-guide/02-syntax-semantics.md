# 第二章：语法与语义完整参考

> Go 1.26 语言语法、语义规则和形式化描述

---

## 2.1 词法元素 (Lexical Elements)

### 2.1.1 标识符与关键字

```go
// 标识符规则
identifier = letter { letter | unicode_digit }

// 关键字（保留字）
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```

### 2.1.2 运算符与优先级

```text
优先级    运算符
────────────────────────────────────────
5         *  /  %  <<  >>  &  &^
4         +  -  |  ^
3         ==  !=  <  <=  >  >=
2         &&
1         ||
```

### 2.1.3 字面量

```go
// 整数字面量
42          // 十进制
0b101010    // 二进制 (Go 1.13+)
0o52        // 八进制 (Go 1.13+)
0x2A        // 十六进制

// 浮点数字面量
3.14
1e10
0x1.ffp+8   // 十六进制浮点数 (Go 1.13+)

// 复数字面量
1i
3.14i

// 符文字面量
'a'
'\n'
'\u4E2D'    // Unicode

// 字符串字面量
" interpreted \\nstring"
` raw
multiline
string `
```

---

## 2.2 类型系统语法

### 2.2.1 类型声明

```go
// 类型定义 - 创建新类型
type MyInt int
type Point struct { X, Y float64 }
type Predicate func(T) bool

// 类型别名 - 等价类型
type Integer = int

// 泛型类型 (Go 1.18+)
type Stack[T any] struct {
    items []T
}

type Number interface {
    ~int | ~int64 | ~float64
}
```

### 2.2.2 接口定义

```go
// 基本接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

// 接口组合
type ReadWriter interface {
    Reader
    Writer
}

// 泛型接口
type Comparable[T any] interface {
    Compare(other T) int
}

// 类型集合接口
type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

---

## 2.3 声明与作用域

### 2.3.1 变量声明

```go
// 短变量声明（函数内）
x := 42

// 完整变量声明
var x int = 42
var x, y int = 1, 2
var (
    a = 1
    b = "hello"
    c float64
)

// 常量声明
const Pi = 3.14159
const (
    Monday = iota
    Tuesday
    Wednesday
)

// 组声明
const (
    Size = 1024
    MaxSize = Size * 4
)
```

### 2.3.2 作用域规则

```text
作用域层级（从内到外）：
1. 预声明标识符（universe block）
2. 包作用域（package block）
3. 文件作用域（file block）- import
4. 函数作用域（function block）
5. 语句作用域（statement block）- if, for, switch
6. 子语句作用域 - case, clause
```

```go
// 作用域遮蔽示例
package main

import "fmt"

var x = "package"

func main() {
    fmt.Println(x) // "package"

    x := "function" // 遮蔽包级x
    fmt.Println(x)  // "function"

    if true {
        x := "block" // 遮蔽函数级x
        fmt.Println(x) // "block"
    }

    fmt.Println(x) // "function"
}
```

---

## 2.4 控制流语句

### 2.4.1 条件语句

```go
// if 语句
if x > 0 {
    return x
} else if x < 0 {
    return -x
} else {
    return 0
}

// if 带短变量声明
if err := doSomething(); err != nil {
    return err
}

// switch 语句
switch day {
case Monday, Tuesday, Wednesday, Thursday, Friday:
    fmt.Println("Weekday")
case Saturday, Sunday:
    fmt.Println("Weekend")
default:
    fmt.Println("Invalid")
}

// 无表达式 switch（替代 if-else 链）
switch {
case x < 0:
    return -1
case x == 0:
    return 0
default:
    return 1
}

// type switch
switch v := i.(type) {
case int:
    fmt.Printf("Integer: %d\n", v)
case string:
    fmt.Printf("String: %s\n", v)
default:
    fmt.Printf("Unknown type: %T\n", v)
}
```

### 2.4.2 循环语句

```go
// 标准 for 循环
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// 条件循环（while 风格）
for x < 100 {
    x *= 2
}

// 无限循环
for {
    // 等价于 for true {}
}

// range 循环
for i, v := range slice {
    // i = 索引, v = 值
}

for k, v := range map {
    // k = 键, v = 值
}

for i := range slice {
    // 只获取索引
}

for range slice {
    // 只执行 len(slice) 次
}

// Go 1.22: 整数 range
for i := range 10 {  // 0 到 9
    fmt.Println(i)
}
```

### 2.4.3 跳转语句

```go
// break
for i := 0; i < 10; i++ {
    if i == 5 {
        break // 退出循环
    }
}

// continue
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue // 跳过偶数
    }
}

// goto（慎用）
for i := 0; i < 10; i++ {
    if i == 5 {
        goto done
    }
}
done:
    fmt.Println("Done")

// return
func sum(a, b int) int {
    return a + b
}

// defer
func readFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close() // 函数返回时执行

    // 处理文件...
    return nil
}
```

---

## 2.5 函数语义

### 2.5.1 函数声明

```go
// 基本函数
func add(a, b int) int {
    return a + b
}

// 多返回值
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// 命名返回值
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // 裸返回
}

// 变参函数
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}

// 泛型函数
func Max[T comparable](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 方法
type Vertex struct {
    X, Y float64
}

func (v Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
    v.X *= f
    v.Y *= f
}
```

### 2.5.2 函数值与闭包

```go
// 函数类型
func apply(f func(int) int, x int) int {
    return f(x)
}

// 匿名函数
func main() {
    square := func(x int) int {
        return x * x
    }
    fmt.Println(apply(square, 5)) // 25
}

// 闭包
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    counter := makeCounter()
    fmt.Println(counter()) // 1
    fmt.Println(counter()) // 2
    fmt.Println(counter()) // 3
}
```

---

## 2.6 复合类型语义

### 2.6.1 数组与切片

```go
// 数组：固定长度，值类型
var a [5]int          // [0 0 0 0 0]
b := [3]int{1, 2, 3}  // [1 2 3]
c := [...]int{1, 2}   // 编译器推导长度 [1 2]

// 切片：动态长度，引用类型
var s []int           // nil 切片
s := []int{1, 2, 3}   // 字面量
s := make([]int, 5)   // [0 0 0 0 0]，len=5, cap=5
s := make([]int, 3, 10) // len=3, cap=10

// 切片操作
s := []int{0, 1, 2, 3, 4, 5}
s[1:4]    // [1 2 3]（左闭右开）
s[:3]     // [0 1 2]
s[3:]     // [3 4 5]
s[:]      // 完整切片

// append
s = append(s, 6)
s = append(s, 7, 8, 9)
s = append(s, anotherSlice...)

// copy
dst := make([]int, len(src))
copy(dst, src)
```

### 2.6.2 Map

```go
// 声明与创建
var m map[string]int  // nil map
m := make(map[string]int)
m := map[string]int{
    "alice": 25,
    "bob":   30,
}

// 操作
m["charlie"] = 35     // 插入/更新
age := m["alice"]     // 查找（key不存在返回零值）
age, ok := m["dave"] // 查找，ok=false 表示不存在
delete(m, "bob")      // 删除

// 遍历
for k, v := range m {
    fmt.Printf("%s: %d\n", k, v)
}
```

### 2.6.3 结构体

```go
// 定义
type Person struct {
    Name    string
    Age     int
    private string // 小写：包私有
}

// 创建
p := Person{Name: "Alice", Age: 25}
p := Person{"Alice", 25} // 位置参数（不推荐）
p := new(Person)         // 返回 *Person，零值初始化

// 匿名字段（嵌入）
type Employee struct {
    Person    // 嵌入，继承方法
    Salary    int
}

e := Employee{
    Person: Person{Name: "Bob", Age: 30},
    Salary: 5000,
}
fmt.Println(e.Name) // 直接访问嵌入字段
```

---

## 2.7 指针语义

```go
// 指针声明
var p *int
i := 42
p = &i

// 解引用
fmt.Println(*p) // 读取
*p = 21         // 写入

// 与 C/C++ 的区别：
// 1. 无指针运算
// 2. 不能获取临时值的地址
// 3. nil 指针可安全比较，但解引用会 panic

// new 函数
p := new(int)    // *int，初始化为 0
s := new(string) // *string，初始化为 ""

// Go 1.26: new 支持表达式
ptr := new(yearsSince(born)) // *int，初始化为表达式结果
```

---

## 2.8 接口语义

### 2.8.1 隐式实现

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}

// 无需显式声明，自动实现
type MyWriter struct{}

func (m MyWriter) Write(p []byte) (n int, err error) {
    return len(p), nil
}

var _ Writer = MyWriter{} // 编译期检查
```

### 2.8.2 空接口与类型断言

```go
// 空接口：可接受任何类型
var i interface{} = "hello"

// 类型断言
s := i.(string)           // 断言为 string，失败 panic
s, ok := i.(string)       // 安全断言，ok=false 表示失败

// 类型 switch
switch v := i.(type) {
case int:
    fmt.Printf("int: %d\n", v)
case string:
    fmt.Printf("string: %s\n", v)
case nil:
    fmt.Println("nil")
default:
    fmt.Printf("unknown: %T\n", v)
}
```

---

## 2.9 并发语法

### 2.9.1 Goroutine

```go
// 启动 goroutine
go function()
go func() { ... }()
go func(param) { ... }(arg)

// 注意：goroutine 在函数返回后继续运行
func main() {
    go func() {
        time.Sleep(time.Second)
        fmt.Println("Hello from goroutine")
    }()
    // 如果 main 立即退出，goroutine 不会执行完
    time.Sleep(2 * time.Second)
}
```

### 2.9.2 Channel

```go
// 创建
ch := make(chan int)      // 无缓冲（同步）
ch := make(chan int, 10)  // 有缓冲（异步）

// 发送
ch <- v

// 接收
v := <-ch
v, ok := <-ch // ok=false 表示 channel 已关闭

// 关闭
close(ch)

// 只读/只写通道
func producer() <-chan int { ... }
func consumer(ch <-chan int) { ... }
func sender(ch chan<- int) { ... }

// select
select {
case v1 := <-ch1:
    fmt.Println("ch1:", v1)
case v2 := <-ch2:
    fmt.Println("ch2:", v2)
case ch3 <- 100:
    fmt.Println("sent to ch3")
default:
    fmt.Println("no channel ready")
}

// 超时模式
select {
case result := <-ch:
    return result
case <-time.After(5 * time.Second):
    return nil, errors.New("timeout")
}
```

---

## 2.10 语义规则总结

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go 语义核心规则                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. 零值初始化                                                              │
│     - 所有变量初始化为其类型的零值                                            │
│     - int: 0, string: "", bool: false, pointer: nil                         │
│     - slice/map/channel: nil                                                │
│     - interface: nil                                                        │
│                                                                             │
│  2. 值语义 vs 引用语义                                                      │
│     - 值类型: int, float, bool, string, array, struct                       │
│     - 引用类型: slice, map, channel, pointer, function, interface           │
│                                                                             │
│  3. 可见性规则                                                              │
│     - 大写开头: 包外可见 (public)                                           │
│     - 小写开头: 包内可见 (private)                                          │
│                                                                             │
│  4. 接口动态分派                                                            │
│     - 运行时确定具体实现                                                    │
│     - 通过 itable 实现高效调用                                              │
│                                                                             │
│  5. 并发同步点                                                              │
│     - channel 发送/接收是同步点                                             │
│     - select case 随机选择就绪 channel                                      │
│     - 无缓冲 channel 保证 happens-before                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

*本章提供了 Go 语言完整的语法和语义参考，是编写正确 Go 代码的基础。*
