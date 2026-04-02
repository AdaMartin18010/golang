# 内核级任务调度 (Kernel-Level Task Scheduling)

> **分类**: 工程与云原生
> **标签**: #kernel #syscall #epoll #io-uring
> **参考**: Linux Kernel, epoll, io_uring

---

## 目录

- [内核级任务调度 (Kernel-Level Task Scheduling)](#内核级任务调度-kernel-level-task-scheduling)
  - [目录](#目录)
  - [内核调度架构](#内核调度架构)
  - [epoll 实现](#epoll-实现)
  - [io\_uring 实现](#io_uring-实现)
  - [futex 实现](#futex-实现)
  - [内核调度策略](#内核调度策略)
  - [性能对比](#性能对比)
  - [完整调度器集成](#完整调度器集成)

## 内核调度架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Kernel-Level Scheduling Architecture                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  User Space                    Kernel Space                    Hardware     │
│  ───────────                   ───────────                     ─────────    │
│                                                                              │
│  ┌─────────────┐               ┌─────────────┐               ┌─────────┐   │
│  │   Task      │  syscall      │   epoll     │               │  NIC    │   │
│  │   Queue     │──────────────►│   Wait      │◄─────────────►│  IRQ    │   │
│  └─────────────┘               └─────────────┘               └─────────┘   │
│         │                             │                                    │
│         │      futex/                 │      io_uring                        │
│         ▼      fcntl                  ▼      submit                          │
│  ┌─────────────┐               ┌─────────────┐               ┌─────────┐   │
│  │  Worker     │               │  io_uring   │◄─────────────►│  Disk   │   │
│  │  Pool       │               │  Queue      │   DMA         │  I/O    │   │
│  └─────────────┘               └─────────────┘               └─────────┘   │
│                                                                              │
│  System Calls: epoll_create, epoll_ctl, epoll_wait                          │
│  io_uring: io_uring_setup, io_uring_enter, io_uring_submit                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## epoll 实现

```go
package kernelsched

import (
    "sync"
    "syscall"
    "unsafe"
)

// EpollPoller epoll 轮询器
type EpollPoller struct {
    epfd   int
    events []syscall.EpollEvent
    mu     sync.Mutex
}

// NewEpollPoller 创建 epoll 轮询器
func NewEpollPoller() (*EpollPoller, error) {
    // epoll_create1(EPOLL_CLOEXEC)
    epfd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
    if err != nil {
        return nil, err
    }

    return &EpollPoller{
        epfd:   epfd,
        events: make([]syscall.EpollEvent, 1024),
    }, nil
}

// Add 添加文件描述符
func (ep *EpollPoller) Add(fd int, events uint32) error {
    ep.mu.Lock()
    defer ep.mu.Unlock()

    event := syscall.EpollEvent{
        Events: events,
        Fd:     int32(fd),
    }

    // EPOLL_CTL_ADD
    return syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_ADD, fd, &event)
}

// Modify 修改事件
func (ep *EpollPoller) Modify(fd int, events uint32) error {
    event := syscall.EpollEvent{
        Events: events,
        Fd:     int32(fd),
    }

    return syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_MOD, fd, &event)
}

// Remove 移除文件描述符
func (ep *EpollPoller) Remove(fd int) error {
    return syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_DEL, fd, nil)
}

// Wait 等待事件
func (ep *EpollPoller) Wait(timeout int) (int, error) {
    // epoll_wait(epfd, events, maxevents, timeout)
    n, err := syscall.EpollWait(epfd, ep.events, timeout)
    if err != nil {
        return 0, err
    }
    return n, nil
}

// Close 关闭
func (ep *EpollPoller) Close() error {
    return syscall.Close(ep.epfd)
}

// Event types
const (
    EPOLLIN    = syscall.EPOLLIN
    EPOLLOUT   = syscall.EPOLLOUT
    EPOLLERR   = syscall.EPOLLERR
    EPOLLHUP   = syscall.EPOLLHUP
    EPOLLET    = syscall.EPOLLET    // 边缘触发
    EPOLLONESHOT = syscall.EPOLLONESHOT
)
```

---

## io_uring 实现

```go
// IoUring io_uring 接口（Linux 5.1+）
type IoUring struct {
    ringFd   int
    sq       *SubmissionQueue
    cq       *CompletionQueue

    sqRing   []uint32
    cqRing   []uint32
}

// SubmissionQueue 提交队列
type SubmissionQueue struct {
    head    *uint32
    tail    *uint32
    mask    uint32
    entries []*IoUringSqe
}

// IoUringSqe 提交队列条目
type IoUringSqe struct {
    Opcode   uint8
    Flags    uint8
    Ioprio   uint16
    Fd       int32
    Off      uint64
    Addr     uint64
    Len      uint32
    RwFlags  uint32
    UserData uint64
}

// Setup 设置 io_uring
func Setup(entries uint32, params *IoUringParams) (*IoUring, error) {
    // syscall io_uring_setup(entries, params)
    ringFd, _, errno := syscall.Syscall(
        SYS_IO_URING_SETUP,
        uintptr(entries),
        uintptr(unsafe.Pointer(params)),
        0,
    )

    if errno != 0 {
        return nil, errno
    }

    ring := &IoUring{
        ringFd: int(ringFd),
    }

    // 内存映射队列
    if err := ring.mmapQueues(params); err != nil {
        syscall.Close(int(ringFd))
        return nil, err
    }

    return ring, nil
}

// SubmitRead 提交读请求
func (ring *IoUring) SubmitRead(fd int, buf []byte, offset uint64) error {
    sqe := ring.GetSqe()
    if sqe == nil {
        return fmt.Errorf("submission queue full")
    }

    sqe.Opcode = IORING_OP_READ
    sqe.Fd = int32(fd)
    sqe.Addr = uint64(uintptr(unsafe.Pointer(&buf[0])))
    sqe.Len = uint32(len(buf))
    sqe.Off = offset

    return ring.Submit()
}

// Submit 提交到内核
func (ring *IoUring) Submit() error {
    // io_uring_enter(ringFd, toSubmit, minComplete, flags, sigset)
    submitted := ring.sq.Flush()

    _, _, errno := syscall.Syscall6(
        SYS_IO_URING_ENTER,
        uintptr(ring.ringFd),
        uintptr(submitted),
        0,
        0,
        0,
        0,
    )

    if errno != 0 {
        return errno
    }

    return nil
}

// WaitCqe 等待完成事件
func (ring *IoUring) WaitCqe() (*IoUringCqe, error) {
    // 轮询或等待完成队列
    for {
        cqe := ring.cq.Peek()
        if cqe != nil {
            ring.cq.Advance()
            return cqe, nil
        }

        // 提交并等待
        ring.SubmitAndWait(1)
    }
}

// IoUringCqe 完成队列条目
type IoUringCqe struct {
    UserData uint64
    Res      int32
    Flags    uint32
}

// 操作码
const (
    IORING_OP_NOP        = 0
    IORING_OP_READV      = 1
    IORING_OP_WRITEV     = 2
    IORING_OP_FSYNC      = 3
    IORING_OP_READ_FIXED = 4
    IORING_OP_POLL_ADD   = 6
    IORING_OP_POLL_REMOVE = 7
    IORING_OP_CONNECT    = 16
    IORING_OP_ACCEPT     = 13
    IORING_OP_ASYNC_CANCEL = 14
)

// Syscall numbers
const (
    SYS_IO_URING_SETUP = 425
    SYS_IO_URING_ENTER = 426
    SYS_IO_URING_REGISTER = 427
)
```

---

## futex 实现

```go
// Futex 快速用户空间互斥
type Futex struct {
    state uint32
}

const (
    futexUnlocked = 0
    futexLocked   = 1
    futexWaiters  = 2
)

// Lock 获取锁
func (f *Futex) Lock() {
    // 尝试原子获取锁
    if atomic.CompareAndSwapUint32(&f.state, futexUnlocked, futexLocked) {
        return
    }

    // 慢路径：需要等待
    f.lockSlow()
}

func (f *Futex) lockSlow() {
    for {
        // 检查是否有等待者
        old := atomic.LoadUint32(&f.state)

        var new uint32
        if old == futexLocked {
            new = futexLocked | futexWaiters
        } else {
            new = futexLocked
        }

        if atomic.CompareAndSwapUint32(&f.state, old, new) {
            if old == futexLocked {
                // 已经有锁持有者，进入内核等待
                f.futexWait(futexLocked | futexWaiters)
            } else {
                // 获取成功
                return
            }
        }
    }
}

// Unlock 释放锁
func (f *Futex) Unlock() {
    old := atomic.SwapUint32(&f.state, futexUnlocked)

    // 如果有等待者，唤醒一个
    if old == futexLocked|futexWaiters {
        f.futexWake(1)
    }
}

// futexWait 内核等待
func (f *Futex) futexWait(val uint32) {
    // syscall(SYS_futex, &state, FUTEX_WAIT, val, NULL, NULL, 0)
    syscall.Syscall6(
        syscall.SYS_FUTEX,
        uintptr(unsafe.Pointer(&f.state)),
        uintptr(syscall.FUTEX_WAIT),
        uintptr(val),
        0, 0, 0,
    )
}

// futexWake 内核唤醒
func (f *Futex) futexWake(n int) {
    syscall.Syscall6(
        syscall.SYS_FUTEX,
        uintptr(unsafe.Pointer(&f.state)),
        uintptr(syscall.FUTEX_WAKE),
        uintptr(n),
        0, 0, 0,
    )
}
```

---

## 内核调度策略

```go
// SetScheduler 设置实时调度策略
func SetScheduler(pid int, policy int, priority int) error {
    // SCHED_FIFO: 先进先出实时调度
    // SCHED_RR:   轮询实时调度
    // SCHED_OTHER: 普通分时调度

    param := &syscall.SchedParam{
        SchedPriority: priority,
    }

    return syscall.SchedSetscheduler(pid, policy, param)
}

// 使用示例：设置任务处理器为实时优先级
func setRealTimePriority() error {
    // 获取当前进程
    pid := syscall.Getpid()

    // 设置 SCHED_FIFO，优先级 50 (1-99)
    // 需要 root 权限或 CAP_SYS_NICE
    return SetScheduler(pid, syscall.SCHED_FIFO, 50)
}

// CPU 亲和性设置
func SetCPUAffinity(cpus []int) error {
    var mask syscall.CPUSet

    for _, cpu := range cpus {
        mask.Set(cpu)
    }

    return syscall.SchedSetaffinity(0, &mask)
}

// 内存锁定（防止交换到磁盘）
func LockMemory() error {
    // mlockall(MCL_CURRENT | MCL_FUTURE)
    return syscall.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE)
}
```

---

## 性能对比

| 机制 | 延迟 | 系统调用 | 适用场景 |
|-----|------|---------|---------|
| poll | O(n) | 每次检查 | 少量 fd |
| select | O(n) | 每次检查 | 少量 fd (< 1024) |
| epoll | O(1) | 仅事件时 | 大量 fd |
| io_uring | O(1) | 批量提交 | 高性能 I/O |

---

## 完整调度器集成

```go
// KernelOptimizedScheduler 内核优化调度器
type KernelOptimizedScheduler struct {
    epoll    *EpollPoller
    ioUring  *IoUring

    tasks    chan *Task
    workers  []*KernelWorker
}

// KernelWorker 内核优化工作线程
type KernelWorker struct {
    id       int
    epoll    *EpollPoller
    taskChan chan *Task
}

func (kw *KernelWorker) Run() {
    // 设置实时优先级
    SetScheduler(0, syscall.SCHED_FIFO, 40)

    // 绑定到特定 CPU
    SetCPUAffinity([]int{kw.id % runtime.NumCPU()})

    // 锁定内存
    LockMemory()

    for {
        // 使用 epoll 等待任务和网络事件
        n, _ := kw.epoll.Wait(-1)

        for i := 0; i < n; i++ {
            // 处理事件
            event := kw.epoll.events[i]
            kw.handleEvent(event)
        }

        // 处理任务队列
        select {
        case task := <-kw.taskChan:
            kw.executeTask(task)
        default:
        }
    }
}
```
