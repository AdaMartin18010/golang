# 2.1 select与context高级用法

<!-- TOC START -->
- [2.1 select与context高级用法](#21-select与context高级用法)
  - [2.1.1 1. 理论基础](#211-1-理论基础)
    - [2.1.1.1 select语句](#2111-select语句)
    - [2.1.1.2 context包](#2112-context包)
  - [2.1.2 2. 典型用法](#212-2-典型用法)
    - [2.1.2.1 select实现超时控制](#2121-select实现超时控制)
    - [2.1.2.2 select实现多路复用](#2122-select实现多路复用)
    - [2.1.2.3 context实现取消](#2123-context实现取消)
    - [2.1.2.4 context实现超时](#2124-context实现超时)
  - [2.1.3 3. 工程分析与最佳实践](#213-3-工程分析与最佳实践)
  - [2.1.4 4. 常见陷阱](#214-4-常见陷阱)
  - [2.1.5 5. 单元测试建议](#215-5-单元测试建议)
  - [2.1.6 6. 参考文献](#216-6-参考文献)
<!-- TOC END -->

## 2.1.1 1. 理论基础

### 2.1.1.1 select语句

select语句用于监听多个channel操作，实现多路复用、超时、取消等高级控制。

- 形式化描述：
  \[
    \text{select} \{ c_1, c_2, ..., c_n \}
  \]
  表示等待多个channel中的任意一个可用。

### 2.1.1.2 context包

context用于跨Goroutine传递取消信号、超时、元数据，是Go并发控制的标准方式。

- 典型结构：
  - context.Background()
  - context.WithCancel(parent)
  - context.WithTimeout(parent, duration)
  - context.WithValue(parent, key, value)

---

## 2.1.2 2. 典型用法

### 2.1.2.1 select实现超时控制

```go
ch := make(chan int)
select {
case v := <-ch:
    fmt.Println("received", v)
case <-time.After(time.Second):
    fmt.Println("timeout")
}
```

### 2.1.2.2 select实现多路复用

```go
select {
case v1 := <-ch1:
    fmt.Println("ch1:", v1)
case v2 := <-ch2:
    fmt.Println("ch2:", v2)
}
```

### 2.1.2.3 context实现取消

```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    <-ctx.Done()
    fmt.Println("cancelled")
}()
cancel()
```

### 2.1.2.4 context实现超时

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
select {
case <-ctx.Done():
    fmt.Println("timeout or cancelled")
}
```

---

## 2.1.3 3. 工程分析与最佳实践

- select可优雅处理channel超时、取消、优先级等复杂场景。
- context应作为函数参数首选，便于链式传递。
- 推荐用context统一管理Goroutine生命周期，防止泄漏。
- select+context组合是高并发服务的标配。
- 注意select分支顺序无优先级，随机选择可用分支。

---

## 2.1.4 4. 常见陷阱

- 忘记cancel context会导致资源泄漏。
- select所有分支都阻塞时会死锁。
- context.Value仅用于传递请求范围内的元数据，勿滥用。

---

## 2.1.5 5. 单元测试建议

- 测试超时、取消、并发场景下的正确性。
- 覆盖边界与异常情况。

---

## 2.1.6 6. 参考文献

- Go官方文档：<https://golang.org/doc/>
- Go Blog: <https://blog.golang.org/context>
- 《Go语言高级编程》
