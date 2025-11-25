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
	"github.com/yourusername/golang/internal/infrastructure/database/postgres"
	"github.com/yourusername/golang/internal/infrastructure/observability/otlp"
	chirouter "github.com/yourusername/golang/internal/interfaces/http/chi"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 初始化日志
	logger := otlp.NewLogger()
	slog.SetDefault(logger.Logger)

	// 初始化 OpenTelemetry
	ctx := context.Background()
	tracerProvider, err := otlp.NewTracerProvider(ctx, cfg.Observability.OTLP.Endpoint, cfg.Observability.OTLP.Insecure)
	if err != nil {
		slog.Warn("Failed to initialize tracer", "error", err)
	} else {
		defer tracerProvider.Shutdown(ctx)
	}

	metricsProvider, err := otlp.NewMetricsProvider(ctx, cfg.Observability.OTLP.Endpoint, cfg.Observability.OTLP.Insecure)
	if err != nil {
		slog.Warn("Failed to initialize metrics", "error", err)
	} else {
		defer metricsProvider.Shutdown(ctx)
	}

	// 初始化数据库（TODO: 使用 Ent 客户端）
	// 临时使用 nil，后续集成 Ent 后替换
	var db interface{} = nil
	userRepo := postgres.NewUserRepository(db)
	userService := appuser.NewService(userRepo)

	// 创建路由
	router := chirouter.NewRouter(userService)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router.Handler(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		slog.Info("Server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server shutting down...")

	// 优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	slog.Info("Server stopped")
}
