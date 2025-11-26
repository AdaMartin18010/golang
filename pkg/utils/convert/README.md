# 类型转换工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [类型转换工具](#类型转换工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

类型转换工具提供了各种类型之间的便捷转换功能，简化类型转换任务，支持基本类型、切片、映射等常用类型的转换。

---

## 2. 功能特性

### 2.1 基本类型转换

- `ToString`: 转换为字符串
- `ToInt`: 转换为int
- `ToInt64`: 转换为int64
- `ToFloat64`: 转换为float64
- `ToBool`: 转换为bool
- `ToBytes`: 转换为[]byte

### 2.2 强制转换（Must函数）

- `MustInt`: 转换为int，失败时panic
- `MustInt64`: 转换为int64，失败时panic
- `MustFloat64`: 转换为float64，失败时panic
- `MustBool`: 转换为bool，失败时panic

### 2.3 默认值转换

- `ToIntDefault`: 转换为int，失败时返回默认值
- `ToInt64Default`: 转换为int64，失败时返回默认值
- `ToFloat64Default`: 转换为float64，失败时返回默认值
- `ToBoolDefault`: 转换为bool，失败时返回默认值

### 2.4 切片转换

- `ToStringSlice`: 转换为[]string
- `ToIntSlice`: 转换为[]int
- `ToInt64Slice`: 转换为[]int64
- `ToFloat64Slice`: 转换为[]float64
- `ToBoolSlice`: 转换为[]bool

### 2.5 映射转换

- `ToMapStringInterface`: 转换为map[string]interface{}
- `ToMapStringString`: 转换为map[string]string

### 2.6 类型检查

- `IsNumeric`: 检查是否为数字类型
- `IsInteger`: 检查是否为整数类型
- `IsFloat`: 检查是否为浮点数类型

---

## 3. 使用示例

### 3.1 基本类型转换

```go
import "github.com/yourusername/golang/pkg/utils/convert"

// 转换为字符串
str := convert.ToString(42)        // "42"
str = convert.ToString(3.14)       // "3.14"
str = convert.ToString(true)       // "true"

// 转换为int
num, err := convert.ToInt("42")    // 42, nil
num, err = convert.ToInt(42.5)     // 42, nil

// 转换为int64
num64, err := convert.ToInt64("42") // 42, nil

// 转换为float64
f, err := convert.ToFloat64("3.14") // 3.14, nil

// 转换为bool
b, err := convert.ToBool("true")    // true, nil
b, err = convert.ToBool(1)          // true, nil
b, err = convert.ToBool(0)          // false, nil

// 转换为[]byte
bytes := convert.ToBytes("hello")   // []byte("hello")
```

### 3.2 强制转换（Must函数）

```go
// 转换为int，失败时panic
num := convert.MustInt("42")        // 42

// 转换为int64，失败时panic
num64 := convert.MustInt64("42")    // 42

// 转换为float64，失败时panic
f := convert.MustFloat64("3.14")    // 3.14

// 转换为bool，失败时panic
b := convert.MustBool("true")       // true
```

### 3.3 默认值转换

```go
// 转换为int，失败时返回默认值
num := convert.ToIntDefault("42", 0)        // 42
num = convert.ToIntDefault("invalid", 100)  // 100

// 转换为int64，失败时返回默认值
num64 := convert.ToInt64Default("42", 0)    // 42
num64 = convert.ToInt64Default("invalid", 100) // 100

// 转换为float64，失败时返回默认值
f := convert.ToFloat64Default("3.14", 0.0)  // 3.14
f = convert.ToFloat64Default("invalid", 0.0) // 0.0

// 转换为bool，失败时返回默认值
b := convert.ToBoolDefault("true", false)   // true
b = convert.ToBoolDefault("invalid", false) // false
```

### 3.4 切片转换

```go
// 转换为[]string
strSlice := convert.ToStringSlice([]int{1, 2, 3})  // ["1", "2", "3"]

// 转换为[]int
intSlice, err := convert.ToIntSlice([]string{"1", "2", "3"})  // [1, 2, 3], nil

// 转换为[]int64
int64Slice, err := convert.ToInt64Slice([]string{"1", "2", "3"})  // [1, 2, 3], nil

// 转换为[]float64
floatSlice, err := convert.ToFloat64Slice([]string{"1.1", "2.2", "3.3"})  // [1.1, 2.2, 3.3], nil

// 转换为[]bool
boolSlice, err := convert.ToBoolSlice([]string{"true", "false", "true"})  // [true, false, true], nil
```

### 3.5 映射转换

```go
// 转换为map[string]interface{}
m1 := map[string]interface{}{
    "key1": "value1",
    "key2": 42,
}
result1, err := convert.ToMapStringInterface(m1)  // map[string]interface{}, nil

// 转换为map[string]string
m2 := map[string]interface{}{
    "key1": "value1",
    "key2": 42,
}
result2, err := convert.ToMapStringString(m2)  // map[string]string{"key1": "value1", "key2": "42"}, nil
```

### 3.6 类型检查

```go
// 检查是否为数字类型
isNum := convert.IsNumeric(42)      // true
isNum = convert.IsNumeric("42")     // true
isNum = convert.IsNumeric("hello")  // false

// 检查是否为整数类型
isInt := convert.IsInteger(42)      // true
isInt = convert.IsInteger("42")     // true
isInt = convert.IsInteger("42.5")   // false

// 检查是否为浮点数类型
isFloat := convert.IsFloat(3.14)    // true
isFloat = convert.IsFloat("3.14")   // true
isFloat = convert.IsFloat("42")     // false
```

### 3.7 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/convert"
)

func main() {
    // 基本类型转换
    str := convert.ToString(42)
    fmt.Printf("String: %s\n", str)

    num, err := convert.ToInt("42")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Int: %d\n", num)
    }

    // 使用默认值
    num = convert.ToIntDefault("invalid", 100)
    fmt.Printf("Int with default: %d\n", num)

    // 切片转换
    strSlice := convert.ToStringSlice([]int{1, 2, 3})
    fmt.Printf("String slice: %v\n", strSlice)

    intSlice, err := convert.ToIntSlice([]string{"1", "2", "3"})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Int slice: %v\n", intSlice)
    }

    // 类型检查
    if convert.IsNumeric("42") {
        fmt.Println("'42' is numeric")
    }

    if convert.IsInteger("42") {
        fmt.Println("'42' is integer")
    }
}
```

---

**更新日期**: 2025-11-11
