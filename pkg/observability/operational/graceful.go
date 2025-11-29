package operational

import (
	"context"
	"fmt"
	"time"
)

// GracefulServer 优雅服务器接口
type GracefulServer interface {
	Shutdown(ctx context.Context) error
}

// GracefulShutdown 执行优雅关闭
// 提供统一的优雅关闭接口
func GracefulShutdown(ctx context.Context, timeout time.Duration, servers ...GracefulServer) error {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var lastErr error
	for _, server := range servers {
		if err := server.Shutdown(shutdownCtx); err != nil {
			lastErr = err
			fmt.Printf("Error shutting down server: %v\n", err)
		}
	}

	if lastErr != nil {
		return fmt.Errorf("graceful shutdown completed with errors: %w", lastErr)
	}

	return nil
}

// GracefulShutdownWithCallback 带回调的优雅关闭
func GracefulShutdownWithCallback(
	ctx context.Context,
	timeout time.Duration,
	shutdownFn func(ctx context.Context) error,
	onShutdown func(),
) error {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 执行关闭函数
	if err := shutdownFn(shutdownCtx); err != nil {
		return err
	}

	// 执行回调
	if onShutdown != nil {
		onShutdown()
	}

	return nil
}

// ShutdownTimeout 关闭超时配置
type ShutdownTimeout struct {
	Default    time.Duration
	HTTP       time.Duration
	Database   time.Duration
	Observability time.Duration
}

// DefaultShutdownTimeout 返回默认关闭超时配置
func DefaultShutdownTimeout() ShutdownTimeout {
	return ShutdownTimeout{
		Default:       30 * time.Second,
		HTTP:          10 * time.Second,
		Database:      15 * time.Second,
		Observability: 10 * time.Second,
	}
}
