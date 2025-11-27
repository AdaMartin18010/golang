# eBPF 收集器

框架级别的 eBPF 收集器，提供基于 eBPF 的系统级可观测性数据收集。

## 📋 功能特性

- ✅ **系统调用追踪**: 追踪系统调用
- ✅ **网络监控**: 监控网络包
- ✅ **性能分析**: 系统级性能分析
- ✅ **OpenTelemetry 集成**: 与 OTLP 集成

## 🚀 快速开始

### 基本使用

```go
import "github.com/yourusername/golang/pkg/observability/ebpf"

// 创建 eBPF 收集器
collector := ebpf.NewCollector(ebpf.Config{
    Tracer:  tracer,
    Meter:   meter,
    Enabled: true,
})

// 启动收集器
if err := collector.Start(); err != nil {
    log.Fatal(err)
}
defer collector.Stop()

// 收集指标
collector.CollectSyscallMetrics(ctx)
collector.CollectNetworkMetrics(ctx)
```

## ⚠️ 注意事项

实际的 eBPF 程序实现需要：

1. **编写 eBPF C 程序** (`.bpf.c` 文件)
2. **使用 cilium/ebpf 加载程序**
3. **从 eBPF map 读取数据**
4. **转换为 OpenTelemetry 指标和追踪**

## 📚 相关文档

- [eBPF 深度解析](../../../docs/architecture/tech-stack/observability/ebpf.md)
- [OpenTelemetry 集成](../otlp/README.md)

## 🔗 参考实现

实际的 eBPF 程序实现请参考：

- `internal/infrastructure/observability/ebpf/programs/`
- `docs/architecture/tech-stack/observability/ebpf.md`
