# Go 1.23 语言特性全面指南

> 本文档基于 **Go 1.23** 最新版本，全面梳理Go语言的所有语法特性、类型系统、控制结构。
> 涵盖 Go 1.23 新特性：迭代器（range-over-func）、iter包、unique包、structs包、Timer改进等。
> 每个特性包含：概念定义、语法形式、属性关系、形式论证、完整示例、反例说明、最佳实践。

---

## 目录

- [Go 1.23 语言特性全面指南](#go-123-语言特性全面指南)
  - [目录](#目录)
  - [1. 基础语法特性](#1-基础语法特性)
    - [1.1 变量声明](#11-变量声明)
      - [概念定义](#概念定义)
      - [语法形式](#语法形式)
      - [属性关系](#属性关系)
      - [形式论证](#形式论证)
      - [完整示例](#完整示例)
      - [反例说明](#反例说明)
      - [最佳实践](#最佳实践)
    - [1.2 常量与iota](#12-常量与iota)
      - [概念定义](#概念定义-1)
      - [语法形式](#语法形式-1)
      - [属性关系](#属性关系-1)
      - [形式论证](#形式论证-1)
      - [完整示例](#完整示例-1)
      - [反例说明](#反例说明-1)
      - [最佳实践](#最佳实践-1)
    - [1.3 基本数据类型及零值](#13-基本数据类型及零值)
      - [概念定义](#概念定义-2)
      - [类型体系](#类型体系)
      - [零值表](#零值表)
      - [完整示例](#完整示例-2)
      - [反例说明](#反例说明-2)
      - [最佳实践](#最佳实践-2)
    - [1.4 类型推断机制](#14-类型推断机制)
      - [概念定义](#概念定义-3)
      - [推断规则](#推断规则)
      - [完整示例](#完整示例-3)
      - [反例说明](#反例说明-3)
      - [最佳实践](#最佳实践-3)
    - [1.5 类型别名与类型定义](#15-类型别名与类型定义)
      - [概念定义](#概念定义-4)
      - [语法形式](#语法形式-2)
      - [属性对比](#属性对比)
      - [完整示例](#完整示例-4)
      - [反例说明](#反例说明-4)
      - [最佳实践](#最佳实践-4)
  - [2. 复合类型](#2-复合类型)
    - [2.1 数组与切片（Slice）](#21-数组与切片slice)
      - [概念定义](#概念定义-5)
      - [语法形式](#语法形式-3)
      - [内部结构](#内部结构)
      - [数组vs切片对比](#数组vs切片对比)
      - [完整示例](#完整示例-5)
      - [反例说明](#反例说明-5)
      - [最佳实践](#最佳实践-5)
    - [2.2 映射（Map）](#22-映射map)
      - [概念定义](#概念定义-6)
      - [语法形式](#语法形式-4)
      - [内部结构](#内部结构-1)
      - [可比较类型（可作为map键）](#可比较类型可作为map键)
      - [完整示例](#完整示例-6)
      - [反例说明](#反例说明-6)
      - [最佳实践](#最佳实践-6)
    - [2.3 结构体（Struct）](#23-结构体struct)
      - [概念定义](#概念定义-7)
      - [语法形式](#语法形式-5)
      - [字段标签（Tag）](#字段标签tag)
      - [完整示例](#完整示例-7)
      - [反例说明](#反例说明-7)
      - [最佳实践](#最佳实践-7)
    - [2.4 嵌套类型与匿名类型](#24-嵌套类型与匿名类型)
      - [概念定义](#概念定义-8)
      - [语法形式](#语法形式-6)
      - [完整示例](#完整示例-8)
      - [反例说明](#反例说明-8)
      - [最佳实践](#最佳实践-8)
  - [3. 控制结构](#3-控制结构)
    - [3.1 if/else 及初始化语句](#31-ifelse-及初始化语句)
      - [概念定义](#概念定义-9)
      - [语法形式](#语法形式-7)
      - [完整示例](#完整示例-9)
      - [反例说明](#反例说明-9)
      - [最佳实践](#最佳实践-9)
    - [3.2 switch语句](#32-switch语句)
      - [概念定义](#概念定义-10)
      - [语法形式](#语法形式-8)
      - [完整示例](#完整示例-10)
      - [反例说明](#反例说明-10)
      - [最佳实践](#最佳实践-10)
    - [3.3 for循环（三种形式）](#33-for循环三种形式)
      - [概念定义](#概念定义-11)
      - [语法形式](#语法形式-9)
      - [完整示例](#完整示例-11)
      - [反例说明](#反例说明-11)
      - [最佳实践](#最佳实践-11)
    - [3.4 range遍历](#34-range遍历)
      - [概念定义](#概念定义-12)
      - [语法形式](#语法形式-10)
      - [遍历规则表](#遍历规则表)
      - [完整示例](#完整示例-12)
      - [反例说明](#反例说明-12)
      - [最佳实践](#最佳实践-12)
    - [3.5 goto、break、continue](#35-gotobreakcontinue)
      - [概念定义](#概念定义-13)
      - [语法形式](#语法形式-11)
      - [完整示例](#完整示例-13)
      - [反例说明](#反例说明-13)
      - [最佳实践](#最佳实践-13)
  - [4. 函数与方法](#4-函数与方法)
    - [4.1 函数定义与多返回值](#41-函数定义与多返回值)
      - [概念定义](#概念定义-14)
      - [语法形式](#语法形式-12)
      - [完整示例](#完整示例-14)
      - [反例说明](#反例说明-14)
      - [最佳实践](#最佳实践-14)
    - [4.2 变参函数](#42-变参函数)
      - [概念定义](#概念定义-15)
      - [语法形式](#语法形式-13)
      - [完整示例](#完整示例-15)
      - [反例说明](#反例说明-15)
      - [最佳实践](#最佳实践-15)
    - [4.3 匿名函数与闭包](#43-匿名函数与闭包)
      - [概念定义](#概念定义-16)
      - [语法形式](#语法形式-14)
      - [完整示例](#完整示例-16)
      - [反例说明](#反例说明-16)
      - [最佳实践](#最佳实践-16)
    - [4.4 方法接收者](#44-方法接收者)
      - [概念定义](#概念定义-17)
      - [语法形式](#语法形式-15)
      - [值接收者 vs 指针接收者](#值接收者-vs-指针接收者)
      - [完整示例](#完整示例-17)
      - [反例说明](#反例说明-17)
      - [最佳实践](#最佳实践-17)
    - [4.5 方法集](#45-方法集)
      - [概念定义](#概念定义-18)
      - [方法集规则](#方法集规则)
      - [完整示例](#完整示例-18)
      - [反例说明](#反例说明-18)
      - [最佳实践](#最佳实践-18)
  - [5. 接口系统](#5-接口系统)
    - [5.1 接口定义与实现](#51-接口定义与实现)
      - [概念定义](#概念定义-19)
      - [语法形式](#语法形式-16)
      - [实现规则](#实现规则)
      - [完整示例](#完整示例-19)
      - [反例说明](#反例说明-19)
      - [最佳实践](#最佳实践-19)
    - [5.2 空接口interface{}](#52-空接口interface)
      - [概念定义](#概念定义-20)
      - [内部表示](#内部表示)
      - [完整示例](#完整示例-20)
      - [反例说明](#反例说明-20)
      - [最佳实践](#最佳实践-20)
    - [5.3 类型断言与类型切换](#53-类型断言与类型切换)
      - [概念定义](#概念定义-21)
      - [语法形式](#语法形式-17)
      - [完整示例](#完整示例-21)
      - [反例说明](#反例说明-21)
      - [最佳实践](#最佳实践-21)
    - [5.4 接口嵌套与组合](#54-接口嵌套与组合)
      - [概念定义](#概念定义-22)
      - [语法形式](#语法形式-18)
      - [完整示例](#完整示例-22)
      - [反例说明](#反例说明-22)
      - [最佳实践](#最佳实践-22)
  - [6. 指针与内存](#6-指针与内存)
    - [6.1 指针操作](#61-指针操作)
      - [概念定义](#概念定义-23)
      - [语法形式](#语法形式-19)
      - [完整示例](#完整示例-23)
      - [反例说明](#反例说明-23)
      - [最佳实践](#最佳实践-23)
    - [6.2 new与make区别](#62-new与make区别)
      - [概念定义](#概念定义-24)
      - [对比表](#对比表)
      - [完整示例](#完整示例-24)
      - [反例说明](#反例说明-24)
      - [最佳实践](#最佳实践-24)
    - [6.3 逃逸分析](#63-逃逸分析)
      - [概念定义](#概念定义-25)
      - [逃逸场景](#逃逸场景)
      - [完整示例](#完整示例-25)
      - [反例说明](#反例说明-25)
      - [最佳实践](#最佳实践-25)
  - [7. 错误处理](#7-错误处理)
    - [7.1 error接口](#71-error接口)
      - [概念定义](#概念定义-26)
      - [语法形式](#语法形式-20)
      - [完整示例](#完整示例-26)
      - [反例说明](#反例说明-26)
      - [最佳实践](#最佳实践-26)
    - [7.2 panic与recover](#72-panic与recover)
      - [概念定义](#概念定义-27)
      - [语法形式](#语法形式-21)
      - [完整示例](#完整示例-27)
      - [反例说明](#反例说明-27)
      - [最佳实践](#最佳实践-27)
    - [7.3 错误链与包装](#73-错误链与包装)
      - [概念定义](#概念定义-28)
      - [语法形式](#语法形式-22)
      - [完整示例](#完整示例-28)
      - [反例说明](#反例说明-28)
      - [最佳实践](#最佳实践-28)
  - [8. 包管理与模块](#8-包管理与模块)
    - [8.1 GOPATH到Go Modules演进](#81-gopath到go-modules演进)
      - [概念定义](#概念定义-29)
      - [完整示例](#完整示例-29)
      - [反例说明](#反例说明-29)
      - [最佳实践](#最佳实践-29)
    - [8.2 go.mod与go.sum](#82-gomod与gosum)
      - [概念定义](#概念定义-30)
      - [go.mod指令](#gomod指令)
      - [完整示例](#完整示例-30)
      - [go.sum示例](#gosum示例)
      - [常用命令](#常用命令)
      - [反例说明](#反例说明-30)
      - [最佳实践](#最佳实践-30)
    - [8.3 语义化版本控制](#83-语义化版本控制)
      - [概念定义](#概念定义-31)
      - [版本选择规则](#版本选择规则)
      - [完整示例](#完整示例-31)
      - [版本兼容性](#版本兼容性)
      - [反例说明](#反例说明-31)
      - [最佳实践](#最佳实践-31)
    - [8.4 私有模块与代理](#84-私有模块与代理)
      - [概念定义](#概念定义-32)
      - [配置示例](#配置示例)
      - [私有模块配置](#私有模块配置)
      - [完整示例](#完整示例-32)
      - [反例说明](#反例说明-32)
      - [最佳实践](#最佳实践-32)
    - [8.5 工作区模式（Workspace）](#85-工作区模式workspace)
      - [概念定义](#概念定义-33)
      - [语法形式](#语法形式-23)
      - [完整示例](#完整示例-33)
      - [常用命令](#常用命令-1)
      - [反例说明](#反例说明-33)
      - [最佳实践](#最佳实践-33)
  - [9. 标准库核心](#9-标准库核心)
    - [9.1 fmt、os、io包](#91-fmtosio包)
      - [概念定义](#概念定义-34)
      - [fmt包核心函数](#fmt包核心函数)
      - [os包核心功能](#os包核心功能)
      - [io包核心接口](#io包核心接口)
      - [完整示例](#完整示例-34)
      - [反例说明](#反例说明-34)
      - [最佳实践](#最佳实践-34)
    - [9.2 net/http包架构](#92-nethttp包架构)
      - [概念定义](#概念定义-35)
      - [核心接口](#核心接口)
      - [完整示例](#完整示例-35)
      - [反例说明](#反例说明-35)
      - [最佳实践](#最佳实践-35)
    - [9.3 encoding/json](#93-encodingjson)
      - [概念定义](#概念定义-36)
      - [核心函数](#核心函数)
      - [结构体标签](#结构体标签)
      - [完整示例](#完整示例-36)
      - [反例说明](#反例说明-36)
      - [最佳实践](#最佳实践-36)
    - [9.4 context包](#94-context包)
      - [概念定义](#概念定义-37)
      - [核心接口和函数](#核心接口和函数)
      - [完整示例](#完整示例-37)
      - [反例说明](#反例说明-37)
      - [最佳实践](#最佳实践-37)
    - [9.5 sync包基础](#95-sync包基础)
      - [概念定义](#概念定义-38)
      - [核心类型](#核心类型)
      - [完整示例](#完整示例-38)
      - [反例说明](#反例说明-38)
      - [最佳实践](#最佳实践-38)
  - [附录：Go 1.22+ 新特性](#附录go-122-新特性)
    - [整数范围循环（Go 1.22）](#整数范围循环go-122)
    - [泛型增强（Go 1.18+）](#泛型增强go-118)
    - [标准库新增](#标准库新增)
  - [10. Go 1.23 新特性](#10-go-123-新特性)
    - [10.1 迭代器（range-over-func）](#101-迭代器range-over-func)
      - [概念定义](#概念定义-39)
      - [语法形式](#语法形式-24)
      - [完整示例](#完整示例-39)
      - [反例说明](#反例说明-39)
      - [最佳实践](#最佳实践-39)
    - [10.2 iter包](#102-iter包)
      - [概念定义](#概念定义-40)
      - [核心类型](#核心类型-1)
      - [完整示例](#完整示例-40)
      - [与slices和maps包集成](#与slices和maps包集成)
      - [最佳实践](#最佳实践-40)
    - [10.3 unique包](#103-unique包)
      - [概念定义](#概念定义-41)
      - [核心概念](#核心概念)
      - [完整示例](#完整示例-41)
      - [反例说明](#反例说明-40)
      - [最佳实践](#最佳实践-41)
    - [10.4 structs包](#104-structs包)
      - [概念定义](#概念定义-42)
      - [HostLayout类型](#hostlayout类型)
      - [完整示例](#完整示例-42)
      - [使用场景](#使用场景)
    - [10.5 Timer和Ticker改进](#105-timer和ticker改进)
      - [概念定义](#概念定义-43)
      - [改进内容](#改进内容)
      - [完整示例](#完整示例-43)
      - [迁移注意事项](#迁移注意事项)
      - [最佳实践](#最佳实践-42)
    - [10.6 泛型类型别名](#106-泛型类型别名)
      - [概念定义](#概念定义-44)
      - [启用方式](#启用方式)
      - [完整示例](#完整示例-44)
      - [限制说明](#限制说明)
      - [最佳实践](#最佳实践-43)

---

## 1. 基础语法特性

### 1.1 变量声明

#### 概念定义

变量是程序运行期间可以修改其值的命名存储位置。Go语言采用静态类型系统，所有变量在编译期必须确定类型。

#### 语法形式

```go
// 完整声明形式
var name type = expression

// 类型推断形式
var name = expression

// 零值初始化形式
var name type

// 多变量声明
var (
    name1 type1 = expr1
    name2 type2 = expr2
)

// 短变量声明（仅函数内部可用）
name := expression
```

#### 属性关系

| 特性 | `var`声明 | `:=`短声明 |
|------|-----------|------------|
| 作用域 | 包级/函数级 | 仅函数级 |
| 类型推断 | 支持 | 强制推断 |
| 重新赋值 | 否（需显式类型） | 部分变量可 |
| 零值初始化 | 支持 | 不支持 |

#### 形式论证

**定理 1.1**：短声明 `:=` 至少声明一个新变量。

**证明**：设 `a, b := expr1, expr2`，若 `a` 和 `b` 均已声明，则编译错误。因此至少有一个是新变量。

```go
func demonstrateShortDecl() {
    x := 10        // 声明新变量x
    x, y := 20, 30 // x重新赋值，y是新变量（合法）
    // x, y := 40, 50 // 编译错误：无新变量
    _ = y
}
```

#### 完整示例

```go
package main

import "fmt"

// 包级变量声明
var globalVar int = 100
var (
    packageName    string = "main"
    packageVersion int    = 1
)

// 零值初始化
var defaultInt    int     // 0
var defaultString string  // ""（空字符串）
var defaultBool   bool    // false

func main() {
    // 完整声明
    var a int = 10

    // 类型推断
    var b = 20  // 推断为int
    var c = 3.14 // 推断为float64

    // 短声明
    d := "hello"
    e, f := 1, 2

    // 短声明的重新赋值特性
    f, g := 3, 4  // f重新赋值，g是新变量

    fmt.Printf("a=%d, b=%d, c=%f, d=%s, e=%d, f=%d, g=%d\n",
        a, b, c, d, e, f, g)
    fmt.Printf("零值: int=%d, string=%q, bool=%t\n",
        defaultInt, defaultString, defaultBool)
}
```

#### 反例说明

```go
// ❌ 错误：短声明在包级
packageVar := 10  // 编译错误：unexpected :=

// ❌ 错误：类型不匹配
var x int = "hello"  // 编译错误：cannot use "hello" as int

// ❌ 错误：无新变量的短声明
func badExample() {
    a := 1
    a := 2  // 编译错误：no new variables on left side of :=
}

// ❌ 错误：变量声明后未使用
func unusedVar() {
    unused := 10  // 编译错误：unused declared and not used
}
```

#### 最佳实践

1. **包级变量**使用`var`声明，保持代码清晰
2. **函数内部**优先使用`:=`短声明，减少冗余
3. **需要显式类型**时使用`var`（如需要特定数值类型）
4. **多变量声明**使用分组形式，提高可读性

---

### 1.2 常量与iota

#### 概念定义

常量是编译期确定的不可变值。Go常量可以是基本类型（布尔、数值、字符串）或类型化常量。

#### 语法形式

```go
// 无类型常量
const name = expression

// 类型化常量
const name type = expression

// 多常量声明
const (
    name1 = expr1
    name2 = expr2
)

// iota枚举（从0开始，每行递增1）
const (
    A = iota  // 0
    B         // 1
    C         // 2
)
```

#### 属性关系

**无类型常量的精度特性**：

- 无类型整数常量可以存储任意精度（仅受内存限制）
- 无类型浮点常量使用高精度表示
- 赋值时才确定具体类型

```go
const Big = 1 << 100  // 合法，无类型常量
// var x int = Big    // 编译错误：溢出
var y float64 = Big   // 合法，float64可以表示
```

#### 形式论证

**定理 1.2**：iota在const块中的值等于其行索引。

```go
const (
    A = iota      // iota = 0
    B = iota * 10 // iota = 1
    C             // iota = 2（继承表达式）
    _             // iota = 3（跳过）
    D = iota      // iota = 4
)
```

#### 完整示例

```go
package main

import "fmt"

// 无类型常量
const Pi = 3.14159265358979323846
const Greeting = "Hello, Go!"

// 类型化常量
const MaxInt32 int32 = 2147483647

// iota枚举
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)

// iota位掩码模式
const (
    Read Permission = 1 << iota   // 1 (001)
    Write                         // 2 (010)
    Execute                       // 4 (100)
)
type Permission int

// iota跳值模式
const (
    Low = iota * 10    // 0
    Medium             // 10
    High               // 20
)

// 多iota块
const (
    _ = iota          // 跳过0
    KB = 1 << (10 * iota)  // 1024
    MB                     // 1048576
    GB                     // 1073741824
)

func main() {
    fmt.Println("星期枚举:", Sunday, Monday, Tuesday)
    fmt.Println("权限:", Read, Write, Execute)
    fmt.Println("级别:", Low, Medium, High)
    fmt.Println("存储单位:", KB, MB, GB)

    // 位运算
    p := Read | Write
    fmt.Printf("权限组合 %d: 可读=%v, 可写=%v, 可执行=%v\n",
        p, p&Read != 0, p&Write != 0, p&Execute != 0)
}
```

#### 反例说明

```go
// ❌ 错误：常量不能修改
const Max = 100
Max = 200  // 编译错误：cannot assign to Max

// ❌ 错误：常量必须是编译期可确定
func getValue() int { return 10 }
const X = getValue()  // 编译错误：not constant

// ❌ 错误：iota在const块外
x := iota  // 编译错误：undefined: iota

// ❌ 错误：常量不能是切片、map、函数等
const s = []int{1, 2, 3}  // 编译错误
```

#### 最佳实践

1. **枚举值**使用iota，从0开始连续递增
2. **位标志**使用`1 << iota`模式
3. **跳过值**使用`_`占位
4. **需要特定类型**时显式指定类型
5. **大数值**使用无类型常量避免溢出

---

### 1.3 基本数据类型及零值

#### 概念定义

Go的基本数据类型分为四大类：布尔型、数值型、字符串型、派生类型。

#### 类型体系

```
基本类型
├── 布尔型: bool
├── 数值型
│   ├── 整数
│   │   ├── 有符号: int8, int16, int32, int64, int
│   │   └── 无符号: uint8, uint16, uint32, uint64, uint, uintptr
│   ├── 浮点: float32, float64
│   └── 复数: complex64, complex128
├── 字符串: string
└── 派生类型
    ├── 指针: *T
    ├── 数组: [n]T
    ├── 切片: []T
    ├── 映射: map[K]V
    ├── 通道: chan T
    ├── 函数: func
    ├── 结构体: struct
    └── 接口: interface
```

#### 零值表

| 类型 | 零值 | 说明 |
|------|------|------|
| bool | false | 布尔假 |
| 数值 | 0 | 所有数值类型 |
| string | "" | 空字符串 |
| 指针 | nil | 空指针 |
| 切片 | nil | 空切片 |
| map | nil | 空映射 |
| 通道 | nil | 空通道 |
| 函数 | nil | 空函数 |
| 接口 | nil | 空接口 |

#### 完整示例

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    // 布尔型
    var b bool
    fmt.Printf("bool: %v (size: %d bytes)\n", b, unsafe.Sizeof(b))

    // 整数型
    var i8 int8
    var i16 int16
    var i32 int32
    var i64 int64
    var i int  // 平台相关：32位系统4字节，64位系统8字节
    fmt.Printf("int8: %v, int16: %v, int32: %v, int64: %v, int: %v\n",
        i8, i16, i32, i64, i)
    fmt.Printf("int size: %d bytes\n", unsafe.Sizeof(i))

    // 无符号整数
    var ui8 uint8
    var ui16 uint16
    var ui32 uint32
    var ui64 uint64
    var ui uint
    fmt.Printf("uint8: %v, uint16: %v, uint32: %v, uint64: %v, uint: %v\n",
        ui8, ui16, ui32, ui64, ui)

    // 浮点型
    var f32 float32
    var f64 float64
    fmt.Printf("float32: %v, float64: %v\n", f32, f64)
    fmt.Printf("float32 size: %d, float64 size: %d\n",
        unsafe.Sizeof(f32), unsafe.Sizeof(f64))

    // 复数型
    var c64 complex64
    var c128 complex128
    fmt.Printf("complex64: %v, complex128: %v\n", c64, c128)

    // 字符串
    var s string
    fmt.Printf("string: %q (len: %d)\n", s, len(s))

    // 派生类型零值
    var ptr *int
    var slice []int
    var m map[string]int
    var ch chan int
    var fn func()
    var iface interface{}

    fmt.Printf("pointer: %v, is nil: %v\n", ptr, ptr == nil)
    fmt.Printf("slice: %v, is nil: %v\n", slice, slice == nil)
    fmt.Printf("map: %v, is nil: %v\n", m, m == nil)
    fmt.Printf("channel: %v, is nil: %v\n", ch, ch == nil)
    fmt.Printf("function: %v, is nil: %v\n", fn, fn == nil)
    fmt.Printf("interface: %v, is nil: %v\n", iface, iface == nil)
}
```

#### 反例说明

```go
// ❌ 错误：不同类型不能直接运算
var a int32 = 10
var b int64 = 20
c := a + b  // 编译错误：mismatched types

// ❌ 错误：整数除法截断
result := 5 / 2  // 结果是2，不是2.5

// ❌ 错误：未初始化的map不能写入
var m map[string]int
m["key"] = 1  // panic: assignment to entry in nil map

// ❌ 错误：nil切片可以读取但不能写入某些操作
var s []int
_ = s[0]  // panic: index out of range
```

#### 最佳实践

1. **整数默认**使用`int`，除非有特定需求
2. **浮点默认**使用`float64`，精度更高
3. **明确大小需求**时使用定长类型（如int32）
4. **类型转换**必须显式，避免隐式转换陷阱
5. **检查nil**后再操作指针、map、切片

---

### 1.4 类型推断机制

#### 概念定义

类型推断是编译器根据上下文自动确定变量类型的机制。Go的类型推断遵循明确的规则。

#### 推断规则

```go
// 规则1: 整数字面量默认为int
x := 10        // int

// 规则2: 浮点字面量默认为float64
y := 3.14      // float64

// 规则3: 复数字面量默认为complex128
z := 1 + 2i    // complex128

// 规则4: 字符串字面量为string
s := "hello"   // string

// 规则5: 字符字面量为rune（int32别名）
r := 'a'       // rune

// 规则6: 根据右值类型推断
var a int32 = 10
b := a         // int32（继承类型）
```

#### 完整示例

```go
package main

import "fmt"

func main() {
    // 字面量推断
    intVal := 42           // int
    floatVal := 3.14159    // float64
    complexVal := 1 + 2i   // complex128
    stringVal := "Go"      // string
    runeVal := '中'         // rune (int32)

    fmt.Printf("int: %T = %v\n", intVal, intVal)
    fmt.Printf("float64: %T = %v\n", floatVal, floatVal)
    fmt.Printf("complex128: %T = %v\n", complexVal, complexVal)
    fmt.Printf("string: %T = %v\n", stringVal, stringVal)
    fmt.Printf("rune: %T = %v\n", runeVal, runeVal)

    // 表达式推断
    a := 10
    b := 20.5
    c := a + int(b)  // int（显式转换）
    fmt.Printf("c: %T = %v\n", c, c)

    // 函数返回值推断
    result := add(1, 2)  // 推断为int
    fmt.Printf("result: %T = %v\n", result, result)

    // 复合类型推断
    slice := []int{1, 2, 3}           // []int
    m := map[string]int{"a": 1}       // map[string]int
    ch := make(chan string)           // chan string

    fmt.Printf("slice: %T\n", slice)
    fmt.Printf("map: %T\n", m)
    fmt.Printf("channel: %T\n", ch)
}

func add(a, b int) int {
    return a + b
}
```

#### 反例说明

```go
// ❌ 错误：类型推断可能导致意外
func problem() {
    x := 1 << 8   // 推断为int = 256
    // 如果期望uint8，会溢出
    var y uint8 = uint8(x)  // 需要显式转换，但可能溢出
    _ = y
}

// ❌ 错误：无类型常量精度问题
const Big = 1 << 100
// var i int = Big  // 编译错误：溢出

// ❌ 错误：类型推断与预期不符
func unexpected() {
    a := 10    // int
    b := 10.0  // float64（不是int!）
    // c := a + b  // 编译错误：类型不匹配
    _ = a
    _ = b
}
```

#### 最佳实践

1. **需要特定类型**时显式声明，不要依赖推断
2. **数值运算**注意类型一致性
3. **大数值**使用无类型常量，赋值时再转换
4. **接口实现**检查时，确保类型正确

---

### 1.5 类型别名与类型定义

#### 概念定义

- **类型定义**：创建全新的类型，与原类型完全不同
- **类型别名**：为现有类型创建另一个名称，完全等价

#### 语法形式

```go
// 类型定义（创建新类型）
type NewType ExistingType

// 类型别名（创建等价名称）
type NewType = ExistingType
```

#### 属性对比

| 特性 | 类型定义 | 类型别名 |
|------|----------|----------|
| 类型标识 | 全新类型 | 与原类型相同 |
| 方法集 | 独立 | 共享 |
| 可赋值性 | 需显式转换 | 直接赋值 |
| 用途 | 领域建模 | 兼容性/重构 |

#### 完整示例

```go
package main

import "fmt"

// ========== 类型定义 ==========
// 创建新类型，与原类型完全不同
type UserID int
type ProductID int
type Email string

// 为新类型定义方法
func (id UserID) String() string {
    return fmt.Sprintf("UserID(%d)", id)
}

func (id ProductID) String() string {
    return fmt.Sprintf("ProductID(%d)", id)
}

func (e Email) Domain() string {
    for i := 0; i < len(e); i++ {
        if e[i] == '@' {
            return string(e[i+1:])
        }
    }
    return ""
}

// ========== 类型别名 ==========
// 完全等价于原类型
type Integer = int
type StringMap = map[string]string

// byte和rune是内置类型别名
type byte = uint8
type rune = int32

func main() {
    // 类型定义的使用
    var uid UserID = 100
    var pid ProductID = 100
    var email Email = "user@example.com"

    fmt.Println(uid.String())      // UserID(100)
    fmt.Println(pid.String())      // ProductID(100)
    fmt.Println(email.Domain())    // example.com

    // 类型安全：不同类型不能直接赋值
    // uid = pid  // 编译错误：cannot use pid as UserID
    uid = UserID(pid)  // 需要显式转换

    // 类型别名的使用
    var i Integer = 10
    var j int = i  // 直接赋值，完全等价
    fmt.Printf("Integer=%d, int=%d\n", i, j)

    var sm StringMap = map[string]string{
        "key": "value",
    }
    var m map[string]string = sm  // 直接赋值
    fmt.Println(m)

    // 验证byte和rune
    var b byte = 255
    var u8 uint8 = b  // 直接赋值
    fmt.Printf("byte=%v, uint8=%v\n", b, u8)
}
```

#### 反例说明

```go
// ❌ 错误：混淆类型定义和别名
type MyInt int        // 新类型
type YourInt = int    // 别名

func confusion() {
    var a MyInt = 10
    var b int = a     // 编译错误：cannot use a as int
    var c YourInt = 10
    var d int = c     // 合法
    _ = b
    _ = d
}

// ❌ 错误：为类型别名定义方法（如果原类型已定义）
type MyString = string
// func (s MyString) Method() {}  // 编译错误：cannot define method on alias

// ❌ 错误：循环类型定义
type T1 T2  // 编译错误
type T2 T1
```

#### 最佳实践

1. **领域建模**使用类型定义，增强类型安全
2. **API兼容性**使用类型别名，平滑迁移
3. **单位类型**（如UserID、Money）使用类型定义
4. **重构代码**使用类型别名过渡
5. **不要滥用**类型别名，会降低代码清晰度

---

## 2. 复合类型

### 2.1 数组与切片（Slice）

#### 概念定义

- **数组**：固定长度、同类型元素的序列，值类型
- **切片**：动态长度、同类型元素的序列，引用类型

#### 语法形式

```go
// 数组声明
var arr [5]int                    // 长度为5的int数组
arr := [5]int{1, 2, 3, 4, 5}     // 初始化
arr := [...]int{1, 2, 3}         // 长度推断
arr := [5]int{1: 10, 3: 30}      // 索引初始化

// 切片声明
var s []int                       // nil切片
s := []int{1, 2, 3}              // 字面量初始化
s := make([]int, 5)              // 长度5，容量5
s := make([]int, 5, 10)          // 长度5，容量10
s := arr[1:4]                    // 从数组切片

// 切片操作
s = append(s, 1, 2, 3)           // 追加元素
copy(dst, src)                   // 复制切片
```

#### 内部结构

```
切片头结构（runtime.slice）
┌─────────────┬─────────────┬─────────────┐
│   指针 ptr   │   长度 len   │   容量 cap   │
│  8 bytes    │   8 bytes   │   8 bytes   │
└─────────────┴─────────────┴─────────────┘
```

#### 数组vs切片对比

| 特性 | 数组 | 切片 |
|------|------|------|
| 长度 | 固定 | 动态 |
| 类型 | 值类型 | 引用类型 |
| 传递 | 复制整个数组 | 复制切片头（24字节） |
| 比较 | 可比较（==） | 不可比较 |
| 灵活性 | 低 | 高 |

#### 完整示例

```go
package main

import "fmt"

func main() {
    // ========== 数组 ==========
    // 基本声明
    var arr1 [5]int
    fmt.Println("零值数组:", arr1)  // [0 0 0 0 0]

    // 初始化
    arr2 := [5]int{1, 2, 3, 4, 5}
    fmt.Println("初始化数组:", arr2)

    // 长度推断
    arr3 := [...]int{10, 20, 30}
    fmt.Printf("推断长度: %d, 值: %v\n", len(arr3), arr3)

    // 索引初始化
    arr4 := [5]int{1: 100, 3: 300}
    fmt.Println("索引初始化:", arr4)  // [0 100 0 300 0]

    // 多维数组
    matrix := [2][3]int{
        {1, 2, 3},
        {4, 5, 6},
    }
    fmt.Println("二维数组:", matrix)

    // ========== 切片 ==========
    // nil切片
    var nilSlice []int
    fmt.Printf("nil切片: %v, len=%d, cap=%d, is nil=%v\n",
        nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)

    // 字面量初始化
    slice1 := []int{1, 2, 3}
    fmt.Printf("字面量切片: %v, len=%d, cap=%d\n",
        slice1, len(slice1), cap(slice1))

    // make创建
    slice2 := make([]int, 5)
    fmt.Printf("make切片: %v, len=%d, cap=%d\n",
        slice2, len(slice2), cap(slice2))

    slice3 := make([]int, 3, 10)
    fmt.Printf("make切片(指定容量): %v, len=%d, cap=%d\n",
        slice3, len(slice3), cap(slice3))

    // 从数组切片
    arr := [5]int{10, 20, 30, 40, 50}
    s1 := arr[1:3]   // [20, 30]
    s2 := arr[:3]    // [10, 20, 30]
    s3 := arr[2:]    // [30, 40, 50]
    s4 := arr[:]     // [10, 20, 30, 40, 50]
    fmt.Printf("切片操作: s1=%v, s2=%v, s3=%v, s4=%v\n", s1, s2, s3, s4)

    // append操作
    s := []int{1, 2}
    s = append(s, 3)       // 追加单个
    s = append(s, 4, 5)    // 追加多个
    s = append(s, []int{6, 7}...)  // 追加切片
    fmt.Println("append后:", s)

    // copy操作
    src := []int{1, 2, 3, 4, 5}
    dst := make([]int, 3)
    n := copy(dst, src)
    fmt.Printf("copy: 复制了%d个元素, dst=%v\n", n, dst)

    // 切片共享底层数组
    original := []int{1, 2, 3, 4, 5}
    shared := original[1:3]
    shared[0] = 100
    fmt.Printf("共享底层数组: original=%v, shared=%v\n", original, shared)

    // 扩容触发新数组分配
    small := make([]int, 2, 2)
    small[0], small[1] = 1, 2
    originalPtr := &small[0]
    small = append(small, 3)  // 触发扩容
    newPtr := &small[0]
    fmt.Printf("扩容前后指针相同? %v\n", originalPtr == newPtr)
}
```

#### 反例说明

```go
// ❌ 错误：数组长度是类型的一部分
var a [3]int
var b [5]int
// a = b  // 编译错误：不同长度数组是不同类型

// ❌ 错误：切片不能直接比较
s1 := []int{1, 2}
s2 := []int{1, 2}
// fmt.Println(s1 == s2)  // 编译错误：invalid operation

// ❌ 错误：越界访问
arr := [3]int{1, 2, 3}
// _ = arr[5]  // panic: index out of range

// ❌ 错误：向nil切片追加（可以），但访问元素不行
var s []int
s = append(s, 1)  // 合法，分配底层数组
// _ = s[0]  // 如果s还是nil，会panic

// ❌ 错误：切片共享导致的意外修改
func modifySlice(s []int) {
    s[0] = 100  // 会修改原切片！
}

// ❌ 错误：append后继续使用旧切片
s := []int{1, 2, 3}
shared := s
s = append(s, 4)  // 可能触发扩容
// shared和s可能指向不同数组
```

#### 最佳实践

1. **优先使用切片**，除非需要固定长度
2. **预估容量**使用make，减少扩容开销
3. **注意共享底层数组**，必要时使用copy
4. **append后重新赋值**：`s = append(s, x)`
5. **大切片传递**使用切片而非数组，避免复制

---

### 2.2 映射（Map）

#### 概念定义

映射是无序的键值对集合，基于哈希表实现。键必须是可比较类型，值可以是任意类型。

#### 语法形式

```go
// 声明
var m map[keyType]valueType

// 初始化
m := make(map[keyType]valueType)
m := make(map[keyType]valueType, capacity)
m := map[keyType]valueType{key1: val1, key2: val2}

// 操作
m[key] = value        // 设置/更新
val := m[key]         // 获取（key不存在返回零值）
val, ok := m[key]     // 安全获取
delete(m, key)        // 删除
len(m)                // 元素个数
```

#### 内部结构

```
Map内部使用哈希表实现
┌─────────────────────────────────────────┐
│  桶数组（bucket array）                  │
│  ┌─────┬─────┬─────┬─────┬─────┐       │
│  │桶0  │桶1  │桶2  │桶3  │ ... │       │
│  └─────┴─────┴─────┴─────┴─────┘       │
│     │                                   │
│     ▼                                   │
│  ┌─────────────────────────┐           │
│  │ 键1 │ 值1 │ 键2 │ 值2 │...│          │
│  └─────────────────────────┘           │
└─────────────────────────────────────────┘
```

#### 可比较类型（可作为map键）

- 布尔、数字、字符串、指针、通道、接口
- 只包含上述类型的数组、结构体
- 切片、map、函数**不可比较**，不能作为键

#### 完整示例

```go
package main

import "fmt"

func main() {
    // nil map
    var nilMap map[string]int
    fmt.Printf("nil map: %v, is nil=%v\n", nilMap, nilMap == nil)

    // make创建
    m1 := make(map[string]int)
    m1["a"] = 1
    m1["b"] = 2
    fmt.Println("make创建的map:", m1)

    // 字面量初始化
    m2 := map[string]int{
        "one": 1,
        "two": 2,
        "three": 3,
    }
    fmt.Println("字面量map:", m2)

    // 获取值
    val := m2["two"]
    fmt.Printf("获取值: m2[\"two\"]=%d\n", val)

    // 安全获取（检查key是否存在）
    if val, ok := m2["four"]; ok {
        fmt.Printf("key存在，值=%d\n", val)
    } else {
        fmt.Println("key不存在")
    }

    // 删除
    delete(m2, "two")
    fmt.Println("删除后:", m2)

    // 遍历（无序）
    scores := map[string]int{
        "Alice": 90,
        "Bob":   85,
        "Carol": 95,
    }
    fmt.Println("遍历map:")
    for name, score := range scores {
        fmt.Printf("  %s: %d\n", name, score)
    }

    // 只遍历key
    fmt.Println("只遍历keys:")
    for name := range scores {
        fmt.Printf("  %s\n", name)
    }

    // 复杂value类型
    users := map[int]map[string]string{
        1: {"name": "Alice", "email": "alice@example.com"},
        2: {"name": "Bob", "email": "bob@example.com"},
    }
    fmt.Println("嵌套map:", users)

    // 结构体作为key
    type Point struct{ X, Y int }
    points := map[Point]string{
        {0, 0}: "origin",
        {1, 1}: "diagonal",
    }
    fmt.Println("结构体key:", points)

    // 检查key是否存在（删除前）
    if _, ok := scores["Alice"]; ok {
        delete(scores, "Alice")
        fmt.Println("删除Alice后:", scores)
    }
}
```

#### 反例说明

```go
// ❌ 错误：向nil map写入
var m map[string]int
m["key"] = 1  // panic: assignment to entry in nil map

// ❌ 错误：使用不可比较类型作为key
// badMap := map[[]int]string{}  // 编译错误：invalid map key type

// ❌ 错误：并发读写map（需要同步）
var m = make(map[string]int)
go func() { m["a"] = 1 }()  // 写
go func() { _ = m["a"] }()  // 读
// 可能导致：fatal error: concurrent map read and write

// ❌ 错误：依赖map遍历顺序
m := map[string]int{"a": 1, "b": 2}
for k, v := range m {
    // 顺序不确定！
}

// ❌ 错误：获取不存在的key不报错
m := map[string]int{}
val := m["notexist"]  // val = 0，不会报错！
// 必须用ok模式检查
```

#### 最佳实践

1. **初始化检查**：使用前确保map已make
2. **安全获取**：使用`val, ok := m[key]`模式
3. **并发访问**：使用`sync.RWMutex`或`sync.Map`
4. **预分配容量**：`make(map[K]V, hint)`减少rehash
5. **不依赖顺序**：map遍历顺序随机

---

### 2.3 结构体（Struct）

#### 概念定义

结构体是字段的集合，用于组合不同类型的数据。Go的结构体是值类型。

#### 语法形式

```go
// 定义
type StructName struct {
    Field1 Type1
    Field2 Type2
    // ...
}

// 匿名结构体
s := struct {
    Name string
    Age  int
}{"Alice", 30}

// 初始化
s := StructName{Field1: val1, Field2: val2}  // 字段名初始化
s := StructName{val1, val2}                  // 位置初始化
s := StructName{}                            // 零值
```

#### 字段标签（Tag）

```go
type User struct {
    Name  string `json:"name" db:"user_name"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"age,string"`
}
```

#### 完整示例

```go
package main

import (
    "encoding/json"
    "fmt"
)

// 基本结构体
type Person struct {
    Name string
    Age  int
}

// 带标签的结构体
type Employee struct {
    ID        int     `json:"id"`
    Name      string  `json:"name"`
    Salary    float64 `json:"salary"`
    IsActive  bool    `json:"is_active"`
    ManagerID *int    `json:"manager_id,omitempty"`
}

// 嵌套结构体
type Address struct {
    Street  string
    City    string
    ZipCode string
}

type Contact struct {
    Person
    Address
    Email string
    Phone string
}

// 指针字段结构体
type Config struct {
    Name    *string
    Timeout *int
}

func main() {
    // 基本结构体初始化
    p1 := Person{Name: "Alice", Age: 30}
    p2 := Person{"Bob", 25}  // 位置初始化
    p3 := Person{}           // 零值
    fmt.Printf("p1=%+v\n", p1)
    fmt.Printf("p2=%+v\n", p2)
    fmt.Printf("p3=%+v\n", p3)

    // 字段访问
    fmt.Printf("%s is %d years old\n", p1.Name, p1.Age)
    p1.Age = 31
    fmt.Printf("Updated age: %d\n", p1.Age)

    // 结构体比较（所有字段可比较）
    p4 := Person{Name: "Alice", Age: 30}
    fmt.Printf("p1 == p4: %v\n", p1 == p4)  // false（p1.Age已修改）

    // 带标签的结构体
    emp := Employee{
        ID:       1,
        Name:     "John",
        Salary:   50000,
        IsActive: true,
    }

    // JSON序列化
    jsonData, _ := json.Marshal(emp)
    fmt.Printf("JSON: %s\n", string(jsonData))

    // JSON反序列化
    jsonStr := `{"id":2,"name":"Jane","salary":60000,"is_active":true}`
    var emp2 Employee
    json.Unmarshal([]byte(jsonStr), &emp2)
    fmt.Printf("Parsed: %+v\n", emp2)

    // 嵌套结构体（嵌入）
    contact := Contact{
        Person:  Person{Name: "Charlie", Age: 35},
        Address: Address{Street: "123 Main St", City: "NYC", ZipCode: "10001"},
        Email:   "charlie@example.com",
        Phone:   "555-1234",
    }

    // 嵌入字段的提升访问
    fmt.Printf("Name: %s (embedded)\n", contact.Name)  // 等价于 contact.Person.Name
    fmt.Printf("City: %s (embedded)\n", contact.City)  // 等价于 contact.Address.City

    // 匿名结构体
    point := struct {
        X, Y int
    }{10, 20}
    fmt.Printf("匿名结构体: %+v\n", point)

    // 指针字段
    name := "test"
    timeout := 30
    cfg := Config{
        Name:    &name,
        Timeout: &timeout,
    }
    fmt.Printf("Config: Name=%s, Timeout=%d\n", *cfg.Name, *cfg.Timeout)

    // 检查可选字段
    cfg2 := Config{}
    if cfg2.Name == nil {
        fmt.Println("Name not set")
    }
}
```

#### 反例说明

```go
// ❌ 错误：循环嵌套结构体
type Node struct {
    Value int
    Next  *Node  // 指针可以，值类型不行
}
// type BadNode struct {
//     Value int
//     Next  BadNode  // 编译错误：invalid recursive type
// }

// ❌ 错误：不可比较字段导致结构体不可比较
type BadStruct struct {
    Data []int  // 切片不可比较
}
// s1 := BadStruct{}
// s2 := BadStruct{}
// _ = s1 == s2  // 编译错误

// ❌ 错误：嵌入字段冲突
type A struct { Name string }
type B struct { Name string }
type C struct {
    A
    B
}
// var c C
// fmt.Println(c.Name)  // 编译错误：ambiguous selector

// ❌ 错误：未导出的字段无法跨包访问
// 在另一个包中：
// p := pkg.Person{}
// p.name = "test"  // 编译错误：name未导出
```

#### 最佳实践

1. **使用字段名初始化**，提高可读性和可维护性
2. **JSON标签**使用camelCase，配合omitempty
3. **嵌入**用于组合和代码复用，注意命名冲突
4. **指针字段**用于可选字段，可区分为设置和零值
5. **方法接收者**：修改用指针`*T`，只读用值`T`

---

### 2.4 嵌套类型与匿名类型

#### 概念定义

- **嵌套类型**：在函数或方法内部定义的类型
- **匿名类型**：没有显式名称的类型（匿名结构体、匿名接口）
- **嵌入类型**：将类型作为字段嵌入到结构体中

#### 语法形式

```go
// 嵌入类型（匿名字段）
type Outer struct {
    Inner      // 嵌入，字段名是Inner
    *Pointer   // 嵌入指针类型
}

// 匿名结构体
var s = struct{ Name string }{Name: "test"}

// 匿名接口
var i interface{ Method() } = obj
```

#### 完整示例

```go
package main

import "fmt"

// ========== 嵌入类型 ==========
type Inner struct {
    Value int
}

func (i Inner) Method() {
    fmt.Println("Inner.Method, Value:", i.Value)
}

type Outer struct {
    Inner       // 嵌入Inner
    Name string
}

// 嵌入指针
type Container struct {
    *Inner
    Data string
}

// ========== 嵌套类型 ==========
func nestedTypes() {
    // 函数内定义类型
    type Local struct {
        X int
    }
    l := Local{X: 10}
    fmt.Println("嵌套类型:", l)
}

// ========== 匿名类型 ==========
func anonymousTypes() {
    // 匿名结构体
    person := struct {
        Name    string
        Age     int
        Address struct {
            City string
        }
    }{
        Name: "Alice",
        Age:  30,
        Address: struct {
            City string
        }{City: "NYC"},
    }
    fmt.Printf("匿名结构体: %+v\n", person)

    // 匿名接口
    var printer interface{ Print() } = &myPrinter{}
    printer.Print()
}

type myPrinter struct{}

func (m *myPrinter) Print() {
    fmt.Println("Printing...")
}

// ========== 嵌入接口 ==========
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}

func main() {
    // 嵌入类型的方法提升
    outer := Outer{
        Inner: Inner{Value: 100},
        Name:  "outer",
    }
    outer.Method()  // 调用嵌入类型的方法
    fmt.Println("访问嵌入字段:", outer.Value)  // 等价于 outer.Inner.Value

    // 嵌入指针
    container := Container{
        Inner: &Inner{Value: 200},
        Data:  "container",
    }
    fmt.Println("嵌入指针:", container.Value)

    // 嵌套类型
    nestedTypes()

    // 匿名类型
    anonymousTypes()

    // 接口组合
    fmt.Println("ReadWriter接口组合了Reader和Writer")
}
```

#### 反例说明

```go
// ❌ 错误：嵌入类型的方法被覆盖
type Base struct{}
func (b Base) Method() { fmt.Println("Base") }

type Derived struct {
    Base
}
func (d Derived) Method() { fmt.Println("Derived") }
// d.Method() 调用Derived的方法，Base.Method被覆盖

// ❌ 错误：嵌入指针类型nil检查
var c Container  // Inner为nil
// c.Method()  // panic: nil pointer dereference

// ❌ 错误：多层嵌入的命名冲突
type A struct{ Name string }
type B struct{ Name string }
type C struct {
    A
    B
}
// var c C
// fmt.Println(c.Name)  // 编译错误：ambiguous
```

#### 最佳实践

1. **嵌入**用于"has-a"关系的代码复用
2. **嵌入接口**用于接口组合
3. **嵌入指针**时注意nil检查
4. **避免多层嵌入**导致命名冲突
5. **匿名类型**仅用于临时/局部场景

---

## 3. 控制结构

### 3.1 if/else 及初始化语句

#### 概念定义

if语句用于条件执行，支持在条件前执行初始化语句，初始化语句中声明的变量作用域仅限于if-else块。

#### 语法形式

```go
// 基本形式
if condition {
    // ...
}

// if-else
if condition {
    // ...
} else {
    // ...
}

// if-else if-else
if condition1 {
    // ...
} else if condition2 {
    // ...
} else {
    // ...
}

// 带初始化语句
if initialization; condition {
    // initialization中声明的变量可在此使用
}
```

#### 完整示例

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    // 基本if
    x := 10
    if x > 5 {
        fmt.Println("x > 5")
    }

    // if-else
    score := 75
    if score >= 60 {
        fmt.Println("及格")
    } else {
        fmt.Println("不及格")
    }

    // if-else if-else
    grade := "B"
    if grade == "A" {
        fmt.Println("优秀")
    } else if grade == "B" {
        fmt.Println("良好")
    } else if grade == "C" {
        fmt.Println("及格")
    } else {
        fmt.Println("不及格")
    }

    // 带初始化语句（最常见用法：错误检查）
    if file, err := os.Open("test.txt"); err != nil {
        fmt.Println("打开文件失败:", err)
    } else {
        defer file.Close()
        fmt.Println("文件打开成功")
    }
    // file和err在此处不可见

    // 数值解析示例
    if n, err := strconv.Atoi("42"); err != nil {
        fmt.Println("解析失败:", err)
    } else {
        fmt.Printf("解析成功: %d\n", n)
    }

    // 复杂条件
    age := 25
    hasID := true
    if age >= 18 && hasID {
        fmt.Println("可以进入")
    }

    // 嵌套if（尽量避免）
    value := 100
    if value > 0 {
        if value < 200 {
            fmt.Println("0 < value < 200")
        }
    }
    // 更好的写法
    if value > 0 && value < 200 {
        fmt.Println("0 < value < 200")
    }
}
```

#### 反例说明

```go
// ❌ 错误：else换行问题（Go要求else在同一行）
// if x > 0 {
//     // ...
// }
// else {  // 编译错误：unexpected else
//     // ...
// }

// ❌ 错误：初始化变量作用域
def scopeIssue() {
    if x := 10; x > 5 {
        fmt.Println(x)  // 可以访问
    }
    // fmt.Println(x)  // 编译错误：undefined: x
}

// ❌ 错误：与初始化语句混淆
// if x = 10; x > 5 {  // 这是赋值不是声明
//     // ...
// }

// ❌ 错误：布尔表达式冗余
if isValid == true {  // 应该直接写 if isValid
    // ...
}
```

#### 最佳实践

1. **优先使用初始化语句**模式处理错误和资源
2. **减少嵌套**，使用提前返回
3. **布尔条件**直接写，不要与true/false比较
4. **else块**尽量简短，或考虑提前返回

---

### 3.2 switch语句

#### 概念定义

switch是多分支条件语句，Go的switch比C更灵活：自动break、支持表达式、类型switch等。

#### 语法形式

```go
// 表达式switch
switch expression {
case value1:
    // ...
case value2, value3:  // 多值匹配
    // ...
default:
    // ...
}

// 无条件switch（替代长if-else链）
switch {
case condition1:
    // ...
case condition2:
    // ...
}

// 带初始化语句
switch initialization; expression {
case ...
}

// 类型switch
switch v := i.(type) {
case Type1:
    // v是Type1类型
}
```

#### 完整示例

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 基本switch
    day := "Monday"
    switch day {
    case "Monday":
        fmt.Println("星期一")
    case "Tuesday":
        fmt.Println("星期二")
    case "Wednesday", "Thursday", "Friday":
        fmt.Println("工作日")
    default:
        fmt.Println("周末")
    }

    // 无条件switch（更清晰的多条件判断）
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

    // fallthrough（显式贯穿）
    n := 2
    switch n {
    case 1:
        fmt.Println("1")
        fallthrough
    case 2:
        fmt.Println("2")
        fallthrough
    case 3:
        fmt.Println("3")  // 2和3都会执行
    }

    // 带初始化语句
    switch hour := time.Now().Hour(); {
    case hour < 12:
        fmt.Println("上午")
    case hour < 18:
        fmt.Println("下午")
    default:
        fmt.Println("晚上")
    }

    // 类型switch
    var i interface{} = "hello"
    switch v := i.(type) {
    case int:
        fmt.Printf("整数: %d\n", v)
    case string:
        fmt.Printf("字符串: %s\n", v)
    case bool:
        fmt.Printf("布尔: %v\n", v)
    case nil:
        fmt.Println("nil")
    default:
        fmt.Printf("未知类型: %T\n", v)
    }

    // 类型switch处理多种类型
    checkType := func(x interface{}) {
        switch v := x.(type) {
        case int, int8, int16, int32, int64:
            fmt.Printf("有符号整数: %v\n", v)
        case uint, uint8, uint16, uint32, uint64:
            fmt.Printf("无符号整数: %v\n", v)
        case float32, float64:
            fmt.Printf("浮点数: %v\n", v)
        default:
            fmt.Printf("其他: %T = %v\n", v, v)
        }
    }
    checkType(42)
    checkType(3.14)
    checkType("test")
}
```

#### 反例说明

```go
// ❌ 错误：忘记Go自动break（与C不同）
switch x {
case 1:
    doSomething()
    // 不需要break！
case 2:
    doSomethingElse()
}

// ❌ 错误：case重复
def duplicateCase() {
    x := 1
    switch x {
    case 1:
        fmt.Println("1")
    case 1:  // 编译错误：duplicate case
        fmt.Println("one")
    }
}

// ❌ 错误：case不是常量
// switch x {
// case getValue():  // 编译错误：case expression must be constant
//     // ...
// }

// ❌ 错误：类型switch的变量作用域
var i interface{} = 10
switch v := i.(type) {
case int:
    fmt.Println(v)  // v是int
}
// fmt.Println(v)  // 编译错误：undefined: v
```

#### 最佳实践

1. **优先使用switch**替代长if-else链
2. **不需要显式break**，Go自动处理
3. **fallthrough**谨慎使用，容易出错
4. **类型switch**用于处理多种类型
5. **default**处理未预料的情况

---

### 3.3 for循环（三种形式）

#### 概念定义

Go只有一种循环语句`for`，但有三种形式：完整C风格、条件-only、无限循环。

#### 语法形式

```go
// 完整形式（C风格）
for initialization; condition; post {
    // ...
}

// 条件-only（类似while）
for condition {
    // ...
}

// 无限循环
for {
    // ...
}

// 带range（见3.4节）
for key, value := range collection {
    // ...
}
```

#### 完整示例

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // ========== 完整形式 ==========
    // 基本计数循环
    for i := 0; i < 5; i++ {
        fmt.Printf("i = %d\n", i)
    }

    // 多变量
    for i, j := 0, 10; i < j; i, j = i+1, j-1 {
        fmt.Printf("i=%d, j=%d\n", i, j)
    }

    // 省略初始化
    i := 0
    for ; i < 5; i++ {
        fmt.Printf("i = %d\n", i)
    }

    // 省略post
    for j := 0; j < 5; {
        fmt.Printf("j = %d\n", j)
        j++
    }

    // ========== 条件-only形式 ==========
    // 类似while
    count := 0
    for count < 5 {
        fmt.Printf("count = %d\n", count)
        count++
    }

    // 读取通道
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)

    for v := range ch {  // 等价于 for v, ok := <-ch; ok; v, ok = <-ch
        fmt.Printf("从通道接收: %d\n", v)
    }

    // ========== 无限循环 ==========
    // 配合break使用
    k := 0
    for {
        if k >= 5 {
            break
        }
        fmt.Printf("无限循环 k = %d\n", k)
        k++
    }

    // 定时器示例
    timeout := time.After(100 * time.Millisecond)
    for {
        select {
        case <-timeout:
            fmt.Println("超时!")
            return
        default:
            fmt.Println("工作中...")
            time.Sleep(30 * time.Millisecond)
        }
    }
}
```

#### 反例说明

```go
// ❌ 错误：循环变量作用域陷阱（Go 1.22之前）
// for i := 0; i < 3; i++ {
//     go func() {
//         fmt.Println(i)  // 可能都打印3！
//     }()
// }
// Go 1.22+已修复，每次迭代创建新变量

// ❌ 错误：循环变量捕获（Go 1.22之前版本）
funcs := []func(){}
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)  // 全部输出3
    })
}
// 修复：传递参数
for i := 0; i < 3; i++ {
    i := i  // 创建新变量（Go 1.22+不需要）
    funcs = append(funcs, func() {
        fmt.Println(i)
    })
}

// ❌ 错误：死循环
// for {
//     // 没有break或return
// }

// ❌ 错误：循环变量在goroutine中共享（旧版本）
for _, v := range []int{1, 2, 3} {
    go func() {
        fmt.Println(v)  // 竞态条件
    }()
}
```

#### 最佳实践

1. **优先使用for range**遍历集合
2. **Go 1.22+**循环变量每次迭代都是新的
3. **无限循环**确保有退出条件
4. **复杂循环**考虑提取为函数
5. **嵌套循环**使用标签配合break/continue

---

### 3.4 range遍历

#### 概念定义

range用于遍历数组、切片、字符串、map、通道，返回索引/键和值。

#### 语法形式

```go
// 数组/切片
for index, value := range slice {
    // ...
}

// 字符串（按rune遍历）
for index, rune := range string {
    // ...
}

// map
for key, value := range map {
    // ...
}

// 通道
for value := range channel {
    // ...
}
```

#### 遍历规则表

| 类型 | 第一个值 | 第二个值 |
|------|----------|----------|
| 数组/切片 | 索引(int) | 元素值 |
| 字符串 | 字节索引(int) | rune(int32) |
| map | 键 | 值 |
| 通道 | 元素值 | - |

#### 完整示例

```go
package main

import "fmt"

func main() {
    // ========== 数组/切片 ==========
    nums := []int{10, 20, 30, 40, 50}

    // 索引和值
    for i, v := range nums {
        fmt.Printf("索引 %d: 值 %d\n", i, v)
    }

    // 只遍历索引
    for i := range nums {
        fmt.Printf("索引: %d\n", i)
    }

    // 只遍历值（使用空白标识符）
    for _, v := range nums {
        fmt.Printf("值: %d\n", v)
    }

    // ========== 字符串 ==========
    s := "Hello, 世界"

    // 按字节遍历
    for i := 0; i < len(s); i++ {
        fmt.Printf("字节 %d: %x\n", i, s[i])
    }

    // 按rune遍历（推荐）
    for i, r := range s {
        fmt.Printf("索引 %d: rune %c (码点 %d)\n", i, r, r)
    }

    // ========== map ==========
    scores := map[string]int{
        "Alice": 90,
        "Bob":   85,
        "Carol": 95,
    }

    // 遍历map（顺序随机）
    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
    }

    // 只遍历key
    for name := range scores {
        fmt.Printf("姓名: %s\n", name)
    }

    // ========== 通道 ==========
    ch := make(chan int, 3)
    ch <- 1
    ch <- 2
    ch <- 3
    close(ch)

    for v := range ch {
        fmt.Printf("从通道接收: %d\n", v)
    }

    // ========== 修改遍历中的值 ==========
    // 值是副本，修改不影响原切片
    for _, v := range nums {
        v *= 2  // 只修改副本
    }
    fmt.Println("修改副本后:", nums)  // 原切片未变

    // 通过索引修改
    for i := range nums {
        nums[i] *= 2  // 修改原切片
    }
    fmt.Println("通过索引修改后:", nums)
}
```

#### 反例说明

```go
// ❌ 错误：遍历map时修改map（某些情况下安全，但不推荐）
m := map[string]int{"a": 1, "b": 2}
for k := range m {
    if k == "a" {
        delete(m, k)  // 可以删除当前key
        // m["c"] = 3  // 可能panic或不确定行为
    }
}

// ❌ 错误：range遍历nil
var nilSlice []int
for _, v := range nilSlice {
    fmt.Println(v)  // 不会执行，但合法
}

var nilMap map[string]int
for k, v := range nilMap {
    fmt.Println(k, v)  // 不会执行，但合法
}

// ❌ 错误：遍历通道时未关闭会导致死锁
// ch := make(chan int)
// for v := range ch {  // 永远不会结束
//     fmt.Println(v)
// }

// ❌ 错误：字符串索引误解
s := "Hello"
for i := range s {
    // i是字节索引，不是字符索引
    // 对于多字节字符，索引会跳跃
}
```

#### 最佳实践

1. **字符串遍历**使用range按rune遍历
2. **map遍历**不依赖顺序
3. **修改元素**通过索引而非值
4. **通道遍历**确保通道会被关闭
5. **只需要索引**时省略值：`for i := range s`

---

### 3.5 goto、break、continue

#### 概念定义

- **break**：跳出最内层循环或switch
- **continue**：跳过当前迭代，进入下一次
- **goto**：无条件跳转到标签（限制使用）

#### 语法形式

```go
// break
for {
    break  // 跳出循环
}

// 带标签的break
OuterLoop:
for i := 0; i < 10; i++ {
    for j := 0; j < 10; j++ {
        break OuterLoop  // 跳出外层循环
    }
}

// continue
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue  // 跳过偶数
    }
}

// 带标签的continue
Outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        continue Outer  // 跳到外层下一次迭代
    }
}

// goto
for i := 0; i < 10; i++ {
    if i == 5 {
        goto End
    }
}
End:
    fmt.Println("结束")
```

#### 完整示例

```go
package main

import "fmt"

func main() {
    // ========== break ==========
    // 基本break
    fmt.Println("基本break:")
    for i := 0; i < 10; i++ {
        if i == 5 {
            break
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()

    // 带标签的break（跳出外层循环）
    fmt.Println("\n带标签的break:")
OuterBreak:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                break OuterBreak
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    fmt.Println()

    // ========== continue ==========
    // 基本continue
    fmt.Println("\n基本continue（跳过偶数）:")
    for i := 0; i < 10; i++ {
        if i%2 == 0 {
            continue
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()

    // 带标签的continue
    fmt.Println("\n带标签的continue:")
OuterContinue:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                continue OuterContinue
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
    fmt.Println()

    // ========== goto ==========
    // 错误处理模式
    fmt.Println("\ngoto错误处理:")
    err := doSomething()
    if err != nil {
        goto Error
    }
    err = doAnotherThing()
    if err != nil {
        goto Error
    }
    fmt.Println("成功!")
    return

Error:
    fmt.Printf("错误: %v\n", err)

    // 循环中的goto
    fmt.Println("\ngoto跳出嵌套:")
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                goto Done
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
Done:
    fmt.Println("\n完成")
}

func doSomething() error {
    return nil
}

func doAnotherThing() error {
    return nil
}
```

#### 反例说明

```go
// ❌ 错误：goto跳入变量作用域
def badGoto() {
    goto Skip
    x := 10  // 编译错误：goto跳过变量声明
Skip:
    fmt.Println(x)
}

// ❌ 错误：goto跳入代码块
def badGoto2() {
    goto Inside
    {
        x := 10
    Inside:
        fmt.Println(x)  // 编译错误
    }
}

// ❌ 错误：滥用goto导致 spaghetti code
// 避免使用goto实现循环逻辑

// ❌ 错误：break在switch中只跳出switch，不是循环
for {
    switch x {
    case 1:
        break  // 只跳出switch，不是循环！
    }
}
// 需要使用带标签的break
Loop:
for {
    switch x {
    case 1:
        break Loop  // 跳出循环
    }
}
```

#### 最佳实践

1. **避免goto**，除非用于错误处理清理
2. **优先使用函数**替代复杂goto逻辑
3. **带标签的break/continue**用于嵌套循环
4. **switch中的break**注意只跳出switch
5. **defer**通常比goto更适合资源清理

---

## 4. 函数与方法

### 4.1 函数定义与多返回值

#### 概念定义

函数是可重用的代码块，Go函数支持多返回值，这是Go错误处理的基础。

#### 语法形式

```go
// 基本函数
func functionName(param1 Type1, param2 Type2) ReturnType {
    // ...
    return value
}

// 多返回值
func functionName() (Type1, Type2) {
    return value1, value2
}

// 命名返回值
func functionName() (result Type) {
    // ...
    return  // 裸return，返回result
}

// 相同类型参数简写
func functionName(a, b int) int {
    return a + b
}
```

#### 完整示例

```go
package main

import (
    "errors"
    "fmt"
)

// 基本函数
func add(a, b int) int {
    return a + b
}

// 多返回值（错误处理模式）
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("除数不能为零")
    }
    return a / b, nil
}

// 命名返回值
func rectangle(width, height float64) (area, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return  // 裸return，返回命名值
}

// 多返回值（结果+布尔标志）
func find(slice []int, target int) (int, bool) {
    for i, v := range slice {
        if v == target {
            return i, true
        }
    }
    return -1, false
}

// 返回函数
func makeMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

// 递归函数
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// 尾递归优化（Go不保证）
func factorialTail(n, acc int) int {
    if n <= 1 {
        return acc
    }
    return factorialTail(n-1, n*acc)
}

func main() {
    // 基本函数调用
    sum := add(3, 5)
    fmt.Printf("3 + 5 = %d\n", sum)

    // 多返回值处理
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("错误:", err)
    } else {
        fmt.Printf("10 / 2 = %f\n", result)
    }

    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("错误:", err)
    }

    // 命名返回值
    area, perimeter := rectangle(5, 3)
    fmt.Printf("矩形: 面积=%.2f, 周长=%.2f\n", area, perimeter)

    // 只使用部分返回值
    area2, _ := rectangle(4, 6)
    fmt.Printf("面积=%.2f\n", area2)

    // 查找函数
    nums := []int{10, 20, 30, 40, 50}
    if idx, found := find(nums, 30); found {
        fmt.Printf("找到30在索引%d\n", idx)
    }

    // 返回的函数
    double := makeMultiplier(2)
    triple := makeMultiplier(3)
    fmt.Printf("double(5)=%d, triple(5)=%d\n", double(5), triple(5))

    // 递归
    fmt.Printf("5! = %d\n", factorial(5))
    fmt.Printf("5! (尾递归) = %d\n", factorialTail(5, 1))
}
```

#### 反例说明

```go
// ❌ 错误：命名返回值与return值冲突
func badNamedReturn() (result int) {
    result = 10
    return 20  // 返回20，不是result
}

// ❌ 错误：裸return在复杂函数中降低可读性
func confusing() (a, b, c int) {
    // 很多代码...
    a = 1
    // 更多代码...
    b = 2
    // 更多代码...
    return  // 难以追踪返回值
}

// ❌ 错误：忽略重要的错误返回值
result, _ := divide(10, 0)  // 危险！
fmt.Println(result)  // 结果是0，但不知道出错了

// ❌ 错误：递归没有终止条件
func infinite(n int) int {
    return n + infinite(n+1)  // 栈溢出
}
```

#### 最佳实践

1. **错误处理**：`(result, error)`模式
2. **命名返回值**：简单函数使用，复杂函数避免
3. **裸return**：仅用于简单函数
4. **忽略返回值**：使用`_`，但不要忽略错误
5. **文档**：为导出函数编写注释

---

### 4.2 变参函数

#### 概念定义

变参函数接受可变数量的参数，使用`...Type`语法，在函数内部作为切片处理。

#### 语法形式

```go
// 变参函数定义
func functionName(params ...Type) {
    // params是[]Type类型
}

// 调用变参函数
functionName(1, 2, 3)
functionName(slice...)  // 展开切片
```

#### 完整示例

```go
package main

import "fmt"

// 基本变参函数
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// 混合参数
func printf(format string, args ...interface{}) {
    fmt.Printf(format, args...)
}

// 变参+其他参数
func join(sep string, parts ...string) string {
    result := ""
    for i, part := range parts {
        if i > 0 {
            result += sep
        }
        result += part
    }
    return result
}

// 递归处理变参
func max(first int, rest ...int) int {
    m := first
    for _, v := range rest {
        if v > m {
            m = v
        }
    }
    return m
}

func main() {
    // 基本调用
    fmt.Println("sum() =", sum())
    fmt.Println("sum(1) =", sum(1))
    fmt.Println("sum(1, 2, 3, 4, 5) =", sum(1, 2, 3, 4, 5))

    // 传递切片
    nums := []int{10, 20, 30}
    fmt.Println("sum(slice...) =", sum(nums...))

    // 混合参数
    printf("Hello, %s! You are %d years old.\n", "Alice", 30)

    // 连接字符串
    fmt.Println(join(", ", "a", "b", "c"))
    fmt.Println(join("-", "2024", "01", "15"))

    // 最大值
    fmt.Println("max(3, 1, 4, 1, 5, 9) =", max(3, 1, 4, 1, 5, 9))

    // 空变参
    fmt.Println("join(",\") =", join(", "))
}
```

#### 反例说明

```go
// ❌ 错误：变参必须在最后
// func bad(a ...int, b string)  // 编译错误

// ❌ 错误：只能有一个变参
// func bad(a ...int, b ...string)  // 编译错误

// ❌ 错误：类型不匹配
nums := []int{1, 2, 3}
// join(", ", nums...)  // 编译错误：需要string，得到int

// ❌ 错误：变参切片与原切片共享
func modify(nums ...int) {
    nums[0] = 100  // 会修改原切片！
}

// ❌ 错误：nil变参
func withNil(items ...string) {
    fmt.Println(items == nil)  // false，变参不会是nil
}
withNil()  // items是空切片[]string{}，不是nil
```

#### 最佳实践

1. **变参放在最后**，只能有一个
2. **切片展开**使用`slice...`
3. **空变参**产生空切片，不是nil
4. **注意共享**：变参切片可能与原切片共享底层数组
5. **interface{}变参**用于通用格式化函数

---

### 4.3 匿名函数与闭包

#### 概念定义

- **匿名函数**：没有名称的函数，可以赋值给变量或立即执行
- **闭包**：函数及其引用的外部变量组合，变量在函数间共享

#### 语法形式

```go
// 匿名函数赋值
fn := func(x, y int) int {
    return x + y
}

// 立即执行
result := func(x, y int) int {
    return x + y
}(3, 5)

// 闭包
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
```

#### 完整示例

```go
package main

import "fmt"

func main() {
    // ========== 匿名函数 ==========
    // 赋值给变量
    add := func(a, b int) int {
        return a + b
    }
    fmt.Println("匿名函数:", add(3, 5))

    // 立即执行（IIFE）
    result := func(x int) int {
        return x * x
    }(4)
    fmt.Println("立即执行:", result)

    // 作为参数
    apply := func(x int, fn func(int) int) int {
        return fn(x)
    }
    square := func(x int) int { return x * x }
    fmt.Println("作为参数:", apply(5, square))

    // 作为返回值
    makeMultiplier := func(factor int) func(int) int {
        return func(x int) int {
            return x * factor
        }
    }
    double := makeMultiplier(2)
    fmt.Println("double(7) =", double(7))

    // ========== 闭包 ==========
    // 计数器闭包
    counter := func() func() int {
        count := 0
        return func() int {
            count++
            return count
        }
    }()
    fmt.Println("counter:", counter())  // 1
    fmt.Println("counter:", counter())  // 2
    fmt.Println("counter:", counter())  // 3

    // 多个独立闭包
    counter1 := makeCounter()
    counter2 := makeCounter()
    fmt.Println("counter1:", counter1())  // 1
    fmt.Println("counter1:", counter1())  // 2
    fmt.Println("counter2:", counter2())  // 1（独立状态）

    // 闭包修改外部变量
    sum := 0
    nums := []int{1, 2, 3, 4, 5}
    forEach(nums, func(n int) {
        sum += n  // 修改外部sum
    })
    fmt.Println("sum =", sum)  // 15

    // 函数工厂
    adder := makeAdder(10)
    fmt.Println("adder(5) =", adder(5))   // 15
    fmt.Println("adder(3) =", adder(3))   // 13（10+3）

    // 带状态的处理器
    handlers := makeHandlers()
    handlers["increment"](5)
    handlers["decrement"](3)
}

// 返回闭包的函数
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// 函数工厂
func makeAdder(base int) func(int) int {
    return func(x int) int {
        return base + x
    }
}

// 接受函数作为参数
func forEach(slice []int, fn func(int)) {
    for _, v := range slice {
        fn(v)
    }
}

// 返回多个闭包
func makeHandlers() map[string]func(int) {
    count := 0
    return map[string]func(int){
        "increment": func(n int) {
            count += n
            fmt.Printf("Incremented: count = %d\n", count)
        },
        "decrement": func(n int) {
            count -= n
            fmt.Printf("Decremented: count = %d\n", count)
        },
    }
}
```

#### 反例说明

```go
// ❌ 错误：循环变量捕获（Go 1.22之前）
funcs := []func(){}
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)  // 全部输出3
    })
}

// ❌ 错误：闭包意外共享变量
var handlers []func()
for i := 0; i < 3; i++ {
    handlers = append(handlers, func() {
        fmt.Println(i)  // 问题同上
    })
}

// 修复：创建局部副本
for i := 0; i < 3; i++ {
    i := i  // 创建新变量（Go 1.22+不需要）
    handlers = append(handlers, func() {
        fmt.Println(i)
    })
}

// ❌ 错误：闭包持有大对象导致内存泄漏
func process() func() {
    bigData := make([]byte, 1<<30)  // 1GB
    return func() {
        _ = bigData[0]  // 闭包引用bigData，无法释放
    }
}

// ❌ 错误：goroutine中的闭包问题
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // 竞态条件
    }()
}
// 修复
for i := 0; i < 10; i++ {
    go func(i int) {
        fmt.Println(i)
    }(i)
}
```

#### 最佳实践

1. **Go 1.22+**循环变量每次迭代都是新的
2. **goroutine参数**显式传递避免竞态
3. **注意内存泄漏**，闭包可能持有大对象
4. **闭包状态**清晰文档化
5. **函数工厂**用于创建配置化的函数

---

### 4.4 方法接收者

#### 概念定义

方法是绑定到类型的函数，接收者指定方法操作的对象。接收者可以是值或指针。

#### 语法形式

```go
// 值接收者（复制接收者）
func (r ReceiverType) MethodName() ReturnType {
    // r是接收者的副本
}

// 指针接收者（修改原对象）
func (r *ReceiverType) MethodName() ReturnType {
    // r指向原对象
}
```

#### 值接收者 vs 指针接收者

| 特性 | 值接收者 | 指针接收者 |
|------|----------|------------|
| 修改原对象 | 否 | 是 |
| 调用方式 | 值和指针都可 | 值和指针都可 |
| 复制开销 | 复制整个对象 | 复制指针（8字节） |
| nil处理 | 不适用 | 需要检查 |
| 一致性 | 简单类型 | 复杂类型 |

#### 完整示例

```go
package main

import "fmt"

// ========== 值接收者 ==========
type Point struct {
    X, Y float64
}

// 值接收者：不修改原对象
func (p Point) Distance(q Point) float64 {
    return sqrt((p.X-q.X)*(p.X-q.X) + (p.Y-q.Y)*(p.Y-q.Y))
}

func (p Point) String() string {
    return fmt.Sprintf("(%.2f, %.2f)", p.X, p.Y)
}

// 尝试修改（不会影响原对象）
func (p Point) Move(dx, dy float64) {
    p.X += dx  // 修改副本
    p.Y += dy
}

// ========== 指针接收者 ==========
type Counter struct {
    value int
}

// 指针接收者：修改原对象
func (c *Counter) Increment() {
    c.value++
}

func (c *Counter) Value() int {
    return c.value
}

// 需要检查nil
func (c *Counter) SafeValue() int {
    if c == nil {
        return 0
    }
    return c.value
}

// ========== 混合使用 ==========
type Rectangle struct {
    Width, Height float64
}

// 值接收者（只读）
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// 指针接收者（修改）
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    // 值接收者
    p1 := Point{3, 4}
    p2 := Point{0, 0}
    fmt.Printf("p1=%s, p2=%s\n", p1, p2)
    fmt.Printf("Distance=%.2f\n", p1.Distance(p2))

    // 值接收者修改无效
    p1.Move(10, 10)
    fmt.Printf("Move后 p1=%s（未改变）\n", p1)

    // 指针接收者
    c := &Counter{}
    c.Increment()
    c.Increment()
    fmt.Printf("Counter value=%d\n", c.Value())

    // nil指针检查
    var nilCounter *Counter
    fmt.Printf("nil counter value=%d\n", nilCounter.SafeValue())

    // 混合使用
    rect := Rectangle{Width: 10, Height: 5}
    fmt.Printf("Area=%.2f\n", rect.Area())
    rect.Scale(2)
    fmt.Printf("Scale后 Area=%.2f\n", rect.Area())

    // 自动解引用
    rect2 := Rectangle{Width: 3, Height: 4}
    rect2.Scale(2)  // 等价于 (&rect2).Scale(2)
    fmt.Printf("rect2 Area=%.2f\n", rect2.Area())
}

func sqrt(x float64) float64 {
    if x < 0 {
        return 0
    }
    z := 1.0
    for i := 0; i < 10; i++ {
        z -= (z*z - x) / (2 * z)
    }
    return z
}
```

#### 反例说明

```go
// ❌ 错误：值接收者修改期望生效
type Buffer struct {
    data []byte
}
func (b Buffer) Write(p []byte) {
    b.data = append(b.data, p...)  // 修改副本，原对象不变！
}

// ❌ 错误：指针接收者nil检查
type Tree struct {
    Value int
    Left  *Tree
    Right *Tree
}
func (t *Tree) Sum() int {
    // 忘记检查nil
    return t.Value + t.Left.Sum() + t.Right.Sum()  // 可能panic
}
// 修复
func (t *Tree) SafeSum() int {
    if t == nil {
        return 0
    }
    return t.Value + t.Left.SafeSum() + t.Right.SafeSum()
}

// ❌ 错误：不一致的接收者类型
type MyType struct{}
func (m MyType) Method1() {}   // 值
func (m *MyType) Method2() {}  // 指针
// 可能导致接口实现问题

// ❌ 错误：大对象使用值接收者
type BigStruct struct {
    data [1000000]int  // 大数组
}
func (b BigStruct) Process() {  // 复制开销大
    // ...
}
```

#### 最佳实践

1. **修改对象**使用指针接收者
2. **大对象**使用指针接收者避免复制
3. **一致性**：一个类型的方法要么全用值，要么全用指针
4. **nil检查**：指针接收者检查nil
5. **不可变对象**使用值接收者

---

### 4.5 方法集

#### 概念定义

方法集是类型上定义的所有方法的集合，决定类型实现了哪些接口。

#### 方法集规则

| 类型 | 方法集 |
|------|--------|
| T | 所有值接收者方法 |
| *T | 所有值接收者方法 + 所有指针接收者方法 |

```
方法集关系：
┌─────────────────────────────────────┐
│           *T 的方法集                │
│  ┌─────────────────────────────┐    │
│  │       T 的方法集             │    │
│  │  (值接收者方法)              │    │
│  │                             │    │
│  │                             │    │
│  └─────────────────────────────┘    │
│  (额外包含指针接收者方法)              │
└─────────────────────────────────────┘
```

#### 完整示例

```go
package main

import "fmt"

// ========== 定义类型和方法 ==========
type MyInt int

// 值接收者方法
func (m MyInt) ValueMethod() string {
    return fmt.Sprintf("ValueMethod: %d", m)
}

// 指针接收者方法
func (m *MyInt) PointerMethod() {
    *m = *m * 2
}

// ========== 接口定义 ==========
type ValueInterface interface {
    ValueMethod() string
}

type PointerInterface interface {
    ValueMethod() string
    PointerMethod()
}

// ========== 结构体方法集 ==========
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    // ========== 方法集演示 ==========
    var i MyInt = 10

    // 值类型的方法集
    fmt.Println("值类型调用值方法:", i.ValueMethod())
    // fmt.Println(i.PointerMethod())  // 编译错误

    // 指针类型的方法集
    p := &i
    fmt.Println("指针调用值方法:", p.ValueMethod())  // 自动解引用
    p.PointerMethod()  // 可以调用指针方法
    fmt.Println("PointerMethod后:", i)

    // ========== 接口实现 ==========
    // MyInt实现ValueInterface（有ValueMethod）
    var vi ValueInterface = i
    fmt.Println("ValueInterface:", vi.ValueMethod())

    // *MyInt实现PointerInterface（有ValueMethod和PointerMethod）
    var pi PointerInterface = &i
    pi.PointerMethod()
    fmt.Println("PointerInterface后:", i)

    // ========== 结构体方法集 ==========
    rect := Rectangle{Width: 10, Height: 5}

    // 值类型可以调用值方法
    fmt.Printf("Area=%.2f\n", rect.Area())

    // 值类型不能调用指针方法（但Go自动取地址）
    rect.Scale(2)  // 等价于 (&rect).Scale(2)
    fmt.Printf("Scale后 Area=%.2f\n", rect.Area())

    // ========== 接口赋值规则 ==========
    // 值类型赋值给接口
    var areaer interface{ Area() float64 } = rect
    fmt.Printf("通过接口调用Area=%.2f\n", areaer.Area())

    // 指针类型赋值给接口
    var scaler interface{ Scale(float64) } = &rect
    scaler.Scale(2)
    fmt.Printf("通过接口Scale后 Area=%.2f\n", rect.Area())

    // ========== 方法表达式 ==========
    // 值方法表达式
    valueExpr := MyInt.ValueMethod
    fmt.Println("方法表达式:", valueExpr(42))

    // 指针方法表达式
    pointerExpr := (*MyInt).PointerMethod
    pointerExpr(&i)
    fmt.Println("PointerMethod表达式后:", i)
}
```

#### 反例说明

```go
// ❌ 错误：值类型不满足需要指针方法的接口
type Interface interface {
    Modify()
}
type MyType struct{}
func (m *MyType) Modify() {}

var i Interface = MyType{}  // 编译错误：MyType没有Modify方法

// ❌ 错误：接口断言时方法集不匹配
var i interface{} = MyType{}
// m := i.(Interface)  // panic

// ❌ 错误：nil指针调用方法
type Tree struct{}
func (t *Tree) Method() {}
var t *Tree = nil
t.Method()  // 可以调用，但如果方法访问字段会panic

// ❌ 错误：混淆方法集和可调用性
// 值类型可以"调用"指针方法（Go自动取地址）
// 但值类型不满足需要指针方法的接口
```

#### 最佳实践

1. **接口实现**：检查方法集是否满足接口
2. **一致性**：类型方法要么全值，要么全指针
3. **指针方法**：修改状态或避免大对象复制
4. **nil接收者**：指针方法检查nil
5. **方法表达式**：用于函数式编程场景

---

## 5. 接口系统

### 5.1 接口定义与实现

#### 概念定义

接口定义了一组方法签名，任何类型只要实现了这些方法就隐式实现了该接口。Go采用**结构化类型系统**（鸭子类型）。

#### 语法形式

```go
// 接口定义
type InterfaceName interface {
    Method1(param Type) ReturnType
    Method2(param Type) ReturnType
}

// 空接口（所有类型都实现）
type EmptyInterface interface{}
```

#### 实现规则

```
类型T实现接口I，当且仅当：
1. T拥有I中定义的所有方法
2. 方法签名完全匹配
3. 方法可以是值接收者或指针接收者
```

#### 完整示例

```go
package main

import (
    "fmt"
    "math"
)

// ========== 接口定义 ==========
// 几何形状接口
type Shape interface {
    Area() float64
    Perimeter() float64
}

// 可绘制接口
type Drawable interface {
    Draw()
}

// 组合接口
type ShapeDrawable interface {
    Shape
    Drawable
}

// ========== 类型实现 ==========
// 矩形
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r Rectangle) Draw() {
    fmt.Printf("绘制矩形: %.2f x %.2f\n", r.Width, r.Height)
}

// 圆形
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

func (c Circle) Draw() {
    fmt.Printf("绘制圆形: 半径=%.2f\n", c.Radius)
}

// ========== 使用接口 ==========
func PrintShapeInfo(s Shape) {
    fmt.Printf("面积: %.2f, 周长: %.2f\n", s.Area(), s.Perimeter())
}

func DrawAll(drawables []Drawable) {
    for _, d := range drawables {
        d.Draw()
    }
}

func main() {
    // 创建具体类型
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}

    // 赋值给接口
    var s Shape = rect
    fmt.Println("矩形:")
    PrintShapeInfo(s)

    s = circle
    fmt.Println("圆形:")
    PrintShapeInfo(s)

    // 组合接口
    var sd ShapeDrawable = rect
    fmt.Println("\n组合接口:")
    PrintShapeInfo(sd)
    sd.Draw()

    // 接口切片
    shapes := []Shape{rect, circle}
    fmt.Println("\n所有形状:")
    for _, shape := range shapes {
        PrintShapeInfo(shape)
    }

    // 类型判断
    fmt.Println("\n类型判断:")
    for _, shape := range shapes {
        switch v := shape.(type) {
        case Rectangle:
            fmt.Printf("矩形: 宽=%.2f, 高=%.2f\n", v.Width, v.Height)
        case Circle:
            fmt.Printf("圆形: 半径=%.2f\n", v.Radius)
        }
    }
}
```

#### 反例说明

```go
// ❌ 错误：方法签名不匹配
type Interface interface {
    Method(int) string
}
type MyType struct{}
func (m MyType) Method(x int64) string { return "" }  // 参数类型不同
// MyType不实现Interface

// ❌ 错误：缺少方法
type Reader interface {
    Read([]byte) (int, error)
    Close() error
}
type MyReader struct{}
func (m MyReader) Read(p []byte) (int, error) { return 0, nil }
// MyReader不实现Reader（缺少Close）

// ❌ 错误：值类型和指针类型的方法集
type Interface interface {
    Modify()
}
type MyType struct{}
func (m *MyType) Modify() {}
// var i Interface = MyType{}  // 编译错误
var i Interface = &MyType{}  // 正确

// ❌ 错误：接口循环定义
// type I interface {
//     Method() I
// }
```

#### 最佳实践

1. **小接口**：接口应该小而专注（Go标准库推荐）
2. **隐式实现**：不需要显式声明`implements`
3. **接口组合**：用小接口组合大接口
4. **接口定义在使用方**：依赖抽象而非具体
5. **文档**：说明接口的语义和契约

---

### 5.2 空接口interface{}

#### 概念定义

空接口不包含任何方法，因此**所有类型都实现了空接口**。空接口可以存储任何类型的值。

#### 内部表示

```
空接口内部结构（runtime.eface）
┌─────────────────┬─────────────────┐
│   类型指针 _type │   数据指针 data  │
│    8 bytes      │    8 bytes      │
└─────────────────┴─────────────────┘
```

#### 完整示例

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    // ========== 空接口存储任意类型 ==========
    var i interface{}

    i = 42
    fmt.Printf("int: %v, type: %T\n", i, i)

    i = "hello"
    fmt.Printf("string: %v, type: %T\n", i, i)

    i = 3.14
    fmt.Printf("float64: %v, type: %T\n", i, i)

    i = []int{1, 2, 3}
    fmt.Printf("[]int: %v, type: %T\n", i, i)

    i = map[string]int{"a": 1}
    fmt.Printf("map: %v, type: %T\n", i, i)

    i = struct{ Name string }{Name: "test"}
    fmt.Printf("struct: %v, type: %T\n", i, i)

    // ========== 空接口切片 ==========
    mixed := []interface{}{
        1,
        "two",
        3.0,
        true,
        []int{1, 2},
        map[string]string{"key": "value"},
    }

    fmt.Println("\n混合切片:")
    for _, v := range mixed {
        fmt.Printf("  %v (type: %T)\n", v, v)
    }

    // ========== 类型断言 ==========
    fmt.Println("\n类型断言:")
    var x interface{} = "hello"

    // 安全断言
    if s, ok := x.(string); ok {
        fmt.Printf("是字符串: %s\n", s)
    }

    // 不安全断言（可能panic）
    // n := x.(int)  // panic

    // ========== 类型switch ==========
    fmt.Println("\n类型switch:")
    checkType := func(v interface{}) {
        switch val := v.(type) {
        case int:
            fmt.Printf("整数: %d\n", val)
        case string:
            fmt.Printf("字符串: %s\n", val)
        case float64:
            fmt.Printf("浮点数: %f\n", val)
        case bool:
            fmt.Printf("布尔: %v\n", val)
        case nil:
            fmt.Println("nil")
        default:
            fmt.Printf("其他类型: %T\n", val)
        }
    }

    checkType(42)
    checkType("test")
    checkType(3.14)
    checkType(nil)

    // ========== reflect包 ==========
    fmt.Println("\n使用reflect:")
    inspect := func(v interface{}) {
        t := reflect.TypeOf(v)
        val := reflect.ValueOf(v)
        fmt.Printf("类型: %s, 值: %v, 种类: %s\n",
            t.Name(), val, t.Kind())
    }

    inspect(100)
    inspect("hello")
    inspect([]int{1, 2, 3})

    // ========== 空接口nil陷阱 ==========
    fmt.Println("\nnil陷阱:")
    var p *int = nil
    var iface interface{} = p
    fmt.Printf("p == nil: %v\n", p == nil)
    fmt.Printf("iface == nil: %v\n", iface == nil)  // false！

    // 检查空接口是否为nil
    if iface == nil {
        fmt.Println("iface是nil")
    } else {
        fmt.Println("iface不是nil（包含类型信息）")
        if reflect.ValueOf(iface).IsNil() {
            fmt.Println("但值是nil")
        }
    }
}
```

#### 反例说明

```go
// ❌ 错误：空接口nil检查
def nilTrap() {
    var p *int = nil
    var i interface{} = p
    if i == nil {
        // 不会执行，因为i包含类型信息
    }
    // 应该这样检查
    if p == nil {
        // ...
    }
}

// ❌ 错误：类型断言panic
var i interface{} = "hello"
n := i.(int)  // panic: interface conversion

// ❌ 错误：直接操作空接口值
var i interface{} = 10
// i++  // 编译错误：invalid operation

// ❌ 错误：空接口切片类型转换
mixed := []interface{}{1, 2, 3}
// ints := []int(mixed)  // 编译错误

// ❌ 错误：滥用空接口失去类型安全
func process(data interface{}) {
    // 运行时才能知道类型，容易出错
}
```

#### 最佳实践

1. **避免滥用空接口**，失去类型安全
2. **类型断言**使用ok模式
3. **类型switch**处理多种类型
4. **nil检查**注意空接口的特殊性
5. **泛型**（Go 1.18+）替代部分空接口场景

---

### 5.3 类型断言与类型切换

#### 概念定义

- **类型断言**：从接口中提取具体类型的值
- **类型切换**：根据接口值的实际类型执行不同分支

#### 语法形式

```go
// 类型断言（可能panic）
value := interfaceValue.(ConcreteType)

// 安全类型断言
value, ok := interfaceValue.(ConcreteType)

// 类型switch
switch v := interfaceValue.(type) {
case Type1:
    // v是Type1类型
case Type2:
    // v是Type2类型
default:
    // ...
}
```

#### 完整示例

```go
package main

import (
    "fmt"
    "io"
    "os"
)

// ========== 自定义类型 ==========
type Stringer interface {
    String() string
}

type MyString string

func (s MyString) String() string {
    return string(s)
}

// ========== 类型断言示例 ==========
func assertString(v interface{}) string {
    // 安全断言
    if s, ok := v.(string); ok {
        return s
    }
    if s, ok := v.(fmt.Stringer); ok {
        return s.String()
    }
    return fmt.Sprintf("%v", v)
}

// ========== 类型切换示例 ==========
func describe(v interface{}) string {
    switch val := v.(type) {
    case nil:
        return "nil"
    case int:
        return fmt.Sprintf("整数: %d", val)
    case int8, int16, int32, int64:
        return fmt.Sprintf("有符号整数: %v", val)
    case uint, uint8, uint16, uint32, uint64:
        return fmt.Sprintf("无符号整数: %v", val)
    case float32:
        return fmt.Sprintf("float32: %f", val)
    case float64:
        return fmt.Sprintf("float64: %f", val)
    case string:
        return fmt.Sprintf("字符串(长度%d): %s", len(val), val)
    case bool:
        return fmt.Sprintf("布尔: %v", val)
    case []interface{}:
        return fmt.Sprintf("interface切片(长度%d)", len(val))
    case map[string]interface{}:
        return fmt.Sprintf("string->interface映射(长度%d)", len(val))
    case error:
        return fmt.Sprintf("错误: %s", val.Error())
    case fmt.Stringer:
        return fmt.Sprintf("Stringer: %s", val.String())
    default:
        return fmt.Sprintf("未知类型 %T", val)
    }
}

// ========== 接口类型切换 ==========
func processReader(r io.Reader) {
    switch v := r.(type) {
    case *os.File:
        fmt.Printf("文件: %v\n", v.Name())
    case *io.LimitedReader:
        fmt.Printf("限制读取器: 限制%d\n", v.N)
    case io.ReadCloser:
        fmt.Println("可读可关闭")
        v.Close()
    case io.Reader:
        fmt.Println("普通读取器")
    default:
        fmt.Println("未知读取器类型")
    }
}

// ========== 多类型case ==========
func isNumeric(v interface{}) bool {
    switch v.(type) {
    case int, int8, int16, int32, int64:
        return true
    case uint, uint8, uint16, uint32, uint64, uintptr:
        return true
    case float32, float64:
        return true
    case complex64, complex128:
        return true
    default:
        return false
    }
}

func main() {
    // 类型断言
    fmt.Println("=== 类型断言 ===")
    var i interface{} = "hello"

    // 安全断言
    if s, ok := i.(string); ok {
        fmt.Printf("字符串: %s\n", s)
    }

    // 不安全断言（确定类型时使用）
    s := i.(string)
    fmt.Printf("确定是字符串: %s\n", s)

    // 类型切换
    fmt.Println("\n=== 类型切换 ===")
    values := []interface{}{
        42,
        "hello",
        3.14,
        true,
        nil,
        []int{1, 2, 3},
        map[string]int{"a": 1},
        MyString("custom"),
    }

    for _, v := range values {
        fmt.Println(describe(v))
    }

    // 数值检查
    fmt.Println("\n=== 数值检查 ===")
    fmt.Printf("42是数值: %v\n", isNumeric(42))
    fmt.Printf("\"hello\"是数值: %v\n", isNumeric("hello"))
    fmt.Printf("3.14是数值: %v\n", isNumeric(3.14))

    // 接口类型切换
    fmt.Println("\n=== 接口类型切换 ===")
    processReader(os.Stdin)
}
```

#### 反例说明

```go
// ❌ 错误：类型断言到接口
type Reader interface{ Read([]byte) (int, error) }
var i interface{} = &os.File{}
// r := i.(Reader)  // 可以，但通常不需要

// ❌ 错误：类型断言失败panic
var i interface{} = "hello"
// n := i.(int)  // panic

// ❌ 错误：类型switch中的变量作用域
switch v := i.(type) {
case int:
    fmt.Println(v)  // v是int
case string:
    fmt.Println(v)  // v是string
}
// fmt.Println(v)  // 编译错误：v未定义

// ❌ 错误：case顺序导致问题
switch i.(type) {
case interface{}:  // 匹配所有类型
    fmt.Println("任意类型")
case int:  // 永远不会执行
    fmt.Println("整数")
}

// ❌ 错误：nil接口的类型断言
var i interface{} = nil
// s := i.(string)  // panic
s, ok := i.(string)  // ok=false，s=""
```

#### 最佳实践

1. **优先使用ok模式**进行安全断言
2. **类型switch**case按特异性排序
3. **nil检查**在类型断言之前
4. **接口类型**case放在具体类型之后
5. **避免过度使用**类型断言，考虑重新设计

---

### 5.4 接口嵌套与组合

#### 概念定义

接口可以嵌入其他接口，形成组合接口。这是Go实现接口复用的主要方式。

#### 语法形式

```go
// 接口嵌入
type Combined interface {
    Interface1
    Interface2
    AdditionalMethod()
}

// 标准库示例
type ReadWriter interface {
    Reader
    Writer
}
```

#### 完整示例

```go
package main

import (
    "fmt"
    "io"
)

// ========== 基础接口 ==========
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}

// ========== 组合接口 ==========
// 读写接口
type ReadWriter interface {
    Reader
    Writer
}

// 读写关闭接口
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// 完整文件接口
type FileInterface interface {
    Reader
    Writer
    Closer
    Seeker
    Name() string
    Stat() (FileInfo, error)
}

type FileInfo interface {
    Name() string
    Size() int64
}

// ========== 实现组合接口 ==========
type MyReadWriter struct {
    data []byte
    pos  int
}

func (m *MyReadWriter) Read(p []byte) (n int, err error) {
    if m.pos >= len(m.data) {
        return 0, io.EOF
    }
    n = copy(p, m.data[m.pos:])
    m.pos += n
    return n, nil
}

func (m *MyReadWriter) Write(p []byte) (n int, err error) {
    m.data = append(m.data, p...)
    return len(p), nil
}

// ========== 接口组合与类型转换 ==========
func useReader(r Reader) {
    buf := make([]byte, 100)
    n, _ := r.Read(buf)
    fmt.Printf("读取了 %d 字节\n", n)
}

func useWriter(w Writer) {
    n, _ := w.Write([]byte("hello"))
    fmt.Printf("写入了 %d 字节\n", n)
}

func useReadWriter(rw ReadWriter) {
    useReader(rw)
    useWriter(rw)
}

// ========== 接口兼容性 ==========
func demonstrateCompatibility() {
    var rw ReadWriter = &MyReadWriter{}

    // ReadWriter可以赋值给Reader
    var r Reader = rw
    fmt.Printf("ReadWriter -> Reader: %v\n", r != nil)

    // ReadWriter可以赋值给Writer
    var w Writer = rw
    fmt.Printf("ReadWriter -> Writer: %v\n", w != nil)

    // 反向不行（需要类型断言）
    // var rw2 ReadWriter = r  // 编译错误
}

// ========== 自定义组合接口 ==========
type Logger interface {
    Log(message string)
}

type Metrics interface {
    Record(name string, value float64)
}

type InstrumentedService interface {
    Logger
    Metrics
    Process(data []byte) error
}

type MyService struct{}

func (s *MyService) Log(message string) {
    fmt.Printf("[LOG] %s\n", message)
}

func (s *MyService) Record(name string, value float64) {
    fmt.Printf("[METRIC] %s = %f\n", name, value)
}

func (s *MyService) Process(data []byte) error {
    s.Log("Processing data")
    s.Record("data_size", float64(len(data)))
    return nil
}

func main() {
    // 使用组合接口
    fmt.Println("=== 组合接口 ===")
    rw := &MyReadWriter{}
    rw.Write([]byte("hello, world"))
    useReadWriter(rw)

    // 接口兼容性
    fmt.Println("\n=== 接口兼容性 ===")
    demonstrateCompatibility()

    // 自定义组合接口
    fmt.Println("\n=== 自定义组合接口 ===")
    var service InstrumentedService = &MyService{}
    service.Process([]byte("test data"))

    // 嵌入标准库接口
    fmt.Println("\n=== 标准库接口组合 ===")
    fmt.Printf("io.ReadWriter包含Reader和Writer\n")
    fmt.Printf("io.ReadWriteCloser包含Reader、Writer和Closer\n")
    fmt.Printf("io.ReadWriteSeeker包含Reader、Writer和Seeker\n")
}
```

#### 反例说明

```go
// ❌ 错误：循环接口嵌入
// type A interface {
//     B
// }
// type B interface {
//     A  // 编译错误：interface type loop
// }

// ❌ 错误：方法冲突
type I1 interface { Method() }
type I2 interface { Method() int }  // 同名但签名不同
// type I3 interface {
//     I1
//     I2  // 编译错误：duplicate method Method
// }

// ❌ 错误：接口嵌入具体类型
// type Bad interface {
//     MyStruct  // 编译错误：interface contains embedded non-interface
// }

// ❌ 错误：接口赋值时的方法集不匹配
type Small interface { A() }
type Large interface { A(); B() }
var s Small = &MyType{}
// var l Large = s  // 编译错误
```

#### 最佳实践

1. **小接口原则**：接口应该小而专注
2. **组合优于继承**：用接口组合实现复用
3. **标准库接口**：优先使用io、fmt等标准接口
4. **接口命名**：组合接口用`ReadWriter`形式
5. **文档说明**：说明接口的语义和使用场景

---

## 6. 指针与内存

### 6.1 指针操作

#### 概念定义

指针存储变量的内存地址。Go的指针比C更安全：不支持指针运算、不能转换为任意指针类型。

#### 语法形式

```go
// 声明指针
var p *int

// 取地址
p = &x

// 解引用
value := *p

// 通过指针修改
*p = 100
```

#### 完整示例

```go
package main

import "fmt"

func main() {
    // ========== 基本指针操作 ==========
    x := 42
    p := &x  // p是*int，存储x的地址

    fmt.Printf("x = %d\n", x)
    fmt.Printf("p = %p\n", p)      // 地址
    fmt.Printf("*p = %d\n", *p)    // 解引用

    // 通过指针修改
    *p = 100
    fmt.Printf("修改后 x = %d\n", x)

    // ========== nil指针 ==========
    var nilPtr *int
    fmt.Printf("nil指针: %v\n", nilPtr)
    // fmt.Printf("*nilPtr = %d\n", *nilPtr)  // panic: nil pointer dereference

    // 检查nil
    if nilPtr != nil {
        fmt.Println(*nilPtr)
    } else {
        fmt.Println("指针为nil")
    }

    // ========== 指针作为参数 ==========
    a, b := 10, 20
    fmt.Printf("交换前: a=%d, b=%d\n", a, b)
    swap(&a, &b)
    fmt.Printf("交换后: a=%d, b=%d\n", a, b)

    // ========== 指针返回值 ==========
    ptr := createInt(999)
    fmt.Printf("创建的整数: %d\n", *ptr)

    // ========== 指针与数组 ==========
    arr := [3]int{1, 2, 3}
    arrPtr := &arr
    fmt.Printf("数组指针: %v\n", (*arrPtr)[0])
    fmt.Printf("简化语法: %v\n", arrPtr[0])  // 自动解引用

    // ========== 指针与结构体 ==========
    type Point struct{ X, Y int }
    pt := Point{10, 20}
    ptPtr := &pt
    fmt.Printf("点: (%d, %d)\n", ptPtr.X, ptPtr.Y)  // 自动解引用

    // 修改
    ptPtr.X = 100
    fmt.Printf("修改后: (%d, %d)\n", pt.X, pt.Y)

    // ========== 多级指针 ==========
    x2 := 5
    p2 := &x2
    pp := &p2
    fmt.Printf("x2 = %d\n", x2)
    fmt.Printf("p2 = %p, *p2 = %d\n", p2, *p2)
    fmt.Printf("pp = %p, *pp = %p, **pp = %d\n", pp, *pp, **pp)

    // ========== 指针比较 ==========
    p1 := &x
    p2 = &x
    fmt.Printf("p1 == p2: %v\n", p1 == p2)  // true，指向同一地址

    y := 42
    p3 := &y
    fmt.Printf("p1 == p3: %v\n", p1 == p3)  // false，不同地址
}

func swap(a, b *int) {
    *a, *b = *b, *a
}

func createInt(v int) *int {
    p := new(int)
    *p = v
    return p
}
```

#### 反例说明

```go
// ❌ 错误：指针运算（Go不支持）
// arr := [5]int{1, 2, 3, 4, 5}
// p := &arr[0]
// p++  // 编译错误：invalid operation

// ❌ 错误：任意类型转换
// var p *int
// var q *float64 = (*float64)(p)  // 编译错误

// ❌ 错误：悬垂指针（Go有GC，但注意闭包）
func dangling() *int {
    x := 10
    return &x  // 合法，Go会逃逸分析分配到堆上
}

// ❌ 错误：nil指针调用方法
type MyStruct struct{}
func (m *MyStruct) Method() { /* ... */ }
var s *MyStruct = nil
s.Method()  // 可以调用，但如果访问字段会panic

// ❌ 错误：指针与值混淆
func modify(p *int) {
    p = nil  // 只修改局部副本
}
func main() {
    x := 10
    p := &x
    modify(p)
    fmt.Println(p)  // 不是nil！
}
```

#### 最佳实践

1. **检查nil指针**再解引用
2. **指针接收者**方法检查nil
3. **避免多级指针**，降低复杂度
4. **指针参数**用于修改调用者数据
5. **返回值指针**注意生命周期

---

### 6.2 new与make区别

#### 概念定义

- **new(T)**：分配零值内存，返回*T
- **make(T, args)**：初始化slice、map、channel，返回T（不是指针）

#### 对比表

| 特性 | new | make |
|------|-----|------|
| 用途 | 分配零值内存 | 初始化slice/map/channel |
| 返回类型 | *T | T（slice/map/channel本身） |
| 初始化 | 零值 | 初始化内部结构 |
| 适用类型 | 任意类型 | 仅slice、map、channel |

#### 完整示例

```go
package main

import "fmt"

func main() {
    // ========== new ==========
    // new分配零值内存，返回指针
    p1 := new(int)       // *int，值为0
    p2 := new(string)    // *string，值为""
    p3 := new([3]int)    // *[3]int，值为[0,0,0]

    fmt.Printf("new(int): %v, *p = %d\n", p1, *p1)
    fmt.Printf("new(string): %v, *p = %q\n", p2, *p2)
    fmt.Printf("new([3]int): %v, *p = %v\n", p3, *p3)

    // new结构体
    type Person struct {
        Name string
        Age  int
    }
    p4 := new(Person)  // *Person，值为{"", 0}
    fmt.Printf("new(Person): %+v\n", *p4)

    // 立即初始化
    p5 := &Person{Name: "Alice", Age: 30}  // 等价于new+初始化
    fmt.Printf("&Person{}: %+v\n", *p5)

    // ========== make ==========
    // make初始化slice、map、channel

    // make slice
    s1 := make([]int, 5)       // 长度5，容量5
    s2 := make([]int, 3, 10)   // 长度3，容量10
    fmt.Printf("make slice: len=%d, cap=%d, %v\n", len(s1), cap(s1), s1)
    fmt.Printf("make slice with cap: len=%d, cap=%d, %v\n", len(s2), cap(s2), s2)

    // make map
    m1 := make(map[string]int)
    m1["key"] = 100  // 可以写入
    fmt.Printf("make map: %v\n", m1)

    m2 := make(map[string]int, 100)  // 预分配容量
    fmt.Printf("make map with hint: len=%d\n", len(m2))

    // make channel
    ch1 := make(chan int)       // 无缓冲
    ch2 := make(chan int, 10)   // 有缓冲，容量10
    fmt.Printf("make channel: %v, %v\n", ch1, ch2)

    // ========== 对比 ==========
    // new vs &T{}
    p6 := new(Person)
    p7 := &Person{}
    fmt.Printf("new vs &: %v, %v\n", *p6, *p7)  // 相同

    // new slice vs make slice
    p8 := new([]int)   // *[]int，值为nil
    s3 := make([]int, 5)  // []int，值为[0,0,0,0,0]
    fmt.Printf("new slice: %v, *p = %v\n", p8, *p8)
    fmt.Printf("make slice: %v\n", s3)

    // new map vs make map
    p9 := new(map[string]int)   // *map[string]int，值为nil
    m3 := make(map[string]int)  // map[string]int，可写入
    fmt.Printf("new map: %v, *p = %v\n", p9, *p9)
    fmt.Printf("make map: %v\n", m3)

    // 使用new的map（需要解引用）
    *p9 = make(map[string]int)
    (*p9)["key"] = 200
    fmt.Printf("new+make map: %v\n", *p9)
}
```

#### 反例说明

```go
// ❌ 错误：make用于非slice/map/channel
// x := make(int)  // 编译错误
// y := make(Person)  // 编译错误

// ❌ 错误：new用于slice/map/channel后直接使用
var s *[]int = new([]int)
// (*s)[0] = 1  // panic: index out of range（*s是nil切片）

var m *map[string]int = new(map[string]int)
// (*m)["key"] = 1  // panic: assignment to entry in nil map

// ❌ 错误：混淆new和make的返回值
s := make([]int, 5)  // 返回[]int
p := new([]int)      // 返回*[]int
// s = p  // 编译错误

// ❌ 错误：make的容量参数误解
s := make([]int, 5, 3)  // 编译错误：len > cap
```

#### 最佳实践

1. **slice/map/channel**使用make
2. **其他类型**使用new或&Type{}
3. **预分配容量**减少重新分配
4. **优先使用字面量**如`&Person{}`而非`new(Person)`
5. **nil检查**new后的指针

---

### 6.3 逃逸分析

#### 概念定义

逃逸分析是编译器决定变量分配在栈上还是堆上的过程。如果变量在函数返回后仍被引用，则分配到堆上（"逃逸"）。

#### 逃逸场景

| 场景 | 是否逃逸 | 说明 |
|------|----------|------|
| 返回局部变量指针 | 是 | 函数外仍需访问 |
| 闭包引用局部变量 | 是 | 闭包可能在外部执行 |
| 大对象 | 可能 | 栈空间限制 |
| 接口值包含指针 | 是 | 接口动态类型 |
| 切片/Map扩容 | 是 | 需要连续内存 |

#### 完整示例

```go
package main

import "fmt"

// ========== 逃逸场景 ==========

// 1. 返回局部变量指针 - 逃逸到堆
func escapePtr() *int {
    x := 10  // x逃逸到堆
    return &x
}

// 2. 闭包引用 - 逃逸到堆
func makeClosure() func() int {
    x := 10  // x逃逸到堆
    return func() int {
        x++
        return x
    }
}

// 3. 接口值 - 逃逸到堆
func escapeInterface() interface{} {
    x := 10
    return x  // x逃逸（装箱到接口）
}

// 4. 大对象 - 可能逃逸
func largeObject() [1000000]int {
    var arr [1000000]int  // 大数组可能逃逸
    return arr
}

// 5. 切片扩容 - 逃逸
func sliceEscape() []int {
    s := make([]int, 0, 2)
    for i := 0; i < 100; i++ {
        s = append(s, i)  // 扩容后可能逃逸
    }
    return s
}

// ========== 不逃逸场景 ==========

// 1. 局部变量，不返回
func noEscape() {
    x := 10  // 分配到栈
    fmt.Println(x)
}

// 2. 传值返回
func noEscapeValue() int {
    x := 10  // 分配到栈
    return x  // 复制值
}

// 3. 指针参数，不返回
func noEscapePtr(p *int) {
    *p = 10  // p指向调用者的变量
}

// ========== 优化示例 ==========

// 不好的：返回指针导致逃逸
func badNewPerson() *Person {
    p := &Person{Name: "test"}  // p逃逸
    return p
}

// 好的：传指针参数，不逃逸
func goodNewPerson(p *Person) {
    p.Name = "test"
    p.Age = 30
}

// 不好的：接口参数导致装箱
func badProcess(v interface{}) {
    fmt.Println(v)
}

// 好的：泛型避免装箱（Go 1.18+）
func goodProcess[T any](v T) {
    fmt.Println(v)
}

type Person struct {
    Name string
    Age  int
}

func main() {
    // 逃逸示例
    p := escapePtr()
    fmt.Println(*p)

    counter := makeClosure()
    fmt.Println(counter())
    fmt.Println(counter())

    iface := escapeInterface()
    fmt.Println(iface)

    // 不逃逸示例
    noEscape()
    v := noEscapeValue()
    fmt.Println(v)

    x := 0
    noEscapePtr(&x)
    fmt.Println(x)

    // 优化示例
    person := &Person{}
    goodNewPerson(person)
    fmt.Printf("%+v\n", person)
}
```

#### 反例说明

```go
// ❌ 错误：不必要的指针返回
type Config struct {
    Debug bool
}
func getConfig() *Config {
    return &Config{Debug: true}  // 逃逸到堆
}
// 应该返回值
func getConfigValue() Config {
    return Config{Debug: true}  // 可能不逃逸
}

// ❌ 错误：接口参数导致装箱
func process(items []interface{}) {
    for _, item := range items {
        _ = item
    }
}
// 应该使用泛型
func processGeneric[T any](items []T) {
    for _, item := range items {
        _ = item
    }
}

// ❌ 错误：闭包捕获大对象
func badClosure() func() {
    bigData := make([]byte, 1<<30)  // 1GB
    return func() {
        _ = bigData[0]  // 闭包引用，无法释放
    }
}

// ❌ 错误：全局变量持有引用
var globalCache map[string]*bigStruct
func cacheData(key string, data *bigStruct) {
    globalCache[key] = data  // 永远不会释放
}
```

#### 最佳实践

1. **使用`-m`标志**查看逃逸分析：`go build -gcflags='-m'`
2. **避免不必要的指针返回**
3. **使用值类型**减少堆分配
4. **注意闭包**捕获的变量
5. **泛型替代空接口**减少装箱

---

## 7. 错误处理

### 7.1 error接口

#### 概念定义

error是Go的内置接口，用于表示错误状态。任何实现了`Error() string`方法的类型都可以作为错误。

#### 语法形式

```go
// error接口定义
type error interface {
    Error() string
}

// 创建错误
err := errors.New("error message")
err := fmt.Errorf("formatted error: %v", value)
```

#### 完整示例

```go
package main

import (
    "errors"
    "fmt"
)

// ========== 自定义错误类型 ==========
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// 带错误码的错误
type CodedError struct {
    Code    int
    Message string
}

func (e *CodedError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// ========== 错误常量 ==========
var (
    ErrNotFound     = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
    ErrTimeout      = errors.New("operation timed out")
)

// ========== 错误处理函数 ==========
func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field:   "age",
            Message: "cannot be negative",
        }
    }
    if age > 150 {
        return &ValidationError{
            Field:   "age",
            Message: "unrealistic value",
        }
    }
    return nil
}

func findUser(id int) (string, error) {
    if id <= 0 {
        return "", ErrInvalidInput
    }
    if id > 100 {
        return "", ErrNotFound
    }
    return fmt.Sprintf("User%d", id), nil
}

// ========== 错误检查 ==========
func isValidationError(err error) bool {
    _, ok := err.(*ValidationError)
    return ok
}

func getErrorCode(err error) int {
    if coded, ok := err.(*CodedError); ok {
        return coded.Code
    }
    return 0
}

func main() {
    // ========== 基本错误创建 ==========
    err1 := errors.New("simple error")
    fmt.Println("基本错误:", err1)

    err2 := fmt.Errorf("formatted error: %d", 404)
    fmt.Println("格式化错误:", err2)

    // ========== 自定义错误 ==========
    if err := validateAge(-5); err != nil {
        fmt.Println("验证错误:", err)
    }

    if err := validateAge(200); err != nil {
        fmt.Println("验证错误:", err)
        if isValidationError(err) {
            fmt.Println("是验证错误类型")
        }
    }

    // ========== 错误常量 ==========
    if _, err := findUser(-1); err != nil {
        fmt.Println("查找错误:", err)
        if errors.Is(err, ErrInvalidInput) {
            fmt.Println("是输入错误")
        }
    }

    if _, err := findUser(999); err != nil {
        fmt.Println("查找错误:", err)
        if errors.Is(err, ErrNotFound) {
            fmt.Println("是未找到错误")
        }
    }

    // ========== 错误链（Go 1.13+） ==========
    wrappedErr := fmt.Errorf("operation failed: %w", ErrNotFound)
    fmt.Println("包装错误:", wrappedErr)

    if errors.Is(wrappedErr, ErrNotFound) {
        fmt.Println("原始错误是ErrNotFound")
    }

    // ========== 多错误（Go 1.20+） ==========
    var errs []error
    errs = append(errs, errors.New("error 1"))
    errs = append(errs, errors.New("error 2"))
    errs = append(errs, errors.New("error 3"))

    combined := errors.Join(errs...)
    fmt.Println("合并错误:", combined)
}
```

#### 反例说明

```go
// ❌ 错误：忽略错误返回值
file, _ := os.Open("file.txt")  // 危险！

// ❌ 错误：比较错误字符串
if err.Error() == "not found" {  // 不可靠
    // ...
}

// ❌ 错误：自定义错误不实现Error()
type MyError struct{}
// func (e MyError) Error() string { return "" }  // 忘记实现

// ❌ 错误：返回nil错误但仍有返回值
func badFunc() (int, error) {
    return 0, nil  // 如果出错应该返回错误
}

// ❌ 错误：panic代替错误返回
func badPanic(id int) string {
    if id <= 0 {
        panic("invalid id")  // 不要这样做
    }
    return fmt.Sprintf("User%d", id)
}
```

#### 最佳实践

1. **检查所有错误**，不要忽略
2. **错误常量**用于可比较的错误
3. **自定义错误类型**携带上下文信息
4. **错误包装**使用`%w`（Go 1.13+）
5. **errors.Is/errors.As**检查错误（Go 1.13+）

---

### 7.2 panic与recover

#### 概念定义

- **panic**：停止正常执行，开始展开栈
- **recover**：捕获panic，恢复执行（只能在defer中使用）

#### 语法形式

```go
// 触发panic
panic(value)

// 恢复panic
func protected() {
    defer func() {
        if r := recover(); r != nil {
            // 处理panic
        }
    }()
    // 可能panic的代码
}
```

#### 完整示例

```go
package main

import (
    "fmt"
    "runtime/debug"
)

// ========== panic示例 ==========
func mayPanic(shouldPanic bool) {
    if shouldPanic {
        panic("something went wrong")
    }
    fmt.Println("No panic")
}

// ========== recover示例 ==========
func safeCall(fn func()) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
            fmt.Printf("Stack trace:\n%s\n", debug.Stack())
        }
    }()
    fn()
}

// ========== 保护函数 ==========
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    if b == 0 {
        panic("division by zero")
    }
    return a / b, nil
}

// ========== 多层panic/recover ==========
func level3() {
    panic("panic from level 3")
}

func level2() {
    level3()
}

func level1() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered at level 1: %v\n", r)
        }
    }()
    level2()
}

// ========== panic值类型 ==========
func panicWithError() {
    panic(errors.New("this is an error"))
}

func panicWithString() {
    panic("this is a string")
}

func panicWithInt() {
    panic(42)
}

// ========== 不可恢复的panic ==========
func unrecoverable() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
            panic("re-panic")  // 重新panic
        }
    }()
    panic("original panic")
}

func main() {
    // ========== 基本panic/recover ==========
    fmt.Println("=== 基本panic/recover ===")
    safeCall(func() {
        mayPanic(true)
    })
    fmt.Println("继续执行")

    // ========== 安全除法 ==========
    fmt.Println("\n=== 安全除法 ===")
    result, err := safeDivide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("10 / 2 = %d\n", result)
    }

    result, err = safeDivide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    }

    // ========== 多层调用 ==========
    fmt.Println("\n=== 多层调用 ===")
    level1()
    fmt.Println("level1后继续执行")

    // ========== 不同类型的panic值 ==========
    fmt.Println("\n=== panic值类型 ===")

    safeCall(panicWithError)
    safeCall(panicWithString)
    safeCall(panicWithInt)

    // ========== 重新panic ==========
    fmt.Println("\n=== 重新panic ===")
    safeCall(unrecoverable)

    fmt.Println("\n程序正常结束")
}

// 错误类型
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}

func errors.New(s string) error {
    return &errorString{s}
}
```

#### 反例说明

```go
// ❌ 错误：在普通函数中使用recover
func badRecover() {
    if r := recover(); r != nil {  // 无效，recover返回nil
        // ...
    }
}

// ❌ 错误：recover后没有处理
func uselessRecover() {
    defer func() {
        recover()  // 捕获但忽略
    }()
    panic("error")
}

// ❌ 错误：用panic代替错误返回
func badAPI(id int) string {
    if id <= 0 {
        panic("invalid id")  // 不要这样做
    }
    return fmt.Sprintf("User%d", id)
}

// ❌ 错误：goroutine中的panic无法被外部recover
go func() {
    panic("goroutine panic")  // 会导致整个程序崩溃
}()

// ❌ 错误：defer中的panic
func deferPanic() {
    defer func() {
        panic("defer panic")  // 会覆盖原panic
    }()
    panic("original panic")
}
```

#### 最佳实践

1. **panic用于不可恢复的错误**（程序bug）
2. **recover在defer中使用**
3. **不要忽略recover的值**
4. **goroutine中的panic**需要内部recover
5. **错误返回**优于panic

---

### 7.3 错误链与包装

#### 概念定义

错误链通过包装错误添加上下文，形成错误链路。Go 1.13+提供了`errors.Is`和`errors.As`支持错误链操作。

#### 语法形式

```go
// 错误包装（Go 1.13+）
err := fmt.Errorf("context: %w", originalErr)

// 检查错误链中是否包含特定错误
if errors.Is(err, targetErr) { ... }

// 提取错误链中的特定错误类型
var target *MyError
if errors.As(err, &target) { ... }
```

#### 完整示例

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

// ========== 自定义错误类型 ==========
type NotFoundError struct {
    Resource string
    ID       int
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with id %d not found", e.Resource, e.ID)
}

func (e *NotFoundError) Is(target error) bool {
    t, ok := target.(*NotFoundError)
    if !ok {
        return false
    }
    return e.Resource == t.Resource && e.ID == t.ID
}

type PermissionError struct {
    User   string
    Action string
}

func (e *PermissionError) Error() string {
    return fmt.Sprintf("user %s has no permission to %s", e.User, e.Action)
}

// ========== 错误包装示例 ==========
var (
    ErrDatabase = errors.New("database error")
    ErrNetwork  = errors.New("network error")
    ErrCache    = errors.New("cache error")
)

func queryDatabase(id int) (string, error) {
    if id < 0 {
        return "", fmt.Errorf("invalid id %d: %w", id, ErrDatabase)
    }
    if id == 0 {
        return "", &NotFoundError{Resource: "user", ID: id}
    }
    return fmt.Sprintf("User%d", id), nil
}

func getFromCache(id int) (string, error) {
    result, err := queryDatabase(id)
    if err != nil {
        return "", fmt.Errorf("cache miss for id %d: %w", id, err)
    }
    return result, nil
}

func getUser(id int) (string, error) {
    result, err := getFromCache(id)
    if err != nil {
        return "", fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return result, nil
}

// ========== 错误检查 ==========
func checkError(err error) {
    fmt.Printf("\n检查错误: %v\n", err)

    // 检查是否是特定错误
    if errors.Is(err, ErrDatabase) {
        fmt.Println("  -> 包含数据库错误")
    }

    if errors.Is(err, ErrNetwork) {
        fmt.Println("  -> 包含网络错误")
    }

    // 检查是否是特定类型
    var notFound *NotFoundError
    if errors.As(err, &notFound) {
        fmt.Printf("  -> 是NotFoundError: %s (ID=%d)\n",
            notFound.Resource, notFound.ID)
    }

    var perm *PermissionError
    if errors.As(err, &perm) {
        fmt.Printf("  -> 是PermissionError: %s cannot %s\n",
            perm.User, perm.Action)
    }

    // 解包错误链
    fmt.Println("  错误链:")
    for err != nil {
        fmt.Printf("    - %v\n", err)
        err = errors.Unwrap(err)
    }
}

// ========== 多错误（Go 1.20+） ==========
func validateUser(name string, age int) error {
    var errs []error

    if name == "" {
        errs = append(errs, errors.New("name is required"))
    }
    if age < 0 {
        errs = append(errs, errors.New("age cannot be negative"))
    }
    if age > 150 {
        errs = append(errs, errors.New("age is unrealistic"))
    }

    return errors.Join(errs...)
}

// ========== 标准库错误 ==========
func fileOperations() {
    _, err := os.Open("nonexistent.txt")
    if err != nil {
        fmt.Printf("\n文件错误: %v\n", err)

        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("  -> 文件不存在")
        }

        if errors.Is(err, os.ErrPermission) {
            fmt.Println("  -> 权限不足")
        }
    }

    // io.EOF
    fmt.Printf("\nio.EOF检查: %v\n", errors.Is(io.EOF, io.EOF))
}

func main() {
    // ========== 错误链 ==========
    fmt.Println("=== 错误链 ===")

    err := getUser(-1)
    checkError(err)

    err = getUser(0)
    checkError(err)

    // ========== 多错误 ==========
    fmt.Println("\n=== 多错误 ===")
    err = validateUser("", -5)
    fmt.Printf("验证错误: %v\n", err)

    // 遍历多错误
    if joined, ok := err.(interface{ Unwrap() []error }); ok {
        fmt.Println("  包含的错误:")
        for _, e := range joined.Unwrap() {
            fmt.Printf("    - %v\n", e)
        }
    }

    // ========== 标准库错误 ==========
    fmt.Println("\n=== 标准库错误 ===")
    fileOperations()
}
```

#### 反例说明

```go
// ❌ 错误：使用%v而不是%w包装
err := fmt.Errorf("context: %v", originalErr)  // 不能解包
// 应该用 %w
err := fmt.Errorf("context: %w", originalErr)  // 可以解包

// ❌ 错误：自定义Is方法不正确
type MyError struct{ Code int }
func (e *MyError) Is(target error) bool {
    // 忘记检查target类型
    return e.Code == target.(*MyError).Code  // 可能panic
}
// 修复
func (e *MyError) Is(target error) bool {
    t, ok := target.(*MyError)
    if !ok {
        return false
    }
    return e.Code == t.Code
}

// ❌ 错误：忽略错误链中的信息
if err != nil {
    return err  // 丢失上下文
}
// 应该
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// ❌ 错误：直接比较错误而不是用errors.Is
if err == io.EOF {  // 可能错过包装后的EOF
    // ...
}
// 应该
if errors.Is(err, io.EOF) {
    // ...
}
```

#### 最佳实践

1. **包装错误**使用`%w`保留原始错误
2. **errors.Is**检查错误链中的特定错误
3. **errors.As**提取特定错误类型
4. **自定义错误**实现Is方法支持精确匹配
5. **添加上下文**帮助定位问题

---

## 8. 包管理与模块

### 8.1 GOPATH到Go Modules演进

#### 概念定义

Go的包管理经历了从GOPATH到Go Modules的演进：

| 特性 | GOPATH (Go 1.0-1.10) | Go Modules (Go 1.11+) |
|------|---------------------|----------------------|
| 依赖位置 | $GOPATH/src | 项目目录下 |
| 版本控制 | 无 | 语义化版本 |
| 依赖隔离 | 全局 | 项目级 |
| 可重现构建 | 否 | 是 |
| 代理支持 | 否 | 是 |

#### 完整示例

```go
// 项目结构（Go Modules）
// myproject/
// ├── go.mod          # 模块定义
// ├── go.sum          # 依赖校验和
// ├── main.go         # 主程序
// ├── internal/       # 内部包
// │   └── helper.go
// └── pkg/            # 可导入包
//     └── utils.go

// go.mod 示例
module github.com/example/myproject

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/stretchr/testify v1.8.4
)

require (
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
)
```

#### 反例说明

```go
// ❌ 旧式GOPATH项目结构
// $GOPATH/src/github.com/user/project/
// 现代Go应该使用Go Modules

// ❌ 混合GOPATH和Modules
// 设置GO111MODULE=on强制使用Modules

// ❌ 手动修改go.mod
// 应该使用go mod edit命令
```

#### 最佳实践

1. **始终使用Go Modules**（Go 1.16+默认启用）
2. **模块路径**使用代码托管地址
3. **语义化版本**发布模块
4. **go.sum提交**到版本控制

---

### 8.2 go.mod与go.sum

#### 概念定义

- **go.mod**：模块定义文件，包含模块路径、Go版本、依赖列表
- **go.sum**：依赖校验和文件，确保依赖完整性

#### go.mod指令

```go
module example.com/myproject      // 模块路径
go 1.22                           // Go版本要求

require (                         // 直接依赖
    github.com/pkg/errors v0.9.1
    github.com/sirupsen/logrus v1.9.3
)

require (                         // 间接依赖
    golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace (                         // 替换依赖
    github.com/old/pkg => github.com/new/pkg v2.0.0
    example.com/local => ./local
)

exclude (                         // 排除版本
    github.com/problem/pkg v1.2.3
)
```

#### 完整示例

```go
// 完整的go.mod示例
module github.com/example/webapp

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-redis/redis/v8 v8.11.5
    github.com/jmoiron/sqlx v1.3.5
    github.com/lib/pq v1.10.9
    github.com/spf13/viper v1.16.0
    github.com/stretchr/testify v1.8.4
    go.uber.org/zap v1.24.0
    gorm.io/driver/postgres v1.5.2
    gorm.io/gorm v1.25.2
)

require (
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/cespare/xxhash/v2 v2.2.0 // indirect
    github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
    github.com/davecgh/go-spew v1.1.1 // indirect
    github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
    github.com/fsnotify/fsnotify v1.6.0 // indirect
    github.com/gabriel-vasile/mimetype v1.4.2 // indirect
    github.com/gin-contrib/sse v0.1.0 // indirect
    github.com/go-playground/locales v0.14.1 // indirect
    github.com/go-playground/universal-translator v0.18.1 // indirect
    github.com/go-playground/validator/v10 v10.14.0 // indirect
    github.com/goccy/go-json v0.10.2 // indirect
    github.com/hashicorp/hcl v1.0.0 // indirect
    github.com/jinzhu/inflection v1.0.0 // indirect
    github.com/jinzhu/now v1.1.5 // indirect
    github.com/json-iterator/go v1.1.12 // indirect
    github.com/klauspost/cpuid/v2 v2.2.4 // indirect
    github.com/leodido/go-urn v1.2.4 // indirect
    github.com/magiconair/properties v1.8.7 // indirect
    github.com/mattn/go-isatty v0.0.19 // indirect
    github.com/mitchellh/mapstructure v1.5.0 // indirect
    github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
    github.com/modern-go/reflect2 v1.0.2 // indirect
    github.com/pelletier/go-toml/v2 v2.0.8 // indirect
    github.com/pmezard/go-difflib v1.0.0 // indirect
    github.com/spf13/afero v1.9.5 // indirect
    github.com/spf13/cast v1.5.1 // indirect
    github.com/spf13/jwalterweatherman v1.1.0 // indirect
    github.com/spf13/pflag v1.0.5 // indirect
    github.com/subosito/gotenv v1.4.2 // indirect
    github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
    github.com/ugorji/go/codec v1.2.11 // indirect
    go.uber.org/atomic v1.11.0 // indirect
    go.uber.org/multierr v1.11.0 // indirect
    golang.org/x/arch v0.3.0 // indirect
    golang.org/x/crypto v0.9.0 // indirect
    golang.org/x/net v0.10.0 // indirect
    golang.org/x/sys v0.8.0 // indirect
    golang.org/x/text v0.9.0 // indirect
    google.golang.org/protobuf v1.30.0 // indirect
    gopkg.in/ini.v1 v1.67.0 // indirect
    gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

#### go.sum示例

```
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/gin-contrib/sse v0.1.0 h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=
github.com/gin-contrib/sse v0.1.0/go.mod h1:RHrZQHXnP2xjPF+u1gW/2HnVO7nvIa9PG3Gm+fLHvGI=
```

#### 常用命令

```bash
# 初始化模块
go mod init example.com/myproject

# 下载依赖
go mod download

# 整理依赖（添加缺失、删除未使用）
go mod tidy

# 编辑go.mod
go mod edit -require=github.com/pkg/errors@v0.9.1
go mod edit -replace=github.com/old/pkg=github.com/new/pkg@v2.0.0

# 查看依赖
go list -m all
go list -m -versions github.com/gin-gonic/gin

# 验证依赖
go mod verify

# 清理缓存
go clean -modcache
```

#### 反例说明

```go
// ❌ 错误：手动编辑go.mod导致格式错误
// 使用 go mod edit 命令

// ❌ 错误：不提交go.sum
// go.sum应该提交到版本控制

// ❌ 错误：使用伪版本不规范
// github.com/user/pkg v0.0.0-20230101120000-abcdef123456
// 应该使用语义化版本

// ❌ 错误：循环依赖
// module A 依赖 B，B 依赖 A
```

#### 最佳实践

1. **提交go.sum**确保构建可重现
2. **go mod tidy**保持依赖整洁
3. **语义化版本**发布模块
4. **最小依赖**原则
5. **定期更新**依赖版本

---

### 8.3 语义化版本控制

#### 概念定义

语义化版本（Semantic Versioning）格式：`MAJOR.MINOR.PATCH`

| 版本变化 | 含义 | 兼容性 |
|----------|------|--------|
| MAJOR | 不兼容的API修改 | 不向后兼容 |
| MINOR | 向后兼容的功能新增 | 向后兼容 |
| PATCH | 向后兼容的问题修复 | 向后兼容 |

#### 版本选择规则

```
require github.com/example/pkg v1.2.3  # 精确版本
require github.com/example/pkg v1.2.0  # 最低版本，允许补丁更新
require github.com/example/pkg v1.0.0  # 最低版本，允许小版本更新
```

#### 完整示例

```go
// go.mod 版本指定
module example.com/myproject

go 1.22

require (
    // 精确版本
    github.com/pkgA/pkg v1.2.3

    // 最低版本（Go会选择满足所有依赖的最高版本）
    github.com/pkgB/pkg v1.0.0

    // 预发布版本
    github.com/pkgC/pkg v1.0.0-beta.1

    // 伪版本（基于commit）
    github.com/pkgD/pkg v0.0.0-20230101120000-abcdef123456
)
```

#### 版本兼容性

```go
// v0.x.x - 不稳定版本，API可能随时变化
// v1.x.x - 稳定版本，保证向后兼容
// v2.x.x - 不兼容v1的重大更新

// 导入v2+模块
import "github.com/example/pkg/v2"
```

#### 反例说明

```go
// ❌ 错误：v1和v2混用
import (
    "github.com/example/pkg"    // v1
    "github.com/example/pkg/v2" // v2
)
// 应该统一使用一个版本

// ❌ 错误：不兼容的版本更新
// v1.2.3 删除已导出函数
// 应该是 v2.0.0

// ❌ 错误：预发布版本用于生产
require github.com/pkg v1.0.0-alpha.1  // 不推荐用于生产
```

#### 最佳实践

1. **v1.0.0后**保持向后兼容
2. **破坏性变更**升级主版本号
3. **预发布版本**标记为alpha/beta/rc
4. **版本标签**使用语义化版本
5. **CHANGELOG**记录版本变更

---

### 8.4 私有模块与代理

#### 概念定义

- **GOPROXY**：模块代理服务器，缓存和加速模块下载
- **GOPRIVATE**：指定不通过代理的私有模块
- **GONOPROXY/GONOSUMDB**：更细粒度的控制

#### 配置示例

```bash
# 使用公共代理（中国）
export GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct

# 私有模块配置
export GOPRIVATE=*.company.com,github.com/myorg/*

# 不使用sumdb验证的模块
export GONOSUMDB=*.company.com

# 完整配置
export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=gitlab.company.com
export GONOPROXY=gitlab.company.com
export GONOSUMDB=gitlab.company.com
```

#### 私有模块配置

```bash
# .gitconfig 配置私有Git服务器
[url "git@gitlab.company.com:"]
    insteadOf = https://gitlab.company.com/

# 或者环境变量
export GIT_CONFIG_GLOBAL=~/.gitconfig
```

#### 完整示例

```go
// 使用私有模块
module github.com/company/myproject

go 1.22

require (
    github.com/company/private-lib v1.0.0
    github.com/public/lib v2.0.0
)

replace github.com/company/private-lib => gitlab.company.com/team/private-lib v1.0.0
```

#### 反例说明

```go
// ❌ 错误：私有模块通过公共代理
// 会导致404或泄露信息

// ❌ 错误：GOPRIVATE配置不完整
// export GOPRIVATE=github.com/mycompany
// 应该 export GOPRIVATE=github.com/mycompany/*

// ❌ 错误：混合http和https
// 应该统一使用https
```

#### 最佳实践

1. **私有模块**配置GOPRIVATE
2. **代理选择**使用可靠代理
3. **HTTPS优先**安全传输
4. **SSH密钥**认证私有仓库
5. **CI/CD**配置相同环境变量

---

### 8.5 工作区模式（Workspace）

#### 概念定义

工作区模式（Go 1.18+）允许同时开发多个相关模块，无需发布到版本控制或本地替换。

#### 语法形式

```go
// go.work 文件
go 1.22

use (
    ./module-a
    ./module-b
    ./module-c
)

replace github.com/example/module-a => ./module-a
```

#### 完整示例

```go
// 项目结构
// myworkspace/
// ├── go.work
// ├── module-a/
// │   ├── go.mod
// │   └── a.go
// ├── module-b/
// │   ├── go.mod
// │   └── b.go
// └── module-c/
//     ├── go.mod
//     └── c.go

// go.work
use (
    ./module-a
    ./module-b
    ./module-c
)

// module-a/go.mod
module github.com/example/module-a

go 1.22

// module-a/a.go
package a

func Hello() string {
    return "Hello from module A"
}

// module-b/go.mod
module github.com/example/module-b

go 1.22

require github.com/example/module-a v1.0.0

// module-b/b.go
package b

import "github.com/example/module-a"

func Greet() string {
    return a.Hello() + " via module B"
}
```

#### 常用命令

```bash
# 初始化工作区
go work init ./module-a ./module-b

# 添加模块到工作区
go work use ./module-c

# 编辑工作区
go work edit -replace=github.com/example/pkg=./local/pkg

# 查看工作区
go work sync
```

#### 反例说明

```go
// ❌ 错误：提交go.work到版本控制
// go.work是本地开发配置，不应该提交
// .gitignore: go.work go.work.sum

// ❌ 错误：工作区与replace冲突
// go.work中的replace会覆盖go.mod中的replace

// ❌ 错误：循环依赖
// module-a依赖module-b，module-b依赖module-a
```

#### 最佳实践

1. **不提交go.work**到版本控制
2. **go.mod保持**独立可构建
3. **工作区用于**本地开发
4. **CI/CD**不使用工作区
5. **文档说明**工作区使用方式

---

## 9. 标准库核心

### 9.1 fmt、os、io包

#### 概念定义

- **fmt**：格式化I/O，提供打印、扫描功能
- **os**：操作系统接口，文件、环境变量、进程
- **io**：基本I/O接口，Reader/Writer

#### fmt包核心函数

```go
// 输出到标准输出
fmt.Print(a ...interface{})
fmt.Println(a ...interface{})
fmt.Printf(format string, a ...interface{})

// 输出到字符串
fmt.Sprint(a ...interface{}) string
fmt.Sprintf(format string, a ...interface{}) string
fmt.Sprintln(a ...interface{}) string

// 输出到Writer
fmt.Fprint(w io.Writer, a ...interface{})
fmt.Fprintf(w io.Writer, format string, a ...interface{})
fmt.Fprintln(w io.Writer, a ...interface{})

// 从字符串扫描
fmt.Sscan(str string, a ...interface{}) (n int, err error)
fmt.Sscanf(str string, format string, a ...interface{}) (n int, err error)
```

#### os包核心功能

```go
// 文件操作
os.Create(name string) (*File, error)
os.Open(name string) (*File, error)
os.OpenFile(name string, flag int, perm FileMode) (*File, error)
os.Remove(name string) error
os.Rename(oldpath, newpath string) error

// 目录操作
os.Mkdir(name string, perm FileMode) error
os.MkdirAll(path string, perm FileMode) error
os.RemoveAll(path string) error
os.ReadDir(name string) ([]DirEntry, error)

// 环境变量
os.Getenv(key string) string
os.Setenv(key, value string) error
os.LookupEnv(key string) (string, bool)
os.Environ() []string

// 进程
os.Exit(code int)
os.Getenv(key string) string
os.Getwd() (string, error)
os.Chdir(dir string) error
```

#### io包核心接口

```go
// 基本接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

type ReadWriter interface {
    Reader
    Writer
}

type ReadCloser interface {
    Reader
    Closer
}

// 组合接口
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

#### 完整示例

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

func main() {
    // ========== fmt包 ==========
    fmt.Println("=== fmt包 ===")

    // 基本输出
    fmt.Print("Hello ")
    fmt.Println("World")

    // 格式化输出
    name := "Alice"
    age := 30
    fmt.Printf("Name: %s, Age: %d\n", name, age)

    // 格式化动词
    fmt.Printf("%%v: %v\n", name)       // 默认格式
    fmt.Printf("%%+v: %+v\n", struct{ X int }{10})  // 结构体字段名
    fmt.Printf("%%#v: %#v\n", name)     // Go语法格式
    fmt.Printf("%%T: %T\n", name)       // 类型
    fmt.Printf("%%t: %t\n", true)       // 布尔
    fmt.Printf("%%d: %d\n", 42)         // 十进制整数
    fmt.Printf("%%b: %b\n", 42)         // 二进制
    fmt.Printf("%%x: %x\n", 255)        // 十六进制
    fmt.Printf("%%f: %f\n", 3.14159)    // 浮点数
    fmt.Printf("%%e: %e\n", 3.14159)    // 科学计数法
    fmt.Printf("%%s: %s\n", name)       // 字符串
    fmt.Printf("%%q: %q\n", name)       // 带引号字符串
    fmt.Printf("%%p: %p\n", &age)       // 指针

    // 格式化到字符串
    str := fmt.Sprintf("Name: %s", name)
    fmt.Println("格式化字符串:", str)

    // ========== os包 ==========
    fmt.Println("\n=== os包 ===")

    // 环境变量
    home := os.Getenv("HOME")
    fmt.Printf("HOME: %s\n", home)

    // 工作目录
    wd, _ := os.Getwd()
    fmt.Printf("工作目录: %s\n", wd)

    // 创建临时文件
    tmpFile, err := os.CreateTemp("", "example-*.txt")
    if err != nil {
        fmt.Println("创建临时文件失败:", err)
        return
    }
    defer os.Remove(tmpFile.Name())
    defer tmpFile.Close()

    fmt.Printf("临时文件: %s\n", tmpFile.Name())

    // 写入文件
    content := []byte("Hello, Go!\n")
    n, err := tmpFile.Write(content)
    if err != nil {
        fmt.Println("写入失败:", err)
        return
    }
    fmt.Printf("写入 %d 字节\n", n)

    // 读取文件
    tmpFile.Seek(0, 0)
    data, err := os.ReadFile(tmpFile.Name())
    if err != nil {
        fmt.Println("读取失败:", err)
        return
    }
    fmt.Printf("文件内容: %s", data)

    // ========== io包 ==========
    fmt.Println("\n=== io包 ===")

    // 使用io.Copy
    src := "Hello, io.Copy!"
    reader := strings.NewReader(src)

    var builder strings.Builder
    written, err := io.Copy(&builder, reader)
    if err != nil {
        fmt.Println("Copy失败:", err)
        return
    }
    fmt.Printf("Copy了 %d 字节: %s\n", written, builder.String())

    // 使用io.ReadAll
    reader2 := strings.NewReader("ReadAll test")
    all, err := io.ReadAll(reader2)
    if err != nil {
        fmt.Println("ReadAll失败:", err)
        return
    }
    fmt.Printf("ReadAll: %s\n", all)

    // 使用bufio
    scanner := bufio.NewScanner(strings.NewReader("line1\nline2\nline3"))
    fmt.Println("Scanner读取:")
    for scanner.Scan() {
        fmt.Printf("  %s\n", scanner.Text())
    }

    // 带缓冲的写入
    var bufWriter strings.Builder
    buffered := bufio.NewWriter(&bufWriter)
    buffered.WriteString("Buffered ")
    buffered.WriteString("write")
    buffered.Flush()
    fmt.Printf("缓冲写入: %s\n", bufWriter.String())
}

// 需要导入
import "strings"
```

#### 反例说明

```go
// ❌ 错误：忘记关闭文件
file, _ := os.Open("file.txt")
// defer file.Close()  // 忘记关闭

// ❌ 错误：忽略fmt.Errorf返回值
fmt.Printf("%d", "string")  // 输出: %!d(string=string)

// ❌ 错误：io.EOF处理不当
n, err := reader.Read(buf)
if err != nil {  // io.EOF也是错误！
    return err
}
// 应该
if err != nil && err != io.EOF {
    return err
}

// ❌ 错误：bufio忘记Flush
writer := bufio.NewWriter(file)
writer.WriteString("data")
// writer.Flush()  // 忘记刷新，数据可能丢失
```

#### 最佳实践

1. **defer关闭**文件和资源
2. **检查错误**包括io.EOF
3. **bufio**用于频繁小读写
4. **fmt.Errorf**添加上下文
5. **io.Copy**用于大文件复制

---

### 9.2 net/http包架构

#### 概念定义

net/http提供HTTP客户端和服务器实现，基于Handler接口构建。

#### 核心接口

```go
// Handler接口
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Handler函数类型
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

#### 完整示例

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
)

// ========== 基本Handler ==========
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!\n", r.URL.Path[1:])
}

// ========== 自定义Handler ==========
type timeHandler struct {
    format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Current time: %s\n", time.Now().Format(th.format))
}

// ========== 中间件模式 ==========
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// ========== REST API示例 ==========
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var users = []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    users = append(users, user)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// ========== HTTP客户端 ==========
func httpClientExample() {
    // 创建客户端
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // GET请求
    resp, err := client.Get("https://api.github.com/users/github")
    if err != nil {
        log.Println("请求失败:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("响应状态: %s\n", resp.Status)
    fmt.Printf("响应长度: %d\n", len(body))

    // POST请求
    jsonData := `{"name":"test"}`
    resp2, err := client.Post(
        "https://httpbin.org/post",
        "application/json",
        strings.NewReader(jsonData),
    )
    if err != nil {
        log.Println("POST失败:", err)
        return
    }
    defer resp2.Body.Close()
}

// ========== 路由设置 ==========
func setupRoutes() *http.ServeMux {
    mux := http.NewServeMux()

    // 基本路由
    mux.HandleFunc("/", helloHandler)
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getUsers(w, r)
        case http.MethodPost:
            createUser(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // 自定义Handler
    mux.Handle("/time", timeHandler{format: time.RFC3339})

    return mux
}

func main() {
    // 设置路由
    mux := setupRoutes()

    // 应用中间件
    handler := loggingMiddleware(mux)

    // 创建服务器
    server := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Println("服务器启动在 :8080")
    log.Fatal(server.ListenAndServe())
}

// 需要导入
import "strings"
```

#### 反例说明

```go
// ❌ 错误：忘记关闭响应Body
resp, _ := http.Get("http://example.com")
// defer resp.Body.Close()  // 忘记关闭，资源泄漏

// ❌ 错误：不设置超时
client := &http.Client{}  // 默认无超时，可能永远阻塞
// 应该
client := &http.Client{
    Timeout: 10 * time.Second,
}

// ❌ 错误：并发访问map
var handlers = map[string]http.Handler{}
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    handlers[r.URL.Path] = nil  // 并发不安全！
})

// ❌ 错误：在Handler中panic
func badHandler(w http.ResponseWriter, r *http.Request) {
    panic("something wrong")  // 会导致服务器崩溃
}
```

#### 最佳实践

1. **defer关闭**响应Body
2. **设置超时**防止阻塞
3. **使用中间件**处理横切关注点
4. **优雅关闭**服务器
5. **限制请求大小**防止攻击

---

### 9.3 encoding/json

#### 概念定义

encoding/json提供JSON编码和解码功能，支持结构体标签自定义。

#### 核心函数

```go
// 编码（Go -> JSON）
func Marshal(v interface{}) ([]byte, error)
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
func NewEncoder(w io.Writer) *Encoder

// 解码（JSON -> Go）
func Unmarshal(data []byte, v interface{}) error
func NewDecoder(r io.Reader) *Decoder
```

#### 结构体标签

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email,omitempty"`  // 空值省略
    Password  string    `json:"-"`                // 忽略字段
    Age       int       `json:"age,string"`       // 编码为字符串
    CreatedAt time.Time `json:"created_at"`       // 字段名映射
}
```

#### 完整示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "time"
)

// ========== 基本结构体 ==========
type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Address string `json:"address,omitempty"`
}

// ========== 高级标签 ==========
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`                    // 忽略
    Email     string    `json:"email,omitempty"`      // 空值省略
    Age       int       `json:"age,string"`           // 字符串编码
    IsActive  bool      `json:"is_active"`            // 蛇形命名
    Score     float64   `json:"score"`
    Tags      []string  `json:"tags"`
    Metadata  Metadata  `json:"metadata"`
    CreatedAt time.Time `json:"created_at"`
}

type Metadata struct {
    Source string `json:"source"`
    IP     string `json:"ip"`
}

// ========== 自定义Marshal/Unmarshal ==========
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    duration, err := time.ParseDuration(s)
    if err != nil {
        return err
    }
    *d = Duration(duration)
    return nil
}

type Task struct {
    Name     string   `json:"name"`
    Duration Duration `json:"duration"`
}

func main() {
    // ========== 基本编码 ==========
    fmt.Println("=== 基本编码 ===")
    person := Person{Name: "Alice", Age: 30}

    // Marshal
    data, err := json.Marshal(person)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Marshal: %s\n", data)

    // MarshalIndent
    prettyData, _ := json.MarshalIndent(person, "", "  ")
    fmt.Printf("MarshalIndent:\n%s\n", prettyData)

    // ========== 基本解码 ==========
    fmt.Println("\n=== 基本解码 ===")
    jsonStr := `{"name":"Bob","age":25}`
    var p Person
    if err := json.Unmarshal([]byte(jsonStr), &p); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Unmarshal: %+v\n", p)

    // ========== 高级标签 ==========
    fmt.Println("\n=== 高级标签 ===")
    user := User{
        ID:        1,
        Username:  "alice",
        Password:  "secret",  // 会被忽略
        Email:     "alice@example.com",
        Age:       30,
        IsActive:  true,
        Score:     95.5,
        Tags:      []string{"admin", "user"},
        Metadata:  Metadata{Source: "web", IP: "192.168.1.1"},
        CreatedAt: time.Now(),
    }

    userData, _ := json.MarshalIndent(user, "", "  ")
    fmt.Printf("User JSON:\n%s\n", userData)

    // omitempty演示
    user2 := User{
        ID:       2,
        Username: "bob",
        // Email为空，会被省略
    }
    user2Data, _ := json.Marshal(user2)
    fmt.Printf("omitempty: %s\n", user2Data)

    // ========== 流式处理 ==========
    fmt.Println("\n=== 流式处理 ===")

    // Encoder
    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    encoder.SetIndent("", "  ")
    encoder.Encode(person)
    fmt.Printf("Encoder: %s", buf.String())

    // Decoder
    decoder := json.NewDecoder(strings.NewReader(jsonStr))
    var decodedPerson Person
    decoder.Decode(&decodedPerson)
    fmt.Printf("Decoder: %+v\n", decodedPerson)

    // 处理JSON流
    stream := `
{"name":"user1","age":20}
{"name":"user2","age":25}
{"name":"user3","age":30}
`
    decoder2 := json.NewDecoder(strings.NewReader(stream))
    fmt.Println("JSON流:")
    for decoder2.More() || decoder2.InputOffset() < int64(len(stream)) {
        var p Person
        if err := decoder2.Decode(&p); err != nil {
            break
        }
        fmt.Printf("  %+v\n", p)
    }

    // ========== 自定义Marshal ==========
    fmt.Println("\n=== 自定义Marshal ===")
    task := Task{
        Name:     "Backup",
        Duration: Duration(2 * time.Hour + 30 * time.Minute),
    }
    taskData, _ := json.MarshalIndent(task, "", "  ")
    fmt.Printf("Task JSON:\n%s\n", taskData)

    // 解码
    taskJSON := `{"name":"Restore","duration":"1h30m"}`
    var decodedTask Task
    json.Unmarshal([]byte(taskJSON), &decodedTask)
    fmt.Printf("Decoded Task: %+v, Duration: %v\n",
        decodedTask, time.Duration(decodedTask.Duration))

    // ========== 未知结构JSON ==========
    fmt.Println("\n=== 未知结构JSON ===")
    unknownJSON := `{
        "name": "test",
        "value": 42,
        "nested": {"a": 1, "b": 2},
        "items": [1, 2, 3]
    }`

    var unknown map[string]interface{}
    json.Unmarshal([]byte(unknownJSON), &unknown)
    fmt.Printf("Unknown: %+v\n", unknown)

    // 使用json.RawMessage延迟解析
    type Message struct {
        Type string          `json:"type"`
        Data json.RawMessage `json:"data"`
    }

    msgJSON := `{"type":"user","data":{"name":"Alice"}}`
    var msg Message
    json.Unmarshal([]byte(msgJSON), &msg)
    fmt.Printf("Message type: %s, data: %s\n", msg.Type, msg.Data)
}
```

#### 反例说明

```go
// ❌ 错误：字段未导出
type Bad struct {
    name string  // 小写未导出，无法编码
}

// ❌ 错误：循环引用
type Node struct {
    Value int
    Next  *Node  // 循环引用会导致无限递归
}

// ❌ 错误：忽略Unmarshal错误
var p Person
json.Unmarshal(data, &p)  // 忽略错误
// 应该
if err := json.Unmarshal(data, &p); err != nil {
    log.Fatal(err)
}

// ❌ 错误：数字精度丢失
var data map[string]interface{}
json.Unmarshal([]byte(`{"num": 9007199254740993}`), &data)
// data["num"] 变成 float64，精度丢失
```

#### 最佳实践

1. **字段导出**（首字母大写）
2. **omitempty**省略空值
3. **自定义Marshal**处理特殊类型
4. **流式处理**大JSON
5. **检查错误**不要忽略

---

### 9.4 context包

#### 概念定义

context用于传递截止日期、取消信号和请求范围值。是Go并发编程的核心工具。

#### 核心接口和函数

```go
// Context接口
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}

// 创建函数
func Background() Context                    // 根context
func TODO() Context                          // 占位context
func WithCancel(parent Context) (Context, CancelFunc)
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key, val interface{}) Context
```

#### 完整示例

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

// ========== 取消信号 ==========
func doWork(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: cancelled (%v)\n", id, ctx.Err())
            return
        default:
            fmt.Printf("Worker %d: working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func cancelExample() {
    ctx, cancel := context.WithCancel(context.Background())

    // 启动多个worker
    for i := 1; i <= 3; i++ {
        go doWork(ctx, i)
    }

    time.Sleep(2 * time.Second)
    fmt.Println("Main: cancelling...")
    cancel()  // 取消所有worker

    time.Sleep(500 * time.Millisecond)
}

// ========== 超时控制 ==========
func slowOperation(ctx context.Context) (string, error) {
    select {
    case <-time.After(3 * time.Second):
        return "completed", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func timeoutExample() {
    // 设置2秒超时
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    result, err := slowOperation(ctx)
    if err != nil {
        fmt.Printf("Operation failed: %v\n", err)
    } else {
        fmt.Printf("Result: %s\n", result)
    }
}

// ========== 截止时间 ==========
func deadlineExample() {
    deadline := time.Now().Add(2 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()

    if d, ok := ctx.Deadline(); ok {
        fmt.Printf("Deadline: %v\n", d)
    }

    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Completed")
    case <-ctx.Done():
        fmt.Printf("Deadline exceeded: %v\n", ctx.Err())
    }
}

// ========== 传递值 ==========
type key string

const userKey key = "user"

func withValueExample() {
    ctx := context.WithValue(context.Background(), userKey, "Alice")

    // 传递值
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    if user, ok := ctx.Value(userKey).(string); ok {
        fmt.Printf("Processing request for user: %s\n", user)
    }
}

// ========== HTTP服务示例 ==========
func httpHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 检查取消
    select {
    case <-ctx.Done():
        http.Error(w, "Request cancelled", http.StatusRequestTimeout)
        return
    default:
    }

    // 模拟处理
    select {
    case <-time.After(2 * time.Second):
        fmt.Fprintln(w, "Request processed")
    case <-ctx.Done():
        http.Error(w, "Request timeout", http.StatusRequestTimeout)
    }
}

// ========== 链式context ==========
func chainExample() {
    // 根context
    ctx := context.Background()

    // 添加超时
    ctx, cancel1 := context.WithTimeout(ctx, 10*time.Second)
    defer cancel1()

    // 添加值
    ctx = context.WithValue(ctx, userKey, "Bob")

    // 添加更短的超时（子context）
    ctx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
    defer cancel2()

    // 使用
    processWithContext(ctx)
}

func processWithContext(ctx context.Context) {
    fmt.Printf("User: %v\n", ctx.Value(userKey))
    if deadline, ok := ctx.Deadline(); ok {
        fmt.Printf("Deadline: %v\n", deadline)
    }
}

// ========== 最佳实践示例 ==========
func bestPractice(ctx context.Context) error {
    // 1. 总是将context作为第一个参数
    // 2. 不要存储context在结构体中
    // 3. 传递而不是存储

    // 检查取消
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // 创建子context用于特定操作
    subCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    // 使用子context
    return doSubOperation(subCtx)
}

func doSubOperation(ctx context.Context) error {
    select {
    case <-time.After(1 * time.Second):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func main() {
    fmt.Println("=== 取消示例 ===")
    cancelExample()

    fmt.Println("\n=== 超时示例 ===")
    timeoutExample()

    fmt.Println("\n=== 截止时间示例 ===")
    deadlineExample()

    fmt.Println("\n=== 值传递示例 ===")
    withValueExample()

    fmt.Println("\n=== 链式示例 ===")
    chainExample()

    fmt.Println("\n=== HTTP服务 ===")
    http.HandleFunc("/", httpHandler)
    fmt.Println("启动HTTP服务 :8080 (按Ctrl+C停止)")
    // go http.ListenAndServe(":8080", nil)
    time.Sleep(100 * time.Millisecond)
}
```

#### 反例说明

```go
// ❌ 错误：存储context在结构体中
type BadService struct {
    ctx context.Context  // 不要这样做
}

// ❌ 错误：context作为非第一个参数
func badFunc(data string, ctx context.Context) {  // 应该是第一个
    // ...
}

// ❌ 错误：传递nil context
result := doSomething(nil)  // 应该传递context.Background()

// ❌ 错误：忽略Done检查
func badOperation(ctx context.Context) {
    time.Sleep(10 * time.Second)  // 不检查ctx.Done()
}

// ❌ 错误：使用string作为key
data := context.WithValue(ctx, "user", "Alice")  // 容易冲突
// 应该使用自定义类型
type key string
const userKey key = "user"
data := context.WithValue(ctx, userKey, "Alice")
```

#### 最佳实践

1. **context作为第一个参数**
2. **不要存储**在结构体中
3. **传递而不是存储**
4. **使用自定义类型**作为key
5. **及时取消**释放资源

---

### 9.5 sync包基础

#### 概念定义

sync包提供基本的同步原语：互斥锁、读写锁、条件变量、WaitGroup、Once、Map、Pool。

#### 核心类型

```go
// 互斥锁
type Mutex struct { ... }
func (m *Mutex) Lock()
func (m *Mutex) Unlock()

// 读写锁
type RWMutex struct { ... }
func (rw *RWMutex) RLock()
func (rw *RWMutex) RUnlock()
func (rw *RWMutex) Lock()
func (rw *RWMutex) Unlock()

// WaitGroup
type WaitGroup struct { ... }
func (wg *WaitGroup) Add(delta int)
func (wg *WaitGroup) Done()
func (wg *WaitGroup) Wait()

// Once
type Once struct { ... }
func (o *Once) Do(f func())

// Map（并发安全）
type Map struct { ... }
func (m *Map) Load(key interface{}) (value interface{}, ok bool)
func (m *Map) Store(key, value interface{})
func (m *Map) Delete(key interface{})

// Pool（对象池）
type Pool struct { ... }
func (p *Pool) Get() interface{}
func (p *Pool) Put(x interface{})
```

#### 完整示例

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// ========== Mutex ==========
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *Counter) Get() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}

// ========== RWMutex ==========
type Cache struct {
    mu    sync.RWMutex
    data  map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

// ========== WaitGroup ==========
func waitGroupExample() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Duration(id) * 100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }

    wg.Wait()
    fmt.Println("All workers done")
}

// ========== Once ==========
func onceExample() {
    var once sync.Once
    var counter int

    for i := 0; i < 10; i++ {
        once.Do(func() {
            counter++
            fmt.Println("Initialized")
        })
    }

    fmt.Printf("Counter: %d (should be 1)\n", counter)
}

// ========== Map ==========
func mapExample() {
    var m sync.Map

    // 存储
    m.Store("key1", "value1")
    m.Store("key2", "value2")
    m.Store(123, 456)  // key可以是任意类型

    // 读取
    if val, ok := m.Load("key1"); ok {
        fmt.Printf("key1: %v\n", val)
    }

    // 删除
    m.Delete("key2")

    // 遍历
    fmt.Println("Map contents:")
    m.Range(func(key, value interface{}) bool {
        fmt.Printf("  %v: %v\n", key, value)
        return true  // 继续遍历
    })

    // LoadOrStore
    if val, loaded := m.LoadOrStore("key1", "new"); loaded {
        fmt.Printf("key1 already exists: %v\n", val)
    }
}

// ========== Pool ==========
func poolExample() {
    var pool = sync.Pool{
        New: func() interface{} {
            fmt.Println("Creating new object")
            return make([]byte, 1024)
        },
    }

    // 获取对象
    obj1 := pool.Get().([]byte)
    fmt.Printf("Got object, len=%d\n", len(obj1))

    // 使用对象
    copy(obj1, []byte("hello"))

    // 归还对象
    pool.Put(obj1)

    // 再次获取（可能复用）
    obj2 := pool.Get().([]byte)
    fmt.Printf("Got object again, len=%d, content=%s\n",
        len(obj2), string(obj2[:5]))

    pool.Put(obj2)
}

// ========== 原子操作 ==========
func atomicExample() {
    var counter int64

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                atomic.AddInt64(&counter, 1)
            }
        }()
    }

    wg.Wait()
    fmt.Printf("Counter: %d (should be 10000)\n", counter)

    // 原子读写
    var value int64 = 42
    loaded := atomic.LoadInt64(&value)
    fmt.Printf("Loaded: %d\n", loaded)

    atomic.StoreInt64(&value, 100)
    fmt.Printf("After store: %d\n", atomic.LoadInt64(&value))

    // CAS操作
    swapped := atomic.CompareAndSwapInt64(&value, 100, 200)
    fmt.Printf("Swapped: %v, Value: %d\n", swapped, value)
}

// ========== 单例模式 ==========
type singleton struct {
    data string
}

var (
    instance *singleton
    once     sync.Once
)

func GetInstance() *singleton {
    once.Do(func() {
        instance = &singleton{data: "initialized"}
        fmt.Println("Singleton initialized")
    })
    return instance
}

func main() {
    fmt.Println("=== Mutex ===")
    counter := &Counter{}
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Inc()
        }()
    }
    wg.Wait()
    fmt.Printf("Counter: %d\n", counter.Get())

    fmt.Println("\n=== RWMutex ===")
    cache := &Cache{data: make(map[string]string)}
    cache.Set("name", "Alice")
    if val, ok := cache.Get("name"); ok {
        fmt.Printf("Cached name: %s\n", val)
    }

    fmt.Println("\n=== WaitGroup ===")
    waitGroupExample()

    fmt.Println("\n=== Once ===")
    onceExample()

    fmt.Println("\n=== Map ===")
    mapExample()

    fmt.Println("\n=== Pool ===")
    poolExample()

    fmt.Println("\n=== Atomic ===")
    atomicExample()

    fmt.Println("\n=== Singleton ===")
    s1 := GetInstance()
    s2 := GetInstance()
    fmt.Printf("Same instance: %v\n", s1 == s2)
}
```

#### 反例说明

```go
// ❌ 错误：复制Mutex
var mu sync.Mutex
func bad() sync.Mutex {  // 返回Mutex副本
    return mu
}

// ❌ 错误：忘记Unlock
counter.mu.Lock()
counter.count++  // 如果panic，不会Unlock
// 应该用 defer

// ❌ 错误：WaitGroup计数错误
var wg sync.WaitGroup
wg.Add(1)
go func() {
    wg.Done()
    wg.Done()  // 错误：多次Done
}()

// ❌ 错误：WaitGroup在Add前Wait
var wg sync.WaitGroup
go func() {
    wg.Wait()  // 立即返回，因为计数为0
    wg.Add(1)  // 太晚
    // ...
}()

// ❌ 错误：普通map并发访问
var m = make(map[string]int)
go func() { m["a"] = 1 }()
go func() { m["b"] = 2 }()  // 竞态条件，可能panic

// ❌ 错误：sync.Map类型不安全
var m sync.Map
m.Store("key", "value")
val := m.Load("key").(int)  // panic: interface conversion
```

#### 最佳实践

1. **defer Unlock**确保释放
2. **RWMutex**读多写少场景
3. **WaitGroup**等待goroutine完成
4. **Once**确保只执行一次
5. **原子操作**简单计数场景

---

## 附录：Go 1.22+ 新特性

### 整数范围循环（Go 1.22）

```go
// 新语法：整数范围
for i := range 10 {
    fmt.Println(i)  // 0-9
}

// 带索引
for i := range N {
    // i从0到N-1
}
```

### 泛型增强（Go 1.18+）

```go
// 类型参数
func Min[T comparable](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 类型约束
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64
}
```

### 标准库新增

```go
// slices包（Go 1.21+）
import "slices"
slices.Sort(arr)
slices.BinarySearch(arr, target)
slices.Contains(arr, item)

// maps包（Go 1.21+）
import "maps"
maps.Copy(dst, src)
maps.Equal(m1, m2)
maps.Keys(m)
maps.Values(m)
```

---

## 10. Go 1.23 新特性

Go 1.23于2024年8月发布，带来了多项重要的语言特性和标准库更新。

### 10.1 迭代器（range-over-func）

#### 概念定义

Go 1.23将`range-over-func`实验性特性正式纳入语言规范，允许在`for-range`循环中使用迭代器函数作为range表达式。

#### 语法形式

```go
// 三种支持的迭代器函数类型
func(yield func() bool)           // 无返回值
func(yield func(K) bool)          // 单值返回
func(yield func(K, V) bool)       // 键值对返回
```

#### 完整示例

```go
package main

import "fmt"

// 自定义切片迭代器
func SliceIter[T any](s []T) func(yield func(int, T) bool) {
    return func(yield func(int, T) bool) {
        for i, v := range s {
            if !yield(i, v) {
                return
            }
        }
    }
}

// 自定义树形结构迭代器
type Tree[T any] struct {
    Value T
    Left  *Tree[T]
    Right *Tree[T]
}

func (t *Tree[T]) InOrder(yield func(T) bool) {
    if t == nil {
        return
    }
    t.Left.InOrder(yield)
    if !yield(t.Value) {
        return
    }
    t.Right.InOrder(yield)
}

func main() {
    // 使用自定义迭代器
    nums := []int{10, 20, 30, 40, 50}

    // 传统方式
    for i, v := range SliceIter(nums) {
        fmt.Printf("%d: %d\n", i, v)
    }

    // 树的中序遍历
    root := &Tree[int]{Value: 10,
        Left:  &Tree[int]{Value: 5},
        Right: &Tree[int]{Value: 15},
    }

    for v := range root.InOrder {
        fmt.Println(v)
    }
}
```

#### 反例说明

```go
// 错误：yield返回false时未正确终止
func BadIter() func(yield func(int) bool) {
    return func(yield func(int) bool) {
        for i := 0; i < 10; i++ {
            yield(i) // 忽略返回值，可能导致不必要的计算
        }
    }
}

// 正确：检查yield返回值
func GoodIter() func(yield func(int) bool) {
    return func(yield func(int) bool) {
        for i := 0; i < 10; i++ {
            if !yield(i) {
                return // 及时终止
            }
        }
    }
}
```

#### 最佳实践

1. **及时终止**：始终检查`yield`的返回值，当为`false`时立即返回
2. **清理资源**：在迭代器中分配的资源需要确保在终止时释放
3. **避免副作用**：迭代器函数应该是纯函数，避免修改外部状态

---

### 10.2 iter包

#### 概念定义

`iter`包提供了用于定义和使用迭代器的基本类型和函数，是Go 1.23引入的标准库新包。

#### 核心类型

```go
// Seq是单值迭代器类型
type Seq[V any] func(yield func(V) bool)

// Seq2是键值对迭代器类型
type Seq2[K, V any] func(yield func(K, V) bool)
```

#### 完整示例

```go
package main

import (
    "fmt"
    "iter"
)

// 生成斐波那契数列迭代器
func Fibonacci(n int) iter.Seq[int] {
    return func(yield func(int) bool) {
        a, b := 0, 1
        for i := 0; i < n; i++ {
            if !yield(a) {
                return
            }
            a, b = b, a+b
        }
    }
}

// 过滤迭代器
func Filter[V any](seq iter.Seq[V], pred func(V) bool) iter.Seq[V] {
    return func(yield func(V) bool) {
        for v := range seq {
            if pred(v) && !yield(v) {
                return
            }
        }
    }
}

// 映射迭代器
func Map[T, R any](seq iter.Seq[T], fn func(T) R) iter.Seq[R] {
    return func(yield func(R) bool) {
        for v := range seq {
            if !yield(fn(v)) {
                return
            }
        }
    }
}

func main() {
    // 使用斐波那契迭代器
    fmt.Println("斐波那契数列:")
    for v := range Fibonacci(10) {
        fmt.Printf("%d ", v)
    }
    fmt.Println()

    // 链式操作
    nums := func(yield func(int) bool) {
        for i := 1; i <= 10; i++ {
            if !yield(i) {
                return
            }
        }
    }

    // 过滤偶数并映射为平方
    evenSquares := Map(Filter(nums, func(n int) bool {
        return n%2 == 0
    }), func(n int) int {
        return n * n
    })

    fmt.Println("偶数的平方:")
    for v := range evenSquares {
        fmt.Printf("%d ", v)
    }
    fmt.Println()
}
```

#### 与slices和maps包集成

```go
package main

import (
    "fmt"
    "slices"
    "maps"
)

func main() {
    // slices包的迭代器函数
    nums := []int{3, 1, 4, 1, 5, 9, 2, 6}

    // All: 返回索引和值的迭代器
    for i, v := range slices.All(nums) {
        fmt.Printf("nums[%d] = %d\n", i, v)
    }

    // Values: 只返回值
    for v := range slices.Values(nums) {
        fmt.Println(v)
    }

    // Backward: 反向迭代
    for i, v := range slices.Backward(nums) {
        fmt.Printf("反向: nums[%d] = %d\n", i, v)
    }

    // Collect: 从迭代器收集到切片
    seq := func(yield func(int) bool) {
        for i := 0; i < 5; i++ {
            if !yield(i * 10) {
                return
            }
        }
    }
    collected := slices.Collect(seq)
    fmt.Println("收集结果:", collected)

    // maps包的迭代器函数
    m := map[string]int{"a": 1, "b": 2, "c": 3}

    // All: 返回键值对
    for k, v := range maps.All(m) {
        fmt.Printf("%s: %d\n", k, v)
    }

    // Keys: 只返回键
    for k := range maps.Keys(m) {
        fmt.Println("键:", k)
    }

    // Values: 只返回值
    for v := range maps.Values(m) {
        fmt.Println("值:", v)
    }
}
```

#### 最佳实践

1. **使用标准类型**：优先使用`iter.Seq`和`iter.Seq2`而不是自定义函数类型
2. **链式操作**：利用迭代器组合实现复杂的数据处理流程
3. **延迟求值**：迭代器是惰性求值的，避免不必要的计算

---

### 10.3 unique包

#### 概念定义

`unique`包提供了值规范化（canonicalization）功能，类似于"interning"或"hash-consing"，用于减少内存占用和提高比较效率。

#### 核心概念

```go
// Handle是指向规范化值的引用
type Handle[T comparable] struct {
    value *T
}

// Make创建规范化值
func Make[T comparable](value T) Handle[T]

// Value返回原始值
func (h Handle[T]) Value() T
```

#### 完整示例

```go
package main

import (
    "fmt"
    "unique"
)

func main() {
    // 创建大量重复的字符串
    str1 := "hello world"
    str2 := "hello world"
    str3 := "hello world"

    // 规范化
    h1 := unique.Make(str1)
    h2 := unique.Make(str2)
    h3 := unique.Make(str3)

    // 比较Handle（指针比较，O(1)）
    fmt.Println("h1 == h2:", h1 == h2) // true
    fmt.Println("h1 == h3:", h1 == h3) // true

    // 获取原始值
    fmt.Println("h1.Value():", h1.Value())

    // 不同值的Handle不相等
    other := unique.Make("different")
    fmt.Println("h1 == other:", h1 == other) // false

    // 实际应用场景：大量重复数据的内存优化
    type User struct {
        ID   int
        Name unique.Handle[string] // 使用Handle减少内存
    }

    users := make([]User, 1000)
    for i := range users {
        users[i] = User{
            ID:   i,
            Name: unique.Make("Common Name"), // 共享同一份内存
        }
    }

    fmt.Printf("创建了%d个用户\n", len(users))
}
```

#### 反例说明

```go
// 错误：对不经常重复的值使用unique
func BadUsage() {
    // 每个值都是唯一的，使用unique反而增加开销
    for i := 0; i < 1000000; i++ {
        _ = unique.Make(fmt.Sprintf("unique-id-%d", i))
    }
}

// 正确：对大量重复的值使用unique
func GoodUsage() {
    statusOK := unique.Make("OK")
    statusError := unique.Make("ERROR")

    for i := 0; i < 1000000; i++ {
        if i%2 == 0 {
            _ = statusOK
        } else {
            _ = statusError
        }
    }
}
```

#### 最佳实践

1. **适用场景**：大量重复的小值（字符串、小结构体）
2. **比较操作**：使用Handle比较代替值比较，获得O(1)性能
3. **内存权衡**：对于不重复的值，使用unique会增加额外开销

---

### 10.4 structs包

#### 概念定义

`structs`包提供了用于修改结构体属性的标记类型（marker types），目前主要用于指示结构体布局符合主机平台期望。

#### HostLayout类型

```go
// HostLayout表示结构体具有符合主机平台期望的布局
type HostLayout struct{}
```

#### 完整示例

```go
package main

import (
    "fmt"
    "structs"
)

// 用于与C库交互的结构体
type CCompatibleStruct struct {
    _ structs.HostLayout // 标记：使用主机平台布局

    X int32
    Y int32
    Z float64
}

// 普通Go结构体（可能有不同的内存布局）
type NormalStruct struct {
    X int32
    Y int32
    Z float64
}

func main() {
    // 在需要与主机API交互时使用HostLayout
    // 例如：传递给C函数、系统调用等

    cs := CCompatibleStruct{X: 1, Y: 2, Z: 3.14}
    ns := NormalStruct{X: 1, Y: 2, Z: 3.14}

    fmt.Printf("CCompatibleStruct大小: %d\n", unsafe.Sizeof(cs))
    fmt.Printf("NormalStruct大小: %d\n", unsafe.Sizeof(ns))

    // 注意：HostLayout不保证特定的布局，只表示"使用主机平台布局"
    // 具体布局仍然依赖于编译器和平台
}
```

#### 使用场景

1. **Cgo交互**：与C库传递结构体时确保布局兼容
2. **系统调用**：需要特定内存布局的系统调用
3. **硬件接口**：与硬件设备通信时的数据结构

---

### 10.5 Timer和Ticker改进

#### 概念定义

Go 1.23对`time.Timer`和`time.Ticker`的实现进行了两项重要改进，解决了历史遗留问题。

#### 改进内容

1. **立即垃圾回收**：不再被引用的Timer和Ticker可以立即被GC回收，无需调用Stop()
2. **无缓冲Channel**：Timer和Ticker的channel现在是无缓冲的（容量0），保证Reset/Stop后不会收到旧值

#### 完整示例

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 改进1：Timer可以被GC回收
    func() {
        t := time.NewTimer(time.Hour)
        _ = t
        // 函数返回后，t不再被引用，可以被GC回收
        // 在Go 1.23之前，这个Timer会一直存在直到触发
    }()

    // 改进2：无缓冲channel保证不收到旧值
    timer := time.NewTimer(100 * time.Millisecond)

    // 等待一段时间
    time.Sleep(200 * time.Millisecond)

    // 重置timer
    timer.Reset(100 * time.Millisecond)

    // 在Go 1.23中，保证不会收到重置前准备好的旧值
    // 因为channel是无缓冲的，旧值已经被丢弃
    <-timer.C
    fmt.Println("收到新值")

    // 验证channel容量
    fmt.Printf("Timer channel容量: %d\n", cap(timer.C)) // 0
    fmt.Printf("Timer channel长度: %d\n", len(timer.C)) // 0

    // 使用Ticker
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    count := 0
    for t := range ticker.C {
        fmt.Println("Tick:", t)
        count++
        if count >= 3 {
            break
        }
    }
}
```

#### 迁移注意事项

```go
// 旧代码：依赖channel容量的代码需要修改
func OldStyle() {
    timer := time.NewTimer(time.Second)

    // Go 1.22及之前：cap(timer.C) == 1
    // Go 1.23：cap(timer.C) == 0

    // 错误：依赖channel长度判断是否可以接收
    if len(timer.C) > 0 {
        <-timer.C
    }

    // 正确：使用非阻塞接收
    select {
    case <-timer.C:
        // 有值可接收
    default:
        // 无值
    }
}
```

#### 最佳实践

1. **不依赖channel容量**：代码不应依赖Timer/Ticker channel的容量
2. **使用select非阻塞接收**：判断是否有值可接收时，使用select语句
3. **及时Stop**：虽然可以被GC，但显式Stop仍然是好习惯

---

### 10.6 泛型类型别名

#### 概念定义

Go 1.23包含对泛型类型别名的预览支持，允许类型别名带有类型参数。

#### 启用方式

```bash
GOEXPERIMENT=aliastypeparams go run main.go
```

#### 完整示例

```go
package main

import "fmt"

// 泛型类型别名（需要GOEXPERIMENT=aliastypeparams）
type MySlice[T any] = []T

type StringMap[V any] = map[string]V

type Pair[T, U any] = struct {
    First  T
    Second U
}

// 使用泛型别名
func ProcessSlice[T any](s MySlice[T]) {
    for _, v := range s {
        fmt.Println(v)
    }
}

func main() {
    // 使用泛型类型别名
    nums := MySlice[int]{1, 2, 3, 4, 5}
    ProcessSlice(nums)

    scores := StringMap[int]{
        "Alice": 95,
        "Bob":   87,
    }
    fmt.Println(scores)

    pair := Pair[string, int]{"age", 25}
    fmt.Println(pair)
}
```

#### 限制说明

1. **实验特性**：需要设置`GOEXPERIMENT=aliastypeparams`
2. **包内使用**：Go 1.23中仅支持在定义包内使用泛型别名
3. **跨包支持**：预计Go 1.24将支持跨包使用

#### 最佳实践

1. **谨慎使用**：目前仍是实验特性，生产环境需谨慎
2. **关注演进**：密切关注Go 1.24的正式发布
3. **向后兼容**：确保代码在特性不可用时仍有备选方案

---

> 文档完成。本文档涵盖了**Go 1.23**的所有核心语言特性及最新更新，每个特性都包含了概念定义、语法形式、属性关系、形式论证、完整示例、反例说明和最佳实践。
