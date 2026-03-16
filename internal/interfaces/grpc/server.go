package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/yourusername/golang/internal/app/user"
	healthpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/healthpb"
	userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/userpb"

	"github.com/yourusername/golang/internal/interfaces/grpc/handlers"
	"github.com/yourusername/golang/internal/interfaces/grpc/interceptors"
)

// Server gRPC 服务器
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	logger     *slog.Logger
	config     *Config

	// 服务处理器
	userHandler   *handlers.UserHandler
	healthHandler *handlers.HealthHandler
}

// Config gRPC 服务器配置
type Config struct {
	// Host 服务器主机地址（默认 "0.0.0.0"）
	Host string
	// Port 服务器端口（默认 50051）
	Port int
	// MaxConnectionIdle 连接最大空闲时间
	MaxConnectionIdle time.Duration
	// MaxConnectionAge 连接最大生命周期
	MaxConnectionAge time.Duration
	// MaxConnectionAgeGrace 连接关闭宽限期
	MaxConnectionAgeGrace time.Duration
	// Time ping 周期
	Time time.Duration
	// Timeout ping 超时时间
	Timeout time.Duration
	// EnableReflection 是否启用 gRPC 反射（开发环境建议启用）
	EnableReflection bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:                  "0.0.0.0",
		Port:                  50051,
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Minute,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
		EnableReflection:      true,
	}
}

// Address 返回服务器地址
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// ServerOption 服务器选项函数
type ServerOption func(*Server)

// WithLogger 设置日志记录器
func WithLogger(logger *slog.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithConfig 设置服务器配置
func WithConfig(config *Config) ServerOption {
	return func(s *Server) {
		s.config = config
	}
}

// WithHealthChecker 添加健康检查器
func WithHealthChecker(checker handlers.HealthChecker) ServerOption {
	return func(s *Server) {
		if s.healthHandler != nil {
			// 重新创建 healthHandler 以包含新的 checker
			// 注意：这需要在所有选项处理完成后再创建 handler
		}
	}
}

// NewServer 创建 gRPC 服务器
//
// 参数:
//   - userService: 用户应用服务
//   - opts: 服务器选项
//
// 返回:
//   - *Server: gRPC 服务器实例
//   - error: 错误信息
func NewServer(userService *user.Service, opts ...ServerOption) (*Server, error) {
	s := &Server{
		logger: slog.Default(),
		config: DefaultConfig(),
	}

	// 应用选项
	for _, opt := range opts {
		opt(s)
	}

	// 创建处理器
	s.userHandler = handlers.NewUserHandler(userService, s.logger)
	s.healthHandler = handlers.NewHealthHandler(s.logger)

	// 创建 gRPC 服务器选项
	serverOpts := []grpc.ServerOption{
		// 连接 Keepalive 参数
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     s.config.MaxConnectionIdle,
			MaxConnectionAge:      s.config.MaxConnectionAge,
			MaxConnectionAgeGrace: s.config.MaxConnectionAgeGrace,
			Time:                  s.config.Time,
			Timeout:               s.config.Timeout,
		}),
		// 拦截器链
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingUnaryInterceptor,
			interceptors.TracingUnaryInterceptor,
		),
	}

	// 创建 gRPC 服务器
	s.grpcServer = grpc.NewServer(serverOpts...)

	// 注册服务
	userpb.RegisterUserServiceServer(s.grpcServer, s.userHandler)
	healthpb.RegisterHealthServiceServer(s.grpcServer, s.healthHandler)

	// 启用反射（开发环境）
	if s.config.EnableReflection {
		reflection.Register(s.grpcServer)
		s.logger.Info("gRPC reflection enabled")
	}

	return s, nil
}

// Start 启动 gRPC 服务器
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - error: 错误信息
func (s *Server) Start(ctx context.Context) error {
	addr := s.config.Address()

	// 创建监听器
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}
	s.listener = listener

	s.logger.Info("gRPC server starting",
		"address", addr,
		"reflection_enabled", s.config.EnableReflection,
	)

	// 在 goroutine 中启动服务器
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.logger.Error("gRPC server serve error", "error", err)
		}
	}()

	return nil
}

// Stop 停止 gRPC 服务器（优雅关闭）
//
// 参数:
//   - ctx: 上下文，用于控制关闭超时
//
// 返回:
//   - error: 错误信息
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("gRPC server stopping...")

	// 创建停止完成通道
	stopped := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	// 等待优雅关闭完成或超时
	select {
	case <-ctx.Done():
		s.logger.Warn("gRPC server graceful stop timeout, forcing stop")
		s.grpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	}
}

// Shutdown 完全关闭服务器（立即关闭）
func (s *Server) Shutdown() {
	s.logger.Info("gRPC server shutting down...")
	s.grpcServer.Stop()
	s.logger.Info("gRPC server shut down")
}

// Addr 返回服务器监听地址
func (s *Server) Addr() net.Addr {
	if s.listener != nil {
		return s.listener.Addr()
	}
	return nil
}

// GetUserHandler 返回用户服务处理器（用于测试）
func (s *Server) GetUserHandler() *handlers.UserHandler {
	return s.userHandler
}

// GetHealthHandler 返回健康检查处理器（用于测试）
func (s *Server) GetHealthHandler() *handlers.HealthHandler {
	return s.healthHandler
}

// SetReadyFunc 设置健康检查的就绪函数
func (s *Server) SetReadyFunc(isReady func() bool) {
	// 重新创建 healthHandler 以包含就绪函数
	s.healthHandler = handlers.NewHealthHandlerWithReadyFunc(
		s.logger,
		isReady,
	)
	// 重新注册健康服务
	healthpb.RegisterHealthServiceServer(s.grpcServer, s.healthHandler)
}

// RegisterHealthChecker 注册健康检查器
func (s *Server) RegisterHealthChecker(checker handlers.HealthChecker) {
	// 注意：此方法需要在服务器启动前调用
	// 实际实现可能需要重新创建 healthHandler
	s.logger.Warn("RegisterHealthChecker should be called before Start")
}

