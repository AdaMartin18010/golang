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
	"github.com/yourusername/golang/pkg/control"
	"github.com/yourusername/golang/pkg/database"
	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/sampling"
	"github.com/yourusername/golang/pkg/observability/tracing"
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
	// 注意：NewFeatureController 返回 Controller 接口，使用类型断言获取具体类型
	controller := control.NewFeatureController()
	featureController, ok := controller.(*control.FeatureController)
	if !ok {
		log.Fatal("Failed to get FeatureController")
	}

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

	// 6. 配置基础中间件
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	// CORS 中间件使用 chi/cors 包或自定义实现
	// r.Use(cors.Handler(cors.Options{...}))

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
			// 检查功能开关
			if !featureController.IsEnabled("experimental-feature") {
				http.Error(w, "Feature is disabled", http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Experimental feature is enabled"}`))
		})
	})
}

