// Package handlers provides mock-based tests for gRPC user handlers.
//
// 本文件展示了如何使用 UserService 接口进行单元测试，
// 通过 mock 实现替代真实的服务依赖。
package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	domainuser "github.com/yourusername/golang/internal/domain/user"
	userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/userpb"
)

// mockUserService 是 UserService 接口的 mock 实现
type mockUserService struct {
	mock.Mock
}

// 确保 mockUserService 实现了 UserService 接口
var _ UserService = (*mockUserService)(nil)

func (m *mockUserService) GetUser(ctx context.Context, id string) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *mockUserService) CreateUser(ctx context.Context, email, name string) (*domainuser.User, error) {
	args := m.Called(ctx, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *mockUserService) UpdateUserName(ctx context.Context, id, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserService) ListUsers(ctx context.Context, limit, offset int) ([]*domainuser.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainuser.User), args.Error(1)
}

// createTestDomainUser 创建测试用的领域用户
func createTestDomainUserForHandler(id, email, name string) *domainuser.User {
	return &domainuser.User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestUserHandler_WithMockService_GetUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望
	testUser := createTestDomainUserForHandler("123", "test@example.com", "Test User")
	mockSvc.On("GetUser", mock.Anything, "123").Return(testUser, nil)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.GetUserRequest{Id: "123"}
	resp, err := handler.GetUser(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "123", resp.User.Id)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "Test User", resp.User.Name)

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_GetUser_NotFound(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望 - 用户不存在
	mockSvc.On("GetUser", mock.Anything, "999").Return(nil, domainuser.ErrUserNotFound)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.GetUserRequest{Id: "999"}
	resp, err := handler.GetUser(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)

	// 验证 gRPC 状态码
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, statusErr.Code())

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_GetUser_EmptyID(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试 - 空的用户 ID
	ctx := context.Background()
	req := &userpb.GetUserRequest{Id: ""}
	resp, err := handler.GetUser(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)

	// 验证 gRPC 状态码
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())

	// mock 不应该被调用
	mockSvc.AssertNotCalled(t, "GetUser")
}

func TestUserHandler_WithMockService_CreateUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望
	testUser := createTestDomainUserForHandler("new-id", "new@example.com", "New User")
	mockSvc.On("CreateUser", mock.Anything, "new@example.com", "New User").Return(testUser, nil)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.CreateUserRequest{Email: "new@example.com", Name: "New User"}
	resp, err := handler.CreateUser(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new-id", resp.User.Id)
	assert.Equal(t, "new@example.com", resp.User.Email)
	assert.Equal(t, "New User", resp.User.Name)

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_CreateUser_AlreadyExists(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望 - 用户已存在
	mockSvc.On("CreateUser", mock.Anything, "exists@example.com", "Existing User").
		Return(nil, domainuser.ErrUserAlreadyExists)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.CreateUserRequest{Email: "exists@example.com", Name: "Existing User"}
	resp, err := handler.CreateUser(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)

	// 验证 gRPC 状态码
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, statusErr.Code())

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_CreateUser_InvalidEmail(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望 - 邮箱格式无效
	mockSvc.On("CreateUser", mock.Anything, "invalid-email", "Test User").
		Return(nil, domainuser.ErrInvalidEmailFormat)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.CreateUserRequest{Email: "invalid-email", Name: "Test User"}
	resp, err := handler.CreateUser(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)

	// 验证 gRPC 状态码
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_UpdateUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望
	testUser := createTestDomainUserForHandler("123", "test@example.com", "Original Name")
	mockSvc.On("GetUser", mock.Anything, "123").Return(testUser, nil)
	mockSvc.On("UpdateUserName", mock.Anything, "123", "Updated Name").Return(nil)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.UpdateUserRequest{Id: "123", Name: "Updated Name"}
	resp, err := handler.UpdateUser(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Updated Name", resp.User.Name)

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_DeleteUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望
	mockSvc.On("DeleteUser", mock.Anything, "123").Return(nil)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.DeleteUserRequest{Id: "123"}
	resp, err := handler.DeleteUser(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_DeleteUser_NotFound(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望 - 用户不存在
	mockSvc.On("DeleteUser", mock.Anything, "999").Return(domainuser.ErrUserNotFound)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 执行测试
	ctx := context.Background()
	req := &userpb.DeleteUserRequest{Id: "999"}
	resp, err := handler.DeleteUser(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)

	// 验证 gRPC 状态码
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, statusErr.Code())

	mockSvc.AssertExpectations(t)
}

func TestUserHandler_WithMockService_ListUsers(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望 - 使用分页参数计算出的 limit 和 offset
	testUsers := []*domainuser.User{
		createTestDomainUserForHandler("1", "user1@example.com", "User 1"),
		createTestDomainUserForHandler("2", "user2@example.com", "User 2"),
	}
	mockSvc.On("ListUsers", mock.Anything, 10, 0).Return(testUsers, nil)

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 创建 mock 流
	mockStream := new(mockListUsersStream)
	mockStream.On("Context").Return(context.Background())
	mockStream.On("Send", mock.AnythingOfType("*userpb.User")).Return(nil).Twice()

	// 执行测试
	req := &userpb.ListUsersRequest{Page: 1, PageSize: 10}
	err := handler.ListUsers(req, mockStream)

	// 验证结果
	assert.NoError(t, err)

	mockSvc.AssertExpectations(t)
	mockStream.AssertExpectations(t)
}

func TestUserHandler_WithMockService_ListUsers_ServiceError(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 设置期望 - 服务返回错误
	mockSvc.On("ListUsers", mock.Anything, 10, 0).Return(nil, errors.New("database connection failed"))

	// 创建 Handler
	handler := NewUserHandler(mockSvc, logger)

	// 创建 mock 流
	mockStream := new(mockListUsersStream)
	mockStream.On("Context").Return(context.Background())

	// 执行测试
	req := &userpb.ListUsersRequest{Page: 1, PageSize: 10}
	err := handler.ListUsers(req, mockStream)

	// 验证结果
	assert.Error(t, err)

	// 验证 gRPC 状态码
	statusErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, statusErr.Code())

	mockSvc.AssertExpectations(t)
}

// mockListUsersStream 是 UserService_ListUsersServer 接口的 mock 实现
type mockListUsersStream struct {
	mock.Mock
}

func (m *mockListUsersStream) Send(user *userpb.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockListUsersStream) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// 其他必需的方法（Protobuf 接口要求）
func (m *mockListUsersStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockListUsersStream) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockListUsersStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockListUsersStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *mockListUsersStream) SetTrailer(md metadata.MD) {
}

func TestUserHandler_InterfaceCompliance(t *testing.T) {
	// 验证 mock 实现了 UserService 接口
	var _ UserService = (*mockUserService)(nil)

	// 验证真实服务也可以被使用（如果可导入）
	// 这里主要验证接口设计的正确性
	assert.True(t, true, "Interface compliance verified")
}

func TestUserHandler_NilLogger(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 创建 Handler - 传入 nil logger，应该使用默认日志
	handler := NewUserHandler(mockSvc, nil)

	// 验证 handler 被创建且 logger 不为 nil
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
}
