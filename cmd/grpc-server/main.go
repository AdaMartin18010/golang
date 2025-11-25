package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/yourusername/golang/internal/config"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/infrastructure/database/postgres"
	// "github.com/yourusername/golang/internal/interfaces/grpc/handlers" // TODO: 生成 gRPC 代码后启用
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 初始化日志
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// 初始化数据库（TODO: 使用 Ent 客户端）
	var db interface{} = nil
	userRepo := postgres.NewUserRepository(db)
	userService := appuser.NewService(userRepo)

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	// TODO: 注册 gRPC 服务
	// userHandler := handlers.NewUserHandler(userService)
	// userpb.RegisterUserServiceServer(grpcServer, userHandler)
	_ = userService // 临时使用，避免未使用变量错误

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port+1) // gRPC 使用不同端口
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Failed to listen", "error", err, "addr", addr)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server starting", "addr", addr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("gRPC server shutting down...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		slog.Info("gRPC server stopped")
	case <-ctx.Done():
		slog.Warn("gRPC server shutdown timeout, forcing stop")
		grpcServer.Stop()
	}
}
