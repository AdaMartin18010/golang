# 优雅关闭实现 (Graceful Shutdown Implementation)

> **分类**: 工程与云原生
> **标签**: #graceful-shutdown #lifecycle #signal-handling
> **参考**: Go HTTP Server Shutdown, Kubernetes Pod Lifecycle

---

## 关闭流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Graceful Shutdown Sequence                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. 接收信号                                                                 │
│     ┌──────────┐                                                             │
│     │ SIGTERM │  (Kubernetes default)                                        │
│     │ SIGINT  │  (Ctrl+C)                                                    │
│     └────┬─────┘                                                             │
│          │                                                                   │
│          ▼                                                                   │
│  2. 通知组件开始关闭                                                          │
│     ┌─────────────────────────────────────────────────────────────────┐     │
│     │  Cancel Context ──► Stop Accepting New Connections               │     │
│     └─────────────────────────────────────────────────────────────────┘     │
│          │                                                                   │
│          ▼                                                                   │
│  3. 等待活跃请求完成                                                          │
│     ┌─────────────────────────────────────────────────────────────────┐     │
│     │  ┌──────────┐  ┌──────────┐  ┌──────────┐                       │     │
│     │  │ Request 1│  │ Request 2│  │ Request N│  ...                   │     │
│     │  │  (2s)    │  │  (5s)    │  │  (1s)    │                       │     │
│     │  └────┬─────┘  └────┬─────┘  └────┬─────┘                       │     │
│     │       │            │            │                               │     │
│     │       └────────────┴────────────┘                               │     │
│     │                  │                                              │     │
│     │                  ▼                                              │     │
│     │           Wait for completion                                   │     │
│     │           (with timeout)                                        │     │
│     └─────────────────────────────────────────────────────────────────┘     │
│          │                                                                   │
│          ▼                                                                   │
│  4. 关闭资源                                                                 │
│     ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐                 │
│     │ Database │  │  Redis   │  │  Queue   │  │  Worker  │                 │
│     │  Close   │  │  Close   │  │  Close   │  │  Stop    │                 │
│     └──────────┘  └──────────┘  └──────────┘  └──────────┘                 │
│          │                                                                   │
│          ▼                                                                   │
│  5. 退出                                                                     │
│     ┌──────────┐                                                             │
│     │  os.Exit │  (with exit code)                                           │
│     └──────────┘                                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心实现

```go
package graceful

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

// ShutdownFunc 关闭函数类型
type ShutdownFunc func(ctx context.Context) error

// Manager 优雅关闭管理器
type Manager struct {
    timeout     time.Duration
    shutdownFuncs []ShutdownFunc
    mu          sync.RWMutex

    // 信号处理
    signals     []os.Signal
    signalChan  chan os.Signal

    // 状态
    shuttingDown bool
    wg           sync.WaitGroup
}

// NewManager 创建关闭管理器
func NewManager(timeout time.Duration) *Manager {
    return &Manager{
        timeout:    timeout,
        signals:    []os.Signal{syscall.SIGTERM, syscall.SIGINT},
        signalChan: make(chan os.Signal, 1),
    }
}

// AddShutdownFunc 添加关闭函数
func (m *Manager) AddShutdownFunc(fn ShutdownFunc) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.shutdownFuncs = append(m.shutdownFuncs, fn)
}

// Listen 监听信号
func (m *Manager) Listen(ctx context.Context) error {
    // 注册信号处理
    signal.Notify(m.signalChan, m.signals...)
    defer signal.Stop(m.signalChan)

    select {
    case sig := <-m.signalChan:
        fmt.Printf("Received signal: %v\n", sig)
        return m.Shutdown(ctx)
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Shutdown 执行关闭
func (m *Manager) Shutdown(ctx context.Context) error {
    m.mu.Lock()
    if m.shuttingDown {
        m.mu.Unlock()
        return fmt.Errorf("already shutting down")
    }
    m.shuttingDown = true
    funcs := make([]ShutdownFunc, len(m.shutdownFuncs))
    copy(funcs, m.shutdownFuncs)
    m.mu.Unlock()

    // 创建带超时的上下文
    shutdownCtx, cancel := context.WithTimeout(ctx, m.timeout)
    defer cancel()

    fmt.Printf("Starting graceful shutdown (timeout: %v)...\n", m.timeout)

    // 逆序执行关闭函数
    var wg sync.WaitGroup
    errChan := make(chan error, len(funcs))

    for i := len(funcs) - 1; i >= 0; i-- {
        wg.Add(1)
        go func(fn ShutdownFunc) {
            defer wg.Done()
            if err := fn(shutdownCtx); err != nil {
                errChan <- err
            }
        }(funcs[i])
    }

    // 等待关闭完成或超时
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        fmt.Println("Graceful shutdown completed")
    case <-shutdownCtx.Done():
        fmt.Printf("Shutdown timeout exceeded: %v\n", shutdownCtx.Err())
    }

    close(errChan)

    // 收集错误
    var errs []error
    for err := range errChan {
        errs = append(errs, err)
    }

    if len(errs) > 0 {
        return fmt.Errorf("shutdown errors: %v", errs)
    }

    return nil
}

// HTTPServerShutdown HTTP 服务器关闭
func HTTPServerShutdown(server *http.Server) ShutdownFunc {
    return func(ctx context.Context) error {
        fmt.Println("Shutting down HTTP server...")
        return server.Shutdown(ctx)
    }
}

// WorkerPoolShutdown 工作池关闭
func WorkerPoolShutdown(stop chan struct{}, wg *sync.WaitGroup) ShutdownFunc {
    return func(ctx context.Context) error {
        fmt.Println("Shutting down worker pool...")
        close(stop)

        done := make(chan struct{})
        go func() {
            wg.Wait()
            close(done)
        }()

        select {
        case <-done:
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// DatabaseShutdown 数据库连接关闭
func DatabaseShutdown(close func() error) ShutdownFunc {
    return func(ctx context.Context) error {
        fmt.Println("Closing database connections...")
        return close()
    }
}
```

---

## HTTP 服务器优雅关闭

```go
package graceful

import (
    "context"
    "fmt"
    "net"
    "net/http"
    "sync"
    "sync/atomic"
    "time"
)

// GracefulHTTPServer 优雅关闭的 HTTP 服务器
type GracefulHTTPServer struct {
    server       *http.Server
    listener     net.Listener

    // 活跃连接管理
    activeConns  int64
    connWG       sync.WaitGroup

    // 关闭控制
    shutdownChan chan struct{}
    shutdownOnce sync.Once
}

// NewGracefulHTTPServer 创建优雅关闭的 HTTP 服务器
func NewGracefulHTTPServer(handler http.Handler) *GracefulHTTPServer {
    return &GracefulHTTPServer{
        server: &http.Server{
            Handler: handler,
            ConnState: func(conn net.Conn, state http.ConnState) {
                // 连接状态管理
            },
        },
        shutdownChan: make(chan struct{}),
    }
}

// ListenAndServe 启动服务器
func (s *GracefulHTTPServer) ListenAndServe(addr string) error {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    s.listener = listener

    // 包装监听器以追踪连接
    s.listener = &gracefulListener{
        Listener:     listener,
        server:       s,
    }

    return s.server.Serve(s.listener)
}

// Shutdown 优雅关闭
func (s *GracefulHTTPServer) Shutdown(ctx context.Context) error {
    s.shutdownOnce.Do(func() {
        close(s.shutdownChan)
    })

    // 关闭监听器（停止接受新连接）
    if s.listener != nil {
        s.listener.Close()
    }

    // 等待活跃连接完成
    done := make(chan struct{})
    go func() {
        s.connWG.Wait()
        close(done)
    }()

    select {
    case <-done:
        fmt.Println("All connections closed gracefully")
    case <-ctx.Done():
        fmt.Println("Shutdown timeout, forcing close")
    }

    // 关闭服务器
    return s.server.Shutdown(ctx)
}

// trackConnection 追踪连接
func (s *GracefulHTTPServer) trackConnection() {
    atomic.AddInt64(&s.activeConns, 1)
    s.connWG.Add(1)
}

// untrackConnection 取消追踪连接
func (s *GracefulHTTPServer) untrackConnection() {
    atomic.AddInt64(&s.activeConns, -1)
    s.connWG.Done()
}

// ActiveConnections 获取活跃连接数
func (s *GracefulHTTPServer) ActiveConnections() int64 {
    return atomic.LoadInt64(&s.activeConns)
}

// gracefulListener 包装监听器
type gracefulListener struct {
    net.Listener
    server *GracefulHTTPServer
}

func (l *gracefulListener) Accept() (net.Conn, error) {
    for {
        conn, err := l.Listener.Accept()
        if err != nil {
            select {
            case <-l.server.shutdownChan:
                return nil, fmt.Errorf("server shutting down")
            default:
            }
            return nil, err
        }

        // 追踪连接
        l.server.trackConnection()

        return &gracefulConn{
            Conn:   conn,
            server: l.server,
        }, nil
    }
}

// gracefulConn 包装连接
type gracefulConn struct {
    net.Conn
    server *GracefulHTTPServer
    once   sync.Once
}

func (c *gracefulConn) Close() error {
    c.once.Do(c.server.untrackConnection)
    return c.Conn.Close()
}
```

---

## 综合应用

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "sync"
    "time"

    "graceful"
)

func main() {
    // 创建关闭管理器
    manager := graceful.NewManager(30 * time.Second)

    // 创建 HTTP 服务器
    httpServer := graceful.NewGracefulHTTPServer(http.HandlerFunc(handler))

    // 创建工作池
    workerStop := make(chan struct{})
    var workerWG sync.WaitGroup
    startWorkers(workerStop, &workerWG)

    // 数据库连接
    db, err := sql.Open("postgres", "...")
    if err != nil {
        panic(err)
    }

    // 注册关闭函数（按依赖顺序逆序）
    manager.AddShutdownFunc(graceful.HTTPServerShutdown(httpServer.server))
    manager.AddShutdownFunc(graceful.WorkerPoolShutdown(workerStop, &workerWG))
    manager.AddShutdownFunc(graceful.DatabaseShutdown(db.Close))

    // 启动 HTTP 服务器（在 goroutine 中）
    go func() {
        if err := httpServer.ListenAndServe(":8080"); err != nil && err != http.ErrServerClosed {
            fmt.Printf("HTTP server error: %v\n", err)
        }
    }()

    fmt.Println("Server started on :8080")

    // 监听关闭信号
    if err := manager.Listen(context.Background()); err != nil {
        fmt.Printf("Shutdown error: %v\n", err)
    }

    fmt.Println("Application exited")
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 模拟长时间运行的请求
    time.Sleep(5 * time.Second)
    fmt.Fprintln(w, "OK")
}

func startWorkers(stop chan struct{}, wg *sync.WaitGroup) {
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go worker(stop, wg, i)
    }
}

func worker(stop chan struct{}, wg *sync.WaitGroup, id int) {
    defer wg.Done()

    for {
        select {
        case <-stop:
            fmt.Printf("Worker %d stopping\n", id)
            // 完成当前工作
            time.Sleep(1 * time.Second)
            fmt.Printf("Worker %d stopped\n", id)
            return
        default:
            // 处理工作
            time.Sleep(100 * time.Millisecond)
        }
    }
}
```

---

## Kubernetes 健康检查

```go
package graceful

import (
    "net/http"
    "sync/atomic"
)

// HealthChecker 健康检查器
type HealthChecker struct {
    ready   int32
    alive   int32
}

func NewHealthChecker() *HealthChecker {
    hc := &HealthChecker{}
    atomic.StoreInt32(&hc.alive, 1)
    return hc
}

// Ready 设置就绪状态
func (hc *HealthChecker) Ready(ready bool) {
    if ready {
        atomic.StoreInt32(&hc.ready, 1)
    } else {
        atomic.StoreInt32(&hc.ready, 0)
    }
}

// Shutdown 设置关闭状态
func (hc *HealthChecker) Shutdown() {
    atomic.StoreInt32(&hc.ready, 0)
    atomic.StoreInt32(&hc.alive, 0)
}

// LivenessHandler 存活检查
func (hc *HealthChecker) LivenessHandler(w http.ResponseWriter, r *http.Request) {
    if atomic.LoadInt32(&hc.alive) == 1 {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("alive"))
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("not alive"))
    }
}

// ReadinessHandler 就绪检查
func (hc *HealthChecker) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
    if atomic.LoadInt32(&hc.ready) == 1 {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ready"))
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("not ready"))
    }
}
```
