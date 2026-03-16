package handlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	healthpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/healthpb"
)

// HealthChecker 定义健康检查接口
// 实现此接口的类型可以提供自定义的健康检查逻辑
type HealthChecker interface {
	// Check 执行健康检查
	// 返回 nil 表示健康，否则返回错误
	Check(ctx context.Context) error
}

// HealthHandler gRPC 健康检查处理器
// 提供服务的健康状态检查功能
type HealthHandler struct {
	healthpb.UnimplementedHealthServiceServer
	logger    *slog.Logger
	checkers  []HealthChecker
	isReady   func() bool
}

// NewHealthHandler 创建健康检查处理器
//
// 参数:
//   - logger: 日志记录器（可为 nil，将使用默认日志）
//   - checkers: 健康检查器列表（可选）
//
// 返回:
//   - *HealthHandler: 健康检查处理器实例
func NewHealthHandler(logger *slog.Logger, checkers ...HealthChecker) *HealthHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &HealthHandler{
		logger:   logger,
		checkers: checkers,
		isReady:  func() bool { return true }, // 默认始终就绪
	}
}

// NewHealthHandlerWithReadyFunc 创建带就绪检查函数的健康检查处理器
//
// 参数:
//   - logger: 日志记录器（可为 nil，将使用默认日志）
//   - isReady: 就绪检查函数，返回 true 表示服务就绪
//   - checkers: 健康检查器列表（可选）
//
// 返回:
//   - *HealthHandler: 健康检查处理器实例
func NewHealthHandlerWithReadyFunc(
	logger *slog.Logger,
	isReady func() bool,
	checkers ...HealthChecker,
) *HealthHandler {
	h := NewHealthHandler(logger, checkers...)
	if isReady != nil {
		h.isReady = isReady
	}
	return h
}

// Check 健康检查
//
// 检查服务健康状态，可以包含以下检查：
//   - 基本服务状态
//   - 数据库连接状态
//   - 依赖服务状态
//   - 资源使用情况
//
// 返回:
//   - SERVING: 服务正常
//   - NOT_SERVING: 服务异常
func (h *HealthHandler) Check(ctx context.Context, req *emptypb.Empty) (*healthpb.HealthResponse, error) {
	h.logger.Debug("gRPC Health Check request")

	// 检查服务是否就绪
	if !h.isReady() {
		h.logger.Warn("Health Check: service not ready")
		return &healthpb.HealthResponse{
			Status: healthpb.HealthResponse_NOT_SERVING,
		}, nil
	}

	// 执行自定义健康检查
	for i, checker := range h.checkers {
		if err := checker.Check(ctx); err != nil {
			h.logger.Error("Health Check failed",
				"checker_index", i,
				"error", err,
			)
			return &healthpb.HealthResponse{
				Status: healthpb.HealthResponse_NOT_SERVING,
			}, nil
		}
	}

	h.logger.Debug("Health Check: service is healthy")
	return &healthpb.HealthResponse{
		Status: healthpb.HealthResponse_SERVING,
	}, nil
}

// CheckWithDetails 带详细信息的健康检查（扩展方法）
//
// 提供更详细的健康检查信息，包括各个组件的状态
func (h *HealthHandler) CheckWithDetails(ctx context.Context) (*HealthDetails, error) {
	details := &HealthDetails{
		Status:    healthpb.HealthResponse_SERVING,
		Component: make(map[string]HealthComponentStatus),
	}

	// 检查服务就绪状态
	if !h.isReady() {
		details.Status = healthpb.HealthResponse_NOT_SERVING
		details.Message = "service not ready"
		return details, nil
	}

	// 执行各个检查器的健康检查
	for i, checker := range h.checkers {
		componentName := getCheckerName(checker, i)
		if err := checker.Check(ctx); err != nil {
			details.Status = healthpb.HealthResponse_NOT_SERVING
			details.Component[componentName] = HealthComponentStatus{
				Status:  healthpb.HealthResponse_NOT_SERVING,
				Message: err.Error(),
			}
		} else {
			details.Component[componentName] = HealthComponentStatus{
				Status: healthpb.HealthResponse_SERVING,
			}
		}
	}

	return details, nil
}

// HealthDetails 健康检查详细信息
type HealthDetails struct {
	Status    healthpb.HealthResponse_Status
	Message   string
	Component map[string]HealthComponentStatus
}

// HealthComponentStatus 组件健康状态
type HealthComponentStatus struct {
	Status  healthpb.HealthResponse_Status
	Message string
}

// getCheckerName 获取检查器名称
func getCheckerName(checker HealthChecker, index int) string {
	// 尝试获取具体类型名称
	// 这里简化处理，实际可以使用反射获取类型名
	return "checker_" + string(rune('0'+index))
}

// SimpleHealthChecker 简单的健康检查器实现
type SimpleHealthChecker struct {
	name  string
	check func(ctx context.Context) error
}

// NewSimpleHealthChecker 创建简单健康检查器
//
// 参数:
//   - name: 检查器名称
//   - check: 检查函数
//
// 返回:
//   - *SimpleHealthChecker: 简单健康检查器实例
func NewSimpleHealthChecker(name string, check func(ctx context.Context) error) *SimpleHealthChecker {
	return &SimpleHealthChecker{
		name:  name,
		check: check,
	}
}

// Check 执行健康检查
func (c *SimpleHealthChecker) Check(ctx context.Context) error {
	if c.check != nil {
		return c.check(ctx)
	}
	return nil
}

// DatabaseHealthChecker 数据库健康检查器
type DatabaseHealthChecker struct {
	ping func(ctx context.Context) error
}

// NewDatabaseHealthChecker 创建数据库健康检查器
//
// 参数:
//   - ping: 数据库 ping 函数
//
// 返回:
//   - *DatabaseHealthChecker: 数据库健康检查器实例
func NewDatabaseHealthChecker(ping func(ctx context.Context) error) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{ping: ping}
}

// Check 执行数据库健康检查
func (c *DatabaseHealthChecker) Check(ctx context.Context) error {
	if c.ping == nil {
		return status.Error(codes.Internal, "database ping function not set")
	}
	return c.ping(ctx)
}

// Ensure HealthHandler implements HealthServiceServer
var _ healthpb.HealthServiceServer = (*HealthHandler)(nil)
