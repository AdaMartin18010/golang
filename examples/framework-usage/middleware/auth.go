package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/security/jwt"
)

// contextKey 是上下文键的类型
type contextKey string

const (
	userIDKey contextKey = "user_id"
	rolesKey  contextKey = "roles"
)

// AuthConfig 认证中间件配置
type AuthConfig struct {
	TokenManager *jwt.TokenManager
	SkipPaths    []string
}

// AuthMiddleware 创建JWT认证中间件
func AuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
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

			// 获取Authorization头
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			// 解析Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// 验证token
			claims, err := config.TokenManager.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			// 检查token是否过期
			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
				http.Error(w, `{"error":"token expired"}`, http.StatusUnauthorized)
				return
			}

			// 将用户信息存入上下文
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, rolesKey, claims.Roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetRoles 从上下文获取用户角色
func GetRoles(ctx context.Context) []string {
	if roles, ok := ctx.Value(rolesKey).([]string); ok {
		return roles
	}
	return nil
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx context.Context) (*jwt.Claims, error) {
	userID := GetUserID(ctx)
	if userID == "" {
		return nil, errors.NewUnauthorizedError("user not found in context")
	}

	return &jwt.Claims{
		UserID: userID,
		Roles:  GetRoles(ctx),
	}, nil
}
