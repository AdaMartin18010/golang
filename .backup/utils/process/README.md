# 进程工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [进程工具](#进程工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

进程工具提供了进程相关的功能，包括进程信息获取、命令执行、进程管理、信号处理等。

---

## 2. 功能特性

### 2.1 进程信息

- `GetPID`: 获取当前进程PID
- `GetPPID`: 获取父进程PID
- `GetExecutable`: 获取可执行文件路径
- `GetArgs`: 获取命令行参数
- `GetEnv`: 获取环境变量
- `SetEnv`: 设置环境变量
- `GetEnvAll`: 获取所有环境变量
- `GetWorkingDir`: 获取工作目录
- `ChangeDir`: 改变工作目录
- `GetProcessInfo`: 获取进程信息

### 2.2 命令执行

- `RunCommand`: 运行命令
- `RunCommandWithDir`: 在指定目录运行命令
- `RunCommandWithEnv`: 使用指定环境变量运行命令
- `RunCommandWithTimeout`: 带超时运行命令
- `StartCommand`: 启动命令（不等待完成）
- `WaitCommand`: 等待命令完成

### 2.3 进程管理

- `KillProcess`: 杀死进程
- `SignalProcess`: 向进程发送信号
- `IsProcessRunning`: 检查进程是否运行
- `WaitForProcess`: 等待进程结束

### 2.4 程序退出

- `Exit`: 退出程序
- `ExitSuccess`: 成功退出
- `ExitError`: 错误退出

### 2.5 信号处理

- `HandleSignals`: 处理信号
- `WaitForInterrupt`: 等待中断信号

### 2.6 守护进程

- `Daemonize`: 守护进程化（简单实现）

### 2.7 用户和权限

- `IsRoot`: 检查是否以root权限运行
- `GetUserID`: 获取用户ID
- `GetEffectiveUserID`: 获取有效用户ID
- `GetGroupID`: 获取组ID
- `GetEffectiveGroupID`: 获取有效组ID

### 2.8 系统信息

- `GetHostname`: 获取主机名
- `GetTempDir`: 获取临时目录
- `CreateTempFile`: 创建临时文件
- `CreateTempDir`: 创建临时目录

---

## 3. 使用示例

### 3.1 进程信息

```go
import "github.com/yourusername/golang/pkg/utils/process"

// 获取进程ID
pid := process.GetPID()
fmt.Printf("PID: %d\n", pid)

// 获取可执行文件路径
executable, err := process.GetExecutable()
if err == nil {
    fmt.Printf("Executable: %s\n", executable)
}

// 获取命令行参数
args := process.GetArgs()
fmt.Printf("Args: %v\n", args)

// 获取工作目录
dir, err := process.GetWorkingDir()
if err == nil {
    fmt.Printf("Working Dir: %s\n", dir)
}

// 获取进程信息
info, err := process.GetProcessInfo()
if err == nil {
    fmt.Printf("PID: %d, PPID: %d\n", info.PID, info.PPID)
}
```

### 3.2 命令执行

```go
// 运行命令
output, err := process.RunCommand("echo", "hello")
if err == nil {
    fmt.Printf("Output: %s\n", output)
}

// 在指定目录运行命令
output, err = process.RunCommandWithDir("/tmp", "ls", "-la")

// 使用指定环境变量运行命令
env := []string{"PATH=/usr/bin", "HOME=/home/user"}
output, err = process.RunCommandWithEnv(env, "echo", "$HOME")

// 带超时运行命令
output, err = process.RunCommandWithTimeout(5*time.Second, "sleep", "10")

// 启动命令（不等待完成）
cmd, err := process.StartCommand("sleep", "5")
if err == nil {
    // 执行其他操作
    err = process.WaitCommand(cmd)
}
```

### 3.3 进程管理

```go
// 检查进程是否运行
isRunning := process.IsProcessRunning(12345)

// 向进程发送信号
err := process.SignalProcess(12345, syscall.SIGTERM)

// 杀死进程
err = process.KillProcess(12345)

// 等待进程结束
err = process.WaitForProcess(12345)
```

### 3.4 程序退出

```go
// 退出程序
process.Exit(0)  // 成功退出
process.Exit(1)  // 错误退出

// 快捷函数
process.ExitSuccess()  // 成功退出
process.ExitError()    // 错误退出
```

### 3.5 信号处理

```go
// 处理信号
process.HandleSignals(func(sig os.Signal) {
    fmt.Printf("Received signal: %s\n", sig)
    // 执行清理操作
    process.ExitSuccess()
}, syscall.SIGINT, syscall.SIGTERM)

// 等待中断信号
sig := process.WaitForInterrupt()
fmt.Printf("Received signal: %s\n", sig)
```

### 3.6 用户和权限

```go
// 检查是否以root权限运行
if process.IsRoot() {
    fmt.Println("Running as root")
}

// 获取用户ID
uid := process.GetUserID()
fmt.Printf("User ID: %d\n", uid)

// 获取组ID
gid := process.GetGroupID()
fmt.Printf("Group ID: %d\n", gid)
```

### 3.7 系统信息

```go
// 获取主机名
hostname, err := process.GetHostname()
if err == nil {
    fmt.Printf("Hostname: %s\n", hostname)
}

// 获取临时目录
tempDir := process.GetTempDir()
fmt.Printf("Temp Dir: %s\n", tempDir)

// 创建临时文件
file, err := process.CreateTempFile("test-*.txt")
if err == nil {
    defer os.Remove(file.Name())
    // 使用文件
}

// 创建临时目录
dir, err := process.CreateTempDir("test-*")
if err == nil {
    defer os.RemoveAll(dir)
    // 使用目录
}
```

### 3.8 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/process"
)

func main() {
    // 获取进程信息
    pid := process.GetPID()
    fmt.Printf("PID: %d\n", pid)

    // 运行命令
    output, err := process.RunCommand("echo", "hello")
    if err == nil {
        fmt.Printf("Output: %s\n", output)
    }

    // 处理信号
    process.HandleSignals(func(sig os.Signal) {
        fmt.Printf("Received signal: %s\n", sig)
        process.ExitSuccess()
    })

    // 等待中断
    sig := process.WaitForInterrupt()
    fmt.Printf("Exiting with signal: %s\n", sig)
}
```

---

**更新日期**: 2025-11-11
