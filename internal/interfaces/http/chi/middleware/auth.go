package middleware

import (
	"net/http"

	"github.com/yourusername/golang/pkg/security/jwt"
	"github.com/yourusername/golang/pkg/security/rbac"
)

// AuthMiddleware 认证中间件配置
type AuthMiddleware struct {
	jwtMiddleware  *jwt.Middleware
	rbacMiddleware *rbac.Middleware
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtMw *jwt.Middleware, rbacMw *rbac.Middleware) *AuthMiddleware {
	return &AuthMiddleware{
		jwtMiddleware:  jwtMw,
		rbacMiddleware: rbacMw,
	}
}

// Authenticate JWT 认证
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return am.jwtMiddleware.Authenticate(next)
}

// RequirePermission 要求特定权限
func (am *AuthMiddleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 先认证
			am.jwtMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 再授权
				am.rbacMiddleware.RequirePermission(resource, action)(next).ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}

// RequireRole 要求特定角色
func (am *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return am.jwtMiddleware.RequireRoles(roles...)
}

// 辅助函数：提取 JWT Claims 和转换为 RBAC 上下文
func convertJWTClaimsToRBACContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从 JWT claims 获取角色
		claims, ok := jwt.GetClaims(r.Context())
		if ok && len(claims.Roles) > 0 {
			// 将角色添加到 RBAC 上下文
			ctx := rbac.WithUserRoles(r.Context(), claims.Roles)
			ctx = rbac.WithUserID(ctx, claims.UserID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
