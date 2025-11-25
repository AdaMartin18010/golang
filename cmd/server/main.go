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
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	logger := otlp.NewLogger()
	slog.SetDefault(logger.Logger)
	logger.Info("Application starting...")

	// 3. 初始化 OpenTelemetry Tracer
	ctx := context.Background()
	shutdownTracer, err := otlp.NewTracerProvider(ctx, cfg.Observability.TraceEndpoint, true)
	if err != nil {
		logger.Warn("Failed to initialize tracer", "error", err)
	} else {
		defer func() {
			if err := shutdownTracer(ctx); err != nil {
				logger.Error("Failed to shutdown tracer", "error", err)
			}
		}()
	}

	// 4. 初始化数据库（使用 Ent）
	// 注意：需要先运行 make generate-ent 生成 Ent 代码
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
		logger.Warn("Failed to initialize database", "error", err)
		logger.Info("Note: To use Ent, run 'make generate-ent' first to generate code")
		// 使用临时的内存实现（仅用于开发/测试）
		// 在实际使用中，应该确保数据库连接成功
		logger.Error("Database connection failed. Please check your database configuration and run 'make generate-ent'.")
		os.Exit(1)
	} else {
		defer entClient.Close()
		logger.Info("Database connected successfully")
		userRepo := entrepo.NewUserRepository(entClient)
		userService = appuser.NewService(userRepo)
	}

	// 5. 初始化 Temporal 客户端（可选）
	var temporalHandler *temporalhandler.Handler
	if cfg.Workflow.Temporal.Address != "" {
		temporalClient, err := temporal.NewClient(cfg.Workflow.Temporal.Address)
		if err != nil {
			logger.Warn("Failed to create temporal client", "error", err)
		} else {
			defer temporalClient.Close()
			temporalHandler = temporalhandler.NewHandler(temporalClient.Client())
			logger.Info("Temporal client initialized", "address", cfg.Workflow.Temporal.Address)
		}
	}

	// 6. 创建 HTTP 路由器
	router := chiRouter.NewRouter(userService, temporalHandler)

	// 7. 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router.Handler(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 8. 启动服务器
	go func() {
		logger.Info("HTTP server starting", "addr", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 9. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 10. 优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Application gracefully stopped.")
}
