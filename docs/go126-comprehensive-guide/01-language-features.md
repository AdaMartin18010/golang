# Go 1.26 语言基础

> 从类型论与计算理论视角理解Go语言核心

---

## 一、计算模型基础

### 1.1 命令式与函数式融合

```
Go的计算模型:
────────────────────────────────────────
Go = 命令式核心 + 函数式特性 + 并发原语

命令式基础:
├─ 可变状态 (Variables)
├─ 控制流 (if/for/switch)
├─ 过程抽象 (Functions)
└─ 结构化数据 (Structs)

函数式特性:
├─ 一等函数 (First-class functions)
├─ 闭包 (Closures)
├─ 函数组合 (Composition)
└─ 不可变值传递

并发支持:
├─ Goroutines (轻量级线程)
├─ Channels (通信同步)
└─ Select (非确定性选择)

理论对应:
├─ 命令式: Turing Machine
├─ 函数式: Lambda Calculus
└─ 并发: CSP (Communicating Sequential Processes)
```

### 1.2 值语义与引用语义

```
值语义的形式化:
────────────────────────────────────────
赋值: x = y 创建y的副本
相等: x == y 比较值内容

Go的值类型:
├─ 基本类型: int, float, bool, string
├─ 数组: [N]T (值语义)
├─ 结构体: struct (值语义)
└─ 函数: func (闭包引用捕获)

引用类型:
├─ 切片: []T (header值 + 底层数组引用)
├─ Map: map[K]V (引用哈希表)
├─ Channel: chan T (引用通信通道)
├─ 接口: interface (类型对 + 值对)
└─ 指针: *T (内存地址引用)

内存模型:
值类型 ──► 栈分配 (通常)
引用类型 ──► 堆分配
逃逸分析决定实际分配位置

代码示例:
// 值类型: 赋值创建副本
func valueSemantics() {
    a := [3]int{1, 2, 3}
    b := a  // 创建完整副本
    b[0] = 99

    fmt.Println(a) // [1 2 3] - 原数组不变
    fmt.Println(b) // [99 2 3] - 副本修改
}

// 引用类型: 赋值共享底层数据
func referenceSemantics() {
    s1 := []int{1, 2, 3}
    s2 := s1  // 共享底层数组
    s2[0] = 99

    fmt.Println(s1) // [99 2 3] - 原切片也被修改!
    fmt.Println(s2) // [99 2 3]
}

// 结构体值语义
func structSemantics() {
    type Point struct{ X, Y int }

    p1 := Point{1, 2}
    p2 := p1  // 副本
    p2.X = 99

    fmt.Println(p1) // {1 2} - 不变
    fmt.Println(p2) // {99 2}
}
```

### 1.3 内存布局与对齐

```
Go内存布局:
────────────────────────────────────────

基本类型大小:
├─ bool: 1 byte
├─ int/uint: 4 or 8 bytes (32/64位)
├─ float32: 4 bytes
├─ float64: 8 bytes
├─ string: 16 bytes (指针+长度)
├─ slice: 24 bytes (指针+长度+容量)
├─ interface: 16 bytes (类型+值)
└─ map/channel/func: 8 bytes (指针)

结构体内存对齐:
├─ 字段按最大对齐要求对齐
├─ 填充优化字段顺序
└─ unsafe.Sizeof 查看实际大小

代码示例:
// 内存对齐示例
func memoryAlignment() {
    type BadLayout struct {
        A bool   // 1 byte + 7 padding
        B int64  // 8 bytes
        C bool   // 1 byte + 7 padding
    }  // 总: 24 bytes

    type GoodLayout struct {
        B int64  // 8 bytes
        A bool   // 1 byte
        C bool   // 1 byte + 6 padding
    }  // 总: 16 bytes

    fmt.Println(unsafe.Sizeof(BadLayout{}))  // 24
    fmt.Println(unsafe.Sizeof(GoodLayout{})) // 16
}

// 使用 unsafe 查看内存布局
func inspectMemory() {
    type MyStruct struct {
        A int32
        B int64
        C int16
    }

    s := MyStruct{A: 1, B: 2, C: 3}

    fmt.Printf("Size: %d\n", unsafe.Sizeof(s))
    fmt.Printf("Align: %d\n", unsafe.Alignof(s))
    fmt.Printf("A offset: %d\n", unsafe.Offsetof(s.A))
    fmt.Printf("B offset: %d\n", unsafe.Offsetof(s.B))
    fmt.Printf("C offset: %d\n", unsafe.Offsetof(s.C))
}
```

---

## 二、声明与作用域

### 2.1 声明的形式化

```
声明的静态语义:
────────────────────────────────────────
声明引入绑定: 标识符 → 实体

实体类别:
├─ 变量: var x T = expr  (可寻址存储)
├─ 常量: const x T = expr (编译期确定)
├─ 类型: type T = ...    (类型别名或定义)
├─ 函数: func f(...) ... (代码+环境)
└─ 包: package p         (命名空间)

作用域规则:
├─ 块作用域: { ... }
├─ 包作用域: 文件顶层
├─ 文件作用域: import
├─ 函数作用域: 参数+结果
└─ 预声明: 全局universe

遮蔽 (Shadowing):
内层声明遮蔽外层同名声明
Go 1.22修复: loop变量每次迭代新实例

代码示例:
// 作用域示例
func scopeExample() {
    x := "outer"
    {
        x := "inner"  // 遮蔽外层x
        fmt.Println(x) // "inner"
    }
    fmt.Println(x) // "outer" - 外层x不变
}

// 包级变量vs局部变量
var global = "package scope"

func variableScopes() {
    local := "function scope"

    if true {
        block := "block scope"
        fmt.Println(global)  // 访问包级变量
        fmt.Println(local)   // 访问函数级变量
        fmt.Println(block)   // 块级变量
    }

    // fmt.Println(block)  // 编译错误: undefined
}
```

### 2.2 短变量声明

```
短声明的形式化:
────────────────────────────────────────
x := expr  ≡  var x = expr

多变量:
x, y := expr1, expr2

重新声明规则:
├─ 至少一个变量是新声明的
├─ 已有变量在相同作用域
└─ 类型必须与已有变量一致

代码示例:
// 短声明基础
func shortDeclarations() {
    name := "Go"           // 推断为string
    version := 1.26        // 推断为float64

    // 多变量
    x, y := 1, 2
    fmt.Println(x, y)
}

// 重新声明规则
func redeclaration() {
    x := 1
    fmt.Println(x)  // 1

    // x已存在，但y是新变量 - 合法
    x, y := 2, 3
    fmt.Println(x, y)  // 2, 3

    // 编译错误: 没有新变量在左边
    // x, y := 4, 5
}

// if语句中的短声明
func ifShortDecl(x int) {
    if v := x * 2; v > 10 {
        fmt.Println("v > 10:", v)  // v在此块内可用
    } else {
        fmt.Println("v <= 10:", v) // v在else块也可用
    }
    // fmt.Println(v)  // 编译错误: v未定义
}
```

---

## 三、控制流

### 3.1 条件与循环

```
if的形式化:
────────────────────────────────────────
if cond { body } else { alt }
语义: cond为真执行body，否则执行alt

Go特性: 简短声明
if v := expr; cond(v) {
    // v在if-else块内可用
}

switch的完备性:
switch expr {
case v1: ...
case v2: ...
default: ...  // 可选但推荐
}

switch true {  // 可省略表达式
 case cond1: ...
 case cond2: ...
}

循环形式:
────────────────────────────────────────
for init; cond; post { ... }  // C风格
for cond { ... }              // while风格
for { ... }                   // 无限循环
for range expr { ... }        // 迭代

Go 1.22改进:
for i := range n { ... }  // range over int
loop变量每次迭代独立

代码示例:
// if 简短声明
func ifShortExample(filename string) error {
    if f, err := os.Open(filename); err != nil {
        return err
    } else {
        defer f.Close()
        // 使用f
    }
    return nil
}

// switch 多值匹配
func switchExample(day string) {
    switch day {
    case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
        fmt.Println("工作日")
    case "Saturday", "Sunday":
        fmt.Println("周末")
    default:
        fmt.Println("无效")
    }
}

// switch 条件表达式 (替代if-else链)
func grade(score int) string {
    switch {
    case score >= 90:
        return "A"
    case score >= 80:
        return "B"
    case score >= 70:
        return "C"
    case score >= 60:
        return "D"
    default:
        return "F"
    }
}

// range 多种用法
func rangeExamples() {
    // slice - 索引和值
    nums := []int{10, 20, 30}
    for i, v := range nums {
        fmt.Printf("[%d] = %d\n", i, v)
    }

    // 只取值 (使用空白标识符)
    for _, v := range nums {
        fmt.Println(v)
    }

    // 只取索引
    for i := range nums {
        fmt.Println(i)
    }

    // map - 键和值
    m := map[string]int{"a": 1, "b": 2}
    for k, v := range m {
        fmt.Printf("%s: %d\n", k, v)
    }

    // string - rune索引和值
    for i, r := range "Hello, 世界" {
        fmt.Printf("[%d] = %c\n", i, r)
    }

    // channel
    ch := make(chan int, 3)
    ch <- 1; ch <- 2; ch <- 3
    close(ch)
    for v := range ch {
        fmt.Println(v)
    }

    // Go 1.22: range over int
    for i := range 5 {
        fmt.Println(i)  // 0, 1, 2, 3, 4
    }
}
```

### 3.2 跳转语句

```
跳转的形式化语义:
────────────────────────────────────────

break:
├─ 终止最内层for/switch/select
├─ break Label: 终止指定外层
└─ 语义: 跳出到标签后

continue:
├─ 跳过当前迭代，进入下一次
├─ continue Label: 指定外层循环
└─ 语义: 跳到循环post或range

goto:
├─ 跳转到同一函数内标签
├─ 不能跳入变量作用域
└─ 语义: 控制流重定向

return:
├─ 函数返回
├─ 执行defer
└─ 可选带返回值 (命名返回值可省略)

panic/recover:
├─ panic: 异常传播，展开栈
├─ recover: defer中捕获panic
└─ 语义: 结构化异常处理

代码示例:
// break 带标签
func breakLabel() {
Outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i*j > 2 {
                break Outer  // 跳出外层循环
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    // 输出: (0,0) (0,1) (0,2) (1,0) (1,1) (1,2) (2,0) (2,1)
}

// continue 带标签
func continueLabel() {
Outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                continue Outer  // 跳到外层下一次迭代
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    // 输出: (0,0) (1,0) (2,0)
}

// panic/recover
func panicRecover() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()

    fmt.Println("Before panic")
    panic("something went wrong")
    fmt.Println("After panic - never reached")
}

// defer 堆栈行为
func deferStack() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    fmt.Println("done")
    // 输出: done, 3, 2, 1 (LIFO)
}
```

---

## 四、函数与闭包

### 4.1 函数的形式化

```
函数类型:
────────────────────────────────────────
FunctionType = func (Parameters) Results

参数传递:
├─ 值传递: 复制实参到形参
├─ 指针传递: 复制地址值
└─ 无引用传递 (不同于C++、Java)

多返回值:
func f() (int, error)
x, err := f()

命名返回值:
func f() (result int, err error) {
    result = 42
    return  // 裸return返回命名变量
}

变长参数:
func printf(format string, a ...interface{})
内部转换为slice

代码示例:
// 多返回值
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// 命名返回值
func getConfig() (host string, port int, err error) {
    host = "localhost"
    port = 8080
    // err = nil 隐式
    return  // 裸return，等同于 return host, port, err
}

// 变长参数
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 使用
fmt.Println(sum(1, 2, 3))     // 6
fmt.Println(sum())             // 0
nums := []int{1, 2, 3, 4}
fmt.Println(sum(nums...))      // 展开slice
```

### 4.2 闭包理论

```
闭包的形式化:
────────────────────────────────────────
闭包 = 函数代码 + 捕获环境

环境捕获:
├─ 值捕获: 复制变量值 (循环变量陷阱，已修复)
├─ 引用捕获: 捕获变量地址
└─ 逃逸分析: 决定分配位置

Go 1.22之前的问题:
funcs := []func(){}
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() { fmt.Println(i) })
}
// 全部输出3 (共享i)

Go 1.22修复:
// 每次迭代i是新变量
// 输出0,1,2

闭包应用:
├─ 装饰器模式
├─ 回调函数
├─ 延迟执行
└─ 函数工厂

代码示例:
// 闭包计数器
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

counter1 := makeCounter()
counter2 := makeCounter()

fmt.Println(counter1())  // 1
fmt.Println(counter1())  // 2
fmt.Println(counter2())  // 1 (独立的count)

// 闭包工厂
func makeMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

double := makeMultiplier(2)
triple := makeMultiplier(3)

fmt.Println(double(5))  // 10
fmt.Println(triple(5))  // 15

// Go 1.22修复: 循环变量
func loopVariableFixed() {
    funcs := []func(){}

    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() {
            fmt.Println(i)  // 每个i是独立的
        })
    }

    for _, f := range funcs {
        f()  // 输出 0, 1, 2
    }
}

// 闭包与并发
func closureConcurrency() {
    var wg sync.WaitGroup

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {  // 传递参数避免闭包问题
            defer wg.Done()
            fmt.Println("Worker", id)
        }(i)
    }

    wg.Wait()
}
```

---

## 五、复合类型详解

### 5.1 数组与切片

```
数组的形式化:
────────────────────────────────────────
Array = [N]T  // 固定长度序列
性质:
├─ 值语义 (赋值复制整个数组)
├─ 长度是类型的一部分
└─ 可比较 (元素可比较时)

切片的内部结构:
type SliceHeader struct {
    Data unsafe.Pointer  // 指向数组
    Len  int             // 长度
    Cap  int             // 容量
}

切片操作复杂度:
├─ s[i]: O(1)
├─ s[i:j]: O(1) (共享底层)
├─ append: 摊销O(1)，扩容时O(n)
└─ copy: O(min(len(dst), len(src)))

代码示例:
// 数组 vs 切片
func arrayVsSlice() {
    // 数组: 值类型，编译时确定长度
    var arr [3]int = [3]int{1, 2, 3}
    arr2 := arr  // 完整复制
    arr2[0] = 99
    fmt.Println(arr)  // [1 2 3] - 不变

    // 切片: 引用类型
    s := []int{1, 2, 3}
    s2 := s  // 共享底层
    s2[0] = 99
    fmt.Println(s)  // [99 2 3] - 被修改!
}

// 切片操作
func sliceOperations() {
    // make创建
    s := make([]int, 5, 10)  // len=5, cap=10

    // 字面量
    s2 := []int{1, 2, 3, 4, 5}

    // 切片操作 (共享底层!)
    sub := s2[1:3]  // [2, 3]
    sub[0] = 99
    fmt.Println(s2)  // [1 99 3 4 5] - 原切片被修改

    // append
    s3 := append(s2, 6)  // 添加元素

    // 扩容触发新底层数组
    s4 := make([]int, 3, 3)
    s5 := append(s4, 4)  // 容量不足，分配新数组
    s5[0] = 99
    fmt.Println(s4)  // [0 0 0] - 不变

    // copy
    dst := make([]int, 3)
    n := copy(dst, s2)  // 复制min(len(dst), len(src))个元素
    fmt.Println(n, dst)
}

// 常见陷阱
func sliceTraps() {
    // 陷阱1: 切片共享底层
    a := []int{1, 2, 3}
    b := a[:2]  // [1, 2]
    b = append(b, 99)
    fmt.Println(a)  // [1 2 99] - a也被修改!

    // 陷阱2: range的副本
    items := []struct{ Value int }{{1}, {2}, {3}}
    for _, item := range items {
        item.Value *= 2  // 修改的是副本!
    }
    fmt.Println(items)  // [{1} {2} {3}]

    // 正确做法
    for i := range items {
        items[i].Value *= 2  // 通过索引修改
    }
}
```

### 5.2 Map的实现原理

```
Map的形式化:
────────────────────────────────────────
Map = K → V  (部分函数，key存在性)

内部结构:
type hmap struct {
    count     int       // 元素数
    flags     uint8
    B         uint8     // log_2(桶数)
    noverflow uint16
    hash0     uint32    // 哈希种子
    buckets   unsafe.Pointer
    oldbuckets unsafe.Pointer
    nevacuate uintptr
}

哈希冲突解决:
├─ 链地址法: 每个桶8个槽位
├─ 溢出桶: 超过8个时链式溢出
└─ 渐进式扩容: 写入时迁移

使用约束:
├─ key必须可比
├─ 非并发安全 (需外部同步)
└─ 遍历顺序随机

代码示例:
// map基本操作
func mapOperations() {
    // 创建
    m := make(map[string]int)

    // 增/改
    m["one"] = 1
    m["two"] = 2

    // 查
    v, ok := m["one"]  // 1, true
    v2, ok2 := m["three"]  // 0, false

    // 删
    delete(m, "two")

    // 遍历 (顺序随机!)
    for k, v := range m {
        fmt.Printf("%s: %d\n", k, v)
    }

    // 长度
    fmt.Println(len(m))
}

// map并发安全
func mapConcurrency() {
    m := make(map[int]int)
    var mu sync.RWMutex

    // 写
    go func() {
        for i := 0; i < 100; i++ {
            mu.Lock()
            m[i] = i
            mu.Unlock()
        }
    }()

    // 读
    go func() {
        for i := 0; i < 100; i++ {
            mu.RLock()
            _ = m[i]
            mu.RUnlock()
        }
    }()
}

// sync.Map (并发安全)
func syncMapExample() {
    var m sync.Map

    m.Store("key", "value")

    if v, ok := m.Load("key"); ok {
        fmt.Println(v)
    }

    m.Delete("key")

    // 范围遍历
    m.Range(func(key, value interface{}) bool {
        fmt.Println(key, value)
        return true  // 继续遍历
    })
}

// map作为set
func mapAsSet() {
    set := make(map[string]struct{})

    // 添加
    set["a"] = struct{}{}
    set["b"] = struct{}{}

    // 检查存在
    if _, ok := set["a"]; ok {
        fmt.Println("a exists")
    }

    // 删除
    delete(set, "b")
}
```

### 5.3 结构体与嵌入

```
结构体的代数性质:
────────────────────────────────────────
Struct = Product Type (积类型)
struct { A int; B string } ≅ int × string

嵌入类型 (Embedding):
type Inner struct { X int }
type Outer struct {
    Inner      // 嵌入，继承方法集
    Y int
}

方法提升:
Outer自动获得Inner的方法
outer.X 等价于 outer.Inner.X

嵌入vs继承:
├─ Go: 组合 (has-a)
├─ OOP: 继承 (is-a)
└─ Go的方式更灵活，无层次耦合

代码示例:
// 结构体基础
func structBasic() {
    type Person struct {
        Name string
        Age  int
    }

    // 创建方式
    p1 := Person{"Alice", 30}
    p2 := Person{Name: "Bob", Age: 25}
    p3 := new(Person)  // 返回*Person
    p3.Name = "Charlie"

    // 匿名结构体
    point := struct {
        X, Y int
    }{10, 20}

    fmt.Println(p1, p2, point)
}

// 嵌入类型
func embeddingTypes() {
    type Animal struct {
        Name string
    }

    func (a *Animal) Speak() {
        fmt.Println(a.Name, "makes a sound")
    }

    type Dog struct {
        Animal  // 嵌入
        Breed string
    }

    dog := Dog{
        Animal: Animal{Name: "Buddy"},
        Breed:  "Golden Retriever",
    }

    // 方法提升
    dog.Speak()  // "Buddy makes a sound"

    // 直接访问嵌入字段
    fmt.Println(dog.Name)        // Buddy
    fmt.Println(dog.Animal.Name) // 同上
}

// 嵌入接口
func embeddingInterface() {
    type ReadWriter interface {
        io.Reader
        io.Writer
    }

    // 相当于:
    // type ReadWriter interface {
    //     Read(p []byte) (n int, err error)
    //     Write(p []byte) (n int, err error)
    // }
}
```

---

## 六、Go 1.26 新特性详解

### 6.1 new(expr) 语法

```
Go 1.26新特性:
────────────────────────────────────────

传统方式:
s := &Some{Field: value}

新方式:
s := new(Some{Field: value})

适用场景:
├─ 复杂结构的堆分配
├─ 可选字段的清晰表达
└─ 减少&和类型重复

语义等价性:
new(T{...}) ≡ &T{...}
但可读性提升

代码示例:
// 传统方式
type Config struct {
    Host string
    Port int
}

func oldWay() {
    // 必须写两遍类型名
    cfg := &Config{
        Host: "localhost",
        Port: 8080,
    }
    _ = cfg
}

// Go 1.26新方式
func newWay() {
    // 更清晰，不需要写类型名两次
    cfg := new(Config{
        Host: "localhost",
        Port: 8080,
    })
    _ = cfg
}

// 复杂结构更有优势
func complexExample() {
    type Server struct {
        Name    string
        Address string
        Config  Config
    }

    // 传统
    s1 := &Server{
        Name:    "api",
        Address: ":8080",
        Config: Config{
            Host: "localhost",
            Port: 8080,
        },
    }

    // Go 1.26
    s2 := new(Server{
        Name:    "api",
        Address: ":8080",
        Config: Config{
            Host: "localhost",
            Port: 8080,
        },
    })

    _ = s1
    _ = s2
}
```

### 6.2 递归泛型约束

```
Go 1.26增强:
────────────────────────────────────────

类型参数自引用:
type Ordered[T Ordered[T]] interface {
    comparable
    Less(T) bool
}

用途:
├─ 定义通用比较接口
├─ 支持树结构的递归定义
└─ 类型安全的自引用约束

示例:
type Node[T Node[T]] interface {
    Children() []T
    Value() int
}

代码示例:
// 递归约束用于树遍历
type TreeNode[T TreeNode[T]] interface {
    Value() int
    Children() []T
}

// 二叉树实现
type BinaryNode struct {
    val   int
    left  *BinaryNode
    right *BinaryNode
}

func (n *BinaryNode) Value() int {
    return n.val
}

func (n *BinaryNode) Children() []*BinaryNode {
    var children []*BinaryNode
    if n.left != nil {
        children = append(children, n.left)
    }
    if n.right != nil {
        children = append(children, n.right)
    }
    return children
}

// 通用树遍历
func TraverseTree[T TreeNode[T]](root T, visit func(int)) {
    if root == nil {
        return
    }
    visit(root.Value())
    for _, child := range root.Children() {
        TraverseTree(child, visit)
    }
}

// 使用
tree := &BinaryNode{
    val: 1,
    left: &BinaryNode{val: 2},
    right: &BinaryNode{val: 3},
}

TraverseTree(tree, func(v int) {
    fmt.Println(v)
})
```

### 6.3 Green Tea GC

```
Go 1.26 GC改进:
────────────────────────────────────────

特性:
├─ 更低延迟 (10-40%改善)
├─ 更高吞吐
├─ 更好内存利用
└─ 自动启用，无需配置

原理:
├─ 并发标记优化
├─ 写屏障改进
└─ 堆大小自适应

监控:
├─ runtime.ReadMemStats()
├─ GC CPU占比 < 10%为健康
└─ pprof分析GC行为

代码示例:
// GC监控
func monitorGC() {
    var m1, m2 runtime.MemStats

    runtime.GC()  // 强制GC
    runtime.ReadMemStats(&m1)

    // 执行一些操作...
    allocateMemory()

    runtime.GC()
    runtime.ReadMemStats(&m2)

    fmt.Printf("GC次数: %d\n", m2.NumGC-m1.NumGC)
    fmt.Printf("GC暂停: %d ns\n", m2.PauseNs[(m2.NumGC+255)%256])
    fmt.Printf("堆内存: %d KB\n", m2.HeapAlloc/1024)
}

func allocateMemory() {
    _ = make([]byte, 10*1024*1024)  // 10MB
}

// 调优GOGC
func tuneGC() {
    // GOGC=100 (默认): 堆增长到2倍时触发GC
    // GOGC=200: 更少的GC，更多内存使用
    // GOGC=50: 更频繁的GC，更少内存使用

    debug.SetGCPercent(200)
}
```

### 6.4 Goroutine Leak检测

```
Go 1.26新API:
────────────────────────────────────────

运行时检测:
runtime.SetGoroutineLeakCallback(
    func(gid uint64, stack []byte) {
        log.Printf("Leak: goroutine %d\n%s", gid, stack)
    },
)

测试集成:
func TestNoLeak(t *testing.T) {
    before := runtime.NumGoroutine()
    defer func() {
        time.Sleep(100 * time.Millisecond)
        if after := runtime.NumGoroutine(); after > before {
            t.Errorf("Leak: %d -> %d", before, after)
        }
    }()

    // 测试代码
}

代码示例:
// goroutine泄露示例
func leakyFunction() {
    ch := make(chan int)

    go func() {
        // 这个goroutine永远阻塞
        val := <-ch  // 无人发送
        fmt.Println(val)
    }()

    // 函数返回，但goroutine仍在运行
}

// 修复: 使用context
func fixedFunction(ctx context.Context) {
    ch := make(chan int)

    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return  // 正确退出
        }
    }()
}

// 使用goleak检测
import "go.uber.org/goleak"

func TestWithGoleak(t *testing.T) {
    defer goleak.VerifyNone(t)

    // 测试代码
    runFunction()
}
```

---

*本章从计算理论视角建立了Go语言基础的形式化理解，涵盖值语义、内存布局、控制流、函数闭包、复合类型和Go 1.26新特性，提供了丰富的代码示例和实战应用。*
