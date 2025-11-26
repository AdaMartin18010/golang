package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/http/response"
	"github.com/yourusername/golang/pkg/logger"
	"log/slog"
)

// RecoveryConfig 恢复中间件配置
type RecoveryConfig struct {
	Logger     *logger.Logger
	StackAll   bool   // 是否在所有情况下输出堆栈
	StackSize  int    // 堆栈大小限制
	LogRequest bool   // 是否记录请求信息
}

// RecoveryMiddleware 恢复中间件（捕获panic）
func RecoveryMiddleware(config RecoveryConfig) func(http.Handler) http.Handler {
	if config.Logger == nil {
		config.Logger = logger.NewLogger(slog.LevelError)
	}
	if config.StackSize == 0 {
		config.StackSize = 1024 * 8
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// 获取堆栈信息
					stack := debug.Stack()
					stackSize := len(stack)
					if stackSize > config.StackSize {
						stack = stack[:config.StackSize]
					}

					// 记录错误
					config.Logger.Error(
						"Panic recovered",
						"error", err,
						"method", r.Method,
						"path", r.URL.Path,
						"stack", string(stack),
					)

					// 返回错误响应
					appErr := errors.NewInternalError(
						"Internal server error",
						fmt.Errorf("%v", err),
					)
					response.Error(w, http.StatusInternalServerError, appErr)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
