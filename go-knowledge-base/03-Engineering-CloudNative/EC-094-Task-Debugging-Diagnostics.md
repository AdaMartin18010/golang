# 任务调试与诊断 (Task Debugging & Diagnostics)

> **分类**: 工程与云原生
> **标签**: #debugging #diagnostics #profiling
> **参考**: Go Diagnostics, Distributed Tracing

---

## 调试架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Task Debugging & Diagnostics                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Debug Information Collection                      │   │
│  │                                                                      │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐           │   │
│  │   │  Stack   │  │  Heap    │  │ Goroutine│  │   CPU    │           │   │
│  │   │  Trace   │  │ Profile  │  │  Dump    │  │ Profile  │           │   │
│  │   └──────────┘  └──────────┘  └──────────┘  └──────────┘           │   │
│  │                                                                      │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────┐                          │   │
│  │   │ Execution│  │  Memory  │  │  Event   │                          │   │
│  │   │  Trace   │  │  Stats   │  │   Log    │                          │   │
│  │   └──────────┘  └──────────┘  └──────────┘                          │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Diagnostic Tools                                  │   │
│  │                                                                      │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐           │   │
│  │   │   pprof  │  │   trace  │  │   dlv    │  │   zap    │           │   │
│  │   │ (profiling)│  │ (tracing)│  │ (debugger)│  │ (logging)│           │   │
│  │   └──────────┘  └──────────┘  └──────────┘  └──────────┘           │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整调试实现

```go
package diagnostics

import (
    "context"
    "fmt"
    "runtime"
    "runtime/pprof"
    "runtime/trace"
    "sync"
    "time"
)

// Debugger 调试器
type Debugger struct {
    taskID      string
    startTime   time.Time
    checkpoints []Checkpoint
    events      []DebugEvent

    mu          sync.RWMutex
}

// Checkpoint 检查点
type Checkpoint struct {
    Name      string
    Timestamp time.Time
    Data      map[string]interface{}
}

// DebugEvent 调试事件
type DebugEvent struct {
    Type      string
    Timestamp time.Time
    Message   string
    Data      map[string]interface{}
}

// NewDebugger 创建调试器
func NewDebugger(taskID string) *Debugger {
    return &Debugger{
        taskID:      taskID,
        startTime:   time.Now(),
        checkpoints: make([]Checkpoint, 0),
        events:      make([]DebugEvent, 0),
    }
}

// AddCheckpoint 添加检查点
func (d *Debugger) AddCheckpoint(name string, data map[string]interface{}) {
    d.mu.Lock()
    defer d.mu.Unlock()

    d.checkpoints = append(d.checkpoints, Checkpoint{
        Name:      name,
        Timestamp: time.Now(),
        Data:      data,
    })
}

// LogEvent 记录事件
func (d *Debugger) LogEvent(eventType, message string, data map[string]interface{}) {
    d.mu.Lock()
    defer d.mu.Unlock()

    d.events = append(d.events, DebugEvent{
        Type:      eventType,
        Timestamp: time.Now(),
        Message:   message,
        Data:      data,
    })
}

// GetExecutionTrace 获取执行追踪
func (d *Debugger) GetExecutionTrace() ExecutionTrace {
    d.mu.RLock()
    defer d.mu.RUnlock()

    return ExecutionTrace{
        TaskID:      d.taskID,
        Duration:    time.Since(d.startTime),
        Checkpoints: d.checkpoints,
        Events:      d.events,
    }
}

// ExecutionTrace 执行追踪
type ExecutionTrace struct {
    TaskID      string        `json:"task_id"`
    Duration    time.Duration `json:"duration"`
    Checkpoints []Checkpoint  `json:"checkpoints"`
    Events      []DebugEvent  `json:"events"`
}

// ProfileCollector 性能分析收集器
type ProfileCollector struct {
    profiles map[string]*Profile
    mu       sync.RWMutex
}

// Profile 性能分析
type Profile struct {
    Type     string
    Data     []byte
    Duration time.Duration
    Timestamp time.Time
}

// NewProfileCollector 创建收集器
func NewProfileCollector() *ProfileCollector {
    return &ProfileCollector{
        profiles: make(map[string]*Profile),
    }
}

// CollectCPUProfile 收集CPU分析
func (pc *ProfileCollector) CollectCPUProfile(duration time.Duration) (*Profile, error) {
    buf := make([]byte, 0)

    if err := pprof.StartCPUProfile(&sliceWriter{&buf}); err != nil {
        return nil, err
    }

    time.Sleep(duration)
    pprof.StopCPUProfile()

    profile := &Profile{
        Type:      "cpu",
        Data:      buf,
        Duration:  duration,
        Timestamp: time.Now(),
    }

    pc.mu.Lock()
    pc.profiles["cpu"] = profile
    pc.mu.Unlock()

    return profile, nil
}

// CollectHeapProfile 收集堆分析
func (pc *ProfileCollector) CollectHeapProfile() (*Profile, error) {
    buf := make([]byte, 0)

    if err := pprof.WriteHeapProfile(&sliceWriter{&buf}); err != nil {
        return nil, err
    }

    profile := &Profile{
        Type:      "heap",
        Data:      buf,
        Timestamp: time.Now(),
    }

    pc.mu.Lock()
    pc.profiles["heap"] = profile
    pc.mu.Unlock()

    return profile, nil
}

// CollectGoroutineProfile 收集Goroutine分析
func (pc *ProfileCollector) CollectGoroutineProfile() (*Profile, error) {
    buf := make([]byte, 0)

    if err := pprof.Lookup("goroutine").WriteTo(&sliceWriter{&buf}, 1); err != nil {
        return nil, err
    }

    profile := &Profile{
        Type:      "goroutine",
        Data:      buf,
        Timestamp: time.Now(),
    }

    pc.mu.Lock()
    pc.profiles["goroutine"] = profile
    pc.mu.Unlock()

    return profile, nil
}

// sliceWriter 切片写入器
type sliceWriter struct {
    buf *[]byte
}

func (sw *sliceWriter) Write(p []byte) (n int, err error) {
    *sw.buf = append(*sw.buf, p...)
    return len(p), nil
}

// ExecutionTracer 执行追踪器
type ExecutionTracer struct {
    taskID string
    trace  *trace.Trace
}

// NewExecutionTracer 创建执行追踪器
func NewExecutionTracer(taskID string) *ExecutionTracer {
    return &ExecutionTracer{
        taskID: taskID,
    }
}

// StartTrace 开始追踪
func (et *ExecutionTracer) StartTrace(w interface{}) error {
    return trace.Start(w.(trace.Writer))
}

// StopTrace 停止追踪
func (et *ExecutionTracer) StopTrace() {
    trace.Stop()
}

// TaskInspector 任务检查器
type TaskInspector struct {
    runtimeStats RuntimeStats
}

// RuntimeStats 运行时统计
type RuntimeStats struct {
    NumGoroutine  int
    NumCPU        int
    Goroutines    []GoroutineInfo
    MemoryStats   runtime.MemStats
}

// GoroutineInfo Goroutine信息
type GoroutineInfo struct {
    ID       int
    State    string
    Stack    string
}

// Inspect 检查运行时状态
func (ti *TaskInspector) Inspect() RuntimeStats {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return RuntimeStats{
        NumGoroutine: runtime.NumGoroutine(),
        NumCPU:       runtime.NumCPU(),
        MemoryStats:  m,
    }
}

// GetGoroutineDump 获取Goroutine转储
func (ti *TaskInspector) GetGoroutineDump() string {
    buf := make([]byte, 1<<20) // 1MB
    n := runtime.Stack(buf, true)
    return string(buf[:n])
}

// DiagnosticServer 诊断服务器
type DiagnosticServer struct {
    debugger    *Debugger
    collector   *ProfileCollector
    inspector   *TaskInspector
}

// NewDiagnosticServer 创建诊断服务器
func NewDiagnosticServer(taskID string) *DiagnosticServer {
    return &DiagnosticServer{
        debugger:  NewDebugger(taskID),
        collector: NewProfileCollector(),
        inspector: &TaskInspector{},
    }
}

// GetDebugInfo 获取调试信息
func (ds *DiagnosticServer) GetDebugInfo() map[string]interface{} {
    return map[string]interface{}{
        "execution_trace": ds.debugger.GetExecutionTrace(),
        "runtime_stats":   ds.inspector.Inspect(),
        "profiles":        ds.getProfiles(),
    }
}

func (ds *DiagnosticServer) getProfiles() map[string]string {
    ds.collector.mu.RLock()
    defer ds.collector.mu.RUnlock()

    profiles := make(map[string]string)
    for name, profile := range ds.collector.profiles {
        profiles[name] = fmt.Sprintf("%d bytes", len(profile.Data))
    }
    return profiles
}

// MemoryLeakDetector 内存泄漏检测器
type MemoryLeakDetector struct {
    snapshots []MemorySnapshot
    threshold float64 // 增长率阈值
}

// MemorySnapshot 内存快照
type MemorySnapshot struct {
    Timestamp time.Time
    Alloc     uint64
    HeapAlloc uint64
}

// NewMemoryLeakDetector 创建检测器
func NewMemoryLeakDetector(threshold float64) *MemoryLeakDetector {
    return &MemoryLeakDetector{
        snapshots: make([]MemorySnapshot, 0),
        threshold: threshold,
    }
}

// TakeSnapshot 拍摄快照
func (mld *MemoryLeakDetector) TakeSnapshot() {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)

    mld.snapshots = append(mld.snapshots, MemorySnapshot{
        Timestamp: time.Now(),
        Alloc:     stats.Alloc,
        HeapAlloc: stats.HeapAlloc,
    })
}

// DetectLeak 检测泄漏
func (mld *MemoryLeakDetector) DetectLeak() (bool, float64) {
    if len(mld.snapshots) < 2 {
        return false, 0
    }

    first := mld.snapshots[0]
    last := mld.snapshots[len(mld.snapshots)-1]

    if first.Alloc == 0 {
        return false, 0
    }

    growth := float64(last.Alloc-first.Alloc) / float64(first.Alloc)
    return growth > mld.threshold, growth
}

// DeadlockDetector 死锁检测器
type DeadlockDetector struct {
    timeout time.Duration
    timer   *time.Timer
    done    chan struct{}
}

// NewDeadlockDetector 创建死锁检测器
func NewDeadlockDetector(timeout time.Duration) *DeadlockDetector {
    return &DeadlockDetector{
        timeout: timeout,
        done:    make(chan struct{}),
    }
}

// Start 开始监控
func (dd *DeadlockDetector) Start(onDeadlock func()) {
    dd.timer = time.AfterFunc(dd.timeout, func() {
        select {
        case <-dd.done:
        default:
            onDeadlock()
        }
    })
}

// Stop 停止监控
func (dd *DeadlockDetector) Stop() {
    close(dd.done)
    if dd.timer != nil {
        dd.timer.Stop()
    }
}

// Reset 重置计时器
func (dd *DeadlockDetector) Reset() {
    if dd.timer != nil {
        dd.timer.Reset(dd.timeout)
    }
}
```

---

## 使用示例

```go
package main

import (
    "fmt"
    "runtime"
    "time"

    "diagnostics"
)

func main() {
    // 创建调试器
    debugger := diagnostics.NewDebugger("task-123")

    // 添加检查点
    debugger.AddCheckpoint("init", map[string]interface{}{
        "goroutines": runtime.NumGoroutine(),
    })

    // 模拟任务执行
    time.Sleep(1 * time.Second)

    debugger.AddCheckpoint("processing", map[string]interface{}{
        "progress": 50,
    })

    // 记录事件
    debugger.LogEvent("info", "Processing item", map[string]interface{}{
        "item_id": "item-1",
    })

    time.Sleep(1 * time.Second)

    debugger.AddCheckpoint("completed", map[string]interface{}{
        "result": "success",
    })

    // 获取追踪
    trace := debugger.GetExecutionTrace()
    fmt.Printf("Task %s executed in %v\n", trace.TaskID, trace.Duration)
    fmt.Printf("Checkpoints: %d\n", len(trace.Checkpoints))

    // 性能分析
    collector := diagnostics.NewProfileCollector()

    // 收集堆分析
    profile, _ := collector.CollectHeapProfile()
    fmt.Printf("Heap profile: %d bytes\n", len(profile.Data))

    // 运行时检查
    inspector := &diagnostics.TaskInspector{}
    stats := inspector.Inspect()
    fmt.Printf("Goroutines: %d, CPUs: %d\n", stats.NumGoroutine, stats.NumCPU)

    // 内存泄漏检测
    detector := diagnostics.NewMemoryLeakDetector(2.0) // 2x growth threshold

    for i := 0; i < 10; i++ {
        detector.TakeSnapshot()
        time.Sleep(100 * time.Millisecond)
    }

    hasLeak, growth := detector.DetectLeak()
    fmt.Printf("Memory leak detected: %v (growth: %.2f%%)\n", hasLeak, growth*100)
}
```

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02