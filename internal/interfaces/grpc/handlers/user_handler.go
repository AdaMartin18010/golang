package handlers

import (
	appuser "github.com/yourusername/golang/internal/application/user"
	// TODO: 导入生成的 proto 代码
	// userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/user"
)

// UserHandler gRPC 用户处理器
type UserHandler struct {
	service *appuser.Service
	// TODO: 实现 gRPC 服务接口
	// userpb.UnimplementedUserServiceServer
}

// NewUserHandler 创建用户处理器
func NewUserHandler(service *appuser.Service) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser 创建用户
// TODO: 实现 gRPC 方法
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
