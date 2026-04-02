# Go 运行时 (Runtime)

> **分类**: 语言设计

---

## 核心组件

```
┌─────────────────────────────────────┐
│            Go Runtime               │
├─────────────┬─────────────┬─────────┤
│  调度器     │   内存管理   │   GC    │
│ Scheduler   │   Memory    │ Collector│
├─────────────┴─────────────┴─────────┤
│        系统调用接口 (Syscall)        │
└─────────────────────────────────────┘
```

---

## 调度器

### G-M-P 模型

```
G: Goroutine - 用户级线程
M: Machine - OS 线程
P: Processor - 逻辑处理器
```

```go
// 设置 GOMAXPROCS
runtime.GOMAXPROCS(4)

// 获取当前 goroutine ID
// (runtime 不直接提供，可通过 hack 获取)
```

---

## 内存分配

### 分级分配

| 级别 | 大小 |
|------|------|
| Tiny | < 16B |
| Small | < 32KB |
| Large | >= 32KB |

```go
// 查看内存统计
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
```

---

## GC 调优

```go
// 设置 GC 目标百分比
GOGC=100  // 默认，堆增长 100% 触发 GC

// 设置内存限制 (Go 1.19+)
GOMEMLIMIT=1GiB
```

---

## 调试

```go
// 打印调度器信息
GODEBUG=schedtrace=1000 ./program

// GC  trace
GODEBUG=gctrace=1 ./program
```
