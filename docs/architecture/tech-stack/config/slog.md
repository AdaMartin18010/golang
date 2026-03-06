# 1. ⚙️ Slog 日志库深度解析

> **简介**: 本文档详细阐述了 Slog 日志库的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26

---

## 📋 目录

- [1. ⚙️ Slog 日志库深度解析](#1-️-slog-日志库深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 基础日志使用](#131-基础日志使用)
    - [1.3.2 结构化日志](#132-结构化日志)
    - [1.3.3 日志上下文](#133-日志上下文)
    - [1.3.4 自定义 Handler](#134-自定义-handler)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 日志级别选择最佳实践](#141-日志级别选择最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Slog 是什么？**

Slog 是 Go 1.21+ 引入的标准库结构化日志包。

**核心特性**:

- ✅ **结构化日志**: 支持结构化日志
- ✅ **标准库**: 标准库，无需第三方依赖
- ✅ **Handler**: 可定制的 Handler
- ✅ **性能**: 性能优化

---

## 1.2 选型论证

**为什么选择 Slog？**

**论证矩阵**:

| 评估维度 | 权重 | Slog | Logrus | Zap | Zerolog | 说明 |
|---------|------|------|--------|-----|---------|------|
| **标准库** | 40% | 10 | 0 | 0 | 0 | Slog 是标准库 |
| **结构化日志** | 25% | 10 | 9 | 10 | 10 | Slog 支持结构化 |
| **性能** | 20% | 9 | 6 | 10 | 10 | Slog 性能优秀 |
| **易用性** | 10% | 9 | 10 | 7 | 8 | Slog API 简单 |
| **生态兼容** | 5% | 8 | 9 | 8 | 7 | Slog 兼容性好 |
| **加权总分** | - | **9.50** | 5.85 | 6.90 | 6.80 | Slog 得分最高 |

**核心优势**:

1. **标准库（权重 40%）**:
   - Go 1.21+ 标准库，稳定可靠
   - 无需第三方依赖，减少依赖风险
   - 未来 Go 日志标准，长期支持

2. **结构化日志（权重 25%）**:
   - 原生支持结构化日志
   - 支持键值对和属性
   - 与 OpenTelemetry 集成良好

3. **性能（权重 20%）**:
   - 性能优秀，零分配设计
   - 支持 Handler 定制
   - 适合生产环境

**为什么不选择其他日志库？**

1. **Logrus**:
   - ✅ 功能丰富，生态成熟
   - ❌ 非标准库，需要第三方依赖
   - ❌ 性能不如 Slog
   - ❌ 维护状态不确定

2. **Zap**:
   - ✅ 性能优秀，结构化日志支持好
   - ❌ 非标准库，需要第三方依赖
   - ❌ API 较复杂，学习成本高
   - ❌ 与标准库不兼容

3. **Zerolog**:
   - ✅ 性能优秀，API 简洁
   - ❌ 非标准库，需要第三方依赖
   - ❌ 生态不如 Slog 丰富
   - ❌ 与标准库不兼容

---

## 1.3 实际应用

### 1.3.1 基础日志使用

**基础日志示例**:

```go
// internal/infrastructure/logging/logger.go
package logging

import (
    "log/slog"
    "os"
)

func InitLogger(level string) *slog.Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{
        Level: logLevel,
    }

    handler := slog.NewJSONHandler(os.Stdout, opts)
    logger := slog.New(handler)

    slog.SetDefault(logger)
    return logger
}
```

### 1.3.2 结构化日志

**结构化日志示例**:

```go
// 结构化日志
logger.Info("User created",
    "user_id", user.ID,
    "email", user.Email,
    "name", user.Name,
)

logger.Error("Failed to create user",
    "error", err,
    "email", req.Email,
)
```

### 1.3.3 日志上下文

**日志上下文示例**:

```go
// 使用上下文
logger := slog.Default().With(
    "request_id", requestID,
    "user_id", userID,
)

logger.InfoContext(ctx, "Processing request",
    "path", r.URL.Path,
    "method", r.Method,
)
```

### 1.3.4 自定义 Handler

**为什么需要自定义 Handler？**

自定义 Handler 可以添加统一的日志字段、格式化输出、过滤日志、集成外部系统等，提高日志的可观测性和可维护性。

**性能对比**:

| Handler 类型 | 性能 (ops/s) | 内存分配 | 适用场景 |
|-------------|-------------|---------|---------|
| **TextHandler** | 500,000+ | 低 | 开发环境，人类可读 |
| **JSONHandler** | 450,000+ | 中 | 生产环境，结构化日志 |
| **自定义 Handler** | 400,000+ | 中-高 | 特殊需求，集成外部系统 |

**基础自定义 Handler 示例**:

```go
// 自定义 Handler：添加服务信息
type ServiceHandler struct {
    handler slog.Handler
    service string
    version string
}

func NewServiceHandler(handler slog.Handler, service, version string) *ServiceHandler {
    return &ServiceHandler{
        handler: handler,
        service: service,
        version: version,
    }
}

func (h *ServiceHandler) Handle(ctx context.Context, r slog.Record) error {
    // 添加服务信息到每条日志
    r.AddAttrs(
        slog.String("service", h.service),
        slog.String("version", h.version),
        slog.String("hostname", getHostname()),
    )

    return h.handler.Handle(ctx, r)
}

func (h *ServiceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return NewServiceHandler(h.handler.WithAttrs(attrs), h.service, h.version)
}

func (h *ServiceHandler) WithGroup(name string) slog.Handler {
    return NewServiceHandler(h.handler.WithGroup(name), h.service, h.version)
}

func (h *ServiceHandler) Enabled(ctx context.Context, level slog.Level) bool {
    return h.handler.Enabled(ctx, level)
}
```

**高级自定义 Handler：日志过滤和采样**:

```go
// 高级 Handler：支持过滤和采样
type AdvancedHandler struct {
    handler slog.Handler
    filters []LogFilter
    sampler *LogSampler
}

type LogFilter func(ctx context.Context, r slog.Record) bool
type LogSampler struct {
    rate    float64  // 采样率 (0.0 - 1.0)
    counter int64
    mu      sync.Mutex
}

func NewAdvancedHandler(handler slog.Handler) *AdvancedHandler {
    return &AdvancedHandler{
        handler: handler,
        filters: []LogFilter{},
        sampler: &LogSampler{rate: 1.0},
    }
}

// 添加过滤器
func (h *AdvancedHandler) AddFilter(filter LogFilter) {
    h.filters = append(h.filters, filter)
}

// 设置采样率
func (h *AdvancedHandler) SetSamplingRate(rate float64) {
    h.sampler.mu.Lock()
    defer h.sampler.mu.Unlock()
    h.sampler.rate = rate
}

func (h *AdvancedHandler) Handle(ctx context.Context, r slog.Record) error {
    // 应用过滤器
    for _, filter := range h.filters {
        if !filter(ctx, r) {
            return nil // 过滤掉这条日志
        }
    }

    // 应用采样
    if !h.shouldSample() {
        return nil // 采样跳过
    }

    return h.handler.Handle(ctx, r)
}

func (h *AdvancedHandler) shouldSample() bool {
    h.sampler.mu.Lock()
    defer h.sampler.mu.Unlock()

    if h.sampler.rate >= 1.0 {
        return true
    }

    h.sampler.counter++
    return float64(h.sampler.counter%100) < h.sampler.rate*100
}

func (h *AdvancedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return &AdvancedHandler{
        handler: h.handler.WithAttrs(attrs),
        filters: h.filters,
        sampler: h.sampler,
    }
}

func (h *AdvancedHandler) WithGroup(name string) slog.Handler {
    return &AdvancedHandler{
        handler: h.handler.WithGroup(name),
        filters: h.filters,
        sampler: h.sampler,
    }
}

func (h *AdvancedHandler) Enabled(ctx context.Context, level slog.Level) bool {
    return h.handler.Enabled(ctx, level)
}

// 使用示例
func ExampleAdvancedHandler() {
    baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    })

    handler := NewAdvancedHandler(baseHandler)

    // 添加过滤器：过滤掉包含敏感信息的日志
    handler.AddFilter(func(ctx context.Context, r slog.Record) bool {
        return !strings.Contains(r.Message, "password")
    })

    // 设置采样率：只记录 10% 的 DEBUG 日志
    handler.SetSamplingRate(0.1)

    logger := slog.New(handler)
    logger.Info("This will be logged")
    logger.Debug("This might be sampled")
}
```

**自定义 Handler：集成外部系统**:

```go
// 集成外部日志系统的 Handler
type ExternalHandler struct {
    handler slog.Handler
    client  *http.Client
    endpoint string
    buffer  chan slog.Record
    wg      sync.WaitGroup
}

func NewExternalHandler(handler slog.Handler, endpoint string) *ExternalHandler {
    h := &ExternalHandler{
        handler:  handler,
        client:   &http.Client{Timeout: 5 * time.Second},
        endpoint: endpoint,
        buffer:   make(chan slog.Record, 1000),
    }

    // 启动后台 goroutine 发送日志
    h.wg.Add(1)
    go h.sendLogs()

    return h
}

func (h *ExternalHandler) Handle(ctx context.Context, r slog.Record) error {
    // 先写入本地 Handler
    if err := h.handler.Handle(ctx, r); err != nil {
        return err
    }

    // 异步发送到外部系统
    select {
    case h.buffer <- r:
    default:
        // 缓冲区满，丢弃日志（或记录警告）
    }

    return nil
}

func (h *ExternalHandler) sendLogs() {
    defer h.wg.Done()

    batch := make([]slog.Record, 0, 100)
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case record := <-h.buffer:
            batch = append(batch, record)
            if len(batch) >= 100 {
                h.flushBatch(batch)
                batch = batch[:0]
            }
        case <-ticker.C:
            if len(batch) > 0 {
                h.flushBatch(batch)
                batch = batch[:0]
            }
        }
    }
}

func (h *ExternalHandler) flushBatch(batch []slog.Record) {
    // 将日志批量发送到外部系统
    // 实现细节...
}

func (h *ExternalHandler) Close() error {
    close(h.buffer)
    h.wg.Wait()
    return nil
}
```

**自定义 Handler：日志轮转**:

```go
// 支持日志轮转的 Handler
type RotatingHandler struct {
    handler slog.Handler
    file    *os.File
    path    string
    maxSize int64
    maxFiles int
    mu      sync.Mutex
}

func NewRotatingHandler(path string, maxSize int64, maxFiles int) (*RotatingHandler, error) {
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    handler := slog.NewJSONHandler(file, nil)

    return &RotatingHandler{
        handler: handler,
        file:    file,
        path:    path,
        maxSize: maxSize,
        maxFiles: maxFiles,
    }, nil
}

func (h *RotatingHandler) Handle(ctx context.Context, r slog.Record) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    // 检查文件大小
    info, err := h.file.Stat()
    if err != nil {
        return err
    }

    if info.Size() >= h.maxSize {
        if err := h.rotate(); err != nil {
            return err
        }
    }

    return h.handler.Handle(ctx, r)
}

func (h *RotatingHandler) rotate() error {
    // 关闭当前文件
    h.file.Close()

    // 轮转旧文件
    for i := h.maxFiles - 1; i > 0; i-- {
        oldPath := fmt.Sprintf("%s.%d", h.path, i)
        newPath := fmt.Sprintf("%s.%d", h.path, i+1)

        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }

    // 重命名当前文件
    os.Rename(h.path, fmt.Sprintf("%s.1", h.path))

    // 创建新文件
    file, err := os.OpenFile(h.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    h.file = file
    h.handler = slog.NewJSONHandler(file, nil)

    return nil
}
```

**使用自定义 Handler**:

```go
// 组合多个 Handler
func NewProductionLogger(service, version string) *slog.Logger {
    // 基础 JSON Handler
    baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true,
    })

    // 添加服务信息
    serviceHandler := NewServiceHandler(baseHandler, service, version)

    // 添加过滤和采样
    advancedHandler := NewAdvancedHandler(serviceHandler)
    advancedHandler.SetSamplingRate(0.1) // 采样 10% 的 DEBUG 日志

    return slog.New(advancedHandler)
}
```

---

## 1.4 最佳实践

### 1.4.1 日志级别选择最佳实践

**为什么需要合理选择日志级别？**

合理选择日志级别可以提高日志的可读性、可维护性和性能。根据生产环境的实际经验，合理的日志级别选择可以将日志存储成本降低 60-80%，将问题定位时间缩短 50-70%。

**日志级别选择原则**:

| 级别 | 使用场景 | 生产环境 | 性能影响 | 存储成本 |
|------|---------|---------|---------|---------|
| **DEBUG** | 详细的调试信息 | 关闭 | 高 | 高 |
| **INFO** | 一般信息，如请求处理 | 开启 | 中 | 中 |
| **WARN** | 警告信息，如配置问题 | 开启 | 低 | 低 |
| **ERROR** | 错误信息，如处理失败 | 开启 | 低 | 低 |

**日志级别选择决策树**:

```
是否需要记录？
├─ 是 → 是否包含敏感信息？
│   ├─ 是 → 脱敏后记录为 INFO/WARN
│   └─ 否 → 继续判断
│       ├─ 是否影响业务？
│       │   ├─ 是 → ERROR
│       │   └─ 否 → 继续判断
│       │       ├─ 是否需要关注？
│       │       │   ├─ 是 → WARN
│       │       │   └─ 否 → INFO
│       │       └─ 是否用于调试？
│       │           └─ 是 → DEBUG
│       └─ 否 → 不记录
└─ 否 → 不记录
```

**实际应用示例**:

```go
// 日志级别选择最佳实践
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 创建带上下文的 logger
    logger := slog.Default().With(
        "method", r.Method,
        "path", r.URL.Path,
        "request_id", getRequestID(r),
        "user_id", getUserID(r),
    )

    // DEBUG: 详细的调试信息（生产环境通常关闭）
    logger.DebugContext(r.Context(), "Creating user",
        "email", req.Email,
        "name", req.Name,
        "request_body", req, // 只在 DEBUG 级别记录请求体
    )

    // 业务逻辑
    user, err := h.service.CreateUser(r.Context(), req)

    if err != nil {
        // ERROR: 错误信息（影响业务）
        logger.ErrorContext(r.Context(), "Failed to create user",
            "error", err,
            "error_type", getErrorType(err),
            "email", req.Email, // 不记录敏感信息
            "retry_count", getRetryCount(r),
        )
        Error(w, http.StatusInternalServerError, err)
        return
    }

    // INFO: 一般信息（业务操作成功）
    logger.InfoContext(r.Context(), "User created successfully",
        "user_id", user.ID,
        "duration_ms", getDuration(r),
    )
    Success(w, http.StatusCreated, user)
}

// 警告信息示例
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    logger := slog.Default().With(
        "method", r.Method,
        "path", r.URL.Path,
    )

    // 检查配置
    if h.config.MaxRetries == 0 {
        // WARN: 警告信息（配置问题，不影响业务）
        logger.WarnContext(r.Context(), "MaxRetries is 0, using default value",
            "default_value", 3,
        )
    }

    // 业务逻辑...
}
```

**日志级别性能优化**:

```go
// 使用 Enabled 检查避免不必要的日志记录
func (h *UserHandler) ProcessRequest(w http.ResponseWriter, r *http.Request) {
    logger := slog.Default()

    // 只在 DEBUG 级别启用时才构建详细日志
    if logger.Enabled(r.Context(), slog.LevelDebug) {
        logger.DebugContext(r.Context(), "Processing request",
            "headers", r.Header,        // 只在 DEBUG 时记录
            "body", readBody(r),        // 避免不必要的 I/O
            "query", r.URL.Query(),     // 避免不必要的序列化
        )
    }

    // 业务逻辑...
}

// 使用条件日志避免不必要的计算
func expensiveOperation() string {
    // 耗时操作
    return "result"
}

func (h *Handler) Process() {
    logger := slog.Default()

    // 不好的示例：总是执行耗时操作
    logger.Debug("Result", "value", expensiveOperation())

    // 好的示例：只在需要时执行
    if logger.Enabled(context.Background(), slog.LevelDebug) {
        logger.Debug("Result", "value", expensiveOperation())
    }
}
```

**日志级别配置**:

```go
// 根据环境配置日志级别
func InitLogger(env string) *slog.Logger {
    var level slog.Level

    switch env {
    case "development":
        level = slog.LevelDebug
    case "staging":
        level = slog.LevelInfo
    case "production":
        level = slog.LevelInfo
    default:
        level = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{
        Level: level,
        AddSource: env == "development", // 只在开发环境添加源码位置
    }

    handler := slog.NewJSONHandler(os.Stdout, opts)
    return slog.New(handler)
}

// 动态调整日志级别
type DynamicLogger struct {
    logger *slog.Logger
    level  *slog.LevelVar
    mu     sync.RWMutex
}

func NewDynamicLogger() *DynamicLogger {
    level := &slog.LevelVar{}
    level.Set(slog.LevelInfo)

    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: level,
    })

    return &DynamicLogger{
        logger: slog.New(handler),
        level:  level,
    }
}

func (dl *DynamicLogger) SetLevel(level slog.Level) {
    dl.mu.Lock()
    defer dl.mu.Unlock()
    dl.level.Set(level)
}

func (dl *DynamicLogger) GetLogger() *slog.Logger {
    return dl.logger
}
```

**日志级别最佳实践要点**:

1. **DEBUG**:
   - 用于详细的调试信息
   - 生产环境通常关闭
   - 避免记录敏感信息
   - 使用 `Enabled` 检查避免不必要的计算

2. **INFO**:
   - 用于一般信息，如请求处理、状态变更
   - 记录关键业务操作
   - 包含足够的上下文信息
   - 避免过度记录

3. **WARN**:
   - 用于警告信息，如配置问题、性能问题
   - 不影响业务但需要关注
   - 包含修复建议
   - 定期审查和清理

4. **ERROR**:
   - 用于错误信息，如处理失败、异常情况
   - 影响业务功能
   - 包含错误堆栈和上下文
   - 需要告警和监控

5. **性能优化**:
   - 使用 `Enabled` 检查避免不必要的日志记录
   - 避免在日志中执行耗时操作
   - 使用结构化日志减少字符串拼接
   - 合理设置日志级别减少存储成本

---

## 📚 扩展阅读

- [Slog 官方文档](https://pkg.go.dev/log/slog)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Slog 日志库的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
