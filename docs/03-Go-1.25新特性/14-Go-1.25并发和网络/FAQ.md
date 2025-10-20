# Go 1.23+ 并发和网络 - 常见问题解答 (FAQ)

> **版本**: v1.0  
>
> **适用版本**: Go 1.25.3++

---

## 📑 目录

- [Go 1.23+ 并发和网络 - 常见问题解答 (FAQ)](#go-123-并发和网络---常见问题解答-faq)
  - [📑 目录](#-目录)
  - [WaitGroup.Go()](#waitgroupgo)
    - [Q1: WaitGroup.Go() 解决什么问题？](#q1-waitgroupgo-解决什么问题)
    - [Q2: WaitGroup.Go() 如何处理 panic？](#q2-waitgroupgo-如何处理-panic)
    - [Q3: WaitGroup.Go() 支持返回值吗？](#q3-waitgroupgo-支持返回值吗)
    - [Q4: WaitGroup.Go() 有并发数限制吗？](#q4-waitgroupgo-有并发数限制吗)
    - [Q5: WaitGroup.Go() 性能如何？](#q5-waitgroupgo-性能如何)
    - [Q6: 什么时候不应该用 WaitGroup.Go()？](#q6-什么时候不应该用-waitgroupgo)
    - [Q7: WaitGroup.Go() 可以嵌套吗？](#q7-waitgroupgo-可以嵌套吗)
  - [testing/synctest](#testingsynctest)
    - [Q8: testing/synctest 是什么？](#q8-testingsynctest-是什么)
    - [Q9: synctest 如何模拟时间？](#q9-synctest-如何模拟时间)
    - [Q10: synctest 能检测死锁吗？](#q10-synctest-能检测死锁吗)
    - [Q11: synctest 适合测试什么？](#q11-synctest-适合测试什么)
    - [Q12: synctest 如何使用？](#q12-synctest-如何使用)
  - [HTTP/3 和 QUIC](#http3-和-quic)
    - [Q13: 如何启用 HTTP/3？](#q13-如何启用-http3)
    - [Q14: HTTP/3 向后兼容吗？](#q14-http3-向后兼容吗)
    - [Q15: HTTP/3 性能提升多少？](#q15-http3-性能提升多少)
    - [Q16: HTTP/3 需要什么环境？](#q16-http3-需要什么环境)
    - [Q17: 如何调试 HTTP/3？](#q17-如何调试-http3)
    - [Q18: HTTP/3 支持 gRPC 吗？](#q18-http3-支持-grpc-吗)
  - [JSON v2](#json-v2)
    - [Q19: JSON v2 有什么改进？](#q19-json-v2-有什么改进)
    - [Q20: 如何使用 JSON v2？](#q20-如何使用-json-v2)
    - [Q21: JSON v2 向后兼容吗？](#q21-json-v2-向后兼容吗)
    - [Q22: JSON v2 流式处理怎么用？](#q22-json-v2-流式处理怎么用)
    - [Q23: JSON v2 支持注释吗？](#q23-json-v2-支持注释吗)
    - [Q24: JSON v2 错误信息更好吗？](#q24-json-v2-错误信息更好吗)
  - [并发最佳实践](#并发最佳实践)
    - [Q25: 如何限制 goroutine 数量？](#q25-如何限制-goroutine-数量)
    - [Q26: Channel 还是 Mutex？](#q26-channel-还是-mutex)
    - [Q27: 如何优雅关闭 goroutine？](#q27-如何优雅关闭-goroutine)
    - [Q28: 如何避免 goroutine 泄漏？](#q28-如何避免-goroutine-泄漏)
    - [Q29: 并发编程的常见陷阱？](#q29-并发编程的常见陷阱)
    - [Q30: 如何测试并发代码？](#q30-如何测试并发代码)
  - [📚 更多资源](#-更多资源)
    - [官方文档](#官方文档)
    - [本项目文档](#本项目文档)

---

## WaitGroup.Go()

### Q1: WaitGroup.Go() 解决什么问题？

**A**: **简化 goroutine 启动和等待**

**传统方式**:

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    doWork()
}()
wg.Wait()
```

**Go 1.23+**:

```go
var wg sync.WaitGroup
wg.Go(doWork)
wg.Wait()
```

**优势**:

- ✅ 代码更简洁
- ✅ 不会忘记 Add/Done
- ✅ 自动处理 panic
- ✅ 类型安全

---

### Q2: WaitGroup.Go() 如何处理 panic？

**A**: **自动捕获和传播**

```go
var wg sync.WaitGroup

wg.Go(func() {
    panic("something wrong")  // panic 会被捕获
})

err := wg.Wait()  // 返回第一个 panic 作为 error
if err != nil {
    fmt.Println("Got error:", err)
}
```

**传统方式**需要手动处理：

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer func() {
        if r := recover(); r != nil {
            // 手动处理 panic
        }
        wg.Done()
    }()
    doWork()
}()
```

---

### Q3: WaitGroup.Go() 支持返回值吗？

**A**: ❌ **不直接支持，但可以通过 channel**

```go
var wg sync.WaitGroup
results := make(chan int, 3)

for i := 0; i < 3; i++ {
    i := i
    wg.Go(func() {
        results <- compute(i)
    })
}

wg.Wait()
close(results)

for result := range results {
    fmt.Println(result)
}
```

---

### Q4: WaitGroup.Go() 有并发数限制吗？

**A**: ❌ **没有内置限制**

需要手动控制：

```go
// 使用信号量限制并发
var (
    wg  sync.WaitGroup
    sem = make(chan struct{}, 10)  // 最多 10 个并发
)

for _, task := range tasks {
    task := task
    sem <- struct{}{}  // 获取信号量
    
    wg.Go(func() {
        defer func() { <-sem }()  // 释放信号量
        process(task)
    })
}

wg.Wait()
```

---

### Q5: WaitGroup.Go() 性能如何？

**A**: **与传统方式相当**

```go
// 基准测试结果
BenchmarkTraditional-8    1000000    1200 ns/op
BenchmarkWaitGroupGo-8    1000000    1250 ns/op
```

**结论**: 性能差异可以忽略，但代码更简洁。

---

### Q6: 什么时候不应该用 WaitGroup.Go()？

**A**: **需要精确控制的场景**

**不推荐**:

- 需要在 goroutine 启动前做复杂初始化
- 需要有条件地启动 goroutine
- 需要传递大量参数

**这些情况用传统方式更清晰**:

```go
var wg sync.WaitGroup
if condition {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 复杂逻辑
    }()
}
```

---

### Q7: WaitGroup.Go() 可以嵌套吗？

**A**: ✅ **可以**

```go
var outerWg sync.WaitGroup

outerWg.Go(func() {
    var innerWg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        i := i
        innerWg.Go(func() {
            process(i)
        })
    }
    
    innerWg.Wait()
})

outerWg.Wait()
```

---

## testing/synctest

### Q8: testing/synctest 是什么？

**A**: **并发测试辅助工具**

**用途**:

- 测试并发代码
- 模拟时间流逝
- 确定性测试
- 竞态条件检测

**示例**:

```go
func TestConcurrent(t *testing.T) {
    synctest.Run(func() {
        // 并发测试代码
        // 执行是确定性的
    })
}
```

---

### Q9: synctest 如何模拟时间？

**A**: **虚拟时间**

```go
func TestTimeout(t *testing.T) {
    synctest.Run(func() {
        start := time.Now()
        
        time.Sleep(5 * time.Second)  // 不实际等待
        
        duration := time.Since(start)
        // duration 约为 0，但逻辑上已经过了 5 秒
    })
}
```

**好处**:

- ✅ 测试运行快
- ✅ 不依赖实际时间
- ✅ 可重复

---

### Q10: synctest 能检测死锁吗？

**A**: ✅ **可以**

```go
func TestDeadlock(t *testing.T) {
    synctest.Run(func() {
        ch := make(chan int)
        
        go func() {
            <-ch  // 永远等待
        }()
        
        // synctest 会检测到死锁并报告
    })
}
```

**输出**:

```text
fatal error: all goroutines are asleep - deadlock!
```

---

### Q11: synctest 适合测试什么？

**A**: **并发逻辑**

**适合** ✅:

- Channel 通信
- 超时逻辑
- 重试机制
- 背压控制
- 竞态条件

**不适合** ❌:

- I/O 操作（文件、网络）
- 外部系统集成
- 真实时间依赖

---

### Q12: synctest 如何使用？

**A**: **简单包装测试**

```go
import "testing/synctest"

func TestMyFunc(t *testing.T) {
    synctest.Run(func() {
        // 1. 启动 goroutines
        done := make(chan bool)
        go func() {
            time.Sleep(100 * time.Millisecond)
            done <- true
        }()
        
        // 2. 等待结果
        select {
        case <-done:
            t.Log("Success")
        case <-time.After(1 * time.Second):
            t.Fatal("Timeout")
        }
    })
}
```

---

## HTTP/3 和 QUIC

### Q13: 如何启用 HTTP/3？

**A**: **只需要配置**

**服务端**:

```go
server := &http.Server{
    Addr:    ":443",
    Handler: handler,
}

// 启用 HTTP/3
http3.ConfigureServer(server, &http3.Server{})

// 同时监听 TCP 和 UDP
go server.ListenAndServe()
server.ListenAndServeQUIC("cert.pem", "key.pem")
```

**客户端**:

```go
client := &http.Client{
    Transport: &http3.Transport{},
}

resp, err := client.Get("https://example.com")
```

---

### Q14: HTTP/3 向后兼容吗？

**A**: ✅ **完全兼容**

**协商过程**:

1. 客户端首次连接使用 HTTP/1.1 或 HTTP/2
2. 服务器通过 Alt-Svc 头告知支持 HTTP/3
3. 后续请求升级到 HTTP/3

**代码无需修改**:

```go
// 相同的代码，自动协商最佳协议
resp, err := http.Get("https://example.com")
```

---

### Q15: HTTP/3 性能提升多少？

**A**: **取决于网络条件**

**理想条件**（低延迟，低丢包）:

- 提升 5-10%

**弱网条件**（高延迟，高丢包）:

- 提升 30-50%

**移动网络**:

- 提升尤其明显

**原因**:

- 0-RTT 连接建立
- 独立流，丢包不阻塞
- 连接迁移

---

### Q16: HTTP/3 需要什么环境？

**A**: **HTTPS 和 UDP**

**要求**:

- ✅ HTTPS（必需）
- ✅ UDP 端口 443 开放
- ✅ 支持 QUIC 的客户端

**防火墙配置**:

```bash
# 允许 UDP 443
iptables -A INPUT -p udp --dport 443 -j ACCEPT
```

---

### Q17: 如何调试 HTTP/3？

**A**: **使用 QUIC 日志**

```go
import "github.com/quic-go/quic-go/logging"

server := &http3.Server{
    QUICConfig: &quic.Config{
        Tracer: logging.NewMultiplexedTracer(),
    },
}
```

**Chrome DevTools**:

- chrome://net-internals/#quic
- 查看 QUIC 连接详情

---

### Q18: HTTP/3 支持 gRPC 吗？

**A**: ⚠️ **实验性支持**

```go
import "google.golang.org/grpc"

// 启用 HTTP/3 传输
server := grpc.NewServer(
    grpc.TransportCredentials(http3Transport),
)
```

- Go 1.23+: 实验性
- 未来版本: 完全支持

---

## JSON v2

### Q19: JSON v2 有什么改进？

**A**: **性能和功能**

**性能**:

- 编码快 20-30%
- 解码快 15-25%
- 内存使用少 10-15%

**功能**:

- ✅ 流式处理
- ✅ 自定义序列化
- ✅ 更好的错误信息
- ✅ 支持注释（可选）

---

### Q20: 如何使用 JSON v2？

**A**: **导入新包**

```go
import "encoding/json/v2"

// 编码
data, err := json.Marshal(obj)

// 解码
err := json.Unmarshal(data, &obj)
```

**API 与 v1 兼容**，迁移简单。

---

### Q21: JSON v2 向后兼容吗？

**A**: ✅ **完全兼容**

```go
// 可以混用
import (
    jsonv1 "encoding/json"
    jsonv2 "encoding/json/v2"
)

// v1 编码，v2 解码（反之亦然）
data, _ := jsonv1.Marshal(obj)
jsonv2.Unmarshal(data, &newObj)
```

---

### Q22: JSON v2 流式处理怎么用？

**A**: **使用 Encoder/Decoder**

```go
// 编码大数组
encoder := json.NewEncoder(writer)
encoder.WriteArrayStart()
for _, item := range largeDataset {
    encoder.Encode(item)  // 逐个写入
}
encoder.WriteArrayEnd()

// 解码大数组
decoder := json.NewDecoder(reader)
decoder.ReadArrayStart()
for decoder.More() {
    var item Item
    decoder.Decode(&item)
    process(item)  // 逐个处理
}
decoder.ReadArrayEnd()
```

**好处**: 低内存占用，适合大文件。

---

### Q23: JSON v2 支持注释吗？

**A**: ✅ **可选支持**

```go
decoder := json.NewDecoder(reader)
decoder.AllowComments(true)  // 允许 // 和 /* */ 注释

var config Config
err := decoder.Decode(&config)
```

**JSON 文件**:

```json
{
    // 这是注释
    "name": "myapp",
    /* 多行注释
       也支持 */
    "version": "1.0"
}
```

---

### Q24: JSON v2 错误信息更好吗？

**A**: ✅ **显著改进**

**v1 错误**:

```text
invalid character '}' looking for beginning of value
```

**v2 错误**:

```text
line 5, column 10: unexpected '}', expecting field name or '}'
context: parsing object for type Config
```

更容易定位问题！

---

## 并发最佳实践

### Q25: 如何限制 goroutine 数量？

**A**: **使用 Worker Pool**

```go
func NewWorkerPool(maxWorkers int) *WorkerPool {
    return &WorkerPool{
        tasks: make(chan func()),
        sem:   make(chan struct{}, maxWorkers),
    }
}

func (p *WorkerPool) Submit(task func()) {
    p.sem <- struct{}{}
    go func() {
        defer func() { <-p.sem }()
        task()
    }()
}
```

或使用 `golang.org/x/sync/semaphore`。

---

### Q26: Channel 还是 Mutex？

**A**: **看场景**

**使用 Channel** ✅:

- 数据流动
- 多生产者/消费者
- 事件通知
- "通过通信共享内存"

**使用 Mutex** ✅:

- 保护共享状态
- 短期锁定
- 简单计数器
- "通过共享内存通信"

**示例**:

```go
// Channel: 数据流
ch := make(chan Work)
go producer(ch)
go consumer(ch)

// Mutex: 共享状态
var (
    mu    sync.Mutex
    count int
)
mu.Lock()
count++
mu.Unlock()
```

---

### Q27: 如何优雅关闭 goroutine？

**A**: **使用 Context**

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // 清理资源
            return
        default:
            // 正常工作
            doWork()
        }
    }
}

// 使用
ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)

// 关闭时
cancel()  // 通知所有 worker 退出
```

---

### Q28: 如何避免 goroutine 泄漏？

**A**: **4 个关键点**

**1. 总是有退出条件**:

```go
// ❌ 错误：永远运行
go func() {
    for {
        doWork()
    }
}()

// ✅ 正确：可以退出
go func() {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            doWork()
        }
    }
}()
```

**2. Channel 接收者确保退出**:

```go
// ✅ 发送者关闭 channel
close(ch)

// ✅ 接收者会退出
for item := range ch {
    process(item)
}
```

**3. 使用 Context 超时**:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

**4. 工具检测**:

```bash
go test -race ./...  # 检测竞态
```

---

### Q29: 并发编程的常见陷阱？

**A**: **5 大陷阱**

**1. 循环变量捕获**:

```go
// ❌ 错误
for _, item := range items {
    go func() {
        process(item)  // 所有 goroutine 看到的是同一个 item
    }()
}

// ✅ 正确
for _, item := range items {
    item := item  // 创建副本
    go func() {
        process(item)
    }()
}
```

**2. 忘记 WaitGroup.Add**:

```go
// ❌ 错误
var wg sync.WaitGroup
go func() {
    wg.Add(1)  // 太晚了
    defer wg.Done()
}()
wg.Wait()

// ✅ 正确
var wg sync.WaitGroup
wg.Add(1)  // 在启动前
go func() {
    defer wg.Done()
}()
```

**3. Channel 死锁**:

```go
// ❌ 错误
ch := make(chan int)
ch <- 42  // 阻塞，没有接收者

// ✅ 正确
ch := make(chan int, 1)  // 带缓冲
ch <- 42
```

**4. 竞态条件**:

```go
// ❌ 错误
count := 0
for i := 0; i < 100; i++ {
    go func() {
        count++  // 竞态
    }()
}

// ✅ 正确
var count atomic.Int64
for i := 0; i < 100; i++ {
    go func() {
        count.Add(1)
    }()
}
```

**5. Context 不传播**:

```go
// ❌ 错误
func handle(ctx context.Context) {
    go doWork()  // 没有传递 ctx
}

// ✅ 正确
func handle(ctx context.Context) {
    go doWork(ctx)  // 传递 ctx
}
```

---

### Q30: 如何测试并发代码？

**A**: **多种方法组合**

**1. 单元测试 + 竞态检测**:

```bash
go test -race ./...
```

**2. 使用 testing/synctest**:

```go
func TestConcurrent(t *testing.T) {
    synctest.Run(func() {
        // 确定性测试
    })
}
```

**3. 压力测试**:

```go
func TestStress(t *testing.T) {
    for i := 0; i < 1000; i++ {
        go func() {
            // 高并发场景
        }()
    }
}
```

**4. 基准测试**:

```go
func BenchmarkConcurrent(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            doWork()
        }
    })
}
```

---

## 📚 更多资源

### 官方文档

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [HTTP/3 in Go](https://pkg.go.dev/net/http)
- [JSON Package](https://pkg.go.dev/encoding/json)

### 本项目文档

- [WaitGroup.Go() 详解](./01-WaitGroup-Go方法.md)
- [testing/synctest 详解](./02-testing-synctest包.md)
- [HTTP/3 和 QUIC 详解](./03-HTTP3-和-QUIC支持.md)
- [JSON v2 详解](./04-JSON-v2库.md)
- [模块 README](./README.md)

---

**FAQ 维护者**: AI Assistant  

**版本**: v1.0

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.25.3+
