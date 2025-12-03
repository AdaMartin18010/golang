package jwt

import (
	"context"
	"net/http"
	"strings"
)

// ContextKey JWT 上下文键
type ContextKey string

const (
	// ClaimsKey JWT Claims 上下文键
	ClaimsKey ContextKey = "jwt_claims"
)

// Middleware JWT 认证中间件
type Middleware struct {
	tokenManager *TokenManager
	skipPaths    map[string]bool
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	TokenManager *TokenManager
	SkipPaths    []string // 跳过认证的路径
}

// NewMiddleware 创建 JWT 中间件
func NewMiddleware(cfg MiddlewareConfig) *Middleware {
	skipPaths := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skipPaths[path] = true
	}

	return &Middleware{
		tokenManager: cfg.TokenManager,
		skipPaths:    skipPaths,
	}
}

// Authenticate JWT 认证中间件
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查是否跳过认证
		if m.skipPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// 从 Authorization header 提取令牌
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing authorization header", http.StatusUnauthorized)
			return
		}

		// 检查 Bearer scheme
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 验证令牌
		claims, err := m.tokenManager.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// 将 claims 添加到上下文
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRoles 要求特定角色的中间件
func (m *Middleware) RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从上下文获取 claims
			claims, ok := r.Context().Value(ClaimsKey).(*Claims)
			if !ok {
				http.Error(w, "Unauthorized: no claims found", http.StatusUnauthorized)
				return
			}

			// 检查角色
			hasRole := false
			for _, userRole := range claims.Roles {
				for _, requiredRole := range roles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: required role not found", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClaims 从上下文获取 JWT Claims
func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*Claims)
	return claims, ok
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) (string, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return "", false
	}
	return claims.UserID, true
}

// GetUserRoles 从上下文获取用户角色
func GetUserRoles(ctx context.Context) ([]string, bool) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return nil, false
	}
	return claims.Roles, true
}
