# 工程最佳实践

> 基于形式化原则的生产级Go开发指南

---

## 一、代码组织原则

### 1.1 项目结构公理

```
公理: 代码结构应反映领域模型
─────────────────────────────
证明:
  若结构反映领域，则:
  1. 新人可通过目录理解业务
  2. 修改局限在相关包内
  3. 测试可对应业务场景
∴ 内聚度高，耦合度低

标准项目结构:
project/
├── cmd/              # 入口点
│   ├── api/          # HTTP API服务
│   └── worker/       # 后台任务
├── internal/         # 私有代码
│   ├── domain/       # 领域模型 (核心业务)
│   ├── application/  # 应用服务 (编排)
│   ├── infrastructure/# 技术实现
│   └── interfaces/   # 接口适配
├── pkg/              # 可复用库
└── api/              # API定义(proto/openapi)

原则: internal保护领域，pkg共享通用
```

### 1.2 包设计原则

```
包内聚性公理:
─────────────────────────────
若包内所有类型服务于单一变化原因，
则该包具有高内聚性。

反模式检测:
├── utils包 - 违背单一职责
├── common包 - 成为垃圾场
└── 按层分包 - 导致循环依赖

正确做法:
├── 按领域分包: user/, order/, payment/
├── 接口定义在消费方
└── 依赖注入解耦

示例:
// 不良: 通用错误包
package errors

// 良好: 领域特定错误
package user

type NotFoundError struct { UserID string }
func (e NotFoundError) Error() string { ... }
```

---

## 二、并发最佳实践

### 2.1 Goroutine生命周期管理

```
定理: 每个goroutine必须有明确的退出路径
─────────────────────────────
违反后果: goroutine泄露 → 内存耗尽

Go 1.26增强:
• runtime.SetGoroutineLeakCallback
• pprof leak检测模式

实践模式:
1. Context取消传播
2. WaitGroup等待
3. Error Group错误聚合
4. Graceful shutdown

代码模板:
func worker(ctx context.Context, jobs <-chan Job) error {
    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                return nil // 正常退出
            }
            if err := process(job); err != nil {
                return err
            }
        case <-ctx.Done():
            return ctx.Err() // 取消退出
        }
    }
}
```

### 2.2 Channel使用规范

```
公理: Channel所有权决定关闭责任
─────────────────────────────
规则:
- 发送方不应关闭channel
- 只有接收方知道何时不再需要数据
- 或创建者负责生命周期

最佳实践:
┌────────────────────────────────────────┐
│ 1. 优先使用for-range接收               │
│    for v := range ch { ... }           │
│                                        │
│ 2. 检查ok判断channel关闭                │
│    v, ok := <-ch                       │
│                                        │
│ 3. 使用select处理多个channel            │
│    select { case v:=<-ch1: ... }       │
│                                        │
│ 4. 带缓冲channel用于解耦                │
│    ch := make(chan T, bufferSize)      │
│                                        │
│ 5. nil channel在select中禁用            │
│    用于动态启用/禁用case                 │
└────────────────────────────────────────┘
```

---

## 三、错误处理策略

### 3.1 错误传播公理

```
公理: 错误应在边界处被处理或转换
─────────────────────────────
推论:
1. 内部错误 → 包装添加上下文
2. 边界错误 → 转换为领域错误
3. 外部错误 → 归类为标准类型

错误处理层次:
Repository层: 数据库错误 → 领域错误
Service层:    业务规则错误 → 应用错误
Handler层:    应用错误 → HTTP状态码

代码示例:
// 领域错误定义
var ErrUserNotFound = errors.New("user not found")

// Repository层转换
func (r *UserRepo) Get(ctx context.Context, id string) (*User, error) {
    user, err := r.db.QueryContext(ctx, ...)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("%w: id=%s", ErrUserNotFound, id)
    }
    return user, err
}

// Handler层映射
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.service.Get(r.Context(), id)
    if errors.Is(err, ErrUserNotFound) {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    // ...
}
```

### 3.2 错误包装原则

```
决策树:
是否需要包装错误?
├── 调用链需要上下文?
│   └── 是 → fmt.Errorf("...: %w", err)
├── 需要错误分类?
│   └── 是 → 自定义错误类型
└── 只是简单传递?
    └── 否 → return err

Go 1.13+ 错误处理:
- errors.Is: 检查错误链中是否存在
- errors.As: 提取特定错误类型
- %w: 包装保留错误链
- %v: 格式化不保留链
```

---

## 四、性能优化准则

### 4.1 内存优化公理

```
公理: 减少分配 = 降低GC压力 = 更好性能
─────────────────────────────
量化指标:
• 每个请求的alloc次数
• heap_alloc / request_count
• GC CPU占比 (< 10%为健康)

优化策略:
┌────────────────────────────────────────┐
│ 1. 对象池复用 (sync.Pool)              │
│    var pool = sync.Pool{               │
│        New: func() interface{} {       │
│            return new(Buffer)          │
│        },                              │
│    }                                   │
│                                        │
│ 2. 预分配切片容量                       │
│    make([]T, 0, estimatedSize)         │
│                                        │
│ 3. 字符串处理用strings.Builder          │
│                                        │
│ 4. 避免在热路径装箱                     │
│    接口值、反射、any类型                │
│                                        │
│ 5. 值传递 vs 指针传递                   │
│    小对象(<64B): 值传递                 │
│    大对象/需修改: 指针                   │
└────────────────────────────────────────┘
```

### 4.2 并发优化

```
并行度公式:
─────────────────────────────
最优Goroutine数 = CPU核心数 × (1 + 等待时间/计算时间)

当I/O密集型: 可增加goroutine数 (网络等待)
当CPU密集型: 接近CPU核心数

限流机制:
├── 固定worker池
│   workers := make(chan struct{}, maxConcurrency)
├── 动态信号量
│   sem := semaphore.NewWeighted(int64(n))
└── 自适应限流
    基于延迟反馈调整
```

---

## 五、测试策略

### 5.1 测试金字塔

```
                    /\
                   /  \
                  / E2E \          少量: 完整流程验证
                 /─────────\
                / Integration \    中等: 组件协作验证
               /─────────────────\
              /     Unit Tests     \  大量: 单函数验证
             /─────────────────────────\

Go测试原则:
- 表驱动测试覆盖边界
- 子测试组织相关场景
- 并行测试加速执行
- Mock隔离外部依赖

示例:
func TestParse(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
        wantErr  bool
    }{
        {"valid", "123", 123, false},
        {"empty", "", 0, true},
        {"invalid", "abc", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v", err)
            }
            if got != tt.expected {
                t.Errorf("Parse() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### 5.2 测试覆盖率目标

```
覆盖率层级:
├── 单元测试: > 80%
├── 集成测试: 核心流程覆盖
└── E2E测试: 关键路径覆盖

Go 1.26工具:
• go test -coverprofile
• go tool cover -html
• fuzzing原生支持
```

---

## 六、安全实践

### 6.1 输入验证公理

```
公理: 所有外部输入都是不可信的
─────────────────────────────
验证层次:
1. 语法验证: 格式、类型、范围
2. 语义验证: 业务规则
3. 授权验证: 权限检查

安全编码:
┌────────────────────────────────────────┐
│ • SQL使用参数化查询                    │
│ • HTTP头正确转义                       │
│ • 密码用bcrypt/scrypt                  │
│ • 敏感数据不记日志                      │
│ • 依赖定期扫描 (govulncheck)           │
└────────────────────────────────────────┘
```

---

*本章基于形式化原则提炼的工程实践，涵盖代码组织、并发、错误处理、性能、测试和安全六大维度。*
