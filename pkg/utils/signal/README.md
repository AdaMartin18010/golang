# 信号处理工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [信号处理工具](#信号处理工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

信号处理工具提供了信号捕获和处理的便捷封装，简化信号处理任务，特别适用于优雅关闭等场景。

---

## 2. 功能特性

### 2.1 信号通知

- `Notify`: 注册信号通知
- `NotifyAll`: 注册所有信号通知
- `NotifyInterrupt`: 注册中断信号通知（SIGINT, SIGTERM）
- `NotifyQuit`: 注册退出信号通知（SIGQUIT）
- `NotifyHangup`: 注册挂起信号通知（SIGHUP）
- `NotifyUser1`: 注册用户信号1通知（SIGUSR1）
- `NotifyUser2`: 注册用户信号2通知（SIGUSR2）

### 2.2 信号处理

- `Handle`: 处理信号
- `HandleInterrupt`: 处理中断信号
- `HandleQuit`: 处理退出信号
- `HandleHangup`: 处理挂起信号

### 2.3 信号等待

- `Wait`: 等待信号
- `WaitInterrupt`: 等待中断信号
- `WaitQuit`: 等待退出信号
- `WaitHangup`: 等待挂起信号
- `WaitWithContext`: 使用context等待信号
- `WaitInterruptWithContext`: 使用context等待中断信号

### 2.4 优雅关闭

- `GracefulShutdown`: 优雅关闭处理器
- `NewGracefulShutdown`: 创建优雅关闭处理器
- `AddHandler`: 添加关闭处理函数
- `Wait`: 等待关闭信号并执行处理函数
- `WaitWithContext`: 使用context等待关闭信号
- `Start`: 启动优雅关闭处理（异步）

### 2.5 信号操作

- `Ignore`: 忽略信号
- `Reset`: 重置信号处理
- `Stop`: 停止信号通知
- `Send`: 发送信号到进程
- `SendInterrupt`: 发送中断信号到进程
- `SendTerminate`: 发送终止信号到进程
- `SendKill`: 发送杀死信号到进程

### 2.6 信号检查

- `IsInterrupt`: 检查信号是否为中断信号
- `IsQuit`: 检查信号是否为退出信号
- `IsHangup`: 检查信号是否为挂起信号
- `IsUser1`: 检查信号是否为用户信号1
- `IsUser2`: 检查信号是否为用户信号2
- `SignalName`: 获取信号名称

### 2.7 Context集成

- `WithContext`: 创建带信号取消的context
- `WithInterruptContext`: 创建带中断信号取消的context

---

## 3. 使用示例

### 3.1 信号通知

```go
import "github.com/yourusername/golang/pkg/utils/signal"

// 注册中断信号通知
sigChan := signal.NotifyInterrupt()

// 等待信号
sig := <-sigChan
fmt.Printf("Received signal: %s\n", signal.SignalName(sig))
```

### 3.2 信号处理

```go
// 处理中断信号
signal.HandleInterrupt(func(sig os.Signal) {
    fmt.Printf("Received interrupt signal: %s\n", signal.SignalName(sig))
    // 执行清理操作
})

// 处理多个信号
signal.Handle(func(sig os.Signal) {
    fmt.Printf("Received signal: %s\n", signal.SignalName(sig))
}, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
```

### 3.3 信号等待

```go
// 等待中断信号
sig := signal.WaitInterrupt()
fmt.Printf("Received signal: %s\n", signal.SignalName(sig))

// 使用context等待信号
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

sig, err := signal.WaitInterruptWithContext(ctx)
if err != nil {
    // 处理错误（可能是超时）
}
```

### 3.4 优雅关闭

```go
// 创建优雅关闭处理器
gs := signal.NewGracefulShutdown()

// 添加关闭处理函数
gs.AddHandler(func() {
    fmt.Println("Closing database connections...")
    // 关闭数据库连接
})

gs.AddHandler(func() {
    fmt.Println("Saving state...")
    // 保存状态
})

gs.AddHandler(func() {
    fmt.Println("Cleaning up resources...")
    // 清理资源
})

// 等待关闭信号并执行处理函数
sig := gs.Wait()
fmt.Printf("Shutdown complete, received signal: %s\n", signal.SignalName(sig))
```

### 3.5 异步优雅关闭

```go
// 启动优雅关闭处理（异步）
gs := signal.NewGracefulShutdown()
gs.AddHandler(func() {
    // 清理操作
})

done := gs.Start()

// 继续执行其他操作
// ...

// 等待关闭完成
sig := <-done
fmt.Printf("Shutdown complete, received signal: %s\n", signal.SignalName(sig))
```

### 3.6 信号检查

```go
sig := <-signal.NotifyInterrupt()

// 检查信号类型
if signal.IsInterrupt(sig) {
    fmt.Println("Received interrupt signal")
}

// 获取信号名称
name := signal.SignalName(sig)
fmt.Printf("Signal name: %s\n", name)
```

### 3.7 Context集成

```go
// 创建带中断信号取消的context
ctx, cancel := signal.WithInterruptContext(context.Background())
defer cancel()

// 在goroutine中使用
go func() {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Context cancelled by signal")
            return
        default:
            // 执行操作
        }
    }
}()
```

### 3.8 发送信号

```go
// 发送中断信号到进程
pid := 12345
err := signal.SendInterrupt(pid)
if err != nil {
    // 处理错误
}

// 发送终止信号
err = signal.SendTerminate(pid)
```

### 3.9 完整示例

```go
package main

import (
    "fmt"
    "os"
    "time"
    "github.com/yourusername/golang/pkg/utils/signal"
)

func main() {
    // 创建优雅关闭处理器
    gs := signal.NewGracefulShutdown()

    // 添加关闭处理函数
    gs.AddHandler(func() {
        fmt.Println("Shutting down gracefully...")
        // 执行清理操作
    })

    // 启动服务器
    go func() {
        // 服务器逻辑
        for {
            time.Sleep(time.Second)
            fmt.Println("Server running...")
        }
    }()

    // 等待关闭信号
    sig := gs.Wait()
    fmt.Printf("Server stopped, received signal: %s\n", signal.SignalName(sig))
    os.Exit(0)
}
```

---

**更新日期**: 2025-11-11
