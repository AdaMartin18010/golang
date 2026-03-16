package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/domain/user"
)

// MockUserRepository 是用户仓储的mock实现
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

// TestNewService 测试创建服务
func TestNewService(t *testing.T) {
	repo := new(MockUserRepository)
	service := NewService(repo)

	assert.NotNil(t, service)
}

// ==================== GetUser 测试 ====================

// TestService_GetUser 测试获取用户
func TestService_GetUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 准备测试数据
	testUser := user.NewUser("test@example.com", "Test User")

	// 设置mock期望
	repo.On("FindByID", ctx, testUser.ID).Return(testUser, nil)

	// 执行测试
	result, err := service.GetUser(ctx, testUser.ID)

	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, result.ID)
	assert.Equal(t, testUser.Email, result.Email)

	// 验证mock被调用
	repo.AssertExpectations(t)
}

// TestService_GetUser_NotFound 测试获取不存在的用户
func TestService_GetUser_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 设置mock期望 - 返回错误
	repo.On("FindByID", ctx, "non-existent").Return(nil, errors.New("user not found"))

	// 执行测试
	result, err := service.GetUser(ctx, "non-existent")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

// TestService_GetUser_RepositoryError 测试仓储错误
func TestService_GetUser_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 设置mock期望 - 数据库错误
	repo.On("FindByID", ctx, "user-123").Return(nil, errors.New("database connection failed"))

	// 执行测试
	result, err := service.GetUser(ctx, "user-123")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection failed")

	repo.AssertExpectations(t)
}

// ==================== CreateUser 测试 ====================

// TestService_CreateUser 测试创建用户
func TestService_CreateUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "newuser@example.com"
	name := "New User"

	// 设置mock期望 - 邮箱不存在
	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))
	// 设置mock期望 - 保存成功
	repo.On("Save", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	// 执行测试
	result, err := service.CreateUser(ctx, email, name)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, name, result.Name)

	repo.AssertExpectations(t)
}

// TestService_CreateUser_DuplicateEmail 测试创建重复邮箱用户
func TestService_CreateUser_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "existing@example.com"
	existingUser := user.NewUser(email, "Existing User")

	// 设置mock期望 - 邮箱已存在
	repo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	// 执行测试
	result, err := service.CreateUser(ctx, email, "New User")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")

	repo.AssertExpectations(t)
}

// TestService_CreateUser_FindByEmailError 测试查找邮箱时出错
func TestService_CreateUser_FindByEmailError(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "test@example.com"

	// 设置mock期望 - 查询邮箱时出错（非"not found"错误）
	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("database error"))

	// 执行测试 - 会继续创建用户（因为不是"已存在"的情况）
	repo.On("Save", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	result, err := service.CreateUser(ctx, email, "Test User")

	// 验证结果 - 应该成功创建（因为错误不是"用户已存在"）
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, email, result.Email)

	repo.AssertExpectations(t)
}

// TestService_CreateUser_SaveError 测试保存用户失败
func TestService_CreateUser_SaveError(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "newuser@example.com"
	name := "New User"

	// 设置mock期望 - 邮箱不存在
	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))
	// 设置mock期望 - 保存失败
	repo.On("Save", ctx, mock.AnythingOfType("*user.User")).Return(errors.New("save failed"))

	// 执行测试
	result, err := service.CreateUser(ctx, email, name)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to save user")

	repo.AssertExpectations(t)
}

// TestService_CreateUser_InvalidEmail 测试无效的邮箱格式
func TestService_CreateUser_InvalidEmail(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 无效邮箱格式（缺少 @ 或 .）
	email := "invalid-email"
	name := "Test User"

	// 设置mock期望 - 邮箱不存在
	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))

	// 执行测试
	result, err := service.CreateUser(ctx, email, name)

	// 验证结果 - 邮箱格式验证应该失败
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user")

	repo.AssertExpectations(t)
}

// TestService_CreateUser_EmptyName 测试空名称
func TestService_CreateUser_EmptyName(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "test@example.com"
	name := "" // 空名称

	// 设置mock期望 - 邮箱不存在
	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))

	// 执行测试
	result, err := service.CreateUser(ctx, email, name)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user")

	repo.AssertExpectations(t)
}

// TestService_CreateUser_NameTooShort 测试名称太短
func TestService_CreateUser_NameTooShort(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "test@example.com"
	name := "A" // 名称长度小于2

	// 设置mock期望 - 邮箱不存在
	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))

	// 执行测试
	result, err := service.CreateUser(ctx, email, name)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user")

	repo.AssertExpectations(t)
}

// ==================== UpdateUserName 测试 ====================

// TestService_UpdateUser 测试更新用户
func TestService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	testUser := user.NewUser("test@example.com", "Old Name")

	// 设置mock期望
	repo.On("FindByID", ctx, testUser.ID).Return(testUser, nil)
	repo.On("Update", ctx, testUser).Return(nil)

	// 执行测试
	err := service.UpdateUserName(ctx, testUser.ID, "New Name")

	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, "New Name", testUser.Name)

	repo.AssertExpectations(t)
}

// TestService_UpdateUser_NotFound 测试更新不存在的用户
func TestService_UpdateUser_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 设置mock期望 - 用户不存在
	repo.On("FindByID", ctx, "non-existent").Return(nil, errors.New("user not found"))

	// 执行测试
	err := service.UpdateUserName(ctx, "non-existent", "New Name")

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	repo.AssertExpectations(t)
}

// TestService_UpdateUser_UpdateError 测试更新失败
func TestService_UpdateUser_UpdateError(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	testUser := user.NewUser("test@example.com", "Old Name")

	// 设置mock期望
	repo.On("FindByID", ctx, testUser.ID).Return(testUser, nil)
	repo.On("Update", ctx, testUser).Return(errors.New("update failed"))

	// 执行测试
	err := service.UpdateUserName(ctx, testUser.ID, "New Name")

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user")

	repo.AssertExpectations(t)
}

// ==================== DeleteUser 测试 ====================

// TestService_DeleteUser 测试删除用户
func TestService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	userID := "user-to-delete"

	// 设置mock期望
	repo.On("Delete", ctx, userID).Return(nil)

	// 执行测试
	err := service.DeleteUser(ctx, userID)

	// 验证结果
	require.NoError(t, err)

	repo.AssertExpectations(t)
}

// TestService_DeleteUser_Error 测试删除失败
func TestService_DeleteUser_Error(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	userID := "user-to-delete"

	// 设置mock期望 - 删除失败
	repo.On("Delete", ctx, userID).Return(errors.New("delete failed"))

	// 执行测试
	err := service.DeleteUser(ctx, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete failed")

	repo.AssertExpectations(t)
}

// TestService_DeleteUser_NotFound 测试删除不存在的用户
func TestService_DeleteUser_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	userID := "non-existent-user"

	// 设置mock期望
	repo.On("Delete", ctx, userID).Return(errors.New("user not found"))

	// 执行测试
	err := service.DeleteUser(ctx, userID)

	// 验证结果
	assert.Error(t, err)

	repo.AssertExpectations(t)
}

// ==================== ListUsers 测试 ====================

// TestService_ListUsers 测试列出用户
func TestService_ListUsers(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	users := []*user.User{
		user.NewUser("user1@example.com", "User 1"),
		user.NewUser("user2@example.com", "User 2"),
		user.NewUser("user3@example.com", "User 3"),
	}

	// 设置mock期望
	repo.On("List", ctx, 10, 0).Return(users, nil)

	// 执行测试
	result, err := service.ListUsers(ctx, 10, 0)

	// 验证结果
	require.NoError(t, err)
	assert.Len(t, result, 3)

	repo.AssertExpectations(t)
}

// TestService_ListUsers_InvalidLimit 测试无效的 limit
func TestService_ListUsers_InvalidLimit(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 测试 limit <= 0
	result, err := service.ListUsers(ctx, 0, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")
	assert.Nil(t, result)

	result, err = service.ListUsers(ctx, -1, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")
	assert.Nil(t, result)
}

// TestService_ListUsers_InvalidOffset 测试无效的 offset
func TestService_ListUsers_InvalidOffset(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 测试 offset < 0
	result, err := service.ListUsers(ctx, 10, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offset cannot be negative")
	assert.Nil(t, result)
}

// TestService_ListUsers_RepositoryError 测试仓储错误
func TestService_ListUsers_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 设置mock期望 - 查询失败
	repo.On("List", ctx, 10, 0).Return(nil, errors.New("database error"))

	// 执行测试
	result, err := service.ListUsers(ctx, 10, 0)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	repo.AssertExpectations(t)
}

// TestService_ListUsers_EmptyResult 测试空结果
func TestService_ListUsers_EmptyResult(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 设置mock期望 - 返回空列表
	repo.On("List", ctx, 10, 0).Return([]*user.User{}, nil)

	// 执行测试
	result, err := service.ListUsers(ctx, 10, 0)

	// 验证结果
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)

	repo.AssertExpectations(t)
}

// TestService_ListUsers_Pagination 测试分页
func TestService_ListUsers_Pagination(t *testing.T) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	// 第一页
	usersPage1 := []*user.User{
		user.NewUser("user1@example.com", "User 1"),
		user.NewUser("user2@example.com", "User 2"),
	}
	repo.On("List", ctx, 2, 0).Return(usersPage1, nil)

	result, err := service.ListUsers(ctx, 2, 0)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// 第二页
	usersPage2 := []*user.User{
		user.NewUser("user3@example.com", "User 3"),
		user.NewUser("user4@example.com", "User 4"),
	}
	repo.On("List", ctx, 2, 2).Return(usersPage2, nil)

	result, err = service.ListUsers(ctx, 2, 2)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	repo.AssertExpectations(t)
}

// ==================== 性能测试 ====================

// BenchmarkService_GetUser 性能测试 - 获取用户
func BenchmarkService_GetUser(b *testing.B) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	testUser := user.NewUser("test@example.com", "Test User")
	repo.On("FindByID", ctx, testUser.ID).Return(testUser, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetUser(ctx, testUser.ID)
	}
}

// BenchmarkService_CreateUser 性能测试 - 创建用户
func BenchmarkService_CreateUser(b *testing.B) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	email := "newuser@example.com"
	name := "New User"

	repo.On("FindByEmail", ctx, email).Return(nil, errors.New("not found"))
	repo.On("Save", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateUser(ctx, email, name)
	}
}

// BenchmarkService_ListUsers 性能测试 - 列出用户
func BenchmarkService_ListUsers(b *testing.B) {
	ctx := context.Background()
	repo := new(MockUserRepository)
	service := NewService(repo)

	users := []*user.User{
		user.NewUser("user1@example.com", "User 1"),
		user.NewUser("user2@example.com", "User 2"),
		user.NewUser("user3@example.com", "User 3"),
	}

	repo.On("List", ctx, 10, 0).Return(users, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ListUsers(ctx, 10, 0)
	}
}
