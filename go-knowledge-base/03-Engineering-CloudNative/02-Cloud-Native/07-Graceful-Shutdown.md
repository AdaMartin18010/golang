# 优雅关闭 (Graceful Shutdown)

> **分类**: 工程与云原生
> **标签**: #graceful-shutdown #context #signal

---

## 基础实现

### HTTP 服务优雅关闭

```go
func main() {
    srv := &http.Server{
        Addr:    ":8080",
        Handler: router(),
    }

    // 启动服务
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // 超时上下文
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}
```

---

## 多服务协调关闭

```go
type App struct {
    services []Service
}

type Service interface {
    Start() error
    Stop(ctx context.Context) error
}

func (a *App) Run() error {
    // 启动所有服务
    for _, s := range a.services {
        if err := s.Start(); err != nil {
            return err
        }
    }

    // 等待信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // 优雅关闭
    return a.shutdown()
}

func (a *App) shutdown() error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    var wg sync.WaitGroup
    errChan := make(chan error, len(a.services))

    for _, s := range a.services {
        wg.Add(1)
        go func(svc Service) {
            defer wg.Done()
            if err := svc.Stop(ctx); err != nil {
                errChan <- err
            }
        }(s)
    }

    wg.Wait()
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
```

---

## 连接池关闭

```go
type Database struct {
    *sql.DB
}

func (db *Database) Stop(ctx context.Context) error {
    done := make(chan struct{})

    go func() {
        db.Close()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

## Worker 池关闭

```go
type WorkerPool struct {
    workers int
    jobs    chan Job
    wg      sync.WaitGroup
}

func (wp *WorkerPool) Stop(ctx context.Context) error {
    // 停止接收新任务
    close(wp.jobs)

    // 等待完成或超时
    done := make(chan struct{})
    go func() {
        wp.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return fmt.Errorf("timeout waiting for workers: %w", ctx.Err())
    }
}
```

---

## Kubernetes 优雅关闭

```go
func main() {
    // 监听预停止钩子
    http.HandleFunc("/pre-stop", func(w http.ResponseWriter, r *http.Request) {
        // 标记不再接收新连接
        healthz.SetReady(false)

        // 等待现有请求完成
        time.Sleep(10 * time.Second)

        w.WriteHeader(http.StatusOK)
    })
}
```
