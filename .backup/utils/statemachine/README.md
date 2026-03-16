# 状态机工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [状态机工具](#状态机工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

状态机工具提供了状态机实现，支持状态转换、回调函数、状态验证等，帮助开发者管理复杂的状态流转逻辑。

---

## 2. 功能特性

### 2.1 状态机

- `StateMachine`: 状态机实现
- `NewStateMachine`: 创建状态机
- `AddTransition`: 添加状态转换
- `AddTransitions`: 批量添加状态转换
- `OnEnter`: 设置进入状态时的回调
- `OnExit`: 设置离开状态时的回调
- `OnTransition`: 设置状态转换时的回调
- `Trigger`: 触发事件
- `Current`: 获取当前状态
- `CanTrigger`: 检查是否可以触发事件
- `Reset`: 重置状态机
- `GetTransitions`: 获取所有状态转换
- `GetAvailableEvents`: 获取当前状态可用的所有事件

### 2.2 简单状态机

- `SimpleStateMachine`: 简单状态机实现（使用字符串）
- `NewSimpleStateMachine`: 创建简单状态机

---

## 3. 使用示例

### 3.1 基本使用

```go
import "github.com/yourusername/golang/pkg/utils/statemachine"

// 创建状态机
sm := statemachine.NewStateMachine[string, string]("idle")

// 添加状态转换
sm.AddTransition("idle", "start", "running")
sm.AddTransition("running", "stop", "idle")
sm.AddTransition("running", "pause", "paused")
sm.AddTransition("paused", "resume", "running")

// 触发事件
err := sm.Trigger("start")
if err != nil {
    fmt.Printf("Error: %v\n", err)
}

// 获取当前状态
current := sm.Current()  // "running"

// 检查是否可以触发事件
if sm.CanTrigger("stop") {
    sm.Trigger("stop")
}
```

### 3.2 回调函数

```go
sm := statemachine.NewStateMachine[string, string]("idle")

sm.AddTransition("idle", "start", "running")

// 设置进入状态时的回调
sm.OnEnter("running", func() {
    fmt.Println("Entered running state")
})

// 设置离开状态时的回调
sm.OnExit("idle", func() {
    fmt.Println("Exited idle state")
})

// 设置状态转换时的回调
sm.OnTransition("idle", "start", func() {
    fmt.Println("Transitioning from idle to running")
})

sm.Trigger("start")
```

### 3.3 批量添加转换

```go
sm := statemachine.NewStateMachine[string, string]("idle")

transitions := []statemachine.Transition[string, string]{
    {From: "idle", Event: "start", To: "running"},
    {From: "running", Event: "stop", To: "idle"},
    {From: "running", Event: "pause", To: "paused"},
}

sm.AddTransitions(transitions)
```

### 3.4 获取可用事件

```go
// 获取当前状态可用的所有事件
events := sm.GetAvailableEvents()
fmt.Printf("Available events: %v\n", events)

// 获取所有状态转换
transitions := sm.GetTransitions()
for _, t := range transitions {
    fmt.Printf("%v -> %v -> %v\n", t.From, t.Event, t.To)
}
```

### 3.5 简单状态机

```go
// 创建简单状态机（使用字符串）
sm := statemachine.NewSimpleStateMachine("idle")

sm.AddTransition("idle", "start", "running")
sm.Trigger("start")
```

### 3.6 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/statemachine"
)

type State string
type Event string

const (
    StateIdle    State = "idle"
    StateRunning State = "running"
    StatePaused  State = "paused"
)

const (
    EventStart  Event = "start"
    EventStop   Event = "stop"
    EventPause  Event = "pause"
    EventResume Event = "resume"
)

func main() {
    sm := statemachine.NewStateMachine[State, Event](StateIdle)
    
    // 添加状态转换
    sm.AddTransition(StateIdle, EventStart, StateRunning)
    sm.AddTransition(StateRunning, EventStop, StateIdle)
    sm.AddTransition(StateRunning, EventPause, StatePaused)
    sm.AddTransition(StatePaused, EventResume, StateRunning)
    
    // 设置回调
    sm.OnEnter(StateRunning, func() {
        fmt.Println("Started")
    })
    
    sm.OnExit(StateRunning, func() {
        fmt.Println("Stopped")
    })
    
    // 触发事件
    sm.Trigger(EventStart)
    fmt.Printf("Current state: %s\n", sm.Current())
    
    sm.Trigger(EventPause)
    fmt.Printf("Current state: %s\n", sm.Current())
}
```

---

**更新日期**: 2025-11-11

