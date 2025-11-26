package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yourusername/golang/pkg/auth/jwt"
	"github.com/yourusername/golang/pkg/errors"
	"github.com/yourusername/golang/pkg/http/response"
)

// AuthConfig 认证配置
type AuthConfig struct {
	JWT          *jwt.JWT
	SkipPaths    []string // 跳过认证的路径
	OptionalAuth bool     // 是否可选认证（不强制）
}

// AuthMiddleware 认证中间件
func AuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否跳过认证
			if shouldSkipAuth(r.URL.Path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// 提取Token
			tokenString := extractToken(r)
			if tokenString == "" {
				if config.OptionalAuth {
					next.ServeHTTP(w, r)
					return
				}
				response.Error(w, http.StatusUnauthorized,
					errors.NewUnauthorizedError("missing or invalid authorization token"))
				return
			}

			// 验证Token
			claims, err := config.JWT.ValidateToken(tokenString)
			if err != nil {
				if config.OptionalAuth {
					next.ServeHTTP(w, r)
					return
				}
				if err == jwt.ErrExpiredToken {
					response.Error(w, http.StatusUnauthorized,
						errors.NewUnauthorizedError("token expired"))
					return
				}
				response.Error(w, http.StatusUnauthorized,
					errors.NewUnauthorizedError("invalid token"))
				return
			}

			// 将用户信息添加到context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "roles", claims.Roles)
			ctx = context.WithValue(ctx, "email", claims.Email)
			ctx = context.WithValue(ctx, "claims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole 要求特定角色
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value("roles").([]string)
			if !ok {
				response.Error(w, http.StatusForbidden,
					errors.NewForbiddenError("user roles not found"))
				return
			}

			// 检查用户是否有任一所需角色
			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range userRoles {
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
				response.Error(w, http.StatusForbidden,
					errors.NewForbiddenError("insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole 要求任一角色
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return RequireRole(roles...)
}

// RequireAllRoles 要求所有角色
func RequireAllRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value("roles").([]string)
			if !ok {
				response.Error(w, http.StatusForbidden,
					errors.NewForbiddenError("user roles not found"))
				return
			}

			// 检查用户是否有所有所需角色
			roleMap := make(map[string]bool)
			for _, role := range userRoles {
				roleMap[role] = true
			}

			for _, requiredRole := range roles {
				if !roleMap[requiredRole] {
					response.Error(w, http.StatusForbidden,
						errors.NewForbiddenError("insufficient permissions"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractToken 从请求中提取Token
func extractToken(r *http.Request) string {
	// 1. 从 Authorization header 提取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// 2. 从 Cookie 提取
	cookie, err := r.Cookie("access_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// 3. 从查询参数提取
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
}

// shouldSkipAuth 检查是否应该跳过认证
func shouldSkipAuth(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// GetUserID 从context中获取UserID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetUsername 从context中获取Username
func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value("username").(string); ok {
		return username
	}
	return ""
}

// GetRoles 从context中获取Roles
func GetRoles(ctx context.Context) []string {
	if roles, ok := ctx.Value("roles").([]string); ok {
		return roles
	}
	return nil
}

// GetEmail 从context中获取Email
func GetEmail(ctx context.Context) string {
	if email, ok := ctx.Value("email").(string); ok {
		return email
	}
	return ""
}

// GetClaims 从context中获取Claims
func GetClaims(ctx context.Context) *jwt.Claims {
	if claims, ok := ctx.Value("claims").(*jwt.Claims); ok {
		return claims
	}
	return nil
}
