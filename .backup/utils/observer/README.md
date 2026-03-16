# 观察者模式工具

**版本**: v1.0  
**更新日期**: 2025-11-11  
**适用于**: Go 1.25.3

---

## 📋 目录

- [观察者模式工具](#观察者模式工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)

---

## 1. 概述

观察者模式工具提供了观察者模式的实现，包括简单主题、异步主题、过滤主题、事件总线等，帮助开发者实现发布-订阅模式。

---

## 2. 功能特性

### 2.1 观察者接口

- `Observer`: 观察者接口
- `ObserverFunc`: 观察者函数类型

### 2.2 简单主题

- `SimpleSubject`: 简单主题实现
- `NewSimpleSubject`: 创建简单主题
- `Subscribe`: 订阅
- `Unsubscribe`: 取消订阅
- `Notify`: 通知所有观察者
- `Count`: 获取观察者数量
- `Clear`: 清空所有观察者

### 2.3 异步主题

- `AsyncSubject`: 异步主题实现
- `NewAsyncSubject`: 创建异步主题
- `Subscribe`: 订阅
- `Unsubscribe`: 取消订阅
- `Notify`: 异步通知所有观察者
- `Count`: 获取观察者数量
- `Clear`: 清空所有观察者

### 2.4 过滤主题

- `FilteredSubject`: 过滤主题实现
- `NewFilteredSubject`: 创建过滤主题
- `Subscribe`: 订阅（带过滤条件）
- `Unsubscribe`: 取消订阅
- `Notify`: 通知所有观察者（根据过滤条件）
- `Count`: 获取观察者数量
- `Clear`: 清空所有观察者

### 2.5 事件总线

- `EventBus`: 事件总线实现
- `NewEventBus`: 创建事件总线
- `Subscribe`: 订阅事件
- `Publish`: 发布事件
- `Unsubscribe`: 取消订阅
- `Clear`: 清空指定事件的所有观察者
- `ClearAll`: 清空所有事件
- `Count`: 获取指定事件的观察者数量

---

## 3. 使用示例

### 3.1 简单主题

```go
import "github.com/yourusername/golang/pkg/utils/observer"

// 创建主题
subject := observer.NewSimpleSubject[string]()

// 创建观察者
observer1 := observer.ObserverFunc[string](func(data string) {
    fmt.Printf("Observer1 received: %s\n", data)
})

observer2 := observer.ObserverFunc[string](func(data string) {
    fmt.Printf("Observer2 received: %s\n", data)
})

// 订阅
unsubscribe1 := subject.Subscribe(observer1)
unsubscribe2 := subject.Subscribe(observer2)

// 通知
subject.Notify("Hello")

// 取消订阅
unsubscribe1()

// 再次通知
subject.Notify("World")
```

### 3.2 异步主题

```go
// 创建异步主题
subject := observer.NewAsyncSubject[string]()

// 订阅
subject.Subscribe(observer.ObserverFunc[string](func(data string) {
    fmt.Printf("Received: %s\n", data)
}))

// 异步通知
subject.Notify("Hello")
```

### 3.3 过滤主题

```go
// 创建过滤主题
subject := observer.NewFilteredSubject[int]()

// 订阅（只接收大于5的值）
subject.Subscribe(
    observer.ObserverFunc[int](func(data int) {
        fmt.Printf("Received: %d\n", data)
    }),
    func(data int) bool {
        return data > 5
    },
)

// 通知
subject.Notify(3)  // 不会触发
subject.Notify(10) // 会触发
```

### 3.4 事件总线

```go
// 创建事件总线
bus := observer.NewEventBus()

// 订阅事件
bus.Subscribe("user.created", observer.ObserverFunc[interface{}](func(data interface{}) {
    fmt.Printf("User created: %v\n", data)
}))

bus.Subscribe("user.updated", observer.ObserverFunc[interface{}](func(data interface{}) {
    fmt.Printf("User updated: %v\n", data)
}))

// 发布事件
bus.Publish("user.created", map[string]string{"id": "1", "name": "Alice"})
bus.Publish("user.updated", map[string]string{"id": "1", "name": "Bob"})
```

### 3.5 完整示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/golang/pkg/utils/observer"
)

func main() {
    // 简单主题
    subject := observer.NewSimpleSubject[string]()
    
    observer1 := observer.ObserverFunc[string](func(data string) {
        fmt.Printf("Observer1: %s\n", data)
    })
    
    unsubscribe := subject.Subscribe(observer1)
    subject.Notify("Hello")
    unsubscribe()
    
    // 事件总线
    bus := observer.NewEventBus()
    bus.Subscribe("event1", observer.ObserverFunc[interface{}](func(data interface{}) {
        fmt.Printf("Event1: %v\n", data)
    }))
    bus.Publish("event1", "data")
}
```

---

**更新日期**: 2025-11-11

