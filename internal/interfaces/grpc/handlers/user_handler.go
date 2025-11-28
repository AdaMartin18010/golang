package handlers

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yourusername/golang/internal/application/user"
	userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/userpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserHandler gRPC 用户服务处理器
type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	service *user.Service
}

// NewUserHandler 创建用户服务处理器
func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetUser 获取用户
func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	resp, err := h.service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &userpb.GetUserResponse{
		User: toProtoUserFromResponse(resp),
	}, nil
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	createReq := user.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}
	resp, err := h.service.CreateUser(ctx, createReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.CreateUserResponse{
		User: toProtoUserFromResponse(resp),
	}, nil
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	updateReq := user.UpdateUserRequest{}
	if req.Name != "" {
		updateReq.Name = &req.Name
	}
	if req.Email != "" {
		updateReq.Email = &req.Email
	}

	resp, err := h.service.UpdateUser(ctx, req.Id, updateReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.UpdateUserResponse{
		User: toProtoUserFromResponse(resp),
	}, nil
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	err := h.service.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil
}

// ListUsers 列出用户（流式）
func (h *UserHandler) ListUsers(req *userpb.ListUsersRequest, stream userpb.UserService_ListUsersServer) error {
	// TODO: 实现流式列表
	return status.Error(codes.Unimplemented, "ListUsers not implemented")
}

// toProtoUserFromResponse 从 UserResponse 转换为 Proto 消息
func toProtoUserFromResponse(resp *user.UserResponse) *userpb.User {
	createdAt, _ := time.Parse(time.RFC3339, resp.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, resp.UpdatedAt)

	return &userpb.User{
		Id:        resp.ID,
		Name:      resp.Name,
		Email:     resp.Email,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}
}
