# 优雅关闭完整实现 (Graceful Shutdown Complete Implementation)

> **分类**: 工程与云原生
> **标签**: #graceful-shutdown #signal-handling #cleanup #zero-downtime
> **参考**: Kubernetes Pod Lifecycle, Systemd, Go 1.8+ Shutdown Patterns

---

## 关闭信号流程

```
OS/K8s                    Application
  │                            │
  ├──── SIGTERM ─────────────► │
  │                            │ 1. 停止接受新请求
  │                            │ 2. 等待活跃请求完成
  │                            │ 3. 关闭数据库连接
  │                            │ 4. 刷新缓冲区
  │                            │ 5. 退出
  │◄──── 退出代码 0 ─────────── ┤
  │                            │
  │ (如超时未退出)              │
  ├──── SIGKILL ─────────────► │ 强制终止
```

---

## 完整优雅关闭实现

```go
package graceful

import (
 "context"
 "errors"
 "net/http"
 "os"
 "os/signal"
 "sync"
 "syscall"
 "time"

 "go.uber.org/zap"
)

// ShutdownManager 关闭管理器
type ShutdownManager struct {
 logger *zap.Logger

 // 超时配置
 hooksTimeout    time.Duration // Hook 执行超时
 forceExitDelay  time.Duration // 强制退出延迟

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

  // 包装 ResponseWriter 以检测客户端断开
  wrapped := &responseWriter{ResponseWriter: w, ctx: r.Context()}

  next.ServeHTTP(wrapped, r)
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

## Kubernetes 健康检查集成

```go
// HealthCheckHandler 健康检查处理器
type HealthCheckHandler struct {
 shutdownManager *ShutdownManager

 // 就绪检查
 readyChecks map[string]HealthCheck

 // 存活检查
 liveChecks map[string]HealthCheck
}

type HealthCheck func() error

// ServeHTTP 处理健康检查
func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 switch r.URL.Path {
 case "/healthz":
  // 存活检查
  h.handleLiveness(w, r)
 case "/readyz":
  // 就绪检查
  h.handleReadiness(w, r)
 }
}

func (h *HealthCheckHandler) handleReadiness(w http.ResponseWriter, r *http.Request) {
 // 如果正在关闭，返回未就绪
 if h.shutdownManager.IsShuttingDown() {
  w.WriteHeader(http.StatusServiceUnavailable)
  w.Write([]byte("shutting down"))
  return
 }

 // 执行就绪检查
 for name, check := range h.readyChecks {
  if err := check(); err != nil {
   w.WriteHeader(http.StatusServiceUnavailable)
   fmt.Fprintf(w, "%s: %v", name, err)
   return
  }
 }

 w.WriteHeader(http.StatusOK)
 w.Write([]byte("ok"))
}
```

---

## 完整示例

```go
func main() {
 logger, _ := zap.NewProduction()
 defer logger.Sync()

 // 创建关闭管理器
 shutdownManager := graceful.NewShutdownManager(logger)

 // 注册数据库关闭钩子
 shutdownManager.RegisterHook(graceful.ShutdownHook{
  Name:     "database",
  Priority: 1,
  Fn: func(ctx context.Context) error {
   return db.Close()
  },
 })

 // 注册缓存关闭钩子
 shutdownManager.RegisterHook(graceful.ShutdownHook{
  Name:     "cache",
  Priority: 2,
  Fn: func(ctx context.Context) error {
   return cache.Close()
  },
 })

 // 创建 HTTP 服务器
 mux := http.NewServeMux()
 mux.HandleFunc("/api/tasks", taskHandler)

 // 包装处理函数
 handler := shutdownManager.WrapHandler(mux)

 server := &http.Server{
  Addr:    ":8080",
  Handler: handler,
 }

 shutdownManager.RegisterServer(server)

 // 启动信号监听
 shutdownManager.Start()

 // 启动服务器
 logger.Info("starting server", zap.String("addr", server.Addr))
 if err := server.ListenAndServe(); err != http.ErrServerClosed {
  logger.Fatal("server failed", zap.Error(err))
 }

 // 等待关闭完成
 <-shutdownManager.Done()
}
```
