# 系统监控

系统资源监控包，提供 CPU、内存、IO、网络等系统资源的监控，以及容器、操作系统、虚拟化环境的检测。

## 📋 功能特性

- ✅ **CPU 监控**: CPU 使用率监控
- ✅ **内存监控**: 内存使用量、GC 统计
- ✅ **IO 监控**: 读写字节数、操作数
- ✅ **网络监控**: 网络流量、连接数
- ✅ **平台检测**: 操作系统、架构、Go 版本
- ✅ **容器检测**: Docker、Kubernetes 环境检测
- ✅ **虚拟化检测**: 检测虚拟化环境（VMware、KVM、AWS、GCP、Azure 等）
- ✅ **OTLP 集成**: 自动导出指标到 OpenTelemetry

## 🚀 快速开始

### 基本使用

```go
import (
    "github.com/yourusername/golang/pkg/observability/system"
    "github.com/yourusername/golang/pkg/observability/otlp"
)

// 初始化 OTLP
otlpClient, _ := otlp.NewEnhancedOTLP(otlp.Config{
    ServiceName:    "my-service",
    ServiceVersion: "v1.0.0",
    Endpoint:       "localhost:4317",
    Insecure:       true,
})

// 创建系统监控器
systemMonitor, err := system.NewSystemMonitor(system.SystemConfig{
    Meter:           otlpClient.Meter("system"),
    Enabled:         true,
    CollectInterval: 5 * time.Second,
})
if err != nil {
    log.Fatal(err)
}

// 启动监控
if err := systemMonitor.Start(); err != nil {
    log.Fatal(err)
}
defer systemMonitor.Stop()
```

### 获取平台信息

```go
// 获取平台信息
info := systemMonitor.GetPlatformInfo()
fmt.Printf("OS: %s\n", info.OS)
fmt.Printf("Arch: %s\n", info.Arch)
fmt.Printf("Container: %s\n", info.ContainerID)
fmt.Printf("K8s Pod: %s\n", info.KubernetesPod)
fmt.Printf("Virtualization: %s\n", info.Virtualization)

// 检查环境
if systemMonitor.IsContainer() {
    fmt.Println("Running in container")
}
if systemMonitor.IsKubernetes() {
    fmt.Println("Running in Kubernetes")
}
if systemMonitor.IsVirtualized() {
    fmt.Println("Running in virtualized environment")
}
```

### 获取资源统计

```go
// 获取内存统计
memStats := systemMonitor.GetMemoryStats()
fmt.Printf("Memory Alloc: %d bytes\n", memStats.Alloc)
fmt.Printf("GC Count: %d\n", memStats.NumGC)

// 获取 Goroutine 数量
goroutines := systemMonitor.GetGoroutineCount()
fmt.Printf("Goroutines: %d\n", goroutines)
```

## 📊 导出的指标

### 系统资源指标

- `system.cpu.usage` - CPU 使用率（%）
- `system.memory.usage` - 内存使用量（字节）
- `system.memory.total` - 总内存（字节）
- `system.goroutines` - Goroutine 数量
- `system.gc.count` - GC 次数
- `system.gc.duration` - GC 持续时间（秒）

### IO 指标

- `system.io.read.bytes` - 读取字节数
- `system.io.write.bytes` - 写入字节数
- `system.io.read.ops` - 读取操作数
- `system.io.write.ops` - 写入操作数

### 网络指标

- `system.network.bytes.sent` - 发送字节数
- `system.network.bytes.received` - 接收字节数
- `system.network.packets.sent` - 发送包数
- `system.network.packets.received` - 接收包数
- `system.network.connections` - 连接数

### 平台信息指标

- `system.platform.info` - 平台信息（通过属性传递）

## 🔍 平台检测

### 容器检测

支持检测以下容器环境：
- Docker
- Kubernetes Pod
- systemd-nspawn
- LXC

### 虚拟化检测

支持检测以下虚拟化环境：
- VMware
- VirtualBox
- KVM/QEMU
- Xen
- AWS
- GCP
- Azure
- 裸机（bare-metal）

## 📚 API 参考

### SystemMonitor

```go
type SystemMonitor struct {
    // ...
}

func NewSystemMonitor(cfg SystemConfig) (*SystemMonitor, error)
func (sm *SystemMonitor) Start() error
func (sm *SystemMonitor) Stop() error
func (sm *SystemMonitor) GetPlatformInfo() PlatformInfo
func (sm *SystemMonitor) GetMemoryStats() MemoryStats
func (sm *SystemMonitor) GetGoroutineCount() int
func (sm *SystemMonitor) IsContainer() bool
func (sm *SystemMonitor) IsKubernetes() bool
func (sm *SystemMonitor) IsVirtualized() bool
```

### PlatformInfo

```go
type PlatformInfo struct {
    OS              string
    Arch            string
    GoVersion       string
    Hostname        string
    ContainerID     string
    ContainerName   string
    KubernetesPod   string
    KubernetesNode  string
    Virtualization  string
    CPUs            int
}
```

## ⚠️ 注意事项

1. **权限要求**: 某些系统信息可能需要特定权限
2. **平台限制**: 部分功能仅在 Linux 上可用
3. **性能影响**: 监控会消耗少量系统资源
4. **精度**: 某些指标（如 CPU 使用率）是近似值

## 🔗 相关文档

- [OTLP 集成](../otlp/README.md)
- [使用指南](../../../docs/usage-guide.md)
