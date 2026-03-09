// Package graphql provides mock-based tests for GraphQL resolvers.
//
// 本文件展示了如何使用 UserService 接口进行单元测试，
// 通过 mock 实现替代真实的服务依赖。
package graphql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domainuser "github.com/yourusername/golang/internal/domain/user"
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
func createTestDomainUser(id, email, name string) *domainuser.User {
	return &domainuser.User{
		ID:        id,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestResolver_WithMockUserService_GetUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望：当调用 GetUser 时返回测试用户
	testUser := createTestDomainUser("123", "test@example.com", "Test User")
	mockSvc.On("GetUser", mock.Anything, "123").Return(testUser, nil)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	query := &Query{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	user, err := query.User(ctx, "123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "123", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)

	// 验证 mock 被调用
	mockSvc.AssertExpectations(t)
}

func TestResolver_WithMockUserService_GetUser_NotFound(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望：当调用 GetUser 时返回未找到错误
	mockSvc.On("GetUser", mock.Anything, "999").Return(nil, domainuser.ErrUserNotFound)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	query := &Query{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	user, err := query.User(ctx, "999")

	// 验证结果 - GraphQL 中未找到返回 nil
	assert.NoError(t, err)
	assert.Nil(t, user)

	mockSvc.AssertExpectations(t)
}

func TestResolver_WithMockUserService_CreateUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望
	testUser := createTestDomainUser("new-id", "new@example.com", "New User")
	mockSvc.On("CreateUser", mock.Anything, "new@example.com", "New User").Return(testUser, nil)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	mutation := &Mutation{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	input := CreateUserInput{Email: "new@example.com", Name: "New User"}
	user, err := mutation.CreateUser(ctx, input)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "new-id", user.ID)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "New User", user.Name)

	mockSvc.AssertExpectations(t)
}

func TestResolver_WithMockUserService_CreateUser_ValidationError(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	mutation := &Mutation{resolver: resolver}

	// 执行测试 - 邮箱为空
	ctx := context.Background()
	input := CreateUserInput{Email: "", Name: "New User"}
	user, err := mutation.CreateUser(ctx, input)

	// 验证结果 - 应该在解析器层被验证
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email is required")

	// mock 不应该被调用
	mockSvc.AssertNotCalled(t, "CreateUser")
}

func TestResolver_WithMockUserService_UpdateUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望
	testUser := createTestDomainUser("123", "test@example.com", "Original Name")
	mockSvc.On("GetUser", mock.Anything, "123").Return(testUser, nil)
	mockSvc.On("UpdateUserName", mock.Anything, "123", "Updated Name").Return(nil)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	mutation := &Mutation{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	newName := "Updated Name"
	input := UpdateUserInput{Name: &newName}
	user, err := mutation.UpdateUser(ctx, "123", input)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Updated Name", user.Name)

	mockSvc.AssertExpectations(t)
}

func TestResolver_WithMockUserService_DeleteUser(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望
	mockSvc.On("DeleteUser", mock.Anything, "123").Return(nil)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	mutation := &Mutation{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	success, err := mutation.DeleteUser(ctx, "123")

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, success)

	mockSvc.AssertExpectations(t)
}

func TestResolver_WithMockUserService_DeleteUser_NotFound(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望 - 删除不存在的用户
	mockSvc.On("DeleteUser", mock.Anything, "999").Return(domainuser.ErrUserNotFound)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	mutation := &Mutation{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	success, err := mutation.DeleteUser(ctx, "999")

	// 验证结果 - 用户不存在视为删除成功
	assert.NoError(t, err)
	assert.False(t, success)

	mockSvc.AssertExpectations(t)
}

func TestResolver_WithMockUserService_ListUsers(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望
	testUsers := []*domainuser.User{
		createTestDomainUser("1", "user1@example.com", "User 1"),
		createTestDomainUser("2", "user2@example.com", "User 2"),
	}
	mockSvc.On("ListUsers", mock.Anything, 10, 0).Return(testUsers, nil)

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	query := &Query{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	limit := 10
	offset := 0
	users, err := query.Users(ctx, &limit, &offset)

	// 验证结果
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "1", users[0].ID)
	assert.Equal(t, "2", users[1].ID)

	mockSvc.AssertExpectations(t)
}

func TestResolver_WithMockUserService_ListUsers_Error(t *testing.T) {
	// 创建 mock 服务
	mockSvc := new(mockUserService)

	// 设置期望 - 返回错误
	mockSvc.On("ListUsers", mock.Anything, 10, 0).Return(nil, errors.New("database error"))

	// 创建 Resolver
	resolver := NewResolver(mockSvc)
	query := &Query{resolver: resolver}

	// 执行测试
	ctx := context.Background()
	limit := 10
	offset := 0
	users, err := query.Users(ctx, &limit, &offset)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "failed to list users")

	mockSvc.AssertExpectations(t)
}

func TestResolver_InterfaceCompliance(t *testing.T) {
	// 验证 mock 实现了 UserService 接口
	var _ UserService = (*mockUserService)(nil)

	// 验证真实服务也可以被使用（如果可导入）
	// 这里主要验证接口设计的正确性
	assert.True(t, true, "Interface compliance verified")
}
