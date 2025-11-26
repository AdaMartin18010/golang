# 系统信息工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [系统信息工具](#系统信息工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 系统信息](#21-系统信息)
    - [2.2 内存信息](#22-内存信息)
    - [2.3 CPU信息](#23-cpu信息)
    - [2.4 堆栈信息](#24-堆栈信息)
    - [2.5 系统检查](#25-系统检查)
    - [2.6 其他信息](#26-其他信息)
    - [2.7 系统监控](#27-系统监控)
  - [3. 使用示例](#3-使用示例)
    - [3.1 系统信息](#31-系统信息)
    - [3.2 内存信息](#32-内存信息)
    - [3.3 CPU信息](#33-cpu信息)
    - [3.4 堆栈信息](#34-堆栈信息)
    - [3.5 系统检查](#35-系统检查)
    - [3.6 系统监控](#36-系统监控)
    - [3.7 完整示例](#37-完整示例)

---

## 1. 概述

系统信息工具提供了系统资源信息获取功能，包括系统信息、内存信息、CPU信息等，适用于系统监控和性能分析场景。

---

## 2. 功能特性

### 2.1 系统信息

- `GetSystemInfo`: 获取系统信息
- `GetOS`: 获取操作系统
- `GetArch`: 获取架构
- `GetGoVersion`: 获取Go版本
- `GetNumCPU`: 获取CPU核心数
- `GetNumGoroutine`: 获取Goroutine数量
- `GetGOMAXPROCS`: 获取GOMAXPROCS设置
- `SetGOMAXPROCS`: 设置GOMAXPROCS

### 2.2 内存信息

- `GetMemoryInfo`: 获取内存信息
- `GetAllocMemory`: 获取已分配的内存（字节）
- `GetTotalAllocMemory`: 获取累计分配的内存（字节）
- `GetSysMemory`: 获取系统内存（字节）
- `GetNumGC`: 获取GC次数
- `GetMemoryUsagePercent`: 获取内存使用率（百分比）
- `FormatBytes`: 格式化字节数
- `GC`: 执行GC
- `FreeOSMemory`: 释放OS内存

### 2.3 CPU信息

- `GetCPUInfo`: 获取CPU信息

### 2.4 堆栈信息

- `GetStack`: 获取当前goroutine的堆栈信息
- `GetAllStacks`: 获取所有goroutine的堆栈信息
- `GetCaller`: 获取调用者信息
- `GetCallers`: 获取调用栈
- `GetFuncName`: 获取函数名
- `GetFileLine`: 获取文件和行号

### 2.5 系统检查

- `IsWindows`: 检查是否为Windows系统
- `IsLinux`: 检查是否为Linux系统
- `IsDarwin`: 检查是否为Darwin系统（macOS）
- `IsUnix`: 检查是否为Unix系统
- `IsAMD64`: 检查是否为AMD64架构
- `IsARM64`: 检查是否为ARM64架构
- `Is386`: 检查是否为386架构

### 2.6 其他信息

- `GetCompiler`: 获取编译器信息
- `GetNumCgoCall`: 获取CGO调用次数

### 2.7 系统监控

- `Monitor`: 系统监控器
- `NewMonitor`: 创建系统监控器
- `Start`: 启动监控
- `Stop`: 停止监控

---

## 3. 使用示例

### 3.1 系统信息

```go
import "github.com/yourusername/golang/pkg/utils/sysinfo"

// 获取系统信息
info := sysinfo.GetSystemInfo()
fmt.Printf("OS: %s\n", info.OS)
fmt.Printf("Arch: %s\n", info.Arch)
fmt.Printf("Go Version: %s\n", info.GoVersion)
fmt.Printf("CPU Cores: %d\n", info.NumCPU)
fmt.Printf("Goroutines: %d\n", info.NumGoroutine)

// 获取操作系统
os := sysinfo.GetOS()
fmt.Printf("OS: %s\n", os)

// 获取架构
arch := sysinfo.GetArch()
fmt.Printf("Arch: %s\n", arch)

// 获取Go版本
version := sysinfo.GetGoVersion()
fmt.Printf("Go Version: %s\n", version)

// 获取CPU核心数
numCPU := sysinfo.GetNumCPU()
fmt.Printf("CPU Cores: %d\n", numCPU)

// 获取Goroutine数量
numGoroutine := sysinfo.GetNumGoroutine()
fmt.Printf("Goroutines: %d\n", numGoroutine)
```

### 3.2 内存信息

```go
// 获取内存信息
memInfo := sysinfo.GetMemoryInfo()
fmt.Printf("Alloc: %d bytes\n", memInfo.Alloc)
fmt.Printf("Total Alloc: %d bytes\n", memInfo.TotalAlloc)
fmt.Printf("Sys: %d bytes\n", memInfo.Sys)
fmt.Printf("GC Count: %d\n", memInfo.NumGC)

// 获取已分配的内存
alloc := sysinfo.GetAllocMemory()
fmt.Printf("Alloc: %s\n", sysinfo.FormatBytes(alloc))

// 获取累计分配的内存
totalAlloc := sysinfo.GetTotalAllocMemory()
fmt.Printf("Total Alloc: %s\n", sysinfo.FormatBytes(totalAlloc))

// 获取系统内存
sys := sysinfo.GetSysMemory()
fmt.Printf("Sys: %s\n", sysinfo.FormatBytes(sys))

// 获取内存使用率
usage := sysinfo.GetMemoryUsagePercent()
fmt.Printf("Memory Usage: %.2f%%\n", usage)

// 格式化字节数
formatted := sysinfo.FormatBytes(1024 * 1024)  // "1.0 MB"
fmt.Printf("Formatted: %s\n", formatted)

// 执行GC
sysinfo.GC()

// 释放OS内存
sysinfo.FreeOSMemory()
```

### 3.3 CPU信息

```go
// 获取CPU信息
cpuInfo := sysinfo.GetCPUInfo()
fmt.Printf("CPU Cores: %d\n", cpuInfo.NumCPU)
fmt.Printf("Goroutines: %d\n", cpuInfo.NumGoroutine)
fmt.Printf("GOMAXPROCS: %d\n", cpuInfo.GOMAXPROCS)

// 设置GOMAXPROCS
old := sysinfo.SetGOMAXPROCS(4)
fmt.Printf("Old GOMAXPROCS: %d\n", old)
```

### 3.4 堆栈信息

```go
// 获取当前goroutine的堆栈信息
stack := sysinfo.GetStack()
fmt.Printf("Stack: %s\n", string(stack))

// 获取所有goroutine的堆栈信息
allStacks := sysinfo.GetAllStacks()
fmt.Printf("All Stacks: %s\n", string(allStacks))

// 获取调用者信息
pc, file, line, ok := sysinfo.GetCaller(0)
if ok {
    fmt.Printf("Caller: %s:%d\n", file, line)
    funcName := sysinfo.GetFuncName(pc)
    fmt.Printf("Function: %s\n", funcName)
}
```

### 3.5 系统检查

```go
// 检查操作系统
if sysinfo.IsWindows() {
    fmt.Println("Running on Windows")
}

if sysinfo.IsLinux() {
    fmt.Println("Running on Linux")
}

if sysinfo.IsDarwin() {
    fmt.Println("Running on macOS")
}

if sysinfo.IsUnix() {
    fmt.Println("Running on Unix-like system")
}

// 检查架构
if sysinfo.IsAMD64() {
    fmt.Println("AMD64 architecture")
}

if sysinfo.IsARM64() {
    fmt.Println("ARM64 architecture")
}
```

### 3.6 系统监控

```go
// 创建系统监控器
monitor := sysinfo.NewMonitor(5*time.Second, func(sysInfo *sysinfo.SystemInfo, memInfo *sysinfo.MemoryInfo, cpuInfo *sysinfo.CPUInfo) {
    fmt.Printf("OS: %s, Arch: %s\n", sysInfo.OS, sysInfo.Arch)
    fmt.Printf("Alloc: %s, Total Alloc: %s\n",
        sysinfo.FormatBytes(memInfo.Alloc),
        sysinfo.FormatBytes(memInfo.TotalAlloc))
    fmt.Printf("Goroutines: %d\n", cpuInfo.NumGoroutine)
})

// 启动监控
go monitor.Start()

// 执行其他操作
// ...

// 停止监控
monitor.Stop()
```

### 3.7 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/golang/pkg/utils/sysinfo"
)

func main() {
    // 获取系统信息
    info := sysinfo.GetSystemInfo()
    fmt.Printf("System Info:\n")
    fmt.Printf("  OS: %s\n", info.OS)
    fmt.Printf("  Arch: %s\n", info.Arch)
    fmt.Printf("  Go Version: %s\n", info.GoVersion)
    fmt.Printf("  CPU Cores: %d\n", info.NumCPU)
    fmt.Printf("  Goroutines: %d\n", info.NumGoroutine)

    // 获取内存信息
    memInfo := sysinfo.GetMemoryInfo()
    fmt.Printf("\nMemory Info:\n")
    fmt.Printf("  Alloc: %s\n", sysinfo.FormatBytes(memInfo.Alloc))
    fmt.Printf("  Total Alloc: %s\n", sysinfo.FormatBytes(memInfo.TotalAlloc))
    fmt.Printf("  Sys: %s\n", sysinfo.FormatBytes(memInfo.Sys))
    fmt.Printf("  GC Count: %d\n", memInfo.NumGC)
    fmt.Printf("  Memory Usage: %.2f%%\n", sysinfo.GetMemoryUsagePercent())

    // 创建监控器
    monitor := sysinfo.NewMonitor(2*time.Second, func(sysInfo *sysinfo.SystemInfo, memInfo *sysinfo.MemoryInfo, cpuInfo *sysinfo.CPUInfo) {
        fmt.Printf("\n[Monitor] Goroutines: %d, Alloc: %s\n",
            cpuInfo.NumGoroutine,
            sysinfo.FormatBytes(memInfo.Alloc))
    })

    // 启动监控
    go monitor.Start()

    // 等待一段时间
    time.Sleep(10 * time.Second)

    // 停止监控
    monitor.Stop()
    fmt.Println("\nMonitor stopped")
}
```

---

**更新日期**: 2025-11-11
