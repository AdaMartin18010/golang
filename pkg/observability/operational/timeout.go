package operational

import (
	"context"
	"time"
)

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Default    time.Duration
	HTTP       time.Duration
	Database   time.Duration
	External   time.Duration
	Observability time.Duration
}

// DefaultTimeoutConfig 返回默认超时配置
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Default:        30 * time.Second,
		HTTP:           10 * time.Second,
		Database:       15 * time.Second,
		External:       10 * time.Second,
		Observability:  5 * time.Second,
	}
}

// WithTimeout 为操作添加超时
func WithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	if timeout <= 0 {
		return fn(ctx)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(timeoutCtx)
}

// TimeoutMiddleware 超时中间件
// 用于 HTTP 请求的超时控制
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
