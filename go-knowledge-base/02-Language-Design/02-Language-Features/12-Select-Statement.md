# Select 语句

> **分类**: 语言设计

---

## 基本用法

```go
select {
case v1 := <-ch1:
    fmt.Println("ch1:", v1)
case v2 := <-ch2:
    fmt.Println("ch2:", v2)
case ch3 <- 100:
    fmt.Println("sent to ch3")
default:
    fmt.Println("no channel ready")
}
```

---

## 非阻塞选择

```go
// 带 default 的非阻塞
select {
case v := <-ch:
    fmt.Println("received:", v)
default:
    fmt.Println("no data available")
}
```

---

## 超时模式

```go
func withTimeout(ch chan string, timeout time.Duration) (string, bool) {
    select {
    case v := <-ch:
        return v, true
    case <-time.After(timeout):
        return "", false
    }
}
```

---

## 永久等待

```go
// 不带 default 会阻塞直到有 case 可用
select {
case v := <-ch:
    fmt.Println(v)
}
```

---

## 多路复用

```go
func fanIn(inputs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            out <- n
        }
    }

    wg.Add(len(inputs))
    for _, in := range inputs {
        go output(in)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

---

## 优先级选择

```go
// 优先从 priorityCh 接收
for {
    select {
    case v := <-priorityCh:
        handlePriority(v)
    default:
        select {
        case v := <-normalCh:
            handleNormal(v)
        case v := <-priorityCh:
            handlePriority(v)
        }
    }
}
```

---

## 优雅退出

```go
func worker(ctx context.Context, jobs <-chan int) {
    for {
        select {
        case <-ctx.Done():
            return
        case job := <-jobs:
            process(job)
        }
    }
}
```
