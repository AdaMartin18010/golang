# 表格工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [表格工具](#表格工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 表格](#21-表格)
    - [2.2 简单表格](#22-简单表格)
    - [2.3 快捷函数](#23-快捷函数)
  - [3. 使用示例](#3-使用示例)
    - [3.1 基本表格](#31-基本表格)
    - [3.2 添加多行](#32-添加多行)
    - [3.3 获取表格字符串](#33-获取表格字符串)
    - [3.4 简单表格](#34-简单表格)
    - [3.5 快捷函数](#35-快捷函数)
    - [3.6 完整示例](#36-完整示例)

---

## 1. 概述

表格工具提供了终端表格输出功能，包括带边框的表格和简单表格，帮助开发者创建格式化的数据展示界面。

---

## 2. 功能特性

### 2.1 表格

- `Table`: 带边框的表格
- `NewTable`: 创建表格
- `AddRow`: 添加行
- `AddRows`: 添加多行
- `Render`: 渲染表格
- `Print`: 打印表格
- `String`: 返回表格字符串

### 2.2 简单表格

- `SimpleTable`: 简单表格（无边框）
- `NewSimpleTable`: 创建简单表格
- `SetSeparator`: 设置分隔符
- `AddRow`: 添加行
- `Render`: 渲染表格
- `Print`: 打印表格
- `String`: 返回表格字符串

### 2.3 快捷函数

- `FormatTable`: 格式化表格（快捷函数）
- `PrintTable`: 打印表格（快捷函数）
- `FormatSimpleTable`: 格式化简单表格（快捷函数）
- `PrintSimpleTable`: 打印简单表格（快捷函数）

---

## 3. 使用示例

### 3.1 基本表格

```go
import "github.com/yourusername/golang/pkg/utils/table"

// 创建表格
tbl := table.NewTable("Name", "Age", "City")

// 添加行
tbl.AddRow("Alice", "30", "Beijing")
tbl.AddRow("Bob", "25", "Shanghai")
tbl.AddRow("Charlie", "35", "Guangzhou")

// 打印表格
tbl.Print()
```

### 3.2 添加多行

```go
tbl := table.NewTable("Name", "Age", "City")
rows := [][]string{
    {"Alice", "30", "Beijing"},
    {"Bob", "25", "Shanghai"},
    {"Charlie", "35", "Guangzhou"},
}
tbl.AddRows(rows)
tbl.Print()
```

### 3.3 获取表格字符串

```go
tbl := table.NewTable("Name", "Age")
tbl.AddRow("Alice", "30")
tbl.AddRow("Bob", "25")
tableStr := tbl.String()
fmt.Print(tableStr)
```

### 3.4 简单表格

```go
// 创建简单表格
st := table.NewSimpleTable("Name", "Age", "City")

// 设置分隔符
st.SetSeparator("  ")

// 添加行
st.AddRow("Alice", "30", "Beijing")
st.AddRow("Bob", "25", "Shanghai")

// 打印表格
st.Print()
```

### 3.5 快捷函数

```go
// 格式化表格
headers := []string{"Name", "Age", "City"}
rows := [][]string{
    {"Alice", "30", "Beijing"},
    {"Bob", "25", "Shanghai"},
}
tableStr := table.FormatTable(headers, rows)
fmt.Print(tableStr)

// 打印表格
table.PrintTable(headers, rows)

// 格式化简单表格
simpleStr := table.FormatSimpleTable(headers, rows)
fmt.Print(simpleStr)

// 打印简单表格
table.PrintSimpleTable(headers, rows)
```

### 3.6 完整示例

```go
package main

import (
    "github.com/yourusername/golang/pkg/utils/table"
)

func main() {
    // 创建表格
    tbl := table.NewTable("Name", "Age", "City", "Email")

    // 添加数据
    tbl.AddRow("Alice", "30", "Beijing", "alice@example.com")
    tbl.AddRow("Bob", "25", "Shanghai", "bob@example.com")
    tbl.AddRow("Charlie", "35", "Guangzhou", "charlie@example.com")

    // 打印表格
    tbl.Print()

    // 使用简单表格
    st := table.NewSimpleTable("Name", "Age")
    st.AddRow("Alice", "30")
    st.AddRow("Bob", "25")
    st.Print()
}
```

---

**更新日期**: 2025-11-11
