package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string // 允许的源
	AllowedMethods   []string // 允许的方法
	AllowedHeaders   []string // 允许的请求头
	ExposedHeaders   []string // 暴露的响应头
	AllowCredentials bool     // 是否允许凭证
	MaxAge           int      // 预检请求缓存时间（秒）
}

// CORSMiddleware CORS中间件
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
