package middleware

import (
	"net/http"
)

// TracingConfig 追踪中间件配置
type TracingConfig struct {
	TracerName     string
	ServiceName    string
	ServiceVersion string
	SkipPaths      []string
}

// TracingMiddleware 创建追踪中间件（简化版，仅添加基本头信息）
func TracingMiddleware(config TracingConfig) func(http.Handler) http.Handler {
	skipMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipMap[path] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 跳过指定路径
			if skipMap[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			// 添加服务信息到响应头
			w.Header().Set("X-Service-Name", config.ServiceName)
			w.Header().Set("X-Service-Version", config.ServiceVersion)

			next.ServeHTTP(w, r)
		})
	}
}
