# 颜色工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [颜色工具](#颜色工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 文本颜色](#21-文本颜色)
    - [2.2 背景色](#22-背景色)
    - [2.3 文本样式](#23-文本样式)
    - [2.4 快捷函数](#24-快捷函数)
    - [2.5 打印函数](#25-打印函数)
    - [2.6 控制函数](#26-控制函数)
    - [2.7 通用函数](#27-通用函数)
  - [3. 使用示例](#3-使用示例)
    - [3.1 文本颜色](#31-文本颜色)
    - [3.2 背景色](#32-背景色)
    - [3.3 文本样式](#33-文本样式)
    - [3.4 快捷函数](#34-快捷函数)
    - [3.5 打印函数](#35-打印函数)
    - [3.6 控制函数](#36-控制函数)
    - [3.7 组合使用](#37-组合使用)
    - [3.8 完整示例](#38-完整示例)

---

## 1. 概述

颜色工具提供了终端颜色输出功能，包括文本颜色、背景色、文本样式等，帮助开发者创建更友好的命令行界面。

---

## 2. 功能特性

### 2.1 文本颜色

- `Black`: 黑色文本
- `Red`: 红色文本
- `Green`: 绿色文本
- `Yellow`: 黄色文本
- `Blue`: 蓝色文本
- `Magenta`: 洋红色文本
- `Cyan`: 青色文本
- `White`: 白色文本
- `BrightBlack`: 高亮黑色文本
- `BrightRed`: 高亮红色文本
- `BrightGreen`: 高亮绿色文本
- `BrightYellow`: 高亮黄色文本
- `BrightBlue`: 高亮蓝色文本
- `BrightMagenta`: 高亮洋红色文本
- `BrightCyan`: 高亮青色文本
- `BrightWhite`: 高亮白色文本

### 2.2 背景色

- `BgBlack`: 黑色背景
- `BgRed`: 红色背景
- `BgGreen`: 绿色背景
- `BgYellow`: 黄色背景
- `BgBlue`: 蓝色背景
- `BgMagenta`: 洋红色背景
- `BgCyan`: 青色背景
- `BgWhite`: 白色背景

### 2.3 文本样式

- `Bold`: 粗体文本
- `Dim`: 暗淡文本
- `Italic`: 斜体文本
- `Underline`: 下划线文本
- `Blink`: 闪烁文本
- `Reverse`: 反转文本
- `Hidden`: 隐藏文本

### 2.4 快捷函数

- `Success`: 成功消息（绿色）
- `Error`: 错误消息（红色）
- `Warning`: 警告消息（黄色）
- `Info`: 信息消息（蓝色）
- `Debug`: 调试消息（青色）

### 2.5 打印函数

- `Print`: 打印彩色文本
- `Println`: 打印彩色文本（换行）
- `Printf`: 格式化打印彩色文本
- `PrintSuccess`: 打印成功消息
- `PrintError`: 打印错误消息
- `PrintWarning`: 打印警告消息
- `PrintInfo`: 打印信息消息
- `PrintDebug`: 打印调试消息
- `PrintlnSuccess`: 打印成功消息（换行）
- `PrintlnError`: 打印错误消息（换行）
- `PrintlnWarning`: 打印警告消息（换行）
- `PrintlnInfo`: 打印信息消息（换行）
- `PrintlnDebug`: 打印调试消息（换行）

### 2.6 控制函数

- `Enable`: 启用颜色
- `Disable`: 禁用颜色
- `SetEnabled`: 设置是否启用颜色
- `IsEnabled`: 检查是否启用颜色
- `SetAutoDetect`: 设置是否自动检测终端支持

### 2.7 通用函数

- `Colorize`: 为文本添加颜色
- `ColorizeWithStyle`: 为文本添加颜色和样式

---

## 3. 使用示例

### 3.1 文本颜色

```go
import "github.com/yourusername/golang/pkg/utils/color"

// 基本颜色
fmt.Println(color.Red("红色文本"))
fmt.Println(color.Green("绿色文本"))
fmt.Println(color.Blue("蓝色文本"))
fmt.Println(color.Yellow("黄色文本"))

// 高亮颜色
fmt.Println(color.BrightRed("高亮红色文本"))
fmt.Println(color.BrightGreen("高亮绿色文本"))
```

### 3.2 背景色

```go
// 背景色
fmt.Println(color.BgRed("红色背景"))
fmt.Println(color.BgGreen("绿色背景"))
fmt.Println(color.BgBlue("蓝色背景"))
```

### 3.3 文本样式

```go
// 文本样式
fmt.Println(color.Bold("粗体文本"))
fmt.Println(color.Underline("下划线文本"))
fmt.Println(color.Italic("斜体文本"))

// 组合样式
fmt.Println(color.Bold(color.Red("粗体红色文本")))
```

### 3.4 快捷函数

```go
// 快捷函数
fmt.Println(color.Success("操作成功"))
fmt.Println(color.Error("操作失败"))
fmt.Println(color.Warning("警告信息"))
fmt.Println(color.Info("提示信息"))
fmt.Println(color.Debug("调试信息"))
```

### 3.5 打印函数

```go
// 打印函数
color.Print("红色文本", color.Red)
color.Println("绿色文本", color.Green)
color.Printf("蓝色文本: %s\n", color.Blue, "value")

// 快捷打印
color.PrintSuccess("操作成功")
color.PrintError("操作失败")
color.PrintWarning("警告信息")
color.PrintInfo("提示信息")
color.PrintDebug("调试信息")

// 换行打印
color.PrintlnSuccess("操作成功")
color.PrintlnError("操作失败")
color.PrintlnWarning("警告信息")
color.PrintlnInfo("提示信息")
color.PrintlnDebug("调试信息")
```

### 3.6 控制函数

```go
// 启用/禁用颜色
color.Enable()
color.Disable()
color.SetEnabled(true)

// 检查是否启用
if color.IsEnabled() {
    fmt.Println("颜色已启用")
}

// 设置自动检测
color.SetAutoDetect(true)
```

### 3.7 组合使用

```go
// 组合颜色和样式
text := color.Bold(color.Red("粗体红色文本"))
fmt.Println(text)

// 组合背景色和文本颜色
text = color.BgYellow(color.Black("黄色背景黑色文本"))
fmt.Println(text)

// 使用ColorizeWithStyle
text = color.ColorizeWithStyle("文本", color.Red, color.Bold, color.Underline)
fmt.Println(text)
```

### 3.8 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/color"
)

func main() {
    // 启用颜色
    color.Enable()

    // 打印不同颜色的文本
    fmt.Println(color.Red("错误信息"))
    fmt.Println(color.Green("成功信息"))
    fmt.Println(color.Yellow("警告信息"))
    fmt.Println(color.Blue("提示信息"))

    // 使用快捷函数
    color.PrintlnSuccess("操作成功")
    color.PrintlnError("操作失败")
    color.PrintlnWarning("警告信息")
    color.PrintlnInfo("提示信息")

    // 组合样式
    fmt.Println(color.Bold(color.Red("粗体红色文本")))
    fmt.Println(color.Underline(color.Blue("下划线蓝色文本")))
}
```

---

**更新日期**: 2025-11-11
