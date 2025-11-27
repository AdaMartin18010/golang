// Package main HTTP 服务器应用入口
//
// 设计原理：
// 1. 这是 HTTP REST API 服务器的入口点
// 2. 遵循 Clean Architecture 架构，使用依赖注入组装组件
// 3. 实现优雅关闭，确保请求处理完成后再退出
//
// 启动流程：
// 1. 加载配置
// 2. 初始化日志
// 3. 初始化 OpenTelemetry Tracer
// 4. 初始化数据库连接
// 5. 运行数据库迁移
// 6. 组装依赖（仓储、服务、路由）
// 7. 初始化 Temporal 客户端（可选）
// 8. 创建 HTTP 服务器
// 9. 启动服务器（在 goroutine 中）
// 10. 等待中断信号
// 11. 优雅关闭
//
// 依赖注入：
// 当前使用手动依赖注入，推荐迁移到 Wire 依赖注入
// 详见：scripts/wire/wire.go 和 docs/architecture/00-架构模型与依赖注入完整说明.md
//
// 优雅关闭：
// - 监听 SIGINT（Ctrl+C）和 SIGTERM（kill 命令）
// - 停止接收新请求
// - 等待正在处理的请求完成（最多 30 秒）
// - 关闭数据库连接和其他资源
//
// 使用方式：
//   go run cmd/server/main.go
//   # 或
//   go build -o bin/server cmd/server/main.go
//   ./bin/server
//
// 配置：
// - 配置文件：configs/config.yaml
// - 环境变量：支持通过环境变量覆盖配置
// - 默认端口：8080
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/internal/config"
	appuser "github.com/yourusername/golang/internal/application/user"
	entdb "github.com/yourusername/golang/internal/infrastructure/database/ent"
	entrepo "github.com/yourusername/golang/internal/infrastructure/database/ent/repository"
	"github.com/yourusername/golang/internal/infrastructure/observability/otlp"
	"github.com/yourusername/golang/internal/infrastructure/workflow/temporal"
	temporalhandler "github.com/yourusername/golang/internal/interfaces/workflow/temporal"
	chiRouter "github.com/yourusername/golang/internal/interfaces/http/chi"
)

func main() {
	// 步骤 1: 加载配置
	//
	// 配置加载说明：
	// - 从 configs/config.yaml 加载配置
	// - 支持环境变量覆盖（APP_* 前缀）
	// - 如果配置文件不存在，使用默认配置
	//
	// 配置优先级：
	// 1. 环境变量（最高）
	// 2. 配置文件
	// 3. 默认值（最低）
	//
	// 示例环境变量：
	//   APP_SERVER_PORT=8080
	//   APP_DATABASE_HOST=localhost
	//   APP_DATABASE_PORT=5432
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 步骤 2: 初始化日志
	//
	// 日志初始化说明：
	// - 使用结构化日志（slog）
	// - 支持 JSON 和 Text 格式
	// - 集成 OpenTelemetry，支持 TraceID 和 SpanID
	//
	// 日志级别：
	// - DEBUG: 调试信息
	// - INFO: 一般信息
	// - WARN: 警告信息
	// - ERROR: 错误信息
	logger := otlp.NewLogger()
	slog.SetDefault(logger.Logger)
	logger.Info("Application starting...")

	// 步骤 3: 初始化 OpenTelemetry Tracer
	//
	// OpenTelemetry 说明：
	// - 用于分布式追踪
	// - 支持 Jaeger、Zipkin 等后端
	// - 如果初始化失败，记录警告但继续运行
	//
	// 配置：
	// - TraceEndpoint: OpenTelemetry Collector 地址
	// - 如果为空，则不启用追踪
	ctx := context.Background()
	shutdownTracer, err := otlp.NewTracerProvider(ctx, cfg.Observability.TraceEndpoint, true)
	if err != nil {
		logger.Warn("Failed to initialize tracer", "error", err)
	} else {
		// 确保在程序退出时关闭 Tracer
		defer func() {
			if err := shutdownTracer(ctx); err != nil {
				logger.Error("Failed to shutdown tracer", "error", err)
			}
		}()
	}

	// 步骤 4: 初始化数据库（使用 Ent ORM）
	//
	// 数据库初始化说明：
	// - 使用 Ent ORM 连接 PostgreSQL
	// - 支持连接池管理
	// - 自动运行数据库迁移
	//
	// 数据库配置：
	// - Host: 数据库主机地址
	// - Port: 数据库端口（默认 5432）
	// - User: 数据库用户名
	// - Password: 数据库密码
	// - DBName: 数据库名称
	// - SSLMode: SSL 模式（disable、require、verify-full 等）
	//
	// 连接池配置：
	// - MaxOpenConns: 最大打开连接数（默认 25）
	// - MaxIdleConns: 最大空闲连接数（默认 5）
	var userService appuser.Service
	entClient, err := entdb.NewClientFromConfig(
		ctx,
		cfg.Database.Host,
		fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		logger.Info("Please ensure PostgreSQL is running and configuration is correct")
		os.Exit(1)
	}
	// 确保在程序退出时关闭数据库连接
	defer entClient.Close()
	logger.Info("Database connected successfully")

	// 步骤 4.1: 运行数据库迁移
	//
	// 数据库迁移说明：
	// - 自动创建或更新数据库表结构
	// - 如果迁移失败，记录警告但继续运行（可能表已存在）
	// - 生产环境建议在部署前手动运行迁移
	if err := entClient.Migrate(ctx); err != nil {
		logger.Warn("Failed to run database migrations", "error", err)
		// 继续运行，可能表已存在
	} else {
		logger.Info("Database migrations completed")
	}

	// 步骤 4.2: 组装依赖（手动依赖注入）
	//
	// 依赖注入说明：
	// - 当前使用手动依赖注入
	// - 推荐迁移到 Wire 依赖注入（详见 scripts/wire/wire.go）
	//
	// 依赖关系：
	//   EntClient → UserRepository → UserService → Router
	//
	// 依赖注入流程：
	// 1. 创建仓储（Repository）：依赖数据库客户端
	// 2. 创建应用服务（Service）：依赖仓储
	// 3. 创建路由（Router）：依赖应用服务
	userRepo := entrepo.NewUserRepository(entClient)
	userService = appuser.NewService(userRepo)

	// 步骤 5: 初始化 Temporal 客户端（可选）
	//
	// Temporal 说明：
	// - Temporal 是工作流编排引擎
	// - 用于处理长时间运行的任务和复杂业务流程
	// - 如果配置中没有 Temporal 地址，则跳过初始化
	//
	// 使用场景：
	// - 订单处理流程
	// - 用户注册流程
	// - 异步任务处理
	var temporalHandler *temporalhandler.Handler
	if cfg.Workflow.Temporal.Address != "" {
		temporalClient, err := temporal.NewClient(cfg.Workflow.Temporal.Address)
		if err != nil {
			logger.Warn("Failed to create temporal client", "error", err)
		} else {
			// 确保在程序退出时关闭 Temporal 客户端
			defer temporalClient.Close()
			temporalHandler = temporalhandler.NewHandler(temporalClient.Client())
			logger.Info("Temporal client initialized", "address", cfg.Workflow.Temporal.Address)
		}
	}

	// 步骤 6: 创建 HTTP 路由器
	//
	// 路由创建说明：
	// - 使用 Chi Router
	// - 注册所有 HTTP 路由和中间件
	// - 支持 REST API、健康检查等
	//
	// 中间件：
	// - 认证授权（JWT）
	// - 限流（Rate Limiting）
	// - 熔断器（Circuit Breaker）
	// - 请求追踪（Tracing）
	// - 性能监控（Metrics）
	// - CORS
	// - 恢复（Recovery）
	router := chiRouter.NewRouter(userService, temporalHandler)

	// 步骤 7: 创建 HTTP 服务器
	//
	// HTTP 服务器配置：
	// - Addr: 监听地址（格式：host:port）
	// - Handler: HTTP 处理器（Chi Router）
	// - ReadTimeout: 读取超时时间（防止慢客户端）
	// - WriteTimeout: 写入超时时间（防止慢客户端）
	// - IdleTimeout: 空闲连接超时时间（连接复用）
	//
	// 超时设置建议：
	// - ReadTimeout: 30-60 秒
	// - WriteTimeout: 30-60 秒
	// - IdleTimeout: 120 秒
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router.Handler(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 步骤 8: 启动服务器
	//
	// 启动说明：
	// - 在独立的 goroutine 中启动服务器，避免阻塞主线程
	// - ListenAndServe() 会阻塞直到服务器关闭
	// - 如果启动失败，记录错误并退出程序
	//
	// 错误处理：
	// - http.ErrServerClosed 是正常关闭，不需要处理
	// - 其他错误表示启动失败，需要退出程序
	go func() {
		logger.Info("HTTP server starting", "addr", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 步骤 9: 等待中断信号
	//
	// 信号监听说明：
	// - SIGINT: 用户按下 Ctrl+C
	// - SIGTERM: kill 命令发送的终止信号
	// - 收到信号后，开始优雅关闭流程
	//
	// 优雅关闭流程：
	// 1. 停止接收新请求
	// 2. 等待正在处理的请求完成（最多 30 秒）
	// 3. 关闭数据库连接
	// 4. 关闭其他资源
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 步骤 10: 优雅关闭
	//
	// 优雅关闭说明：
	// - 设置 30 秒超时，如果超时则强制关闭
	// - Shutdown() 会停止接收新请求，等待正在处理的请求完成
	// - 确保所有资源都被正确关闭
	//
	// 关闭顺序：
	// 1. HTTP 服务器（停止接收新请求，等待请求完成）
	// 2. 数据库连接（通过 defer entClient.Close()）
	// 3. Temporal 客户端（通过 defer temporalClient.Close()）
	// 4. OpenTelemetry Tracer（通过 defer shutdownTracer()）
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Application gracefully stopped.")
}
