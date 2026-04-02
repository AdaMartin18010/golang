# EC-007: 优雅关闭完整实现 (Graceful Shutdown Complete)

> **维度**: Engineering CloudNative
> **级别**: S (15+ KB)
> **标签**: #graceful-shutdown #signal-handling #kubernetes #zero-downtime
> **相关**: EC-042, EC-109, FT-012

---

## 整合说明

本文档合并了以下历史文档：

- `07-Graceful-Shutdown.md` (3.4 KB) - 基础概念
- `120-Task-Graceful-Shutdown-Complete.md` (8.8 KB) - 生产实现

---

## 核心问题

分布式系统中，如何在不中断活跃请求的情况下安全退出进程？

```
┌─────────────────────────────────────────────────────────────────────┐
│                       优雅关闭流程                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  SIGTERM                                                           │
│     │                                                               │
│     ▼                                                               │
│  ┌──────────────┐                                                  │
│  │ 停止接受新请求 │ ◄── HTTP Server Shutdown                        │
│  └──────────────┘                                                  │
│     │                                                               │
│     ▼                                                               │
│  ┌──────────────┐                                                  │
│  │ 等待活跃请求完成│ ◄── Context Cancellation + WaitGroup            │
│  └──────────────┘                                                  │
│     │                                                               │
│     ▼                                                               │
│  ┌──────────────┐                                                  │
│  │ 执行关闭钩子  │ ◄── 数据库、缓存、消息队列                        │
│  └──────────────┘                                                  │
│     │                                                               │
│     ▼                                                               │
│  ┌──────────────┐                                                  │
│  │ 刷新缓冲区   │ ◄── 日志、指标                                    │
│  └──────────────┘                                                  │
│     │                                                               │
│     ▼                                                               │
│   退出码 0                                                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 完整实现

```go
package graceful

import (
 "context"
 "errors"
 "net/http"
 "os"
 "os/signal"
 "sort"
 "sync"
 "sync/atomic"
 "syscall"
 "time"

 "go.uber.org/zap"
)

// ShutdownManager 关闭管理器
type ShutdownManager struct {
 logger *zap.Logger

 // 超时配置
 hooksTimeout   time.Duration // Hook 执行超时
 forceExitDelay time.Duration // 强制退出延迟

 // 注册的关闭钩子
 hooks   []ShutdownHook
 hooksMu sync.RWMutex

 // HTTP 服务器
 servers []*http.Server

 // 活跃请求计数
 activeRequests atomic.Int64

 // 状态
 shuttingDown atomic.Bool
 done         chan struct{}
}

// ShutdownHook 关闭钩子
type ShutdownHook struct {
 Name     string
 Priority int // 优先级（越小越早执行）
 Fn       func(ctx context.Context) error
}

// NewShutdownManager 创建关闭管理器
func NewShutdownManager(logger *zap.Logger) *ShutdownManager {
 return &ShutdownManager{
  logger:         logger,
  hooksTimeout:   30 * time.Second,
  forceExitDelay: 60 * time.Second,
  done:           make(chan struct{}),
 }
}

// RegisterHook 注册关闭钩子
func (sm *ShutdownManager) RegisterHook(hook ShutdownHook) {
 sm.hooksMu.Lock()
 defer sm.hooksMu.Unlock()

 sm.hooks = append(sm.hooks, hook)

 // 按优先级排序
 sort.Slice(sm.hooks, func(i, j int) bool {
  return sm.hooks[i].Priority < sm.hooks[j].Priority
 })
}

// RegisterServer 注册 HTTP 服务器
func (sm *ShutdownManager) RegisterServer(server *http.Server) {
 sm.servers = append(sm.servers, server)
}

// WrapHandler 包装 HTTP Handler（跟踪活跃请求）
func (sm *ShutdownManager) WrapHandler(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  // 检查是否正在关闭
  if sm.shuttingDown.Load() {
   w.WriteHeader(http.StatusServiceUnavailable)
   w.Write([]byte("Server is shutting down"))
   return
  }

  // 增加计数
  sm.activeRequests.Add(1)
  defer sm.activeRequests.Add(-1)

  next.ServeHTTP(w, r)
 })
}

// Start 启动信号监听
func (sm *ShutdownManager) Start() {
 sigCh := make(chan os.Signal, 1)
 signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

 go func() {
  sig := <-sigCh
  sm.logger.Info("received shutdown signal", zap.String("signal", sig.String()))

  if err := sm.Shutdown(context.Background()); err != nil {
   sm.logger.Error("shutdown failed", zap.Error(err))
   os.Exit(1)
  }

  os.Exit(0)
 }()
}

// Shutdown 执行关闭
func (sm *ShutdownManager) Shutdown(ctx context.Context) error {
 // 标记正在关闭
 if !sm.shuttingDown.CompareAndSwap(false, true) {
  return errors.New("already shutting down")
 }

 close(sm.done)

 ctx, cancel := context.WithTimeout(ctx, sm.hooksTimeout)
 defer cancel()

 sm.logger.Info("starting graceful shutdown")

 // 1. 停止 HTTP 服务器（停止接受新连接）
 if err := sm.shutdownServers(ctx); err != nil {
  sm.logger.Error("server shutdown failed", zap.Error(err))
 }

 // 2. 等待活跃请求完成
 if err := sm.waitForRequests(ctx); err != nil {
  sm.logger.Error("wait for requests failed", zap.Error(err))
 }

 // 3. 执行关闭钩子
 if err := sm.executeHooks(ctx); err != nil {
  sm.logger.Error("hooks execution failed", zap.Error(err))
 }

 sm.logger.Info("graceful shutdown completed")
 return nil
}

// shutdownServers 关闭 HTTP 服务器
func (sm *ShutdownManager) shutdownServers(ctx context.Context) error {
 var wg sync.WaitGroup
 errCh := make(chan error, len(sm.servers))

 for _, server := range sm.servers {
  wg.Add(1)
  go func(srv *http.Server) {
   defer wg.Done()

   sm.logger.Info("shutting down server", zap.String("addr", srv.Addr))

   if err := srv.Shutdown(ctx); err != nil {
    errCh <- err
   }
  }(server)
 }

 wg.Wait()
 close(errCh)

 for err := range errCh {
  if err != nil {
   return err
  }
 }

 return nil
}

// waitForRequests 等待活跃请求完成
func (sm *ShutdownManager) waitForRequests(ctx context.Context) error {
 ticker := time.NewTicker(100 * time.Millisecond)
 defer ticker.Stop()

 for {
  select {
  case <-ctx.Done():
   return ctx.Err()
  case <-ticker.C:
   active := sm.activeRequests.Load()
   sm.logger.Info("waiting for requests", zap.Int64("active", active))

   if active == 0 {
    return nil
   }
  }
 }
}

// executeHooks 执行关闭钩子
func (sm *ShutdownManager) executeHooks(ctx context.Context) error {
 sm.hooksMu.RLock()
 hooks := make([]ShutdownHook, len(sm.hooks))
 copy(hooks, sm.hooks)
 sm.hooksMu.RUnlock()

 var wg sync.WaitGroup
 errCh := make(chan error, len(hooks))

 for _, hook := range hooks {
  wg.Add(1)
  go func(h ShutdownHook) {
   defer wg.Done()

   sm.logger.Info("executing shutdown hook", zap.String("name", h.Name))

   hookCtx, cancel := context.WithTimeout(ctx, sm.hooksTimeout/2)
   defer cancel()

   if err := h.Fn(hookCtx); err != nil {
    sm.logger.Error("hook failed", zap.String("name", h.Name), zap.Error(err))
    errCh <- err
   }
  }(hook)
 }

 wg.Wait()
 close(errCh)

 for err := range errCh {
  if err != nil {
   return err
  }
 }

 return nil
}

// Done 返回关闭信号通道
func (sm *ShutdownManager) Done() <-chan struct{} {
 return sm.done
}

// IsShuttingDown 检查是否正在关闭
func (sm *ShutdownManager) IsShuttingDown() bool {
 return sm.shuttingDown.Load()
}
```

---

## Kubernetes 集成

```go
// HealthCheckHandler 健康检查处理器
type HealthCheckHandler struct {
 shutdownManager *ShutdownManager
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.URL.Path {
 case "/healthz":
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("ok"))
 case "/readyz":
  if h.shutdownManager.IsShuttingDown() {
   w.WriteHeader(http.StatusServiceUnavailable)
   w.Write([]byte("shutting down"))
   return
  }
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("ready"))
 }
}
```

---

## 使用示例

```go
func main() {
 logger, _ := zap.NewProduction()
 defer logger.Sync()

 shutdownManager := graceful.NewShutdownManager(logger)

 // 注册钩子
 shutdownManager.RegisterHook(graceful.ShutdownHook{
  Name:     "database",
  Priority: 1,
  Fn: func(ctx context.Context) error {
   return db.Close()
  },
 })

 // 创建服务器
 mux := http.NewServeMux()
 mux.HandleFunc("/api/tasks", taskHandler)

 server := &http.Server{
  Addr:    ":8080",
  Handler: shutdownManager.WrapHandler(mux),
 }

 shutdownManager.RegisterServer(server)
 shutdownManager.Start()

 logger.Info("starting server", zap.String("addr", server.Addr))
 if err := server.ListenAndServe(); err != http.ErrServerClosed {
  logger.Fatal("server failed", zap.Error(err))
 }

 <-shutdownManager.Done()
}
```

---

## 关键设计决策

| 决策 | 理由 |
|------|------|
| Priority 排序钩子 | 确保依赖顺序（如先关 DB 再关缓存） |
| 原子状态标记 | 避免竞态条件 |
| 超时控制 | 防止无限等待 |
| WaitGroup 并行 | 加速多服务器关闭 |
