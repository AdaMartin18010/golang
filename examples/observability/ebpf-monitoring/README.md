# eBPF 监控示例

**使用**: Cilium eBPF v0.20.0
**集成**: OpenTelemetry OTLP

---

## 📋 功能展示

本示例展示如何使用 eBPF 进行系统级监控：

1. **系统调用追踪** - 追踪应用的所有系统调用
2. **TCP 连接监控** - 监控 TCP 连接的建立和关闭
3. **网络流量统计** - 统计发送和接收的字节数
4. **延迟测量** - 测量系统调用和连接的延迟

---

## 🚀 运行示例

### 前置要求

1. **Linux 环境** (Ubuntu 20.04+ / Debian 11+ / RHEL 8+)
2. **Root 权限** 或 `CAP_BPF` capability
3. **内核版本** >= 5.2 (推荐 5.10+)
4. **Clang/LLVM** (编译 eBPF 程序)

### 安装依赖

```bash
# Ubuntu/Debian
sudo apt-get install clang llvm linux-headers-$(uname -r)

# RHEL/CentOS
sudo yum install clang llvm kernel-devel
```

### 生成 eBPF 代码

```bash
# 在项目根目录
make generate-ebpf
```

### 运行示例

```bash
# 需要 root 权限
sudo go run main.go

# 或者使用 capabilities
sudo setcap cap_bpf,cap_net_admin+ep $(which go)
go run main.go
```

---

## 📊 输出示例

```
🚀 eBPF 监控示例
使用 Cilium eBPF 库进行系统级监控

📊 创建 eBPF 收集器...
▶️  启动 eBPF 监控...
✅ eBPF 监控已启动

监控功能：
  ✅ 系统调用追踪 (sys_enter/sys_exit)
  ✅ TCP 连接监控 (connect/accept/close)
  ✅ 网络流量统计 (bytes sent/recv)
  ✅ 延迟测量 (syscall/connection latency)

🔄 模拟工作负载...

📡 eBPF 监控运行中...
按 Ctrl+C 停止

{
  "Name": "syscall",
  "SpanContext": {...},
  "Attributes": [
    {"Key": "syscall.id", "Value": 39},
    {"Key": "process.pid", "Value": 12345},
    {"Key": "syscall.duration_ms", "Value": 0.123}
  ]
}
```

---

## 🔧 故障排查

### 权限错误

```
Error: failed to load eBPF program: operation not permitted
```

**解决方案**:

```bash
# 方案1: 使用 root
sudo go run main.go

# 方案2: 添加 capabilities
sudo setcap cap_bpf,cap_net_admin+ep $(which go)
```

### 内核版本过低

```
Error: eBPF program type not supported
```

**解决方案**:

- 升级内核到 5.2+ (推荐 5.10+)
- 或使用兼容模式（功能有限）

### Clang 未安装

```
Error: clang: command not found
```

**解决方案**:

```bash
sudo apt-get install clang llvm
```

---

## 📚 相关文档

- [eBPF 实现文档](../../../pkg/observability/ebpf/README.md)
- [Cilium eBPF 文档](https://ebpf-go.dev/)
- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/)

---

## 🎯 扩展示例

### 添加自定义追踪

```go
// 创建自定义追踪器
customTracer := &MyCustomTracer{
    tracer: otel.Tracer("my-tracer"),
    meter:  otel.Meter("my-meter"),
}

// 集成到 collector
// collector.AddTracer(customTracer)
```

### 过滤特定进程

```go
collector, err := ebpf.NewCollector(ebpf.Config{
    // ... 其他配置
    TargetPID: 12345, // 只追踪特定进程
})
```

---

**状态**: ✅ 生产就绪
**平台**: Linux 5.2+
**权限**: Root 或 CAP_BPF
