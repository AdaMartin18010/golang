# Go调试与性能分析

> Go程序调试、性能分析和问题排查技术

---

## 一、调试工具

### 1.1 Delve调试器

```text
Delve安装与使用:
────────────────────────────────────────

安装:
go install github.com/go-delve/delve/cmd/dlv@latest

基本命令:
dlv debug                    # 调试当前包
dlv test                     # 调试测试
dlv attach <pid>             # 附加到进程
dlv exec <binary>            # 调试可执行文件

常用命令:
├─ break (b): 设置断点
├─ continue (c): 继续执行
├─ next (n): 单步跳过
├─ step (s): 单步进入
├─ print (p): 打印变量
├─ locals: 显示本地变量
├─ goroutines: 显示所有goroutine
└─ exit: 退出

代码示例:
// 程序
package main

func main() {
    x := 10
    y := 20
    z := add(x, y)
    println(z)
}

func add(a, b int) int {
    return a + b
}

调试会话:
$ dlv debug
(dlv) break main.add
(dlv) continue
(dlv) print a
10
(dlv) print b
20
(dlv) next
(dlv) print $result
30
```

### 1.2 VSCode调试配置

```text
launch.json配置:
────────────────────────────────────────

{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}"
        },
        {
            "name": "Launch Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${fileDirname}"
        },
        {
            "name": "Attach to Process",
            "type": "go",
            "request": "attach",
            "mode": "local",
            "processId": 0
        }
    ]
}
```

---

## 二、性能分析

### 2.1 CPU分析

```text
CPU Profiling:
────────────────────────────────────────

import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // 应用代码
}

分析命令:
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

交互命令:
├─ top: 显示热点函数
├─ list FuncName: 查看函数源码
├─ web: 生成调用图
├─ svg: 生成SVG
└─ pdf: 生成PDF

代码示例:
// 程序内收集
func writeCPUProfile(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    if err := pprof.StartCPUProfile(f); err != nil {
        return err
    }
    defer pprof.StopCPUProfile()

    // 执行操作...
    time.Sleep(30 * time.Second)

    return nil
}
```

### 2.2 内存分析

```text
Memory Profiling:
────────────────────────────────────────

分析类型:
├─ inuse_space: 当前占用内存
├─ inuse_objects: 当前占用对象数
├─ alloc_space: 累计分配内存
└─ alloc_objects: 累计分配对象数

分析命令:
go tool pprof http://localhost:6060/debug/pprof/heap
go tool pprof -inuse_space heap.out
go tool pprof -alloc_space heap.out

查找泄漏:
1. 在不同时间点收集heap profile
2. 对比差异
go tool pprof -base base.heap current.heap

代码示例:
// 触发GC后收集
func captureHeap(filename string) error {
    runtime.GC()  // 先GC，获得准确数据

    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    return pprof.WriteHeapProfile(f)
}
```

### 2.3 Goroutine分析

```text
Goroutine Profiling:
────────────────────────────────────────

查看所有goroutine:
go tool pprof http://localhost:6060/debug/pprof/goroutine

分析阻塞:
go tool pprof http://localhost:6060/debug/pprof/block

分析锁竞争:
go tool pprof http://localhost:6060/debug/pprof/mutex

代码示例:
// 检测goroutine泄漏
func detectGoroutineLeak() {
    baseline := runtime.NumGoroutine()

    // 执行操作...

    time.Sleep(100 * time.Millisecond)
    current := runtime.NumGoroutine()

    if current > baseline+10 {
        log.Printf("Potential leak: %d -> %d goroutines", baseline, current)

        // 输出所有goroutine堆栈
        pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
    }
}
```

---

## 三、执行追踪

### 3.1 Trace工具

```text
Execution Tracing:
────────────────────────────────────────

收集trace:
import "runtime/trace"

func main() {
    f, err := os.Create("trace.out")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    if err := trace.Start(f); err != nil {
        log.Fatal(err)
    }
    defer trace.Stop()

    // 应用代码
}

分析trace:
go tool trace trace.out

查看内容:
├─ View trace: 时间线视图
├─ Goroutine analysis: goroutine统计
├─ Network blocking profile: 网络阻塞
├─ Synchronization blocking profile: 同步阻塞
├─ Syscall blocking profile: 系统调用阻塞
└─ Scheduler latency profile: 调度延迟

代码示例:
// 创建任务区域
type Task struct {
    ctx context.Context
}

func (t *Task) Run() {
    ctx, task := trace.NewTask(t.ctx, "taskName")
    defer task.End()

    // 子区域
    trace.WithRegion(ctx, "subTask", func() {
        // 执行工作
    })
}
```

---

## 四、日志诊断

### 4.1 结构化日志

```text
结构化日志实现:
────────────────────────────────────────

使用zap:
import "go.uber.org/zap"

var logger *zap.Logger

func init() {
    var err error
    logger, err = zap.NewProduction()
    if err != nil {
        log.Fatal(err)
    }
}

func example() {
    logger.Info("request processed",
        zap.String("method", "GET"),
        zap.String("path", "/api/users"),
        zap.Int("status", 200),
        zap.Duration("latency", time.Millisecond*45),
    )
}

// 带上下文的日志
type contextKey struct{}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
    return context.WithValue(ctx, contextKey{}, logger)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
    if logger, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
        return logger
    }
    return zap.NewNop()
}

// 使用
func handleRequest(ctx context.Context, req *Request) {
    logger := LoggerFromContext(ctx)
    logger.Info("handling request",
        zap.String("request_id", req.ID),
    )
}
```

### 4.2 分布式追踪

```text
OpenTelemetry实现:
────────────────────────────────────────

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("my-service")

func processRequest(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "process-request")
    defer span.End()

    span.SetAttributes(
        attribute.String("user.id", "123"),
        attribute.Int("items.count", 10),
    )

    // 子span
    ctx, dbSpan := tracer.Start(ctx, "db-query")
    result, err := db.Query(ctx, "SELECT * FROM users")
    if err != nil {
        dbSpan.RecordError(err)
    }
    dbSpan.End()

    span.SetAttributes(attribute.Int("result.count", len(result)))
}
```

---

*本章提供了Go程序调试、性能分析和问题排查的实用技术。*
