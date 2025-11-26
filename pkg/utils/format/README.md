# 格式化工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [格式化工具](#格式化工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 数字格式化](#21-数字格式化)
    - [2.2 时间格式化](#22-时间格式化)
    - [2.3 字节格式化](#23-字节格式化)
    - [2.4 证件格式化](#24-证件格式化)
    - [2.5 掩码格式化](#25-掩码格式化)
    - [2.6 字符串格式化](#26-字符串格式化)
  - [3. 使用示例](#3-使用示例)
    - [3.1 数字格式化](#31-数字格式化)
    - [3.2 时间格式化](#32-时间格式化)
    - [3.3 字节格式化](#33-字节格式化)
    - [3.4 证件格式化](#34-证件格式化)
    - [3.5 掩码格式化](#35-掩码格式化)
    - [3.6 字符串格式化](#36-字符串格式化)
    - [3.7 完整示例](#37-完整示例)

---

## 1. 概述

格式化工具提供了各种数据格式化功能，包括数字、时间、字节、电话号码、身份证号、银行卡号等的格式化，以及字符串处理功能。

---

## 2. 功能特性

### 2.1 数字格式化

- `FormatNumber`: 格式化数字（添加千分位分隔符）
- `FormatFloat`: 格式化浮点数（添加千分位分隔符）
- `FormatPercent`: 格式化百分比

### 2.2 时间格式化

- `FormatTime`: 格式化时间
- `FormatTimeRFC3339`: 格式化时间为RFC3339格式
- `FormatTimeISO8601`: 格式化时间为ISO8601格式
- `FormatTimeHuman`: 格式化时间为人类可读格式
- `FormatTimeRelative`: 格式化相对时间
- `FormatDuration`: 格式化持续时间（人类可读）
- `FormatDurationShort`: 格式化持续时间（简短格式）

### 2.3 字节格式化

- `FormatBytes`: 格式化字节数
- `FormatBytesShort`: 格式化字节数（简短格式）

### 2.4 证件格式化

- `FormatPhone`: 格式化电话号码
- `FormatIDCard`: 格式化身份证号
- `FormatBankCard`: 格式化银行卡号

### 2.5 掩码格式化

- `FormatMask`: 格式化掩码（隐藏部分信息）
- `FormatMaskPhone`: 格式化手机号（中间4位掩码）
- `FormatMaskEmail`: 格式化邮箱（用户名部分掩码）
- `FormatMaskIDCard`: 格式化身份证号（中间掩码）
- `FormatMaskBankCard`: 格式化银行卡号（中间掩码）

### 2.6 字符串格式化

- `FormatPlural`: 格式化复数形式
- `FormatList`: 格式化列表
- `FormatListWithAnd`: 格式化列表（最后一项用"和"连接）
- `FormatListWithOr`: 格式化列表（最后一项用"或"连接）
- `FormatTruncate`: 截断字符串
- `FormatPadLeft`: 左填充
- `FormatPadRight`: 右填充
- `FormatPadCenter`: 居中填充
- `FormatIndent`: 缩进
- `FormatWrap`: 换行

---

## 3. 使用示例

### 3.1 数字格式化

```go
import "github.com/yourusername/golang/pkg/utils/format"

// 格式化数字（添加千分位分隔符）
num := format.FormatNumber(1234567)  // "1,234,567"

// 格式化浮点数
f := format.FormatFloat(1234567.89, 2)  // "1,234,567.89"

// 格式化百分比
percent := format.FormatPercent(25, 100)  // "25.00%"
```

### 3.2 时间格式化

```go
// 格式化时间
t := time.Now()
formatted := format.FormatTime(t, "2006-01-02 15:04:05")

// 格式化时间为RFC3339格式
rfc3339 := format.FormatTimeRFC3339(t)  // "2006-01-02T15:04:05Z07:00"

// 格式化时间为ISO8601格式
iso8601 := format.FormatTimeISO8601(t)  // "2006-01-02T15:04:05Z07:00"

// 格式化时间为人类可读格式
human := format.FormatTimeHuman(t.Add(-5 * time.Minute))  // "5分钟前"

// 格式化相对时间
relative := format.FormatTimeRelative(t.Add(-2 * time.Hour))  // "2小时前"

// 格式化持续时间
duration := format.FormatDuration(2*time.Hour + 30*time.Minute)  // "2h 30m"

// 格式化持续时间（简短格式）
short := format.FormatDurationShort(2*time.Hour + 30*time.Minute)  // "2.5h"
```

### 3.3 字节格式化

```go
// 格式化字节数
bytes := format.FormatBytes(1024 * 1024)  // "1.00 MB"

// 格式化字节数（简短格式）
short := format.FormatBytesShort(1024 * 1024)  // "1.0MB"
```

### 3.4 证件格式化

```go
// 格式化电话号码
phone := format.FormatPhone("13800138000")  // "138 0013 8000"

// 格式化身份证号
idCard := format.FormatIDCard("123456199001011234")  // "123456 19900101 1234"

// 格式化银行卡号
bankCard := format.FormatBankCard("1234567890123456")  // "1234 5678 9012 3456"
```

### 3.5 掩码格式化

```go
// 格式化掩码（隐藏部分信息）
masked := format.FormatMask("1234567890", 3, 7, '*')  // "123***7890"

// 格式化手机号（中间4位掩码）
phone := format.FormatMaskPhone("13800138000")  // "138****8000"

// 格式化邮箱（用户名部分掩码）
email := format.FormatMaskEmail("user@example.com")  // "u***r@example.com"

// 格式化身份证号（中间掩码）
idCard := format.FormatMaskIDCard("123456199001011234")  // "123456********1234"

// 格式化银行卡号（中间掩码）
bankCard := format.FormatMaskBankCard("1234567890123456")  // "1234****3456"
```

### 3.6 字符串格式化

```go
// 格式化复数形式
plural := format.FormatPlural(1, "item", "items")  // "1 item"
plural = format.FormatPlural(2, "item", "items")   // "2 items"

// 格式化列表
list := format.FormatList([]string{"a", "b", "c"}, ", ")  // "a, b, c"

// 格式化列表（最后一项用"和"连接）
listAnd := format.FormatListWithAnd([]string{"apple", "banana", "orange"})  // "apple、banana和orange"

// 格式化列表（最后一项用"或"连接）
listOr := format.FormatListWithOr([]string{"apple", "banana", "orange"})  // "apple、banana或orange"

// 截断字符串
truncated := format.FormatTruncate("hello world", 8, "...")  // "hello..."

// 左填充
paddedLeft := format.FormatPadLeft("123", 5, '0')  // "00123"

// 右填充
paddedRight := format.FormatPadRight("123", 5, '0')  // "12300"

// 居中填充
paddedCenter := format.FormatPadCenter("123", 7, ' ')  // "  123  "

// 缩进
indented := format.FormatIndent("line1\nline2", "  ")  // "  line1\n  line2"

// 换行
wrapped := format.FormatWrap("hello world", 5)  // "hello\n worl\nd"
```

### 3.7 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/golang/pkg/utils/format"
)

func main() {
    // 数字格式化
    fmt.Printf("Number: %s\n", format.FormatNumber(1234567))
    fmt.Printf("Float: %s\n", format.FormatFloat(1234567.89, 2))
    fmt.Printf("Percent: %s\n", format.FormatPercent(25, 100))

    // 时间格式化
    now := time.Now()
    fmt.Printf("RFC3339: %s\n", format.FormatTimeRFC3339(now))
    fmt.Printf("Human: %s\n", format.FormatTimeHuman(now.Add(-5 * time.Minute)))
    fmt.Printf("Duration: %s\n", format.FormatDuration(2*time.Hour + 30*time.Minute))

    // 字节格式化
    fmt.Printf("Bytes: %s\n", format.FormatBytes(1024 * 1024))

    // 证件格式化
    fmt.Printf("Phone: %s\n", format.FormatPhone("13800138000"))
    fmt.Printf("ID Card: %s\n", format.FormatIDCard("123456199001011234"))

    // 掩码格式化
    fmt.Printf("Masked Phone: %s\n", format.FormatMaskPhone("13800138000"))
    fmt.Printf("Masked Email: %s\n", format.FormatMaskEmail("user@example.com"))

    // 字符串格式化
    fmt.Printf("Plural: %s\n", format.FormatPlural(2, "item", "items"))
    fmt.Printf("List: %s\n", format.FormatListWithAnd([]string{"apple", "banana", "orange"}))
    fmt.Printf("Truncate: %s\n", format.FormatTruncate("hello world", 8, "..."))
}
```

---

**更新日期**: 2025-11-11
