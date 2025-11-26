# 1. 📊 eBPF 深度解析

> **简介**: 本文档详细阐述了 eBPF 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 📊 eBPF 深度解析](#1--ebpf-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 eBPF 程序编写](#131-ebpf-程序编写)
    - [1.3.2 使用 cilium/ebpf 加载程序](#132-使用-ciliumebpf-加载程序)
    - [1.3.3 系统调用追踪](#133-系统调用追踪)
    - [1.3.4 网络监控](#134-网络监控)
    - [1.3.5 性能分析](#135-性能分析)
    - [1.3.6 与 OpenTelemetry 集成](#136-与-opentelemetry-集成)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 eBPF 程序设计最佳实践](#141-ebpf-程序设计最佳实践)
    - [1.4.2 性能优化最佳实践](#142-性能优化最佳实践)
    - [1.4.3 安全最佳实践](#143-安全最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**eBPF 是什么？**

eBPF (extended Berkeley Packet Filter) 是一个用于在 Linux 内核中运行沙箱程序的技术，允许在内核空间安全地执行用户定义的代码。

**核心特性**:

- ✅ **内核空间执行**: 在内核中运行，性能优秀
- ✅ **安全性**: 通过验证器确保程序安全
- ✅ **动态加载**: 无需重启内核即可加载程序
- ✅ **低开销**: 高效的事件处理机制
- ✅ **可编程性**: 支持复杂的过滤和处理逻辑

---

## 1.2 选型论证

**为什么选择 eBPF？**

**论证矩阵**:

| 评估维度 | 权重 | eBPF | ptrace | perf | SystemTap | 说明 |
|---------|------|------|--------|------|-----------|------|
| **性能** | 30% | 10 | 5 | 8 | 7 | eBPF 内核执行，性能最优 |
| **安全性** | 25% | 10 | 6 | 8 | 7 | eBPF 验证器保证安全 |
| **灵活性** | 20% | 10 | 7 | 6 | 9 | eBPF 可编程性强 |
| **易用性** | 15% | 7 | 6 | 8 | 6 | eBPF 学习曲线适中 |
| **生态支持** | 10% | 9 | 7 | 9 | 7 | eBPF 生态活跃 |
| **加权总分** | - | **9.30** | 6.20 | 7.80 | 7.30 | eBPF 得分最高 |

**核心优势**:

1. **性能（权重 30%）**:
   - 在内核空间执行，避免用户态-内核态切换
   - 低开销的事件处理
   - 适合高频事件监控

2. **安全性（权重 25%）**:
   - 通过验证器确保程序安全
   - 防止无限循环和内存访问错误
   - 沙箱执行环境

3. **灵活性（权重 20%）**:
   - 支持复杂的过滤和处理逻辑
   - 可以访问内核数据结构
   - 支持多种事件类型

**为什么不选择其他技术？**

1. **ptrace**:
   - ✅ 功能强大，可以追踪进程
   - ❌ 性能开销大
   - ❌ 需要停止目标进程
   - ❌ 不适合生产环境

2. **perf**:
   - ✅ 性能优秀，功能丰富
   - ❌ 灵活性不如 eBPF
   - ❌ 需要 root 权限
   - ❌ 配置复杂

3. **SystemTap**:
   - ✅ 功能强大，脚本灵活
   - ❌ 需要编译内核模块
   - ❌ 稳定性不如 eBPF
   - ❌ 学习曲线陡峭

**适用场景**:

- ✅ 系统调用追踪
- ✅ 网络包过滤和监控
- ✅ 性能分析和优化
- ✅ 安全监控和审计
- ✅ 实时指标收集
- ✅ 故障诊断和调试

**不适用场景**:

- ❌ Windows 系统（仅支持 Linux）
- ❌ 需要修改内核逻辑的场景
- ❌ 需要访问所有内核数据的场景

---

## 1.3 实际应用

### 1.3.1 eBPF 程序编写

**基础 eBPF 程序示例**:

```c
// internal/infrastructure/observability/ebpf/programs/trace_syscall.bpf.c
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

// 定义 map 用于存储数据
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, u32);
    __type(value, u64);
} syscall_count SEC(".maps");

// 追踪系统调用
SEC("tracepoint/syscalls/sys_enter_openat")
int trace_syscall_openat(struct trace_event_raw_sys_enter *ctx) {
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    u64 *count = bpf_map_lookup_elem(&syscall_count, &pid);

    if (count) {
        (*count)++;
    } else {
        u64 init = 1;
        bpf_map_update_elem(&syscall_count, &pid, &init, BPF_ANY);
    }

    return 0;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
```

### 1.3.2 使用 cilium/ebpf 加载程序

**加载 eBPF 程序**:

```go
// internal/infrastructure/observability/ebpf/collector.go
package ebpf

import (
    "context"
    "fmt"
    "os"

    "github.com/cilium/ebpf"
    "github.com/cilium/ebpf/link"
    "github.com/cilium/ebpf/rlimit"
)

type Collector struct {
    collection *ebpf.Collection
    links      []link.Link
}

func NewCollector() (*Collector, error) {
    // 移除内存限制
    if err := rlimit.RemoveMemlock(); err != nil {
        return nil, fmt.Errorf("failed to remove memlock: %w", err)
    }

    // 加载编译后的 eBPF 程序
    spec, err := ebpf.LoadCollectionSpec("trace_syscall.bpf.o")
    if err != nil {
        return nil, fmt.Errorf("failed to load collection spec: %w", err)
    }

    // 创建 collection
    collection, err := ebpf.NewCollection(spec)
    if err != nil {
        return nil, fmt.Errorf("failed to create collection: %w", err)
    }

    return &Collector{
        collection: collection,
        links:      make([]link.Link, 0),
    }, nil
}

func (c *Collector) AttachTracepoint(tpName string) error {
    // 附加到 tracepoint
    tp, err := link.OpenTracepoint(link.TracepointOptions{
        Tracepoint: tpName,
        Program:    c.collection.Programs["trace_syscall_openat"],
    })
    if err != nil {
        return fmt.Errorf("failed to attach tracepoint: %w", err)
    }

    c.links = append(c.links, tp)
    return nil
}

func (c *Collector) Close() error {
    // 关闭所有链接
    for _, l := range c.links {
        l.Close()
    }

    // 关闭 collection
    return c.collection.Close()
}
```

### 1.3.3 系统调用追踪

**追踪系统调用示例**:

```go
// 追踪系统调用
func (c *Collector) TraceSyscalls(ctx context.Context) error {
    // 附加到系统调用 tracepoint
    syscalls := []string{
        "sys_enter_openat",
        "sys_enter_read",
        "sys_enter_write",
        "sys_enter_connect",
    }

    for _, syscall := range syscalls {
        if err := c.AttachTracepoint(syscall); err != nil {
            return fmt.Errorf("failed to attach %s: %w", syscall, err)
        }
    }

    return nil
}

// 读取统计数据
func (c *Collector) GetSyscallStats() (map[uint32]uint64, error) {
    syscallCount := c.collection.Maps["syscall_count"]

    stats := make(map[uint32]uint64)
    var key uint32
    var value uint64

    iter := syscallCount.Iterate()
    for iter.Next(&key, &value) {
        stats[key] = value
    }

    if err := iter.Err(); err != nil {
        return nil, fmt.Errorf("failed to iterate map: %w", err)
    }

    return stats, nil
}
```

### 1.3.4 网络监控

**网络包监控示例**:

```c
// network_monitor.bpf.c
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, u32);
    __type(value, u64);
} packet_count SEC(".maps");

SEC("xdp")
int xdp_prog(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) {
        return XDP_PASS;
    }

    if (eth->h_proto != __constant_htons(ETH_P_IP)) {
        return XDP_PASS;
    }

    struct iphdr *ip = (struct iphdr *)(eth + 1);
    if ((void *)(ip + 1) > data_end) {
        return XDP_PASS;
    }

    u32 protocol = ip->protocol;
    u64 *count = bpf_map_lookup_elem(&packet_count, &protocol);

    if (count) {
        (*count)++;
    } else {
        u64 init = 1;
        bpf_map_update_elem(&packet_count, &protocol, &init, BPF_ANY);
    }

    return XDP_PASS;
}
```

### 1.3.5 性能分析

**CPU 性能分析示例**:

```go
// CPU 性能分析
func (c *Collector) ProfileCPU(ctx context.Context, duration time.Duration) error {
    // 附加到 perf event
    pe, err := link.OpenPerfEvent(link.PerfEventOptions{
        Fd:        -1, // CPU
        PerfType:  unix.PERF_TYPE_SOFTWARE,
        Config:    unix.PERF_COUNT_SW_CPU_CLOCK,
        SampleFreq: 100, // 100 Hz
    })
    if err != nil {
        return fmt.Errorf("failed to open perf event: %w", err)
    }
    defer pe.Close()

    // 读取性能数据
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            // 读取性能数据
            stats, err := c.GetPerformanceStats()
            if err != nil {
                return err
            }
            // 处理统计数据
            c.processStats(stats)
        }
    }
}
```

### 1.3.6 与 OpenTelemetry 集成

**集成 OpenTelemetry**:

```go
// 与 OpenTelemetry 集成
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
)

func (c *Collector) ExportToOpenTelemetry(meter metric.Meter) error {
    // 创建指标
    syscallCounter, err := meter.Int64Counter(
        "ebpf.syscall.count",
        metric.WithDescription("Number of system calls"),
    )
    if err != nil {
        return err
    }

    // 定期导出数据
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            stats, err := c.GetSyscallStats()
            if err != nil {
                continue
            }

            for pid, count := range stats {
                syscallCounter.Add(context.Background(), int64(count),
                    metric.WithAttributes(
                        attribute.Int("pid", int(pid)),
                    ),
                )
            }
        }
    }()

    return nil
}
```

---

## 1.4 最佳实践

### 1.4.1 eBPF 程序设计最佳实践

**为什么需要良好的 eBPF 程序设计？**

良好的 eBPF 程序设计可以提高程序的可维护性、安全性和性能。

**程序设计原则**:

1. **简化逻辑**: 保持程序逻辑简单，避免复杂计算
2. **内存安全**: 始终检查边界，避免内存访问错误
3. **错误处理**: 正确处理错误情况
4. **性能优化**: 减少 map 查找和更新次数

**实际应用示例**:

```c
// eBPF 程序设计最佳实践
SEC("tracepoint/syscalls/sys_enter_openat")
int trace_syscall_openat_safe(struct trace_event_raw_sys_enter *ctx) {
    // 1. 边界检查
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    if (data + sizeof(struct trace_event_raw_sys_enter) > data_end) {
        return 0; // 安全返回
    }

    // 2. 获取进程 ID
    u32 pid = bpf_get_current_pid_tgid() >> 32;

    // 3. 查找 map（带错误处理）
    u64 *count = bpf_map_lookup_elem(&syscall_count, &pid);
    if (count) {
        // 4. 原子更新
        __sync_fetch_and_add(count, 1);
    } else {
        // 5. 初始化新条目
        u64 init = 1;
        bpf_map_update_elem(&syscall_count, &pid, &init, BPF_NOEXIST);
    }

    return 0;
}
```

**最佳实践要点**:

1. **边界检查**: 始终检查数据边界，避免越界访问
2. **错误处理**: 正确处理 map 查找失败等情况
3. **原子操作**: 使用原子操作更新共享数据
4. **简化逻辑**: 保持程序逻辑简单，避免复杂计算
5. **性能优化**: 减少 map 操作，使用局部变量

### 1.4.2 性能优化最佳实践

**为什么需要性能优化？**

eBPF 程序在内核中执行，性能优化可以减少系统开销。

**性能优化原则**:

1. **减少 map 操作**: 最小化 map 查找和更新
2. **使用局部变量**: 减少内存访问
3. **避免循环**: 避免复杂循环，保持程序简单
4. **批量处理**: 批量处理事件，减少开销
5. **采样**: 对高频事件进行采样

**实际应用示例**:

```c
// 性能优化最佳实践
SEC("tracepoint/syscalls/sys_enter_openat")
int trace_syscall_openat_optimized(struct trace_event_raw_sys_enter *ctx) {
    // 1. 使用局部变量
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    u64 count = 1;

    // 2. 采样（每 100 个事件处理一次）
    if (pid % 100 != 0) {
        return 0;
    }

    // 3. 批量更新（使用 per-CPU map）
    u64 *count_ptr = bpf_map_lookup_elem(&per_cpu_count, &pid);
    if (count_ptr) {
        *count_ptr += count;
    } else {
        bpf_map_update_elem(&per_cpu_count, &pid, &count, BPF_ANY);
    }

    return 0;
}
```

**最佳实践要点**:

1. **采样**: 对高频事件进行采样，减少处理开销
2. **per-CPU map**: 使用 per-CPU map 避免锁竞争
3. **局部变量**: 使用局部变量减少内存访问
4. **批量处理**: 批量处理事件，减少系统调用
5. **简化逻辑**: 保持程序逻辑简单，提高执行效率

### 1.4.3 安全最佳实践

**为什么需要安全最佳实践？**

eBPF 程序在内核中执行，安全问题可能导致系统崩溃或安全漏洞。

**安全最佳实践**:

1. **验证器检查**: 确保程序通过验证器检查
2. **边界检查**: 始终检查数据边界
3. **权限控制**: 限制 eBPF 程序的使用权限
4. **代码审查**: 仔细审查 eBPF 程序代码
5. **测试**: 充分测试 eBPF 程序

**实际应用示例**:

```go
// 安全最佳实践
func (c *Collector) LoadProgramSafely(specPath string) error {
    // 1. 验证文件权限
    info, err := os.Stat(specPath)
    if err != nil {
        return fmt.Errorf("failed to stat file: %w", err)
    }

    if info.Mode().Perm()&0077 != 0 {
        return fmt.Errorf("file has insecure permissions")
    }

    // 2. 加载并验证程序
    spec, err := ebpf.LoadCollectionSpec(specPath)
    if err != nil {
        return fmt.Errorf("failed to load spec: %w", err)
    }

    // 3. 验证程序大小
    for name, prog := range spec.Programs {
        if len(prog.Instructions) > 1000000 {
            return fmt.Errorf("program %s too large", name)
        }
    }

    // 4. 创建 collection（验证器会自动检查）
    collection, err := ebpf.NewCollection(spec)
    if err != nil {
        return fmt.Errorf("failed to create collection (verifier error): %w", err)
    }

    c.collection = collection
    return nil
}
```

**最佳实践要点**:

1. **验证器**: 依赖 eBPF 验证器确保程序安全
2. **边界检查**: 在 eBPF 程序中始终检查边界
3. **权限控制**: 限制 eBPF 程序的使用权限
4. **代码审查**: 仔细审查 eBPF 程序代码
5. **测试**: 在测试环境中充分测试程序

---

## 📚 扩展阅读

- [eBPF 官方文档](https://ebpf.io/)
- [cilium/ebpf 官方文档](https://github.com/cilium/ebpf)
- [eBPF 和 Go](https://github.com/cilium/ebpf)
- [OpenTelemetry eBPF](https://opentelemetry.io/docs/instrumentation/ebpf/)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 eBPF 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
