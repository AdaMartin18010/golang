package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// 可观测性
	"github.com/yourusername/golang/pkg/observability/ebpf"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/observability/system"

	// 安全
	"github.com/yourusername/golang/pkg/security/jwt"
	"github.com/yourusername/golang/pkg/security/rbac"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	log.Println("🚀 完整集成示例 - Go Clean Architecture 框架")
	log.Println("展示所有核心功能的集成使用")
	log.Println("")

	ctx := context.Background()

	// ============================================
	// 1. 初始化可观测性 (OTLP)
	// ============================================
	log.Println("📊 初始化 OpenTelemetry...")
	otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:    "complete-integration-example",
		ServiceVersion: "1.0.0",
		Endpoint:       "localhost:4317",
		Insecure:       true,
		SampleRate:     1.0, // 演示用，100% 采样
	})
	if err != nil {
		log.Printf("⚠️  OTLP 初始化失败: %v (继续运行)", err)
	} else {
		log.Println("✅ OTLP 初始化成功")
		defer otlpClient.Shutdown(ctx)
	}

	// ============================================
	// 2. 初始化系统监控
	// ============================================
	log.Println("📈 初始化系统监控...")
	meter := otel.Meter("complete-example")

	systemMonitor, err := system.NewMonitor(system.Config{
		Meter:           meter,
		CollectInterval: 5 * time.Second,
	})
	if err != nil {
		log.Printf("⚠️  系统监控初始化失败: %v", err)
	} else {
		log.Println("✅ 系统监控初始化成功")
		go systemMonitor.Start()
		defer systemMonitor.Stop()
	}

	// ============================================
	// 3. 初始化平台检测
	// ============================================
	log.Println("🖥️  检测运行环境...")
	platformMonitor, err := system.NewPlatformMonitor(meter)
	if err != nil {
		log.Printf("⚠️  平台检测失败: %v", err)
	} else {
		info := platformMonitor.GetInfo()
		log.Printf("✅ 运行环境检测:")
		log.Printf("   - OS: %s (%s)", info.OS, info.Arch)
		log.Printf("   - 容器: %s", getContainerStatus(info))
		log.Printf("   - Kubernetes: %s", getK8sStatus(info))
		log.Printf("   - 云厂商: %s", getCloudStatus(info))
	}

	// ============================================
	// 4. 初始化 eBPF 监控 (可选，需要 Linux)
	// ============================================
	log.Println("🔍 初始化 eBPF 监控...")
	tracer := otel.Tracer("complete-example")

	ebpfCollector, err := ebpf.NewCollector(ebpf.Config{
		Tracer:                  tracer,
		Meter:                   meter,
		Enabled:                 true,
		EnableSyscallTracking:   true,
		EnableNetworkMonitoring: true,
		CollectInterval:         5 * time.Second,
	})
	if err != nil {
		log.Printf("⚠️  eBPF 初始化失败: %v (需要 Linux + Root)", err)
	} else {
		if err := ebpfCollector.Start(); err != nil {
			log.Printf("⚠️  eBPF 启动失败: %v (需要 Linux + Root)", err)
		} else {
			log.Println("✅ eBPF 监控启动成功")
			defer ebpfCollector.Stop()
		}
	}

	// ============================================
	// 5. 初始化安全模块 (JWT + RBAC)
	// ============================================
	log.Println("🔐 初始化安全模块...")

	// JWT Token Manager
	tokenManager, err := jwt.NewTokenManager(jwt.Config{
		Issuer:          "complete-example",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		SigningMethod:   "RS256",
	})
	if err != nil {
		log.Fatalf("❌ JWT 初始化失败: %v", err)
	}
	log.Println("✅ JWT Token Manager 初始化成功")

	// RBAC 系统
	rbacSystem := rbac.NewRBAC()
	if err := rbacSystem.InitializeDefaultRoles(); err != nil {
		log.Fatalf("❌ RBAC 初始化失败: %v", err)
	}
	log.Println("✅ RBAC 系统初始化成功")

	// 创建中间件
	jwtMiddleware := jwt.NewMiddleware(jwt.MiddlewareConfig{
		TokenManager: tokenManager,
		SkipPaths:    []string{"/", "/health", "/login", "/metrics"},
	})
	rbacMiddleware := rbac.NewMiddleware(rbacSystem)

	// ============================================
	// 6. 创建 HTTP 服务
	// ============================================
	log.Println("🌐 创建 HTTP 服务...")
	r := chi.NewRouter()

	// 基础中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// 公开端点
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Complete Integration Example!"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// 登录端点
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "login")
		defer span.End()

		// 简化示例：直接生成令牌
		tokenPair, err := tokenManager.GenerateTokenPair(
			"user-123",
			"john.doe",
			"john@example.com",
			[]string{"user", "moderator"},
		)
		if err != nil {
			span.RecordError(err)
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		span.SetAttributes(attribute.String("user.id", "user-123"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token":"` + tokenPair.AccessToken + `","expires_in":` + formatInt64(tokenPair.ExpiresIn) + `}`))
	})

	// 需要认证的端点组
	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware.Authenticate)

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			_, span := tracer.Start(r.Context(), "get-profile")
			defer span.End()

			claims, _ := jwt.GetClaims(r.Context())
			span.SetAttributes(attribute.String("user.id", claims.UserID))

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"user_id":"` + claims.UserID + `","username":"` + claims.Username + `"}`))
		})

		// 需要权限的端点
		r.Group(func(r chi.Router) {
			r.Use(rbacMiddleware.RequirePermission("user", "read"))

			r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
				_, span := tracer.Start(r.Context(), "list-users")
				defer span.End()

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"users":[],"total":0}`))
			})
		})

		// 需要管理员角色
		r.Group(func(r chi.Router) {
			r.Use(rbacMiddleware.RequireRole("admin"))

			r.Post("/admin/users", func(w http.ResponseWriter, r *http.Request) {
				_, span := tracer.Start(r.Context(), "admin-create-user")
				defer span.End()

				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"message":"User created"}`))
			})
		})
	})

	// ============================================
	// 7. 启动服务器
	// ============================================
	addr := ":8080"
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("🚀 服务器启动在 %s\n", addr)
		log.Println("")
		log.Println("📖 API 端点:")
		log.Println("  GET  /              - 欢迎页面（公开）")
		log.Println("  GET  /health        - 健康检查（公开）")
		log.Println("  POST /login         - 登录获取令牌（公开）")
		log.Println("  GET  /profile       - 用户资料（需要认证）")
		log.Println("  GET  /users         - 用户列表（需要 user:read 权限）")
		log.Println("  POST /admin/users   - 创建用户（需要 admin 角色）")
		log.Println("")
		log.Println("💡 测试命令:")
		log.Println("  # 1. 登录")
		log.Println("  curl -X POST http://localhost:8080/login")
		log.Println("")
		log.Println("  # 2. 访问资料")
		log.Println("  curl -H 'Authorization: Bearer <token>' http://localhost:8080/profile")
		log.Println("")
		log.Println("📊 可观测性:")
		log.Println("  - OpenTelemetry 追踪已启用")
		log.Println("  - 系统监控每 5 秒采集")
		if ebpfCollector.IsEnabled() {
			log.Println("  - eBPF 监控已启用")
		}
		log.Println("")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// ============================================
	// 8. 优雅关闭
	// ============================================
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("")
	log.Println("🛑 收到停止信号，正在优雅关闭...")

	// 创建关闭上下文
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// 停止 HTTP 服务器
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  HTTP 服务器关闭错误: %v", err)
	}

	log.Println("✅ 应用已优雅关闭")
}

// 辅助函数

func formatInt64(n int64) string {
	return string(rune('0' + n%10))
}

func getContainerStatus(info system.PlatformInfo) string {
	if info.ContainerID != "" {
		return "是 (ID: " + info.ContainerID[:12] + ")"
	}
	return "否"
}

func getK8sStatus(info system.PlatformInfo) string {
	if info.KubernetesPod != "" {
		return "是 (Pod: " + info.KubernetesPod + ", Node: " + info.KubernetesNode + ")"
	}
	return "否"
}

func getCloudStatus(info system.PlatformInfo) string {
	if info.CloudProvider != "" {
		region := info.CloudRegion
		if region == "" {
			region = "unknown"
		}
		return info.CloudProvider + " (" + region + ")"
	}
	return "本地/裸机"
}
