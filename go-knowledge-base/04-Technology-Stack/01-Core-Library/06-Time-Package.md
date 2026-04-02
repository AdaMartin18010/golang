# time 包详解

> **分类**: 开源技术堆栈

---

## 基本使用

### 获取时间

```go
now := time.Now()
fmt.Println(now.Format("2006-01-02 15:04:05"))
```

### 格式化

```go
// Go 特殊时间格式参考
const layout = "2006-01-02 15:04:05"
fmt.Println(now.Format(layout))
```

---

## 时间操作

### 加减

```go
future := now.Add(24 * time.Hour)
past := now.Add(-7 * 24 * time.Hour)
```

### 差值

```go
duration := future.Sub(now)
fmt.Println(duration.Hours())
```

---

## Timer 和 Ticker

```go
// Timer: 一次性
timer := time.NewTimer(5 * time.Second)
<-timer.C
fmt.Println("Timer expired")

// Ticker: 重复
ticker := time.NewTicker(1 * time.Second)
for t := range ticker.C {
    fmt.Println("Tick at", t)
}
```

---

## Context 超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```
