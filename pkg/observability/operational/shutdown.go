package operational

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownManager 优雅关闭管理器
type ShutdownManager struct {
	shutdownFuncs []ShutdownFunc
	mu            sync.Mutex
	timeout       time.Duration
	signals       []os.Signal
}

// ShutdownFunc 关闭函数类型
type ShutdownFunc func(ctx context.Context) error

// NewShutdownManager 创建优雅关闭管理器
func NewShutdownManager(timeout time.Duration) *ShutdownManager {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &ShutdownManager{
		shutdownFuncs: make([]ShutdownFunc, 0),
		timeout:       timeout,
		signals:       []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	}
}

// Register 注册关闭函数
func (sm *ShutdownManager) Register(fn ShutdownFunc) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.shutdownFuncs = append(sm.shutdownFuncs, fn)
}

// WaitForShutdown 等待关闭信号并执行优雅关闭
func (sm *ShutdownManager) WaitForShutdown() error {
	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sm.signals...)

	// 等待信号
	sig := <-sigChan
	fmt.Printf("Received signal: %v, starting graceful shutdown...\n", sig)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), sm.timeout)
	defer cancel()

	// 执行所有关闭函数
	return sm.Shutdown(ctx)
}

// Shutdown 执行优雅关闭
func (sm *ShutdownManager) Shutdown(ctx context.Context) error {
	sm.mu.Lock()
	funcs := make([]ShutdownFunc, len(sm.shutdownFuncs))
	copy(funcs, sm.shutdownFuncs)
	sm.mu.Unlock()

	// 使用 WaitGroup 等待所有关闭函数完成
	var wg sync.WaitGroup
	errChan := make(chan error, len(funcs))

	for _, fn := range funcs {
		wg.Add(1)
		go func(f ShutdownFunc) {
			defer wg.Done()
			if err := f(ctx); err != nil {
				errChan <- err
			}
		}(fn)
	}

	// 等待所有关闭函数完成或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 所有关闭函数完成
		close(errChan)
		var errs []error
		for err := range errChan {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			return fmt.Errorf("shutdown errors: %v", errs)
		}
		return nil
	case <-ctx.Done():
		// 超时
		return fmt.Errorf("shutdown timeout after %v", sm.timeout)
	}
}

// GracefulShutdown 优雅关闭辅助函数
// 用于包装需要优雅关闭的服务
func GracefulShutdown(serviceName string, shutdownFn func(ctx context.Context) error) ShutdownFunc {
	return func(ctx context.Context) error {
		fmt.Printf("Shutting down %s...\n", serviceName)
		if err := shutdownFn(ctx); err != nil {
			return fmt.Errorf("failed to shutdown %s: %w", serviceName, err)
		}
		fmt.Printf("%s shutdown complete\n", serviceName)
		return nil
	}
}
