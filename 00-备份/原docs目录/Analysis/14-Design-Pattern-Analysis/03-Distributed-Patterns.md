# 分布式模式分析

## 概述

分布式模式是构建大规模分布式系统的核心设计模式。本文档基于Golang技术栈，深入分析各种分布式模式的设计、实现和性能特征。

## 1. 服务发现模式 (Service Discovery)

### 1.1 定义

允许服务动态注册和发现，实现服务的自动发现和负载均衡。

### 1.2 形式化定义

$$\text{ServiceDiscovery} = (S, R, L, H)$$

其中：

- $S$ 是服务集合
- $R$ 是注册机制
- $L$ 是负载均衡
- $H$ 是健康检查

### 1.3 Golang实现

```go
package servicediscovery

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Service 服务信息
type Service struct {
    ID       string
    Name     string
    Address  string
    Port     int
    Metadata map[string]string
    Status   ServiceStatus
    LastSeen time.Time
}

// ServiceStatus 服务状态
type ServiceStatus string

const (
    StatusHealthy   ServiceStatus = "healthy"
    StatusUnhealthy ServiceStatus = "unhealthy"
    StatusUnknown   ServiceStatus = "unknown"
)

// ServiceRegistry 服务注册表
type ServiceRegistry struct {
    services map[string]*Service
    mu       sync.RWMutex
    watchers map[string][]chan *ServiceEvent
    ctx      context.Context
    cancel   context.CancelFunc
}

// ServiceEvent 服务事件
type ServiceEvent struct {
    Type    string
    Service *Service
}

// NewServiceRegistry 创建服务注册表
func NewServiceRegistry() *ServiceRegistry {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &ServiceRegistry{
        services: make(map[string]*Service),
        watchers: make(map[string][]chan *ServiceEvent),
        ctx:      ctx,
        cancel:   cancel,
    }
}

// Register 注册服务
func (r *ServiceRegistry) Register(service *Service) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.services[service.ID]; exists {
        return fmt.Errorf("service already exists: %s", service.ID)
    }
    
    service.LastSeen = time.Now()
    service.Status = StatusHealthy
    r.services[service.ID] = service
    
    // 通知观察者
    r.notifyWatchers("register", service)
    
    return nil
}

// Deregister 注销服务
func (r *ServiceRegistry) Deregister(serviceID string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    service, exists := r.services[serviceID]
    if !exists {
        return fmt.Errorf("service not found: %s", serviceID)
    }
    
    delete(r.services, serviceID)
    
    // 通知观察者
    r.notifyWatchers("deregister", service)
    
    return nil
}

// GetService 获取服务
func (r *ServiceRegistry) GetService(serviceID string) (*Service, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    service, exists := r.services[serviceID]
    return service, exists
}

// GetServicesByName 根据名称获取服务
func (r *ServiceRegistry) GetServicesByName(name string) []*Service {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var services []*Service
    for _, service := range r.services {
        if service.Name == name {
            services = append(services, service)
        }
    }
    return services
}

// UpdateHealth 更新健康状态
func (r *ServiceRegistry) UpdateHealth(serviceID string, status ServiceStatus) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    service, exists := r.services[serviceID]
    if !exists {
        return fmt.Errorf("service not found: %s", serviceID)
    }
    
    service.Status = status
    service.LastSeen = time.Now()
    
    // 通知观察者
    r.notifyWatchers("health_update", service)
    
    return nil
}

// Watch 监听服务变化
func (r *ServiceRegistry) Watch(serviceName string) chan *ServiceEvent {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    eventChan := make(chan *ServiceEvent, 10)
    r.watchers[serviceName] = append(r.watchers[serviceName], eventChan)
    
    return eventChan
}

// notifyWatchers 通知观察者
func (r *ServiceRegistry) notifyWatchers(eventType string, service *Service) {
    watchers := r.watchers[service.Name]
    for _, watcher := range watchers {
        select {
        case watcher <- &ServiceEvent{Type: eventType, Service: service}:
        default:
            // 通道已满，跳过
        }
    }
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
    strategy LoadBalancingStrategy
    registry *ServiceRegistry
}

// LoadBalancingStrategy 负载均衡策略
type LoadBalancingStrategy interface {
    Select(services []*Service) *Service
}

// RoundRobinStrategy 轮询策略
type RoundRobinStrategy struct {
    current int
    mu      sync.Mutex
}

func (s *RoundRobinStrategy) Select(services []*Service) *Service {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if len(services) == 0 {
        return nil
    }
    
    service := services[s.current%len(services)]
    s.current++
    return service
}

// RandomStrategy 随机策略
type RandomStrategy struct{}

func (s *RandomStrategy) Select(services []*Service) *Service {
    if len(services) == 0 {
        return nil
    }
    
    // 简化的随机选择
    return services[0]
}

// LeastConnectionsStrategy 最少连接策略
type LeastConnectionsStrategy struct{}

func (s *LeastConnectionsStrategy) Select(services []*Service) *Service {
    if len(services) == 0 {
        return nil
    }
    
    // 选择连接数最少的服务
    var selected *Service
    minConnections := int(^uint(0) >> 1)
    
    for _, service := range services {
        if service.Status == StatusHealthy {
            connections := 0 // 这里应该从服务获取实际连接数
            if connections < minConnections {
                minConnections = connections
                selected = service
            }
        }
    }
    
    return selected
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer(registry *ServiceRegistry, strategy LoadBalancingStrategy) *LoadBalancer {
    return &LoadBalancer{
        strategy: strategy,
        registry: registry,
    }
}

// SelectService 选择服务
func (lb *LoadBalancer) SelectService(serviceName string) (*Service, error) {
    services := lb.registry.GetServicesByName(serviceName)
    
    // 过滤健康服务
    healthyServices := make([]*Service, 0)
    for _, service := range services {
        if service.Status == StatusHealthy {
            healthyServices = append(healthyServices, service)
        }
    }
    
    if len(healthyServices) == 0 {
        return nil, fmt.Errorf("no healthy services found for: %s", serviceName)
    }
    
    return lb.strategy.Select(healthyServices), nil
}

// ServiceDiscoveryClient 服务发现客户端
type ServiceDiscoveryClient struct {
    registry    *ServiceRegistry
    loadBalancer *LoadBalancer
    cache       map[string]*Service
    cacheTTL    time.Duration
    mu          sync.RWMutex
}

// NewServiceDiscoveryClient 创建服务发现客户端
func NewServiceDiscoveryClient(registry *ServiceRegistry, strategy LoadBalancingStrategy) *ServiceDiscoveryClient {
    return &ServiceDiscoveryClient{
        registry:     registry,
        loadBalancer: NewLoadBalancer(registry, strategy),
        cache:        make(map[string]*Service),
        cacheTTL:     30 * time.Second,
    }
}

// Discover 发现服务
func (c *ServiceDiscoveryClient) Discover(serviceName string) (*Service, error) {
    // 检查缓存
    c.mu.RLock()
    if cached, exists := c.cache[serviceName]; exists {
        c.mu.RUnlock()
        return cached, nil
    }
    c.mu.RUnlock()
    
    // 从负载均衡器选择服务
    service, err := c.loadBalancer.SelectService(serviceName)
    if err != nil {
        return nil, err
    }
    
    // 更新缓存
    c.mu.Lock()
    c.cache[serviceName] = service
    c.mu.Unlock()
    
    return service, nil
}

// ClearCache 清除缓存
func (c *ServiceDiscoveryClient) ClearCache() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache = make(map[string]*Service)
}

```

## 2. 熔断器模式 (Circuit Breaker)

### 2.1 定义

防止系统级联故障，通过监控失败率自动断开故障服务。

### 2.2 状态机

- **关闭状态**: 正常处理请求
- **开启状态**: 快速失败，不处理请求
- **半开状态**: 允许部分请求尝试

### 2.3 Golang实现

```go
package circuitbreaker

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// State 熔断器状态
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

// CircuitBreaker 熔断器
type CircuitBreaker struct {
    name           string
    state          State
    failureCount   int64
    successCount   int64
    lastFailureTime time.Time
    threshold      int64
    timeout        time.Duration
    mu             sync.RWMutex
    stats          *CircuitBreakerStats
}

// CircuitBreakerStats 熔断器统计
type CircuitBreakerStats struct {
    TotalRequests   int64
    SuccessfulRequests int64
    FailedRequests  int64
    StateChanges    int64
    mu              sync.RWMutex
}

func (s *CircuitBreakerStats) IncrementTotal() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.TotalRequests++
}

func (s *CircuitBreakerStats) IncrementSuccess() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.SuccessfulRequests++
}

func (s *CircuitBreakerStats) IncrementFailure() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.FailedRequests++
}

func (s *CircuitBreakerStats) IncrementStateChange() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.StateChanges++
}

func (s *CircuitBreakerStats) GetStats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    failureRate := float64(0)
    if s.TotalRequests > 0 {
        failureRate = float64(s.FailedRequests) / float64(s.TotalRequests)
    }
    
    return map[string]interface{}{
        "total_requests":     s.TotalRequests,
        "successful_requests": s.SuccessfulRequests,
        "failed_requests":    s.FailedRequests,
        "failure_rate":       failureRate,
        "state_changes":      s.StateChanges,
    }
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(name string, threshold int64, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        name:      name,
        state:     StateClosed,
        threshold: threshold,
        timeout:   timeout,
        stats:     &CircuitBreakerStats{},
    }
}

// Execute 执行操作
func (cb *CircuitBreaker) Execute(operation func() error) error {
    if !cb.canExecute() {
        return fmt.Errorf("circuit breaker is open")
    }
    
    cb.stats.IncrementTotal()
    
    err := operation()
    
    if err != nil {
        cb.onFailure()
        cb.stats.IncrementFailure()
        return err
    }
    
    cb.onSuccess()
    cb.stats.IncrementSuccess()
    return nil
}

// ExecuteWithContext 带上下文的执行
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, operation func() error) error {
    if !cb.canExecute() {
        return fmt.Errorf("circuit breaker is open")
    }
    
    cb.stats.IncrementTotal()
    
    // 创建带超时的上下文
    timeoutCtx, cancel := context.WithTimeout(ctx, cb.timeout)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        done <- operation()
    }()
    
    select {
    case err := <-done:
        if err != nil {
            cb.onFailure()
            cb.stats.IncrementFailure()
            return err
        }
        cb.onSuccess()
        cb.stats.IncrementSuccess()
        return nil
    case <-timeoutCtx.Done():
        cb.onFailure()
        cb.stats.IncrementFailure()
        return fmt.Errorf("operation timeout")
    }
}

// canExecute 检查是否可以执行
func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        // 检查是否超时
        if time.Since(cb.lastFailureTime) >= cb.timeout {
            cb.transitionToHalfOpen()
            return true
        }
        return false
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

// onSuccess 成功回调
func (cb *CircuitBreaker) onSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    cb.successCount++
    
    if cb.state == StateHalfOpen {
        if cb.successCount >= cb.threshold {
            cb.transitionToClosed()
        }
    }
}

// onFailure 失败回调
func (cb *CircuitBreaker) onFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    cb.failureCount++
    cb.lastFailureTime = time.Now()
    
    switch cb.state {
    case StateClosed:
        if cb.failureCount >= cb.threshold {
            cb.transitionToOpen()
        }
    case StateHalfOpen:
        cb.transitionToOpen()
    }
}

// transitionToOpen 转换到开启状态
func (cb *CircuitBreaker) transitionToOpen() {
    if cb.state != StateOpen {
        cb.state = StateOpen
        cb.stats.IncrementStateChange()
    }
}

// transitionToHalfOpen 转换到半开状态
func (cb *CircuitBreaker) transitionToHalfOpen() {
    if cb.state != StateHalfOpen {
        cb.state = StateHalfOpen
        cb.failureCount = 0
        cb.successCount = 0
        cb.stats.IncrementStateChange()
    }
}

// transitionToClosed 转换到关闭状态
func (cb *CircuitBreaker) transitionToClosed() {
    if cb.state != StateClosed {
        cb.state = StateClosed
        cb.failureCount = 0
        cb.successCount = 0
        cb.stats.IncrementStateChange()
    }
}

// GetState 获取状态
func (cb *CircuitBreaker) GetState() State {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    return cb.state
}

// GetStats 获取统计信息
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
    stats := cb.stats.GetStats()
    stats["state"] = cb.GetState()
    stats["failure_count"] = cb.failureCount
    stats["success_count"] = cb.successCount
    return stats
}

// CircuitBreakerManager 熔断器管理器
type CircuitBreakerManager struct {
    breakers map[string]*CircuitBreaker
    mu       sync.RWMutex
}

// NewCircuitBreakerManager 创建熔断器管理器
func NewCircuitBreakerManager() *CircuitBreakerManager {
    return &CircuitBreakerManager{
        breakers: make(map[string]*CircuitBreaker),
    }
}

// GetOrCreate 获取或创建熔断器
func (m *CircuitBreakerManager) GetOrCreate(name string, threshold int64, timeout time.Duration) *CircuitBreaker {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if breaker, exists := m.breakers[name]; exists {
        return breaker
    }
    
    breaker := NewCircuitBreaker(name, threshold, timeout)
    m.breakers[name] = breaker
    return breaker
}

// Get 获取熔断器
func (m *CircuitBreakerManager) Get(name string) (*CircuitBreaker, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    breaker, exists := m.breakers[name]
    return breaker, exists
}

// GetAllStats 获取所有统计信息
func (m *CircuitBreakerManager) GetAllStats() map[string]interface{} {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    stats := make(map[string]interface{})
    for name, breaker := range m.breakers {
        stats[name] = breaker.GetStats()
    }
    return stats
}

```

## 3. API网关模式 (API Gateway)

### 3.1 定义

为客户端提供统一的API入口，处理路由、认证、限流等功能。

### 3.2 Golang实现

```go
package apigateway

import (
    "context"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// Route 路由定义
type Route struct {
    Path        string
    Method      string
    ServiceName string
    Middleware  []Middleware
}

// Middleware 中间件接口
type Middleware interface {
    Process(ctx context.Context, req *http.Request) (*http.Request, error)
}

// AuthenticationMiddleware 认证中间件
type AuthenticationMiddleware struct {
    tokenValidator func(string) bool
}

func (m *AuthenticationMiddleware) Process(ctx context.Context, req *http.Request) (*http.Request, error) {
    token := req.Header.Get("Authorization")
    if token == "" {
        return nil, fmt.Errorf("missing authorization header")
    }
    
    if !m.tokenValidator(token) {
        return nil, fmt.Errorf("invalid token")
    }
    
    return req, nil
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
    limiter *RateLimiter
}

func (m *RateLimitMiddleware) Process(ctx context.Context, req *http.Request) (*http.Request, error) {
    clientIP := req.RemoteAddr
    if !m.limiter.Allow(clientIP) {
        return nil, fmt.Errorf("rate limit exceeded")
    }
    return req, nil
}

// RateLimiter 限流器
type RateLimiter struct {
    limits map[string]*TokenBucket
    mu     sync.RWMutex
}

// TokenBucket 令牌桶
type TokenBucket struct {
    tokens    int
    capacity  int
    rate      int
    lastRefill time.Time
    mu        sync.Mutex
}

func NewTokenBucket(capacity, rate int) *TokenBucket {
    return &TokenBucket{
        tokens:     capacity,
        capacity:   capacity,
        rate:       rate,
        lastRefill: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    
    // 补充令牌
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill)
    tokensToAdd := int(elapsed.Seconds()) * tb.rate
    
    if tokensToAdd > 0 {
        tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
        tb.lastRefill = now
    }
    
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    
    return false
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// NewRateLimiter 创建限流器
func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        limits: make(map[string]*TokenBucket),
    }
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(clientID string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    bucket, exists := rl.limits[clientID]
    if !exists {
        bucket = NewTokenBucket(100, 10) // 100个令牌，每秒10个
        rl.limits[clientID] = bucket
    }
    
    return bucket.Allow()
}

// APIGateway API网关
type APIGateway struct {
    routes   []*Route
    registry *ServiceRegistry
    breaker  *CircuitBreakerManager
    limiter  *RateLimiter
    mu       sync.RWMutex
}

// NewAPIGateway 创建API网关
func NewAPIGateway(registry *ServiceRegistry) *APIGateway {
    return &APIGateway{
        registry: registry,
        breaker:  NewCircuitBreakerManager(),
        limiter:  NewRateLimiter(),
    }
}

// AddRoute 添加路由
func (g *APIGateway) AddRoute(route *Route) {
    g.mu.Lock()
    defer g.mu.Unlock()
    g.routes = append(g.routes, route)
}

// HandleRequest 处理请求
func (g *APIGateway) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // 查找路由
    route := g.findRoute(r.URL.Path, r.Method)
    if route == nil {
        http.NotFound(w, r)
        return
    }
    
    // 执行中间件
    ctx := context.Background()
    for _, middleware := range route.Middleware {
        var err error
        r, err = middleware.Process(ctx, r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }
    
    // 获取服务
    service, err := g.registry.GetService(route.ServiceName)
    if err != nil {
        http.Error(w, "service not found", http.StatusServiceUnavailable)
        return
    }
    
    // 执行请求
    err = g.executeRequest(route.ServiceName, service, r, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// findRoute 查找路由
func (g *APIGateway) findRoute(path, method string) *Route {
    g.mu.RLock()
    defer g.mu.RUnlock()
    
    for _, route := range g.routes {
        if route.Path == path && route.Method == method {
            return route
        }
    }
    return nil
}

// executeRequest 执行请求
func (g *APIGateway) executeRequest(serviceName string, service *Service, req *http.Request, w http.ResponseWriter) error {
    // 获取熔断器
    breaker := g.breaker.GetOrCreate(serviceName, 5, 30*time.Second)
    
    // 执行请求
    return breaker.Execute(func() error {
        // 这里应该实际转发请求到后端服务
        // 简化实现，直接返回成功
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf("Response from %s", serviceName)))
        return nil
    })
}

// GatewayServer 网关服务器
type GatewayServer struct {
    gateway *APIGateway
    server  *http.Server
}

// NewGatewayServer 创建网关服务器
func NewGatewayServer(gateway *APIGateway, addr string) *GatewayServer {
    mux := http.NewServeMux()
    mux.HandleFunc("/", gateway.HandleRequest)
    
    server := &http.Server{
        Addr:    addr,
        Handler: mux,
    }
    
    return &GatewayServer{
        gateway: gateway,
        server:  server,
    }
}

// Start 启动服务器
func (s *GatewayServer) Start() error {
    return s.server.ListenAndServe()
}

// Shutdown 关闭服务器
func (s *GatewayServer) Shutdown(ctx context.Context) error {
    return s.server.Shutdown(ctx)
}

```

## 4. Saga模式

### 4.1 定义

管理分布式事务，通过补偿操作确保最终一致性。

### 4.2 Golang实现

```go
package saga

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// SagaStep Saga步骤
type SagaStep struct {
    ID          string
    Name        string
    Execute     func(ctx context.Context) error
    Compensate  func(ctx context.Context) error
    Status      StepStatus
    Error       error
}

// StepStatus 步骤状态
type StepStatus int

const (
    StepPending StepStatus = iota
    StepExecuting
    StepCompleted
    StepFailed
    StepCompensated
)

// Saga Saga事务
type Saga struct {
    ID        string
    Steps     []*SagaStep
    Status    SagaStatus
    mu        sync.RWMutex
    ctx       context.Context
    cancel    context.CancelFunc
}

// SagaStatus Saga状态
type SagaStatus int

const (
    SagaRunning SagaStatus = iota
    SagaCompleted
    SagaFailed
    SagaCompensated
)

// NewSaga 创建Saga
func NewSaga(id string) *Saga {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &Saga{
        ID:     id,
        Steps:  make([]*SagaStep, 0),
        Status: SagaRunning,
        ctx:    ctx,
        cancel: cancel,
    }
}

// AddStep 添加步骤
func (s *Saga) AddStep(step *SagaStep) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Steps = append(s.Steps, step)
}

// Execute 执行Saga
func (s *Saga) Execute() error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for i, step := range s.Steps {
        // 执行步骤
        step.Status = StepExecuting
        err := step.Execute(s.ctx)
        
        if err != nil {
            step.Status = StepFailed
            step.Error = err
            s.Status = SagaFailed
            
            // 补偿前面的步骤
            return s.compensate(i - 1)
        }
        
        step.Status = StepCompleted
    }
    
    s.Status = SagaCompleted
    return nil
}

// compensate 补偿
func (s *Saga) compensate(fromIndex int) error {
    for i := fromIndex; i >= 0; i-- {
        step := s.Steps[i]
        if step.Status == StepCompleted {
            step.Status = StepCompensated
            if err := step.Compensate(s.ctx); err != nil {
                return fmt.Errorf("compensation failed for step %s: %v", step.ID, err)
            }
        }
    }
    
    s.Status = SagaCompensated
    return nil
}

// GetStatus 获取状态
func (s *Saga) GetStatus() SagaStatus {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.Status
}

// SagaManager Saga管理器
type SagaManager struct {
    sagas map[string]*Saga
    mu    sync.RWMutex
}

// NewSagaManager 创建Saga管理器
func NewSagaManager() *SagaManager {
    return &SagaManager{
        sagas: make(map[string]*Saga),
    }
}

// CreateSaga 创建Saga
func (m *SagaManager) CreateSaga(id string) *Saga {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    saga := NewSaga(id)
    m.sagas[id] = saga
    return saga
}

// GetSaga 获取Saga
func (m *SagaManager) GetSaga(id string) (*Saga, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    saga, exists := m.sagas[id]
    return saga, exists
}

// ExecuteSaga 执行Saga
func (m *SagaManager) ExecuteSaga(id string) error {
    saga, exists := m.GetSaga(id)
    if !exists {
        return fmt.Errorf("saga not found: %s", id)
    }
    
    return saga.Execute()
}

```

## 5. 性能分析

### 5.1 模式性能对比

| 模式 | 延迟 | 吞吐量 | 可用性 | 一致性 | 复杂度 |
|------|------|--------|--------|--------|--------|
| 服务发现 | 低 | 高 | 高 | 最终 | 中 |
| 熔断器 | 低 | 高 | 高 | 弱 | 低 |
| API网关 | 中 | 高 | 高 | 强 | 高 |
| Saga | 高 | 中 | 中 | 最终 | 高 |

### 5.2 性能指标

**服务发现**:
$$\text{DiscoveryTime} = \text{LookupTime} + \text{HealthCheckTime}$$

**熔断器**:
$$\text{FailureRate} = \frac{\text{FailedRequests}}{\text{TotalRequests}}$$

**API网关**:
$$\text{GatewayLatency} = \text{RoutingTime} + \text{MiddlewareTime} + \text{ServiceTime}$$

### 5.3 容量规划

**服务发现容量**:
$$C_{discovery} = \frac{\text{ServiceCount} \times \text{UpdateRate}}{\text{ProcessingCapacity}}$$

**熔断器容量**:
$$C_{breaker} = \text{RequestRate} \times \text{Timeout} \times \text{Threshold}$$

## 6. 最佳实践

### 6.1 设计原则

1. **容错设计**: 假设故障会发生
2. **降级策略**: 提供降级方案
3. **监控告警**: 全面的监控体系
4. **自动化**: 减少人工干预

### 6.2 实现建议

1. **异步处理**: 使用异步模式提高性能
2. **缓存策略**: 合理使用缓存
3. **重试机制**: 实现智能重试
4. **超时控制**: 设置合理的超时时间

### 6.3 常见陷阱

1. **级联故障**: 避免服务间相互影响
2. **数据不一致**: 处理分布式一致性问题
3. **性能瓶颈**: 避免热点资源
4. **配置错误**: 确保配置正确性

## 7. 应用场景

### 7.1 服务发现

- 微服务架构
- 容器编排
- 云原生应用
- 动态扩缩容

### 7.2 熔断器

- 外部API调用
- 数据库访问
- 第三方服务
- 网络请求

### 7.3 API网关

- 统一入口
- 认证授权
- 限流控制
- 监控日志

### 7.4 Saga模式

- 分布式事务
- 订单处理
- 支付流程
- 库存管理

## 8. 总结

分布式模式为构建大规模分布式系统提供了重要的设计指导。通过合理应用这些模式，可以构建出高可用、高性能的分布式系统。

### 关键优势

- **高可用**: 提高系统可靠性
- **高性能**: 支持高并发处理
- **可扩展**: 支持水平扩展
- **可维护**: 清晰的架构设计

### 成功要素

1. **合理选择**: 根据需求选择合适的模式
2. **性能优化**: 持续的性能优化
3. **监控告警**: 完善的监控体系
4. **测试验证**: 全面的测试覆盖

通过合理应用分布式模式，可以构建出高质量的分布式系统，为业务发展提供强有力的技术支撑。
