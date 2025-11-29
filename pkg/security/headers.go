package security

import (
	"net/http"
)

// SecurityHeaders 安全头部中间件
type SecurityHeaders struct {
	config SecurityHeadersConfig
}

// SecurityHeadersConfig 安全头部配置
type SecurityHeadersConfig struct {
	// Content-Security-Policy
	CSP string

	// X-Content-Type-Options
	ContentTypeOptions string

	// X-Frame-Options
	FrameOptions string

	// X-XSS-Protection
	XSSProtection string

	// Strict-Transport-Security
	HSTS string

	// Referrer-Policy
	ReferrerPolicy string

	// Permissions-Policy
	PermissionsPolicy string
}

// DefaultSecurityHeadersConfig 默认安全头部配置
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		CSP:                "default-src 'self'",
		ContentTypeOptions: "nosniff",
		FrameOptions:       "DENY",
		XSSProtection:      "1; mode=block",
		HSTS:               "max-age=31536000; includeSubDomains",
		ReferrerPolicy:     "strict-origin-when-cross-origin",
		PermissionsPolicy:  "geolocation=(), microphone=(), camera=()",
	}
}

// NewSecurityHeaders 创建安全头部中间件
func NewSecurityHeaders(config SecurityHeadersConfig) *SecurityHeaders {
	if config.CSP == "" {
		config = DefaultSecurityHeadersConfig()
	}

	return &SecurityHeaders{
		config: config,
	}
}

// Middleware 返回 HTTP 中间件函数
func (s *SecurityHeaders) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置安全头部
		if s.config.CSP != "" {
			w.Header().Set("Content-Security-Policy", s.config.CSP)
		}

		if s.config.ContentTypeOptions != "" {
			w.Header().Set("X-Content-Type-Options", s.config.ContentTypeOptions)
		}

		if s.config.FrameOptions != "" {
			w.Header().Set("X-Frame-Options", s.config.FrameOptions)
		}

		if s.config.XSSProtection != "" {
			w.Header().Set("X-XSS-Protection", s.config.XSSProtection)
		}

		if s.config.HSTS != "" && r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", s.config.HSTS)
		}

		if s.config.ReferrerPolicy != "" {
			w.Header().Set("Referrer-Policy", s.config.ReferrerPolicy)
		}

		if s.config.PermissionsPolicy != "" {
			w.Header().Set("Permissions-Policy", s.config.PermissionsPolicy)
		}

		// 移除服务器信息
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")

		next.ServeHTTP(w, r)
	})
}

// HandlerFunc 返回 HTTP 处理函数
func (s *SecurityHeaders) HandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return s.Middleware(http.HandlerFunc(fn)).ServeHTTP
}
