# FV Tool Complete Tutorial

**Go Formal Verifier (FV)** 完整实战教程

---

## 📋 目录

1. [教程概述](#教程概述)
2. [准备工作](#准备工作)
3. [实战案例1: Web API服务分析](#实战案例1-web-api服务分析)
4. [实战案例2: 并发密集型应用](#实战案例2-并发密集型应用)
5. [实战案例3: 遗留代码重构](#实战案例3-遗留代码重构)
6. [高级功能](#高级功能)
7. [故障排查](#故障排查)

---

## 教程概述

本教程将通过真实案例，带你深入了解 FV 工具的使用方法和最佳实践。

**你将学到**:

- 如何分析不同类型的Go项目
- 如何解读和处理分析报告
- 如何配置FV以适应项目需求
- 如何集成到开发流程

**前置要求**:

- 已安装 FV 工具
- 熟悉 Go 语言基础
- 了解基本的命令行操作

---

## 准备工作

### 创建示例项目

我们将创建一个包含常见问题的示例项目：

```bash
mkdir fv-tutorial
cd fv-tutorial
go mod init github.com/example/fv-tutorial
```

### 示例代码

创建 `main.go`:

```go
package main

import (
 "fmt"
 "sync"
 "time"
)

// 问题1: 复杂度过高的函数
func processOrder(orderID int, items []string, discount float64, 
                  shipping bool, priority int, customer map[string]string) error {
 if orderID <= 0 {
  return fmt.Errorf("invalid order ID")
 }
 
 if len(items) == 0 {
  return fmt.Errorf("no items")
 }
 
 total := 0.0
 for _, item := range items {
  if item == "book" {
   total += 10.0
  } else if item == "pen" {
   total += 2.0
  } else if item == "notebook" {
   total += 5.0
  } else if item == "laptop" {
   total += 1000.0
  } else if item == "mouse" {
   total += 20.0
  } else {
   total += 1.0
  }
 }
 
 if discount > 0 {
  if discount < 0.1 {
   total *= 0.95
  } else if discount < 0.2 {
   total *= 0.9
  } else if discount < 0.3 {
   total *= 0.85
  } else {
   total *= 0.8
  }
 }
 
 if shipping {
  if priority == 1 {
   total += 20.0
  } else if priority == 2 {
   total += 10.0
  } else {
   total += 5.0
  }
 }
 
 fmt.Printf("Order %d total: $%.2f\n", orderID, total)
 return nil
}

// 问题2: Goroutine泄漏
func startWorker() {
 ch := make(chan int)
 
 go func() {
  for v := range ch {
   fmt.Println(v)
  }
 }()
 
 // 忘记关闭channel，导致goroutine泄漏
 ch <- 1
 ch <- 2
}

// 问题3: 数据竞争
var counter int

func incrementCounter() {
 var wg sync.WaitGroup
 
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter++ // 数据竞争
  }()
 }
 
 wg.Wait()
}

// 问题4: 不安全的类型断言
func processValue(v interface{}) {
 str := v.(string) // 没有安全检查
 fmt.Println(str)
}

func main() {
 processOrder(1, []string{"book", "pen"}, 0.1, true, 1, nil)
 startWorker()
 incrementCounter()
 processValue("hello")
}
```

---

## 实战案例1: Web API服务分析

### 场景描述

你正在开发一个 RESTful API 服务，需要确保代码质量和并发安全。

### 第1步: 初始分析

```bash
fv analyze
```

**输出示例**:

```
========================================
📊 分析报告
========================================

项目: ./fv-tutorial
文件数: 1
总行数: 95
问题数: 8
质量评分: 62/100

----------------------------------------
问题统计:
  ❌ 错误: 3
  ⚠️  警告: 5
  ℹ️  提示: 0
----------------------------------------

❌ 错误:
  [concurrency] main.go:42:2
    Potential goroutine leak: channel never closed
    💡 建议: Add defer close(ch) after channel creation

  [concurrency] main.go:54:4
    Data race detected: unsynchronized access to shared variable
    💡 建议: Use sync.Mutex or atomic operations

  [type] main.go:67:10
    Unsafe type assertion without check
    💡 建议: Use v, ok := v.(string) pattern

⚠️  警告:
  [complexity] main.go:12:1
    Function processOrder has cyclomatic complexity 15 (threshold: 10)
    💡 建议: Break down into smaller functions

  [complexity] main.go:12:1
    Function processOrder has 6 parameters (threshold: 5)
    💡 建议: Consider using a struct for parameters
```

### 第2步: 生成详细报告

```bash
fv analyze --format=html --output=api-analysis.html
```

打开 HTML 报告，你会看到：

- 🎯 质量评分仪表板
- 📊 问题分布图表
- 🔍 每个问题的详细信息和代码位置

### 第3步: 修复问题

#### 修复 Goroutine 泄漏

**原代码**:

```go
func startWorker() {
 ch := make(chan int)
 
 go func() {
  for v := range ch {
   fmt.Println(v)
  }
 }()
 
 ch <- 1
 ch <- 2
}
```

**修复后**:

```go
func startWorker() {
 ch := make(chan int)
 defer close(ch) // 确保关闭channel
 
 go func() {
  for v := range ch {
   fmt.Println(v)
  }
 }()
 
 ch <- 1
 ch <- 2
 
 time.Sleep(10 * time.Millisecond) // 等待goroutine处理
}
```

#### 修复数据竞争

**原代码**:

```go
var counter int

func incrementCounter() {
 var wg sync.WaitGroup
 
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter++ // 数据竞争
  }()
 }
 
 wg.Wait()
}
```

**修复后**:

```go
type SafeCounter struct {
 mu    sync.Mutex
 value int
}

func (c *SafeCounter) Increment() {
 c.mu.Lock()
 defer c.mu.Unlock()
 c.value++
}

func incrementCounter() {
 var wg sync.WaitGroup
 counter := &SafeCounter{}
 
 for i := 0; i < 100; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   counter.Increment() // 安全的并发访问
  }()
 }
 
 wg.Wait()
}
```

#### 修复类型断言

**原代码**:

```go
func processValue(v interface{}) {
 str := v.(string) // 不安全
 fmt.Println(str)
}
```

**修复后**:

```go
func processValue(v interface{}) {
 str, ok := v.(string)
 if !ok {
  fmt.Println("Error: value is not a string")
  return
 }
 fmt.Println(str)
}
```

#### 重构复杂函数

**原代码**: 95行的 `processOrder` 函数

**修复后**:

```go
type Order struct {
 ID       int
 Items    []string
 Discount float64
 Shipping bool
 Priority int
 Customer map[string]string
}

func (o *Order) Calculate() (float64, error) {
 if err := o.validate(); err != nil {
  return 0, err
 }
 
 total := o.calculateItems()
 total = o.applyDiscount(total)
 total = o.addShipping(total)
 
 return total, nil
}

func (o *Order) validate() error {
 if o.ID <= 0 {
  return fmt.Errorf("invalid order ID")
 }
 if len(o.Items) == 0 {
  return fmt.Errorf("no items")
 }
 return nil
}

func (o *Order) calculateItems() float64 {
 prices := map[string]float64{
  "book":     10.0,
  "pen":      2.0,
  "notebook": 5.0,
  "laptop":   1000.0,
  "mouse":    20.0,
 }
 
 total := 0.0
 for _, item := range o.Items {
  price, ok := prices[item]
  if !ok {
   price = 1.0
  }
  total += price
 }
 return total
}

func (o *Order) applyDiscount(total float64) float64 {
 if o.Discount <= 0 {
  return total
 }
 
 discountRates := []struct {
  threshold float64
  rate      float64
 }{
  {0.3, 0.8},
  {0.2, 0.85},
  {0.1, 0.9},
  {0.0, 0.95},
 }
 
 for _, dr := range discountRates {
  if o.Discount >= dr.threshold {
   return total * dr.rate
  }
 }
 return total
}

func (o *Order) addShipping(total float64) float64 {
 if !o.Shipping {
  return total
 }
 
 shippingCosts := map[int]float64{
  1: 20.0,
  2: 10.0,
 }
 
 cost, ok := shippingCosts[o.Priority]
 if !ok {
  cost = 5.0
 }
 
 return total + cost
}
```

### 第4步: 重新分析

```bash
fv analyze
```

**新输出**:

```
========================================
📊 分析报告
========================================

项目: ./fv-tutorial
文件数: 1
总行数: 120
问题数: 0
质量评分: 98/100

----------------------------------------
✅ Excellent code quality!
----------------------------------------
```

### 第5步: 配置质量门槛

为了保持代码质量，创建配置文件：

```bash
fv init-config --output=.fv.yaml
```

编辑 `.fv.yaml`:

```yaml
output:
  fail_on_error: true
  min_quality_score: 85
```

现在运行：

```bash
fv analyze --config=.fv.yaml
```

如果质量低于85分，FV将返回非零退出码。

---

## 实战案例2: 并发密集型应用

### 场景描述

你正在开发一个高并发的数据处理服务，需要特别关注并发问题。

### 示例代码

创建 `worker_pool.go`:

```go
package main

import (
 "context"
 "fmt"
 "sync"
 "time"
)

// WorkerPool 工作池
type WorkerPool struct {
 workers   int
 taskQueue chan Task
 wg        sync.WaitGroup
 ctx       context.Context
 cancel    context.CancelFunc
}

// Task 任务接口
type Task interface {
 Execute() error
}

// NewWorkerPool 创建工作池
func NewWorkerPool(workers int) *WorkerPool {
 ctx, cancel := context.WithCancel(context.Background())
 
 return &WorkerPool{
  workers:   workers,
  taskQueue: make(chan Task, 100),
  ctx:       ctx,
  cancel:    cancel,
 }
}

// Start 启动工作池
func (wp *WorkerPool) Start() {
 for i := 0; i < wp.workers; i++ {
  wp.wg.Add(1)
  go wp.worker(i)
 }
}

// worker 工作协程
func (wp *WorkerPool) worker(id int) {
 defer wp.wg.Done()
 
 for {
  select {
  case <-wp.ctx.Done():
   fmt.Printf("Worker %d stopping\n", id)
   return
  case task, ok := <-wp.taskQueue:
   if !ok {
    return
   }
   if err := task.Execute(); err != nil {
    fmt.Printf("Worker %d: task failed: %v\n", id, err)
   }
  }
 }
}

// Submit 提交任务
func (wp *WorkerPool) Submit(task Task) error {
 select {
 case <-wp.ctx.Done():
  return fmt.Errorf("worker pool is shutting down")
 case wp.taskQueue <- task:
  return nil
 }
}

// Shutdown 关闭工作池
func (wp *WorkerPool) Shutdown(timeout time.Duration) error {
 // 停止接收新任务
 wp.cancel()
 
 // 关闭任务队列
 close(wp.taskQueue)
 
 // 等待所有worker完成
 done := make(chan struct{})
 go func() {
  wp.wg.Wait()
  close(done)
 }()
 
 // 等待超时
 select {
 case <-done:
  return nil
 case <-time.After(timeout):
  return fmt.Errorf("shutdown timeout")
 }
}
```

### 分析并发问题

```bash
# 使用严格模式配置
fv init-config --output=.fv-concurrency.yaml --strict

# 编辑配置，启用所有并发检查
# rules:
#   concurrency:
#     enabled: true
#     check_goroutine_leak: true
#     check_data_race: true
#     check_deadlock: true
#     check_channel: true

# 运行分析
fv analyze --config=.fv-concurrency.yaml
```

FV 会检查：

- ✅ Goroutine是否正确关闭
- ✅ Channel是否正确关闭
- ✅ 是否存在数据竞争
- ✅ 是否可能死锁
- ✅ Context是否正确使用

---

## 实战案例3: 遗留代码重构

### 场景描述

你继承了一个遗留项目，需要评估代码质量并制定重构计划。

### 第1步: 全面扫描

```bash
# 包含测试文件
fv analyze --include-tests --format=json --output=legacy-analysis.json
```

### 第2步: 分析结果

```bash
# 查看质量分数
jq -r '.stats.quality_score' legacy-analysis.json

# 统计问题类型
jq -r '.stats' legacy-analysis.json

# 找出问题最多的文件
jq -r '.issues[] | .file' legacy-analysis.json | sort | uniq -c | sort -rn | head -10
```

### 第3步: 制定重构计划

基于分析结果，按优先级处理：

1. **高优先级**: 错误级别的问题

   ```bash
   jq -r '.issues[] | select(.severity=="error")' legacy-analysis.json
   ```

2. **中优先级**: 高复杂度函数

   ```bash
   jq -r '.issues[] | select(.category=="complexity")' legacy-analysis.json
   ```

3. **低优先级**: 优化建议

   ```bash
   jq -r '.issues[] | select(.severity=="info")' legacy-analysis.json
   ```

### 第4步: 跟踪进度

创建基准报告：

```bash
# 初始状态
fv analyze --format=json --output=baseline.json

# 重构后
fv analyze --format=json --output=current.json

# 比较质量分数
echo "Baseline: $(jq -r '.stats.quality_score' baseline.json)"
echo "Current:  $(jq -r '.stats.quality_score' current.json)"
```

---

## 高级功能

### 1. 自定义规则阈值

```yaml
# .fv-custom.yaml
rules:
  complexity:
    cyclomatic_threshold: 5      # 更严格
    max_function_lines: 30
    max_parameters: 3
  
  performance:
    enabled: true
    check_allocation: true
```

### 2. 选择性分析

```bash
# 只分析特定包
fv analyze --dir=./api

# 排除生成的代码
fv analyze --exclude="*_gen.go,*.pb.go"
```

### 3. 集成到Git Hooks

创建 `.git/hooks/pre-commit`:

```bash
#!/bin/bash
echo "Running FV analysis..."

fv analyze --config=.fv.yaml --no-color --fail-on-error

if [ $? -ne 0 ]; then
    echo "❌ FV analysis failed. Please fix the issues before committing."
    exit 1
fi

echo "✅ FV analysis passed"
```

### 4. 生成徽章

```bash
# 生成质量分数
SCORE=$(fv analyze --format=json | jq -r '.stats.quality_score')

# 创建徽章
echo "[![FV Quality](https://img.shields.io/badge/FV%20Quality-${SCORE}%25-green)](./fv-report.html)"
```

---

## 故障排查

### 问题1: 误报

**现象**: FV报告了不存在的问题

**解决**:

```yaml
# 调整检查灵敏度
rules:
  concurrency:
    check_goroutine_leak: false  # 如果误报过多
```

### 问题2: 性能问题

**现象**: 分析大项目很慢

**解决**:

```yaml
analysis:
  workers: 8               # 增加并发数
  max_file_size: 512      # 跳过大文件
  timeout: 600            # 增加超时时间
```

### 问题3: 配置不生效

**现象**: 修改配置后没有变化

**解决**:

```bash
# 确保使用了正确的配置文件
fv analyze --config=.fv.yaml -v

# 验证配置文件格式
cat .fv.yaml | yaml-lint
```

---

## 总结

通过本教程，你学会了：

1. ✅ 如何分析不同类型的Go项目
2. ✅ 如何解读FV报告
3. ✅ 如何修复常见问题
4. ✅ 如何配置FV以适应需求
5. ✅ 如何集成到开发流程

### 下一步

- 🔧 将FV集成到你的CI/CD流程
- 📚 阅读[最佳实践](Best-Practices.md)
- 🚀 开始提升代码质量！

---

**Happy Coding!** 🎉
