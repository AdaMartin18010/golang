# 验证工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [验证工具](#验证工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 格式验证](#21-格式验证)
    - [2.2 字符串验证](#22-字符串验证)
    - [2.3 长度验证](#23-长度验证)
    - [2.4 字符串操作验证](#24-字符串操作验证)
    - [2.5 集合验证](#25-集合验证)
    - [2.6 数值验证](#26-数值验证)
    - [2.7 密码验证](#27-密码验证)
  - [3. 使用示例](#3-使用示例)
    - [3.1 格式验证](#31-格式验证)
    - [3.2 字符串验证](#32-字符串验证)
    - [3.3 集合验证](#33-集合验证)
    - [3.4 数值验证](#34-数值验证)
    - [3.5 密码验证](#35-密码验证)
    - [3.6 正则表达式验证](#36-正则表达式验证)
    - [3.7 日期时间验证](#37-日期时间验证)

---

## 1. 概述

验证工具提供了丰富的数据验证函数，简化常见的数据验证任务。

---

## 2. 功能特性

### 2.1 格式验证

- `IsEmail`: 验证邮箱
- `IsPhone`: 验证手机号（中国）
- `IsIDCard`: 验证身份证号（中国）
- `IsURL`: 验证URL
- `IsIPv4`: 验证IPv4地址
- `IsIPv6`: 验证IPv6地址
- `IsIP`: 验证IP地址（IPv4或IPv6）
- `IsUUID`: 验证UUID
- `IsCreditCard`: 验证信用卡号
- `IsDate`: 验证日期格式（YYYY-MM-DD）
- `IsTime`: 验证时间格式（HH:MM:SS）
- `IsDateTime`: 验证日期时间格式（YYYY-MM-DD HH:MM:SS）

### 2.2 字符串验证

- `IsEmpty`: 检查字符串是否为空
- `IsNotEmpty`: 检查字符串是否非空
- `IsNumeric`: 检查字符串是否为数字
- `IsAlpha`: 检查字符串是否只包含字母
- `IsAlphanumeric`: 检查字符串是否只包含字母和数字
- `IsLower`: 检查字符串是否全为小写
- `IsUpper`: 检查字符串是否全为大写
- `IsChinese`: 检查字符串是否只包含中文字符
- `HasChinese`: 检查字符串是否包含中文字符

### 2.3 长度验证

- `HasMinLength`: 检查字符串长度是否大于等于最小值
- `HasMaxLength`: 检查字符串长度是否小于等于最大值
- `HasLength`: 检查字符串长度是否在指定范围内

### 2.4 字符串操作验证

- `Contains`: 检查字符串是否包含子串
- `StartsWith`: 检查字符串是否以指定前缀开始
- `EndsWith`: 检查字符串是否以指定后缀结束
- `Matches`: 检查字符串是否匹配正则表达式

### 2.5 集合验证

- `IsIn`: 检查值是否在切片中
- `IsNotIn`: 检查值是否不在切片中

### 2.6 数值验证

- `IsBetween`: 检查数值是否在指定范围内
- `IsPositive`: 检查数值是否为正数
- `IsNegative`: 检查数值是否为负数
- `IsZero`: 检查数值是否为零
- `IsNonZero`: 检查数值是否非零

### 2.7 密码验证

- `IsStrongPassword`: 检查是否为强密码（至少8位，包含大小写字母、数字和特殊字符）
- `IsWeakPassword`: 检查是否为弱密码

---

## 3. 使用示例

### 3.1 格式验证

```go
import "github.com/yourusername/golang/pkg/utils/validator"

// 验证邮箱
if validator.IsEmail("user@example.com") {
    // 邮箱有效
}

// 验证手机号
if validator.IsPhone("13800138000") {
    // 手机号有效
}

// 验证身份证号
if validator.IsIDCard("110101199003075132") {
    // 身份证号有效
}

// 验证URL
if validator.IsURL("https://example.com") {
    // URL有效
}

// 验证IP地址
if validator.IsIPv4("192.168.1.1") {
    // IPv4地址有效
}

// 验证UUID
if validator.IsUUID("550e8400-e29b-41d4-a716-446655440000") {
    // UUID有效
}
```

### 3.2 字符串验证

```go
// 检查是否为空
if validator.IsEmpty("") {
    // 字符串为空
}

// 检查是否为数字
if validator.IsNumeric("123") {
    // 字符串为数字
}

// 检查是否只包含字母
if validator.IsAlpha("abc") {
    // 字符串只包含字母
}

// 检查长度
if validator.HasLength("test", 3, 5) {
    // 长度在范围内
}
```

### 3.3 集合验证

```go
// 检查值是否在切片中
if validator.IsIn(1, []int{1, 2, 3}) {
    // 值在切片中
}

// 检查值是否不在切片中
if validator.IsNotIn(4, []int{1, 2, 3}) {
    // 值不在切片中
}
```

### 3.4 数值验证

```go
// 检查数值是否在范围内
if validator.IsBetween(5, 1, 10) {
    // 数值在范围内
}

// 检查是否为正数
if validator.IsPositive(10) {
    // 数值为正数
}

// 检查是否为负数
if validator.IsNegative(-10) {
    // 数值为负数
}
```

### 3.5 密码验证

```go
// 检查是否为强密码
if validator.IsStrongPassword("Password123!") {
    // 强密码
}

// 检查是否为弱密码
if validator.IsWeakPassword("password") {
    // 弱密码
}
```

### 3.6 正则表达式验证

```go
// 检查字符串是否匹配正则表达式
if validator.Matches("test@example.com", `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`) {
    // 匹配成功
}
```

### 3.7 日期时间验证

```go
// 验证日期格式
if validator.IsDate("2025-11-11") {
    // 日期格式有效
}

// 验证时间格式
if validator.IsTime("12:30:45") {
    // 时间格式有效
}

// 验证日期时间格式
if validator.IsDateTime("2025-11-11 12:30:45") {
    // 日期时间格式有效
}
```

---

**更新日期**: 2025-11-11
