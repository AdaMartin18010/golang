# TS-033-Linux-Kernel-6-15-Features

> **Dimension**: 04-Technology-Stack
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: Linux Kernel 6.15 (May 2025)
> **Size**: >20KB

---

## 1. Linux 6.15 概览

### 1.1 发布信息

- **发布日期**: 2025年5月25日
- **代号**: "Timbernetes" (与K8s 1.35共享代号)
- **开发周期**: 2025年3月-5月
- **维护者**: Linus Torvalds / Greg Kroah-Hartman

### 1.2 主要特性分类

| 类别 | 特性数量 | 关键改进 |
|------|---------|---------|
| VFS/文件系统 | 5+ | Mount通知、idmapped mounts |
| 内存管理 | 3+ | AMD INVLPGB、HMM改进 |
| 硬件支持 | 10+ | Apple Touch Bar、AMD Zen 6 |
| 安全 | 4+ | Landlock审计、内存映射密封 |
| 性能 | 5+ | 调度器优化、CRC加速 |

---

## 2. VFS (Virtual File System) 改进

### 2.1 Mount Notifications

**功能**: 无需轮询/proc/mountinfo即可监听挂载拓扑变化

**API**:

```c
// 使用fanotify监听挂载事件
int fanotify_init(unsigned int flags, unsigned int event_f_flags);
int fanotify_mark(int fanotify_fd, unsigned int flags,
                  uint64_t mask, int dirfd, const char *pathname);

// 新事件类型
FANOTIFY_EV_MOUNT    // 挂载事件
FANOTIFY_EV_UMOUNT   // 卸载事件
FANOTIFY_EV_MOVE     // 移动挂载点
```

**使用场景**:

- 容器运行时监控挂载变化
- 文件系统管理工具
- 安全审计

**Go示例**:

```go
package main

import (
    "fmt"
    "os"
    "syscall"
    "golang.org/x/sys/unix"
)

func monitorMounts() error {
    // 创建fanotify实例
    fd, err := unix.FanotifyInit(unix.FAN_CLASS_NOTIF, 0)
    if err != nil {
        return err
    }
    defer unix.Close(fd)

    // 注册挂载事件监听
    err = unix.FanotifyMark(fd, unix.FAN_MARK_ADD|unix.FAN_MARK_MOUNT,
        unix.FANOTIFY_EV_MOUNT|unix.FANOTIFY_EV_UMOUNT,
        unix.AT_FDCWD, "/")
    if err != nil {
        return err
    }

    // 读取事件
    buf := make([]byte, 4096)
    for {
        n, err := unix.Read(fd, buf)
        if err != nil {
            return err
        }

        // 解析事件
        event := (*unix.FanotifyEventMetadata)(unsafe.Pointer(&buf[0]))
        fmt.Printf("Mount event: fd=%d, pid=%d\n", event.Fd, event.Pid)
    }
}
```

### 2.2 Idmapped Mounts 改进

**open_tree_attr() 系统调用**:

```c
// 从已有idmapped mount创建新的idmapped mount
int open_tree_attr(int dirfd, const char *pathname, unsigned int flags,
                   struct mount_attr *attr, size_t usize);
```

**优势**:

- 支持嵌套idmapped mounts
- 更灵活的ID映射管理
- 容器场景优化

### 2.3 Detached Mounts 增强

**新能力**:

1. 从detached mount创建detached mount
2. 在detached mounts上挂载其他detached mounts
3. OverlayFS支持detached mounts

**容器用例**:

```bash
# 创建私有rootfs而无需暴露整个文件系统
unshare -m bash -c '
    # 创建detached mount
    mount --make-private /
    mkdir -p /tmp/newroot
    mount -t tmpfs tmpfs /tmp/newroot

    # 组装新的rootfs
    mount --bind /bin /tmp/newroot/bin
    mount --bind /lib /tmp/newroot/lib
    # ... 其他挂载

    # 切换到新的rootfs
    pivot_root /tmp/newroot /tmp/newroot/oldroot
'
```

---

## 3. 性能优化

### 3.1 AMD INVLPGB 支持

**功能**: 广播TLB失效指令

**技术细节**:

- 支持CPU: AMD Zen 3+
- 功能: 无需IPI即可使远程CPU的TLB条目失效
- 优势: 减少中断、提高性能

**性能提升**:

```
TLB Shootdown 延迟:
- 传统方式 (IPI): ~2-5μs
- INVLPGB: ~0.5-1μs (60-80%↓)

多核扩展性:
- 64核系统: 显著降低开销
- 内存密集型应用受益
```

### 3.2 CRC 加速 (Intel/AMD AVX-512)

**优化**:

- 使用AVX-512指令加速CRC计算
- 受益场景: 网络数据包校验、存储完整性检查

**基准测试**:

```
CRC32 1GB数据:
- 标量: 1.2 GB/s
- AVX-512: 4.8 GB/s (4x faster)
```

### 3.3 调度器改进

**新默认空闲CPU选择逻辑**:

- 更好的负载均衡
- 降低调度延迟
- 提高吞吐量

---

## 4. 硬件支持

### 4.1 Apple Touch Bar 支持

**驱动**: Apple Touch Bar Backlight & Keyboard Mode

- 背光控制
- 键盘模式切换
- 显示控制

### 4.2 AMD Zen 6 CPU 识别

**支持**:

- CPUID识别
- 性能监控单元(PMU)
- 电源管理

### 4.3 新平台支持

| 平台 | 支持状态 | 说明 |
|------|---------|------|
| Google Pixel Pro 6 | ✓ | 智能手机 |
| Huawei Matebook E Go | ✓ | ARM笔记本 |
| Milk-V Jupiter RISC-V | ✓ | 开发板 |
| Snapdragon X1 | ✓ | ASUS Zenbook A14 |

---

## 5. 安全增强

### 5.1 Landlock 审计机制

**功能**: 更细粒度的访问审计

```c
// 新的Landlock规则类型
LANDLOCK_RULE_ACCESS_FS    // 文件系统访问
LANDLOCK_RULE_ACCESS_NET   // 网络访问
```

**审计日志**:

```
landlock: access denied { write } for pid=1234
          path=/etc/passwd
          rule=custom_policy_1
```

### 5.2 内存映射密封 (Memory Mapping Sealing)

**功能**: 密封内存映射防止修改

```c
// 新的mmap标志
MAP_SEAL    // 创建密封映射

// 密封操作
mprotect_seal(addr, len, PROT_NONE);  // 禁止所有访问
mremap_seal(...);  // 禁止重新映射
```

**安全用例**:

- 防止恶意代码修改敏感数据
- 强化沙箱环境

### 5.3 fwctl 子系统

**功能**: 用户空间与设备固件安全通信

```c
// 创建固件控制通道
int fwctl_open(const char *device);
int fwctl_ioctl(int fd, unsigned int cmd, void *arg);
```

---

## 6. 文件系统改进

### 6.1 Bcachefs 稳定化

**进展**:

- On-disk format冻结
- 自我修复能力提升
- 快照删除性能优化

### 6.2 XFS Zoned Storage 支持

**功能**: 支持分区存储设备

```bash
# 创建XFS on zoned device
mkfs.xfs -d zonesize=256M /dev/sdz
```

### 6.3 F2FS 增强

- 更好的垃圾回收
- 性能优化
- 移动设备优化

---

## 7. 网络改进

### 7.1 调度器事件计数 (sched_ext)

```c
// 使用BPF程序计数调度事件
SEC("tp/sched/sched_switch")
int trace_sched_switch(void *ctx) {
    // 统计上下文切换
    return 0;
}
```

### 7.2 网络性能优化

- TCP性能改进
- 网络栈优化
- 减少延迟

---

## 8. eBPF 增强

### 8.1 新特性

| 特性 | 描述 | 内核版本 |
|------|------|---------|
| BPF tokens | 安全非特权使用 | 6.15+ |
| BPF arenas | 共享内存 | 6.15+ |
| BPF exceptions | 异常处理 | 6.15+ |

### 8.2 调度器支持

**sched_ext**:

- 可扩展调度器类
- 使用BPF编写调度策略
- 动态加载/卸载

---

## 9. Rust 集成

### 9.1 新Rust抽象

- hrtimer支持
- ARMv7支持
- 更多驱动程序用Rust编写

### 9.2 内核单元测试

```rust
// 内核Rust单元测试宏
#[kernel_test]
fn test_my_module() {
    assert_eq!(my_function(), expected_value);
}
```

---

## 10. 其他重要变化

### 10.1 移除32位x86大内存支持

**变化**:

- 不再支持>8 CPUs的32位x86
- 不再支持>4GB RAM的32位x86

**影响**: 推动向64位迁移

### 10.2 Python 3.9+ 要求

**变化**: 内核构建和文档需要Python 3.9+

### 10.3 Perf 子系统

**新功能**:

- 延迟分析能力
- 更精确的采样

---

## 11. 升级指南

### 11.1 从 6.14 升级到 6.15

```bash
# 下载内核
curl -O https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.15.tar.xz

# 解压
tar -xf linux-6.15.tar.xz
cd linux-6.15

# 配置
make menuconfig

# 编译
make -j$(nproc)

# 安装
sudo make modules_install
sudo make install
```

### 11.2 重要配置变更

```
CONFIG_MOUNT_NOTIFICATIONS=y
CONFIG_AMD_INVLPGB=y
CONFIG_LANDLOCK_AUDIT=y
CONFIG_FWCTL=y
```

---

## 12. 性能基准

### 12.1 编译时间

```
内核编译 (make -j32):
- 6.14: 4分32秒
- 6.15: 4分28秒 (1.5% faster)
```

### 12.2 调度性能

```
调度延迟 P99:
- 6.14: 45μs
- 6.15: 38μs (15% improvement)
```

---

## 13. 参考文献

1. Linux 6.15 Release Notes (kernel.org)
2. LWN: Linux 6.15 Merge Window
3. Phoronix Linux 6.15 Benchmarks
4. KernelNewbies: Linux 6.15
5. AMD INVLPGB Whitepaper

---

*Last Updated: 2026-04-03*
