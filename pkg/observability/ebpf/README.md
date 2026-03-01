# eBPF 可观测性实现

**版本**: v2.0
**更新日期**: 2025-12-03
**使用库**: github.com/cilium/ebpf (最成熟的Go eBPF库)

---

## 📋 功能概述

本模块使用 **Cilium eBPF** 库实现系统级可观测性，提供：

1. **系统调用追踪** - 追踪应用的系统调用（execve, open, connect 等）
2. **网络监控** - 监控 TCP 连接和网络流量
3. **性能分析** - CPU、内存、I/O性能
4. **安全监控** - 检测异常行为

---

## 🏗️ 架构设计

### 分层结构

```text
┌─────────────────────────────────────┐
│   OpenTelemetry Metrics/Traces     │  ← 输出层
├─────────────────────────────────────┤
│   eBPF Collector (Go)              │  ← 数据处理层
├─────────────────────────────────────┤
│   eBPF Maps (Kernel)               │  ← 数据缓冲层
├─────────────────────────────────────┤
│   eBPF Programs (C)                │  ← 数据采集层
├─────────────────────────────────────┤
│   Kernel Events                     │  ← 事件源
└─────────────────────────────────────┘
```

### 数据流

```text
Kernel Event → eBPF Program → eBPF Map → Go Collector → OTLP → Backend
```

---

## 📁 文件结构

```
pkg/observability/ebpf/
├── collector.go              # 主收集器
├── syscall_tracer.go         # 系统调用追踪器
├── network_tracer.go         # 网络追踪器
├── syscall_linux.go          # Linux 平台实现
├── syscall_stub.go           # 非 Linux 平台 stub
├── network_linux.go          # Linux 平台实现
├── network_stub.go           # 非 Linux 平台 stub
├── programs/
│   ├── syscall.bpf.c        # 系统调用 eBPF C 程序
│   └── network.bpf.c        # 网络监控 eBPF C 程序
└── README.md                # 本文档
```

---

## 🚀 快速开始

### 1. 环境要求

#### Linux 内核要求

- **最低版本**: Linux 4.18+
- **推荐版本**: Linux 5.10+ (LTS)
- **最佳版本**: Linux 6.1+ (最新特性)

#### 开发工具

```bash
# 安装 LLVM/Clang (编译 eBPF 程序)
sudo apt-get install clang llvm

# 安装 bpf2go (生成 Go 绑定)
go install github.com/cilium/ebpf/cmd/bpf2go@latest

# 安装内核头文件
sudo apt-get install linux-headers-$(uname -r)
```

### 2. 编译 eBPF 程序

```bash
# 生成 Go 绑定（在 Linux 系统上执行）
go generate ./pkg/observability/ebpf/...

# 或者在项目根目录执行
make generate-ebpf
```

这会生成以下文件：
- `syscall_bpfel.go` / `syscall_bpfeb.go` - 系统调用程序绑定（小端/大端）
- `network_bpfel.go` / `network_bpfeb.go` - 网络程序绑定（小端/大端）

### 3. 使用示例

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/yourusername/golang/pkg/observability/ebpf"
    "go.opentelemetry.io/otel"
)

func main() {
    // 创建 OTLP tracer 和 meter
    tracer := otel.Tracer("ebpf-collector")
    meter := otel.Meter("ebpf-collector")

    // 创建 eBPF 收集器
    collector, err := ebpf.NewCollector(ebpf.Config{
        Tracer:                  tracer,
        Meter:                   meter,
        Enabled:                 true,
        CollectInterval:         5 * time.Second,
        EnableSyscallTracking:   true,
        EnableNetworkMonitoring: true,
    })
    if err != nil {
        log.Fatal(err)
    }

    // 启动收集
    if err := collector.Start(); err != nil {
        log.Fatal(err)
    }
    defer collector.Stop()

    // 获取统计数据
    ctx := context.Background()
    
    // 获取系统调用统计
    syscallStats, _ := collector.GetSyscallStats(ctx)
    log.Printf("Syscall stats: %+v", syscallStats)
    
    // 获取活跃连接数
    activeConns, _ := collector.GetActiveConnections(ctx)
    log.Printf("Active connections: %d", activeConns)
    
    // 获取连接详情
    connDetails, _ := collector.GetConnectionDetails(ctx)
    for _, conn := range connDetails {
        log.Printf("Connection: %s:%d -> %s:%d", 
            conn.SrcIP, conn.SrcPort, conn.DstIP, conn.DstPort)
    }

    // 应用运行...
    select {}
}
```

---

## 📖 API 参考

### Collector (收集器)

```go
// 创建收集器
collector, err := ebpf.NewCollector(ebpf.Config{
    Enabled:                    true,           // 启用 eBPF
    Tracer:                     tracer,         // OpenTelemetry Tracer
    Meter:                      meter,          // OpenTelemetry Meter
    CollectInterval:            5 * time.Second,// 收集间隔
    EnableSyscallTracking:      true,           // 启用系统调用追踪
    EnableNetworkMonitoring:    true,           // 启用网络监控
    EnablePerformanceProfiling: false,          // 启用性能分析
})

// 启动收集器
err = collector.Start()

// 停止收集器
err = collector.Stop()

// 获取状态
status := collector.GetStatus()

// 获取系统调用统计
stats, err := collector.GetSyscallStats(ctx)

// 获取活跃连接数
count, err := collector.GetActiveConnections(ctx)

// 获取连接详情
details, err := collector.GetConnectionDetails(ctx)
```

### SyscallTracer (系统调用追踪器)

```go
// 创建追踪器
tracer, err := ebpf.NewSyscallTracer(ebpf.SyscallTracerConfig{
    Enabled:    true,
    Tracer:     tracer,
    Meter:      meter,
    TargetPID:  0,          // 0 = 所有进程
    BufferSize: 65536,      // perf buffer 大小
})

// 启动追踪
err = tracer.Start()

// 停止追踪
err = tracer.Stop()

// 获取统计
stats, err := tracer.GetSyscallStats(ctx)
```

### NetworkTracer (网络追踪器)

```go
// 创建追踪器
tracer, err := ebpf.NewNetworkTracer(ebpf.NetworkTracerConfig{
    Enabled:       true,
    Tracer:        tracer,
    Meter:         meter,
    TrackInbound:  true,    // 追踪入站连接
    TrackOutbound: true,    // 追踪出站连接
    BufferSize:    65536,
})

// 启动追踪
err = tracer.Start()

// 停止追踪
err = tracer.Stop()

// 获取活跃连接数
count, err := tracer.GetActiveConnections(ctx)

// 获取连接统计
stats, err := tracer.GetConnectionStats(ctx)

// 获取连接详情
details, err := tracer.GetConnectionDetails(ctx)
```

---

## 🔧 实现详情

### 系统调用追踪 (syscall_tracer.go)

**功能**:
- 使用 tracepoint 追踪 sys_enter/sys_exit
- 记录系统调用 ID、进程 ID、执行时间
- 通过 perf buffer 发送事件到用户空间
- 支持按系统调用类型统计

**eBPF Maps**:
- `syscall_events`: Perf event array，发送事件到用户空间
- `syscall_stats`: Hash map，统计各系统调用次数
- `syscall_start_time`: Hash map，记录系统调用开始时间

### 网络监控 (network_tracer.go)

**功能**:
- 使用 kprobe 追踪 TCP 连接生命周期
- 监控 connect/accept/close 事件
- 记录源/目的 IP 和端口
- 统计发送/接收字节数
- 计算连接持续时间

**eBPF Maps**:
- `tcp_events`: Perf event array，发送事件到用户空间
- `tcp_connections`: Hash map，存储活跃连接信息
- `tcp_stats`: Hash map，按 PID 统计连接数

**追踪的内核函数**:
- `tcp_connect`: 出站连接
- `tcp_v4_connect`: IPv4 连接（返回值追踪）
- `inet_csk_accept`: 入站连接
- `tcp_close`: 连接关闭
- `tcp_sendmsg`: 发送数据
- `tcp_recvmsg`: 接收数据

---

## 📊 指标说明

### 系统调用指标

| 指标名称 | 类型 | 描述 |
|---------|------|------|
| `ebpf.syscall.count` | Counter | 系统调用总数 |
| `ebpf.syscall.duration` | Histogram | 系统调用执行时间（毫秒） |

**属性**:
- `syscall.id`: 系统调用 ID
- `process.pid`: 进程 ID
- `syscall.return`: 返回值

### 网络指标

| 指标名称 | 类型 | 描述 |
|---------|------|------|
| `ebpf.tcp.connections` | Counter | TCP 连接总数 |
| `ebpf.tcp.connections.active` | UpDownCounter | 活跃连接数 |
| `ebpf.tcp.bytes` | Counter | 传输字节数 |
| `ebpf.tcp.connection.duration` | Histogram | 连接持续时间（毫秒） |

**属性**:
- `tcp.event_type`: 事件类型 (connect/accept/close)
- `process.pid`: 进程 ID
- `net.peer.ip`: 对端 IP
- `net.peer.port`: 对端端口
- `net.host.ip`: 本地 IP
- `net.host.port`: 本地端口

---

## ⚠️ 注意事项

### 权限要求

运行 eBPF 程序需要以下权限之一：

1. **root 用户**
2. **CAP_BPF** capability
3. **CAP_SYS_ADMIN** capability (旧内核)

```bash
# 添加 capability（不需要 root）
sudo setcap cap_bpf,cap_net_admin,cap_perfmon+ep /path/to/your/binary
```

### 内核兼容性

不同内核版本支持的 eBPF 特性：

| 内核版本 | 特性支持 |
|---------|---------|
| 4.18+ | 基础 eBPF，tracepoint，kprobe |
| 5.2+ | BPF_PROG_TYPE_SOCK_OPS |
| 5.8+ | BPF_PROG_TYPE_EXT，BPF LSM |
| 6.1+ | 完整特性 |

### 性能影响

eBPF 程序在内核中执行，开销很小：

- **系统调用追踪**: ~0.1-0.5% CPU 开销
- **网络监控**: ~0.5-1% CPU 开销

可以通过以下方式优化：
- 增加 perf buffer 大小
- 减少采样频率
- 只追踪特定 PID

---

## 🐛 故障排查

### 常见问题

**1. 编译错误：找不到 bpf 头文件**

```bash
# 安装内核头文件
sudo apt-get install linux-headers-$(uname -r)

# 或者指定头文件路径
go generate -I/usr/src/linux-headers-$(uname -r)/include
```

**2. 运行时错误：Operation not permitted**

```bash
# 以 root 运行
sudo ./your-program

# 或添加 capability
sudo setcap cap_bpf,cap_net_admin,cap_perfmon+ep ./your-program
```

**3. 无法附加 kprobe**

可能是内核函数名称变化或内核版本不支持：

```bash
# 查看可用的 kprobe
sudo cat /sys/kernel/debug/tracing/available_filter_functions | grep tcp
```

**4. perf buffer 读取不到数据**

检查：
- eBPF 程序是否正确加载
- 是否有事件产生
- buffer 大小是否合适

---

## 📚 参考资源

- [Cilium eBPF 官方文档](https://ebpf-go.dev/)
- [eBPF 入门教程](https://ebpf.io/what-is-ebpf/)
- [BPF CO-RE](https://nakryiko.com/posts/bpf-portability-and-co-re/)
- [Linux Kernel eBPF Documentation](https://www.kernel.org/doc/html/latest/bpf/)

---

## 🔄 更新日志

### v2.0 (2025-12-03)

- ✅ 完善系统调用追踪器实现
- ✅ 完善网络监控实现
- ✅ 实现 perf buffer 读取
- ✅ 实现 eBPF map 读取
- ✅ 添加优雅启停机制
- ✅ 添加平台兼容性支持（Linux/非 Linux）
- ✅ 完善 C 程序地址获取实现

---

**状态**: ✅ 已完成
**目标**: 使用 Cilium eBPF 实现真正的系统级监控
