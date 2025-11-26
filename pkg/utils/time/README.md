# 时间工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [时间工具](#时间工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

时间工具提供了丰富的时间操作函数，简化常见的时间处理任务。

---

## 2. 功能特性

### 2.1 时间戳获取

- `Unix`: 获取Unix时间戳（秒）
- `UnixMilli`: 获取Unix时间戳（毫秒）
- `UnixMicro`: 获取Unix时间戳（微秒）
- `UnixNano`: 获取Unix时间戳（纳秒）

### 2.2 时间格式化

- `Format`: 格式化时间
- `FormatDefault`: 使用默认格式格式化时间
- `FormatDate`: 格式化日期
- `FormatTime`: 格式化时间

### 2.3 时间解析

- `Parse`: 解析时间字符串
- `ParseDefault`: 使用默认格式解析时间字符串
- `ParseDate`: 解析日期字符串

### 2.4 时间计算

- `AddDays`: 添加天数
- `AddMonths`: 添加月数
- `AddYears`: 添加年数

### 2.5 时间范围

- `StartOfDay`: 获取一天的开始时间
- `EndOfDay`: 获取一天的结束时间
- `StartOfWeek`: 获取一周的开始时间
- `EndOfWeek`: 获取一周的结束时间
- `StartOfMonth`: 获取一月的开始时间
- `EndOfMonth`: 获取一月的结束时间
- `StartOfYear`: 获取一年的开始时间
- `EndOfYear`: 获取一年的结束时间

### 2.6 时间比较

- `DaysBetween`: 计算两个时间之间的天数
- `HoursBetween`: 计算两个时间之间的小时数
- `MinutesBetween`: 计算两个时间之间的分钟数
- `SecondsBetween`: 计算两个时间之间的秒数
- `IsToday`: 判断是否是今天
- `IsYesterday`: 判断是否是昨天
- `IsTomorrow`: 判断是否是明天
- `IsSameDay`: 判断是否是同一天
- `IsSameWeek`: 判断是否是同一周
- `IsSameMonth`: 判断是否是同一月
- `IsSameYear`: 判断是否是同一年

### 2.7 人性化显示

- `HumanizeDuration`: 人性化显示时长
- `HumanizeTime`: 人性化显示时间

---

## 3. 使用示例

### 3.1 时间戳获取

```go
import "github.com/yourusername/golang/pkg/utils/time"

// 获取Unix时间戳（秒）
timestamp := time.Unix()

// 获取Unix时间戳（毫秒）
timestampMs := time.UnixMilli()
```

### 3.2 时间格式化

```go
now := time.Now()

// 默认格式
formatted := time.FormatDefault(now) // "2023-01-02 15:04:05"

// 日期格式
date := time.FormatDate(now) // "2023-01-02"

// 时间格式
tm := time.FormatTime(now) // "15:04:05"
```

### 3.3 时间计算

```go
now := time.Now()

// 添加天数
tomorrow := time.AddDays(now, 1)

// 添加月数
nextMonth := time.AddMonths(now, 1)

// 添加年数
nextYear := time.AddYears(now, 1)
```

### 3.4 时间范围

```go
now := time.Now()

// 一天的开始和结束
start := time.StartOfDay(now)
end := time.EndOfDay(now)

// 一周的开始和结束
weekStart := time.StartOfWeek(now)
weekEnd := time.EndOfWeek(now)

// 一月的开始和结束
monthStart := time.StartOfMonth(now)
monthEnd := time.EndOfMonth(now)
```

### 3.5 时间比较

```go
t1 := time.Now()
t2 := time.AddDays(t1, 5)

// 计算天数差
days := time.DaysBetween(t1, t2) // 5

// 判断是否是今天
if time.IsToday(t1) {
    // 处理今天
}

// 判断是否是同一天
if time.IsSameDay(t1, t2) {
    // 处理同一天
}
```

### 3.6 人性化显示

```go
duration := 2 * time.Hour
humanized := time.HumanizeDuration(duration) // "2小时"

pastTime := time.AddHours(time.Now(), -3)
humanized := time.HumanizeTime(pastTime) // "3小时前"
```

---

**更新日期**: 2025-11-11
