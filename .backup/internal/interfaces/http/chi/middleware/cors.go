package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// CORSConfig 是 CORS（跨域资源共享）中间件的配置结构。
//
// 功能说明：
// - 配置跨域请求的处理规则
// - 支持预检请求（OPTIONS）处理
//
// 字段说明：
// - AllowedOrigins: 允许的源列表（默认：["*"]）
//   - "*": 允许所有源
//   - 具体域名：只允许指定的源
// - AllowedMethods: 允许的 HTTP 方法（默认：GET、POST、PUT、DELETE、OPTIONS、PATCH）
// - AllowedHeaders: 允许的请求头（默认：Content-Type、Authorization、X-Requested-With）
// - ExposedHeaders: 暴露的响应头（客户端可以访问的响应头）
// - AllowCredentials: 是否允许携带凭证（Cookie、Authorization 等）
//   - true: 允许携带凭证
//   - false: 不允许携带凭证
//   - 注意：如果设置为 true，AllowedOrigins 不能包含 "*"
// - MaxAge: 预检请求缓存时间（秒，默认：3600）
//   浏览器会缓存预检请求的结果，在此时间内不会再次发送预检请求
//
// 使用示例：
//
//	config := middleware.CORSConfig{
//	    AllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
//	    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
//	    AllowedHeaders:   []string{"Content-Type", "Authorization"},
//	    AllowCredentials: true,
//	    MaxAge:           3600,
//	}
//	router.Use(middleware.CORSMiddleware(config))
type CORSConfig struct {
	AllowedOrigins   []string // 允许的源
	AllowedMethods   []string // 允许的方法
	AllowedHeaders   []string // 允许的请求头
	ExposedHeaders   []string // 暴露的响应头
	AllowCredentials bool     // 是否允许凭证
	MaxAge           int      // 预检请求缓存时间（秒）
}

// CORSMiddleware 创建 CORS 中间件。
//
// 功能说明：
// - 处理跨域请求
// - 设置 CORS 响应头
// - 处理预检请求（OPTIONS）
//
// 工作流程：
// 1. 检查请求的 Origin 是否在允许列表中
// 2. 如果不在允许列表中，继续处理请求（不设置 CORS 头）
// 3. 设置 CORS 响应头：
//    - Access-Control-Allow-Origin: 允许的源
//    - Access-Control-Allow-Credentials: 是否允许凭证
//    - Access-Control-Allow-Methods: 允许的方法
//    - Access-Control-Allow-Headers: 允许的请求头
//    - Access-Control-Expose-Headers: 暴露的响应头
//    - Access-Control-Max-Age: 预检请求缓存时间
// 4. 如果是预检请求（OPTIONS），返回 204 No Content
// 5. 否则继续处理请求
//
// CORS 响应头说明：
// - Access-Control-Allow-Origin: 允许的源（必须）
// - Access-Control-Allow-Credentials: 是否允许凭证
// - Access-Control-Allow-Methods: 允许的 HTTP 方法
// - Access-Control-Allow-Headers: 允许的请求头
// - Access-Control-Expose-Headers: 暴露的响应头（客户端可以访问）
// - Access-Control-Max-Age: 预检请求缓存时间
//
// 预检请求：
// - 浏览器在发送跨域请求前会先发送 OPTIONS 请求
// - 中间件会直接返回 204 No Content，不继续处理
//
// 参数：
// - config: CORS 配置
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	// 允许所有源
//	config := middleware.CORSConfig{
//	    AllowedOrigins: []string{"*"},
//	}
//	router.Use(middleware.CORSMiddleware(config))
//
//	// 允许特定源
//	config := middleware.CORSConfig{
//	    AllowedOrigins:   []string{"https://example.com"},
//	    AllowCredentials: true,
//	}
//	router.Use(middleware.CORSMiddleware(config))
//
// 注意事项：
// - 如果 AllowCredentials 为 true，AllowedOrigins 不能包含 "*"
// - 应该在中间件链的早期使用（在其他中间件之前）
// - 生产环境应该明确指定允许的源，避免使用 "*"
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	// 默认配置
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{"Content-Type", "Authorization", "X-Requested-With"}
	}
	if config.MaxAge == 0 {
		config.MaxAge = 3600
	}

	allowedMethodsStr := strings.Join(config.AllowedMethods, ", ")
	allowedHeadersStr := strings.Join(config.AllowedHeaders, ", ")
	exposedHeadersStr := strings.Join(config.ExposedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// 检查源是否允许
			if len(config.AllowedOrigins) > 0 {
				allowed := false
				for _, allowedOrigin := range config.AllowedOrigins {
					if allowedOrigin == "*" || allowedOrigin == origin {
						allowed = true
						break
					}
				}
				if !allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			// 设置CORS头
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(config.AllowedOrigins) > 0 && config.AllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", allowedMethodsStr)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeadersStr)

			if len(config.ExposedHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", exposedHeadersStr)
			}

			w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
