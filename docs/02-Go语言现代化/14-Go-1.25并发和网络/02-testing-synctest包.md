# testing/synctest 包（Go 1.25）

> **版本要求**: Go 1.25+  
> **包路径**: `testing/synctest`  
> **实验性**: 否（正式特性）  
> **最后更新**: 2025年10月18日

---

## 📚 目录

- [概述](#概述)
- [为什么需要 synctest](#为什么需要-synctest)
- [核心功能](#核心功能)
- [基本使用](#基本使用)
- [高级特性](#高级特性)
- [实践案例](#实践案例)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)
- [参考资料](#参考资料)

---

## 概述

`testing/synctest` 是 Go 1.25 新增的测试包,专门用于测试并发代码,提供确定性的并发行为模拟和死锁检测。

### 核心价值

- ✅ **确定性测试**: 并发代码可重现
- ✅ **死锁检测**: 自动检测死锁
- ✅ **时间控制**: 模拟时间流逝
- ✅ **简化调试**: 降低并发测试复杂度

---

## 为什么需要 synctest?

### 传统并发测试的问题

#### 问题 1: 不确定性

```go
// 传统测试: 结果不确定
func TestRace(t *testing.T) {
    var count int
    
    go func() { count++ }()
    go func() { count++ }()
    
    time.Sleep(100 * time.Millisecond)
    
    // 可能是 1 或 2,测试不稳定!
    if count != 2 {
        t.Error("race condition")
    }
}
```

#### 问题 2: 死锁难检测

```go
// 死锁可能需要很久才能发现
func TestDeadlock(t *testing.T) {
    ch := make(chan int)
    
    go func() {
        // 忘记发送
    }()
    
    <-ch  // 永远阻塞,测试超时
}
```

#### 问题 3: 时间依赖

```go
// 依赖 sleep,测试慢且不可靠
func TestTimeout(t *testing.T) {
    done := make(chan bool)
    
    go func() {
        time.Sleep(1 * time.Second)
        done <- true
    }()
    
    time.Sleep(2 * time.Second)  // 浪费时间
    // ...
}
```

### synctest 的解决方案

```go
import "testing/synctest"

func TestWithSynctest(t *testing.T) {
    synctest.Run(func() {
        // 确定性并发执行
        // 自动死锁检测
        // 快速时间模拟
    })
}
```

---

## 核心功能

### 1. Run() - 确定性执行

```go
func Run(f func())
```

在确定性环境中运行函数:

- 所有 goroutine 按可预测顺序执行
- 自动检测死锁
- 可重现的测试结果

---

### 2. Wait() - 等待所有 goroutine

```go
func Wait()
```

等待当前环境中所有 goroutine 完成。

---

### 3. 时间控制

模拟时间流逝,无需实际等待:

```go
// 在 synctest 环境中
time.Sleep(1 * time.Hour)  // 立即返回
time.After(1 * time.Hour)  // 立即触发
```

---

## 基本使用

### 示例 1: 简单并发测试

```go
package main

import (
    "testing"
    "testing/synctest"
)

func TestConcurrent(t *testing.T) {
    synctest.Run(func() {
        count := 0
        done := make(chan bool, 2)
        
        // 启动两个 goroutine
        go func() {
            count++
            done <- true
        }()
        
        go func() {
            count++
            done <- true
        }()
        
        // 等待完成
        <-done
        <-done
        
        // 确定性结果
        if count != 2 {
            t.Errorf("expected 2, got %d", count)
        }
    })
}
```

---

### 示例 2: 死锁检测

```go
func TestDeadlock(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Error("expected panic on deadlock")
        }
    }()
    
    synctest.Run(func() {
        ch := make(chan int)
        
        go func() {
            // 忘记发送数据
        }()
        
        <-ch  // 死锁! synctest 会 panic
    })
}
```

---

### 示例 3: 时间控制

```go
func TestTimeout(t *testing.T) {
    synctest.Run(func() {
        start := time.Now()
        
        // 模拟 1 小时延迟,但立即返回
        time.Sleep(1 * time.Hour)
        
        elapsed := time.Since(start)
        
        // 几乎没有实际时间流逝
        if elapsed > 1*time.Second {
            t.Errorf("took too long: %v", elapsed)
        }
    })
}
```

---

### 示例 4: Channel 通信

```go
func TestChannel(t *testing.T) {
    synctest.Run(func() {
        ch := make(chan int)
        
        go func() {
            ch <- 42
        }()
        
        result := <-ch
        
        if result != 42 {
            t.Errorf("expected 42, got %d", result)
        }
    })
}
```

---

## 高级特性

### 1. 嵌套 goroutine

```go
func TestNested(t *testing.T) {
    synctest.Run(func() {
        results := make(chan int, 3)
        
        go func() {
            // 嵌套 goroutine
            go func() {
                results <- 1
            }()
            
            go func() {
                results <- 2
            }()
            
            results <- 3
        }()
        
        // 收集所有结果
        sum := 0
        for i := 0; i < 3; i++ {
            sum += <-results
        }
        
        if sum != 6 {
            t.Errorf("expected 6, got %d", sum)
        }
    })
}
```

---

### 2. select 语句

```go
func TestSelect(t *testing.T) {
    synctest.Run(func() {
        ch1 := make(chan int)
        ch2 := make(chan int)
        
        go func() {
            ch1 <- 1
        }()
        
        go func() {
            ch2 <- 2
        }()
        
        // select 行为确定性
        results := []int{}
        for i := 0; i < 2; i++ {
            select {
            case v := <-ch1:
                results = append(results, v)
            case v := <-ch2:
                results = append(results, v)
            }
        }
        
        // 验证结果
        if len(results) != 2 {
            t.Errorf("expected 2 results, got %d", len(results))
        }
    })
}
```

---

### 3. 互斥锁测试

```go
func TestMutex(t *testing.T) {
    synctest.Run(func() {
        var mu sync.Mutex
        count := 0
        
        for i := 0; i < 10; i++ {
            go func() {
                mu.Lock()
                count++
                mu.Unlock()
            }()
        }
        
        synctest.Wait()  // 等待所有 goroutine
        
        if count != 10 {
            t.Errorf("expected 10, got %d", count)
        }
    })
}
```

---

### 4. Context 测试

```go
func TestContext(t *testing.T) {
    synctest.Run(func() {
        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
        defer cancel()
        
        done := make(chan bool)
        
        go func() {
            select {
            case <-ctx.Done():
                done <- true
            }
        }()
        
        // 立即取消
        cancel()
        
        // 验证取消
        <-done
    })
}
```

---

## 实践案例

### 案例 1: 生产者-消费者测试

```go
func TestProducerConsumer(t *testing.T) {
    synctest.Run(func() {
        buffer := make(chan int, 10)
        produced := 0
        consumed := 0
        
        // 生产者
        go func() {
            for i := 0; i < 100; i++ {
                buffer <- i
                produced++
            }
            close(buffer)
        }()
        
        // 消费者
        go func() {
            for range buffer {
                consumed++
            }
        }()
        
        synctest.Wait()
        
        if produced != 100 || consumed != 100 {
            t.Errorf("produced %d, consumed %d", produced, consumed)
        }
    })
}
```

---

### 案例 2: Worker Pool 测试

```go
func TestWorkerPool(t *testing.T) {
    synctest.Run(func() {
        jobs := make(chan int, 100)
        results := make(chan int, 100)
        
        // 启动 5 个 worker
        for w := 0; w < 5; w++ {
            go func() {
                for job := range jobs {
                    results <- job * 2
                }
            }()
        }
        
        // 发送 100 个任务
        go func() {
            for i := 0; i < 100; i++ {
                jobs <- i
            }
            close(jobs)
        }()
        
        // 收集结果
        sum := 0
        for i := 0; i < 100; i++ {
            sum += <-results
        }
        
        expected := 100 * 99  // 0+1+...+99 = 4950, *2 = 9900
        if sum != expected {
            t.Errorf("expected %d, got %d", expected, sum)
        }
    })
}
```

---

### 案例 3: 超时测试

```go
func TestOperationTimeout(t *testing.T) {
    synctest.Run(func() {
        result := make(chan string, 1)
        
        go func() {
            // 模拟慢操作
            time.Sleep(2 * time.Second)
            result <- "done"
        }()
        
        // 1 秒超时
        select {
        case <-result:
            t.Error("should have timed out")
        case <-time.After(1 * time.Second):
            // 预期超时 (瞬间完成)
        }
    })
}
```

---

### 案例 4: 重试逻辑测试

```go
func TestRetry(t *testing.T) {
    synctest.Run(func() {
        attempts := 0
        
        // 重试函数
        retry := func() error {
            attempts++
            if attempts < 3 {
                return errors.New("temporary failure")
            }
            return nil
        }
        
        // 重试逻辑
        for i := 0; i < 5; i++ {
            if err := retry(); err == nil {
                break
            }
            time.Sleep(1 * time.Second)  // 瞬间完成
        }
        
        if attempts != 3 {
            t.Errorf("expected 3 attempts, got %d", attempts)
        }
    })
}
```

---

## 最佳实践

### 1. 始终在 synctest.Run() 中测试并发

```go
// ✅ 推荐
func TestConcurrent(t *testing.T) {
    synctest.Run(func() {
        // 并发测试代码
    })
}

// ❌ 不推荐: 不确定性
func TestConcurrentBad(t *testing.T) {
    // 直接测试并发,结果不确定
}
```

---

### 2. 使用 synctest.Wait() 等待

```go
// ✅ 推荐
synctest.Run(func() {
    for i := 0; i < 10; i++ {
        go work()
    }
    synctest.Wait()  // 等待所有 goroutine
})

// ❌ 不推荐
synctest.Run(func() {
    for i := 0; i < 10; i++ {
        go work()
    }
    time.Sleep(1 * time.Second)  // 不可靠
})
```

---

### 3. 验证死锁行为

```go
func TestExpectedDeadlock(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Error("expected deadlock panic")
        }
    }()
    
    synctest.Run(func() {
        ch := make(chan int)
        <-ch  // 应该死锁
    })
}
```

---

### 4. 测试 Context 取消

```go
func TestContextCancellation(t *testing.T) {
    synctest.Run(func() {
        ctx, cancel := context.WithCancel(context.Background())
        
        done := make(chan bool)
        go func() {
            <-ctx.Done()
            done <- true
        }()
        
        cancel()
        <-done  // 验证取消生效
    })
}
```

---

### 5. 避免真实 I/O 操作

```go
// ❌ 不推荐: 真实网络请求
synctest.Run(func() {
    http.Get("https://example.com")  // 不确定性
})

// ✅ 推荐: 使用 mock
synctest.Run(func() {
    mockClient.Get("https://example.com")
})
```

---

## 常见问题

### Q1: synctest 支持所有并发原语吗?

**A**: ✅ 支持大部分

- ✅ goroutine
- ✅ channel
- ✅ select
- ✅ sync.Mutex, sync.RWMutex
- ✅ sync.WaitGroup
- ✅ time.Sleep, time.After
- ⚠️ 真实 I/O 操作可能不确定

---

### Q2: synctest 会影响性能吗?

**A**: ⚡ 反而更快!

并发测试在 synctest 中通常更快,因为:

- 时间模拟 (无需真实等待)
- 确定性执行 (减少重试)

---

### Q3: 如何调试 synctest 测试?

**A**: 添加日志

```go
synctest.Run(func() {
    log.Println("Starting goroutine 1")
    go func() {
        log.Println("Goroutine 1 running")
    }()
    
    log.Println("Starting goroutine 2")
    go func() {
        log.Println("Goroutine 2 running")
    }()
    
    synctest.Wait()
})
```

---

### Q4: 可以嵌套 synctest.Run() 吗?

**A**: ❌ 不支持

```go
// ❌ 错误
synctest.Run(func() {
    synctest.Run(func() {  // 嵌套不支持
        // ...
    })
})
```

---

### Q5: 如何测试竞态条件?

**A**: 使用 `-race` 标志

```bash
# synctest + race detector
go test -race -v ./...
```

synctest 提供确定性,race detector 检测数据竞争。

---

## 参考资料

### 官方文档

- 📘 [Go 1.25 Release Notes](https://go.dev/doc/go1.25#testing)
- 📘 [testing/synctest](https://pkg.go.dev/testing/synctest)
- 📘 [Testing Guide](https://go.dev/doc/tutorial/add-a-test)

### 相关章节

- 🔗 [Go 1.25 并发和网络](./README.md)
- 🔗 [并发编程](../../03-并发编程/README.md)
- 🔗 [测试最佳实践](../../08-最佳实践/测试.md)

---

## 更新日志

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2025-10-18 | v1.0 | 初始版本,完整的 synctest 指南 |

---

**编写者**: AI Assistant  
**审核者**: [待审核]  
**最后更新**: 2025年10月18日

---

<p align="center">
  <b>🧪 使用 testing/synctest 让并发测试更可靠! ✅</b>
</p>
