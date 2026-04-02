# Goroutines

> **分类**: 语言设计

---

## 定义

Goroutine 是 Go 的轻量级线程，由 Go 运行时管理。

```go
go function()  // 启动 goroutine
```

---

## 特性

| 特性 | 说明 |
|------|------|
| 轻量 | 初始栈 2KB，动态增长 |
| 快速创建 | 比 OS 线程快得多 |
| 由运行时调度 | G-M-P 模型 |
| 非抢占式 (Go 1.14 前) | 现在可抢占 |

---

## 调度模型

### G-M-P

```
G: Goroutine (执行单元)
M: Machine (OS 线程)
P: Processor (逻辑处理器)
```

```
全局队列
    ↓
P[0] - M[0] - G
P[1] - M[1] - G
P[2] - M[2] - G
```

---

## 创建与退出

```go
// 创建
func main() {
    go say("hello")
    go say("world")
    time.Sleep(time.Second)
}

func say(s string) {
    for i := 0; i < 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println(s)
    }
}
```

---

## 同步

### WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        work()
    }()
}

wg.Wait()  // 等待所有完成
```

### 退出 goroutine

```go
// 正常退出
return

// 异常退出
runtime.Goexit()

// 整个程序退出
os.Exit(0)
```

---

## 最佳实践

### 不要过度创建

```go
// 不好: 无限创建
for _, item := range items {
    go process(item)  // 可能创建数百万 goroutine
}

// 好: 使用工作池
pool := make(chan struct{}, 10)  // 限制10个
for _, item := range items {
    pool <- struct{}{}  // 获取令牌
    go func(i Item) {
        defer func() { <-pool }()  // 释放令牌
        process(i)
    }(item)
}
```

---

## 与线程对比

| 特性 | Goroutine | OS Thread |
|------|-----------|-----------|
| 栈大小 | 2KB (动态) | 1-2MB (固定) |
| 创建时间 | ~μs | ~ms |
| 切换开销 | 小 | 大 |
| 数量 | 百万级 | 千级 |
| 调度 | Go 运行时 | OS 内核 |
