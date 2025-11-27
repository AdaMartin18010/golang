# 随机数工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [随机数工具](#随机数工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 随机数生成](#21-随机数生成)
    - [2.2 随机字符串生成](#22-随机字符串生成)
    - [2.3 随机字节数组生成](#23-随机字节数组生成)
    - [2.4 加密安全随机数](#24-加密安全随机数)
    - [2.5 随机选择](#25-随机选择)
    - [2.6 随机操作](#26-随机操作)
  - [3. 使用示例](#3-使用示例)
    - [3.1 随机数生成](#31-随机数生成)
    - [3.2 随机字符串生成](#32-随机字符串生成)
    - [3.3 随机字节数组生成](#33-随机字节数组生成)
    - [3.4 加密安全随机数](#34-加密安全随机数)
    - [3.5 随机选择](#35-随机选择)
    - [3.6 随机操作](#36-随机操作)

---

## 1. 概述

随机数工具提供了丰富的随机数生成函数，支持整数、浮点数、字符串、字节数组等类型的随机生成，以及随机选择、打乱、加权选择等功能。

---

## 2. 功能特性

### 2.1 随机数生成

- `Int`: 生成随机整数 [0, max)
- `IntRange`: 生成指定范围内的随机整数 [min, max)
- `Int64`: 生成随机64位整数 [0, max)
- `Int64Range`: 生成指定范围内的随机64位整数 [min, max)
- `Float64`: 生成随机浮点数 [0.0, 1.0)
- `Float64Range`: 生成指定范围内的随机浮点数 [min, max)
- `Bool`: 生成随机布尔值
- `Probability`: 根据概率返回true

### 2.2 随机字符串生成

- `String`: 生成指定长度的随机字符串（字母+数字）
- `StringWithCharset`: 使用指定字符集生成随机字符串
- `LettersString`: 生成随机字母字符串
- `DigitsString`: 生成随机数字字符串
- `Hex`: 生成随机十六进制字符串
- `HexUpper`: 生成随机十六进制字符串（大写）
- `UUID`: 生成随机UUID字符串（简化版）
- `UUIDWithDashes`: 生成带连字符的UUID字符串（简化版）
- `FastString`: 快速生成随机字符串（使用unsafe，性能更高）

### 2.3 随机字节数组生成

- `Bytes`: 生成指定长度的随机字节数组
- `SecureBytes`: 使用加密安全的随机数生成器生成随机字节数组

### 2.4 加密安全随机数

- `SecureInt`: 使用加密安全的随机数生成器生成整数
- `SecureIntRange`: 使用加密安全的随机数生成器生成指定范围内的整数
- `SecureString`: 使用加密安全的随机数生成器生成随机字符串
- `SecureStringWithCharset`: 使用加密安全的随机数生成器和指定字符集生成随机字符串

### 2.5 随机选择

- `Choice`: 从切片中随机选择一个元素
- `Choices`: 从切片中随机选择n个元素（允许重复）
- `Sample`: 从切片中随机选择n个不重复的元素
- `WeightedChoice`: 根据权重随机选择元素

### 2.6 随机操作

- `Shuffle`: 随机打乱切片（Fisher-Yates算法）
- `Duration`: 生成指定范围内的随机时间间隔
- `Time`: 生成指定时间范围内的随机时间
- `Seed`: 设置随机数生成器的种子

---

## 3. 使用示例

### 3.1 随机数生成

```go
import "github.com/yourusername/golang/pkg/utils/random"

// 生成随机整数
val := random.Int(100) // [0, 100)

// 生成指定范围内的随机整数
val := random.IntRange(10, 20) // [10, 20)

// 生成随机浮点数
val := random.Float64() // [0.0, 1.0)

// 生成指定范围内的随机浮点数
val := random.Float64Range(1.0, 2.0) // [1.0, 2.0)

// 生成随机布尔值
b := random.Bool()

// 根据概率返回true
b := random.Probability(0.7) // 70%概率返回true
```

### 3.2 随机字符串生成

```go
// 生成随机字符串（字母+数字）
str := random.String(10)

// 使用指定字符集生成随机字符串
str := random.StringWithCharset(10, "abc123")

// 生成随机字母字符串
str := random.LettersString(10)

// 生成随机数字字符串
str := random.DigitsString(10)

// 生成随机十六进制字符串
str := random.Hex(16)

// 生成UUID字符串
uuid := random.UUID()
uuidWithDashes := random.UUIDWithDashes()
```

### 3.3 随机字节数组生成

```go
// 生成随机字节数组
bytes := random.Bytes(16)

// 使用加密安全的随机数生成器
secureBytes, err := random.SecureBytes(16)
```

### 3.4 加密安全随机数

```go
// 生成加密安全的随机整数
val, err := random.SecureInt(100)

// 生成加密安全的随机字符串
str, err := random.SecureString(10)

// 使用指定字符集生成加密安全的随机字符串
str, err := random.SecureStringWithCharset(10, "abc123")
```

### 3.5 随机选择

```go
// 从切片中随机选择一个元素
slice := []int{1, 2, 3, 4, 5}
val, ok := random.Choice(slice)

// 从切片中随机选择n个元素（允许重复）
result := random.Choices(slice, 3)

// 从切片中随机选择n个不重复的元素
result := random.Sample(slice, 3)

// 根据权重随机选择元素
items := []string{"a", "b", "c"}
weights := []float64{0.1, 0.2, 0.7}
val, ok := random.WeightedChoice(items, weights)
```

### 3.6 随机操作

```go
// 随机打乱切片
slice := []int{1, 2, 3, 4, 5}
random.Shuffle(slice)

// 生成随机时间间隔
duration := random.Duration(time.Second, 10*time.Second)

// 生成随机时间
start := time.Now()
end := start.Add(24 * time.Hour)
tm := random.Time(start, end)

// 设置随机数生成器的种子
random.Seed(12345)
```

---

**更新日期**: 2025-11-11
