# 附录: 参考资源与速查表

> Go 1.26 完全技术指南的参考集合

---

## 一、快速参考

### 1.1 Go 1.26 新特性速查

```
特性                语法                          版本    状态
────────────────────────────────────────────────────────
new(expr)           new(Some{Field: value})       1.26    正式
递归泛型约束         type N[T N[T]] interface{}    1.26    正式
Green Tea GC        自动启用                      1.26    默认
Goroutine Leak检测  runtime.SetLeakCallback       1.26    新API
```

### 1.2 类型系统速查

```
类型分类:
├── 基础类型: bool, int, float, string, complex
├── 复合类型: array, slice, map, struct, pointer, function, interface, channel
├── 类型参数: type Parameter[T Constraint]
└── 类型别名: type Alias = Original

接口实现规则:
├── 隐式实现: 无需声明 implements
├── 方法集: 值方法集 ⊂ 指针方法集
└── 空接口: interface{} = any (1.18+)
```

### 1.3 并发原语速查

```
Goroutine管理:
├── go func()        - 启动goroutine
├── runtime.GOMAXPROCS(n) - 设置P数量
└── runtime.NumGoroutine() - 当前goroutine数

Channel操作:
├── ch <- v          - 发送
├── v := <-ch        - 接收
├── close(ch)        - 关闭
└── select {case...} - 多路复用

同步原语:
├── sync.Mutex       - 互斥锁
├── sync.RWMutex     - 读写锁
├── sync.WaitGroup   - 等待组
├── sync.Once        - 一次性执行
├── sync.Pool        - 对象池
└── sync/atomic      - 原子操作

Context控制:
├── context.Background() - 根context
├── context.TODO()       - 占位context
├── context.WithCancel() - 可取消
├── context.WithTimeout() - 超时控制
└── context.WithValue()   - 传值
```

---

## 二、学术参考

### 2.1 经典论文

```
语言基础:
├─ Hoare, C.A.R. "Communicating Sequential Processes" (1978, 1985)
├─ Pike, R. "Go at Google: Language Design in the Service of Software Engineering"
├─ Griesemer et al. "Featherweight Go" (POPL 2020)
└─ Griesemer et al. "Hoare Logic for Go"

类型系统:
├─ Cardelli, L. "Type Systems" (ACM Computing Surveys)
├─ Pierce, B.C. "Types and Programming Languages" (教科书)
└─ Kennedy, A. "Types for Units-of-Measure"

并发理论:
├─ Herlihy, M. & Shavit, N. "The Art of Multiprocessor Programming"
├─ Go Memory Model Specification
└─ Adve & Gharachorloo "Shared Memory Consistency Models"

分布式:
├─ Lamport, L. "Time, Clocks, and the Ordering of Events"
├─ Lamport, L. "Paxos Made Simple"
└─ Burrows, M. "The Chubby Lock Service"
```

### 2.2 推荐课程

```
在线资源:
├─ Stanford CS242: Programming Languages
├─ MIT 6.006: Introduction to Algorithms
├─ MIT 6.824: Distributed Systems
├─ CMU 15-312: Foundations of Programming Languages
├─ "Software Foundations" (Coq教程)
└─ "Functional Programming in Lean"

书籍:
├─ "The Go Programming Language" (Donovan & Kernighan)
├─ "Concurrency in Go" (Katherine Cox-Buday)
├─ "Learning Go" (Jon Bodner)
└─ "100 Go Mistakes" (Teiva Harsanyi)
```

---

## 三、工具参考

### 3.1 开发工具

```
编辑器:
├─ VS Code + Go扩展
├─ GoLand (JetBrains)
├─ Vim/Neovim + vim-go
└─ Emacs + go-mode

构建工具:
├─ go build (标准)
├─ Makefile
├─ Bazel
└─ Mage (Go-based)

测试工具:
├─ go test (标准)
├─ testify (断言库)
├─ ginkgo (BDD框架)
└─ gomock (Mock生成)
```

### 3.2 性能工具

```
分析工具:
├─ go test -bench=.        - 基准测试
├─ go test -benchmem       - 内存分析
├─ go test -cpuprofile     - CPU分析
├─ go test -memprofile     - 堆分析
├─ go tool pprof           - 可视化分析
└─ go tool trace           - 执行追踪

运行时:
├─ runtime.ReadMemStats()  - 内存统计
├─ runtime/pprof           - 运行时分析
└─ net/http/pprof          - HTTP分析端点
```

### 3.3 验证工具

```
静态分析:
├─ go vet                   - 基础检查
├─ staticcheck              - 高级分析
├─ golangci-lint            - 综合检查
├─ gosec                    - 安全检查
└─ govulncheck              - 漏洞检查

并发检测:
├─ go test -race            - 数据竞争
├─ go test -deadlock (实验) - 死锁检测
└─ 1.26 leak检测            - goroutine泄露

形式化:
├─ go contracts (实验)      - 契约检查
├─ Iris-Go (研究)           - 分离逻辑
└─ Promela建模              - 模型检测
```

---

## 四、标准库速查

### 4.1 核心包

```
基础:
├── fmt      - 格式化I/O
├── os       - 操作系统接口
├── io       - I/O原语
├── bufio    - 缓冲I/O
├── bytes    - 字节切片操作
├── strings  - 字符串操作
├── strconv  - 字符串转换
├── math     - 数学函数
├── time     - 时间处理
├── errors   - 错误处理
└── context  - 请求上下文

数据:
├── sort     - 排序算法
├── container/heap    - 堆实现
├── container/list    - 双向链表
├── container/ring    - 环形链表
├── encoding/json     - JSON编解码
├── encoding/xml      - XML编解码
└── database/sql      - 数据库接口

网络:
├── net      - 网络I/O
├── net/http - HTTP客户端/服务器
├── net/rpc  - RPC框架
└── html/template - HTML模板

并发:
├── sync     - 同步原语
├── sync/atomic - 原子操作
└── context  - Context
```

---

## 五、常见模式

### 5.1 错误处理模式

```go
// 包装错误
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// 错误类型判断
if errors.Is(err, sql.ErrNoRows) {
    // 处理未找到
}

var notFound *NotFoundError
if errors.As(err, &notFound) {
    // 处理特定错误类型
}
```

### 5.2 并发模式

```go
// Worker Pool
func workerPool(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    for j := range jobs {
        results <- process(j)
    }
    wg.Done()
}

// Pipeline
func pipeline(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            out <- process(v)
        }
    }()
    return out
}

// Fan-out/Fan-in
func fanOut(input <-chan int, n int) []<-chan int {
    outputs := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        outputs[i] = processChannel(input)
    }
    return outputs
}
```

---

## 六、版本历史

```
Go版本演进:
Go 1.0  (2012) - 初始稳定版
Go 1.1  (2013) - 性能改进
Go 1.5  (2015) - 自举编译器
Go 1.11 (2018) - Modules引入
Go 1.13 (2019) - 错误包装
Go 1.18 (2022) - 泛型
Go 1.20 (2023) - PGO, 错误连接
Go 1.21 (2023) - 内置min/max/clear
Go 1.22 (2024) - 循环变量语义修复
Go 1.23 (2024) - iter包, 改进panic信息
Go 1.24 (2025) - 泛型类型别名
Go 1.25 (2025) - 运行时优化
Go 1.26 (2026) - new(), 递归泛型, Green Tea GC
```

---

*附录提供了快速参考资源，支持日常开发和学习查询。*
