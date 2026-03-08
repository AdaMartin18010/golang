# Go 1.26 高级编程模式实战

**版本**: Go 1.26
**性质**: 实战编程模式
**目标**: 掌握Go的高级编程技巧

---

## 目录

- [Go 1.26 高级编程模式实战](#go-126-高级编程模式实战)
  - [目录](#目录)
  - [1. 函数选项模式进阶](#1-函数选项模式进阶)
    - [带验证的选项模式](#带验证的选项模式)
    - [使用示例](#使用示例)
  - [2. 依赖注入模式](#2-依赖注入模式)
    - [构造函数注入](#构造函数注入)
    - [使用 Wire 生成依赖注入代码](#使用-wire-生成依赖注入代码)
  - [3. 仓储模式实现](#3-仓储模式实现)
    - [泛型仓储接口](#泛型仓储接口)
    - [实现示例](#实现示例)
  - [4. 工作单元模式](#4-工作单元模式)
  - [5. 事件驱动模式](#5-事件驱动模式)
  - [6. 断路器模式](#6-断路器模式)
  - [7. 限流器模式](#7-限流器模式)
  - [8. 连接池模式](#8-连接池模式)
  - [总结](#总结)

---

## 1. 函数选项模式进阶

### 带验证的选项模式

```go
package main

import (
    "fmt"
    "net"
    "time"
)

// ServerConfig 服务器配置
type ServerConfig struct {
    Host         string
    Port         int
    Timeout      time.Duration
    MaxConns     int
    TLS          bool
    CertFile     string
    KeyFile      string
}

// Option 配置选项类型
type Option func(*ServerConfig) error

// WithHost 设置主机
func WithHost(host string) Option {
    return func(c *ServerConfig) error {
        if host == "" {
            return fmt.Errorf("host cannot be empty")
        }
        c.Host = host
        return nil
    }
}

// WithPort 设置端口
func WithPort(port int) Option {
    return func(c *ServerConfig) error {
        if port <= 0 || port > 65535 {
            return fmt.Errorf("invalid port: %d", port)
        }
        c.Port = port
        return nil
    }
}

// WithTimeout 设置超时
func WithTimeout(d time.Duration) Option {
    return func(c *ServerConfig) error {
        if d <= 0 {
            return fmt.Errorf("timeout must be positive")
        }
        c.Timeout = d
        return nil
    }
}

// WithMaxConns 设置最大连接数
func WithMaxConns(n int) Option {
    return func(c *ServerConfig) error {
        if n <= 0 {
            return fmt.Errorf("max connections must be positive")
        }
        c.MaxConns = n
        return nil
    }
}

// WithTLS 启用TLS
func WithTLS(certFile, keyFile string) Option {
    return func(c *ServerConfig) error {
        if certFile == "" || keyFile == "" {
            return fmt.Errorf("cert and key files are required for TLS")
        }
        c.TLS = true
        c.CertFile = certFile
        c.KeyFile = keyFile
        return nil
    }
}

// NewServerConfig 创建配置
func NewServerConfig(opts ...Option) (*ServerConfig, error) {
    // 默认值
    cfg := &ServerConfig{
        Host:     "0.0.0.0",
        Port:     8080,
        Timeout:  30 * time.Second,
        MaxConns: 100,
    }

    // 应用选项
    for _, opt := range opts {
        if err := opt(cfg); err != nil {
            return nil, fmt.Errorf("invalid option: %w", err)
        }
    }

    return cfg, nil
}

// Addr 返回监听地址
func (c *ServerConfig) Addr() string {
    return net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))
}
```

### 使用示例

```go
// 基础配置
 cfg1, err := NewServerConfig()

// 自定义配置
cfg2, err := NewServerConfig(
    WithHost("localhost"),
    WithPort(9090),
    WithTimeout(60*time.Second),
    WithMaxConns(1000),
)

// TLS配置
cfg3, err := NewServerConfig(
    WithHost("example.com"),
    WithPort(443),
    WithTLS("cert.pem", "key.pem"),
)
```

---

## 2. 依赖注入模式

### 构造函数注入

```go
package main

import (
    "context"
    "database/sql"
)

// 接口定义
type UserRepository interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    Save(ctx context.Context, user *User) error
}

type Cache interface {
    Get(key string) ([]byte, bool)
    Set(key string, value []byte, ttl int)
}

type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, err error, fields ...Field)
}

// UserService 依赖通过构造函数注入
type UserService struct {
    repo   UserRepository
    cache  Cache
    logger Logger
}

// NewUserService 构造函数
func NewUserService(
    repo UserRepository,
    cache Cache,
    logger Logger,
) *UserService {
    return &UserService{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
    // 实现...
    return s.repo.GetByID(ctx, id)
}
```

### 使用 Wire 生成依赖注入代码

```go
// wire.go
//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

// InitializeUserService 由Wire生成
func InitializeUserService(db *sql.DB, redis *redis.Client) (*UserService, error) {
    wire.Build(
        NewUserRepository,
        NewCache,
        NewLogger,
        NewUserService,
    )
    return nil, nil
}
```

---

## 3. 仓储模式实现

### 泛型仓储接口

```go
package repository

import (
    "context"
)

// Entity 实体接口约束
type Entity interface {
    GetID() int64
    SetID(id int64)
}

// Repository 泛型仓储接口
type Repository[T Entity] interface {
    GetByID(ctx context.Context, id int64) (T, error)
    GetAll(ctx context.Context) ([]T, error)
    Create(ctx context.Context, entity T) error
    Update(ctx context.Context, entity T) error
    Delete(ctx context.Context, id int64) error
    Exists(ctx context.Context, id int64) (bool, error)
}

// QueryOptions 查询选项
type QueryOptions struct {
    Limit  int
    Offset int
    OrderBy string
    Filters map[string]interface{}
}

// QueryBuilder 查询构建器
type QueryBuilder[T Entity] interface {
    WithOptions(opts QueryOptions) QueryBuilder[T]
    Execute(ctx context.Context) ([]T, error)
    Count(ctx context.Context) (int64, error)
}
```

### 实现示例

```go
package repository

import (
    "context"
    "database/sql"
    "fmt"
)

// SQLRepository SQL实现
type SQLRepository[T Entity] struct {
    db        *sql.DB
    table     string
    scanFunc  func(*sql.Rows) (T, error)
}

// NewSQLRepository 创建SQL仓储
func NewSQLRepository[T Entity](
    db *sql.DB,
    table string,
    scanFunc func(*sql.Rows) (T, error),
) *SQLRepository[T] {
    return &SQLRepository[T]{
        db:       db,
        table:    table,
        scanFunc: scanFunc,
    }
}

func (r *SQLRepository[T]) GetByID(ctx context.Context, id int64) (T, error) {
    var zero T
    query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", r.table)

    row := r.db.QueryRowContext(ctx, query, id)
    entity, err := r.scanRow(row)
    if err == sql.ErrNoRows {
        return zero, ErrNotFound
    }
    if err != nil {
        return zero, fmt.Errorf("query failed: %w", err)
    }

    return entity, nil
}

func (r *SQLRepository[T]) scanRow(row *sql.Row) (T, error) {
    // 具体实现依赖scanFunc
    var zero T
    return zero, nil
}
```

---

## 4. 工作单元模式

```go
package unitofwork

import (
    "context"
    "database/sql"
    "fmt"
)

// UnitOfWork 工作单元
type UnitOfWork struct {
    db     *sql.DB
    tx     *sql.Tx
    ctx    context.Context
}

// NewUnitOfWork 创建工作单元
func NewUnitOfWork(db *sql.DB) *UnitOfWork {
    return &UnitOfWork{db: db}
}

// Begin 开始事务
func (u *UnitOfWork) Begin(ctx context.Context) error {
    if u.tx != nil {
        return fmt.Errorf("transaction already started")
    }

    tx, err := u.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }

    u.tx = tx
    u.ctx = ctx
    return nil
}

// Commit 提交事务
func (u *UnitOfWork) Commit() error {
    if u.tx == nil {
        return fmt.Errorf("no active transaction")
    }

    if err := u.tx.Commit(); err != nil {
        return fmt.Errorf("commit failed: %w", err)
    }

    u.tx = nil
    return nil
}

// Rollback 回滚事务
func (u *UnitOfWork) Rollback() error {
    if u.tx == nil {
        return nil
    }

    if err := u.tx.Rollback(); err != nil {
        return fmt.Errorf("rollback failed: %w", err)
    }

    u.tx = nil
    return nil
}

// Exec 执行SQL
func (u *UnitOfWork) Exec(query string, args ...interface{}) (sql.Result, error) {
    if u.tx == nil {
        return nil, fmt.Errorf("no active transaction")
    }
    return u.tx.ExecContext(u.ctx, query, args...)
}

// QueryRow 查询单行
func (u *UnitOfWork) QueryRow(query string, args ...interface{}) *sql.Row {
    if u.tx == nil {
        return nil
    }
    return u.tx.QueryRowContext(u.ctx, query, args...)
}

// Do 在事务中执行操作
func (u *UnitOfWork) Do(ctx context.Context, fn func(*UnitOfWork) error) error {
    if err := u.Begin(ctx); err != nil {
        return err
    }

    defer u.Rollback()  // 如果成功提交，Rollback不做任何事

    if err := fn(u); err != nil {
        return err
    }

    return u.Commit()
}
```

---

## 5. 事件驱动模式

```go
package eventbus

import (
    "context"
    "reflect"
    "sync"
)

// Event 事件接口
type Event interface {
    EventName() string
}

// Handler 事件处理器
type Handler func(ctx context.Context, event Event) error

// EventBus 事件总线
type EventBus struct {
    handlers map[string][]Handler
    mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]Handler),
    }
}

// Subscribe 订阅事件
func (b *EventBus) Subscribe(eventName string, handler Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Publish 发布事件
func (b *EventBus) Publish(ctx context.Context, event Event) error {
    b.mu.RLock()
    handlers := b.handlers[event.EventName()]
    b.mu.RUnlock()

    for _, h := range handlers {
        if err := h(ctx, event); err != nil {
            return err  // 或收集所有错误
        }
    }

    return nil
}

// PublishAsync 异步发布
func (b *EventBus) PublishAsync(ctx context.Context, event Event) {
    go b.Publish(ctx, event)
}

// 类型安全的事件注册
func SubscribeEvent[T Event](bus *EventBus, handler func(ctx context.Context, event T) error) {
    var zero T
    eventName := zero.EventName()

    bus.Subscribe(eventName, func(ctx context.Context, e Event) error {
        // 类型断言
        typed, ok := e.(T)
        if !ok {
            return nil  // 忽略类型不匹配
        }
        return handler(ctx, typed)
    })
}
```

---

## 6. 断路器模式

```go
package circuitbreaker

import (
    "errors"
    "sync"
    "time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type State int

const (
    StateClosed State = iota   // 正常
    StateOpen                  // 断路
    StateHalfOpen              // 半开
)

// CircuitBreaker 断路器
type CircuitBreaker struct {
    mu                sync.Mutex
    state             State
    failureCount      int
    successCount      int
    lastFailureTime   time.Time

    // 配置
    failureThreshold  int           // 触发断路的失败次数
    successThreshold  int           // 半开状态恢复成功的次数
    timeout           time.Duration // 断路后尝试恢复的时间
}

// New 创建断路器
func New(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failureThreshold: failureThreshold,
        successThreshold: successThreshold,
        timeout:          timeout,
        state:            StateClosed,
    }
}

// Execute 执行操作
func (cb *CircuitBreaker) Execute(fn func() error) error {
    if err := cb.beforeRequest(); err != nil {
        return err
    }

    err := fn()
    cb.afterRequest(err)
    return err
}

func (cb *CircuitBreaker) beforeRequest() error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    switch cb.state {
    case StateClosed:
        return nil
    case StateOpen:
        if time.Since(cb.lastFailureTime) > cb.timeout {
            cb.state = StateHalfOpen
            cb.successCount = 0
            return nil
        }
        return ErrCircuitOpen
    case StateHalfOpen:
        return nil
    }
    return nil
}

func (cb *CircuitBreaker) afterRequest(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err == nil {
        cb.onSuccess()
    } else {
        cb.onFailure()
    }
}

func (cb *CircuitBreaker) onSuccess() {
    switch cb.state {
    case StateClosed:
        cb.failureCount = 0
    case StateHalfOpen:
        cb.successCount++
        if cb.successCount >= cb.successThreshold {
            cb.state = StateClosed
            cb.failureCount = 0
        }
    }
}

func (cb *CircuitBreaker) onFailure() {
    cb.failureCount++
    cb.lastFailureTime = time.Now()

    switch cb.state {
    case StateClosed:
        if cb.failureCount >= cb.failureThreshold {
            cb.state = StateOpen
        }
    case StateHalfOpen:
        cb.state = StateOpen
    }
}

// State 获取当前状态
func (cb *CircuitBreaker) State() State {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    return cb.state
}
```

---

## 7. 限流器模式

```go
package ratelimiter

import (
    "context"
    "sync"
    "time"
)

// TokenBucket 令牌桶限流器
type TokenBucket struct {
    capacity int           // 桶容量
    tokens   float64       // 当前令牌数
    rate     float64       // 每秒产生令牌数
    lastTime time.Time     // 上次更新时间
    mu       sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity int, rate float64) *TokenBucket {
    return &TokenBucket{
        capacity: capacity,
        tokens:   float64(capacity),
        rate:     rate,
        lastTime: time.Now(),
    }
}

// Allow 检查是否允许通过
func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastTime).Seconds()
    tb.lastTime = now

    // 添加新令牌
    tb.tokens += elapsed * tb.rate
    if tb.tokens > float64(tb.capacity) {
        tb.tokens = float64(tb.capacity)
    }

    // 消费令牌
    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }

    return false
}

// Wait 等待直到允许通过
func (tb *TokenBucket) Wait(ctx context.Context) error {
    for {
        if tb.Allow() {
            return nil
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(10 * time.Millisecond):
            // 继续尝试
        }
    }
}

// LeakyBucket 漏桶限流器
type LeakyBucket struct {
    capacity  int
    water     int
    leakRate  time.Duration  // 漏水间隔
    lastLeak  time.Time
    mu        sync.Mutex
}

// NewLeakyBucket 创建漏桶
func NewLeakyBucket(capacity int, leakRate time.Duration) *LeakyBucket {
    return &LeakyBucket{
        capacity: capacity,
        leakRate: leakRate,
        lastLeak: time.Now(),
    }
}

// Add 添加水量
func (lb *LeakyBucket) Add(n int) bool {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    // 漏水
    now := time.Now()
    elapsed := now.Sub(lb.lastLeak)
    leaks := int(elapsed / lb.leakRate)
    lb.water -= leaks
    if lb.water < 0 {
        lb.water = 0
    }
    lb.lastLeak = now.Add(-time.Duration(elapsed%lb.leakRate))

    // 检查容量
    if lb.water+n > lb.capacity {
        return false
    }

    lb.water += n
    return true
}
```

---

## 8. 连接池模式

```go
package pool

import (
    "context"
    "errors"
    "sync"
    "sync/atomic"
    "time"
)

var (
    ErrPoolClosed    = errors.New("pool is closed")
    ErrConnBusy      = errors.New("connection is busy")
    ErrConnTimeout   = errors.New("get connection timeout")
)

// Conn 连接接口
type Conn interface {
    Close() error
    IsValid() bool
}

// Factory 连接工厂
type Factory func() (Conn, error)

// Pool 连接池
type Pool struct {
    factory     Factory
    minIdle     int
    maxActive   int
    maxIdleTime time.Duration
    maxLifetime time.Duration

    idleConns   chan *poolConn
    activeCount int32
    closed      int32
    wg          sync.WaitGroup
}

type poolConn struct {
    Conn
    createdAt time.Time
    lastUsed  time.Time
}

// New 创建连接池
func New(factory Factory, minIdle, maxActive int, maxIdleTime, maxLifetime time.Duration) *Pool {
    p := &Pool{
        factory:     factory,
        minIdle:     minIdle,
        maxActive:   maxActive,
        maxIdleTime: maxIdleTime,
        maxLifetime: maxLifetime,
        idleConns:   make(chan *poolConn, maxActive),
    }

    // 初始化最小连接
    for i := 0; i < minIdle; i++ {
        if conn, err := p.createConn(); err == nil {
            p.idleConns <- conn
        }
    }

    // 启动清理协程
    p.wg.Add(1)
    go p.maintain()

    return p
}

// Get 获取连接
func (p *Pool) Get(ctx context.Context) (Conn, error) {
    if atomic.LoadInt32(&p.closed) == 1 {
        return nil, ErrPoolClosed
    }

    select {
    case conn := <-p.idleConns:
        if p.isValid(conn) {
            conn.lastUsed = time.Now()
            return &PooledConn{pool: p, conn: conn}, nil
        }
        // 无效连接，关闭并创建新连接
        conn.Close()
        atomic.AddInt32(&p.activeCount, -1)
        return p.Get(ctx)

    default:
        // 没有空闲连接，创建新连接
        if int(atomic.LoadInt32(&p.activeCount)) < p.maxActive {
            conn, err := p.createConn()
            if err != nil {
                return nil, err
            }
            return &PooledConn{pool: p, conn: conn}, nil
        }

        // 等待空闲连接
        select {
        case conn := <-p.idleConns:
            conn.lastUsed = time.Now()
            return &PooledConn{pool: p, conn: conn}, nil
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
}

// Put 归还连接
func (p *Pool) Put(conn *poolConn) {
    if atomic.LoadInt32(&p.closed) == 1 {
        conn.Close()
        return
    }

    select {
    case p.idleConns <- conn:
        // 成功放回
    default:
        // 池已满，关闭连接
        conn.Close()
        atomic.AddInt32(&p.activeCount, -1)
    }
}

// Close 关闭连接池
func (p *Pool) Close() error {
    if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
        return nil
    }

    close(p.idleConns)
    for conn := range p.idleConns {
        conn.Close()
    }

    p.wg.Wait()
    return nil
}

func (p *Pool) createConn() (*poolConn, error) {
    conn, err := p.factory()
    if err != nil {
        return nil, err
    }

    atomic.AddInt32(&p.activeCount, 1)
    now := time.Now()
    return &poolConn{
        Conn:      conn,
        createdAt: now,
        lastUsed:  now,
    }, nil
}

func (p *Pool) isValid(conn *poolConn) bool {
    if time.Since(conn.createdAt) > p.maxLifetime {
        return false
    }
    if time.Since(conn.lastUsed) > p.maxIdleTime {
        return false
    }
    return conn.IsValid()
}

func (p *Pool) maintain() {
    defer p.wg.Done()
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // 清理过期连接
            // 实现省略...
        }
    }
}

// PooledConn 包装连接，拦截Close
type PooledConn struct {
    pool *Pool
    conn *poolConn
}

func (pc *PooledConn) Close() error {
    pc.pool.Put(pc.conn)
    return nil
}

func (pc *PooledConn) IsValid() bool {
    return pc.conn.IsValid()
}
```

---

## 总结

这些高级编程模式在Go中广泛应用：

| 模式 | 使用场景 | 关键实现 |
|-----|---------|---------|
| 函数选项 | 复杂配置 | 验证 + 默认值 |
| 依赖注入 | 组件解耦 | 构造函数 + Wire |
| 仓储模式 | 数据访问 | 泛型接口 |
| 工作单元 | 事务管理 | 延迟提交 |
| 事件驱动 | 解耦通信 | 订阅/发布 |
| 断路器 | 容错保护 | 状态机 |
| 限流器 | 流量控制 | 令牌桶/漏桶 |
| 连接池 | 资源复用 | 生命周期管理 |

---

**文档版本**: 1.0
**最后更新**: 2026-03-08
