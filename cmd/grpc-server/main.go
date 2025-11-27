// Package main gRPC 服务器应用入口
//
// 设计原理：
// 1. 这是 gRPC 服务器的入口点
// 2. 提供高性能的 RPC 服务，使用 Protocol Buffers 序列化
// 3. 支持流式 RPC（Streaming RPC）
// 4. 实现优雅关闭，确保请求处理完成后再退出
//
// 启动流程：
// 1. 加载配置
// 2. 初始化日志
// 3. 初始化数据库连接
// 4. 组装依赖（仓储、服务）
// 5. 创建 gRPC 服务器
// 6. 注册 gRPC 服务
// 7. 启动服务器（在 goroutine 中）
// 8. 等待中断信号
// 9. 优雅关闭
//
// gRPC 特点：
// - 高性能：使用 HTTP/2 和 Protocol Buffers
// - 类型安全：通过 .proto 文件定义接口
// - 流式支持：支持客户端流、服务端流、双向流
// - 跨语言：支持多种编程语言
//
// 使用方式：
//   go run cmd/grpc-server/main.go
//   # 或
//   go build -o bin/grpc-server cmd/grpc-server/main.go
//   ./bin/grpc-server
//
// 配置：
// - 配置文件：configs/config.yaml
// - 默认端口：8081（Server.Port + 1）
// - 需要先运行 protoc 生成 gRPC 代码
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
	// 步骤 1: 加载配置
	//
	// 配置加载说明：
	// - 从 configs/config.yaml 加载配置
	// - 支持环境变量覆盖（APP_* 前缀）
	// - 如果配置文件不存在，使用默认配置
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 步骤 2: 初始化日志
	//
	// 日志初始化说明：
	// - 使用结构化日志（slog）
	// - JSON 格式输出，便于日志收集和分析
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// 步骤 3: 初始化数据库（TODO: 使用 Ent 客户端）
	//
	// 数据库初始化说明：
	// - 当前使用占位符，需要替换为实际的 Ent 客户端
	// - 参考 cmd/server/main.go 中的数据库初始化代码
	// - 需要创建数据库连接、运行迁移等
	var db interface{} = nil
	userRepo := postgres.NewUserRepository(db)
	userService := appuser.NewService(userRepo)

	// 步骤 4: 创建 gRPC 服务器
	//
	// gRPC 服务器说明：
	// - 使用 google.golang.org/grpc 创建服务器
	// - 可以配置拦截器（Interceptor）用于认证、日志、追踪等
	// - 可以配置选项（Options）用于压缩、超时等
	//
	// 示例配置：
	//   grpcServer := grpc.NewServer(
	//       grpc.UnaryInterceptor(authInterceptor),
	//       grpc.StreamInterceptor(streamInterceptor),
	//   )
	grpcServer := grpc.NewServer()

	// 步骤 5: 注册 gRPC 服务
	//
	// 服务注册说明：
	// - 需要先运行 protoc 生成 gRPC 代码
	// - 生成代码位置：internal/interfaces/grpc/proto/
	// - 注册服务处理器，将 gRPC 请求路由到应用服务
	//
	// 注册流程：
	// 1. 创建 Handler（实现 gRPC 服务接口）
	// 2. 注册服务到 gRPC 服务器
	//
	// TODO: 生成 gRPC 代码后启用
	// userHandler := handlers.NewUserHandler(userService)
	// userpb.RegisterUserServiceServer(grpcServer, userHandler)
	_ = userService // 临时使用，避免未使用变量错误

	// 步骤 6: 启动服务器
	//
	// 启动说明：
	// - gRPC 使用不同的端口（Server.Port + 1），避免与 HTTP 服务器冲突
	// - 在独立的 goroutine 中启动服务器，避免阻塞主线程
	// - 如果启动失败，记录错误并退出程序
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

	// 步骤 7: 等待中断信号
	//
	// 信号监听说明：
	// - 监听 SIGINT（Ctrl+C）和 SIGTERM（kill 命令）
	// - 收到信号后，开始优雅关闭流程
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("gRPC server shutting down...")

	// 步骤 8: 优雅关闭
	//
	// 优雅关闭说明：
	// - 设置 30 秒超时，如果超时则强制关闭
	// - GracefulStop() 会停止接收新请求，等待正在处理的请求完成
	// - 如果超时，使用 Stop() 强制关闭
	//
	// 关闭流程：
	// 1. 停止接收新请求
	// 2. 等待正在处理的请求完成（最多 30 秒）
	// 3. 关闭监听器
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
