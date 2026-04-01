# Go 1.26.1 语法特性全面分析

> 本文档对 Go 1.26.1 的所有语法特性进行系统性分析，包括新引入的 `new()` 表达式特性和自引用泛型约束。

---

## 目录

- [Go 1.26.1 语法特性全面分析](#go-1261-语法特性全面分析)
  - [目录](#目录)
  - [1. 词法元素](#1-词法元素)
    - [1.1 标识符 (Identifiers)](#11-标识符-identifiers)
      - [概念定义](#概念定义)
      - [属性特征](#属性特征)
      - [关系依赖](#关系依赖)
      - [详细示例代码](#详细示例代码)
      - [反例说明](#反例说明)
    - [1.2 关键字 (Keywords)](#12-关键字-keywords)
      - [概念定义](#概念定义-1)
      - [属性特征](#属性特征-1)
      - [关键字分类](#关键字分类)
      - [详细示例代码](#详细示例代码-1)
    - [1.3 运算符 (Operators)](#13-运算符-operators)
      - [概念定义](#概念定义-2)
      - [运算符分类与优先级](#运算符分类与优先级)
      - [详细运算符说明](#详细运算符说明)
      - [反例说明](#反例说明-1)
    - [1.4 字面量 (Literals)](#14-字面量-literals)
      - [概念定义](#概念定义-3)
      - [字面量分类](#字面量分类)
      - [详细示例代码](#详细示例代码-2)
  - [2. 类型系统](#2-类型系统)
    - [2.1 基本类型](#21-基本类型)
      - [概念定义](#概念定义-4)
      - [类型分类](#类型分类)
      - [详细示例代码](#详细示例代码-3)
      - [反例说明](#反例说明-2)
    - [2.2 复合类型](#22-复合类型)
      - [概念定义](#概念定义-5)
      - [类型关系图](#类型关系图)
      - [详细示例代码](#详细示例代码-4)
    - [2.3 接口类型](#23-接口类型)
      - [概念定义](#概念定义-6)
      - [接口分类](#接口分类)
      - [详细示例代码](#详细示例代码-5)
    - [2.4 泛型类型 (Go 1.18+)](#24-泛型类型-go-118)
      - [概念定义](#概念定义-7)
      - [泛型组件](#泛型组件)
      - [详细示例代码](#详细示例代码-6)
  - [3. 变量和常量声明](#3-变量和常量声明)
    - [3.1 变量声明](#31-变量声明)
      - [概念定义](#概念定义-8)
      - [声明方式](#声明方式)
      - [详细示例代码](#详细示例代码-7)
    - [3.2 常量声明](#32-常量声明)
      - [概念定义](#概念定义-9)
      - [常量特性](#常量特性)
      - [详细示例代码](#详细示例代码-8)
    - [3.3 new() 内置函数（Go 1.26 新特性）](#33-new-内置函数go-126-新特性)
      - [概念定义](#概念定义-10)
      - [特性对比](#特性对比)
      - [详细示例代码](#详细示例代码-9)
      - [反例说明](#反例说明-3)
  - [4. 控制流语句](#4-控制流语句)
    - [4.1 if 语句](#41-if-语句)
      - [概念定义](#概念定义-11)
      - [语法结构](#语法结构)
      - [详细示例代码](#详细示例代码-10)
    - [4.2 for 语句](#42-for-语句)
      - [概念定义](#概念定义-12)
      - [循环形式](#循环形式)
      - [详细示例代码](#详细示例代码-11)
    - [4.3 switch 语句](#43-switch-语句)
      - [概念定义](#概念定义-13)
      - [switch 形式](#switch-形式)
      - [详细示例代码](#详细示例代码-12)
    - [4.4 select 语句](#44-select-语句)
      - [概念定义](#概念定义-14)
      - [select 特性](#select-特性)
      - [详细示例代码](#详细示例代码-13)
  - [5. 函数和方法声明](#5-函数和方法声明)
    - [5.1 函数声明](#51-函数声明)
      - [概念定义](#概念定义-15)
      - [函数特性](#函数特性)
      - [详细示例代码](#详细示例代码-14)
    - [5.2 方法声明](#52-方法声明)
      - [概念定义](#概念定义-16)
      - [方法特性](#方法特性)
      - [详细示例代码](#详细示例代码-15)
  - [6. 结构体和方法](#6-结构体和方法)
    - [6.1 结构体声明](#61-结构体声明)
      - [概念定义](#概念定义-17)
      - [结构体特性](#结构体特性)
      - [详细示例代码](#详细示例代码-16)
    - [6.2 结构体嵌入与方法集](#62-结构体嵌入与方法集)
      - [概念定义](#概念定义-18)
      - [详细示例代码](#详细示例代码-17)
  - [7. 包和导入机制](#7-包和导入机制)
    - [7.1 包声明](#71-包声明)
      - [概念定义](#概念定义-19)
      - [包类型](#包类型)
      - [详细示例代码](#详细示例代码-18)
    - [7.2 导入机制](#72-导入机制)
      - [概念定义](#概念定义-20)
      - [导入方式](#导入方式)
      - [详细示例代码](#详细示例代码-19)
    - [7.3 包初始化](#73-包初始化)
      - [概念定义](#概念定义-21)
      - [初始化顺序](#初始化顺序)
      - [详细示例代码](#详细示例代码-20)
  - [8. 新特性详解](#8-新特性详解)
    - [8.1 new() 表达式特性（Go 1.26）](#81-new-表达式特性go-126)
      - [概念定义](#概念定义-22)
      - [特性详解](#特性详解)
      - [详细示例代码](#详细示例代码-21)
    - [8.2 自引用泛型约束（Go 1.26）](#82-自引用泛型约束go-126)
      - [概念定义](#概念定义-23)
      - [特性详解](#特性详解-1)
      - [详细示例代码](#详细示例代码-22)
  - [9. 思维导图与决策树](#9-思维导图与决策树)
    - [9.1 Go 语法特性思维导图](#91-go-语法特性思维导图)
    - [9.2 类型选择决策树](#92-类型选择决策树)
    - [9.3 控制流选择决策树](#93-控制流选择决策树)
  - [10. 总结](#10-总结)
    - [Go 1.26.1 语法特性总览](#go-1261-语法特性总览)
  - [附录：Go 关键字速查表](#附录go-关键字速查表)

---

## 1. 词法元素

### 1.1 标识符 (Identifiers)

#### 概念定义

标识符是用于命名变量、类型、函数、包等程序实体的名称。

#### 属性特征

- **首字符规则**：必须以字母（Unicode 字母）或下划线 `_` 开头
- **后续字符**：可以是字母、数字或下划线
- **大小写敏感**：`MyVar` 和 `myvar` 是不同的标识符
- **长度限制**：无长度限制
- **Unicode 支持**：支持 Unicode 字符作为标识符

#### 关系依赖

```
标识符
├── 预声明标识符（内置）
│   ├── 类型：int, string, bool, error 等
│   ├── 常量：true, false, iota, nil
│   └── 函数：append, cap, close, complex, copy, delete, imag, len,
│           make, new, panic, print, println, real, recover
├── 包级标识符
│   ├── 包名
│   ├── 导出标识符（大写开头）
│   └── 非导出标识符（小写开头）
└── 局部标识符
    ├── 函数参数
    ├── 返回值命名
    └── 局部变量
```

#### 详细示例代码

```go
package main

// 合法标识符示例
var myVariable int           // 小写字母开头
var MyVariable int           // 大写字母开头（导出）
var _private int             // 下划线开头
var var123 int               // 字母开头后跟数字
var 变量名 int                // Unicode 支持（中文）
var π float64                // Unicode 希腊字母
var _ int                    // 空白标识符（特殊用途）

// 包级标识符
const MaxSize = 100          // 导出常量
const maxSize = 100          // 非导出常量

type MyType struct{}         // 导出类型
type myType struct{}         // 非导出类型

func PublicFunc() {}         // 导出函数
func privateFunc() {}        // 非导出函数
```

#### 反例说明

```go
// 非法标识符示例
var 123var int       // 错误：不能以数字开头
var my-var int       // 错误：不能包含连字符
var my var int       // 错误：不能包含空格
var type int         // 错误：不能使用关键字
var break int        // 错误：不能使用关键字
```

---

### 1.2 关键字 (Keywords)

#### 概念定义

关键字是 Go 语言预留的具有特殊含义的单词，不能用作标识符。

#### 属性特征

- Go 1.26.1 共有 **25 个关键字**
- 关键字全部小写
- 分为四类：声明、复合、控制、其他

#### 关键字分类

```
┌─────────────────────────────────────────────────────────────┐
│                      Go 关键字分类                           │
├──────────────┬──────────────────────────────────────────────┤
│   类别       │              关键字                          │
├──────────────┼──────────────────────────────────────────────┤
│ 声明相关     │ const, func, import, package, type, var      │
├──────────────┼──────────────────────────────────────────────┤
│ 复合类型     │ chan, interface, map, struct                 │
├──────────────┼──────────────────────────────────────────────┤
│ 控制流       │ break, case, continue, default, else,        │
│              │ fallthrough, for, goto, if, range, return,   │
│              │ switch                                       │
├──────────────┼──────────────────────────────────────────────┤
│ 其他         │ defer, go, select                            │
└──────────────┴──────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import "fmt"

// 声明关键字使用
const Pi = 3.14159          // const
var count int               // var

type Person struct {        // type, struct
    Name string
}

func main() {               // func
    // 控制流关键字
    if count > 0 {          // if, else
        fmt.Println("positive")
    } else {
        fmt.Println("non-positive")
    }

    // for 循环
    for i := 0; i < 10; i++ {  // for
        if i == 5 {
            break           // break
        }
        if i == 3 {
            continue        // continue
        }
    }

    // switch 语句
    switch count {          // switch, case, default
    case 0:
        fmt.Println("zero")
    default:
        fmt.Println("other")
    }

    // select 语句
    ch := make(chan int)    // chan, make
    go func() {             // go
        ch <- 1
    }()

    select {                // select
    case v := <-ch:
        fmt.Println(v)
    default:
        fmt.Println("no data")
    }

    // defer
    defer fmt.Println("deferred")  // defer

    // goto（谨慎使用）
    goto End
End:
    fmt.Println("end")
}
```

---

### 1.3 运算符 (Operators)

#### 概念定义

运算符是用于执行操作的符号，Go 运算符分为算术、比较、逻辑、位运算、赋值和其他运算符。

#### 运算符分类与优先级

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 运算符优先级表                             │
├──────────┬──────────────────────────────────────────────────────┤
│ 优先级   │ 运算符                                               │
├──────────┼──────────────────────────────────────────────────────┤
│    5     │ *  /  %  <<  >>  &  &^                               │
│    4     │ +  -  |  ^                                           │
│    3     │ ==  !=  <  <=  >  >=                                 │
│    2     │ &&                                                   │
│    1     │ ||                                                   │
└──────────┴──────────────────────────────────────────────────────┘
```

#### 详细运算符说明

```go
package main

func main() {
    a, b := 10, 3

    // ========== 算术运算符 ==========
    sum := a + b        // 加法: 13
    diff := a - b       // 减法: 7
    prod := a * b       // 乘法: 30
    quot := a / b       // 整数除法: 3
    rem := a % b        // 取模: 1

    // ========== 位运算符 ==========
    and := a & b        // 按位与: 1010 & 0011 = 0010 (2)
    or := a | b         // 按位或: 1010 | 0011 = 1011 (11)
    xor := a ^ b        // 按位异或: 1010 ^ 0011 = 1001 (9)
    not := ^a           // 按位取反
    leftShift := a << 2 // 左移: 1010 << 2 = 101000 (40)
    rightShift := a >> 1// 右移: 1010 >> 1 = 0101 (5)
    andNot := a &^ b    // 按位清除: 1010 &^ 0011 = 1000 (8)

    // ========== 比较运算符 ==========
    eq := a == b        // 等于: false
    ne := a != b        // 不等于: true
    lt := a < b         // 小于: false
    le := a <= b        // 小于等于: false
    gt := a > b         // 大于: true
    ge := a >= b        // 大于等于: true

    // ========== 逻辑运算符 ==========
    x, y := true, false
    land := x && y      // 逻辑与: false
    lor := x || y       // 逻辑或: true
    lnot := !x          // 逻辑非: false

    // ========== 赋值运算符 ==========
    c := 10
    c += 5              // c = c + 5
    c -= 3              // c = c - 3
    c *= 2              // c = c * 2
    c /= 4              // c = c / 4
    c %= 3              // c = c % 3
    c &= 1              // c = c & 1
    c |= 2              // c = c | 2
    c ^= 3              // c = c ^ 3
    c <<= 1             // c = c << 1
    c >>= 1             // c = c >> 1
    c &^= 1             // c = c &^ 1

    // ========== 其他运算符 ==========
    // 取地址
    ptr := &a           // 取 a 的地址
    val := *ptr         // 解引用，获取指针指向的值

    // 通道操作
    ch := make(chan int)
    ch <- 1             // 发送操作
    <-ch                // 接收操作（丢弃值）
    v := <-ch           // 接收操作（保存值）

    _ = sum; _ = diff; _ = prod; _ = quot; _ = rem
    _ = and; _ = or; _ = xor; _ = not; _ = leftShift; _ = rightShift; _ = andNot
    _ = eq; _ = ne; _ = lt; _ = le; _ = gt; _ = ge
    _ = land; _ = lor; _ = lnot
    _ = v
}
```

#### 反例说明

```go
package main

func main() {
    a, b := 10, 3

    // 错误：Go 不支持 ++ 和 -- 作为表达式
    // c := a++           // 编译错误

    // 错误：不支持三元运算符
    // max := a > b ? a : b  // 编译错误

    // 错误：不支持运算符重载
    // 必须使用函数来实现自定义类型的运算

    // 错误：不支持隐式类型转换
    var f float64 = 3.14
    // i := f + a         // 编译错误：类型不匹配
    i := int(f) + a       // 正确：显式转换
    _ = i
}
```

---

### 1.4 字面量 (Literals)

#### 概念定义

字面量是源代码中表示固定值的符号。

#### 字面量分类

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go 字面量类型                               │
├────────────────┬────────────────────────────────────────────────┤
│ 类型           │ 说明                                           │
├────────────────┼────────────────────────────────────────────────┤
│ 整数字面量     │ 十进制、八进制(0前缀)、十六进制(0x前缀)、二进制(0b前缀) │
├────────────────┼────────────────────────────────────────────────┤
│ 浮点数字面量   │ 小数形式、指数形式                             │
├────────────────┼────────────────────────────────────────────────┤
│ 复数字面量     │ 实部+虚部i                                     │
├────────────────┼────────────────────────────────────────────────┤
│ 字符串字面量   │ 双引号解释字符串、反引号原始字符串             │
├────────────────┼────────────────────────────────────────────────┤
│ 符文字面量     │ 单引号 Unicode 码点                            │
├────────────────┼────────────────────────────────────────────────┤
│ 布尔字面量     │ true, false                                    │
├────────────────┼────────────────────────────────────────────────┤
│ 复合字面量     │ 数组、切片、map、结构体、函数字面量            │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

func main() {
    // ========== 整数字面量 ==========
    dec := 42           // 十进制
    oct := 052          // 八进制 (以0开头) = 42
    hex := 0x2A         // 十六进制 = 42
    bin := 0b101010     // 二进制 = 42

    // 带下划线的数字字面量（提高可读性）
    million := 1_000_000
    hexWithSep := 0x_FF_FF

    // ========== 浮点数字面量 ==========
    f1 := 3.14159       // 小数形式
    f2 := 3.14e2        // 指数形式 = 314.0
    f3 := 3.14e-2       // 负指数 = 0.0314
    f4 := 1_234.567_89  // 带分隔符

    // ========== 复数字面量 ==========
    c1 := 3 + 4i        // 复数
    c2 := 1.5 + 2.5i    // 浮点复数

    // ========== 字符串字面量 ==========
    // 解释字符串（支持转义序列）
    s1 := "Hello, World!\n"
    s2 := "Tab\there"
    s3 := "Quote: \"quoted\""
    s4 := "Unicode: \u4e2d\u6587"  // 中文

    // 原始字符串（反引号，不解释转义）
    s5 := `Raw string with \n not interpreted`
    s6 := `Multi-line
    string preserved
    as-is`

    // 字符串连接
    s7 := "Hello, " + "World!"

    // ========== 符文字面量 ==========
    r1 := 'A'           // ASCII 字符
    r2 := '中'          // Unicode 字符
    r3 := '\n'          // 转义序列
    r4 := '\u4e2d'      // Unicode 码点（中）
    r5 := '\U00004e2d'  // Unicode 码点（完整形式）
    r6 := '\x41'        // 字节值

    // ========== 布尔字面量 ==========
    t := true
    f := false

    // ========== 复合字面量 ==========
    // 数组字面量
    arr := [3]int{1, 2, 3}
    arrWithIndex := [5]int{0: 1, 4: 5}  // 指定索引

    // 切片字面量
    slice := []int{1, 2, 3}

    // Map 字面量
    m := map[string]int{"one": 1, "two": 2}

    // 结构体字面量
    type Point struct{ X, Y int }
    p := Point{X: 10, Y: 20}
    p2 := Point{10, 20}  // 按位置

    // 函数字面量（匿名函数）
    fn := func(x, y int) int { return x + y }

    // nil 字面量
    var ptr *int = nil
    var iface interface{} = nil
    var sl []int = nil
    var mp map[string]int = nil
    var ch chan int = nil

    _ = dec; _ = oct; _ = hex; _ = bin; _ = million; _ = hexWithSep
    _ = f1; _ = f2; _ = f3; _ = f4; _ = c1; _ = c2
    _ = s1; _ = s2; _ = s3; _ = s4; _ = s5; _ = s6; _ = s7
    _ = r1; _ = r2; _ = r3; _ = r4; _ = r5; _ = r6
    _ = t; _ = f
    _ = arr; _ = arrWithIndex; _ = slice; _ = m; _ = p; _ = p2; _ = fn
    _ = ptr; _ = iface; _ = sl; _ = mp; _ = ch
}
```

---

## 2. 类型系统

### 2.1 基本类型

#### 概念定义

Go 的基本类型是语言内置的、不可再分解的原子类型。

#### 类型分类

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go 基本类型分类                             │
├──────────────┬──────────────────────────────────────────────────┤
│ 布尔类型     │ bool                                             │
├──────────────┼──────────────────────────────────────────────────┤
│ 整数类型     │ int, int8, int16, int32, int64                   │
│              │ uint, uint8, uint16, uint32, uint64, uintptr     │
├──────────────┼──────────────────────────────────────────────────┤
│ 浮点类型     │ float32, float64                                 │
├──────────────┼──────────────────────────────────────────────────┤
│ 复数类型     │ complex64, complex128                            │
├──────────────┼──────────────────────────────────────────────────┤
│ 字符串类型   │ string                                           │
├──────────────┼──────────────────────────────────────────────────┤
│ 符文类型     │ rune (int32 的别名)                              │
├──────────────┼──────────────────────────────────────────────────┤
│ 字节类型     │ byte (uint8 的别名)                              │
└──────────────┴──────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    // ========== 布尔类型 ==========
    var b bool = true
    b = false
    fmt.Printf("bool size: %d bytes\n", unsafe.Sizeof(b))

    // ========== 整数类型 ==========
    // 有符号整数
    var i int = -100          // 平台相关（32或64位）
    var i8 int8 = -128        // -128 到 127
    var i16 int16 = -32768    // -32768 到 32767
    var i32 int32 = -2147483648
    var i64 int64 = -9223372036854775808

    // 无符号整数
    var u uint = 100
    var u8 uint8 = 255        // 0 到 255
    var u16 uint16 = 65535    // 0 到 65535
    var u32 uint32 = 4294967295
    var u64 uint64 = 18446744073709551615

    // 指针大小的无符号整数
    var uptr uintptr          // 用于存储指针地址

    // 类型别名
    var r rune = '中'         // rune = int32，用于 Unicode 码点
    var by byte = 'A'         // byte = uint8，用于原始字节

    fmt.Printf("int size: %d bytes\n", unsafe.Sizeof(i))
    fmt.Printf("int8: %d, int16: %d, int32: %d, int64: %d\n",
        unsafe.Sizeof(i8), unsafe.Sizeof(i16), unsafe.Sizeof(i32), unsafe.Sizeof(i64))

    // ========== 浮点类型 ==========
    var f32 float32 = 3.14159265359  // 约6-7位有效数字
    var f64 float64 = 3.14159265358979323846  // 约15-16位有效数字

    fmt.Printf("float32 size: %d bytes, float64 size: %d bytes\n",
        unsafe.Sizeof(f32), unsafe.Sizeof(f64))

    // ========== 复数类型 ==========
    var c64 complex64 = 3 + 4i       // float32 实部和虚部
    var c128 complex128 = 3 + 4i     // float64 实部和虚部

    fmt.Printf("real: %f, imag: %f\n", real(c128), imag(c128))

    // ========== 字符串类型 ==========
    var s string = "Hello, Go!"
    s = "字符串是不可变的字节序列"

    // 字符串操作
    fmt.Printf("Length: %d bytes\n", len(s))
    fmt.Printf("First char: %c\n", s[0])

    // 字符串遍历
    for i, ch := range s {
        fmt.Printf("Index: %d, Rune: %c\n", i, ch)
    }

    // ========== 类型转换 ==========
    // Go 要求显式类型转换
    var x int = 10
    var y float64 = float64(x)    // int -> float64
    var z int = int(y)            // float64 -> int（截断小数）

    _ = i; _ = i8; _ = i16; _ = i32; _ = i64
    _ = u; _ = u8; _ = u16; _ = u32; _ = u64; _ = uptr
    _ = r; _ = by
    _ = c64; _ = z
}
```

#### 反例说明

```go
package main

func main() {
    // 错误：隐式类型转换不被允许
    var i int = 10
    // var f float64 = i    // 编译错误
    var f float64 = float64(i)  // 正确：显式转换

    // 错误：不同类型不能直接比较
    var a int = 10
    var b int32 = 10
    // if a == b { }        // 编译错误
    if a == int(b) { }      // 正确：先转换

    // 错误：溢出不会报错，但结果不正确
    var u8 uint8 = 255
    u8 = u8 + 1             // 溢出，结果为 0

    _ = f
}
```

---

### 2.2 复合类型

#### 概念定义

复合类型是由其他类型组合而成的类型，包括数组、切片、map、通道、指针、函数和接口。

#### 类型关系图

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go 复合类型                                 │
├────────────────┬────────────────────────────────────────────────┤
│ 数组 [n]T      │ 固定长度，值类型，可比较                       │
├────────────────┼────────────────────────────────────────────────┤
│ 切片 []T       │ 动态长度，引用类型，底层数组                   │
├────────────────┼────────────────────────────────────────────────┤
│ Map map[K]V    │ 键值对，引用类型，无序                         │
├────────────────┼────────────────────────────────────────────────┤
│ 通道 chan T    │ 用于 goroutine 间通信                          │
├────────────────┼────────────────────────────────────────────────┤
│ 指针 *T        │ 存储变量地址                                   │
├────────────────┼────────────────────────────────────────────────┤
│ 函数 func      │ 一等公民，可作为参数和返回值                   │
├────────────────┼────────────────────────────────────────────────┤
│ 结构体 struct  │ 字段集合，值类型                               │
├────────────────┼────────────────────────────────────────────────┤
│ 接口 interface │ 方法集合，实现隐式接口                         │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import "fmt"

func main() {
    // ========== 数组 ==========
    // 声明和初始化
    var arr1 [5]int                    // 零值初始化
    arr2 := [3]int{1, 2, 3}            // 字面量初始化
    arr3 := [...]int{1, 2, 3, 4, 5}    // 自动推断长度
    arr4 := [5]int{0: 1, 4: 5}         // 指定索引初始化

    // 多维数组
    var matrix [3][3]int
    matrix[0][0] = 1

    // 数组是值类型
    arrCopy := arr2     // 复制整个数组
    arrCopy[0] = 100    // 不影响原数组

    fmt.Println("Array:", arr1, arr2, arr3, arr4)

    // ========== 切片 ==========
    // 声明和初始化
    var s1 []int                    // nil 切片
    s2 := []int{1, 2, 3}            // 字面量
    s3 := make([]int, 5)            // 长度5，容量5
    s4 := make([]int, 3, 10)        // 长度3，容量10

    // 从数组创建切片
    arr := [5]int{1, 2, 3, 4, 5}
    slice1 := arr[1:4]              // [2, 3, 4]
    slice2 := arr[:3]               // [1, 2, 3]
    slice3 := arr[2:]               // [3, 4, 5]

    // 切片操作
    s2 = append(s2, 4, 5)           // 追加元素
    s2 = append(s2, []int{6, 7}...) // 追加切片

    // 复制切片
    dest := make([]int, len(s2))
    copy(dest, s2)

    fmt.Println("Slice:", s1, s2, s3, s4, slice1, slice2, slice3)

    // ========== Map ==========
    // 声明和初始化
    var m1 map[string]int           // nil map
    m2 := map[string]int{"a": 1, "b": 2}
    m3 := make(map[string]int)      // 空 map
    m4 := make(map[string]int, 10)  // 预分配空间

    // 操作
    m2["c"] = 3                     // 设置
    val := m2["a"]                  // 获取
    val2, ok := m2["d"]             // 获取并检查存在
    delete(m2, "b")                 // 删除

    // 遍历
    for k, v := range m2 {
        fmt.Printf("%s: %d\n", k, v)
    }

    fmt.Println("Map:", m1, m2, m3, m4, val, val2, ok)

    // ========== 通道 ==========
    // 无缓冲通道
    ch1 := make(chan int)

    // 有缓冲通道
    ch2 := make(chan int, 10)

    // 单向通道（用于函数参数）
    // sendOnly chan<- int
    // recvOnly <-chan int

    // 使用
    go func() {
        ch1 <- 42           // 发送
    }()

    v := <-ch1              // 接收

    // 关闭通道
    close(ch2)

    // 检查通道状态
    val3, ok2 := <-ch2      // ok2 为 false 表示通道已关闭且无数据

    fmt.Println("Channel:", v, val3, ok2)

    // ========== 指针 ==========
    x := 42
    p := &x                 // 取地址
    fmt.Println(*p)         // 解引用: 42
    *p = 100                // 通过指针修改
    fmt.Println(x)          // 100

    // 结构体指针
    type Point struct{ X, Y int }
    pt := &Point{X: 10, Y: 20}
    fmt.Println(pt.X)       // 自动解引用，等价于 (*pt).X

    // new 函数
    ptr := new(int)         // 分配内存并返回指针
    *ptr = 50

    // ========== 函数类型 ==========
    // 函数作为变量
    var add func(int, int) int
    add = func(a, b int) int { return a + b }

    // 函数作为参数
    apply := func(f func(int) int, x int) int {
        return f(x)
    }

    result := apply(func(x int) int { return x * x }, 5)
    fmt.Println("Function:", add(2, 3), result)

    _ = arrCopy; _ = dest; _ = m1; _ = m3; _ = m4
    _ = ch1; _ = ptr
}
```

---

### 2.3 接口类型

#### 概念定义

接口是方法签名的集合，任何实现了接口所有方法的类型都隐式实现了该接口。

#### 接口分类

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go 接口类型                                 │
├────────────────┬────────────────────────────────────────────────┤
│ 空接口         │ interface{} - 可存储任意类型                   │
├────────────────┼────────────────────────────────────────────────┤
│ 方法接口       │ 包含一个或多个方法签名                         │
├────────────────┼────────────────────────────────────────────────┤
│ 嵌入接口       │ 接口可以嵌入其他接口                           │
├────────────────┼────────────────────────────────────────────────┤
│ 类型约束接口   │ 用于泛型约束（Go 1.18+）                       │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import "fmt"

// ========== 基本接口 ==========
type Stringer interface {
    String() string
}

type Writer interface {
    Write([]byte) (int, error)
}

// ========== 嵌入接口 ==========
type ReadWriter interface {
    Reader          // 嵌入 Reader 接口
    Writer          // 嵌入 Writer 接口
}

type Reader interface {
    Read([]byte) (int, error)
}

// ========== 实现接口 ==========
type MyString string

func (s MyString) String() string {
    return string(s)
}

// ========== 空接口 ==========
func printAny(v interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}

// ========== 类型断言 ==========
func processValue(v interface{}) {
    // 类型断言
    if s, ok := v.(string); ok {
        fmt.Println("It's a string:", s)
        return
    }

    // 类型 switch
    switch x := v.(type) {
    case int:
        fmt.Println("It's an int:", x)
    case float64:
        fmt.Println("It's a float64:", x)
    case Stringer:
        fmt.Println("It has String():", x.String())
    default:
        fmt.Println("Unknown type:", x)
    }
}

// ========== 接口值结构 ==========
func inspectInterface(v interface{}) {
    fmt.Printf("Interface value: %v\n", v)
    fmt.Printf("Dynamic type: %T\n", v)
}

// ========== nil 接口 ==========
func describeNilInterface() {
    var i interface{}       // nil 接口值
    fmt.Println(i == nil)   // true

    var p *int = nil
    i = p                   // i 不是 nil，因为它有类型信息
    fmt.Println(i == nil)   // false!
}

func main() {
    // 接口使用
    var s Stringer = MyString("Hello")
    fmt.Println(s.String())

    // 空接口
    printAny(42)
    printAny("hello")
    printAny(3.14)

    // 类型断言和类型 switch
    processValue(100)
    processValue("test")
    processValue(MyString("custom"))

    // 接口值检查
    inspectInterface(42)
    inspectInterface("hello")

    // nil 接口
    describeNilInterface()
}
```

---

### 2.4 泛型类型 (Go 1.18+)

#### 概念定义

泛型允许编写适用于多种类型的代码，而不需要为每种类型重复编写。

#### 泛型组件

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go 泛型组件                                 │
├────────────────┬────────────────────────────────────────────────┤
│ 类型参数       │ [T any] - 函数或类型的泛型参数列表             │
├────────────────┼────────────────────────────────────────────────┤
│ 类型约束       │ interface 定义的类型集合                       │
├────────────────┼────────────────────────────────────────────────┤
│ 类型推断       │ 编译器自动推断类型参数                         │
├────────────────┼────────────────────────────────────────────────┤
│ 类型集合       │ ~T | U 形式的类型并集                          │
├────────────────┼────────────────────────────────────────────────┤
│ 近似元素       │ ~int 表示底层类型为 int 的所有类型             │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "cmp"
    "fmt"
)

// ========== 基本泛型函数 ==========
// 类型参数 T 约束为 any（任意类型）
func Print[T any](v T) {
    fmt.Println(v)
}

// 类型参数带约束
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// ========== 泛型类型 ==========
// 泛型栈
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    v := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return v, true
}

// 泛型链表节点
type Node[T any] struct {
    Value T
    Next  *Node[T]
}

// ========== 自定义类型约束 ==========
// 数值类型约束
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
        ~float32 | ~float64
}

// 使用自定义约束
func Sum[T Number](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}

// 近似元素的使用
type MyInt int  // 自定义类型

// ~int 包含 MyInt，因为 MyInt 的底层类型是 int
func Double[T ~int](x T) T {
    return x * 2
}

// ========== 多类型参数 ==========
func Map[K comparable, V any](keys []K, values []V) map[K]V {
    m := make(map[K]V)
    for i, k := range keys {
        if i < len(values) {
            m[k] = values[i]
        }
    }
    return m
}

// ========== 泛型接口 ==========
type Container[T any] interface {
    Add(T)
    Get() T
    Len() int
}

// ========== 类型推断 ==========
func main() {
    // 显式指定类型参数
    Print[string]("Hello")
    Print[int](42)

    // 类型推断
    Print("World")          // 推断为 string
    Print(100)              // 推断为 int

    // 泛型函数使用
    fmt.Println(Max(10, 20))        // int
    fmt.Println(Max(3.14, 2.71))    // float64
    fmt.Println(Max("a", "b"))      // string

    // 泛型类型使用
    intStack := Stack[int]{}
    intStack.Push(1)
    intStack.Push(2)
    intStack.Push(3)
    v, _ := intStack.Pop()
    fmt.Println("Popped:", v)

    stringStack := Stack[string]{}
    stringStack.Push("hello")
    stringStack.Push("world")

    // 自定义类型约束
    ints := []int{1, 2, 3, 4, 5}
    fmt.Println("Sum:", Sum(ints))

    floats := []float64{1.5, 2.5, 3.5}
    fmt.Println("Sum floats:", Sum(floats))

    // 近似元素
    var mi MyInt = 10
    fmt.Println("Double MyInt:", Double(mi))

    // 多类型参数
    keys := []string{"a", "b", "c"}
    vals := []int{1, 2, 3}
    m := Map(keys, vals)
    fmt.Println("Map:", m)

    // 泛型链表
    head := &Node[int]{Value: 1}
    head.Next = &Node[int]{Value: 2}
    head.Next.Next = &Node[int]{Value: 3}

    for n := head; n != nil; n = n.Next {
        fmt.Println("Node:", n.Value)
    }
}
```

---

## 3. 变量和常量声明

### 3.1 变量声明

#### 概念定义

变量是存储数据的命名位置，Go 提供多种声明方式。

#### 声明方式

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 变量声明方式                               │
├────────────────┬────────────────────────────────────────────────┤
│ var 声明       │ var x int = 10                                 │
│                │ var x = 10                                     │
│                │ var x int                                      │
├────────────────┼────────────────────────────────────────────────┤
│ 短变量声明     │ x := 10                                        │
│                │ 只能在函数内部使用                             │
├────────────────┼────────────────────────────────────────────────┤
│ 多变量声明     │ var a, b, c int                                │
│                │ a, b, c := 1, 2, 3                             │
├────────────────┼────────────────────────────────────────────────┤
│ 分组声明       │ var (                                          │
│                │     x int = 10                                 │
│                │     y string = "hello"                         │
│                │ )                                              │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import "fmt"

// 包级变量声明
var packageVar int = 100

// 分组声明
var (
    name    string = "Go"
    version float64 = 1.26
    isReady bool   = true
)

func main() {
    // ========== var 声明 ==========
    var a int = 10              // 完整形式
    var b = 20                  // 类型推断
    var c int                   // 零值初始化

    fmt.Printf("a=%d, b=%d, c=%d\n", a, b, c)

    // ========== 短变量声明 ==========
    x := 30                     // 只能在函数内
    y := "hello"
    z := 3.14

    fmt.Printf("x=%d, y=%s, z=%f\n", x, y, z)

    // ========== 多变量声明 ==========
    var m, n, p int = 1, 2, 3
    q, r, s := 4, 5, 6

    fmt.Printf("m=%d, n=%d, p=%d\n", m, n, p)
    fmt.Printf("q=%d, r=%d, s=%d\n", q, r, s)

    // 不同类型多变量
    var (
        i   int    = 1
        str string = "test"
        f   float64 = 2.5
    )

    // 短变量声明多变量
    ii, ss, ff := 2, "test2", 3.5

    // ========== 重新声明 ==========
    // 短变量声明中，至少一个变量必须是新的
    x, newVar := 35, 100        // x 被重新赋值，newVar 是新变量
    fmt.Printf("x=%d, newVar=%d\n", x, newVar)

    // ========== 变量作用域 ==========
    {
        inner := "inner scope"
        fmt.Println(inner)
    }
    // fmt.Println(inner)  // 错误：inner 超出作用域

    // ========== 变量零值 ==========
    var intZero int             // 0
    var stringZero string       // ""
    var boolZero bool           // false
    var sliceZero []int         // nil
    var mapZero map[string]int  // nil
    var ptrZero *int            // nil
    var interfaceZero interface{} // nil

    fmt.Printf("Zero values: %v, %q, %v, %v, %v, %v, %v\n",
        intZero, stringZero, boolZero, sliceZero, mapZero, ptrZero, interfaceZero)

    // ========== 匿名变量 ==========
    // 使用 _ 忽略不需要的值
    _, err := someFunc()        // 忽略第一个返回值
    if err != nil {
        fmt.Println("Error:", err)
    }

    _ = i; _ = str; _ = f; _ = ii; _ = ss; _ = ff
}

func someFunc() (int, error) {
    return 42, nil
}
```

---

### 3.2 常量声明

#### 概念定义

常量是在编译时确定且不可修改的值。

#### 常量特性

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 常量特性                                   │
├────────────────┬────────────────────────────────────────────────┤
│ 无类型常量     │ 可以表示任意精度，直到赋值时才确定类型         │
├────────────────┼────────────────────────────────────────────────┤
│ 有类型常量     │ 声明时指定类型，有范围限制                     │
├────────────────┼────────────────────────────────────────────────┤
│ iota           │ 常量生成器，从0开始递增                        │
├────────────────┼────────────────────────────────────────────────┤
│ 常量表达式     │ 可以在编译时计算的表达式                       │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import "fmt"

// ========== 基本常量 ==========
const Pi = 3.14159265358979323846
const Greeting = "Hello, World!"
const MaxRetries = 3

// 有类型常量
const TypedPi float64 = 3.14159

// 分组常量声明
const (
    Monday = iota    // 0
    Tuesday          // 1
    Wednesday        // 2
    Thursday         // 3
    Friday           // 4
    Saturday         // 5
    Sunday           // 6
)

// iota 高级用法
const (
    _ = iota                    // 跳过 0
    KB = 1 << (10 * iota)       // 1 << 10 = 1024
    MB = 1 << (10 * iota)       // 1 << 20 = 1048576
    GB = 1 << (10 * iota)       // 1 << 30 = 1073741824
    TB = 1 << (10 * iota)       // 1 << 40
)

// iota 位掩码
const (
    Read Permission = 1 << iota    // 1
    Write                          // 2
    Execute                        // 4
)

type Permission int

// ========== 无类型常量的灵活性 ==========
const Big = 100000000000000000000  // 无类型整数，任意精度

func main() {
    // 基本常量使用
    fmt.Println("Pi:", Pi)
    fmt.Println("Greeting:", Greeting)

    // 枚举值
    fmt.Println("Days:", Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday)

    // 字节单位
    fmt.Printf("KB: %d, MB: %d, GB: %d\n", KB, MB, GB)

    // 权限位掩码
    var p Permission = Read | Write
    fmt.Printf("Permission: %b (%d)\n", p, p)
    fmt.Println("Can read:", p&Read != 0)
    fmt.Println("Can execute:", p&Execute != 0)

    // 无类型常量的灵活性
    var f float64 = Big      // 转换为 float64
    var i int64 = Big / 100  // 转换为 int64
    fmt.Printf("Big as float64: %e\n", f)
    fmt.Printf("Big/100 as int64: %d\n", i)

    // 常量表达式
    const (
        SecondsPerMinute = 60
        SecondsPerHour   = 60 * SecondsPerMinute
        SecondsPerDay    = 24 * SecondsPerHour
    )
    fmt.Println("Seconds per day:", SecondsPerDay)

    // 复杂常量表达式
    const (
        a = 1
        b = a + 2           // 3
        c = b * 3           // 9
        d = c << 1          // 18
    )
    fmt.Println("Expression constants:", a, b, c, d)
}
```

---

### 3.3 new() 内置函数（Go 1.26 新特性）

#### 概念定义

Go 1.26 中，`new()` 现在可以接受表达式作为操作数，而不仅仅是类型。

#### 特性对比

```
┌─────────────────────────────────────────────────────────────────┐
│                    new() 函数变化                                │
├────────────────┬────────────────────────────────────────────────┤
│ Go 1.25 及之前 │ new(T) - T 必须是类型                          │
├────────────────┼────────────────────────────────────────────────┤
│ Go 1.26+       │ new(T) - T 可以是类型或表达式                  │
│                │ new(expr) 等价于 &expr                         │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    // ========== 传统用法（仍然有效）==========
    ptr1 := new(int)            // 创建 int 指针
    *ptr1 = 42
    fmt.Println("Traditional new():", *ptr1)

    // ========== Go 1.26 新特性：表达式操作数 ==========
    // 可以直接创建指向值的指针
    ptr2 := new(int64(300))     // 创建指向 int64(300) 的指针
    fmt.Println("new(int64(300)):", *ptr2)

    // 等价于
    val := int64(300)
    ptr3 := &val
    fmt.Println("&val:", *ptr3)

    // ========== 实际应用场景：可选字段 ==========
    // JSON 序列化中的可选字段
    type Config struct {
        Name    *string `json:"name,omitempty"`
        Timeout *int    `json:"timeout,omitempty"`
        Enabled *bool   `json:"enabled,omitempty"`
    }

    // 传统方式（繁琐）
    timeout1 := 30
    config1 := Config{
        Name:    stringPtr("server1"),
        Timeout: &timeout1,
    }

    // Go 1.26 新方式（简洁）
    config2 := Config{
        Name:    new(string("server2")),
        Timeout: new(int(60)),
        Enabled: new(bool(true)),
    }

    // 序列化
    json1, _ := json.Marshal(config1)
    json2, _ := json.Marshal(config2)

    fmt.Println("Config1 JSON:", string(json1))
    fmt.Println("Config2 JSON:", string(json2))

    // ========== 更多表达式示例 ==========
    // 算术表达式
    ptr4 := new(10 + 20)        // 指向 30
    fmt.Println("new(10 + 20):", *ptr4)

    // 函数调用结果
    ptr5 := new(len("hello"))   // 指向 5
    fmt.Println("new(len(\"hello\")):", *ptr5)

    // 类型转换表达式
    ptr6 := new(float64(42))    // 指向 42.0
    fmt.Println("new(float64(42)):", *ptr6)

    // 复合字面量
    ptr7 := new([3]int{1, 2, 3})
    fmt.Println("new([3]int{1,2,3}):", *ptr7)

    _ = config1
}

// 辅助函数（传统方式需要）
func stringPtr(s string) *string {
    return &s
}
```

#### 反例说明

```go
package main

func main() {
    // 错误：new() 不能用于变量
    x := 10
    // ptr := new(x)        // 编译错误：x 不是类型或常量表达式

    // 正确：使用 & 操作符
    ptr := &x
    _ = ptr

    // 错误：new() 表达式必须是可寻址的
    // ptr2 := new(10)      // 编译错误：10 不是可寻址的

    // 正确：使用变量
    val := 10
    ptr2 := &val
    _ = ptr2
}
```

---

## 4. 控制流语句

### 4.1 if 语句

#### 概念定义

if 语句根据条件执行代码块，支持前置语句。

#### 语法结构

```
┌─────────────────────────────────────────────────────────────────┐
│                    if 语句结构                                   │
├────────────────┬────────────────────────────────────────────────┤
│ 基本形式       │ if condition { ... }                           │
├────────────────┼────────────────────────────────────────────────┤
│ if-else        │ if condition { ... } else { ... }              │
├────────────────┼────────────────────────────────────────────────┤
│ if-else-if     │ if c1 { ... } else if c2 { ... } else { ... }  │
├────────────────┼────────────────────────────────────────────────┤
│ 带前置语句     │ if v := expr; condition { ... }                │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    // ========== 基本 if ==========
    x := 10
    if x > 5 {
        fmt.Println("x is greater than 5")
    }

    // ========== if-else ==========
    if x%2 == 0 {
        fmt.Println("x is even")
    } else {
        fmt.Println("x is odd")
    }

    // ========== if-else-if ==========
    score := 85
    if score >= 90 {
        fmt.Println("Grade: A")
    } else if score >= 80 {
        fmt.Println("Grade: B")
    } else if score >= 70 {
        fmt.Println("Grade: C")
    } else if score >= 60 {
        fmt.Println("Grade: D")
    } else {
        fmt.Println("Grade: F")
    }

    // ========== 带前置语句的 if ==========
    // 前置语句中声明的变量只在 if 块内有效
    if err := doSomething(); err != nil {
        fmt.Println("Error:", err)
    }
    // err 在这里不可见

    // 常见模式：文件操作
    if file, err := os.Open("test.txt"); err == nil {
        defer file.Close()
        // 使用 file
        fmt.Println("File opened successfully")
    } else {
        fmt.Println("Failed to open file:", err)
    }
    // file 在这里不可见

    // ========== 嵌套 if ==========
    age := 25
    hasID := true
    if age >= 18 {
        if hasID {
            fmt.Println("Access granted")
        } else {
            fmt.Println("Please show your ID")
        }
    } else {
        fmt.Println("Access denied: underage")
    }

    // ========== 逻辑组合 ==========
    if age >= 18 && hasID {
        fmt.Println("Access granted (combined)")
    }

    // ========== 短路求值 ==========
    // && 和 || 是短路运算符
    if x > 0 && 10/x > 1 {  // 如果 x <= 0，不会执行 10/x
        fmt.Println("Both conditions true")
    }
}

func doSomething() error {
    return errors.New("something went wrong")
}
```

---

### 4.2 for 语句

#### 概念定义

for 是 Go 中唯一的循环语句，但支持多种形式。

#### 循环形式

```
┌─────────────────────────────────────────────────────────────────┐
│                    for 循环形式                                  │
├────────────────┬────────────────────────────────────────────────┤
│ 完整形式       │ for init; condition; post { ... }              │
├────────────────┼────────────────────────────────────────────────┤
│ 条件形式       │ for condition { ... }  // 类似 while           │
├────────────────┼────────────────────────────────────────────────┤
│ 无限循环       │ for { ... }                                    │
├────────────────┼────────────────────────────────────────────────┤
│ range 循环     │ for k, v := range collection { ... }           │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import "fmt"

func main() {
    // ========== 完整 for 循环 ==========
    for i := 0; i < 5; i++ {
        fmt.Println("Iteration:", i)
    }

    // ========== 条件形式（类似 while）==========
    count := 0
    for count < 5 {
        fmt.Println("Count:", count)
        count++
    }

    // ========== 无限循环 ==========
    // 必须包含 break 或 return
    counter := 0
    for {
        if counter >= 3 {
            break
        }
        fmt.Println("Infinite loop iteration:", counter)
        counter++
    }

    // ========== range 循环 - 数组/切片 ==========
    nums := []int{10, 20, 30, 40, 50}

    // 索引和值
    for i, v := range nums {
        fmt.Printf("Index: %d, Value: %d\n", i, v)
    }

    // 只获取索引
    for i := range nums {
        fmt.Printf("Index: %d\n", i)
    }

    // 只获取值（使用空白标识符）
    for _, v := range nums {
        fmt.Printf("Value: %d\n", v)
    }

    // ========== range 循环 - Map ==========
    scores := map[string]int{"Alice": 90, "Bob": 85, "Charlie": 95}

    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
    }

    // Map 遍历顺序是随机的
    // 如果需要固定顺序，需要先获取键并排序

    // ========== range 循环 - 字符串 ==========
    str := "Hello, 世界"

    // 按字节遍历
    for i := 0; i < len(str); i++ {
        fmt.Printf("Byte %d: %c\n", i, str[i])
    }

    // 按符文遍历（推荐）
    for i, r := range str {
        fmt.Printf("Index: %d, Rune: %c\n", i, r)
    }

    // ========== range 循环 - 通道 ==========
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)

    for v := range ch {
        fmt.Println("Received:", v)
    }

    // ========== 控制语句 ==========
    // break - 跳出循环
    for i := 0; i < 10; i++ {
        if i == 5 {
            break  // 跳出整个循环
        }
        fmt.Println("Break at:", i)
    }

    // continue - 跳过当前迭代
    for i := 0; i < 5; i++ {
        if i == 2 {
            continue  // 跳过 i == 2
        }
        fmt.Println("Continue at:", i)
    }

    // 标签和 break（跳出外层循环）
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                break outer  // 跳出 outer 标签的循环
            }
            fmt.Printf("(%d, %d)\n", i, j)
        }
    }

    // ========== 复杂示例 ==========
    // 嵌套循环
    matrix := [][]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }

    for i, row := range matrix {
        for j, val := range row {
            fmt.Printf("matrix[%d][%d] = %d\n", i, j, val)
        }
    }
}
```

---

### 4.3 switch 语句

#### 概念定义

switch 语句提供多路分支选择，比多个 if-else 更清晰。

#### switch 形式

```
┌─────────────────────────────────────────────────────────────────┐
│                    switch 语句形式                               │
├────────────────┬────────────────────────────────────────────────┤
│ 表达式 switch  │ switch expr { case v1: ... case v2: ... }      │
├────────────────┼────────────────────────────────────────────────┤
│ 类型 switch    │ switch v := x.(type) { case T1: ... }          │
├────────────────┼────────────────────────────────────────────────┤
│ 无表达式 switch│ switch { case cond1: ... case cond2: ... }     │
│                │ （相当于 if-else-if 链）                       │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // ========== 基本 switch ==========
    day := 3
    switch day {
    case 1:
        fmt.Println("Monday")
    case 2:
        fmt.Println("Tuesday")
    case 3:
        fmt.Println("Wednesday")
    case 4:
        fmt.Println("Thursday")
    case 5:
        fmt.Println("Friday")
    case 6, 7:  // 多个 case 值
        fmt.Println("Weekend")
    default:
        fmt.Println("Invalid day")
    }

    // ========== 带前置语句的 switch ==========
    switch hour := time.Now().Hour(); {
    case hour < 12:
        fmt.Println("Good morning!")
    case hour < 18:
        fmt.Println("Good afternoon!")
    default:
        fmt.Println("Good evening!")
    }

    // ========== 无表达式 switch ==========
    score := 85
    switch {
    case score >= 90:
        fmt.Println("A")
    case score >= 80:
        fmt.Println("B")
    case score >= 70:
        fmt.Println("C")
    case score >= 60:
        fmt.Println("D")
    default:
        fmt.Println("F")
    }

    // ========== fallthrough ==========
    n := 2
    switch n {
    case 1:
        fmt.Println("One")
        fallthrough  // 继续执行下一个 case
    case 2:
        fmt.Println("Two")
        fallthrough  // 继续执行下一个 case
    case 3:
        fmt.Println("Three")
        // 没有 fallthrough，不会执行 default
    default:
        fmt.Println("Default")
    }
    // 输出: Two
    //       Three

    // ========== 类型 switch ==========
    var x interface{} = "hello"

    switch v := x.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s (length: %d)\n", v, len(v))
    case float64:
        fmt.Printf("Float64: %f\n", v)
    case bool:
        fmt.Printf("Bool: %t\n", v)
    case nil:
        fmt.Println("Nil")
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }

    // ========== 类型 switch 与多个类型 ==========
    var y interface{} = 42

    switch y.(type) {
    case int, int8, int16, int32, int64:
        fmt.Println("Signed integer")
    case uint, uint8, uint16, uint32, uint64:
        fmt.Println("Unsigned integer")
    case float32, float64:
        fmt.Println("Float")
    default:
        fmt.Println("Other type")
    }

    // ========== 复杂示例 ==========
    describe := func(i interface{}) {
        switch v := i.(type) {
        case int:
            fmt.Printf("Twice %d is %d\n", v, v*2)
        case string:
            fmt.Printf("%q is %d bytes long\n", v, len(v))
        case []int:
            fmt.Printf("Slice with %d elements\n", len(v))
        case map[string]int:
            fmt.Printf("Map with %d entries\n", len(v))
        default:
            fmt.Printf("I don't know about type %T!\n", v)
        }
    }

    describe(21)
    describe("hello")
    describe([]int{1, 2, 3})
    describe(map[string]int{"a": 1})
    describe(3.14)
}
```

---

### 4.4 select 语句

#### 概念定义

select 语句用于在多个通道操作中进行选择，类似于 switch，但用于通道。

#### select 特性

```
┌─────────────────────────────────────────────────────────────────┐
│                    select 语句特性                               │
├────────────────┬────────────────────────────────────────────────┤
│ 随机选择       │ 多个 case 就绪时随机选择一个                   │
├────────────────┼────────────────────────────────────────────────┤
│ 阻塞           │ 没有 case 就绪时阻塞                           │
├────────────────┼────────────────────────────────────────────────┤
│ default        │ 没有 case 就绪时立即执行 default               │
├────────────────┼────────────────────────────────────────────────┤
│ 非阻塞操作     │ 配合 default 实现非阻塞通道操作                │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // ========== 基本 select ==========
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from ch1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from ch2"
    }()

    // 接收两个消息
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received:", msg2)
        }
    }

    // ========== 带 default 的非阻塞 select ==========
    ch := make(chan int, 1)

    // 非阻塞发送
    select {
    case ch <- 42:
        fmt.Println("Sent 42")
    default:
        fmt.Println("Channel full, couldn't send")
    }

    // 非阻塞接收
    select {
    case v := <-ch:
        fmt.Println("Received:", v)
    default:
        fmt.Println("No data available")
    }

    // ========== 超时处理 ==========
    slowCh := make(chan string)

    go func() {
        time.Sleep(3 * time.Second)
        slowCh <- "slow response"
    }()

    select {
    case msg := <-slowCh:
        fmt.Println("Received:", msg)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }

    // ========== 多 case 就绪（随机选择）==========
    readyCh1 := make(chan int, 1)
    readyCh2 := make(chan int, 1)

    readyCh1 <- 1
    readyCh2 <- 2

    // 两个通道都有数据，随机选择一个
    select {
    case v := <-readyCh1:
        fmt.Println("From ch1:", v)
    case v := <-readyCh2:
        fmt.Println("From ch2:", v)
    }

    // ========== 使用 select 实现定时器 ==========
    ticker := time.NewTicker(500 * time.Millisecond)
    done := make(chan bool)

    go func() {
        time.Sleep(2 * time.Second)
        done <- true
    }()

    for {
        select {
        case <-ticker.C:
            fmt.Println("Tick at", time.Now())
        case <-done:
            fmt.Println("Done!")
            ticker.Stop()
            return
        }
    }
}
```

---

## 5. 函数和方法声明

### 5.1 函数声明

#### 概念定义

函数是执行特定任务的代码块，Go 函数支持多返回值、变长参数、匿名函数等特性。

#### 函数特性

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 函数特性                                   │
├────────────────┬────────────────────────────────────────────────┤
│ 多返回值       │ func f() (int, error)                          │
├────────────────┼────────────────────────────────────────────────┤
│ 命名返回值     │ func f() (result int, err error)               │
├────────────────┼────────────────────────────────────────────────┤
│ 变长参数       │ func f(args ...int)                            │
├────────────────┼────────────────────────────────────────────────┤
│ 匿名函数       │ func() { ... }                                 │
├────────────────┼────────────────────────────────────────────────┤
│ 闭包           │ 函数可以捕获外部变量                           │
├────────────────┼────────────────────────────────────────────────┤
│ 递归           │ 函数可以调用自身                               │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "errors"
    "fmt"
)

// ========== 基本函数 ==========
func add(a, b int) int {
    return a + b
}

// ========== 多返回值 ==========
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// ========== 命名返回值 ==========
func rectangle(width, height float64) (area, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return  // 裸 return，返回命名值
}

// ========== 变长参数 ==========
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 变长参数 + 普通参数
func printf(format string, args ...interface{}) {
    fmt.Printf(format, args...)
}

// ========== 函数作为参数 ==========
func applyOperation(a, b int, op func(int, int) int) int {
    return op(a, b)
}

// ========== 函数作为返回值 ==========
func makeMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

// ========== 递归函数 ==========
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// 尾递归优化（Go 不自动优化，需要手动）
func factorialTail(n, acc int) int {
    if n <= 1 {
        return acc
    }
    return factorialTail(n-1, n*acc)
}

// ========== defer 语句 ==========
func deferExample() {
    defer fmt.Println("First defer")
    defer fmt.Println("Second defer")
    fmt.Println("Function body")
    // 输出顺序：
    // Function body
    // Second defer
    // First defer
}

// defer 用于资源清理
func processFile(filename string) error {
    file, err := openFile(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // 确保文件关闭

    // 处理文件...
    return nil
}

type File struct{}

func (f *File) Close() error { return nil }

func openFile(name string) (*File, error) {
    return &File{}, nil
}

// ========== panic 和 recover ==========
func mayPanic() {
    panic("something went wrong")
}

func safeCall() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()
    mayPanic()
    fmt.Println("This won't be printed")
}

// ========== 泛型函数 ==========
func GenericMax[T comparable](a, b T) T {
    // 注意：comparable 不支持 > 操作
    // 这里仅作示例
    return a
}

func main() {
    // 基本函数
    fmt.Println("add(3, 4):", add(3, 4))

    // 多返回值
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }

    // 命名返回值
    area, perimeter := rectangle(5, 3)
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", area, perimeter)

    // 变长参数
    fmt.Println("sum(1, 2, 3, 4, 5):", sum(1, 2, 3, 4, 5))
    nums := []int{1, 2, 3, 4, 5}
    fmt.Println("sum(nums...):", sum(nums...))  // 展开切片

    // 函数作为参数
    multiply := func(a, b int) int { return a * b }
    fmt.Println("applyOperation(3, 4, multiply):", applyOperation(3, 4, multiply))

    // 函数作为返回值（闭包）
    double := makeMultiplier(2)
    triple := makeMultiplier(3)
    fmt.Println("double(5):", double(5))
    fmt.Println("triple(5):", triple(5))

    // 递归
    fmt.Println("factorial(5):", factorial(5))
    fmt.Println("factorialTail(5, 1):", factorialTail(5, 1))

    // defer
    deferExample()

    // panic/recover
    safeCall()
    fmt.Println("Program continues after recovery")

    // 匿名函数
    func() {
        fmt.Println("Anonymous function")
    }()

    // 立即执行函数
    result2 := func(a, b int) int {
        return a + b
    }(3, 4)
    fmt.Println("IIFE result:", result2)
}
```

---

### 5.2 方法声明

#### 概念定义

方法是带有接收者的函数，接收者可以是值类型或指针类型。

#### 方法特性

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 方法特性                                   │
├────────────────┬────────────────────────────────────────────────┤
│ 值接收者       │ func (t T) Method() - 接收值的副本             │
├────────────────┼────────────────────────────────────────────────┤
│ 指针接收者     │ func (t *T) Method() - 接收指针，可修改原值    │
├────────────────┼────────────────────────────────────────────────┤
│ 方法集         │ 类型 T 的方法集：所有值接收者方法              │
│                │ 类型 *T 的方法集：所有值和指针接收者方法       │
├────────────────┼────────────────────────────────────────────────┤
│ 嵌入类型       │ 可以嵌入其他类型来继承方法                     │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "fmt"
    "math"
)

// ========== 基本方法 ==========
type Rectangle struct {
    Width, Height float64
}

// 值接收者方法
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// 值接收者方法（不修改原值）
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// 指针接收者方法（可修改原值）
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

// ========== 自定义类型方法 ==========
type MyInt int

func (m MyInt) IsPositive() bool {
    return m > 0
}

func (m MyInt) Double() MyInt {
    return m * 2
}

// ========== 嵌入类型（继承方法）==========
type Point struct {
    X, Y float64
}

func (p Point) Distance() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

type ColoredPoint struct {
    Point           // 嵌入 Point，继承其方法
    Color string
}

// ========== 接口实现 ==========
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle 已经实现了 Shape 接口

// ========== 方法值和方法表达式 ==========
type Calculator struct {
    Value int
}

func (c Calculator) Add(n int) int {
    return c.Value + n
}

func (c *Calculator) AddToValue(n int) {
    c.Value += n
}

func main() {
    // ========== 基本方法使用 ==========
    rect := Rectangle{Width: 10, Height: 5}

    fmt.Printf("Area: %.2f\n", rect.Area())
    fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())

    // 指针接收者方法
    rect.Scale(2)
    fmt.Printf("After scale: Width=%.2f, Height=%.2f\n", rect.Width, rect.Height)

    // ========== 自定义类型方法 ==========
    num := MyInt(42)
    fmt.Println("IsPositive:", num.IsPositive())
    fmt.Println("Double:", num.Double())

    // ========== 嵌入类型 ==========
    cp := ColoredPoint{
        Point: Point{X: 3, Y: 4},
        Color: "red",
    }

    // 可以直接调用嵌入类型的方法
    fmt.Printf("ColoredPoint Distance: %.2f\n", cp.Distance())
    // 也可以这样调用
    fmt.Printf("ColoredPoint Distance: %.2f\n", cp.Point.Distance())

    // ========== 方法值 ==========
    calc := Calculator{Value: 10}

    // 方法值 - 绑定接收者
    addFive := calc.Add  // 方法值，接收者绑定为 calc
    fmt.Println("addFive(3):", addFive(3))
    fmt.Println("addFive(7):", addFive(7))

    // ========== 方法表达式 ==========
    // 方法表达式 - 不绑定接收者
    add := Calculator.Add  // 需要显式传递接收者
    fmt.Println("add(calc, 5):", add(calc, 5))

    // 指针方法表达式
    addToValue := (*Calculator).AddToValue
    addToValue(&calc, 10)
    fmt.Println("After AddToValue:", calc.Value)

    // ========== 值 vs 指针接收者 ==========
    r1 := Rectangle{Width: 5, Height: 3}
    r2 := &Rectangle{Width: 5, Height: 3}

    // 值接收者方法可以被值和指针调用
    fmt.Println("r1.Area():", r1.Area())
    fmt.Println("r2.Area():", r2.Area())

    // 指针接收者方法可以被值和指针调用（编译器自动取地址）
    r1.Scale(2)  // 等价于 (&r1).Scale(2)
    r2.Scale(2)

    fmt.Printf("r1 after scale: %.2f x %.2f\n", r1.Width, r1.Height)
    fmt.Printf("r2 after scale: %.2f x %.2f\n", r2.Width, r2.Height)

    // ========== 接口赋值 ==========
    var shape Shape = rect  // Rectangle 实现了 Shape
    fmt.Printf("Shape Area: %.2f\n", shape.Area())
}
```

---

## 6. 结构体和方法

### 6.1 结构体声明

#### 概念定义

结构体是字段的集合，用于组合相关数据。

#### 结构体特性

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 结构体特性                                 │
├────────────────┬────────────────────────────────────────────────┤
│ 命名字段       │ Name string                                    │
├────────────────┼────────────────────────────────────────────────┤
│ 嵌入字段       │ 只写类型名，继承其方法                         │
├────────────────┼────────────────────────────────────────────────┤
│ 标签           │ `json:"name"` 用于反射                         │
├────────────────┼────────────────────────────────────────────────┤
│ 可比较性       │ 所有字段可比较时，结构体可比较                 │
├────────────────┼────────────────────────────────────────────────┤
│ 匿名字段       │ 可以省略字段名                                 │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "encoding/json"
    "fmt"
    "reflect"
)

// ========== 基本结构体 ==========
type Person struct {
    Name    string
    Age     int
    Email   string
}

// ========== 带标签的结构体 ==========
type User struct {
    ID        int       `json:"id" db:"user_id"`
    Username  string    `json:"username" validate:"required"`
    Password  string    `json:"-"`  // 忽略此字段
    Email     string    `json:"email,omitempty"`
    CreatedAt string    `json:"created_at"`
}

// ========== 嵌入结构体 ==========
type Address struct {
    Street  string
    City    string
    Country string
}

type Employee struct {
    Person          // 嵌入 Person
    Address         // 嵌入 Address
    EmployeeID int
    Department string
}

// ========== 指针字段 ==========
type Config struct {
    Name   *string
    Port   *int
    Debug  *bool
}

// ========== 递归结构体 ==========
type TreeNode struct {
    Value    int
    Children []*TreeNode
}

// ========== 结构体方法 ==========
func (p Person) String() string {
    return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

func (p *Person) Birthday() {
    p.Age++
}

func main() {
    // ========== 结构体初始化 ==========
    // 零值初始化
    var p1 Person
    fmt.Printf("Zero value: %+v\n", p1)

    // 字面量初始化
    p2 := Person{
        Name:  "Alice",
        Age:   30,
        Email: "alice@example.com",
    }
    fmt.Printf("Literal: %+v\n", p2)

    // 按位置初始化（不推荐）
    p3 := Person{"Bob", 25, "bob@example.com"}
    fmt.Printf("Positional: %+v\n", p3)

    // new 创建
    p4 := new(Person)
    p4.Name = "Charlie"
    fmt.Printf("new(): %+v\n", *p4)

    // Go 1.26 new() 表达式
    p5 := new(Person{Name: "David", Age: 35})
    fmt.Printf("new(expr): %+v\n", *p5)

    // ========== 结构体字段访问 ==========
    fmt.Println("Name:", p2.Name)
    p2.Age = 31
    fmt.Printf("Updated: %+v\n", p2)

    // ========== 嵌入结构体 ==========
    emp := Employee{
        Person: Person{
            Name:  "Eve",
            Age:   28,
            Email: "eve@company.com",
        },
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
        },
        EmployeeID: 1001,
        Department: "Engineering",
    }

    // 直接访问嵌入字段
    fmt.Println("Employee Name:", emp.Name)
    fmt.Println("Employee City:", emp.City)

    // 也可以显式指定
    fmt.Println("Employee Name:", emp.Person.Name)

    // 调用嵌入类型的方法
    fmt.Println("Employee String():", emp.Person.String())

    // ========== 结构体标签 ==========
    user := User{
        ID:        1,
        Username:  "john_doe",
        Password:  "secret",
        Email:     "john@example.com",
        CreatedAt: "2024-01-01",
    }

    // 反射读取标签
    t := reflect.TypeOf(user)
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fmt.Printf("Field: %s, JSON tag: %s\n", field.Name, field.Tag.Get("json"))
    }

    // JSON 序列化
    jsonData, _ := json.Marshal(user)
    fmt.Println("JSON:", string(jsonData))

    // ========== 结构体比较 ==========
    s1 := Person{Name: "Alice", Age: 30}
    s2 := Person{Name: "Alice", Age: 30}
    s3 := Person{Name: "Bob", Age: 30}

    fmt.Println("s1 == s2:", s1 == s2)  // true
    fmt.Println("s1 == s3:", s1 == s3)  // false

    // ========== 结构体指针 ==========
    ptr := &Person{Name: "Frank", Age: 40}
    fmt.Println("Pointer access:", ptr.Name)  // 自动解引用

    // ========== 递归结构体 ==========
    root := &TreeNode{Value: 1}
    root.Children = []*TreeNode{
        {Value: 2, Children: []*TreeNode{{Value: 4}, {Value: 5}}},
        {Value: 3, Children: []*TreeNode{{Value: 6}}},
    }

    // 遍历树
    var traverse func(*TreeNode, int)
    traverse = func(node *TreeNode, depth int) {
        if node == nil {
            return
        }
        indent := ""
        for i := 0; i < depth; i++ {
            indent += "  "
        }
        fmt.Printf("%s%d\n", indent, node.Value)
        for _, child := range node.Children {
            traverse(child, depth+1)
        }
    }

    fmt.Println("Tree structure:")
    traverse(root, 0)

    // ========== 结构体方法 ==========
    person := Person{Name: "Grace", Age: 25}
    fmt.Println("Person:", person.String())
    person.Birthday()
    fmt.Println("After birthday:", person)
}
```

---

### 6.2 结构体嵌入与方法集

#### 概念定义

结构体嵌入允许一个结构体包含另一个结构体作为匿名字段，从而继承其方法。

#### 详细示例代码

```go
package main

import "fmt"

// ========== 基础类型 ==========
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "Some sound"
}

func (a Animal) Move() string {
    return "Moving"
}

// ========== 嵌入类型 ==========
type Dog struct {
    Animal          // 嵌入 Animal
    Breed   string
}

// 覆盖父类方法
func (d Dog) Speak() string {
    return "Woof!"
}

// 新方法
func (d Dog) Fetch() string {
    return "Fetching the ball"
}

// ========== 多层嵌入 ==========
type Puppy struct {
    Dog             // 嵌入 Dog（间接嵌入 Animal）
    Age     int
}

// ========== 指针嵌入 ==========
type Cat struct {
    *Animal         // 指针嵌入
    Color   string
}

func (c Cat) Speak() string {
    return "Meow!"
}

// ========== 方法提升 ==========
type Engine struct {
    Horsepower int
}

func (e Engine) Start() string {
    return "Engine starting"
}

type Car struct {
    Engine          // 嵌入 Engine
    Brand   string
}

func main() {
    // ========== 基本嵌入 ==========
    dog := Dog{
        Animal: Animal{Name: "Buddy"},
        Breed:  "Golden Retriever",
    }

    // 直接访问嵌入字段的方法
    fmt.Println("Dog Name:", dog.Name)        // 从 Animal 提升
    fmt.Println("Dog Speak:", dog.Speak())    // Dog 的 Speak
    fmt.Println("Dog Move:", dog.Move())      // 从 Animal 提升
    fmt.Println("Dog Fetch:", dog.Fetch())    // Dog 自己的方法

    // 显式访问嵌入字段
    fmt.Println("Animal Speak:", dog.Animal.Speak())  // Animal 的 Speak

    // ========== 多层嵌入 ==========
    puppy := Puppy{
        Dog: Dog{
            Animal: Animal{Name: "Tiny"},
            Breed:  "Labrador",
        },
        Age: 1,
    }

    fmt.Println("\nPuppy Name:", puppy.Name)      // 从 Animal 提升（通过 Dog）
    fmt.Println("Puppy Breed:", puppy.Breed)      // 从 Dog 提升
    fmt.Println("Puppy Age:", puppy.Age)          // 自己的字段
    fmt.Println("Puppy Speak:", puppy.Speak())    // Dog 的 Speak

    // ========== 指针嵌入 ==========
    cat := Cat{
        Animal: &Animal{Name: "Whiskers"},
        Color:  "Orange",
    }

    fmt.Println("\nCat Name:", cat.Name)
    fmt.Println("Cat Speak:", cat.Speak())
    fmt.Println("Cat Move:", cat.Move())

    // ========== 方法提升 ==========
    car := Car{
        Engine: Engine{Horsepower: 200},
        Brand:  "Toyota",
    }

    fmt.Println("\nCar Brand:", car.Brand)
    fmt.Println("Car Horsepower:", car.Horsepower)  // 从 Engine 提升
    fmt.Println("Car Start:", car.Start())          // 从 Engine 提升

    // ========== 方法集 ==========
    // 值接收者方法可以被值和指针调用
    // 指针接收者方法只能被指针调用（编译器会自动取地址）

    // ========== 接口实现 ==========
    type Speaker interface {
        Speak() string
    }

    var s Speaker = dog  // Dog 实现了 Speaker
    fmt.Println("\nSpeaker:", s.Speak())

    // 嵌入类型的接口实现也被提升
    var s2 Speaker = puppy  // Puppy 也实现了 Speaker（通过 Dog）
    fmt.Println("Speaker (puppy):", s2.Speak())
}
```

---

## 7. 包和导入机制

### 7.1 包声明

#### 概念定义

包是 Go 代码的组织单位，每个文件必须以 package 声明开头。

#### 包类型

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 包类型                                     │
├────────────────┬────────────────────────────────────────────────┤
│ main 包        │ 可执行程序的入口包                             │
├────────────────┼────────────────────────────────────────────────┤
│ 普通包         │ 以目录名命名的库包                             │
├────────────────┼────────────────────────────────────────────────┤
│ 内部包         │ internal 目录下的包，只能被父目录导入          │
├────────────────┼────────────────────────────────────────────────┤
│ 标准库包       │ Go 内置的包，如 fmt, os, net 等                │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
// ========== main 包示例 ==========
// file: main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

// ========== 普通包示例 ==========
// file: utils/utils.go
package utils

import "strings"

// 导出函数（大写开头）
func ToUpper(s string) string {
    return strings.ToUpper(s)
}

// 非导出函数（小写开头）
func helper() string {
    return "helper"
}

// ========== 内部包示例 ==========
// file: mypkg/internal/helper.go
package internal

// 只能被 mypkg 或其子包导入
func InternalHelper() string {
    return "internal helper"
}
```

---

### 7.2 导入机制

#### 概念定义

import 语句用于导入其他包，支持多种导入方式。

#### 导入方式

```
┌─────────────────────────────────────────────────────────────────┐
│                    Go 导入方式                                   │
├────────────────┬────────────────────────────────────────────────┤
│ 标准导入       │ import "fmt"                                   │
├────────────────┼────────────────────────────────────────────────┤
│ 别名导入       │ import f "fmt"                                 │
├────────────────┼────────────────────────────────────────────────┤
│ 点导入         │ import . "fmt"  // 直接使用包内标识符          │
├────────────────┼────────────────────────────────────────────────┤
│ 空白导入       │ import _ "fmt"  // 仅执行包的 init 函数         │
├────────────────┼────────────────────────────────────────────────┤
│ 分组导入       │ import ( "fmt" "os" )                          │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

// ========== 标准导入 ==========
import "fmt"
import "os"

// ========== 分组导入（推荐）==========
import (
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "time"
)

// ========== 别名导入 ==========
import (
    str "strings"           // 别名 str
    f "fmt"                 // 别名 f（不推荐，除非冲突）
)

// ========== 点导入（谨慎使用）==========
import (
    . "math"                // 直接使用 Pi, Sin 等
)

// ========== 空白导入 ==========
import (
    _ "net/http/pprof"      // 注册 pprof 处理器
    _ "image/png"           // 注册 PNG 解码器
)

// ========== 第三方包导入 ==========
import (
    // "github.com/gin-gonic/gin"
    // "gorm.io/gorm"
)

func main() {
    // 标准导入使用
    fmt.Println("Standard import")

    // 别名导入使用
    result := str.ToUpper("hello")
    f.Println("Upper:", result)

    // 点导入使用
    fmt.Println("Pi:", Pi)
    fmt.Println("Sin(0):", Sin(0))

    // 使用其他导入的包
    now := time.Now()
    fmt.Println("Time:", now)

    // 错误处理
    err := errors.New("example error")
    fmt.Println("Error:", err)

    // JSON 处理
    data := map[string]interface{}{"key": "value"}
    jsonBytes, _ := json.Marshal(data)
    fmt.Println("JSON:", string(jsonBytes))

    // 文件操作
    file, err := os.CreateTemp("", "example")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // IO 操作
    _, _ = io.WriteString(file, "Hello")

    _ = http.ListenAndServe
}
```

---

### 7.3 包初始化

#### 概念定义

包初始化在导入时执行，包括变量初始化和 init 函数。

#### 初始化顺序

```
┌─────────────────────────────────────────────────────────────────┐
│                    包初始化顺序                                  │
├────────────────┬────────────────────────────────────────────────┤
│ 1. 导入包      │ 递归初始化导入的包                             │
├────────────────┼────────────────────────────────────────────────┤
│ 2. 变量初始化  │ 按声明顺序初始化包级变量                       │
├────────────────┼────────────────────────────────────────────────┤
│ 3. init 函数   │ 按声明顺序执行 init 函数                       │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
// ========== 包初始化示例 ==========
// file: config/config.go
package config

import "fmt"

// 包级变量初始化
var Config = loadConfig()

func loadConfig() map[string]string {
    fmt.Println("Loading config...")
    return map[string]string{
        "host": "localhost",
        "port": "8080",
    }
}

// init 函数
func init() {
    fmt.Println("config package init")
}

func init() {
    fmt.Println("config package init 2")
}

// ========== 主包 ==========
// file: main.go
package main

import (
    "fmt"
    _ "config"  // 空白导入触发初始化
)

var globalVar = initGlobal()

func initGlobal() int {
    fmt.Println("Initializing globalVar")
    return 42
}

func init() {
    fmt.Println("main package init")
}

func main() {
    fmt.Println("main function")
}

// 输出顺序：
// Loading config...
// config package init
// config package init 2
// Initializing globalVar
// main package init
// main function
```

---

## 8. 新特性详解

### 8.1 new() 表达式特性（Go 1.26）

#### 概念定义

Go 1.26 允许 `new()` 接受表达式作为操作数，而不仅仅是类型。

#### 特性详解

```
┌─────────────────────────────────────────────────────────────────┐
│                    new() 表达式特性                              │
├────────────────┬────────────────────────────────────────────────┤
│ 传统用法       │ new(T) 返回 *T，值为零值                       │
├────────────────┼────────────────────────────────────────────────┤
│ 新用法         │ new(expr) 返回 *expr.Type，值为 expr 的值      │
├────────────────┼────────────────────────────────────────────────┤
│ 等价形式       │ new(expr) ≡ &expr（当 expr 是可寻址的）        │
├────────────────┼────────────────────────────────────────────────┤
│ 主要用途       │ 简化可选字段初始化，JSON/protobuf 序列化       │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "encoding/json"
    "fmt"
)

// ========== 传统方式 vs 新方式对比 ==========

// 传统方式创建指针
func traditionalWay() {
    // 需要中间变量
    val := int64(300)
    ptr := &val
    fmt.Println("Traditional:", *ptr)
}

// Go 1.26 新方式
func newWay() {
    // 直接创建指向值的指针
    ptr := new(int64(300))
    fmt.Println("New way:", *ptr)
}

// ========== JSON 序列化场景 ==========
type User struct {
    Name    *string `json:"name,omitempty"`
    Age     *int    `json:"age,omitempty"`
    Active  *bool   `json:"active,omitempty"`
}

// 传统方式创建 User（需要辅助函数或中间变量）
func createUserTraditional() *User {
    name := "Alice"
    age := 30
    active := true

    return &User{
        Name:   &name,
        Age:    &age,
        Active: &active,
    }
}

// Go 1.26 新方式（简洁）
func createUserNew() *User {
    return &User{
        Name:   new(string("Alice")),
        Age:    new(int(30)),
        Active: new(bool(true)),
    }
}

// ========== Protobuf 场景 ==========
type Message struct {
    ID       *int64  `json:"id,omitempty"`
    Content  *string `json:"content,omitempty"`
    Priority *int32  `json:"priority,omitempty"`
}

func createMessageTraditional(id int64, content string, priority int32) *Message {
    return &Message{
        ID:       &id,
        Content:  &content,
        Priority: &priority,
    }
}

func createMessageNew(id int64, content string, priority int32) *Message {
    return &Message{
        ID:       new(int64(id)),
        Content:  new(string(content)),
        Priority: new(int32(priority)),
    }
}

// ========== 更多表达式示例 ==========
func expressionExamples() {
    // 算术表达式
    ptr1 := new(10 + 20)
    fmt.Println("new(10 + 20):", *ptr1)  // 30

    // 类型转换
    ptr2 := new(float64(42))
    fmt.Println("new(float64(42)):", *ptr2)  // 42.0

    // 复合字面量
    ptr3 := new([3]int{1, 2, 3})
    fmt.Println("new([3]int{1,2,3}):", *ptr3)

    // 结构体字面量
    type Point struct{ X, Y int }
    ptr4 := new(Point{X: 10, Y: 20})
    fmt.Println("new(Point{X:10,Y:20}):", *ptr4)

    // 函数调用结果
    ptr5 := new(len("hello"))
    fmt.Println("new(len(\"hello\")):", *ptr5)  // 5
}

// ========== 复杂场景：配置构建 ==========
type ServerConfig struct {
    Host     *string `json:"host,omitempty"`
    Port     *int    `json:"port,omitempty"`
    Timeout  *int    `json:"timeout,omitempty"`
    MaxConns *int    `json:"max_conns,omitempty"`
    TLS      *bool   `json:"tls,omitempty"`
}

// 传统方式（冗长）
func buildConfigTraditional() *ServerConfig {
    host := "localhost"
    port := 8080
    timeout := 30
    maxConns := 100
    tls := true

    return &ServerConfig{
        Host:     &host,
        Port:     &port,
        Timeout:  &timeout,
        MaxConns: &maxConns,
        TLS:      &tls,
    }
}

// Go 1.26 新方式（简洁）
func buildConfigNew() *ServerConfig {
    return &ServerConfig{
        Host:     new(string("localhost")),
        Port:     new(int(8080)),
        Timeout:  new(int(30)),
        MaxConns: new(int(100)),
        TLS:      new(bool(true)),
    }
}

// ========== 泛型辅助函数 ==========
// 可以创建泛型辅助函数来简化代码
func Ptr[T any](v T) *T {
    return &v
}

// 使用泛型辅助函数
func buildConfigWithHelper() *ServerConfig {
    return &ServerConfig{
        Host:     Ptr("localhost"),
        Port:     Ptr(8080),
        Timeout:  Ptr(30),
        MaxConns: Ptr(100),
        TLS:      Ptr(true),
    }
}

func main() {
    fmt.Println("=== new() Expression Feature Demo ===\n")

    traditionalWay()
    newWay()

    fmt.Println("\n=== JSON Serialization ===")
    user1 := createUserTraditional()
    user2 := createUserNew()

    json1, _ := json.Marshal(user1)
    json2, _ := json.Marshal(user2)

    fmt.Println("Traditional:", string(json1))
    fmt.Println("New way:", string(json2))

    fmt.Println("\n=== Expression Examples ===")
    expressionExamples()

    fmt.Println("\n=== Config Building ===")
    config1 := buildConfigTraditional()
    config2 := buildConfigNew()
    config3 := buildConfigWithHelper()

    json3, _ := json.Marshal(config1)
    json4, _ := json.Marshal(config2)
    json5, _ := json.Marshal(config3)

    fmt.Println("Traditional:", string(json3))
    fmt.Println("New way:", string(json4))
    fmt.Println("With helper:", string(json5))
}
```

---

### 8.2 自引用泛型约束（Go 1.26）

#### 概念定义

Go 1.26 允许泛型类型在自己的类型参数列表中引用自己，实现 F-bounded 多态性。

#### 特性详解

```
┌─────────────────────────────────────────────────────────────────┐
│                    自引用泛型约束                                │
├────────────────┬────────────────────────────────────────────────┤
│ 概念           │ 类型参数约束可以引用包含该类型参数的接口       │
├────────────────┼────────────────────────────────────────────────┤
│ 用途           │ 实现 F-bounded 多态性，递归类型约束            │
├────────────────┼────────────────────────────────────────────────┤
│ 典型模式       │ type I[T I[T]] interface { ... }               │
├────────────────┼────────────────────────────────────────────────┤
│ 应用场景       │ 比较器、加法器、树结构、递归数据结构           │
└────────────────┴────────────────────────────────────────────────┘
```

#### 详细示例代码

```go
package main

import (
    "fmt"
)

// ========== 基本自引用约束 ==========
// Adder 接口：可以与自己相加的类型
type Adder[A Adder[A]] interface {
    Add(other A) A
}

// MyInt 实现 Adder
type MyInt int

func (a MyInt) Add(other MyInt) MyInt {
    return a + other
}

// ========== 比较器接口 ==========
// Comparer 接口：可以与自己比较的类型
type Comparer[C Comparer[C]] interface {
    Compare(other C) int  // 返回 -1, 0, 1
    Equal(other C) bool
}

// Ordered 扩展 Comparer
type Ordered[O Ordered[O]] interface {
    Comparer[O]
    Less(other O) bool
}

// Int 实现 Ordered
type Int int

func (a Int) Compare(b Int) int {
    if a < b {
        return -1
    } else if a > b {
        return 1
    }
    return 0
}

func (a Int) Equal(b Int) bool {
    return a == b
}

func (a Int) Less(b Int) bool {
    return a < b
}

// ========== 使用自引用约束的泛型函数 ==========
func SumAdders[T Adder[T]](items []T) T {
    var sum T
    for _, item := range items {
        sum = sum.Add(item)
    }
    return sum
}

func MaxOrdered[T Ordered[T]](items []T) T {
    if len(items) == 0 {
        var zero T
        return zero
    }
    max := items[0]
    for _, item := range items[1:] {
        if max.Less(item) {
            max = item
        }
    }
    return max
}

// ========== 树结构示例 ==========
// TreeNode 接口：可以形成树的类型
type TreeNode[T TreeNode[T]] interface {
    GetChildren() []T
    GetValue() int
}

// BinaryTree 实现 TreeNode
type BinaryTree struct {
    Value int
    Left  *BinaryTree
    Right *BinaryTree
}

func (b *BinaryTree) GetChildren() []*BinaryTree {
    children := []*BinaryTree{}
    if b.Left != nil {
        children = append(children, b.Left)
    }
    if b.Right != nil {
        children = append(children, b.Right)
    }
    return children
}

func (b *BinaryTree) GetValue() int {
    return b.Value
}

// 泛型树遍历
func TraverseTree[T TreeNode[T]](node T, depth int) {
    indent := ""
    for i := 0; i < depth; i++ {
        indent += "  "
    }
    fmt.Printf("%s%d\n", indent, node.GetValue())

    // 这里需要类型断言，因为 GetChildren 返回 []T
    // 实际使用时可能需要更复杂的处理
}

// ========== 表达式求值示例 ==========
// Expr 接口：表达式类型
type Expr[E Expr[E]] interface {
    Eval() int
    String() string
}

// 常量表达式
type Const struct {
    Value int
}

func (c Const) Eval() int    { return c.Value }
func (c Const) String() string { return fmt.Sprintf("%d", c.Value) }

// 加法表达式
type Add struct {
    Left, Right Expr[Expr[any]]  // 简化处理
}

func (a Add) Eval() int {
    return a.Left.Eval() + a.Right.Eval()
}

func (a Add) String() string {
    return fmt.Sprintf("(%s + %s)", a.Left.String(), a.Right.String())
}

// ========== 克隆接口示例 ==========
// Cloner 接口：可以克隆自己的类型
type Cloner[C Cloner[C]] interface {
    Clone() C
}

// CloneSlice 克隆切片
type CloneSlice[T Cloner[T]] []T

func (cs CloneSlice[T]) Clone() CloneSlice[T] {
    result := make(CloneSlice[T], len(cs))
    for i, item := range cs {
        result[i] = item.Clone()
    }
    return result
}

// ========== 更复杂的自引用约束 ==========
// Container 接口：可以包含相同类型元素的容器
type Container[C Container[C], E any] interface {
    Add(element E)
    Get(index int) E
    Size() int
}

// 泛型列表
type List[T any] struct {
    items []T
}

func (l *List[T]) Add(item T) {
    l.items = append(l.items, item)
}

func (l *List[T]) Get(index int) T {
    return l.items[index]
}

func (l *List[T]) Size() int {
    return len(l.items)
}

// ========== 使用示例 ==========
func main() {
    fmt.Println("=== Self-Referential Generic Constraints ===\n")

    // ========== Adder 示例 ==========
    fmt.Println("--- Adder Example ---")
    myInts := []MyInt{1, 2, 3, 4, 5}
    sum := SumAdders(myInts)
    fmt.Printf("Sum of %v = %d\n", myInts, sum)

    // ========== Ordered 示例 ==========
    fmt.Println("\n--- Ordered Example ---")
    ints := []Int{3, 1, 4, 1, 5, 9, 2, 6}
    max := MaxOrdered(ints)
    fmt.Printf("Max of %v = %d\n", ints, max)

    // ========== 树遍历示例 ==========
    fmt.Println("\n--- Tree Example ---")
    tree := &BinaryTree{
        Value: 1,
        Left: &BinaryTree{
            Value: 2,
            Left:  &BinaryTree{Value: 4},
            Right: &BinaryTree{Value: 5},
        },
        Right: &BinaryTree{
            Value: 3,
            Left:  &BinaryTree{Value: 6},
        },
    }

    TraverseTree(tree, 0)

    // ========== 表达式求值示例 ==========
    fmt.Println("\n--- Expression Example ---")
    expr := Add{
        Left:  Const{Value: 10},
        Right: Const{Value: 20},
    }
    fmt.Printf("Expression: %s = %d\n", expr.String(), expr.Eval())

    // ========== 列表示例 ==========
    fmt.Println("\n--- List Example ---")
    list := &List[int]{}
    list.Add(1)
    list.Add(2)
    list.Add(3)
    fmt.Printf("List size: %d\n", list.Size())
    fmt.Printf("List[1]: %d\n", list.Get(1))

    // ========== 复杂自引用约束说明 ==========
    fmt.Println("\n--- Complex Self-Reference ---")
    fmt.Println("type Adder[A Adder[A]] interface { Add(A) A }")
    fmt.Println("This allows A to reference itself in the constraint")
    fmt.Println("Enabling F-bounded polymorphism in Go!")
}
```

---

## 9. 思维导图与决策树

### 9.1 Go 语法特性思维导图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Go 1.26.1 语法特性思维导图                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐                                                            │
│  │  词法元素   │                                                            │
│  ├─────────────┤                                                            │
│  │ • 标识符    │──→ 字母/下划线开头，大小写敏感，Unicode支持               │
│  │ • 关键字    │──→ 25个，分四类（声明/复合/控制/其他）                    │
│  │ • 运算符    │──→ 算术/比较/逻辑/位运算/赋值/其他                        │
│  │ • 字面量    │──→ 整数/浮点/复数/字符串/符文/布尔/复合                   │
│  └─────────────┘                                                            │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                            │
│  │  类型系统   │                                                            │
│  ├─────────────┤                                                            │
│  │ • 基本类型  │──→ bool/int/float/complex/string/rune/byte                │
│  │ • 复合类型  │──→ array/slice/map/chan/ptr/func/struct/interface         │
│  │ • 接口类型  │──→ 方法集，隐式实现，类型断言，类型switch                 │
│  │ • 泛型类型  │──→ 类型参数，类型约束，类型推断（Go 1.18+）               │
│  │ • 自引用约束│──→ F-bounded 多态性（Go 1.26 新特性）                     │
│  └─────────────┘                                                            │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                            │
│  │  声明语句   │                                                            │
│  ├─────────────┤                                                            │
│  │ • 变量声明  │──→ var/:=，多变量，分组声明                               │
│  │ • 常量声明  │──→ const，iota，无类型常量                                │
│  │ • new()新特性│──→ 接受表达式操作数（Go 1.26 新特性）                    │
│  └─────────────┘                                                            │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                            │
│  │  控制流     │                                                            │
│  ├─────────────┤                                                            │
│  │ • if        │──→ 条件判断，前置语句，嵌套                               │
│  │ • for       │──→ 完整/条件/无限/range，break/continue/标签              │
│  │ • switch    │──→ 表达式/类型/无表达式，fallthrough                      │
│  │ • select    │──→ 通道选择，随机选择，超时处理                           │
│  └─────────────┘                                                            │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                            │
│  │  函数/方法  │                                                            │
│  ├─────────────┤                                                            │
│  │ • 函数      │──→ 多返回值，命名返回值，变长参数，匿名函数               │
│  │ • 方法      │──→ 值/指针接收者，方法集，嵌入提升                        │
│  │ • defer     │──→ 延迟执行，LIFO顺序                                     │
│  │ • panic/    │──→ 异常处理，恢复机制                                     │
│  │   recover   │                                                            │
│  └─────────────┘                                                            │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                            │
│  │  结构体     │                                                            │
│  ├─────────────┤                                                            │
│  │ • 声明      │──→ 命名字段，嵌入字段，标签                               │
│  │ • 初始化    │──→ 零值，字面量，new()                                    │
│  │ • 方法      │──→ 值/指针接收者，方法提升                                │
│  │ • 嵌入      │──→ 多层嵌入，方法继承                                     │
│  └─────────────┘                                                            │
│         │                                                                   │
│         ▼                                                                   │
│  ┌─────────────┐                                                            │
│  │  包/导入    │                                                            │
│  ├─────────────┤                                                            │
│  │ • 包声明    │──→ main/普通/internal/标准库                              │
│  │ • 导入      │──→ 标准/别名/点/空白/分组                                 │
│  │ • 初始化    │──→ 导入顺序，变量初始化，init函数                         │
│  └─────────────┘                                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 9.2 类型选择决策树

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Go 类型选择决策树                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  需要存储什么类型的数据？                                                    │
│       │                                                                     │
│       ├──→ 简单值                                                          │
│       │     │                                                              │
│       │     ├──→ 布尔值 ────────────────────────────→ bool                 │
│       │     │                                                              │
│       │     ├──→ 整数                                                       │
│       │     │     ├──→ 有符号 ──→ 范围小？─→ 是 ──→ int8/int16            │
│       │     │     │                 │                                      │
│       │     │     │                 └──→ 否 ──→ int32/int64/int           │
│       │     │     │                                                        │
│       │     │     └──→ 无符号 ──→ 范围小？─→ 是 ──→ uint8/uint16          │
│       │     │                       │                                      │
│       │     │                       └──→ 否 ──→ uint32/uint64/uint        │
│       │     │                                                              │
│       │     ├──→ 浮点数 ──→ 精度要求高？─→ 是 ──→ float64                 │
│       │     │                 │                                            │
│       │     │                 └──→ 否 ──→ float32                          │
│       │     │                                                              │
│       │     ├──→ 复数 ──→ 精度要求高？─→ 是 ──→ complex128                │
│       │     │               │                                              │
│       │     │               └──→ 否 ──→ complex64                          │
│       │     │                                                              │
│       │     └──→ 文本                                                       │
│       │           ├──→ 单个字符 ──→ rune (Unicode) / byte (ASCII)         │
│       │           │                                                        │
│       │           └──→ 字符串 ──→ string                                   │
│       │                                                                    │
│       ├──→ 复合数据结构                                                    │
│       │     │                                                              │
│       │     ├──→ 固定长度序列 ──→ array [n]T                              │
│       │     │                                                              │
│       │     ├──→ 动态长度序列 ──→ slice []T                               │
│       │     │                                                              │
│       │     ├──→ 键值对集合 ──→ map[K]V                                   │
│       │     │                                                              │
│       │     ├──→ 结构化数据 ──→ struct                                    │
│       │     │                                                              │
│       │     └──→ 异构数据 ──→ interface                                   │
│       │                                                                    │
│       ├──→ 并发通信                                                        │
│       │     └──→ 通道 ──→ 缓冲？─→ 是 ──→ chan T (带缓冲)                 │
│       │                   │                                                │
│       │                   └──→ 否 ──→ chan T (无缓冲)                      │
│       │                                                                    │
│       ├──→ 引用/间接访问                                                   │
│       │     └──→ 指针 ──→ *T                                              │
│       │                                                                    │
│       ├──→ 可调用实体                                                      │
│       │     └──→ 函数 ──→ func(...) ...                                   │
│       │                                                                    │
│       └──→ 多种类型通用代码                                                │
│             └──→ 泛型 ──→ [T Constraint]                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

### 9.3 控制流选择决策树

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go 控制流选择决策树                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  需要实现什么控制逻辑？                                                      │
│       │                                                                     │
│       ├──→ 条件分支                                                        │
│       │     │                                                              │
│       │     ├──→ 单一条件 ──────────────────────────→ if                  │
│       │     │                                                              │
│       │     ├──→ 两个分支 ──────────────────────────→ if-else             │
│       │     │                                                              │
│       │     ├──→ 多分支（离散值）                                               │
│       │     │     ├──→ 需要穿透？─→ 是 ──→ switch + fallthrough           │
│       │     │     │                                                            │
│       │     │     └──→ 否 ──→ switch                                         │
│       │     │                                                              │
│       │     └──→ 多分支（范围/复杂条件）                                          │
│       │           └──→ switch（无表达式）或 if-else-if                       │
│       │                                                                    │
│       ├──→ 循环                                                            │
│       │     │                                                              │
│       │     ├──→ 已知迭代次数 ──────────────────────→ for i := 0; ...     │
│       │     │                                                              │
│       │     ├──→ 条件循环 ──────────────────────────→ for condition       │
│       │     │                                                              │
│       │     ├──→ 无限循环 ──────────────────────────→ for { ... }         │
│       │     │                                                              │
│       │     └──→ 遍历集合                                                    │
│       │           ├──→ 数组/切片/字符串/map/通道 ──→ range                 │
│       │           │                                                        │
│       │           └──→ 需要索引？─→ 是 ──→ for i, v := range              │
│       │                         │                                          │
│       │                         └──→ 否 ──→ for _, v := range             │
│       │                                                                    │
│       ├──→ 并发通信                                                        │
│       │     └──→ 多通道操作 ────────────────────────→ select              │
│       │                                                                    │
│       └──→ 跳转                                                            │
│             ├──→ 跳出循环 ──→ break                                        │
│             ├──→ 跳出多层循环 ──→ break label                              │
│             ├──→ 跳过迭代 ──→ continue                                     │
│             └──→ 无条件跳转 ──→ goto（不推荐）                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 10. 总结

### Go 1.26.1 语法特性总览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go 1.26.1 核心特性总结                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  【新特性】                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ 1. new() 表达式特性                                                 │   │
│  │    • new(expr) 接受表达式作为操作数                                 │   │
│  │    • 简化可选字段初始化                                             │   │
│  │    • 特别适用于 JSON/protobuf 序列化                                │   │
│  │                                                                     │   │
│  │ 2. 自引用泛型约束                                                   │   │
│  │    • type I[T I[T]] interface { ... }                               │   │
│  │    • 实现 F-bounded 多态性                                          │   │
│  │    • 支持递归类型约束                                               │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  【核心语法】                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ • 简洁的语法设计，25 个关键字                                       │   │
│  │ • 强大的类型系统，支持泛型编程                                      │   │
│  │ • 显式错误处理，多返回值                                            │   │
│  │ • 内置并发支持（goroutine + channel）                               │   │
│  │ • 垃圾回收，内存安全                                                │   │
│  │ • 快速编译，静态链接                                                │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  【最佳实践】                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │ • 使用 gofmt 格式化代码                                             │   │
│  │ • 使用 golint/vet 检查代码                                          │   │
│  │ • 错误处理：if err != nil { return err }                            │   │
│  │ • 接口设计：小接口优于大接口                                        │   │
│  │ • 并发安全：使用 channel 而非共享内存                               │   │
│  │ • 文档注释：导出标识符必须注释                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 附录：Go 关键字速查表

| 类别 | 关键字 | 用途 |
|------|--------|------|
| 声明 | `const` | 声明常量 |
| 声明 | `func` | 声明函数 |
| 声明 | `import` | 导入包 |
| 声明 | `package` | 声明包 |
| 声明 | `type` | 声明类型 |
| 声明 | `var` | 声明变量 |
| 复合类型 | `chan` | 声明通道 |
| 复合类型 | `interface` | 声明接口 |
| 复合类型 | `map` | 声明映射 |
| 复合类型 | `struct` | 声明结构体 |
| 控制流 | `break` | 跳出循环/switch |
| 控制流 | `case` | switch/select 分支 |
| 控制流 | `continue` | 跳过迭代 |
| 控制流 | `default` | 默认分支 |
| 控制流 | `else` | if 的否定分支 |
| 控制流 | `fallthrough` | switch 穿透 |
| 控制流 | `for` | 循环 |
| 控制流 | `goto` | 跳转（不推荐） |
| 控制流 | `if` | 条件判断 |
| 控制流 | `range` | 遍历 |
| 控制流 | `return` | 返回 |
| 控制流 | `switch` | 多路分支 |
| 其他 | `defer` | 延迟执行 |
| 其他 | `go` | 启动 goroutine |
| 其他 | `select` | 通道选择 |

---

*文档生成时间：基于 Go 1.26.1 语法规范*
