# 网络工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [网络工具](#网络工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

网络工具提供了网络相关的功能，包括IP地址验证、网络接口信息获取、端口检查、主机名解析、CIDR操作等。

---

## 2. 功能特性

### 2.1 IP地址验证

- `IsValidIP`: 检查是否为有效的IP地址
- `IsIPv4`: 检查是否为IPv4地址
- `IsIPv6`: 检查是否为IPv6地址
- `IsPrivateIP`: 检查是否为私有IP地址
- `IsLoopback`: 检查是否为回环地址
- `IsMulticast`: 检查是否为多播地址
- `IsUnspecified`: 检查是否为未指定地址

### 2.2 IP地址操作

- `ParseIP`: 解析IP地址
- `IPToInt`: 将IPv4地址转换为整数
- `IntToIP`: 将整数转换为IPv4地址

### 2.3 本地网络信息

- `GetLocalIP`: 获取本地IP地址
- `GetLocalIPs`: 获取所有本地IP地址
- `GetHostname`: 获取主机名
- `GetNetworkInfo`: 获取网络接口信息

### 2.4 主机名解析

- `ResolveIP`: 解析主机名到IP地址
- `ResolveHostname`: 解析IP地址到主机名

### 2.5 端口操作

- `IsPortOpen`: 检查端口是否开放
- `IsPortOpenTimeout`: 检查端口是否开放（带超时）
- `ValidatePort`: 验证端口号
- `GetFreePort`: 获取可用端口

### 2.6 CIDR操作

- `ParseCIDR`: 解析CIDR
- `IsIPInCIDR`: 检查IP是否在CIDR范围内

### 2.7 主机验证

- `ValidateHost`: 验证主机名或IP

### 2.8 MAC地址

- `FormatMAC`: 格式化MAC地址
- `IsValidMAC`: 检查是否为有效的MAC地址

### 2.9 网络连接

- `Ping`: 简单的ping实现（TCP连接测试）
- `IsReachable`: 检查主机是否可达

---

## 3. 使用示例

### 3.1 IP地址验证

```go
import "github.com/yourusername/golang/pkg/utils/network"

// 检查是否为有效的IP地址
isValid := network.IsValidIP("192.168.1.1")  // true

// 检查是否为IPv4地址
isIPv4 := network.IsIPv4("192.168.1.1")  // true

// 检查是否为IPv6地址
isIPv6 := network.IsIPv6("2001:db8::1")  // true

// 检查是否为私有IP地址
isPrivate := network.IsPrivateIP("192.168.1.1")  // true

// 检查是否为回环地址
isLoopback := network.IsLoopback("127.0.0.1")  // true
```

### 3.2 IP地址操作

```go
// 解析IP地址
ip := network.ParseIP("192.168.1.1")

// 将IPv4地址转换为整数
ipInt, err := network.IPToInt("192.168.1.1")  // 3232235777

// 将整数转换为IPv4地址
ipStr := network.IntToIP(3232235777)  // "192.168.1.1"
```

### 3.3 本地网络信息

```go
// 获取本地IP地址
localIP, err := network.GetLocalIP()

// 获取所有本地IP地址
localIPs, err := network.GetLocalIPs()

// 获取主机名
hostname, err := network.GetHostname()

// 获取网络接口信息
interfaces, err := network.GetNetworkInfo()
for _, iface := range interfaces {
    fmt.Printf("Interface: %s, IPs: %v\n", iface.Name, iface.IPs)
}
```

### 3.4 主机名解析

```go
// 解析主机名到IP地址
ips, err := network.ResolveIP("google.com")

// 解析IP地址到主机名
hostname, err := network.ResolveHostname("8.8.8.8")
```

### 3.5 端口操作

```go
// 检查端口是否开放
isOpen := network.IsPortOpen("localhost", 8080)

// 检查端口是否开放（带超时）
isOpen = network.IsPortOpenTimeout("localhost", 8080, 5)

// 验证端口号
isValid := network.ValidatePort(8080)  // true

// 获取可用端口
port, err := network.GetFreePort()
```

### 3.6 CIDR操作

```go
// 解析CIDR
ipnet, err := network.ParseCIDR("192.168.1.0/24")

// 检查IP是否在CIDR范围内
inRange := network.IsIPInCIDR("192.168.1.1", "192.168.1.0/24")  // true
```

### 3.7 主机验证

```go
// 验证主机名或IP
isValid := network.ValidateHost("localhost")  // true
isValid = network.ValidateHost("192.168.1.1")  // true
```

### 3.8 MAC地址

```go
// 格式化MAC地址
formatted := network.FormatMAC("00:11:22:33:44:55")

// 检查是否为有效的MAC地址
isValid := network.IsValidMAC("00:11:22:33:44:55")  // true
```

### 3.9 网络连接

```go
// Ping（TCP连接测试）
isReachable := network.Ping("google.com", 80, 5)

// 检查主机是否可达
isReachable = network.IsReachable("google.com", 80)
```

### 3.10 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/network"
)

func main() {
    // IP地址验证
    fmt.Printf("IsValidIP: %v\n", network.IsValidIP("192.168.1.1"))
    fmt.Printf("IsIPv4: %v\n", network.IsIPv4("192.168.1.1"))
    fmt.Printf("IsPrivateIP: %v\n", network.IsPrivateIP("192.168.1.1"))

    // 获取本地IP
    localIP, err := network.GetLocalIP()
    if err == nil {
        fmt.Printf("Local IP: %s\n", localIP)
    }

    // 检查端口
    isOpen := network.IsPortOpen("localhost", 8080)
    fmt.Printf("Port 8080 open: %v\n", isOpen)

    // 获取可用端口
    port, err := network.GetFreePort()
    if err == nil {
        fmt.Printf("Free port: %d\n", port)
    }
}
```

---

**更新日期**: 2025-11-11
