# 正则表达式工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [正则表达式工具](#正则表达式工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 匹配操作](#21-匹配操作)
    - [2.2 查找操作](#22-查找操作)
    - [2.3 替换操作](#23-替换操作)
    - [2.4 分割操作](#24-分割操作)
    - [2.5 统计操作](#25-统计操作)
    - [2.6 位置操作](#26-位置操作)
    - [2.7 工具函数](#27-工具函数)
    - [2.8 常用模式](#28-常用模式)
    - [2.9 快捷匹配函数](#29-快捷匹配函数)
  - [3. 使用示例](#3-使用示例)
    - [3.1 匹配操作](#31-匹配操作)
    - [3.2 查找操作](#32-查找操作)
    - [3.3 替换操作](#33-替换操作)
    - [3.4 分割操作](#34-分割操作)
    - [3.5 统计操作](#35-统计操作)
    - [3.6 位置操作](#36-位置操作)
    - [3.7 工具函数](#37-工具函数)
    - [3.8 常用模式](#38-常用模式)
    - [3.9 带选项的操作](#39-带选项的操作)

---

## 1. 概述

正则表达式工具提供了丰富的正则表达式操作函数，简化字符串匹配、查找、替换、分割等常见任务。

---

## 2. 功能特性

### 2.1 匹配操作

- `Match`: 检查字符串是否匹配正则表达式
- `MatchString`: 检查字符串是否匹配正则表达式（忽略错误）
- `HasMatch`: 检查是否有匹配
- `HasMatchString`: 检查是否有匹配（忽略错误）

### 2.2 查找操作

- `Find`: 查找第一个匹配的子串
- `FindString`: 查找第一个匹配的子串（忽略错误）
- `FindAll`: 查找所有匹配的子串
- `FindAllString`: 查找所有匹配的子串（忽略错误）
- `FindSubmatch`: 查找第一个匹配的子串和子组
- `FindAllSubmatch`: 查找所有匹配的子串和子组
- `Extract`: 提取匹配的子串
- `ExtractString`: 提取匹配的子串（忽略错误）
- `ExtractFirst`: 提取第一个匹配的子串
- `ExtractFirstString`: 提取第一个匹配的子串（忽略错误）
- `ExtractLast`: 提取最后一个匹配的子串
- `ExtractLastString`: 提取最后一个匹配的子串（忽略错误）
- `ExtractGroups`: 提取匹配的子组

### 2.3 替换操作

- `Replace`: 替换匹配的子串
- `ReplaceString`: 替换匹配的子串（忽略错误）
- `ReplaceAll`: 替换所有匹配的子串
- `ReplaceFunc`: 使用函数替换匹配的子串
- `ReplaceN`: 替换前n个匹配的子串
- `ReplaceNString`: 替换前n个匹配的子串（忽略错误）
- `ReplaceWithCallback`: 使用回调函数替换匹配的子串
- `Remove`: 移除匹配的子串
- `RemoveString`: 移除匹配的子串（忽略错误）

### 2.4 分割操作

- `Split`: 使用正则表达式分割字符串
- `SplitString`: 使用正则表达式分割字符串（忽略错误）

### 2.5 统计操作

- `Count`: 统计匹配的数量
- `CountString`: 统计匹配的数量（忽略错误）

### 2.6 位置操作

- `GetMatches`: 获取所有匹配的位置
- `GetMatchPositions`: 获取所有匹配的位置（返回开始和结束位置）

### 2.7 工具函数

- `IsValid`: 检查正则表达式是否有效
- `Escape`: 转义正则表达式特殊字符
- `Compile`: 编译正则表达式
- `MustCompile`: 编译正则表达式（失败则panic）
- `CompileWithOptions`: 编译正则表达式（带选项）
- `Validate`: 验证字符串是否匹配正则表达式

### 2.8 常用模式

- `EmailPattern`: 邮箱正则表达式
- `PhonePattern`: 手机号正则表达式（中国）
- `URLPattern`: URL正则表达式
- `IPv4Pattern`: IPv4地址正则表达式
- `IPv6Pattern`: IPv6地址正则表达式
- `UUIDPattern`: UUID正则表达式
- `DatePattern`: 日期正则表达式
- `TimePattern`: 时间正则表达式
- `DateTimePattern`: 日期时间正则表达式
- `ChinesePattern`: 中文字符正则表达式
- `NumberPattern`: 数字正则表达式
- `LetterPattern`: 字母正则表达式
- `AlphanumericPattern`: 字母数字正则表达式

### 2.9 快捷匹配函数

- `MatchEmail`: 检查是否为邮箱
- `MatchPhone`: 检查是否为手机号（中国）
- `MatchURL`: 检查是否为URL
- `MatchIPv4`: 检查是否为IPv4地址
- `MatchIPv6`: 检查是否为IPv6地址
- `MatchUUID`: 检查是否为UUID
- `MatchDate`: 检查是否为日期格式
- `MatchTime`: 检查是否为时间格式
- `MatchDateTime`: 检查是否为日期时间格式
- `MatchChinese`: 检查是否包含中文字符
- `MatchNumber`: 检查是否为数字
- `MatchLetter`: 检查是否为字母
- `MatchAlphanumeric`: 检查是否为字母数字

---

## 3. 使用示例

### 3.1 匹配操作

```go
import "github.com/yourusername/golang/pkg/utils/regex"

// 检查是否匹配
matched, err := regex.Match(`\d+`, "123")
if matched {
    // 匹配成功
}

// 检查是否匹配（忽略错误）
if regex.MatchString(`\d+`, "123") {
    // 匹配成功
}
```

### 3.2 查找操作

```go
// 查找第一个匹配
match, err := regex.Find(`\d+`, "abc123def")

// 查找所有匹配
matches, err := regex.FindAll(`\d+`, "abc123def456", -1)

// 提取匹配的子串
extracted := regex.ExtractString(`\d+`, "abc123def456")

// 提取第一个匹配
first := regex.ExtractFirstString(`\d+`, "abc123def456")

// 提取最后一个匹配
last := regex.ExtractLastString(`\d+`, "abc123def456")

// 提取命名子组
groups, err := regex.ExtractGroups(`(?P<name>\w+)\s+(?P<age>\d+)`, "John 30")
// groups: map[string]string{"name": "John", "age": "30"}
```

### 3.3 替换操作

```go
// 替换匹配的子串
result, err := regex.Replace(`\d+`, "abc123def", "XXX")

// 替换所有匹配
result, err := regex.ReplaceAll(`\d+`, "abc123def456", "XXX")

// 使用函数替换
result, err := regex.ReplaceFunc(`\d+`, "abc123def", func(match string) string {
    return "XXX"
})

// 替换前n个匹配
result, err := regex.ReplaceN(`\d+`, "abc123def456", "XXX", 1)

// 移除匹配的子串
result, err := regex.Remove(`\d+`, "abc123def")
```

### 3.4 分割操作

```go
// 使用正则表达式分割
parts, err := regex.Split(`\s+`, "a b c", -1)

// 分割字符串（忽略错误）
parts := regex.SplitString(`\s+`, "a b c", -1)
```

### 3.5 统计操作

```go
// 统计匹配数量
count, err := regex.Count(`\d+`, "abc123def456")

// 统计匹配数量（忽略错误）
count := regex.CountString(`\d+`, "abc123def456")
```

### 3.6 位置操作

```go
// 获取所有匹配的位置
matches, err := regex.GetMatches(`\d+`, "abc123def456")

// 获取所有匹配的位置（包含匹配内容）
positions, err := regex.GetMatchPositions(`\d+`, "abc123def456")
```

### 3.7 工具函数

```go
// 检查正则表达式是否有效
if regex.IsValid(`\d+`) {
    // 有效的正则表达式
}

// 转义特殊字符
escaped := regex.Escape(`.*+?^${}()|[]\`)

// 编译正则表达式
re, err := regex.Compile(`\d+`)

// 验证字符串是否匹配
err := regex.Validate(`\d+`, "123")
```

### 3.8 常用模式

```go
// 使用预定义模式
if regex.MatchEmail("test@example.com") {
    // 有效的邮箱
}

if regex.MatchPhone("13800138000") {
    // 有效的手机号
}

if regex.MatchURL("https://example.com") {
    // 有效的URL
}

// 使用预定义模式常量
matched, err := regex.Match(regex.EmailPattern, "test@example.com")
```

### 3.9 带选项的操作

```go
// 编译正则表达式（不区分大小写）
re, err := regex.CompileWithOptions(`abc`, false)

// 匹配（不区分大小写）
matched, err := regex.MatchWithOptions(`abc`, "ABC", false)

// 替换（不区分大小写）
result, err := regex.ReplaceWithOptions(`abc`, "ABC", "xyz", false)
```

---

**更新日期**: 2025-11-11
