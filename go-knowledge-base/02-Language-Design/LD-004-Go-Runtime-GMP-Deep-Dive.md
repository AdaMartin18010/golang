# LD-004: Go 运行时调度器深度分析 (Go Runtime GMP Scheduler Deep Dive)

> **维度**: Language Design
> **级别**: S (30+ KB)
> **标签**: #go-runtime #gmp-scheduler #goroutine #work-stealing
> **权威来源**: [Go Scheduler](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html), [Go Runtime Source](https://github.com/golang/go/tree/master/src/runtime)

---

## GMP 模型

```
G (Goroutine)    M (Machine/OS Thread)    P (Processor)
    │                    │                      │
    │  ~2KB stack       │  ~8MB stack         │  Run Queue
    │  growable         │  fixed              │  Local queue
    │                   │                     │
    └───────────────────┴─────────────────────┘
                        │
              ┌─────────┴─────────┐
              │   Global RunQ     │
              │   Work Stealing   │
              │   Network Poller  │
              └───────────────────┘
```

---

## 核心数据结构

```go
// G (Goroutine)
type g struct {
    stack      stack     // 栈
    sched      gobuf     // 调度上下文
    goid       int64     // ID
    atomicstatus uint32  // 状态
    m          *m         // 绑定的 M
    p          uintptr   // 绑定的 P
}

// M (Machine)
type m struct {
    g0      *g      // 调度 goroutine
    curg    *g      // 当前 goroutine
    p       puintptr // 绑定的 P
    id      int64   // 线程 ID
}

// P (Processor)
type p struct {
    id       int32
    runq     [256]guintptr  // 本地队列
    runnext  guintptr      // 下一个运行的 G
    mcache   *mcache       // 内存分配缓存
}
```

---

## 调度循环

```go
func schedule() {
    _g_ := getg()
    _p_ := _g_.m.p.ptr()

    // 1. 从本地队列获取
    if gp, _ := runqget(_p_); gp != nil {
        execute(gp)
        return
    }

    // 2. 从全局队列获取
    if gp := globrunqget(_p_, 0); gp != nil {
        execute(gp)
        return
    }

    // 3. 网络轮询
    if list := netpoll(0); !list.empty() {
        gp := list.pop()
        injectglist(&list)
        execute(gp)
        return
    }

    // 4. Work Stealing
    if gp := stealWork(); gp != nil {
        execute(gp)
        return
    }

    // 5. 空闲
    stopm()
}
```

---

## Work Stealing 算法

```go
func stealWork() *g {
    // 随机顺序尝试从其他 P 偷取
    for i := 0; i < 4; i++ {
        p2 := allp[stealOrder[_p_.id][i]]
        if gp := runqsteal(_p_, p2); gp != nil {
            return gp
        }
    }
    return nil
}

// 偷取半个队列
func runqsteal(_p_, p2 *p) *g {
    n := runqgrabsize(p2)
    batch := runqgrab(p2, &p2.runq, n)

    for i := 0; i < n; i++ {
        runqput(_p_, batch[i], false)
    }
    return batch[n-1]
}
```

---

## 系统调用处理

```
G1 (running)
   │ syscall
   ▼
entersyscall()
   │
   ├── 保存状态
   ├── 释放 P
   └── 其他 M 可以获取 P 继续工作

syscall return
   │
exitsyscall()
   │
   ├── 尝试获取原来的 P
   ├── 如果失败，获取新的 P
   └── 恢复执行
```

---

## 抢占式调度

```go
// 信号驱动的抢占
func preemptM(mp *m) {
    signalM(mp, sigPreempt)  // 发送 SIGURG
}

// 在函数调用边界检查抢占标志
func morestack() {
    if gp.preempt {
        goschedImpl(gp)
    }
}
```

---

## 性能数据

| 操作 | 时间 |
|------|------|
| goroutine 创建 | ~200ns |
| goroutine 切换 | ~200ns |
| thread 创建 | ~1-2μs |
| thread 切换 | ~1-2μs |

---

## 参考文献

1. [Scheduling In Go](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)
2. [Go Runtime Source](https://github.com/golang/go/tree/master/src/runtime)
3. [Analysis of Go Runtime Scheduler](http://www1.cs.columbia.edu/~aho/cs6998/reports/12-12-11_DeshpandeSponslerWeiss_GO.pdf)
