# 事件总线框架

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [事件总线框架](#事件总线框架)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 核心功能](#2-核心功能)
  - [3. 使用示例](#3-使用示例)
  - [4. 最佳实践](#4-最佳实践)

---

## 1. 概述

事件总线框架提供了发布-订阅模式的事件处理功能：

- ✅ **发布订阅**: 支持事件发布和订阅
- ✅ **事件过滤**: 支持事件过滤器
- ✅ **异步处理**: 异步事件处理
- ✅ **指标统计**: 事件处理指标统计
- ✅ **线程安全**: 完全线程安全

---

## 2. 核心功能

### 2.1 Event 接口

```go
type Event interface {
    Type() string
    Data() interface{}
    Timestamp() time.Time
}
```

### 2.2 事件总线

- **订阅管理**: 支持订阅和取消订阅
- **事件过滤**: 支持自定义事件过滤器
- **异步处理**: 事件异步处理，不阻塞发布者
- **指标统计**: 提供事件处理指标

---

## 3. 使用示例

### 3.1 基本使用

```go
import (
    "context"
    "github.com/yourusername/golang/pkg/eventbus"
)

// 创建事件总线
eb := eventbus.NewEventBus(100)
eb.Start()
defer eb.Stop()

// 订阅事件
handler := func(ctx context.Context, event eventbus.Event) error {
    // 处理事件
    fmt.Printf("Received event: %s, data: %v\n", event.Type(), event.Data())
    return nil
}

_, err := eb.Subscribe("user.created", handler)
if err != nil {
    // 处理错误
}

// 发布事件
event := eventbus.NewEvent("user.created", map[string]interface{}{
    "user_id": "123",
    "name":    "John",
})
eb.Publish(event)
```

### 3.2 带过滤器

```go
// 只处理特定用户的事件
filter := func(event eventbus.Event) bool {
    if data, ok := event.Data().(map[string]interface{}); ok {
        if userID, ok := data["user_id"].(string); ok {
            return userID == "123"
        }
    }
    return false
}

_, err := eb.SubscribeWithFilter("user.created", handler, filter)
```

### 3.3 异步发布

```go
// 异步发布，不阻塞
eb.PublishAsync(event)
```

### 3.4 获取指标

```go
metrics := eb.GetMetrics()
fmt.Printf("Total events: %d\n", metrics.TotalEvents)
fmt.Printf("Handled events: %d\n", metrics.HandledEvents)
fmt.Printf("Failed events: %d\n", metrics.FailedEvents)
```

### 3.5 在领域事件中使用

```go
// 领域事件
type UserCreatedEvent struct {
    *eventbus.BaseEvent
    UserID string
    Name   string
}

func NewUserCreatedEvent(userID, name string) *UserCreatedEvent {
    event := &UserCreatedEvent{
        BaseEvent: eventbus.NewEvent("user.created", nil),
        UserID:    userID,
        Name:      name,
    }
    event.SetMetadata("user_id", userID)
    return event
}

// 发布领域事件
event := NewUserCreatedEvent("123", "John")
eb.Publish(event)
```

---

## 4. 最佳实践

### 4.1 DO's ✅

1. **使用Start/Stop**: 启动和停止事件总线
2. **错误处理**: 在Handler中正确处理错误
3. **事件类型**: 使用清晰的事件类型命名
4. **过滤器**: 使用过滤器减少不必要的处理
5. **指标监控**: 定期检查指标以监控系统健康

### 4.2 DON'Ts ❌

1. **不要阻塞**: Handler不应该长时间阻塞
2. **不要忽略错误**: 正确处理Handler错误
3. **不要忘记停止**: 应用关闭时停止事件总线
4. **不要过度订阅**: 避免订阅过多事件类型

---

## 5. 相关资源

- [框架拓展计划](../../docs/00-框架拓展计划.md)

---

**更新日期**: 2025-11-11
