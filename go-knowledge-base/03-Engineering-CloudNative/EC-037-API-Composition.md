# EC-037: API Composition Pattern (API 组合模式)

> **维度**: Engineering-CloudNative
> **级别**: S (>15KB)
> **标签**: #api-composition #query #aggregator #microservices
> **权威来源**:
>
> - [API Composition Pattern](https://microservices.io/patterns/data/api-composition.html) - Chris Richardson
> - [Backend for Frontend Pattern](https://samnewman.io/patterns/architectural/bff/) - Sam Newman
> - [GraphQL](https://graphql.org/) - Facebook

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在微服务架构中，每个服务有自己的数据库，当客户端需要聚合来自多个服务的数据时，如何避免客户端直接调用多个服务（导致紧耦合和复杂性）？

**形式化描述**:

```
给定: 服务集合 S = {S₁, S₂, ..., Sₙ}，每个服务提供查询接口 Qᵢ
给定: 客户端查询需求 R，需要从多个服务获取数据
约束:
  - 最小化客户端复杂度
  - 优化响应时间
  - 保持服务松耦合
目标: 设计组合函数 C: Q₁ × Q₂ × ... × Qₙ → R
```

**直接访问的问题**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Client Direct Access Anti-Pattern                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│     Client                                                              │
│        │                                                                │
│        ├──────────────► Order Service ──► Order DB                     │
│        │                                                                │
│        ├──────────────► Payment Service ──► Payment DB                 │
│        │           (需要处理多个连接、错误、超时)                         │
│        ├──────────────► Inventory Service ──► Inventory DB             │
│        │                                                                │
│        ├──────────────► Shipping Service ──► Shipping DB               │
│        │                                                                │
│        └──────────────► Customer Service ──► Customer DB               │
│                                                                         │
│  问题:                                                                   │
│  • 客户端紧耦合多个服务                                                  │
│  • 网络开销大（多次往返）                                                │
│  • 错误处理复杂                                                          │
│  • 客户端需要知道服务拓扑                                                │
│  • 不同客户端重复实现聚合逻辑                                             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.2 解决方案形式化

**定义 1.1 (API 组合器)**
API 组合器是一个中间层组件，负责：

1. 接收客户端聚合查询请求
2. 并行或串行调用多个服务
3. 合并响应结果
4. 返回统一的聚合视图

**形式化表示**:

```
组合器 C:
  C(request) = merge(f₁(Q₁(request)), f₂(Q₂(request)), ..., fₙ(Qₙ(request)))

其中:
  - fᵢ: 服务 Sᵢ 的查询适配函数
  - merge: 结果合并函数
  - 可并行执行: ∀i≠j: Qᵢ ∥ Qⱼ
```

**执行模式**:

```
并行模式: 所有查询同时执行
  ┌────────────────────────────────────────┐
  │           Parallel Execution           │
  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐     │
  │  │ Q₁  │ │ Q₂  │ │ Q₃  │ │ Q₄  │     │
  │  │  │  │ │  │  │ │  │  │ │  │  │     │
  │  └──┼──┘ └──┼──┘ └──┼──┘ └──┼──┘     │
  │     └───────┴───────┴───────┘        │
  │                  │                    │
  │                  ▼                    │
  │              merge()                  │
  └────────────────────────────────────────┘

串行模式: 存在依赖关系的查询顺序执行
  ┌────────────────────────────────────────┐
  │          Sequential Execution          │
  │  ┌─────┐    ┌─────┐    ┌─────┐       │
  │  │ Q₁  │───►│ Q₂  │───►│ Q₃  │       │
  │  └─────┘    └─────┘    └─────┘       │
  │     (Q₂ 依赖 Q₁ 的结果)                │
  └────────────────────────────────────────┘
```

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    API Composition Architecture                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                          Client                                  │   │
│  │  (Web App / Mobile App / Third Party)                           │   │
│  └───────────────────────────┬─────────────────────────────────────┘   │
│                              │ Single Request                          │
│                              ▼                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    API Composer / BFF                            │   │
│  │                                                                  │   │
│  │  ┌───────────────────────────────────────────────────────────┐  │   │
│  │  │              Query Planner                                │  │   │
│  │  │  • 解析查询需求                                           │  │   │
│  │  │  • 识别依赖关系                                           │  │   │
│  │  │  • 优化执行计划                                           │  │   │
│  │  └───────────────────────┬───────────────────────────────────┘  │   │
│  │                          │                                       │   │
│  │  ┌───────────────────────┴───────────────────────────────────┐  │   │
│  │  │              Request Executor (Parallel/Sequential)       │  │   │
│  │  │                                                           │  │   │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐      │  │   │
│  │  │  │ Call    │  │ Call    │  │ Call    │  │ Call    │      │  │   │
│  │  │  │ Svc 1   │  │ Svc 2   │  │ Svc 3   │  │ Svc N   │      │  │   │
│  │  │  │ (async) │  │ (async) │  │ (async) │  │ (async) │      │  │   │
│  │  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘      │  │   │
│  │  │       └────────────┴────────────┴────────────┘            │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  │                          │                                       │   │
│  │  ┌───────────────────────┴───────────────────────────────────┐  │   │
│  │  │              Response Aggregator                            │  │   │
│  │  │  • 合并响应                                               │  │   │
│  │  │  • 数据转换                                               │  │   │
│  │  │  • 错误处理                                               │  │   │
│  │  └───────────────────────────────────────────────────────────┘  │   │
│  │                                                                  │   │
│  └───────────────────────────┬─────────────────────────────────────┘   │
│                              │                                          │
│          ┌───────────────────┼───────────────────┐                      │
│          ▼                   ▼                   ▼                      │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐               │
│  │ Order Service │  │Payment Service│  │Inventory Svc  │               │
│  │               │  │               │  │               │               │
│  │ • Order DB    │  │ • Payment DB  │  │ • InventoryDB │               │
│  └───────────────┘  └───────────────┘  └───────────────┘               │
│                                                                         │
│  优化策略:                                                               │
│  • 缓存: 对不常变化的数据进行缓存                                          │
│  • 批量: 合并对同一服务的多个请求                                           │
│  • 短路: 某个服务失败时使用默认值                                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 核心组合器实现

```go
// apicomposition/core.go
package apicomposition

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// ServiceClient 服务客户端接口
type ServiceClient interface {
    // Query 查询服务
    Query(ctx context.Context, query ServiceQuery) (*ServiceResponse, error)

    // Name 返回服务名称
    Name() string
}

// ServiceQuery 服务查询
type ServiceQuery struct {
    Endpoint    string
    Method      string
    Parameters  map[string]interface{}
    Headers     map[string]string
}

// ServiceResponse 服务响应
type ServiceResponse struct {
    Data       interface{}
    StatusCode int
    Error      error
    Latency    time.Duration
}

// CompositionRequest 组合请求
type CompositionRequest struct {
    ID          string
    Queries     []QuerySpec
    Timeout     time.Duration
    Strategy    ExecutionStrategy
    Fallback    FallbackStrategy
}

// QuerySpec 查询规格
type QuerySpec struct {
    ServiceName string
    Query       ServiceQuery
    Required    bool                    // 是否必需
    DependsOn   []string                // 依赖的其他查询
    Transformer ResultTransformer       // 结果转换器
}

// ResultTransformer 结果转换器
type ResultTransformer func(data interface{}) (interface{}, error)

// ExecutionStrategy 执行策略
type ExecutionStrategy int

const (
    StrategyParallel ExecutionStrategy = iota   // 并行执行
    StrategySequential                          // 串行执行
    StrategyAdaptive                            // 自适应（基于依赖）
)

// FallbackStrategy 降级策略
type FallbackStrategy int

const (
    FallbackFailAll FallbackStrategy = iota     // 全部失败
    FallbackPartial                             // 部分成功
    FallbackCache                               // 使用缓存
    FallbackDefault                             // 使用默认值
)

// CompositionResult 组合结果
type CompositionResult struct {
    ID          string
    Data        map[string]interface{}
    Errors      map[string]error
    Partial     bool
    Latency     time.Duration
    Timestamp   time.Time
}

// Composer API 组合器
type Composer struct {
    clients      map[string]ServiceClient
    cache        Cache
    timeout      time.Duration
    maxParallel  int
    logger       Logger
}

// Cache 缓存接口
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, bool)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration)
}

// Logger 日志接口
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Debug(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}

// NewComposer 创建组合器
func NewComposer(timeout time.Duration, maxParallel int, logger Logger) *Composer {
    return &Composer{
        clients:     make(map[string]ServiceClient),
        timeout:     timeout,
        maxParallel: maxParallel,
        logger:      logger,
    }
}

// RegisterClient 注册服务客户端
func (c *Composer) RegisterClient(client ServiceClient) {
    c.clients[client.Name()] = client
}

// SetCache 设置缓存
func (c *Composer) SetCache(cache Cache) {
    c.cache = cache
}

// Compose 执行组合查询
func (c *Composer) Compose(ctx context.Context, request *CompositionRequest) (*CompositionResult, error) {
    start := time.Now()

    c.logger.Info("starting composition",
        Field{"request_id", request.ID},
        Field{"query_count", len(request.Queries)})

    result := &CompositionResult{
        ID:        request.ID,
        Data:      make(map[string]interface{}),
        Errors:    make(map[string]error),
        Timestamp: time.Now(),
    }

    // 构建依赖图
    dag := c.buildDAG(request.Queries)

    // 执行查询
    switch request.Strategy {
    case StrategyParallel:
        c.executeParallel(ctx, request, dag, result)
    case StrategySequential:
        c.executeSequential(ctx, request, dag, result)
    case StrategyAdaptive:
        c.executeAdaptive(ctx, request, dag, result)
    }

    result.Latency = time.Since(start)

    // 检查是否需要降级
    if len(result.Errors) > 0 {
        result.Partial = true
        if request.Fallback == FallbackFailAll && hasRequiredError(request, result) {
            return nil, fmt.Errorf("composition failed: %v", result.Errors)
        }
    }

    c.logger.Info("composition completed",
        Field{"request_id", request.ID},
        Field{"latency_ms", result.Latency.Milliseconds()})

    return result, nil
}

// buildDAG 构建依赖图
func (c *Composer) buildDAG(queries []QuerySpec) map[string][]string {
    dag := make(map[string][]string)

    for _, q := range queries {
        dag[q.ServiceName] = q.DependsOn
    }

    return dag
}

// executeParallel 并行执行
func (c *Composer) executeParallel(ctx context.Context, request *CompositionRequest, dag map[string][]string, result *CompositionResult) {
    ctx, cancel := context.WithTimeout(ctx, request.Timeout)
    defer cancel()

    var wg sync.WaitGroup
    sem := make(chan struct{}, c.maxParallel)

    for _, query := range request.Queries {
        // 检查依赖是否满足
        if !c.dependenciesMet(query, result) {
            result.Errors[query.ServiceName] = fmt.Errorf("dependencies not met")
            continue
        }

        wg.Add(1)
        sem <- struct{}{} // 获取信号量

        go func(q QuerySpec) {
            defer wg.Done()
            defer func() { <-sem }() // 释放信号量

            data, err := c.executeQuery(ctx, q)
            if err != nil {
                result.Errors[q.ServiceName] = err
                return
            }

            result.Data[q.ServiceName] = data
        }(query)
    }

    wg.Wait()
}

// executeSequential 串行执行
func (c *Composer) executeSequential(ctx context.Context, request *CompositionRequest, dag map[string][]string, result *CompositionResult) {
    ctx, cancel := context.WithTimeout(ctx, request.Timeout)
    defer cancel()

    // 拓扑排序（简化版）
    executed := make(map[string]bool)

    for len(executed) < len(request.Queries) {
        for _, query := range request.Queries {
            if executed[query.ServiceName] {
                continue
            }

            // 检查依赖
            depsMet := true
            for _, dep := range query.DependsOn {
                if !executed[dep] {
                    depsMet = false
                    break
                }
            }

            if !depsMet {
                continue
            }

            data, err := c.executeQuery(ctx, query)
            if err != nil {
                result.Errors[query.ServiceName] = err
                if query.Required {
                    return // 必需查询失败，终止
                }
            } else {
                result.Data[query.ServiceName] = data
            }

            executed[query.ServiceName] = true
        }
    }
}

// executeAdaptive 自适应执行
func (c *Composer) executeAdaptive(ctx context.Context, request *CompositionRequest, dag map[string][]string, result *CompositionResult) {
    // 根据是否有依赖决定执行方式
    hasDeps := false
    for _, deps := range dag {
        if len(deps) > 0 {
            hasDeps = true
            break
        }
    }

    if hasDeps {
        c.executeSequential(ctx, request, dag, result)
    } else {
        c.executeParallel(ctx, request, dag, result)
    }
}

// executeQuery 执行单个查询
func (c *Composer) executeQuery(ctx context.Context, spec QuerySpec) (interface{}, error) {
    client, exists := c.clients[spec.ServiceName]
    if !exists {
        return nil, fmt.Errorf("unknown service: %s", spec.ServiceName)
    }

    // 检查缓存
    if c.cache != nil {
        cacheKey := c.buildCacheKey(spec)
        if cached, found := c.cache.Get(ctx, cacheKey); found {
            return cached, nil
        }
    }

    // 执行查询
    start := time.Now()
    resp, err := client.Query(ctx, spec.Query)
    latency := time.Since(start)

    if err != nil {
        c.logger.Error("query failed",
            Field{"service", spec.ServiceName},
            Field{"error", err})
        return nil, err
    }

    if resp.Error != nil {
        return nil, resp.Error
    }

    c.logger.Debug("query succeeded",
        Field{"service", spec.ServiceName},
        Field{"latency_ms", latency.Milliseconds()})

    // 数据转换
    data := resp.Data
    if spec.Transformer != nil {
        data, err = spec.Transformer(data)
        if err != nil {
            return nil, fmt.Errorf("transform failed: %w", err)
        }
    }

    // 写入缓存
    if c.cache != nil {
        c.cache.Set(ctx, c.buildCacheKey(spec), data, 5*time.Minute)
    }

    return data, nil
}

// dependenciesMet 检查依赖是否满足
func (c *Composer) dependenciesMet(query QuerySpec, result *CompositionResult) bool {
    for _, dep := range query.DependsOn {
        if _, exists := result.Data[dep]; !exists {
            return false
        }
    }
    return true
}

// buildCacheKey 构建缓存键
func (c *Composer) buildCacheKey(spec QuerySpec) string {
    params, _ := json.Marshal(spec.Query.Parameters)
    return fmt.Sprintf("%s:%s:%s", spec.ServiceName, spec.Query.Endpoint, params)
}

// hasRequiredError 检查是否有必需服务错误
func hasRequiredError(request *CompositionRequest, result *CompositionResult) bool {
    for _, q := range request.Queries {
        if q.Required {
            if _, exists := result.Errors[q.ServiceName]; exists {
                return true
            }
        }
    }
    return false
}
```

### 2.2 HTTP 服务客户端实现

```go
// apicomposition/http_client.go
package apicomposition

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// HTTPServiceClient HTTP 服务客户端
type HTTPServiceClient struct {
    name       string
    baseURL    string
    httpClient *http.Client
}

// NewHTTPServiceClient 创建 HTTP 服务客户端
func NewHTTPServiceClient(name, baseURL string, timeout time.Duration) *HTTPServiceClient {
    if timeout == 0 {
        timeout = 30 * time.Second
    }

    return &HTTPServiceClient{
        name:    name,
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: timeout,
        },
    }
}

// Name 返回服务名称
func (c *HTTPServiceClient) Name() string {
    return c.name
}

// Query 执行查询
func (c *HTTPServiceClient) Query(ctx context.Context, query ServiceQuery) (*ServiceResponse, error) {
    start := time.Now()

    url := c.baseURL + query.Endpoint

    var body io.Reader
    if query.Method == "POST" || query.Method == "PUT" {
        jsonData, _ := json.Marshal(query.Parameters)
        body = bytes.NewReader(jsonData)
    }

    req, err := http.NewRequestWithContext(ctx, query.Method, url, body)
    if err != nil {
        return nil, err
    }

    // 添加 headers
    req.Header.Set("Content-Type", "application/json")
    for k, v := range query.Headers {
        req.Header.Set(k, v)
    }

    // 执行请求
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return &ServiceResponse{
            Error:   err,
            Latency: time.Since(start),
        }, nil
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return &ServiceResponse{
            Error:   err,
            Latency: time.Since(start),
        }, nil
    }

    // 解析响应
    var data interface{}
    if err := json.Unmarshal(respBody, &data); err != nil {
        data = string(respBody)
    }

    result := &ServiceResponse{
        Data:       data,
        StatusCode: resp.StatusCode,
        Latency:    time.Since(start),
    }

    if resp.StatusCode >= 400 {
        result.Error = fmt.Errorf("service returned status %d", resp.StatusCode)
    }

    return result, nil
}
```

### 2.3 应用示例

```go
// examples/composition/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go-knowledge-base/apicomposition"
)

type logger struct{}

func (l *logger) Info(msg string, fields ...apicomposition.Field) {
    log.Printf("[INFO] %s %+v", msg, fields)
}

func (l *logger) Error(msg string, fields ...apicomposition.Field) {
    log.Printf("[ERROR] %s %+v", msg, fields)
}

func (l *logger) Debug(msg string, fields ...apicomposition.Field) {
    log.Printf("[DEBUG] %s %+v", msg, fields)
}

func main() {
    logger := &logger{}
    composer := apicomposition.NewComposer(30*time.Second, 10, logger)

    // 注册服务客户端
    composer.RegisterClient(apicomposition.NewHTTPServiceClient(
        "order-service", "http://localhost:8081", 5*time.Second))
    composer.RegisterClient(apicomposition.NewHTTPServiceClient(
        "payment-service", "http://localhost:8082", 5*time.Second))
    composer.RegisterClient(apicomposition.NewHTTPServiceClient(
        "inventory-service", "http://localhost:8083", 5*time.Second))

    // 构建组合请求
    request := &apicomposition.CompositionRequest{
        ID:      "req-001",
        Timeout: 10 * time.Second,
        Strategy: apicomposition.StrategyParallel,
        Fallback: apicomposition.FallbackPartial,
        Queries: []apicomposition.QuerySpec{
            {
                ServiceName: "order-service",
                Query: apicomposition.ServiceQuery{
                    Endpoint: "/orders/123",
                    Method:   "GET",
                },
                Required: true,
            },
            {
                ServiceName: "payment-service",
                Query: apicomposition.ServiceQuery{
                    Endpoint: "/payments/order/123",
                    Method:   "GET",
                },
                Required: false,
            },
            {
                ServiceName: "inventory-service",
                Query: apicomposition.ServiceQuery{
                    Endpoint: "/inventory/reservations/order/123",
                    Method:   "GET",
                },
                Required: false,
            },
        },
    }

    // 执行组合查询
    result, err := composer.Compose(context.Background(), request)
    if err != nil {
        log.Fatal(err)
    }

    // 输出结果
    fmt.Printf("Composition completed in %v\n", result.Latency)
    fmt.Printf("Partial: %v\n", result.Partial)
    fmt.Printf("Data: %+v\n", result.Data)
    fmt.Printf("Errors: %+v\n", result.Errors)
}
```

---

## 3. 测试策略

### 3.1 单元测试

```go
// apicomposition/core_test.go
package apicomposition

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type mockServiceClient struct {
    mock.Mock
    name string
}

func (m *mockServiceClient) Name() string {
    return m.name
}

func (m *mockServiceClient) Query(ctx context.Context, query ServiceQuery) (*ServiceResponse, error) {
    args := m.Called(ctx, query)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*ServiceResponse), args.Error(1)
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...Field)  {}
func (m *mockLogger) Error(msg string, fields ...Field) {}
func (m *mockLogger) Debug(msg string, fields ...Field) {}

func TestComposer_ExecuteParallel(t *testing.T) {
    logger := &mockLogger{}
    composer := NewComposer(10*time.Second, 5, logger)

    // 创建 mock 客户端
    client1 := &mockServiceClient{name: "service1"}
    client2 := &mockServiceClient{name: "service2"}

    composer.RegisterClient(client1)
    composer.RegisterClient(client2)

    // 设置期望
    client1.On("Query", mock.Anything, mock.Anything).Return(&ServiceResponse{
        Data: map[string]string{"result": "data1"},
    }, nil)
    client2.On("Query", mock.Anything, mock.Anything).Return(&ServiceResponse{
        Data: map[string]string{"result": "data2"},
    }, nil)

    request := &CompositionRequest{
        ID:       "test-001",
        Timeout:  5 * time.Second,
        Strategy: StrategyParallel,
        Fallback: FallbackPartial,
        Queries: []QuerySpec{
            {ServiceName: "service1", Required: true},
            {ServiceName: "service2", Required: true},
        },
    }

    result, err := composer.Compose(context.Background(), request)

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Len(t, result.Data, 2)
    assert.False(t, result.Partial)
}

func TestComposer_WithDependencies(t *testing.T) {
    logger := &mockLogger{}
    composer := NewComposer(10*time.Second, 5, logger)

    client1 := &mockServiceClient{name: "service1"}
    client2 := &mockServiceClient{name: "service2"}

    composer.RegisterClient(client1)
    composer.RegisterClient(client2)

    // service2 依赖 service1
    client1.On("Query", mock.Anything, mock.Anything).Return(&ServiceResponse{
        Data: map[string]string{"id": "123"},
    }, nil)
    client2.On("Query", mock.Anything, mock.Anything).Return(&ServiceResponse{
        Data: map[string]string{"detail": "data"},
    }, nil)

    request := &CompositionRequest{
        ID:       "test-002",
        Timeout:  5 * time.Second,
        Strategy: StrategyAdaptive,
        Queries: []QuerySpec{
            {ServiceName: "service1"},
            {ServiceName: "service2", DependsOn: []string{"service1"}},
        },
    }

    result, err := composer.Compose(context.Background(), request)

    assert.NoError(t, err)
    assert.Len(t, result.Data, 2)
}

func TestComposer_RequiredFailure(t *testing.T) {
    logger := &mockLogger{}
    composer := NewComposer(10*time.Second, 5, logger)

    client1 := &mockServiceClient{name: "service1"}
    composer.RegisterClient(client1)

    client1.On("Query", mock.Anything, mock.Anything).Return(&ServiceResponse{
        Error: errors.New("service unavailable"),
    }, nil)

    request := &CompositionRequest{
        ID:       "test-003",
        Timeout:  5 * time.Second,
        Strategy: StrategyParallel,
        Fallback: FallbackFailAll,
        Queries: []QuerySpec{
            {ServiceName: "service1", Required: true},
        },
    }

    result, err := composer.Compose(context.Background(), request)

    assert.Error(t, err)
    assert.Nil(t, result)
}
```

---

## 4. 与其他模式的集成

### 4.1 与 CQRS 的集成

```
┌─────────────────────────────────────────────────────────────────────────┐
│              API Composition with CQRS Integration                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────┐                                                    │
│  │     Client      │                                                    │
│  └────────┬────────┘                                                    │
│           │                                                             │
│           ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    API Composer (Read Side)                      │   │
│  │                                                                  │   │
│  │  ┌─────────────────────────────────────────────────────────────┐  │   │
│  │  │              Query from Read Models                          │  │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │  │   │
│  │  │  │ Order View  │  │Payment View │  │InventoryView│         │  │   │
│  │  │  │  (Redis)    │  │  (ES)       │  │  (Mongo)    │         │  │   │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘         │  │   │
│  │  └─────────────────────────────────────────────────────────────┘  │   │
│  │                                                                  │   │
│  │  ┌─────────────────────────────────────────────────────────────┐  │   │
│  │  │              Fallback to Services                            │  │   │
│  │  │  (if read model stale or missing)                            │  │   │
│  │  └─────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  优势:                                                                   │
│  • 读模型优化查询性能                                                    │
│  • API Composer 可以优先查询物化视图                                     │
│  • 写操作仍通过 Command 到聚合根                                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 与 Backend-for-Frontend (BFF) 的关系

| API Composition | BFF |
|-----------------|-----|
| 通用聚合逻辑 | 特定客户端逻辑 |
| 多个客户端复用 | 一个 BFF 对应一个客户端类型 |
| 关注数据聚合 | 关注用户体验 |
| 可以组合使用：BFF → API Composition → Services |

---

## 5. 决策标准

### 5.1 何时使用 API Composition

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    API Composition Decision Tree                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  客户端需要从多个服务获取数据？ ──否──► 直接调用单个服务                    │
│            │                                                            │
│            ▼                                                            │
│  数据关联复杂（需要 join）？     ──是──► 考虑 CQRS / 物化视图               │
│            │                                                            │
│            ▼                                                            │
│  有多个不同类型的客户端？        ──是──► 考虑 BFF + API Composition         │
│            │                                                            │
│            ▼                                                            │
│  使用 API Composition                                                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    API Composition Checklist                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  设计阶段:                                                               │
│  □ 识别需要聚合的查询场景                                                │
│  □ 分析服务间依赖关系                                                    │
│  □ 确定缓存策略                                                          │
│  □ 设计降级策略                                                          │
│                                                                         │
│  实现阶段:                                                               │
│  □ 实现并行查询执行                                                      │
│  □ 实现超时和熔断机制                                                    │
│  □ 添加监控和指标                                                        │
│  □ 实现缓存层                                                            │
│                                                                         │
│  注意事项:                                                               │
│  ❌ 避免过度聚合（太多服务）                                               │
│  ❌ 避免嵌套组合（Composer 调用 Composer）                                 │
│  ❌ 不要忘记处理部分失败                                                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>17KB, 完整形式化 + Go 实现 + 测试)

**相关文档**:

- [EC-038-Command-Query-Responsibility.md](./EC-038-Command-Query-Responsibility.md)
- [EC-031-Choreography-Pattern.md](./EC-031-Choreography-Pattern.md)
