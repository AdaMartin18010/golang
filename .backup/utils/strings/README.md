# 字符串工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [字符串工具](#字符串工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

字符串工具提供了丰富的字符串操作函数，简化常见的字符串处理任务。

---

## 2. 功能特性

### 2.1 基础操作

- `IsEmpty`: 检查字符串是否为空
- `IsNotEmpty`: 检查字符串是否非空
- `Truncate`: 截断字符串
- `Reverse`: 反转字符串

### 2.2 包含检查

- `ContainsAny`: 检查是否包含任意一个子字符串
- `ContainsAll`: 检查是否包含所有子字符串

### 2.3 格式化

- `PadLeft`: 左侧填充
- `PadRight`: 右侧填充
- `PadCenter`: 居中填充
- `RemoveWhitespace`: 移除空白字符

### 2.4 命名转换

- `CamelToSnake`: 驼峰转蛇形
- `SnakeToCamel`: 蛇形转驼峰
- `FirstUpper`: 首字母大写
- `FirstLower`: 首字母小写

### 2.5 掩码处理

- `Mask`: 通用掩码函数
- `MaskEmail`: 掩码邮箱
- `MaskPhone`: 掩码手机号

### 2.6 随机生成

- `RandomString`: 生成随机字符串
- `RandomStringWithCharset`: 使用指定字符集生成随机字符串

---

## 3. 使用示例

### 3.1 基础操作

```go
import "github.com/yourusername/golang/pkg/utils/strings"

// 检查是否为空
if strings.IsEmpty(str) {
    // 处理空字符串
}

// 截断字符串
truncated := strings.Truncate("hello world", 5) // "he..."

// 反转字符串
reversed := strings.Reverse("hello") // "olleh"
```

### 3.2 格式化

```go
// 左侧填充
padded := strings.PadLeft("123", 5, '0') // "00123"

// 右侧填充
padded := strings.PadRight("123", 5, '0') // "12300"

// 居中填充
padded := strings.PadCenter("123", 7, '0') // "0012300"
```

### 3.3 命名转换

```go
// 驼峰转蛇形
snake := strings.CamelToSnake("HelloWorld") // "hello_world"

// 蛇形转驼峰
camel := strings.SnakeToCamel("hello_world") // "helloWorld"

// 首字母大写
upper := strings.FirstUpper("hello") // "Hello"
```

### 3.4 掩码处理

```go
// 掩码邮箱
masked := strings.MaskEmail("test@example.com") // "t***t@example.com"

// 掩码手机号
masked := strings.MaskPhone("13812345678") // "138****5678"

// 通用掩码
masked := strings.Mask("1234567890", 3, 7, '*') // "123***7890"
```

### 3.5 随机生成

```go
// 生成随机字符串
random, err := strings.RandomString(16)

// 使用指定字符集
random, err := strings.RandomStringWithCharset(16, "abcdefghijklmnopqrstuvwxyz")
```

---

**更新日期**: 2025-11-11
