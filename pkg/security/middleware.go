package security

import (
	"context"
	"net/http"
)

// SecurityMiddleware 安全中间件集合
type SecurityMiddleware struct {
	csrfProtection *CSRFProtection
	rateLimiter    *RateLimiter
	securityHeaders *SecurityHeaders
	xssProtection  *XSSProtection
}

// SecurityMiddlewareConfig 安全中间件配置
type SecurityMiddlewareConfig struct {
	CSRF           *CSRFConfig
	RateLimit      *RateLimiterConfig
	SecurityHeaders *SecurityHeadersConfig
	EnableXSS      bool
}

// NewSecurityMiddleware 创建安全中间件
func NewSecurityMiddleware(config SecurityMiddlewareConfig) *SecurityMiddleware {
	sm := &SecurityMiddleware{}

	if config.CSRF != nil {
		sm.csrfProtection = NewCSRFProtection(*config.CSRF)
	}

	if config.RateLimit != nil {
		sm.rateLimiter = NewRateLimiter(*config.RateLimit)
	}

	if config.SecurityHeaders != nil {
		sm.securityHeaders = NewSecurityHeaders(*config.SecurityHeaders)
	}

	if config.EnableXSS {
		sm.xssProtection = NewXSSProtection()
	}

	return sm
}

// Middleware 返回组合的安全中间件
func (sm *SecurityMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 安全头部
		if sm.securityHeaders != nil {
			sm.securityHeaders.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 继续处理
			})).ServeHTTP(w, r)
		}

		// 2. 速率限制
		if sm.rateLimiter != nil {
			key := sm.getRateLimitKey(r)
			allowed, err := sm.rateLimiter.Allow(r.Context(), key)
			if err != nil || !allowed {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		}

		// 3. CSRF 保护（仅 POST/PUT/DELETE/PATCH）
		if sm.csrfProtection != nil && sm.requiresCSRF(r.Method) {
			sessionID := sm.getSessionID(r)
			token := r.Header.Get("X-CSRF-Token")
			if token == "" {
				token = r.FormValue("csrf_token")
			}

			if err := sm.csrfProtection.ValidateToken(sessionID, token); err != nil {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		// 4. XSS 保护（清理输入）
		if sm.xssProtection != nil {
			// 清理查询参数
			for key, values := range r.URL.Query() {
				for i, value := range values {
					values[i] = sm.xssProtection.Sanitize(value)
				}
				r.URL.Query()[key] = values
			}
		}

		next.ServeHTTP(w, r)
	})
}

// getRateLimitKey 获取速率限制键
func (sm *SecurityMiddleware) getRateLimitKey(r *http.Request) string {
	// 优先使用 IP
	ip := sm.getClientIP(r)
	if ip != "" {
		return "ip:" + ip
	}

	// 其次使用用户 ID（如果已认证）
	userID := sm.getUserID(r)
	if userID != "" {
		return "user:" + userID
	}

	// 最后使用端点
	return "endpoint:" + r.URL.Path
}

// getClientIP 获取客户端 IP
func (sm *SecurityMiddleware) getClientIP(r *http.Request) string {
	// 检查 X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// 检查 X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 使用 RemoteAddr
	return r.RemoteAddr
}

// getSessionID 获取会话 ID
func (sm *SecurityMiddleware) getSessionID(r *http.Request) string {
	// 从 Cookie 获取
	cookie, err := r.Cookie("session_id")
	if err == nil {
		return cookie.Value
	}

	// 从 Header 获取
	return r.Header.Get("X-Session-ID")
}

// getUserID 获取用户 ID
func (sm *SecurityMiddleware) getUserID(r *http.Request) string {
	// 从上下文或 Header 获取（需要根据实际实现调整）
	return r.Header.Get("X-User-ID")
}

// requiresCSRF 检查是否需要 CSRF 保护
func (sm *SecurityMiddleware) requiresCSRF(method string) bool {
	csrfMethods := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, m := range csrfMethods {
		if method == m {
			return true
		}
	}
	return false
}

// CSRFMiddleware 仅 CSRF 中间件
func (sm *SecurityMiddleware) CSRFMiddleware(next http.Handler) http.Handler {
	if sm.csrfProtection == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sm.requiresCSRF(r.Method) {
			sessionID := sm.getSessionID(r)
			token := r.Header.Get("X-CSRF-Token")
			if token == "" {
				token = r.FormValue("csrf_token")
			}

			if err := sm.csrfProtection.ValidateToken(sessionID, token); err != nil {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware 仅速率限制中间件
func (sm *SecurityMiddleware) RateLimitMiddleware(next http.Handler) http.Handler {
	if sm.rateLimiter == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := sm.getRateLimitKey(r)
		allowed, err := sm.rateLimiter.Allow(r.Context(), key)
		if err != nil || !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware 仅安全头部中间件
func (sm *SecurityMiddleware) SecurityHeadersMiddleware(next http.Handler) http.Handler {
	if sm.securityHeaders == nil {
		return next
	}

	return sm.securityHeaders.Middleware(next)
}

// Shutdown 关闭安全中间件
func (sm *SecurityMiddleware) Shutdown() error {
	if sm.csrfProtection != nil {
		sm.csrfProtection.Shutdown()
	}

	if sm.rateLimiter != nil {
		sm.rateLimiter.Shutdown(context.Background())
	}

	return nil
}
