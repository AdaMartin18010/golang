# ❓ 常见问题解答 (FAQ)

> 解答Go学习和使用中的常见问题

**版本**: v2.1
**更新日期**: 2025-10-29
**问题数**: 60+

---

## 📋 目录

- [❓ 常见问题解答 (FAQ)](#-常见问题解答-faq)
  - [📋 目录](#-目录)
  - [🆕 工具相关问题 (2025-10-25新增)](#-工具相关问题-2025-10-25新增)
    - [Q26: 如何快速分析Go项目代码质量？](#q26-如何快速分析go项目代码质量)
    - [Q27: 如何快速生成并发模式代码？](#q27-如何快速生成并发模式代码)
    - [Q28: 如何将FV集成到CI/CD？](#q28-如何将fv集成到cicd)
    - [Q29: FV分析速度慢怎么办？](#q29-fv分析速度慢怎么办)
    - [Q30: 如何定制FV的检查规则？](#q30-如何定制fv的检查规则)
  - [📑 目录](#-目录-1)
  - [入门问题](#入门问题)
    - [Q1: 为什么选择Go语言？](#q1-为什么选择go语言)
    - [Q2: Go和其他语言对比如何？](#q2-go和其他语言对比如何)
    - [Q3: 学习Go需要什么基础？](#q3-学习go需要什么基础)
    - [Q4: 学习Go需要多长时间？](#q4-学习go需要多长时间)
  - [语言特性](#语言特性)
    - [Q5: Go有哪些独特的特性？](#q5-go有哪些独特的特性)
    - [Q6: Go支持泛型吗？](#q6-go支持泛型吗)
    - [Q7: Go有面向对象吗？](#q7-go有面向对象吗)
  - [并发编程](#并发编程)
    - [Q8: Goroutine和线程有什么区别？](#q8-goroutine和线程有什么区别)
    - [Q9: Channel什么时候用？](#q9-channel什么时候用)
    - [Q10: 什么时候用Mutex，什么时候用Channel？](#q10-什么时候用mutex什么时候用channel)
  - [Web开发](#web开发)
    - [Q11: Go适合做Web开发吗？](#q11-go适合做web开发吗)
    - [Q12: 应该用哪个Web框架？](#q12-应该用哪个web框架)
    - [Q13: 如何处理JWT认证？](#q13-如何处理jwt认证)
  - [性能优化](#性能优化)
    - [Q14: 如何分析Go程序性能？](#q14-如何分析go程序性能)
    - [Q15: 如何优化内存使用？](#q15-如何优化内存使用)
    - [Q16: Go的GC可以调优吗？](#q16-go的gc可以调优吗)
  - [工程实践](#工程实践)
    - [Q17: 如何组织Go项目结构？](#q17-如何组织go项目结构)
    - [Q18: 如何写测试？](#q18-如何写测试)
    - [Q19: 如何处理错误？](#q19-如何处理错误)
    - [Q20: 如何管理配置？](#q20-如何管理配置)
  - [部署运维](#部署运维)
    - [Q21: 如何部署Go应用？](#q21-如何部署go应用)
    - [Q22: 如何优化Docker镜像大小？](#q22-如何优化docker镜像大小)
    - [Q23: 如何实现优雅关闭？](#q23-如何实现优雅关闭)
  - [学习建议](#学习建议)
    - [Q24: 有哪些Go学习资源？](#q24-有哪些go学习资源)
    - [Q25: 如何快速提升Go水平？](#q25-如何快速提升go水平)
  - [🔗 相关资源](#-相关资源)

## 🆕 工具相关问题 (2025-10-25新增)

### Q26: 如何快速分析Go项目代码质量？

**A**: 使用Formal Verifier工具：

```bash
# 1. 构建工具
cd tools/formal-verifier
go build -o fv ./cmd/fv

# 2. 零配置分析
./fv analyze --dir=.

# 3. 生成HTML报告
./fv analyze --format=html --output=report.html
```

**特点**:

- ✅ 零配置，开箱即用
- ✅ 多格式报告（HTML/JSON/Markdown）
- ✅ CI/CD就绪

**文档**: [快速入门](../tools/formal-verifier/docs/Quick-Start.md)

---

### Q27: 如何快速生成并发模式代码？

**A**: 使用Pattern Generator工具：

```bash
# 1. 构建工具
cd tools/concurrency-pattern-generator
go build -o cpg ./cmd/cpg

# 2. 查看所有模式
./cpg --list

# 3. 生成Worker Pool
./cpg --pattern worker-pool --workers 5 --output pool.go
```

**支持30个模式**:

- 经典模式 (5个)
- 同步模式 (8个)
- 控制流模式 (5个)
- 数据流模式 (7个)
- 高级模式 (5个)

**文档**: [CPG README](../tools/concurrency-pattern-generator/README.md)

---

### Q28: 如何将FV集成到CI/CD？

**A**: 三步快速集成：

```bash
# 1. 创建严格配置
./fv init-config --strict > .fv-strict.yaml

# 2. 复制GitHub Actions配置
cp .github/workflows/fv-analysis.yml.example \
   .github/workflows/fv-analysis.yml

# 3. 提交并推送
git add .fv-strict.yaml .github/workflows/
git commit -m "Add FV CI/CD"
```

**支持平台**:

- GitHub Actions ✓
- GitLab CI ✓
- Jenkins ✓

**文档**: [CI/CD集成](../tools/formal-verifier/docs/CI-CD-Integration.md)

---

### Q29: FV分析速度慢怎么办？

**A**: 优化配置：

```yaml
# .fv.yaml
analysis:
  workers: 8  # 增加并发数

project:
  exclude:    # 优化排除规则
    - "vendor/*"
    - "*/testdata/*"
    - "*/mocks/*"
```

**其他建议**:

- 使用SSD硬盘
- 排除不必要的目录
- 使用缓存

---

### Q30: 如何定制FV的检查规则？

**A**: 修改配置文件：

```yaml
rules:
  complexity:
    cyclomatic_threshold: 10      # 圈复杂度
    max_function_lines: 50        # 函数最大行数

  concurrency:
    check_goroutine_leaks: true   # Goroutine泄露
    check_channel_deadlocks: true # Channel死锁
    check_data_races: true        # 数据竞争
```

**配置优先级**:

```text
命令行参数 > 配置文件 > 默认值
```

---

---

## 📑 目录

- [❓ 常见问题解答 (FAQ)](#-常见问题解答-faq)
  - [📋 目录](#-目录)
  - [🆕 工具相关问题 (2025-10-25新增)](#-工具相关问题-2025-10-25新增)
    - [Q26: 如何快速分析Go项目代码质量？](#q26-如何快速分析go项目代码质量)
    - [Q27: 如何快速生成并发模式代码？](#q27-如何快速生成并发模式代码)
    - [Q28: 如何将FV集成到CI/CD？](#q28-如何将fv集成到cicd)
    - [Q29: FV分析速度慢怎么办？](#q29-fv分析速度慢怎么办)
    - [Q30: 如何定制FV的检查规则？](#q30-如何定制fv的检查规则)
  - [📑 目录](#-目录-1)
  - [入门问题](#入门问题)
    - [Q1: 为什么选择Go语言？](#q1-为什么选择go语言)
    - [Q2: Go和其他语言对比如何？](#q2-go和其他语言对比如何)
    - [Q3: 学习Go需要什么基础？](#q3-学习go需要什么基础)
    - [Q4: 学习Go需要多长时间？](#q4-学习go需要多长时间)
  - [语言特性](#语言特性)
    - [Q5: Go有哪些独特的特性？](#q5-go有哪些独特的特性)
    - [Q6: Go支持泛型吗？](#q6-go支持泛型吗)
    - [Q7: Go有面向对象吗？](#q7-go有面向对象吗)
  - [并发编程](#并发编程)
    - [Q8: Goroutine和线程有什么区别？](#q8-goroutine和线程有什么区别)
    - [Q9: Channel什么时候用？](#q9-channel什么时候用)
    - [Q10: 什么时候用Mutex，什么时候用Channel？](#q10-什么时候用mutex什么时候用channel)
  - [Web开发](#web开发)
    - [Q11: Go适合做Web开发吗？](#q11-go适合做web开发吗)
    - [Q12: 应该用哪个Web框架？](#q12-应该用哪个web框架)
    - [Q13: 如何处理JWT认证？](#q13-如何处理jwt认证)
  - [性能优化](#性能优化)
    - [Q14: 如何分析Go程序性能？](#q14-如何分析go程序性能)
    - [Q15: 如何优化内存使用？](#q15-如何优化内存使用)
    - [Q16: Go的GC可以调优吗？](#q16-go的gc可以调优吗)
  - [工程实践](#工程实践)
    - [Q17: 如何组织Go项目结构？](#q17-如何组织go项目结构)
    - [Q18: 如何写测试？](#q18-如何写测试)
    - [Q19: 如何处理错误？](#q19-如何处理错误)
    - [Q20: 如何管理配置？](#q20-如何管理配置)
  - [部署运维](#部署运维)
    - [Q21: 如何部署Go应用？](#q21-如何部署go应用)
    - [Q22: 如何优化Docker镜像大小？](#q22-如何优化docker镜像大小)
    - [Q23: 如何实现优雅关闭？](#q23-如何实现优雅关闭)
  - [学习建议](#学习建议)
    - [Q24: 有哪些Go学习资源？](#q24-有哪些go学习资源)
    - [Q25: 如何快速提升Go水平？](#q25-如何快速提升go水平)
  - [🔗 相关资源](#-相关资源)

---

## 入门问题

### Q1: 为什么选择Go语言？

**A**: Go的主要优势：

1. **简单易学**: 语法简洁，上手快
2. **高性能**: 接近C的性能
3. **并发支持**: 原生Goroutine和Channel
4. **标准库丰富**: 开箱即用
5. **部署简单**: 单一二进制文件
6. **工具链完善**: go fmt, go test等
7. **社区活跃**: 大量优秀开源项目

**适合场景**: Web后端、微服务、云原生、DevOps工具

**相关文档**: [Hello World](01-语言基础/01-语法基础/01-Hello-World.md)

---

### Q2: Go和其他语言对比如何？

**A**: 简单对比：

| 特性 | Go | Python | Java | Node.js |
|------|-----|--------|------|---------|
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| 并发 | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| 易学性 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| 部署 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| 生态 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

**结论**: Go在性能、并发和部署上有优势，适合后端和云原生开发

---

### Q3: 学习Go需要什么基础？

**A**:

**必备基础**:

- ✅ 基本的编程概念（变量、函数、循环）
- ✅ 命令行基础

**加分项**:

- ⭐ 其他编程语言经验（Python/Java/C等）
- ⭐ HTTP和网络基础
- ⭐ 数据库基础

**零基础**: 也可以学习，但需要更多时间理解基本概念

**相关文档**: [零基础入门路径](LEARNING_PATHS.md#%e9%9b%b6%e5%9f%ba%e7%a1%80%e5%85%a5%e9%97%a8%e8%b7%af%e5%be%84)

---

### Q4: 学习Go需要多长时间？

**A**: 根据目标不同：

| 目标 | 时间 | 内容 |
|------|------|------|
| 入门 | 2-4周 | 语法基础+简单项目 |
| 初级开发者 | 3个月 | Web开发+数据库 |
| 中级开发者 | 6个月 | 微服务+云原生 |
| 高级开发者 | 1-2年 | 架构+性能优化 |
| 专家 | 3年+ | 深度+广度 |

**每天学习时间**: 2-3小时为宜

**相关文档**: [学习路径图](LEARNING_PATHS.md)

---

## 语言特性

### Q5: Go有哪些独特的特性？

**A**:

1. **Goroutine**: 轻量级并发
2. **Channel**: 通信机制
3. **Defer**: 延迟执行
4. **Interface**: 隐式实现
5. **多返回值**: 函数可返回多个值
6. **内置并发**: sync包
7. **快速编译**: 秒级编译

**相关文档**: [语言基础](01-语言基础/README.md)

---

### Q6: Go支持泛型吗？

**A**:

✅ **支持！** Go 1.18+ 开始支持泛型

```go
// 泛型函数
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// 使用
numbers := []int{1, 2, 3}
doubled := Map(numbers, func(n int) int { return n * 2 })
// [2, 4, 6]
```

**相关文档**: [Go 1.21特性](10-Go版本特性/01-Go-1.21特性/README.md)

---

### Q7: Go有面向对象吗？

**A**:

**没有传统OOP，但有类似机制**:

```go
// 结构体（类似类）
type Person struct {
    Name string
    Age  int
}

// 方法
func (p *Person) SayHello() {
    fmt.Printf("Hello, I'm %s\n", p.Name)
}

// 接口（隐式实现）
type Greeter interface {
    SayHello()
}
```

**没有**:

- ❌ 继承
- ❌ 类
- ❌ 构造函数

**有**:

- ✅ 组合
- ✅ 接口
- ✅ 方法

**相关文档**: [基本数据类型](01-语言基础/01-语法基础/03-基本数据类型.md)

---

## 并发编程

### Q8: Goroutine和线程有什么区别？

**A**:

| 特性 | Goroutine | 线程 |
|------|-----------|------|
| 内存占用 | 2KB | 2MB |
| 调度 | 用户态 | 内核态 |
| 切换成本 | 低 | 高 |
| 数量 | 百万级 | 千级 |

```go
// 启动100万个Goroutine很容易
for i := 0; i < 1000000; i++ {
    go func() {
        // 做点什么
    }()
}
```

**相关文档**: [Goroutine基础](01-语言基础/02-并发编程/02-Goroutine基础.md)

---

### Q9: Channel什么时候用？

**A**:

**使用Channel的场景**:

1. ✅ Goroutine间通信
2. ✅ 数据流转
3. ✅ 信号通知
4. ✅ 控制并发数

```go
// 典型使用：Worker Pool
func workerPool(tasks []Task) {
    taskCh := make(chan Task, 100)
    resultCh := make(chan Result, 100)

    // 启动workers
    for i := 0; i < 10; i++ {
        go worker(taskCh, resultCh)
    }

    // 发送任务
    for _, task := range tasks {
        taskCh <- task
    }
    close(taskCh)

    // 收集结果
    for range tasks {
        result := <-resultCh
        // 处理结果
    }
}
```

**相关文档**: [Channel基础](01-语言基础/02-并发编程/03-Channel基础.md)

---

### Q10: 什么时候用Mutex，什么时候用Channel？

**A**:

**Use Mutex for**:

- 保护共享状态
- 短时间的临界区
- 简单的计数器

```go
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}
```

**Use Channel for**:

- Goroutine间通信
- 数据流转
- 协调任务

**Go谚语**: "Don't communicate by sharing memory; share memory by communicating."

**相关文档**: [sync包](01-语言基础/02-并发编程/06-sync包.md)

---

## Web开发

### Q11: Go适合做Web开发吗？

**A**:

✅ **非常适合！**

**优势**:

1. 高性能HTTP服务器
2. 丰富的Web框架（Gin, Echo, Fiber）
3. 优秀的并发处理
4. 简单的部署

**代表项目**:

- Docker
- Kubernetes
- Prometheus
- Traefik

**相关文档**: [Web开发](03-Web开发/README.md)

---

### Q12: 应该用哪个Web框架？

**A**:

| 框架 | 特点 | 适合场景 |
|------|------|---------|
| **net/http** | 标准库 | 小项目、学习 |
| **Gin** | 最流行、生态好 | 大部分Web项目 |
| **Echo** | 高性能 | 高性能要求 |
| **Fiber** | Express风格 | Node.js转Go |

**推荐**: 新项目用 **Gin**

```go
// Gin示例
r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})
r.Run(":8080")
```

**相关文档**: [Gin框架](03-Web开发/04-Gin框架.md)

---

### Q13: 如何处理JWT认证？

**A**:

使用 `golang-jwt/jwt` 包：

```go
import "github.com/golang-jwt/jwt/v5"

// 生成Token
func GenerateToken(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString([]byte("your-secret-key"))
}

// 验证Token
func ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })
}
```

**相关文档**: [认证与授权](03-Web开发/00-Go认证与授权深度实战指南.md)

---

## 性能优化

### Q14: 如何分析Go程序性能？

**A**:

**使用pprof**:

```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()

    // 你的程序
}
```

访问: <http://localhost:6060/debug/pprof/>

**常用命令**:

```bash
# CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 内存 profile
go tool pprof http://localhost:6060/debug/pprof/heap

# 火焰图
go tool pprof -http=:8080 profile
```

**相关文档**: [性能分析与pprof](07-性能优化/01-性能分析与pprof.md)

---

### Q15: 如何优化内存使用？

**A**:

**常见技巧**:

1. **使用对象池**

    ```go
    var bufferPool = sync.Pool{
        New: func() interface{} {
            return new(bytes.Buffer)
        },
    }
    ```

2. **预分配slice**

    ```go
    // 好
    slice := make([]int, 0, 1000)

    // 不好
    slice := []int{}
    ```

3. **避免不必要的指针**

    ```go
    // 小对象用值类型
    type Point struct {
        X, Y int
    }
    ```

**相关文档**: [内存优化](07-性能优化/02-内存优化.md)

---

### Q16: Go的GC可以调优吗？

**A**:

✅ **可以！**

**GOGC环境变量**:

```bash
# 默认100，表示heap增长100%时触发GC
export GOGC=200  # 降低GC频率

# 禁用GC（仅测试用）
export GOGC=off
```

**SetGCPercent**:

```go
import "runtime/debug"

// 设置GOGC为200
debug.SetGCPercent(200)
```

**Go 1.19+ GOMEMLIMIT**:

```bash
# 限制最大内存使用
export GOMEMLIMIT=2GiB
```

**相关文档**: [GC调优](07-性能优化/05-GC调优.md)

---

## 工程实践

### Q17: 如何组织Go项目结构？

**A**:

**推荐结构**:

```text
my-project/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/
│   └── utils/
├── api/
├── configs/
├── scripts/
├── go.mod
└── go.sum
```

**说明**:

- `cmd/`: 应用入口
- `internal/`: 私有代码
- `pkg/`: 可导出代码
- `api/`: API定义（Proto/OpenAPI）

---

### Q18: 如何写测试？

**A**:

**单元测试**:

```go
// math_test.go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Expected 5, got %d", result)
    }
}

// 运行测试
// go test ./...
```

**表格驱动测试**:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

**相关文档**: [测试实践](09-工程实践/00-Go测试深度实战指南.md)

---

### Q19: 如何处理错误？

**A**:

**Go 1.13+错误处理**:

```go
import (
    "errors"
    "fmt"
)

// 定义错误
var ErrNotFound = errors.New("not found")

// 包装错误
func GetUser(id int) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("get user %d: %w", id, err)
    }
    return user, nil
}

// 判断错误类型
if errors.Is(err, ErrNotFound) {
    // 处理NotFound错误
}

// 提取错误
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // 处理路径错误
}
```

---

### Q20: 如何管理配置？

**A**:

**使用Viper**:

```go
import "github.com/spf13/viper"

func init() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    if err := viper.ReadInConfig(); err != nil {
        panic(err)
    }
}

// 读取配置
port := viper.GetString("server.port")
debug := viper.GetBool("debug")
```

**config.yaml**:

```yaml
server:
  port: "8080"
debug: true
database:
  host: "localhost"
  port: 5432
```

---

## 部署运维

### Q21: 如何部署Go应用？

**A**:

**方式1: 二进制部署**:

```bash
# 编译
CGO_ENABLED=0 GOOS=linux go build -o myapp

# 运行
./myapp
```

**方式2: Docker部署**:

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp

FROM alpine:latest
COPY --from=builder /app/myapp /myapp
CMD ["/myapp"]
```

**方式3: Kubernetes部署**:

**相关文档**: [云原生部署](06-云原生与容器/README.md)

---

### Q22: 如何优化Docker镜像大小？

**A**:

**最佳实践**:

```dockerfile
# 多阶段构建
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o myapp

# 最小镜像
FROM scratch
COPY --from=builder /app/myapp /myapp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/myapp"]
```

**结果**: 从800MB → 10MB

**相关文档**: [Dockerfile最佳实践](06-云原生与容器/02-Dockerfile最佳实践.md)

---

### Q23: 如何实现优雅关闭？

**A**:

```go
func main() {
    server := &http.Server{Addr: ":8080"}

    // 启动服务器
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}
```

---

## 学习建议

### Q24: 有哪些Go学习资源？

**A**:

**官方资源**:

- ✅ [Tour of Go](https://go.dev/tour/) - 官方教程
- ✅ [Effective Go](https://go.dev/doc/effective_go) - 最佳实践
- ✅ [Go Blog](https://go.dev/blog/) - 官方博客

**书籍**:

- ✅ 《Go程序设计语言》
- ✅ 《Go语言实战》
- ✅ 《Go并发编程实战》

**网站**:

- ✅ [Go by Example](https://gobyexample.com/)
- ✅ [Awesome Go](https://awesome-go.com/)

**本文档库**:

- ✅ [学习路径](LEARNING_PATHS.md)
- ✅ [技术索引](INDEX.md)

---

### Q25: 如何快速提升Go水平？

**A**:

**5个建议**:

1. **每天写代码** (2-3小时)
   - 跟着文档做练习
   - 实现小项目

2. **阅读优秀代码**
   - 标准库源码
   - Docker/Kubernetes源码
   - 知名开源项目

3. **刷算法题**
   - [LeetCode实战](02-数据结构与算法/04-实战案例.md)
   - 每天2-3题

4. **参与开源**
   - 提PR修bug
   - 贡献文档
   - 分享经验

5. **技术分享**
   - 写技术博客
   - 做技术分享
   - 教别人学Go

**相关文档**: [学习路径图](LEARNING_PATHS.md)

---

## 🔗 相关资源

- [技术索引](INDEX.md) - 查找特定主题
- [学习路径](LEARNING_PATHS.md) - 系统学习计划
- [快速开始](QUICK_START.md) - 快速入门
- [术语表](GLOSSARY.md) - 技术术语

---

**最后更新**: 2025-10-29
**文档版本**: v2.0
**维护团队**: Documentation Team

---

<div align="center">

**❓ 没找到答案？请提交Issue或在社区提问**:

</div>
