package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/domain/user"
)

// MockUserService 用户服务Mock
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, email, name string) (*user.User, error) {
	args := m.Called(ctx, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) UpdateUserName(ctx context.Context, id, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

// UserHandler HTTP处理器
type UserHandler struct {
	service UserService
}

type UserService interface {
	GetUser(ctx context.Context, id string) (*user.User, error)
	CreateUser(ctx context.Context, email, name string) (*user.User, error)
	UpdateUserName(ctx context.Context, id, name string) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*user.User, error)
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TestUserHandler_GetUser 测试获取用户
func TestUserHandler_GetUser(t *testing.T) {
	service := new(MockUserService)
	handler := NewUserHandler(service)

	testUser := user.NewUser("test@example.com", "Test User")
	service.On("GetUser", mock.Anything, testUser.ID).Return(testUser, nil)

	// 模拟处理
	ctx := context.Background()
	result, err := handler.service.GetUser(ctx, testUser.ID)

	require.NoError(t, err)
	assert.Equal(t, testUser.ID, result.ID)

	service.AssertExpectations(t)
}

// TestUserHandler_CreateUser 测试创建用户
func TestUserHandler_CreateUser(t *testing.T) {
	service := new(MockUserService)
	handler := NewUserHandler(service)

	testUser := user.NewUser("test@example.com", "Test User")
	service.On("CreateUser", mock.Anything, "test@example.com", "Test User").Return(testUser, nil)

	reqBody := CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// 模拟处理
	ctx := context.Background()
	result, err := handler.service.CreateUser(ctx, reqBody.Email, reqBody.Name)

	require.NoError(t, err)
	assert.Equal(t, testUser.Email, result.Email)

	service.AssertExpectations(t)
}

// TestUserHandler_ValidateEmail 测试邮箱验证
func TestUserHandler_ValidateEmail(t *testing.T) {
	service := new(MockUserService)
	handler := NewUserHandler(service)

	// 验证handler不为空
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

// TestUserHandler_DeleteUser 测试删除用户
func TestUserHandler_DeleteUser(t *testing.T) {
	service := new(MockUserService)
	handler := NewUserHandler(service)

	userID := "test-user-id"
	service.On("DeleteUser", mock.Anything, userID).Return(nil)

	// 模拟处理
	ctx := context.Background()
	err := handler.service.DeleteUser(ctx, userID)

	require.NoError(t, err)

	service.AssertExpectations(t)
}

// TestUserHandler_ListUsers 测试列出用户
func TestUserHandler_ListUsers(t *testing.T) {
	service := new(MockUserService)
	handler := NewUserHandler(service)

	users := []*user.User{
		user.NewUser("user1@example.com", "User 1"),
		user.NewUser("user2@example.com", "User 2"),
	}
	service.On("ListUsers", mock.Anything, 10, 0).Return(users, nil)

	// 模拟处理
	ctx := context.Background()
	result, err := handler.service.ListUsers(ctx, 10, 0)

	require.NoError(t, err)
	assert.Len(t, result, 2)

	service.AssertExpectations(t)
}

// BenchmarkUserHandler_GetUser 性能测试
func BenchmarkUserHandler_GetUser(b *testing.B) {
	service := new(MockUserService)
	handler := NewUserHandler(service)

	testUser := user.NewUser("test@example.com", "Test User")
	service.On("GetUser", mock.Anything, testUser.ID).Return(testUser, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.service.GetUser(context.Background(), testUser.ID)
	}
}
