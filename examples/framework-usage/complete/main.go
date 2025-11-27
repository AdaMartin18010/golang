// Package main 展示如何完整使用框架的所有核心能力
//
// 本示例展示：
// 1. 如何初始化框架的各种能力（数据库、可观测性、精细控制等）
// 2. 如何配置 HTTP 中间件
// 3. 如何在实际业务中使用这些能力
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
	chiMiddleware "github.com/yourusername/golang/internal/interfaces/http/chi/middleware"
	"github.com/yourusername/golang/pkg/control"
	"github.com/yourusername/golang/pkg/database"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/sampling"
	"github.com/yourusername/golang/pkg/tracing"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. 初始化可观测性（OTLP）
	log.Println("Initializing observability...")
	sampler, err := sampling.NewProbabilisticSampler(0.5)
	if err != nil {
		log.Fatalf("Failed to create sampler: %v", err)
	}

	otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:    "example-service",
		ServiceVersion: "v1.0.0",
		Endpoint:       "localhost:4317",
		Insecure:       true,
		Sampler:        sampler,
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize OTLP: %v (continuing without OTLP)", err)
	} else {
		defer otlpClient.Shutdown(ctx)
		log.Println("OTLP initialized successfully")
	}

	// 2. 初始化追踪器
	tracer := tracing.NewTracer("example-service")
	log.Println("Tracer initialized")

	// 3. 初始化数据库
	log.Println("Initializing database...")
	db, err := database.NewDatabase(database.Config{
		Driver:       database.DriverSQLite3,
		DSN:          "file:example.db?cache=shared&mode=memory",
		MaxOpenConns: 25,
		MaxIdleConns: 5,
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	// 4. 初始化精细控制
	log.Println("Initializing control mechanisms...")
	featureController := control.NewFeatureController()
	rateController := control.NewRateController()
	circuitController := control.NewCircuitController()

	// 注册功能开关
	featureController.Register("experimental-feature", "Experimental feature for testing", true, map[string]interface{}{
		"max_requests": 100,
	})

	// 注册速率限制
	rateController.SetRateLimit("api-calls", 100.0, time.Second)

	// 注册熔断器
	circuitController.RegisterCircuit("external-api", 10, 5, 30*time.Second)

	log.Println("Control mechanisms initialized")

	// 5. 创建 HTTP 路由器
	r := chi.NewRouter()

	// 6. 配置中间件
	setupMiddleware(r, sampler, featureController, rateController, circuitController)

	// 7. 配置路由
	setupRoutes(r, db, tracer, featureController)

	// 8. 启动服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 优雅关闭
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// setupMiddleware 配置所有中间件
func setupMiddleware(
	r *chi.Mux,
	sampler sampling.Sampler,
	featureController *control.FeatureController,
	rateController *control.RateController,
	circuitController *control.CircuitController,
) {
	// 基础中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// 采样中间件
	r.Use(chiMiddleware.SamplingMiddleware(chiMiddleware.SamplingConfig{
		Sampler:             sampler,
		SkipPaths:           []string{"/health", "/metrics"},
		AddSamplingDecision: true,
	}))

	// 追踪中间件
	r.Use(chiMiddleware.TracingMiddleware(chiMiddleware.TracingConfig{
		ServiceName:    "example-service",
		ServiceVersion: "v1.0.0",
		SkipPaths:      []string{"/health", "/metrics"},
		AddRequestID:   true,
	}))

	// 反射中间件
	r.Use(chiMiddleware.ReflectMiddleware(chiMiddleware.ReflectConfig{
		EnableMetadata:     true,
		EnableSelfDescribe: true,
		SkipPaths:          []string{"/health", "/metrics"},
	}))

	// 数据转换中间件
	r.Use(chiMiddleware.ConverterMiddleware(chiMiddleware.ConverterConfig{
		EnableRequestConversion:  true,
		EnableResponseConversion: true,
		DefaultResponseFormat:    "json",
	}))

	// 精细控制中间件
	r.Use(chiMiddleware.ControlMiddleware(chiMiddleware.ControlConfig{
		FeatureController: featureController,
		RateController:    rateController,
		CircuitController: circuitController,
		FeatureFlags: map[string]string{
			"/api/v1/experimental": "experimental-feature",
		},
		RateLimits: map[string]string{
			"/api/v1/users": "api-calls",
		},
		CircuitBreakers: map[string]string{
			"/api/v1/external": "external-api",
		},
		SkipPaths: []string{"/health", "/metrics"},
	}))

	// 其他中间件
	r.Use(chiMiddleware.LoggingMiddleware)
	r.Use(chiMiddleware.RecoveryMiddleware)
	r.Use(chiMiddleware.TimeoutMiddleware(60 * time.Second))
	r.Use(chiMiddleware.CORSMiddleware)
}

// setupRoutes 配置路由
func setupRoutes(
	r *chi.Mux,
	db database.Database,
	tracer *tracing.Tracer,
	featureController *control.FeatureController,
) {
	// 健康检查
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API 路由
	r.Route("/api/v1", func(r chi.Router) {
		// 示例：用户列表
		r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 开始追踪
			ctx, span := tracer.StartSpan(ctx, "get-users")
			defer span.End()

			// 检查采样决策
			if chiMiddleware.IsSampled(ctx) {
				log.Printf("Request sampled: %s", r.URL.Path)
			}

			// 使用数据转换
			data := chiMiddleware.GetRequestData(ctx)
			if data != nil {
				log.Printf("Request data: %v", data)
			}

			// 使用反射检查器
			inspector := chiMiddleware.GetInspector(ctx)
			if inspector != nil {
				type User struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}
				user := User{ID: 1, Name: "Test"}
				metadata := inspector.InspectType(user)
				log.Printf("Type metadata: %+v", metadata)
			}

			// 执行数据库操作
			rows, err := db.Query(ctx, "SELECT id, name FROM users LIMIT 10")
			if err != nil {
				tracer.LocateError(ctx, err, map[string]interface{}{
					"endpoint": "/api/v1/users",
				})
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// 处理结果
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"users": []}`))
		})

		// 示例：实验性功能
		r.Get("/experimental", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 检查功能开关
			if !chiMiddleware.GetFeatureFlag(ctx, "experimental-feature") {
				http.Error(w, "Feature is disabled", http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Experimental feature is enabled"}`))
		})
	})
}
