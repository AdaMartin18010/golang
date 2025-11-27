// Package handlers provides gRPC request handlers for the interfaces layer.
//
// gRPC 处理器负责：
// 1. 接收 gRPC 请求
// 2. 将 Protocol Buffers 消息转换为应用层 DTO
// 3. 调用应用层服务处理业务逻辑
// 4. 将应用层结果转换为 Protocol Buffers 响应
// 5. 处理错误并返回适当的 gRPC 状态码
//
// 设计原则：
// 1. 薄层设计：处理器只负责协议转换，不包含业务逻辑
// 2. 错误映射：将应用层错误映射为 gRPC 状态码
// 3. 类型安全：使用生成的 Protocol Buffers 代码确保类型安全
// 4. 可测试性：通过依赖注入便于单元测试
//
// 工作流程：
// 1. 接收 gRPC 请求（Protocol Buffers 格式）
// 2. 验证请求参数
// 3. 转换为应用层 DTO
// 4. 调用应用层服务
// 5. 将结果转换为 Protocol Buffers 响应
// 6. 返回响应或错误
//
// 示例：
//
//	// 创建处理器
//	userService := appuser.NewService(userRepo)
//	handler := handlers.NewUserHandler(userService)
//
//	// 注册到 gRPC 服务器
//	userpb.RegisterUserServiceServer(grpcServer, handler)
package handlers

import (
	appuser "github.com/yourusername/golang/internal/application/user"
	// TODO: 导入生成的 proto 代码
	// 生成命令：protoc --go_out=. --go-grpc_out=. proto/user.proto
	// userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/user"
)

// UserHandler 是 gRPC 用户服务的处理器实现。
//
// 功能说明：
// - 实现 gRPC 服务接口（由 Protocol Buffers 生成）
// - 处理用户相关的 gRPC 请求（创建、查询、更新、删除等）
// - 将 gRPC 请求转换为应用层调用
//
// 设计说明：
// - 依赖应用层服务（appuser.Service）处理业务逻辑
// - 实现 userpb.UserServiceServer 接口（待生成）
// - 通过依赖注入接收服务实例
//
// 使用示例：
//
//	// 创建应用层服务
//	userService := appuser.NewService(userRepo)
//
//	// 创建 gRPC 处理器
//	handler := handlers.NewUserHandler(userService)
//
//	// 注册到 gRPC 服务器
//	userpb.RegisterUserServiceServer(grpcServer, handler)
//
// 注意事项：
// - 需要先运行 protoc 生成 gRPC 代码
// - 处理器方法应处理错误并返回适当的 gRPC 状态码
// - 可以使用拦截器（Interceptor）进行认证、日志、追踪等
type UserHandler struct {
	service *appuser.Service
	// TODO: 实现 gRPC 服务接口
	// 生成 proto 代码后，取消注释以下行：
	// userpb.UnimplementedUserServiceServer
	//
	// 这确保实现了所有必需的方法，即使某些方法暂时未实现
}

// NewUserHandler 创建并初始化 gRPC 用户处理器。
//
// 功能说明：
// - 接收应用层服务实例
// - 创建并返回配置好的处理器
//
// 参数：
// - service: 应用层用户服务实例
//   用于处理实际的业务逻辑
//
// 返回：
// - *UserHandler: 配置好的处理器实例
//
// 使用示例：
//
//	userService := appuser.NewService(userRepo)
//	handler := handlers.NewUserHandler(userService)
//
// 注意事项：
// - 服务实例应通过依赖注入提供
// - 处理器是无状态的，可以在多个请求间共享
func NewUserHandler(service *appuser.Service) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser 创建用户的 gRPC 方法实现。
//
// 功能说明：
// - 接收 CreateUserRequest（Protocol Buffers 格式）
// - 转换为应用层 CreateUserRequest
// - 调用应用层服务创建用户
// - 将结果转换为 UserResponse（Protocol Buffers 格式）
//
// 参数：
// - ctx: 上下文，包含请求元数据（认证信息、超时等）
// - req: 创建用户请求（Protocol Buffers 格式）
//
// 返回：
// - *userpb.UserResponse: 用户响应（Protocol Buffers 格式）
// - error: 如果创建失败，返回 gRPC 错误
//   错误应使用 status.Error() 包装，包含适当的错误码
//
// 错误处理：
// - InvalidArgument: 请求参数无效
// - AlreadyExists: 用户已存在
// - Internal: 内部服务器错误
//
// 使用示例（待实现）：
//
//	func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
//	    // 1. 参数验证
//	    if req.Email == "" {
//	        return nil, status.Error(codes.InvalidArgument, "email is required")
//	    }
//
//	    // 2. 转换为应用层请求
//	    appReq := appuser.CreateUserRequest{
//	        Email: req.Email,
//	        Name:  req.Name,
//	    }
//
//	    // 3. 调用应用层服务
//	    user, err := h.service.CreateUser(ctx, appReq)
//	    if err != nil {
//	        // 4. 错误映射
//	        return nil, mapAppErrorToGRPCError(err)
//	    }
//
//	    // 5. 转换为响应
//	    return toProtoUser(user), nil
//	}
//
// TODO: 实现 gRPC 方法
// 1. 生成 proto 代码：protoc --go_out=. --go-grpc_out=. proto/user.proto
// 2. 取消注释并实现以下方法
// 3. 实现错误映射函数 mapAppErrorToGRPCError
// 4. 实现转换函数 toProtoUser
//
// func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
// 	user, err := h.service.CreateUser(ctx, appuser.CreateUserRequest{
// 		Email: req.Email,
// 		Name:  req.Name,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return toProtoUser(user), nil
// }
