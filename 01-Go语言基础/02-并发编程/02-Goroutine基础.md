# 1.2.1 Goroutine基础

<!-- TOC START -->
- [1.2.1 Goroutine基础](#goroutine基础)
  - [1.2.1.1 📚 **理论分析**](#📚-**理论分析**)
    - [1.2.1.1.1 **Goroutine定义与原理**](#**goroutine定义与原理**)
      - [1.2.1.1.1.1 **形式化描述**](#**形式化描述**)
    - [1.2.1.1.2 **生命周期与调度**](#**生命周期与调度**)
    - [1.2.1.1.3 **调度机制**](#**调度机制**)
  - [1.2.1.2 💻 **代码示例**](#💻-**代码示例**)
    - [1.2.1.2.1 **创建与运行Goroutine**](#**创建与运行goroutine**)
    - [1.2.1.2.2 **Goroutine并发执行**](#**goroutine并发执行**)
    - [1.2.1.2.3 **使用sync.WaitGroup等待Goroutine完成**](#**使用syncwaitgroup等待goroutine完成**)
    - [1.2.1.2.4 **Goroutine泄漏示例与避免**](#**goroutine泄漏示例与避免**)
  - [1.2.1.3 📊 **性能分析**](#📊-**性能分析**)
  - [1.2.1.4 🧪 **测试代码**](#🧪-**测试代码**)
  - [1.2.1.5 🎯 **最佳实践**](#🎯-**最佳实践**)
  - [1.2.1.6 🔍 **常见问题**](#🔍-**常见问题**)
  - [1.2.1.7 📚 **扩展阅读**](#📚-**扩展阅读**)
<!-- TOC END -->

## 1.2.1.1 📚 **理论分析**

### 1.2.1.1.1 **Goroutine定义与原理**

- Goroutine是Go语言的核心并发单元，极轻量级，由Go运行时调度。
- 每个Goroutine初始栈空间约2KB，可动态扩展。
- Goroutine通过`go`关键字创建，运行在Go调度器（G-P-M模型）之上。

#### 1.2.1.1.1.1 **形式化描述**

```text
Goroutine ::= go FunctionCall
FunctionCall ::= identifier(arguments)

```

### 1.2.1.1.2 **生命周期与调度**

- 创建：`go f()`立即返回，f在新Goroutine中异步执行
- 运行：由Go调度器分配到内核线程
- 阻塞：遇到IO、Channel、锁等操作时自动挂起
- 退出：函数返回或运行时panic

### 1.2.1.1.3 **调度机制**

- Go运行时采用协作式抢占调度，避免长时间阻塞
- G-P-M模型：G（Goroutine）、P（Processor）、M（Machine）
- Goroutine间切换开销极低

## 1.2.1.2 💻 **代码示例**

### 1.2.1.2.1 **创建与运行Goroutine**

```go
package main
import (
    "fmt"
    "time"
)
func sayHello() {
    fmt.Println("Hello from goroutine!")
}
func main() {
    go sayHello()
    fmt.Println("Hello from main!")
    time.Sleep(time.Millisecond * 50)
}

```

### 1.2.1.2.2 **Goroutine并发执行**

```go
package main
import (
    "fmt"
    "time"
)
func worker(id int) {
    fmt.Printf("Worker %d start\n", id)
    time.Sleep(time.Millisecond * 100)
    fmt.Printf("Worker %d done\n", id)
}
func main() {
    for i := 1; i <= 3; i++ {
        go worker(i)
    }
    time.Sleep(time.Millisecond * 200)
    fmt.Println("All workers launched")
}

```

### 1.2.1.2.3 **使用sync.WaitGroup等待Goroutine完成**

```go
package main
import (
    "fmt"
    "sync"
)
func main() {
    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Worker %d\n", id)
        }(i)
    }
    wg.Wait()
    fmt.Println("All done")
}

```

### 1.2.1.2.4 **Goroutine泄漏示例与避免**

```go
package main
import "time"
func leak() {
    ch := make(chan int)
    go func() {
        // 永远阻塞，泄漏
        ch <- 1
    }()
    time.Sleep(time.Millisecond * 10)
}
func main() {
    leak()
    // 正确做法：确保有接收方或关闭通道
}

```

## 1.2.1.3 📊 **性能分析**

- Goroutine创建和切换开销远低于系统线程
- 可轻松创建数万~百万级Goroutine
- 过多Goroutine会导致调度压力和内存消耗

## 1.2.1.4 🧪 **测试代码**

```go
package main
import (
    "sync"
    "testing"
)
func TestGoroutine(t *testing.T) {
    var wg sync.WaitGroup
    sum := 0
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(x int) {
            defer wg.Done()
            sum += x // 存在竞态，演示用
        }(i)
    }
    wg.Wait()
    // 结果不确定，演示竞态
}

```

## 1.2.1.5 🎯 **最佳实践**

- 使用`sync.WaitGroup`等待所有Goroutine完成
- 避免在Goroutine中访问共享变量（用Channel或锁）
- 控制Goroutine数量，防止资源耗尽
- 保证每个Goroutine都有退出路径，避免泄漏

## 1.2.1.6 🔍 **常见问题**

- Q: Goroutine和线程有何区别？
  A: Goroutine由Go运行时调度，远比线程轻量
- Q: 如何安全地在Goroutine间传递数据？
  A: 使用Channel或同步原语
- Q: Goroutine泄漏如何排查？
  A: 使用pprof、trace等工具分析

## 1.2.1.7 📚 **扩展阅读**

- [Go官方文档-Goroutine](https://golang.org/doc/effective_go.html#goroutines)
- [Go by Example: Goroutines](https://gobyexample.com/goroutines)
- [Go Scheduler: G-P-M模型](https://blog.golang.org/2018/09/go-scheduler.html)

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成
