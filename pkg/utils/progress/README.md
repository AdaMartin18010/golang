# 进度条工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [进度条工具](#进度条工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

进度条工具提供了进度条和旋转器功能，帮助开发者创建友好的命令行进度显示界面。

---

## 2. 功能特性

### 2.1 进度条

- `ProgressBar`: 完整功能的进度条
- `NewProgressBar`: 创建进度条
- `Add`: 增加进度
- `Set`: 设置进度
- `Increment`: 增加1
- `SetTotal`: 设置总数
- `Current`: 获取当前进度
- `Total`: 获取总数
- `Percent`: 获取百分比
- `Finish`: 完成
- `Reset`: 重置

### 2.2 进度条选项

- `WithWidth`: 设置进度条宽度
- `WithShowPercent`: 设置是否显示百分比
- `WithShowSpeed`: 设置是否显示速度
- `WithShowETA`: 设置是否显示预计剩余时间
- `WithWriter`: 设置输出写入器
- `WithPrefix`: 设置前缀
- `WithSuffix`: 设置后缀

### 2.3 简单进度条

- `SimpleProgressBar`: 简单进度条
- `NewSimpleProgressBar`: 创建简单进度条

### 2.4 旋转器

- `Spinner`: 旋转器
- `NewSpinner`: 创建旋转器
- `Start`: 启动旋转器
- `Stop`: 停止旋转器
- `StopWithMessage`: 停止并显示消息

---

## 3. 使用示例

### 3.1 基本进度条

```go
import "github.com/yourusername/golang/pkg/utils/progress"

// 创建进度条
pb := progress.NewProgressBar(100)

// 更新进度
for i := 0; i < 100; i++ {
    pb.Increment()
    time.Sleep(10 * time.Millisecond)
}

// 完成
pb.Finish()
```

### 3.2 自定义进度条

```go
// 创建自定义进度条
pb := progress.NewProgressBar(1000,
    progress.WithWidth(60),
    progress.WithShowPercent(true),
    progress.WithShowSpeed(true),
    progress.WithShowETA(true),
    progress.WithPrefix("Processing"),
    progress.WithSuffix("items"),
)

// 更新进度
for i := 0; i < 1000; i++ {
    pb.Add(1)
    time.Sleep(1 * time.Millisecond)
}

// 完成
pb.Finish()
```

### 3.3 简单进度条

```go
// 创建简单进度条
spb := progress.NewSimpleProgressBar(100)

// 更新进度
for i := 0; i < 100; i++ {
    spb.Increment()
    time.Sleep(10 * time.Millisecond)
}

// 完成
spb.Finish()
```

### 3.4 旋转器

```go
// 创建旋转器
spinner := progress.NewSpinner("Loading")

// 启动旋转器
spinner.Start()

// 执行操作
time.Sleep(2 * time.Second)

// 停止旋转器
spinner.Stop()

// 或停止并显示消息
spinner.StopWithMessage("Done!")
```

### 3.5 完整示例

```go
package main

import (
    "time"
    "github.com/yourusername/golang/pkg/utils/progress"
)

func main() {
    // 创建进度条
    pb := progress.NewProgressBar(100,
        progress.WithShowPercent(true),
        progress.WithShowSpeed(true),
        progress.WithShowETA(true),
        progress.WithPrefix("Downloading"),
    )

    // 模拟下载
    for i := 0; i < 100; i++ {
        pb.Increment()
        time.Sleep(50 * time.Millisecond)
    }

    // 完成
    pb.Finish()

    // 使用旋转器
    spinner := progress.NewSpinner("Processing")
    spinner.Start()
    time.Sleep(2 * time.Second)
    spinner.StopWithMessage("Complete!")
}
```

---

**更新日期**: 2025-11-11
