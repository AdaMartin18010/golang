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

// RecoveryConfig 是恢复中间件的配置结构。
//
// 功能说明：
// - 配置 Panic 恢复行为
// - 控制堆栈信息输出
//
// 字段说明：
// - Logger: 日志记录器（用于记录 Panic 信息）
//   如果为 nil，使用默认的错误级别日志记录器
// - StackAll: 是否在所有情况下输出堆栈（当前未使用）
// - StackSize: 堆栈大小限制（默认：8KB）
//   超过此大小的堆栈会被截断
// - LogRequest: 是否记录请求信息（当前未使用）
//
// 使用示例：
//
//	config := middleware.RecoveryConfig{
//	    Logger:    logger,
//	    StackSize: 1024 * 16, // 16KB
//	}
//	router.Use(middleware.RecoveryMiddleware(config))
type RecoveryConfig struct {
	Logger     *logger.Logger
	StackAll   bool   // 是否在所有情况下输出堆栈
	StackSize  int    // 堆栈大小限制
	LogRequest bool   // 是否记录请求信息
}

// RecoveryMiddleware 创建 Panic 恢复中间件。
//
// 功能说明：
// - 捕获 HTTP 请求处理过程中的 Panic
// - 记录 Panic 信息和堆栈跟踪
// - 返回 500 错误响应，避免程序崩溃
//
// 工作流程：
// 1. 使用 defer 和 recover 捕获 Panic
// 2. 获取堆栈信息
// 3. 截断堆栈（如果超过限制）
// 4. 记录错误日志（包含 Panic 信息、请求信息、堆栈）
// 5. 返回 500 错误响应
//
// 记录的日志信息：
// - error: Panic 错误信息
// - method: HTTP 方法
// - path: 请求路径
// - stack: 堆栈跟踪信息
//
// 参数：
// - config: 恢复配置
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	config := middleware.RecoveryConfig{
//	    Logger:    logger,
//	    StackSize: 1024 * 8,
//	}
//	router.Use(middleware.RecoveryMiddleware(config))
//
// 注意事项：
// - 应该在中间件链的最外层使用
// - 确保日志记录器已正确配置
// - 堆栈信息可能包含敏感信息，生产环境应谨慎处理
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
