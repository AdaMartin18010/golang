# 1. ⏰ Cron 深度解析

> **简介**: 本文档详细阐述了 Cron 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. ⏰ Cron 深度解析](#1--cron-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 Cron 表达式](#131-cron-表达式)
    - [1.3.2 任务调度](#132-任务调度)
    - [1.3.3 任务管理](#133-任务管理)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 Cron 使用最佳实践](#141-cron-使用最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Cron 是什么？**

Cron 是一个基于时间的任务调度器，用于在指定时间执行任务。

**核心特性**:

- ✅ **时间调度**: 支持复杂的 Cron 表达式
- ✅ **并发安全**: 支持并发执行任务
- ✅ **任务管理**: 支持添加、删除、暂停任务
- ✅ **时区支持**: 支持时区配置

---

## 1.2 选型论证

**为什么选择 Cron？**

**论证矩阵**:

| 评估维度 | 权重 | robfig/cron | gocron | go-cron | 标准库 time.Ticker | 说明 |
|---------|------|-------------|--------|---------|-------------------|------|
| **功能完整性** | 30% | 10 | 8 | 7 | 5 | robfig/cron 功能最完整 |
| **易用性** | 25% | 9 | 9 | 8 | 7 | robfig/cron 易用性好 |
| **性能** | 20% | 9 | 8 | 8 | 10 | robfig/cron 性能优秀 |
| **维护性** | 15% | 10 | 7 | 6 | 8 | robfig/cron 维护性好 |
| **社区支持** | 10% | 10 | 7 | 6 | 10 | robfig/cron 社区活跃 |
| **加权总分** | - | **9.50** | 8.00 | 7.20 | 7.40 | robfig/cron 得分最高 |

**核心优势**:

1. **功能完整性（权重 30%）**:
   - 支持标准 Cron 表达式
   - 支持秒级精度
   - 支持时区配置

2. **易用性（权重 25%）**:
   - API 简洁，易于使用
   - 文档完善
   - 示例丰富

---

## 1.3 实际应用

### 1.3.1 Cron 表达式

**Cron 表达式格式**:

```go
// Cron 表达式格式（支持秒）
// 秒 分 时 日 月 周
// * * * * * *
// | | | | | |
// | | | | | +-- 周 (0-6, 0=Sunday)
// | | | | +---- 月 (1-12)
// | | | +------ 日 (1-31)
// | | +-------- 时 (0-23)
// | +---------- 分 (0-59)
// +------------ 秒 (0-59)

// 示例
// "0 0 * * * *"     - 每小时执行
// "0 0 0 * * *"     - 每天 0 点执行
// "0 0 0 * * 0"     - 每周日 0 点执行
// "0 */5 * * * *"   - 每 5 分钟执行
// "*/30 * * * * *"  - 每 30 秒执行
```

### 1.3.2 任务调度

**创建 Cron 调度器**:

```go
// internal/infrastructure/scheduler/cron.go
package scheduler

import (
    "github.com/robfig/cron/v3"
)

type Scheduler struct {
    cron *cron.Cron
}

func NewScheduler() *Scheduler {
    // 支持秒级精度
    c := cron.New(cron.WithSeconds())
    return &Scheduler{cron: c}
}

// 添加任务
func (s *Scheduler) AddJob(spec string, cmd func()) (cron.EntryID, error) {
    return s.cron.AddFunc(spec, cmd)
}

// 启动调度器
func (s *Scheduler) Start() {
    s.cron.Start()
}

// 停止调度器
func (s *Scheduler) Stop() {
    s.cron.Stop()
}

// 使用示例
func Example() {
    scheduler := NewScheduler()

    // 每天 0 点执行
    scheduler.AddJob("0 0 0 * * *", func() {
        logger.Info("Daily task executed")
    })

    // 每 5 分钟执行
    scheduler.AddJob("0 */5 * * * *", func() {
        logger.Info("Periodic task executed")
    })

    scheduler.Start()
    defer scheduler.Stop()
}
```

### 1.3.3 任务管理

**任务管理**:

```go
// 任务管理
type TaskManager struct {
    cron    *cron.Cron
    entries map[string]cron.EntryID
}

func NewTaskManager() *TaskManager {
    return &TaskManager{
        cron:    cron.New(cron.WithSeconds()),
        entries: make(map[string]cron.EntryID),
    }
}

// 添加任务
func (tm *TaskManager) AddTask(name, spec string, cmd func()) error {
    id, err := tm.cron.AddFunc(spec, cmd)
    if err != nil {
        return err
    }
    tm.entries[name] = id
    return nil
}

// 删除任务
func (tm *TaskManager) RemoveTask(name string) error {
    id, ok := tm.entries[name]
    if !ok {
        return errors.New("task not found")
    }
    tm.cron.Remove(id)
    delete(tm.entries, name)
    return nil
}

// 列出所有任务
func (tm *TaskManager) ListTasks() []string {
    var tasks []string
    for name := range tm.entries {
        tasks = append(tasks, name)
    }
    return tasks
}
```

---

## 1.4 最佳实践

### 1.4.1 Cron 使用最佳实践

**为什么需要最佳实践？**

合理的 Cron 使用可以提高系统的稳定性和可维护性。

**最佳实践原则**:

1. **错误处理**: 完善的错误处理和日志记录
2. **并发控制**: 避免任务重叠执行
3. **资源管理**: 合理管理资源，避免泄漏
4. **监控告警**: 监控任务执行状态

**实际应用示例**:

```go
// Cron 最佳实践
type SafeScheduler struct {
    cron    *cron.Cron
    running map[string]bool
    mu      sync.Mutex
}

func NewSafeScheduler() *SafeScheduler {
    return &SafeScheduler{
        cron:    cron.New(cron.WithSeconds()),
        running: make(map[string]bool),
    }
}

// 安全执行任务（防止重叠）
func (s *SafeScheduler) AddSafeJob(name, spec string, cmd func()) error {
    return s.cron.AddFunc(spec, func() {
        s.mu.Lock()
        if s.running[name] {
            s.mu.Unlock()
            logger.Warn("Task already running", "name", name)
            return
        }
        s.running[name] = true
        s.mu.Unlock()

        defer func() {
            s.mu.Lock()
            delete(s.running, name)
            s.mu.Unlock()
        }()

        // 执行任务
        start := time.Now()
        defer func() {
            if r := recover(); r != nil {
                logger.Error("Task panicked", "name", name, "error", r)
            }
            logger.Info("Task completed", "name", name, "duration", time.Since(start))
        }()

        cmd()
    })
}

// 带超时的任务
func (s *SafeScheduler) AddJobWithTimeout(name, spec string, timeout time.Duration, cmd func()) error {
    return s.cron.AddFunc(spec, func() {
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()

        done := make(chan error, 1)
        go func() {
            defer func() {
                if r := recover(); r != nil {
                    done <- fmt.Errorf("panic: %v", r)
                }
            }()
            cmd()
            done <- nil
        }()

        select {
        case err := <-done:
            if err != nil {
                logger.Error("Task failed", "name", name, "error", err)
            }
        case <-ctx.Done():
            logger.Error("Task timeout", "name", name, "timeout", timeout)
        }
    })
}
```

**最佳实践要点**:

1. **错误处理**: 完善的错误处理和日志记录
2. **并发控制**: 使用锁防止任务重叠执行
3. **超时控制**: 为任务设置超时，避免长时间阻塞
4. **监控告警**: 监控任务执行状态，及时发现问题

---

## 📚 扩展阅读

- [robfig/cron 官方文档](https://github.com/robfig/cron)
- [Cron 表达式指南](https://en.wikipedia.org/wiki/Cron)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Cron 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
