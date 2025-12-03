package rbac

import (
	"context"
	"net/http"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	// UserRolesKey 用户角色上下文键
	UserRolesKey ContextKey = "user_roles"
	// UserIDKey 用户ID上下文键
	UserIDKey ContextKey = "user_id"
)

// Middleware RBAC 中间件
type Middleware struct {
	rbac *RBAC
}

// NewMiddleware 创建 RBAC 中间件
func NewMiddleware(rbac *RBAC) *Middleware {
	return &Middleware{rbac: rbac}
}

// RequirePermission 要求特定权限的中间件
func (m *Middleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从上下文获取用户角色
			userRoles, ok := r.Context().Value(UserRolesKey).([]string)
			if !ok || len(userRoles) == 0 {
				http.Error(w, "Unauthorized: no roles found", http.StatusUnauthorized)
				return
			}

			// 检查权限
			hasPermission, err := m.rbac.CheckPermission(r.Context(), userRoles, resource, action)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			// 权限检查通过，继续处理
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole 要求特定角色的中间件
func (m *Middleware) RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从上下文获取用户角色
			userRoles, ok := r.Context().Value(UserRolesKey).([]string)
			if !ok || len(userRoles) == 0 {
				http.Error(w, "Unauthorized: no roles found", http.StatusUnauthorized)
				return
			}

			// 检查是否有任一所需角色
			hasRole := false
			for _, userRole := range userRoles {
				for _, requiredRole := range requiredRoles {
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

			// 角色检查通过，继续处理
			next.ServeHTTP(w, r)
		})
	}
}

// WithUserRoles 将用户角色添加到上下文
func WithUserRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, UserRolesKey, roles)
}

// GetUserRoles 从上下文获取用户角色
func GetUserRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(UserRolesKey).([]string)
	return roles, ok
}

// WithUserID 将用户ID添加到上下文
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
