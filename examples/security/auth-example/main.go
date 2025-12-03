package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/yourusername/golang/pkg/security/jwt"
	"github.com/yourusername/golang/pkg/security/rbac"
)

func main() {
	log.Println("🔐 安全认证授权示例")
	log.Println("展示 JWT + RBAC 完整集成")
	log.Println("")

	// 1. 创建 JWT Token Manager
	log.Println("📝 创建 JWT Token Manager...")
	tokenManager, err := jwt.NewTokenManager(jwt.Config{
		Issuer:          "auth-example",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		SigningMethod:   "RS256",
	})
	if err != nil {
		log.Fatal("Failed to create token manager:", err)
	}

	// 2. 创建 RBAC 系统
	log.Println("🔒 创建 RBAC 系统...")
	rbacSystem := rbac.NewRBAC()

	// 初始化默认角色
	if err := rbacSystem.InitializeDefaultRoles(); err != nil {
		log.Fatal("Failed to initialize RBAC:", err)
	}

	// 3. 创建中间件
	jwtMiddleware := jwt.NewMiddleware(jwt.MiddlewareConfig{
		TokenManager: tokenManager,
		SkipPaths:    []string{"/login", "/health"},
	})

	rbacMiddleware := rbac.NewMiddleware(rbacSystem)

	// 4. 创建路由
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 公开端点
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 登录端点
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		// 简化示例：直接生成令牌
		// 实际应该验证用户名密码

		tokenPair, err := tokenManager.GenerateTokenPair(
			"user-123",
			"john.doe",
			"john@example.com",
			[]string{"user", "moderator"},
		)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// 返回令牌
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"access_token": "%s",
			"refresh_token": "%s",
			"token_type": "Bearer",
			"expires_in": %d
		}`, tokenPair.AccessToken, tokenPair.RefreshToken, tokenPair.ExpiresIn)
	})

	// 需要认证的端点组
	r.Group(func(r chi.Router) {
		// 应用 JWT 认证
		r.Use(jwtMiddleware.Authenticate)

		// 所有认证用户都可以访问
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			claims, _ := jwt.GetClaims(r.Context())
			w.Write([]byte("Welcome, " + claims.Username))
		})

		// 需要特定权限
		r.Group(func(r chi.Router) {
			r.Use(rbacMiddleware.RequirePermission("user", "read"))
			r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("User list"))
			})
		})

		// 需要特定角色
		r.Group(func(r chi.Router) {
			r.Use(rbacMiddleware.RequireRole("admin"))
			r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("User created"))
			})
		})
	})

	// 启动服务器
	addr := ":8080"
	log.Printf("🚀 服务器启动在 %s\n", addr)
	log.Println("")
	log.Println("📖 API 端点:")
	log.Println("  POST /login           - 登录获取令牌（公开）")
	log.Println("  GET  /health          - 健康检查（公开）")
	log.Println("  GET  /profile         - 用户资料（需要认证）")
	log.Println("  GET  /users           - 用户列表（需要 user:read 权限）")
	log.Println("  POST /users           - 创建用户（需要 admin 角色）")
	log.Println("")
	log.Println("💡 测试命令:")
	log.Println("  # 1. 登录获取令牌")
	log.Println("  curl -X POST http://localhost:8080/login")
	log.Println("")
	log.Println("  # 2. 使用令牌访问")
	log.Println("  curl -H 'Authorization: Bearer <token>' http://localhost:8080/profile")
	log.Println("")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

// 需要导入 fmt
import "fmt"
